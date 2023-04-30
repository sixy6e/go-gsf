package decode

import (
    "os"
    //"bytes"
    "encoding/binary"
)

// Record contains information about a given record stored within the GSF file.
// It contains the record identifier, the size of the data within the record,
// a byte index within the file of where the data starts for the record
// as well as an indicator as to whether or not a checksum is given for the record.
type Record struct {
    Id RecordID
    Datasize uint32
    Byte_index int64
    Checksum_flag bool
    // Record_index uint64
}

// RecordHdr decodes the header part of any given record.
// Each record has a small header that defines the type of record, the size
// of the data within the record, and whether the record contains a checksum
func RecordHdr(stream *os.File) Record {
    
    blob := [2]uint32{}
    _ = binary.Read(stream, binary.BigEndian, &blob)
    data_size := blob[0]
    record_id := RecordID(blob[1])
    checksum_flag := int64(record_id) & 0x80000000 == 1

    pos, _ := Tell(stream)

    rec_hdr := Record{
        Id: record_id,
        Datasize: data_size,
        Byte_index: pos,
        Checksum_flag: checksum_flag,
    }

    return rec_hdr
}
