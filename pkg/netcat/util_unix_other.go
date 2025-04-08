// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows && !tamago && !((arm64 || riscv64) && linux)

package netcat

import "syscall"

// dup2 wraps syscall.Dup2 on anything POSIX-like that is neither ARM64 nor
// RISC-V64 Linux. (ARM64 and RISC-V64 Linux only offer syscall.Dup3; see e.g.
// <https://github.com/golang/go/issues/11981>.)
func dup2(oldfd, newfd int) error {
	return syscall.Dup2(oldfd, newfd)
}
