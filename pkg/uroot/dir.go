package uroot

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

// DirArchiver implements Archiver for a directory.
type DirArchiver struct{}

// Reader implements Archiver.Reader.
//
// Currently unsupported for directories.
func (da DirArchiver) Reader(io.ReaderAt) ArchiveReader {
	return nil
}

// OpenWriter implements Archiver.OpenWriter.
func (da DirArchiver) OpenWriter(path, goos, goarch string) (ArchiveWriter, error) {
	if len(path) == 0 {
		var err error
		path, err = ioutil.TempDir("", "u-root")
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := os.Stat(path); os.IsExist(err) {
			return nil, fmt.Errorf("path %q already exists", path)
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}
	log.Printf("Path is %s", path)
	return dirWriter{path}, nil
}

// dirWriter implements ArchiveWriter.
type dirWriter struct {
	dir string
}

// WriteRecord implements ArchiveWriter.WriteRecord.
func (dw dirWriter) WriteRecord(r cpio.Record) error {
	return cpio.CreateFileInRoot(r, dw.dir)
}

// Finish implements ArchiveWriter.Finish.
func (dw dirWriter) Finish() error {
	return nil
}
