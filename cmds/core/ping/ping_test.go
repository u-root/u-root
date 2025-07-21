// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/guest"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type testConn struct {
	lastMessage []byte
}

func (tc *testConn) ReadFrom(b []byte) (int, net.Addr, error) {
	m, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), tc.lastMessage)
	if err != nil {
		return 0, nil, err
	}

	body := m.Body.(*icmp.Echo)

	respone := icmp.Message{
		Type: ipv4.ICMPTypeEchoReply,
		Code: 0,
		Body: &icmp.Echo{
			ID:   body.ID,
			Seq:  body.Seq,
			Data: body.Data,
		},
	}

	resp, err := respone.Marshal(nil)
	if err != nil {
		return 0, nil, err
	}

	n := copy(b, resp)
	return n, nil, nil
}

func (tc *testConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	tc.lastMessage = b
	return len(b), nil
}

func (tc *testConn) Close() error {
	return nil
}

func (tc *testConn) LocalAddr() net.Addr {
	return nil
}

func (tc *testConn) SetDeadline(t time.Time) error {
	return nil
}

func (tc *testConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (tc *testConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type pingOutputLine struct {
	addr    string
	size    int
	seq     int
	audible bool
}

func parsePingLines(t *testing.T, output []byte) []pingOutputLine {
	t.Helper()
	var lines []pingOutputLine

	for line := range bytes.SplitSeq(output, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		var pl pingOutputLine
		sp := strings.Split(string(line), " ")
		if sp[0][0] == '\a' {
			pl.audible = true
			size, err := strconv.Atoi(sp[0][1:])
			if err != nil {
				t.Fatalf("size atoi failed: %v", err)
			}
			pl.size = size
		} else {
			size, err := strconv.Atoi(sp[0])
			if err != nil {
				t.Fatalf("size atoi failed: %v", err)
			}
			pl.size = size
		}

		pl.addr = sp[3][:len(sp[3])-1]
		sp = strings.Split(sp[4], "=")
		sqn, err := strconv.Atoi(sp[1])
		if err != nil {
			t.Fatalf("address atoi failed: %v", err)
		}
		pl.seq = sqn

		lines = append(lines, pl)
	}

	return lines
}

func TestPing(t *testing.T) {
	var paConn net.PacketConn = &testConn{}
	paConn.Close()

	stdout := &bytes.Buffer{}
	cmd := &cmd{
		stdout: stdout,
		conn:   paConn,
		params: params{
			host:       "1.1.1.1",
			packetSize: 56,
			intv:       1000,
			wtf:        100,
			iter:       1,
			net6:       false,
			audible:    true,
		},
	}

	err := cmd.run()
	if err != nil {
		t.Error(err)
	}

	lines := parsePingLines(t, stdout.Bytes())
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}

	if lines[0].size != 64 {
		t.Errorf("expected size 64 (56 + header), got %d", lines[0].size)
	}
	if !lines[0].audible {
		t.Errorf("expected audible, got %v", lines[0].audible)
	}
	if lines[0].seq != 1 {
		t.Errorf("expected seq 1, got %d", lines[0].seq)
	}
	if lines[0].addr != "1.1.1.1" {
		t.Errorf("expected addr, got %s", lines[0].addr)
	}
}

func TestRawPing(t *testing.T) {
	guest.SkipIfNotInVM(t)
	stdout := &bytes.Buffer{}
	cmd, err := command(stdout, params{
		packetSize: 56,
		host:       "127.0.0.1",
		intv:       1000,
		wtf:        100,
		iter:       1,
	})
	if err != nil {
		t.Errorf("command() failed: %v", err)
	}

	err = cmd.run()
	if err != nil {
		t.Errorf("run() failed: %v", err)
	}

	lines := parsePingLines(t, stdout.Bytes())

	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}

	if lines[0].size != 64 {
		t.Errorf("expected size 64 (56 + header), got %d", lines[0].size)
	}
	if lines[0].audible {
		t.Errorf("expected no audible, got %v", lines[0].audible)
	}
	if lines[0].seq != 1 {
		t.Errorf("expected seq 1, got %d", lines[0].seq)
	}
	if lines[0].addr != "127.0.0.1" {
		t.Errorf("expected addr, got %s", lines[0].addr)
	}
}
