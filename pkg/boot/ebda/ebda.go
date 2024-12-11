// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ebda looks for the Extended Bios Data Area (EBDA) pointer in /dev/mem,
// and provides access to the EBDA. This is useful for us to read or write to the EBDA,
// for example to copy the RSDP into the EBDA.
//   - The address 0x40E contains the pointer to the start of the EBDA, shifted right by 4 bits
//   - We take that and find the EBDA, where the first byte usually encodes the length of the
//     the area in KiB.
//   - If the pointer is not set, there may be no EBDA.
package ebda

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// EBDAAddressOffset is the traditional offset where the location of the EBDA is stored.
const EBDAAddressOffset = 0x40E

// EBDA represents an EBDA region in memory
type EBDA struct {
	BaseOffset int64
	Length     int64
	Data       []byte
}

func findEBDAOffset(f io.ReadSeeker) (int64, error) {
	var ep uint16
	if _, err := f.Seek(EBDAAddressOffset, io.SeekStart); err != nil {
		return 0, fmt.Errorf("error seeking to memory offset %#X for the EBDA pointer, got: %w", EBDAAddressOffset, err)
	}

	// read the ebda pointer
	if err := binary.Read(f, binary.LittleEndian, &ep); err != nil {
		return 0, fmt.Errorf("unable to read EBDA Pointer: %w", err)
	}

	if ep == 0 {
		return 0, errors.New("ebda offset is 0! unable to proceed")
	}
	// The EBDA pointer is stored shifted 4 bits right.
	// We've never seen one that's not that way, though osdev implies there are some.
	// For reference: https://wiki.osdev.org/Memory_Map_(x86)
	return int64(ep) << 4, nil
}

// ReadEBDA assumes it's been passed in /dev/mem and searches for and reads the EBDA.
func ReadEBDA(f io.ReadSeeker) (*EBDA, error) {
	var err error
	var eSize uint8
	e := &EBDA{}

	if e.BaseOffset, err = findEBDAOffset(f); err != nil {
		return nil, err
	}

	if _, err = f.Seek(e.BaseOffset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("error seeking to memory offset %#X for the start of the EBDA, got: %w", e.BaseOffset, err)
	}

	// Read length
	if err := binary.Read(f, binary.LittleEndian, &eSize); err != nil {
		return nil, fmt.Errorf("error reading EBDA length, got: %w", err)
	}
	e.Length = int64(eSize) << 10

	// Find the size of this segment if the first byte is not set.
	// We assume that the EBDA never crosses a segment for now.
	// This is a reasonable assumption in most cases, since most EBDAs I've seen are in the 0x9XXXX range,
	// and the 0xA0000 segment is off limits.
	if e.Length == 0 {
		e.Length = 0x10000 - (e.BaseOffset & 0xffff)
	}
	e.Data = make([]byte, e.Length)

	if _, err = f.Seek(e.BaseOffset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("error seeking to memory offset %#X for the start of the EBDA, got: %w", e.BaseOffset, err)
	}
	if err = binary.Read(f, binary.LittleEndian, e.Data); err != nil {
		return nil, fmt.Errorf("error reading EBDA region, tried to read from %#X of size %#X, got %w", e.BaseOffset, e.Length, err)
	}

	return e, nil
}

// WriteEBDA assumes it's been passed in /dev/mem and searches for and writes to the EBDA.
func WriteEBDA(e *EBDA, f io.ReadWriteSeeker) error {
	var err error
	if e.BaseOffset == 0 {
		// Don't know why this is 0, but try to find one from /dev/mem
		if e.BaseOffset, err = findEBDAOffset(f); err != nil {
			return err
		}
	}

	if e.Length&0x3ff != 0 {
		// Length can't be encoded as KiB.
		return fmt.Errorf("length is not an integer multiple of 1 KiB, got %#X", e.Length)
	}
	if e.Length != int64(len(e.Data)) {
		return fmt.Errorf("length field is not equal to buffer length. length field: %#X, buffer length: %#X", e.Length, len(e.Data))
	}

	// Update length field in buffer
	eLen := e.Length >> 10
	if eLen > 0xff {
		// Length is too big to be written into one byte. fail.
		return fmt.Errorf("length is greater than 255 KiB, cannot be marshalled into one byte. Length: %dKiB", eLen)
	}

	// Cheat since it's only one byte
	e.Data[0] = byte(eLen)

	if _, err = f.Seek(e.BaseOffset, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to memory offset %#X for the start of the EBDA, got: %w", e.BaseOffset, err)
	}

	return binary.Write(f, binary.LittleEndian, e.Data)
}
