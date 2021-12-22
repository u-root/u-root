// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/ubinary"
)

type MemIOReader interface {
	Read(UintN, int64) error
}

type MemIOWriter interface {
	Write(UintN, int64) error
}

type MemIO interface {
	MemIOReader
	MemIOWriter
	io.Closer
}

type MemIOPort struct {
	*os.File
}

func (m *MemIOPort) Read(out UintN, addr int64) error {
	if _, err := m.File.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Read(m.File, ubinary.NativeEndian, out)
}

func (m *MemIOPort) Write(in UintN, addr int64) error {
	if _, err := m.File.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Write(m.File, ubinary.NativeEndian, in)
}

func (m *MemIOPort) Close() error {
	return m.File.Close()
}

func NewMemIOPort(f *os.File) (*MemIOPort, error) {
	return &MemIOPort{
		File: f,
	}, nil
}
