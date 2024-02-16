// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qcoverage allows collecting kernel and Go integration test coverage.
package qcoverage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
)

// ShareGOCOVERDIR shares VMTEST_GOCOVERDIR with the guest if it's available in
// the environment.
//
// Use the vmmount command to mount the directory before calling any commands
// that should have GOCOVERDIR coverage, or mount a virtio-9p directory with
// tag "gocov" at /mount/9p/gocov.
func ShareGOCOVERDIR() qemu.Fn {
	goCov := os.Getenv("VMTEST_GOCOVERDIR")
	if goCov == "" {
		return nil
	}
	return qemu.All(
		qemu.P9Directory(goCov, "gocov"),
		qemu.WithAppendKernel("GOCOVERDIR=/mount/9p/gocov"),
	)
}

// CollectKernelCoverage collects kernel coverage files for each test to
// VMTEST_KERNEL_COVERAGE_DIR/{testName}/{instance}, where instance is a number
// starting at 0.
//
// If VMTEST_KERNEL_COVERAGE_DIR is empty, collection is skipped.
func CollectKernelCoverage(tb testing.TB) qemu.Fn {
	if os.Getenv("VMTEST_KERNEL_COVERAGE_DIR") == "" {
		tb.Logf("Skipping kernel coverage collection since VMTEST_KERNEL_COVERAGE_DIR is not set")
		return nil
	}

	coverageDir := os.Getenv("VMTEST_KERNEL_COVERAGE_DIR")
	if err := os.MkdirAll(coverageDir, 0o770); err != nil {
		tb.Fatalf("Could not create VMTEST_KERNEL_COVERAGE_DIR: %v", err)
	}

	sharedDir := testtmp.TempDir(tb)
	return qemu.All(
		qemu.P9Directory(sharedDir, "kcoverage"),
		qemu.WithTask(qemu.Cleanup(func() error {
			if err := saveCoverage(tb, filepath.Join(sharedDir, "kernel_coverage.tar"), coverageDir); err != nil {
				return fmt.Errorf("error saving kernel coverage: %v", err)
			}
			return nil
		})),
	)
}

// Keeps track of the number of instances per test so we do not overlap
// coverage reports.
var instance = map[string]int{}

func saveCoverage(tb testing.TB, coverageFile, coverageDir string) error {
	// Coverage may not have been collected, for example if the kernel is
	// not built with CONFIG_GCOV_KERNEL.
	if fi, err := os.Stat(coverageFile); err != nil {
		return fmt.Errorf("could not access result kernel coverage file (is your kernel built with CONFIG_GCOV_KERNEL?): %w", err)
	} else if !fi.Mode().IsRegular() {
		return fmt.Errorf("kernel coverage file is not a regular file")
	}

	// Move coverage to common directory.
	uniqueCoveragePath := filepath.Join(coverageDir, tb.Name(), fmt.Sprintf("%d", instance[tb.Name()]))
	instance[tb.Name()]++
	if err := os.MkdirAll(uniqueCoveragePath, 0o770); err != nil {
		return err
	}

	dest := filepath.Join(uniqueCoveragePath, filepath.Base(coverageFile))
	tb.Logf("Kernel coverage file for this test: %s", dest)
	return os.Rename(coverageFile, dest)
}
