// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestDmesg(t *testing.T) {
	cmd := testutil.Command(t)
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		t.Fatalf("Error: %v, Output: %v", err, string(out))
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
