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
//      --append string             Append to the kernel command line
//  -c, --cmdline string            Append to the kernel command line
//  -d, --debug                     Print debug info (default true)
//      --dtb string                FILE used as the flatten device tree blob
//  -e, --exec                      Execute a currently loaded kernel
//  -x, --extra string              Add a cpio containing extra files
//      --initramfs string          Use file as the kernel's initial ramdisk
//  -i, --initrd string             Use file as the kernel's initial ramdisk
//  -l, --load                      Load the new kernel into the current kernel
//  -L, --loadsyscall               Use the kexec load syscall (not file_load) (default true)
//      --mmap-initrd               Mmap initrd file into virtual buffer, other than directly reading it
//      --mmap-kernel               Mmap kernel file into virtual buffer, other than directly reading it
//      --module stringArray        Load multiboot module with command line args (e.g --module="mod arg1")
//  -p, --purgatory string          pick a purgatory, use '-p xyz' to get a list (default "default")
//      --purgatory-serial string   Name of the console used for purgatory printing
//      --reuse-cmdline             Use the kernel command line from running system

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
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/boot/purgatory"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
)

type options struct {
	cmdline         string
	debug           bool
	dtb             string
	exec            bool
	extra           string
	initramfs       string
	load            bool
	loadSyscall     bool
	mmapInitrd      bool
	mmapKernel      bool
	modules         []string
	purgatory       string
	purgatorySerial string
	reuseCmdline    bool
}

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Append to the kernel command line")
	flag.StringVar(&o.cmdline, "append", "", "Append to the kernel command line")
	flag.BoolVarP(&o.debug, "debug", "d", false, "Print debug info")
	flag.StringVar(&o.dtb, "dtb", "", "FILE used as the flatten device tree blob")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.StringVarP(&o.extra, "extra", "x", "", "Add a cpio containing extra files")
	flag.StringVar(&o.initramfs, "initramfs", "", "Use file as the kernel's initial ramdisk")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.loadSyscall, "loadsyscall", "L", false, "Use the kexec_load syscall (not kexec_file_load)")
	flag.BoolVar(&o.mmapInitrd, "mmap-initrd", true, "Mmap initrd file into virtual buffer, other than directly reading it (Only supported in Arm64 classic load mode for now)")
	flag.BoolVar(&o.mmapKernel, "mmap-kernel", true, "Mmap kernel file into virtual buffer, other than directly reading it (Only supported in Arm64 classi load mode for now)")
	flag.StringArrayVar(&o.modules, "module", nil, `Load multiboot module with command line args (e.g --module="mod arg1")`)
	// This "purgatory" flag is broken out as it is almost never to be used. But it is valueable, nonetheless.
	flag.StringVarP(&o.purgatory, "purgatory", "p", "default", "picks a purgatory only if loading a Linux kernel with kexec_load, use '-p xyz' to get a list")
	flag.StringVarP(&o.purgatorySerial, "purgatory-serial", "s", "", "Name of the console used for purgatory printing")
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

	if err := purgatory.Select(opts.purgatory); err != nil {
		log.Fatal(err)
	}
	if opts.load {
		kernelpath := flag.Arg(0)
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
				Kernel:      uio.NewLazyFile(kernelpath),
				Initrd:      i,
				Cmdline:     newCmdline,
				LoadSyscall: opts.loadSyscall,
				KexecOpts: linux.KexecOptions{
					Serial:     opts.purgatorySerial,
					DTB:        opts.dtb,
					MmapKernel: opts.mmapKernel,
					MmapRamfs:  opts.mmapInitrd,
				},
			}
		}
		if err := image.Load(opts.debug); err != nil {
			log.Fatal(err)
		}
		log.Printf("DH: image load finised w/o error")
	}

	log.Printf("DH: start checking exec")
	if opts.exec {
		log.Printf("DH: do exec")
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
