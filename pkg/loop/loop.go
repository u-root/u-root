// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"github.com/u-root/u-root/pkg/mount"
	"golang.org/x/sys/unix"
)

// Loop implements mount.Mount
type Loop struct {
	Dev     string
	Source  string
	Dir     string
	FStype  string
	Flags   uintptr
	Data    string
	Mounted bool
}

// New initializes a Loop struct and allocates a loodevice to it.
func New(source, target, fstype string, flags uintptr, data string) (mount.Mounter, error) {
	devicename, err := FindDevice()
	if err != nil {
		return nil, err
	}
	if err := SetFdFiles(devicename, source); err != nil {
		return nil, err
	}
	l := &Loop{Dev: devicename, Dir: target, Source: source, FStype: fstype, Flags: flags, Data: data}
	return l, nil
}

// Mount mounts the provided source file, with type fstype, and flags and data options
// (which are usually 0 and ""), using any available loop device.
func (l *Loop) Mount() error {
	if err := unix.Mount(l.Dev, l.Dir, l.FStype, l.Flags, l.Data); err != nil {
		return err
	}
	l.Mounted = true
	return nil
}
