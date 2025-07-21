// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package fpdt

import (
	"io"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
)

// TestFPDT verifies that all data within ACPI FBPT
// is not corrupted by verifiying its checksum.
func TestFPDTChecksum(t *testing.T) {
	var acpiFPDT acpi.Table
	var err error
	if acpiFPDT, err = ReadACPIFPDTTable(); err != nil {
		t.Skip("Skipping test, no ACPI FPDT table in /sys")
	}
	var checksum uint8
	for _, b := range acpiFPDT.Data() {
		checksum += b
	}
	if checksum != 0 {
		t.Fatalf("FPDT Table data corrupted, the Checksum does not equal 0")
	}
}

// Tests if FindFBPTTableAdrr() correctly locates
// the FPBT table within FPDT
func TestFBPTAddress(t *testing.T) {
	var acpiFPDT acpi.Table
	var err error
	if acpiFPDT, err = ReadACPIFPDTTable(); err != nil {
		t.Skip("Skipping test, no ACPI FPDT table in /sys")
	}

	// Get FBPT Pointer from FPDT Table
	var FBPTAddr uint64
	if FBPTAddr, err = FindFBPTTableAdrr(acpiFPDT); err != nil {
		t.Skipf("Skipping test, expected nil but got: %v", err)
	}

	var f *os.File
	if _, err := f.Seek(int64(FBPTAddr), io.SeekStart); err != nil {
		t.Fatal(err)
	}

	var fbptSig [4]byte
	if _, err := io.ReadFull(f, fbptSig[:]); err != nil {
		t.Fatal(err)
	}

	if string(fbptSig[:]) != "FBPT" {
		t.Fatalf("FBPT structure signature check failed. Expected: FBPT, Got: %s", string(fbptSig[:]))
	}
}
