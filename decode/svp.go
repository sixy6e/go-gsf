package decode

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
)

// SoundVelocityProfile contains the values of sound velocoty used in estimating
// individual sounding locations.
// It consists of; the time the prodile was observed, the time it was introduced into the
// sounding location procedure, the position of the observation, the number of points in
// the profile, and the individual points expressed as depth and sound velocity.
type SoundVelocityProfile struct {
    Observation_timestamp time.Time
    Applied_timestamp time.Time
    Longitude float64
    Latitude float64
    Depth []float32
    Sound_velocity []float32
}

// NewSoundVelocityProfile is a constructor for SoundVelocityProfile by decoding
// a SOUND_VELOCITY_PROFILE Record.
// It contains the values of sound velocity used in estimating individual sounding locations.
// Note: The provided samples appear to not store the position. It has been described that
// the position could be retrieved from the closest matching timestamp with that of a
// ping timestamp (within some acceptable tolerance).
func NewSoundVelocityProfile(buffer []byte) *SoundVelocityProfile {
    var (
        base1 svp_base1
        base2 svp_base2
        base1 struct {
            Obs_seconds int32
            Obs_nano_seconds int32
            App_seconds int32
            App_nano_seconds int32
            Longitude int32
            Latitude int32
            N_points int32
        }
        base2 struct {
            Depth []int32
            Sound_velocity []int32
        }
        svp SoundVelocityProfile
        i int32
    )

    reader := bytes.NewReader(buffer)

    _ = binary.Read(reader, binary.BigEndian, &base1)

    // 7 * 4bytes have now been read
    idx := 28

    // A previous implementation created arrays for all vars (lon, lat etc)
    // it might be better to create a single point where depth/sound velocity
    // are single elements containing an array of data
    base2.Depth = make([]int32, base1.N_points)
    base2.Sound_velocity = make([]int32, base1.N_points)

    reader = bytes.NewReader(buffer[idx:])
    _ = binary.Read(reader, binary.BigEndian, &base2)

    // it's not quite clear from the spec as to whether UTC is enforced
    // high potential that someone has stored local time
    svp.Observation_timestamp = time.Unix(int64(base1.Obs_seconds), int64(base1.Obs_nano_seconds)).UTC()
    svp.Applied_timestamp = time.Unix(int64(base1.App_seconds), int64(base1.App_nano_seconds)).UTC()

    // all the provided sample files have 0.0 for the lon and lat; WTH‽
    svp.Longitude = float64(float32(base1.Longitude) / SCALE2)
    svp.Latitude = float64(float32(base1.Latitude) / SCALE2)

    for i = 0; i < base1.N_points; i++ {
        svp.Depth[i] = float32(base2.Depth[i]) / SCALE1
        svp.Sound_velocity[i] = float32(base2.Sound_velocity[i]) / SCALE1
    }

    return new(svp)
}
