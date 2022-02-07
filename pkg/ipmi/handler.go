// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type syscalls interface {
	syscall(uintptr, uintptr, uintptr, uintptr) (uintptr, uintptr, unix.Errno)
	fileSyscallConn(f *os.File) (syscall.RawConn, error)
	fileSetReadDeadline(f *os.File, t time.Duration) error
	connRead(f func(fd uintptr) bool, conn syscall.RawConn) error
}

type realSyscalls struct{}

func (r *realSyscalls) syscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, unix.Errno) {
	return unix.Syscall(trap, a1, a2, a3)
}

func (r *realSyscalls) fileSyscallConn(f *os.File) (syscall.RawConn, error) {
	return f.SyscallConn()
}

func (r *realSyscalls) fileSetReadDeadline(f *os.File, t time.Duration) error {
	return f.SetReadDeadline(time.Now().Add(t))
}

func (r *realSyscalls) connRead(f func(fd uintptr) bool, conn syscall.RawConn) error {
	return conn.Read(f)
}
