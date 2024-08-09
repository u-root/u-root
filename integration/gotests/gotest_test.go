// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/govmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/mkuimage/uimage"
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
	cmd.Env = append(os.Environ(), "GOARCH="+string(qemu.GuestArch()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	allPkgs := strings.Fields(strings.TrimSpace(string(out)))

	// TODO: Some tests do not run properly in QEMU at the moment. They are
	// blocklisted. These tests fail for mostly two reasons:
	// 1. either it requires networking (not enabled in the kernel)
	// 2. or it depends on some test files (for example /bin/sleep)
	blocklist := []string{
		"github.com/u-root/u-root/cmds/core/cmp",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/fusermount",
		"github.com/u-root/u-root/cmds/core/gosh",
		"github.com/u-root/u-root/cmds/core/wget",
		"github.com/u-root/u-root/cmds/core/netcat",
		"github.com/u-root/u-root/cmds/core/which",
		// Some of TestEdCommands do not exit properly and end up left running. No idea how to fix this yet.
		"github.com/u-root/u-root/cmds/exp/ed",
		"github.com/u-root/u-root/cmds/exp/pox",
		"github.com/u-root/u-root/pkg/crypto",
		"github.com/u-root/u-root/pkg/tarutil",
		"github.com/u-root/u-root/pkg/ldd",

		// These have special configuration.
		"github.com/u-root/u-root/pkg/gpio",
		"github.com/u-root/u-root/pkg/mount",
		"github.com/u-root/u-root/pkg/mount/block",
		"github.com/u-root/u-root/pkg/mount/loop",
		"github.com/u-root/u-root/pkg/ipmi",
		"github.com/u-root/u-root/pkg/smbios",

		// Missing xzcat in VM.
		"github.com/u-root/u-root/cmds/exp/bzimage",
		"github.com/u-root/u-root/pkg/boot/bzimage",

		// ??
		"github.com/u-root/u-root/pkg/tss",
		"github.com/u-root/u-root/pkg/syscallfilter",
	}
	switch qemu.GuestArch() {
	case qemu.ArchArm64:
		blocklist = append(blocklist,
			"github.com/u-root/u-root/pkg/strace",

			// These tests run in 1-2 seconds on x86, but run
			// beyond their huge timeout under arm64 in the VM. Not
			// sure why. Slow emulation?
			"github.com/u-root/u-root/cmds/core/pci",
			"github.com/u-root/u-root/cmds/exp/cbmem",
			"github.com/u-root/u-root/pkg/vfile",
		)

	case qemu.ArchArm:
		blocklist = append(blocklist,
			"github.com/u-root/u-root/cmds/exp/cbmem",

			// These 4 tests do not compile on arm.
			"github.com/u-root/u-root/pkg/boot/kexec",
			"github.com/u-root/u-root/pkg/flash/chips",
			"github.com/u-root/u-root/pkg/mount/gpt",
			"github.com/u-root/u-root/pkg/mount/mtd",
		)
	}

	var pkgs []string
	for _, p := range allPkgs {
		if !slices.Contains(blocklist, p) {
			pkgs = append(pkgs, p)
		}
	}
	return pkgs
}

// TestGoTest effectively runs "go test ./..." inside a QEMU instance. The
// tests run as root and can do all sorts of things not possible otherwise.
func TestGoTest(t *testing.T) {
	pkgs := testPkgs(t)

	govmtest.Run(t, "vm",
		govmtest.WithPackageToTest(pkgs...),
		govmtest.WithUimage(
			uimage.WithShell("gosh"),
			uimage.WithBusyboxCommands("github.com/u-root/u-root/cmds/core/*"),
			uimage.WithFiles(
				"/etc/group",
				"/etc/passwd",
			),
		),
		govmtest.WithQEMUFn(
			qemu.WithVMTimeout(15*time.Minute),
			qemu.VirtioRandom(),

			// Bump this up so that some unit tests can happily
			// and questionably pre-claim large bytes slices.
			//
			// e.g. pkg/mount/gpt/gpt_test.go need to claim 4.29G
			//
			//     disk = make([]byte, 0x100000000)
			qemu.IfNotArch(qemu.ArchArm, qemu.ArbitraryArgs("-m", "6G")),
			qemu.IfArch(qemu.ArchArm, qemu.ArbitraryArgs("-m", "3G")),

			// aarch64 VMs start at 1970-01-01 without RTC explicitly set.
			qemu.ArbitraryArgs("-rtc", "base=localtime,clock=vm"),
		),
	)
}
