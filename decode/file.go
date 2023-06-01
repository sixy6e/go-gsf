package decode

import (
    // "os"
    "encoding/binary"
    "bytes"
    "time"
    // "fmt"

    "github.com/samber/lo"
    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// Tell is a small helper fucntion for telling the current position within a
// binary file opened for reading.
func Tell(stream *tiledb.VFSfh) (int64, error) {
    pos, err := stream.Seek(0, 1)

    return pos, err
}

// Padding is a small helper function for padding a GSF record.
// The GSF specification mentions that a records complete length has to be
// a multiple of 4.
// Most likely not needed for reading. Padding should be applied when writing a record.
func Padding(stream *tiledb.VFSfh) {
    pos, _ := Tell(stream)
    pad := pos % 4
    pad, _ = stream.Seek(pad, 1)

    return 
}

// GSFv contains the version of GSF used to construct the GSF file.
type GSFv struct {
    Version string
}

// FileRec decodes the HEADER record from a GSF file.
// It contains the version of GSF used to create the file.
func FileRec(buffer []byte, rec Record) GSFv {
    // buffer := make([]byte, rec.Datasize)

    // _ = binary.Read(stream, binary.BigEndian, &buffer)

    file_hdr := GSFv{Version: string(buffer)}

    return file_hdr
}

var rec_arr = [12]RecordID{
    HEADER,
    SWATH_BATHYMETRY_PING,
    SOUND_VELOCITY_PROFILE,
    PROCESSING_PARAMETERS,
    SENSOR_PARAMETERS,
    COMMENT,
    HISTORY,
    NAVIGATION_ERROR,
    SWATH_BATHY_SUMMARY,
    SINGLE_BEAM_PING,
    HV_NAVIGATION_ERROR,
    ATTITUDE,
}

var subrec_arr = [32]SubRecordID{
    DEPTH,
    ACROSS_TRACK,
    ALONG_TRACK,
    TRAVEL_TIME,
    BEAM_ANGLE,
    MEAN_CAL_AMPLITUDE,
    MEAN_REL_AMPLITUDE,
    ECHO_WIDTH,
    QUALITY_FACTOR,
    RECEIVE_HEAVE,
    DEPTH_ERROR,
    ACROSS_TRACK_ERROR,
    ALONG_TRACK_ERROR,
    NOMINAL_DEPTH,
    QUALITY_FLAGS,
    BEAM_FLAGS,
    SIGNAL_TO_NOISE,
    BEAM_ANGLE_FORWARD,
    VERTICAL_ERROR,
    HORIZONTAL_ERROR,
    INTENSITY_SERIES,
    SECTOR_NUMBER,
    DETECTION_INFO,
    INCIDENT_BEAM_ADJ,
    SYSTEM_CLEANING,
    DOPPLER_CORRECTION,
    SONAR_VERT_UNCERTAINTY,
    SONAR_HORZ_UNCERTAINTY,
    DETECTION_WINDOW,
    MEAN_ABS_COEF,
    UNKNOWN,
    SCALE_FACTORS,
}

// FileInfo is the overarching structure containing basic info about the GSF file.
// Items include file location, file size, counts of each record (main and subrecords),
// as well as basic info about the pings such as number of beams and schema for each
// ping contained within the file.
type FileInfo struct {
    GSF_URI string
    Size uint64
    Min_Max_Beams []uint16
    Consistent_Beams bool
    Duplicate_Pings bool
    Duplicates []time.Time
    // Record_Counts map[RecordID]uint64
    Record_Counts map[string]uint64
    // SubRecord_Counts map[SubRecordID]uint64
    SubRecord_Counts map[string]uint64
    // Record_Index map[RecordID][]Record
    Record_Index map[string][]Record
    Ping_Info []PingInfo
}

// Index, as the name implies, builds a file index of all Record types.
// Each Record contains the record ID, record size, byte index and checksum flag.
func Index(gsf_uri string, config_uri string) FileInfo {

    var (
        // rec_idx map[RecordID][]Record
        rec_idx map[string][]Record
        // rec_counts map[RecordID]uint64
        rec_counts map[string]uint64
        // sub_rec_counts map[SubRecordID]uint64
        sub_rec_counts map[string]uint64
        // val1 RecordID  // used for zeroing initial state
        // val2 SubRecordID  // used for zeroing initial state
        rec Record
        pinfo PingInfo
        pings []PingInfo
        finfo FileInfo
        config *tiledb.Config
        err error
        nbeams []uint16
        timestamps []time.Time
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

    stream, err := vfs.Open(gsf_uri, tiledb.TILEDB_VFS_READ)
    if err != nil {
        panic(err)
    }
    defer stream.Close()
    // defer stream.Free()

    // rec_idx = make(map[RecordID][]Record)
    rec_idx = make(map[string][]Record)
    // rec_counts = make(map[RecordID]uint64)
    rec_counts = make(map[string]uint64)
    // sub_rec_counts = make(map[SubRecordID]uint64)
    sub_rec_counts = make(map[string]uint64)
    // nbeams = make([]uint64, 0)  // could be faster to declare after

    one := uint64(1)
    // zero := uint64(0)  // used for zeroing initial state

    // potentially superfluous; only required if there is desire to show records that
    // don't exist within the file, i.e. records with count=0
    // for _, val1 = range rec_arr {
    //     rec_idx[val1] = nil
    //     rec_counts[val1] = zero
    // }

    // potentially superfluous; only required if there is desire to show records that
    // don't exist within the file, i.e. records with count=0
    // for _, val2 = range subrec_arr {
    //     sub_rec_counts[val2] = zero
    // }

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(stream)

    // filesize is used as an EOF indicator when streaming the raw bytes
    // filestat, _ := stream.Stat()
    filesize, _ := vfs.FileSize(gsf_uri)
    // filename := filestat.Name()

    // start at front of the stream
    pos, _ := stream.Seek(0, 0)

    // reading the bytestream and build record index information
    for uint64(pos) < filesize {
        // TODO; test that pos moves after we read a header
        rec = RecordHdr(stream)

        // increment record count
        // rec_counts[rec.Id] += one
        rec_counts[RecordNames[rec.Id]] += one

        rec_idx[RecordNames[rec.Id]] = append(rec_idx[RecordNames[rec.Id]], rec)

        if rec.Id == SWATH_BATHYMETRY_PING {
            // need to do some sub record decoding
            buffer := make([]byte, rec.Datasize)
            _ = binary.Read(stream, binary.BigEndian, &buffer)
            reader := bytes.NewReader(buffer)

            pinfo = ping_info(reader, rec)
            pings = append(pings, pinfo)
            // nbeams = append(nbeams, ping.Number_Beams)

            // increment sub-record count
            for _, sid := range(pinfo.Sub_Records) {
                // sub_rec_counts[sid] += one
                sub_rec_counts[SubRecordNames[sid]] += one
            }

            pos, _ = Tell(stream)
        } else {
            // seek over the record and loop to the next
            pos, _ = stream.Seek(int64(rec.Datasize), 1)
        }

    }

    // reset file posistion
    _, _ = stream.Seek(original_pos, 0)

    // there have been instances where the number of beams was inconsistent between pings
    // the general idea is to know whether we're dealing with a consistent number of beams
    nbeams = make([]uint16, len(pings))
    for i, ping := range(pings) {
        nbeams[i] = ping.Number_Beams
    }

    // domain for number of beams
    max := lo.Max(nbeams)
    min := lo.Min(nbeams)
    min_max_beams := []uint16{min, max}
    consistent_beams := min == max

    // duplicate pings. one of the samples had duplicate timestamps
    timestamps = make([]time.Time, len(pings))
    for i, ping := range(pings) {
        timestamps[i] = ping.Timestamp
    }

    duplicates := lo.FindDuplicates(timestamps)

    finfo.GSF_URI = gsf_uri
    finfo.Size = filesize
    finfo.Min_Max_Beams = min_max_beams
    finfo.Consistent_Beams = consistent_beams
    finfo.Duplicate_Pings = len(duplicates) > 0
    finfo.Duplicates = duplicates
    finfo.Record_Counts = rec_counts
    finfo.SubRecord_Counts = sub_rec_counts
    finfo.Record_Index = rec_idx
    finfo.Ping_Info = pings

    return finfo
}
