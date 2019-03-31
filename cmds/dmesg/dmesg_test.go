// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// TODO(https://github.com/u-root/u-root/issues/1160): This test has been disabled
func testDmesg(t *testing.T) {
	cmd := testutil.Command(t)
	out, err := cmd.CombinedOutput()
	if err != nil || len(out) == 0 {
		t.Fatalf("Error: %v, Output: %v", err, string(out))
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
