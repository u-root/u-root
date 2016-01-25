// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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
	iter       = flag.Int64("c", 0, "# iterations")
	intv       = flag.Int("i", 1000, "interval in milliseconds")
	version    = flag.Bool("V", false, "version")
	wtf        = flag.Int("w", 100, "wait time in milliseconds")
	help       = flag.Bool("h", false, "help")
)

func usage() {
	fmt.Fprintf(os.Stdout, "ping [-hV] [-c count] [-i interval] [-s packetsize [-w deadline] destination\n")
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

func main() {
	flag.Parse()

	// options without parameters (right now just: -hV)
	if flag.NArg() < 1 {
		optwithoutparam()
	}

	var i int64 = 1
	netname := "ip4:icmp"
	interval := time.Duration(*intv)
	waitFor := time.Duration(*wtf) * time.Millisecond
	host := flag.Args()[0]
	// todo: just figure out if it's an ip6 address and go from there.
	if *net6 {
		netname = "ip6:ipv6-icmp"
	}
	msg := make([]byte, *packetSize)

	// ping needs to run forever, except if '*iter' is not zero
	for {
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
		if *iter != 0 {
			if i < *iter {
				i++
			} else {
				break
			}
		}
	}
}
