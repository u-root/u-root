// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"encoding/binary"
	"io"
	"os"
)

// Reader is the interface for reading from memory and IO ports.
type Reader interface {
	Read(UintN, int64) error
}

// Writer is the interface for writing to memory and IO ports.
type Writer interface {
	Write(UintN, int64) error
}

// ReadWriteCloser implements io.ReadWriteCloser
type ReadWriteCloser interface {
	Reader
	Writer
	io.Closer
}

// Port implements memory and IO port access via an os.File.
type Port struct {
	*os.File
}

var _ ReadWriteCloser = &Port{}

// Read implements Reader for a Port
func (m *Port) Read(out UintN, addr int64) error {
	if _, err := m.File.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Read(m.File, binary.NativeEndian, out)
}

// Write implements Writer for a Port
func (m *Port) Write(in UintN, addr int64) error {
	if _, err := m.File.Seek(addr, io.SeekStart); err != nil {
		return err
	}
	return binary.Write(m.File, binary.NativeEndian, in)
}

// Close implements Close.
func (m *Port) Close() error {
	return m.File.Close()
}

// NewMemIOPort returns a Port, given an os.File.
func NewMemIOPort(f *os.File) *Port {
	return &Port{
		File: f,
	}
}
