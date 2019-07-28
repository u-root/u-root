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
	"strings"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/json2test"
	"github.com/u-root/u-root/pkg/uio"
)

// testPkgs returns a slice of tests to run.
func testPkgs(t *testing.T) []string {
	// Packages which do not contain tests (or do not contain tests for the
	// build target) will still compile a test binary which vacuously pass.
	cmd := exec.Command("go", "list",
		"github.com/u-root/u-root/cmds/core/...",
		"github.com/u-root/u-root/cmds/boot/...",
		"github.com/u-root/u-root/pkg/...",
		"github.com/u-root/u-root/cmds/exp/...",
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
		//"github.com/u-root/u-root/pkg/pty",

		// Missing xzcat in VM.
		"github.com/u-root/u-root/cmds/exp/bzimage",
		"github.com/u-root/u-root/pkg/bzimage",

		// Missing /dev/mem and /sys/firmware/efi
		"github.com/u-root/u-root/pkg/acpi",

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

	tc := json2test.NewTestCollector()

	// Create the CPIO and start QEMU.
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/gotest/uinit",
			"github.com/u-root/u-root/cmds/core/init",
			// Used by different tests.
			"github.com/u-root/u-root/cmds/core/ls",
			"github.com/u-root/u-root/cmds/core/sleep",
			"github.com/u-root/u-root/cmds/core/echo",
		},
		TmpDir: tmpDir,
		SerialOutput: uio.ClosingMultiWriter(
			// Collect JSON test events in tc.
			json2test.EventParser(tc),
			// Write non-JSON output to log.
			JSONLessTestLineWriter(t, "serial"),
		),
	})
	defer cleanup()

	if err := q.ExpectTimeout("GoTest Done", 120*time.Second); err != nil {
		t.Errorf("Waiting for GoTest Done: %v", err)
	}

	for pkg, test := range tc.Tests {
		switch test.State {
		case json2test.StateFail:
			t.Errorf("Test %v failed:\n%v", pkg, test.FullOutput)
		case json2test.StateSkip:
			t.Logf("Test %v skipped", pkg)
		case json2test.StatePass:
			// Nothing.
		default:
			t.Errorf("Test %v left in state %v:\n%v", pkg, test.State, test.FullOutput)
		}
	}
}
