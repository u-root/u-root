// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestHelloWorld runs an init which prints the string "HELLO WORLD" and exits.
func RunTestHelloWorld(t *testing.T, initramfs string) {
	if len(initramfs) == 0 {
		f, err := vmtest.CreateTestInitramfs(
			uroot.Opts{}, "github.com/u-root/u-root/integration/testcmd/helloworld/uinit", "")
		if err != nil {
			t.Errorf("failed to create test initramfs: %v", err)
		}
		defer os.Remove(f)
		initramfs = f
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Initramfs: initramfs,
	})
	defer cleanup()

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(`expected "HELLO WORLD", got error: `, err)
	}
}

// TestHelloWorldNegative runs an init which does not print the string "HELLO WORLD".
func RunTestHelloWorldNegative(t *testing.T, initramfs string) {
	if len(initramfs) == 0 {
		f, err := vmtest.CreateTestInitramfs(
			uroot.Opts{}, "github.com/u-root/u-root/integration/testcmd/helloworld/uinit", "")
		if err != nil {
			t.Errorf("failed to create test initramfs: %v", err)
		}
		defer os.Remove(f)
		initramfs = f
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Initramfs: initramfs,
	})
	defer cleanup()

	if err := q.Expect("GOODBYE WORLD"); err == nil {
		t.Fatal(`expected error, but matched "GOODBYE WORLD"`)
	}

}

func RunTestScript(t *testing.T, initramfs string) {
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "ShellScript",
		TestCmds: []string{
			"echo HELLO WORLD",
			"shutdown -h",
		},
		Initramfs: initramfs,
	})
	defer cleanup()

	if err := q.Expect("HELLO WORLD"); err != nil {
		t.Fatal(`expected "HELLO WORLD", got error: `, err)
	}
}
