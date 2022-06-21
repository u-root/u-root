// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

// LinuxImage implements OSImage for a Linux kernel + initramfs.
type LinuxImage struct {
	Name string

	Kernel      io.ReaderAt
	Initrd      io.ReaderAt
	Cmdline     string
	BootRank    int
	LoadSyscall bool
	DeviceTree  io.ReaderAt

	KexecOpts linux.KexecOptions
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
	labelInfo := []string{
		fmt.Sprintf("kernel=%s", stringer(li.Kernel)),
	}
	if li.Initrd != nil {
		labelInfo = append(
			labelInfo,
			fmt.Sprintf("initrd=%s", stringer(li.Initrd)),
		)
	}
	if li.DeviceTree != nil {
		labelInfo = append(
			labelInfo,
			fmt.Sprintf("dtb=%s", stringer(li.DeviceTree)),
		)
	}

	return fmt.Sprintf("Linux(%s)", strings.Join(labelInfo, " "))
}

// Rank for the boot menu order
func (li *LinuxImage) Rank() int {
	return li.BootRank
}

// String prints a human-readable version of this linux image.
func (li *LinuxImage) String() string {
	return fmt.Sprintf(
		"LinuxImage(\n  Name: %s\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n  Dtb: %s\n)\n",
		li.Name, stringer(li.Kernel), stringer(li.Initrd), li.Cmdline, stringer(li.DeviceTree),
	)
}

// copyToFileIfNotRegular copies given io.ReadAt to a tmpfs file when
// necessary. It skips copying when source file is a regular file under
// tmpfs or ramfs, and it is not opened for writing.
//
// Copy is necessary for other cases, such as when the reader is an io.File
// but not sufficient for kexec, as os.File could be a socket, a pipe or
// some other strange thing. Also kexec_file_load will fail (similar to
// execve) if anything has the file opened for writing. That's unfortunately
// something we can't guarantee here - unless we make a copy of the file
// and dump it somewhere.
func copyToFileIfNotRegular(r io.ReaderAt, verbose bool) (*os.File, error) {
	// If source is a regular file in tmpfs, simply re-use that than copy.
	//
	// The assumption (bad?) is original local file was opened as a type
	// conforming to os.File. We then can derive file descriptor, and the
	// name.
	if f, ok := r.(*os.File); ok {
		if fi, err := f.Stat(); err == nil && fi.Mode().IsRegular() {
			if r, _ := mount.IsTmpRamfs(f.Name()); r {
				// Check if original file is opened for write. Perform copy to
				// get a read only version in that case, than directly return
				// as kexec load would fail on a file opened for write.
				wr := unix.O_RDWR | unix.O_WRONLY
				if am, err := unix.FcntlInt(f.Fd(), unix.F_GETFL, 0); err == nil && am&wr == 0 {
					return f, nil
				}
				// Original file is either opened for write, or it failed to
				// check (possibly, current kernel is too old and not supporting
				// the sys call cmd)
			}
			// Original file is neither on a tmpfs, nor a ramfs.
		}
		// Not a regular file, or could not confirm it is a regular file.
	}

	rdr := uio.Reader(r)

	if verbose {
		// In verbose mode, print a dot every 5MiB. It is not pretty,
		// but it at least proves the files are still downloading.
		progress := func(r io.Reader) io.Reader {
			return &uio.ProgressReadCloser{
				RC:       io.NopCloser(r),
				Symbol:   ".",
				Interval: 5 * 1024 * 1024,
				W:        os.Stdout,
			}
		}
		rdr = progress(rdr)
	}

	f, err := os.CreateTemp("", "kexec-image")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, rdr); err != nil {
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

	k, err := copyToFileIfNotRegular(util.TryGzipFilter(li.Kernel), verbose)
	if err != nil {
		return err
	}
	defer k.Close()

	// Append device-tree file to the end of initrd
	if li.DeviceTree != nil {
		if li.Initrd != nil {
			li.Initrd = CatInitrds(li.Initrd, li.DeviceTree)
		} else {
			li.Initrd = li.DeviceTree
		}
	}

	var i *os.File
	if li.Initrd != nil {
		i, err = copyToFileIfNotRegular(li.Initrd, verbose)
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
		if li.DeviceTree != nil {
			log.Print("Device tree loaded: true")
		}
	}

	if li.LoadSyscall {
		return linux.KexecLoad(k, i, li.Cmdline, li.KexecOpts)
	}
	return kexec.FileLoad(k, i, li.Cmdline)
}
