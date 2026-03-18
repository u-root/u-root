// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela < manoel_vilela@engineer.com >

package main

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type test struct {
	expect string
	err    error
	args   []string
}

func testseq(t *testing.T, tests []test, format, sep string, width bool) {
	for _, tt := range tests {
		b := bytes.Buffer{}
		w := io.Writer(&b)
		err := seq(w, format, sep, width, tt.args)

		if !errors.Is(err, tt.err) {
			t.Fatalf("expected %v, got %v", tt.err, err)
		}

		if diff := cmp.Diff(b.String(), tt.expect); diff != "" {
			t.Errorf("unexpected result (-want +got):\n%s", diff)
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
		{
			args:   []string{"3"},
			expect: "1\n2\n3\n",
		},
		{
			args: []string{"1", "2", "3", "4"},
			err:  errUsage,
		},
		{
			args:   []string{"3", "1"},
			expect: "3\n2\n1\n",
		},
		{
			args:   []string{"5", "-2", "1"},
			expect: "5\n3\n1\n",
		},
		{
			args:   []string{"-2"},
			expect: "1\n0\n-1\n-2\n",
		},
		{
			args: []string{"1", "-1", "2"},
			err:  errPositiveInc,
		},
		{
			args: []string{"10", "1", "2"},
			err:  errNegativeDec,
		},
		{
			args: []string{"10", "0", "2"},
			err:  errZeroDec,
		},
		{
			args:   []string{"5", "-1", "5"},
			expect: "5\n",
		},
		{
			args:   []string{"5", "1", "5"},
			expect: "5\n",
		},
		{
			args: []string{"hello", "2"},
			err:  strconv.ErrSyntax,
		},
		{
			args: []string{"2", "hello"},
			err:  strconv.ErrSyntax,
		},
		{
			args: []string{"hello"},
			err:  strconv.ErrSyntax,
		},
		{
			args: []string{"hello", "1", "2"},
			err:  strconv.ErrSyntax,
		},
		{
			args: []string{"1", "hello", "2"},
			err:  strconv.ErrSyntax,
		},
		{
			args: []string{"1", "2", "hello"},
			err:  strconv.ErrSyntax,
		},
	}

	testseq(t, tests, "%v", "\n", false)
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

	testseq(t, tests, "%v", "\n", true)
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

	testseq(t, tests, "%.2f", "\n", true)
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

	testseq(t, tests, "%v", "->", false)
}

func TestMaxDecimalPlaces(t *testing.T) {
	tests := []struct {
		args     []string
		expected int
	}{
		{
			args:     []string{"10", "2", "8"},
			expected: 0,
		},
		{
			args:     []string{"10.1", "2", "8"},
			expected: 1,
		},
		{
			args:     []string{"10.1", "2.02", "8.003"},
			expected: 3,
		},
	}

	for _, tt := range tests {
		res := maxDecimalPlaces(tt.args)
		if res != tt.expected {
			t.Errorf("expected %v got %v", tt.expected, res)
		}
	}
}
