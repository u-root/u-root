// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

type kexecLoadFunc func(entry uintptr, segments kexec.Segments, flags uint64) error

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
	// TODO(chengchieh): refactor kexec pkg and create a real mock function. A real
	// kexec mock load should include segments and alignment check.
	defer func(old kexecLoadFunc) { kexecLoad = old }(kexecLoad)
	kexecLoad = func(entry uintptr, segments kexec.Segments, flags uint64) error {
		t.Log("mock kexec.Load()")
		return nil
	}

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
