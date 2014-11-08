// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net"
	"time"
)

var (
	net6       = flag.Bool("6", false, "use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)")
	packetSize = flag.Int("s", 128, "Data size")
	iter       = flag.Int("c", 8, "# iterations")
	intv       = flag.Int("i", 1000, "interval in milliseconds")
	wtf        = flag.Int("w", 100, "wait time in milliseconds")
)

func main() {
	netname := "ip4:icmp"
	flag.Parse()
	interval := time.Duration(*intv)
	waitFor := time.Duration(*wtf) * time.Millisecond
	host := flag.Args()[0]
	// todo: just figure out if it's an ip6 address and go from there.
	if *net6 {
		netname = "ip6:ipv6-icmp"
	}
	msg := make([]byte, *packetSize)

	for i := 0; i < *iter; i++ {
		c, err := net.Dial(netname, host)
		if err != nil {
			log.Fatalf("net.Dial(%v %v) failed: %v", netname, host, err)
		}

		c.SetDeadline(time.Now().Add(waitFor))
		defer c.Close()
		msg[0] = byte(i)
		if _, err := c.Write(msg[:]); err != nil {
			log.Printf("Write failed: %v", err)
		} else {
			c.SetDeadline(time.Now().Add(waitFor))
			if amt, err := c.Read(msg[:]); err == nil {
				log.Printf("%v(%d bytes): %v", i, amt, time.Now())
			} else {
				log.Printf("Read failed: %v", err)
			}
		}
		time.Sleep(time.Millisecond * interval)
	}
}
