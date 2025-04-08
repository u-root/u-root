// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (arm64 || riscv64) && linux

package netcat

import "syscall"

// dup2 wraps syscall.Dup3 on ARM64 and RISC-V64 Linux. The wrapping isn't 100%
// faithful, because syscall.Dup2 and syscall.Dup3 behave differently if oldfd
// equals newfd. But that difference shouldn't be visible in our use case.
func dup2(oldfd, newfd int) error {
	return syscall.Dup3(oldfd, newfd, 0)
}
