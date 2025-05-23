package gsf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"time"
)

// Em3Imagery caters for generation 3 EM sensors. Specifically:
// EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW,
// EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D,
// EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW.
type Em3Imagery struct {
	RangeNorm      []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartTvgRamp   []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StopTvgRamp    []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BackscatterN   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	BackscatterO   []uint8   `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
	MeanAbsorption []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
}

// nullEm3Imagery initialises Em3Imagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullEm3Imagery() (imd Em3Imagery) {
	imd = Em3Imagery{
		[]uint16{NULL_UINT16_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[]uint8{NULL_UINT8_ZERO},
		[]uint8{NULL_UINT8_ZERO},
		[]float32{float32(math.NaN())},
	}
	return imd
}

// DecodeEm3Imagery decodes generation 3 EM INTENSITY_SERIES SubRecord and
// constructs the Em3Imagery type.
func DecodeEm3Imagery(reader *bytes.Reader) (img_md Em3Imagery, scl_off ScaleOffset, err error) {
	var (
		base struct {
			RangeNorm      uint16
			StartTvgRamp   uint16
			StopTvgRamp    uint16
			BackscatterN   uint8
			BackscatterO   uint8
			MeanAbsorption uint16
			Scale          int16
			Offset         int16
			Spare          [4]byte
		}
	)

	err = binary.Read(reader, binary.BigEndian, &base)
	if err != nil {
		errn := errors.New("EM3 sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, scl_off, err
	}

	img_md.RangeNorm = []uint16{base.RangeNorm}
	img_md.StartTvgRamp = []uint16{base.StartTvgRamp}
	img_md.StopTvgRamp = []uint16{base.StopTvgRamp}
	img_md.BackscatterN = []uint8{base.BackscatterN}
	img_md.BackscatterO = []uint8{base.BackscatterO}
	img_md.MeanAbsorption = []float32{float32(base.MeanAbsorption) / SCALE_2_F32}

	// The gsf spec mentions that the scale factor is 2 for EM3 based sensors
	// Ideally the stored value should be used, unfortunately, some of the sample
	// files had incorrect scale factors due to a bug in the source software that
	// generated the file.
	// The other issue, is that some source software put different values again
	// in the SCALE_FACTORS SubRecord, and were using those.
	// So if at any point, the intensity timeseries data doesn't look right,
	// potentially the code needs to be adjusted to read the correct factors.

	// scl_off = ScaleOffset{float64(2), float64(base.Offset)}
	scl_off = ScaleOffset{float64(base.Scale), float64(base.Offset)}

	return img_md, scl_off, err
}

// Em4Imagery caters for generation 4 EM sensors. Specifically:
// EM710, EM302, EM122, EM2040.
type Em4Imagery struct {
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
}

// NullEm4Imagery initialises Em4Imagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullEm4Imagery() (imd Em4Imagery) {
	imd = Em4Imagery{
		[]float64{math.NaN()},
		[]float32{float32(math.NaN())},
		[]float32{float32(math.NaN())},
		[]uint16{NULL_UINT16_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[]float32{float32(math.NaN())},
		[]float32{float32(math.NaN())},
		[]float32{float32(math.NaN())},
		[]float32{float32(math.NaN())},
	}
	return imd
}

// DecodeEm4Imagery decodes generation 4 EM INTENSITY_SERIES SubRecords and constructs the Em4Imagery type.
func DecodeEm4Imagery(reader *bytes.Reader) (img_md Em4Imagery, scl_off ScaleOffset, err error) {

	var (
		base struct {
			SamplingFrequency1  uint32
			SamplingFrequency2  uint32
			MeanAbsorption      uint16
			TransmitPulseLength uint16
			RangeNorm           uint16
			StartTvgRamp        uint16
			StopTvgRamp         uint16
			BackscatterN        uint16
			BackscatterO        uint16
			TransmitBeamWidth   uint16
			TvgCrossOver        uint16
			Offset              int16
			Scale               int16
			Spare               [20]byte // 20 bytes spare
		} // 50 bytes
	)

	err = binary.Read(reader, binary.BigEndian, &base)
	if err != nil {
		errn := errors.New("EM4 sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, scl_off, err
	}

	img_md.SamplingFrequency = []float64{
		float64(base.SamplingFrequency1) +
			float64(base.SamplingFrequency2)/
				float64(4_000_000_000)}
	img_md.MeanAbsorption = []float32{float32(base.MeanAbsorption) / SCALE_2_F32}
	img_md.TransmitPulseLength = []float32{float32(base.TransmitPulseLength)}
	img_md.RangeNorm = []uint16{base.RangeNorm}
	img_md.StartTvgRamp = []uint16{base.StartTvgRamp}
	img_md.StopTvgRamp = []uint16{base.StopTvgRamp}
	img_md.BackscatterN = []float32{float32(int16(base.BackscatterN)) / SCALE_1_F32}
	img_md.BackscatterO = []float32{float32(int16(base.BackscatterO)) / SCALE_1_F32}
	img_md.TransmitBeamWidth = []float32{float32(base.TransmitBeamWidth) / SCALE_1_F32}
	img_md.TvgCrossOver = []float32{float32(base.TvgCrossOver) / SCALE_1_F32}

	// The gsf spec mentions that the scale factor is 10 for EM4 based sensors
	// Ideally the stored value should be used, unfortunately, some of the sample
	// files had incorrect scale factors due to a bug in the source software that
	// generated the file.
	// The other issue, is that some source software put different values again
	// in the SCALE_FACTORS SubRecord, and were using those.
	// So if at any point, the intensity timeseries data doesn't look right,
	// potentially the code needs to be adjusted to read the correct factors.

	// scl_off = ScaleOffset{float64(10), float64(base.Offset)}
	scl_off = ScaleOffset{float64(base.Scale), float64(base.Offset)}

	return img_md, scl_off, err
}

// Reson7100Imagery caters for the Reson7100 series sensors, specifically:
// RESON_7125.
type Reson7100Imagery struct {
	Null []uint8 `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

// nullReson7100Imagery initialises Reson7100Imagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullReson7100Imagery() (imd Reson7100Imagery) {
	imd = Reson7100Imagery{
		[]uint8{NULL_UINT8_ZERO},
	}
	return imd
}

