// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type nopCloser struct {
	*bytes.Buffer
}

func (nopCloser) Close() error { return nil }

func TestTransformCopy(t *testing.T) {
	cfg := config{
		transforms: transforms{
			transform{from: "old", to: "new"},
		},
		inplace: false,
	}

	input := strings.NewReader("old\n")
	output := &nopCloser{bytes.NewBuffer(nil)}

	readStreams := []io.ReadCloser{io.NopCloser(input)}
	writeStreams := []io.WriteCloser{output}

	transformCopy(cfg, readStreams, writeStreams)

	result := output.String()
	if strings.TrimSpace(result) != "new" {
		t.Errorf("Expected 'new', got '%s'", result)
	}
}
