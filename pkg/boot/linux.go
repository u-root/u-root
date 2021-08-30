// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/uio"
)

// LinuxImage implements OSImage for a Linux kernel + initramfs.
type LinuxImage struct {
	Name string

	Kernel   io.ReaderAt
	Initrd   io.ReaderAt
	Cmdline  string
	BootRank int
}

var _ OSImage = &LinuxImage{}

// named is satisifed by both *os.File and *vfile.File. Hack hack hack.
type named interface {
	Name() string
}

func stringer(mod io.ReaderAt) string {
	if s, ok := mod.(fmt.Stringer); ok {
		return s.String()
	}
	if f, ok := mod.(named); ok {
		return f.Name()
	}
	return fmt.Sprintf("%T", mod)
}

// Label returns either the Name or a short description.
func (li *LinuxImage) Label() string {
	if len(li.Name) > 0 {
		return li.Name
	}
	return fmt.Sprintf("Linux(kernel=%s initrd=%s)", stringer(li.Kernel), stringer(li.Initrd))
}

// Rank for the boot menu order
func (li *LinuxImage) Rank() int {
	return li.BootRank
}

// String prints a human-readable version of this linux image.
func (li *LinuxImage) String() string {
	return fmt.Sprintf("LinuxImage(\n  Name: %s\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n)\n", li.Name, stringer(li.Kernel), stringer(li.Initrd), li.Cmdline)
}

func copyToFileIfNotRegular(r io.Reader) (*os.File, error) {
	// If source is a regular file in tmpfs, simply re-use that than copy.
	//
	// The assumption (bad?) is original local file was opened as a type
	// conforming to os.File. We then can derive file descriptor, and the
	// name.
	if f, ok := r.(*os.File); ok {
		if fi, err := f.Stat(); err == nil && fi.Mode().IsRegular() { // Name is the path.
			if r, _ := mount.IsTmpfs(f.Name()); r {
				return f, nil
				// This is not checking if same file is linked by another
				// file descriptor for writing. It is tempting to check or
				// guarantee, by this point, the file passed in for kexec
				// is not opened for writing. Caller in the same routine
				// probably has done its work on the file before it decides
				// to kexec file load it.
				//
				// Same file maybe opened by another process for writing
				// which is theoretically possible, but may not need to be
				// paranoid on.
				//
				// But kexec file load will propogate up an error if the
				// file is opened for writting, which is not too different
				// from we handling it, and return an userspace error to
				// caller, or not call into file load. Here, it simply leaps
				// before it looks, and let file load system call errors up.
			}
		}
		// Skip scenarios where it does not exist, and we don't know
		// if it exists.
	}

	f, err := ioutil.TempFile("", "kexec-image")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	if err := f.Sync(); err != nil {
		return nil, err
	}

	readOnlyF, err := os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	return readOnlyF, nil
}

// Edit the kernel command line.
func (li *LinuxImage) Edit(f func(cmdline string) string) {
	li.Cmdline = f(li.Cmdline)
}

// Load implements OSImage.Load and kexec_load's the kernel with its initramfs.
func (li *LinuxImage) Load(verbose bool) error {
	if li.Kernel == nil {
		return errors.New("LinuxImage.Kernel must be non-nil")
	}

	kernel, initrd := uio.Reader(util.TryGzipFilter(li.Kernel)), uio.Reader(li.Initrd)
	if verbose {
		// In verbose mode, print a dot every 5MiB. It is not pretty,
		// but it at least proves the files are still downloading.
		progress := func(r io.Reader) io.Reader {
			return &uio.ProgressReadCloser{
				RC:       ioutil.NopCloser(r),
				Symbol:   ".",
				Interval: 5 * 1024 * 1024,
				W:        os.Stdout,
			}
		}
		kernel = progress(kernel)
		initrd = progress(initrd)
	}

	// It seems inefficient to always copy, in particular when the reader
	// is an io.File but that's not sufficient, os.File could be a socket,
	// a pipe or some other strange thing. Also kexec_file_load will fail
	// (similar to execve) if anything as the file opened for writing.
	// That's unfortunately something we can't guarantee here - unless we
	// make a copy of the file and dump it somewhere.
	k, err := copyToFileIfNotRegular(kernel)
	if err != nil {
		return err
	}
	defer k.Close()

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFileIfNotRegular(initrd)
		if err != nil {
			return err
		}
		defer i.Close()
	}

	if verbose {
		log.Printf("Kernel: %s", k.Name())
		if i != nil {
			log.Printf("Initrd: %s", i.Name())
		}
		log.Printf("Command line: %s", li.Cmdline)
	}
	return kexec.FileLoad(k, i, li.Cmdline)
}
