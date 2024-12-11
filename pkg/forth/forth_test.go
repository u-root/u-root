// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forth

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

type forthTest struct {
	val   string
	res   Cell
	err   error
	empty bool
}

var forthTests = []forthTest{
	{val: "hostname", res: "", empty: true},
	{val: "2", res: "2", empty: true},
	{val: "2 2 +", res: "4", empty: true},
	{val: "4 2 -", res: "2", empty: true},
	{val: "4 2 *", res: "8", empty: true},
	{val: "4 2 /", res: "2", empty: true},
	{val: "5 2 %", res: "1", empty: true},
	{val: "347 32 mod", res: "27", empty: true},
	{val: "32 347 mod", res: "32", empty: true},
	{val: "sb43 hostbase", res: "43", empty: true},
	{val: "sb43 hostbase dup 20 / swap 20 % dup ifelse", res: "3", empty: true},
	{val: "sb40 hostbase dup 20 / swap 20 % dup ifelse", res: "2", empty: true},
	{val: "2 4 swap /", res: "2", empty: true},
	{val: "0 1 1 ifelse", res: "1", empty: true},
	{val: "0 1 0 ifelse", res: "0", empty: true},
	{val: "str cat strcat", res: "strcat", empty: true},
	{val: "1 dup +", res: "2", empty: true},
	{val: "4095 4096 roundup", res: "4096", empty: true},
	{val: "4097 8192 roundup", res: "8192", empty: true},
	{val: "2 x +", res: "2", err: strconv.ErrSyntax, empty: false},
	{val: "1 dd +", res: "2", empty: true},
	{val: "1 d3d", res: "3", empty: true},
	{val: "drop", res: "", err: ErrEmptyStack, empty: true},
	{val: "5 5 + hostbase", res: "", err: strconv.ErrSyntax, empty: true},
	{val: "1 1 '+ newword 1 1 '+ newword", res: "", err: ErrWordExist, empty: true},
	{val: "1 4 bad newword", res: "", err: ErrNotEnoughElements, empty: false},
	{val: "typeof", res: "", err: ErrEmptyStack, empty: true},
	{val: "1 typeof", res: "string", err: nil, empty: true},
	{val: "zardoz typeof", res: "", err: nil, empty: true},
	{val: "1 %d printf", res: "1\n", empty: true},
	{val: "%d printf", res: "", err: ErrEmptyStack, empty: true},
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
	for i, tt := range forthTests {
		var err error
		res, err := EvalPop(f, tt.val)
		t.Logf("tt %v res %v err %v", tt, res, err)
		if res == tt.res || errors.Is(err, tt.err) {
			t.Logf("Test: '%v' '%v' '%v': Pass\n", tt.val, res, err)
		} else {
			t.Errorf("Test %d: '%v' got (%v, %v): want (%v, %v): Fail\n", i, tt.val, res, err, tt.res, tt.err)
			t.Logf("ops %v\n", Ops())
			continue
		}
		if f.Empty() != tt.empty {
			t.Errorf("Test %d: %v: stack is %v and should be %v", i, tt, f.Empty(), tt.empty)
		}
		/* stack may not be right; reset it. */
		f.Reset()
	}
}

func TestBadPop(t *testing.T) {
	var b [3]byte
	f := New()
	f.Push(b)
	res, err := EvalPop(f, "2 +")
	t.Logf("%v, %v", res, err)
	if err == nil {
		t.Errorf("err: got %v, want %v", err, strconv.ErrSyntax)
	}
	if res != nil {
		t.Errorf("res: got %v, want nil", res)
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
	if err = Eval(f, "0", uint8(0), "+"); err != nil {
		t.Fatalf("got %v, want nil", err)
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
