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
	symlink  = flag.Bool("s", false, "make symbolic links instead of hard links")
	verbose  = flag.Bool("v", false, "print name of each linked file")
	force    = flag.Bool("f", false, "remove destination files")
	nondir   = flag.Bool("T", false, "treat linkname operand as a non-dir always")
	prompt   = flag.Bool("i", false, "prompt if user wants overwrite")
	logical  = flag.Bool("L", false, "dereference targets if are symbolic links")
	physical = flag.Bool("P", false, "make hard links directly to symbolic links")
	relative = flag.Bool("r", false, "create symlinks relative to link location")
	dirtgt   = flag.String("t", "", "specify the directory to put the links")
	version  = flag.Bool("version", false, "output version information and exit")
	input    = bufio.NewScanner(os.Stdin)
)

func printVersion() {
	fmt.Fprintf(os.Stderr, "ln u-root version 0.1 by Manoel Vilela\n")
	os.Exit(0)
}

// ask for overwrite destination
func promptOverwrite(fname string) bool {
	fmt.Printf("ln: overwrite '%v'? ", fname)
	input.Scan()
	if input.Text()[0] != 'y' {
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
	if *dirtgt == "" && len(args) > 1 {
		targets = args[:len(args)-1]
		lastArg := args[len(args)-1]

		if lf, err := os.Stat(lastArg); !*nondir && err == nil && lf.IsDir() {
			*dirtgt = lastArg
		} else {
			linkName = lastArg
		}
	} else {
		targets = args
	}

	return
}

// get the relative link path
// between a newlink operand and the target
func relLink(target, linkName string) (string, error) {
	linkPath, _ := path.Split(linkName)
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
	newLink := linkName
	if newLink == "" {
		_, newLink = path.Split(target)
	}
	return newLink
}

func ln(args []string) error {
	var (
		linkFunc func(target, linkName string) error
		remove   bool
	)
	if *symlink {
		linkFunc = os.Symlink
	} else {
		linkFunc = os.Link
	}

	targets, linkName := evalArgs(args)
	for _, target := range targets {
		// dereference symlinks
		if *logical || *physical {
			if newTarget, err := filepath.EvalSymlinks(target); err != nil {
				return err
			} else if newTarget != target {
				target = newTarget
				if *physical { // if is a symlink, create other symlink
					// we have a problem on this
					// linkFunc will be used on other targets at iteration
					// that maybe not be a symlink
					linkFunc = os.Symlink
				}
			}
		}

		newLink := inferLinkName(target, linkName)
		if *dirtgt != "" {
			newLink = path.Join(*dirtgt, newLink)
		}

		if exists(newLink) {
			if *prompt && !*force {
				remove = promptOverwrite(newLink)
			}

			if *force || remove {
				if err := os.Remove(newLink); err != nil {
					return err
				}
			}
		}

		// make relative paths with symlinks
		if *relative && *symlink {
			if relTarget, err := relLink(target, linkName); err != nil {
				return err
			} else {
				target = relTarget
			}
		} else if *relative {
			return fmt.Errorf("cannot do -r without -s")
		}

		if err := linkFunc(target, newLink); err != nil {
			return err
		}

		if *verbose {
			fmt.Printf("'%v' -> '%v'\n", newLink, target)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if *version {
		printVersion()
	} else if len(args) < 1 {
		log.Printf("missing file operand")
		flag.Usage()
	}

	if err := ln(args); err != nil {
		log.Fatalf("link creation failed: %v", err)
	}
}
