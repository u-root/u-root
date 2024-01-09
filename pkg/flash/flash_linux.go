// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flash provides higher-level functions such as reading, erasing,
// writing and programming the flash chip.
package flash

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/flash/chips"
	"github.com/u-root/u-root/pkg/flash/op"
	"github.com/u-root/u-root/pkg/flash/sfdp"
	"github.com/u-root/u-root/pkg/spidev"
)

// sfdpMaxAddress is the highest possible SFDP address (24 bit address space).
const sfdpMaxAddress = (1 << 24) - 1

// SPI interface for the underlying calls to the SPI driver.
type SPI interface {
	Transfer(transfers []spidev.Transfer) error
	ID() (int, error)
}

// Flash provides operations for SPI flash chips.
type Flash struct {
	// spi is the underlying SPI device.
	spi SPI

	chip *chips.Chip

	// these are filled in from SFDP or the chip.
	// is4ba is true if 4-byte addressing mode is enabled.
	is4ba bool

	// size is the size of the flash chip in bytes.
	size int64

	// sfdp is cached.
	sfdp *sfdp.SFDP

	// JEDEC ID is cached.
	id uint32

	pageSize   int64
	sectorSize int64
	blockSize  int64
}

// New creates a new flash device from a SPI interface.
func New(spi SPI) (*Flash, error) {
	f := &Flash{
		spi: spi,
	}

	var err error
	f.sfdp, err = sfdp.Read(f.SFDPReader())
	if err == nil {
		if err = f.FillFromSFDP(); err == nil {
			return f, nil
		}
	}

	id, e := f.spi.ID()
	if e != nil {
		return nil, errors.Join(err, e)
	}

	if f.chip, e = chips.New(id); e != nil {
		return nil, errors.Join(err, e)
	}

	f.is4ba = f.chip.Is4BA
	f.size = f.chip.Size
	f.pageSize = f.chip.PageSize
	f.sectorSize = f.chip.SectorSize
	f.blockSize = f.chip.BlockSize
	return f, nil
}

func (f *Flash) FillFromSFDP() error {
	density, err := f.SFDP().Param(sfdp.ParamFlashMemoryDensity)
	if err != nil {
		return fmt.Errorf("flash chip SFDP does not have density param")
	}
	if density >= 0x80000000 {
		return fmt.Errorf("unsupported flash density: %#x", density)
	}
	f.size = (density + 1) / 8

	// Assume 4ba if address if size requires 4 bytes.
	if f.size >= 0x1000000 {
		f.is4ba = true
	}

	// TODO
	f.pageSize = 256
	f.sectorSize = 4096
	f.blockSize = 65536

	return nil
}

// Size returns the size of the flash chip in bytes.
func (f *Flash) Size() int64 {
	return f.size
}

const maxTransferSize = 4096

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// prepareAddress converts an address to the 3- or 4-byte addressing mode.
func (f *Flash) prepareAddress(addr int64) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(addr))
	if f.is4ba {
		return data
	}
	return data[1:]
}

