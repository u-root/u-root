// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Find finds files. It is similar to the Unix command. It uses REs, not globs,
// for matching.
//
// OPTIONS:
//     -d: enable debugging in the find package
//     -mode integer-arg: match against mode, e.g. -mode 0755
//     -type: match against a file type, e.g. -type f will match files
//     -name: glob to match against file
//     -l: long listing. It's not very good, yet, but it's useful enough.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/find"
)

const cmd = "find [opts] starting-at-path"

var (
	perm      = flag.Int("mode", -1, "Permissions")
	fileType  = flag.String("type", "", "File type")
	name      = flag.String("name", "", "glob for name")
	long      = flag.Bool("l", false, "long listing")
	debug     = flag.Bool("d", false, "Enable debugging in the find package")
	fileTypes = map[string]os.FileMode{
		"f":         0,
		"file":      0,
		"d":         os.ModeDir,
		"directory": os.ModeDir,
	}
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) != 1 {
		flag.Usage()
	}
	root := a[0]
	var mask, mode os.FileMode
	if *perm != -1 {
		mask = os.ModePerm
		mode = os.FileMode(*perm)
	}
	if *fileType != "" {
		intType, ok := fileTypes[*fileType]
		if !ok {
			var keys []string
			for key := range fileTypes {
				keys = append(keys, key)
			}
			log.Fatalf("%v is not a valid file type\n valid types are %v", *fileType, strings.Join(keys, ","))
		}
		mode |= intType
		mask |= os.ModeType
	}

	f, err := find.New(func(f *find.Finder) error {
		f.Root = root
		f.Pattern = *name
		f.ModeMask = mask
		f.Mode = mode
		if *debug {
			f.Debug = log.Printf
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	go f.Find()
	for l := range f.Names {
		if l.Err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", l.Name, l.Err)
			continue
		}
		// TODO: get long listing formats out of ls and into a package.
		if *long {
			fmt.Printf("%v\n", l.FileInfo)
			continue
		}
		fmt.Printf("%s\n", l.Name)
	}
}
