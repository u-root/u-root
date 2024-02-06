// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/uio/uio"
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
	DTB         io.ReaderAt

	// ReservedRanges are additional physical memory pieces that will be
	// avoided when allocating kexec segments. Only used for LoadSyscall.
	//
	// ReservedRanges will not be shared with the next kernel, which is
	// free to use this memory unless some other mechanism (such as
	// memmap=) reserves it.
	ReservedRanges kexec.Ranges
}

var _ OSImage = &LinuxImage{}

var errNilKernel = errors.New("kernel image is empty, nothing to execute")

// named is satisifed by *os.File.
type named interface {
	Name() string
}

func stringer(mod interface{}) string {
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
	if li.DTB != nil {
		labelInfo = append(
			labelInfo,
			fmt.Sprintf("dtb=%s", stringer(li.DTB)),
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
		"LinuxImage(\n  Name: %s\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n  DTB: %v\n)\n",
		li.Name, stringer(li.Kernel), stringer(li.Initrd), li.Cmdline, stringer(li.DTB),
	)
}

func isTmpfsReadOnlyFile(f *os.File) bool {
	if f == nil {
		return false
	}

	if fi, err := f.Stat(); err == nil && fi.Mode().IsRegular() {
		if r, _ := mount.IsTmpRamfs(f.Name()); r {
			// Check if original file is opened for write. Perform copy to
			// get a read only version in that case, than directly return
			// as kexec load would fail on a file opened for write.
			wr := unix.O_RDWR | unix.O_WRONLY
			if am, err := unix.FcntlInt(f.Fd(), unix.F_GETFL, 0); err == nil && am&wr == 0 {
				return true
			}
			// Original file is either opened for write, or it failed to
			// check (possibly, current kernel is too old and not supporting
			// the sys call cmd)
		}
		// Original file is neither on a tmpfs, nor a ramfs.
	}
	// Not a regular file, or could not confirm it is a regular file.
	return false
}

// CopyToFileIfNotRegular copies given io.ReadAt to a tmpfs file when
// necessary. It skips copying when source file is a regular file under
// tmpfs or ramfs, and it is not opened for writing.
//
// Copy is necessary for other cases, such as when the reader is an io.File
// but not sufficient for kexec, as os.File could be a socket, a pipe or
// some other strange thing. Also kexec_file_load will fail (similar to
// execve) if anything has the file opened for writing. That's unfortunately
// something we can't guarantee here - unless we make a copy of the file
// and dump it somewhere.
func CopyToFileIfNotRegular(r io.ReaderAt, verbose bool) (*os.File, error) {
	// If source is a regular file in tmpfs, simply re-use that than copy.
	//
	// The assumption (bad?) is original local file was opened as a type
	// conforming to os.File. We then can derive file descriptor, and the
	// name.
	if f, ok := r.(*os.File); ok {
		if isTmpfsReadOnlyFile(f) {
			return f, nil
		}
	}

	// For boot entries whose image is lazily downloaded, try booting the
	// backend file directly if possible.
	if lor, ok := r.(*uio.LazyOpenerAt); ok {
		// It is lazy, so read one byte to make sure backend file
		// is created.
		if err := uio.ReadOneByte(r); err != nil {
			return nil, err
		}
		if isTmpfsReadOnlyFile(lor.File()) {
			return lor.File(), nil
		}
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

func (li *LinuxImage) loadImage(loadOpts *loadOptions) (*os.File, *os.File, error) {
	if li.Kernel == nil {
		return nil, nil, errNilKernel
	}

	k, err := CopyToFileIfNotRegular(util.TryGzipFilter(li.Kernel), loadOpts.verbose)
	if err != nil {
		return nil, nil, err
	}

	// Append device-tree file to the end of initrd.
	if li.DTB != nil {
		if li.Initrd != nil {
			li.Initrd = CatInitrds(li.Initrd, li.DTB)
		} else {
			li.Initrd = li.DTB
		}
	}

	var i *os.File
	if li.Initrd != nil {
		i, err = CopyToFileIfNotRegular(li.Initrd, loadOpts.verbose)
		if err != nil {
			k.Close()
			return nil, nil, err
		}
	}
	return k, i, nil
}

// Load implements OSImage.Load and kexec_load's the kernel with its initramfs.
func (li *LinuxImage) Load(opts ...LoadOption) error {
	loadOpts := defaultLoadOptions()
	for _, opt := range opts {
		opt(loadOpts)
	}

	k, i, err := li.loadImage(loadOpts)
	if err != nil {
		return err
	}
	defer k.Close()
	if i != nil {
		defer i.Close()
	}

	loadOpts.logger.Printf("Kernel: %s", k.Name())
	if i != nil {
		loadOpts.logger.Printf("Initrd: %s", i.Name())
	}
	loadOpts.logger.Printf("Command line: %s", li.Cmdline)
	loadOpts.logger.Printf("DTB: %#v", li.DTB)

	if !loadOpts.callKexecLoad {
		return nil
	}
	if li.LoadSyscall {
		return linux.KexecLoad(k, i, li.Cmdline, li.DTB, li.ReservedRanges)
	}
	return kexec.FileLoad(k, i, li.Cmdline)
}
