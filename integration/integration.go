// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

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
	"time"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/builder"
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
	// Env is the Go environment to use to build u-root.
	Env *golang.Environ

	// Name is the test's name.
	//
	// If name is left empty, the calling function's function name will be
	// used as determined by runtime.Caller
	Name string

	// Go commands to include in the initramfs for the VM.
	//
	// If left empty, all u-root commands will be included.
	Cmds []string

	// Uinit are commands to execute after init.
	//
	// If populated, a uinit.go will be generated from these.
	Uinit []string

	// Files are files to include in the VMs initramfs.
	Files []string

	// TmpDir is a temporary directory for build artifacts.
	TmpDir string

	// LogFile is a file to log serial output to.
	//
	// The default is serial/$Name.log
	LogFile string

	// Logger logs build statements.
	Logger logger.Logger

	// Timeout is the timeout for expect statements.
	Timeout time.Duration

	// Network is the VM's network.
	Network *qemu.Network

	// Extra environment variables to set when building (used by u-bmc)
	ExtraBuildEnv []string

	// Serial Output
	SerialOutput io.WriteCloser
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

	qOpts, err := QEMU(o)
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
		dir := vm.Options.Devices[0].(qemu.ReadOnlyDirectory).Dir
		if t.Failed() {
			t.Log("keeping temp dir: ", dir)
		} else if len(o.TmpDir) == 0 {
			if err := os.RemoveAll(dir); err != nil {
				t.Logf("failed to remove temporary directory %s: %v", dir, err)
			}
		}
	}
}

func QEMU(o *Options) (*qemu.Options, error) {
	if len(o.Name) == 0 {
		o.Name = callerName(2)
	}

	if o.Env == nil {
		env := golang.Default()
		o.Env = &env
		o.Env.CgoEnabled = false
		env.GOARCH = TestArch()
	}

	if len(o.LogFile) == 0 {
		// Create file for serial logs.
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create serial log directory: %v", err)
		}

		o.LogFile = filepath.Join(logDir, fmt.Sprintf("%s.log", o.Name))
	}

	var cmds []string
	if len(o.Cmds) == 0 {
		cmds = append(cmds, "github.com/u-root/u-root/cmds/*")
	} else {
		cmds = append(cmds, o.Cmds...)
	}
	// Create a uinit from the commands given.
	if len(o.Uinit) > 0 {
		urootPkg, err := o.Env.Package("github.com/u-root/u-root/integration")
		if err != nil {
			return nil, err
		}
		testDir := filepath.Join(urootPkg.Dir, "testcmd")

		dirpath, err := ioutil.TempDir(testDir, "uinit-")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(dirpath)

		if err := os.MkdirAll(filepath.Join(dirpath, "uinit"), 0755); err != nil {
			return nil, err
		}

		var realUinit [][]string
		for _, cmd := range o.Uinit {
			realUinit = append(realUinit, fields(cmd))
		}

		if err := ioutil.WriteFile(
			filepath.Join(dirpath, "uinit", "uinit.go"),
			[]byte(fmt.Sprintf(template, realUinit)),
			0755); err != nil {
			return nil, err
		}
		cmds = append(cmds, path.Join("github.com/u-root/u-root/integration/testcmd", filepath.Base(dirpath), "uinit"))
	}

	// Create or reuse a temporary directory.
	tmpDir := o.TmpDir
	if len(tmpDir) == 0 {
		var err error
		tmpDir, err = ioutil.TempDir("", "uroot-integration")
		if err != nil {
			return nil, err
		}
	}

	if o.Logger == nil {
		o.Logger = log.New(os.Stderr, "", 0)
	}

	// OutputFile
	outputFile := filepath.Join(tmpDir, "initramfs.cpio")
	w, err := initramfs.CPIO.OpenWriter(o.Logger, outputFile, "", "")
	if err != nil {
		return nil, err
	}

	// Build u-root
	opts := uroot.Opts{
		Env: *o.Env,
		Commands: []uroot.Commands{
			{
				Builder:  builder.BusyBox,
				Packages: cmds,
			},
		},
		ExtraFiles:   o.Files,
		TempDir:      tmpDir,
		BaseArchive:  uroot.DefaultRamfs.Reader(),
		OutputFile:   w,
		InitCmd:      "init",
		DefaultShell: "elvish",
	}
	if err := uroot.CreateInitramfs(o.Logger, opts); err != nil {
		return nil, err
	}

	// Copy kernel to tmpDir.
	bzImage := filepath.Join(tmpDir, "bzImage")
	if err := cp.Copy(os.Getenv("UROOT_KERNEL"), bzImage); err != nil {
		return nil, err
	}

	logFile := o.SerialOutput
	if logFile == nil {
		if o.LogFile != "" {
			logFile, err = os.Create(o.LogFile)
			if err != nil {
				return nil, fmt.Errorf("could not create log file: %v", err)
			}
		}
	}

	kernelArgs := ""
	switch TestArch() {
	case "amd64":
		kernelArgs = "console=ttyS0 earlyprintk=ttyS0"
	case "arm":
		kernelArgs = "console=ttyAMA0"
	}

	return &qemu.Options{
		Initramfs:    outputFile,
		Kernel:       bzImage,
		KernelArgs:   kernelArgs,
		SerialOutput: logFile,
		Timeout:      o.Timeout,
		Devices: []qemu.Device{
			qemu.ReadOnlyDirectory{Dir: tmpDir},
			qemu.VirtioRandom{},
			o.Network,
		},
	}, nil
}
