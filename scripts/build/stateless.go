// stateless implements stateless files.

package build

import (
	"fmt"
	"io"
	"os"
)

type stateless struct {
	name string
	off  int64
}

var (
	files int
	opens int
)

// Read implements io.Reader for stateless files.
// Unlike many Read functions it can return errors
// from Open as well as ReadAt.
func (s *stateless) Read(b []byte) (int, error) {
	f, err := os.Open(s.name)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	n, err := f.ReadAt(b, s.off)
	if n > 0 {
		s.off += int64(n)
	}
	opens++
	return n, err
}

// Stateless open tries to open a file. If it succeeds,
// it returns a stateless struct, else the error from os.Open.
// The file is closed on return so we don't get EMFILE errors
// on some OSes.
func StatelessOpen(n string) (io.Reader, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	files++
	defer f.Close()
	return &stateless{name: n}, nil
}

// StatelessStats returns a string containing information about
// stateless file activity.
func StatelessStats() string {
	return fmt.Sprintf("Number of files used %d Number of opens %d", files, opens)
}
