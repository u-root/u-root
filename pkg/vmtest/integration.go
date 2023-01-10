// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/testutil"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/ulog/ulogtest"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// Options are integration test options.
type Options struct {
	// BuildOpts are u-root initramfs options.
	//
	// They are used if the test needs to generate an initramfs.
	// Fields that are not set are populated by QEMU and QEMUTest as
	// possible.
	BuildOpts uroot.Opts

	// QEMUOpts are QEMU VM options for the test.
	//
	// Fields that are not set are populated by QEMU and QEMUTest as
	// possible.
	QEMUOpts qemu.Options

	// Name is the test's name.
	//
	// If name is left empty, the calling function's function name will be
	// used as determined by runtime.Caller.
	Name string

	// Uinit is the uinit that should be added to a generated initramfs.
	//
	// If none is specified, the generic uinit will be used, which searches for
	// and runs the script generated from TestCmds.
	Uinit string

	// TestCmds are commands to execute after init.
	//
	// QEMUTest generates an Elvish script with these commands. The script is
	// shared with the VM, and is run from the generic uinit.
	TestCmds []string

	// TmpDir is the temporary directory exposed to the QEMU VM.
	TmpDir string

	// Logger logs build statements.
	Logger ulog.Logger

	// Extra environment variables to set when building (used by u-bmc)
	ExtraBuildEnv []string

	// Use virtual vfat rather than 9pfs
	UseVVFAT bool

	// By default, if your kernel has CONFIG_DEBUG_FS=y and
	// CONFIG_GCOV_KERNEL=y enabled, the kernel's coverage will be
	// collected and saved to:
	//   u-root/integration/coverage/{{testname}}/{{instance}}/kernel_coverage.tar
	NoKernelCoverage bool
}

// Tests are run from u-root/integration/{gotests,generic-tests}/
const coveragePath = "../coverage"

// Keeps track of the number of instances per test so we do not overlap
// coverage reports.
var instance = map[string]int{}

func last(s string) string {
	l := strings.Split(s, ".")
	return l[len(l)-1]
}

func callerName(depth int) string {
	// Use the test name as the serial log's file name.
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		panic("runtime caller failed")
	}
	f := runtime.FuncForPC(pc)
	return last(f.Name())
}

// TestLineWriter is an io.Writer that logs full lines of serial to tb.
func TestLineWriter(tb testing.TB, prefix string) io.WriteCloser {
	return uio.FullLineWriter(&testLineWriter{tb: tb, prefix: prefix})
}

type jsonStripper struct {
	uio.LineWriter
}

func (j jsonStripper) OneLine(p []byte) {
	// Poor man's JSON detector.
	if len(p) == 0 || p[0] == '{' {
		return
	}
	j.LineWriter.OneLine(p)
}

func JSONLessTestLineWriter(tb testing.TB, prefix string) io.WriteCloser {
	return uio.FullLineWriter(jsonStripper{&testLineWriter{tb: tb, prefix: prefix}})
}

// testLineWriter is an io.Writer that logs full lines of serial to tb.
type testLineWriter struct {
	tb     testing.TB
	prefix string
}

func replaceCtl(str []byte) []byte {
	for i, c := range str {
		if c == 9 || c == 10 {
		} else if c < 32 || c == 127 {
			str[i] = '~'
		}
	}
	return str
}

func (tsw *testLineWriter) OneLine(p []byte) {
	tsw.tb.Logf("%s %s: %s", testutil.NowLog(), tsw.prefix, string(replaceCtl(p)))
}

// TestArch returns the architecture under test. Pass this as GOARCH when
// building Go programs to be run in the QEMU environment.
func TestArch() string {
	if env := os.Getenv("UROOT_TESTARCH"); env != "" {
		return env
	}
	return "amd64"
}

// SkipWithoutQEMU skips the test when the QEMU environment variables are not
// set. This is already called by QEMUTest(), so use if some expensive
// operations are performed before calling QEMUTest().
func SkipWithoutQEMU(t *testing.T) {
	if _, ok := os.LookupEnv("UROOT_QEMU"); !ok {
		t.Skip("QEMU test is skipped unless UROOT_QEMU is set")
	}
	if _, ok := os.LookupEnv("UROOT_KERNEL"); !ok {
		t.Skip("QEMU test is skipped unless UROOT_KERNEL is set")
	}
}

