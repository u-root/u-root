// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"os"
	guser "os/user"
	"path/filepath"
	"strings"
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
	if err := os.WriteFile(confPath, conf, 0o600); err != nil {
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
// return errInvalidArgs
func TestBadArgs(t *testing.T) {
	if err := run([]string{"sshtest"}, os.Stdin, io.Discard, io.Discard); err != errInvalidArgs {
		t.Fatalf(`run(["sshtest"], ...) = %v, want %v`, err, errInvalidArgs)
	}
}

// This attempts to connect to git@github.com and run a command. It will fail but that's ok.
// TODO: restore this test, but first we need to add better support for locating known_hosts files with this in it:
// github.com,140.82.121.4 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
func TestSshCommand(t *testing.T) {
	t.Skipf("Skipping for now, until we can relocate the known_hosts file")
	kf := genPrivKey(t)
	if err := run([]string{"sshtest", "-i", kf, "git@github.com", "pwd"}, os.Stdin, io.Discard, io.Discard); err == nil || !strings.Contains(err.Error(), "unable to connect") {
		t.Fatalf(`run(["sshtest"], ...) = %v, want "...unable to connect..."`, err)
	}
}

// This attempts to connect to git@github.com and start a shell. It will fail but that's ok.
func TestSshShell(t *testing.T) {
	t.Skipf("Skipping for now, until we can relocate the known_hosts file")
	kf := genPrivKey(t)
	if err := run([]string{"sshtest", "-i", kf, "git@github.com"}, os.Stdin, io.Discard, io.Discard); err == nil || !strings.Contains(err.Error(), "unable to connect") {
		t.Fatalf(`run(["sshtest"], ...) = %v, want "...unable to connect..."`, err)
	}
}

// returns the path containing a private key
func genPrivKey(t *testing.T) string {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// Trying to throw off simple-minded scanners looking for key files
	x := "PRIVATE"
	block := &pem.Block{
		Type:  "RSA " + x + " KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}

	dir := t.TempDir()
	kf := filepath.Join(dir, "kf")
	f, err := os.OpenFile(kf, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := pem.Encode(f, block); err != nil {
		log.Fatal(err)
	}
	return kf
}
