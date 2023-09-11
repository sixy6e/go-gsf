package decode

import (
    // "os"
    // "bytes"
    "encoding/binary"
    "time"
    "strings"
)

// History type contains the documentation or processing that has been applied to the data.
// It contains the timestamp the step occurred, operator name, machine name,
// any command line arguments or relevant parameters, as well any comments to summarise
// the processing that occurred.
type History struct {
    Processing_timestamp time.Time
    Machine_name string
    Operator_name string
    Command string
    Value string
}

// DecodeHistory is a constructor for History by decoding the HISTORY record which
// contains any documentation or processing
// which has been applied to the data.
// Captured information; time the step occurred, operator name, computer name
// program being used and any command line args or relevant parameters, as well as any
// comments to summarise the processing that occurred.
func DecodeHistory(buffer []byte) History {
    // timestamp
    seconds := int64(binary.BigEndian.Uint32(buffer[0:4]))
    nano_seconds := int64(binary.BigEndian.Uint32(buffer[4:8]))

    // start stop index markers (first 8 bytes was for the timestamp)
    start_idx := int64(8)
    end_idx := int64(10)

    // essentially each item comprises of two pieces.
    // a value of type int16 which indicates the length of the following string
    // TODO find a better method of deconstructing this blob

    // machine name
    size := int64(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    start_idx += 2
    end_idx += size
    machine_name := strings.Trim(string(buffer[start_idx:end_idx]), "\x00")
    start_idx += size
    end_idx += 2

    // operator name
    size = int64(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    start_idx += 2
    end_idx += size
    operator_name := strings.Trim(string(buffer[start_idx:end_idx]), "\x00")
    start_idx += size
    end_idx += 2

    // command
    size = int64(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    start_idx += 2
    end_idx += size
    command := strings.Trim(string(buffer[start_idx:end_idx]), "\x00")
    start_idx += size
    end_idx += 2

    // comment
    size = int64(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    start_idx += 2
    end_idx += size
    comment := strings.Trim(string(buffer[start_idx:end_idx]), "\x00")

    history := History{
        Processing_timestamp: time.Unix(seconds, nano_seconds).UTC(),
        Machine_name: machine_name,
        Operator_name: operator_name,
        Command: command,
        Value: comment,
    }

    return history
}

// HistoryRecords decodes all HISTORY records.
func (g *GsfFile) HistoryRecords(fi *FileInfo) (history []History) {
    var (
        buffer []byte
    )
    history = make([]History, fi.Record_Counts["HISTORY"])

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    for _, rec := range(fi.Record_Index["HISTORY"]) {
        buffer = g.RecBuf(rec)
        hist := DecodeHistory(buffer)
        history = append(history, hist)
    }

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

    return history
}
