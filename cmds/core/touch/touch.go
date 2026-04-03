// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// touch changes file access and modification times.
//
// Synopsis:
//
//	touch [-amc] [-d datetime] file...
//
// Description:
//
//	If a file does not exist, it will be created unless -c is specified.
//
// Options:
//
//	-a: change only the access time
//	-m: change only the modification time
//	-c: do not create any file if it does not exist
//	-d: use specified time instead of current time (RFC3339 format)
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type params struct {
	time         time.Time
	access       bool
	modification bool
	create       bool
}

type cmd struct {
	params
	args []string
}

var errNoFiles = errors.New("no files specified")

func command(args ...string) (*cmd, error) {
	var c cmd
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	var access bool
	var modification bool
	var create bool
	var dateTime string

	fs.BoolVar(&access, "a", false, "change only the access time")
	fs.BoolVar(&modification, "m", false, "change only the modification time")
	fs.BoolVar(&create, "c", false, "do not create any file if it does not exist")
	fs.StringVar(&dateTime, "d", "", "use specified time instead of current time RFC3339")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: touch [-amc] [-d datetime] file...\n\n")
		fmt.Fprintf(fs.Output(), "touch changes file access and modification times.\n")
		fmt.Fprintf(fs.Output(), "If a file does not exist, it will be created unless -c is specified.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args[1:])); err != nil {
		return nil, err
	}

	if len(fs.Args()) == 0 {
		fs.Usage()
		return nil, errNoFiles
	}
	c.args = fs.Args()

	var err error
	c.params, err = parseParams(dateTime, access, modification, create)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func parseParams(dateTime string, access, modification, create bool) (params, error) {
	t := time.Now()
	if dateTime != "" {
		var err error
		t, err = time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return params{}, err
		}
	}
	return params{
		access:       access || (!access && !modification),
		modification: modification || (!access && !modification),
		create:       create,
		time:         t,
	}, nil
}

func (c *cmd) run() error {
	var errs error
	for _, arg := range c.args {
		_, existsErr := os.Stat(arg)
		notExist := os.IsNotExist(existsErr)
		if notExist {
			if c.create {
				continue
			}

			f, err := os.Create(arg)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			f.Close()
		}

		accessTime := time.Time{}
		if c.access || notExist {
			accessTime = c.time
		}
		modificationTime := time.Time{}
		if c.modification || notExist {
			modificationTime = c.time
		}

		err := os.Chtimes(arg, accessTime, modificationTime)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func main() {
	c, err := command(os.Args...)
	if err != nil {
		log.Fatalf("touch: %v", err)
	}
	if err := c.run(); err != nil {
		log.Fatalf("touch: %v", err)
	}
}
