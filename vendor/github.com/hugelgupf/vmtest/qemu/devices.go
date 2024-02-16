// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ErrInvalidDir is used when no directory is specified for file sharing.
var ErrInvalidDir = errors.New("no directory specified")

// ErrInvalidTag is used when no tag is specified for 9P file system sharing.
var ErrInvalidTag = errors.New("no tag specified for 9P file system")

// ErrIsNotDir is used when the directory specified for file sharing is not a directory.
var ErrIsNotDir = errors.New("file system sharing requires directory")

// IDAllocator is used to ensure no overlapping QEMU option IDs.
type IDAllocator struct {
	// maps a prefix to the maximum used suffix number.
	idx map[string]uint32
}

// NewIDAllocator returns a new ID allocator for QEMU option IDs.
func NewIDAllocator() *IDAllocator {
	return &IDAllocator{
		idx: make(map[string]uint32),
	}
}

// ID returns the next available ID for the given prefix.
func (a *IDAllocator) ID(prefix string) string {
	prefix = strings.TrimRight(prefix, "0123456789")
	idx := a.idx[prefix]
	a.idx[prefix]++
	return fmt.Sprintf("%s%d", prefix, idx)
}

// ReadOnlyDirectory adds args that expose a directory as a /dev/sda1 readonly
// vfat partition in the VM guest.
func ReadOnlyDirectory(dir string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		if len(dir) == 0 {
			return ErrInvalidDir
		}
		if fi, err := os.Stat(dir); err != nil {
			return fmt.Errorf("cannot access directory %s to be shared with guest: %w", dir, err)
		} else if !fi.IsDir() {
			return &os.PathError{
				Op:   "9P-directory-sharing",
				Path: dir,
				Err:  fmt.Errorf("%w: is %s", ErrIsNotDir, fi.Mode().Type()),
			}
		}

		drive := alloc.ID("drive")
		ahci := alloc.ID("ahci")

		// Expose the temp directory to QEMU as /dev/sda1
		opts.AppendQEMU(
			"-drive", fmt.Sprintf("file=fat:rw:%s,if=none,id=%s", dir, drive),
			"-device", fmt.Sprintf("ich9-ahci,id=%s", ahci),
			"-device", fmt.Sprintf("ide-hd,drive=%s,bus=%s.0", drive, ahci),
		)
		return nil
	}
}

// IDEBlockDevice emulates an AHCI/IDE block device.
func IDEBlockDevice(file string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		if _, err := os.Stat(file); err != nil {
			return fmt.Errorf("cannot access file %s to be shared with guest: %w", file, err)
		}

		drive := alloc.ID("drive")
		ahci := alloc.ID("ahci")

		opts.AppendQEMU(
			"-drive", fmt.Sprintf("file=%s,if=none,id=%s", file, drive),
			"-device", fmt.Sprintf("ich9-ahci,id=%s", ahci),
			"-device", fmt.Sprintf("ide-hd,drive=%s,bus=%s.0", drive, ahci),
		)
		return nil
	}
}

// P9Directory adds QEMU args that expose a directory as a Plan9 (9p)
// read-write filesystem in the VM.
//
// dir is the directory to expose as read-write 9p filesystem.
//
// tag is an identifier that is used within the VM when mounting an fs, e.g.
// 'mount -t 9p my-vol-ident mountpoint'. The tag must be unique for each dir.
//
// P9Directory will add a kernel cmdline argument in the style of
// VMTEST_MOUNT9P_$qemuID=$tag. Likely this is only useful on Linux. The
// vmmount command in vminit/vmmount can be used to mount 9P directories passed
// to the VM this way at /mount/9p/$tag in the guest. See the example in
// ./examples/shareddir.
func P9Directory(dir string, tag string) Fn {
	return p9Directory(dir, false, tag)
}

// P9BootDirectory adds QEMU args that expose a directory as a Plan9 (9p)
// read-write filesystem in the VM as the boot device.
//
// The directory will be used as the root volume. There can only be one boot
// 9pfs at a time. The tag used will be /dev/root, and Linux kernel args will
// be appended to mount it as the root file system.
func P9BootDirectory(dir string) Fn {
	return p9Directory(dir, true, "/dev/root")
}

func p9Directory(dir string, boot bool, tag string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		if len(dir) == 0 {
			return fmt.Errorf("%w for shared 9P file system", ErrInvalidDir)
		}
		if len(tag) == 0 {
			return ErrInvalidTag
		}
		if fi, err := os.Stat(dir); err != nil {
			return fmt.Errorf("cannot access directory %s to be shared with guest: %w", dir, err)
		} else if !fi.IsDir() {
			return &os.PathError{
				Op:   "9P-directory-sharing",
				Path: dir,
				Err:  fmt.Errorf("%w: is %s", ErrIsNotDir, fi.Mode().Type()),
			}
		}

		var id string
		if boot {
			id = "rootdrv"
		} else {
			id = alloc.ID("fsdev")
		}

		// Expose the temp directory to QEMU
		var deviceArgs string
		switch opts.Arch() {
		case ArchArm:
			deviceArgs = fmt.Sprintf("virtio-9p-device,fsdev=%s,mount_tag=%s", id, tag)
		default:
			deviceArgs = fmt.Sprintf("virtio-9p-pci,fsdev=%s,mount_tag=%s", id, tag)
		}

		opts.AppendQEMU(
			// security_model=mapped-file seems to be the best choice. It gives
			// us control over uid/gid/mode seen in the guest, without requiring
			// elevated perms on the host.
			"-fsdev", fmt.Sprintf("local,id=%s,path=%s,security_model=mapped-file", id, dir),
			"-device", deviceArgs,
		)
		if boot {
			opts.AppendKernel(
				"devtmpfs.mount=1",
				"root=/dev/root",
				"rootfstype=9p",
				"rootflags=trans=virtio,version=9p2000.L",
			)
		} else {
			opts.AppendKernel(fmt.Sprintf("VMTEST_MOUNT9P_%s=%s", id, tag))
		}
		return nil
	}
}

