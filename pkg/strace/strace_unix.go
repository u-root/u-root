// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strace

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func Wait(pid int) (unix.WaitStatus, error) {
	var w syscall.WaitStatus
	_, err := syscall.Wait4(pid, &w, 0, nil)
	uw := unix.WaitStatus(w)
	return uw, err
}
