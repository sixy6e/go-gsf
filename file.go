package gsf

import (
    "encoding/binary"
    "bytes"
    "fmt"

    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// Tell is a small helper function for telling the current position within a
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

type GsfFile struct {
    Uri string
    filesize uint64
    config *tiledb.Config
    ctx *tiledb.Context
    vfs *tiledb.VFS
    handler *tiledb.VFSfh
    Stream
}

func OpenGSF(gsf_uri string, config_uri string, in_memory bool) GsfFile {
    var (
        gsf GsfFile
        config *tiledb.Config
        err error
    )

    gsf.Uri = gsf_uri
    
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

    // defer config.Free()
    gsf.config = config

    ctx, err := tiledb.NewContext(config)
    if err != nil {
        panic(err)
    }
    // defer ctx.Free()
    gsf.ctx = ctx

    vfs, err := tiledb.NewVFS(ctx, config)
    if err != nil {
        panic(err)
    }
    // defer vfs.Free()
    gsf.vfs = vfs

    handler, err := vfs.Open(gsf_uri, tiledb.TILEDB_VFS_READ)
    if err != nil {
        panic(err)
    }
    // defer handler.Close()
    gsf.handler = handler

    filesize, _ := vfs.FileSize(gsf_uri)
    gsf.filesize = filesize

    // generic stream
    stream, err := GenericStream(handler, filesize, in_memory)

    gsf.Stream = stream

    return gsf
}

// Releases the open tiledb file handler connections.
func (g *GsfFile) Close() {
    g.handler.Close()
    g.vfs.Free()
    g.ctx.Free()
    g.config.Free()
}

// RecBuf reads the bytes from an opened GsfFile specified by the RecordHdr.
func (g *GsfFile) RecBuf(r RecordHdr) (buffer []byte) {
    var err error

    buffer = make([]byte, r.Datasize)
    _, err = g.Stream.Seek(r.Byte_index, 0)

    if err != nil {
        panic(err)
    }

    _ = binary.Read(g.Stream, binary.BigEndian, &buffer)

    return buffer
}

func (g *GsfFile) ProcInfo(fi *FileInfo) (proc_info ProcessingInfo) {
    proc_info.Histories = g.HistoryRecords(fi)
    proc_info.Comments = g.CommentRecords(fi)

    buffer := g.RecBuf(fi.Index.Record_Index["PROCESSING_PARAMETERS"][0])
    proc_info.Processing_Parameters = DecodeProcessingParameters(buffer)

    return proc_info
}

// func (g *GsfFile) AttInfo(fi *FileInfo) (Attitude, Attitude) {
//     buffer := g.RecBuf(fi.Index.Record_Index["ATTITUDE"][0])
//     at1 := DecodeAttitude(buffer)
//     last := fi.Metadata.Record_Counts["ATTITUDE"] - 1
//     buffer = g.RecBuf(fi.Index.Record_Index["ATTITUDE"][last])
//     at2 := DecodeAttitude(buffer)
// 
//     return at1, at2
// }

type GsfDetails struct {
    GSF_URI string
    GSF_Version string
    Size uint64
}

type SensorInfo struct {
    Sensor_ID int32
    Sensor_Name string
}

type Metadata struct {
    GSF_Details GsfDetails
    Sensor_Info SensorInfo
    CRS Crs
    SubRecord_Schema []string
    Quality_Info QualityInfo
    Record_Counts map[string]uint64
    SubRecord_Counts map[string]uint64
    Measurement_Counts map[string]uint64
    Swath_Summary SwathBathySummary
}

type Index struct {
    Ping_Groups []PingGroup
    Record_Index map[string][]RecordHdr
}

type ProcessingInfo struct {
    Histories []History
    Comments []Comment
    Processing_Parameters map[string]interface{}
}
    
// FileInfo is the overarching structure containing basic info about the GSF file.
// Items include file location, file size, counts of each record (main and subrecords),
// as well as basic info about the pings such as number of beams and schema for each
// ping contained within the file.
type FileInfo struct {
    Metadata
    Index
    // Processing_Parameters map[string]interface{}
    Ping_Info []PingInfo
}

// Info builds a file index of all Record types as well generic information
// and metadata such as CRS, sensor, schema, record counts, and basic QA.
func (g *GsfFile) Info() FileInfo {
    var (
        rec_idx map[string][]RecordHdr
        rec_counts map[string]uint64
        sub_rec_counts map[SubRecordID]uint64
        sub_rec_counts_str map[string]uint64
        sensor_id int32 
        sensor_name string
        rec RecordHdr
        pinfo PingInfo
        pings []PingInfo
        finfo FileInfo
        // err error
        crs Crs
        buffer []byte
        reader *bytes.Reader
        version Header
        swath_sum SwathBathySummary
        params map[string]interface{}
        meas_counts map[string]uint64
        // att_measurements uint64
        // att_time []time.Time
    )

    rec_idx = make(map[string][]RecordHdr)
    rec_counts = make(map[string]uint64)
    sub_rec_counts = make(map[SubRecordID]uint64)
    sub_rec_counts_str = make(map[string]uint64)
    meas_counts = make(map[string]uint64)

    one := uint64(1)

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    // start at front of the stream
    pos, _ := g.Stream.Seek(0, 0)

    // reading the byte stream and build record index information
    for uint64(pos) < g.filesize {
        // TODO; test that pos moves after we read a header
        rec = DecodeRecordHdr(g.Stream)

        // increment record count
        rec_counts[RecordNames[rec.Id]] += one

        rec_idx[RecordNames[rec.Id]] = append(rec_idx[RecordNames[rec.Id]], rec)

        switch rec.Id {
            case SWATH_BATHYMETRY_PING:
                // need to do some sub record decoding
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)
                reader = bytes.NewReader(buffer)

                pinfo = ping_info(reader, rec)
                pings = append(pings, pinfo)

                // increment sub-record count
                for _, sid := range(pinfo.Sub_Records) {
                    sub_rec_counts[sid] += one
                }

                // increment total point (measurement/observation) count
                meas_counts[RecordNames[rec.Id]] += uint64(pinfo.Number_Beams)

                pos, _ = Tell(g.Stream)
            case PROCESSING_PARAMETERS:
                // should only be one of these records in the gsf file
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)

                params = DecodeProcessingParameters(buffer)

                // TODO; change params rec to be a defined struct to avoid this type assertion
                hd, ok := params["geoid"].(string)
                if ok {
                    crs.Horizontal_Datum = fmt.Sprint(hd)
                }
                vd, ok := params["tidal_datum"].(string)
                if ok {
                    crs.Vertical_Datum = fmt.Sprint(vd)
                }
            case SWATH_BATHY_SUMMARY:
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)
                reader = bytes.NewReader(buffer)

                swath_sum = DecodeSwathBathySummary(reader)
            case HEADER:
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)

                version = DecodeHeader(buffer)
            case ATTITUDE:
                // at this stage, only interested in the total observation count
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)
                reader = bytes.NewReader(buffer)
                att_hdr := attitude_header(reader)
                meas_counts[RecordNames[rec.Id]] += att_hdr.Measurements
            case SOUND_VELOCITY_PROFILE:
                // at this stage, only interested in the total observation count
                buffer = make([]byte, rec.Datasize)
                _ = binary.Read(g.Stream, binary.BigEndian, &buffer)
                reader = bytes.NewReader(buffer)
                s_hdr := svp_header(reader)
                meas_counts[RecordNames[rec.Id]] += s_hdr.N_points
            default:
                // seek over the record and loop to the next
                pos, _ = g.Stream.Seek(int64(rec.Datasize), 1)
        }

    }

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

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

    finfo.Metadata.GSF_Details = GsfDetails{GSF_URI: g.Uri, GSF_Version: version.Version, Size: g.filesize}
    finfo.Metadata.Sensor_Info = SensorInfo{Sensor_ID: sensor_id, Sensor_Name: sensor_name}
    finfo.Metadata.CRS = crs
    finfo.Metadata.SubRecord_Schema = sr_schema
    finfo.Metadata.Record_Counts = rec_counts
    finfo.Metadata.SubRecord_Counts = sub_rec_counts_str
    finfo.Metadata.Measurement_Counts = meas_counts
    finfo.Metadata.Swath_Summary = swath_sum

    finfo.Index.Record_Index = rec_idx

    finfo.Ping_Info = pings
    // finfo.Processing_Parameters = params

    finfo.PGroups()
    finfo.QInfo()

    return finfo
}
