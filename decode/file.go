package decode

import (
    // "os"
    "encoding/binary"
    "bytes"
    // "time"
    "fmt"

    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// Tell is a small helper fucntion for telling the current position within a
// binary file opened for reading.
func Tell(stream Stream) (int64, error) {
    pos, err := stream.Seek(0, 1)

    return pos, err
}

// Padding is a small helper function for padding a GSF record.
// The GSF specification mentions that a records complete length has to be
// a multiple of 4.
// Most likely not needed for reading. Padding should be applied when writing a record.
func Padding(stream Stream) {
    pos, _ := Tell(stream)
    pad := pos % 4
    pad, _ = stream.Seek(pad, 1)

    return 
}

// Should contain the whole CRS not just horizontal and vertical datums
type Crs struct {
    Horizontal_Datum string
    Vertical_Datum string
}

type PingGroup struct {
    Start uint64
    Stop uint64
}

// FileInfo is the overarching structure containing basic info about the GSF file.
// Items include file location, file size, counts of each record (main and subrecords),
// as well as basic info about the pings such as number of beams and schema for each
// ping contained within the file.
type FileInfo struct {
    GSF_URI string
    Size uint64
    Sensor_ID int32
    Sensor_Name string
    CRS Crs
    SubRecord_Schema []string
    Quality_Info QualityInfo
    Record_Counts map[string]uint64
    SubRecord_Counts map[string]uint64
    Ping_Groups []PingGroup
    Record_Index map[string][]RecordHdr
    Ping_Info []PingInfo
}

// PingGroups combines pings together based on their presence or absence of
// scale factors. It is a forward linear search, and if a given ping is missing
// scale factors, then it included as part of the ping group where the previous
// set of scale factors were found.
// For example; [0, 10] indicates that the ping group contains pings 0 up to and
// including ping 9. It is a [start, stop) index based on the linear ordering
// of pings found in the GSF file.
func (fi *FileInfo) PGroups() {
    var (
        start int
        ping_group PingGroup
        groups []PingGroup
    )

    groups = make([]PingGroup, 0)

    for i, ping := range(fi.Ping_Info) {
        if ping.Scale_Factors {
            if i > 0 {
                // new group
                ping_group = PingGroup{uint64(start), uint64(i)}
                groups = append(groups, ping_group)
            }
            start = i
        }
    }

    fi.Ping_Groups = groups
}

// Index, as the name implies, builds a file index of all Record types.
// Each Record contains the record ID, record size, byte index and checksum flag.
func Index(gsf_uri string, config_uri string, in_memory bool) FileInfo {

    var (
        rec_idx map[string][]RecordHdr
        rec_counts map[string]uint64
        sub_rec_counts map[SubRecordID]uint64
        sub_rec_counts_str map[string]uint64
        sensor_id int32 
        sensor_name string
        rec Record
        pinfo PingInfo
        pings []PingInfo
        finfo FileInfo
        config *tiledb.Config
        err error
        crs Crs
    )

    // get a generic config if no path provided
    if config_uri == "" {
        config, err = tiledb.NewConfig()
        if err != nil {
            panic(err)
        }
    } else {
        config, err = tiledb.LoadConfig(config_uri)
        if err != nil {
            panic(err)
        }
    }

    defer config.Free()

    ctx, err := tiledb.NewContext(config)
    if err != nil {
        panic(err)
    }
    defer ctx.Free()

    vfs, err := tiledb.NewVFS(ctx, config)
    if err != nil {
        panic(err)
    }
    defer vfs.Free()

    handler, err := vfs.Open(gsf_uri, tiledb.TILEDB_VFS_READ)
    if err != nil {
        panic(err)
    }
    defer handler.Close()
    // defer stream.Free()

    rec_idx = make(map[string][]Record)
    rec_counts = make(map[string]uint64)
    sub_rec_counts = make(map[SubRecordID]uint64)
    sub_rec_counts_str = make(map[string]uint64)

    one := uint64(1)

    // filesize is used as an EOF indicator when streaming the raw bytes
    filesize, _ := vfs.FileSize(gsf_uri)

    // create a generic stream
    stream, err := GenericStream(handler, filesize, in_memory)

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(stream)

    // start at front of the stream
    pos, _ := stream.Seek(0, 0)

    // reading the bytestream and build record index information
    for uint64(pos) < filesize {
        // TODO; test that pos moves after we read a header
        rec = RecordHdr(stream)

        // increment record count
        rec_counts[RecordNames[rec.Id]] += one

        rec_idx[RecordNames[rec.Id]] = append(rec_idx[RecordNames[rec.Id]], rec)

        // TODO; convert to case switch
        if rec.Id == SWATH_BATHYMETRY_PING {
            // need to do some sub record decoding
            buffer := make([]byte, rec.Datasize)
            _ = binary.Read(stream, binary.BigEndian, &buffer)
            reader := bytes.NewReader(buffer)

            pinfo = ping_info(reader, rec)
            pings = append(pings, pinfo)

            // increment sub-record count
            for _, sid := range(pinfo.Sub_Records) {
                sub_rec_counts[sid] += one
            }

            pos, _ = Tell(stream)
        } else if rec.Id == PROCESSING_PARAMETERS {
            // should only be one of these records in the gsf file
            buffer := make([]byte, rec.Datasize)
            _ = binary.Read(stream, binary.BigEndian, &buffer)
            params := ProcessingParametersRec(buffer, rec)

            // TODO; change params rec to be a defined struct to avoid this type assertion
            hd, ok := params["geoid"].(string)
            if ok {
                crs.Horizontal_Datum = fmt.Sprint(hd)
            }
            vd, ok := params["tidal_datum"].(string)
            if ok {
                crs.Vertical_Datum = fmt.Sprint(vd)
            }
        } else {
            // seek over the record and loop to the next
            pos, _ = stream.Seek(int64(rec.Datasize), 1)
        }

    }

    // reset file posistion
    _, _ = stream.Seek(original_pos, 0)

    // consistent schema; we've had cases where the schema is inconsistent between pings
    sr_schema := make([]string, 0)
    for key, val := range sub_rec_counts {
        sub_rec_counts_str[SubRecordNames[key]] = val
        if key > 100 {
            sensor_id = int32(key)
            sensor_name = SubRecordNames[key]
        } else if key < 100 {
            sr_schema = append(sr_schema, SubRecordNames[key])
        }
    }

    finfo.GSF_URI = gsf_uri
    finfo.Size = filesize
    finfo.Sensor_ID = sensor_id
    finfo.Sensor_Name = sensor_name
    finfo.CRS = crs
    finfo.SubRecord_Schema = sr_schema
    finfo.Record_Counts = rec_counts
    finfo.SubRecord_Counts = sub_rec_counts_str
    finfo.Record_Index = rec_idx
    finfo.Ping_Info = pings

    finfo.PGroups()
    finfo.QInfo()

    return finfo
}
