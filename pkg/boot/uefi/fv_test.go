// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"os"
	"testing"
)

func TestFindSecurityCorePEEntry(t *testing.T) {
	dat, err := os.ReadFile("testdata/fv_with_sec.fd")
	if err != nil {
		t.Fatalf("fail to read firmware volume: %v", err)
	}
	offset, err := findSecurityCorePEEntry(dat)
	if err != nil {
		t.Fatalf("fail to find SEC in Firmware Volume: %v", err)
	}
	want := 0x15da0
	if offset != want {
		t.Fatalf("want '%x', get '%x'", want, offset)
	}
}

func TestFindSecurityCorePEEntryNestedSec(t *testing.T) {
	dat, err := os.ReadFile("testdata/fv_with_nested_sec.fd")
	if err != nil {
		t.Fatalf("fail to read firmware volume: %v", err)
	}
	offset, err := findSecurityCorePEEntry(dat)
	if err != nil {
		t.Fatalf("fail to find SEC in Firmware Volume: %v", err)
	}
	want := 0x160b4
	if offset != want {
		t.Fatalf("want '%x', get '%x'", want, offset)
	}
}

func TestFindSecurityCorePEEntryNotFound(t *testing.T) {
	dat, err := os.ReadFile("testdata/fv_without_sec.fd")
	if err != nil {
		t.Fatalf("fail to read firmware volume: %v", err)
	}
	_, err = findSecurityCorePEEntry(dat)
	if err == nil {
		t.Fatalf("should not found a sec in uefi_no_sec.fd")
	}
}

func TestFFSHeaderUnmarshalBinaryFailForSize(t *testing.T) {
	var fh EFIFFSFileHeader
	err := fh.UnmarshalBinary([]byte{0x0})
	want := "invalid entry point stucture length 1"
	if err.Error() != want {
		t.Fatalf("Should be '%s', but get '%v'", want, err)
	}
}

func TestUnmarshalBinaryFailForSize(t *testing.T) {
	var fvh EFIFirmwareVolumeHeader
	err := fvh.UnmarshalBinary([]byte{0x0})
	want := "invalid entry point stucture length 1"
	if err.Error() != want {
		t.Fatalf("Should be '%s', but get '%v'", want, err)
	}
}

func TestIncorrectFVHSignature(t *testing.T) {
	dat := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x78, 0xe5, 0x8c, 0x8c, 0x3d, 0x8a, 0x1c, 0x4f,
		0x99, 0x35, 0x89, 0x61, 0x85, 0xc3, 0x2d, 0xd3,
		0x00, 0xd0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x5f, 0x56, 0x56, 0x48, 0xff, 0xfe, 0x07, 0x00,
		0x48, 0x00, 0x4e, 0x16, 0x60, 0x00, 0x00, 0x02,
		0x1d, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	var fvh EFIFirmwareVolumeHeader
	err := fvh.UnmarshalBinary(dat)
	want := "invalid Signature string \"_VVH\""
	if err.Error() != want {
		t.Fatalf("Should be '%s', but get '%v'", want, err)
	}
}
