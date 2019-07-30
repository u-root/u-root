// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"github.com/u-root/u-root/pkg/uroot/logger"
)

// Serial output is written to this directory and picked up by circleci, or
// you, if you want to read the serial logs.
const logDir = "serial"

const template = `
package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	for _, cmds := range %#v {
		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
`

// Options are integration test options.
type Options struct {
	// BuildOpts are u-root initramfs options.
	//
	// Fields that are not set are populated by QEMU and QEMUTest as
	// possible.
	BuildOpts uroot.Opts

	// QEMUOpts are QEMU VM options for the test.
	//
	// Fields that are not set are populated by QEMU and QEMUTest as
	// possible.
	QEMUOpts qemu.Options

	// DontSetEnv doesn't set the BuildOpts.Env and uses the user-supplied one.
	//
	// TODO: make uroot.Opts.Env a pointer?
	DontSetEnv bool

	// Name is the test's name.
	//
	// If name is left empty, the calling function's function name will be
	// used as determined by runtime.Caller.
	Name string

	// Uinit are commands to execute after init.
	//
	// If populated, a uinit.go will be generated from these and added to
	// the busybox generated in BuildOpts.Commands.
	Uinit []string

	// Logger logs build statements.
	Logger logger.Logger

	// Extra environment variables to set when building (used by u-bmc)
	ExtraBuildEnv []string

	// Use virtual vfat rather than 9pfs
	UseVVFAT bool
}

func last(s string) string {
	l := strings.Split(s, ".")
	return l[len(l)-1]
}

type testLogger struct {
	t *testing.T
}

func (tl testLogger) Printf(format string, v ...interface{}) {
	tl.t.Logf(format, v...)
}

