// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// http://kurtqiao.github.io/uefi/2015/01/13/uefi-boot-manager.html

// EfiVarDir is the older kernel sysfs /sys/firmware/efi/vars directory, which can be overridden for testing. Each variable is represented by a directory. Superceded by efivarfs (see below)
var EfiVarDir = "/sys/firmware/efi/vars"

// EFIVarfsDir is the kernel efivarfs /sys/firmware/efi/efivars directory, which can be overridden for testing. Each variable is represented by a file.
var EfiVarfsDir = "/sys/firmware/efi/efivars"

// EfiVar is a generic efi var.
type EfiVar struct {
	UUID, Name string
	Attributes [4]byte
	Data       []byte
}
type EfiVars []EfiVar

func ReadVar(uuid, name string) (e EfiVar, err error) {
	var path string
	_, err = os.Stat(EfiVarfsDir)
	if err == nil {
		path = filepath.Join(EfiVarfsDir, name+"-"+uuid)
		e.Data, err = os.ReadFile(path)
		if err == nil && len(e.Data) > 4 {
			// first 4 bytes are UEFI variable attributes
			e.Attributes = [4]byte(e.Data[:4])
			e.Data = e.Data[4:]
			e.UUID = uuid
			e.Name = name
			return
		}
	}
	// fallback to efivars
	path = filepath.Join(EfiVarDir, name+"-"+uuid, "data")
	e.Data, err = os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("could not find EFI variable in either sysfs or efivarfs kernel interface: %w", err)
		return
	}
	e.UUID = uuid
	e.Name = name

	// read attributes
	var attr uint32
	path = filepath.Join(EfiVarDir, name+"-"+uuid, "attributes")
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("could not read EFI variable attributes for variable %s-%s through sysfs interface: %w", name, uuid, err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// see section 8.2.1, GetVariable() of the UEFI specification
		switch line {
		case "EFI_VARIABLE_NON_VOLATILE":
			attr |= 0x00000001
		case "EFI_VARIABLE_BOOTSERVICE_ACCESS":
			attr |= 0x00000002
		case "EFI_VARIABLE_RUNTIME_ACCESS":
			attr |= 0x00000004
		case "EFI_VARIABLE_HARDWARE_ERROR_RECORD":
			attr |= 0x00000008
		case "EFI_VARIABLE_TIME_BASED_AUTHENTICATED_WRITE_ACCESS":
			attr |= 0x00000020
		case "EFI_VARIABLE_APPEND_WRITE":
			attr |= 0x00000040
		case "EFI_VARIABLE_ENHANCED_AUTHENTICATED_ACCESS":
			attr |= 0x00000080
		}
	}
	binary.LittleEndian.PutUint32(e.Attributes[:], attr)
	return
}

// AllVars returns all efi variables
func AllVars() (vars EfiVars) { return ReadVars(nil) }

// ReadVars returns efi variables matching filter
func ReadVars(filt VarFilter) (vars EfiVars) {
	entries, _ := filepath.Glob(filepath.Join(EfiVarfsDir, "*-*"))
	if len(entries) == 0 {
		// fallback to efivar
		entries, _ = filepath.Glob(filepath.Join(EfiVarDir, "*-*"))
	}
	for _, entry := range entries {
		base := filepath.Base(entry)
		n := strings.Count(base, "-")
		if n < 5 {
			continue
		}
		components := strings.SplitN(base, "-", 2)
		if filt != nil && !filt(components[1], components[0]) {
			continue
		}
		v, err := ReadVar(components[1], components[0])
		if err != nil {
			continue
		}
		vars = append(vars, v)
	}
	return
}

// VarFilter is a type of function used to filter efi vars
type VarFilter func(uuid, name string) bool

// NotFilter returns a filter negating the given filter.
func NotFilter(f VarFilter) VarFilter {
	return func(u, n string) bool { return !f(u, n) }
}

// AndFilter returns true only if all given filters return true.
func AndFilter(filters ...VarFilter) VarFilter {
	return func(u, n string) bool {
		for _, f := range filters {
			if !f(u, n) {
				return false
			}
		}
		return true
	}
}

// Filter returns the elements of the list for which the filter function
// returns true.
func (vars EfiVars) Filter(filt VarFilter) EfiVars {
	var res EfiVars
	for _, v := range vars {
		if filt(v.UUID, v.Name) {
			res = append(res, v)
		}
	}
	return res
}

// DecodeUTF16 decodes the input as a utf16 string.
// https://gist.github.com/bradleypeabody/185b1d7ed6c0c2ab6cec
func DecodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("must have even length byte slice")
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = BytesToU16(b[i : i+2])
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}

// BytesToU16 converts a []byte of length 2 to a uint16.
func BytesToU16(b []byte) uint16 {
	if len(b) != 2 {
		log.Fatalf("bytesToU16: bad len %d (%x)", len(b), b)
	}
	return uint16(b[0]) + (uint16(b[1]) << 8)
}
