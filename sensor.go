package gsf

import (
	"errors"
	"strconv"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
)

var ErrSensor = errors.New("Sensor not supported")
var ErrWriteSensorMd = errors.New("Error writing sensor metadata")

// type ImageryMetadata interface {
// 	Serialise() bool
// }
//
// type SensorMetadata interface {
// 	Append() error
// }

type SensorMetadata struct {
	EM_4 EM4
}

// writeSensorMetadata handles the serialisation of specific sensor related
// metadata to the already setup TileDB array.
// Pushes the buffers to TileDB, doesn't setup the schema or establish the array.
func (sm *SensorMetadata) writeSensorMetadata(query *tiledb.Query, sensor_id SubRecordID) error {
	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		err := setStructFieldBuffers(query, &sm.EM_4)
		if err != nil {
			return errors.Join(err, ErrWriteSensorMd)
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
		em4 := EM4{}
		chunkedStructSlices(&em4, number_pings)
		sen_md.EM_4 = em4
	default:
		// TODO; update return sig to allow return of an err rather than simply panic
		panic(errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id)))))
	}

	return sen_md
}

type SensorImageryMetadata struct {
	EM3_imagery EM3Imagery
	EM4_imagery EM4Imagery
}

// writeSensorImageryMetadata handles the serialisation of specific sensor related
// imagery metadata to the already setup TileDB array.
// Pushes the buffers to TileDB, doesn't setup the schema or establish the array.
func (sim *SensorImageryMetadata) writeSensorImageryMetadata(query *tiledb.Query, sensor_id SubRecordID) error {
	switch sensor_id {
	case EM710, EM302, EM122, EM2040:
		err := setStructFieldBuffers(query, &sim.EM4_imagery)
		if err != nil {
			return errors.Join(err, ErrWriteSensorMd)
		}
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		err := setStructFieldBuffers(query, &sim.EM3_imagery)
		if err != nil {
			return errors.Join(err, ErrWriteSensorMd)
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
		em4i := EM4Imagery{}
		chunkedStructSlices(&em4i, number_pings)
		sen_img_md.EM4_imagery = em4i
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		// EM3
		em3i := EM3Imagery{}
		chunkedStructSlices(&em3i, number_pings)
		sen_img_md.EM3_imagery = em3i
	default:
		// TODO; update return sig to allow return of an err rather than simply panic
		panic(errors.Join(ErrSensor, errors.New(strconv.Itoa(int(sensor_id)))))
	}

	return sen_img_md
}
