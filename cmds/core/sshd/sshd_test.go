// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !windows

package main

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestParseParmams(t *testing.T) {
	params := parseParams()
	if params.debug {
		t.Error("expected default debug to be false, got true")
	}
	if params.keys != "authorized_keys" {
		t.Errorf("expected default keys to be authorized_keys, got %q", params.keys)
	}
	if params.privkey != "id_rsa" {
		t.Errorf("expected default privatekey to be id_rsa, got %q", params.privkey)
	}
	if params.net != "tcp" {
		t.Errorf("expected default net to be tcp, got %q", params.net)
	}
	if params.ip != "0.0.0.0" {
		t.Errorf("expected default ip to be 0.0.0.0, got %q", params.ip)
	}
	if params.port != "2022" {
		t.Errorf("expected default port to be 2022, got %q", params.port)
	}
}

func TestListen(t *testing.T) {
	sockPath := filepath.Join(t.TempDir(), "sshd.sock")

	for _, tc := range []struct {
		name     string
		network  string
		ip       string
		port     string
		wantNet  string
		wantAddr string
		wantErr  bool
	}{
		{
			name:    "tcp listener",
			network: "tcp",
			ip:      "127.0.0.1",
			port:    "0",
			wantNet: "tcp",
		},
		{
			name:    "empty network defaults to tcp",
			network: "",
			ip:      "127.0.0.1",
			port:    "0",
			wantNet: "tcp",
		},
		{
			name:     "unix socket listener",
			network:  "unix",
			ip:       "ignored",
			port:     sockPath,
			wantNet:  "unix",
			wantAddr: sockPath,
		},
		{
			name:    "vsock invalid non-numeric port",
			network: "vsock",
			port:    "not-a-port",
			wantErr: true,
		},
		{
			name:    "vsock port exceeding 32 bits",
			network: "vsock",
			port:    "4294967296",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ln, err := listen(tc.network, tc.ip, tc.port)
			if ln != nil {
				defer ln.Close()
			}
			if (err != nil) != tc.wantErr {
				t.Fatalf("listen(%q, %q, %q) err = %v, wantErr %v", tc.network, tc.ip, tc.port, err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if got := ln.Addr().Network(); got != tc.wantNet {
				t.Errorf("listener network = %q, want %q", got, tc.wantNet)
			}
			if tc.wantAddr != "" && ln.Addr().String() != tc.wantAddr {
				t.Errorf("listener address = %q, want %q", ln.Addr().String(), tc.wantAddr)
			}
		})
	}

	t.Run("vsock 32-bit port above 65535", func(t *testing.T) {
		ln, err := listen("vsock", "", "70000")
		if ln != nil {
			defer ln.Close()
		}
		var numErr *strconv.NumError
		if errors.As(err, &numErr) {
			t.Fatalf("expected 32-bit vsock port 70000 to pass port parsing, got %v", err)
		}
	})
}

func TestConfigErrors(t *testing.T) {
	t.Run("authorized_keys file does not exist", func(t *testing.T) {
		cmd := command(params{
			keys: "filenotexist",
		})

		err := cmd.run()
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("authorized_keys file don't have public keys", func(t *testing.T) {
		tf, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tf.Name())
		tf.WriteString("keys")

		cmd := command(params{
			keys: tf.Name(),
		})

		err = cmd.run()
		if err == nil {
			t.Errorf("expected ssh: no key found, got nil")
		}
	})
	t.Run("private key file does not exist", func(t *testing.T) {
		cmd := command(params{
			privkey: "filenotexist",
			keys:    "./testdata/id_rsa.pub",
		})

		err := cmd.run()
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("privete key file does have private key", func(t *testing.T) {
		tf, err := os.CreateTemp("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tf.Name())
		tf.WriteString("privatekey")

		cmd := command(params{
			privkey: tf.Name(),
			keys:    "./testdata/id_rsa.pub",
		})

		err = cmd.run()
		if err == nil {
			t.Error("expected ssh: no key found, got nil")
		}
	})
	t.Run("wrong port host config", func(t *testing.T) {
		cmd := command(params{
			privkey: "./testdata/id_rsa",
			keys:    "./testdata/id_rsa.pub",
			ip:      "host",
			port:    "port",
		})

		err := cmd.run()
		if err == nil {
			t.Error("expected listen tcp: lookup tcp/port: unknown port, got nil")
		}
	})
}

func connect(t *testing.T, addr string, cfg *ssh.ClientConfig) *ssh.Client {
	t.Helper()
	var count int
	for {
		if count > 10 {
			t.Fatal("can't connect to sshd server")
		}
		clt, err := ssh.Dial("tcp", addr, cfg)

		if err == nil {
			return clt
		}

		time.Sleep(100 * time.Millisecond)
		count++
	}
}

func TestSessionRun(t *testing.T) {
	cmd := command(params{
		privkey: "./testdata/id_rsa",
		keys:    "./testdata/id_rsa.pub",
		ip:      "127.0.0.1",
		port:    "2022",
		debug:   true,
	})

	go cmd.run()

	pk, err := os.ReadFile("./testdata/id_rsa")
	if err != nil {
		t.Fatalf("can't read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		t.Fatalf("can't parse private key: %v", err)
	}

	cfg := ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second,
	}

	clt := connect(t, net.JoinHostPort(cmd.ip, cmd.port), &cfg)

	session, err := clt.NewSession()
	if err != nil {
		t.Fatalf("can't create session: %v", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		t.Fatalf("can't pipe stdin: %v", err)
	}
	stdin.Close()

	output, err := session.StdoutPipe()
	if err != nil {
		t.Fatalf("can't pipe output: %v", err)
	}

	session.Run("echo hello u-root")

	b := make([]byte, 128)
	n, err := output.Read(b)
	if err != nil {
		t.Fatalf("can't read output: %v", err)
	}

	if string(b[:n]) != "hello u-root\n" {
		t.Errorf("expected hello u-root, got %q", string(b[:n]))
	}
}
