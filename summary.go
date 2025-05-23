package gsf

import (
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
	End_datetime   time.Time
	Min_longitude  float64
	Max_longitude  float64
	Min_latitude   float64
	Max_latitude   float64
	Min_depth      float64
	Max_depth      float64
}

// DecodeSwathBathySummary acts as the constructor for SwathBathySummary by decoding
// the SWATH_BATHY_SUMMARY record.
// It contains the geometrical and temporal extent of the swath data contained
// within the GSF file.
func DecodeSwathBathySummary(reader *bytes.Reader) SwathBathySummary {
	var buffer struct {
		First_ping_sec      uint32
		First_ping_nano_sec uint32
		Last_ping_sec       uint32
		Last_ping_nano_sec  uint32
		Min_lat             uint32
		Min_lon             uint32
		Max_lat             uint32
		Max_lon             uint32
		Min_depth           uint32
		Max_depth           int32
	}

	_ = binary.Read(reader, binary.BigEndian, &buffer)

	// should look at storing the scale factors as const's or a struct
	summary := SwathBathySummary{
		Start_datetime: time.Unix(int64(buffer.First_ping_sec), int64(buffer.First_ping_nano_sec)).UTC(),
		End_datetime:   time.Unix(int64(buffer.Last_ping_sec), int64(buffer.Last_ping_nano_sec)).UTC(),
		Min_longitude:  float64(int32(buffer.Min_lon)) / SCALE_7_F64,
		Max_longitude:  float64(int32(buffer.Max_lon)) / SCALE_7_F64,
		Min_latitude:   float64(int32(buffer.Min_lat)) / SCALE_7_F64,
		Max_latitude:   float64(int32(buffer.Max_lat)) / SCALE_7_F64,
		Min_depth:      float64(buffer.Min_depth) / SCALE_2_F64,
		Max_depth:      float64(buffer.Max_depth) / SCALE_2_F64,
	}

	return summary
}
