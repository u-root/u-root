// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ip manipulates network addresses, interfaces, routing, and other config.
package main

import (
	"bufio"
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
)

// the pattern:
// at each level parse off arg[0]. If it matches, continue. If it does not, all error with how far you got, what arg you saw,
// and why it did not work out.

func (cmd cmd) usage() error {
	return fmt.Errorf("this was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		cmd.args[0:cmd.cursor], cmd.args[cmd.cursor:], cmd.args[cmd.cursor], cmd.expectedValues)
}

func parseFlags(args []string, out io.Writer) (*cmd, error) {
	cmd := cmd{
		out: out,
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&cmd.opts.family, "f", "", "Specify family (inet, inet6, mpls, link)")
	fs.StringVar(&cmd.opts.family, "family", "", "Specify family (inet, inet6, mpls, link)")
	fs.BoolVar(&cmd.opts.resolve, "r", false, "Use system resolver to display DNS names")
	fs.BoolVar(&cmd.opts.resolve, "resolve", false, "Use system resolver to display DNS names")
	fs.BoolVar(&cmd.opts.inet4, "4", false, "Set protocol family to inet")
	fs.BoolVar(&cmd.opts.inet6, "6", false, "Set protocol family to inet6")
	fs.BoolVar(&cmd.opts.bridge, "B", false, "Set protocol family to bridge")
	fs.BoolVar(&cmd.opts.mpls, "M", false, "Set protocol family to mpls")
	fs.BoolVar(&cmd.opts.link, "0", false, "Set protocol family to link")
	fs.BoolVar(&cmd.opts.details, "d", false, "Display details")
	fs.BoolVar(&cmd.opts.details, "details", false, "Display details")
	fs.BoolVar(&cmd.opts.stats, "s", false, "Display statistics")
	fs.BoolVar(&cmd.opts.stats, "statistics", false, "Display statistics")
	fs.IntVar(&cmd.opts.loops, "l", 0, "Set maximum number of attempts to flush all addresses")
	fs.IntVar(&cmd.opts.loops, "loops", 1, "Set maximum number of attempts to flush all addresses")
	fs.BoolVar(&cmd.opts.humanReadable, "h", false, "Display timings and sizes in human readable format")
	fs.BoolVar(&cmd.opts.humanReadable, "humanreadable", false, "Display timings and sizes in human-readable format")
	fs.BoolVar(&cmd.opts.iec, "iec", false, "Use 1024-based block sizes for human-readable sizes")
	fs.BoolVar(&cmd.opts.brief, "br", false, "Brief output")
	fs.BoolVar(&cmd.opts.brief, "brief", false, "Brief output")
	fs.BoolVar(&cmd.opts.json, "j", false, "Output in JSON format")
	fs.BoolVar(&cmd.opts.json, "json", false, "Output in JSON format")
	fs.BoolVar(&cmd.opts.prettify, "p", false, "Make JSON output pretty")
	fs.BoolVar(&cmd.opts.prettify, "pretty", false, "Make JSON output pretty")
	fs.StringVar(&cmd.opts.color, "c", "", "Use color output")
	fs.StringVar(&cmd.opts.color, "color", "", "Use color output")
	fs.StringVar(&cmd.opts.rcvBuf, "rc", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.StringVar(&cmd.opts.rcvBuf, "rcvbuf", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.BoolVar(&cmd.opts.timeStamp, "t", false, "Display time stamps")
	fs.BoolVar(&cmd.opts.timeStamp, "timestamp", false, "Display time stamps")
	fs.BoolVar(&cmd.opts.timeStampShort, "ts", false, "Display short time stamps")
	fs.BoolVar(&cmd.opts.timeStampShort, "tshort", false, "Display short time stamps")
	fs.BoolVar(&cmd.opts.all, "a", false, "Display all information")
	fs.BoolVar(&cmd.opts.all, "all", false, "Display all information")
	fs.BoolVar(&cmd.opts.numeric, "N", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.BoolVar(&cmd.opts.numeric, "numeric", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.StringVar(&cmd.opts.batch, "b", "", "Read commands from a file")
	fs.StringVar(&cmd.opts.batch, "batch", "", "Read commands from a file")
	fs.BoolVar(&cmd.opts.force, "force", false, "Don't terminate ip on errors in batch mode.  If there were any errors during execution of the commands, the application return code will be non zero.")
	fs.BoolVar(&cmd.opts.oneline, "o", false, "Output each record on a single line")
	fs.BoolVar(&cmd.opts.oneline, "oneline", false, "Output each record on a single line")
	fs.StringVar(&cmd.opts.netns, "n", "", "Switch to network namespace")
	fs.StringVar(&cmd.opts.netns, "netns", "", "Switch to network namespace")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", ipHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	cmd.args = fs.Args()

	cmd.family = netlink.FAMILY_ALL

	var families []int

	if cmd.opts.inet4 {
		cmd.family = netlink.FAMILY_V4
		families = append(families, netlink.FAMILY_V4)
	}
	if cmd.opts.inet6 {
		cmd.family = netlink.FAMILY_V6
		families = append(families, netlink.FAMILY_V6)
	}

	if cmd.opts.mpls {
		return nil, fmt.Errorf("protocol family MPLS is not yet supported")
	}

	if cmd.opts.bridge {
		return nil, fmt.Errorf("protocol family bridge is not yet supported")
	}

	if cmd.opts.link {
		cmd.family = netlink.FAMILY_ALL
		families = append(families, netlink.FAMILY_ALL)
	}

	if cmd.opts.resolve {
		return nil, fmt.Errorf("resolving DNS names is unsupported")
	}

	if cmd.opts.color != "" {
		return nil, fmt.Errorf("color output is unsupported")
	}

	if cmd.opts.oneline {
		return nil, fmt.Errorf("outputting each record on a single line is unsupported")
	}

	var (
		err    error
		handle *netlink.Handle
	)

	if cmd.opts.netns != "" {
		nsHandle, err := netns.GetFromName(cmd.opts.netns)
		if err != nil {
			return nil, fmt.Errorf("failed to find network namespace %q: %v", cmd.opts.netns, err)
		}
		defer nsHandle.Close()

		handle, err = netlink.NewHandleAt(nsHandle, families...)
		if err != nil {
			return nil, fmt.Errorf("failed to create netlink handle in network namespace %q: %v", cmd.opts.netns, err)
		}
	} else {
		handle, err = netlink.NewHandle(families...)
		if err != nil {
			return nil, fmt.Errorf("failed to create netlink handle: %v", err)
		}
	}

	if cmd.opts.rcvBuf != "" {
		bufSize, err := strconv.Atoi(cmd.opts.rcvBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rcvbuf flag: %v", err)
		}

		handle.SetSocketReceiveBufferSize(bufSize, true)
	}

	cmd.handle = handle

	return &cmd, nil
}

type cmd struct {
	// Output writer
	out io.Writer
	// Netlink handle for all netlink ops
	handle *netlink.Handle
	// Cursor is our next token pointer
	cursor int
	// Options
	opts flags
	// Arguments
	args []string
	// Expected values for current token placement
	expectedValues []string
	// Selected protocol family
	family int
}

func (cmd cmd) run() error {
	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", cmd.args, cmd.cursor, cmd.expectedValues)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", cmd.args, cmd.cursor, cmd.expectedValues)
			}
			log.Fatalf("ip: %v", err)
		default:
			log.Fatalf("ip: unexpected panic value: %T(%v)", err, err)
		}

		return
	}()

	if cmd.opts.batch != "" {
		return cmd.batchCmds()
	}

	return cmd.runSubCommand()
}

func (cmd cmd) batchCmds() error {
	file, err := os.Open(cmd.opts.batch)
	if err != nil {
		log.Fatalf("Failed to open batch file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd.args = strings.Fields(line)

		if len(cmd.args) == 0 { // Skip empty lines
			continue
		}

		err := cmd.runSubCommand()
		if err != nil {
			if cmd.opts.force {
				log.Printf("Error (force mode on, continuing): Failed to run command '%s': %v", line, err)
			} else {
				return fmt.Errorf("failed to run command '%s': %v", line, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading batch file: %v", err)
	}

	return nil
}

func (cmd cmd) runSubCommand() error {
	cmd.cursor = -1

	switch c := cmd.findPrefix("address", "route", "link", "monitor", "neigh", "tunnel", "tuntap", "tap", "tcp_metrics", "tcpmetrics", "vrf", "xfrm", "help"); c {
	case "address":
		return cmd.address()
	case "link":
		return cmd.link()
	case "route":
		return cmd.route()
	case "neigh":
		return cmd.neigh()
	case "monitor":
		return cmd.monitor()
	case "tunnel":
		return cmd.tunnel()
	case "tuntap", "tap":
		return cmd.tuntap()
	case "tcpmetrics", "tcp_metrics":
		return cmd.tcpMetrics()
	case "vrf":
		return cmd.vrf()
	case "xfrm":
		return cmd.xfrm()
	case "help":
		fmt.Fprint(cmd.out, ipHelp)

		return nil
	default:
		return cmd.usage()
	}
}

func main() {
	cmd, err := parseFlags(os.Args, os.Stdout)
	if err != nil {
		log.Fatalf("ip: %v", err)
	}

	err = cmd.run()
	if err != nil {
		log.Fatalf("ip: %v", err)
	}
}
