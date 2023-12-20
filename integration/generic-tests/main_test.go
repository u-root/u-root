// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if len(os.Getenv("VMTEST_KERNEL")) == 0 {
		log.Fatalf("Failed to run tests: no kernel provided: use runvmtest tool to supply it")
	}
	if len(os.Getenv("VMTEST_QEMU")) == 0 {
		log.Fatalf("Failed to run tests: no QEMU binary provided: use runvmtest tool to supply it")
	}

	log.Printf("Starting generic tests...")

	os.Exit(m.Run())
}
