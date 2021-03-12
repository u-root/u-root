// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !arm,!arm64

package memio

// Read reads data from physical memory at address addr. On x86 platforms,
// this uses the seek+read syscalls. On arm platforms, this uses mmap.
func Read(addr int64, data UintN) error {
	return pathRead(memPath, addr, data)
}

// Write writes data to physical memory at address addr. On x86 platforms, this
// uses the seek+read syscalls. On arm platforms, this uses mmap.
func Write(addr int64, data UintN) error {
	return pathWrite(memPath, addr, data)
}
