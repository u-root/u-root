// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"bytes"
	"testing"
)

func TestProgressReader(t *testing.T) {
	input := bytes.NewBufferString("01234567890123456789")
	stdout := &bytes.Buffer{}
	pr := ProgressReader{
		R:        input,
		Symbol:   "#",
		Interval: 4,
		W:        stdout,
	}

	// Read one byte at a time.
	output := make([]byte, 1)
	pr.Read(output)
	pr.Read(output)
	pr.Read(output)
	if len(stdout.Bytes()) != 0 {
		t.Errorf("found %q, but expected no bytes to be written", stdout)
	}
	pr.Read(output)
	if stdout.String() != "#" {
		t.Errorf("found %q, expected %q to be written", stdout.String(), "#")
	}

	// Read 9 bytes all at once.
	output = make([]byte, 9)
	pr.Read(output)
	if stdout.String() != "###" {
		t.Errorf("found %q, expected %q to be written", stdout.String(), "###")
	}
	if string(output) != "456789012" {
		t.Errorf("found %q, expected %q to be written", string(output), "456789012")
	}
}
