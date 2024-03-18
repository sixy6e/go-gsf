package gsf

// TODO; confirm the fields with the spec
// fields were simply copied from EM4Imagery
type EM3Imagery struct {
	SamplingFrequency   []float64
	MeanAbsorption      []float32
	TransmitPulseLength []float32
	RangeNorm           []uint16
	StartTvgRamp        []uint16
	StopTvgRamp         []uint16
	BackscatterN        []float32
	BackscatterO        []float32
	TransmitBeamWidth   []float32
	TvgCrossOver        []float32
	offset              []int16 // TODO: replace with ScaleOffset, but as a separate type not embedded
	scale               []int16 // TODO: replace with ScaleOffset, but as a separate type not embedded
}
