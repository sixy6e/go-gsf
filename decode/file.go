package decode

import "os"

// Tell is a small helper fucntion for telling the current position within a
// binary file opened for reading.
func Tell(stream *os.File) int64 {
    pos, _ := stream.Seek(0, 1)

    return pos
}

// Padding is a small helper function for padding a GSF record.
// The GSF specification mentions that a records complete length has to be
// a multiple of 4.
func Padding(stream *os.File) {
    pos := tell(stream)
    pad := pos % 4
    pad, _ = stream.Seek(pad, 1)

    return 
}

// GSFv contains the version of GSF used to construct the GSF file.
type GSFv struct {
    Version string
}

// FileRec decodes the HEADER record from a GSF file.
// It contains the version of GSF used to create the file.
func FileRec(stream *os.File, rec Record) GSFv {
    buffer := make([]byte, rec.Datasize)

    _ = binary.Read(stream, binary.BigEndian, &buffer)

    file_hdr := GSFv{Version: string(buffer)}

    return file_hdr
}

// PingInfo contains some basic information regarding the ping such as
// the number of beams, what sub-records are populated.
// The initial reasoning behind why, is to provide a basic descriptor
// to inform a global schema across all pings, and derive max(n_beams) to
// inform a global [ping, beam] dimensional array structure.
type PingInfo struct {
    Number_Beams uint16
    Sub_Records []SubRecordID
    Scale_Factors bool
}

// Index, as the name implies, builds a file index of all Record types.
// Each Record contains the record ID, record size, byte index and checksum flag.
func Index(stream, *os.File) any {  // TODO return type(s)
    // probably return map[RecordId]uint64, map[SubRecordId][uint64], []PingInfo
}