// DecodeReson7100Imagery decodes Reson7100 INTENSITY_SERIES SubRecords and constructs
// the Reson7100Imagery type.
func DecodeReson7100Imagery(reader *bytes.Reader) (img_md Reson7100Imagery, err error) {
	var buffer struct {
		Size  uint16
		Spare [64]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Reson7100 sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}

	img_md.Null = []uint8{0}

	return img_md, err
}

// Reson7100Imagery caters for the ResonTSeries sensor.
type ResonTSeriesImagery struct {
	Null []uint8 `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

// nullResonTSeriesImagery initialises ResonTSeriesImagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullResonTSeriesImagery() (imd ResonTSeriesImagery) {
	imd = ResonTSeriesImagery{
		[]uint8{NULL_UINT8_ZERO},
	}
	return imd
}

// DecodeResonTSeriesImagery decodes the ResonTSeries INTENSITY_SERIES SubRecord
// and constructs the ResonTSeriesImagery type.
func DecodeResonTSeriesImagery(reader *bytes.Reader) (img_md ResonTSeriesImagery, err error) {
	var buffer struct {
		Size  uint16
		Spare [64]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("ResonTSeries sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}

	img_md.Null = []uint8{0}

	return img_md, err
}

// Reson8100Imagery caters for Reson8100 series sensors specifically:
// RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160.
type Reson8100Imagery struct {
	Null []uint8 `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

// nullReson8100Imagery initialises Reson8100Imagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullReson8100Imagery() (imd Reson8100Imagery) {
	imd = Reson8100Imagery{
		[]uint8{NULL_UINT8_ZERO},
	}
	return imd
}

// DecodeReson8100Imagery decodes Reson8100 INTENSITY_SERIES SubRecord and constructs
// the Reson8100Imagery type.
func DecodeReson8100Imagery(reader *bytes.Reader) (img_md Reson8100Imagery, err error) {
	var buffer struct {
		Spare [8]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Reson8100 sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}

	img_md.Null = []uint8{0}

	return img_md, err
}

