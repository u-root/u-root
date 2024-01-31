// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/util"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/uio"
	"github.com/u-root/uio/ulog"
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

	KexecOpts linux.KexecOptions
}

// LoadedLinuxImage is a processed version of LinuxImage.
//
// Main difference being that kernel and initrd is made as
// a read-only *os.File. There is also additional processing
// such as DTB, if available under KexecOpts, will be appended
// to Initrd.
type LoadedLinuxImage struct {
	Name        string
	Kernel      *os.File
	Initrd      *os.File
	Cmdline     string
	LoadSyscall bool
	KexecOpts   linux.KexecOptions
}

// loadedLinuxImageJSON is same as LoadedLinuxImage, but with transformed fields to help with serialization of LoadedLinuxImage.
type loadedLinuxImageJSON struct {
	Name        string
	KernelPath  string
	InitrdPath  string
	Cmdline     string
	LoadSyscall bool
	KexecOpts   linux.KexecOptions
}

var _ OSImage = &LinuxImage{}

var errNilKernel = errors.New("kernel image is empty, nothing to execute")

// MarshalJSON customizes marshaling for LoadedLinuxImage. It handles serializations
// for *os.File, so that kernel and initrd can be unmarshalled properly.
func (lli *LoadedLinuxImage) MarshalJSON() ([]byte, error) {
	lliJSON := loadedLinuxImageJSON{}
	// Sync and close kernel and initrd File object, and marshal paths to the files.
	if lli.Kernel != nil {
		if err := lli.Kernel.Sync(); err != nil {
			return nil, err
		}
		lliJSON.KernelPath = lli.Kernel.Name()
		if err := lli.Kernel.Close(); err != nil {
			return nil, err
		}
	}
	if lli.Initrd != nil {
		if err := lli.Initrd.Sync(); err != nil {
			return nil, err
		}
		lliJSON.InitrdPath = lli.Initrd.Name()
		if err := lli.Initrd.Close(); err != nil {
			return nil, err
		}
	}
	lliJSON.Name = lli.Name
	lliJSON.Cmdline = lli.Cmdline
	lliJSON.LoadSyscall = lli.LoadSyscall
	lliJSON.KexecOpts = lli.KexecOpts

	return json.Marshal(lliJSON)
}

// UnmarshalJSON customizes unmarshaling for LoadedLinuxImage. It processes kernel
// and initrd file by name, and opens a read-only copies for further execution.
func (lli *LoadedLinuxImage) UnmarshalJSON(b []byte) error {
	lliJSON := loadedLinuxImageJSON{}
	if err := json.Unmarshal(b, &lliJSON); err != nil {
		return err
	}
	if len(strings.TrimSpace(lliJSON.KernelPath)) > 0 {
		readOnlyK, err := os.Open(lliJSON.KernelPath)
		if err != nil {
			return err
		}
		lli.Kernel = readOnlyK
	}
	if len(strings.TrimSpace(lliJSON.InitrdPath)) > 0 {
		readOnlyI, err := os.Open(lliJSON.InitrdPath)
		if err != nil {
			return err
		}
		lli.Initrd = readOnlyI
	}
	lli.Name = lliJSON.Name
	lli.Cmdline = lliJSON.Cmdline
	lli.LoadSyscall = lliJSON.LoadSyscall
	lli.KexecOpts = lliJSON.KexecOpts
	return nil
}

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
	if li.KexecOpts.DTB != nil {
		labelInfo = append(
			labelInfo,
			fmt.Sprintf("dtb=%s", stringer(li.KexecOpts.DTB)),
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
		"LinuxImage(\n  Name: %s\n  Kernel: %s\n  Initrd: %s\n  Cmdline: %s\n  KexecOpts: %v\n)\n",
		li.Name, stringer(li.Kernel), stringer(li.Initrd), li.Cmdline, li.KexecOpts,
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

// loadLinuxImage processes given LinuxImage, and make it ready for kexec.
//
// For example:
//
//   - Acquiring a read-only copy of kernel and initrd as kernel
//     don't like them being opened for writting by anyone while
//     executing.
//   - Append DTB, if present to end of initrd.
func loadLinuxImage(li *LinuxImage, logger ulog.Logger, verbose bool) (*LoadedLinuxImage, func(), error) {
	if li.Kernel == nil {
		return nil, nil, errNilKernel
	}

	k, err := CopyToFileIfNotRegular(util.TryGzipFilter(li.Kernel), verbose)
	if err != nil {
		return nil, nil, err
	}

	// Append device-tree file to the end of initrd.
	if li.KexecOpts.DTB != nil {
		if li.Initrd != nil {
			li.Initrd = CatInitrds(li.Initrd, li.KexecOpts.DTB)
		} else {
			li.Initrd = li.KexecOpts.DTB
		}
	}

	var i *os.File
	if li.Initrd != nil {
		i, err = CopyToFileIfNotRegular(li.Initrd, verbose)
		if err != nil {
			return nil, nil, err
		}
	}

	logger.Printf("Kernel: %s", k.Name())
	if i != nil {
		logger.Printf("Initrd: %s", i.Name())
	}
	logger.Printf("Command line: %s", li.Cmdline)
	logger.Printf("KexecOpts: %#v", li.KexecOpts)

	cleanup := func() {
		k.Close()
		i.Close()
	}

	return &LoadedLinuxImage{
		Name:        li.Name,
		Kernel:      k,
		Initrd:      i,
		Cmdline:     li.Cmdline,
		LoadSyscall: li.LoadSyscall,
		KexecOpts:   li.KexecOpts,
	}, cleanup, nil
}

// saveLoadedLinuxImage marshals a LoadedLinuxImage to a json file.
//
// Marshalling it to a json versus a binary format, makes it readable
// and easier to tamper in case when we need change a field or two
// for experiments and debug.
//
// With that said, it is obvious that this saved info can then be
// loaded by a later kexec for further execution from an already loaded
// linux image and original load options.
func saveLoadedLinuxImage(lli *LoadedLinuxImage, p string) error {
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	out, err := json.Marshal(lli)
	if err != nil {
		return err
	}
	nw, err := f.Write(out)
	if err != nil {
		return err
	}
	if nw != len(out) {
		return fmt.Errorf("written %d bytes, want %d bytes", nw, len(out))
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return f.Close()
}

// Edit the kernel command line.
func (li *LinuxImage) Edit(f func(cmdline string) string) {
	li.Cmdline = f(li.Cmdline)
}

// Load implements OSImage.Load and kexec_load's the kernel with its initramfs.
func (li *LinuxImage) Load(opts ...LoadOption) error {
	loadOpts := defaultLoadOptions()
	for _, opt := range opts {
		opt(loadOpts)
	}

	loadedImage, cleanup, err := loadLinuxImage(li, loadOpts.logger, loadOpts.verbose)
	if err != nil {
		return err
	}
	defer cleanup()

	if !loadOpts.callKexecLoad {
		// If dryRun, serializes previously loaded linuxImage info to a file in tmpfs.
		// The info can be re-loaded for later kexec execution. It works b/c kernel and
		// initrd are already downloaded and saved into tmpfs.
		return saveLoadedLinuxImage(loadedImage, loadOpts.linuxImageCfgFile)
	}
	if li.LoadSyscall {
		return linux.KexecLoad(loadedImage.Kernel, loadedImage.Initrd, loadedImage.Cmdline, loadedImage.KexecOpts)
	}
	return kexec.FileLoad(loadedImage.Kernel, loadedImage.Initrd, loadedImage.Cmdline)
}
