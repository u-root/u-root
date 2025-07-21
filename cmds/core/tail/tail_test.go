// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func tempFile(t *testing.T, content string) *os.File {
	t.Helper()
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "tailRunTest")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	return f
}

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
	input := "a\nb\nc\n"
	f := tempFile(t, input)

	var b bytes.Buffer
	err := run(os.Stdin, &b, false, 10, time.Second, []string{f.Name()})
	if err != nil {
		t.Error(err)
	}

	if b.String() != input {
		t.Errorf("tail output does not match, want %q, got %q", input, b.String())
	}

	err = run(nil, nil, true, 10, time.Second, []string{"a", "b"})
	if err == nil {
		t.Error("tail should return an error if more than one file specified if follow true")
	}

	b.Truncate(0)
	err = run(f, &b, false, -1, time.Second, nil)
	if err != nil {
		t.Error(err)
	}

	if b.String() != "c\n" {
		t.Errorf("tail output does not match, want %q, got %q", input, b.String())
	}
}

func TestTailMultipleFiles(t *testing.T) {
	f1 := tempFile(t, "f1")
	f2 := tempFile(t, "f2")

	var b bytes.Buffer
	err := run(os.Stdin, &b, false, 10, time.Second, []string{
		f1.Name(), f2.Name(),
	})
	if err != nil {
		t.Error(err)
	}

	ex := fmt.Sprintf("==> %s <==\nf1\n==> %s <==\nf2", f1.Name(), f2.Name())
	if ex != b.String() {
		t.Errorf("%v != %v", ex, b.String())
	}
}

type syncWriter struct {
	ch chan []byte
}

func (sw *syncWriter) Write(b []byte) (int, error) {
	// ignore writes with len zero
	if len(b) == 0 {
		return 0, nil
	}
	sw.ch <- b
	return len(b), nil
}

func TestTailFollow(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "follow")
	if err != nil {
		t.Fatalf("can't create temp file: %v", err)
	}
	defer f.Close()

	sw := &syncWriter{
		ch: make(chan []byte),
	}

	go func() {
		run(f, sw, true, 10, 100*time.Millisecond, nil)
	}()
	ff, err := os.OpenFile(f.Name(), os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("can't open temp file: %v", err)
	}

	// wait a bit before writting to file
	time.Sleep(300 * time.Millisecond)

	firstLine := []byte("hello\n")

	_, err = ff.Write(firstLine)
	if err != nil {
		t.Fatalf("can't write to file: %v", err)
	}

	ff.Sync()

	r1 := <-sw.ch
	if !bytes.Equal(r1, firstLine) {
		t.Fatalf("expected %q, got %q", string(firstLine), string(r1))
	}
}

func TestLastNLines(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
		n      int
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
