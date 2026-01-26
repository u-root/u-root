// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"reflect"
	"testing"
)

// func ReadVar(uuid, name string) (e EfiVar, err error) {
func TestVars(t *testing.T) {
	testcases := []struct {
		testname string
		varname  string
		uuid     string
		wantvar  EfiVar
		wanterr  bool
	}{
		{
			testname: "empty input",
			varname:  "",
			uuid:     "",
			wanterr:  true,
		},
		{
			testname: "malformed uuid",
			varname:  "Boot0001",
			uuid:     "8be4df6193ca11d2aa0d00e098032b8c",
			wanterr:  true,
		},
		{
			testname: "BootOrder test",
			varname:  "BootOrder",
			uuid:     "8be4df61-93ca-11d2-aa0d-00e098032b8c",
			wanterr:  false,
			wantvar: EfiVar{
				UUID:       "8be4df61-93ca-11d2-aa0d-00e098032b8c",
				Name:       "BootOrder",
				Attributes: [4]byte{0x07, 0x00, 0x00, 0x00},
				Data:       []byte{0x0a, 0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x05, 0x00, 0x09, 0x00, 0x04, 0x00, 0x06, 0x00},
			},
		},
		{
			testname: "Boot0000 test",
			varname:  "Boot0000",
			uuid:     "8be4df61-93ca-11d2-aa0d-00e098032b8c",
			wanterr:  false,
			wantvar: EfiVar{
				UUID:       "8be4df61-93ca-11d2-aa0d-00e098032b8c",
				Name:       "Boot0000",
				Attributes: [4]byte{0x07, 0x00, 0x00, 0x00},
				Data: []byte{
					0x09, 0x01, 0x00, 0x00, 0x2c, 0x00, 0x55, 0x00, 0x69, 0x00, 0x41, 0x00, 0x70, 0x00, 0x70, 0x00,
					0x00, 0x00, 0x04, 0x07, 0x14, 0x00, 0xc9, 0xbd, 0xb8, 0x7c, 0xeb, 0xf8, 0x34, 0x4f, 0xaa, 0xea,
					0x3e, 0xe4, 0xaf, 0x65, 0x16, 0xa1, 0x04, 0x06, 0x14, 0x00, 0x21, 0xaa, 0x2c, 0x46, 0x14, 0x76,
					0x03, 0x45, 0x83, 0x6e, 0x8a, 0xb6, 0xf4, 0x66, 0x23, 0x31, 0x7f, 0xff, 0x04, 0x00},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.testname, func(t *testing.T) {
			got, err := ReadVar(tc.uuid, tc.varname)
			if (err != nil) != tc.wanterr {
				t.Errorf("Readvar(%v, %v) returned err=%v, but wanterr=%v", tc.uuid, tc.varname, err, tc.wanterr)
			}
			if (err == nil) && !reflect.DeepEqual(got, tc.wantvar) {
				t.Errorf("Readvar(%v, %v) returned %v, wanted %v", tc.uuid, tc.varname, got, tc.wantvar)
			}
		})
	}
}

// func AllVars() EfiVars
func TestAllVars(t *testing.T) {
	n := 33
	vars := AllVars()
	if len(vars) != n {
		t.Errorf("expect %d vars, got %d", n, len(vars))
	}
}

// func DecodeUTF16(b []byte) (string, error)
func TestDecodeUTF16(t *testing.T) {
	want := "TEST"
	got, err := DecodeUTF16([]byte{84, 0, 69, 0, 83, 0, 84, 0})
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

// func (vars EfiVars) Filter(filt VarFilter) EfiVars
func TestFilter(t *testing.T) {
	filt := func(_, _ string) bool { return true }
	v := AllVars()
	matches := v.Filter(AndFilter(filt, NotFilter(filt)))
	if len(matches) != 0 {
		t.Errorf("should be no matches but got\n%#v", matches)
	}
}
