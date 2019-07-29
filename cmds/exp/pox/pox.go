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
//     pox [-[-debug]|d] -[-run|r file] [-[-create]|c]  [-[-file]|f tcz-file] file [...file]
//
// Description:
//     pox makes portable executables in squashfs format compatible with
//     tcz format. We don't build in the execution code, rather, we set it
//     up so we can use the command itself. You can create, create and run a command,
//     mount a pox, or mount a pox and run a command in it.
//
// Options:
//     debug|d: verbose
//     file|f file: file name (default /tmp/pox.tcz)
//     run|r: run a file by loopback mounting the squashfs and using the first arg as a command to run in a chroot
//     create|c: create the file.
//     both -c and -r can be used on the same command.
//
// Example:
//	pox -d -r /bin/bash /bin/cat /bin/ls /etc/hosts
//	Will build a squashfs, mount it, and drop you into it running bash.
//	You can use ls and cat on /etc/hosts.
//	Simpler example:
//	pox -d -r /bin/ls /etc/hosts
//	will run ls and exit.
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

const usage = "pox [-[-debug]|d] -[-run|r file] [-[-create]|c]  [-[-file]|f tcz-file] file [...file]"

var (
	debug  = flag.BoolP("debug", "d", false, "enable debug prints")
	run    = flag.BoolP("run", "r", false, "run a test with the first argument")
	create = flag.BoolP("create", "c", true, "create it")
	file   = flag.StringP("output", "f", "/tmp/pox.tcz", "Output file")
	v      = func(string, ...interface{}) {}
)

func pox() error {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	names := flag.Args()
	if len(names) == 0 {
		return fmt.Errorf(usage)
	}

	if *create {
		l, err := ldd.Ldd(names)
		if err != nil {
			return fmt.Errorf("Running ldd on %v: %v", names, err)
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
			out, err := os.OpenFile(dfile, os.O_WRONLY|os.O_CREATE, fi.Mode().Perm())
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
	}

	if !*run {
		return nil
	}
	dir, err := ioutil.TempDir("", "pox")
	if err != nil {
		return err
	}
	if !*debug {
		defer os.RemoveAll(dir)
	}
	m, err := loop.New(*file, "squashfs", "")
	if err != nil {
		return err
	}
	if err := m.Mount(dir, 0); err != nil {
		return err
	}
	c := exec.Command(names[0])
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	c.SysProcAttr = &syscall.SysProcAttr{
		Chroot: dir,
	}

	err = c.Run()
	if err != nil {
		log.Printf("Running test: %v", err)
	}
	if err := m.Unmount(0); err != nil {
		v("Unmounting and freeing %v: %v", m, err)
	}

	v("Done, your pox is in %v", *file)
	return err
}

func main() {
	if err := pox(); err != nil {
		log.Fatal(err)
	}
}
