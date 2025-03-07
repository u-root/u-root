// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || (linux && !arm && !386 && !mips && !mipsle)

package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		errCode := run(io.Discard, io.Discard, "./stat.go")
		if errCode != 0 {
			t.Errorf("expected 0 got %d", errCode)
		}
	})

	t.Run("file not exists", func(t *testing.T) {
		buf := &bytes.Buffer{}
		errCode := run(io.Discard, buf, "filenotexists")
		if errCode != 1 {
			t.Errorf("expected 1 got %d", errCode)
		}

		stderr := buf.String()
		if !strings.Contains(stderr, "filenotexists") {
			t.Errorf("expected filenotexists in stderr, got %q", stderr)
		}
	})
}
