// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"golang.org/x/sys/unix"
)

var (
	fileSystemUID, fileSystemGID int
)

func dropPrivs() error {
	if fileSystemUID = unix.Getuid(); fileSystemUID == 0 {
		return nil
	}
	fileSystemGID = unix.Getgid()
	if err := unix.Setfsuid(fileSystemUID); err != nil {
		return err
	}
	return unix.Setfsgid(fileSystemGID)
}

func restorePrivs() {
	if os.Getuid() == 0 {
		return
	}
	// We're exiting, if there's an error, not much to do.
	unix.Setfsuid(fileSystemUID)
	unix.Setfsgid(fileSystemGID)
}

func preMount() error {
	// I guess this umask is the thing to do.
	unix.Umask(033)
	return nil
}
