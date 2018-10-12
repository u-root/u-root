// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/builder"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
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

// Returns temporary directory and QEMU instance.
func testWithQEMU(t *testing.T, o options) (string, *qemu.QEMU) {
	if _, ok := os.LookupEnv("UROOT_QEMU"); !ok {
		t.Skip("test is skipped unless UROOT_QEMU is set")
	}
	if _, ok := os.LookupEnv("UROOT_KERNEL"); !ok {
		t.Skip("test is skipped unless UROOT_KERNEL is set")
	}

	// Create or reuse a temporary directory.
	tmpDir := o.tmpDir
	if tmpDir == "" {
		var err error
		tmpDir, err = ioutil.TempDir("", "uroot-integration")
		if err != nil {
			t.Fatal(err)
		}
	}

	// Env
	env := golang.Default()
	env.CgoEnabled = false

	// OutputFile
	outputFile := filepath.Join(tmpDir, "initramfs.cpio")
	w, err := initramfs.CPIO.OpenWriter(outputFile, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Build u-root
	opts := uroot.Opts{
		Env: env,
		Commands: []uroot.Commands{
			{
				Builder: builder.BusyBox,
				Packages: []string{
					"github.com/u-root/u-root/cmds/*",
					path.Join("github.com/u-root/u-root/integration/testdata", o.uinitName, "uinit"),
				},
			},
		},
		TempDir:      tmpDir,
		BaseArchive:  uroot.DefaultRamfs.Reader(),
		OutputFile:   w,
		InitCmd:      "init",
		DefaultShell: "elvish",
	}
	logger := log.New(os.Stderr, "", log.LstdFlags)
	if err := uroot.CreateInitramfs(logger, opts); err != nil {
		t.Fatal(err)
	}

	// Copy kernel to tmpDir.
	bzImage := filepath.Join(tmpDir, "bzImage")
	if err := cp.Copy(os.Getenv("UROOT_KERNEL"), bzImage); err != nil {
		t.Fatal(err)
	}

	// Expose the temp directory to QEMU as /dev/sda1
	args := []string{}
	args = append(args, "-drive", "file=fat:ro:"+tmpDir+",if=none,id=tmpdir")
	args = append(args, "-device", "ich9-ahci,id=ahci")
	args = append(args, "-device", "ide-drive,drive=tmpdir,bus=ahci.0")
	args = append(args, o.extraArgs...)

	// Create file for serial logs.
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("could not create serial log directory: %v", err)
	}
	logFile, err := os.Create(path.Join(logDir, o.uinitName+".log"))
	if err != nil {
		t.Fatalf("could not create log file: %v", err)
	}

	// Start QEMU
	q := &qemu.QEMU{
		InitRAMFS:    outputFile,
		Kernel:       bzImage,
		ExtraArgs:    args,
		SerialOutput: logFile,
	}
	t.Logf("command line:\n%s", q.CmdLineQuoted())
	if err := q.Start(); err != nil {
		t.Fatal("could not spawn QEMU: ", err)
	}
	return tmpDir, q
}

func cleanup(t *testing.T, tmpDir string, q *qemu.QEMU) {
	q.Close()
	if t.Failed() {
		t.Log("keeping temp dir: ", tmpDir)
	} else {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temporary directory %s", tmpDir)
		}
	}
}
