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
	Seabeam          Seabeam
	Em12             Em12
	Em100            Em100
	Em950            Em950
	Em121A           Em121A
	Em121            Em121
	Sass             Sass
	SeaMap           SeaMap
	SeaBat           SeaBat
	Em1000           Em1000
	TypeIIISeabeam   TypeIIISeabeam
	SbAmp            SbAmp
	SeaBatII         SeaBatII
	SeaBat8101       SeaBat8101
	Seabeam2112      Seabeam2112
	ElacMkII         ElacMkII
	CmpSass          CmpSass
	Reson8100        Reson8100
	Em3              Em3
	Em4              Em4
	GeoSwathPlus     GeoSwathPlus
	Klein5410Bss     Klein5410Bss
	Reson7100        Reson7100
	Em3Raw           Em3Raw
	DeltaT           DeltaT
	R2Sonic          R2Sonic
	ResonTSeries     ResonTSeries
	Kmall            Kmall
	SwathSbEchotrac  SwathSbEchotrac
	SwathSbMgd77     SwathSbMgd77
	SwathSbBdb       SwathSbBdb
	SwathSbNoShDb    SwathSbNoShDb
	SwathSbNavisound SwathSbNavisound
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
		errn := errors.New("Error setting subarray query for writing SensorMetadata")
		return errors.Join(err, errn)
	}

	switch sensor_id {
	case SEABEAM:
		err := setStructFieldBuffers(query, &sm.Seabeam)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Seabeam metadata")
			return errors.Join(err, errn)
		}
	case EM12:
		err := setStructFieldBuffers(query, &sm.Em12)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em12 metadata")
			return errors.Join(err, errn)
		}
	case EM100:
		err := setStructFieldBuffers(query, &sm.Em100)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em100 metadata")
			return errors.Join(err, errn)
		}
	case EM950:
		err := setStructFieldBuffers(query, &sm.Em950)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em950 metadata")
			return errors.Join(err, errn)
		}
	case EM121A:
		err := setStructFieldBuffers(query, &sm.Em121A)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em121A metadata")
			return errors.Join(err, errn)
		}
	case EM121:
		err := setStructFieldBuffers(query, &sm.Em121)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em121 metadata")
			return errors.Join(err, errn)
		}
	case SASS: // obsolete
		err := setStructFieldBuffers(query, &sm.Sass)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Sass metadata")
			return errors.Join(err, errn)
		}
	case SEAMAP:
		err := setStructFieldBuffers(query, &sm.SeaMap)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SeaMap metadata")
			return errors.Join(err, errn)
		}
	case SEABAT:
		err := setStructFieldBuffers(query, &sm.SeaBat)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SeaBat metadata")
			return errors.Join(err, errn)
		}
	case EM1000:
		err := setStructFieldBuffers(query, &sm.Em1000)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em1000 metadata")
			return errors.Join(err, errn)
		}
	case TYPEIII_SEABEAM: // obsolete
		err := setStructFieldBuffers(query, &sm.TypeIIISeabeam)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.TypeIIISeabeam metadata")
			return errors.Join(err, errn)
		}
	case SB_AMP:
		err := setStructFieldBuffers(query, &sm.SbAmp)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SbAmp metadata")
			return errors.Join(err, errn)
		}
	case SEABAT_II:
		err := setStructFieldBuffers(query, &sm.SeaBatII)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SeaBatII metadata")
			return errors.Join(err, errn)
		}
	case SEABAT_8101:
		err := setStructFieldBuffers(query, &sm.SeaBat8101)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SeaBat8101 metadata")
			return errors.Join(err, errn)
		}
	case SEABEAM_2112:
		err := setStructFieldBuffers(query, &sm.Seabeam2112)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Seabeam2112 metadata")
			return errors.Join(err, errn)
		}
	case ELAC_MKII:
		err := setStructFieldBuffers(query, &sm.ElacMkII)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.ElacMkII metadata")
			return errors.Join(err, errn)
		}
	case CMP_SAAS: // CMP (compressed), should be used in place of SASS
		err := setStructFieldBuffers(query, &sm.CmpSass)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.CmpSass metadata")
			return errors.Join(err, errn)
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		err := setStructFieldBuffers(query, &sm.Reson8100)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Reson8100 metadata")
			return errors.Join(err, errn)
		}
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		err := setStructFieldBuffers(query, &sm.Em3)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em3 metadata")
			return errors.Join(err, errn)
		}
	case EM710, EM302, EM122, EM2040:
		err := setStructFieldBuffers(query, &sm.Em4)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em4 metadata")
			return errors.Join(err, errn)
		}
	case GEOSWATH_PLUS:
		err := setStructFieldBuffers(query, &sm.GeoSwathPlus)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.GeoSwathPlus metadata")
			return errors.Join(err, errn)
		}
	case KLEIN_5410_BSS:
		err := setStructFieldBuffers(query, &sm.Klein5410Bss)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Klein5410Bss metadata")
			return errors.Join(err, errn)
		}
	case RESON_7125:
		err := setStructFieldBuffers(query, &sm.Reson7100)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Reson7100 metadata")
			return errors.Join(err, errn)
		}
	case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
		err := setStructFieldBuffers(query, &sm.Em3Raw)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Em3Raw metadata")
			return errors.Join(err, errn)
		}
	case DELTA_T:
		err := setStructFieldBuffers(query, &sm.DeltaT)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.DeltaT metadata")
			return errors.Join(err, errn)
		}
	case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
		err := setStructFieldBuffers(query, &sm.R2Sonic)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.R2Sonic metadata")
			return errors.Join(err, errn)
		}
	case RESON_TSERIES:
		err := setStructFieldBuffers(query, &sm.ResonTSeries)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.ResonTSeries metadata")
			return errors.Join(err, errn)
		}
	case KMALL:
		err := setStructFieldBuffers(query, &sm.Kmall)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.Kmall metadata")
			return errors.Join(err, errn)
		}
	case SWATH_SB_ECHOTRAC, SWATH_SB_BATHY2000, SWATH_SB_PDD:
		// they use the same struct, so pushing all to the one sensor
		err := setStructFieldBuffers(query, &sm.SwathSbEchotrac)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SwathSbEchotrac metadata")
			return errors.Join(err, errn)
		}
	case SWATH_SB_MGD77:
		err := setStructFieldBuffers(query, &sm.SwathSbMgd77)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SwathSbMgd77 metadata")
			return errors.Join(err, errn)
		}
	case SWATH_SB_BDB:
		err := setStructFieldBuffers(query, &sm.SwathSbBdb)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SwathSbBdb metadata")
			return errors.Join(err, errn)
		}
	case SWATH_SB_NOSHDB:
		err := setStructFieldBuffers(query, &sm.SwathSbNoShDb)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SwathSbNoShDb metadata")
			return errors.Join(err, errn)
		}
	case SWATH_SB_NAVISOUND:
		err := setStructFieldBuffers(query, &sm.SwathSbNavisound)
		if err != nil {
			errn := errors.New("Error writing SensorMetadata.SwathSbNavisound metadata")
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
		err = schemaAttrs(&sm.Seabeam, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Seabeam attributes")
			return errors.Join(err, err_md)
		}
	case EM12:
		err = schemaAttrs(&sm.Em12, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em12 attributes")
			return errors.Join(err, err_md)
		}
	case EM100:
		err = schemaAttrs(&sm.Em100, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em100 attributes")
			return errors.Join(err, err_md)
		}
	case EM950:
		err = schemaAttrs(&sm.Em950, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em950 attributes")
			return errors.Join(err, err_md)
		}
	case EM121A:
		err = schemaAttrs(&sm.Em121A, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em121A attributes")
			return errors.Join(err, err_md)
		}
	case EM121:
		err = schemaAttrs(&sm.Em121, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em121 attributes")
			return errors.Join(err, err_md)
		}
	case SASS: // obsolete
		err = schemaAttrs(&sm.Sass, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Sass attributes")
			return errors.Join(err, err_md)
		}
	case SEAMAP:
		err = schemaAttrs(&sm.SeaMap, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SeaMap attributes")
			return errors.Join(err, err_md)
		}
	case SEABAT:
		err = schemaAttrs(&sm.SeaBat, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SeaBat attributes")
			return errors.Join(err, err_md)
		}
	case EM1000:
		err = schemaAttrs(&sm.Em1000, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em1000 attributes")
			return errors.Join(err, err_md)
		}
	case TYPEIII_SEABEAM: // obsolete
		err = schemaAttrs(&sm.TypeIIISeabeam, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.TypeIIISeabeam attributes")
			return errors.Join(err, err_md)
		}
	case SB_AMP:
		err = schemaAttrs(&sm.SbAmp, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SbAmp attributes")
			return errors.Join(err, err_md)
		}
	case SEABAT_II:
		err = schemaAttrs(&sm.SeaBatII, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SeaBatII attributes")
			return errors.Join(err, err_md)
		}
	case SEABAT_8101:
		err = schemaAttrs(&sm.SeaBat8101, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SeaBat8101 attributes")
			return errors.Join(err, err_md)
		}
	case SEABEAM_2112:
		err = schemaAttrs(&sm.Seabeam2112, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Seabeam2112 attributes")
			return errors.Join(err, err_md)
		}
	case ELAC_MKII:
		err = schemaAttrs(&sm.ElacMkII, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.ElacMkII attributes")
			return errors.Join(err, err_md)
		}
	case CMP_SAAS: // CMP (compressed), should be used in place of SASS
		err = schemaAttrs(&sm.CmpSass, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.CmpSass attributes")
			return errors.Join(err, err_md)
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		err = schemaAttrs(&sm.Reson8100, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Reson8100 attributes")
			return errors.Join(err, err_md)
		}
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		err = schemaAttrs(&sm.Em3, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em3 attributes")
			return errors.Join(err, err_md)
		}
	case EM710, EM302, EM122, EM2040:
		err = schemaAttrs(&sm.Em4, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em4 attributes")
			return errors.Join(err, err_md)
		}
	case GEOSWATH_PLUS:
		err = schemaAttrs(&sm.GeoSwathPlus, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.GeoSwathPlus attributes")
			return errors.Join(err, err_md)
		}
	case KLEIN_5410_BSS:
		err = schemaAttrs(&sm.Klein5410Bss, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Klein5410Bss attributes")
			return errors.Join(err, err_md)
		}
	case RESON_7125:
		err = schemaAttrs(&sm.Reson7100, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Reson7100 attributes")
			return errors.Join(err, err_md)
		}
	case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
		err = schemaAttrs(&sm.Em3Raw, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Em3Raw attributes")
			return errors.Join(err, err_md)
		}
	case DELTA_T:
		err = schemaAttrs(&sm.DeltaT, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.DeltaT attributes")
			return errors.Join(err, err_md)
		}
	case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
		err = schemaAttrs(&sm.R2Sonic, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.R2Sonic attributes")
			return errors.Join(err, err_md)
		}
	case SR_NOT_DEFINED: // the spec makes no mention of ID 154
		panic("Subrecord ID 154 is not defined.")
	case RESON_TSERIES:
		err = schemaAttrs(&sm.ResonTSeries, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.ResonTSeries attributes")
			return errors.Join(err, err_md)
		}
	case KMALL:
		err = schemaAttrs(&sm.Kmall, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.Kmall attributes")
			return errors.Join(err, err_md)
		}

		// single beam swath sensor specific subrecords
	case SWATH_SB_ECHOTRAC, SWATH_SB_BATHY2000, SWATH_SB_PDD:
		err = schemaAttrs(&sm.SwathSbEchotrac, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SwathSbEchotrac attributes")
			return errors.Join(err, err_md)
		}
	case SWATH_SB_MGD77:
		err = schemaAttrs(&sm.SwathSbMgd77, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SwathSbMgd77 attributes")
			return errors.Join(err, err_md)
		}
	case SWATH_SB_BDB:
		err = schemaAttrs(&sm.SwathSbBdb, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SwathSbBdb attributes")
			return errors.Join(err, err_md)
		}
	case SWATH_SB_NOSHDB:
		err = schemaAttrs(&sm.SwathSbNoShDb, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SwathSbNoShDb attributes")
			return errors.Join(err, err_md)
		}
	case SWATH_SB_NAVISOUND:
		err = schemaAttrs(&sm.SwathSbNavisound, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorMetadata.SwathSbNavisound attributes")
			return errors.Join(err, err_md)
		}
	}

	return nil
}

