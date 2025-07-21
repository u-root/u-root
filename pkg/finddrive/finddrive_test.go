// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package finddrive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/smbios"
)

const (
	matchingSlotType    = 0xBB
	nonMatchingSlotType = 0x66
	missingSlotType     = 0x43
	devName             = "nvme0n1"
	matchedSlotPath     = "/dev/" + devName
)

var (
	matchingSlot = smbios.SystemSlots{
		SlotType:             matchingSlotType,
		SegmentGroupNumber:   0x0A5A,
		BusNumber:            0x44,
		DeviceFunctionNumber: 0x96,
	}
	nonMatchingSlot = smbios.SystemSlots{
		SlotType: nonMatchingSlotType,
	}
)

func mockSysDir(t *testing.T) string {
	sysDir := t.TempDir()

	devicesPath := filepath.Join(sysDir, "devices/pci0a5a:44/0a5a:44:12.6/0a5a:00:00.0/nvme/nvme0/")
	err := os.MkdirAll(devicesPath, 0o777)
	if err != nil {
		t.Errorf("Error creating path %s: %v", devicesPath, err)
	}

	err = os.MkdirAll(filepath.Join(sysDir, "block/"), 0o777)
	if err != nil {
		t.Errorf("Error creating path %s: %v", filepath.Join(sysDir, "block/"), err)
	}

	devicesFile := filepath.Join(devicesPath, devName)
	f, err := os.Create(devicesFile)
	if err != nil {
		t.Errorf("Error creating file %s: %v", devicesFile, err)
	}
	f.Close()

	err = os.Symlink(devicesFile, filepath.Join(sysDir, "block/nvme0n1"))
	if err != nil {
		t.Errorf("Error creating symlink: %v", err)
	}

	return sysDir
}

func TestFindSlotType(t *testing.T) {
	sysDir := mockSysDir(t)
	slots := []*smbios.SystemSlots{&nonMatchingSlot, &matchingSlot, &nonMatchingSlot, &matchingSlot, &nonMatchingSlot}

	paths, err := findSlotType(sysDir, slots, matchingSlotType)
	if err != nil {
		t.Errorf("findSlotType(%v, %v, %v) returned error: %v, want: nil", sysDir, slots, matchingSlotType, err)
	}
	if len(paths) != 2 || paths[0] != matchedSlotPath || paths[1] != matchedSlotPath {
		t.Errorf("findSlotType(%v, %v, %v) returned paths: %v, want: %v", sysDir, slots, matchingSlotType, paths, []string{matchedSlotPath, matchedSlotPath})
	}
}

func TestFindSlotTypeMissing(t *testing.T) {
	sysDir := mockSysDir(t)
	slots := []*smbios.SystemSlots{&nonMatchingSlot, &matchingSlot, &nonMatchingSlot, &matchingSlot, &nonMatchingSlot}

	paths, err := findSlotType(sysDir, slots, missingSlotType)
	if err != nil {
		t.Errorf("findSlotType(%v, %v, %v) returned error: %v, want: nil", sysDir, slots, missingSlotType, err)
	}
	if len(paths) != 0 {
		t.Errorf("findSlotType(%v, %v, %v) returned paths: %v, want: []", sysDir, slots, missingSlotType, paths)
	}
}

func TestFindSlotTypeNoSlots(t *testing.T) {
	sysDir := mockSysDir(t)
	slots := []*smbios.SystemSlots{}

	paths, err := findSlotType(sysDir, slots, matchingSlotType)
	if err != nil {
		t.Errorf("findSlotType(%v, %v, %v) returned error: %v, want: nil", sysDir, slots, matchingSlotType, err)
	}
	if len(paths) != 0 {
		t.Errorf("findSlotType(%v, %v, %v) returned paths: %v, want: []", sysDir, slots, matchingSlotType, paths)
	}
}
