// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ln makes links to files.
//
// Synopsis:
//
//	ln [-svfTiLPrt] TARGET LINK
//
// Options:
//
//	-s: make symbolic links instead of hard links
//	-v: print name of each linked file
//	-f: remove destination files
//	-T: treat linkname operand as a non-dir always
//	-i: prompt if the user wants overwrite
//	-L: dereference targets if are symbolic links
//	-P: make hard links directly to symbolic links
//	-r: create symlinks relative to link location
//	-t: specify the directory to put the links
//
// Author:
//
//	Manoel Vilela <manoel_vilela@engineer.com>
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	symlink  bool
	verbose  bool
	force    bool
	nondir   bool
	prompt   bool
	logical  bool
	physical bool
	relative bool
	dirtgt   string
}

// promptOverwrite ask for overwrite destination
func promptOverwrite(fname string) bool {
	fmt.Printf("ln: overwrite '%v'? ", fname)
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return err == nil && strings.ToLower(answer)[0] == 'y'
}

// exists verify if a fname exists
// the IsExists don't works fine, sorry for messy
func exists(fname string) bool {
	_, err := os.Lstat(fname)
	return !os.IsNotExist(err)
}

// evalArgs based of the four type of uses describe at ln
// get the targets and linknames operands
// attention: linkName can be "", which latter will be inferred (see inferLinkName function)
func (conf *config) evalArgs(args []string) (targets []string, linkName string) {
	if conf.dirtgt != "" || len(args) <= 1 {
		return args, ""
	}

	targets = args[:len(args)-1]
	lastArg := args[len(args)-1]

	if lf, err := os.Stat(lastArg); !conf.nondir && err == nil && lf.IsDir() {
		conf.dirtgt = lastArg
	} else {
		linkName = lastArg
	}

	return targets, linkName
}

// relLink get the relative link path between
// a target and linkName fpath
// between a linkName operand and the target
// HMM, i don't have sure if that works well...
func relLink(target, linkName string) (string, error) {
	base := filepath.Dir(linkName)
	if newTarget, err := filepath.Rel(base, target); err == nil {
		return newTarget, nil
	} else if absLink, err := filepath.Abs(linkName); err == nil {
		return filepath.Rel(absLink, target)
	}

	return "", nil
}

// inferLinkname infers the linkName if don't passed ("")
// otherwhise preserves the linkName
// e.g.:
// $ ln -s -v /usr/bin/cp
// cp -> /usr/bin/cp
func inferLinkName(target, linkName string) string {
	if linkName == "" {
		linkName = filepath.Base(target)
	}
	return linkName
}

// dereferTarget treat symlinks according the flags -P and -L
// if conf.logical or conf.physical follow the links
// if conf.physical create a hard link instead symbolink
func (conf config) dereferTarget(target string, linkFunc *func(string, string) error) (string, error) {
	if conf.logical || conf.physical {
		if newTarget, err := filepath.EvalSymlinks(target); err != nil {
			return "", err
		} else if newTarget != target {
			target = newTarget
			if conf.physical {
				*linkFunc = os.Symlink
			}
		}
	}
	return target, nil
}

// ln is a general procedure for controlling the
// flow of links creation, handling the flags and other stuffs.
func (conf config) ln(args []string) error {
	var remove bool

	linkFunc := os.Link
	if conf.symlink {
		linkFunc = os.Symlink
	}

	originalPath, err := os.Getwd()
	if err != nil {
		return err
	}

	targets, linkName := conf.evalArgs(args)
	for _, target := range targets {
		linkFunc := linkFunc // back-overwrite possibility

		// dereference symlinks
		t, err := conf.dereferTarget(target, &linkFunc)
		if err != nil {
			return err
		}
		target = t

		linkName := inferLinkName(target, linkName)
		if conf.dirtgt != "" {
			linkName = filepath.Join(conf.dirtgt, linkName)
		}

		if exists(linkName) {
			if conf.prompt && !conf.force {
				remove = promptOverwrite(linkName)
				if !remove {
					continue
				}
			}

			if conf.force || remove {
				if err := os.Remove(linkName); err != nil {
					return err
				}
			}
		}

		if conf.relative && !conf.symlink {
			return fmt.Errorf("cannot do -r without -s")
		}

		// make relative paths with symlinks
		if conf.relative {
			relTarget, err := relLink(target, linkName)
			if err != nil {
				return err
			}
			target = relTarget

			if dir := filepath.Dir(linkName); dir != "" {
				linkName = filepath.Base(linkName)
				if err := os.Chdir(dir); err != nil {
					return err
				}
			}

		}

		if err := linkFunc(target, linkName); err != nil {
			return err
		}

		if conf.relative {
			if err := os.Chdir(originalPath); err != nil {
				return err
			}
		}

		if conf.verbose {
			fmt.Printf("%q -> %q\n", linkName, target)
		}
	}

	return nil
}

func main() {
	var conf config
	flag.BoolVar(&conf.symlink, "s", false, "make symbolic links instead of hard links")
	flag.BoolVar(&conf.verbose, "v", false, "print name of each linked file")
	flag.BoolVar(&conf.force, "f", false, "remove destination files")
	flag.BoolVar(&conf.nondir, "T", false, "treat linkname operand as a non-dir always")
	flag.BoolVar(&conf.prompt, "i", false, "prompt if the user wants overwrite")
	flag.BoolVar(&conf.logical, "L", false, "dereference targets if are symbolic links")
	flag.BoolVar(&conf.physical, "P", false, "make hard links directly to symbolic links")
	flag.BoolVar(&conf.relative, "r", false, "create symlinks relative to link location")
	flag.StringVar(&conf.dirtgt, "t", "", "specify the directory to put the links")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Printf("ln: missing file operand")
		flag.Usage()
	}

	if err := conf.ln(args); err != nil {
		log.Fatalf("ln: link creation failed: %v", err)
	}
}
