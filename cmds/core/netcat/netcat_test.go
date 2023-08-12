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

func setupEchoServerUDP(t *testing.T) string {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 0})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer conn.Close()
		buf := make([]byte, 64)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}

		if _, err := conn.WriteToUDP(buf[:n], raddr); err != nil {
			return
		}
	}()

	return conn.LocalAddr().String()
}

type testBuffer struct {
	buf []byte
	ch  chan string
}

func (t *testBuffer) Write(p []byte) (int, error) {
	t.buf = append(t.buf, p...)
	t.ch <- string(p)
	return len(p), nil
}

func TestDialUDP(t *testing.T) {
	addr := setupEchoServerUDP(t)

	stdin := strings.NewReader("hello world")
	stdout := &testBuffer{ch: make(chan string)}
	stderr := &bytes.Buffer{}

	cmd, err := command(stdin, stdout, stderr, params{network: "udp", verbose: true}, []string{addr})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_ = cmd.run()
	}()

	res := <-stdout.ch
	if res != "hello world" {
		t.Errorf("expected 'hello world', got %q", res)
	}
}

func TestListenUDP(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &testBuffer{ch: make(chan string)}
	stderr := &testBuffer{ch: make(chan string)}

	cmd, err := command(stdin, stdout, stderr, params{network: "udp", listen: true, verbose: true}, []string{"127.0.0.1:0"})
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		_ = cmd.run()
	}()

	listenOn := <-stderr.ch // consume listening on message
	srvAddr := strings.TrimSpace(string(listenOn[13:]))

	conn, err := net.Dial("udp", srvAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	<-stderr.ch // consume connected to message
	res := <-stdout.ch
	if res != "hello world" {
		t.Errorf("expected 'hello world', got %q", res)
	}

	_, err = conn.Write([]byte("bye"))
	if err != nil {
		t.Fatal(err)
	}
	res = <-stdout.ch
	if res != "bye" {
		t.Errorf("expected 'bye', got %q", res)
	}
}

func TestWrongNetwork(t *testing.T) {
	cmd, err := command(nil, nil, nil, params{network: "quic"}, []string{"127.0.0.1:8080"})
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.run()
	if err == nil {
		t.Error("quic is not a valid network, expected error")
	}
}
