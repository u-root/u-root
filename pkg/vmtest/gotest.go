// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest/internal/json2test"
	"golang.org/x/tools/go/packages"
)

func lookupPkgs(env *gbbgolang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles,
		Env:   append(os.Environ(), env.Env()...),
		Dir:   dir,
		Tests: true,
	}
	return packages.Load(cfg, patterns...)
}

// GolangTest compiles the unit tests found in pkgs and runs them in a QEMU VM.
func GolangTest(t *testing.T, pkgs []string, o *Options) {
	SkipWithoutQEMU(t)
	// TODO: support arm
	if TestArch() != "amd64" && TestArch() != "arm64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	vmCoverProfile, ok := os.LookupEnv("UROOT_QEMU_COVERPROFILE")
	if !ok {
		t.Log("QEMU test coverage is not collected unless UROOT_QEMU_COVERPROFILE is set")
	}

	if o == nil {
		o = &Options{}
	}

	// Create a temporary directory.
	if len(o.TmpDir) == 0 {
		tmpDir, err := os.MkdirTemp("", "uroot-integration")
		if err != nil {
			t.Fatal(err)
		}
		o.TmpDir = tmpDir
	}

	if o.BuildOpts.UrootSource == "" {
		sourcePath, ok := os.LookupEnv("UROOT_SOURCE")
		if !ok {
			t.Fatal("This test needs UROOT_SOURCE set to the absolute path of the checked out u-root source")
		}
		o.BuildOpts.UrootSource = sourcePath
	}

	// Set up u-root build options.
	env := gbbgolang.Default(gbbgolang.DisableCGO(), gbbgolang.WithGOARCH(TestArch()))
	o.BuildOpts.Env = env

	// Statically build tests and add them to the temporary directory.
	var tests []string
	testDir := filepath.Join(o.TmpDir, "tests")

	if len(vmCoverProfile) > 0 {
		f, err := os.Create(filepath.Join(o.TmpDir, "coverage.profile"))
		if err != nil {
			t.Fatalf("Could not create coverage file %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("Could not close coverage.profile: %v", err)
		}
	}

	for _, pkg := range pkgs {
		pkgDir := filepath.Join(testDir, pkg)
		if err := os.MkdirAll(pkgDir, 0o755); err != nil {
			t.Fatal(err)
		}

		testFile := filepath.Join(pkgDir, fmt.Sprintf("%s.test", path.Base(pkg)))

		args := []string{
			"test",
			"-gcflags=all=-l",
			"-ldflags", "-s -w",
			"-c", pkg,
			"-o", testFile,
		}
		if len(vmCoverProfile) > 0 {
			args = append(args, "-covermode=atomic")
		}

		cmd := env.GoCmd(args...)
		if stderr, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("could not build %s: %v\n%s", pkg, err, string(stderr))
		}

		// When a package does not contain any tests, the test
		// executable is not generated, so it is not included in the
		// `tests` list.
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			tests = append(tests, pkg)

			pkgs, err := lookupPkgs(o.BuildOpts.Env, "", pkg)
			if err != nil {
				t.Fatalf("Failed to look up package %q: %v", pkg, err)
			}

			// One directory = one package in standard Go, so
			// finding the first file's parent directory should
			// find us the package directory.
			var dir string
			for _, p := range pkgs {
				if len(p.GoFiles) > 0 {
					dir = filepath.Dir(p.GoFiles[0])
				}
			}
			if dir == "" {
				t.Fatalf("Could not find package directory for %q", pkg)
			}

			// Optimistically copy any files in the pkg's
			// directory, in case e.g. a testdata dir is there.
			if err := copyRelativeFiles(dir, filepath.Join(testDir, pkg)); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Create the CPIO and start QEMU.
	o.BuildOpts.AddBusyBoxCommands("github.com/u-root/u-root/cmds/core/*")
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
	if len(vmCoverProfile) > 0 {
		o.QEMUOpts.KernelArgs += " uroot.uinitargs=-coverprofile=/testdata/coverage.profile"
	}

	q, cleanup := QEMUTest(t, o)
	defer cleanup()

	if err := q.Expect("TESTS PASSED MARKER"); err != nil {
		t.Errorf("Waiting for 'TESTS PASSED MARKER' signal: %v", err)
	}

	if len(vmCoverProfile) > 0 {
		cov, err := os.Open(filepath.Join(o.TmpDir, "coverage.profile"))
		if err != nil {
			t.Fatalf("No coverage file shared from VM: %v", err)
		}

		out, err := os.OpenFile(vmCoverProfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			t.Fatalf("Could not open vmcoverageprofile: %v", err)
		}

		if _, err := io.Copy(out, cov); err != nil {
			t.Fatalf("Error copying coverage: %s", err)
		}
		if err := out.Close(); err != nil {
			t.Fatalf("Could not close vmcoverageprofile: %v", err)
		}
		if err := cov.Close(); err != nil {
			t.Fatalf("Could not close coverage.profile: %v", err)
		}
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
