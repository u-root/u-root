// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ntpdate uses NTP to adjust the system clock.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/beevik/ntp"
)

var (
	config  = flag.String("config", "/etc/ntp.conf", "NTP config file.")
	verbose = flag.Bool("verbose", false, "Verbose output")
	debug   = func(string, ...interface{}) {}
)

const (
	fallback = "time.google.com"
)

func parseServers(r *bufio.Reader) []string {
	var uri []string
	var l string
	var err error

	debug("Reading config file")
	for err == nil {
		// This handles the case where the last line doesn't end in \n
		l, err = r.ReadString('\n')
		debug("%v", l)
		if w := strings.Fields(l); len(w) > 1 && w[0] == "server" {
			// We look only for the server lines, we ignore options like iburst
			// TODO(ganshun): figure out what options we want to support.
			uri = append(uri, w[1])
		}
	}

	return uri
}

func getTime(servers []string) (t time.Time, err error) {
	for _, s := range servers {
		debug("Getting time from %v", s)
		if t, err = ntp.Time(s); err == nil {
			// Right now we return on the first valid time.
			// We can implement better heuristics here.
			debug("Got time %v", t)
			return
		}
		debug("Error getting time: %v", err)
	}

	err = fmt.Errorf("unable to get any time from servers %v", servers)
	return
}

func main() {
	var servers []string
	flag.Parse()
	if *verbose {
		debug = log.Printf
	}

	debug("Reading NTP servers from config file: %v", *config)
	f, err := os.Open(*config)
	if err == nil {
		defer f.Close()
		servers = parseServers(bufio.NewReader(f))
		debug("Found %v servers", len(servers))
	} else {
		log.Printf("Unable to open config file: %v\nFalling back to : %v", err, fallback)
		servers = []string{fallback}
	}

	t, err := getTime(servers)
	if err != nil {
		log.Fatalf("Unable to get time: %v", err)
	}

	tv := syscall.NsecToTimeval(t.UnixNano())
	if err = syscall.Settimeofday(&tv); err != nil {
		log.Fatalf("Unable to set system time: %v", err)
	}
}
