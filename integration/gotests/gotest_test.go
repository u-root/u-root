// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
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
	// blocklisted. These tests fail for mostly two reasons:
	// 1. either it requires networking (not enabled in the kernel)
	// 2. or it depends on some test files (for example /bin/sleep)
	blocklist := []string{
		"github.com/u-root/u-root/cmds/core/cmp",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/fusermount",
		"github.com/u-root/u-root/cmds/core/wget",
		"github.com/u-root/u-root/cmds/core/which",
		"github.com/u-root/u-root/cmds/exp/gosh",
		"github.com/u-root/u-root/cmds/exp/pox",
		"github.com/u-root/u-root/pkg/crypto",
		"github.com/u-root/u-root/pkg/tarutil",
		"github.com/u-root/u-root/pkg/ldd",

		// These have special configuration.
		"github.com/u-root/u-root/pkg/gpio",
		"github.com/u-root/u-root/pkg/mount",
		"github.com/u-root/u-root/pkg/mount/loop",

		// Missing xzcat in VM.
		"github.com/u-root/u-root/cmds/exp/bzimage",
		"github.com/u-root/u-root/pkg/boot/bzimage",

		// No Go compiler in VM.
		"github.com/u-root/u-root/pkg/bb",
		"github.com/u-root/u-root/pkg/uroot",
		"github.com/u-root/u-root/pkg/uroot/builder",

		// ??
		"github.com/u-root/u-root/pkg/tss",
	}
	if vmtest.TestArch() == "arm64" {
		blocklist = append(
			blocklist,
			"github.com/u-root/u-root/pkg/strace",
		)
	}
	for i := 0; i < len(pkgs); i++ {
		for _, b := range blocklist {
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
	pkgs := testPkgs(t)

	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Timeout: 120 * time.Second,
			Devices: []qemu.Device{
				// Bump this up so that some unit tests can happily
				// and questionably pre-claim large bytes slices.
				//
				// e.g. pkg/mount/gpt/gpt_test.go need to claim 4.29G
				//
				//     disk = make([]byte, 0x100000000)
				qemu.ArbitraryArgs{"-m", "6G"},
			},
		},
		BuildOpts: uroot.Opts{
			ExtraFiles: []string{
				"/etc/group",
			},
		},
	}
	vmtest.GolangTest(t, pkgs, o)
}
