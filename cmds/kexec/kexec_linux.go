// Copyright 2015-2018 the u-root Authors. All rights reserved
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
//     -e or --exec:                  Execute a currently loaded kernel
package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot"
	"github.com/u-root/u-root/pkg/uio"
)

type options struct {
	cmdline      string
	reuseCmdline bool
	initramfs    string
	load         bool
	exec         bool
	debug        bool
	modules      []string
}

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Set the kernel command line")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.BoolVarP(&o.debug, "debug", "d", false, "Print debug info")
	flag.StringSliceVar(&o.modules, "module", nil, `Load module with command line args (e.g --module="mod arg1")`)
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

	newCmdline := opts.cmdline
	if opts.reuseCmdline {
		procCmdLine := cmdline.NewCmdLine()
		if procCmdLine.Err != nil {
			log.Fatal("Couldn't read /proc/cmdline")
		} else {
			newCmdline = procCmdLine.Raw
		}
	}

	// mbk indicates that we are a multiboot kernel
	var mbk bool
	var kernelpath string
	if opts.load {
		kernelpath = flag.Arg(0)
		log.Printf("Loading %s for kernel.", kernelpath)

		if err := multiboot.Probe(kernelpath); err == nil {
			mbk = true
		} else if err == multiboot.ErrFlagsNotSupported {
			log.Fatal(err)
		}
	}
	if opts.load {
		var image boot.OSImage
		image = &boot.LinuxImage{
			Kernel:  uio.NewLazyFile(kernelpath),
			Initrd:  uio.NewLazyFile(opts.initramfs),
			Cmdline: newCmdline,
		}
		if mbk {
			log.Printf("%s is a multiboot v1 kernel.", kernelpath)
			image = &boot.MultibootImage{
				Debug:   opts.debug,
				Modules: opts.modules,
				Path:    kernelpath,
				Cmdline: newCmdline,
			}
		}
		if err := image.Load(); err != nil {
			log.Fatal(err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
