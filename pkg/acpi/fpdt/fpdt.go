// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fpdt reads FPDT ACPI table and gets FPDT record information
package fpdt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/acpi"
)

const (
	// ACPI FPDT table
	acpiFPDTSig = "FPDT"
)

// ReadACPIFPDTTable finds which ACPI table is FPDT and returns it
func ReadACPIFPDTTable() (acpi.Table, error) {
	var AcpiFPDT acpi.Table

	// Get ACPI tables via the acpi packege
	_, tables, err := acpi.GetTable()
	if err != nil {
		return AcpiFPDT, err
	}

	// Scan the tables to find ACPI.FPDT
	//	Signature = "FPDT"

	for _, t := range tables {
		sigMatch := (t.Sig() == acpiFPDTSig)
		if sigMatch {
			return t, nil
		}
	}
	return nil, errors.New("unable to find FPDT")
}

// FindFBPTTableAdrr finds FPBT Table Address
func FindFBPTTableAdrr(t acpi.Table) (uint64, error) {
	var addr uint64

	if t.Sig() != "FPDT" {
		return addr, fmt.Errorf("wrong table type passed. Table Signature %s", t.Sig())
	}

	for i := 0; i < len(t.TableData()); i += int(t.TableData()[i+2]) {
		// Find Firmware Basic Boot Performance Pointer Record
		// see ACPI Table Spec: https://uefi.org/sites/default/files/resources/ACPI%206_2_A_Sept29.pdf (page 210)
		if t.TableData()[i] == 0x00 && t.TableData()[i+1] == 0x00 {
			addr = binary.NativeEndian.Uint64(t.TableData()[i+8 : i+16])
			return addr, nil
		}
	}
	return addr, errors.New("unable to find FPBT Address")
}

// ReadFPDTRecordHeader reads Header for records
// found in FPDT Table as found in ACPI spec.
// returns (HeaderType, HeaderLength, HeaderRevision)
// ACPI Table Spec https://uefi.org/sites/default/files/resources/ACPI%206_2_A_Sept29.pdf (page 208)
func ReadFPDTRecordHeader(mem io.ReadSeeker) (uint16, uint8, uint8, error) {
	var HeaderType [2]byte
	var HeaderLength [1]byte
	var HeaderRevision [1]byte

	if _, err := io.ReadFull(mem, HeaderType[:]); err != nil {
		return 0, 0, 0, err
	}
	if _, err := io.ReadFull(mem, HeaderLength[:]); err != nil {
		return 0, 0, 0, err
	}
	if _, err := io.ReadFull(mem, HeaderRevision[:]); err != nil {
		return 0, 0, 0, err
	}

	return binary.NativeEndian.Uint16(HeaderType[:]), uint8(HeaderLength[0]), uint8(HeaderRevision[0]), nil
}
