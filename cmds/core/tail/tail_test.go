// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestTailReadBackwards(t *testing.T) {
	input, err := os.Open("./testdata/read_backwards.txt")
	if err != nil {
		t.Error(err)
	}
	output := &bytes.Buffer{}
	err = readLastLinesBackwards(input, output, 2)
	if err != nil {
		t.Error(err)
	}
	expected := []byte("second\nthird\n")
	got := output.Bytes()
	if !bytes.Equal(got, expected) {
		t.Fatalf("Invalid result reading backwards. Got %v; want %v", got, expected)
	}
	// try reading more, which should return EOF
	buf := make([]byte, 16)
	n, err := input.Read(buf)
	if err == nil {
		t.Fatalf("Expected EOF, got more bytes instead: %v", string(buf[:n]))
	}
	if err != io.EOF {
		t.Fatalf("Expected EOF, got another error instead: %v", err)
	}
}

func TestTailReadFromBeginning(t *testing.T) {
	input, err := os.Open("./testdata/read_from_beginning.txt")
	if err != nil {
		t.Error(err)
	}
	output := &bytes.Buffer{}
	err = readLastLinesFromBeginning(input, output, 3)
	if err != nil {
		t.Error(err)
	}
	expected := []byte("eight\nnine\nten\n")
	got := make([]byte, 4096) // anything larger than the expected result
	n, err := output.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got[:n], expected) {
		t.Fatalf("Invalid data while reading from the beginning. Got %v; want %v", string(got[:n]), string(expected))
	}
	// try reading more, which should return EOF
	buf := make([]byte, 16)
	n, err = input.Read(buf)
	if err == nil {
		t.Fatalf("Expected EOF, got more bytes instead: %v", string(buf[:n]))
	}
	if err != io.EOF {
		t.Fatalf("Expected EOF, got another error instead: %v", err)
	}
}

func TestTailRun(t *testing.T) {
	f, err := os.CreateTemp("", "tailRunTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	input := "a\nb\nc\n"
	_, err = f.WriteString(input)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	err = run(os.Stdin, &b, false, 10, []string{f.Name()})
	if err != nil {
		t.Error(err)
	}

	if b.String() != input {
		t.Errorf("tail output does not match, want %q, got %q", input, b.String())
	}

	err = run(nil, nil, false, 10, []string{"a", "b"})
	if err == nil {
		t.Error("tail should return an error if more than one file specified")
	}

	b.Truncate(0)
	err = run(f, &b, false, -1, nil)
	if err != nil {
		t.Error(err)
	}

	if b.String() != "c\n" {
		t.Errorf("tail output does not match, want %q, got %q", input, b.String())
	}
}

func TestLastNLines(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
		n      uint
	}{
		{
			input:  []byte{'a', '\n', '\n', 'b', '\n'},
			output: []byte{'a', '\n', '\n', 'b', '\n'},
			n:      4,
		},
		{
			input:  []byte{'a', '\n', '\n', 'b', '\n'},
			output: []byte{'a', '\n', '\n', 'b', '\n'},
			n:      3,
		},
		{
			input:  []byte{'a', '\n', '\n', 'b', '\n'},
			output: []byte{'\n', 'b', '\n'},
			n:      2,
		},
		{
			input:  []byte{'a', '\n', '\n', 'b', '\n'},
			output: []byte{'b', '\n'},
			n:      1,
		},
		{
			input:  []byte{'a', '\n', 'b', '\n', 'c', '\n'},
			output: []byte{'c', '\n'},
			n:      1,
		},
		{
			input:  []byte{'a', '\n', 'b', '\n', 'c'},
			output: []byte{'c'},
			n:      1,
		},
		{
			input:  []byte{'\n'},
			output: []byte{'\n'},
			n:      1,
		},
	}

	for _, test := range tests {
		r := lastNLines(test.input, test.n)
		if !bytes.Equal(r, test.output) {
			t.Errorf("want: %q, got: %q", string(test.output), string(r))
		}
	}
}
