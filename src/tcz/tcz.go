// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

const (
	/*
	 * IOCTL commands --- we will commandeer 0x4C ('L')
	 */
	LOOP_SET_CAPACITY = 0x4C07
	LOOP_CHANGE_FD    = 0x4C06
	LOOP_GET_STATUS64 = 0x4C05
	LOOP_SET_STATUS64 = 0x4C04
	LOOP_GET_STATUS   = 0x4C03
	LOOP_SET_STATUS   = 0x4C02
	LOOP_CLR_FD       = 0x4C01
	LOOP_SET_FD       = 0x4C00
	LO_NAME_SIZE      = 64
	LO_KEY_SIZE       = 32
	/* /dev/loop-control interface */
	LOOP_CTL_ADD      = 0x4C80
	LOOP_CTL_REMOVE   = 0x4C81
	LOOP_CTL_GET_FREE = 0x4C82
	SYS_ioctl         = 16
)

const tcz = "/tinycorelinux.net/5.x/x86_64/tcz"

var l = log.New(os.Stdout, "tcz: ", 0)

// consider making this a goroutine which pushes the string down the channel.
func findloop() (name string, err error) {
	cfd, err := syscall.Open("/dev/loop-control", syscall.O_RDWR, 0)
	if err != nil {
		log.Fatalf("/dev/loop-control: %v", err)
	}
	defer syscall.Close(cfd)
	a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(cfd), LOOP_CTL_GET_FREE, 0)
	if errno != 0 {
		log.Fatalf("ioctl: %v\n", err)
	}
	log.Printf("a %v b %v err %v\n", a, b, err)
	name = fmt.Sprintf("/dev/loop%d", a)
	return name, nil
}
func linkone(p string, i os.FileInfo, err error) error {
	l.Printf("symtree: p %v\n", p)
	if err != nil {
		return err
	}

	// the tree of symlinks starts at /tcz
	packagel := filepath.SplitList(p)
	// surely there's a better way.
	n := append([]string{"/"}, packagel[2:]...)
	to := path.Join(n...)

	l.Printf("symtree: symlink %v to %v\n", p, to)
	return os.Symlink(p, to)
}
func clonetree(tree string) error {
	lt := len(tree)
	err := filepath.Walk(tree, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
		   l.Printf("walking, dir %v\n", path)
			os.MkdirAll(path[lt:], 0600)
			return nil
		}
		// all else gets a symlink.
		l.Printf("Need to symlnk %v to %v\n", path, path)
		os.Symlink(path, path[lt:])
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return err
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	cmdName := os.Args[1]

	if err := os.MkdirAll(tcz, 0600); err != nil {
		l.Fatal(err)
	}

	packagePath := path.Join("/tcz", cmdName)
	if err := os.MkdirAll(packagePath, 0600); err != nil {
		l.Fatal(err)
	}

	// path.Join doesn't quite work here.
	filepath := path.Join(tcz, cmdName+".tcz")
	if _, err := os.Stat(filepath); err != nil {
		cmd := "http:/" + filepath

		resp, err := http.Get(cmd)
		if err != nil {
			l.Fatalf("Get of %v failed: %v\n", cmd, err)
		}
		defer resp.Body.Close()

		if resp.Status != "200 OK" {
			l.Fatalf("Not OK! %v\n", resp.Status)
		}

		l.Printf("resp %v err %v\n", resp, err)
		// we've to the whole tcz in resp.Body.
		// First, save it to /tcz/name
		f, err := os.Create(filepath)
		if err != nil {
			l.Fatal("Create of :%v: failed: %v\n", filepath, err)
		} else {
			l.Printf("created %v f %v\n", filepath, f)
		}

		if c, err := io.Copy(f, resp.Body); err != nil {
			l.Fatal(err)
		} else {
			/* OK, these are compressed tars ... */
			l.Printf("c %v err %v\n", c, err)
		}
		f.Close()
	}
	/* now walk it, and fetch anything we need. */
	loopname, err := findloop()
	if err != nil {
		l.Fatal(err)
	}
	l.Printf("findloop gets %v err %v\n", loopname, err)

	ffd, err := syscall.Open(filepath, syscall.O_RDONLY, 0)
	lfd, err := syscall.Open(loopname, syscall.O_RDONLY, 0)
	l.Printf("ffd %v lfd %v\n", ffd, lfd)
	a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(lfd), LOOP_SET_FD, uintptr(ffd))
	if errno != 0 {
		l.Fatalf("loop set fd ioctl: %v, %v, %v\n", a, b, errno)
	}
	/* now mount it. The convention is the mount is in /tcz/packagename */
	if err := syscall.Mount(loopname, packagePath, "squashfs", syscall.MS_MGC_VAL|syscall.MS_RDONLY, ""); err != nil {
		l.Fatalf("Mount %s on %s: %v\n", loopname, packagePath, err)
	}
	err = clonetree(packagePath)
	l.Printf("%v\n", err)
}
