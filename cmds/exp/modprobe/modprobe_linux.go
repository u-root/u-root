// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modprobe - Add and remove modules from the Linux Kernel.
//
// Synopsis:
//
//	modprobe [-nqvrb] [-d DIR] [-S VERSION] [--] modulename [parameters...]
//	modprobe [-nqvrb] [-d DIR] [-S VERSION] -a modulename...
//	modprobe [-nqv] [-d DIR] [-S VERSION] -r modulename...
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/u-root/u-root/pkg/kmodule"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errUsage = errors.New("one module and optional module options required")

type flags struct {
	dryRun        bool
	all           bool
	verbose       bool
	quiet         bool
	remove        bool
	useBlacklist  bool
	ignoreInstall bool
	rootDir       string
	kernelVer     string
}

func run(args []string, stderr io.Writer) error {
	var f flags
	fs := flag.NewFlagSet("modprobe", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&f.dryRun, "n", false, "Dry run")
	fs.BoolVar(&f.dryRun, "dry-run", false, "Dry run")
	fs.BoolVar(&f.dryRun, "show", false, "Dry run")
	fs.BoolVar(&f.all, "a", false, "Insert all module names on the command line")
	fs.BoolVar(&f.all, "all", false, "Insert all module names on the command line")
	fs.BoolVar(&f.verbose, "v", false, "Print messages about what the program is doing")
	fs.BoolVar(&f.verbose, "verbose", false, "Print messages about what the program is doing")
	fs.BoolVar(&f.quiet, "q", false, "Disable error messages")
	fs.BoolVar(&f.quiet, "quiet", false, "Disable error messages")
	fs.BoolVar(&f.remove, "r", false, "Remove modules instead of inserting")
	fs.BoolVar(&f.remove, "remove", false, "Remove modules instead of inserting")
	fs.BoolVar(&f.useBlacklist, "b", false, "Apply blacklist to module names too")
	fs.BoolVar(&f.useBlacklist, "use-blacklist", false, "Apply blacklist to module names too")
	fs.BoolVar(&f.ignoreInstall, "i", false, "Ignore install and remove commands")
	fs.BoolVar(&f.ignoreInstall, "ignore-install", false, "Ignore install commands")
	fs.BoolVar(&f.ignoreInstall, "ignore-remove", false, "Ignore remove commands")
	fs.StringVar(&f.rootDir, "d", "/", "Root directory for modules")
	fs.StringVar(&f.rootDir, "dirname", "/", "Root directory for modules")
	fs.StringVar(&f.kernelVer, "S", "", "Set kernel version instead of using uname")
	fs.StringVar(&f.kernelVer, "set-version", "", "Set kernel version instead of using uname")

	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: modprobe [-anqvrb] [-d DIR] [-S VERSION] [--] modulename[s] [parameters...]")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if f.quiet {
		stderr = io.Discard
		fs.SetOutput(stderr)
	}
	log.SetOutput(stderr)

	if fs.NArg() == 0 {
		fmt.Fprintln(stderr, "Usage: ERROR: one module and optional module options.")
		fs.Usage()
		return errUsage
	}

	opts := kmodule.ProbeOpts{
		RootDir:      f.rootDir,
		KVer:         f.kernelVer,
		UseBlacklist: f.useBlacklist,
	}
	if f.dryRun {
		if f.verbose {
			fmt.Fprintln(stderr, "Unique dependencies in order, already handled ones get skipped:")
		}
		var dryRunMu sync.Mutex
		opts.DryRunCB = func(modPath string) {
			dryRunMu.Lock()
			defer dryRunMu.Unlock()
			fmt.Fprintln(stderr, modPath)
		}
	}

	prober, err := kmodule.NewProber(opts)
	if err != nil {
		return err
	}

	if f.remove {
		var retErr error
		for _, modName := range fs.Args() {
			if err := prober.Remove(modName); err != nil {
				retErr = errors.Join(retErr, fmt.Errorf("could not remove module %q: %w", modName, err))
			}
		}
		return retErr
	}

	if f.all {
		var retErr error
		for _, modName := range fs.Args() {
			if err := prober.Probe(modName, ""); err != nil {
				retErr = errors.Join(retErr, fmt.Errorf("could not load module %q: %w", modName, err))
			}
		}
		return retErr
	}

	modName := fs.Args()[0]
	modOptions := strings.Join(fs.Args()[1:], " ")

	if err := prober.Probe(modName, modOptions); err != nil {
		return fmt.Errorf("could not load module %q: %w", modName, err)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		log.Fatal(err)
	}
}
