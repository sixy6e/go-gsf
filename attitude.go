package gsf

import (
    "bytes"
    "encoding/binary"
    "time"
    "errors"
    "math"

    tiledb "github.com/TileDB-Inc/TileDB-Go"
)

var ErrCreateAttitudeTdb = errors.New("Error Creating Attitude TileDB Array")
var ErrWriteAttitudeTdb = errors.New("Error Writing Attitude TileDB Array")

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
    timestamp := make([]time.Time, 0, n)
    pitch := make([]float32, 0, n)
    roll := make([]float32, 0, n)
    heave := make([]float32, 0, n)
    heading := make([]float32, 0, n)
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

// attitude_tiledb_array establishes the schema and array on disk/object store.
// Timestamp could be a dimension, but for the time being it'll be a dense array
// with row (row_id) as the queryable dimension.
// At this stage, it is assumed that requests for attitude data will be the whole
// thing anyway.
func attitude_tiledb_array(file_uri string, ctx *tiledb.Context, nrows uint64) error {
    // an arbitrary choice; maybe at a future date we evaluate a good number
    tile_sz := uint64(math.Min(float64(50000), float64(nrows)))

    // array domain
    domain, err := tiledb.NewDomain(ctx)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer domain.Free()

    // setup dimension options
    // using a combination of delta filter (ascending rows) and zstandard
    dim, err := tiledb.NewDimension(ctx, "__tiledb_rows", tiledb.TILEDB_UINT64, []uint64{0, nrows - uint64(1)}, tile_sz)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer dim.Free()

    dim_filters, err := tiledb.NewFilterList(ctx)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer dim_filters.Free()

    // TODO; might be worth setting a window size
    dim_f1, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_POSITIVE_DELTA)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer dim_f1.Free()

    level := int32(16)
    dim_f2, err := ZstdFilter(ctx, level)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer dim_f2.Free()

    // attach filters to the pipeline
    err = AddFilters(dim_filters, dim_f1, dim_f2)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    err = dim.SetFilterList(dim_filters)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    err = domain.AddDimensions(dim)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // setup schema
    schema, err := tiledb.NewArraySchema(ctx, tiledb.TILEDB_DENSE)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer schema.Free()

    err = schema.SetDomain(domain)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // cell and tile ordering was an arbitrary choice
    err = schema.SetCellOrder(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    err = schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // setup attributes; timestamp, pitch, roll, heave, heading
    // just using zstd for compression. timestamp could benefit from positive delta
    // Attitude records should be in ascending, but no guarantee from these GSF files.
    zstd, err := ZstdFilter(ctx, level)
    if err != nil {
        return errors.Join(ErrCreateSvpTdb, err)
    }
    defer zstd.Free()

    ts, err := tiledb.NewAttribute(ctx, "timestamp", tiledb.TILEDB_DATETIME_NS)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer ts.Free()

    pitch, err := tiledb.NewAttribute(ctx, "pitch", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer pitch.Free()

    roll, err := tiledb.NewAttribute(ctx, "roll", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer roll.Free()

    heave, err := tiledb.NewAttribute(ctx, "heave", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer heave.Free()

    heading, err := tiledb.NewAttribute(ctx, "heading", tiledb.TILEDB_FLOAT32)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer heading.Free()

    // compression filter pipeline
    attr_filts, err := tiledb.NewFilterList(ctx)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer attr_filts.Free()

    err = attr_filts.AddFilter(zstd)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // attach filter pipeline to attrs
    err = AttachFilters(attr_filts, ts, pitch, roll, heave, heading)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // attach attrs to the schema
    err = schema.AddAttributes(ts, pitch, roll, heave, heading)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    // finally, create the empty array on disk, object store, etc
    array, err := tiledb.NewArray(ctx, file_uri)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }
    defer array.Free()

    err = array.Create(schema)
    if err != nil {
        return errors.Join(ErrCreateAttitudeTdb, err)
    }

    return nil
}

// ToTileDB writes the Attitude data to a TileDB array.
// Timestamp could be a dimension, but for the time being it'll be a dense array
// with row (row_id) as the queryable dimension.
// At this stage, it is assumed that requests for attitude data will be the whole
// thing anyway.
// Column structure:
// [__tiledb_rows (dim), timestamp (attr), pitch (attr), roll (attr), heave (attr), heading (attr)].
func (a *Attitude) ToTileDB(file_uri string, config_uri string) error {
    var config *tiledb.Config
    var err error

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

    nrows := uint64(len(a.Timestamp))
    err = attitude_tiledb_array(file_uri, ctx, nrows)
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
        return errors.Join(ErrWriteAttitudeTdb, err)
    }
    defer query.Free()

    err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    temp_data := make([]int64, nrows)
    for i := uint64(0); i < nrows; i++ {
        temp_data[i] = a.Timestamp[i].UnixNano()
    }
    _, err = query.SetDataBuffer("timestamp", temp_data)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    _, err = query.SetDataBuffer("pitch", a.Pitch)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    _, err = query.SetDataBuffer("roll", a.Roll)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    _, err = query.SetDataBuffer("heave", a.Heave)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    _, err = query.SetDataBuffer("heading", a.Heading)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    // define the subarray (dim coordinates that we'll write into)
    subarr, err := array.NewSubarray()
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }
    defer subarr.Free()

    rng := tiledb.MakeRange(uint64(0), nrows - uint64(1))
    subarr.AddRangeByName("__tiledb_rows", rng)
    err = query.SetSubarray(subarr)
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    // write the data flush
    err = query.Submit()
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
    }

    err = query.Finalize()
    if err != nil {
        return errors.Join(ErrWriteAttitudeTdb, err)
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
