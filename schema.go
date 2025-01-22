package gsf

import (
	"errors"
	"math"
	"reflect"
	"strings"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
	"github.com/samber/lo"
	stgpsr "github.com/yuin/stagparser"
)

var ErrCreateAttributeTdb = errors.New("Error Creating Attribute for TileDB Array")
var ErrCreateMdDenseTdb = errors.New("Error Creating Dense Metadata TileDB Array")
var ErrCreateBeamSparseTdb = errors.New("Error Creating Beam Sparse TileDB Array")
var ErrCreateSchemaTdb = errors.New("Error Creating TileDB Schema")
var ErrCreateDimTdb = errors.New("Error Creating TileDB Dimension")

// pascalCase convert a string separated by underscores into
// PascalCase. For example, ALPHA_BETA_GAMMA -> AlphaBetaGamma.
func pascalCase(name string) (result string) {
	result = ""
	split := strings.Split(name, "_")

	for _, v := range split {
		low := strings.ToLower(v)
		result += strings.ToUpper(string(low[0])) + low[1:]
	}

	return result
}

func fieldNames(t any) (names []string) {
	names = make([]string, 0, 10)

	btype := reflect.TypeOf(t)
	for i := 0; i < btype.NumField(); i++ {
		if btype.Field(i).IsExported() {
			names = append(names, btype.Field(i).Name)
		}
	}
	return names
}

// chunkedStructSlices is a helper func for initialising structs containing
// slices to a defined capacity. For example PingData where the slices will be of
// total number of beams in capacity. Or for SensorMetadata which will be of
// npings in capacity. This ideally should reduce any overhead in reallocation
// during appending.
// However, unexported fields won't be handled. Will need to handle those outside
// on a case by case basis.
func chunkedStructSlices(t any, length int) error {
	values := reflect.ValueOf(t).Elem()
	types := reflect.TypeOf(t).Elem()
	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		t := field.Type()
		if types.Field(i).IsExported() {
			field.Set(reflect.MakeSlice(t, 0, length))
		}
	}

	return nil
}

// chunkedBeamArray is a helper func for initialising structs containing
// slices to a defined capacity. For example PingData where the slices will be of
// total number of beams in capacity.
// This ideally should reduce any overhead in reallocation during appending.
// Only those fields listed in the parameter beam_names will be set.
func chunkedBeamArray(t any, length int, beam_names []string) error {
	values := reflect.ValueOf(t).Elem()
	for _, v := range beam_names {
		field := values.FieldByName(v)
		ftype := field.Type()
		field.Set(reflect.MakeSlice(ftype, 0, length))
	}
	return nil
}

func schemaAttrs(t any, schema *tiledb.ArraySchema, ctx *tiledb.Context) error {
	var (
		field_tdb_defs map[string]stgpsr.Definition
		def            stgpsr.Definition
		status         bool
	)
	values := reflect.ValueOf(t).Elem()
	types := values.Type()
	filt_defs, _ := stgpsr.ParseStruct(t, "filters")
	tdb_defs, _ := stgpsr.ParseStruct(t, "tiledb")

	// process every field in the struct
	for i := 0; i < values.NumField(); i++ {
		name := types.Field(i).Name

		if !types.Field(i).IsExported() {
			continue
		}

		field_filt_defs := filt_defs[name]

		// a mapping just seemed easier to pull required defs
		// rather than a simple listing
		field_tdb_defs = make(map[string]stgpsr.Definition)
		for _, v := range tdb_defs[name] {
			field_tdb_defs[v.Name()] = v
		}

		// pull the field type and ignore dimension fields
		def, status = field_tdb_defs["ftype"]
		if status == false {
			errf := errors.New("Field: " + name)
			return errors.Join(ErrCreateAttributeTdb, errors.New("ftype tag not found"), errf)
		}
		ftype, _ := def.Attribute("ftype")
		if ftype == "dim" {
			// ignore dimensions
			continue
		}

		err := CreateAttr(name, field_filt_defs, field_tdb_defs, schema, ctx)
		if err != nil {
			return errors.Join(ErrCreateAttributeTdb, err)
		}
	}
	return nil
}

