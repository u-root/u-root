// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flash provides higher-level functions such as reading, erasing,
// writing and programming the flash chip.
package flash

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

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
	ID() (chips.ID, error)
}

// Flash provides operations for SPI flash chips.
type Flash struct {
	// spi is the underlying SPI device.
	spi SPI

	// Chip is derived from SFDP or looking up
	// the chip via the ID.
	chips.Chip

	// sfdp is cached.
	sfdp *sfdp.SFDP
}

// New creates a new flash device from a SPI interface.
func New(spi SPI) (*Flash, error) {
	f := &Flash{
		spi: spi,
	}

	var err error
	var id chips.ID
	id, err = f.spi.ID()
	if err != nil {
		return nil, fmt.Errorf("can not ID chip: %w", err)
	}
	c, err := chips.Lookup(id)
	if err == nil {
		f.Chip = *c
		// Even when we have a chip, there is still
		// benefit to trying to get an SFDP.
		// Further, the package as written wants
		// some sort of empty sfdp to exist, and this
		// is the easiest way to do it.
		f.sfdp, _ = sfdp.Read(f.SFDPReader())
		return f, nil
	}

	f.sfdp, err = sfdp.Read(f.SFDPReader())
	if err != nil {
		return nil, fmt.Errorf("chip %#x: chip not known, and no SFDP: %w", id, err)
	}
	if err = f.FillFromSFDP(); err != nil {
		return nil, fmt.Errorf("chip %#x: chip not known, and no SFDP: %w", id, err)
	}

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
	f.ArraySize = (density + 1) / 8

	// Assume 4ba if address if size requires 4 bytes.
	if f.ArraySize >= 0x1000000 {
		f.Is4BA = true
	}

	// TODO
	f.PageSize = 256
	f.SectorSize = 4096
	f.BlockSize = 65536

	return nil
}

// Size returns the size of the flash chip in bytes.
func (f *Flash) Size() int64 {
	return f.ArraySize
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
	if f.Is4BA {
		return data
	}
	return data[1:]
}

// ReadAt reads from the flash chip.
func (f *Flash) ReadAt(p []byte, off int64) (int, error) {
	// This is a valid implementation of io.ReaderAt.
	if off < 0 || off > f.ArraySize {
		return 0, io.EOF
	}
	p = p[:min(int64(len(p)), f.ArraySize-off)]

	// Split the transfer into maxTransferSize chunks.
	for i := 0; i < len(p); i += maxTransferSize {
		if err := f.spi.Transfer([]spidev.Transfer{
			{Tx: append(op.Read.Bytes(), f.prepareAddress(off+int64(i))...)},
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
		{Tx: op.PRDRES.Bytes(), CSChange: true},
		{Tx: op.WriteEnable.Bytes(), CSChange: true},
		// Send the address.
		{Tx: append(op.PageProgram.Bytes(), f.prepareAddress(off)...)},
		// Send the data.
		{Tx: p},
	}); err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Microsecond)
	if err := f.spi.Transfer([]spidev.Transfer{
		// The meaning of CSChange is ... odd.
		// IF CSChange is set true here, then CE# never goes
		// high. If CSChange is left unchanged,
		// CE# is properly deasserted from the data write above,
		// asserted for this command, and deasserted
		// at the end.
		{Tx: op.WriteDisable.Bytes()},
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
	if off < 0 || off > f.ArraySize {
		return 0, io.EOF
	}
	p = p[:min(int64(len(p)), f.ArraySize-off)]

	// Special case where no page boundaries are crossed.
	if off%f.PageSize+int64(len(p)) <= f.PageSize {
		return f.writeAt(p, off)
	}

	// Otherwise, there are three regions:
	// 1. A partial page before the first aligned offset. (optional)
	// 2. All the aligned pages in the middle.
	// 3. A partial page after the last aligned offset. (optional)
	firstAlignedOff := (off + f.PageSize - 1) / f.PageSize * f.PageSize
	lastAlignedOff := (off + int64(len(p))) / f.PageSize * f.PageSize

	if off != firstAlignedOff {
		if n, err := f.writeAt(p[:firstAlignedOff-off], off); err != nil {
			return n, err
		}
	}
	for i := firstAlignedOff; i < lastAlignedOff; i += f.PageSize {
		if _, err := f.writeAt(p[i:i+f.PageSize], off+i); err != nil {
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
	if off < 0 || off > f.ArraySize || off+n > f.ArraySize {
		return 0, io.EOF
	}

	if (off%f.SectorSize != 0) || (n%f.SectorSize != 0) {
		return 0, fmt.Errorf("len(p) and off must be multiple of the sector size")
	}

	for i := int64(0); i < n; {
		opcode := op.SectorErase
		eraseSize := f.SectorSize

		// Optimization to erase faster.
		if i%f.BlockSize == 0 && n-i > f.BlockSize {
			opcode = op.BlockErase
			eraseSize = f.BlockSize
		}

		if err := f.spi.Transfer([]spidev.Transfer{
			// Enable writing.
			{
				Tx:       op.WriteEnable.Bytes(),
				CSChange: true,
			},
			// Send the address.
			{Tx: append(opcode.Bytes(), f.prepareAddress(off+i)...)},
		}); err != nil {
			return i, err
		}

		i += eraseSize
	}
	return n, nil
}

// ReadJEDECID reads the flash chip's JEDEC ID.
func (f *Flash) ReadJEDECID() (uint32, error) {
	id, err := f.spi.ID()
	return uint32(id), err
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
