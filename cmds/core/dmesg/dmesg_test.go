// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestDmesg(t *testing.T) {
	if os.Getenv("IS_CIRCLECI") == "1" {
		t.Skipf("test skipped on circleci, which doesn't allow access to dmesg at all")
	}

	cmd := testutil.Command(t)
	out, err := cmd.CombinedOutput()
	if err != nil || len(out) == 0 {
		t.Logf("Output: %s", string(out))
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
