// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/ubinary"
)

func pathRead(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Read(f, ubinary.NativeEndian, data)
}

func pathWrite(path string, addr int64, data UintN) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Write(f, ubinary.NativeEndian, data)
}
