// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"fmt"
	"os"
)

// Dev is a Linux device to be used for operations on MTDs.
type Dev struct {
	*os.File
}

// NewFlasher returns a Flasher or an error
func NewFlasher(n string) (Flasher, error) {
	return nil, fmt.Errorf("not yet")
}

// QueueWrite adds a []byte to the pending write queue.
func (m *Dev) QueueWrite(b []byte, off int64) (int, error) {
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
func (m *Dev) DevName() string {
	return m.File.Name()
}
