// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func TestTTY(t *testing.T) {
	stdout := &bytes.Buffer{}
	err := run(stdout)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
}
