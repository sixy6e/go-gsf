package gsf

import (
	// "os"
	"bytes"
	"encoding/binary"
	// "reflect"
	// "fmt"
	// stgpsr "github.com/yuin/stagparser"
)

// Removing BrbIntensity.BottomDetect for now. Need to get more info on what the
// BottomDetectIndex means. I thought it would be the index location of the TimeSeries
// slice, but index values have included values that are outside the array.
// i.e TimeSeries samples = 6, BottomDetectIndex = 6
// This would be ok if it is 1-based indexing stored in the GSF file for this particular
// piece of data, except that the BottomDetectIndex also contains values of 0.

// BrbIntensity
type BrbIntensity struct {
	TimeSeries []float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	// BottomDetect      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	BottomDetectIndex []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartRange        []uint16 `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
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
		make([]float32, 0, number_beams*66), // 66 ... just becasuse
		// make([]float32, 0, number_beams),
		make([]uint16, 0, number_beams),
		make([]uint16, 0, number_beams),
		make([]uint16, 0, number_beams),
	}
	return brb_int
}

// DecodeBrbIntensity decodes the timeseries intensity sub-record. Each beam will have
// a variable length of intensity samples, the index for the bottom detect sample, and the
// sample itself.
func DecodeBrbIntensity(reader *bytes.Reader, nbeams uint16, sensor_id SubRecordID) (intensity BrbIntensity, img_md SensorImageryMetadata) {

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
		timeseries   []float32
		samples_u1   []uint8
		samples_u2   []uint16
		samples_u4   []uint32
		samples_f32  []float32
		three_bytes  [3]byte
		unpack_bytes [4]byte
		scl_off      ScaleOffset
		// n_bytes      int64
		// img_md      SensorImageryMetadata
		// timeseries [][]float32
	)

	count = make([]uint16, 0, nbeams)
	detect = make([]uint16, 0, nbeams)
	// detect_val = make([]float32, 0, nbeams)
	st_rng = make([]uint16, 0, nbeams)
	timeseries = make([]float32, 0, nbeams*66) // 66 ... just becasuse
	// timeseries = make([][]float32, nbeams)
	// nbytes = 0

	// reader := bytes.NewReader(buffer)

	_ = binary.Read(reader, binary.BigEndian, &base)
	// nbytes += 21

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
		em4img, scl__off := DecodeEM4Imagery(reader)
		img_md.EM4_imagery = em4img
		scl_off.Scale = scl__off.Scale
		scl_off.Offset = scl__off.Offset
		// nbytes += n_bytes
	case KLEIN_5410_BSS:
		// DecodeKlein5410BssImagery
	case KMALL:
		// DecodeKMALLImagery
	case R2SONIC_2020, R2SONIC_2022, R2SONIC_2024:
		// DecodeR2SonicImagery
	}

	bytes_per_sample := base.Bits_per_sample / 8

	for i := uint16(0); i < nbeams; i++ {
		_ = binary.Read(reader, binary.BigEndian, &base2)
		// nbytes += 12

		count = append(count, base2.Sample_count)
		detect = append(detect, base2.Detect_sample)
		st_rng = append(st_rng, base2.Start_range)

		// the spec only mentioned 1 byte per sample, so the following
		// is an implementation (of sorts) based on the gsf code base
		if base.Bits_per_sample == 12 {
			// TODO
			samples_f32 = make([]float32, base2.Sample_count)
			samples_u4 = make([]uint32, base2.Sample_count)

			// 3 bytes of data are bit compacted and decompress into 2 samples
			for j := uint16(0); j < nbeams; j += 2 {
				_ = binary.Read(reader, binary.BigEndian, &three_bytes)
				// nbytes += 3

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
				samples_f32[k] = float32(v)
			}
		} else {
			samples_f32 = make([]float32, base2.Sample_count)
			switch bytes_per_sample {
			case 1:
				samples_u1 = make([]uint8, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u1)
				// n_bytes += 1 * int64(base2.Sample_count)
				for k, v := range samples_u1 {
					samples_f32[k] = float32(v)
				}
			case 2:
				samples_u2 = make([]uint16, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u2)
				// n_bytes += 2 * int64(base2.Sample_count)
				for k, v := range samples_u2 {
					samples_f32[k] = float32(v)
				}
			case 4:
				samples_u4 = make([]uint32, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u4)
				// n_bytes += 4 * int64(base2.Sample_count)
				for k, v := range samples_u4 {
					samples_f32[k] = float32(v)
				}
			}

		}
		// apply scale and offset
		// The gsf spec mentions that the scale factor is 2 for EM3 and
		// 10 for EM4 based sensors. Ideally the stored value should be
		// used, unfortunately, some of the sample files had incorrect
		// scale factors due to a bug in the source software that
		// generated the file.
		switch sensor_id {
		case EM120, EM120_RAW, EM300, EM300_RAW, EM1002, EM1002_RAW, EM2000, EM2000_RAW, EM3000, EM3000_RAW, EM3002, EM3002_RAW, EM3000D, EM3000D_RAW, EM3002D, EM3002D_RAW, EM121A_SIS, EM121A_SIS_RAW:
			// TODO; loop over length, as range will copy the array
			for k, v := range samples_f32 {
				samples_f32[k] = (v - float32(img_md.EM3_imagery.offset[0])) / float32(2)
			}
		case EM122, EM302, EM710, EM2040:
			// TODO; loop over length, as range will copy the array
			for k, v := range samples_f32 {
				samples_f32[k] = (v - scl_off.Offset) / float32(10)
			}
		}
		// append
		// detect_val = append(detect_val, samples_f32[base2.Detect_sample])
		// timeseries[i] = samples_f32
		timeseries = append(timeseries, samples_f32...)
	}

	intensity.TimeSeries = timeseries
	// intensity.BottomDetect = detect_val
	intensity.StartRange = st_rng
	intensity.BottomDetectIndex = detect
	intensity.sample_count = count

	return intensity, img_md // , n_bytes
}
