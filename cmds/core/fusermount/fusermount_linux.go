// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

var fileSystemUID, fileSystemGID int

func dropPrivs() error {
	uid := unix.Getuid()
	if uid == 0 {
		return nil
	}

	var err error
	fileSystemUID, err = unix.SetfsuidRetUid(uid)
	if err != nil {
		return err
	}
	fileSystemGID, err = unix.SetfsgidRetGid(unix.Getgid())
	return err
}

func restorePrivs() error {
	if unix.Getuid() == 0 {
		return nil
	}
	// We're exiting, if there's an error, not much to do.
	if err := unix.Setfsuid(fileSystemUID); err != nil {
		return fmt.Errorf("restorePrivs Setfsuid: %w", err)
	}
	if err := unix.Setfsgid(fileSystemGID); err != nil {
		return fmt.Errorf("restorePrivs Setfsgid: %w", err)
	}

	return nil
}

func preMount() error {
	// I guess this umask is the thing to do.
	unix.Umask(0o33)
	return nil
}
