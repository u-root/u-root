// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"testing"
)

func TestLSFD(t *testing.T) {
	t.Run("all files", func(t *testing.T) {
		err := run(io.Discard, proc, nil)
		if err != nil {
			t.Errorf("expected nil got %v", err)
		}
	})

	t.Run("one file", func(t *testing.T) {
		err := run(io.Discard, proc, []string{"1"})
		if err != nil {
			t.Errorf("expected nil got %v", err)
		}
	})
}
