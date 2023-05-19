package decode

import (
    // "os"
    // "bytes"
    "encoding/binary"
    "time"
)

type summary_base struct {
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

// SwathBathySummaryRec decodes the SWATH_BATHY_SUMMARY record.
// It contains the geometrical and temporal extent of the swath data contained
// within the GSF file.
func SwathBathySummaryRec(buffer []byte, rec Record) SwathBathySummary {
    var buffer2 summary_base

    // for s3 reads we might need to look at reading the entire record into a buffer
    // then convert to a reader to pass through the binary.Read
    // buffer = make([]byte, rec.Datasize) // if wanting generality;
    // _ , _ = stream.Read(buffer)
    // reader := bytes.NewReader(buffer)
    // buffer2 := struct{...}
    // _ = binary.Read(reader, binary.BigEndian, &buffer2)

    // TODO; look for alternate way of processing; potentially json.unmarshall???
    reader := bytes.NewReader(buffer)
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
