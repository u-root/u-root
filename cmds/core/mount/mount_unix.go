// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd

// mount mounts a filesystem at the specified path.
//
// Synopsis:
//
//	mount [-r] [-o options] [-t FSTYPE] DEV PATH
//
// Options:
//
//	-r: read only
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/mount"
	"golang.org/x/sys/unix"
)

var (
	opts           = map[string]uintptr{}
	errUsage       = errors.New("")
	mountflagsbits = map[uint64]string{
		unix.MNT_RDONLY:      "Rdonly",
		unix.MNT_NOEXEC:      "Noexec",
		unix.MNT_NOSUID:      "Nosuid",
		unix.MNT_NOATIME:     "Noatime",
		unix.MNT_SNAPSHOT:    "Snapshot",
		unix.MNT_SUIDDIR:     "Suiddir",
		unix.MNT_SYNCHRONOUS: "Synchronous",
		unix.MNT_ASYNC:       "Async",
		unix.MNT_FORCE:       "Force",
		unix.MNT_NOCLUSTERR:  "Noclusterr",
		unix.MNT_NOCLUSTERW:  "Noclusterw",
		// unix.MNT_NOCOVER: "Nocover",
		// unix.MNT_EMPTYDIR: "Emptydir",
	}
)

type mountOptions []string

func mountflags(flags uint64) (s string) {
	if flags == 0 {
		return ""
	}
	for k, v := range mountflagsbits {
		if (flags & k) == 0 {
			continue
		}
		s += "," + v
		flags &= ^k
	}
	if flags != 0 {
		s += fmt.Sprintf(",%#x", flags)
	}

	return
}

func i8tostring(i []int8) string {
	x := slices.Index(i, 0)
	b := *(*[]byte)(unsafe.Pointer(&i))
	return string(b[:x])
}

func (o *mountOptions) String() string {
	return strings.Join(*o, ",")
}

func (o *mountOptions) Set(value string) error {
	for _, option := range strings.Split(value, ",") {
		*o = append(*o, option)
	}
	return nil
}

type cmd struct {
	stdout  io.Writer
	stderr  io.Writer
	fsType  string
	options mountOptions
	ro      bool
}

func command(stdout, stderr io.Writer, ro bool, fsType string, opts mountOptions) *cmd {
	return &cmd{
		stdout:  stdout,
		stderr:  stderr,
		ro:      ro,
		options: opts,
		fsType:  fsType,
	}
}

func (c *cmd) run(args ...string) error {
	if len(args) == 0 {
		// The freebsd design for getting mounts is to do a
		// getfsstat with a NULL *statfs, which will return the number of mounts;
		// then to allocate that number of struct, and call getfsstat again.
		// This is inherently racy; mounts can come and go.
		// I prefer putting this stuff in a synthetic, a la
		// v8..10, linux, and Plan 9. But nobody asked me :-)
		n, err := syscall.Getfsstat(nil, 1)
		if err != nil {
			return err
		}

		fs := make([]syscall.Statfs_t, n)
		n, err = syscall.Getfsstat(fs, 1)
		if err != nil {
			return err
		}

		fs = fs[:n]
		for _, f := range fs {
			fmt.Fprintf(c.stdout, "%s on %s (%s%s)\n", i8tostring(f.Mntfromname[:]), i8tostring(f.Mntonname[:]), i8tostring(f.Fstypename[:]), mountflags(f.Flags))
		}
		return nil
	}

	if len(args) < 2 {
		return errUsage
	}

	dev := args[0]
	path := args[1]
	var flags uintptr
	var data []string
	for _, option := range c.options {
		switch option {
		default:
			if f, ok := opts[option]; ok {
				flags |= f
			} else {
				data = append(data, option)
			}
		}
	}

	if c.ro {
		flags |= unix.MNT_RDONLY
	}
	if _, err := mount.Mount(dev, path, c.fsType, strings.Join(data, ","), flags); err != nil {
		return err
	}

	return nil
}

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	ro := fs.Bool("r", false, "Read only mount")
	fsType := fs.String("t", "", "File system type")

	var options mountOptions
	fs.Var(&options, "o", "Comma separated list of mount options")
	fs.Parse(os.Args[1:])

	cmd := command(os.Stdout, os.Stderr, *ro, *fsType, options)

	if err := cmd.run(fs.Args()...); err != nil {
		if errors.Is(err, errUsage) {
			fs.Usage()
		}
		log.Fatal(err)
	}
}
