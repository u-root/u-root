// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pox builds a portable executable as a squashfs image.
// It is intended to create files compatible with tinycore
// tcz files. One of more of the files can be programs
// but that is not required.
// This could have been a simple program but mksquashfs does not
// preserve path information.
// Yeah.
//
// Synopsis:
//     pox [-[-debug]|d] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]
//
// Description:
//     pox makes portable executables in squashfs format compatible with
//     tcz format. We don't build in the execution code, rather, we set it
//     up so we can use the command itself. You can either create the TCZ image
//     or run a command within an image that was previously created.
//
// Options:
//     debug|d: verbose
//     file|f file: file name (default /tmp/pox.tcz)
//     run|r: run a file by loopback mounting the squashfs and using the first arg to pox as the command to run.
//     create|c: create the TCZ file.
//     Exactly one of -c and -r must be used on the same command.
//
// Example:
//	$ pox -c /bin/bash /bin/cat /bin/ls /etc/hosts
//	Will build a squashfs, which will be /tmp/pox.tcz
//
//	$ sudo pox -r /bin/bash
//	Will drop you into the /tmp/pox.tcz running bash
//	You can use ls and cat on /etc/hosts.
//
//	Simpler example:
//	$ sudo pox -r /bin/ls
//	will run ls and exit.
//
// Notes:
// - When running a pox, you likely need sudo to access /dev/loop*.
//
// - Your binaries and programs show up in the TCZ using whatever path you
// provided to pox.  For instance, if you are in /home/you/somedir/ and have
// ./bin/foo, and you pox -c bin/foo, your TCZ will contain bin/foo.  If you
// pox /home/you/somedir/bin/foo, your TCZ will contain the full path.
//
// - When creating a pox with an executable with shared libraries that are not
// installed on your system, such as for a project installed in your home
// directory, run pox from the installation prefix directory, such that the
// libraries and binaries are below pox's working directory.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/loop"
)

const usage = "pox [-[-debug]|d] -[-run|r] | -[-create]|c  [-[-file]|f tcz-file] file [...file]"

var (
	debug  = flag.BoolP("debug", "d", false, "enable debug prints")
	run    = flag.BoolP("run", "r", false, "Run the first file argument")
	create = flag.BoolP("create", "c", false, "create it")
	file   = flag.StringP("output", "f", "/tmp/pox.tcz", "Output file")
	v      = func(string, ...interface{}) {}
)

func poxCreate(names []string) error {
	if len(names) == 0 {
		return fmt.Errorf(usage)
	}
	l, err := ldd.Ldd(names)
	if err != nil {
		var stderr []byte
		if eerr, ok := err.(*exec.ExitError); ok {
			stderr = eerr.Stderr
		}
		return fmt.Errorf("Running ldd on %v: %v %s", names,
			err, stderr)
	}

	for _, dep := range l {
		v("%s", dep.FullName)
		names = append(names, dep.FullName)
	}
	// Now we need to make a template file hierarchy and put
	// the stuff we want in there.
	dir, err := ioutil.TempDir("", "pox")
	if err != nil {
		return err
	}
	if !*debug {
		defer os.RemoveAll(dir)
	}
	// We don't use defer() here to close files as
	// that can cause open failures with a large enough number.
	for _, f := range names {
		v("Process %v", f)
		fi, err := os.Stat(f)
		if err != nil {
			return err
		}
		in, err := os.Open(f)
		if err != nil {
			return err
		}
		dfile := filepath.Join(dir, f)
		d := filepath.Dir(dfile)
		if err := os.MkdirAll(d, 0755); err != nil {
			in.Close()
			return err
		}
		out, err := os.OpenFile(dfile, os.O_WRONLY|os.O_CREATE,
			fi.Mode().Perm())
		if err != nil {
			in.Close()
			return err
		}
		_, err = io.Copy(out, in)
		in.Close()
		out.Close()
		if err != nil {
			return err
		}

	}
	c := exec.Command("mksquashfs", dir, *file, "-noappend")
	o, err := c.CombinedOutput()
	v("%v", string(o))
	if err != nil {
		return fmt.Errorf("%v: %v: %v", c.Args, string(o), err)
	}
	v("Done, your pox is in %v", *file)

	return nil
}

func poxRun(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(usage)
	}
	dir, err := ioutil.TempDir("", "pox")
	if err != nil {
		return err
	}
	if !*debug {
		defer os.RemoveAll(dir)
	}

	lo, err := loop.New(*file, "squashfs", "")
	if err != nil {
		return err
	}
	defer lo.Free() //nolint:errcheck

	mountPoint, err := lo.Mount(dir, 0)
	if err != nil {
		return err
	}
	defer mountPoint.Unmount(0) //nolint:errcheck

	c := exec.Command(args[0])
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{
		Chroot: dir,
	}

	if err = c.Run(); err != nil {
		log.Printf("Running test: %v", err)
	}

	return err
}

func pox() error {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	if (*create && *run) || (!*create && !*run) {
		return fmt.Errorf(usage)
	}
	if *create {
		return poxCreate(flag.Args())
	}
	if *run {
		return poxRun(flag.Args())
	}
	return fmt.Errorf(usage)
}

func main() {
	if err := pox(); err != nil {
		log.Fatal(err)
	}
}
