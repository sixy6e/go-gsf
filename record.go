package gsf

import (
	"bytes"
	"encoding/binary"
)

// RecordHdr contains information about a given record stored within the GSF file.
// It contains the record identifier, the size of the data within the record,
// a byte index within the file of where the data starts for the record
// as well as an indicator as to whether or not a checksum is given for the record.
type RecordHdr struct {
	Id            RecordID
	Datasize      uint32
	Byte_index    int64
	Checksum_flag bool
	Reserved      uint32
}

// DecodeRecordHdr acts as the constructor for RecordHdr by decoding the header of
// a records byte stream.
// Each record has a small header that defines the type of record, the size
// of the data within the record, and whether the record contains a checksum
func DecodeRecordHdr(stream Stream) RecordHdr {

	blob := [2]uint32{}
	_ = binary.Read(stream, binary.BigEndian, &blob)
	data_size := blob[0]
	record_id := RecordID(blob[1] & 0x003FFFFF)
	bits := blob[1] & 0x80000000
	checksum_flag := bits == 1
	reserved := (blob[1] & 0x7FC00000) >> 22

	pos, _ := Tell(stream)

	rec_hdr := RecordHdr{
		Id:            record_id,
		Datasize:      data_size,
		Byte_index:    pos,
		Checksum_flag: checksum_flag,
		Reserved:      reserved,
	}

	return rec_hdr
}

// SubRecord contains information pertaining to the SubRecord, such as the ID,
// the size in bytes of the record, and where does the SubRecord start as a byte
// index location.
type SubRecord struct {
	Id         SubRecordID
	Datasize   uint32
	Byte_index int64
}

// apply_scale_factor is a helper function that applies the scale offset
// inversion for data encoded into GSF.
// unscaled = value / scale - offset
func apply_scale_factor(value float64, scl_off ScaleFactor) (unscaled float64) {
	unscaled = value/scl_off.Scale - scl_off.Offset
	return unscaled
}

// f64_to_f32 is a helper func for converting a slice of type float64 to float32.
func f64_to_f32(f64 []float64) (f32 []float32) {
	n := len(f64)
	f32 = make([]float32, n)

	for i := 0; i < n; i++ {
		f32[i] = float32(f64[i])
	}

	return f32
}

// f64_to_u16 is a helper func for converting a slice of type float64 to uint16.
// Further testing is required on whether this is the correct approach,
// or shall we instead convert the scale and offset to uint16?
func f64_to_u16(f64 []float64) (u16 []uint16) {
	n := len(f64)
	u16 = make([]uint16, n)

	for i := 0; i < n; i++ {
		u16[i] = uint16(f64[i])
	}

	return u16
}

// DecodeByteArray decodes the beam array data stored as uint8.
func (sr *SubRecord) DecodeByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float64) {
	var (
		data []uint8
	)

	data = make([]uint8, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
	}

	return scaled_data
}

// DecodeSignedByteArray decodes the beam array data stored as int8.
func (sr *SubRecord) DecodeSignedByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float64) {
	var (
		data []int
	)

	data = make([]int, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
	}

	return scaled_data
}

// DecodeTwoByteArray decodes the beam array data stored as uint16.
func (sr *SubRecord) DecodeTwoByteArray(reader *bytes.Reader, number_beams uint16, scale_factor ScaleFactor) (scaled_data []float64) {
	var (
		data []uint16
	)

	data = make([]uint16, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
	}

	return scaled_data
}

// DecodeSignedTwoByteArray decodes the beam array data stored as int16.
func (sr *SubRecord) DecodeSignedTwoByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float64) {
	var (
		data []uint16
	)

	data = make([]uint16, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(int16(v)), scale_factor)
	}

	return scaled_data
}

// DecodeFourByteArray decodes the beam array data stored as uint32.
func (sr *SubRecord) DecodeFourByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float64) {
	var (
		data []uint32
	)

	data = make([]uint32, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
	}

	return scaled_data
}

// DecodeSignedFourTwoByteArray decodes the beam array data stored as int32.
func (sr *SubRecord) DecodeSignedFourByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float64) {
	var (
		data []uint32
	)

	data = make([]uint32, number_beams)
	scaled_data = make([]float64, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	for k, v := range data {
		scaled_data[k] = apply_scale_factor(float64(int32(v)), scale_factor)
	}

	return scaled_data
}

// DecodeSubRecArray decodes the beam array data.
// The scaled data is unscaled into float64.
// Whilst float64 might be overkill, (definitely is for some beam array types),
// it more closely follows the original C-code.
// In some sense, as the decoding could be 1, 2 or 4 byte, it makes code simpler
// to simple decode all into float64 and apply strict datatype conversion elsewhere.
// Additional work is required in order to determine more appropriate datatypes
// for the differing beam arrays.
// Rather than handle it here, it will be left up to the caller to convert as
// where necessary.
// This method may be overkill in handling {1, 2, 4}byte, {signed, unsigned}
// arrays. Individual funcs may be more suitable.
func (sr *SubRecord) DecodeSubRecArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
	bytes_per_beam uint32,
	signed bool,
) (scaled_data []float64) {
	scaled_data = make([]float64, number_beams)

	switch signed {
	case true:
		switch bytes_per_beam {
		case BYTES_PER_BEAM_ONE:
			data := make([]int8, number_beams) // reference decoded direct into int8
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
			}
		case BYTES_PER_BEAM_TWO:
			data := make([]uint16, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(int32(v)), scale_factor)
			}
		case BYTES_PER_BEAM_FOUR:
			data := make([]uint32, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(int32(v)), scale_factor)
			}
		}
	case false:
		switch bytes_per_beam {
		case BYTES_PER_BEAM_ONE:
			data := make([]uint8, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
			}
		case BYTES_PER_BEAM_TWO:
			data := make([]uint16, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
			}
		case BYTES_PER_BEAM_FOUR:
			data := make([]uint32, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = apply_scale_factor(float64(v), scale_factor)
			}
		}
	}

	return scaled_data
}
