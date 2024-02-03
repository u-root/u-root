// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package spimock provides a fake SPI flash part for unit testing. This
// simulates an MX66L51235F.
package spimock

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/flash/chips"
	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/spidev"
	"golang.org/x/sys/unix"
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

// FakeSize is the size of the mocked flash chip.
const FakeSize = 64 * 1024 * 1024

// WriteWaitStates is the number WritePending is set to after a write.
const WriteWaitStates = 5

// MockSPI is an implementation of flash.SPI which records the transfers.
type MockSPI struct {
	// Data contains fake contents of the flash chip.
	Data []byte
	// isMmap is set to true if Data is memory mapped to a file.
	isMmap bool
	// SFDPData contains fake SFDP.
	SFDP []byte

	SpeedHz uint32

	// Is4BA is set to true if the chip is currently in 4-byte addressing
	// mode.
	Is4BA bool
	// IsWriteEnabled is set to true if data can be written.
	IsWriteEnabled bool
	// WritePending is non-zero while a write is pending. It decreases on
	// every read of the status register.
	WritePending int

	// Transfers is a recording of the transfers.
	Transfers []spidev.Transfer
	// ForceTransferError is returned by Transfer when set.
	ForceTransferErr error
	// ForceSetSpeedHzError is returned by SetSpeedHz when set.
	ForceSetSpeedHzErr error
}

// New returns a new MockSPI in memory.
func New() *MockSPI {
	return &MockSPI{
		Data: make([]byte, FakeSize),
		SFDP: FakeSFDP,
	}
}

// NewFromFile returns a new MockSPI which is backed by a file. If the file
// does not exist, it will be created. Ideally, the file's size should match
// FlashSize.
func NewFromFile(filename string) (*MockSPI, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := unix.Mmap(int(f.Fd()), 0, FakeSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return &MockSPI{
		Data:   data,
		isMmap: true,
		SFDP:   FakeSFDP,
	}, nil
}

// Close closes the mock.
func (s *MockSPI) Close() error {
	if s.isMmap {
		return unix.Munmap(s.Data)
	}
	return nil
}

// tlen returns the length of a single transfers.
func tlen(t *spidev.Transfer) int {
	if len(t.Tx) > len(t.Rx) {
		return len(t.Tx)
	}
	return len(t.Rx)
}

// tx reads byte n from the transfers.
func tx(transfers []spidev.Transfer, n int) (byte, error) {
	for i := range transfers {
		l := tlen(&transfers[i])
		if n >= l {
			n -= l
			continue
		}
		if n >= len(transfers[i].Tx) {
			return 0, nil
		}
		return transfers[i].Tx[n], nil
	}
	return 0, io.EOF
}

// rx writes byte n from the transfers.
func rx(transfers []spidev.Transfer, n int, val byte) error {
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

// address returns an address from the given tx offset in transfers and given
// addressing mode. The second return value is the 3 or 4 for the addressing mode.
func address(transfers []spidev.Transfer, off int, is4BA bool) (int64, int64) {
	tx0 := func(n int) int64 {
		b, _ := tx(transfers, n)
		return int64(b)
	}

	if is4BA {
		// Big-endian
		return (tx0(off) << 24) | (tx0(off+1) << 16) | (tx0(off+2) << 8) | tx0(off+3), 4
	}
	// Big-endian
	return (tx0(off) << 16) | (tx0(off+1) << 8) | tx0(off+2), 3
}

// Transfer implements flash.SPI.
func (s *MockSPI) Transfer(transfers []spidev.Transfer) error {
	if s.ForceTransferErr != nil {
		return s.ForceTransferErr
	}

	s.Transfers = append(s.Transfers, transfers...)

	o, err := tx(transfers, 0)
	if err != nil {
		return err
	}
	switch op.OpCode(o) {
	case op.PageProgram:
		if !s.IsWriteEnabled {
			break
		}
		addr, addrLen := address(transfers, 1, s.Is4BA)
		// Copy each byte from tx to data with wrap-around within the page.
		for i := 0; ; i++ {
			b, err := tx(transfers, 1+int(addrLen)+i)
			if err == io.EOF {
				break
			}
			s.Data[addr&^255|(addr+int64(i))&255] &= b
		}
		s.IsWriteEnabled = false
		s.WritePending = WriteWaitStates
	case op.Read:
		addr, addrLen := address(transfers, 1, s.Is4BA)
		// Copy each byte from data to rx.
		for i := int64(0); rx(transfers, int(1+addrLen+i), s.Data[addr+i]) == nil; i++ {
		}
	case op.WriteDisable:
		s.IsWriteEnabled = false
	case op.ReadStatus:
		var statusReg uint8
		if s.WritePending != 0 {
			statusReg |= 1
		}
		if s.IsWriteEnabled {
			statusReg |= 2
		}
		rx(transfers, 1, statusReg)
	case op.WriteEnable:
		s.IsWriteEnabled = true
	case op.SectorErase:
		if !s.IsWriteEnabled {
			break
		}
		addr, _ := address(transfers, 1, s.Is4BA)
		addr &= ^0xfff
		copy(s.Data[addr:], bytes.Repeat([]byte{0xff}, 0x1000))
		s.IsWriteEnabled = false
	case op.ReadSFDP:
		// ReadSFDP is always 3-byte addressing.
		addr, addrLen := address(transfers, 1, false)
		// Copy each byte from sfdp to rx.
		for i := int64(0); rx(transfers, int(1+addrLen+1+i), s.SFDP[addr+i]) == nil; i++ {
		}
	case op.ReadJEDECID:
		rx(transfers, 1, 0xc2)
		rx(transfers, 2, 0x20)
		rx(transfers, 3, 0x1a)
	case op.Enter4BA:
		s.Is4BA = true
	case op.BlockErase:
		if !s.IsWriteEnabled {
			break
		}
		addr, _ := address(transfers, 1, s.Is4BA)
		addr &= ^0x7fff
		copy(s.Data[addr:], bytes.Repeat([]byte{0xff}, 0x8000))
		s.IsWriteEnabled = false
	case op.PRDRES:
	case op.Exit4BA:
		s.Is4BA = false
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

func (s *MockSPI) ID() (chips.ID, error) {
	if s.ForceTransferErr != nil {
		return -1, s.ForceTransferErr
	}
	return 0xbf2541, nil
}

func (s *MockSPI) Status() (op.Status, error) {
	if s.ForceTransferErr != nil {
		return op.Status(0xff), s.ForceTransferErr
	}
	return 0, nil
}
