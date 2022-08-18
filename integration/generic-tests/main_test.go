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
	if len(os.Getenv("UROOT_KERNEL")) == 0 {
		log.Fatalf("Failed to run tests: no kernel provide: source integration/GET_KERNEL_QEMU to get a kernel")
	}
	if len(os.Getenv("UROOT_QEMU")) == 0 {
		log.Fatalf("Failed to run tests: no QEMU binary provided: source integration/GET_KERNEL_QEMU to get a qemu")
	}

	log.Printf("Starting generic tests...")

	os.Exit(m.Run())
}
