// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package namespace

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	// REPL Replace the old file by the new one.
	// Henceforth, an evaluation of old will be translated to the new file.
	// If they are directories (for mount, this condition is true by definition),
	// old becomes a union directory consisting of one directory (the new file).
	REPL mountflag = syscall.MREPL
	// BEFORE Both the old and new files must be directories.
	// Add the constituent files of the new directory to the
	// union directory at old so its contents appear first in the union.
	// After an BEFORE bind or mount, the new directory will be
	// searched first when evaluating file names in the union directory.
	BEFORE mountflag = syscall.MBEFORE
	// AFTER Like MBEFORE but the new directory goes at the end of the union.
	AFTER mountflag = syscall.MAFTER
	// CREATE flag that can be OR'd with any of the above.
	// When a create system call (see open(2)) attempts to create in a union directory,
	// and the file does not exist, the elements of the union are searched in order until
	// one is found with CREATE set. The file is created in that directory;
	// if that attempt fails, the create fails.
	CREATE mountflag = syscall.MCREATE
	// CACHE flag, valid for mount only, turns on caching for files made available by the mount.
	// By default, file contents are always retrieved from the server.
	// With caching enabled, the kernel may instead use a local cache
	// to satisfy read(5) requests for files accessible through this mount point.
	CACHE mountflag = syscall.MCACHE
)

const (
	// These are copied over from the syscall pkg for plan9 https://go.plan9.io/pkg/syscall/

	// BIND is the plan9 bind syscall. https://9p.io/magic/man2html/2/bind
	BIND syzcall = syscall.SYS_BIND
	// CHDIR is the plan9 bind syscall. https://9p.io/magic/man2html/2/chdir
	CHDIR syzcall = syscall.SYS_CHDIR
	// UNMOUNT is the plan9 unmount syscall. https://9p.io/magic/man2html/2/bind
	UNMOUNT syzcall = syscall.SYS_UNMOUNT
	// MOUNT is the plan9 MOUNT syscall. https://9p.io/magic/man2html/2/bind
	MOUNT syzcall = syscall.SYS_MOUNT
	// RFORK is the plan9 rfork() syscall. https://9p.io/magic/man2html/2/fork
	// used to perform clear
	RFORK syzcall = syscall.SYS_RFORK
	// IMPORT is not a syscall. https://9p.io/magic/man2html/4/import
	IMPORT syzcall = syscall.SYS_EXEC // this is kinda handvavy cause this isn't an actual syscall
	// but right now it shells out to actual import.
	// TODO(sevki): implement the actual import
	// INCLUDE is not a syscall
	INCLUDE syzcall = syscall.SYS_OPEN // this is also handvavy
)

var DefaultNamespace = &namespace{}

type namespace struct{}

func (n *namespace) Bind(new string, old string, flag mountflag) error {
	return syscall.Bind(new, old, int(flag))
}
func (n *namespace) Chdir(dir string) error { return syscall.Chdir(dir) }

// Unmount unmounts
func (n *namespace) Unmount(new string, old string) error { return syscall.Unmount(new, old) }

// Clear clears the name space with rfork(RFCNAMEG).
func (n *namespace) Clear() error {
	r1, _, _ := syscall.RawSyscall(syscall.SYS_RFORK, uintptr(syscall.RFCNAMEG), 0, 0)
	if r1 != 0 {
		if int32(r1) == -1 {
			return errors.New(errstr())
		}
		// parent; return PID
		return nil
	}
	return nil
}

// Import imports a name space from a remote system
// This currently shells out to import.
//
// TODO(sevki): implement native import in cmds
// https://github.com/Harvey-OS/harvey/blob/master/sys/src/cmd/import.c
func (n *namespace) Import(host string, remotepath string, mountpoint string, f mountflag) error {
	flag := mountflag(f)
	args := []string{host}
	if remotepath != "" {
		args = append(args, remotepath)
	}
	args = append(args, mountpoint)
	flg := ""
	if flag&AFTER != 0 {
		flg += "a"
	}
	if flag&BEFORE != 0 {
		flg += "b"
	}
	if flag&CREATE != 0 {
		flg += "c"
	}
	if len(flg) > 0 {
		args = append([]string{flg}, args...)
	}
	return syscall.Exec("import", args, nil)
}

// Mount opens a fd with the server name and mounts the open fd to
// old
func (n *namespace) Mount(servername string, old, spec string, flag mountflag) error {
	fd, err := syscall.Open(servername, syscall.O_RDWR)
	if err != nil {
		return fmt.Errorf("open failed: %v", err)
	}
	return syscall.Mount(fd, -1, old, int(flag), spec)
}

func errstr() string {
	var buf [syscall.ERRMAX]byte

	syscall.RawSyscall(syscall.SYS_ERRSTR, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)

	buf[len(buf)-1] = 0
	return cstring(buf[:])
}

func cstring(s []byte) string {
	for i := range s {
		if s[i] == 0 {
			return string(s[0:i])
		}
	}
	return string(s)
}
