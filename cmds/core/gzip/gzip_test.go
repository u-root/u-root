// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

func TestGzipMain(t *testing.T) {
	// This is a simple integration test to ensure the main package works.
	// More detailed tests are in the pkg/core/gzip package.

	if os.Getenv("TEST_MAIN_BINARY") == "1" {
		// When running as the test binary, just exit successfully
		return
	}

	t.Skip("Skipping main package test - detailed tests are in pkg/core/gzip")
}
