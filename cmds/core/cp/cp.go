// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cp copies files.
//
// Synopsis:
//     cp [-rRfivwP] FROM... TO
//
// Options:
//     -w n: number of worker goroutines
//     -R: copy file hierarchies
//     -r: alias to -R recursive mode
//     -i: prompt about overwriting file
//     -f: force overwrite files
//     -v: verbose copy mode
//     -P: don't follow symlinks
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/cp"
)

var (
	flags struct {
		recursive        bool
		ask              bool
		force            bool
		verbose          bool
		noFollowSymlinks bool
	}
	input = bufio.NewReader(os.Stdin)
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "cp [-wRrifvP] file[s] ... dest"
		defUsage()
	}
	flag.BoolVarP(&flags.recursive, "RECURSIVE", "R", false, "copy file hierarchies")
	flag.BoolVarP(&flags.recursive, "recursive", "r", false, "alias to -R recursive mode")
	flag.BoolVarP(&flags.ask, "interactive", "i", false, "prompt about overwriting file")
	flag.BoolVarP(&flags.force, "force", "f", false, "force overwrite files")
	flag.BoolVarP(&flags.verbose, "verbose", "v", false, "verbose copy mode")
	flag.BoolVarP(&flags.noFollowSymlinks, "no-dereference", "P", false, "don't follow symlinks")
}

// promptOverwrite ask if the user wants overwrite file
func promptOverwrite(dst string) (bool, error) {
	fmt.Printf("cp: overwrite %q? ", dst)
	answer, err := input.ReadString('\n')
	if err != nil {
		return false, err
	}

	if strings.ToLower(answer)[0] != 'y' {
		return false, nil
	}

	return true, nil
}

// cpArgs is a function whose eval the args
// and make decisions for copyfiles
func cpArgs(args []string) error {
	todir := false
	from, to := args[:len(args)-1], args[len(args)-1]
	toStat, err := os.Stat(to)
	if err == nil {
		todir = toStat.IsDir()
	}
	if flag.NArg() > 2 && !todir {
		log.Fatalf("is not a directory: %s\n", to)
	}

	opts := cp.Options{
		NoFollowSymlinks: flags.noFollowSymlinks,

		// cp the command makes sure that
		//
		// (1) the files it's copying aren't already the same,
		// (2) the user is asked about overwriting an existing file if
		//     one is already there.
		PreCallback: func(src, dst string, srcfi os.FileInfo) error {
			// check if src is dir
			if !flags.recursive && srcfi.IsDir() {
				log.Printf("cp: -r not specified, omitting directory %s", src)
				return cp.ErrSkip
			}

			dstfi, err := os.Stat(dst)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("cp: %q: can't handle error %v", dst, err)
				return cp.ErrSkip
			} else if err != nil {
				// dst does not exist.
				return nil
			}

			// dst does exist.

			if os.SameFile(srcfi, dstfi) {
				log.Printf("cp: %q and %q are the same file", src, dst)
				return cp.ErrSkip
			}
			if flags.ask && !flags.force {
				overwrite, err := promptOverwrite(dst)
				if err != nil {
					return err
				}
				if !overwrite {
					return cp.ErrSkip
				}
			}
			return nil
		},

		PostCallback: func(src, dst string) {
			if flags.verbose {
				fmt.Printf("%q -> %q\n", src, dst)
			}
		},
	}

	var lastErr error
	for _, file := range from {
		dst := to
		if todir {
			dst = filepath.Join(dst, filepath.Base(file))
		}
		if flags.recursive {
			err = opts.CopyTree(file, dst)
		} else {
			err = opts.Copy(file, dst)
		}
		if err != nil {
			log.Printf("cp: %v\n", err)
			lastErr = err
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

	if err := cpArgs(flag.Args()); err != nil {
		os.Exit(1)
	}
}
