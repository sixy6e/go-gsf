package gsf

import (
	// "os"
	"bytes"
	"encoding/binary"
	// tiledb "github.com/TileDB-Inc/TileDB-Go"
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
}

// DecodeRecordHdr acts as the constructor for RecordHdr by decoding the header of
// a records byte stream.
// Each record has a small header that defines the type of record, the size
// of the data within the record, and whether the record contains a checksum
func DecodeRecordHdr(stream Stream) RecordHdr {

	blob := [2]uint32{}
	_ = binary.Read(stream, binary.BigEndian, &blob)
	data_size := blob[0]
	record_id := RecordID(blob[1])
	bits := int64(record_id) & 0x80000000
	checksum_flag := bits == 1

	pos, _ := Tell(stream)

	rec_hdr := RecordHdr{
		Id:            record_id,
		Datasize:      data_size,
		Byte_index:    pos,
		Checksum_flag: checksum_flag,
	}

	return rec_hdr
}

type SubRecord struct {
	Id         SubRecordID
	Datasize   uint32
	Byte_index int64
}

// DecodeByteArray decodes the beam array data stored as uint8.
func (sr *SubRecord) DecodeByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float32) {
	var (
		data []uint8
	)

	data = make([]uint8, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 1 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeSignedByteArray decodes the beam array data stored as int8.
func (sr *SubRecord) DecodeSignedByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float32) {
	var (
		data []int
	)

	data = make([]int, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 1 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeTwoByteArray decodes the beam array data stored as uint16.
func (sr *SubRecord) DecodeTwoByteArray(reader *bytes.Reader, number_beams uint16, scale_factor ScaleFactor) (scaled_data []float32) {
	var (
		data []uint16
	)

	data = make([]uint16, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 2 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeSignedTwoByteArray decodes the beam array data stored as int16.
func (sr *SubRecord) DecodeSignedTwoByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float32) {
	var (
		data []int16
	)

	data = make([]int16, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 2 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeFourByteArray decodes the beam array data stored as uint32.
func (sr *SubRecord) DecodeFourByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float32) {
	var (
		data []uint32
	)

	data = make([]uint32, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 4 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeSignedFourTwoByteArray decodes the beam array data stored as int32.
func (sr *SubRecord) DecodeSignedFourByteArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
) (scaled_data []float32) {
	var (
		data []int32
	)

	data = make([]int32, number_beams)
	scaled_data = make([]float32, number_beams)

	_ = binary.Read(reader, binary.BigEndian, &data)
	// n_bytes = 4 * int64(number_beams)
	// if err != nil {
	//
	// }

	for k, v := range data {
		scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
	}

	return scaled_data // , n_bytes
}

// DecodeSubRecArray decodes the beam array data.
// For the time being, the data is scaled to float32.
// Further testing is required to determine if float64 is required.
// Float32 was chosen as the sample data provided (several dozen TB)
// didn't exhibit a need for float64 precision.
// This method may be overkill in handling {1, 2, 4}byte, {signed, unsigned}
// arrays. Individual funcs may be more suitable.
func (sr *SubRecord) DecodeSubRecArray(
	reader *bytes.Reader,
	number_beams uint16,
	scale_factor ScaleFactor,
	bytes_per_beam uint32,
	signed bool,
) (scaled_data []float32) {
	scaled_data = make([]float32, number_beams) // will float32 be enough???

	switch signed {
	case true:
		switch bytes_per_beam {
		case BYTES_PER_BEAM_ONE:
			data := make([]int8, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			// n_bytes = 1 * int64(number_beams)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		case BYTES_PER_BEAM_TWO:
			data := make([]int16, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			// n_bytes = 2 * int64(number_beams)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		case BYTES_PER_BEAM_FOUR:
			data := make([]int32, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			// n_bytes = 4 * int64(number_beams)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		}
	case false:
		switch bytes_per_beam {
		case BYTES_PER_BEAM_ONE:
			data := make([]uint8, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		case BYTES_PER_BEAM_TWO:
			data := make([]uint16, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		case BYTES_PER_BEAM_FOUR:
			data := make([]uint32, number_beams)
			_ = binary.Read(reader, binary.BigEndian, &data)
			for k, v := range data {
				scaled_data[k] = float32(v)/scale_factor.Scale - scale_factor.Offset
			}
		}
	}

	return scaled_data // , n_bytes
}
