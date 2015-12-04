// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// created by Rafael Campos Nunes <rafaelnunes@engineer.com>

package main

import (
	"os"
	"os/exec"
	"testing"
)

type commands struct {
	name      string
	pathOnSys []byte
}

// in setup we fill the pathOnSys variables with their corresponding path on the system.
var (
	tests = []commands{
		{
			"cat",
			[]byte{},
		},
		{
			"which",
			[]byte{},
		},
		{
			"sed",
			[]byte{},
		},
		{
			"ldd",
			[]byte{},
		},
	}

	p = os.Getenv("PATH")
)

func setup() error {
	var err error

	for i := range tests {
		tests[i].pathOnSys, err = exec.Command("which", tests[i].name).Output()
		if err != nil {
			return err
		}
	}

	return nil
}

/* Test_which_1 tests `which` command against one POSIX command that are included in Linux.
 * The output of which.go has to be exactly equal to the output of which itself.
 */
func Test_which_1(t *testing.T) {
	err := setup()

	if err != nil {
		t.Fatal("setup has failed, %v", err)
	}

	commands := [1]string{"cat"}
	which(p, commands[:])
}

/* Test_which_1 tests `which` command against the three POSIX commands that are included in Linux.
 * The output of which.go has to be exactly equal to the output of which itself. If it works with
 * three, it should work with more commands as well.
 */
func Test_which_2(t *testing.T) {
	err := setup()

	if err != nil {
		t.Fatal("setup has failed, %v", err)
	}

	commands := [3]string{"which", "ldd", "sed"}
	which(p, commands[:])
}
