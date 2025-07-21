// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// kexec executes a new kernel over the running kernel (u-root).
//
// Synopsis:
//     kexec [--initrd=FILE] [--command-line=STRING] [-l] [-e] [KERNELIMAGE]
//
// Description:
//		 Loads a kernel for later execution.
//
// Options:
//      --append string        Append to the kernel command line. Implies --reuse-cmdline
//  -c, --cmdline string       Set the kernel command line
//  -d, --debug                Print debug info (default true)
//  -e, --exec                 Execute a currently loaded kernel
//  -x, --extra string         Add one or more files to the initrd
//      --initramfs string     Use file as the kernel's initial ramdisk
//  -i, --initrd string        Use file as the kernel's initial ramdisk
//  -l, --load                 Load the new kernel into the current kernel
//  -L, --loadsyscall          Use the kexec load syscall (not file_load) (default true)
//      --module stringArray   Load multiboot module with command line args (e.g --module="mod arg1")
//  -p, --purgatory string     pick a purgatory, use '-p xyz' to get a list (default "default")
//      --reuse-cmdline        Use the kernel command line from running system

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/boot/linux"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/boot/purgatory"
	"github.com/u-root/u-root/pkg/boot/universalpayload"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
	"github.com/u-root/uio/uio"
)

type options struct {
	kernelpath string

	// Flags
	cmdline       string
	appendCmdline string
	debug         bool
	dtb           string
	exec          bool
	extra         string
	initramfs     string
	load          bool
	loadSyscall   bool
	modules       []string
	purgatory     string
	reuseCmdline  bool
}

func (o *options) parseCmdline(args []string, f *flag.FlagSet) {
	f.StringVar(&o.cmdline, "cmdline", "", "Set the kernel command line")
	f.StringVar(&o.cmdline, "c", "", "Set the kernel command line")

	f.StringVar(&o.appendCmdline, "append", "", "Append to the kernel command line. Implies --reuse-cmdline")

	f.BoolVar(&o.debug, "debug", false, "Print debug info")
	f.BoolVar(&o.debug, "d", false, "Print debug info")

	f.StringVar(&o.dtb, "dtb", "", "FILE used as the flatten device tree blob")

	f.BoolVar(&o.exec, "exec", false, "Execute a currently loaded kernel")
	f.BoolVar(&o.exec, "e", false, "Execute a currently loaded kernel")

	f.StringVar(&o.extra, "extra", "", "Add one or more files to the initrd")
	f.StringVar(&o.extra, "x", "", "Add one or more files to the initrd")

	f.StringVar(&o.initramfs, "initrd", "", "Use file as the kernel's initial ramdisk")
	f.StringVar(&o.initramfs, "i", "", "Use file as the kernel's initial ramdisk")

	f.StringVar(&o.initramfs, "initramfs", "", "Use file as the kernel's initial ramdisk")

	// Although -l or--load is actually a switch, traditional kexec command allows the kernel
	// to be passed as the value right after the l flag. This is why we have to handle this
	// as a special case and define the flag as a string. Together with some checks below,
	// we can handle both:
	// - kexec -l [other flags] /path/to/kernel
	// - kexec -l /path/to/kernel [other flags]
	var loadFlagPath string
	f.StringVar(&loadFlagPath, "load", "", "Load the new kernel into the current kernel")
	f.StringVar(&loadFlagPath, "l", "", "Load the new kernel into the current kernel (shorthand)")

	f.BoolVar(&o.loadSyscall, "loadsyscall", false, "Use the kexec_load syscall (not kexec_file_load)")
	f.BoolVar(&o.loadSyscall, "L", false, "Use the kexec_load syscall (not kexec_file_load) (shorthand)")

	f.Var((*unixflag.StringArray)(&o.modules), "module", `Load multiboot module with command line args (e.g --module="mod arg1")`)

	// This is broken out as it is almost never to be used. But it is valueable, nonetheless.
	f.StringVar(&o.purgatory, "purgatory", "default", "picks a purgatory only if loading a Linux kernel with kexec_load, use '-p xyz' to get a list")
	f.StringVar(&o.purgatory, "p", "default", "picks a purgatory only if loading a Linux kernel with kexec_load, use '-p xyz' to get a list (shorthand)")

	f.BoolVar(&o.reuseCmdline, "reuse-cmdline", false, "Use the kernel command line from running system")

	unixargs := unixflag.ArgsToGoArgs(args[1:])
	hackedArgs := hackLoadFlagValue(unixargs)

	f.Parse(hackedArgs)

	if loadFlagPath != "" {
		o.load = true
		if loadFlagPath == setButEmpty {
			loadFlagPath = ""
		}
	}

	// Allow the kernel argument to appear eitheras the value of the -l flag
	// or as an argument at the end of the command line.
	if f.NArg() > 0 && loadFlagPath == "" {
		o.kernelpath = f.Arg(0)
	} else if loadFlagPath != "" {
		o.kernelpath = loadFlagPath
	}
}

