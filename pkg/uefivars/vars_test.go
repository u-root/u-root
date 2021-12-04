// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"testing"
)

// func AllVars() EfiVars
func TestAllVars(t *testing.T) {
	n := 32
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
