// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela < manoel_vilela@engineer.com >

package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type test struct {
	args   []string
	expect string
}

func testseq(tests []test, t *testing.T) {
	for _, tst := range tests {
		b := bytes.Buffer{}
		w := io.Writer(&b)
		if err := seq(w, tst.args); err != nil {
			t.Error(err)
		}

		got := b.Bytes()
		want := []byte(tst.expect)
		t.Logf("Got: \n%v\n", string(got))
		t.Logf("Expect: \n%v\n", tst.expect)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Mismatching output; got %v, want %v", got, want)
		}
	}
}

// test default behavior without flags
func TestSeqDefault(t *testing.T) {
	var tests = []test{
		{
			[]string{"1", "3"},
			"1\n2\n3\n",
		},
		{
			[]string{"1", "0.5", "3"},
			"1\n1.5\n2\n2.5\n3\n",
		},
	}

	testseq(tests, t)
}

// test seq fixed width with leading zeros
func TestSeqWidthEqual(t *testing.T) {
	flags.widthEqual = true
	var tests = []test{
		{
			[]string{"8", "10"},
			"08\n09\n10\n",
		},
		{
			[]string{"8", "0.5", "10"},
			"08\n8.5\n09\n9.5\n10\n",
		},
	}

	testseq(tests, t)

}
