// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

// This file is inserted here so the lack of these variables/implmenetations doesn't break
// IDEs and tooling on other platforms.

package namespace

// DefaultNamespace is the default namespace
var DefaultNamespace = &unixnamespace{}

const (
	// REPL Replace the old file by the new one.
	// Henceforth, an evaluation of old will be translated to the new file.
	// If they are directories (for mount, this condition is true by definition),
	// old becomes a union directory consisting of one directory (the new file).
	REPL mountflag = 0x0000
	// BEFORE Both the old and new files must be directories.
	// Add the constituent files of the new directory to the
	// union directory at old so its contents appear first in the union.
	// After an BEFORE bind or mount, the new directory will be
	// searched first when evaluating file names in the union directory.
	BEFORE mountflag = 0x0001
	// AFTER Like MBEFORE but the new directory goes at the end of the union.
	AFTER mountflag = 0x0002
	// CREATE flag that can be OR'd with any of the above.
	// When a create system call (see open(2)) attempts to create in a union directory,
	// and the file does not exist, the elements of the union are searched in order until
	// one is found with CREATE set. The file is created in that directory;
	// if that attempt fails, the create fails.
	CREATE mountflag = 0x0004
	// CACHE flag, valid for mount only, turns on caching for files made available by the mount.
	// By default, file contents are always retrieved from the server.
	// With caching enabled, the kernel may instead use a local cache
	// to satisfy read(5) requests for files accessible through this mount point.
	CACHE mountflag = 0x0010
)

const (
	// These are copied over from the syscall pkg for plan9 https://golang.org/pkg/syscall/?GOOS=plan9

	// BIND is the plan9 bind syscall. https://9p.io/magic/man2html/2/bind
	BIND syzcall = 2
	// CHDIR is the plan9 bind syscall. https://9p.io/magic/man2html/2/chdir
	CHDIR syzcall = 3
	// UNMOUNT is the plan9 unmount syscall. https://9p.io/magic/man2html/2/bind
	UNMOUNT syzcall = 35
	// MOUNT is the plan9 MOUNT syscall. https://9p.io/magic/man2html/2/bind
	MOUNT syzcall = 46
	// RFORK is the plan9 rfork() syscall. https://9p.io/magic/man2html/2/fork
	// used to perform clear
	RFORK syzcall = 19
	// IMPORT is not a syscall. https://9p.io/magic/man2html/4/import
	IMPORT syzcall = 7
	// INCLUDE is not a syscall
	INCLUDE syzcall = 14
)

type unixnamespace struct{}

// Bind binds new on old.
func (ns *unixnamespace) Bind(name string, old string, flag mountflag) error {
	panic("not implemented") // TODO: Implement
}

// Mount mounts servename on old.
func (ns *unixnamespace) Mount(servername string, old string, spec string, flag mountflag) error {
	panic("not implemented") // TODO: Implement
}

// Unmount unmounts new from old, or everything mounted on old if new is missing.
func (ns *unixnamespace) Unmount(name string, old string) error {
	panic("not implemented") // TODO: Implement
}

// Clear clears the name space with rfork(RFCNAMEG).
func (ns *unixnamespace) Clear() error {
	panic("not implemented") // TODO: Implement
}

// Chdir changes the working directory to dir.
func (ns *unixnamespace) Chdir(dir string) error {
	panic("not implemented") // TODO: Implement
}

// Import imports a name space from a remote system
func (ns *unixnamespace) Import(host string, remotepath string, mountpoint string, flag mountflag) error {
	panic("not implemented") // TODO: Implement
}
