// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package strace

type iovec struct {
	P Addr   /* Starting address */
	S uint32 /* Number of bytes to transfer */
}

// Addr is an address for use in strace I/O
type Addr uintptr
