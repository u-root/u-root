// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/u-root/u-root/pkg/dt"
)

// Properties to be fetched from device tree.
const (
	FirstLevelNodeName    = "images"
	SecondLevelNodeName   = "tianocore"
	LoadAddrPropertyName  = "load"
	EntryAddrPropertyName = "entry-start"
)

var sysfsCPUInfoPath = "/proc/cpuinfo"

type FdtLoad struct {
	Load       uint64
	EntryStart uint64
}

// Errors returned by utilities
var (
	ErrFailToReadFdtFile       = errors.New("failed to read fdt file")
	ErrNodeImagesNotFound      = fmt.Errorf("failed to find '%s' node", FirstLevelNodeName)
	ErrNodeTianocoreNotFound   = fmt.Errorf("failed to find '%s' node", SecondLevelNodeName)
	ErrNodeLoadNotFound        = fmt.Errorf("failed to find get '%s' property", LoadAddrPropertyName)
	ErrNodeEntryStartNotFound  = fmt.Errorf("failed to find get '%s' property", EntryAddrPropertyName)
	ErrFailToConvertLoad       = fmt.Errorf("failed to convert property '%s' to u64", LoadAddrPropertyName)
	ErrFailToConvertEntryStart = fmt.Errorf("failed to convert property '%s' to u64", EntryAddrPropertyName)
	ErrCPUAddressNotFound      = errors.New("'address sizes' information not found")
	ErrCPUAddressRead          = errors.New("failed to read 'address sizes'")
	ErrCPUAddressConvert       = errors.New("failed to convert physical bits size")
	ErrAlignPadRange           = errors.New("failed to align pad size, out of range")
)

// GetFdtInfo Device Tree Blob resides at the start of FIT binary. In order to
// get the expected load and entry point address, need to walk through
// DTB to get value of properties 'load' and 'entry-start'.
//
// The simplified device tree layout is:
//
//	/{
//	    images {
//	        tianocore {
//	            entry-start = <0x00000000 0x00805ac3>;
//	            load = <0x00000000 0x00800000>;
//	        }
//	    }
//	 }
func GetFdtInfo(name string) (*FdtLoad, error) {
	return getFdtInfo(name, nil)
}

func getFdtInfo(name string, dtb io.ReaderAt) (*FdtLoad, error) {
	var fdt *dt.FDT
	var err error

	if dtb != nil {
		fdt, err = dt.ReadFDT(io.NewSectionReader(dtb, 0, math.MaxInt64))
	} else {
		fdt, err = dt.ReadFile(name)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: fdt file: %s, err: %w", ErrFailToReadFdtFile, name, err)
	}

	firstLevelNode, succeed := fdt.NodeByName(FirstLevelNodeName)
	if succeed != true {
		return nil, ErrNodeImagesNotFound
	}

	secondLevelNode, succeed := firstLevelNode.NodeByName(SecondLevelNodeName)
	if succeed != true {
		return nil, ErrNodeTianocoreNotFound
	}

	loadAddrProp, succeed := secondLevelNode.LookProperty(LoadAddrPropertyName)
	if succeed != true {
		return nil, ErrNodeLoadNotFound
	}

	loadAddr, err := loadAddrProp.AsU64()
	if err != nil {
		return nil, errors.Join(ErrFailToConvertLoad, err)
	}

	entryAddrProp, succeed := secondLevelNode.LookProperty(EntryAddrPropertyName)
	if succeed != true {
		return nil, ErrNodeEntryStartNotFound
	}

	entryAddr, err := entryAddrProp.AsU64()
	if err != nil {
		return nil, errors.Join(ErrFailToConvertEntryStart, err)
	}

	return &FdtLoad{
		Load:       loadAddr,
		EntryStart: entryAddr,
	}, nil
}

// Get Physical Address size from sysfs node /proc/cpuinfo.
// Both Physical and Virtual Address size will be prompted as format:
// "address sizes	: 39 bits physical, 48 bits virtual"
// Use regular expression to fetch the integer of Physical Address
// size before "bits physical" keyword
func getPhysicalAddressSizes() (uint8, error) {
	file, err := os.Open(sysfsCPUInfoPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s: %w", sysfsCPUInfoPath, err)
	}
	defer file.Close()

	// Regular expression to match the address size line
	re := regexp.MustCompile(`address sizes\s*:\s*(\d+)\s+bits physical,\s*(\d+)\s+bits virtual`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match := re.FindStringSubmatch(line); match != nil {
			// Convert the physical bits size to integer
			physicalBits, err := strconv.ParseUint(match[1], 10, 8)
			if err != nil {
				return 0, errors.Join(ErrCPUAddressConvert, err)
			}
			return uint8(physicalBits), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("%w: file: %s, err: %w", ErrCPUAddressRead, sysfsCPUInfoPath, err)
	}

	return 0, ErrCPUAddressNotFound
}

// alignHOBLength writes pad bytes at the end of a HOB buf
// It's because we calculate HOB length with golang, while write bytes to the buf with actual length
func alignHOBLength(expectLen uint64, bufLen int, buf *bytes.Buffer) error {
	if expectLen < uint64(bufLen) {
		return ErrAlignPadRange
	}

	if expectLen > math.MaxInt {
		return ErrAlignPadRange
	}
	if padLen := int(expectLen) - bufLen; padLen > 0 {
		pad := make([]byte, padLen)
		if err := binary.Write(buf, binary.LittleEndian, pad); err != nil {
			return err
		}
	}
	return nil
}
