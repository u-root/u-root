// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions
package tss

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"
)

func readTPM12VendorAttributes(rwc io.ReadWriter) (TPMInfo, error) {
	var vendorInfo string

	_, err := tpm.GetManufacturer(rwc)
	if err != nil {
		return TPMInfo{}, err
	}

	return TPMInfo{
		VendorInfo:           strings.Trim(vendorInfo, "\x00"), // Stubbed
		Manufacturer:         TCGVendorID(uint32(0)),           // Stubbed
		FirmwareVersionMajor: int(0),                           // Stubbed
		FirmwareVersionMinor: int(0),                           // Stubbed
	}, nil
}

func readTPM20VendorAttributes(rwc io.ReadWriter) (TPMInfo, error) {
	var vendorInfo string
	// The Vendor String is split up into 4 sections of 4 bytes,
	// for a maximum length of 16 octets of ASCII text. We iterate
	// through the 4 indexes to get all 16 bytes & construct vendorInfo.
	// See: TPM_PT_VENDOR_STRING_1 in TPM 2.0 Structures reference.
	for i := 0; i < 4; i++ {
		caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtVendorString+uint32(i))
		if err != nil {
			return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_VENDOR_STRING_%d) failed: %v", i+1, err)
		}
		subset, ok := caps[0].(tpm2.TaggedProperty)
		if !ok {
			return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
		}
		// Reconstruct the 4 ASCII octets from the uint32 value.
		vendorInfo += string(subset.Value&0xFF000000) + string(subset.Value&0xFF0000) + string(subset.Value&0xFF00) + string(subset.Value&0xFF)
	}

	caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtManufacturer)
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_MANUFACTURER) failed: %v", err)
	}
	manu, ok := caps[0].(tpm2.TaggedProperty)
	if !ok {
		return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
	}

	caps, _, err = tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, tpmPtFwVersion1)
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_FIRMWARE_VERSION_1) failed: %v", err)
	}
	fw, ok := caps[0].(tpm2.TaggedProperty)
	if !ok {
		return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
	}

	return TPMInfo{
		VendorInfo:           strings.Trim(vendorInfo, "\x00"),
		Manufacturer:         TCGVendorID(manu.Value),
		FirmwareVersionMajor: int((fw.Value & 0xffff0000) >> 16),
		FirmwareVersionMinor: int(fw.Value & 0x0000ffff),
	}, nil
}
