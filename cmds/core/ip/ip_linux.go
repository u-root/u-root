// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ip manipulates network addresses, interfaces, routing, and other config.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
	"github.com/vishvananda/netlink"
)

type flags struct {
	family        string
	inet4         bool
	inet6         bool
	details       bool
	stats         bool
	loops         int
	humanReadable bool
	iec           bool
	json          bool
	prettify      bool
}

const ipHelp = `Usage: ip [ OPTIONS ] OBJECT { COMMAND | help }
where  OBJECT := { address |  help | link | monitor | neighbor | neighbour |
				   route | rule | tap | tcpmetrics |
                   token | tunnel | tuntap | vrf | xfrm }
       OPTIONS := { -s[tatistics] | -d[etails] | -r[esolve] |
                    -h[uman-readable] | -iec | -j[son] | -p[retty] |
                    -f[amily] { inet | inet6 | mpls | bridge | link } |
                    -4 | -6 | -M | -B | -0 |
                    -l[oops] { maximum-addr-flush-attempts } | -br[ief] |
                    -o[neline] | -t[imestamp] | -ts[hort] | -b[atch] [filename] |
                    -rc[vbuf] [size] | -n[etns] name | -N[umeric] | -a[ll] |
                    -c[olor]}`

// The language implemented by the standard 'ip' is not super consistent
// and has lots of convenience shortcuts.
// The BNF the standard ip  shows you doesn't show many of these short cuts, and
// it is wrong in other ways.
// For this ip command:.
// The inputs is just the set of args.
// The input is very short -- it's not a program!
// Each token is just a string and we need not produce terminals with them -- they can
// just be the terminals and we can switch on them.
// The cursor is always our current token pointer. We do a simple recursive descent parser
// and accumulate information into a global set of variables. At any point we can see into the
// whole set of args and see where we are. We can indicate at each point what we're expecting so
// that in usage() or recover() we can tell the user exactly what we wanted, unlike the standard ip,
// which just dumps a whole (incorrect) BNF at you when you do anything wrong.
// To handle errors in too few arguments, we just do a recover block. That lets us blindly
// reference the arg[] array without having to check the length everywhere.

// RE: the use of globals. The reason is simple: we parse one command, do it, and quit.
// It doesn't make sense to write this otherwise.
var (
	// Cursor is out next token pointer.
	// The language of this command doesn't require much more.
	f              flags
	cursor         int
	arg            []string
	expectedValues []string
	family         int // netlink.FAMILY_ALL, netlink.FAMILY_V4, netlink.FAMILY_V6
	addrScopes     = map[netlink.Scope]string{
		netlink.SCOPE_UNIVERSE: "global",
		netlink.SCOPE_HOST:     "host",
		netlink.SCOPE_SITE:     "site",
		netlink.SCOPE_LINK:     "link",
		netlink.SCOPE_NOWHERE:  "nowhere",
	}
)

// the pattern:
// at each level parse off arg[0]. If it matches, continue. If it does not, all error with how far you got, what arg you saw,
// and why it did not work out.

func usage() error {
	return fmt.Errorf("this was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], expectedValues)
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&f.family, "f", "", "Specify family (inet, inet6, mpls, link)")
	fs.StringVar(&f.family, "family", "", "Specify family (inet, inet6, mpls, link)")
	fs.BoolVar(&f.inet4, "4", false, "Display IPv4 addresses")
	fs.BoolVar(&f.inet6, "6", false, "Display IPv6 addresses")
	fs.BoolVar(&f.details, "d", false, "Display details")
	fs.BoolVar(&f.details, "details", false, "Display details")
	fs.BoolVar(&f.stats, "s", false, "Display statistics")
	fs.BoolVar(&f.stats, "statistics", false, "Display statistics")
	fs.IntVar(&f.loops, "l", 0, "Set maximum number of attempts to flush all addresses")
	fs.IntVar(&f.loops, "loops", 1, "Set maximum number of attempts to flush all addresses")
	fs.BoolVar(&f.humanReadable, "h", false, "Display timings and sizes in human readable format")
	fs.BoolVar(&f.humanReadable, "humanreadable", false, "Display timings and sizes in human-readable format")
	fs.BoolVar(&f.iec, "iec", false, "Use 1024-based block sizes for human-readable sizes")
	fs.BoolVar(&f.json, "j", false, "Output in JSON format")
	fs.BoolVar(&f.json, "json", false, "Output in JSON format")
	fs.BoolVar(&f.prettify, "p", false, "Make JSON output pretty")
	fs.BoolVar(&f.prettify, "pretty", false, "Make JSON output pretty")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", ipHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	arg = fs.Args()

	if f.family != "" && (f.inet4 || f.inet6) {
		return fmt.Errorf("cannot specify both -f and -4 or -6")
	}

	family = netlink.FAMILY_ALL
	if f.inet6 {
		family = netlink.FAMILY_V6
	} else if f.inet4 {
		family = netlink.FAMILY_V4
	} else if f.family != "" {
		switch f.family {
		case "inet":
			family = netlink.FAMILY_V4
		case "inet6":
			family = netlink.FAMILY_V6
		case "mpls":
			family = netlink.FAMILY_MPLS
		case "link":
			family = netlink.FAMILY_ALL
		}
	}

	expectedValues = []string{"address", "route", "link", "monitor", "neigh", "tunnel", "tuntap", "tap", "tcp_metrics", "tcpmetrics", "vrf", "xfrm", "help"}
	cursor = 0

	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", arg, cursor, expectedValues)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", arg, cursor, expectedValues)
			}
			log.Fatalf("ip: %v", err)
		default:
			log.Fatalf("ip: unexpected panic value: %T(%v)", err, err)
		}

		return
	}()

	// The ip command doesn't actually follow the BNF it prints on error.
	// There are lots of handy shortcuts that people will expect.
	var err error

	c := findPrefix(arg[cursor], expectedValues)
	switch c {
	case "address":
		err = address(out)
	case "link":
		err = link(out)
	case "route":
		err = route(out)
	case "neigh":
		err = neigh(out)
	case "monitor":
		err = monitor(out)
	case "tunnel":
		err = tunnel(out)
	case "tuntap", "tap":
		err = tuntap(out)
	case "tcpmetrics", "tcp_metrics":
		err = tcpMetrics(out)
	case "vrf":
		err = vrf(out)
	case "xfrm":
		err = xfrm(out)
	case "help":
		fmt.Fprint(out, ipHelp)

		return nil
	default:
		err = usage()
	}
	if err != nil {
		return fmt.Errorf("%v: %v", c, err)
	}
	return nil
}

func main() {
	err := run(os.Args, os.Stdout)
	if err != nil {
		log.Fatalf("ip: %v", err)
	}
}
