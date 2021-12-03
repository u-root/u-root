// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"io"
	"testing"
)

func TestProgressReadCloser(t *testing.T) {
	input := io.NopCloser(bytes.NewBufferString("01234567890123456789"))
	stdout := &bytes.Buffer{}
	prc := ProgressReadCloser{
		RC:       input,
		Symbol:   "#",
		Interval: 4,
		W:        stdout,
	}

	// Read one byte at a time.
	output := make([]byte, 1)
	prc.Read(output)
	prc.Read(output)
	prc.Read(output)
	if len(stdout.Bytes()) != 0 {
		t.Errorf("found %q, but expected no bytes to be written", stdout)
	}
	prc.Read(output)
	if stdout.String() != "#" {
		t.Errorf("found %q, expected %q to be written", stdout.String(), "#")
	}

	// Read 9 bytes all at once.
	output = make([]byte, 9)
	prc.Read(output)
	if stdout.String() != "###" {
		t.Errorf("found %q, expected %q to be written", stdout.String(), "###")
	}
	if string(output) != "456789012" {
		t.Errorf("found %q, expected %q to be written", string(output), "456789012")
	}

	// Read until EOF
	output, err := io.ReadAll(&prc)
	if err != nil {
		t.Errorf("got %v, expected nil error", err)
	}
	if stdout.String() != "#####\n" {
		t.Errorf("found %q, expected %q to be written", stdout.String(), "#####\n")
	}
	if string(output) != "3456789" {
		t.Errorf("found %q, expected %q to be written", string(output), "3456789")
	}

	err = prc.Close()
	if err != nil {
		t.Errorf("got %v, expected nil error", err)
	}
}
