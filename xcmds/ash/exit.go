// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"strconv"
)

func init() {
	addBuiltIn("exit", exitBuiltin)
}

func exitBuiltin(c *Command) error {
	var err error
	if len(c.argv) == 0 {
		os.Exit(0)
	} else if len(c.argv) > 1 {
		err = errors.New("Too many arguments")
	} else if ret, err2 := strconv.Atoi(c.argv[0]); err2 == nil {
		os.Exit(ret)
	} else {
		err = errors.New("Non numeric argument")
	}
	return err
}
