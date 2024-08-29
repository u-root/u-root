// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"os"

	"github.com/vtolstov/go-ioctl"
	"golang.org/x/sys/unix"
)

const (
	_IOC_NVRAM = 'p'
)

// ioctl to clear whole CMOS/NV-RAM
var _NVRAM_INIT = ioctl.IO(_IOC_NVRAM, 0x40)

func cmosClear() error {
	f, err := os.OpenFile("/dev/nvram", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	return unix.IoctlSetPointerInt(int(f.Fd()), uint(_NVRAM_INIT), 0)
}

func reboot() error {
	return unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}
