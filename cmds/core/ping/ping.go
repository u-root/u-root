// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "ping [-V] [-6] [-c count] [-i interval] [-s packetsize] [-w deadline] [-a audible] destination"

var (
	net6       = flag.Bool("6", false, "use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)")
	packetSize = flag.Int("s", 64, "Data size")
	iter       = flag.Uint64("c", math.MaxUint64, "# iterations")
	intv       = flag.Int("i", 1000, "interval in milliseconds")
	wtf        = flag.Int("w", 100, "wait time in milliseconds")
	audible    = flag.Bool("a", false, "Audible rings a bell when a packet is received")
)

const (
	ICMP_TYPE_ECHO_REQUEST             = 8
	ICMP_TYPE_ECHO_REPLY               = 0
	ICMP_ECHO_REPLY_HEADER_IPV4_OFFSET = 20
)

const (
	ICMP6_TYPE_ECHO_REQUEST             = 128
	ICMP6_TYPE_ECHO_REPLY               = 129
	ICMP6_ECHO_REPLY_HEADER_IPV6_OFFSET = 40
)

type Ping struct {
	dial func(string, string) (net.Conn, error)
}

func New() *Ping {
	return &Ping{
		dial: net.Dial,
	}
}

func cksum(bs []byte) uint16 {
	sum := uint32(0)

	for k := 0; k < len(bs)/2; k++ {
		sum += uint32(bs[k*2]) << 8
		sum += uint32(bs[k*2+1])
	}
	if len(bs)%2 != 0 {
		sum += uint32(bs[len(bs)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = (sum >> 16) + (sum & 0xffff)
	if sum == 0xffff {
		sum = 0
	}

	return ^uint16(sum)
}

func (p *Ping) ping1(net6 bool, host string, i uint64, waitFor time.Duration) (string, error) {
	netname := "ip4:icmp"
	// todo: just figure out if it's an ip6 address and go from there.
	if net6 {
		netname = "ip6:ipv6-icmp"
	}
	c, err := p.dial(netname, host)
	if err != nil {
		return "", fmt.Errorf("net.Dial(%v %v) failed: %v", netname, host, err)
	}
	defer c.Close()

	if net6 {
		ipc := c.(*net.IPConn)
		if err := setupICMPv6Socket(ipc); err != nil {
			return "", fmt.Errorf("failed to set up the ICMPv6 connection: %w", err)
		}
	}

	// Send ICMP Echo Request
	c.SetDeadline(time.Now().Add(waitFor))
	msg := make([]byte, *packetSize)
	if net6 {
		msg[0] = ICMP6_TYPE_ECHO_REQUEST
	} else {
		msg[0] = ICMP_TYPE_ECHO_REQUEST
	}
	msg[1] = 0
	binary.BigEndian.PutUint16(msg[6:], uint16(i))
	binary.BigEndian.PutUint16(msg[4:], uint16(i>>16))
	binary.BigEndian.PutUint16(msg[2:], cksum(msg))
	if _, err := c.Write(msg[:]); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	// Get ICMP Echo Reply
	c.SetDeadline(time.Now().Add(waitFor))
	rmsg := make([]byte, *packetSize+256)
	before := time.Now()
	amt, err := c.Read(rmsg[:])
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}
	latency := time.Since(before)
	if !net6 {
		rmsg = rmsg[ICMP_ECHO_REPLY_HEADER_IPV4_OFFSET:]
	}
	if net6 {
		if rmsg[0] != ICMP6_TYPE_ECHO_REPLY {
			return "", fmt.Errorf("bad ICMPv6 echo reply type, got %d, want %d", rmsg[0], ICMP6_TYPE_ECHO_REPLY)
		}
	} else {
		if rmsg[0] != ICMP_TYPE_ECHO_REPLY {
			return "", fmt.Errorf("bad ICMP echo reply type, got %d, want %d", rmsg[0], ICMP_TYPE_ECHO_REPLY)
		}
	}
	cks := binary.BigEndian.Uint16(rmsg[2:])
	binary.BigEndian.PutUint16(rmsg[2:], 0)
	// only validate the checksum for IPv4. For IPv6 this *should* be done by the
	// TCP stack (and do we need to validate the checksum anyway?)
	if !net6 && cks != cksum(rmsg) {
		return "", fmt.Errorf("bad ICMP checksum: %v (expected %v)", cks, cksum(rmsg))
	}
	id := binary.BigEndian.Uint16(rmsg[4:])
	seq := binary.BigEndian.Uint16(rmsg[6:])
	rseq := uint64(id)<<16 + uint64(seq)
	if rseq != i {
		return "", fmt.Errorf("wrong sequence number %v (expected %v)", rseq, i)
	}

	return fmt.Sprintf("%d bytes from %v: icmp_seq=%v, time=%v", amt, host, i, latency), nil
}

func ping(host string) error {
	if *packetSize < 8 {
		return fmt.Errorf("packet size too small (must be >= 8): %v", *packetSize)
	}

	interval := time.Duration(*intv)
	p := New()
	// ping needs to run forever if count is not specified, so default value is MaxUint64
	waitFor := time.Duration(*wtf) * time.Millisecond
	for i := uint64(0); i < *iter; i++ {
		msg, err := p.ping1(*net6, host, i+1, waitFor)
		if err != nil {
			return fmt.Errorf("ping failed: %v", err)
		}
		if *audible {
			msg = "\a" + msg
		}
		log.Print(msg)
		time.Sleep(time.Millisecond * interval)
	}

	return nil
}

func main() {
	flag.Usage = util.Usage(flag.Usage, usage)
	flag.Parse()
	// options without parameters (right now just: -hV)
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	host := flag.Args()[0]
	if err := ping(host); err != nil {
		log.Fatal(err)
	}
}