func saveCoverage(t *testing.T, path string) error {
	// Coverage may not have been collected, for example if the kernel is
	// not built with CONFIG_GCOV_KERNEL.
	if fi, err := os.Stat(path); os.IsNotExist(err) || (err != nil && !fi.Mode().IsRegular()) {
		return nil
	}

	// Move coverage to common directory.
	uniqueCoveragePath := filepath.Join(coveragePath, t.Name(), fmt.Sprintf("%d", instance[t.Name()]))
	instance[t.Name()]++
	if err := os.MkdirAll(uniqueCoveragePath, 0o770); err != nil {
		return err
	}
	if err := os.Rename(path, filepath.Join(uniqueCoveragePath, filepath.Base(path))); err != nil {
		return err
	}
	return nil
}

func QEMUTest(t *testing.T, o *Options) (*qemu.VM, func()) {
	SkipWithoutQEMU(t)

	// Delete any previous coverage data.
	if _, ok := instance[t.Name()]; !ok {
		testCoveragePath := filepath.Join(coveragePath, t.Name())
		if err := os.RemoveAll(testCoveragePath); err != nil && !os.IsNotExist(err) {
			t.Logf("Error erasing previous coverage: %v", err)
		}
	}

	if len(o.Name) == 0 {
		o.Name = callerName(2)
	}
	if o.Logger == nil {
		o.Logger = &ulogtest.Logger{TB: t}
	}
	if o.QEMUOpts.SerialOutput == nil {
		o.QEMUOpts.SerialOutput = TestLineWriter(t, "serial")
	}

	// Create or reuse a temporary directory. This is exposed to the VM.
	if o.TmpDir == "" {
		tmpDir, err := os.MkdirTemp("", "uroot-integration")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		o.TmpDir = tmpDir
	}

	qOpts, err := QEMU(o)
	if err != nil {
		t.Fatalf("Failed to create QEMU VM %s: %v", o.Name, err)
	}

	vm, err := qOpts.Start()
	if err != nil {
		t.Fatalf("Failed to start QEMU VM %s: %v", o.Name, err)
	}

	return vm, func() {
		vm.Close()
		if !o.NoKernelCoverage {
			if err := saveCoverage(t, filepath.Join(o.TmpDir, "kernel_coverage.tar")); err != nil {
				t.Logf("Error saving kernel coverage: %v", err)
			}
		}

		t.Logf("QEMU command line to reproduce %s:\n%s", o.Name, vm.CmdlineQuoted())
		if t.Failed() {
			t.Log("Keeping temp dir: ", o.TmpDir)
		} else if len(o.TmpDir) == 0 {
			if err := os.RemoveAll(o.TmpDir); err != nil {
				t.Logf("failed to remove temporary directory %s: %v", o.TmpDir, err)
			}
		}
	}
}

// QEMU builds the u-root environment and prepares QEMU options given the test
// options and environment variables.
//
// QEMU will augment o.BuildOpts and o.QEMUOpts with configuration that the
// caller either requested (through the Options.Uinit field, for example) or
// that the caller did not set.
//
// QEMU returns the QEMU launch options or an error.
func QEMU(o *Options) (*qemu.Options, error) {
	if len(o.Name) == 0 {
		o.Name = callerName(2)
	}

	// Generate Elvish shell script of test commands in o.TmpDir.
	if len(o.TestCmds) > 0 {
		testFile := filepath.Join(o.TmpDir, "test.elv")

		if err := os.WriteFile(
			testFile, []byte(strings.Join(o.TestCmds, "\n")), 0o777); err != nil {
			return nil, err
		}
	}

	// Set the initramfs.
	if len(o.QEMUOpts.Initramfs) == 0 {
		o.QEMUOpts.Initramfs = filepath.Join(o.TmpDir, "initramfs.cpio")
		if err := ChooseTestInitramfs(o.BuildOpts, o.Uinit, o.QEMUOpts.Initramfs); err != nil {
			return nil, err
		}
	}

	if len(o.QEMUOpts.Kernel) == 0 {
		// Copy kernel to o.TmpDir for tests involving kexec.
		kernel := filepath.Join(o.TmpDir, "kernel")
		if err := cp.Copy(os.Getenv("UROOT_KERNEL"), kernel); err != nil {
			return nil, err
		}
		o.QEMUOpts.Kernel = kernel
	}

	switch TestArch() {
	case "amd64":
		o.QEMUOpts.KernelArgs += " console=ttyS0 earlyprintk=ttyS0"
	case "arm":
		o.QEMUOpts.KernelArgs += " console=ttyAMA0"
	}
	o.QEMUOpts.KernelArgs += " uroot.vmtest"

	var dir qemu.Device
	if o.UseVVFAT {
		dir = qemu.ReadOnlyDirectory{Dir: o.TmpDir}
	} else {
		dir = qemu.P9Directory{Dir: o.TmpDir, Arch: TestArch()}
	}
	o.QEMUOpts.Devices = append(o.QEMUOpts.Devices, qemu.VirtioRandom{}, dir)

	if o.NoKernelCoverage {
		o.QEMUOpts.KernelArgs += " UROOT_NO_KERNEL_COVERAGE=1"
	}

	return &o.QEMUOpts, nil
}

