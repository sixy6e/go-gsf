package gsf

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
		chunkedStuctSlices(&em4, number_pings)
		sen_md.EM_4 = em4
	}

	return sen_md
}

type SensorImageryMetadata struct {
	EM3_imagery EM3Imagery
	EM4_imagery EM4Imagery
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
		chunkedStuctSlices(&em4i, number_pings)
		sen_img_md.EM4_imagery = em4i
	case EM120, EM300, EM1002, EM2000, EM3000, EM3002, EM3000D, EM3002D, EM121A_SIS:
		// EM3
		em3i := EM3Imagery{}
		chunkedStuctSlices(&em3i, number_pings)
		sen_img_md.EM3_imagery = em3i
	}

	return sen_img_md
}
