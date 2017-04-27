// stateless implements stateless Files.

package build

import (
	"io"
	"log"
	"os"
)

type stateless struct {
	name string
	off  int64
}

var (
	Files int
	Opens int
)

func (s *stateless) Read(b []byte) (int, error) {
	f, err := os.Open(s.name)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	n, err := f.ReadAt(b, s.off)
	log.Printf("Read %v at off %v n %v err %v", s.name, s.off, n, err)
	if err != nil {
		s.off += int64(n)
	}
	Files++
	return n, err
}

func StatelessOpen(n string) (io.Reader, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	Opens++
	defer f.Close()
	return &stateless{name: n}, nil
}
