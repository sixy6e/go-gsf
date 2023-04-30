package decode

import (
    "os"
    // "bytes"
    "encoding/binary"
    "time"
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

// HistoryRec decodes the HISTORY record which contains any documentation or processing
// which has been applied to the data.
// Captured information; time the step occurred, operator name, computer name
// program being used and any command line args or relevant parameters, as well as any
// comments to summarise the processing that occurred.
func HistoryRec(stream *os.File, rec Record) History {
    buffer := make([]byte, rec.Datasize)

    _ , _ = stream.Read(buffer)
    // reader := bytes.NewReader(buffer)
    // _ = binary.Read(reader, binary.BigEndian, &buffer2)

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
    size := int16(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))  // TODO; test int64
    size64 := int64(size)
    start_idx += 2
    end_idx += size64
    machine_name := string(buffer[start_idx:end_idx])
    start_idx += size64
    end_idx += 2

    // operator name
    size = int16(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    size64 = int64(size)
    start_idx += 2
    end_idx += size64
    operator_name := string(buffer[start_idx:end_idx])
    start_idx += size64
    end_idx += 2

    // command
    size = int16(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    size64 = int64(size)
    start_idx += 2
    end_idx += size64
    command := string(buffer[start_idx:end_idx])
    start_idx += size64
    end_idx += 2

    // comment
    size = int16(binary.BigEndian.Uint16(buffer[start_idx:end_idx]))
    size64 = int64(size)
    start_idx += 2
    end_idx += size64
    comment := string(buffer[start_idx:end_idx])

    history := History{
        Processing_timestamp: time.Unix(seconds, nano_seconds).UTC(),
        Machine_name: machine_name,
        Operator_name: operator_name,
        Command: command,
        Value: comment,
    }

    return history
}
