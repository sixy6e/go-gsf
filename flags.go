package gsf

import (
	"bytes"
	"encoding/binary"
)

// DecodeBeamFlagsArray decodes the beam flags array subrecord.
// The length of the returned slice is determined by the input
// number of beams.
// Each element indicates whether or not the beam contains usable data.
func DecodeBeamFlagsArray(reader *bytes.Reader, nbeams uint16) ([]uint8, int64) {
	var (
		data    []uint8
		n_bytes int64
	)

	data = make([]uint8, nbeams)
	n_bytes = 0

	_ = binary.Read(reader, binary.BigEndian, &data)
	n_bytes += 1 * int64(nbeams)

	return data, n_bytes
}
