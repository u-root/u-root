// Copyright 2020-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tss

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	tpm2 "github.com/google/go-tpm/legacy/tpm2"
	tpm1 "github.com/google/go-tpm/tpm"
)

func readTPM12Information(rwc io.ReadWriter) (TPMInfo, error) {
	manufacturerRaw, err := tpm1.GetManufacturer(rwc)
	if err != nil {
		return TPMInfo{}, err
	}

	manufacturerID := binary.BigEndian.Uint32(manufacturerRaw)
	return TPMInfo{
		VendorInfo:   TCGVendorID(manufacturerID).String(),
		Manufacturer: TCGVendorID(manufacturerID),
	}, nil
}

func readTPM20Information(rwc io.ReadWriter) (TPMInfo, error) {
	var vendorInfo string
	// The Vendor String is split up into 4 sections of 4 bytes,
	// for a maximum length of 16 octets of ASCII text. We iterate
	// through the 4 indexes to get all 16 bytes & construct vendorInfo.
	// See: TPM_PT_VENDOR_STRING_1 in TPM 2.0 Structures reference.
	for i := 0; i < 4; i++ {
		caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, uint32(tpm2.VendorString1)+uint32(i))
		if err != nil {
			return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_VENDOR_STRING_%d) failed: %w", i+1, err)
		}
		subset, ok := caps[0].(tpm2.TaggedProperty)
		if !ok {
			return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
		}
		// Reconstruct the 4 ASCII octets from the uint32 value.
		vendorInfo += string([]byte{
			byte(subset.Value >> 24),
			byte(subset.Value >> 16),
			byte(subset.Value >> 8),
			byte(subset.Value),
		})
	}

	caps, _, err := tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, uint32(tpm2.Manufacturer))
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_MANUFACTURER) failed: %w", err)
	}
	manu, ok := caps[0].(tpm2.TaggedProperty)
	if !ok {
		return TPMInfo{}, fmt.Errorf("got capability of type %T, want tpm2.TaggedProperty", caps[0])
	}

	caps, _, err = tpm2.GetCapability(rwc, tpm2.CapabilityTPMProperties, 1, uint32(tpm2.FirmwareVersion1))
	if err != nil {
		return TPMInfo{}, fmt.Errorf("tpm2.GetCapability(PT_FIRMWARE_VERSION_1) failed: %w", err)
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

func takeOwnership12(rwc io.ReadWriteCloser, ownerPW, srkPW string) error {
	var ownerAuth [20]byte
	var srkAuth [20]byte

	if ownerPW != "" {
		ownerAuth = sha1.Sum([]byte(ownerPW))
	}

	if srkPW != "" {
		srkAuth = sha1.Sum([]byte(srkPW))
	}

	pubek, err := tpm1.ReadPubEK(rwc)
	if err != nil {
		return err
	}

	if err := tpm1.TakeOwnership(rwc, ownerAuth, srkAuth, pubek); err != nil {
		return err
	}
	return nil
}

func takeOwnership20(rwc io.ReadWriteCloser, ownerPW, srkPW string) error {
	return fmt.Errorf("not supported by go-tpm for TPM2.0")
}

func clearOwnership12(rwc io.ReadWriteCloser, ownerPW string) error {
	var ownerAuth [20]byte

	if ownerPW != "" {
		ownerAuth = sha1.Sum([]byte(ownerPW))
	}

	err := tpm1.OwnerClear(rwc, ownerAuth)
	if err != nil {
		err := tpm1.ForceClear(rwc)
		if err != nil {
			return fmt.Errorf("couldn't clear TPM 1.2 with ownerauth nor force clear")
		}
	}

	return nil
}

func clearOwnership20(rwc io.ReadWriteCloser, ownerPW string) error {
	return fmt.Errorf("not supported by go-tpm for TPM2.0")
}

func readPubEK12(rwc io.ReadWriteCloser, ownerPW string) ([]byte, error) {
	var ownerAuth [20]byte
	if ownerPW != "" {
		ownerAuth = sha1.Sum([]byte(ownerPW))
	}

	ek, err := tpm1.OwnerReadPubEK(rwc, ownerAuth)
	if err != nil {
		return nil, err
	}

	return ek, nil
}

func readPubEK20(rwc io.ReadWriteCloser, ownerPW string) ([]byte, error) {
	return nil, fmt.Errorf("not supported by go-tpm for TPM2.0")
}

func resetLockValue12(rwc io.ReadWriteCloser, ownerPW string) (bool, error) {
	var ownerAuth [20]byte
	if ownerPW != "" {
		ownerAuth = sha1.Sum([]byte(ownerPW))
	}

	if err := tpm1.ResetLockValue(rwc, ownerAuth); err != nil {
		return false, err
	}
	return true, nil
}

func resetLockValue20(rwc io.ReadWriteCloser, ownerPW string) (bool, error) {
	return false, fmt.Errorf("not yet supported by tss")
}