// VirtioRandom adds QEMU args that expose a PCI random number generator to the
// guest VM.
func VirtioRandom() Fn {
	return ArbitraryArgs("-device", "virtio-rng-pci")
}

// ArbitraryArgs adds arbitrary arguments to the QEMU command line.
func ArbitraryArgs(aa ...string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.AppendQEMU(aa...)
		return nil
	}
}

// WithQEMUArgs adds arguments to the QEMU command line.
func WithQEMUArgs(aa ...string) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.AppendQEMU(aa...)
		return nil
	}
}

// HaltOnKernelPanic passes args to QEMU and kernel to halt when the kernel
// panics.
//
// Linux's default behavior is to hang forever, which is not great test
// behavior.
func HaltOnKernelPanic() Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		opts.AppendQEMU("-no-reboot")
		opts.AppendKernel("panic=-1")
		return nil
	}
}

func replaceCtl(str []byte) []byte {
	for i, c := range str {
		if c == 9 || c == 10 {
		} else if c < 32 || c == 127 {
			str[i] = '~'
		}
	}
	return str
}

// LinePrinter prints one line to some output.
type LinePrinter func(line string)

// LogSerialByLine processes serial output from the guest one line at a time
// and calls callback on each full line.
func LogSerialByLine(callback LinePrinter) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		r, w := io.Pipe()
		opts.SerialOutput = append(opts.SerialOutput, w)
		opts.Tasks = append(opts.Tasks, WaitVMStarted(func(ctx context.Context, n *Notifications) error {
			s := bufio.NewScanner(r)
			for s.Scan() {
				callback(string(replaceCtl(s.Bytes())))
			}
			if err := s.Err(); err != nil {
				return fmt.Errorf("error reading serial from VM: %w", err)
			}
			return nil
		}))
		return nil
	}
}

// TS prefixes line printer output with a timestamp since the first log line.
//
// format can be any Time.Format format string. Recommendations are
// time.TimeOnly or time.DateTime.
func TS(format string, printer LinePrinter) LinePrinter {
	return func(line string) {
		printer(fmt.Sprintf("[%s] %s", time.Now().Format(format), line))
	}
}

// DefaultPrint is the default LinePrinter, adding a prefix and relative timestamp.
func DefaultPrint(prefix string, printer func(fmt string, arg ...any)) LinePrinter {
	return RelativeTS(Prefix(prefix, PrintLine(printer)))
}

// RelativeTS prefixes line printer output with "[%06.4fs] " seconds since the
// first log line.
func RelativeTS(printer LinePrinter) LinePrinter {
	start := sync.OnceValue(time.Now)
	return func(line string) {
		printer(fmt.Sprintf("[%06.4fs] %s", time.Since(start()).Seconds(), line))
	}
}

// PrintLine is a LinePrinter that prints to a standard "formatter" like testing.TB.Logf or fmt.Printf.
func PrintLine(printer func(fmt string, arg ...any)) LinePrinter {
	return func(line string) {
		printer("%s", line)
	}
}

// Prefix returns a LinePrinter that prefixes the given LinePrinter with "prefix: ".
func Prefix(prefix string, printer LinePrinter) LinePrinter {
	return func(line string) {
		printer(fmt.Sprintf("%s: %s", prefix, line))
	}
}

// ByArch applies only the Fn config function applicable to the VM guest
// architecture.
func ByArch(m map[Arch]Fn) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		a := opts.Arch()
		fn, ok := m[a]
		if !ok {
			return nil
		}
		return fn(alloc, opts)
	}
}

// IfNotArch applies fn only if the VM guest arch is not the given arch.
func IfNotArch(arch Arch, fn Fn) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		if opts.Arch() == arch {
			return nil
		}
		return fn(alloc, opts)
	}
}

// IfArch applies fn only if the VM guest arch is the given arch.
func IfArch(arch Arch, fn Fn) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		if opts.Arch() == arch {
			return fn(alloc, opts)
		}
		return nil
	}
}

// All applies all given configurators in order. If an error occurs, it returns
// the error early.
func All(fn ...Fn) Fn {
	return func(alloc *IDAllocator, opts *Options) error {
		for _, f := range fn {
			if err := f(alloc, opts); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithVmtestIdent adds VMTEST_IN_GUEST=1 to kernel commmand-line.
//
// Tests may use this env var to identify they are running inside a vmtest
// using guest.SkipIfNotInVM or guest.SkipIfInVM.
func WithVmtestIdent() Fn {
	return WithAppendKernel("VMTEST_IN_GUEST=1")
}
