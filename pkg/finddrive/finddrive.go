// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package finddrive provides functionality to find an NVMe block device associated with
// a particular physical slot on the machine, based on information in the SMBIOS table.
package finddrive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/smbios"
)

const (
	// M2MKeySlotType is the SMBIOS slot type code for M.2 M-key
	M2MKeySlotType = 0x17
)

// Devices connected to a slot with domain number AAAA, bus number BB, device number CC, and
// function number D are represented as subdirectories of /sys/devices/pciAAAA:BB/AAAA:BB:CC.D/
// These paths are accessible via symlinks in /sys/block for NVMe devices, which can be
// inspected to discover the physical slot associated with that block device.
// For instance, the path /sys/block/nvme0n1 will be a symlink to
// /sys/devices/pciAAAA:BB/AAAA:BB:CC.D/<drive BDF>/nvme/nvme0/nvme0n1
// (The BDF of the device itself is unimportant)
func findBlockDevFromSmbios(sysPath string, s smbios.SystemSlots) ([]string, error) {
	// The SMBIOS table uses the upper 5 bits of this byte to store the device number,
	// and the lower 3 bits to store the function number
	dev := (s.DeviceFunctionNumber & 0b11111000) >> 3
	fn := s.DeviceFunctionNumber & 0b111
	domainBusStr := fmt.Sprintf("%04x:%02x", s.SegmentGroupNumber, s.BusNumber)
	slotBDFPrefix := filepath.Join(sysPath, fmt.Sprintf("devices/pci%s/%s:%02x.%x/", domainBusStr, domainBusStr, dev, fn))

	blockPath := filepath.Join(sysPath, "block/")
	dirEntries, err := os.ReadDir(blockPath)
	if err != nil {
		return nil, err
	}
	var devPaths []string
	for _, dirEntry := range dirEntries {
		path := filepath.Join(blockPath, dirEntry.Name())
		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(realPath, slotBDFPrefix) {
			devPaths = append(devPaths, filepath.Join("/dev", dirEntry.Name()))
		}
	}
	return devPaths, nil
}

func findSlotType(sysPath string, slots []*smbios.SystemSlots, slotType uint8) ([]string, error) {
	var paths []string
	for _, s := range slots {
		if s.SlotType == slotType {
			newPaths, err := findBlockDevFromSmbios(sysPath, *s)
			if err == nil {
				paths = append(paths, newPaths...)
			}
		}
	}

	return paths, nil
}

// FindSlotType searches the SMBIOS table for drives inserted in a slot with the specified type
func FindSlotType(slotType uint8) ([]string, error) {
	smbiosTables, err := smbios.FromSysfs()
	if err != nil {
		return nil, err
	}
	slots, err := smbiosTables.GetSystemSlots()
	if err != nil {
		return nil, err
	}

	return findSlotType("/sys", slots, slotType)
}
