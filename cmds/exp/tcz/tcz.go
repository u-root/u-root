// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

const (
	cmd = "tcz [options] package-names"
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
	dirMode           = 0o755
	tinyCoreRoot      = "/TinyCorePackages/tcloop"
)

// http://distro.ibiblio.org/tinycorelinux/5.x/x86_64/tcz/
// The .dep is the name + .dep

var (
	l                  = log.New(os.Stdout, "tcz: ", 0)
	host               = flag.String("h", "tinycorelinux.net", "Host name for packages")
	version            = flag.String("v", "8.x", "tinycore version")
	arch               = flag.String("a", "x86_64", "tinycore architecture")
	port               = flag.String("p", "80", "Host port")
	install            = flag.Bool("i", true, "Install the packages, i.e. mount and create symlinks")
	tczRoot            = flag.String("r", "/tcz", "tcz root directory")
	debugPrint         = flag.Bool("d", false, "Enable debug prints")
	skip               = flag.String("skip", "", "Packages to skip")
	debug              = func(f string, s ...interface{}) {}
	tczServerDir       string
	tczLocalPackageDir string
	ignorePackage      = make(map[string]struct{})
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
	debug("a %v b %v err %v\n", a, b, err)
	name = fmt.Sprintf("/dev/loop%d", a)
	return name, nil
}

func clonetree(tree string) error {
	debug("Clone tree %v", tree)
	lt := len(tree)
	err := filepath.Walk(tree, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		debug("Clone tree with path %s fi %v", path, fi)
		if fi.IsDir() {
			debug("Walking, dir %v\n", path)
			if path[lt:] == "" {
				return nil
			}
			if err := os.MkdirAll(path[lt:], dirMode); err != nil {
				debug("Mkdir of %s failed: %v", path[lt:], err)
				// TODO: EEXIST should not be an error. Ignore
				// err for now. FIXME.
				// return err
			}
			return nil
		}
		// all else gets a symlink.

		// If the link exists
		if target, err := os.Readlink(path[lt:]); err == nil {
			// Confirm that it points to the same path to be symlinked
			if target == path {
				return nil
			}

			// If it does not, return error because tcz packages are inconsistent
			return fmt.Errorf("symlink: need %q -> %q, but %q -> %q is already there", path, path[lt:], path, target)
		}

		debug("Need to symlink %v to %v\n", path, path[lt:])

		if err := os.Symlink(path, path[lt:]); err != nil {
			return fmt.Errorf("symlink: %w", err)
		}

		return nil
	})
	if err != nil {
		l.Fatalf("Clone tree: %v", err)
	}
	return nil
}

func fetch(p string) error {
	fullpath := filepath.Join(tczLocalPackageDir, p)
	packageName := filepath.Join(tczServerDir, p)

	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		debug("package %s is downloaded\n", fullpath)
		return nil
	}

	if _, err := os.Stat(fullpath); err != nil {
		cmd := fmt.Sprintf("http://%s:%s/%s", *host, *port, packageName)
		debug("Fetch %v\n", cmd)

		resp, err := http.Get(cmd)
		if err != nil {
			l.Fatalf("Get of %v failed: %v\n", cmd, err)
		}
		defer resp.Body.Close()

		if resp.Status != "200 OK" {
			debug("%v Not OK! %v\n", cmd, resp.Status)
			return syscall.ENOENT
		}

		debug("resp %v err %v\n", resp, err)
		// we have the whole tcz in resp.Body.
		// First, save it to /tczRoot/name
		f, err := os.Create(fullpath)
		if err != nil {
			l.Fatalf("Create of :%v: failed: %v\n", fullpath, err)
		} else {
			debug("created %v f %v\n", fullpath, f)
		}

		if c, err := io.Copy(f, resp.Body); err != nil {
			l.Fatal(err)
		} else {
			/* OK, these are compressed tars ... */
			debug("c %v err %v\n", c, err)
		}
		f.Close()
	}
	return nil
}

