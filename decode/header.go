package decode

// Header contains the version of GSF used to construct the GSF file.
type Header struct {
    Version string
}

// NewHeader constructs a Header by decoding the HEADER record from a GSF file.
// It contains the version of GSF used to create the file.
func NewHeader(buffer []byte) *Header {
    // buffer := make([]byte, rec.Datasize)

    // _ = binary.Read(stream, binary.BigEndian, &buffer)

    file_hdr := Header{Version: string(buffer)}

    return new(file_hdr)
}
