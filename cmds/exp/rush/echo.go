// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

func init() {
	_ = addBuiltIn("echo", echo)
}

// echo command: print the arguments with space separation
func echo(c *Command) error {
	// Join the arguments with spaces and print them
	fmt.Println(strings.Join(c.argv, " "))
	return nil
}
