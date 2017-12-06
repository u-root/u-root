// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// kexec executes a new kernel over the running kernel (u-root).
//
// Synopsis:
//     kexec [--initrd=FILE] [--command-line=STRING] [-l] [-e] [KERNELIMAGE]
//
// Description:
//		 Loads a kernel for later execution.
//
// Options:
//     --cmdline=STRING:       command line for kernel
//     --command-line=STRING:  command line for kernel
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
	"log"
	"os"

	"github.com/u-root/u-root/pkg/kexec"
)

type options struct {
	cmdline      string
	reuseCmdline bool
	initramfs    string
	load         bool
	exec         bool
}

func registerFlags() *options {
	o := &options{}
	flag.StringVar(&o.cmdline, "cmdline", "", "Set the kernel command line")
	flag.StringVar(&o.cmdline, "command-line", "", "Set the kernel command line")

	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")

	flag.StringVar(&o.initramfs, "i", "", "Use file as the kernel's initial ramdisk")
	flag.StringVar(&o.initramfs, "initrd", "", "Use file as the kernel's initial ramdisk")
	flag.StringVar(&o.initramfs, "ramdisk", "", "Use file as the kernel's initial ramdisk")

	flag.BoolVar(&o.load, "l", false, "Load the new kernel into the current kernel.")
	flag.BoolVar(&o.load, "load", false, "Load the new kernel into the current kernel.")

	flag.BoolVar(&o.exec, "e", false, "Execute a currently loaded kernel.")
	flag.BoolVar(&o.exec, "exec", false, "Execute a currently loaded kernel.")
	return o
}

func main() {
	opts := registerFlags()
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
		procCmdline, err := kexec.CurrentKernelCmdline()
		if err != nil {
			log.Fatalf("%v", err)
		}
		cmdline = procCmdline
	}

	if opts.load {
		kernelpath := flag.Args()[0]
		log.Printf("Loading %s for kernel\n", kernelpath)

		kernel, err := os.OpenFile(kernelpath, os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("open(%q): %v", kernelpath, err)
		}
		defer kernel.Close()

		var ramfs *os.File
		if opts.initramfs != "" {
			ramfs, err = os.OpenFile(opts.initramfs, os.O_RDONLY, 0)
			if err != nil {
				log.Fatalf("open(%q): %v", opts.initramfs, err)
			}
			defer ramfs.Close()
		}

		if err := kexec.FileLoad(kernel, ramfs, cmdline); err != nil {
			log.Fatalf("%v", err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
