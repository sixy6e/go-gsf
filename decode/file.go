package decode

import (
    "os"
    "encoding/binary"
)

// Tell is a small helper fucntion for telling the current position within a
// binary file opened for reading.
func Tell(stream *os.File) int64 {
    pos, err := stream.Seek(0, 1)

    return pos, err
}

// Padding is a small helper function for padding a GSF record.
// The GSF specification mentions that a records complete length has to be
// a multiple of 4.
// Most likely not needed for reading. Padding should be applied when writing a record.
func Padding(stream *os.File) {
    pos := Tell(stream)
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
func FileRec(stream *os.File, rec Record) GSFv {
    buffer := make([]byte, rec.Datasize)

    _ = binary.Read(stream, binary.BigEndian, &buffer)

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
    Record_Counts map[RecordID]uint64
    SubRecord_Counts map[SubRecordID]uint64
    Pings map[RecordID][]PingInfo
}

// Index, as the name implies, builds a file index of all Record types.
// Each Record contains the record ID, record size, byte index and checksum flag.
func Index(stream *os.File) any {  // TODO return type(s)
    // probably return map[RecordId]uint64, map[SubRecordId][uint64], []PingInfo

    var (
        rec_idx = map[RecordID][]Record
        rec_counts = map[RecordID]uint64
        sub_rec_counts = map[SubRecordID]uint64
        val1 = RecordID
        val2 = SubRecordID
        rec = Record
        pinfo = PingInfo
        pings = []PingInfo
    )

    one := uint64(1)
    zero := uint64(0)

    for _, val1 := range rec_arr {
        rec_idx[val1] = nil
        rec_counts[val1] = zero
    }

    for _, val2 = range subrec_arr {
        sub_rec_counts[val2] = zero
    }

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(stream)

    // filesize is used as an EOF indicator when streaming the raw bytes
    filestat, _ := stream.Stat()
    filesize := filestat.Size()
    filename := filestat.Name()

    // start at front of the stream
    pos, _ := stream.Seek(0, 0)

    // reading the bytestream and build record index information
    for pos < filesize {
        // TODO; test that pos moves after we read a header
        rec = RecordHdr(stream)

        // increment record count
        rec_counts[rec.Id] += one

        rec_idx[rec.Id] = append(rec_idx[rec.Id], rec)

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
        } else {
            // seek over the record and loop to the next
            pos, _ = stream.Seek(int64(rec.Datasize), 1)
        }

        // read the record and loop to the next
        // pos, _ = stream.Seek(int64(rec.Datasize), 1)  // careful we don't do double reads (pinfo does seeking too)
        // _ = Padding(stream)

    }
}
