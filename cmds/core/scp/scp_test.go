// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

func TestScpSource(t *testing.T) {
	var w bytes.Buffer
	var r bytes.Buffer

	tf, err := os.CreateTemp("", "TestScpSource")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(tf.Name())
	tf.Write([]byte("test-file-contents"))

	r.Write([]byte{0})
	err = scpSource(&w, &r, tf.Name())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	expected := fmt.Appendf(nil, "C0600 18 %s\ntest-file-contents", path.Base(tf.Name()))
	expected = append(expected, 0)
	if string(expected) != w.String() {
		t.Fatalf("Got: %v\nExpected: %v", w.String(), string(expected))
	}
}

func TestScpSink(t *testing.T) {
	var w bytes.Buffer
	var r bytes.Buffer

	tf, err := os.CreateTemp("", "TestScpSink")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	r.Write([]byte("C0600 18 test\ntest-file-contents"))
	// Post IO-copy success status
	r.Write([]byte{0})

	err = scpSink(&w, &r, tf.Name())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	// 1: Initial SUCCESS post to start scp
	// 2: Success opening file tf.Name()
	// 3: Success writing file
	expected := []byte{0, 0, 0}
	if string(expected) != w.String() {
		t.Fatalf("Got: %v\nExpected: %v", w.Bytes(), expected)
	}

	m := make([]byte, 18)
	n, err := tf.Read(m)
	if err != nil {
		t.Fatalf("IO error: %v", err)
	}
	if n != 18 {
		t.Fatalf("Expected 18 bytes, got %v", n)
	}

	// Ensure EOF
	_, err = tf.Read(m)
	if err != io.EOF {
		t.Fatalf("Expected EOF, got %v", err)
	}

	if string(m) != "test-file-contents" {
		t.Fatalf("Expected 'test-file-contents', got '%v'", string(m))
	}
}
