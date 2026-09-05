// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"bytes"
	"testing"
)

func TestRawWriter_NoNewlines(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("hello world")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	if buf.String() != "hello world" {
		t.Errorf("Output was %q, expected %q", buf.String(), "hello world")
	}
}

func TestRawWriter_SingleNewline(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("hello\nworld")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	expected := "hello\r\nworld"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}

func TestRawWriter_MultipleNewlines(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("line1\nline2\nline3\n")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	expected := "line1\r\nline2\r\nline3\r\n"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}

func TestRawWriter_ConsecutiveNewlines(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("hello\n\n\nworld")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	expected := "hello\r\n\r\n\r\nworld"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}

func TestRawWriter_OnlyNewline(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("\n")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	expected := "\r\n"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}

func TestRawWriter_EmptyWrite(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	input := []byte("")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Write() returned %d, expected 0", n)
	}
	if buf.Len() != 0 {
		t.Errorf("Buffer has %d bytes, expected 0", buf.Len())
	}
}

func TestRawWriter_PreservesCarriageReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	// Test that existing \r\n sequences are not double-translated
	input := []byte("hello\r\nworld")
	n, err := rw.Write(input)
	if err != nil {
		t.Errorf("Write() failed: %v", err)
	}
	if n != len(input) {
		t.Errorf("Write() returned %d, expected %d", n, len(input))
	}
	// Should translate the \n even though \r precedes it
	expected := "hello\r\r\nworld"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}

func TestRawWriter_MultipleWrites(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := NewRawWriter(buf)

	// Write in multiple chunks
	writes := []string{"hello\n", "world\n", "test"}
	for _, w := range writes {
		_, err := rw.Write([]byte(w))
		if err != nil {
			t.Errorf("Write() failed: %v", err)
		}
	}

	expected := "hello\r\nworld\r\ntest"
	if buf.String() != expected {
		t.Errorf("Output was %q, expected %q", buf.String(), expected)
	}
}
