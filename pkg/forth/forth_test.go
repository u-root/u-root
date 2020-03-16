// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forth

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type forthTest struct {
	val string
	res Cell
	err string
}

var forthTests = []forthTest{
	{"hostname", "", ""},
	{"2", "2", ""},
	{"", "", "[]: length is not 1"},
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
	{"2 x +", nil, "parsing \"x\": invalid syntax"},
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
	NewWord(f, "d3d", "dup", "dup", "+", "+")
	for _, tt := range forthTests {
		var err error
		res, err := EvalPop(f, tt.val)
		t.Logf("tt %v res %v err %v", tt, res, err)
		if res == tt.res || (err != nil && err.Error() == tt.err) {
			if err != nil {
				/* stack is not going to be right; reset it. */
				f.Reset()
			}
			t.Logf("Test: '%v' '%v' '%v': Pass\n", tt.val, res, err)
		} else {
			t.Errorf("Test: '%v' got (%v, %v): want (%v, %v): Fail\n", tt.val, res, err, tt.res, tt.err)
			t.Logf("ops %v\n", Ops())
			continue
		}
		if f.Length() != 0 {
			t.Errorf("Test: %v: stack is %d and should be empty", tt, f.Length())
		}
		if f.Empty() != true {
			t.Errorf("Test: %v: stack is %v and should be empty", tt, f.Empty())
		}
	}

}

func TestBadPop(t *testing.T) {
	var b [3]byte
	f := New()
	f.Push(b)
	res, err := EvalPop(f, "2 +")
	t.Logf("%v, %v", res, err)
	nan := fmt.Errorf("NaN: %T", b)
	if !reflect.DeepEqual(err, nan) {
		t.Errorf("got %v, want %v", err, nan)
	}
	if res != nil {
		t.Errorf("got %v, want nil", res)
	}
}

func TestOpmap(t *testing.T) {
	f := New()
	err := Eval(f, "words")
	if err != nil {
		t.Fatalf("words: got %v, nil", err)
	}
	if f.Length() != 1 {
		t.Fatalf("words: got length %d, want 1", f.Length())
	}
	w := f.Pop()
	switch w.(type) {
	case []string:
	default:
		t.Fatalf("words: got %T, want []string", w)
	}
	t.Logf("words are %v", w)
}

func TestNewWord(t *testing.T) {
	Debug = t.Logf
	f := New()
	// This test creates a word, tp, with 3 args, which simply
	// pushes 1 and 3 on the stack and applies +.
	// Note the use of ' so we can evaluate + as a string,
	// not an operator.
	err := Eval(f, "1", "3", "'+", "3", "tp", "newword")
	if err != nil {
		t.Fatalf("newword: got %v, nil", err)
	}
	err = Eval(f, "tp")
	if err != nil {
		t.Fatalf("newword: got %v, want nil", err)
	}
	t.Logf("stack %v", f.Stack())
}

// make sure that if we die in an Eval nested in an Eval, we fall all the way
// back out.
func TestEvalPanic(t *testing.T) {
	f := New()
	Debug = t.Logf
	err := Eval(f, "0", "'+", "2", "p", "newword")
	if err != nil {
		t.Fatalf("newword: got %v, nil", err)
	}
	t.Logf("p created, now try problems")
	err = Eval(f, "0", uint8(0), "+")
	if err == nil {
		t.Fatal("Got nil, want error")
	}
	t.Logf("Test plus with wrong types: %v", err)
	f.Reset()
	err = Eval(f, "p")
	if err == nil {
		t.Fatal("P with too few args: Got nil, want error")
	}
	t.Logf("p with too few args: %v", err)
	err2 := Eval(f, "p", "0", uint8(0), "+")
	if err2.Error() != err.Error() {
		t.Fatalf("Got %v, want %v", err2, err)
	}
}
