// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Find finds files. It is similar to the Unix command. It uses REs, not globs,
// for matching.
//
// OPTIONS:
//
//	-d: enable debugging in the find package
//	-mode integer-arg: match against mode, e.g. -mode 0755
//	-type: match against a file type, e.g. -type f will match files
//	-name: glob to match against file
//	-l: long listing. It's not very good, yet, but it's useful enough.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/find"
)

var errNotValidType = errors.New("not a valid file type")
var errUsage = errors.New("find [opts] starting-at-path")

type flags struct {
	fileType string
	name     string
	perm     int
	long     bool
	debug    bool
}

type cmd struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	flags
}

// reorderArgs reorders arguments so flags are moved to the front, which is the
// way the "flag" package can parse them.
func reorderArgs(args []string) []string {
	var (
		newArgs []string
		i       int
	)

	for i < len(args) {
		var (
			arg          = args[i]
			expectsValue = arg == "-name" || arg == "-type" || arg == "-mode"
			hasNext      = i+1 < len(args)
			isFlag       = strings.HasPrefix(arg, "-")
		)

		prepend := func(args ...string) {
			newArgs = append(args, newArgs...)
		}

		switch {
		case expectsValue && hasNext:
			next := args[i+1]
			prepend(arg, next)
			i += 2
		case isFlag:
			prepend(arg)
			i++
		default:
			newArgs = append(newArgs, arg)
			i++
		}
	}

	return newArgs
}

func command(stdout, stderr io.Writer, args []string) (*cmd, error) {
	var f flags

	fs := flag.NewFlagSet("find", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&f.fileType, "type", "", "file type")
	fs.StringVar(&f.name, "name", "", "glob for name")
	fs.IntVar(&f.perm, "mode", -1, "permissions")
	fs.BoolVar(&f.long, "l", false, "long listing")
	fs.BoolVar(&f.debug, "d", false, "enable debugging in the find package")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: find [opts] starting-at-path\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(reorderArgs(args)); err != nil {
		return nil, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return nil, errUsage
	}

	return &cmd{
		stdout: stdout,
		stderr: stderr,
		args:   fs.Args(),
		flags:  f,
	}, nil
}

func (c *cmd) run() error {
	fileTypes := map[string]os.FileMode{
		"f":         0,
		"file":      0,
		"d":         os.ModeDir,
		"directory": os.ModeDir,
		"s":         os.ModeSocket,
		"p":         os.ModeNamedPipe,
		"l":         os.ModeSymlink,
		"c":         os.ModeCharDevice | os.ModeDevice,
		"b":         os.ModeDevice,
	}

	root := c.args[0]

	var mask, mode os.FileMode
	if c.perm != -1 {
		mask = os.ModePerm
		mode = os.FileMode(c.perm)
	}
	if c.fileType != "" {
		intType, ok := fileTypes[c.fileType]
		if !ok {
			var keys []string
			for key := range fileTypes {
				keys = append(keys, key)
			}
			return fmt.Errorf("%w: %v\n valid types are %v", errNotValidType, c.fileType, strings.Join(keys, ","))
		}
		mode |= intType
		mask |= os.ModeType
	}

	debugLog := func(string, ...any) {}
	if c.debug {
		debugLog = log.Printf
	}
	names := find.Find(context.Background(),
		find.WithRoot(root),
		find.WithModeMatch(mode, mask),
		find.WithFilenameMatch(c.name),
		find.WithDebugLog(debugLog),
	)

	for l := range names {
		if l.Err != nil {
			fmt.Fprintf(c.stderr, "%s: %v\n", l.Name, l.Err)
			continue
		}
		if c.long {
			fmt.Fprintf(c.stdout, "%s\n", l)
			continue
		}
		fmt.Fprintf(c.stdout, "%s\n", l.Name)
	}

	return nil
}

func main() {
	c, err := command(os.Stdout, os.Stderr, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}
