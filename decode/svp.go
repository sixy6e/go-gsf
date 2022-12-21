package decode

import (
    "os"
    "bytes"
    "encoding/binary"
    "time"
)

type svp_base1 struct {
    Obs_seconds int32
    Obs_nano_seconds int32
    App_seconds int32
    App_nano_seconds int32
    Longitude int32
    Latitude int32
    N_points int32
}

type svp_base2 struct {
    Depth []int32
    Sound_velocity []int32
}

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

// SoundVelocityProfileRec decodes a SOUND_VELOCITY_PROFILE Record.
// It contains the values of sound velocity used in estimating individual sounding locations.
// Note: The provided samples appear to not store the position. It has been described that
// the position could be retrieved from the closest matching timestamp with that of a
// ping timestamp (within some acceptable tolerance).
func SoundVelocityProfileRec(stream *os.File, rec Record) SoundVelocityProfile {
    var (
        base1 svp_base1
        base2 svp_base2
        svp SoundVelocityProfile
    )

    buffer := make([]byte, rec.Datasize)
    _ , _ = stream.Read(buffer)
    reader := bytes.NewReader(buffer)

    _ = binary.Read(reader, binary.BigEndian, &base1)

    // 8 * 4bytes have now been read
    idx := 28

    // A previous implementation created arrays for all vars (lon, lat etc)
    // it might be better to create a single point where depth/sound velocity
    // are single elements containing an array of data
    base2.depth = make([]int32, base1.n_points)
    base2.sound_velocity = make([]int32, base1.n_points)

    reader := bytes.NewReader(buffer[idx:])
    _ = binary.Read(reader, binary.BigEndian, &base2)

    // it's not quite clear from the spec as to whether UTC is enforced
    // high potential that someone has stored local time
    svp.Observation_timestamp = time.Unix(int64(base1.obs_seconds), int64(base1.obs_nano_seconds)).UTC()
    svp.Applied_timestamp = time.Unix(int64(base1.app_seconds), int64(base1.app_nano_seconds)).UTC()

    // all the provided sample files have 0.0 for the lon and lat; WTH‽
    svp.Longitude = float64(base1.longitude) / scale2
    svp.Latitude = float64(base1.latitude) / scale2

    for i := 0; i < int(base1.N_points); i++ {
        svp.Depth[i] = float32(base2.depth[i]) / scale1
        svp.Sound_velocity[i] = float32(base2.sound_velocity[i]) / scale1
    }

    return svp
}
