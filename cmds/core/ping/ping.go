// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Send icmp packets to a server to test network connectivity.
//
// Synopsis:
//     ping [-hV] [-c COUNT] [-i INTERVAL] [-s PACKETSIZE] [-w DEADLINE] DESTINATION
//
// Options:
//     -6: use ipv6 (ip6:ipv6-icmp)
//     -s: data size (default: 64)
//     -c: # iterations, 0 to run forever (default)
//     -i: interval in milliseconds (default: 1000)
//     -V: version
//     -w: wait time in milliseconds (default: 100)
//     -a: Audible rings a bell when a packet is received
//     -h: help
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	net6       = flag.Bool("6", false, "use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)")
	packetSize = flag.Int("s", 64, "Data size")
	iter       = flag.Uint64("c", 0, "# iterations")
	intv       = flag.Int("i", 1000, "interval in milliseconds")
	version    = flag.Bool("V", false, "version")
	wtf        = flag.Int("w", 100, "wait time in milliseconds")
	audible    = flag.Bool("a", false, "Audible rings a bell when a packet is received")
)

const (
	ICMP_TYPE_ECHO_REQUEST             = 8
	ICMP_TYPE_ECHO_REPLY               = 0
	ICMP_ECHO_REPLY_HEADER_IPV4_OFFSET = 20
	ICMP_ECHO_REPLY_HEADER_IPV6_OFFSET = 40
)

func usage() {
	fmt.Fprintf(os.Stdout, "ping [-V] [-c count] [-i interval] [-s packetsize] [-w deadline] destination\n")
	os.Exit(0)
}

func showversion() {
	fmt.Fprintf(os.Stdout, "ping utility, Uroot version\n")
	os.Exit(0)
}

func optwithoutparam() {
	if *version {
		showversion()
	}
	// if we reach this point, invalid or help (-h) gets the same result
	usage()
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

func ping1(netname string, host string, i uint64) (string, error) {
	c, derr := net.Dial(netname, host)
	if derr != nil {
		return "", fmt.Errorf("net.Dial(%v %v) failed: %v", netname, host, derr)
	}
	defer c.Close()

	// Send ICMP Echo Request
	waitFor := time.Duration(*wtf) * time.Millisecond
	c.SetDeadline(time.Now().Add(waitFor))
	msg := make([]byte, *packetSize)
	msg[0] = ICMP_TYPE_ECHO_REQUEST
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
	amt, rerr := c.Read(rmsg[:])
	if rerr != nil {
		return "", fmt.Errorf("read failed: %v", rerr)
	}
	latency := time.Since(before)
	if (rmsg[0] & 0x0F) == 6 {
		rmsg = rmsg[ICMP_ECHO_REPLY_HEADER_IPV6_OFFSET:]
	} else {
		rmsg = rmsg[ICMP_ECHO_REPLY_HEADER_IPV4_OFFSET:]
	}
	if rmsg[0] != ICMP_TYPE_ECHO_REPLY {
		return "", fmt.Errorf("bad ICMP echo reply type: %v", msg[0])
	}
	cks := binary.BigEndian.Uint16(rmsg[2:])
	binary.BigEndian.PutUint16(rmsg[2:], 0)
	if cks != cksum(rmsg) {
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

func main() {
	flag.Parse()

	// options without parameters (right now just: -hV)
	if flag.NArg() < 1 {
		optwithoutparam()
	}
	if *packetSize < 8 {
		log.Fatalf("packet size too small (must be >= 8): %v", *packetSize)
	}

	netname := "ip4:icmp"
	interval := time.Duration(*intv)
	host := flag.Args()[0]
	// todo: just figure out if it's an ip6 address and go from there.
	if *net6 {
		netname = "ip6:ipv6-icmp"
	}

	// ping needs to run forever, except if '*iter' is not zero
	var i uint64
	for i = 1; *iter == 0 || i <= *iter; i++ {
		msg, err := ping1(netname, host, i)
		if err != nil {
			log.Fatalf("ping failed: %v", err)
		}
		if *audible {
			msg = "\a" + msg
		}
		log.Print(msg)
		time.Sleep(time.Millisecond * interval)
	}
}
