// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/netstat"
)

// Flags groups the available functionality of netstat.
// More information is given in the main function down below.
type Flags struct {
	// Info source flags
	route      bool
	interfaces bool
	iface      string
	groups     bool
	stats      bool

	// Socket flags
	tcp  bool
	udp  bool
	udpL bool
	raw  bool
	unix bool
	// AF Flags
	ipv4 bool
	ipv6 bool

	// Route type flag
	routecache bool

	// Format flags
	wide      bool
	numeric   bool
	numHost   bool
	numPorts  bool
	numUsers  bool
	symbolic  bool
	extend    bool
	programs  bool
	timers    bool
	contin    bool
	listening bool
	all       bool
}

func main() {
	f := Flags{}
	flag.BoolVarP(&f.route, "route", "r", false, "display routing table")
	flag.BoolVarP(&f.interfaces, "interfaces", "i", false, "display interface table")
	flag.StringVarP(&f.iface, "interface", "I", "", "Display interface table for interface <if>")
	flag.BoolVarP(&f.groups, "groups", "g", false, "display multicast group memberships")
	flag.BoolVarP(&f.stats, "statistics", "s", false, "display networking statistics (like SNMP)")

	flag.BoolVarP(&f.tcp, "tcp", "t", false, "Print TCP sockets")
	flag.BoolVarP(&f.udp, "udp", "u", false, "Print UDP sockets")
	flag.BoolVarP(&f.udpL, "udplite", "U", false, "Print UDPlite sockets")
	flag.BoolVarP(&f.raw, "raw", "w", false, "Print IPv4/IPv6 RAW sockets")
	flag.BoolVarP(&f.unix, "unix", "x", false, "Print UNIX sockets")

	flag.BoolVarP(&f.ipv4, "4", "4", false, "IPv4 flag. default: true")
	flag.BoolVarP(&f.ipv6, "6", "6", false, "IPv6 flag. default: false")

	flag.BoolVarP(&f.routecache, "cache", "C", false, "")

	flag.BoolVarP(&f.wide, "wide", "W", false, "don't truncate IP addresses")
	flag.BoolVarP(&f.numeric, "numeric", "n", false, "don't resolve names")
	flag.BoolVar(&f.numHost, "numeric-hosts", false, "don't resolve host names")
	flag.BoolVar(&f.numPorts, "numeric-ports", false, "don't resolve port names")
	flag.BoolVar(&f.numUsers, "numeric-users", false, "don't resolve user names")
	flag.BoolVarP(&f.symbolic, "symbolic", "N", false, "resolve hardware names")
	flag.BoolVarP(&f.extend, "extend", "e", false, "display other/more information")
	flag.BoolVarP(&f.programs, "programs", "p", false, "display PID/Program name for sockets")
	flag.BoolVarP(&f.timers, "timers", "o", false, "display timers")
	flag.BoolVarP(&f.contin, "continuous", "c", false, "continuous listing")
	flag.BoolVarP(&f.listening, "listening", "l", false, "display listening server sockets")
	flag.BoolVarP(&f.all, "all", "a", false, "display all sockets (default: connected)")

	flag.Parse()
	if err := run(f, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(f Flags, out io.Writer) error {
	afs := make([]netstat.AddressFamily, 0)

	// Validate info source flags
	// none or one allowed to be set
	if !xorFlags(f.route, f.interfaces, f.groups, f.stats, f.iface != "") {
		flag.Usage()
		return nil
	}

	// Can't use default capability of pflags package, have to determine it like this
	// to keep same usage as original netstat tool.
	if !f.ipv4 && !f.ipv6 {
		f.ipv4 = true
	}

	// In case we want to print IP Sockets and give no further information, all protocols shall
	// be printet. The default value of the protocol flags is false, but for that we need to override
	// them to be true.
	if (f.ipv4 || f.ipv6 || f.stats) &&
		!(f.tcp || f.udp || f.udpL || f.raw || f.unix) {
		f.tcp = true
		f.udp = true
		f.udpL = true
		f.raw = true
		f.unix = true
	}

	socks, err := evalProtocols(
		f.tcp,
		f.udp,
		f.udpL,
		f.raw,
		f.unix,
		f.ipv4,
		f.ipv6,
	)
	if err != nil {
		return err
	}

	// numeric groups the format functionality of numeric-hosts, numeric-ports and numeric-users.
	// It overrides the other numeric format flags.
	if f.numeric {
		f.numHost = true
		f.numPorts = true
		f.numUsers = true
	}

	// Evaluate for route cache for IPv6
	if f.routecache && f.route && !f.ipv6 {
		return netstat.ErrRouteCacheIPv6only
	}

	// Set up format flags for route listing and socket listing
	outflags := netstat.FmtFlags{
		Extend:    f.extend,
		Wide:      f.wide,
		NumHosts:  f.numHost,
		NumPorts:  f.numPorts,
		NumUsers:  f.numUsers,
		ProgNames: f.programs,
		Timer:     f.timers,
		Symbolic:  f.symbolic,
	}

	// Set up output generator for route and socket listing
	outfmts, err := netstat.NewOutput(outflags)
	if err != nil {
		return err
	}

	if f.route {
		if f.ipv4 {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if f.ipv6 {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			for {
				str, err := af.RoutesFormatString(f.routecache)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%s\n", str)
				if !f.contin {
					break
				}
				af.ClearOutput()
				time.Sleep(2 * time.Second)
			}

		}

		return err
	}

	if f.interfaces {
		return netstat.PrintInterfaceTable(f.iface, f.contin, out)
	}

	if f.iface != "" {
		return netstat.PrintInterfaceTable(f.iface, f.contin, out)
	}

	if f.groups {
		return netstat.PrintMulticastGroups(f.ipv4, f.ipv6, out)
	}

	if f.stats {
		if f.ipv4 {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if f.ipv6 {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			if err := af.PrintStatistics(out); err != nil {
				return err
			}
		}
		return nil
	}

	for {
		for _, sock := range socks {
			str, err := sock.SocketsString(f.listening, f.all, outfmts)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "%s\n", str)
		}
		if !f.contin {
			break
		}
		outfmts.Builder.Reset()
		time.Sleep(2 * time.Second)
	}

	return nil
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
