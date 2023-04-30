package decode

import (
    "os"
    "bytes"
    "encoding/binary"
    "time"
)

type attitude_base1 struct {
    Seconds int32
    Nano_seconds int32
    Measurements int16
}

type attitude_base2 struct {
    Time_offset int32
    Pitch int32
    Roll int32
    Heave int32
    Heading int32
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

// AttitudeRec decodes an Attitude Record which contains the measurements
// as reported by the vessel attitude sensor.
// Fields include: Timestamp, Pitch, Roll, Heave and Heading.
func AttitudeRec(stream *os.File, rec Record) Attitude {
    var (
        idx int = 0
        base1 attitude_base1
        base2 attitude_base2
        offset time.Duration
    )

    buffer := make([]byte, rec.Datasize)
    _ , _ = stream.Read(buffer)
    reader := bytes.NewReader(buffer)

    _ = binary.Read(reader, binary.BigEndian, &base1)
    idx += 10

    acq_time := time.Unix(int64(base1.Seconds), int64(base1.Nano_seconds)).UTC()

    attitude := Attitude{
        Timestamp: make([]time.Time, base1.Measurements),
        Pitch: make([]float32, base1.Measurements),
        Roll: make([]float32, base1.Measurements),
        Heave: make([]float32, base1.Measurements),
        Heading: make([]float32, base1.Measurements),
    }

    for i:= int16(0); i < base1.Measurements; i++ {
        reader = bytes.NewReader(buffer[idx:])  // probably superfluous in creating a new reader
        _ = binary.Read(reader, binary.BigEndian, &base2)

        // haven't looked deep into why, but 1_000_000 and nanoseconds worked
        // original C code did some funky stuff in an internal function for determining time
        offset = time.Duration(int64(base2.Time_offset) * SCALE4)
        attitude.Timestamp[i] = acq_time.Add(time.Nanosecond * offset)
        attitude.Pitch[i] = float32(base2.Pitch) / SCALE2
        attitude.Roll[i] = float32(base2.Roll) / SCALE2
        attitude.Heave[i] = float32(base2.Heave) / SCALE2
        attitude.Heading[i] = float32(base2.Heading) / SCALE2

        idx += 10
    }

    return attitude
}