// ReadAt reads from the flash chip.
func (f *Flash) ReadAt(p []byte, off int64) (int, error) {
	// This is a valid implementation of io.ReaderAt.
	if off < 0 || off > f.size {
		return 0, io.EOF
	}
	p = p[:min(int64(len(p)), f.size-off)]

	// Split the transfer into maxTransferSize chunks.
	for i := 0; i < len(p); i += maxTransferSize {
		if err := f.spi.Transfer([]spidev.Transfer{
			{Tx: append([]byte{byte(op.Read)}, f.prepareAddress(off+int64(i))...)},
			{Rx: p[i:min(int64(i)+maxTransferSize, int64(len(p)))]},
		}); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// writeAt performs a write operation without any care for page sizes or
// alignment.
func (f *Flash) writeAt(p []byte, off int64) (int, error) {
	if err := f.spi.Transfer([]spidev.Transfer{
		// Enable writing.
		{Tx: []byte{byte(op.PRDRES)}, CSChange: true},
		{Tx: []byte{byte(op.WriteEnable)}, CSChange: true},
		// Send the address.
		{Tx: append([]byte{byte(op.PageProgram)}, f.prepareAddress(off)...)},
		// Send the data.
		{Tx: p, CSChange: true},
		{Tx: []byte{byte(op.WriteDisable)}, CSChange: true},
	}); err != nil {
		return 0, err
	}
	return len(p), nil
}

// WriteAt writes to the flash chip.
//
// For optimal performance, call this function on boundaries of pageSize.
//
// NOTE: This will not erase before writing! The ProgramAt function is probably
// what you want instead!
func (f *Flash) WriteAt(p []byte, off int64) (int, error) {
	// This is a valid implementation of io.WriterAt.
	if off < 0 || off > f.size {
		return 0, io.EOF
	}
	p = p[:min(int64(len(p)), f.size-off)]

	// Special case where no page boundaries are crossed.
	if off%f.pageSize+int64(len(p)) <= f.pageSize {
		return f.writeAt(p, off)
	}

	// Otherwise, there are three regions:
	// 1. A partial page before the first aligned offset. (optional)
	// 2. All the aligned pages in the middle.
	// 3. A partial page after the last aligned offset. (optional)
	firstAlignedOff := (off + f.pageSize - 1) / f.pageSize * f.pageSize
	lastAlignedOff := (off + int64(len(p))) / f.pageSize * f.pageSize

	if off != firstAlignedOff {
		if n, err := f.writeAt(p[:firstAlignedOff-off], off); err != nil {
			return n, err
		}
	}
	for i := firstAlignedOff; i < lastAlignedOff; i += f.pageSize {
		if _, err := f.writeAt(p[i:i+f.pageSize], off+i); err != nil {
			return int(i), err
		}
	}
	if off+int64(len(p)) != lastAlignedOff {
		if _, err := f.writeAt(p[lastAlignedOff-off:], lastAlignedOff); err != nil {
			return int(lastAlignedOff - off), err
		}
	}
	return len(p), nil
}

func (f *Flash) ProgramAt(p []byte, off int64) (int, error) {
	// TODO
	return 0, nil
}

// EraseAt erases n bytes from offset off. Both parameters must be aligned to
// sectorSize.
func (f *Flash) EraseAt(n int64, off int64) (int64, error) {
	if off < 0 || off > f.size || off+n > f.size {
		return 0, io.EOF
	}

	if (off%f.sectorSize != 0) || (n%f.sectorSize != 0) {
		return 0, fmt.Errorf("len(p) and off must be multiple of the sector size")
	}

	for i := int64(0); i < n; {
		opcode := op.SectorErase
		eraseSize := f.sectorSize

		// Optimization to erase faster.
		if i%f.blockSize == 0 && n-i > f.blockSize {
			opcode = op.BlockErase
			eraseSize = f.blockSize
		}

		if err := f.spi.Transfer([]spidev.Transfer{
			// Enable writing.
			{
				Tx:       []byte{byte(op.WriteEnable)},
				CSChange: true,
			},
			// Send the address.
			{Tx: append([]byte{byte(opcode)}, f.prepareAddress(off+i)...)},
		}); err != nil {
			return i, err
		}

		i += eraseSize
	}
	return n, nil
}

// ReadJEDECID reads the flash chip's JEDEC ID.
func (f *Flash) ReadJEDECID() (uint32, error) {
	tx := []byte{byte(op.ReadJEDECID)}
	rx := make([]byte, 3)

	if err := f.spi.Transfer([]spidev.Transfer{
		{Tx: tx},
		{Rx: rx},
	}); err != nil {
		return 0, err
	}

	// Little-endian
	return (uint32(rx[0]) << 16) | (uint32(rx[1]) << 8) | uint32(rx[2]), nil
}

// SFDPReader is used to read from the SFDP address space.
func (f *Flash) SFDPReader() *SFDPReader {
	return (*SFDPReader)(f)
}

// SFDP returns all the SFDP tables from the flash chip. The value is cached.
func (f *Flash) SFDP() *sfdp.SFDP {
	return f.sfdp
}

// SFDPReader is a wrapper around Flash where the ReadAt function reads from
// the SFDP address space.
type SFDPReader Flash

// ReadAt reads from the given offset in the SFDP address space.
func (f *SFDPReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 || off > sfdpMaxAddress {
		return 0, io.EOF
	}
	p = p[:min(int64(len(p)), sfdpMaxAddress-off)]
	tx := []byte{
		byte(op.ReadSFDP),
		// offset, 3-bytes, big-endian
		byte((off >> 16) & 0xff), byte((off >> 8) & 0xff), byte(off & 0xff),
		// dummy 0xff
		0xff,
	}
	if err := f.spi.Transfer([]spidev.Transfer{
		{Tx: tx},
		{Rx: p},
	}); err != nil {
		return 0, err
	}
	return len(p), nil
}
