// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

// fusermount is a very limited replacement for the C fusermount.  It
// is invoked by other programs, or interactively only to unmount.
//
// Synopsis:
//
//	fusermount [-u|--unmount] [-z|--lazy] [-v|--verbose] <mountpoint>
//
// For mounting, per the FUSE model, the environment variable
// _FUSE_COMMFD must have the value of a file descriptor variable on
// which we pass the fuse fd.
//
// There is some checking we don't do, e.g. for the number of active
// mount points.  Last time I checked, that's the kind of stuff
// kernels do.
//
// Description:
//
//	invoke fuse mount operations
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	// CommFD is the environment variable which contains the comms fd.
	CommFD  = "_FUSE_COMMFD"
	fuseDev = "/dev/fuse"
)

var (
	debug = func(string, ...interface{}) {}
	mpt   string
)

const usage = "usage: fusermount [-u|--unmount] [-z|--lazy] [-v|--verbose] <mountpoint>"

func umount(n string, lazy bool) error {
	// we're not doing all the folderol of standard
	// fusermount for euid() == 0.
	// Let's see how that works out.
	flags := 0
	if lazy {
		flags |= unix.MNT_DETACH
	}

	// TODO: anything we need here if unit.Getuid() == 0.
	// So far there is nothing.
	err := unix.Unmount(n, flags)
	return err
}

func openFUSE() (int, error) {
	return unix.Open("/dev/fuse", unix.O_RDWR, 0)
}

// MountPointOK performs validation on the mountpoint.
// Bury all your magic in here.
func MountPointOK(mpt string) error {
	// We wait until we can drop privs to test the mpt
	// parameter, since ability to walk the path can
	// differ for root and the real user id.
	if err := dropPrivs(); err != nil {
		return err
	}
	defer restorePrivs()
	mpt = filepath.Clean(mpt)
	r, err := filepath.EvalSymlinks(mpt)
	if err != nil {
		return err
	}
	if r != mpt {
		return fmt.Errorf("resolved path %q and mountpoint %q are not the same", r, mpt)
	}
	// I'm not sure why fusermount wants to open the mountpoint, so let's mot for now.
	// And, for now, directories only? We don't see a current need to mount
	// FUSE on any other type of file.
	if err := os.Chdir(mpt); err != nil {
		return err
	}

	return nil
}

func getCommFD() (int, error) {
	commfd, ok := os.LookupEnv(CommFD)
	if !ok {
		return -1, fmt.Errorf(CommFD + "was not set and this program can't be used interactively")
	}
	debug("CommFD %v", commfd)

	cfd, err := strconv.Atoi(commfd)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", CommFD, err)
	}
	debug("CFD is %v", cfd)
	var st unix.Stat_t
	if err := unix.Fstat(cfd, &st); err != nil {
		return -1, fmt.Errorf("_FUSE_COMMFD: %d: %w", cfd, err)
	}
	debug("cfd stat is %v", st)

	return cfd, nil
}

func doMount(fd int) error {
	flags := uintptr(unix.MS_NODEV | unix.MS_NOSUID)
	// From the kernel:
	// if (!d->fd_present || !d->rootmode_present ||
	//	!d->user_id_present || !d->group_id_present)
	//		return 0;
	// Yeah. You get EINVAL if any one of these is not set.
	// Docs? what? Docs?
	return unix.Mount("nodev", ".", "fuse", flags, fmt.Sprintf("rootmode=%o,user_id=0,group_id=0,fd=%d", unix.S_IFDIR, fd))
}

// returnResult returns the result from earlier operations.
// It is called with the control fd, a FUSE fd, and an error.
// If the error is not nil, then we are shutting down the cfd;
// If it is nil then we try to send the fd back.
// We return either e or the error result and e
func returnResult(cfd, ffd int, e error) error {
	if e != nil {
		if err := unix.Shutdown(cfd, unix.SHUT_RDWR); err != nil {
			return fmt.Errorf("shutting down after failed mount with %w: %w", e, err)
		}
		return e
	}
	oob := unix.UnixRights(int(ffd))
	if err := unix.Sendmsg(cfd, []byte(""), oob, nil, 0); err != nil {
		return fmt.Errorf("%s: %d: %w", CommFD, cfd, err)
	}
	return nil
}

func run(out io.Writer, args []string) error {
	var unmount, lazy, verbose bool

	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.BoolVar(&unmount, "unmount", false, "unmount")
	f.BoolVar(&unmount, "u", false, "unmount (shorthand)")

	f.BoolVar(&lazy, "lazy", false, "lazy unmount")
	f.BoolVar(&lazy, "z", false, "lazy unmount (shorthand)")

	f.BoolVar(&verbose, "verbose", false, "verbose")
	f.BoolVar(&verbose, "v", false, "verbose (shorthand)")

	f.Usage = func() {
		fmt.Fprintf(f.Output(), usage+"\n")
		f.PrintDefaults()
	}

	f.Parse(args[1:])

	if verbose {
		debug = log.Printf
	}

	if len(f.Args()) != 1 {
		f.Usage()
	}
	mpt = f.Arg(0)
	debug("mountpoint: %v", mpt)

	// We let "ability to open /dev/fuse" stand in as an indicator or
	// "we support FUSE".
	FuseFD, err := openFUSE()
	if err != nil {
		log.Printf("%v", err)
		os.Exit(int(syscall.ENOENT))
	}
	debug("FuseFD %v", FuseFD)

	// Bad design. All they had to do was make a -z and -u and have
	// them both mean unmount. Oh well.
	if lazy && !unmount {
		log.Fatalf("-z can only be used with -u")
	}

	// Fuse has to be seen to be believed.
	// The only interactive use of fusermount is to unmount
	if unmount {
		if err := umount(mpt, lazy); err != nil {
			return err
		}
	}

	if err := MountPointOK(mpt); err != nil {
		return err
	}

	if err := preMount(); err != nil {
		return err
	}

	cfd, err := getCommFD()
	if err != nil {
		return err
	}

	if err := doMount(FuseFD); err != nil {
		return err
	}

	if err := returnResult(cfd, FuseFD, err); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
