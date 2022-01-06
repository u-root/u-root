// Copyright (c) 2017, Andrey Nering <andrey.nering@gmail.com>
// See LICENSE for licensing information

//go:build !windows
// +build !windows

package interp

import (
	"os"
	"os/user"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

func mkfifo(path string, mode uint32) error {
	return unix.Mkfifo(path, mode)
}

// hasPermissionToDir returns if the OS current user has execute permission
// to the given directory
func hasPermissionToDir(info os.FileInfo) bool {
	user, err := user.Current()
	if err != nil {
		return true
	}
	uid, _ := strconv.Atoi(user.Uid)
	// super-user
	if uid == 0 {
		return true
	}

	st, _ := info.Sys().(*syscall.Stat_t)
	if st == nil {
		return true
	}
	perm := info.Mode().Perm()
	// user (u)
	if perm&0o100 != 0 && st.Uid == uint32(uid) {
		return true
	}

	gid, _ := strconv.Atoi(user.Gid)
	// other users in group (g)
	if perm&0o010 != 0 && st.Uid != uint32(uid) && st.Gid == uint32(gid) {
		return true
	}
	// remaining users (o)
	if perm&0o001 != 0 && st.Uid != uint32(uid) && st.Gid != uint32(gid) {
		return true
	}

	return false
}
