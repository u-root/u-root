// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"fmt"
	"os"
	"syscall"
)

const (
	/*
	 * IOCTL commands --- we will commandeer 0x4C ('L')
	 */
	LOOP_SET_CAPACITY = 0x4C07
	LOOP_CHANGE_FD    = 0x4C06
	LOOP_GET_STATUS64 = 0x4C05
	LOOP_SET_STATUS64 = 0x4C04
	LOOP_GET_STATUS   = 0x4C03
	LOOP_SET_STATUS   = 0x4C02
	LOOP_CLR_FD       = 0x4C01
	LOOP_SET_FD       = 0x4C00
	LO_NAME_SIZE      = 64
	LO_KEY_SIZE       = 32
	/* /dev/loop-control interface */
	LOOP_CTL_ADD      = 0x4C80
	LOOP_CTL_REMOVE   = 0x4C81
	LOOP_CTL_GET_FREE = 0x4C82
)

// FindDevice finds an unused loop device.
func FindDevice() (name string, err error) {
	cfd, err := os.OpenFile("/dev/loop-control", os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer cfd.Close()

	number, err := CtlGetFree(cfd.Fd())
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/dev/loop%d", number), nil
}

// ClearFd clears the loop device associated with filedescriptor fd.
func ClearFd(fd uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, LOOP_CLR_FD, 0); err != 0 {
		return err
	}

	return nil
}

// CtlGetFree finds a free loop device querying the loop control device pointed
// by fd. It returns the number of the free loop device /dev/loopX
func CtlGetFree(fd uintptr) (uintptr, error) {
	number, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, LOOP_CTL_GET_FREE, 0)
	if err != 0 {
		return 0, err
	}
	return number, nil
}

// SetFd associates a loop device pointed by lfd with a regular file pointed by ffd.
func SetFd(lfd, ffd uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, lfd, LOOP_SET_FD, ffd)
	if err != 0 {
		return err
	}

	return nil
}

// SetFdFiles associates loop device "devicename" with regular file "filename"
func SetFdFiles(devicename, filename string) error {
	mode := os.O_RDWR
	file, err := os.OpenFile(filename, mode, 0644)
	if err != nil {
		mode = os.O_RDONLY
		file, err = os.OpenFile(filename, mode, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	device, err := os.OpenFile(devicename, mode, 0644)
	if err != nil {
		return err
	}
	defer device.Close()

	return SetFd(device.Fd(), file.Fd())
}

// ClearFdFile clears the loop device "devicename"
func ClearFdFile(devicename string) error {
	device, err := os.Open(devicename)
	if err != nil {
		return err
	}
	defer device.Close()

	return ClearFd(device.Fd())
}
