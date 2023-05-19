package decode

import (
    "os"
    "bytes"
    "encoding/binary"
    "time"
)

type comment_base struct {
    Seconds int32
    Nano_seconds int32
    Comment_length int32
}

// Comment contains the item of interest and the timestamp the user created it.
type Comment struct {
    Timestamp time.Time
    Value string
}

// CommentRec decodes the comment record which is for capturing anything of interest, events etc.
func CommentRec(buffer []byte, rec Record) Comment {
    // buffer1 := make([]byte, rec.Datasize)
    var buffer2 comment_base

    // _ , _ = stream.Read(buffer1)
    reader := bytes.NewReader(buffer)
    _ = binary.Read(reader, binary.BigEndian, &buffer2)

    data := Comment{
        Timestamp: time.Unix(int64(buffer2.Seconds), int64(buffer2.Nano_seconds)).UTC(),
        Value: string(buffer[12:]),
    }

    return data
}
