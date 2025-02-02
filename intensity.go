package gsf

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/samber/lo"
)

// Removing BrbIntensity.BottomDetect for now. Need to get more info on what the
// BottomDetectIndex means. I thought it would be the index location of the TimeSeries
// slice, but index values have included values that are outside the array.
// i.e TimeSeries samples = 6, BottomDetectIndex = 6
// This would be ok if it is 1-based indexing stored in the GSF file for this particular
// piece of data, except that the BottomDetectIndex also contains values of 0.

// BrbIntensity
type BrbIntensity struct {
	TimeSeries []float64 `tiledb:"dtype=float64,ftype=attr,var" filters:"zstd(level=16)"`
	// BottomDetect      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	BottomDetectIndex []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartRange        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	TsMean            []float64 `tiledb:"dtype=float64,ftype=attr" filters:"zstd(level=16)"`
	sample_count      []uint16
	// timeseries        [][]float32
}

// newBrbIntensity is a helper func for when initialising BrbIntensity and
// attached to the PingData type. This func is only utilised when we're
// processing the ping data in chunks, and combining each chunk into
// a single cohesive unit for output into TileDB.
// The TimeSeries field is of variable length, and we don't know the total
// number of samples in each beam until runtime. So the array is set to the
// capacity of number_beams * 66. No thorough investigation on the choice of 66,
// except that for a few sample GSF files, the number of samples for each beam
// was in the 60s.
func newBrbIntensity(number_beams int) (brb_int BrbIntensity) {
	brb_int = BrbIntensity{
		make([]float64, 0, number_beams*66), // 66 ... just becasuse
		// make([]float32, 0, number_beams),
		make([]uint16, 0, number_beams),
		make([]uint16, 0, number_beams),
		make([]float64, 0, number_beams),
		make([]uint16, 0, number_beams),
	}
	return brb_int
}

