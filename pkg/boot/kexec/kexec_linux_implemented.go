// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,arm64

package kexec

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/sys/unix"
)

// FileLoad loads the given kernel as the new kernel with the given ramfs and
// cmdline.
//
// For some architectures (such as arm64), the vmlinuz kernel is not
// self-decompressing and the expectation is for the bootloader to decompress
// it. This automatically detects a gzip image and decompresses it.
//
// The kexec_file_load(2) syscall is x86-64 and arm64 only.
func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	var flags int
	var ramfsfd int
	if ramfs != nil {
		ramfsfd = int(ramfs.Fd())
	} else {
		flags |= unix.KEXEC_FILE_NO_INITRAMFS
	}

	if isGzip(kernel) {
		g, err := gzip.NewReader(kernel)
		if err != nil {
			return fmt.Errorf("could read kernel gzip header: %v", err)
		}
		decompressed, err := ioutil.TempFile("", "kernel")
		if err != nil {
			return fmt.Errorf("could not create temp file: %v", err)
		}
		defer os.Remove(decompressed.Name())
		if _, err := io.Copy(decompressed, g); err != nil {
			return fmt.Errorf("could not decompress kernel: %v", err)
		}
		if err := g.Close(); err != nil {
			return fmt.Errorf("could not close gunzip (bad checksum?): %v", err)
		}
		if err := decompressed.Sync(); err != nil {
			return fmt.Errorf("could not sync decompressed kernel: %v", err)
		}
		kernel = decompressed
	}

	if err := unix.KexecFileLoad(int(kernel.Fd()), ramfsfd, cmdline, flags); err != nil {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = %v", kernel.Fd(), ramfsfd, cmdline, flags, err)
	}
	return nil
}

// isGzip returns true if the file is gzip file format.
func isGzip(kernel *os.File) bool {
	gzipMagic := []byte{0x1F, 0x8B}
	magic := make([]byte, len(gzipMagic))
	n, err := kernel.ReadAt(magic, 0)
	return err == nil && n == len(magic) && bytes.Compare(gzipMagic, magic) == 0
}
