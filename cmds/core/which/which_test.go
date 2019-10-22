// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Rafael Campos Nunes <rafaelnunes@engineer.com>

package main

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"strings"
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

	for i, t := range tests {
		tests[i].pathOnSys, err = exec.Command("which", "-a", t.name).Output()
		if err != nil {
			return err
		}
	}

	return nil
}

// TestWhichUnique tests `which` command against one POSIX command that are included in Linux.
// The output of which.go has to be exactly equal to the output of which itself.
func TestWhichUnique(t *testing.T) {
	err := setup()

	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}

	commands := []string{"cat"}
	var b bytes.Buffer
	if err := which(&b, strings.Split(p, ":"), commands[:], true); err != nil {
		t.Fatal(err)
	}

	// Comparing against only the cat command.
	if !reflect.DeepEqual(b.Bytes(), tests[0].pathOnSys) {
		t.Fatalf("Locating commands has failed, wants: %v, got: %v", string(tests[0].pathOnSys), b.String())
	}
}

// TestWhichMultiple tests `which` command against the three POSIX commands that are included in Linux.
// The output of which.go has to be exactly equal to the output of which itself. If it works with
// three, it should work with more commands as well.
func TestWhichMultiple(t *testing.T) {
	err := setup()

	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}

	pathsCombined := []byte{}
	commands := []string{}
	for _, t := range tests {
		pathsCombined = append(pathsCombined, t.pathOnSys...)
		commands = append(commands, t.name)
	}

	var b bytes.Buffer
	if err := which(&b, strings.Split(p, ":"), commands[:], true); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b.Bytes(), pathsCombined) {
		t.Fatalf("Locating commands has failed, wants: %v, got: %v", string(pathsCombined), b.String())
	}
}
