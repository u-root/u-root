// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package govmtest is an API for running Go unit tests in the guest and
// collecting their results and test coverage.
package govmtest

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
	"github.com/hugelgupf/vmtest/qemu/qcoverage"
	"github.com/hugelgupf/vmtest/qemu/qevent"
	"github.com/hugelgupf/vmtest/qemu/quimage"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
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

// Options configures a Go test.
type Options struct {
	Packages    []string
	QEMUOpts    []qemu.Fn
	Initramfs   []uimage.Modifier
	TestTimeout time.Duration
}

// Modifier is a configurator for Options.
type Modifier func(t testing.TB, o *Options) error

// WithQEMUFn adds QEMU options.
func WithQEMUFn(fn ...qemu.Fn) Modifier {
	return func(_ testing.TB, o *Options) error {
		o.QEMUOpts = append(o.QEMUOpts, fn...)
		return nil
	}
}

// WithUimage merges o with already appended initramfs build options.
func WithUimage(mods ...uimage.Modifier) Modifier {
	return func(_ testing.TB, o *Options) error {
		o.Initramfs = append(o.Initramfs, mods...)
		return nil
	}
}

// WithPackageToTest adds additional packages to the test.
func WithPackageToTest(pkgs ...string) Modifier {
	return func(t testing.TB, o *Options) error {
		o.Packages = append(o.Packages, pkgs...)
		return nil
	}
}

// WithGoTestTimeout sets a timeout for individual Go test binaries.
func WithGoTestTimeout(timeout time.Duration) Modifier {
	return func(t testing.TB, o *Options) error {
		o.TestTimeout = timeout
		return nil
	}
}

// Run compiles the tests added with WithPackageToTest and runs them in a QEMU
// VM configured by mods. It collects the test results and provides a pass/fail
// result of each individual test.
//
// Run runs tests and benchmarks, but not fuzz tests.
//
// The test environment in the VM is very minimal. If a test depends on other
// binaries or specific files to be present, they must be specified with
// additional initramfs commands via [WithUimage].
//
// All files and directories in the same directory as the test package will be
// made available to the test in the guest as well (e.g. testdata/
// directories).
//
// Coverage from the Go tests is collected if a coverage file name is specified
// via the VMTEST_GO_PROFILE env var, as well as integration test coverage if
// VMTEST_GOCOVERDIR is set.
//
//   - TODO: specify test, bench, fuzz filter. Flags for fuzzing.
func Run(t testing.TB, name string, mods ...Modifier) {
	qemu.SkipWithoutQEMU(t)

	goOpts := &Options{}
	for _, mod := range mods {
		if mod != nil {
			if err := mod(t, goOpts); err != nil {
				t.Fatal(err)
			}
		}
	}
	if len(goOpts.Packages) == 0 {
		t.Fatal("No packages specified for govmtest")
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
		uinitArgs = append(uinitArgs, "-coverprofile=/mount/9p/gotestdata/coverage.profile")
	}
	if goOpts.TestTimeout > 0 {
		uinitArgs = append(uinitArgs, fmt.Sprintf("-test_timeout=%s", goOpts.TestTimeout))
	}

	umods := append([]uimage.Modifier{
		uimage.WithBusyboxCommands(
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/hugelgupf/vmtest/vminit/shutdownafter",
			"github.com/hugelgupf/vmtest/vminit/vmmount",
			"github.com/hugelgupf/vmtest/vminit/gouinit",
		),
		uimage.WithBinaryCommands("cmd/test2json"),
		uimage.WithInit("init"),
		uimage.WithUinit("shutdownafter", append([]string{"--", "vmmount", "--", "gouinit"}, uinitArgs...)...),
	}, goOpts.Initramfs...)

	// Create the initramfs and start the VM.
	vm := qemu.StartT(t,
		name,
		qemu.ArchUseEnvv,
		append([]qemu.Fn{
			quimage.WithUimageT(t, umods...),
			qemu.P9Directory(sharedDir, "gotestdata"),
			qcoverage.CollectKernelCoverage(t),
			qcoverage.ShareGOCOVERDIR(),
			qemu.WithVmtestIdent(),
		}, goOpts.QEMUOpts...)...)
	if err := vm.Wait(); err != nil {
		t.Errorf("VM exited with %v", err)
	}

	// Collect Go coverage.
	if len(vmCoverProfile) > 0 {
		if err := cp.Copy(filepath.Join(sharedDir, "coverage.profile"), vmCoverProfile); err != nil {
			t.Errorf("Could not copy coverage file: %v", err)
		}
	}

	errors, err := qevent.ReadFile[testevent.ErrorEvent](filepath.Join(sharedDir, "errors.json"))
	if err != nil {
		t.Errorf("Reading test events: %v", err)
	}
	for _, e := range errors {
		t.Errorf("Binary %s experienced error: %s", e.Binary, e.Error)
	}

	tc := json2test.NewTestCollector()
	events, err := qevent.ReadFile[json2test.TestEvent](filepath.Join(sharedDir, "results.json"))
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
