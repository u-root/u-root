// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netstat"
)

var (
	interfacesFlag = flag.BoolP("interfaces", "i", false, "display interface table")
	ifFlag         = flag.StringP("interface", "I", "", "Display interface table for interface <if>")
	groupsFlag     = flag.BoolP("groups", "g", false, "display multicast group memberships")
	statsFlag      = flag.BoolP("statistics", "s", false, "display networking statistics (like SNMP)")
	// AF Flags
	ipv4Flag = flag.BoolP("4", "4", false, "IPv4 flag. default: true")
	ipv6Flag = flag.BoolP("6", "6", false, "IPv6 flag. default: false")
	// Format flags
	wideFlag      = flag.BoolP("wide", "W", false, "don't truncate IP addresses")
	numericFlag   = flag.BoolP("numeric", "n", false, "don't resolve names")
	numHostFlag   = flag.Bool("numeric-hosts", false, "don't resolve host names")
	numPortsFlag  = flag.Bool("numeric-ports", false, "don't resolve port names")
	numUsersFlag  = flag.Bool("numeric-users", false, "don't resolve user names")
	symbolicFlag  = flag.BoolP("symbolic", "N", false, "resolve hardware names")
	extendFlag    = flag.BoolP("extend", "e", false, "display other/more information")
	programsFlag  = flag.BoolP("programs", "p", false, "display PID/Program name for sockets")
	timersFlag    = flag.BoolP("timers", "o", false, "display timers")
	continFlag    = flag.BoolP("continuous", "c", false, "continuous listing")
)

func evalFlags() error {
	flag.Parse()

	afs := make([]netstat.AddressFamily, 0)
	// Can't use default capability of pflags package, have to determine it like this
	// to keep same usage as original netstat tool.
	if !*ipv4Flag && !*ipv6Flag {
		*ipv4Flag = true

	}
	// Evaluate numeric flags
	if *numericFlag {
		*numHostFlag = true
		*numPortsFlag = true
		*numUsersFlag = true
	}

	// Set up format flags for route listing and socket listing
	outflags := netstat.FmtFlags{
		Extend:    *extendFlag,
		Wide:      *wideFlag,
		NumHosts:  *numHostFlag,
		NumPorts:  *numPortsFlag,
		NumUsers:  *numUsersFlag,
		ProgNames: *programsFlag,
		Timer:     *timersFlag,
		Symbolic:  *symbolicFlag,
	}

	// Set up output generator for route and socket listing
	outfmts, err := netstat.NewOutput(outflags)
	if err != nil {
		log.Fatal(err)
	}

	if *interfacesFlag {
		if err := netstat.PrintInterfaceTable(*ifFlag, *continFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *ifFlag != "" {
		if err := netstat.PrintInterfaceTable(*ifFlag, *continFlag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *groupsFlag {
		if err := netstat.PrintMulticastGroups(*ipv4Flag, *ipv6Flag); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *statsFlag {
		if *ipv4Flag {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if *ipv6Flag {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			if err := af.PrintStatistics(); err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}

	return nil
}

func main() {
	if err := evalFlags(); err != nil {
		log.Fatal(err)
	}
}
