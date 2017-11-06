// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/log"
)

const (
	cmd = "tcz [options] package-names"

	// IOCTL commands --- we will commandeer 0x4C ('L')
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
)

var (
	host               = flag.String("h", "tinycorelinux.net", "Host name for packages")
	version            = flag.String("v", "8.x", "tinycore version")
	arch               = flag.String("a", "x86_64", "tinycore architecture")
	port               = flag.String("p", "80", "Host port")
	install            = flag.Bool("i", true, "Install the packages, i.e. mount and create symlinks")
	tczRoot            = flag.String("r", "/tcz", "tcz root directory")
	skip               = flag.String("skip", "", "Packages to skip")
	tczServerDir       string
	tczLocalPackageDir string
	ignorePackage      = make(map[string]struct{})
)

// consider making this a goroutine which pushes the string down the channel.
func findloop() (name string, err error) {
	cfd, err := os.OpenFile("/dev/loop-control", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("/dev/loop-control: %v", err)
	}
	defer cfd.Close()

	num, _, errno := syscall.Syscall(syscall.SYS_IOCTL, cfd.Fd(), LOOP_CTL_GET_FREE, 0)
	if errno != 0 {
		return "", fmt.Errorf("ioctl cannot get free loop: %v", errno)
	}
	return fmt.Sprintf("/dev/loop%d", num), nil
}

func clonetree(tree string) error {
	log.Printf("Clone tree %v", tree)
	lt := len(tree)
	err := filepath.Walk(tree, func(path string, fi os.FileInfo, err error) error {
		log.Printf("Clone tree with path %q fi %v", path, fi)

		switch {
		case fi.IsDir():
			log.Printf("Walking; dir %q", path)
			if path[lt:] == "" {
				return nil
			}

			if err := os.MkdirAll(path[lt:], 0700); err != nil && !os.IsExist(err) {
				return fmt.Errorf("Mkdir(%q) failed: %v", path[lt:], err)
			}
			return nil

		default:
			// All other file types get a symlink.
			if target, err := os.Readlink(path[lt:]); err == nil {
				// Confirm that it points to the same path to be symlinked
				if target == path {
					return nil
				}

				// If it does not, return error because tcz packages are inconsistent
				return fmt.Errorf("symlink: need %q -> %q, but %q -> %q is already there", path, path[lt:], path, target)
			}

			log.Printf("Need to symlink %q to %q", path, path[lt:])

			return os.Symlink(path, path[lt:])
		}
	})

	if err != nil {
		log.Fatalf("Clone tree: %v", err)
	}
	return nil
}

func fetch(p string) error {
	fullpath := filepath.Join(tczLocalPackageDir, p)
	packageName := filepath.Join(tczServerDir, p)

	if _, err := os.Stat(fullpath); err == nil {
		log.Printf("package %s is already downloaded", fullpath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat(%q) returned unexpected error: %v", fullpath, err)
	}

	// Package path does not exist yet. Let's download the package.
	cmd := fmt.Sprintf("http://%s:%s/%s", *host, *port, packageName)
	log.Printf("Fetch %v", cmd)

	resp, err := http.Get(cmd)
	if err != nil {
		return fmt.Errorf("HTTP GET of %q failed: %v", cmd, err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return fmt.Errorf("HTTP GET of %q: Not OK! %v", cmd, resp.Status)
	}
	log.Printf("resp %v err %v", resp, err)

	// we have the whole tcz in resp.Body.
	// First, save it to /tczRoot/name
	f, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("Create of %q failed: %v", fullpath, err)
	}
	defer f.Close()
	log.Printf("created %q: %v", fullpath, f)

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("copy: %v", err)
	}
	return nil
}