func (tl testLogger) Print(v ...interface{}) {
	tl.t.Log(v...)
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

func (tsw *testLineWriter) OneLine(p []byte) {
	tsw.tb.Logf("%s: %s", tsw.prefix, strings.ReplaceAll(string(p), "\033", "~"))
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

func QEMUTest(t *testing.T, o *Options) (*qemu.VM, func()) {
	SkipWithoutQEMU(t)

	if len(o.Name) == 0 {
		o.Name = callerName(2)
	}
	if o.Logger == nil {
		o.Logger = &testLogger{t}
	}
	if o.QEMUOpts.SerialOutput == nil {
		o.QEMUOpts.SerialOutput = TestLineWriter(t, "serial")
	}
	if TestArch() == "arm" {
		//currently, 9p does not work on arm
		o.UseVVFAT = true
	}

	qOpts, tmpDir, err := QEMU(o)
	if err != nil {
		t.Fatalf("Failed to create QEMU VM %s: %v", o.Name, err)
	}

	vm, err := qOpts.Start()
	if err != nil {
		t.Fatalf("Failed to start QEMU VM %s: %v", o.Name, err)
	}
	t.Logf("QEMU command line for %s:\n%s", o.Name, vm.CmdlineQuoted())

	return vm, func() {
		vm.Close()
		if t.Failed() {
			t.Log("Keeping temp dir: ", tmpDir)
		} else if len(o.BuildOpts.TempDir) == 0 {
			if err := os.RemoveAll(o.BuildOpts.TempDir); err != nil {
				t.Logf("failed to remove temporary directory %s: %v", o.BuildOpts.TempDir, err)
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
// QEMU returns the QEMU launch options, the temporary directory exposed to the
// QEMU VM, or an error.
func QEMU(o *Options) (*qemu.Options, string, error) {
	if len(o.Name) == 0 {
		o.Name = callerName(2)
	}

	if len(o.QEMUOpts.Initramfs) == 0 {
		if !o.DontSetEnv {
			env := golang.Default()
			env.CgoEnabled = false
			env.GOARCH = TestArch()
			o.BuildOpts.Env = env
		}

		var cmds []string
		if len(o.BuildOpts.Commands) == 0 {
			cmds = append(cmds, "github.com/u-root/u-root/cmds/*")
		}
		// Create a uinit from the commands given.
		if len(o.Uinit) > 0 {
			urootPkg, err := o.BuildOpts.Env.Package("github.com/u-root/u-root/integration")
			if err != nil {
				return nil, "", err
			}
			testDir := filepath.Join(urootPkg.Dir, "testcmd")

			dirpath, err := ioutil.TempDir(testDir, "uinit-")
			if err != nil {
				return nil, "", err
			}
			defer os.RemoveAll(dirpath)

			if err := os.MkdirAll(filepath.Join(dirpath, "uinit"), 0755); err != nil {
				return nil, "", err
			}

			var realUinit [][]string
			for _, cmd := range o.Uinit {
				realUinit = append(realUinit, fields(cmd))
			}

			if err := ioutil.WriteFile(
				filepath.Join(dirpath, "uinit", "uinit.go"),
				[]byte(fmt.Sprintf(template, realUinit)),
				0755); err != nil {
				return nil, "", err
			}
			cmds = append(cmds, path.Join("github.com/u-root/u-root/integration/testcmd", filepath.Base(dirpath), "uinit"))
		}
		// Add our commands to the build opts.
		if len(cmds) > 0 {
			o.BuildOpts.AddBusyBoxCommands(cmds...)
		}

		// Create or reuse a temporary directory.
		if len(o.BuildOpts.TempDir) == 0 {
			tmpDir, err := ioutil.TempDir("", "uroot-integration")
			if err != nil {
				return nil, "", err
			}
			o.BuildOpts.TempDir = tmpDir
		}
		if o.BuildOpts.BaseArchive == nil {
			o.BuildOpts.BaseArchive = uroot.DefaultRamfs.Reader()
		}
		if len(o.BuildOpts.InitCmd) == 0 {
			o.BuildOpts.InitCmd = "init"
		}
		if len(o.BuildOpts.DefaultShell) == 0 {
			o.BuildOpts.DefaultShell = "elvish"
		}

		if o.Logger == nil {
			o.Logger = log.New(os.Stderr, "", 0)
		}

		// OutputFile
		var outputFile string
		if o.BuildOpts.OutputFile == nil {
			outputFile = filepath.Join(o.BuildOpts.TempDir, "initramfs.cpio")
			w, err := initramfs.CPIO.OpenWriter(o.Logger, outputFile, "", "")
			if err != nil {
				return nil, "", err
			}
			o.BuildOpts.OutputFile = w
		}

		// Finally, create an initramfs!
		if err := uroot.CreateInitramfs(o.Logger, o.BuildOpts); err != nil {
			return nil, "", err
		}

		o.QEMUOpts.Initramfs = outputFile
	}

	if len(o.QEMUOpts.Kernel) == 0 {
		// Copy kernel to tmpDir for tests involving kexec.
		kernel := filepath.Join(o.BuildOpts.TempDir, "kernel")
		if err := cp.Copy(os.Getenv("UROOT_KERNEL"), kernel); err != nil {
			return nil, "", err
		}
		o.QEMUOpts.Kernel = kernel
	}

	if len(o.QEMUOpts.KernelArgs) == 0 {
		var kernelArgs string
		switch TestArch() {
		case "amd64":
			kernelArgs = "console=ttyS0 earlyprintk=ttyS0"
		case "arm":
			kernelArgs = "console=ttyAMA0"
		}
		o.QEMUOpts.KernelArgs = kernelArgs
	}

	var dir qemu.Device
	if o.UseVVFAT {
		dir = qemu.ReadOnlyDirectory{Dir: o.BuildOpts.TempDir}
	} else {
		dir = qemu.P9Directory{Dir: o.BuildOpts.TempDir}
	}
	o.QEMUOpts.Devices = append(o.QEMUOpts.Devices, qemu.VirtioRandom{}, dir)

	return &o.QEMUOpts, o.BuildOpts.TempDir, nil
}
