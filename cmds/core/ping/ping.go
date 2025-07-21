// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// Send icmp packets to a server to test network connectivity.
//
// Synopsis:
//
//	ping [-hV] [-c COUNT] [-i INTERVAL] [-s PACKETSIZE] [-w DEADLINE] DESTINATION
//
// Options:
//
//	-6: use ipv6 (ip6:ipv6-icmp)
//	-s: data size (default: 64)
//	-c: # iterations, 0 to run forever (default)
//	-i: interval in milliseconds (default: 1000)
//	-V: version
//	-w: wait time in milliseconds (default: 100)
//	-a: Audible rings a bell when a packet is received
//	-h: help
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/uroot/util"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const usage = "ping [-V] [-6] [-c count] [-i interval] [-s packetsize] [-w deadline] [-a audible] destination"

type params struct {
	packetSize int
	intv       int
	wtf        int
	iter       uint64
	host       string
	net6       bool
	audible    bool
}

type cmd struct {
	stdout io.Writer
	conn   net.PacketConn
	params
}

func command(stdin io.Writer, p params) (*cmd, error) {
	netname, address := "ip4:icmp", "0.0.0.0"
	if p.net6 {
		netname, address = "ip6:ipv6-icmp", "::1"
	}
	conn, err := icmp.ListenPacket(netname, address)
	if err != nil {
		return nil, fmt.Errorf("can't setup %s socket on %s: %w", netname, address, err)
	}

	return &cmd{stdin, conn, p}, nil
}

func (c *cmd) run() error {
	defer c.conn.Close()
	if c.packetSize < 8 {
		return fmt.Errorf("packet size too small (must be >= 8): %v", c.packetSize)
	}

	network := "ip4"
	if c.net6 {
		network = "ip6"
	}

	addr, err := net.ResolveIPAddr(network, c.host)
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	interval := time.Duration(c.intv)
	waitFor := time.Duration(c.wtf) * time.Millisecond
	for i := uint64(0); i < c.iter; i++ {
		msg, err := c.ping(addr, i+1, waitFor)
		if err != nil {
			return fmt.Errorf("ping failed: %w", err)
		}
		if c.audible {
			msg = "\a" + msg
		}
		fmt.Fprintf(c.stdout, "%s\n", msg)
		time.Sleep(time.Millisecond * interval)
	}

	return nil
}

func (c *cmd) ping(addr *net.IPAddr, i uint64, waitFor time.Duration) (string, error) {
	c.conn.SetDeadline(time.Now().Add(waitFor))

	var echoRequestType icmp.Type = ipv4.ICMPTypeEcho
	if c.net6 {
		echoRequestType = ipv6.ICMPTypeEchoRequest
	}

	wm := icmp.Message{
		Type: echoRequestType, Code: 0, Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  int(i),
			Data: bytes.Repeat([]byte{1}, c.packetSize),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return "", fmt.Errorf("icmp.Message.Marshal failed: %w", err)
	}

	startTime := time.Now()
	_, err = c.conn.WriteTo(wb, addr)
	if err != nil {
		return "", fmt.Errorf("conn.Write failed: %w", err)
	}

	rb := make([]byte, 1500)
	n, _, err := c.conn.ReadFrom(rb)
	if err != nil {
		return "", fmt.Errorf("conn.Read failed: %w", err)
	}

	latency := time.Since(startTime)

	var echoReplyType icmp.Type = ipv4.ICMPTypeEchoReply
	if c.net6 {
		echoReplyType = ipv6.ICMPTypeEchoReply
	}

	msg, err := icmp.ParseMessage(echoReplyType.Protocol(), rb[:n])
	if err != nil {
		return "", fmt.Errorf("icmp.ParseMessage failed: %w", err)
	}

	echoReply, ok := msg.Body.(*icmp.Echo)
	if !ok {
		return "", fmt.Errorf("got %+v; want echo reply", msg)
	}

	if echoReply.ID != os.Getpid()&0xffff {
		return "", fmt.Errorf("got id %v; want %v", echoReply.ID, os.Getpid()&0xffff)
	}
	if echoReply.Seq != int(i) {
		return "", fmt.Errorf("got seq %v; want %v", echoReply.Seq, i)
	}

	return fmt.Sprintf("%d bytes from %v: icmp_seq=%v time=%v", n, c.host, i, latency), nil
}

func main() {
	var (
		net6       = flag.Bool("6", false, "use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)")
		packetSize = flag.Int("s", 56, "Data size")
		iter       = flag.Uint64("c", math.MaxUint64, "# iterations")
		intv       = flag.Int("i", 1000, "interval in milliseconds")
		wtf        = flag.Int("w", 100, "wait time in milliseconds")
		audible    = flag.Bool("a", false, "Audible rings a bell when a packet is received")
	)

	flag.Usage = util.Usage(flag.Usage, usage)
	flag.Parse()
	// options without parameters (right now just: -hV)
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Args()[0]
	cmd, err := command(os.Stdout, params{*packetSize, *intv, *wtf, *iter, host, *net6, *audible})
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.run(); err != nil {
		log.Fatal(err)
	}
}
