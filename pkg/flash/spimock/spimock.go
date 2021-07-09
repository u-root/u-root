// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package spimock provides a fake SPI flash part for unit testing. This
// simulates an MX66L51235F.
package spimock

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/spi"
)

// FakeSFDP is the SFDP from the MX66L51235F datasheet.
var FakeSFDP = []byte{
	// 0x00: SFDP Table
	0x53, 0x46, 0x44, 0x50,
	0x00, 0x01, 0x01, 0xff,
	0x00, 0x00, 0x01, 0x09,
	0x30, 0x00, 0x00, 0xff,
	0xc2, 0x00, 0x01, 0x04,
	0x60, 0x00, 0x00, 0xff,
	// Padding
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,

	// 0x30: Table 0
	0xe5, 0x20, 0xf3, 0xff,
	0xff, 0xff, 0xff, 0x1f,
	0x44, 0xeb, 0x08, 0x6b,
	0x08, 0x3b, 0x04, 0xbb,
	0xfe, 0xff, 0xff, 0xff,
	0xff, 0xff, 0x00, 0xff,
	0xff, 0xff, 0x44, 0xeb,
	0x0c, 0x20, 0x0f, 0x52,
	0x10, 0xd8, 0x00, 0xff,
	// Padding
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,

	// 0x60: Table 1
	0x00, 0x36, 0x00, 0x27,
	0x9d, 0xf9, 0xc0, 0x64,
	0x85, 0xcb, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
}

// MockSPI is an implementation of flash.SPI which records the transfers.
type MockSPI struct {
	// Data contains fake contents of the flash chip.
	Data []byte
	// SFDPData contains fake SFDP.
	SFDP []byte

	SpeedHz uint32

	// Transfers is a recording of the transfers.
	Transfers []spi.Transfer
	// ForceTransferError is returned by Transfer when set.
	ForceTransferErr error
	// ForceSetSpeedHzError is returned by SetSpeedHz when set.
	ForceSetSpeedHzErr error
}

// New returns a new MockSPI.
func New() *MockSPI {
	return &MockSPI{
		Data: make([]byte, 64*1024*1024),
		SFDP: FakeSFDP,
	}
}

// Close closes the mock. There is nothing for this to actually do.
func (s *MockSPI) Close() error {
	return nil
}

// tlen returns the length of a single transfers.
func tlen(t *spi.Transfer) int {
	if len(t.Tx) > len(t.Rx) {
		return len(t.Tx)
	}
	return len(t.Rx)
}

// tx reads byte n from the transfers.
func tx(transfers []spi.Transfer, n int) byte {
	for i := range transfers {
		l := tlen(&transfers[i])
		if n >= l {
			n -= l
			continue
		}
		if n >= len(transfers[i].Tx) {
			return 0xff
		}
		return transfers[i].Tx[n]
	}
	return 0xff
}

// rx writes byte n from the transfers.
func rx(transfers []spi.Transfer, n int, val byte) error {
	for i := range transfers {
		l := tlen(&transfers[i])
		if n >= l {
			n -= l
			continue
		}
		if n >= len(transfers[i].Rx) {
			return io.EOF
		}
		transfers[i].Rx[n] = val
		return nil
	}
	return io.EOF
}

// Transfer implements flash.SPI.
func (s *MockSPI) Transfer(transfers []spi.Transfer) error {
	if s.ForceTransferErr != nil {
		return s.ForceTransferErr
	}

	s.Transfers = append(s.Transfers, transfers...)

	o := op.Op(tx(transfers, 0))
	switch o {
	case op.ReadSFDP:
		// Big-endian
		off := (int(tx(transfers, 1)) << 16) | (int(tx(transfers, 2)) << 8) | int(tx(transfers, 3))
		// Copy each byte to rx.
		for i := 0; rx(transfers, i+5, s.SFDP[off+i]) == nil; i++ {
		}
	default:
		return fmt.Errorf("unrecognized opcode %#02x", o)
	}

	return nil
}

// SetSpeedHz sets the SPI speed. The value set is recorded in the mock.
func (s *MockSPI) SetSpeedHz(hz uint32) error {
	if s.ForceSetSpeedHzErr != nil {
		return s.ForceSetSpeedHzErr
	}
	s.SpeedHz = hz
	return nil
}
