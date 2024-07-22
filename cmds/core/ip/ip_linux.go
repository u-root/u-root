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
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type flags struct {
	family         string
	families       []int
	inet4          bool
	inet6          bool
	bridge         bool
	mpls           bool
	link           bool
	details        bool
	stats          bool
	loops          int
	humanReadable  bool
	iec            bool
	json           bool
	prettify       bool
	brief          bool
	resolve        bool
	color          string
	rcvBuf         string
	timeStamp      bool
	timeStampShort bool
	all            bool
	numeric        bool
	batch          string
	force          bool
	oneline        bool
	netns          string
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
	families       []int
	family         int
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

func parseFlags(args []string, out io.Writer) (*netlink.Handle, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&f.family, "f", "", "Specify family (inet, inet6, mpls, link)")
	fs.StringVar(&f.family, "family", "", "Specify family (inet, inet6, mpls, link)")
	fs.BoolVar(&f.resolve, "r", false, "Use system resolver to display DNS names")
	fs.BoolVar(&f.resolve, "resolve", false, "Use system resolver to display DNS names")
	fs.BoolVar(&f.inet4, "4", false, "Set protocol family to inet")
	fs.BoolVar(&f.inet6, "6", false, "Set protocol family to inet6")
	fs.BoolVar(&f.bridge, "B", false, "Set protocol family to bridge")
	fs.BoolVar(&f.mpls, "M", false, "Set protocol family to mpls")
	fs.BoolVar(&f.link, "0", false, "Set protocol family to link")
	fs.BoolVar(&f.details, "d", false, "Display details")
	fs.BoolVar(&f.details, "details", false, "Display details")
	fs.BoolVar(&f.stats, "s", false, "Display statistics")
	fs.BoolVar(&f.stats, "statistics", false, "Display statistics")
	fs.IntVar(&f.loops, "l", 0, "Set maximum number of attempts to flush all addresses")
	fs.IntVar(&f.loops, "loops", 1, "Set maximum number of attempts to flush all addresses")
	fs.BoolVar(&f.humanReadable, "h", false, "Display timings and sizes in human readable format")
	fs.BoolVar(&f.humanReadable, "humanreadable", false, "Display timings and sizes in human-readable format")
	fs.BoolVar(&f.iec, "iec", false, "Use 1024-based block sizes for human-readable sizes")
	fs.BoolVar(&f.brief, "br", false, "Brief output")
	fs.BoolVar(&f.brief, "brief", false, "Brief output")
	fs.BoolVar(&f.json, "j", false, "Output in JSON format")
	fs.BoolVar(&f.json, "json", false, "Output in JSON format")
	fs.BoolVar(&f.prettify, "p", false, "Make JSON output pretty")
	fs.BoolVar(&f.prettify, "pretty", false, "Make JSON output pretty")
	fs.StringVar(&f.color, "c", "", "Use color output")
	fs.StringVar(&f.color, "color", "", "Use color output")
	fs.StringVar(&f.rcvBuf, "rc", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.StringVar(&f.rcvBuf, "rcvbuf", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.BoolVar(&f.timeStamp, "t", false, "Display time stamps")
	fs.BoolVar(&f.timeStamp, "timestamp", false, "Display time stamps")
	fs.BoolVar(&f.timeStampShort, "ts", false, "Display short time stamps")
	fs.BoolVar(&f.timeStampShort, "tshort", false, "Display short time stamps")
	fs.BoolVar(&f.all, "a", false, "Display all information")
	fs.BoolVar(&f.all, "all", false, "Display all information")
	fs.BoolVar(&f.numeric, "N", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.BoolVar(&f.numeric, "numeric", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.StringVar(&f.batch, "b", "", "Read commands from a file")
	fs.StringVar(&f.batch, "batch", "", "Read commands from a file")
	fs.BoolVar(&f.force, "force", false, "Don't terminate ip on errors in batch mode.  If there were any errors during execution of the commands, the application return code will be non zero.")
	fs.BoolVar(&f.oneline, "o", false, "Output each record on a single line")
	fs.BoolVar(&f.oneline, "oneline", false, "Output each record on a single line")
	fs.StringVar(&f.netns, "n", "", "Switch to network namespace")
	fs.StringVar(&f.netns, "netns", "", "Switch to network namespace")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", ipHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	arg = fs.Args()

	family = netlink.FAMILY_ALL

	if f.inet4 {
		family = netlink.FAMILY_V4
		families = append(families, netlink.FAMILY_V4)
	}
	if f.inet6 {
		family = netlink.FAMILY_V6
		families = append(families, netlink.FAMILY_V6)
	}

	if f.mpls {
		return nil, fmt.Errorf("protocol family MPLS is not yet supported")
	}

	if f.bridge {
		return nil, fmt.Errorf("protocol family bridge is not yet supported")
	}

	if f.link {
		family = netlink.FAMILY_ALL
		families = append(families, netlink.FAMILY_ALL)
	}

	if f.resolve {
		return nil, fmt.Errorf("resolving DNS names is unsupported")
	}

	if f.color != "" {
		return nil, fmt.Errorf("color output is unsupported")
	}

	if f.oneline {
		return nil, fmt.Errorf("outputting each record on a single line is unsupported")
	}

	if f.batch != "" {
		return nil, fmt.Errorf("batch mode is unsupported")
	}

	if f.force {
		return nil, fmt.Errorf("force mode is unsupported")
	}

	var (
		err    error
		handle *netlink.Handle
	)

	if f.netns != "" {
		nsHandle, err := netns.GetFromName(f.netns)
		if err != nil {
			return nil, fmt.Errorf("failed to find network namespace %q: %v", f.netns, err)
		}
		defer nsHandle.Close()

		handle, err = netlink.NewHandleAt(nsHandle, families...)
		if err != nil {
			return nil, fmt.Errorf("failed to create netlink handle in network namespace %q: %v", f.netns, err)
		}
	} else {
		handle, err = netlink.NewHandle(families...)
		if err != nil {
			return nil, fmt.Errorf("failed to create netlink handle: %v", err)
		}
	}

	if f.rcvBuf != "" {
		bufSize, err := strconv.Atoi(f.rcvBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rcvbuf flag: %v", err)
		}

		handle.SetSocketReceiveBufferSize(bufSize, true)
	}

	return handle, nil
}

type cmd struct {
	out    io.Writer
	handle *netlink.Handle
}

func (cmd cmd) run() error {
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
		err = cmd.address()
	case "link":
		err = cmd.link()
	case "route":
		err = cmd.route()
	case "neigh":
		err = cmd.neigh()
	case "monitor":
		err = cmd.monitor()
	case "tunnel":
		err = cmd.tunnel()
	case "tuntap", "tap":
		err = cmd.tuntap()
	case "tcpmetrics", "tcp_metrics":
		err = cmd.tcpMetrics()
	case "vrf":
		err = cmd.vrf()
	case "xfrm":
		err = cmd.xfrm()
	case "help":
		fmt.Fprint(cmd.out, ipHelp)

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
	handle, err := parseFlags(os.Args, os.Stdout)
	if err != nil {
		log.Fatalf("ip: %v", err)
	}

	cmd := cmd{
		out:    os.Stdout,
		handle: handle,
	}

	err = cmd.run()
	if err != nil {
		log.Fatalf("ip: %v", err)
	}
}
