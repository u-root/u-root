// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestFalse(t *testing.T) {
	// Ensure 1 is returned.
	out, err := testutil.Command(t).CombinedOutput()
	if err := testutil.IsExitCode(err, 1); err != nil {
		t.Error(err)
	}
	if len(out) != 0 {
		t.Fatalf("Expected no output; got %#v", string(out))
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
