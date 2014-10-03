// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// our first builtin: cd
package main

import (
	"fmt"
	"os"
	"time"
)

func init() {
	addBuiltIn("time", runtime)
}

func runtime(cmd string, s []string) error {
	var err error
	start := time.Now()
	if len(s) > 0 {
		err = runit(s[0], s[1:])
	}
	cost := time.Since(start)
	fmt.Fprintf(os.Stderr, "%v\n", cost)
	return err
}
