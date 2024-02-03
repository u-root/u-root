// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
)

// ShareGOCOVERDIR shares VMTEST_GOCOVERDIR with the guest if it's available in the
// environment.
//
// Call guest.GOCOVERDIR to set up the directory in the guest.
func ShareGOCOVERDIR() Opt {
	return func(t testing.TB, v *VMOptions) error {
		goCov := os.Getenv("VMTEST_GOCOVERDIR")
		if goCov == "" {
			return nil
		}
		v.QEMUOpts = append(v.QEMUOpts,
			qemu.P9Directory(goCov, "gocov"),
			qemu.WithAppendKernel("VMTEST_GOCOVERDIR=gocov"),
		)
		return nil
	}
}

// CollectKernelCoverage collects kernel coverage files for each test to
// VMTEST_KERNEL_COVERAGE_DIR/{testName}/{instance}, where instance is a number
// starting at 0.
//
// If VMTEST_KERNEL_COVERAGE_DIR is empty, collection is skipped.
func CollectKernelCoverage() Opt {
	return func(t testing.TB, v *VMOptions) error {
		if os.Getenv("VMTEST_KERNEL_COVERAGE_DIR") == "" {
			t.Logf("Skipping kernel coverage collection since VMTEST_KERNEL_COVERAGE_DIR is not set")
			return nil
		}
		coverageDir := os.Getenv("VMTEST_KERNEL_COVERAGE_DIR")
		if err := os.MkdirAll(coverageDir, 0o770); err != nil {
			return fmt.Errorf("could not create VMTEST_KERNEL_COVERAGE_DIR: %v", err)
		}

		sharedDir := testtmp.TempDir(t)
		v.QEMUOpts = append(v.QEMUOpts,
			qemu.P9Directory(sharedDir, "kcoverage"),
			qemu.WithAppendKernel("VMTEST_KCOVERAGE_TAG=kcoverage"),
			qemu.WithTask(qemu.Cleanup(func() error {
				if err := saveCoverage(t, filepath.Join(sharedDir, "kernel_coverage.tar"), coverageDir); err != nil {
					return fmt.Errorf("error saving kernel coverage: %v", err)
				}
				return nil
			})),
		)
		return nil
	}
}

// Keeps track of the number of instances per test so we do not overlap
// coverage reports.
var instance = map[string]int{}

func saveCoverage(t testing.TB, coverageFile, coverageDir string) error {
	// Coverage may not have been collected, for example if the kernel is
	// not built with CONFIG_GCOV_KERNEL.
	if fi, err := os.Stat(coverageFile); err != nil {
		return fmt.Errorf("could not access result kernel coverage file (is your kernel built with CONFIG_GCOV_KERNEL?): %w", err)
	} else if !fi.Mode().IsRegular() {
		return fmt.Errorf("kernel coverage file is not a regular file")
	}

	// Move coverage to common directory.
	uniqueCoveragePath := filepath.Join(coverageDir, t.Name(), fmt.Sprintf("%d", instance[t.Name()]))
	instance[t.Name()]++
	if err := os.MkdirAll(uniqueCoveragePath, 0o770); err != nil {
		return err
	}

	dest := filepath.Join(uniqueCoveragePath, filepath.Base(coverageFile))
	t.Logf("Kernel coverage file for this test: %s", dest)
	return os.Rename(coverageFile, dest)
}
