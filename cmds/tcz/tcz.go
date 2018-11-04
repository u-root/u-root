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
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/loop"
)

const (
	cmd          = "tcz [options] package-names"
	dirMode      = 0755
	tinyCoreRoot = "/TinyCorePackages/tcloop"
)

//http://distro.ibiblio.org/tinycorelinux/5.x/x86_64/tcz/
//The .dep is the name + .dep

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

func clonetree(tree string) error {
	debug("Clone tree %v", tree)
	lt := len(tree)
	err := filepath.Walk(tree, func(path string, fi os.FileInfo, err error) error {

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
				//return err
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
			return fmt.Errorf("Symlink: %v", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("clone tree: %v", err)
	}
	return nil
}

func fetch(p string) error {
	fullpath := filepath.Join(tczLocalPackageDir, p)
	packageName := filepath.Join(tczServerDir, p)

	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		// Either already exists (already been downloaded) or some
		// unresolvable error.
		return err
	}

	cmd := fmt.Sprintf("http://%s:%s/%s", *host, *port, packageName)
	resp, err := http.Get(cmd)
	if err != nil {
		return fmt.Errorf("http.Get(%s) failed: %v", cmd, err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return syscall.ENOENT
	}

	// we have the whole tcz in resp.Body.
	// First, save it to /tczRoot/name
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("error reading download body of %q: %v", cmd, err)
	}
	return nil
}

// deps is ALL the packages we need fetched or not
// this may even let us work with parallel tcz, ALMOST
func installPackage(tczName string, deps map[string]bool) error {
	debug("installPackage: %v %v\n", tczName, deps)
	depName := tczName + ".dep"
	if err := fetch(tczName); err != nil {
		return err
	}
	deps[tczName] = true

	debug("Fetched %v\n", tczName)

	// now fetch dependencies if any.
	if err := fetch(depName); err == nil {
		debug("Fetched dep ok!\n")
	} else {
		debug("No dep file found\n")
		if err := ioutil.WriteFile(filepath.Join(tczLocalPackageDir, depName), []byte{}, os.FileMode(0444)); err != nil {
			debug("Tried to write Blank file %v, failed %v\n", depName, err)
		}
		return nil
	}

	// read deps file
	depFullPath := filepath.Join(tczLocalPackageDir, depName)
	deplist, err := ioutil.ReadFile(depFullPath)
	if err != nil {
		return fmt.Errorf("read(%q) = %v", depName, err)
	}

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
	if err := ioutil.WriteFile(depFullPath, []byte(realDepList), os.FileMode(0444)); err != nil {
		debug("Tried to write deplist file %v, failed %v\n", depName, err)
		return err
	}
	return nil
}

func setupPackages(tczName string, deps map[string]bool) error {
	for v := range deps {
		cmdName := strings.Split(v, filepath.Ext(v))[0]
		packagePath := filepath.Join(tinyCoreRoot, cmdName)

		if _, err := os.Stat(packagePath); err == nil {
			debug("PackagePath %s exists, skipping mount", packagePath)
			continue
		}

		if err := os.MkdirAll(packagePath, dirMode); err != nil {
			return fmt.Errorf("package directory %s at %s, can not be created: %v", tczName, packagePath, err)
		}

		loopname, err := loop.FindDevice()
		if err != nil {
			return err
		}
		pkgpath := filepath.Join(tczLocalPackageDir, v)
		if err := loop.SetFile(loopname, pkgpath); err != nil {
			return err
		}

		/* now mount it. The convention is the mount is in /tinyCoreRoot/packagename */
		if err := syscall.Mount(loopname, packagePath, "squashfs", syscall.MS_RDONLY, ""); err != nil {
			return err
		}
		if err := clonetree(packagePath); err != nil {
			return err
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
