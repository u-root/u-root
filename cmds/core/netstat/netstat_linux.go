// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/netstat"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var help = `usage: netstat [-WeenNC] [<Af>] -r         netstat {-h|--help}
       netstat [-WnNaeol] [<Socket> ...]
       netstat { [-WeenNa] -I[<Iface>] | [-eenNa] -i | [-cnNe] | -s [-6tuw] } [delay]

        -r, --route             display routing table
        -I, --interface=<Iface> display interface table for <Iface>
        -i, --interfaces        display interface table
        -g, --groups            display multicast group memberships
        -s, --statistics        display networking statistics (like SNMP)

        -W, --wide              don't truncate IP addresses
        -n, --numeric           don't resolve names
        --numeric-hosts         don't resolve host names
        --numeric-ports         don't resolve port names
        --numeric-users         don't resolve user names
        -N, --symbolic          resolve hardware names
        -e, --extend            display other/more information
        -p, --programs          display PID/Program name for sockets
        -o, --timers            display timers
        -c, --continuous        continuous listing

        -l, --listening         display listening server sockets
        -a, --all               display all sockets (default: connected)
        -C, --cache             display routing cache instead of FIB

  <Socket>={-t|--tcp} {-u|--udp} {-U|--udplite} {-w|--raw} {-x|--unix}
  <AF>=Use '-6|-4' or '-A <af>' or '--<af>'; default: inet
  List of possible address families (which support routing):
    inet (DARPA Internet) inet6 (IPv6)`

func printHelp() {
	fmt.Printf("%s\n", help)
}

// cmd groups the available functionality of netstat.
// More information is given in the main function down below.
type cmd struct {
	out io.Writer

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
	af   string

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

var errMutualExcludeFlags = errors.New("only one of route, interfaces, groups, statistics or interface allowed")

func (c cmd) run() error {
	afs := make([]netstat.AddressFamily, 0)

	if !xorFlags(c.route, c.interfaces, c.groups, c.stats, c.iface != "") {
		return errMutualExcludeFlags
	}

	// use IPv4 as default if no address family is given
	if !c.ipv4 && !c.ipv6 {
		c.ipv4 = true
	}

	// In case we want to print IP Sockets and give no further information, all protocols shall
	// be printet. The default value of the protocol flags is false, but for that we need to override
	// them to be true.
	if (c.ipv4 || c.ipv6 || c.stats) &&
		!(c.tcp || c.udp || c.udpL || c.raw || c.unix) {
		c.tcp = true
		c.udp = true
		c.udpL = true
		c.raw = true
		c.unix = true
	}

	socks, err := evalProtocols(
		c.tcp,
		c.udp,
		c.udpL,
		c.raw,
		c.unix,
		c.ipv4,
		c.ipv6,
	)
	if err != nil {
		return err
	}

	// numeric groups the format functionality of numeric-hosts, numeric-ports and numeric-users.
	// It overrides the other numeric format flags.
	if c.numeric {
		c.numHost = true
		c.numPorts = true
		c.numUsers = true
	}

	// Evaluate for route cache for IPv6
	if c.routecache && c.route && !c.ipv6 {
		return netstat.ErrRouteCacheIPv6only
	}

	// Set up format flags for route listing and socket listing
	outflags := netstat.FmtFlags{
		Extend:    c.extend,
		Wide:      c.wide,
		NumHosts:  c.numHost,
		NumPorts:  c.numPorts,
		NumUsers:  c.numUsers,
		ProgNames: c.programs,
		Timer:     c.timers,
		Symbolic:  c.symbolic,
	}

	// Set up output generator for route and socket listing
	outfmts, err := netstat.NewOutput(outflags)
	if err != nil {
		return err
	}

	if c.route {
		if c.ipv4 {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if c.ipv6 {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			for {
				str, err := af.RoutesFormatString(c.routecache)
				if err != nil {
					return err
				}
				fmt.Fprintf(c.out, "%s\n", str)
				if !c.contin {
					break
				}
				af.ClearOutput()
				time.Sleep(2 * time.Second)
			}
		}

		return err
	}

	if c.interfaces {
		return netstat.PrintInterfaceTable(c.iface, c.contin, c.out)
	}

	if c.iface != "" {
		return netstat.PrintInterfaceTable(c.iface, c.contin, c.out)
	}

	if c.groups {
		return netstat.PrintMulticastGroups(true, true, c.out)
	}

	if c.stats {
		if c.ipv4 {
			afs = append(afs, netstat.NewAddressFamily(false, outfmts))
		}

		if c.ipv6 {
			afs = append(afs, netstat.NewAddressFamily(true, outfmts))
		}

		for _, af := range afs {
			if err := af.PrintStatistics(c.out); err != nil {
				return err
			}
		}
		return nil
	}

	for {
		for _, sock := range socks {
			str, err := sock.SocketsString(c.listening, c.all, outfmts)
			if err != nil {
				return err
			}
			fmt.Fprintf(c.out, "%s\n", str)
		}
		if !c.contin {
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
			return nil, fmt.Errorf("evalProtocols creating TCP4 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if tcp && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_TCP6)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating TCP6 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if udp && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_UDP)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating UDP4 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if udp && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_UDP6)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating UDP6 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if udpl && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_UDPL)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating UDPL4 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if udpl && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_UDPL6)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating UDPL6 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if raw && ipv4 {
		t, err := netstat.NewSocket(netstat.PROT_RAW)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating RAW4 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if raw && ipv6 {
		t, err := netstat.NewSocket(netstat.PROT_RAW6)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating RAW6 socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	if unix {
		t, err := netstat.NewSocket(netstat.PROT_UNIX)
		if err != nil {
			return nil, fmt.Errorf("evalProtocols creating UNIX socket: %w", err)
		}
		retProtos = append(retProtos, t)
	}

	return retProtos, nil
}

func command(out io.Writer, args []string) *cmd {
	var c cmd
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.BoolVar(&c.route, "route", false, "display routing table")
	fs.BoolVar(&c.route, "r", false, "display routing table")

	fs.BoolVar(&c.interfaces, "interfaces", false, "display interface table")
	fs.BoolVar(&c.interfaces, "i", false, "display interface table")

	fs.StringVar(&c.iface, "interface", "", "Display interface table for interface <if>")
	fs.StringVar(&c.iface, "I", "", "Display interface table for interface <if>")

	fs.BoolVar(&c.groups, "groups", false, "display multicast group memberships")
	fs.BoolVar(&c.groups, "g", false, "display multicast group memberships")

	fs.BoolVar(&c.stats, "statistics", false, "display networking statistics (like SNMP)")
	fs.BoolVar(&c.stats, "s", false, "display networking statistics (like SNMP)")

	fs.BoolVar(&c.tcp, "tcp", false, "Print TCP sockets")
	fs.BoolVar(&c.tcp, "t", false, "Print TCP sockets")

	fs.BoolVar(&c.udp, "udp", false, "Print UDP sockets")
	fs.BoolVar(&c.udp, "u", false, "Print UDP sockets")

	fs.BoolVar(&c.udpL, "udplite", false, "Print UDPlite sockets")
	fs.BoolVar(&c.udpL, "U", false, "Print UDPlite sockets")

	fs.BoolVar(&c.raw, "raw", false, "Print IPv4/IPv6 RAW sockets")
	fs.BoolVar(&c.raw, "w", false, "Print IPv4/IPv6 RAW sockets")

	fs.BoolVar(&c.unix, "unix", false, "Print UNIX sockets")
	fs.BoolVar(&c.unix, "x", false, "Print UNIX sockets")

	fs.BoolVar(&c.ipv4, "4", false, "IPv4 fs. default: false")
	fs.BoolVar(&c.ipv6, "6", false, "IPv6 fs. default: false")
	fs.BoolVar(&c.ipv4, "inet", false, "IPv4 fs. default: false")     // alternative af setting, see help text
	fs.BoolVar(&c.ipv6, "inet6", false, "IPv6 fs. default: false")    // alternative af setting, see help text
	fs.StringVar(&c.af, "A", "", "Address family, 'inet' or 'inet6'") // alternative af setting, see help text

	fs.BoolVar(&c.routecache, "cache", false, "")
	fs.BoolVar(&c.routecache, "C", false, "")

	fs.BoolVar(&c.wide, "wide", false, "don't truncate IP addresses")
	fs.BoolVar(&c.wide, "W", false, "don't truncate IP addresses")

	fs.BoolVar(&c.numeric, "numeric", false, "don't resolve names")
	fs.BoolVar(&c.numeric, "n", false, "don't resolve names")

	fs.BoolVar(&c.numHost, "numeric-hosts", false, "don't resolve host names")
	fs.BoolVar(&c.numPorts, "numeric-ports", false, "don't resolve port names")
	fs.BoolVar(&c.numUsers, "numeric-users", false, "don't resolve user names")

	fs.BoolVar(&c.symbolic, "symbolic", false, "resolve hardware names")
	fs.BoolVar(&c.symbolic, "N", false, "resolve hardware names")

	fs.BoolVar(&c.extend, "extend", false, "display other/more information")
	fs.BoolVar(&c.extend, "e", false, "display other/more information")

	fs.BoolVar(&c.programs, "programs", false, "display PID/Program name for sockets")
	fs.BoolVar(&c.programs, "p", false, "display PID/Program name for sockets")

	fs.BoolVar(&c.timers, "timers", false, "display timers")
	fs.BoolVar(&c.timers, "o", false, "display timers")

	fs.BoolVar(&c.contin, "continuous", false, "continuous listing")
	fs.BoolVar(&c.contin, "c", false, "continuous listing")

	fs.BoolVar(&c.listening, "listening", false, "display listening server sockets")
	fs.BoolVar(&c.listening, "l", false, "display listening server sockets")

	fs.BoolVar(&c.all, "all", false, "display all sockets (default: connected)")
	fs.BoolVar(&c.all, "a", false, "display all sockets (default: connected)")

	fs.Usage = printHelp
	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	// Apply alternative address family setting
	switch c.af {
	case "inet":
		c.ipv4 = true
	case "inet6":
		c.ipv6 = true
	}

	c.out = out
	return &c
}

func main() {
	switch err := command(os.Stdout, os.Args).run(); err {
	case nil:
	case errMutualExcludeFlags:
		printHelp()
		log.Print(err)
		os.Exit(2)
	default:
		log.Fatal(err)
	}
}
