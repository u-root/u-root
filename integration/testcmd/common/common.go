// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"archive/tar"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/tarutil"
	"golang.org/x/sys/unix"
)

const (
	envUse9P            = "UROOT_USE_9P"
	envNoKernelCoverage = "UROOT_NO_KERNEL_COVERAGE"

	sharedDir          = "/testdata"
	kernelCoverageFile = "/testdata/kernel_coverage.tar"

	// https://wiki.qemu.org/Documentation/9psetup#msize recommends an
	// msize of at least 10MiB. Larger number might give better
	// performance. QEMU will print a warning if it is too small. Linux's
	// default is 8KiB which is way too small.
	msize9P = 10 * 1024 * 1024
)

// gcovFilter filters on all files ending with a gcda or gcno extension.
func gcovFilter(hdr *tar.Header) bool {
	if hdr.Typeflag == tar.TypeDir {
		hdr.Mode = 0o770
		return true
	}
	if (filepath.Ext(hdr.Name) == ".gcda" && hdr.Typeflag == tar.TypeReg) ||
		(filepath.Ext(hdr.Name) == ".gcno" && hdr.Typeflag == tar.TypeSymlink) {
		hdr.Mode = 0o660
		return true
	}
	return false
}

// CollectKernelCoverage saves the kernel coverage report to a tar file.
func CollectKernelCoverage() {
	if err := collectKernelCoverage(kernelCoverageFile); err != nil {
		log.Printf("Failed to collect kernel coverage: %v", err)
	}
}

func collectKernelCoverage(filename string) error {
	// Check if we are collecting kernel coverage.
	if os.Getenv(envNoKernelCoverage) == "1" {
		log.Print("Not collecting kernel coverage")
		return nil
	}
	gcovDir := "/sys/kernel/debug/gcov"
	if _, err := os.Stat(gcovDir); os.IsNotExist(err) {
		log.Printf("Not collecting kernel coverage because %q does not exist", gcovDir)
		return nil
	}
	if os.Getenv(envUse9P) != "1" {
		// 9p is required to rescue the file from the VM.
		return fmt.Errorf("not collecting kernel coverage because filesystem is not 9p")
	}

	// Mount debugfs.
	dfs := "/sys/kernel/debug"
	if err := os.MkdirAll(dfs, 0o666); err != nil {
		return fmt.Errorf("os.MkdirAll(%v): %v != nil", dfs, err)
	}
	if err := unix.Mount("debugfs", dfs, "debugfs", 0, ""); err != nil {
		return fmt.Errorf("failed to mount debugfs: %v", err)
	}

	// Copy out the kernel code coverage.
	log.Print("Collecting kernel coverage...")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := tarutil.CreateTar(f, []string{strings.TrimLeft(gcovDir, "/")}, &tarutil.Opts{
		Filters: []tarutil.Filter{gcovFilter},
		// Make sure the files are not stored absolute; otherwise, they
		// become difficult to extract safely.
		ChangeDirectory: "/",
	}); err != nil {
		f.Close()
		return err
	}
	// Sync to "disk" because we are about to shut down the kernel.
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("error syncing: %v", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing: %v", err)
	}
	return nil
}

// MountSharedDir mounts the directory shared with the VM test. A cleanup
// function is returned to unmount.
func MountSharedDir() (func(), error) {
	// Mount a disk and run the tests within.
	var (
		mp  *mount.MountPoint
		err error
	)

	if err := os.MkdirAll(sharedDir, 0o644); err != nil {
		return nil, err
	}

	if os.Getenv(envUse9P) == "1" {
		mp, err = mount.Mount("tmpdir", sharedDir, "9p", fmt.Sprintf("9P2000.L,msize=%d", msize9P), 0)
	} else {
		mp, err = mount.Mount("/dev/sda1", sharedDir, "vfat", "", unix.MS_RDONLY)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to mount test directory: %v", err)
	}
	return func() { mp.Unmount(0) }, nil
}
