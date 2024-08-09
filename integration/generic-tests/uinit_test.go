// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qcoverage"
	"github.com/hugelgupf/vmtest/qemu/quimage"
	"github.com/u-root/mkuimage/uimage"
)

// TestHelloWorld runs an init which prints the string "HELLO WORLD" and exits.
func TestHelloWorld(t *testing.T) {
	vm := qemu.StartT(t, "vm", qemu.ArchUseEnvv,
		quimage.WithUimageT(t,
			uimage.WithInit("init"),
			uimage.WithUinit("helloworld"),
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/integration/testcmd/helloworld",
				"github.com/u-root/u-root/cmds/core/init",
			),
		),
		qemu.WithVMTimeout(time.Minute),
		qcoverage.CollectKernelCoverage(t),
	)

	if _, err := vm.Console.ExpectString("HELLO WORLD"); err != nil {
		t.Error(`expected "HELLO WORLD", got error: `, err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

// TestHelloWorldNegative runs an init which does not print the string "GOODBYE WORLD".
func TestHelloWorldNegative(t *testing.T) {
	vm := qemu.StartT(t, "vm", qemu.ArchUseEnvv,
		quimage.WithUimageT(t,
			uimage.WithInit("init"),
			uimage.WithUinit("helloworld"),
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/integration/testcmd/helloworld",
				"github.com/u-root/u-root/cmds/core/init",
			),
		),
		qemu.WithVMTimeout(time.Minute),
		qcoverage.CollectKernelCoverage(t),
	)

	if _, err := vm.Console.ExpectString("GOODBYE WORLD"); err == nil {
		t.Error(`expected error, but matched "GOODBYE WORLD"`)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
