// Copyright 2022 the u-root Authors. All rights reserved
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
	"time"

	"github.com/hugelgupf/vmtest/internal/json2test"
	"github.com/hugelgupf/vmtest/internal/testevent"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/uio/cp"
	"golang.org/x/tools/go/packages"
)

func lookupPkgs(env golang.Environ, dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles,
		Env:   append(os.Environ(), env.Env()...),
		Dir:   dir,
		Tests: true,
	}
	return packages.Load(cfg, patterns...)
}

func compileTestAndData(env *golang.Environ, pkg, destDir string, cover bool) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}

	testFile := filepath.Join(destDir, fmt.Sprintf("%s.test", path.Base(pkg)))

	args := []string{
		"-gcflags=all=-l",
		"-ldflags", "-s -w",
		"-c", pkg,
		"-o", testFile,
	}
	if cover {
		args = append(args, "-covermode=atomic")
	}
	cmd := env.GoCmd("test", args...)
	if stderr, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("could not build %s: %v\n%s", pkg, err, string(stderr))
	}

	// When a package does not contain any tests, the test
	// executable is not generated, so it is not included in the
	// `tests` list.
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		pkgs, err := lookupPkgs(*env, "", pkg)
		if err != nil {
			return fmt.Errorf("failed to look up package %q: %v", pkg, err)
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
			return fmt.Errorf("could not find package directory for %q", pkg)
		}

		// Optimistically copy any files in the pkg's
		// directory, in case e.g. a testdata dir is there.
		if err := copyRelativeFiles(dir, destDir); err != nil {
			return err
		}
	}
	return nil
}

// GoTestOptions is configuration for RunGoTestsInVM.
type GoTestOptions struct {
	VMOpts      []Opt
	Packages    []string
	TestTimeout time.Duration
}

// GoTestOpt is a configurator for GoTestOptions.
type GoTestOpt func(t testing.TB, o *GoTestOptions) error

// WithVMOpt appends the VM configurators for use with Go tests.
func WithVMOpt(opts ...Opt) GoTestOpt {
	return func(t testing.TB, o *GoTestOptions) error {
		o.VMOpts = append(o.VMOpts, opts...)
		return nil
	}
}

// AppendPackage adds additional packages to the test.
func AppendPackage(pkgs ...string) GoTestOpt {
	return func(t testing.TB, o *GoTestOptions) error {
		o.Packages = append(o.Packages, pkgs...)
		return nil
	}
}

// WithGoTestTimeout sets a timeout for individual Go test binaries.
func WithGoTestTimeout(timeout time.Duration) GoTestOpt {
	return func(t testing.TB, o *GoTestOptions) error {
		o.TestTimeout = timeout
		return nil
	}
}

// RunGoTestsInVM compiles the tests found in pkgs and runs them in a QEMU VM
// configured in options `o`. It collects the test results and provides a
// pass/fail result of each individual test.
//
// RunGoTestsInVM runs tests and benchmarks, but not fuzz tests. Guest test
// architecture can be set with VMTEST_ARCH.
//
// The test environment in the VM is very minimal. If a test depends on other
// binaries or specific files to be present, they must be specified with
// additional initramfs commands via WithMergedInitramfs.
//
// All files and directories in the same directory as the test package will be
// made available to the test in the guest as well (e.g. testdata/
// directories).
//
// Coverage from the Go tests is collected if a coverage file name is specified
// via the VMTEST_GO_PROFILE env var.
//
//   - TODO: specify test, bench, fuzz filter. Flags for fuzzing.
func RunGoTestsInVM(t testing.TB, pkgs []string, opts ...GoTestOpt) {
	SkipWithoutQEMU(t)

	goOpts := &GoTestOptions{
		Packages: pkgs,
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(t, goOpts); err != nil {
				t.Fatal(err)
			}
		}
	}

	sharedDir := testtmp.TempDir(t)
	vmCoverProfile, ok := os.LookupEnv("VMTEST_GO_PROFILE")
	if !ok {
		t.Log("In-guest Go test coverage is not collected unless VMTEST_GO_PROFILE is set")
	}

	// Set up u-root build options.
	env := golang.Default(golang.DisableCGO(), golang.WithGOARCH(string(qemu.GuestArch())))

	// Statically build tests and add them to the temporary directory.
	testDir := filepath.Join(sharedDir, "tests")

	// Compile the Go tests. Place the test binaries in a directory that
	// will be shared with the VM using 9P.
	for _, pkg := range goOpts.Packages {
		pkgDir := filepath.Join(testDir, pkg)
		if err := compileTestAndData(env, pkg, pkgDir, len(vmCoverProfile) > 0); err != nil {
			t.Fatal(err)
		}
	}

	var uinitArgs []string
	if len(vmCoverProfile) > 0 {
		uinitArgs = append(uinitArgs, "-coverprofile=/gotestdata/coverage.profile")
	}
	if goOpts.TestTimeout > 0 {
		uinitArgs = append(uinitArgs, fmt.Sprintf("-test_timeout=%s", goOpts.TestTimeout))
	}
	initramfs := uroot.Opts{
		Env: env,
		Commands: append(
			uroot.BusyBoxCmds(
				"github.com/u-root/u-root/cmds/core/dhclient",
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/hugelgupf/vmtest/vminit/gouinit",
			),
			uroot.BinaryCmds("cmd/test2json")...),
		InitCmd:   "init",
		UinitCmd:  "gouinit",
		UinitArgs: uinitArgs,
		TempDir:   testtmp.TempDir(t),
	}

	qemuFns := []qemu.Fn{
		qemu.P9Directory(sharedDir, "gotests"),
	}
	goCov := os.Getenv("GOCOVERDIR")
	if goCov != "" {
		qemuFns = append(qemuFns,
			qemu.P9Directory(goCov, "gocov"),
			qemu.WithAppendKernel("VMTEST_GOCOVERDIR=gocov"),
		)
	}
	// Create the initramfs and start the VM.
	vm := StartVM(t, append(
		[]Opt{
			WithMergedInitramfs(initramfs),
			WithQEMUFn(qemuFns...),
			CollectKernelCoverage(),
		}, goOpts.VMOpts...)...)

	if err := vm.Wait(); err != nil {
		t.Errorf("VM exited with %v", err)
	}

	// Collect Go coverage.
	if len(vmCoverProfile) > 0 {
		if err := cp.Copy(filepath.Join(sharedDir, "coverage.profile"), vmCoverProfile); err != nil {
			t.Errorf("Could not copy coverage file: %v", err)
		}
	}

	errors, err := qemu.ReadEventFile[testevent.ErrorEvent](filepath.Join(sharedDir, "errors.json"))
	if err != nil {
		t.Errorf("Reading test events: %v", err)
	}
	for _, e := range errors {
		t.Errorf("Binary %s experienced error: %s", e.Binary, e.Error)
	}

	tc := json2test.NewTestCollector()
	events, err := qemu.ReadEventFile[json2test.TestEvent](filepath.Join(sharedDir, "results.json"))
	if err != nil {
		t.Errorf("Reading Go test events: %v", err)
	}
	for _, event := range events {
		tc.Handle(event)
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
