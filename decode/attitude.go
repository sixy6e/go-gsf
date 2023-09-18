package decode

import (
    "bytes"
    "encoding/binary"
    "time"
)

// The start and end datetimes, might not reflect the true start and end datetimes.
// The records use and offset for each measurement, so the end time will surely
// occur after the base time for the last attitude record.
type AttitudeSummary struct {
    Start_datetime time.Time
    End_datetime time.Time
    Measurement_count uint64
}

// Attitude contains the measurements as reported by the vessel attitude sensor.
// Fields include: Timestamp, Pitch, Roll, Heave and Heading.
type Attitude struct {
    Timestamp []time.Time
    Pitch []float32
    Roll []float32
    Heave []float32
    Heading []float32
}

type attitude_hdr struct {
    Seconds int64
    Nano_seconds int64
    Timestamp time.Time
    Measurements uint64
}

func attitude_header(reader *bytes.Reader) (att_hdr attitude_hdr) {
    var (
        base struct {
            Seconds int32
            Nano_seconds int32
            Measurements int16
        }
    )
    _ = binary.Read(reader, binary.BigEndian, &base)
    acq_time := time.Unix(int64(base.Seconds), int64(base.Nano_seconds)).UTC()
    att_hdr.Seconds = int64(base.Seconds)
    att_hdr.Nano_seconds = int64(base.Nano_seconds)
    att_hdr.Timestamp = acq_time
    att_hdr.Measurements = uint64(base.Measurements)
    return att_hdr
}

// DecodeAttitude is a constructor for Attitude by decoding an ATTITUDE Record
// which contains the measurements
// as reported by the vessel attitude sensor.
// Fields include: Timestamp, Pitch, Roll, Heave and Heading.
func DecodeAttitude(buffer []byte) Attitude {
    var (
        idx int64 = 0
        base struct {
            Time_offset int16
            Pitch int16
            Roll int16
            Heave int16
            Heading int16
        }
        offset time.Duration
    )

    reader := bytes.NewReader(buffer)

    // TODO; create a small func to decode the attitude header and find the total of n-measurements
    att_hdr := attitude_header(reader)
    idx += 10  // TODO; remove, if superfluous

    attitude := Attitude{
        Timestamp: make([]time.Time, att_hdr.Measurements),
        Pitch: make([]float32, att_hdr.Measurements),
        Roll: make([]float32, att_hdr.Measurements),
        Heave: make([]float32, att_hdr.Measurements),
        Heading: make([]float32, att_hdr.Measurements),
    }

    for i:= uint64(0); i < att_hdr.Measurements; i++ {
        _ = binary.Read(reader, binary.BigEndian, &base)

        // the offset is scaled by 1000, indicating the units are now in milliseconds
        offset = time.Duration(base.Time_offset)
        attitude.Timestamp[i] = att_hdr.Timestamp.Add(time.Millisecond * offset)
        attitude.Pitch[i] = float32(base.Pitch) / SCALE2
        attitude.Roll[i] = float32(base.Roll) / SCALE2
        attitude.Heave[i] = float32(base.Heave) / SCALE2
        attitude.Heading[i] = float32(base.Heading) / SCALE2
    }

    return attitude
}

// AttitudeRecords decodes all HISTORY records.
func (g *GsfFile) AttitudeRecords(fi *FileInfo) (attitude Attitude) {
    var (
        buffer []byte
    )
    n := fi.Metadata.Measurement_Counts["ATTITUDE"]
    timestamp := make([]time.Time, n)
    pitch := make([]float32, n)
    roll := make([]float32, n)
    heave := make([]float32, n)
    heading := make([]float32, n)
    //attitude = make([]Attitude, fi.Record_Counts["ATTITUDE"])

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    for _, rec := range(fi.Record_Index["ATTITUDE"]) {
        buffer = g.RecBuf(rec)
        att := DecodeAttitude(buffer)
        // attitude = append(attitude, att)

        timestamp = append(timestamp, att.Timestamp...)
        pitch = append(pitch, att.Pitch...)
        roll = append(roll, att.Roll...)
        heave = append(heave, att.Heave...)
        heading = append(heading, att.Heading...)
    }

    attitude = Attitude{
        Timestamp: timestamp,
        Pitch: pitch,
        Roll: roll,
        Heave: heave,
        Heading: heading,
    }

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

    return attitude
}
