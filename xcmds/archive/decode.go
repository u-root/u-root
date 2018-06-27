// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Archive archives files. The VTOC is at the front; we're not modeling tape drives or
// streams as in tar and cpio. This will greatly speed up listing the archive,
// modifying it, and so on. We think.
// Why a new tool?
package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func doOneFile(f *os.File, v file) error {
	var err error
	s := v.Mode.String()
	fmt.Printf("%v: mode %v", v.Name, s)
	switch s[0] {
	case 'd':
		err = os.MkdirAll(v.Name, v.Mode)
	case '-':
		src := io.LimitReader(f, v.Size)
		dst, err := os.OpenFile(v.Name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, v.Mode)
		if err != nil {
			return err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
	case 'L':
		err = os.Symlink(v.Link, v.Name)
	case 'D':
		m := uint32(v.Mode.Perm())
		switch v.Mode & (os.ModeDevice | os.ModeCharDevice) {
		case os.ModeDevice:
			m = m | syscall.S_IFBLK
		default:
			m = m | syscall.S_IFCHR
		}
		err = syscall.Mknod(v.Name, m, int(v.Dev))
	default:
		// It's not an error to do hae an archive we things we can't do.
		//err = errors.New(fmt.Sprintf("Can't make %v", v))
		fmt.Printf("Can' do %v yet", v)
		return nil
	}
	if err != nil {
		return err
	}
	return os.Lchown(v.Name, v.Uid, v.Gid)
}

func decode(files ...string) error {
	for _, f := range files {
		fd, vtoc, err := loadVTOC(f)
		if err != nil {
			fmt.Printf("%v", err)
		}
		for _, v := range vtoc {
			if err := doOneFile(fd, v); err != nil {
				return err
			}
			debug("%v", v)
		}
	}
	return nil
}
