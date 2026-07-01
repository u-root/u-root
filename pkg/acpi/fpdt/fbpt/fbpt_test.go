// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package fbpt

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/acpi/fpdt"
)

// TestFBPT just verifies that the FBPT table reads OK.
// It does not verify content as content varies all
// the time.
func TestFBPT(t *testing.T) {
	var acpiFPDT acpi.Table
	var err error
	if acpiFPDT, err = fpdt.ReadACPIFPDTTable(); err != nil {
		t.Skip("Skipping test, no ACPI FPDT table in /sys")
	}

	// Get FBPT Pointer from FPDT Table
	var FBPTAddr uint64
	if FBPTAddr, err = fpdt.FindFBPTTableAdrr(acpiFPDT); err != nil {
		t.Skipf("Skipping test, expected nil but got: %v", err)
	}

	if _, _, _, err = FindAllFBPTRecords(FBPTAddr); err != nil {
		t.Fatalf("Unable to read FBPT records: %v", err)
	}
}

// TestDynamicRecordLengthUnderflow checks that a dynamic string record whose
// declared length is below the 34 byte fixed size is rejected instead of
// underflowing recordLength-34 and reading adjacent memory into Description.
func TestDynamicRecordLengthUnderflow(t *testing.T) {
	// 30 fixed bytes followed by bytes standing in for adjacent memory.
	buf := make([]byte, 30+300)
	for i := 30; i < len(buf); i++ {
		buf[i] = 0xAA
	}

	if _, err := readFirmwarePerformanceDataTableDynamicRecord(bytes.NewReader(buf), 30); err == nil {
		t.Fatal("expected error for record length 30, got nil")
	}

	want := "uroot!"
	good := append(append([]byte{}, buf[:30]...), []byte(want)...)
	rec, err := readFirmwarePerformanceDataTableDynamicRecord(bytes.NewReader(good), 40)
	if err != nil {
		t.Fatalf("unexpected error for valid record: %v", err)
	}
	if rec.Description != want {
		t.Fatalf("Description = %q, want %q", rec.Description, want)
	}
}
