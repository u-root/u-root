// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
)

// Returns temporary directory and QEMU instance.
//
// - `uinitName` is the name of a directory containing uinit found at
//   `github.com/u-root/u-root/integration/testdata`.
func testWithQEMU(t *testing.T, uinitName string, extraArgs []string) (string, *qemu.QEMU) {
	if _, ok := os.LookupEnv("UROOT_QEMU"); !ok {
		t.Skip("test is skipped unless UROOT_QEMU is set")
	}
	if _, ok := os.LookupEnv("UROOT_KERNEL"); !ok {
		t.Skip("test is skipped unless UROOT_KERNEL is set")
	}

	// TempDir
	tmpDir, err := ioutil.TempDir("", "uroot-integration")
	if err != nil {
		t.Fatal(err)
	}

	// Env
	env := golang.Default()
	env.CgoEnabled = false

	// Builder
	builder, err := uroot.GetBuilder("bb")
	if err != nil {
		t.Fatal(err)
	}

	// Packages
	pkgs, err := uroot.DefaultPackageImports(env)
	if err != nil {
		t.Fatal(err)
	}
	pkgs = append(pkgs, path.Join("github.com/u-root/u-root/integration/testdata", uinitName, "uinit"))

	// Archiver
	archiver, err := uroot.GetArchiver("cpio")
	if err != nil {
		t.Fatal(err)
	}

	// OutputFile
	outputFile := filepath.Join(tmpDir, "initramfs.cpio")
	w, err := archiver.OpenWriter(outputFile, "", "")
	if err != nil {
		t.Fatal(err)
	}

	// Build u-root
	opts := uroot.Opts{
		TempDir: tmpDir,
		Env:     env,
		Commands: []uroot.Commands{
			{
				Builder:  builder,
				Packages: pkgs,
			},
		},
		Archiver:     archiver,
		OutputFile:   w,
		DefaultShell: "rush",
	}
	if err := uroot.CreateInitramfs(opts); err != nil {
		t.Fatal(err)
	}

	// Copy kernel to tmpDir.
	bzImage := filepath.Join(tmpDir, "bzImage")
	if err := cp.Copy(os.Getenv("UROOT_KERNEL"), bzImage); err != nil {
		t.Fatal(err)
	}

	// Expose the temp directory to QEMU as /dev/sda1
	extraArgs = append(extraArgs, "-drive", "file=fat:ro:"+tmpDir+",if=none,id=tmpdir")
	extraArgs = append(extraArgs, "-device", "ich9-ahci,id=ahci")
	extraArgs = append(extraArgs, "-device", "ide-drive,drive=tmpdir,bus=ahci.0")

	// Start QEMU
	q := &qemu.QEMU{
		InitRAMFS: outputFile,
		Kernel:    bzImage,
		ExtraArgs: extraArgs,
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