// deps is ALL the packages we need fetched or not
// this may even let us work with parallel tcz, ALMOST
func installPackage(tczName string, deps map[string]bool) error {
	log.Printf("installPackage: %v %v", tczName, deps)

	if err := fetch(tczName); err != nil {
		return err
	}
	deps[tczName] = true
	log.Printf("Fetched %v", tczName)

	// http://distro.ibiblio.org/tinycorelinux/5.x/x86_64/tcz/
	// The dependency file name is the name + .dep
	depName := tczName + ".dep"

	// Fetch dependencies if any.
	if err := fetch(depName); err == nil {
		log.Printf("Fetched dep ok!")
	} else {
		log.Printf("No dep file found")
		if err := ioutil.WriteFile(filepath.Join(tczLocalPackageDir, depName), []byte{}, os.FileMode(0444)); err != nil {
			return fmt.Errorf("tried to write blank file %v: %v", depName, err)
		}
		return nil
	}

	// Read deps file.
	depFullPath := filepath.Join(tczLocalPackageDir, depName)
	deplist, err := ioutil.ReadFile(depFullPath)
	if err != nil {
		return fmt.Errorf("fetched dep file %v but can't read it? %v", depName, err)
	}

	depScanner := bytes.NewBuffer(deplist)
	var dependencies []string
	for {
		dependency, err := depScanner.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if _, ok := ignorePackage[dependency]; ok {
			log.Printf("%v is ignored", dependency)
			continue
		}

		dependencies = append(dependencies, dependency)

		log.Printf("For %v get package %v\n", tczName, dependency)
		if deps[dependency] {
			continue
		}
		if err := installPackage(dependency, deps); err != nil {
			return err
		}
	}

	realDepList := strings.Join(dependencies, "\n") + "\n"
	if string(deplist) == realDepList {
		return nil
	}
	if err := ioutil.WriteFile(depFullPath, []byte(realDepList), os.FileMode(0444)); err != nil {
		return fmt.Errorf("tried to write deplist file %v: %v", depName, err)
	}
	return nil

}

func setupPackages(tczName string, deps map[string]bool) error {
	log.Printf("setupPackages: @ %v deps %v", tczName, deps)

	for v := range deps {
		cmdName := strings.Split(v, filepath.Ext(v))[0]
		packagePath := filepath.Join("/tmp/tcloop", cmdName)

		if _, err := os.Stat(packagePath); err == nil {
			fmt.Printf("PackagePath %s exists, skipping mount", packagePath)
			continue
		}

		if err := os.MkdirAll(packagePath, 0700); err != nil {
			return fmt.Errorf("package directory %s at %s, can not be created: %v", tczName, packagePath, err)
		}

		loopname, err := findloop()
		if err != nil {
			return err
		}

		pkgpath := filepath.Join(tczLocalPackageDir, v)
		ffd, err := syscall.Open(pkgpath, syscall.O_RDONLY, 0)
		if err != nil {
			return err
		}
		lfd, err := syscall.Open(loopname, syscall.O_RDONLY, 0)
		if err != nil {
			return err
		}

		log.Printf("ffd %v lfd %v\n", ffd, lfd)

		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(lfd), LOOP_SET_FD, uintptr(ffd)); errno != 0 {
			return errno
		}

		/* now mount it. The convention is the mount is in /tmp/tcloop/packagename */
		if err := syscall.Mount(loopname, packagePath, "squashfs", syscall.MS_RDONLY, ""); err != nil {
			return fmt.Errorf("mount(%q on %q): %v", loopname, packagePath, err)
		}

		if err := clonetree(packagePath); err != nil {
			return err
		}
	}
	return nil

}

func usage() string {
	return "tcz [-v version] [-a architecture] [-h hostname] [-p host port] [-d log.Printf prints] PROGRAM..."
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

	ip := strings.Fields(*skip)
	log.Printf("ignored packages: %v", ip)
	for _, p := range ip {
		ignorePackage[p+".tcz"] = struct{}{}
	}
	needPackages := make(map[string]bool)
	tczServerDir = filepath.Join("/", *version, *arch, "/tcz")
	tczLocalPackageDir = filepath.Join(*tczRoot, tczServerDir)

	packages := flag.Args()

	if len(packages) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(tczLocalPackageDir, 0700); err != nil {
		log.Fatalf("%v", err)
	}

	if err := os.MkdirAll("/tmp/tcloop", 0700); err != nil {
		log.Fatalf("%v", err)
	}

	for _, cmdName := range packages {
		tczName := cmdName + ".tcz"
		if err := installPackage(tczName, needPackages); err != nil {
			log.Fatalf("%v", err)
		}

		log.Printf("After installpackages: needPackages %v\n", needPackages)
		if *install {
			if err := setupPackages(tczName, needPackages); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
