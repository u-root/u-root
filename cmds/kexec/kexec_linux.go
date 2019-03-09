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
	"fmt"
	"io/ioutil"
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
	return o
}

type loader interface {
	Load(path, cmdLine string) error
}

type file struct {
	initramfs string
}

type memory struct {
	initramfs string
	s         []kexec.Segment
}

type mboot struct {
	debug   bool
	modules []string
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

func (me memory) Load(path, cmdLine string) error {
	var m kexec.Memory
	if err := m.ParseMemoryMap(); err != nil {
		return err
	}
	kernel, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open(%q): %v", path, err)
	}
	defer kernel.Close()
	if err := m.LoadElfSegments(kernel); err != nil {
		return err
	}
	if false {
		var ramfs *os.File
		if me.initramfs != "" {
			ramfs, err = os.OpenFile(me.initramfs, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("open(%q): %v", me.initramfs, err)
			}
			defer ramfs.Close()
		}
	}
	return kexec.Load(0x100000, m.Segments, 0)
}

func (mb mboot) Load(path, cmdLine string) error {
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
	segs := m.Segments()
	if err := kexec.Load(m.EntryPoint, segs, 0); err != nil {
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

	if (opts.initramfs != "" || opts.acpi != "") && !opts.load {
		log.Fatal("If you want an initramfs you must set it up at load time")
	}

	if opts.acpi != "" {
		if false {
			var m kexec.Memory
			if err := m.ParseMemoryMap(); err != nil {
				log.Fatal(err)
			}
			log.Printf("acpi load: m is %s", m.String())
			// it's extremely unlikely that we are replacing all acpi tables.
			// For now, assume we are appending.
			b, err := acpi.TablesData()
			if err != nil {
				log.Fatal(err)
			}

			addb, err := ioutil.ReadFile(opts.acpi)
			if err != nil {
				log.Fatal(err)
			}
			b = append(b, addb...)
			addr, err := m.AddKexecSegmentACPI(b)
			if err != nil {
				log.Fatalf("Allocating segment for ACPI: %v", err)
			}

			log.Printf("ACPI loaded: addr is %#x, m is %s", addr, m)
		}
		var l loader = memory{}
		kernelpath := flag.Args()[0]
		if err := l.Load(kernelpath, newCmdLine); err != nil {
			log.Fatal(err)

		}
	}

	if opts.load {
		kernelpath := flag.Args()[0]
		log.Printf("Loading %s for kernel\n", kernelpath)

		var l loader = file{initramfs: opts.initramfs}
		if err := multiboot.Probe(kernelpath); err == nil {
			log.Printf("%s is a multiboot v1 kernel.", kernelpath)
			l = mboot{
				debug:   opts.debug,
				modules: opts.modules,
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
