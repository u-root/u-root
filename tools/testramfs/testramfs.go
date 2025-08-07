// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// testramfs tests things, badly
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/pty"
	"github.com/u-root/u-root/pkg/termios"
)

const (
	unshareFlags = syscall.CLONE_NEWNS
	cloneFlags   = syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWNET |
		// making newpid work will be more tricky,
		// since none of my CLs to fix go runtime for
		// it ever got in.
		// syscall.CLONE_NEWPID |
		syscall.CLONE_NEWUTS
)

func main() {
	var (
		noremove    bool
		interactive bool
	)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.BoolVar(&noremove, "noremove", false, "remove tempdir when done")
	fs.BoolVar(&noremove, "n", false, "remove tempdir when done")
	fs.BoolVar(&interactive, "interactive", false, "interactive mode")
	fs.BoolVar(&interactive, "i", false, "interactive mode")
	fs.Parse(os.Args[1:])

	if fs.NArg() != 1 {
		log.Fatalf("usage: %s <cpio-path>", os.Args[0])
	}

	c := fs.Args()[0]

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

	// Note this is basically a chroot and umask is inherited.
	// The umask has to be zero else some creation will end
	// up with incorrect permissions, a particular problem
	// in device creation.
	u := unix.Umask(0)
	defer unix.Umask(u)

	tempDir, err := os.MkdirTemp("", "u-root")
	if err != nil {
		log.Fatal(err)
	}
	// Don't do a RemoveAll. This should be empty and
	// an error can tell us we got something wrong.
	if !noremove {
		defer func(n string) {
			log.Printf("Removing %v", n)
			if err := os.Remove(n); err != nil {
				log.Fatal(err)
			}
		}(tempDir)
	}
	if err := syscall.Mount("testramfs.tmpfs", tempDir, "tmpfs", 0, ""); err != nil {
		log.Fatal(err)
	}
	if !noremove {
		defer func(n string) {
			log.Printf("Unmounting %v", n)
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
		cpio.CreateFileInRoot(rec, tempDir, false)
	}

	cmd, err := pty.New()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Command("/init")
	cmd.C.SysProcAttr.Chroot = tempDir
	cmd.C.SysProcAttr.Cloneflags = cloneFlags
	cmd.C.SysProcAttr.Unshareflags = unshareFlags
	if interactive {
		t, err := termios.GetTermios(0)
		if err != nil {
			log.Fatal("Getting Termios")
		}
		defer func(t *termios.Termios) {
			if err := termios.SetTermios(0, t); err != nil {
				log.Print(err)
			}
		}(t)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go io.Copy(cmd.TTY, cmd.Ptm)

	// At this point you could use an array of commands/output templates to
	// drive the test, and end with the exit command shown nere.
	for _, c := range []string{"date\n", "exit\n", "exit\n", "exit\n"} {
		if _, err := cmd.Ptm.Write([]byte(c)); err != nil {
			log.Printf("Writing %s: %v", c, err)
		}
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
