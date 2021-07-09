// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flash provides higher-level functions such as reading, erasing,
// writing and programming the flash chip.
//
// TODO: Implement Read function.
// TODO: Implement Erase function.
// TODO: Implement Write function.
// TODO: Implement Program function.
package flash

import (
	"fmt"

	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/flash/sfdp"
	"github.com/u-root/u-root/pkg/spi"
)

// sfdpMaxAddress is the highest possible SFDP address (24 bit address space).
const sfdpMaxAddress = (1 << 24) - 1

// SFDPAddressError is returned if the SFDP address is out of range.
type SFDPAddressError struct {
	Offset int64
}

// Error implements error.
func (e *SFDPAddressError) Error() string {
	return fmt.Sprintf("SFDP address is invalid, %#x is outside range [0, %#x]", e.Offset, sfdpMaxAddress)
}

// SPI interface for the underlying calls to the SPI driver.
type SPI interface {
	Transfer(transfers []spi.Transfer) error
}

// Flash provides operations for SPI flash chips.
type Flash struct {
	// SPI is the underlying SPI device.
	SPI SPI

	// sfdp is cached along with any errors during reading.
	sfdp    *sfdp.SFDP
	sfdpErr error
}

// SFDPReader is used to read from the SFDP address space.
func (f *Flash) SFDPReader() *SFDPReader {
	return (*SFDPReader)(f)
}

// SFDP reads and returns all the SFDP tables from the flash chip. The value is
// cached.
func (f *Flash) SFDP() (*sfdp.SFDP, error) {
	if f.sfdpErr != nil {
		return nil, f.sfdpErr
	}
	if f.sfdp == nil {
		f.sfdp, f.sfdpErr = sfdp.Read(f.SFDPReader())
	}
	return f.sfdp, f.sfdpErr
}

// SFDPReader is a wrapper around Flash where the ReadAt function reads from
// the SFDP address space.
type SFDPReader Flash

// ReadAt reads from the given offset in the SFDP address space.
func (f *SFDPReader) ReadAt(rx []byte, off int64) (int, error) {
	if off < 0 || off > sfdpMaxAddress {
		return 0, &SFDPAddressError{Offset: off}
	}
	tx := []byte{
		byte(op.ReadSFDP),
		// offset, 3-bytes, big-endian
		byte((off >> 16) & 0xff), byte((off >> 8) & 0xff), byte(off & 0xff),
		// dummy 0xff
		0xff,
	}
	if err := f.SPI.Transfer([]spi.Transfer{
		{Tx: tx},
		{Rx: rx},
	}); err != nil {
		return 0, err
	}
	return len(rx), nil
}
