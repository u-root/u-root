// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/boot/ebda"
	"github.com/u-root/u-root/pkg/memio"
)

// " RSD PTR" in hex, 8 bytes.
const rsdpTag = 0x2052545020445352

// GetRSDPEBDA finds the RSDP in the EBDA.
func GetRSDPEBDA() (*RSDP, error) {
	f, err := os.OpenFile("/dev/mem", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	e, err := ebda.ReadEBDA(f)
	if err != nil {
		return nil, err
	}

	return getRSDPMem(int64(e.BaseOffset), int64(e.BaseOffset+e.Length))
}

func getRSDPMem(start, end int64) (*RSDP, error) {
	for base := start; base < end; base += 16 {
		var r memio.Uint64
		if err := memio.Read(int64(base), &r); err != nil {
			continue
		}
		if r != rsdpTag {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			return nil, err
		}
		return rsdp, nil
	}
	return nil, fmt.Errorf("could not find ACPI RSDP via /dev/mem from %#08x to %#08x", start, end)
}

// GetRSDPMem is the option of last choice, it just grovels through
// the e0000-ffff0 area, 16 bytes at a time, trying to find an RSDP.
// These are well-known addresses for 20+ years.
func GetRSDPMem() (*RSDP, error) {
	return getRSDPMem(0xe0000, 0xffff0)
}
