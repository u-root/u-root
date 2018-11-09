// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pox builds a portable executable as a squashfs image.
// It is intended to create files compatible with tinycore
// tcz files.
// This could have been a simple program but mksquashfs does not
// preserve path information.
// Yeah.
//
// Synopsis:
//     pox [-d] -[output|o file]
//
// Description:
//     pox makes portable executables in squashfs format compatible with
//     tcz format. We don't build in the execution code, rather, we set it
//     up so we can use the command itself.
//
// Options:
//     debug|d: verbose
//     output|o file: output file name (default /tmp/pox.tcz)
//     test|t: run a test by loopback mounting the squashfs and using the first arg as a command to run in a chroot
//
// Example:
//	pox -d -t /bin/bash /bin/cat /bin/ls /etc/hosts
//	Will build and squashfs, mount it, and drop you into it running bash.
//	You can use ls and cat on /etc/hosts.
//	Simpler example:
//	pox -d -t /bin/ls /etc/hosts
//	will run ls and exit.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/loop"
)

const usage = "pox [-d] [-f file] command..."

var (
	debug  = flag.BoolP("debug", "d", false, "enable debug prints")
	test   = flag.BoolP("test", "t", false, "run a test with the first argument")
	create = flag.BoolP("create", "c", true, "create it")
	v      = func(string, ...interface{}) {}
	ofile  = flag.StringP("output", "o", "/tmp/pox.tcz", "Output file")
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
			return err
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
		archiver := cpio.InMemArchive()
		for _, f := range names {
			v("Process %v", f)
			rec, err := cpio.GetRecord(f)
			if err != nil {
				return err
			}
			if err := archiver.WriteRecord(rec); err != nil {
				return err
			}
		}
		v("%v", archiver)
		rr := archiver.Reader()
		for {
			r, err := rr.ReadRecord()
			v("%v %v", r, err)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := cpio.CreateFileInRoot(r, dir, true); err != nil {
				return err
			}

		}
		c := exec.Command("mksquashfs", dir, *ofile, "-noappend")
		o, err := c.CombinedOutput()
		v("%v", string(o))
		if err != nil {
			return fmt.Errorf("%v: %v: %v", c.Args, string(o), err)
		}
	}

	if !*test {
		return nil
	}
	dir, err := ioutil.TempDir("", "pox")
	if err != nil {
		return err
	}
	if !*debug {
		defer os.RemoveAll(dir)
	}
	m, err := loop.New(*ofile, dir, "squashfs", 0, "")
	if err != nil {
		return err
	}
	if err := m.Mount(); err != nil {
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

	v("Done, your pox is in %v", *ofile)
	return err
}

func main() {
	if err := pox(); err != nil {
		log.Fatal(err)
	}
}
