// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type params struct {
	time         time.Time
	access       bool
	modification bool
	create       bool
}

type cmd struct {
	stderr io.Writer
	params
	args []string
}

func command(stderr io.Writer, p params, args ...string) *cmd {
	return &cmd{
		args:   args,
		params: p,
		stderr: stderr,
	}
}

func parseParams(dateTime string, access, modification, create bool) (params, error) {
	flag.Parse()
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
	// TODO: move to errors.Join in future
	// right now log error in stderr and return last error
	// to indicate a non zero exit code
	var lastError error
	for _, arg := range c.args {
		_, existsErr := os.Stat(arg)
		notExist := os.IsNotExist(existsErr)
		if notExist {
			if c.create {
				continue
			}

			f, err := os.Create(arg)
			if err != nil {
				lastError = err
				fmt.Fprintf(c.stderr, "touch: %v\n", err)
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
			lastError = err
			fmt.Fprintf(c.stderr, "touch: %v\n", err)
		}
	}

	return lastError
}

func main() {
	access := flag.Bool("a", false, "change only the access time")
	modification := flag.Bool("m", false, "change only the modification time")
	create := flag.Bool("c", false, "do not create any file if it does not exist")
	dateTime := flag.String("d", "", "use specified time instead of current time RFC3339")
	flag.Parse()
	p, err := parseParams(*dateTime, *access, *modification, *create)
	if err != nil || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if err := command(os.Stderr, p, flag.Args()...).run(); err != nil {
		os.Exit(1)
	}
}
