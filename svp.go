package gsf

import (
    // "os"
    "bytes"
    "encoding/binary"
    "time"
    "errors"
    // "fmt"

    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

var ErrCreateSvpTdb = errors.New("Error Creating SVP TileDB Array")
var ErrWriteSvpTdb = errors.New("Error Writing SVP TileDB Array")

// SoundVelocityProfile contains the values of sound velocity used in estimating
// individual sounding locations.
// It consists of; the time the profile was observed, the time it was introduced into the
// sounding location procedure, the position of the observation, the number of points in
// the profile, and the individual points expressed as depth and sound velocity.
// While most of the sample files only contained a single SVP Record, in order to
// cater for generality in that other data providers may include 100's or even 1000's
// of SVP records. We'll construct the data in a way that it mimics a single row of data,
// and depth and sound_velocity are variable length fields.
// The downside, is that when serialising as an n-Dimensional construct, we have a single
// row, indexed by row number (potentially we could use lon/lat dimensional axes) which
// isn't (by some standards), efficient storage use.
// We could replicate lon/lat and timestamps by n times where n = len(depth) field.
// But this isn't efficient either, by compression algorithms will take care of that.
// The individual fields (lon, lat, timestamps) could merely be attached metadata, but
// that means querying capability is impeded.
// So for simplicity, if it is a single record, or a million records, we'll use the same
// data structure; array and array of arrays (variable length arrays).
type SoundVelocityProfile struct {
    Observation_timestamp []time.Time
    Applied_timestamp []time.Time
    Longitude []float64
    Latitude []float64
    Depth [][]float32
    Sound_velocity [][]float32
    depth []float32
    sound_velocity []float32
    n_points uint64
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
        // base struct {
        //     Depth []int32
        //     Sound_velocity []int32
        // }
        base []int32
        depth_f32 []float32
        svp_f32 []float32
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
    // base.Depth = make([]int32, hdr.N_points)
    // base.Sound_velocity = make([]int32, hdr.N_points)
    base = make([]int32, 2 * hdr.N_points)

    reader = bytes.NewReader(buffer[idx:])
    // err := binary.Read(reader, binary.BigEndian, &base.Depth)
    err := binary.Read(reader, binary.BigEndian, &base)
    if err != nil {
        panic(err)
    }

    svp.Observation_timestamp = []time.Time{hdr.Observation_timestamp}
    svp.Applied_timestamp = []time.Time{hdr.Applied_timestamp}

    // all the provided sample files have 0.0 for the lon and lat; WTH‽
    svp.Longitude = []float64{hdr.Longitude}
    svp.Latitude = []float64{hdr.Latitude}

    depth_f32 = make([]float32, 0, hdr.N_points)
    svp_f32 = make([]float32, 0, hdr.N_points)

    for i = 0; i < 2*hdr.N_points; i += 2 {
        depth_f32 = append(depth_f32, float32(base[i]) / SCALE2)
        svp_f32 = append(svp_f32, float32(base[i+1]) / SCALE2)
    }

    svp.depth = depth_f32
    svp.sound_velocity = svp_f32

    svp.Depth = [][]float32{svp.depth}
    svp.Sound_velocity = [][]float32{svp.sound_velocity}

    svp.n_points = hdr.N_points

    return svp
}

// SoundVelocityProfileRecords decodes all SOUND_VELOCITY_PROFILE records.
func (g *GsfFile) SoundVelocityProfileRecords(fi *FileInfo) (svp SoundVelocityProfile) {
    var (
        buffer []byte
        obs_time []time.Time
        app_time []time.Time
        lon []float64
        lat []float64
        depth []float32
        velocity []float32
        count []uint64
        depth_nd_slices [][]float32
        velocity_nd_slices [][]float32
    )
    rec_counts := fi.Record_Counts["SOUND_VELOCITY_PROFILE"]

    obs_time = make([]time.Time, 0, rec_counts)
    app_time = make([]time.Time, 0, rec_counts)
    lon = make([]float64, 0, rec_counts)
    lat = make([]float64, 0, rec_counts)
    count = make([]uint64, 0, rec_counts)

    depth = make([]float32, 0, fi.Metadata.Measurement_Counts["SOUND_VELOCITY_PROFILE"])
    velocity = make([]float32, 0, fi.Metadata.Measurement_Counts["SOUND_VELOCITY_PROFILE"])

    // get the original starting point so we can jump back when done
    original_pos, _ := Tell(g.Stream)

    for _, rec := range(fi.Record_Index["SOUND_VELOCITY_PROFILE"]) {
        buffer = g.RecBuf(rec)
        sv_p := DecodeSoundVelocityProfile(buffer)

        obs_time = append(obs_time, sv_p.Observation_timestamp...)
        app_time = append(app_time, sv_p.Applied_timestamp...)
        lon = append(lon, sv_p.Longitude...)
        lat = append(lat, sv_p.Latitude...)

        // each SVP record will only be a 2D slice containing a single row
        depth = append(depth, sv_p.Depth[0]...)
        velocity = append(velocity, sv_p.Sound_velocity[0]...)

        count = append(count, sv_p.n_points)
    }

    // generate the 2D slices that are ideally views of the 1D slice
    depth_nd_slices = make([][]float32, rec_counts)
    velocity_nd_slices = make([][]float32, rec_counts)
    start_idx := uint64(0)
    end_idx := uint64(0)
    for i, val := range(count) {
        end_idx = start_idx + uint64(val)
        depth_nd_slices[i] = depth[start_idx:end_idx]
        velocity_nd_slices[i] = velocity[start_idx:end_idx]
        start_idx = end_idx
    }

    // exported 1D slices
    svp.Observation_timestamp = obs_time
    svp.Applied_timestamp = app_time
    svp.Longitude = lon
    svp.Latitude = lat

    // unexported backend 1D slices
    svp.depth = depth
    svp.sound_velocity = velocity

    svp.n_points = uint64(len(depth))

    // exported 2D slice containing views on the 1D slice
    svp.Depth = depth_nd_slices
    svp.Sound_velocity = velocity_nd_slices

    // reset file position
    _, _ = g.Stream.Seek(original_pos, 0)

    return svp
}

// svp_tiledb_array establishes the schema and array on disk/object store.
func svp_tiledb_array(file_uri string, ctx *tiledb.Context, nrows uint64) error {
    // an arbitrary choice; maybe at a future date we evaluate a good number
    // the samples provided so far indicate 1 or 2 rows (points of acquisition)
    // so making the tilesize the same as the number of rows will be fine until
    // we start getting hundreds of rows
    // tile_sz := uint64(1)
    tile_sz := nrows

    // array domain
    domain, err := tiledb.NewDomain(ctx)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer domain.Free()

    // setup dimension options
    // using a combination of delta filter (ascending rows) and zstandard
    dim, err := tiledb.NewDimension(ctx, "__tiledb_rows", tiledb.TILEDB_UINT64, []uint64{0, nrows - uint64(1)}, tile_sz)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer dim.Free()

    dim_filters, err := tiledb.NewFilterList(ctx)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer dim_filters.Free()

    // TODO; might be worth setting a window size
    dim_f1, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_POSITIVE_DELTA)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer dim_f1.Free()

    level := int32(16)
    dim_f2, err := ZstdFilter(ctx, level)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer dim_f2.Free()

    // attach dim filters to the pipeline
    err = AddFilters(dim_filters, dim_f1, dim_f2)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    err = dim.SetFilterList(dim_filters)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    err = domain.AddDimensions(dim)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // setup schema
    schema, err := tiledb.NewArraySchema(ctx, tiledb.TILEDB_DENSE)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer schema.Free()

    err = schema.SetDomain(domain)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // cell and tile ordering was an arbitrary choice
    err = schema.SetCellOrder(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    err = schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // setup attributes:
    // observation_timestamp, applied_timestamp, longitude, latitude, depth, sound_velocity
    // just using zstd for compression. timestamps could benefit from positive delta,
    // but additional work is required for evaluation
    zstd, err := ZstdFilter(ctx, level)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer zstd.Free()

    obs_ts, err := tiledb.NewAttribute(ctx, "observation_timestamp", tiledb.TILEDB_DATETIME_NS)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer obs_ts.Free()

    app_ts, err := tiledb.NewAttribute(ctx, "applied_timestamp", tiledb.TILEDB_DATETIME_NS)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer app_ts.Free()

    lon, err := tiledb.NewAttribute(ctx, "longitude", tiledb.TILEDB_FLOAT64)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer lon.Free()

    lat, err := tiledb.NewAttribute(ctx, "latitude", tiledb.TILEDB_FLOAT64)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer lat.Free()

    depth, err := tiledb.NewAttribute(ctx, "depth", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer depth.Free()

    velocity, err := tiledb.NewAttribute(ctx, "sound_velocity", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer velocity.Free()

    // define depth and velocity as variable length arrays
    depth.SetCellValNum(tiledb.TILEDB_VAR_NUM)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    velocity.SetCellValNum(tiledb.TILEDB_VAR_NUM)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // compression filter pipeline
    attr_filts, err := tiledb.NewFilterList(ctx)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer attr_filts.Free()

    err = attr_filts.AddFilter(zstd)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // attach filter pipeline to attrs
    err = AttachFilters(attr_filts, obs_ts, app_ts, lon, lat, depth, velocity)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // attach attrs to the schema
    err = schema.AddAttributes(obs_ts, app_ts, lon, lat, depth, velocity)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    // finally, create the empty array on disk, object store, etc
    array, err := tiledb.NewArray(ctx, file_uri)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer array.Free()

    err = array.Create(schema)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }

    return nil
}

// ToTileDB writes the SoundVelocityProfile data to a TileDB array.
// Without knowing the general access and usage patterns, it is hard to define
// a specific structure. In order to be generic in that the GSF file may contain
// multiple SVP records, we'll output the data as a dense TileDB array, using
// row id's as the dimensional access.
// We could replicate the longitude, latitude, and timestamps for n*depths,
// but it's better to wait and modify the data structure if usage patterns necessitate it.
// Column structure:
// [__tiledb_rows (dim), observation_timestamp (attr), applied_timestamp (attr), longitude (attr), latitude (attr), depth (attr), sound_velocity (attr)].
// The depth and sound_velocity attributes are variable length arrays that contain the
// profile for the specific acquisition defined by observation timestamp, longitude and latitude.
func (s *SoundVelocityProfile) ToTileDB(file_uri string, config_uri string) error {
    var (
        config *tiledb.Config
        err error
        arr_offset []uint64
        offset uint64
        bytes_val uint64
    )

    // get a generic config if no path provided
    if config_uri == "" {
        config, err = tiledb.NewConfig()
        if err != nil {
            return err
        }
    } else {
        config, err = tiledb.LoadConfig(config_uri)
        if err != nil {
            return err
        }
    }

    defer config.Free()

    ctx, err := tiledb.NewContext(config)
    if err != nil {
        return err
    }
    defer ctx.Free()

    nrows := uint64(len(s.Observation_timestamp))
    err = svp_tiledb_array(file_uri, ctx, nrows)
    if err != nil {
        return err
    }

    // open the array for writing the attitude data
    array, err := ArrayOpen(ctx, file_uri, tiledb.TILEDB_WRITE)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }
    defer array.Free()
    defer array.Close()

    // query construction
    query, err := tiledb.NewQuery(ctx, array)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }
    defer query.Free()

    err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    // convert time.Time arrays to int64 UnixNano time
    obs_time := make([]int64, nrows)
    app_time := make([]int64, nrows)
    for i := uint64(0); i < nrows; i++ {
        obs_time[i] = s.Observation_timestamp[i].UnixNano()
        app_time[i] = s.Applied_timestamp[i].UnixNano()
    }

    _, err = query.SetDataBuffer("observation_timestamp", obs_time)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetDataBuffer("applied_timestamp", app_time)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetDataBuffer("longitude", s.Longitude)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetDataBuffer("latitude", s.Latitude)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    // variable length attrs
    // need to define a 1D offset array for variable length attributes
    // eg data = []int32{1, 1, 2, 3, 3, 3, 4}; offset = []uint64{0, 8, 12, 24}
     arr_offset = make([]uint64, nrows)
     offset = uint64(0)
     bytes_val = uint64(4)  // may look confusing with uint64, so 4*bytes for float32
     for i := uint64(0); i < nrows; i++ {
         length := uint64(len(s.Depth[i]))
         arr_offset[i] = offset * bytes_val
         offset += length * bytes_val
     }
    _, err = query.SetOffsetsBuffer("depth", arr_offset)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetOffsetsBuffer("sound_velocity", arr_offset)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetDataBuffer("depth", s.depth)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    _, err = query.SetDataBuffer("sound_velocity", s.sound_velocity)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    // define the subarray (dim coordinates that we'll write into)
    subarr, err := array.NewSubarray()
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }
    defer subarr.Free()

    rng := tiledb.MakeRange(uint64(0), nrows - uint64(1))
    subarr.AddRangeByName("__tiledb_rows", rng)
    err = query.SetSubarray(subarr)
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    // write the data flush
    err = query.Submit()
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    err = query.Finalize()
    if err != nil {
        return errors.Join(ErrWriteSvpTdb, err)
    }

    // attach some metadata to preserve python pandas functionality
    md := map[string]string{"__tiledb_rows": "uint64"}
    jsn, err := JsonDumps(md)
    if err != nil {
        return err
    }
    err = array.PutMetadata("__pandas_index_dims", jsn)

    return nil
}
