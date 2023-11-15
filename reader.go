package gsf

import (
	"bytes"
	"encoding/binary"

	tiledb "github.com/TileDB-Inc/TileDB-Go"
)

// Stream caters for a generic reader type so that we can handle both
// a stream of data from a file on disk or object store, as well as
// an in-memory byte stream.
// This GSF module deals with either a *tiledb.VFSfh or *bytes.Reader,
// and all we care about are two methods, Read and Seek,
// which both implement.
type Stream interface {
	Read(p []byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
}

// function to handle whether we build an in-memory byte stream or leave
// is as stream handled by *tiledb.VFSfh
func GenericStream(stream *tiledb.VFSfh, size uint64, inmem bool) (Stream, error) {
	if inmem {
		buffer := make([]byte, size)
		err := binary.Read(stream, binary.BigEndian, &buffer)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buffer)
		return reader, nil
	} else {
		return stream, nil
	}
}
