// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/u-root/u-root/integration/internal/gotest"
	"github.com/u-root/u-root/pkg/golang"
)

// testPkgs returns a slice of tests to run.
func testPkgs(t *testing.T) []string {
	// Packages which do not contain tests (or do not contain tests for the
	// build target) will still compile a test binary which vacuously pass.
	cmd := exec.Command("go", "list",
		"github.com/u-root/u-root/cmds/core/...",
		"github.com/u-root/u-root/cmds/boot/...",
		// TODO: only running tests in cmds because tests in pkg have
		// duplicate names which confuses the test runner. This should
		// get fixed.
		// "github.com/u-root/u-root/xcmds/...",
		// "github.com/u-root/u-root/pkg/...",
	)
	cmd.Env = append(os.Environ(), "GOARCH="+TestArch())
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
		"github.com/u-root/u-root/cmds/core/dhclient",
		"github.com/u-root/u-root/cmds/core/elvish/eval",
		"github.com/u-root/u-root/cmds/core/fusermount",
		"github.com/u-root/u-root/cmds/core/gpt",
		"github.com/u-root/u-root/cmds/core/kill",
		"github.com/u-root/u-root/cmds/core/mount",
		"github.com/u-root/u-root/cmds/core/tail",
		"github.com/u-root/u-root/cmds/core/wget",
		"github.com/u-root/u-root/cmds/core/which",
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

func copyRelativeFiles(src string, dst string) error {
	return filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if fi.Mode().IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), fi.Mode().Perm())
		} else if fi.Mode().IsRegular() {
			srcf, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcf.Close()
			dstf, err := os.Create(filepath.Join(dst, rel))
			if err != nil {
				return err
			}
			defer dstf.Close()
			_, err = io.Copy(dstf, srcf)
			return err
		}
		return nil
	})
}

// TestGoTest effectively runs "go test ./..." inside a QEMU instance. The
// tests run as root and can do all sorts of things not possible otherwise.
func TestGoTest(t *testing.T) {
	SkipWithoutQEMU(t)

	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "uroot-integration")
	if err != nil {
		t.Fatal(err)
	}

	env := golang.Default()
	env.CgoEnabled = false
	env.GOARCH = TestArch()

	// Statically build tests and add them to the temporary directory.
	pkgs := testPkgs(t)
	var tests []string
	os.Setenv("CGO_ENABLED", "0")
	testDir := filepath.Join(tmpDir, "tests")
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(testDir, pkg)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			t.Fatal(err)
		}

		testFile := filepath.Join(pkgDir, fmt.Sprintf("%s.test", path.Base(pkg)))
		cmd := exec.Command("go", "test", "-ldflags", "-s -w", "-c", pkg, "-o", testFile)
		if err := cmd.Run(); err != nil {
			t.Fatalf("could not build %s: %v", pkg, err)
		}

		// When a package does not contain any tests, the test
		// executable is not generated, so it is not included in the
		// `tests` list.
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			tests = append(tests, pkg)

			p, err := env.Package(pkg)
			if err != nil {
				t.Fatal(err)
			}
			// Optimistically copy any files in the pkg's
			// directory, in case e.g. a testdata dir is there.
			if err := copyRelativeFiles(p.Dir, filepath.Join(testDir, pkg)); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create the CPIO and start QEMU.
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/gotest/uinit",
			"github.com/u-root/u-root/cmds/core/init",
			// Used by gotest/uinit.
			"github.com/u-root/u-root/cmds/core/mkdir",
			"github.com/u-root/u-root/cmds/core/mount",
			// Used by an elvish test.
			"github.com/u-root/u-root/cmds/core/ls",
		},
		TmpDir: tmpDir,
	})
	defer cleanup()

	// Check that each test passed.
	gotest.WalkTests(testDir, func(i int, _ string, base string) {
		runMsg := fmt.Sprintf("TAP: # running %d - %s", i, base)
		passMsg := fmt.Sprintf("TAP: ok %d - %s", i, base)
		failMsg := fmt.Sprintf("TAP: not ok %d - %s", i, base)
		passOrFailMsg := regexp.MustCompile(fmt.Sprintf("TAP: (not )?ok %d - %s", i, base))

		t.Log(runMsg)
		str, err := q.ExpectRETimeout(passOrFailMsg, 30*time.Second)
		if err != nil {
			// If we can neither find the "ok" nor the "not ok" message, the
			// test runner inside QEMU is misbehaving and we fatal early
			// instead of wasting time.
			t.Logf(failMsg)
			t.Fatal("TAP: Bail out! QEMU test runner stopped printing.")
		}

		if strings.Contains(str, passMsg) {
			t.Log(passMsg)
		} else {
			t.Error(failMsg)
		}
	})
}
