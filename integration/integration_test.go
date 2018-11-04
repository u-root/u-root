// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"fmt"
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
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/builder"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
	"github.com/u-root/u-root/pkg/uroot/logger"
)

// Serial output is written to this directory and picked up by circleci, or
// you, if you want to read the serial logs.
const logDir = "serial"

type options struct {
	// uinitName is the name of a directory containing uinit found at
	// `github.com/u-root/u-root/integration/testdata`.
	uinitName string

	// extraArgs are extra arguments passed to QEMU.
	extraArgs []string

	// tmpDir indicates a path to use as a temporary directory for the
	// test. If this is unset, a new directory is created and returned.
	tmpDir string
}

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

func QEMUTestSetup(t *testing.T) {
	if _, ok := os.LookupEnv("UROOT_QEMU"); !ok {
		t.Skip("test is skipped unless UROOT_QEMU is set")
	}
	if _, ok := os.LookupEnv("UROOT_KERNEL"); !ok {
		t.Skip("test is skipped unless UROOT_KERNEL is set")
	}
}

type Options struct {
	// Go commands to include.
	Cmds []string

	// Commands to execute after init.
	//
	// Will override any uinit included in Cmds.
	Uinit   []string
	Files   []string
	TmpDir  string
	LogFile string
	Logger  logger.Logger
}

func last(s string) string {
	l := strings.Split(s, ".")
	return l[len(l)-1]
}

type TestLogger struct {
	t *testing.T
}

func (tl TestLogger) Printf(format string, v ...interface{}) {
	tl.t.Logf(format, v...)
}

func (tl TestLogger) Print(v ...interface{}) {
	tl.t.Log(v...)
}

// Returns temporary directory and QEMU instance.
func QEMUTest(t *testing.T, o *Options) (*qemu.QEMU, func()) {
	QEMUTestSetup(t)

	// Create file for serial logs.
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("could not create serial log directory: %v", err)
	}

	// Use the test name as the serial log's file name.
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		t.Fatal("runtime caller failed")
	}
	f := runtime.FuncForPC(pc)

	if len(o.LogFile) == 0 {
		o.LogFile = filepath.Join(logDir, fmt.Sprintf("%s.log", last(f.Name())))
	}
	if o.Logger == nil {
		o.Logger = &TestLogger{t}
	}

	q, err := QEMU(o)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("QEMU command line:\n%s", q.CmdlineQuoted())
	if err := q.Start(); err != nil {
		t.Fatal(err)
	}

	return q, func() {
		q.Close()
		if t.Failed() {
			t.Log("keeping temp dir: ", q.SharedDir)
		} else if len(o.TmpDir) == 0 {
			if err := os.RemoveAll(q.SharedDir); err != nil {
				t.Logf("failed to remove temporary directory %s: %v", q.SharedDir, err)
			}
		}
	}
}

func QEMU(o *Options) (*qemu.QEMU, error) {
	env := golang.Default()
	env.CgoEnabled = false

	var cmds []string
	cmds = append(cmds, o.Cmds...)
	if len(o.Uinit) > 0 {
		urootPkg, err := env.Package("github.com/u-root/u-root/integration")
		if err != nil {
			return nil, err
		}
		testDir := filepath.Join(urootPkg.Dir, "testdata")

		dirpath, err := ioutil.TempDir(testDir, "uinit")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(dirpath)

		if err := os.MkdirAll(filepath.Join(dirpath, "uinit"), 0755); err != nil {
			return nil, err
		}

		var realUinit [][]string
		for _, cmd := range o.Uinit {
			realUinit = append(realUinit, strings.Fields(cmd))
		}

		if err := ioutil.WriteFile(
			filepath.Join(dirpath, "uinit", "uinit.go"),
			[]byte(fmt.Sprintf(template, realUinit)),
			0755); err != nil {
			return nil, err
		}
		cmds = append(cmds, path.Join("github.com/u-root/u-root/integration/testdata", filepath.Base(dirpath), "uinit"))
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
		Env: env,
		Commands: []uroot.Commands{
			{
				Builder:  builder.BusyBox,
				Packages: append([]string{"github.com/u-root/u-root/cmds/*"}, cmds...),
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

	var logFile *os.File
	if len(o.LogFile) != 0 {
		logFile, err = os.Create(o.LogFile)
		if err != nil {
			return nil, fmt.Errorf("could not create log file: %v", err)
		}
	}

	// Start QEMU
	return &qemu.QEMU{
		Initramfs:    outputFile,
		Kernel:       bzImage,
		SerialOutput: logFile,
		SharedDir:    tmpDir,
	}, nil
}
