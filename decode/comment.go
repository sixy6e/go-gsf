package decode

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
)

// Comment contains the item of interest and the timestamp the user created it.
type Comment struct {
    Timestamp time.Time
    Value string
}

// DecodeComment is a constructor for Comment by decoding the COMMENT record which
// is for capturing anything of interest, events etc.
func DecodeComment(buffer []byte) Comment {
    var buffer2 struct {
        Seconds int32
        Nano_seconds int32
        Comment_length int32
    }

    reader := bytes.NewReader(buffer)
    _ = binary.Read(reader, binary.BigEndian, &buffer2)

    data := Comment{
        Timestamp: time.Unix(int64(buffer2.Seconds), int64(buffer2.Nano_seconds)).UTC(),
        Value: string(buffer[12:]),
    }

    return data
}
