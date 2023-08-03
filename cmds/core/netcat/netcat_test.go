// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	_, err := command(nil, nil, nil, params{}, nil)
	if !errors.Is(err, errMissingHostnamePort) {
		t.Errorf("expected %v, got %v", errMissingHostnamePort, err)
	}
}

func TestParseParams(t *testing.T) {
	p := parseParams()

	// test defaults
	if p.network != "tcp" {
		t.Errorf("expected default network to be tcp, got %s", p.network)
	}

	if p.listen != false {
		t.Errorf("expected default listen to be false, got %t", p.listen)
	}

	if p.verbose != false {
		t.Errorf("expected default verbose to be false, got %t", p.verbose)
	}
}

func setupEchoServer(t *testing.T) string {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: []byte{127, 0, 0, 1}, Port: 0})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := l.AcceptTCP()
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, 64)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		if _, err := conn.Write(buf[:n]); err != nil {
			return
		}
	}()

	return l.Addr().String()
}

func TestTCP(t *testing.T) {
	addr := setupEchoServer(t)

	stdin := strings.NewReader("hello world")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd, err := command(stdin, stdout, stderr, params{network: "tcp", verbose: true}, []string{addr})
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.run()
	if err != nil {
		t.Fatal(err)
	}

	if stdout.String() != "hello world" {
		t.Errorf("expected 'hello world', got %q", stdout.String())
	}

	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "Connected") && !strings.Contains(stderrStr, "Disconnected") {
		t.Errorf("expected 'Connected' and 'Listening' in stderr, got %q", stderrStr)
	}
}
