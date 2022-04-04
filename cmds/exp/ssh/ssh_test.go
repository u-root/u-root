// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"os"
	guser "os/user"
	"path/filepath"
	"testing"
)

type destTest struct {
	input string
	user  string
	host  string
	port  string
}

var destTests = []destTest{
	{"example.org", "", "example.org", "22"},
	{"foo@example.org", "foo", "example.org", "22"},
	{"foo@example.org", "foo", "example.org", "22"},
	{"ssh://192.168.0.2:23", "", "192.168.0.2", "23"},
	{"ssh://x@example.org", "x", "example.org", "22"},
}

func TestParseDest(t *testing.T) {
	for _, x := range destTests {
		if x.user == "" {
			var u *guser.User
			u, _ = guser.Current()
			x.user = u.Username
		}
		user, host, port, err := parseDest(x.input)
		if err != nil {
			t.Fatal(err)
		}
		if user != x.user || host != x.host || port != x.port {
			t.Fatalf("failed on %v: got %v, %v, %v", x, user, host, port)
		}
	}
}

// Load a config file and ask for the keyfile for a host in it
// By populating a real file & reading it, we get to test loadConfig too
func TestGetKeyFile(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, "sshconfig")
	conf := []byte(`Host foo
	IdentityFile bar_key`)
	if err := ioutil.WriteFile(confPath, conf, 0600); err != nil {
		t.Fatal(err)
	}
	if err := loadConfig(confPath); err != nil {
		t.Fatal(err)
	}
	if kf := getKeyFile("foo", ""); kf != "bar_key" {
		t.Fatalf(`getKeyFile("foo", "") = %v, want "bar_key"`, kf)
	}
}

// Test what happens if we pass invalid command-line arguments... should
// return ErrInvalidArgs
func TestBadArgs(t *testing.T) {
	if err := run([]string{"sshtest"}, os.Stdin, io.Discard, io.Discard); err != ErrInvalidArgs {
		t.Fatalf(`run(["sshtest"], ...) = %v, want %v`, err, ErrInvalidArgs)
	}
}
