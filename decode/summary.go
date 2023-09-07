package decode

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
)

// SwathBathySummary type contains the summary information over the entire swath data
// contained within the GSF file.
// Fields include start and end datetime, min/max of longitude, latitude and depth.
// Conceptually a 4 dimensional extent description consisting of (x, y, z, t).
type SwathBathySummary struct {
    Start_datetime time.Time
    End_datetime time.Time
    Min_longitude float64
    Max_longitude float64
    Min_latitude float64
    Max_latitude float64
    Min_depth float32
    Max_depth float32
}

// DecodeSwathBathySummary acts as the constructor for SwathBathySummary by decoding
// the SWATH_BATHY_SUMMARY record.
// It contains the geometrical and temporal extent of the swath data contained
// within the GSF file.
func DecodeSwathBathySummary(buffer []byte) SwathBathySummary {
    var buffer2 struct {
        First_ping_sec int32
        First_ping_nano_sec int32
        Last_ping_sec int32
        Last_ping_nano_sec int32
        Min_lat int32
        Min_lon int32
        Max_lat int32
        Max_lon int32
        Min_depth int32
        Max_depth int32
    }

    reader := bytes.DecodeReader(buffer)
    _ = binary.Read(reader, binary.BigEndian, &buffer2)

    // should look at storing the scale factors as consts or a struct
    summary := SwathBathySummary{
        Start_datetime: time.Unix(int64(buffer2.First_ping_sec), int64(buffer2.First_ping_nano_sec)).UTC(),
        End_datetime: time.Unix(int64(buffer2.Last_ping_sec), int64(buffer2.Last_ping_nano_sec)).UTC(),
        Min_longitude: float64(float32(buffer2.Min_lon) / SCALE1),
        Max_longitude: float64(float32(buffer2.Max_lon) / SCALE1),
        Min_latitude: float64(float32(buffer2.Min_lat) / SCALE1),
        Max_latitude: float64(float32(buffer2.Max_lat) / SCALE1),
        Min_depth: float32(buffer2.Min_depth) / SCALE2,
        Max_depth: float32(buffer2.Max_depth) / SCALE2,
    }

    return summary
}
