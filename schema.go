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
		field.Set(reflect.MakeSlice(ftype, 0, 6))
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
	return nil
}

func mdSchemaAttrs(sensor_id SubRecordID, contains_intensity bool, schema *tiledb.ArraySchema, ctx *tiledb.Context) (err error) {
	// names = make([]string, 0, 10)

	switch sensor_id {

	case SEABEAM:
		// DecodeSeabeam
	case EM12:
		// DecodeEM12
	case EM100:
		// DecodeEM100
	case EM950:
		// DecodeEM950
	case EM121A:
		// DecodeEM121A
	case EM121:
		// DecodeEM121
	case SASS: // obsolete
		// DecodeSASS
	case SEAMAP:
		// DecodeSeaMap
	case SEABAT:
		// DecodeSeaBat
	case EM1000:
		// DecodeEM1000
	case TYPEIII_SEABEAM: // obsolete
		// DecodeTypeIII
	case SB_AMP:
		// DecodeSBAmp
	case SEABAT_II:
		// DecodeSeaBatII
	case SEABAT_8101:
		// DecodeSeaBat8101
	case SEABEAM_2112:
		// DecodeSeaBeam2112
	case ELAC_MKII:
		// DecodeElacMkII
	case CMP_SAAS: // CMP (compressed), should be used in place of SASS
		// DecodeCmpSass
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		// DecodeReson8100
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		// DecodeEM3
	case EM710, EM302, EM122, EM2040:
		// DecodeEM4
		// names = append(names, fieldNames(EM4{})...)
		err = schemaAttrs(&EM4{}, schema, ctx)
	case GEOSWATH_PLUS:
		// DecodeGeoSwathPlus
	case KLEIN_5410_BSS:
		// DecodeKlein5410Bss
	case RESON_7125:
		// DecodeReson7100
	case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
		// DecodeEM3Raw
	case DELTA_T:
		// DecodeDeltaT
	case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
		// DecodeR2Sonic
	case SR_NOT_DEFINED: // the spec makes no mention of ID 154
		panic("Subrecord ID 154 is not defined.")
	case RESON_TSERIES:
		// DecodeResonTSeries
	case KMALL:
		// DecodeKMALL

		// single beam swath sensor specific subrecords
	case SWATH_ECHOTRAC, SWATH_BATHY2000, SWATH_PDD:
		// DecodeSBEchotrac
	case SWATH_MGD77:
		// DecodeSBMGD77
	case SWATH_BDB:
		// DecodeSBBDB
	case SWATH_NOSHDB:
		// DecodeSBNOSHDB
	case SWATH_NAVISOUND:
		// DecodeSBNavisound
	}

	if err != nil {
		return err
	}

	// sensor imagery metadata
	if contains_intensity {

		switch sensor_id {

		case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
			// DecodeEM3Imagery
		case RESON_7125:
			// DecodeReson7100Imagery
		case RESON_TSERIES:
			// DecodeResonTSeriesImagery
		case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
			// DecodeReson8100Imagery
		case EM122, EM302, EM710, EM2040:
			// DecodeEM4Imagery
			// img_md.EM4_imagery = DecodeEM4Imagery(reader)
			// nbytes += n_bytes
			// names = append(names, fieldNames(EM4Imagery{})...)
			err = schemaAttrs(&EM4Imagery{}, schema, ctx)
		case KLEIN_5410_BSS:
			// DecodeKlein5410BssImagery
		case KMALL:
			// DecodeKMALLImagery
		case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
			// DecodeR2SonicImagery
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func pingDenseSchema(ctx *tiledb.Context, sensor_id SubRecordID, npings uint64, contains_intensity bool) (*tiledb.ArraySchema, error) {
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
	// defer schema.Free()

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

	// add the struct fields as tiledb attributes
	// ping header, sensor_metadata, sensor_imagery_metadata
	err = schemaAttrs(&PingHeaders{}, schema, ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = mdSchemaAttrs(sensor_id, contains_intensity, schema, ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	// TODO; look into returning schema and create array elsewhere
	// finally, create the empty array on disk, object store, etc
	// array, err := tiledb.NewArray(ctx, file_uri)
	// if err != nil {
	// 	return nil, errors.Join(ErrCreateAttributeTdb, err)
	// }
	// defer array.Free()

	// err = array.Create(schema)
	// if err != nil {
	// 	return nil, errors.Join(ErrCreateAttributeTdb, err)
	// }

	return schema, nil
}

func beamArrayAttrs(contains_intensity bool, beam_subrecords []string, schema *tiledb.ArraySchema, ctx *tiledb.Context) (err error) {
	var (
		field_tdb_defs map[string]stgpsr.Definition
		def            stgpsr.Definition
		status         bool
	)

	ba := BeamArray{}
	beam_names := make([]string, len(beam_subrecords))

	// cleanup subrecord names to match the BeamArray fields names
	for k, v := range beam_subrecords {
		beam_names[k] = pascalCase(v)
	}

	// values := reflect.ValueOf(ba)
	// types := values.Type()
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
			return errors.Join(ErrCreateAttributeTdb, err)
		}
	}

	// processing the basic ping info (ping id, beam id)
	err = schemaAttrs(&PingBeamNumbers{}, schema, ctx)
	if err != nil {
		return errors.Join(ErrCreateAttributeTdb, err)
	}

	return nil
}

// beamSparseSchema sets up a sparse array schema for the beam array data
// and if it exists, the brb intensity data.
// Longitude and Latitude are the dimensional axes, denoted by X & Y.
// The schema is set to allow duplicates, hilbert for cell ordering, row-major
// for tile ordering.
func beamSparseSchema(contains_intensity bool, beam_subrecords []string, ctx *tiledb.Context) (schema *tiledb.ArraySchema, err error) {
	// array domain
	domain, err := tiledb.NewDomain(ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer domain.Free()

	tile_sz := uint64(1000)
	min_f64 := math.MaxFloat64 * -1

	// setup lon/lat (X/Y) dimensions
	xdim, err := tiledb.NewDimension(ctx, "X", tiledb.TILEDB_FLOAT64, []float64{min_f64, math.MaxFloat64}, tile_sz)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}
	defer xdim.Free()

	ydim, err := tiledb.NewDimension(ctx, "Y", tiledb.TILEDB_FLOAT64, []float64{min_f64, math.MaxFloat64}, tile_sz)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
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
	// defer schema.Free()

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

	err = beamArrayAttrs(contains_intensity, beam_subrecords, schema, ctx)
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	err = schema.Check()
	if err != nil {
		return nil, errors.Join(ErrCreateAttributeTdb, err)
	}

	return schema, nil
}

func (fi *FileInfo) pingSchemas(dense_ctx, sparse_ctx *tiledb.Context) (md_dense_schema, beam_sparse_schema *tiledb.ArraySchema, err error) {
	beam_subrecords := fi.SubRecord_Schema
	contains_intensity := lo.Contains(beam_subrecords, SubRecordNames[INTENSITY_SERIES])

	rec_name := RecordNames[SWATH_BATHYMETRY_PING]
	npings := fi.Record_Counts[rec_name]

	// cleanup subrecord names to match the BeamArray fields names
	// for k, v := range beam_names {
	// 	beam_names[k] = pascalCase(v)
	// }

	// if contains_intensity {
	// 	btype := reflect.TypeOf(BrbIntensity{})
	// 	for i := 0; i < btype.NumField(); i++ {
	// 		if btype.Field(i).IsExported() {
	// 			beam_names = append(beam_names, btype.Field(i).Name)
	// 		}
	// 	}
	// }

	sensor_id := SubRecordID(fi.Metadata.Sensor_Info.Sensor_ID)

	// ping dense array
	md_dense_schema, err = pingDenseSchema(dense_ctx, sensor_id, npings, contains_intensity)
	if err != nil {
		return nil, nil, err
	}
	// defer dense_schema.Free()
	// md_names = md_fields(sensor_id, contains_intensity, schema, ctx)

	beam_sparse_schema, err = beamSparseSchema(contains_intensity, beam_subrecords, sparse_ctx)
	if err != nil {
		return nil, nil, err
	}

	return md_dense_schema, beam_sparse_schema, nil
}

func (fi *FileInfo) PingArrays(dense_file_uri, sparse_file_uri string, dense_ctx, sparse_ctx *tiledb.Context) (beam_names, md_names []string, err error) {
	var (
	// config *tiledb.Config
	)

	// get a generic config if no path provided
	// if config_uri == "" {
	// 	config, err = tiledb.NewConfig()
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// } else {
	// 	config, err = tiledb.LoadConfig(config_uri)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// }

	// defer config.Free()

	// // contexts for both the sparse and dense arrays
	// dense_ctx, err := tiledb.NewContext(config)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// defer dense_ctx.Free()

	// sparse_ctx, err := tiledb.NewContext(config)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// defer sparse_ctx.Free()

	md_dense_schema, beam_sparse_schema, err := fi.pingSchemas(dense_ctx, sparse_ctx)
	if err != nil {
		return nil, nil, err
	}
	defer md_dense_schema.Free()
	defer beam_sparse_schema.Free()

	// create the empty arrays on disk, object store, etc
	md_dense_array, err := tiledb.NewArray(dense_ctx, dense_file_uri)
	if err != nil {
		return nil, nil, errors.Join(ErrCreateMdDenseTdb, err)
	}
	defer md_dense_array.Free()

	err = md_dense_array.Create(md_dense_schema)
	if err != nil {
		return nil, nil, errors.Join(ErrCreateMdDenseTdb, err)
	}

	beam_sparse_array, err := tiledb.NewArray(sparse_ctx, sparse_file_uri)
	if err != nil {
		return nil, nil, errors.Join(ErrCreateBeamSparseTdb, err)
	}
	defer beam_sparse_array.Free()

	err = beam_sparse_array.Create(beam_sparse_schema)
	if err != nil {
		return nil, nil, errors.Join(ErrCreateBeamSparseTdb, err)
	}

	// field names for each array schema
	// sensor metadata
	attrs, err := md_dense_schema.Attributes()
	md_names = make([]string, len(attrs))
	for k, v := range attrs {
		name, err := v.Name()
		if err != nil {
			return nil, nil, err
		}
		md_names[k] = name
	}

	// beam sparse
	attrs, err = beam_sparse_schema.Attributes()
	beam_names = make([]string, len(attrs))
	for k, v := range attrs {
		name, err := v.Name()
		if err != nil {
			return nil, nil, err
		}
		beam_names[k] = name
	}

	return beam_names, md_names, nil
}
