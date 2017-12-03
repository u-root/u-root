package gzip

import (
	"io"

	"github.com/klauspost/pgzip"
)

func compress(r io.Reader, w io.Writer, level int, blocksize int, processes int) error {
	zw, err := pgzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}

	if err := zw.SetConcurrency(blocksize*1024, processes); err != nil {
		zw.Close()
		return err
	}

	if _, err := io.Copy(zw, r); err != nil {
		zw.Close()
		return err
	}

	return zw.Close()
}
