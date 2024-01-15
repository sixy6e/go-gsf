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

type SensorImageryMetadata struct {
	EM3_imagery EM3Imagery
	EM4_imagery EM4Imagery
}
