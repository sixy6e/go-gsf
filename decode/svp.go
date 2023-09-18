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

type svp_hdr struct {
    Observation_timestamp time.Time
    Applied_timestamp time.Time
    Longitude float64
    Latitude float64
    N_points uint64
}

func svp_header(reader *bytes.Reader) (hdr svp_hdr) {
    var (
        base struct {
            Obs_seconds int32
            Obs_nano_seconds int32
            App_seconds int32
            App_nano_seconds int32
            Longitude int32
            Latitude int32
            N_points int32
        }
    )

    _ = binary.Read(reader, binary.BigEndian, &base)

    // it's not quite clear from the spec as to whether UTC is enforced
    // high potential that someone has stored local time
    hdr.Observation_timestamp = time.Unix(int64(base.Obs_seconds), int64(base.Obs_nano_seconds)).UTC()
    hdr.Applied_timestamp = time.Unix(int64(base.App_seconds), int64(base.App_nano_seconds)).UTC()

    // all the provided sample files have 0.0 for the lon and lat; WTH‽
    hdr.Longitude = float64(float32(base.Longitude) / SCALE2)
    hdr.Latitude = float64(float32(base.Latitude) / SCALE2)

    hdr.N_points = uint64(base.N_points)

    return hdr
}

// DecodeSoundVelocityProfile is a constructor for SoundVelocityProfile by decoding
// a SOUND_VELOCITY_PROFILE Record.
// It contains the values of sound velocity used in estimating individual sounding locations.
// Note: The provided samples appear to not store the position. It has been described that
// the position could be retrieved from the closest matching timestamp with that of a
// ping timestamp (within some acceptable tolerance).
func DecodeSoundVelocityProfile(buffer []byte) SoundVelocityProfile {
    var (
        base struct {
            Depth []int32
            Sound_velocity []int32
        }
        svp SoundVelocityProfile
        i uint64
    )

    reader := bytes.NewReader(buffer)

    // _ = binary.Read(reader, binary.BigEndian, &base1)
    hdr := svp_header(reader)

    // 7 * 4bytes have now been read
    idx := 28

    // A previous implementation created arrays for all vars (lon, lat etc)
    // it might be better to create a single point where depth/sound velocity
    // are single elements containing an array of data
    base.Depth = make([]int32, hdr.N_points)
    base.Sound_velocity = make([]int32, hdr.N_points)

    reader = bytes.NewReader(buffer[idx:])
    _ = binary.Read(reader, binary.BigEndian, &base)

    // all the provided sample files have 0.0 for the lon and lat; WTH‽
    svp.Longitude = hdr.Longitude
    svp.Latitude = hdr.Latitude

    for i = 0; i < hdr.N_points; i++ {
        svp.Depth[i] = float32(base.Depth[i]) / SCALE1
        svp.Sound_velocity[i] = float32(base.Sound_velocity[i]) / SCALE1
    }

    return svp
}

// SvpRecords decodes all SOUND_VELOCITY_PROFILE records.
func (g *GsfFile) SoundVelocityProfileRecords(fi *FileInfo) (svp []SoundVelocityProfile) {
    var (
        buffer []byte
    )
    svp = make([]SoundVelocityProfile, fi.Record_Counts["SOUND_VELOCITY_PROFILE"])

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    for _, rec := range(fi.Record_Index["SOUND_VELOCITY_PROFILE"]) {
        buffer = g.RecBuf(rec)
        sv_p := DecodeSoundVelocityProfile(buffer)
        svp = append(svp, sv_p)
    }

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

    return svp
}
