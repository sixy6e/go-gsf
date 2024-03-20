package gsf

// TODO; confirm the fields with the spec
// fields were simply copied from EM4Imagery
type EM3Imagery struct {
	SamplingFrequency   []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	MeanAbsorption      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitPulseLength []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	RangeNorm           []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartTvgRamp        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StopTvgRamp         []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BackscatterN        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	BackscatterO        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TransmitBeamWidth   []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	TvgCrossOver        []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	offset              []int16   // TODO: replace with ScaleOffset, but as a separate type not embedded
	scale               []int16   // TODO: replace with ScaleOffset, but as a separate type not embedded
}
