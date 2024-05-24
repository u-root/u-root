// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netstat"
)

var (
	// Info source flags
	routeFlag      = flag.BoolP("route", "r", false, "display routing table")
	interfacesFlag = flag.BoolP("interfaces", "i", false, "display interface table")
	ifFlag         = flag.StringP("interface", "I", "", "Display interface table for interface <if>")
	groupsFlag     = flag.BoolP("groups", "g", false, "display multicast group memberships")
	statsFlag      = flag.BoolP("statistics", "s", false, "display networking statistics (like SNMP)")

	// Socket flags
	tcpFlag  = flag.BoolP("tcp", "t", false, "TCP")
	udpFlag  = flag.BoolP("udp", "u", false, "UDP")
	udpLFlag = flag.BoolP("udplite", "U", false, "UDPlite")
	rawFlag  = flag.BoolP("raw", "w", false, "RAW")
	unixFlag = flag.BoolP("unix", "x", false, "UNIX")

	// AF Flags
	ipv4Flag = flag.BoolP("4", "4", false, "IPv4 flag. default: true")
	ipv6Flag = flag.BoolP("6", "6", false, "IPv6 flag. default: false")

	// Route type flag
	routecacheFalg = flag.BoolP("cache", "C", false, "")

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
	listeningFlag = flag.BoolP("listening", "l", false, "display listening server sockets")
	allFlag       = flag.BoolP("all", "a", false, "display all sockets (default: connected)")
)

func run() error {
	flag.Parse()

	afs := make([]netstat.AddressFamily, 0)

	// Validate info source flags
	// none or one allowed allowed to be set
	if !xorFlags(*routeFlag, *interfacesFlag, *groupsFlag, *statsFlag, *ifFlag != "") {
		flag.Usage()
		return nil
	}

	// Can't use default capability of pflags package, have to determine it like this
	// to keep same usage as original netstat tool.
	if !*ipv4Flag && !*ipv6Flag {
		*ipv4Flag = true

	}

	// Why do I have to write it like that? It's ugly....MOM!
	if (*ipv4Flag || *ipv6Flag || *statsFlag) &&
		!(*tcpFlag || *udpFlag || *udpLFlag || *rawFlag || *unixFlag) {
		*tcpFlag = true
		*udpFlag = true
		*udpLFlag = true
		*rawFlag = true
		*unixFlag = true
	}

	socks, err := evalProtocols(*tcpFlag, *udpFlag, *udpLFlag, *rawFlag, *unixFlag, *ipv4Flag, *ipv6Flag)
	if err != nil {
		return err
	}

	// Evaluate numeric flags
	if *numericFlag {
		*numHostFlag = true
		*numPortsFlag = true
		*numUsersFlag = true
	}

	// Evaluate for route cache for IPv6
	if *routecacheFalg && *routeFlag && !*ipv6Flag {
		return netstat.ErrRouteCacheIPv6only
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
		return err
	}

	if *routeFlag {
		if *ipv4Flag {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if *ipv6Flag {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			for {
				str, err := af.RoutesFormatString(*routecacheFalg)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", str)
				if !*continFlag {
					break
				}
				af.ClearOutput()
				time.Sleep(2 * time.Second)
			}

		}

		return err
	}

	if *interfacesFlag {
		return netstat.PrintInterfaceTable(*ifFlag, *continFlag)
	}

	if *ifFlag != "" {
		return netstat.PrintInterfaceTable(*ifFlag, *continFlag)
	}

	if *groupsFlag {
		return netstat.PrintMulticastGroups(*ipv4Flag, *ipv6Flag)
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
				return err
			}
		}
		return nil
	}

	for {
		for _, sock := range socks {
			str, err := sock.SocketsString(*listeningFlag, *allFlag, outfmts)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", str)
		}
		if !*continFlag {
			break
		}
		outfmts.Builder.Reset()
		time.Sleep(2 * time.Second)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func xorFlags(flags ...bool) bool {
	c := 0
	for _, flag := range flags {
		if flag {
			c++
		}
	}

	return c <= 1
}

func evalProtocols(tcp, udp, udpl, raw, unix, ipv4, ipv6 bool) ([]netstat.Socket, error) {
	retProtos := make([]netstat.Socket, 0)

	if tcp && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_TCP)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if tcp && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_TCP6)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if udp && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_UDP)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if udp && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_UDP6)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if udpl && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_UDPL)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if udpl && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_UDPL6)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if raw && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_RAW)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if raw && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_RAW6)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	if unix {
		t, err := netstat.NewSocket(netstat.PROT_UNIX)
		if err != nil {
			return nil, err
		}
		retProtos = append(retProtos, t)
	}

	return retProtos, nil
}
