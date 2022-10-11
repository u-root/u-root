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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/cp"
)

type flags struct {
	recursive        bool
	ask              bool
	force            bool
	verbose          bool
	noFollowSymlinks bool
}

var (
	f flags
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "cp [-wRrifvP] file[s] ... dest"
		defUsage()
	}
	flag.BoolVarP(&f.recursive, "RECURSIVE", "R", false, "copy file hierarchies")
	flag.BoolVarP(&f.recursive, "recursive", "r", false, "alias to -R recursive mode")
	flag.BoolVarP(&f.ask, "interactive", "i", false, "prompt about overwriting file")
	flag.BoolVarP(&f.force, "force", "f", false, "force overwrite files")
	flag.BoolVarP(&f.verbose, "verbose", "v", false, "verbose copy mode")
	flag.BoolVarP(&f.noFollowSymlinks, "no-dereference", "P", false, "don't follow symlinks")
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

// run evaluates the args and makes decisions for copyfiles
func run(args []string, f flags, w io.Writer, i *bufio.Reader) error {
	todir := false
	from, to := args[:len(args)-1], args[len(args)-1]
	toStat, err := os.Stat(to)
	if err == nil {
		todir = toStat.IsDir()
	}
	if len(args) > 2 && !todir {
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
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(flag.Args(), f, os.Stderr, bufio.NewReader(os.Stdin)); err != nil {
		log.Fatalf("%q", err)
	}
}
