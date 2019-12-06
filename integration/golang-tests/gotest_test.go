// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

var (
	kernelPath = flag.String("kernel", "", "path to the Linux kernel binary to use for tests")
	qemuPath   = flag.String("qemu", "", "path to the QEMU binary to use for tests")
	testarch   = flag.String("testarch", "", "name of the architecture to use for tests")
)

// testPkgs returns a slice of tests to run.
func testPkgs(t *testing.T) []string {
	// Packages which do not contain tests (or do not contain tests for the
	// build target) will still compile a test binary which vacuously pass.
	cmd := exec.Command("go", "list",
		"github.com/u-root/u-root/cmds/boot/...",
		"github.com/u-root/u-root/cmds/core/...",
		"github.com/u-root/u-root/cmds/exp/...",
		"github.com/u-root/u-root/pkg/...",
	)
	cmd.Env = append(os.Environ(), "GOARCH="+vmtest.TestArch())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	pkgs := strings.Fields(strings.TrimSpace(string(out)))

	// TODO: Some tests do not run properly in QEMU at the moment. They are
	// blacklisted. These tests fail for mostly two reasons:
	// 1. either it requires networking (not enabled in the kernel)
	// 2. or it depends on some test files (for example /bin/sleep)
	blacklist := []string{
		"github.com/u-root/u-root/cmds/core/cmp",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/elvish/eval",
		"github.com/u-root/u-root/cmds/core/elvish/edit/tty",
		"github.com/u-root/u-root/cmds/core/fusermount",
		"github.com/u-root/u-root/cmds/core/wget",
		"github.com/u-root/u-root/cmds/core/which",
		"github.com/u-root/u-root/cmds/exp/rush",
		"github.com/u-root/u-root/cmds/exp/pox",
		"github.com/u-root/u-root/pkg/crypto",
		"github.com/u-root/u-root/pkg/tarutil",
		"github.com/u-root/u-root/pkg/ldd",
		"github.com/u-root/u-root/pkg/loop",
		"github.com/u-root/u-root/pkg/gpio",

		// Missing xzcat in VM.
		"github.com/u-root/u-root/cmds/exp/bzimage",
		"github.com/u-root/u-root/pkg/bzimage",

		// Missing /dev/mem and /sys/firmware/efi
		"github.com/u-root/u-root/pkg/boot/acpi",

		// No Go compiler in VM.
		"github.com/u-root/u-root/pkg/bb",
		"github.com/u-root/u-root/pkg/uroot",
		"github.com/u-root/u-root/pkg/uroot/builder",
	}
	for i := 0; i < len(pkgs); i++ {
		for _, b := range blacklist {
			if pkgs[i] == b {
				pkgs = append(pkgs[:i], pkgs[i+1:]...)
			}
		}
	}

	return pkgs
}

// TestGoTest effectively runs "go test ./..." inside a QEMU instance. The
// tests run as root and can do all sorts of things not possible otherwise.
func TestGoTest(t *testing.T) {
	if len(*kernelPath) > 0 {
		os.Setenv("UROOT_KERNEL", *kernelPath)
	}
	if len(*qemuPath) > 0 {
		os.Setenv("UROOT_QEMU", *qemuPath)
	}
	if len(*testarch) > 0 {
		os.Setenv("UROOT_TESTARCH", *testarch)
	}

	pkgs := testPkgs(t)

	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Timeout: 120 * time.Second,
		},
	}
	vmtest.GolangTest(t, pkgs, o)
}