// ChooseTestInitramfs chooses which initramfs will be used for a given test and
// places it at the location given by outputFile.
// Default to the override initramfs if one is specified in the UROOT_INITRAMFS
// environment variable. Else, build an initramfs with the given parameters.
// If no uinit was provided, the generic one is used.
func ChooseTestInitramfs(o uroot.Opts, uinit, outputFile string) error {
	override := os.Getenv("UROOT_INITRAMFS")
	if len(override) > 0 {
		log.Printf("Overriding with initramfs %q", override)
		return cp.Copy(override, outputFile)
	}

	if len(uinit) == 0 {
		log.Printf("Defaulting to generic initramfs")
		uinit = "github.com/u-root/u-root/integration/testcmd/generic/uinit"
	}

	_, err := CreateTestInitramfs(o, uinit, outputFile)
	return err
}

// CreateTestInitramfs creates an initramfs with the given build options and
// uinit, and writes it to the given output file. If no output file is provided,
// one will be created.
// The output file name is returned. It is the caller's responsibility to remove
// the initramfs file after use.
func CreateTestInitramfs(o uroot.Opts, uinit, outputFile string) (string, error) {
	if o.Env == nil {
		env := gbbgolang.Default()
		env.CgoEnabled = false
		env.GOARCH = TestArch()
		o.Env = &env
	}

	if o.UrootSource == "" {
		sourcePath, ok := os.LookupEnv("UROOT_SOURCE")
		if !ok {
			return "", fmt.Errorf("failed to get u-root source directory, please set UROOT_SOURCE to the absolute path of the u-root source directory")
		}
		o.UrootSource = sourcePath
	}

	logger := log.New(os.Stderr, "", 0)

	// If build opts don't specify any commands, include all commands. Else,
	// always add init and elvish.
	var cmds []string
	if len(o.Commands) == 0 {
		cmds = []string{
			"github.com/u-root/u-root/cmds/core/*",
			"github.com/u-root/u-root/cmds/exp/*",
		}
	}

	if len(uinit) != 0 {
		cmds = append(cmds, uinit)
	}

	// Add our commands to the build opts.
	o.AddBusyBoxCommands(cmds...)

	// Fill in the default build options if not specified.
	if o.BaseArchive == nil {
		o.BaseArchive = uroot.DefaultRamfs().Reader()
	}
	if len(o.InitCmd) == 0 {
		o.InitCmd = "init"
	}
	if len(o.DefaultShell) == 0 {
		o.DefaultShell = "elvish"
	}
	if len(o.TempDir) == 0 {
		tempDir, err := os.MkdirTemp("", "initramfs-tempdir")
		if err != nil {
			return "", fmt.Errorf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)
		o.TempDir = tempDir
	}

	// Create an output file if one was not provided.
	if len(outputFile) == 0 {
		f, err := os.CreateTemp("", "initramfs.cpio")
		if err != nil {
			return "", fmt.Errorf("failed to create output file: %v", err)
		}
		outputFile = f.Name()
	}
	w, err := initramfs.CPIO.OpenWriter(logger, outputFile)
	if err != nil {
		return "", fmt.Errorf("Failed to create initramfs writer: %v", err)
	}
	o.OutputFile = w

	return outputFile, uroot.CreateInitramfs(logger, o)
}
