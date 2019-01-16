// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// kexec executes a new kernel over the running kernel (u-root).
//
// Synopsis:
//     kexec [--initrd=FILE] [--command-line=STRING] [-l] [-e] [KERNELIMAGE]
//
// Description:
//	   Loads a kernel for later execution.
//
// Options:
//     --acpi=STRING or -a=string      Add an ACPI table (only one at present)
//     --cmdline=STRING or -c=STRING: Set the kernel command line
//     --reuse-commandline:           Use the kernel command line from running system
//     --i=FILE or --initrd=FILE:     Use file as the kernel's initial ramdisk
//     -l or --load:                  Load the new kernel into the current kernel
//     -e or --exec:                  Execute a currently loaded kernel
//     -d or --debug:                 Print debug info
//     --module:                      Load module with command line args
//     --dtb FILE:                    Override the device tree with this file
//     --dryrun:                      Print segments, do not load kernel
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/acpi"
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
	debug        bool
	acpi         string
	modules      []string
	dtb          string // device tree
	dryrun       bool
}

func registerFlags() *options {
	o := &options{}
	flag.StringVarP(&o.acpi, "acpi", "a", "", "Add an acpi table")
	flag.StringVarP(&o.cmdline, "cmdline", "c", "", "Set the kernel command line")
	flag.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")
	flag.StringVarP(&o.initramfs, "initrd", "i", "", "Use file as the kernel's initial ramdisk")
	flag.BoolVarP(&o.load, "load", "l", false, "Load the new kernel into the current kernel")
	flag.BoolVarP(&o.exec, "exec", "e", false, "Execute a currently loaded kernel")
	flag.BoolVarP(&o.debug, "debug", "d", false, "Print debug info")
	flag.StringSliceVar(&o.modules, "module", nil, `Load module with command line args (e.g --module="mod arg1")`)
	flag.StringVar(&o.dtb, "dtb", "", "Override the device tree with this file")
	flag.BoolVar(&o.dryrun, "dryrun", false, "Print segments, do not load kernel")
	return o
}

type loader interface {
	Load(path, cmdLine string, o *options) error
}

type file struct {
	initramfs string
}

type mboot struct {
	debug   bool
	modules []string
	acpi    string
}

func (f file) Load(path, cmdLine string, o *options) error {
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

	kOpts := &kexec.LinuxOpts{
		Initramfs: ramfs,
		CmdLine:   cmdLine,
	}
	if o.dryrun {
		kOpts.FileLoadSyscall = kexec.DryrunFileLoad
		kOpts.LoadSyscall = kexec.DryrunLoad
	}
	if o.dtb != "" {
		f, err := os.OpenFile(o.dtb, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("open(%q): %v", path, err)
		}
		defer f.Close()
		kOpts.DTB = f
	}
	return kexec.Load(kOpts)
}

func (mb mboot) Load(path, cmdLine string, o *options) error {
	// Trampoline should be a part of current binary.
	p, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Cannot find current executable path: %v", err)
	}
	trampoline, err := filepath.EvalSymlinks(p)
	if err != nil {
		return fmt.Errorf("Cannot eval symlinks for %v: %v", p, err)
	}
	m := multiboot.New(path, cmdLine, trampoline, mb.modules)
	if err := m.Load(mb.debug); err != nil {
		return fmt.Errorf("Load failed: %v", err)
	}
	if mb.acpi != "" {
		if err := m.ACPI(mb.acpi); err != nil {
			return err
		}
	}

	if err := kexec.RawLoad(m.EntryPoint, m.Segments(), 0); err != nil {
		return fmt.Errorf("kexec.RawLoad() error: %v", err)
	}
	return nil
}

func main() {
	opts := registerFlags()
	flag.Parse()

	if opts.debug {
		acpi.Debug = log.Printf
	}
	if (opts.exec == false && flag.NArg() == 0) || flag.NArg() > 1 {
		flag.PrintDefaults()
		log.Fatalf("usage: kexec [flags] kernelname OR kexec -e")
	}

	if opts.cmdline != "" && opts.reuseCmdline {
		flag.PrintDefaults()
		log.Fatalf("--reuse-cmdline and other command line options are mutually exclusive")
	}

	if opts.load == false && opts.exec == false && opts.acpi == "" {
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

	// mbk indicates that we are a multiboot kernel
	var mbk bool
	var kernelpath string
	if opts.load {
		kernelpath = flag.Arg(0)
		log.Printf("Loading %s for kernel\n", kernelpath)

		if err := multiboot.Probe(kernelpath); err == nil {
			mbk = true
		} else if err == multiboot.ErrFlagsNotSupported {
			log.Fatal(err)
		}
	}
	if opts.acpi != "" && !mbk {
		log.Fatal("You can only specify -a when loading (-l) multiboot kernels")
	}
	if opts.load {
		var l loader = file{initramfs: opts.initramfs}
		if mbk {
			log.Printf("%s is a multiboot v1 kernel.", kernelpath)
			l = mboot{
				debug:   opts.debug,
				modules: opts.modules,
				acpi:    opts.acpi,
			}
		}
		if err := l.Load(kernelpath, newCmdLine, opts); err != nil {
			log.Fatal(err)
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
