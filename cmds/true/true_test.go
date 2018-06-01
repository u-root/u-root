// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// Ensure 0 is returned.
func TestTrue(t *testing.T) {
	out, err := testutil.Command(t).CombinedOutput()
	if err != nil || len(out) != 0 {
		t.Fatal("Expected no output and no error")
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
