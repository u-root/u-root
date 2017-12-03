package gzip

import (
	"io"

	"github.com/klauspost/pgzip"
)

func decompress(r io.Reader, w io.Writer, blocksize int, processes int) error {
	zr, err := pgzip.NewReaderN(r, blocksize*1024, processes)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, zr); err != nil {
		zr.Close()
		return err
	}

	return zr.Close()
}
