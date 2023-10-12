package gsf

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
    "strings"
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
        Value: strings.Trim(string(buffer[12:]), "\x00"),
    }

    return data
}

// CommentRecords decodes all COMMENT records.
func (g *GsfFile) CommentRecords(fi *FileInfo) (comments []Comment) {
    var (
        buffer []byte
    )
    comments = make([]Comment, fi.Record_Counts["COMMENT"])

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    for _, rec := range(fi.Record_Index["COMMENT"]) {
        buffer = g.RecBuf(rec)
        comment := DecodeComment(buffer)
        comments = append(comments, comment)
    }

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

    return comments
}
