// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Opts contains options to load an image for kexec on Linux.
type LinuxOpts struct {
	// The image to load. Required.
	Kernel *os.File

	// The initramfs to pass to the next kernel. The caller is responsible for
	// closing this file. Optional.
	Initramfs *os.File

	// The command line for the next kernel.
	CmdLine string

	// Device tree to use for the next kernel. Optional.
	DTB io.ReadSeeker

	// Function to call to perform the kexec_load syscall. For dryrun and mocking.
	LoadSyscall func(entry uintptr, segments []Segment, flags uint64) error

	// Function to call to perform the kexec_file_load syscall. For dryrun and mocking.
	FileLoadSyscall func(kernelFd int, initrd int, cmdline string, flags int) error
}

func DryrunOpts() *LinuxOpts {
	return &LinuxOpts{
		LoadSyscall:     DryrunLoad,
		FileLoadSyscall: DryrunFileLoad,
	}
}

// SupportedLoaders contains a list of kexec loaders supported by the current
// arch. Even if a loader is listed, it might return a ProbeError at runtime to
// indicate the given set of options and/or kernel type is not supported. Use the
// Load function to find a suitable loader from this list.
var SupportedLoaders = []func(*LinuxOpts) error{}

// ProbeError
type ProbeError struct{ error }

// Load determines which type of kernel you're trying to load and loads it.
func Load(opts *LinuxOpts) error {
	for _, loader := range SupportedLoaders {
		err := loader(opts)
		if _, ok := err.(*ProbeError); !ok {
			return err
		}
	}
	return fmt.Errorf("failed to determine file type")
}

// KexecError is the error type returned by kexec_load.
type KexecError struct {
	Entry    uintptr
	Segments []Segment
	Flags    uint64
	Errno    unix.Errno
}

func (e KexecError) Error() string {
	return fmt.Sprintf("kexec_load error: %v (errno %d, entry %#x, flags %#x, segments %#v)",
		e.Errno, e.Errno, e.Entry, e.Flags, e.Segments)
}

// RawLoad is a wrapper around the kexec_load(2) syscall.
// Preconditions:
// - segments must not overlap
// - segments must be full pages
func RawLoad(entry uintptr, segments []Segment, flags uint64) error {
	if _, _, errno := unix.Syscall6(
		unix.SYS_KEXEC_LOAD,
		entry,
		uintptr(len(segments)),
		uintptr(unsafe.Pointer(&segments[0])),
		uintptr(flags),
		0, 0); errno != 0 {
		return KexecError{
			Entry:    entry,
			Segments: segments,
			Flags:    flags,
			Errno:    errno,
		}
	}
	return nil
}

// Reboot executes a previously loaded kernel.
func Reboot() error {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_KEXEC); err != nil {
		return fmt.Errorf("sys_reboot(..., kexec) = %v", err)
	}
	return nil
}

func DryrunLoad(entry uintptr, segments []Segment, flags uint64) error {
	log.Printf("kexec_load(2):")
	log.Printf("  entry: %#v", entry)
	for i, s := range segments {
		log.Printf("  segments[%d]: %v", i, s)
	}
	log.Printf("  flags: %#v", flags)
	return nil
}

func DryrunFileLoad(kernelFd int, initrd int, cmdline string, flags int) error {
	log.Printf("kexec_file_load(2):")
	log.Printf("  kernelFd: %d", kernelFd)
	log.Printf("  initrd:   %d", initrd)
	log.Printf("  cmdline:  %q", cmdline)
	log.Printf("  flags:    %#v", flags)
	return nil
}
