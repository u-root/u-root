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
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
)

type options struct {
	cmdline      string
	reuseCmdline bool
	initramfs    string
	load         bool
	exec         bool
	debug        bool
	extra        string
	modules      []string
}

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Append to the kernel command line")
	flag.StringVar(&o.cmdline, "append", "", "Append to the kernel command line")
	flag.StringVarP(&o.extra, "extra", "x", "", "Add a cpio containing extra files")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.StringVar(&o.initramfs, "initramfs", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.BoolVarP(&o.debug, "debug", "d", false, "Print debug info")
	flag.StringArrayVar(&o.modules, "module", nil, `Load multiboot module with command line args (e.g --module="mod arg1")`)
	return o
}

func main() {
	opts := registerFlags()
	flag.Parse()

	if (!opts.exec && flag.NArg() == 0) || flag.NArg() > 1 {
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

	if opts.load {
		kernelpath := flag.Arg(0)
		mbkernel, err := os.Open(kernelpath)
		if err != nil {
			log.Fatal(err)
		}
		defer mbkernel.Close()
		var image boot.OSImage
		if err := multiboot.Probe(mbkernel); err == nil {
			image = &boot.MultibootImage{
				Modules: multiboot.LazyOpenModules(opts.modules),
				Kernel:  mbkernel,
				Cmdline: newCmdline,
			}
		} else {
			var files []io.ReaderAt
			if len(opts.extra) > 0 {
				b := &bytes.Buffer{}
				archiver, err := cpio.Format("newc")
				if err != nil {
					log.Fatal(err)
				}
				w := archiver.Writer(b)
				cr := cpio.NewRecorder()
				// to deconflict names, we may want to prepend the names with
				// kexec_extra/ or something.
				for _, n := range strings.Fields(opts.extra) {
					rec, err := cr.GetRecord(n)
					if err != nil {
						log.Fatalf("Getting record of %q failed: %v", n, err)
					}
					if err := w.WriteRecord(rec); err != nil {
						log.Fatalf("Writing record %q failed: %v", n, err)
					}
				}
				if err := cpio.WriteTrailer(w); err != nil {
					log.Fatalf("Error writing trailer record: %v", err)
				}
				files = append(files, bytes.NewReader(b.Bytes()))
			}
			if opts.initramfs != "" {
				for _, n := range strings.Fields(opts.initramfs) {
					files = append(files, uio.NewLazyFile(n))
				}

			}
			var i io.ReaderAt
			if len(files) > 0 {
				i = boot.CatInitrds(files...)
			}
			image = &boot.LinuxImage{
				Kernel:  uio.NewLazyFile(kernelpath),
				Initrd:  i,
				Cmdline: newCmdline,
			}
		}
		if err := image.Load(opts.debug); err != nil {
			log.Fatal(err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