const setButEmpty = "eMptY"

// hackLoadFlagValue makes sure that the value of the -l flag has a string representation for being not set.
// Otherwise, in case of e.g. "kexec -l -d /path/to/kernel" the -d flag would be interpreted as the value of the -l flag.
func hackLoadFlagValue(in []string) []string {
	if len(in) == 0 {
		return in
	}

	var out []string
	for n := 0; n <= len(in)-2; n++ {
		current := n
		// next := n + 1
		if (in[current] == "-l" || in[current] == "--load") && strings.HasPrefix(in[current+1], "-") {
			out = append(out, in[current], setButEmpty)
		} else {
			out = append(out, in[current])
		}
	}

	out = append(out, in[len(in)-1])
	return out
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("%v", err)
	}
}

func run(args []string) error {
	opts := &options{}
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	opts.parseCmdline(args, f)

	if opts.debug {
		linux.Debug = log.Printf
		purgatory.Debug = log.Printf
	}

	if (!opts.exec && len(opts.kernelpath) == 0) || f.NArg() > 1 {
		f.PrintDefaults()
		return fmt.Errorf("usage: kexec [fs] kernelname OR kexec -e")
	}

	if err, warningMsg := universalpayload.Load(opts.kernelpath, linux.Debug); err != nil {
		log.Printf("Failed to load universalpayload (%v), try legacy kernel..", err)
	} else {
		// universalpayload package suppresses warning message, we print messages here.
		if warningMsg != nil {
			log.Printf("Warning messages from universalpayload:\n%v\n", warningMsg)
		}

		if err := universalpayload.Exec(); err != nil {
			log.Printf("Failed to execute universalpayload (%v), try legacy kernel..", err)
		}
	}

	if opts.cmdline != "" && opts.reuseCmdline {
		f.PrintDefaults()
		return fmt.Errorf("--reuse-cmdline and other command line options are mutually exclusive")
	}

	if !opts.load && !opts.exec {
		opts.load = true
		opts.exec = true
	}

	newCmdline := opts.cmdline
	if opts.reuseCmdline || len(opts.appendCmdline) != 0 {
		procCmdLine := cmdline.NewCmdLine()
		if procCmdLine.Err != nil {
			return fmt.Errorf("couldn't read /proc/cmdline")
		}
		newCmdline = procCmdLine.Raw + " " + opts.appendCmdline

	}

	if err := purgatory.Select(opts.purgatory); err != nil {
		return err
	}
	if opts.load {
		kernel, err := os.Open(opts.kernelpath)
		if err != nil {
			return err
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
				initrd, err := boot.CreateInitrd(strings.Fields(opts.extra)...)
				if err != nil {
					return err
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
				i = boot.CatInitrds(files...)
			}

			var dtb io.ReaderAt
			if len(opts.dtb) > 0 {
				dtb, err = os.Open(opts.dtb)
				if err != nil {
					return fmt.Errorf("failed to open dtb file %s: %w", opts.dtb, err)
				}
			}
			image = &boot.LinuxImage{
				Kernel:      uio.NewLazyFile(opts.kernelpath),
				Initrd:      i,
				Cmdline:     newCmdline,
				LoadSyscall: opts.loadSyscall,
				DTB:         dtb,
			}
		}
		if err := image.Load(boot.WithVerbose(opts.debug)); err != nil {
			return err
		}
	}

	if opts.exec {
		if err := kexec.Reboot(); err != nil {
			return err
		}
	}

	return nil
}
