package uio

import (
	"io"
	"math"
)

// Reader generates a Reader from a ReaderAt.
func Reader(r io.ReaderAt) io.Reader {
	return io.NewSectionReader(r, 0, math.MaxInt64)
}
