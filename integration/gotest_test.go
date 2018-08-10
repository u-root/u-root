// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// testPkgs returns a slice of tests to run.
func testPkgs(t *testing.T) []string {
	// Packages which do not contain tests (or do not contain tests for the
	// build target) will still compile a test binary which vacuously pass.
	cmd := exec.Command("go", "list",
		"github.com/u-root/u-root/cmds/...",
		// TODO: only running tests in cmds because tests in pkg have
		// duplicate names which confuses the test runner. This should
		// get fixed.
		// "github.com/u-root/u-root/xcmds/...",
		// "github.com/u-root/u-root/pkg/...",
	)
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
		"github.com/u-root/u-root/cmds/cmp",
		"github.com/u-root/u-root/cmds/dhclient",
		"github.com/u-root/u-root/cmds/fmap",
		"github.com/u-root/u-root/cmds/kill",
		"github.com/u-root/u-root/cmds/tail",
		"github.com/u-root/u-root/cmds/wget",
		"github.com/u-root/u-root/cmds/which",
		"github.com/u-root/u-root/cmds/wifi",
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
	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "uroot-integration")
	if err != nil {
		t.Fatal(err)
	}

	// Statically build tests and add them to the temporary directory.
	pkgs := testPkgs(t)
	os.Setenv("CGO_ENABLED", "0")
	testDir := filepath.Join(tmpDir, "tests")
	for _, pkg := range pkgs {
		testFile := filepath.Join(testDir, path.Base(pkg))
		cmd := exec.Command("go", "test", "-ldflags", "-s -w", "-c", pkg, "-o", testFile)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("could not build %s: %v", pkg, err)
		}
	}

	// Create the CPIO and start QEMU.
	tmpDir, q := testWithQEMU(t, options{
		uinitName: "gotest",
		tmpDir:    tmpDir,
	})
	defer cleanup(t, tmpDir, q)

	// Check that each test passed.
	bases := []string{}
	for _, pkg := range pkgs {
		bases = append(bases, path.Base(pkg))
	}
	sort.Strings(bases) // Tests are run and checked in sorted order.
	t.Log("TAP: TAP version 12")
	t.Logf("TAP: 1..%d", len(bases))
	for i, base := range bases {
		runMsg := fmt.Sprintf("TAP: # running %d - %s", i, base)
		passMsg := fmt.Sprintf("TAP: ok %d - %s", i, base)
		failMsg := fmt.Sprintf("TAP: not ok %d - %s", i, base)

		t.Log(runMsg)
		if err := q.ExpectTimeout(passMsg, 2*time.Second); err == nil {
			t.Logf(passMsg)
		} else {
			t.Errorf(failMsg)
		}
	}
}
