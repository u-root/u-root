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
//     --cmdline=STRING or -c=STRING: Set the kernel command line
//     --reuse-commandline:           Use the kernel command line from running system
//     --i=FILE or --initrd=FILE:     Use file as the kernel's initial ramdisk
//     -l or --load:                  Load the new kernel into the current kernel
//     -e or --exec:		      Execute a currently loaded kernel
package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/cmdline"
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
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Set the kernel command line")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	return o
}

func main() {
	opts := registerFlags()
	flag.Parse()

	if (opts.exec == false && flag.NArg() == 0) || flag.NArg() > 1 {
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

	newCmdLine := opts.cmdline
	if opts.reuseCmdline {
		procCmdLine := cmdline.NewCmdLine()
		if procCmdLine.Err != nil {
			log.Fatal("Couldn't read /proc/cmdline")
		} else {
			newCmdLine = procCmdLine.Raw
		}
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

		if err := kexec.FileLoad(kernel, ramfs, newCmdLine); err != nil {
			log.Fatalf("%v", err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
