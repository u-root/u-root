// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

type kexecLoadFunc func(entry uintptr, segments kexec.Segments, flags uint64) error

// TODO(chengchieh): move this function to kexec package
func mockKexecLoad(entry uintptr, segments kexec.Segments, flags uint64) error {
	if len(segments) > 16 {
		return fmt.Errorf("number of segments should be less than 16 before dedup")
	}

	segments, err := kexec.AlignAndMerge(segments)
	if err != nil {
		return fmt.Errorf("could not align segments: %w", err)
	}

	if !segments.PhysContains(entry) {
		return fmt.Errorf("entry point %#v is not covered by any segment", entry)
	}
	return nil
}

func TestLoadFvImage(t *testing.T) {
	fv, err := New("testdata/uefi.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = func() (*acpi.RSDP, error) {
		t.Log("mock acpi.GetRSDP()")
		return &acpi.RSDP{}, nil
	}
	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = func() (int64, int64, error) {
		t.Log("mock getSMBIOSBase()")
		return 100, 200, nil
	}
	defer func(old kexecLoadFunc) { kexecLoad = old }(kexecLoad)
	kexecLoad = mockKexecLoad

	defer func(old func() (kexec.MemoryMap, error)) { kexecParseMemoryMap = old }(kexecParseMemoryMap)
	kexecParseMemoryMap = func() (kexec.MemoryMap, error) {
		t.Log("mock kexec.ParseMemMap()")
		return kexec.MemoryMap{}, nil
	}

	if err = fv.Load(true); err != nil {
		t.Fatal(err)
	}
}

func TestLoadFvImageNotFound(t *testing.T) {
	_, err := New("testdata/uefi_NOT_FOUND.fd")
	if err == nil {
		t.Fatal("Should not found the payload.")
	}
}
