// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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

//http://distro.ibiblio.org/tinycorelinux/5.x/x86_64/tcz/
//The .dep is the name + .dep

var (
	l       = log.New(os.Stdout, "tcz: ", 0)
	host    = flag.String("h", "tinycorelinux.net", "Host name for packages")
	version = flag.String("v", "5.x", "tinycore version")
	arch    = flag.String("a", "x86_64", "tinycore architecture")
	port    = flag.String("p", "80", "Host port")
	tcz     string
)

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

	// the tree of symlinks starts at /tmp/tcloop
	packagel := filepath.SplitList(p)
	// surely there's a better way.
	n := append([]string{"/"}, packagel[2:]...)
	to := path.Join(n...)

	l.Printf("symtree: remove %v\n", to)
	os.Remove(to)
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

func fetch(p string) error {
	// path.Join doesn't quite work here.
	fullpath := path.Join(tcz, p)
	if _, err := os.Stat(fullpath); err != nil {
		cmd := "http://" + *host + ":" + *port + "/" + fullpath
		l.Printf("Fetch %v\n", cmd)

		resp, err := http.Get(cmd)
		if err != nil {
			l.Fatalf("Get of %v failed: %v\n", cmd, err)
		}
		defer resp.Body.Close()

		if resp.Status != "200 OK" {
			l.Printf("%v Not OK! %v\n", cmd, resp.Status)
			return syscall.ENOENT
		}

		l.Printf("resp %v err %v\n", resp, err)
		// we have the whole tcz in resp.Body.
		// First, save it to /tcz/name
		f, err := os.Create(fullpath)
		if err != nil {
			l.Fatal("Create of :%v: failed: %v\n", fullpath, err)
		} else {
			l.Printf("created %v f %v\n", fullpath, f)
		}

		if c, err := io.Copy(f, resp.Body); err != nil {
			l.Fatal(err)
		} else {
			/* OK, these are compressed tars ... */
			l.Printf("c %v err %v\n", c, err)
		}
		f.Close()
	}
	return nil
}

// deps is ALL the packages we need fetched or not
// this may even let us work with parallel tcz, ALMOST
func installPackage(tczName string, deps map[string]bool) error {
	l.Printf("installPackage: %v %v\n", tczName, deps)
	depName := tczName + ".dep"
	// path.Join doesn't quite work here.
	if err := fetch(tczName); err != nil {
		l.Fatal(err)
	}
	deps[tczName] = true
	l.Printf("Fetched %v\n", tczName)
	// now fetch dependencies if any.
	if err := fetch(depName); err == nil {
		l.Printf("Fetched dep ok!\n")
	} else {
		l.Printf("No dep file found\n")
		if err := ioutil.WriteFile(path.Join(tcz, depName), []byte{}, os.FileMode(0444)); err != nil {
			l.Printf("Tried to write Blank file %v, failed %v\n", depName, err)
		}
		return nil
	}
	// read deps file
	deplist, err := ioutil.ReadFile(path.Join(tcz, depName))
	if err != nil {
		l.Fatalf("Fetched dep file %v but can't read it? %v", depName, err)
	}
	l.Printf("deplist for %v is :%v:\n", depName, deplist)
	for _, v := range strings.Split(string(deplist), "\n") {
		// split("name\n") gets you a 2-element array with second
		// element the empty string
		if len(v) == 0 {
			break
		}
		l.Printf("FOR %v get package %v\n", tczName, v)
		if deps[v] {
			continue
		}
		if err := installPackage(v, deps); err != nil {
			return err
		}
	}
	return nil

}

func setupPackages(tczName string, deps map[string]bool) error {
	l.Printf("setupPackages: @ %v deps %v\n", tczName, deps)
	for v := range deps {
		cmdName := strings.Split(v, ".")[0]
		packagePath := path.Join("/tmp/tcloop", cmdName)
		if err := os.MkdirAll(packagePath, 0600); err != nil {
			l.Fatal(err)
		}

		loopname, err := findloop()
		if err != nil {
			l.Fatal(err)
		}
		l.Printf("findloop gets %v err %v\n", loopname, err)
		pkgpath := path.Join(tcz, v)
		ffd, err := syscall.Open(pkgpath, syscall.O_RDONLY, 0)
		if err != nil {
			l.Fatal("%v: %v\n", pkgpath, err)
		}
		lfd, err := syscall.Open(loopname, syscall.O_RDONLY, 0)
		if err != nil {
			l.Fatal("%v: %v\n", loopname, err)
		}
		l.Printf("ffd %v lfd %v\n", ffd, lfd)
		a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(lfd), LOOP_SET_FD, uintptr(ffd))
		if errno != 0 {
			l.Fatalf("loop set fd ioctl: pkgpath :%v:, loop :%v:, %v, %v, %v\n", pkgpath, loopname, a, b, errno)
		}
		/* now mount it. The convention is the mount is in /tmp/tcloop/packagename */
		if err := syscall.Mount(loopname, packagePath, "squashfs", syscall.MS_MGC_VAL|syscall.MS_RDONLY, ""); err != nil {
			l.Fatalf("Mount :%s: on :%s: %v\n", loopname, packagePath, err)
		}
		err = clonetree(packagePath)
		if err != nil {
			l.Fatalf("clonetree:  %v\n", err)
		}
	}
	return nil

}

func main() {
	flag.Parse()
	needPackages := make(map[string]bool)
	tcz = path.Join("/", *version, *arch, "tcz")
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	cmdName := flag.Args()[0]
	tczName := cmdName + ".tcz"

	if err := os.MkdirAll(tcz, 0600); err != nil {
		l.Fatal(err)
	}

	if err := os.MkdirAll("/tmp/tcloop", 0600); err != nil {
		l.Fatal(err)
	}

	// path.Join doesn't quite work here.
	if err := installPackage(tczName, needPackages); err != nil {
		l.Fatal(err)
	}
	l.Printf("After installpackages: needPackages %v\n", needPackages)

	if err := setupPackages(tczName, needPackages); err != nil {
		l.Fatal(err)
	}
}
