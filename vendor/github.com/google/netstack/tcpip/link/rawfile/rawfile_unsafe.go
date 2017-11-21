// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rawfile contains utilities for using the netstack with raw host
// files on Linux hosts.
package rawfile

import (
	"syscall"
	"unsafe"

	"github.com/google/netstack/tcpip"
)

// TODO: Placed here to avoid breakage caused by coverage
// instrumentation. Any, even unrelated, changes to this file should ensure
// that coverage still work. See bug for details.
//go:noescape
func blockingPoll(fds unsafe.Pointer, nfds int, timeout int64) (n int, err syscall.Errno)

// GetMTU determines the MTU of a network interface device.
func GetMTU(name string) (uint32, error) {
	fd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}

	defer syscall.Close(fd)

	var ifreq struct {
		name [16]byte
		mtu  int32
		_    [20]byte
	}

	copy(ifreq.name[:], name)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCGIFMTU, uintptr(unsafe.Pointer(&ifreq)))
	if errno != 0 {
		return 0, errno
	}

	return uint32(ifreq.mtu), nil
}

// NonBlockingWrite writes the given buffer to a file descriptor. It fails if
// partial data is written.
func NonBlockingWrite(fd int, buf []byte) *tcpip.Error {
	var ptr unsafe.Pointer
	if len(buf) > 0 {
		ptr = unsafe.Pointer(&buf[0])
	}

	_, _, e := syscall.RawSyscall(syscall.SYS_WRITE, uintptr(fd), uintptr(ptr), uintptr(len(buf)))
	if e != 0 {
		return TranslateErrno(e)
	}

	return nil
}

// NonBlockingWrite2 writes up to two byte slices to a file descriptor in a
// single syscall. It fails if partial data is written.
func NonBlockingWrite2(fd int, b1, b2 []byte) *tcpip.Error {
	// If the is no second buffer, issue a regular write.
	if len(b2) == 0 {
		return NonBlockingWrite(fd, b1)
	}

	// We have two buffers. Build the iovec that represents them and issue
	// a writev syscall.
	iovec := [...]syscall.Iovec{
		{
			Base: (*byte)(unsafe.Pointer(&b1[0])),
			Len:  uint64(len(b1)),
		},
		{
			Base: (*byte)(unsafe.Pointer(&b2[0])),
			Len:  uint64(len(b2)),
		},
	}

	_, _, e := syscall.RawSyscall(syscall.SYS_WRITEV, uintptr(fd), uintptr(unsafe.Pointer(&iovec[0])), 2)
	if e != 0 {
		return TranslateErrno(e)
	}

	return nil
}

// BlockingRead reads from a file descriptor that is set up as non-blocking. If
// no data is available, it will block in a poll() syscall until the file
// descirptor becomes readable.
func BlockingRead(fd int, b []byte) (int, *tcpip.Error) {
	for {
		n, _, e := syscall.RawSyscall(syscall.SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
		if e == 0 {
			return int(n), nil
		}

		event := struct {
			fd      int32
			events  int16
			revents int16
		}{
			fd:     int32(fd),
			events: 1, // POLLIN
		}

		_, e = blockingPoll(unsafe.Pointer(&event), 1, -1)
		if e != 0 && e != syscall.EINTR {
			return 0, TranslateErrno(e)
		}
	}
}

// BlockingReadv reads from a file descriptor that is set up as non-blocking and
// stores the data in a list of iovecs buffers. If no data is available, it will
// block in a poll() syscall until the file descirptor becomes readable.
func BlockingReadv(fd int, iovecs []syscall.Iovec) (int, *tcpip.Error) {
	for {
		n, _, e := syscall.RawSyscall(syscall.SYS_READV, uintptr(fd), uintptr(unsafe.Pointer(&iovecs[0])), uintptr(len(iovecs)))
		if e == 0 {
			return int(n), nil
		}

		event := struct {
			fd      int32
			events  int16
			revents int16
		}{
			fd:     int32(fd),
			events: 1, // POLLIN
		}

		_, e = blockingPoll(unsafe.Pointer(&event), 1, -1)
		if e != 0 && e != syscall.EINTR {
			return 0, TranslateErrno(e)
		}
	}
}
