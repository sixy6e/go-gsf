package gsf

import (
	// "os"
	"bytes"
	"encoding/binary"
	// "reflect"
	// "fmt"
	// stgpsr "github.com/yuin/stagparser"
)

type BrbIntensity struct {
	TimeSeries        []float32 `tiledb:"dtype=float32,ftype=attr,var" filters:"zstd(level=16)"`
	BottomDetect      []float32 `tiledb:"dtype=float32,ftype=attr" filters:"zstd(level=16)"`
	BottomDetectIndex []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	StartRange        []uint16  `tiledb:"dtype=uint16,ftype=attr" filters:"zstd(level=16)"`
	sample_count      []uint16
	// timeseries        [][]float32
}

// DecocdeBrbIntensity decodes the timeseries intensity sub-record. Each beam will have
// a variable length of intensity samples, the index for the bottom detect sample, and the
// sample itself.
func DecocdeBrbIntensity(buffer []byte, nbeams uint16, sensor_id SubRecordID) (intensity BrbIntensity) {

	var (
		base struct {
			Bits_per_sample     uint8 // (8, 12, 16, or 32)
			Applied_corrections uint32
			Spare               [4]uint32 // 16 bytes
		}
		base2 struct {
			Sample_count  uint16
			Detect_sample uint16
			Start_range   uint16
			Spare         [3]uint16 // 6 bytes
		}
		count       []uint16
		detect      []uint16
		st_rng      []uint16
		detect_val  []float32
		timeseries  []float32
		samples_u1  []uint8
		samples_u2  []uint16
		samples_u4  []uint32
		samples_f32 []float32
		img_md      SensorImageryMetadata
		// timeseries [][]float32
	)

	count = make([]uint16, 0, nbeams)
	detect = make([]uint16, 0, nbeams)
	detect_val = make([]float32, 0, nbeams)
	st_rng = make([]uint16, 0, nbeams)
	timeseries = make([]float32, 0, nbeams*66) // 66 ... just becasuse
	// timeseries = make([][]float32, nbeams)

	reader := bytes.NewReader(buffer)

	_ = binary.Read(reader, binary.BigEndian, &base)

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
		img_md.EM4Imagery = DecodeEM4Imagery(reader)
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

		count = append(count, base2.Sample_count)
		detect = append(detect, base2.Detect_sample)
		st_rng = append(st_rng, base2.Start_range)

		if base.Bits_per_sample == 12 {
			// TODO
		} else {
			// for j := uint16(0); j < base2.Sample_count; j++ {
			// 	samples_f32 = make([]float32, base2.Sample_count)

			// 	switch bytes_per_sample {
			// 	case 1:
			// 		samples_u1 = make([]uint8, base2.Sample_count)
			// 		_ = binary.Read(reader, binary.BigEndian, samples_u1)
			// 	case 2:
			// 		samples_u2 = make([]uint16, base2.Sample_count)
			// 		_ = binary.Read(reader, binary.BigEndian, samples_u2)
			// 	case 4:
			// 		samples_u4 = make([]uint32, base2.Sample_count)
			// 		_ = binary.Read(reader, binary.BigEndian, samples_u4)
			// 	}

			// }
			samples_f32 = make([]float32, base2.Sample_count)
			switch bytes_per_sample {
			case 1:
				samples_u1 = make([]uint8, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u1)
				for k, v := range samples_u1 {
					samples_f32[k] = float32(v)
				}
			case 2:
				samples_u2 = make([]uint16, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u2)
				for k, v := range samples_u2 {
					samples_f32[k] = float32(v)
				}
			case 4:
				samples_u4 = make([]uint32, base2.Sample_count)
				_ = binary.Read(reader, binary.BigEndian, samples_u4)
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
			for k, v := range samples_f32 {
				samples_f32[k] = (v - float32(img_md.EM3Imagery.offset[0])) / float32(2)
			}
		case EM122, EM302, EM710, EM2040:
			for k, v := range samples_f32 {
				samples_f32[k] = (v - float32(img_md.EM4Imagery.offset[0])) / float32(10)
			}
		}
		// append
		detect_val = append(detect_val, samples_f32[base2.Detect_sample])
		// timeseries[i] = samples_f32
		timeseries = append(timeseries, samples_f32...)
	}

	intensity.TimeSeries = timeseries
	intensity.BottomDetect = detect_val
	intensity.StartRange = st_rng
	intensity.BottomDetectIndex = detect
	intensity.sample_count = count

	return intensity
}
