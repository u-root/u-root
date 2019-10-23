// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest/internal/json2test"
)

// GolangTest compiles the unit tests found in pkgs and runs them in a QEMU VM.
func GolangTest(t *testing.T, pkgs []string, o *Options) {
	SkipWithoutQEMU(t)
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	if o == nil {
		o = &Options{}
	}

	// Create a temporary directory.
	if len(o.TmpDir) == 0 {
		tmpDir, err := ioutil.TempDir("", "uroot-integration")
		if err != nil {
			t.Fatal(err)
		}
		o.TmpDir = tmpDir
	}

	env := golang.Default()
	env.CgoEnabled = false
	env.GOARCH = TestArch()
	o.BuildOpts.Env = env

	// Statically build tests and add them to the temporary directory.
	var tests []string
	os.Setenv("CGO_ENABLED", "0")
	testDir := filepath.Join(o.TmpDir, "tests")
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(testDir, pkg)
		if err := os.MkdirAll(pkgDir, 0755); err != nil {
			t.Fatal(err)
		}

		testFile := filepath.Join(pkgDir, fmt.Sprintf("%s.test", path.Base(pkg)))

		cmd := exec.Command("go", "test",
			"-ldflags", "-s -w",
			"-c", pkg,
			"-o", testFile,
		)
		if stderr, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("could not build %s: %v\n%s", pkg, err, string(stderr))
		}

		// When a package does not contain any tests, the test
		// executable is not generated, so it is not included in the
		// `tests` list.
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			tests = append(tests, pkg)

			p, err := o.BuildOpts.Env.Package(pkg)
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
	o.BuildOpts.AddCommands(uroot.BinaryCmds("cmd/test2json")...)

	// Specify the custom gotest uinit.
	o.Uinit = "github.com/u-root/u-root/integration/testcmd/gotest/uinit"

	tc := json2test.NewTestCollector()
	serial := []io.Writer{
		// Collect JSON test events in tc.
		json2test.EventParser(tc),
		// Write non-JSON output to log.
		JSONLessTestLineWriter(t, "serial"),
	}
	if o.QEMUOpts.SerialOutput != nil {
		serial = append(serial, o.QEMUOpts.SerialOutput)
	}
	o.QEMUOpts.SerialOutput = uio.MultiWriteCloser(serial...)

	q, cleanup := QEMUTest(t, o)
	defer cleanup()

	if err := q.Expect("GoTest Done"); err != nil {
		t.Errorf("Waiting for GoTest Done: %v", err)
	}

	// TODO: check that tc.Tests == tests
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