func (sm *SensorMetadata) appendSensorMetadata(sp *SensorMetadata, sensor_id SubRecordID) error {
	// sp refers to a single pings worth of SensorMetadata
	// whereas sm should be pointing back to the chunks of pings
	switch sensor_id {
	case SEABEAM:
		rf_pd := reflect.ValueOf(&sm.Seabeam).Elem()
		rf_sp := reflect.ValueOf(&sp.Seabeam).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM12:
		rf_pd := reflect.ValueOf(&sm.Em12).Elem()
		rf_sp := reflect.ValueOf(&sp.Em12).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM100:
		rf_pd := reflect.ValueOf(&sm.Em100).Elem()
		rf_sp := reflect.ValueOf(&sp.Em100).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM950:
		rf_pd := reflect.ValueOf(&sm.Em950).Elem()
		rf_sp := reflect.ValueOf(&sp.Em950).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM121A:
		rf_pd := reflect.ValueOf(&sm.Em121A).Elem()
		rf_sp := reflect.ValueOf(&sp.Em121A).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM121:
		rf_pd := reflect.ValueOf(&sm.Em121).Elem()
		rf_sp := reflect.ValueOf(&sp.Em121).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SASS: // obsolete
		rf_pd := reflect.ValueOf(&sm.Sass).Elem()
		rf_sp := reflect.ValueOf(&sp.Sass).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SEAMAP:
		rf_pd := reflect.ValueOf(&sm.SeaMap).Elem()
		rf_sp := reflect.ValueOf(&sp.SeaMap).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SEABAT:
		rf_pd := reflect.ValueOf(&sm.SeaBat).Elem()
		rf_sp := reflect.ValueOf(&sp.SeaBat).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM1000:
		rf_pd := reflect.ValueOf(&sm.Em1000).Elem()
		rf_sp := reflect.ValueOf(&sp.Em1000).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case TYPEIII_SEABEAM: // obsolete
		rf_pd := reflect.ValueOf(&sm.TypeIIISeabeam).Elem()
		rf_sp := reflect.ValueOf(&sp.TypeIIISeabeam).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SB_AMP:
		rf_pd := reflect.ValueOf(&sm.SbAmp).Elem()
		rf_sp := reflect.ValueOf(&sp.SbAmp).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SEABAT_II:
		rf_pd := reflect.ValueOf(&sm.SeaBatII).Elem()
		rf_sp := reflect.ValueOf(&sp.SeaBatII).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SEABAT_8101:
		rf_pd := reflect.ValueOf(&sm.SeaBat8101).Elem()
		rf_sp := reflect.ValueOf(&sp.SeaBat8101).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SEABEAM_2112:
		rf_pd := reflect.ValueOf(&sm.Seabeam2112).Elem()
		rf_sp := reflect.ValueOf(&sp.Seabeam2112).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case ELAC_MKII:
		rf_pd := reflect.ValueOf(&sm.ElacMkII).Elem()
		rf_sp := reflect.ValueOf(&sp.ElacMkII).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case CMP_SAAS: // CMP (compressed), should be used in place of SASS
		rf_pd := reflect.ValueOf(&sm.CmpSass).Elem()
		rf_sp := reflect.ValueOf(&sp.CmpSass).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		rf_pd := reflect.ValueOf(&sm.Reson8100).Elem()
		rf_sp := reflect.ValueOf(&sp.Reson8100).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		rf_pd := reflect.ValueOf(&sm.Em3).Elem()
		rf_sp := reflect.ValueOf(&sp.Em3).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM710, EM302, EM122, EM2040:
		rf_pd := reflect.ValueOf(&sm.Em4).Elem()
		rf_sp := reflect.ValueOf(&sp.Em4).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case GEOSWATH_PLUS:
		rf_pd := reflect.ValueOf(&sm.GeoSwathPlus).Elem()
		rf_sp := reflect.ValueOf(&sp.GeoSwathPlus).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case KLEIN_5410_BSS:
		rf_pd := reflect.ValueOf(&sm.Klein5410Bss).Elem()
		rf_sp := reflect.ValueOf(&sp.Klein5410Bss).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_7125:
		rf_pd := reflect.ValueOf(&sm.Reson7100).Elem()
		rf_sp := reflect.ValueOf(&sp.Reson7100).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
		rf_pd := reflect.ValueOf(&sm.Em3Raw).Elem()
		rf_sp := reflect.ValueOf(&sp.Em3Raw).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case DELTA_T:
		rf_pd := reflect.ValueOf(&sm.DeltaT).Elem()
		rf_sp := reflect.ValueOf(&sp.DeltaT).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
		rf_pd := reflect.ValueOf(&sm.R2Sonic).Elem()
		rf_sp := reflect.ValueOf(&sp.R2Sonic).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_TSERIES:
		rf_pd := reflect.ValueOf(&sm.ResonTSeries).Elem()
		rf_sp := reflect.ValueOf(&sp.ResonTSeries).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case KMALL:
		rf_pd := reflect.ValueOf(&sm.Kmall).Elem()
		rf_sp := reflect.ValueOf(&sp.Kmall).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SWATH_SB_ECHOTRAC, SWATH_SB_BATHY2000, SWATH_SB_PDD:
		rf_pd := reflect.ValueOf(&sm.SwathSbEchotrac).Elem()
		rf_sp := reflect.ValueOf(&sp.SwathSbEchotrac).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SWATH_SB_MGD77:
		rf_pd := reflect.ValueOf(&sm.SwathSbMgd77).Elem()
		rf_sp := reflect.ValueOf(&sp.SwathSbMgd77).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SWATH_SB_BDB:
		rf_pd := reflect.ValueOf(&sm.SwathSbBdb).Elem()
		rf_sp := reflect.ValueOf(&sp.SwathSbBdb).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SWATH_SB_NOSHDB:
		rf_pd := reflect.ValueOf(&sm.SwathSbNoShDb).Elem()
		rf_sp := reflect.ValueOf(&sp.SwathSbNoShDb).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case SWATH_SB_NAVISOUND:
		rf_pd := reflect.ValueOf(&sm.SwathSbNavisound).Elem()
		rf_sp := reflect.ValueOf(&sp.SwathSbNavisound).Elem()
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
	case SEABEAM:
		md := Seabeam{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Seabeam = md
	case EM12:
		md := Em12{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em12 = md
	case EM100:
		md := Em100{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em100 = md
	case EM950:
		md := Em950{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em950 = md
	case EM121A:
		md := Em121A{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em121A = md
	case EM121:
		md := Em121{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em121 = md
	case SASS: // obsolete
		md := Sass{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Sass = md
	case SEAMAP:
		md := SeaMap{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SeaMap = md
	case SEABAT:
		md := SeaBat{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SeaBat = md
	case EM1000:
		md := Em1000{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em1000 = md
	case TYPEIII_SEABEAM: // obsolete
		md := TypeIIISeabeam{}
		chunkedStructSlices(&md, number_pings)
		sen_md.TypeIIISeabeam = md
	case SB_AMP:
		md := SbAmp{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SbAmp = md
	case SEABAT_II:
		md := SeaBatII{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SeaBatII = md
	case SEABAT_8101:
		md := SeaBat8101{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SeaBat8101 = md
	case SEABEAM_2112:
		md := Seabeam2112{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Seabeam2112 = md
	case ELAC_MKII:
		md := ElacMkII{}
		chunkedStructSlices(&md, number_pings)
		sen_md.ElacMkII = md
	case CMP_SAAS: // CMP (compressed), should be used in place of SASS
		md := CmpSass{}
		chunkedStructSlices(&md, number_pings)
		sen_md.CmpSass = md
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		md := Reson8100{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Reson8100 = md
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		md := Em3{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em3 = md
	case EM710, EM302, EM122, EM2040:
		md := Em4{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em4 = md
	case GEOSWATH_PLUS:
		md := GeoSwathPlus{}
		chunkedStructSlices(&md, number_pings)
		sen_md.GeoSwathPlus = md
	case KLEIN_5410_BSS:
		md := Klein5410Bss{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Klein5410Bss = md
	case RESON_7125:
		md := Reson7100{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Reson7100 = md
	case EM300_RAW, EM1002_RAW, EM2000_RAW, EM3000_RAW, EM120_RAW, EM3002_RAW, EM3000D_RAW, EM3002D_RAW, EM121A_SIS_RAW:
		md := Em3Raw{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Em3Raw = md
	case DELTA_T:
		md := DeltaT{}
		chunkedStructSlices(&md, number_pings)
		sen_md.DeltaT = md
	case R2SONIC_2022, R2SONIC_2024, R2SONIC_2020:
		md := R2Sonic{}
		chunkedStructSlices(&md, number_pings)
		sen_md.R2Sonic = md
	case RESON_TSERIES:
		md := ResonTSeries{}
		chunkedStructSlices(&md, number_pings)
		sen_md.ResonTSeries = md
	case KMALL:
		md := Kmall{}
		chunkedStructSlices(&md, number_pings)
		sen_md.Kmall = md
	case SWATH_SB_ECHOTRAC, SWATH_SB_BATHY2000, SWATH_SB_PDD:
		md := SwathSbEchotrac{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SwathSbEchotrac = md
	case SWATH_SB_MGD77:
		md := SwathSbMgd77{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SwathSbMgd77 = md
	case SWATH_SB_BDB:
		md := SwathSbBdb{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SwathSbBdb = md
	case SWATH_SB_NOSHDB:
		md := SwathSbNoShDb{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SwathSbNoShDb = md
	case SWATH_SB_NAVISOUND:
		md := SwathSbNavisound{}
		chunkedStructSlices(&md, number_pings)
		sen_md.SwathSbNavisound = md
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
			errn := errors.New("Error writing SensorImageryMetadata.Em4_imagery metadata")
			return errors.Join(err, errn)
		}
	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		err := setStructFieldBuffers(query, &sim.Em3_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.Em3_imagery metadata")
			return errors.Join(err, errn)
		}
	case RESON_7125:
		err := setStructFieldBuffers(query, &sim.Reson7100_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.Reson7100_imagery metadata")
			return errors.Join(err, errn)
		}
	case RESON_TSERIES:
		err := setStructFieldBuffers(query, &sim.ResonTSeries_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.ResonTSeries_imagery metadata")
			return errors.Join(err, errn)
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		err := setStructFieldBuffers(query, &sim.Reson8100_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.Reson8100_imagery metadata")
			return errors.Join(err, errn)
		}
	case KLEIN_5410_BSS:
		err := setStructFieldBuffers(query, &sim.Klein5410Bss_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.Klein5410Bss_imagery metadata")
			return errors.Join(err, errn)
		}
	case KMALL:
		err := setStructFieldBuffers(query, &sim.Kmall_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.Kmall_imagery metadata")
			return errors.Join(err, errn)
		}
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		err := setStructFieldBuffers(query, &sim.R2Sonic_imagery)
		if err != nil {
			errn := errors.New("Error writing SensorImageryMetadata.R2Sonic_imagery metadata")
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
		err = schemaAttrs(&Em3Imagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Em3_imagery attributes")
			return errors.Join(err, err_md)
		}
	case RESON_7125:
		err = schemaAttrs(&Reson7100Imagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Reson7100_imagery attributes")
			return errors.Join(err, err_md)
		}
	case RESON_TSERIES:
		err = schemaAttrs(&ResonTSeriesImagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.ResonTSeries_imagery attributes")
			return errors.Join(err, err_md)
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		err = schemaAttrs(&Reson8100Imagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Reson8100_imagery attributes")
			return errors.Join(err, err_md)
		}
	case EM122, EM302, EM710, EM2040:
		err = schemaAttrs(&Em4Imagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Em4_imagery attributes")
			return errors.Join(err, err_md)
		}
	case KLEIN_5410_BSS:
		err = schemaAttrs(&Klein5410BssImagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Klein5410Bss_imagery attributes")
			return errors.Join(err, err_md)
		}
	case KMALL:
		err = schemaAttrs(&KmallImagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.Kmall_imagery attributes")
			return errors.Join(err, err_md)
		}
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		err = schemaAttrs(&R2SonicImagery{}, schema, ctx)
		if err != nil {
			err_md := errors.New("Error creating SensorImageryMetadata.R2Sonic_imagery attributes")
			return errors.Join(err, err_md)
		}
	}

	return nil
}

func (sim *SensorImageryMetadata) appendSensorImageryMetadata(sp *SensorImageryMetadata, sensor_id SubRecordID) error {
	// sp refers to a single pings worth of SensorImageryMetadata
	// whereas sim should be pointing back to the chunks of pings
	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
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
		rf_pd := reflect.ValueOf(&sim.Em3_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Em3_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_7125:
		rf_pd := reflect.ValueOf(&sim.Reson7100_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Reson7100_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_TSERIES:
		rf_pd := reflect.ValueOf(&sim.ResonTSeries_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.ResonTSeries_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		rf_pd := reflect.ValueOf(&sim.Reson8100_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Reson8100_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case KLEIN_5410_BSS:
		rf_pd := reflect.ValueOf(&sim.Klein5410Bss_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Klein5410Bss_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case KMALL:
		rf_pd := reflect.ValueOf(&sim.Kmall_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.Kmall_imagery).Elem()
		types := rf_pd.Type()

		for i := 0; i < rf_pd.NumField(); i++ {
			name := types.Field(i).Name
			field_pd := rf_pd.FieldByName(name)
			field_sp := rf_sp.FieldByName(name)
			field_pd.Set(reflect.AppendSlice(field_pd, field_sp))
		}
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		rf_pd := reflect.ValueOf(&sim.R2Sonic_imagery).Elem()
		rf_sp := reflect.ValueOf(&sp.R2Sonic_imagery).Elem()
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
		em4i := Em4Imagery{}
		chunkedStructSlices(&em4i, number_pings)
		sen_img_md.Em4_imagery = em4i
	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		em3i := Em3Imagery{}
		chunkedStructSlices(&em3i, number_pings)
		sen_img_md.Em3_imagery = em3i
	case RESON_7125:
		r7100 := Reson7100Imagery{}
		chunkedStructSlices(&r7100, number_pings)
		sen_img_md.Reson7100_imagery = r7100
	case RESON_TSERIES:
		rtseries := ResonTSeriesImagery{}
		chunkedStructSlices(&rtseries, number_pings)
		sen_img_md.ResonTSeries_imagery = rtseries
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		r8100 := Reson8100Imagery{}
		chunkedStructSlices(&r8100, number_pings)
		sen_img_md.Reson8100_imagery = r8100
	case KLEIN_5410_BSS:
		k5410bss := Klein5410BssImagery{}
		chunkedStructSlices(&k5410bss, number_pings)
		sen_img_md.Klein5410Bss_imagery = k5410bss
	case KMALL:
		kmall := KmallImagery{}
		chunkedStructSlices(&kmall, number_pings)
		sen_img_md.Kmall_imagery = kmall
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		r2sonic := R2SonicImagery{}
		chunkedStructSlices(&r2sonic, number_pings)
		sen_img_md.R2Sonic_imagery = r2sonic
	default:
		// TODO; update return sig to allow return of an err rather than simply panic
		panic(errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id)))))
	}

	return sen_img_md
}
