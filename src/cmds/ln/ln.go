// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	flags struct {
		symlink  bool
		verbose  bool
		force    bool
		nondir   bool
		prompt   bool
		logical  bool
		physical bool
		relative bool
		version  bool
		dirtgt   string
	}

	input = bufio.NewScanner(os.Stdin)
)

func init() {
	flag.BoolVar(&flags.symlink, "s", false, "make symbolic links instead of hard links")
	flag.BoolVar(&flags.verbose, "v", false, "print name of each linked file")
	flag.BoolVar(&flags.force, "f", false, "remove destination files")
	flag.BoolVar(&flags.nondir, "T", false, "treat linkname operand as a non-dir always")
	flag.BoolVar(&flags.prompt, "i", false, "prompt if the user wants overwrite")
	flag.BoolVar(&flags.logical, "L", false, "dereference targets if are symbolic links")
	flag.BoolVar(&flags.physical, "P", false, "make hard links directly to symbolic links")
	flag.BoolVar(&flags.relative, "r", false, "create symlinks relative to link location")
	flag.StringVar(&flags.dirtgt, "t", "", "specify the directory to put the links")
	flag.BoolVar(&flags.version, "version", false, "output version information and exit")
}

// simple version print
func version() {
	fmt.Fprintf(os.Stderr, "ln u-root version 0.1 by Manoel Vilela\n")
	os.Exit(0)
}

// ask for overwrite destination
func promptOverwrite(fname string) bool {
	fmt.Printf("ln: overwrite '%v'? ", fname)
	if input.Scan(); input.Text()[0] != 'y' {
		return false
	}
	return true
}

// the IsExists don't works fine, sorry for messy
func exists(fname string) bool {
	_, err := os.Lstat(fname)
	return !os.IsNotExist(err)
}

// based of the four type of uses describe
// get the targets and linknames operands
// attention: linkName can be "", which latter will be inferred (see inferLinkName)
func evalArgs(args []string) (targets []string, linkName string) {
	if flags.dirtgt != "" || len(args) <= 1 {
		return args, ""
	}

	targets = args[:len(args)-1]
	lastArg := args[len(args)-1]

	if lf, err := os.Stat(lastArg); !flags.nondir && err == nil && lf.IsDir() {
		flags.dirtgt = lastArg
	} else {
		linkName = lastArg
	}

	return targets, linkName
}

// get the relative link path
// between a linkName operand and the target
func relLink(target, linkName string) (string, error) {
	linkPath := filepath.Dir(linkName)
	if absPath, err := filepath.Abs(linkPath); err != nil {
		return "", err
	} else {
		return filepath.Rel(absPath, target)
	}
}

// if linkName don't passed ("") get the fname of target
// e.g.: $ ln -s -v /usr/bin/cp
//         cp -> /usr/bin/cp
func inferLinkName(target, linkName string) string {
	if linkName == "" {
		linkName = filepath.Base(target)
	}
	return linkName
}

// if flags.logical or flags.physical follow the links
// if flag.physical create a hard link instead symbolink
func dereferTarget(target string, linkFunc *func(string, string) error) (string, error) {
	if flags.logical || flags.physical {
		if newTarget, err := filepath.EvalSymlinks(target); err != nil {
			return "", err
		} else if newTarget != target {
			target = newTarget
			if flags.physical {
				*linkFunc = os.Symlink
			}
		}
	}
	return target, nil
}

func ln(args []string) error {
	var remove bool

	linkFunc := os.Link
	if flags.symlink {
		linkFunc = os.Symlink
	}

	targets, linkName := evalArgs(args)
	for _, target := range targets {
		linkFunc := linkFunc // back-overwrite possibilty

		// dereference symlinks
		t, err := dereferTarget(target, &linkFunc)
		if err != nil {
			return err
		}
		target = t

		linkName := inferLinkName(target, linkName)
		if flags.dirtgt != "" {
			linkName = path.Join(flags.dirtgt, linkName)
		}

		if exists(linkName) {
			if flags.prompt && !flags.force {
				remove = promptOverwrite(linkName)
			}

			if flags.force || remove {
				if err := os.Remove(linkName); err != nil {
					return err
				}
			}
		}

		if flags.relative && !flags.symlink {
			return fmt.Errorf("cannot do -r without -s")
		}

		// make relative paths with symlinks
		if flags.relative {
			if relTarget, err := relLink(target, linkName); err != nil {
				return err
			} else {
				target = relTarget
			}

		}

		if err := linkFunc(target, linkName); err != nil {
			return err
		}

		if flags.verbose {
			fmt.Printf("%q -> %q\n", linkName, target)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if flags.version {
		version()
	} else if len(args) == 0 {
		log.Printf("missing file operand")
		flag.Usage()
	}

	if err := ln(args); err != nil {
		log.Fatalf("link creation failed: %v", err)
	}
}
