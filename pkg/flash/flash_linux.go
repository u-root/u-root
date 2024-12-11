// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flash provides higher-level functions such as reading, erasing,
// writing and programming the flash chip.
package flash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
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
	Status() (op.Status, error)
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

// FillFromSFDP fills the Flash struct with parameters from
// the SFDP. Querying the SFDP is a bit messy, for each type of
// parameter, so this code pulls the SFDP parameters into
// struct members. It also makes the creation of chips a bit easier,
// when we do not have an SFDP.
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

	if wer, err := f.SFDP().Param(sfdp.ParamWriteEnableInstructionRequired); err == nil && wer != 0 {
		we, err := f.SFDP().Param(sfdp.ParamWriteEnableOpcodeSelect)
		if err != nil {
			return fmt.Errorf("write enable is required but WriteEnableOpcodeSelect is not in SFDP:%w", err)
		}
		f.WriteEnableInstructionRequired = true
		f.WriteEnableOpcodeSelect = op.OpCode(we)
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
	t := []spidev.Transfer{
		{Tx: op.PRDRES.Bytes(), CSChange: true},
	}
	if f.Chip.WriteEnableInstructionRequired {
		t = append(t, spidev.Transfer{Tx: f.Chip.WriteEnableOpcodeSelect.Bytes(), CSChange: true})
	}
	// AAAAAAARRRRGHHHHH!
	// If this WriteEnable is not here, then the page program request
	// NEVER MAKES IT TO THE SPI BUS.
	// So, ... put it here, even if not requested, until we figure this out.
	// Further, on the macronix part, we can't leave the write disable in, or
	// the write enable command ends up being written to the part?
	// This is a mess.
	t = append(t, spidev.Transfer{Tx: op.WriteEnable.Bytes(), CSChange: true})
	t = append(t, spidev.Transfer{Tx: append(append(op.PageProgram.Bytes(), f.prepareAddress(off)...), p...)})
	if f.Chip.WriteEnableInstructionRequired {
		// The meaning of CSChange is ... odd.
		// IF CSChange is set true here, then CE# never goes
		// high. If CSChange is left unchanged,
		// CE# is properly deasserted from the data write above,
		// asserted for this command, and deasserted
		// at the end.
		t = append(t, spidev.Transfer{Tx: op.WriteDisable.Bytes(), DelayUSecs: 10})
	}
	if err := f.spi.Transfer(t); err != nil {
		return 0, err
	}
	// Hang out for a bit, let the part do its thing.
	time.Sleep(time.Duration(len(p)) * time.Microsecond)
	var i int
	for i = 0; i < len(p); i++ {
		stat, err := f.spi.Status()
		if err != nil {
			return len(p), fmt.Errorf("spi status read fails after writing %d bytes:%w", len(p), err)
		}
		if !stat.Busy() {
			break
		}
		time.Sleep(10 * time.Microsecond)
	}

	if i == len(p) {
		return len(p), fmt.Errorf("spi busy after writing %d bytes", len(p))
	}
	// Between each check for done, sleep about 10 microseconds, to give the part
	// a chance to catch its breath.
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
	// 2. All the aligned pages
	firstAlignedOff := (off + f.PageSize - 1) / f.PageSize * f.PageSize

	b := bytes.NewBuffer(p)
	if off != firstAlignedOff {
		dat := b.Next(int(firstAlignedOff - off))
		if n, err := f.writeAt(dat, off); err != nil {
			return n, err
		}
	}
	for i := firstAlignedOff; b.Len() > 0; i += f.PageSize {
		dat := b.Next(int(f.PageSize))
		if _, err := f.writeAt(dat, i); err != nil {
			return int(i), err
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
		return 0, fmt.Errorf("offset (%#x) is < 0, or > %#x, or off+size (%#x) is > f.ArraySize (%#x):%w", off, f.ArraySize, off+n, f.ArraySize, os.ErrInvalid)
	}

	if (off%f.SectorSize != 0) || (n%f.SectorSize != 0) {
		return 0, fmt.Errorf("offset (%#x) and size (%#x) must be multiple of the sector size(%#x):%w", off, n, f.SectorSize, os.ErrInvalid)
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

		var spin int
		// Give it 100 tries, which is 1000 ms
		for spin = 0; spin <= 100; spin++ {
			time.Sleep(10 * time.Millisecond)
			stat, err := f.spi.Status()
			if err != nil {
				return n, fmt.Errorf("spi status read fails after erasing %#x:%w", off+i, err)
			}
			if !stat.Busy() {
				break
			}
		}

		if spin > 100 {
			return i, fmt.Errorf("spi busy after erasing %d bytes", i)
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
