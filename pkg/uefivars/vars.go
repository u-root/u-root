// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"bytes"
	"fmt"
	"log"
	"os"
	fp "path/filepath"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// http://kurtqiao.github.io/uefi/2015/01/13/uefi-boot-manager.html

// EfiVarDir is the sysfs /sys/firmware/efi/vars directory, which can be overridden for testing.
var EfiVarDir = "/sys/firmware/efi/vars"

// EfiVar is a generic efi var.
type EfiVar struct {
	UUID, Name string
	Data       []byte
}
type EfiVars []EfiVar

func ReadVar(uuid, name string) (e EfiVar, err error) {
	path := fp.Join(EfiVarDir, name+"-"+uuid, "data")
	e.UUID = uuid
	e.Name = name
	e.Data, err = os.ReadFile(path)
	return
}

// AllVars returns all efi variables
func AllVars() (vars EfiVars) { return ReadVars(nil) }

// ReadVars returns efi variables matching filter
func ReadVars(filt VarFilter) (vars EfiVars) {
	entries, err := fp.Glob(fp.Join(EfiVarDir, "*-*"))
	if err != nil {
		log.Printf("error reading efi vars: %s", err)
		return
	}
	for _, entry := range entries {
		base := fp.Base(entry)
		n := strings.Count(base, "-")
		if n < 5 {
			log.Printf("skipping %s - not a valid var?", base)
			continue
		}
		components := strings.SplitN(base, "-", 2)
		if filt != nil && !filt(components[1], components[0]) {
			continue
		}
		info, err := os.Stat(entry)
		if err == nil && info.IsDir() {
			v, err := ReadVar(components[1], components[0])
			if err != nil {
				log.Printf("reading efi var %s: %s", base, err)
				continue
			}
			vars = append(vars, v)
		}
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
