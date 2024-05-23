package gsf

import (
	"errors"
	"reflect"
	"strconv"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
)

var ErrSensor = errors.New("Sensor not supported")
var ErrWriteSensorMd = errors.New("Error writing sensor metadata")

type SensorMetadata struct {
	Seabeam        Seabeam
	Em12           Em12
	Em100          Em100
	Em950          Em950
	Em121A         Em121A
	Em121          Em121
	Sass           Sass
	Seamap         Seamap
	Seabat         Seabat
	Em1000         Em1000
	TypeIIISeabeam TypeIIISeabeam
	SbAmp          SbAmp
	SeabatII       SeabatII
	Seabat8101     Seabat8101
	Seabeam2112    Seabeam2112
	ElacMkII       ElacMkII
	CmpSass        CmpSass
	Reson8100      Reson8100
	Em3            Em3
	Em4            Em4
	GeoSwathPlus   GeoSwathPlus
	Klein5410Bss   Klein5410Bss
	Reson7100      Reson7100
	Em3Raw         Em3Raw
	DeltaT         DeltaT
	R2Sonic        R2Sonic
	ResonTSeries   ResonTSeries
	Kmall          Kmall
}

// writeSensorMetadata handles the serialisation of specific sensor related
// metadata to the already setup TileDB array.
// Pushes the buffers to TileDB, doesn't setup the schema or establish the array.
func (sm *SensorMetadata) writeSensorMetadata(ctx *tiledb.Context, array *tiledb.Array, sensor_id SubRecordID, ping_start, ping_end uint64) error {
	// query construction
	query, err := tiledb.NewQuery(ctx, array)
	if err != nil {
		return err
	}
	defer query.Free()

	err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		errn := errors.New("Error setting tile layout for SensorMetadata")
		return errors.Join(err, errn)
	}

	// define the subarray (dim coordinates that we'll write into)
	subarr, err := array.NewSubarray()
	if err != nil {
		errn := errors.New("Error defining subarray for writing SensorMetadata")
		return errors.Join(err, errn)
	}
	defer subarr.Free()

	rng := tiledb.MakeRange(ping_start, ping_end)
	subarr.AddRangeByName("PING_ID", rng)
	err = query.SetSubarray(subarr)
	if err != nil {
		errn := errors.New("Error setting subarray query for wrting SensorMetadata")
		return errors.Join(err, errn)
	}

	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		// EM4
		err := setStructFieldBuffers(query, &sm.Em_4)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata")
			return errors.Join(err, errn)
		}
	default:
		return errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id))))
	}

	// write the data and flush
	err = query.Submit()
	if err != nil {
		errn := errors.New("Error submitting TileDB query")
		return errors.Join(err, errn)
	}

	err = query.Finalize()
	if err != nil {
		errn := errors.New("Error finalising TileDB query")
		return errors.Join(err, errn)
	}

	return nil
}

func (sm *SensorMetadata) attachAttrs(schema *tiledb.ArraySchema, ctx *tiledb.Context, sensor_id SubRecordID) (err error) {
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
		err = schemaAttrs(&sm.Em_4, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.EM_4 attributes")
			return errors.Join(err, err_md)
		}
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

	return nil
}

func (sm *SensorMetadata) appendSensorMetadata(sp *SensorMetadata, sensor_id SubRecordID) error {
	// sp refers to a single pings worth of SensorMetadata
	// whereas sm should be pointing back to the chunks of pings
	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		// EM4
		rf_pd := reflect.ValueOf(&sm.Em_4).Elem()
		rf_sp := reflect.ValueOf(&sp.Em_4).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	default:
		return errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id))))
	}

	return nil
}

// newSensorMetadata is a helper func for initialising SensorMetadata where
// the specific sensor will contain slices initialised to the number of pings
// required.
// This func is only utilised when processing groups of pings to form a single
// cohesive block of data.
func newSensorMetadata(number_pings int, sensor_id SubRecordID) (sen_md SensorMetadata) {
	sen_md = SensorMetadata{}

	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		em4 := Em4{}
		chunkedStructSlices(&em4, number_pings)
		sen_md.Em_4 = em4
	default:
		// TODO; update return sig to allow return of an err rather than simply panic
		panic(errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id)))))
	}

	return sen_md
}

