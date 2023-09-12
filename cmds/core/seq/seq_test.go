// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela < manoel_vilela@engineer.com >

package main

import (
	"bytes"
	"io"
	"testing"
)

type test struct {
	expect string
	args   []string
}

func testseq(tests []test, format, sep string, width bool, t *testing.T) {
	for _, tst := range tests {
		b := bytes.Buffer{}
		w := io.Writer(&b)
		if err := seq(w, format, sep, width, tst.args); err != nil {
			t.Error(err)
		}

		got := b.Bytes()
		want := []byte(tst.expect)

		if !bytes.Equal(got, want) {
			t.Logf("Got: \n%v\n", string(got))
			t.Logf("Expect: \n%v\n", tst.expect)
			t.Error("Mismatching output")
		}
	}
}

// test default behavior without flags
func TestSeqDefault(t *testing.T) {
	tests := []test{
		{
			args:   []string{"1", "3"},
			expect: "1\n2\n3\n",
		},
		{
			args:   []string{"1", "0.5", "3"},
			expect: "1.0\n1.5\n2.0\n2.5\n3.0\n",
		},
	}

	testseq(tests, "%v", "\n", false, t)
}

// test seq fixed width with leading zeros
func TestSeqWidthEqual(t *testing.T) {
	tests := []test{
		{
			args:   []string{"8", "10"},
			expect: "08\n09\n10\n",
		},
		{
			args:   []string{"8", "0.5", "10"},
			expect: "08.0\n08.5\n09.0\n09.5\n10.0\n",
		},
	}

	testseq(tests, "%v", "\n", true, t)
}

func TestSeqCustomFormat(t *testing.T) {
	tests := []test{
		{
			args:   []string{"8", "10"},
			expect: "08.00\n09.00\n10.00\n",
		},
		{
			args:   []string{"8", "0.5", "10"},
			expect: "08.00\n08.50\n09.00\n09.50\n10.00\n",
		},
	}

	testseq(tests, "%.2f", "\n", true, t)
}

func TestSeqSeparator(t *testing.T) {
	tests := []test{
		{
			args:   []string{"8", "10"},
			expect: "8->9->10\n",
		},
		{
			args:   []string{"8", "0.5", "10"},
			expect: "8.0->8.5->9.0->9.5->10.0\n",
		},
	}

	testseq(tests, "%v", "->", false, t)
}
