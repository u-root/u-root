// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64

package linux

import (
	"io"
	"os"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"golang.org/x/sys/unix"
)

// KexecLoad is not implemented for platforms other than amd64 and arm64.
func KexecLoad(kernel, ramfs *os.File, cmdline string, dtb io.ReaderAt, reservations kexec.Ranges) error {
	return unix.ENOSYS
}
