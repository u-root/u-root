// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// kexec executes a new kernel over the running kernel (u-root).
//
// Synopsis:
//     kexec [--initrd=FILE] [--command-line=STRING] [-l] [-e] [KERNELIMAGE]
//
// Description:
//     This is only intended to be used with kexec_load_files, not the older
//     kexec.
//
// Options:
//     --cmdline=STRING:       command line for kernel
//     --command-line=STRING:  command line for kernel
//     --append=STRING:        command line for kernel
//
//     --reuse-commandline:    reuse command line from running system
//
//     --i=FILE:       initramfs file
//     --initrd=FILE:  initramfs file
//     --ramdisk=FILE: initramfs file
//
//     -l or --load:   only load the kernel
//     -e or --exec:   reboot with the currently loaded kernel
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"syscall"
	"unsafe"
)

type options struct {
	cmdline      string
	reuseCmdline bool
	initramfs    string
	load         bool
	exec         bool
}

func registerFlags(f *flag.FlagSet) *options {
	o := &options{}
	f.StringVar(&o.cmdline, "cmdline", "", "Set the kernel command line")
	f.StringVar(&o.cmdline, "command-line", "", "Set the kernel command line")
	f.StringVar(&o.cmdline, "append", "", "Set the kernel command line")

	f.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")

	f.StringVar(&o.initramfs, "i", "", "Use file as the kernel's initial ramdisk")
	f.StringVar(&o.initramfs, "initrd", "", "Use file as the kernel's initial ramdisk")
	f.StringVar(&o.initramfs, "ramdisk", "", "Use file as the kernel's initial ramdisk")

	f.BoolVar(&o.load, "l", false, "Load the new kernel into the current kernel.")
	f.BoolVar(&o.load, "load", false, "Load the new kernel into the current kernel.")

	f.BoolVar(&o.exec, "e", false, "Execute a currently loaded kernel.")
	f.BoolVar(&o.exec, "exec", false, "Execute a currently loaded kernel.")
	return o
}

// Syscall number for file kexec.
const _SYS_KEXEC_FILE_LOAD = 320

// Linux kexec syscall flags.
const (
	_KEXEC_FILE_UNLOAD       = 0x1
	_KEXEC_FILE_ON_CRASH     = 0x2
	_KEXEC_FILE_NO_INITRAMFS = 0x4
)

func kexec(kernelfd, ramfsfd int, cmdline string, flags int) error {
	cmd := append([]byte(cmdline), 0)
	if e1, e2, errno := syscall.Syscall6(
		_SYS_KEXEC_FILE_LOAD,
		uintptr(kernelfd),
		uintptr(ramfsfd),
		uintptr(len(cmd)),
		uintptr(unsafe.Pointer(&cmd[0])),
		uintptr(flags),
		0); errno != 0 {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = (%d, %d, errno %v)", kernelfd, ramfsfd, cmdline, flags, e1, e2, errno)
	}
	return nil
}

func main() {
	opts := registerFlags(flag.CommandLine)
	flag.Parse()

	if opts.exec == false && len(flag.Args()) == 0 {
		flag.PrintDefaults()
		log.Fatalf("usage: kexec [flags] kernelname OR kexec -e")
	}

	if opts.cmdline != "" && opts.reuseCmdline {
		flag.PrintDefaults()
		log.Fatalf("--reuse-cmdline and other command line options are mutually exclusive")
	}

	if opts.load == false && opts.exec == false {
		opts.load = true
		opts.exec = true
	}

	cmdline := opts.cmdline
	if opts.reuseCmdline {
		procCmdline, err := ioutil.ReadFile("/proc/cmdline")
		if err != nil {
			log.Fatalf("%v", err)
		}
		cmdline = string(procCmdline)
	}

	if opts.load {
		var flags int

		kernel := flag.Args()[0]
		log.Printf("Loading %s for kernel\n", kernel)
		kernelfd, err := syscall.Open(kernel, syscall.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("open(%s): %v", kernel, err)
		}

		ramfsfd, err := syscall.Open(opts.initramfs, syscall.O_RDONLY, 0)
		if err != nil {
			flags |= _KEXEC_FILE_NO_INITRAMFS
		}

		if err := kexec(kernelfd, ramfsfd, cmdline, flags); err != nil {
			log.Fatalf("%v", err)
		}
	}

	if opts.exec {
		if e1, e2, errno := syscall.Syscall(syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, syscall.LINUX_REBOOT_CMD_KEXEC); errno != 0 {
			log.Fatalf("sys_reboot(..., kexec) = (%d, %d, errno %v)", e1, e2, errno)
		}
	}
}
