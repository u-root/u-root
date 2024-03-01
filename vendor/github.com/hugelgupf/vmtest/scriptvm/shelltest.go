// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scriptvm is an API to run shell scripts in a VM guest.
package scriptvm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/qemu/qcoverage"
	"github.com/hugelgupf/vmtest/qemu/quimage"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/mkuimage/uimage"
)

// Options are QEMU VM integration test options.
type Options struct {
	// QEMUOpts are options to the QEMU VM.
	QEMUOpts []qemu.Fn

	// Initramfs is an optional u-root initramfs to build.
	Initramfs []uimage.Modifier
}

// Modifier is used to configure a VM.
type Modifier func(testing.TB, *Options) error

// WithQEMUFn adds QEMU options.
func WithQEMUFn(fn ...qemu.Fn) Modifier {
	return func(_ testing.TB, v *Options) error {
		v.QEMUOpts = append(v.QEMUOpts, fn...)
		return nil
	}
}

// WithUimage merges o with already appended initramfs build options.
func WithUimage(mods ...uimage.Modifier) Modifier {
	return func(_ testing.TB, v *Options) error {
		v.Initramfs = append(v.Initramfs, mods...)
		return nil
	}
}

// Run starts a VM and runs the given script using gosh in the guest.
//
// gosh is based on mvdan.cc/sh and strives to be bash-compatible.
//
// If any command fails, the test fails.
//
//   - TODO: timeouts for individual individual commands.
func Run(t testing.TB, name, script string, mods ...Modifier) {
	vm := Start(t, name, script, mods...)

	if _, err := vm.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("Waiting for 'TESTS PASSED MARKER' failed -- script likely failed: %v", err)
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("VM exited with %v", err)
	}
}

// Start starts a VM and runs the script using gosh in the guest.
// If the commands return, the VM will be shutdown.
func Start(t testing.TB, name, script string, mods ...Modifier) *qemu.VM {
	qemu.SkipWithoutQEMU(t)

	o := &Options{}
	for _, mod := range mods {
		if mod != nil {
			if err := mod(t, o); err != nil {
				t.Fatal(err)
			}
		}
	}

	sharedDir := testtmp.TempDir(t)

	// Generate gosh shell script of test commands in o.SharedDir.
	if len(script) > 0 {
		testFile := filepath.Join(sharedDir, "test.sh")
		if err := os.WriteFile(testFile, []byte(strings.Join([]string{"set -ex", script}, "\n")), 0o777); err != nil {
			t.Fatal(err)
		}
	}

	initramfs := append([]uimage.Modifier{
		uimage.WithBusyboxCommands(
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/gosh",
			"github.com/hugelgupf/vmtest/vminit/shutdownafter",
			"github.com/hugelgupf/vmtest/vminit/vmmount",
			"github.com/hugelgupf/vmtest/vminit/shelluinit",
		),
		uimage.WithInit("init"),
		uimage.WithUinit("shutdownafter", "--", "vmmount", "--", "shelluinit"),
	}, o.Initramfs...)

	qopts := []qemu.Fn{
		quimage.WithUimageT(t, initramfs...),
		qemu.P9Directory(sharedDir, "shelltest"),
		qcoverage.CollectKernelCoverage(t),
		qcoverage.ShareGOCOVERDIR(),
		qemu.WithVmtestIdent(),
	}

	// Prepend our default options so user-supplied o.QEMUOpts supersede.
	return qemu.StartT(t, name, qemu.ArchUseEnvv, append(qopts, o.QEMUOpts...)...)
}