// deps is ALL the packages we need fetched or not
// this may even let us work with parallel tcz, ALMOST
func installPackage(tczName string, deps map[string]bool) error {
	debug("installPackage: %v %v\n", tczName, deps)
	depName := tczName + ".dep"
	if err := fetch(tczName); err != nil {
		l.Fatal(err)
	}
	deps[tczName] = true
	debug("Fetched %v\n", tczName)
	// now fetch dependencies if any.
	if err := fetch(depName); err == nil {
		debug("Fetched dep ok!\n")
	} else {
		debug("No dep file found\n")
		if err := os.WriteFile(filepath.Join(tczLocalPackageDir, depName), []byte{}, os.FileMode(0o444)); err != nil {
			debug("Tried to write Blank file %v, failed %v\n", depName, err)
		}
		return nil
	}
	// read deps file
	depFullPath := filepath.Join(tczLocalPackageDir, depName)
	deplist, err := os.ReadFile(depFullPath)
	if err != nil {
		l.Fatalf("Fetched dep file %v but can't read it? %v", depName, err)
	}
	debug("deplist for %v is :%v:\n", depName, deplist)
	realDepList := ""
	for _, v := range strings.Split(string(deplist), "\n") {
		// split("name\n") gets you a 2-element array with second
		// element the empty string
		if len(v) == 0 {
			break
		}
		if _, ok := ignorePackage[v]; ok {
			debug("%v is ignored", v)
			continue
		}
		realDepList = realDepList + v + "\n"
		debug("FOR %v get package %v\n", tczName, v)
		if deps[v] {
			continue
		}
		if err := installPackage(v, deps); err != nil {
			return err
		}
	}
	if string(deplist) == realDepList {
		return nil
	}
	if err := os.WriteFile(depFullPath, []byte(realDepList), os.FileMode(0o444)); err != nil {
		debug("Tried to write deplist file %v, failed %v\n", depName, err)
		return err
	}
	return nil
}

func setupPackages(tczName string, deps map[string]bool) error {
	debug("setupPackages: @ %v deps %v\n", tczName, deps)
	for v := range deps {
		cmdName := strings.Split(v, filepath.Ext(v))[0]
		packagePath := filepath.Join(tinyCoreRoot, cmdName)

		if _, err := os.Stat(packagePath); err == nil {
			debug("PackagePath %s exists, skipping mount", packagePath)
			continue
		}

		if err := os.MkdirAll(packagePath, dirMode); err != nil {
			l.Fatalf("Package directory %s at %s, can not be created: %v", tczName, packagePath, err)
		}

		loopname, err := findloop()
		if err != nil {
			l.Fatal(err)
		}
		debug("findloop gets %v err %v\n", loopname, err)
		pkgpath := filepath.Join(tczLocalPackageDir, v)
		ffd, err := syscall.Open(pkgpath, syscall.O_RDONLY, 0)
		if err != nil {
			l.Fatalf("%v: %v\n", pkgpath, err)
		}
		lfd, err := syscall.Open(loopname, syscall.O_RDONLY, 0)
		if err != nil {
			l.Fatalf("%v: %v\n", loopname, err)
		}
		debug("ffd %v lfd %v\n", ffd, lfd)

		a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(lfd), LOOP_SET_FD, uintptr(ffd))
		if errno != 0 {
			l.Fatalf("loop set fd ioctl: pkgpath :%v:, loop :%v:, %v, %v, %v\n", pkgpath, loopname, a, b, errno)
		}

		/* now mount it. The convention is the mount is in /tinyCoreRoot/packagename */
		if err := syscall.Mount(loopname, packagePath, "squashfs", syscall.MS_MGC_VAL|syscall.MS_RDONLY, ""); err != nil {
			// how I hate Linux.
			// Since all you ever get back is a USELESS ERRNO
			// note: the open succeeded for both things, the mkdir worked,
			// the loop device open worked. And we can get ENODEV. So,
			// just what went wrong? "We're not telling. Here's an errno."
			li, lerr := os.Stat(loopname)
			pi, perr := os.Stat(pkgpath)
			di, derr := os.Stat(packagePath)
			e := fmt.Sprintf("%v: loop is (%v, %v); package is (%v, %v); dir is (%v, %v). Is squashfs built into your kernel?", err, li, lerr, pi, perr, di, derr)
			l.Fatalf("Mount :%s: on :%s: %v\n", loopname, packagePath, e)
		}
		err = clonetree(packagePath)
		if err != nil {
			l.Fatalf("clonetree:  %v\n", err)
		}
	}
	return nil
}

func usage() string {
	return "tcz [-v version] [-a architecture] [-h hostname] [-p host port] [-d debug prints] PROGRAM..."
}

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()
	if *debugPrint {
		debug = l.Printf
	}

	ip := strings.Fields(*skip)
	debug("ignored packages: %v", ip)
	for _, p := range ip {
		ignorePackage[p+".tcz"] = struct{}{}
	}
	needPackages := make(map[string]bool)
	tczServerDir = filepath.Join("/", *version, *arch, "/tcz")
	tczLocalPackageDir = filepath.Join(*tczRoot, tczServerDir)

	packages := flag.Args()
	debug("tcz: packages %v", packages)

	if len(packages) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(tczLocalPackageDir, dirMode); err != nil {
		l.Fatal(err)
	}

	if *install {
		if err := os.MkdirAll(tinyCoreRoot, dirMode); err != nil {
			l.Fatal(err)
		}
	}

	for _, cmdName := range packages {

		tczName := cmdName + ".tcz"

		if err := installPackage(tczName, needPackages); err != nil {
			l.Fatal(err)
		}

		debug("After installpackages: needPackages %v\n", needPackages)

		if *install {
			if err := setupPackages(tczName, needPackages); err != nil {
				l.Fatal(err)
			}
		}
	}
}
