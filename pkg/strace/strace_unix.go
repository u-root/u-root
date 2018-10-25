// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strace

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// Wait will wait for the specified pid using Wait4.
// Callers may specify the full range of values
// as specified in the waipid man page, though
// we typically use only -1 or a valid pid.
func Wait(wpid int) (int, unix.WaitStatus, error) {
	var w syscall.WaitStatus
	pid, err := syscall.Wait4(wpid, &w, 0, nil)
	uw := unix.WaitStatus(w)
	return pid, uw, err
}
