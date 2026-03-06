// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cp copies files.
//
// Synopsis:
//
//	cp [-rRfivwP] FROM... TO
//
// Options:
//
//	-w n: number of worker goroutines
//	-R: copy file hierarchies
//	-r: alias to -R recursive mode
//	-i: prompt about overwriting file
//	-f: force overwrite files
//	-v: verbose copy mode
//	-P: don't follow symlinks
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type flags struct {
	recursive        bool
	ask              bool
	force            bool
	verbose          bool
	noFollowSymlinks bool
}

// promptOverwrite ask if the user wants overwrite file
func promptOverwrite(dst string, out io.Writer, in *bufio.Reader) (bool, error) {
	fmt.Fprintf(out, "cp: overwrite %q? ", dst)
	answer, err := in.ReadString('\n')
	if err != nil {
		return false, err
	}

	if strings.ToLower(answer)[0] != 'y' {
		return false, nil
	}

	return true, nil
}

func setupPreCallback(recursive, ask, force bool, writer io.Writer, reader bufio.Reader) func(string, string, os.FileInfo) error {
	return func(src, dst string, srcfi os.FileInfo) error {
		// check if src is dir
		if !recursive && srcfi.IsDir() {
			fmt.Fprintf(writer, "cp: -r not specified, omitting directory %s\n", src)
			return cp.ErrSkip
		}

		dstfi, err := os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(writer, "cp: %q: can't handle error %v\n", dst, err)
			return cp.ErrSkip
		} else if err != nil {
			// dst does not exist.
			return nil
		}

		// dst does exist.

		if os.SameFile(srcfi, dstfi) {
			fmt.Fprintf(writer, "cp: %q and %q are the same file\n", src, dst)
			return cp.ErrSkip
		}
		if ask && !force {
			overwrite, err := promptOverwrite(dst, writer, &reader)
			if err != nil {
				return err
			}
			if !overwrite {
				return cp.ErrSkip
			}
		}
		return nil
	}
}

func setupPostCallback(verbose bool, w io.Writer) func(src, dst string) {
	return func(src, dst string) {
		if verbose {
			fmt.Fprintf(w, "%q -> %q\n", src, dst)
		}
	}
}

// run evaluates the falgs and args and makes decisions for copyfiles
func run(args []string, w io.Writer, i *bufio.Reader) error {
	var f flags

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.BoolVar(&f.recursive, "RECURSIVE", false, "copy file hierarchies")
	fs.BoolVar(&f.recursive, "R", false, "copy file hierarchies (shorthand)")

	fs.BoolVar(&f.recursive, "recursive", false, "alias to -R recursive mode")
	fs.BoolVar(&f.recursive, "r", false, "alias to -R recursive mode (shorthand)")

	fs.BoolVar(&f.ask, "interactive", false, "prompt about overwriting file")
	fs.BoolVar(&f.ask, "i", false, "prompt about overwriting file (shorthand)")

	fs.BoolVar(&f.force, "force", false, "force overwrite files")
	fs.BoolVar(&f.force, "f", false, "force overwrite files (shorthand)")

	fs.BoolVar(&f.verbose, "verbose", false, "verbose copy mode")
	fs.BoolVar(&f.verbose, "v", false, "verbose copy mode (shorthand)")

	fs.BoolVar(&f.noFollowSymlinks, "no-dereference", false, "don't follow symlinks")
	fs.BoolVar(&f.noFollowSymlinks, "P", false, "don't follow symlinks (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: cp [-RrifvP] file[s] ... dest\n\n")
		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	if fs.NArg() < 2 {
		fs.Usage()
		os.Exit(1)
	}

	todir := false
	from, to := fs.Args()[:fs.NArg()-1], fs.Args()[fs.NArg()-1]
	toStat, err := os.Stat(to)
	if err == nil {
		todir = toStat.IsDir()
	}
	if fs.NArg() > 2 && !todir {
		return eNotDir
	}

	opts := cp.Options{
		NoFollowSymlinks: f.noFollowSymlinks,

		// cp the command makes sure that
		//
		// (1) the files it's copying aren't already the same,
		// (2) the user is asked about overwriting an existing file if
		//     one is already there.
		PreCallback: setupPreCallback(f.recursive, f.ask, f.force, w, *i),

		PostCallback: setupPostCallback(f.verbose, w),
	}

	var lastErr error
	for _, file := range from {
		dst := to
		if todir {
			dst = filepath.Join(dst, filepath.Base(file))
		}
		if f.recursive {
			lastErr = opts.CopyTree(file, dst)
		} else {
			lastErr = opts.Copy(file, dst)
		}
	}
	return lastErr
}

func main() {
	err := run(os.Args, os.Stderr, bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("%q", err)
	}
}