// basePidSchema sets up a base schema using the Ping ID as the dimensional axis.
// Doesn't attach any attriubtes. Used for the PingHeaders, SensorMetadata and
// SensorImageryMetadata structures.
func basePidSchema(ctx *tiledb.Context, npings uint64) (*tiledb.ArraySchema, error) {
	// an arbitrary choice; maybe at a future date we evaluate a good number
	tile_sz := uint64(math.Min(float64(50000), float64(npings)))

	// array domain
	domain, err := tiledb.NewDomain(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer domain.Free()

	// setup dimension options
	// using a combination of delta filter (ascending rows) and zstandard
	dim, err := tiledb.NewDimension(ctx, "PING_ID", tiledb.TILEDB_UINT64, []uint64{0, npings - uint64(1)}, tile_sz)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim.Free()

	dim_filters, err := tiledb.NewFilterList(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_filters.Free()

	// TODO; might be worth setting a window size
	dim_f1, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_POSITIVE_DELTA)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_f1.Free()

	level := int32(16)
	dim_f2, err := ZstdFilter(ctx, level)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_f2.Free()

	// attach filters to the pipeline
	err = AddFilters(dim_filters, dim_f1, dim_f2)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	err = dim.SetFilterList(dim_filters)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = domain.AddDimensions(dim)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	// setup schema
	schema, err := tiledb.NewArraySchema(ctx, tiledb.TILEDB_DENSE)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = schema.SetDomain(domain)
	if err != nil {
		return nil, errors.Join(ErrCreateAttitudeTdb, err)
	}

	// cell and tile ordering was an arbitrary choice
	err = schema.SetCellOrder(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	return schema, nil
}

// baseLonLatSchema sets up a base schema using X and Y as the dimensional axes,
// where X and Y are longitude and latitude coordinatges.
// Doesn't attach any attributes.
// Used for the beam array data and if it exists, the brb intensity data.
// The schema is set to allow duplicates, hilbert for cell ordering, row-major
// for tile ordering.
func baseLonLatSchema(ctx *tiledb.Context) (schema *tiledb.ArraySchema, err error) {
	// array domain
	domain, err := tiledb.NewDomain(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer domain.Free()

	tile_sz := float64(1000)
	min_f64 := math.MaxFloat64 * -1

	// setup lon/lat (X/Y) dimensions
	xdim, err := tiledb.NewDimension(ctx, "X", tiledb.TILEDB_FLOAT64, []float64{min_f64, math.MaxFloat64}, tile_sz)
	if err != nil {
		errdim := errors.New("Error Creating Dimension X")
		return nil, errors.Join(ErrCreateAttributeTdb, err, errdim)
	}
	defer xdim.Free()

	ydim, err := tiledb.NewDimension(ctx, "Y", tiledb.TILEDB_FLOAT64, []float64{min_f64, math.MaxFloat64}, tile_sz)
	if err != nil {
		errdim := errors.New("Error Creating Dimension Y")
		return nil, errors.Join(ErrCreateAttributeTdb, err, errdim)
	}
	defer ydim.Free()

	dim_filters, err := tiledb.NewFilterList(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_filters.Free()

	level := int32(16)
	dim_filt, err := ZstdFilter(ctx, level)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_filt.Free()

	// attach dimension filters to the pipeline
	err = AddFilters(dim_filters, dim_filt)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = xdim.SetFilterList(dim_filters)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = ydim.SetFilterList(dim_filters)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = domain.AddDimensions(xdim, ydim)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	// setup schema
	schema, err = tiledb.NewArraySchema(ctx, tiledb.TILEDB_SPARSE)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	err = schema.SetDomain(domain)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	err = schema.SetCapacity(100_000)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	err = schema.SetCellOrder(tiledb.TILEDB_HILBERT)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	err = schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	err = schema.SetAllowsDups(true)
	if err != nil {
		return nil, errors.Join(ErrCreateSchemaTdb, err)
	}

	return schema, nil
}

func beamAttachAttrs(schema *tiledb.ArraySchema, ctx *tiledb.Context, beam_subrecords []string, contains_intensity bool) (err error) {
	var (
		field_tdb_defs map[string]stgpsr.Definition
		def            stgpsr.Definition
		status         bool
	)

	// handle X & Y and PingNumber and BeamNumber as attributes depending on whether
	// we're dealing with a dense or sparse array
	dense, err := schema.Type()
	if err != nil {
		return err
	}
	if dense == tiledb.TILEDB_DENSE {
		err = schemaAttrs(&XY{}, schema, ctx)
		if err != nil {
			err_pbn := errors.New("Error attaching X & Y attributes")
			return errors.Join(err, ErrCreateAttributeTdb, err_pbn)
		}
	} else {
		err = schemaAttrs(&PingBeamNumbers{}, schema, ctx)
		if err != nil {
			err_pbn := errors.New("Error attaching PingNumber & BeamNumber attributes")
			return errors.Join(err, ErrCreateAttributeTdb, err_pbn)
		}
	}

	ba := BeamArray{}
	beam_names := make([]string, len(beam_subrecords))

	// cleanup subrecord names to match the BeamArray fields names
	for k, v := range beam_subrecords {
		beam_names[k] = pascalCase(v)
	}

	filt_defs, _ := stgpsr.ParseStruct(ba, "filters")
	tdb_defs, _ := stgpsr.ParseStruct(ba, "tiledb")

	// processing the beam array subrecords
	for _, name := range beam_names {

		// ignore intensity series as it needs to be handled by a separate type
		if name == "IntensitySeries" {
			continue
		}

		field_filt_defs := filt_defs[name]

		field_tdb_defs = make(map[string]stgpsr.Definition)
		for _, v := range tdb_defs[name] {
			field_tdb_defs[v.Name()] = v
		}

		// pull the field type and ignore dimension fields
		def, status = field_tdb_defs["ftype"]
		if status == false {
			return errors.Join(ErrCreateAttributeTdb, errors.New("ftype tag not found"))
		}
		ftype, _ := def.Attribute("ftype")
		if ftype == "dim" {
			// ignore dimensions
			continue
		}

		err := CreateAttr(name, field_filt_defs, field_tdb_defs, schema, ctx)
		if err != nil {
			return errors.Join(ErrCreateAttributeTdb, err)
		}
	}

	// processing the brb intensity data
	if contains_intensity {
		err = schemaAttrs(&BrbIntensity{}, schema, ctx)
		if err != nil {
			err_brb := errors.New("Error attaching BrbIntensity attributes")
			return errors.Join(err, ErrCreateAttributeTdb, err_brb)
		}
	}

	// processing the basic ping info (ping id, beam id)
	// err = schemaAttrs(&PingBeamNumbers{}, schema, ctx)
	// if err != nil {
	// 	err_pbn := errors.New("Error attaching PingBeamNumbers attributes")
	// 	return errors.Join(err, ErrCreateAttributeTdb, err_pbn)
	// }

	return nil
}

func phTdbArray(ctx *tiledb.Context, array_uri string, npings uint64) error {
	schema, err := basePidSchema(ctx, npings)
	if err != nil {
		return err
	}
	defer schema.Free()

	err = schemaAttrs(&PingHeaders{}, schema, ctx)
	if err != nil {
		errn := errors.New("Error creating PingHeaders attributes")
		return errors.Join(err, errn)
	}

	err = schema.Check()
	if err != nil {
		errn := errors.New("Error checking PingHeaders schema")
		return errors.Join(err, errn)
	}

	array, err := tiledb.NewArray(ctx, array_uri)
	if err != nil {
		errn := errors.New("Error creating PingHeaders array")
		return errors.Join(err, errn)
	}
	defer array.Free()

	err = array.Create(schema)
	if err != nil {
		errn := errors.New("Error creating PingHeaders array")
		return errors.Join(err, errn)
	}

	// attach some metadata to preserve python pandas functionality
	md := map[string]string{"PING_ID": "uint64"}
	key := "__pandas_index_dims"
	err = WriteArrayMetadata(ctx, array_uri, key, md)
	if err != nil {
		return err
	}

	return nil
}

func senTdbArray(ctx *tiledb.Context, array_uri string, npings uint64, sensor_id SubRecordID) error {
	schema, err := basePidSchema(ctx, npings)
	if err != nil {
		errn := errors.New("Error creating base schema for SensorMetadata")
		return errors.Join(err, errn)
	}
	defer schema.Free()

	smd := SensorMetadata{}
	err = smd.attachAttrs(schema, ctx, sensor_id)
	if err != nil {
		return err
	}

	err = schema.Check()
	if err != nil {
		errn := errors.New("Error checking SensorMetadata TileDB schema")
		return errors.Join(err, errn)
	}

	array, err := tiledb.NewArray(ctx, array_uri)
	if err != nil {
		errn := errors.New("Error creating SensorMetadata TileDB array")
		return errors.Join(err, errn)
	}
	defer array.Free()

	err = array.Create(schema)
	if err != nil {
		errn := errors.New("Error creating SensorMetadata TileDB array")
		return errors.Join(err, errn)
	}

	// attach some metadata to preserve python pandas functionality
	md := map[string]string{"PING_ID": "uint64"}
	key := "__pandas_index_dims"
	err = WriteArrayMetadata(ctx, array_uri, key, md)
	if err != nil {
		return err
	}

	return nil
}

func senImgTdbArray(ctx *tiledb.Context, array_uri string, npings uint64, sensor_id SubRecordID) error {
	schema, err := basePidSchema(ctx, npings)
	if err != nil {
		err_sen := errors.New("Error creating base schema for SensorImageryMetadata")
		return errors.Join(err, err_sen)
	}
	defer schema.Free()

	simd := SensorImageryMetadata{}
	err = simd.attachAttrs(schema, ctx, sensor_id)
	if err != nil {
		return err
	}

	err = schema.Check()
	if err != nil {
		errn := errors.New("Error checking SensorImageryMetadata TileDB schema")
		return errors.Join(err, errn)
	}

	array, err := tiledb.NewArray(ctx, array_uri)
	if err != nil {
		errn := errors.New("Error creating SensorImageryMetadata TileDB array")
		return errors.Join(err, errn)
	}
	defer array.Free()

	err = array.Create(schema)
	if err != nil {
		errn := errors.New("Error creating SensorImageryMetadata TileDB array")
		return errors.Join(err, errn)
	}

	// attach some metadata to preserve python pandas functionality
	md := map[string]string{"PING_ID": "uint64"}
	key := "__pandas_index_dims"
	err = WriteArrayMetadata(ctx, array_uri, key, md)
	if err != nil {
		return err
	}

	return nil
}

func beamTdbArray(ctx *tiledb.Context, array_uri string, beam_subrecords []string, contains_intensity, dense_bd bool, npings uint64, max_beams uint16) error {
	var (
		schema *tiledb.ArraySchema
		err    error
		md     map[string]string
	)

	if dense_bd {
		schema, err = basePingBeamSchema(ctx, npings, max_beams)
		if err != nil {
			errn := errors.New("Error creating base dense schema for beam array")
			return errors.Join(err, errn)
		}
		md = map[string]string{"PingNumber": "uint64", "BeamNumber": "uint64"}
	} else {
		schema, err = baseLonLatSchema(ctx)
		if err != nil {
			errn := errors.New("Error creating base sparse schema for beam array")
			return errors.Join(err, errn)
		}
		md = map[string]string{"X": "float64", "Y": "float64"}
	}
	defer schema.Free()

	err = beamAttachAttrs(schema, ctx, beam_subrecords, contains_intensity)
	if err != nil {
		errn := errors.New("Error attaching beam data attributes")
		return errors.Join(err, errn)
	}

	err = schema.Check()
	if err != nil {
		errn := errors.New("Error checking beam array TileDB schema")
		return errors.Join(err, errn)
	}

	array, err := tiledb.NewArray(ctx, array_uri)
	if err != nil {
		errn := errors.New("Error creating TileDB beam array")
		return errors.Join(err, errn)
	}
	defer array.Free()

	err = array.Create(schema)
	if err != nil {
		errn := errors.New("Error creating TileDB beam array")
		return errors.Join(err, errn)
	}

	// attach some metadata to preserve python pandas functionality
	// md := map[string]string{"X": "float64", "Y": "float64"}
	key := "__pandas_index_dims"
	err = WriteArrayMetadata(ctx, array_uri, key, md)
	if err != nil {
		return err
	}

	return nil
}

func (fi *FileInfo) pingTdbArrays(ph_ctx, s_md_ctx, si_md_ctx, bd_ctx *tiledb.Context, ph_uri, s_md_uri, si_md_uri, bd_uri string, dense_bd bool) (err error) {
	beam_subrecords := fi.SubRecord_Schema
	contains_intensity := lo.Contains(beam_subrecords, SubRecordNames[INTENSITY_SERIES])
	rec_name := RecordNames[SWATH_BATHYMETRY_PING]
	npings := fi.Record_Counts[rec_name]
	sensor_id := SubRecordID(fi.Metadata.Sensor_Info.Sensor_ID)
	max_beams := fi.Metadata.Quality_Info.Min_Max_Beams[1]

	err = phTdbArray(ph_ctx, ph_uri, npings)
	if err != nil {
		err_ph := errors.New("Error creating PingHeaders TileDB array")
		return errors.Join(err, err_ph)
	}

	err = senTdbArray(s_md_ctx, s_md_uri, npings, sensor_id)
	if err != nil {
		err_s := errors.New("Error creating SensorMetadata TileDB array")
		return errors.Join(err, err_s)
	}

	if contains_intensity {
		err = senImgTdbArray(si_md_ctx, si_md_uri, npings, sensor_id)
		if err != nil {
			err_si := errors.New("Error creating SensorImageryMetadata TileDB array")
			return errors.Join(err, err_si)
		}
	}

	err = beamTdbArray(bd_ctx, bd_uri, beam_subrecords, contains_intensity, dense_bd, npings, max_beams)
	if err != nil {
		err_ba := errors.New("Error creating TileDB beam array")
		return errors.Join(err, err_ba)
	}

	return nil
}

func basePingBeamSchema(ctx *tiledb.Context, npings uint64, max_beams uint16) (schema *tiledb.ArraySchema, err error) {
	// array domain
	domain, err := tiledb.NewDomain(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer domain.Free()

	// want to create blocks of pings. no sense in tiling on the beam axis.
	// better to keep beams as discrete pieces for each ping
	ping_tile_sz := uint64(math.Min(float64(1000), float64(npings)))
	beam_tile_sz := uint64(max_beams)

	// setup dimension options
	// using a combination of delta filter (ascending rows) and zstandard
	pdim, err := tiledb.NewDimension(ctx, "PingNumber", tiledb.TILEDB_UINT64, []uint64{0, npings - uint64(1)}, ping_tile_sz)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer pdim.Free()

	bdim, err := tiledb.NewDimension(ctx, "BeamNumber", tiledb.TILEDB_UINT64, []uint64{0, uint64(max_beams) - uint64(1)}, beam_tile_sz)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer bdim.Free()

	dim_filters, err := tiledb.NewFilterList(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_filters.Free()

	// TODO; might be worth setting a window size
	dim_f1, err := tiledb.NewFilter(ctx, tiledb.TILEDB_FILTER_POSITIVE_DELTA)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_f1.Free()

	level := int32(16)
	dim_f2, err := ZstdFilter(ctx, level)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer dim_f2.Free()

	// attach filters to the pipeline
	err = AddFilters(dim_filters, dim_f1, dim_f2)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	err = pdim.SetFilterList(dim_filters)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = bdim.SetFilterList(dim_filters)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = domain.AddDimensions(pdim, bdim)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	schema, err = tiledb.NewArraySchema(ctx, tiledb.TILEDB_DENSE)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = schema.SetDomain(domain)
	if err != nil {
		return nil, errors.Join(ErrCreateAttitudeTdb, err)
	}

	// cell and tile ordering was an arbitrary choice
	err = schema.SetCellOrder(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = schema.SetTileOrder(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	return schema, nil
}
