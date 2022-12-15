// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// FindDevice finds an unused loop device and returns its /dev/loopN path.
func FindDevice() (string, error) {
	cfd, err := os.OpenFile("/dev/loop-control", os.O_RDWR, 0o644)
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
	return unix.IoctlSetInt(fd, unix.LOOP_CLR_FD, 0)
}

// GetFree finds a free loop device /dev/loopN.
//
// fd must be a loop control device.
//
// It returns the number of the free loop device /dev/loopN.
// The _LOOP_CTL_GET_FREE does not follow the rules. Values
// of 0 or greater are the number of the device; less than
// zero is an error.
// So you can not use unix.IoctlGetInt as it assumes the return
// value is stored in a pointer in the normal style. Yuck.
func GetFree(fd int) (int, error) {
	return unix.IoctlRetInt(fd, unix.LOOP_CTL_GET_FREE)
}

// SetFD associates a loop device lfd with a regular file ffd.
func SetFD(lfd, ffd int) error {
	return unix.IoctlSetInt(lfd, unix.LOOP_SET_FD, ffd)
}

// SetFile associates loop device "devicename" with regular file "filename"
func SetFile(devicename, filename string) error {
	mode := os.O_RDWR
	file, err := os.OpenFile(filename, mode, 0o644)
	if err != nil {
		mode = os.O_RDONLY
		file, err = os.OpenFile(filename, mode, 0o644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	device, err := os.OpenFile(devicename, mode, 0o644)
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
