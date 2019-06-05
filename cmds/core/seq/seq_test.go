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

func resetFlags() {
	flags.format = "%v"
	flags.separator = "\n"
	flags.widthEqual = false
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

		if !reflect.DeepEqual(got, want) {
			t.Logf("Got: \n%v\n", string(got))
			t.Logf("Expect: \n%v\n", tst.expect)
			t.Error("Mismatching output")
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
			"1.0\n1.5\n2.0\n2.5\n3.0\n",
		},
	}

	testseq(tests, t)
}

// test seq fixed width with leading zeros
func TestSeqWidthEqual(t *testing.T) {
	flags.widthEqual = true
	defer resetFlags()
	var tests = []test{
		{
			[]string{"8", "10"},
			"08\n09\n10\n",
		},
		{
			[]string{"8", "0.5", "10"},
			"08.0\n08.5\n09.0\n09.5\n10.0\n",
		},
	}

	testseq(tests, t)
}

func TestSeqCustomFormat(t *testing.T) {
	flags.format = "%.2f"
	flags.widthEqual = true
	defer resetFlags()
	var tests = []test{
		{
			[]string{"8", "10"},
			"08.00\n09.00\n10.00\n",
		},
		{
			[]string{"8", "0.5", "10"},
			"08.00\n08.50\n09.00\n09.50\n10.00\n",
		},
	}

	testseq(tests, t)
}

func TestSeqSeparator(t *testing.T) {
	flags.separator = "->"
	defer resetFlags()
	var tests = []test{
		{
			[]string{"8", "10"},
			"8->9->10\n",
		},
		{
			[]string{"8", "0.5", "10"},
			"8.0->8.5->9.0->9.5->10.0\n",
		},
	}

	testseq(tests, t)

}
