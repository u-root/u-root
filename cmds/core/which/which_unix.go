// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris

package main

import (
	"os"

	"golang.org/x/sys/unix"
)

var (
	// Effective UID and GID. Technically, Linux checks the fsuid and
	// fsgid. But those are in 99.999% of cases the same as EUID and EGID.
	// Nobody uses that moldy stuff anymore.
	eUID = uint32(os.Geteuid())
	eGID = uint32(os.Getegid())
)

func stat(path string) (unix.Stat_t, error) {
	var s unix.Stat_t
	if err := unix.Stat(path, &s); err != nil {
		return unix.Stat_t{}, err
	}
	return s, nil
}

// Execute permission bits.
const (
	otherExec = 1
	groupExec = 1 << 3
	userExec  = 1 << 6
)

func canExecute(path string) bool {
	info, err := stat(path)
	if err != nil {
		return false
	}

	// This is a... first-order approximation.
	//
	// Somebody should write a Go capabilities library. Then we can check
	// CAP_DAC_OVERRIDE. Probably not really necessary, though. Who has
	// Unix shell users with select capabilities? Hopefully nobody.
	if eUID == 0 {
		return true
	} else if info.Uid == eUID {
		return info.Mode&userExec == userExec
	} else if info.Gid == eGID {
		return info.Mode&groupExec == groupExec
	}
	return info.Mode&otherExec == otherExec
}
