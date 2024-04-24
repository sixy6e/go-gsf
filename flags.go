package gsf

import (
	"bytes"
	"encoding/binary"
)

// DecodeBeamFlagsArray decodes the beam flags array subrecord.
// The length of the returned slice is determined by the input
// number of beams.
// Each element indicates whether or not the beam contains usable data.
func DecodeBeamFlagsArray(reader *bytes.Reader, nbeams uint16) []uint8 {
	var (
		data []uint8
	)

	data = make([]uint8, nbeams)

	_ = binary.Read(reader, binary.BigEndian, &data)

	return data
}