// DecodeBrbIntensity decodes the timeseries intensity sub-record. Each beam will have
// a variable length of intensity samples, the index for the bottom detect sample, and the
// sample itself.
func DecodeBrbIntensity(reader *bytes.Reader, nbeams uint16, sensor_id SubRecordID) (intensity BrbIntensity, img_md SensorImageryMetadata, err error) {

	var (
		base struct {
			Bits_per_sample     uint8 // (8, 12, 16, or 32)
			Applied_corrections uint32
			Spare               [4]uint32 // 16 bytes
		} // 21 bytes
		base2 struct {
			Sample_count  uint16
			Detect_sample uint16
			Start_range   uint16
			Spare         [3]uint16 // 6 bytes
		} // 12 bytes
		count  []uint16
		detect []uint16
		st_rng []uint16
		// detect_val   []float32
		timeseries   []float64
		samples_u1   []uint8
		samples_u2   []uint16
		samples_u4   []uint32
		samples_f64  []float64
		three_bytes  [3]byte
		unpack_bytes [4]byte
		scl_off      ScaleOffset
		ts_mean      []float64
		// n_bytes      int64
		// img_md      SensorImageryMetadata
		// timeseries [][]float32
	)

	count = make([]uint16, 0, nbeams)
	detect = make([]uint16, 0, nbeams)
	// detect_val = make([]float32, 0, nbeams)
	st_rng = make([]uint16, 0, nbeams)
	timeseries = make([]float64, 0, nbeams*66) // 66 ... just becasuse
	ts_mean = make([]float64, 0, nbeams)

	_ = binary.Read(reader, binary.BigEndian, &base)

	switch sensor_id {

	case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
		// DecodeEM3Imagery
		em3img, scl__off, err := DecodeEm3Imagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Em3_imagery = em3img
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case RESON_7125:
		// DecodeReson7100Imagery
		reson7100, scl__off, err := DecodeReson7100Imagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Reson7100_imagery = reson7100
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case RESON_TSERIES:
		// DecodeResonTSeriesImagery
		tseries, scl__off, err := DecodeResonTSeriesImagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.ResonTSeries_imagery = tseries
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case RESON_8101, RESON_8111, RESON_8124, RESON_8125, RESON_8150, RESON_8160:
		// DecodeReson8100Imagery
		reson8100, scl__off, err := DecodeReson8100Imagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Reson8100_imagery = reson8100
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case EM122, EM302, EM710, EM2040:
		// DecodeEM4Imagery
		em4img, scl__off, err := DecodeEm4Imagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Em4_imagery = em4img
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case KLEIN_5410_BSS:
		// DecodeKlein5410BssImagery
		klein, scl__off, err := DecodeKlein5410BssImagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Klein5410Bss_imagery = klein
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case KMALL:
		// DecodeKMALLImagery
		kmall, scl__off, err := DecodeKmallImagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.Kmall_imagery = kmall
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		// DecodeR2SonicImagery
		r2sonic, scl__off, err := DecodeR2SonicImagery(reader)
		if err != nil {
			return intensity, img_md, err
		}
		img_md.R2Sonic_imagery = r2sonic
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
	}

	bytes_per_sample := base.Bits_per_sample / 8

	for i := uint16(0); i < nbeams; i++ {
		_ = binary.Read(reader, binary.BigEndian, &base2)

		count = append(count, base2.Sample_count)
		detect = append(detect, base2.Detect_sample)
		st_rng = append(st_rng, base2.Start_range)

		// the spec only mentioned 1 byte per sample, so the following
		// is an implementation (of sorts) based on the gsf code base
		if base.Bits_per_sample == 12 {
			// TODO
			samples_f64 = make([]float64, base2.Sample_count)
			samples_u4 = make([]uint32, base2.Sample_count)

			// 3 bytes of data are bit compacted and decompress into 2 samples
			for j := uint16(0); j < nbeams; j += 2 {
				_ = binary.Read(reader, binary.BigEndian, &three_bytes)

				// unpacking the first sample

				// upper bits of 3b[0] into lowerbits of unpack[2]
				unpack_bytes[2] = three_bytes[0] >> 4

				// lower bits of 3b[1] into upper bits of unpack[3]
				unpack_bytes[3] = (three_bytes[0] & 0x0f) << 4

				// upper bits of 3b[1] combine into lower bits of unpack[3]
				unpack_bytes[3] |= (three_bytes[1] & 0xf0) >> 4

				samples_u4[j] = binary.BigEndian.Uint32(unpack_bytes[:])

				if j+1 < nbeams {
					// unpacking the second sample

					// lower bits of tb[1] into unpack[2]
					unpack_bytes[2] = three_bytes[1] & 0x0f

					// tb[2] into unpack[3]
					unpack_bytes[3] = three_bytes[2]

					samples_u4[j+1] = binary.BigEndian.Uint32(unpack_bytes[:])
				}
			}

			for k, v := range samples_u4 {
				samples_f64[k] = float64(v)
			}
		} else {
			samples_f64 = make([]float64, base2.Sample_count)
			switch bytes_per_sample {
			case 1:
				samples_u1 = make([]uint8, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u1)
				for k, v := range samples_u1 {
					samples_f64[k] = float64(v)
				}
			case 2:
				samples_u2 = make([]uint16, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u2)
				for k, v := range samples_u2 {
					samples_f64[k] = float64(v)
				}
			case 4:
				samples_u4 = make([]uint32, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u4)
				for k, v := range samples_u4 {
					samples_f64[k] = float64(v)
				}
			}
		}

		// apply scale and offset
		// The gsf spec mentions that the scale factor is 2 for EM3 and
		// 10 for EM4 based sensors. Ideally the stored value should be
		// used, unfortunately, some of the sample files had incorrect
		// scale factors due to a bug in the source software that
		// generated the file.
		// dB_value = (value - offset) / scale
		switch sensor_id {
		case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
			for k, v := range samples_f64 {
				samples_f64[k] = (v - scl_off.Offset) / float64(2)
			}
		case EM122, EM302, EM710, EM2040:
			for k, v := range samples_f64 {
				samples_f64[k] = (v - scl_off.Offset) / SCALE_1_F64
			}
		}

		// other sensor types don't contain a scale and offset value in their
		// imagery metadata. so i guess we do nothing

		// need to handle the case where base2.Sample_count == 0
		// not recording a value will upset the variable length offsets
		// as well as the beam locations
		// NaN should be fine to indicate a missing observation
		if base2.Sample_count == 0 {
			samples_f64 = make([]float64, 1)
			samples_f64[0] = math.NaN()
			ts_mean = append(ts_mean, math.NaN())
		} else {
			ts_mean = append(ts_mean, lo.Mean(samples_f64))
		}

		// append
		// detect_val = append(detect_val, samples_f64[base2.Detect_sample])
		timeseries = append(timeseries, samples_f64...)
	}

	intensity.TimeSeries = timeseries
	// intensity.BottomDetect = detect_val
	intensity.StartRange = st_rng
	intensity.BottomDetectIndex = detect
	intensity.TsMean = ts_mean
	intensity.sample_count = count

	return intensity, img_md, err
}
