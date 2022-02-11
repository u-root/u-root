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
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/util"

	"github.com/u-root/u-root/pkg/find"
)

const cmd = "find [opts] starting-at-path"

type flags struct {
	perm     int
	filetype string
	name     string
	long     bool
	debug    bool
}

var (
	fargs     = flags{}
	fileTypes = map[string]os.FileMode{
		"f":         0,
		"file":      0,
		"d":         os.ModeDir,
		"directory": os.ModeDir,
	}
)

func init() {
	flag.IntVar(&fargs.perm, "mode", -1, "Permissions")
	flag.StringVar(&fargs.filetype, "type", "", "File type")
	flag.StringVar(&fargs.name, "name", "", "glob for name")
	flag.BoolVar(&fargs.long, "l", false, "long listing")
	flag.BoolVar(&fargs.debug, "d", false, "Enable debugging in the find package")
}

func runFind(out io.Writer, errOut io.Writer, fargs flags, arg []string) error {
	if len(arg) != 1 {
		flag.Usage()
		return nil
	}
	root := arg[0]

	var mask, mode os.FileMode
	if fargs.perm != -1 {
		mask = os.ModePerm
		mode = os.FileMode(fargs.perm)
	}
	if fargs.filetype != "" {
		intType, ok := fileTypes[fargs.filetype]
		if !ok {
			var keys []string
			for key := range fileTypes {
				keys = append(keys, key)
			}
			return fmt.Errorf("%v is not a valid file type\n valid types are %v", fargs.filetype, strings.Join(keys, ","))
		}
		mode |= intType
		mask |= os.ModeType
	}

	debugLog := func(string, ...interface{}) {}
	if fargs.debug {
		debugLog = log.Printf
	}
	names := find.Find(context.Background(),
		find.WithRoot(root),
		find.WithModeMatch(mode, mask),
		find.WithFilenameMatch(fargs.name),
		find.WithDebugLog(debugLog),
	)
	for l := range names {
		if l.Err != nil {
			fmt.Fprintf(errOut, "%v: %v\n", l.Name, l.Err)
			continue
		}
		if fargs.long {
			fmt.Fprintf(out, "%s\n", l)
			continue
		}
		fmt.Fprintf(out, "%s\n", l.Name)
	}
	return nil
}

func main() {
	flag.Parse()
	util.Usage(cmd)
	if err := runFind(os.Stdout, os.Stderr, fargs, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
