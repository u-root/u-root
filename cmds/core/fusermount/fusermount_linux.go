// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "golang.org/x/sys/unix"

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

func restorePrivs() {
	if unix.Getuid() == 0 {
		return
	}
	// We're exiting, if there's an error, not much to do.
	unix.Setfsuid(fileSystemUID)
	unix.Setfsgid(fileSystemGID)
}

func preMount() error {
	// I guess this umask is the thing to do.
	unix.Umask(0o33)
	return nil
}
