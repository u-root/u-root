// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type makeit struct {
	n string      // name
	m os.FileMode // mode
	c []byte      // content
}

var old = makeit{
	n: "old.txt",
	m: 0777,
	c: []byte("old"),
}

var new = makeit{
	n: "new.txt",
	m: 0777,
	c: []byte("new"),
}

var tests = []makeit{
	{
		n: "hi1.txt",
		m: 0666,
		c: []byte("hi"),
	},
	{
		n: "hi2.txt",
		m: 0777,
		c: []byte("hi"),
	},
	old,
	new,
}

func setup() (string, error) {
	d, err := ioutil.TempDir(os.TempDir(), "hi.dir")
	if err != nil {
		return "", err
	}

	tmpdir := filepath.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0777); err != nil {
		return "", err
	}

	for _, t := range tests {
		if err := ioutil.WriteFile(filepath.Join(d, t.n), []byte(t.c), t.m); err != nil {
			return "", err
		}
	}

	return d, nil
}

func resetFlags() {
	// TODO: Is it possible to extract the default out of the declaration?
	*update = false
}

func TestMv(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("err")
	}
	defer os.RemoveAll(d)

	fmt.Println("Renaming file...")
	files1 := []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi4.txt")}
	if err := mv(files1, false); err != nil {
		t.Error(err)
	}

	dsub := filepath.Join(d, "hi.sub.dir")

	fmt.Println("Moving files to directory...")
	files2 := []string{filepath.Join(d, "hi2.txt"), filepath.Join(d, "hi4.txt"), dsub}
	if err := mv(files2, true); err != nil {
		t.Error(err)
	}
}

func TestMvUpdate(t *testing.T) {
	*update = true
	defer resetFlags()
	d, err := setup()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Testing mv -u...")

	// Ensure that the newer file actually has a newer timestamp
	currentTime := time.Now().Local()
	oldTime := currentTime.Add(-10 * time.Second)
	err = os.Chtimes(filepath.Join(d, old.n), oldTime, oldTime)
	if err != nil {
		t.Error(err)
	}
	err = os.Chtimes(filepath.Join(d, new.n), currentTime, currentTime)
	if err != nil {
		t.Error(err)
	}

	// Check that it doesn't downgrade files with -u switch
	var files1 = []string{filepath.Join(d, old.n), filepath.Join(d, new.n)}
	if err := mv(files1, false); err != nil {
		t.Error(err)
	}
	newContent1, err := ioutil.ReadFile(filepath.Join(d, new.n))
	if err != nil {
		t.Error(err)
	}
	if bytes.Equal(newContent1, old.c) {
		t.Error("Newer file was overwritten by older file. Should not happen with -u.")
	}

	// Check that it does update files with -u switch
	var files2 = []string{filepath.Join(d, new.n), filepath.Join(d, old.n)}
	if err := mv(files2, false); err != nil {
		t.Error(err)
	}
	newContent2, err := ioutil.ReadFile(filepath.Join(d, old.n))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(newContent2, new.c) {
		t.Error("Older file was not overwritten by newer file. Should happen with -u.")
	}
	if _, err := os.Lstat(filepath.Join(d, old.n)); err != nil {
		t.Error("The new file shouldn't be there anymore.")
	}
}