// KmallImagery caters for KMALL sensors.
type KmallImagery struct {
	Null []uint8 `tiledb:"dtype=uint8,ftype=attr" filters:"zstd(level=16)"`
}

// nullKmallImagery initialises KmallImagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullKmallImagery() (imd KmallImagery) {
	imd = KmallImagery{
		[]uint8{NULL_UINT8_ZERO},
	}
	return imd
}

// DecodeKmallImagery decodes KMALL INTENSITY_SERIES SubRecord and constructs the
// KmallImagery type.
func DecodeKmallImagery(reader *bytes.Reader) (img_md KmallImagery, err error) {
	var buffer struct {
		Spare [64]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("KMALL sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}

	img_md.Null = []uint8{0}

	return img_md, err
}

// Klein5410BssImagery caters for the KLEIN_5410_BSS sensor.
type Klein5410BssImagery struct {
	ResolutionMode []uint16   `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TvgPage        []uint16   `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	BeamId         [][]uint16 `tiledb:"dtype=uint16,ftype=attr,var" filters:"zstd(level=16)"`
}

// nullKlein5410BssImagery initialises Klein5410BssImagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullKlein5410BssImagery() (imd Klein5410BssImagery) {
	imd = Klein5410BssImagery{
		[]uint16{NULL_UINT16_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[][]uint16{{NULL_UINT16_ZERO}},
	}
	return imd
}

// DecodeKlein5410BssImagery decodes KLEIN_5410_BSS INTENSITY_SERIES SubRecord
// and constructs the Klein5410BssImagery type.
func DecodeKlein5410BssImagery(reader *bytes.Reader) (img_md Klein5410BssImagery, err error) {
	var buffer struct {
		ResolutionMode uint16
		TvgPage        uint16
		BeamId         [5]uint16
		Spare          [4]byte
	}
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("Klein5410BSS sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}

	img_md.ResolutionMode = []uint16{buffer.ResolutionMode}
	img_md.TvgPage = []uint16{buffer.TvgPage}
	img_md.BeamId = [][]uint16{buffer.BeamId[:]}

	return img_md, err
}

// R2SonicImagery caters for R2Sonic sensors specifically:
// R2SONIC_2020, R2SONIC_2022, R2SONIC_2024.
type R2SonicImagery struct {
	ModelNumber      []string    `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
	SerialNumber     []string    `tiledb:"dtype=string,ftype=attr" filters:"zstd(level=16)"`
	DgTime           []time.Time `tiledb:"dtype=datetime_ns,ftype=attr" filters:"zstd(level=16)"`
	PingNumber       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	PingPeriod       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	SoundSpeed       []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	Frequency        []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPower          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxPulseWidth     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxBeamWidthVert  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxBeamWidthHoriz []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxSteeringVert   []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxSteeringHoriz  []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	TxMiscInfo       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	RxBandwidth      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxSampleRate     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxRange          []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxGain           []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxSpreading      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxAbsorption     []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxMountTilt      []float64   `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	RxMiscInfo       []uint32    `tiledb:"dtype=uint32,ftype=attr" filters:"zstd(level=16)"`
	NumberBeams      []uint16    `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	MoreInfo         [][]float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
}

// nullR2SonicImagery initialises R2SonicImagery with a null/fill value.
// Used for instances where this data isn't acquired for a given ping.
func nullR2SonicImagery() (imd R2SonicImagery) {
	imd = R2SonicImagery{
		[]string{"NULL"},
		[]string{"NULL"},
		[]time.Time{time.Unix(0, 0).UTC()},
		[]uint32{NULL_UINT32_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]uint32{NULL_UINT32_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]float64{NULL_FLOAT64_ZERO},
		[]uint32{NULL_UINT32_ZERO},
		[]uint16{NULL_UINT16_ZERO},
		[][]float64{{NULL_FLOAT64_ZERO}},
	}
	return imd
}

// DecodeR2SonicImagery decodes R2Sonic INTENSITY_SERIES SubRecord and constructs
// the R2SonicImagery type.
func DecodeR2SonicImagery(reader *bytes.Reader) (img_md R2SonicImagery, err error) {
	var (
		buffer struct {
			ModelNumber      [12]byte
			SerialNumber     [12]byte
			TvSec            uint32
			TvNsec           uint32
			PingNumber       uint32
			PingPeriod       uint32
			SoundSpeed       uint32
			Frequency        uint32
			TxPower          uint32
			TxPulseWidth     uint32
			TxBeamWidthVert  uint32
			TxBeamWidthHoriz uint32
			TxSteeringVert   uint32
			TxSteeringHoriz  uint32
			TxMiscInfo       uint32
			RxBandwidth      uint32
			RxSampleRate     uint32
			RxRange          uint32
			RxGain           uint32
			RxSpreading      uint32
			RxAbsorption     uint32
			RxMountTilt      uint32
			RxMiscInfo       uint32
			Spare1           [2]byte
			NumberBeams      uint16
		}
		var_buf struct {
			MoreInfo [6]uint32
			Spare2   [32]byte
		}
	)

	// block one
	err = binary.Read(reader, binary.BigEndian, &buffer)
	if err != nil {
		errn := errors.New("R2Sonic sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}
	img_md.ModelNumber = []string{string(buffer.ModelNumber[:])}
	img_md.SerialNumber = []string{string(buffer.SerialNumber[:])}
	img_md.DgTime = []time.Time{time.Unix(int64(buffer.TvSec), int64(buffer.TvNsec)).UTC()}
	img_md.PingNumber = []uint32{buffer.PingNumber}
	img_md.PingPeriod = []float64{float64(buffer.PingPeriod) / SCALE_6_F64}
	img_md.SoundSpeed = []float64{float64(buffer.SoundSpeed) / SCALE_2_F64}
	img_md.Frequency = []float64{float64(buffer.Frequency) / SCALE_3_F64}
	img_md.TxPower = []float64{float64(buffer.TxPower) / SCALE_2_F64}
	img_md.TxPulseWidth = []float64{float64(buffer.TxPulseWidth) / SCALE_7_F64}
	img_md.TxBeamWidthVert = []float64{float64(buffer.TxBeamWidthVert) / SCALE_6_F64}
	img_md.TxBeamWidthHoriz = []float64{float64(buffer.TxBeamWidthHoriz) / SCALE_6_F64}
	img_md.TxSteeringVert = []float64{float64(int32(buffer.TxSteeringVert)) / SCALE_6_F64}
	img_md.TxSteeringHoriz = []float64{float64(int32(buffer.TxSteeringHoriz)) / SCALE_6_F64}
	img_md.TxMiscInfo = []uint32{buffer.TxMiscInfo}
	img_md.RxBandwidth = []float64{float64(buffer.RxBandwidth) / SCALE_4_F64}
	img_md.RxSampleRate = []float64{float64(buffer.RxSampleRate) / SCALE_3_F64}
	img_md.RxRange = []float64{float64(buffer.RxRange) / SCALE_5_F64}
	img_md.RxGain = []float64{float64(buffer.RxGain) / SCALE_2_F64}
	img_md.RxSpreading = []float64{float64(buffer.RxSpreading) / SCALE_3_F64}
	img_md.RxAbsorption = []float64{float64(buffer.RxAbsorption) / SCALE_3_F64}
	img_md.RxMountTilt = []float64{float64(int32(buffer.RxMountTilt)) / SCALE_6_F64}
	img_md.RxMiscInfo = []uint32{buffer.RxMiscInfo}
	img_md.NumberBeams = []uint16{buffer.NumberBeams}

	// block two (var length array)
	err = binary.Read(reader, binary.BigEndian, &var_buf)
	if err != nil {
		errn := errors.New("R2Sonic sensor")
		err = errors.Join(err, ErrSensorImgMetadata, errn)
		return img_md, err
	}
	minfo := make([]float64, 0, 6)
	for i := 0; i < 6; i++ {
		minfo = append(minfo, float64(int32(var_buf.MoreInfo[i]))/SCALE_6_F64)
	}

	img_md.MoreInfo = [][]float64{minfo}

	return img_md, err
}
