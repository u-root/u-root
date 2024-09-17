// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

func init() {
	_ = addBuiltIn("pwd", pwd)
}

// pwd command: print the current working directory
func pwd(_ *Command) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("pwd: %v", err)
	}
	fmt.Println(dir)
	return nil
}
