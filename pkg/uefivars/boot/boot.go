// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/uefivars"
)

const (
	BootUUID = "8be4df61-93ca-11d2-aa0d-00e098032b8c"
)

// BootEntryVar is a boot entry. It will have the name BootXXXX where XXXX is
// hexadecimal.
type BootEntryVar struct {
	Number uint16 // from the var name
	EfiLoadOption
}

// EfiLoadOption defines the data struct used for vars such as BootXXXX.
//
// As defined in UEFI spec v2.8A:
//
//	typedef struct _EFI_LOAD_OPTION {
//	    UINT32 Attributes;
//	    UINT16 FilePathListLength;
//	    // CHAR16 Description[];
//	    // EFI_DEVICE_PATH_PROTOCOL FilePathList[];
//	    // UINT8 OptionalData[];
//	} EFI_LOAD_OPTION;
type EfiLoadOption struct {
	Attributes         uint32
	FilePathListLength uint16
	Description        string
	FilePathList       EfiDevicePathProtocolList
	OptionalData       []byte
}
type BootEntryVars []*BootEntryVar

// Gets BootXXXX var, if it exists
func ReadBootVar(num uint16) (*BootEntryVar, error) {
	v, err := uefivars.ReadVar(BootUUID, fmt.Sprintf("Boot%04X", num))
	if err != nil {
		return nil, fmt.Errorf("reading var Boot%04X: %w", num, err)
	}
	return BootVar(v), nil
}

// Reads BootCurrent, and from there gets the BootXXXX var referenced.
func ReadCurrentBootVar() (*BootEntryVar, error) {
	curr := ReadBootCurrent()
	if curr == nil {
		return nil, nil
	}
	return ReadBootVar(curr.Current)
}

func (b BootEntryVar) String() string {
	opts, err := uefivars.DecodeUTF16(b.OptionalData)
	if err != nil {
		opts = string(b.OptionalData)
	}
	return fmt.Sprintf("Boot%04X: attrs=0x%x, desc=%q, path=%s, opts=%x", b.Number, b.Attributes, b.Description, b.FilePathList.String(), opts)
}

// AllBootEntryVars returns list of boot entries (BootXXXX). Note that
// BootCurrent, BootOptionSupport, BootNext, BootOrder, etc do not count as
// boot entries.
func AllBootEntryVars() BootEntryVars {
	// BootEntries() is somewhat redundant, but parses the vars into BootEntryVar{}
	return BootEntries(uefivars.ReadVars(BootEntryFilter))
}

// AllBootVars returns all uefi vars that use the boot UUID and whose names begin with Boot.
//
// These include:
//
//   - BootXXXX (individual boot entries, XXXX is hex)
//   - BootCurrent (marks whichever BootXXXX entry was used this boot)
//   - BootOptionSupport
//   - BootNext (can specify a particular entry to use next boot)
//   - BootOrder (the order in which entries are tried)
func AllBootVars() uefivars.EfiVars {
	return uefivars.ReadVars(BootVarFilter)
}

// A VarNameFilter passing boot-related vars. These are a superset of those
// returned by BootEntryFilter.
func BootVarFilter(uuid, name string) bool {
	return uuid == BootUUID && strings.HasPrefix(name, "Boot")
}

// A VarNameFilter passing boot entries. These are a subset of the vars
// returned by BootVarFilter.
func BootEntryFilter(uuid, name string) bool {
	if !BootVarFilter(uuid, name) {
		return false
	}
	// Boot entries begin with BootXXXX-, where XXXX is hex.
	// First, check for the dash.
	if len(name) != 8 {
		return false
	}
	// Try to parse XXXX as hex. If it parses, it's a boot entry.
	_, err := strconv.ParseUint(name[4:], 16, 16)
	return err == nil
}

// BootVar decodes an efivar as a boot entry. use IsBootEntry() to screen first.
func BootVar(v uefivars.EfiVar) (b *BootEntryVar) {
	num, err := strconv.ParseUint(v.Name[4:], 16, 16)
	if err != nil {
		log.Printf("error parsing boot var %s: %s", v.Name, err)
	}
	b = new(BootEntryVar)
	b.Number = uint16(num)
	b.Attributes = binary.LittleEndian.Uint32(v.Data[:4])
	b.FilePathListLength = binary.LittleEndian.Uint16(v.Data[4:6])

	// Description is null-terminated utf16
	var i uint16
	for i = 6; ; i += 2 {
		if v.Data[i] == 0 {
			break
		}
	}
	b.Description, err = uefivars.DecodeUTF16(v.Data[6:i])
	if err != nil {
		log.Printf("reading description: %s (%d -> %x)", err, i, v.Data[6:i])
	}
	b.OptionalData = v.Data[i+2+b.FilePathListLength:]

	b.FilePathList, err = ParseFilePathList(v.Data[i+2 : i+2+b.FilePathListLength])
	if err != nil {
		log.Printf("parsing FilePathList in %s: %s", b.String(), err)
	}
	return
}

// BootCurrentVar represents the UEFI BootCurrent var.
type BootCurrentVar struct {
	uefivars.EfiVar
	Current uint16
}

// BootCurrent returns the BootCurrent var, if any, from the given list.
func BootCurrent(vars uefivars.EfiVars) *BootCurrentVar {
	for _, v := range vars {
		if v.Name == "BootCurrent" {
			return &BootCurrentVar{
				EfiVar:  v,
				Current: uint16(v.Data[1])<<8 | uint16(v.Data[0]),
			}
		}
	}
	return nil
}

// ReadBootCurrent reads and returns the BootCurrent var.
func ReadBootCurrent() *BootCurrentVar {
	v, err := uefivars.ReadVar(BootUUID, "BootCurrent")
	if err != nil {
		log.Printf("reading uefi BootCurrent var: %s", err)
		return nil
	}
	return &BootCurrentVar{
		EfiVar:  v,
		Current: uint16(v.Data[1])<<8 | uint16(v.Data[0]),
	}
}

// BootEntries takes a list of efi vars and parses any that are boot entries,
// returning a list of them.
func BootEntries(vars uefivars.EfiVars) (bootvars BootEntryVars) {
	for _, v := range vars {
		if IsBootEntry(v) {
			bootvars = append(bootvars, BootVar(v))
		}
	}
	return
}

// IsBootEntry returns true if the given var is a boot entry.
func IsBootEntry(e uefivars.EfiVar) bool {
	if e.UUID != BootUUID || len(e.Name) != 8 || e.Name[:4] != "Boot" {
		return false
	}
	_, err := strconv.ParseUint(e.Name[4:], 16, 16)
	return err == nil
}
