// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"fmt"
	"strings"
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

func mockKexecMemoryMapFromSysfsMemmap() (kexec.MemoryMap, error) {
	return kexec.MemoryMap{}, nil
}

func mockGetRSDP() (*acpi.RSDP, error) {
	return &acpi.RSDP{}, nil
}

func mockGetSMBIOSBase() (int64, int64, error) {
	return 100, 200, nil
}

func TestLoadFvImage(t *testing.T) {
	fv, err := New("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = mockGetRSDP

	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = mockGetSMBIOSBase

	defer func(old kexecLoadFunc) { kexecLoad = old }(kexecLoad)
	kexecLoad = mockKexecLoad

	if err = fv.Load(true); err != nil {
		t.Fatal(err)
	}
}

func TestNewNotFound(t *testing.T) {
	_, err := New("testdata/uefi_NOT_FOUND.fd")
	want := "open testdata/uefi_NOT_FOUND.fd: no such file or directory"
	if err.Error() != want {
		t.Fatalf("Should be '%s', but get '%v'", want, err)
	}
}

func TestNewInvalidPayload(t *testing.T) {
	_, err := New("testdata/fv_with_invalid_sec.fd")
	// for golang >= 1.7
	want1 := "unrecognized PE machine"
	// for golang < 1.7
	want2 := "Unrecognised COFF file header"
	if !(strings.Contains(err.Error(), want1) || strings.Contains(err.Error(), want2)) {
		t.Fatalf("Should be '%s' or '%s', but get '%v'", want1, want2, err)
	}
}

func TestLoadFvImageNotFound(t *testing.T) {
	fv := &FVImage{name: "NOT_FOUND"}
	err := fv.Load(true)

	want := "open NOT_FOUND: no such file or directory"

	if err.Error() != want {
		t.Fatalf("Should be '%s', but get '%v'", want, err)
	}
}

func TestLoadFvImageFailAtMemoryMapFromSysfsMemmap(t *testing.T) {
	fv, err := New("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = func() (kexec.MemoryMap, error) {
		return nil, fmt.Errorf("PARSE_MEMORY_MAP_FAILED")
	}

	err = fv.Load(true)

	want := "PARSE_MEMORY_MAP_FAILED"
	if err.Error() != want {
		t.Fatalf("want '%s', get '%v'", want, err)
	}
}

func TestLoadFvImageFailAtGetRSDP(t *testing.T) {
	fv, err := New("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = func() (*acpi.RSDP, error) {
		return nil, fmt.Errorf("RSDP_NOT_FOUND")
	}

	err = fv.Load(true)

	want := "RSDP_NOT_FOUND"
	if err.Error() != want {
		t.Fatalf("want '%s', get '%v'", want, err)
	}
}

func TestLoadFvImageFailAtGetSMBIOS(t *testing.T) {
	fv, err := New("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = mockGetRSDP

	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = func() (int64, int64, error) {
		t.Log("mock getSMBIOSBase()")
		return 100, 200, fmt.Errorf("SMBIOS_NOT_FOUND")
	}

	err = fv.Load(true)

	want := "SMBIOS_NOT_FOUND"
	if err.Error() != want {
		t.Fatalf("want '%s', get '%v'", want, err)
	}
}

func TestLoadFvImageFailAtKexec(t *testing.T) {
	fv, err := New("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatal(err)
	}

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap

	defer func(old func() (*acpi.RSDP, error)) { getRSDP = old }(getRSDP)
	getRSDP = mockGetRSDP

	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = mockGetSMBIOSBase

	defer func(old kexecLoadFunc) { kexecLoad = old }(kexecLoad)
	kexecLoad = func(entry uintptr, segments kexec.Segments, flags uint64) error {
		return fmt.Errorf("KEXEC_FAILED")
	}

	err = fv.Load(true)

	want := "kexec.Load() error: KEXEC_FAILED"
	if err.Error() != want {
		t.Fatalf("want '%s', get '%v'", want, err)
	}
}
