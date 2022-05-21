// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"os"
)

// Dev contains information about ongoing MTD status and operation.
type Dev struct {
	*os.File
	devName string
}

// DevName is the default name for the MTD device.
var DevName = "/dev/mtd0"

// NewDev creates a Dev, returning Flasher or error.
func NewDev(n string) (Flasher, error) {
	f, err := os.OpenFile(n, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &Dev{File: f, devName: n}, nil
}

// QueueWrite adds a []byte to the pending write queue.
func (m *Dev) QueueWriteAt(b []byte, off int64) (int, error) {
	return m.File.WriteAt(b, off)
}

// SyncWrite syncs a pending queue of writes to a device.
func (m *Dev) SyncWrite() error {
	return nil
}

// ReadAt implements io.ReadAT
func (m *Dev) ReadAt(b []byte, off int64) (int, error) {
	return m.File.ReadAt(b, off)
}

// Close implements io.Close
func (m *Dev) Close() error {
	return m.File.Close()
}

// DevName returns the name of the flash device.
func (m *Dev) Name() string {
	return m.devName
}
