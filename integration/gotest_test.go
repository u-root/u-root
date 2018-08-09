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
	"testing"
)

// TestGoTest effectively runs "go test ./..." inside a QEMU instance. The
// tests run as root and can do all sorts of things not possible otherwise.
func TestGoTest(t *testing.T) {
	// Create a temporary directory.
	tmpDir, err := ioutil.TempDir("", "uroot-integration")
	if err != nil {
		t.Fatal(err)
	}

	// Create list of packages to test.
	// TODO: include all packages in cmds
	pkgs := []string{
		"github.com/u-root/u-root/cmds/ls",
		"github.com/u-root/u-root/cmds/mknod",
	}

	// Statically build tests and add them to the temporary directory.
	os.Setenv("CGO_ENABLED", "0")
	testDir := filepath.Join(tmpDir, "tests")
	for _, pkg := range pkgs {
		testFile := filepath.Join(testDir, path.Base(pkg))
		cmd := exec.Command("go", "test", "-c", pkg, "-o", testFile)
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
	sort.Strings(pkgs) // Tests are run and checked in sorted order.
	for _, base := range bases {
		passMsg := fmt.Sprintf("#### %s PASSED ####", base)
		if err := q.Expect(passMsg); err == nil {
			t.Logf("go test '%s' passed", base)
		} else {
			t.Errorf("go test '%s' failed: %v", base, err)
		}
	}
}
