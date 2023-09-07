package decode

// Header contains the version of GSF used to construct the GSF file.
type Header struct {
    Version string
}

// DecodeHeader constructs a Header by decoding the HEADER record from a GSF file.
// It contains the version of GSF used to create the file.
func DecodeHeader(buffer []byte) Header {
    // buffer := make([]byte, rec.Datasize)

    // _ = binary.Read(stream, binary.BigEndian, &buffer)

    file_hdr := Header{Version: string(buffer)}

    return file_hdr
}
