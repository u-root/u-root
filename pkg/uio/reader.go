package uio

import (
	"io"
	"io/ioutil"
	"math"
)

type inMemReaderAt interface {
	Bytes() []byte
}

// ReadAll reads everything that r contains.
//
// Callers *must* not modify bytes in the returned byte slice.
//
// If r is an in-memory representation, ReadAll will attempt to return a
// pointer to those bytes directly.
func ReadAll(r io.ReaderAt) ([]byte, error) {
	if imra, ok := r.(inMemReaderAt); ok {
		return imra.Bytes(), nil
	}
	return ioutil.ReadAll(Reader(r))
}

// Reader generates a Reader from a ReaderAt.
func Reader(r io.ReaderAt) io.Reader {
	return io.NewSectionReader(r, 0, math.MaxInt64)
}
