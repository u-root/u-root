// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const (
	// Loop ioctl commands --- we will commandeer 0x4C ('L')
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

	// /dev/loop-control interface
	LOOP_CTL_ADD      = 0x4C80
	LOOP_CTL_REMOVE   = 0x4C81
	LOOP_CTL_GET_FREE = 0x4C82
)

// FindDevice finds an unused loop device and returns its /dev/loopN path.
func FindDevice() (string, error) {
	cfd, err := os.OpenFile("/dev/loop-control", os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer cfd.Close()

	number, err := GetFree(int(cfd.Fd()))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/dev/loop%d", number), nil
}

// ClearFD clears the loop device associated with file descriptor fd.
func ClearFD(fd int) error {
	return unix.IoctlSetInt(fd, LOOP_CLR_FD, 0)
}

// GetFree finds a free loop device querying the loop control device pointed
// by fd. It returns the number of the free loop device /dev/loopX.
func GetFree(fd int) (int, error) {
	return unix.IoctlGetInt(fd, LOOP_CTL_GET_FREE)
}

// SetFD associates a loop device lfd with a regular file ffd.
func SetFD(lfd, ffd int) error {
	return unix.IoctlSetInt(lfd, LOOP_SET_FD, ffd)
}

// SetFile associates loop device "devicename" with regular file "filename"
func SetFile(devicename, filename string) error {
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

	return SetFD(int(device.Fd()), int(file.Fd()))
}

// ClearFile clears the fd association of the loop device "devicename".
func ClearFile(devicename string) error {
	device, err := os.Open(devicename)
	if err != nil {
		return err
	}
	defer device.Close()

	return ClearFD(int(device.Fd()))
}
