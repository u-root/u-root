// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lineio

import (
	"bytes"
	"io"
	"reflect"
	"regexp"
	"testing"
)

type lineCase struct {
	line    int64
	bufSize int
	err     error
	data    string
	size    int
}

type dataCase struct {
	data  string
	tests []lineCase
}

var cases = []dataCase{
	{
		data: `Hello World!`,
		tests: []lineCase{
			{
				line:    1,
				bufSize: 128,
				err:     io.EOF,
				data:    "Hello World!",
				size:    12,
			},
			{
				line:    2,
				bufSize: 128,
				err:     io.EOF,
				data:    "",
				size:    0,
			},
		},
	},
	{
		data: `Line 1
Line 2
Line 3
Line 4
Line 5

Line 7
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
Line 9`,
		tests: []lineCase{
			{
				line:    1,
				bufSize: 128,
				err:     io.EOF,
				data:    "Line 1",
				size:    6,
			},
			{
				line:    2,
				bufSize: 128,
				err:     io.EOF,
				data:    "Line 2",
				size:    6,
			},
			// Skip some lines before the next one.
			{
				line:    5,
				bufSize: 128,
				err:     io.EOF,
				data:    "Line 5",
				size:    6,
			},
			// Exact size
			{
				line:    1,
				bufSize: 6,
				err:     nil,
				data:    "Line 1",
				size:    6,
			},
			// Less then entire line
			{
				line:    1,
				bufSize: 4,
				err:     nil,
				data:    "Line",
				size:    4,
			},
			// Empty line
			{
				line:    6,
				bufSize: 128,
				err:     io.EOF,
				data:    "",
				size:    0,
			},
			// After empty line
			{
				line:    7,
				bufSize: 128,
				err:     io.EOF,
				data:    "Line 7",
				size:    6,
			},
			// After long line
			{
				line:    9,
				bufSize: 128,
				err:     io.EOF,
				data:    "Line 9",
				size:    6,
			},
		},
	},
}

func TestReadLine(t *testing.T) {
	for _, c := range cases {
		r := NewLineReader(bytes.NewReader([]byte(c.data)))

		for _, l := range c.tests {
			buf := make([]byte, l.bufSize)

			n, err := r.ReadLine(buf, l.line)
			if err != l.err {
				t.Errorf("data: '%s', ReadLine(%d): err got %v want %v", c.data, l.line, err, l.err)
			}

			if n != l.size {
				t.Errorf("data: '%s', ReadLine(%d): n got %d want %d", c.data, l.line, n, l.size)
			}

			s := string(buf[:n])

			if s != l.data {
				t.Errorf("data: '%s', ReadLine(%d): buf got '%s' want '%s'", c.data, l.line, s, l.data)
			}
		}
	}

}

func TestLineExists(t *testing.T) {
	input := `Line 1
Line 2
Line 3`

	r := NewLineReader(bytes.NewReader([]byte(input)))

	if !r.LineExists(1) {
		t.Errorf("LineExists(1) = false want true")
	}

	if !r.LineExists(2) {
		t.Errorf("LineExists(2) = false want true")
	}

	if !r.LineExists(3) {
		t.Errorf("LineExists(3) = false want true")
	}

	if r.LineExists(4) {
		t.Errorf("LineExists(4) = true want false")
	}
}

// Newline at the end of a file shouldn't count as a line.
func TestLineExistsNewline(t *testing.T) {
	input := `Line 1
`

	r := NewLineReader(bytes.NewReader([]byte(input)))

	if !r.LineExists(1) {
		t.Errorf("LineExists(1) = false want true")
	}

	if r.LineExists(2) {
		t.Errorf("LineExists(2) = true want false")
	}
}

func TestSearchLine(t *testing.T) {
	input := `Line 1
aaa bbb ccc
Last line`

	searchCases := []struct {
		line int64
		reg  *regexp.Regexp
		ret  [][]int
		err  error
	}{
		{
			line: 1,
			reg:  regexp.MustCompile(`Line 1`),
			ret:  [][]int{{0, 6}},
			err:  nil,
		},
		{
			line: 1,
			reg:  regexp.MustCompile(`^Line 1$`),
			ret:  [][]int{{0, 6}},
			err:  nil,
		},
		{
			line: 1,
			reg:  regexp.MustCompile(`\d`),
			ret:  [][]int{{5, 6}},
			err:  nil,
		},
		{
			line: 2,
			reg:  regexp.MustCompile(`[a-z]+`),
			ret:  [][]int{{0, 3}, {4, 7}, {8, 11}},
			err:  nil,
		},
		{
			line: 3,
			reg:  regexp.MustCompile(`^Last line$`),
			ret:  [][]int{{0, 9}},
			err:  nil,
		},
	}

	for _, c := range searchCases {
		r := NewLineReader(bytes.NewReader([]byte(input)))

		ret, err := r.SearchLine(c.reg, c.line)
		if err != c.err {
			t.Errorf("SearchLine(%v, %d): err got %v want %v", c.reg, c.line, err, c.err)
		}

		if !reflect.DeepEqual(ret, c.ret) {
			t.Errorf("SearchLine(%v, %d) = %v want %v", c.reg, c.line, ret, c.ret)
		}
	}
}
