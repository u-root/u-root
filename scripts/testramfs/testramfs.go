// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// testramfs tests things, badly
package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/cpio"
)

const cloneFlags = syscall.CLONE_NEWNS |
	syscall.CLONE_NEWIPC |
	syscall.CLONE_NEWNET |
	syscall.CLONE_NEWPID |
	syscall.CLONE_NEWUTS |
	syscall.CLONE_NEWPID
	//| syscall.CLONE_NEWUSER

var (
	unshared = flag.Bool("unshared", false, "whether this instance has an unshared name space")
	noremove = flag.BoolP("noremove", "n", false, "remove tempdir when done")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalln("usage: %s <cpio-path>", os.Args[0])
	}

	c := flag.Args()[0]

	if !*unshared {
		f, err := os.Open(c)
		if err != nil {
			log.Fatal(err)
		}

		// So, what's the plan here?
		//
		// - new mount namespace
		//   - root mount is a tmpfs mount filled with the archive.
		//
		// - new PID namespace
		//   - archive/init actually runs as PID 1.

		tempDir, err := ioutil.TempDir("", "u-root")
		if err != nil {
			log.Fatal(err)
		}
		// Don't do a RemoveAll. This should be empty and
		// an error can tell us we got something wrong.
		if !*noremove {
			defer func(n string) {
				if err := os.RemoveAll(n); err != nil {
					log.Fatal(err)
				}
			}(tempDir)
		}
		if err := syscall.Mount("", tempDir, "tmpfs", 0, ""); err != nil {
			log.Fatal(err)
		}
		if !*noremove {
			defer func(n string) {
				if err := syscall.Unmount(n, syscall.MNT_DETACH); err != nil {
					log.Fatal(err)
				}
			}(tempDir)
		}

		archiver, err := cpio.Format("newc")
		if err != nil {
			log.Fatal(err)
		}

		r := archiver.Reader(f)
		for {
			rec, err := r.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			cpio.CreateFileInRoot(rec, tempDir)
		}

		c := exec.Command("/proc/self/exe", "--unshared", c)
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
		c.Dir = tempDir
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}

		return
	}

	cmd := exec.Command("/init")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	// It's always best to chroot to '.', over the years it's had special meaning.
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: cloneFlags, Chroot: "."}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
