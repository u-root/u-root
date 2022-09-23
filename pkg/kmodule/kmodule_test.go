// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux

package kmodule

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	l, err := New()
	if err != nil {
		t.Skipf("New(): got %v, want nil; skipping tests as modules not enabled", err)
	}
	// At the very minimum, if we have it, we have list.
	var out = &bytes.Buffer{}
	if _, err := io.Copy(out, l); err != nil {
		t.Fatalf("List: got %v, want nil", err)
	}
	t.Logf("List: %s", out.String())
}
