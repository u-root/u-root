// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/u-root/uio/uio"
	"golang.org/x/sys/unix"
)

func mmap(f *os.File) ([]byte, func() error, error) {
	s, err := f.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("stat error: %w", err)
	}
	if s.Size() == 0 {
		return nil, nil, fmt.Errorf("%w: cannot mmap zero-len file", os.ErrInvalid)
	}
	d, err := unix.Mmap(int(f.Fd()), 0, int(s.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, nil, fmt.Errorf("mmap failed: %w", err)
	}

	ummap := func() error {
		if err := unix.Munmap(d); err != nil {
			return fmt.Errorf("failed to unmap %s: %w", f.Name(), err)
		}
		return nil
	}
	return d, ummap, nil
}

func getFile(f *os.File) ([]byte, func() error, error) {
	if d, unmap, err := mmap(f); err == nil {
		return d, unmap, nil
	}
	var d []byte
	var err error
	// Pipes and other files like that will fail to seek.
	if _, serr := f.Seek(0, 0); serr != nil {
		d, err = io.ReadAll(f)
	} else {
		d, err = uio.ReadAll(f)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read kernel file: %w", err)
	}
	return d, func() error { return nil }, nil
}
