// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
)

func TestScript(t *testing.T) {
	script := "echo HELLO WORLD"
	vm := vmtest.StartVMAndRunCmds(t, script, vmtest.WithQEMUFn(qemu.WithVMTimeout(30*time.Second)))
	if _, err := vm.Console.ExpectString("HELLO WORLD"); err != nil {
		t.Errorf("Want HELLO WORLD: %v", err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
