// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forth

import (
	"os"
	"testing"
)

type forthTest struct {
	val string
	res string
	err string
}

var forthTests = []forthTest{
	{"hostname", "", ""},
	{"2", "2", ""},
	{"", "2", "Empty stack"},
	{"2 2 +", "4", ""},
	{"4 2 -", "2", ""},
	{"4 2 *", "8", ""},
	{"4 2 /", "2", ""},
	{"5 2 %", "1", ""},
	{"sb43 hostbase", "43", ""},
	{"sb43 hostbase dup 20 / swap 20 % dup ifelse", "3", ""},
	{"sb40 hostbase dup 20 / swap 20 % dup ifelse", "2", ""},
	{"2 4 swap /", "2", ""},
	{"0 1 1 ifelse", "1", ""},
	{"0 1 0 ifelse", "0", ""},
	{"str cat strcat", "strcat", ""},
	{"1 dup +", "2", ""},
	{"4095 4096 roundup", "4096", ""},
	{"4097 8192 roundup", "8192", ""},
	{"2 x +", "", "parsing \"x\": invalid syntax"},
	{"1 dd +", "2", ""},
	{"1 d3d", "3", ""},
}

func TestForth(t *testing.T) {

	forthTests[0].res, _ = os.Hostname()
	f := New()
	if f.Length() != 0 {
		t.Errorf("Test: stack is %d and should be 0", f.Length())
	}
	if f.Empty() != true {
		t.Errorf("Test: stack is %v and should be false", f.Empty())
	}
	f.Push("test")
	if f.Length() != 1 {
		t.Errorf("Test: stack is %d and should be 1", f.Length())
	}
	if f.Empty() == true {
		t.Errorf("Test: stack is %v and should be false", f.Empty())
	}
	f.Reset()
	if f.Length() != 0 {
		t.Errorf("Test: After Reset(): stack is %d and should be 0", f.Length())
	}
	if f.Empty() != true {
		t.Errorf("Test: After Reset(): stack is %v and should be true", f.Empty())
	}
	NewWord(f, "dd", "dup")
	NewWord(f, "d3d", "dup dup + +")
	for _, tt := range forthTests {
		var err error
		res, err := Eval(f, tt.val)
		if res == tt.res || (err != nil && err.Error() == tt.err) {
			if err != nil {
				/* stack is not going to be right; reset it. */
				f.Reset()
			}
			t.Logf("Test: '%v' '%v' '%v': Pass\n", tt.val, res, err)
		} else {
			t.Errorf("Test: '%v' '%v' '%v': Fail\n", tt.val, res, err)
			t.Logf("ops %v\n", Ops())
		}
		if f.Length() != 0 {
			t.Errorf("Test: %v: stack is %d and should be empty", tt, f.Length())
		}
		if f.Empty() != true {
			t.Errorf("Test: %v: stack is %v and should be empty", tt, f.Empty())
		}
	}

}
