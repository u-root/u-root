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
//     -e or --exec:		      Execute a currently loaded kernel
package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/kexec"
	"github.com/u-root/u-root/pkg/multiboot"
)

type options struct {
	cmdline      string
	reuseCmdline bool
	initramfs    string
	load         bool
	exec         bool
	modules      []string
	trampoline   string
}

var trampoline = fmt.Sprintf(`Path to trampoline (will be removed in future releases).
Trampoline is a executable blob, which should set machine
to a specific state defined by multiboot v1 spec.
https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.

Trampoline should use a long word value stored right after %q byte sequence
as a value to be stored in ebx register and use a quad word value stored right after
%q as kernel entry point.`, multiboot.TrampolineEBX, multiboot.TrampolineEP)

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Set the kernel command line")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.StringSliceVar(&o.modules, "module", nil, `Load module with command line args (e.g --module="mod arg1")`)
	flag.StringVar(&o.trampoline, "trampoline", "", trampoline)
	return o
}

type loader interface {
	Load(path, cmdLine string) error
}

type file struct {
	initramfs string
}

type mboot struct {
	trampoline string
	modules    []string
}

func (f file) Load(path, cmdLine string) error {
	kernel, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open(%q): %v", path, err)
	}
	defer kernel.Close()

	var ramfs *os.File
	if f.initramfs != "" {
		ramfs, err = os.OpenFile(f.initramfs, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("open(%q): %v", f.initramfs, err)
		}
		defer ramfs.Close()
	}

	return kexec.FileLoad(kernel, ramfs, cmdLine)
}

func (mb mboot) Load(path, cmdLine string) error {
	m := multiboot.New(path, cmdLine, mb.trampoline, mb.modules)
	if err := m.Load(); err != nil {
		return fmt.Errorf("Load failed: %v", err)
	}

	err := kexec.Load(m.EntryPoint, m.Segments(), 0)
	if err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}
	return nil
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

		var l loader = file{opts.initramfs}
		if err := multiboot.Probe(kernelpath); err == nil {
			log.Printf("%s is a multiboot v1 kernel.", kernelpath)
			l = mboot{
				modules:    opts.modules,
				trampoline: opts.trampoline,
			}
		} else if err == multiboot.ErrFlagsNotSupported {
			log.Fatal(err)
		}

		if err := l.Load(kernelpath, newCmdLine); err != nil {
			log.Fatal(err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