type SensorImageryMetadata struct {
	Em3_imagery          Em3Imagery
	Em4_imagery          Em4Imagery
	Reson7100_imagery    Reson7100Imagery
	ResonTSeries_imagery ResonTSeriesImagery
	Reson8100_imagery    Reson8100Imagery
	Kmall_imagery        KmallImagery
	Klein5410Bss_imagery Klein5410BssImagery
	R2Sonic_imagery      R2SonicImagery
}

// writeSensorImageryMetadata handles the serialisation of specific sensor related
// imagery metadata to the already setup TileDB array.
// Pushes the buffers to TileDB, doesn't setup the schema or establish the array.
func (sim *SensorImageryMetadata) writeSensorImageryMetadata(ctx *tiledb.Context, array *tiledb.Array, sensor_id SubRecordID, ping_start, ping_end uint64) error {
	// query construction
	query, err := tiledb.NewQuery(ctx, array)
	if err != nil {
		return err
	}
	defer query.Free()

	err = query.SetLayout(tiledb.TILEDB_ROW_MAJOR)
	if err != nil {
		errn := errors.New("Error setting tile layout for SensorImageryMetadata")
		return errors.Join(err, errn)
	}

	// define the subarray (dim coordinates that we'll write into)
	subarr, err := array.NewSubarray()
	if err != nil {
		errn := errors.New("Error defining subarray for writing SensorImageryMetadata")
		return errors.Join(err, errn)
	}
	defer subarr.Free()

	rng := tiledb.MakeRange(ping_start, ping_end)
	subarr.AddRangeByName("PING_ID", rng)
	err = query.SetSubarray(subarr)
	if err != nil {
		errn := errors.New("Error setting subarray query for wrting SensorImageryMetadata")
		return errors.Join(err, errn)
	}

	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		err := setStructFieldBuffers(query, &sim.Em4_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata")
			return errors.Join(err, errn)
		}
	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		err := setStructFieldBuffers(query, &sim.Em3_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata")
			return errors.Join(err, errn)
		}
	default:
		return errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id))))
	}

	// write the data and flush
	err = query.Submit()
	if err != nil {
		errn := errors.New("Error submitting TileDB query")
		return errors.Join(err, errn)
	}

	err = query.Finalize()
	if err != nil {
		errn := errors.New("Error finalising TileDB query")
		return errors.Join(err, errn)
	}

	return nil
}

func (sim *SensorImageryMetadata) attachAttrs(schema *tiledb.ArraySchema, ctx *tiledb.Context, sensor_id SubRecordID) (err error) {
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
		err = schemaAttrs(&Em4Imagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.EM4_imagery attributes")
			return errors.Join(err, err_md)
		}
	case KLEIN_5410_BSS:
		// DecodeKlein5410BssImagery
	case KMALL:
		// DecodeKMALLImagery
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		// DecodeR2SonicImagery
	}

	return nil
}

func (sim *SensorImageryMetadata) appendSensorImageryMetadata(sp *SensorImageryMetadata, sensor_id SubRecordID) error {
	// sp refers to a single pings worth of SensorImageryMetadata
	// whereas sim should be pointing back to the chunks of pings
	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		// EM4
		rf_pd := reflect.ValueOf(&sim.Em4_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Em4_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		// EM3
		rf_pd := reflect.ValueOf(&sim.Em3_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Em3_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	default:
		return errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id))))
	}

	return nil
}

// newSensorImageryMetadata is a helper func for initialising SensorImageryMetadata where
// the specific sensor will contain slices initialised to the number of pings
// required.
// This func is only utilised when processing groups of pings to form a single
// cohesive block of data.
func newSensorImageryMetadata(number_pings int, sensor_id SubRecordID) (sen_img_md SensorImageryMetadata) {
	sen_img_md = SensorImageryMetadata{}

	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		// EM4
		em4i := Em4Imagery{}
		chunkedStructSlices(&em4i, number_pings)
		sen_img_md.Em4_imagery = em4i
	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		// EM3
		em3i := Em3Imagery{}
		chunkedStructSlices(&em3i, number_pings)
		sen_img_md.Em3_imagery = em3i
	default:
		// TODO; update return sig to allow return of an err rather than simply panic
		panic(errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id)))))
	}

	return sen_img_md
}
