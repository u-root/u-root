// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package fbpt

import (
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
