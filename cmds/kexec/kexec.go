// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// kexec command in Go.
// This is only intended to be used with kexec_load_files, not the older kexec.
package main

// N.B. /**/ comments are verbatim from uapi/linux/kexec.h.
/* kexec system call -  It loads the new kernel to boot into.
 * kexec does not sync, or unmount filesystems so if you need
 * that to happen you need to do that yourself.
 */

import (
	"flag"
	"io/ioutil"
	"log"
	"syscall"
	"unsafe"
)

/* kexec flags for different usage scenarios */
const (
	KEXEC_FILE_UNLOAD       = 0x1
	KEXEC_FILE_ON_CRASH     = 0x2
	KEXEC_FILE_NO_INITRAMFS = 0x4
)

var (
	dryrun        = flag.Bool("dryrun", false, "Do not do kexec system calls")
	cmdline       = flag.String("cmdline", "", "Command line for kernel")
	initramfs     = flag.String("i", "", "initramfs")
	kern      int = -1
	ramfs     int = -1
	flags     uintptr
)

func main() {
	var err error
	var b []byte
	var l uintptr

	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.PrintDefaults()
		log.Fatalf("usage: kexec [flags] kernelname")
	}
	kernel := flag.Args()[0]

	if *cmdline != "" {
		b = append(b, []byte(*cmdline)...)
		l = uintptr(len(b)) + 1
	} else {
		b, err = ioutil.ReadFile("/proc/cmdline")
		if err != nil {
			log.Fatalf("%v", err)
		}
		b[len(b)-1] = 0
		l = uintptr(len(b))
	}

	p := uintptr(unsafe.Pointer(&b[0]))

	log.Printf("Loading %v\n", kernel)

	if kern, err = syscall.Open(kernel, syscall.O_RDONLY, 0); err != nil {
		log.Fatalf("%v", err)
	}

	if ramfs, err = syscall.Open(*initramfs, syscall.O_RDONLY, 0); err != nil {
		flags |= KEXEC_FILE_NO_INITRAMFS
	}

	log.Printf("command line: '%v'", string(b))
	log.Printf("%v %v %v %v %v %v", 320, uintptr(kern), uintptr(ramfs), p, l, flags)
	if *dryrun {
		log.Printf("Dry run -- exiting now")
		return
	}
	e1, e2, err := syscall.Syscall6(320, uintptr(kern), uintptr(ramfs), l, p, flags, uintptr(0))
	log.Printf("a %v b %v err %v", e1, e2, err)

	e1, e2, err = syscall.Syscall6(syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, syscall.LINUX_REBOOT_CMD_KEXEC, 0, 0, 0)

	log.Printf("a %v b %v err %v", e1, e2, err)
}
