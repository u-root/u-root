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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/find"
)

type params struct {
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
	params params
}

func command(stdout, stderr io.Writer, params params, args []string) *cmd {
	return &cmd{
		stdout: stdout,
		stderr: stderr,
		args:   args,
		params: params,
	}
}

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "find [opts] starting-at-path"
		defUsage()
		os.Exit(1)
	}
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

	if len(c.args) != 1 {
		flag.Usage()
	}
	root := c.args[0]

	var mask, mode os.FileMode
	if c.params.perm != -1 {
		mask = os.ModePerm
		mode = os.FileMode(c.params.perm)
	}
	if c.params.fileType != "" {
		intType, ok := fileTypes[c.params.fileType]
		if !ok {
			var keys []string
			for key := range fileTypes {
				keys = append(keys, key)
			}
			return fmt.Errorf("%v is not a valid file type\n valid types are %v", c.params.fileType, strings.Join(keys, ","))
		}
		mode |= intType
		mask |= os.ModeType
	}

	debugLog := func(string, ...interface{}) {}
	if c.params.debug {
		debugLog = log.Printf
	}
	names := find.Find(context.Background(),
		find.WithRoot(root),
		find.WithModeMatch(mode, mask),
		find.WithFilenameMatch(c.params.name),
		find.WithDebugLog(debugLog),
	)

	for l := range names {
		if l.Err != nil {
			fmt.Fprintf(c.stderr, "%s: %v\n", l.Name, l.Err)
			continue
		}
		if c.params.long {
			fmt.Fprintf(c.stdout, "%s\n", l)
			continue
		}
		fmt.Fprintf(c.stdout, "%s\n", l.Name)
	}

	return nil
}

func main() {
	perm := flag.Int("mode", -1, "permissions")
	fileType := flag.String("type", "", "file type")
	name := flag.String("name", "", "glob for name")
	long := flag.Bool("l", false, "long listing")
	debug := flag.Bool("d", false, "enable debugging in the find package")
	flag.Parse()
	p := params{perm: *perm, fileType: *fileType, name: *name, long: *long, debug: *debug}
	if err := command(os.Stdout, os.Stderr, p, flag.Args()).run(); err != nil {
		log.Fatalf("find: %v", err)
	}
}
