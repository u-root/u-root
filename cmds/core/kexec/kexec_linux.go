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
//      --append string        Append to the kernel command line
//  -c, --cmdline string       Append to the kernel command line
//  -d, --debug                Print debug info (default true)
//  -e, --exec                 Execute a currently loaded kernel
//  -x, --extra string         Add a cpio containing extra files
//      --initramfs string     Use file as the kernel's initial ramdisk
//  -i, --initrd string        Use file as the kernel's initial ramdisk
//  -l, --load                 Load the new kernel into the current kernel
//  -L, --loadsyscall          Use the kexec load syscall (not file_load) (default true)
//      --module stringArray   Load multiboot module with command line args (e.g --module="mod arg1")
//  -p, --purgatory string     pick a purgatory, use '-p xyz' to get a list (default "default")
//      --reuse-cmdline        Use the kernel command line from running system

package main

import (
	"io"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/boot/purgatory"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/uio"
)

type options struct {
	cmdline      string
	debug        bool
	dtb          string
	exec         bool
	extra        string
	initramfs    string
	load         bool
	loadSyscall  bool
	modules      []string
	purgatory    string
	reuseCmdline bool
}

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Append to the kernel command line")
	flag.StringVar(&o.cmdline, "append", "", "Append to the kernel command line")
	flag.BoolVarP(&o.debug, "debug", "d", false, "Print debug info")
	flag.StringVar(&o.dtb, "dtb", "", "FILE used as the flatten device tree blob")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.StringVarP(&o.extra, "extra", "x", "", "Add a cpio containing extra files")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.StringVar(&o.initramfs, "initramfs", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.loadSyscall, "loadsyscall", "L", false, "Use the kexec_load syscall (not kexec_file_load)")
	flag.StringArrayVar(&o.modules, "module", nil, `Load multiboot module with command line args (e.g --module="mod arg1")`)

	// This is broken out as it is almost never to be used. But it is valueable, nonetheless.
	flag.StringVarP(&o.purgatory, "purgatory", "p", "default", "picks a purgatory only if loading a Linux kernel with kexec_load, use '-p xyz' to get a list")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	return o
}

func main() {
	opts := registerFlags()
	flag.Parse()

	if opts.debug {
		linux.Debug = log.Printf
		purgatory.Debug = log.Printf
	}

	var kernelpath string
	if flag.NArg() > 0 {
		kernelpath = flag.Arg(0)
	}

	if (!opts.exec && len(kernelpath) == 0) || flag.NArg() > 1 {
		flag.PrintDefaults()
		log.Fatalf("usage: kexec [flags] kernelname OR kexec -e")
	}

	if opts.cmdline != "" && opts.reuseCmdline {
		flag.PrintDefaults()
		log.Fatalf("--reuse-cmdline and other command line options are mutually exclusive")
	}

	if !opts.load && !opts.exec {
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

	if err := purgatory.Select(opts.purgatory); err != nil {
		log.Fatal(err)
	}
	if opts.load {
		kernel, err := os.Open(kernelpath)
		if err != nil {
			log.Fatal(err)
		}
		defer kernel.Close()
		var image boot.OSImage
		if err := multiboot.Probe(kernel); err == nil {
			image = &boot.MultibootImage{
				Modules: multiboot.LazyOpenModules(opts.modules),
				Kernel:  kernel,
				Cmdline: newCmdline,
			}
		} else {
			var files []io.ReaderAt
			if len(opts.extra) > 0 {
				initrd, err := linux.CreateInitrd(strings.Fields(opts.extra)...)
				if err != nil {
					log.Fatal(err)
				}
				files = append(files, initrd)
			}
			if opts.initramfs != "" {
				for _, n := range strings.Fields(opts.initramfs) {
					files = append(files, uio.NewLazyFile(n))
				}
			}
			var i io.ReaderAt
			if len(files) > 0 {
				i = linux.CatInitrds(files...)
			}

			var dtb io.ReaderAt
			if len(opts.dtb) > 0 {
				dtb, err = os.Open(opts.dtb)
				if err != nil {
					log.Fatalf("Failed to open dtb file %s: %v", opts.dtb, err)
				}
			}
			image = &linux.Image{
				Kernel:      uio.NewLazyFile(kernelpath),
				Initrd:      i,
				Cmdline:     newCmdline,
				LoadSyscall: opts.loadSyscall,
				DTB:         dtb,
			}
		}
		if err := image.Load(boot.WithVerbose(opts.debug)); err != nil {
			log.Fatal(err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
