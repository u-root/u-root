// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package quimage provides a Go API for creating QEMU VMs with u-root uimage initramfses.
//
// Environment variables:
//
//	VMTEST_INITRAMFS_OVERRIDE (when set, use instead of building an initramfs)
package quimage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
	"github.com/u-root/uio/llog"
)

// ErrOutputFileSpecified is returned when uimage.Opts are supplied that already
// have an initramfs file.
var ErrOutputFileSpecified = errors.New("initramfs output file must be left unspecified")

// WithUimage builds the specified initramfs and attaches it to the QEMU VM.
//
// When VMTEST_INITRAMFS_OVERRIDE is set, it foregoes building an initramfs and
// uses the initramfs path in the env variable.
//
// The arch used to build the initramfs is derived by default from the arch set
// in qemu.Options, which is either explicitly set, VMTEST_ARCH, or if unset,
// runtime.GOARCH (the host GOARCH).
func WithUimage(l *llog.Logger, initrdPath string, mods ...uimage.Modifier) qemu.Fn {
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		if override := os.Getenv("VMTEST_INITRAMFS_OVERRIDE"); len(override) > 0 {
			opts.Initramfs = override
			return nil
		}

		mods = append([]uimage.Modifier{
			uimage.WithEnv(
				golang.DisableCGO(),
				golang.WithGOARCH(string(opts.Arch())),
			),
			uimage.WithCPIOOutput(initrdPath),
		}, mods...)
		if err := uimage.Create(l, mods...); err != nil {
			return fmt.Errorf("error creating initramfs: %w", err)
		}
		opts.Initramfs = initrdPath
		return nil
	}
}

// WithUimageT adds an initramfs to the VM using a logger for t and
// placing the initramfs in a test-created temp dir.
func WithUimageT(t testing.TB, mods ...uimage.Modifier) qemu.Fn {
	l := llog.Test(t)
	initrdPath := filepath.Join(testtmp.TempDir(t), "initramfs.cpio")
	return WithUimage(l, initrdPath, append(mods, uimage.WithTempDir(testtmp.TempDir(t)))...)
}
