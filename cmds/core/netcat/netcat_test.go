// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"
	"testing"
)

func freePort(t *testing.T) string {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

func TestArgs(t *testing.T) {
	_, err := command(nil, nil, nil, params{}, nil)
	if !errors.Is(err, errMissingHostnamePort) {
		t.Errorf("expected %v, got %v", errMissingHostnamePort, err)
	}
}

func TestTCP(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	addr := net.JoinHostPort("localhost", freePort(t))
	nc, err := command(stdin, stdout, stderr, params{
		network: "tcp",
		listen:  true,
		verbose: true,
	}, []string{addr})

	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan error)

	go func(ch chan error) {
		err := nc.run()
		ch <- err
	}(ch)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Write([]byte("hello world")); err != nil {
		t.Fatal(err)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	ncErr := <-ch
	if ncErr != nil {
		t.Error("expected nil, got", ncErr)
	}

	if stdout.String() != "hello world" {
		t.Errorf("expected hello world, got %q", stdout.String())
	}

	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "Listening on") {
		t.Errorf("expected to contain 'Listening on', got %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "Connected to") {
		t.Errorf("expected to contain 'Connected to', got %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "Disconnected") {
		t.Errorf("expected to contain 'Disconnected', got %q", stderrOutput)
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
