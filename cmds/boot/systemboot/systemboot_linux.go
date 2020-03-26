// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/ipmi"
)

const (
	_IOC_NVRAM = 'p'
)

// ioctl to clear whole CMOS/NV-RAM
var _NVRAM_INIT = ipmi.IO(_IOC_NVRAM, 0x40)

func cmosClear() error {
	f, err := os.OpenFile("/dev/nvram", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	var i = 0
	if errno := ipmi.Ioctl(uintptr(f.Fd()), _NVRAM_INIT, unsafe.Pointer(&i)); errno != 0 {
		return errno
	}
	return nil
}

func reboot() error {
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
