// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

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
