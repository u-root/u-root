// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

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
	"golang.org/x/sys/unix"
)

type flags struct {
	Family         string
	Inet4          bool
	Inet6          bool
	Bridge         bool
	MPLS           bool
	Link           bool
	Details        bool
	Stats          bool
	Loops          int
	HumanReadable  bool
	Iec            bool
	JSON           bool
	Prettify       bool
	Brief          bool
	Resolve        bool
	Color          string
	RcvBuf         string
	TimeStamp      bool
	TimeStampShort bool
	All            bool
	Numeric        bool
	Batch          string
	Force          bool
	Oneline        bool
	Netns          string
}

const ipHelp = `Usage: ip [ OPTIONS ] OBJECT { COMMAND | help }
where  OBJECT := { address |  help | link | monitor | neighbor | neighbour |
				   route | tap | tcpmetrics | tunnel | tuntap | vrf | xfrm }
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
//
// the pattern:
// at each level parse off arg[0]. If it matches, continue. If it does not, all error with how far you got, what arg you saw,
// and why it did not work out.

func (cmd *cmd) usage() error {
	return fmt.Errorf("this was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		cmd.Args[0:cmd.Cursor], cmd.Args[cmd.Cursor:], cmd.Args[cmd.Cursor], cmd.ExpectedValues)
}

func parseFlags(args []string, out io.Writer) (cmd, error) {
	cmd := cmd{
		Out: out,
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&cmd.Opts.Family, "f", "", "Specify family (inet, inet6, mpls, link)")
	fs.StringVar(&cmd.Opts.Family, "family", "", "Specify family (inet, inet6, mpls, link)")
	fs.BoolVar(&cmd.Opts.Resolve, "r", false, "Use system resolver to display DNS names")
	fs.BoolVar(&cmd.Opts.Resolve, "resolve", false, "Use system resolver to display DNS names")
	fs.BoolVar(&cmd.Opts.Inet4, "4", false, "Set protocol family to inet")
	fs.BoolVar(&cmd.Opts.Inet6, "6", false, "Set protocol family to inet6")
	fs.BoolVar(&cmd.Opts.Bridge, "B", false, "Set protocol family to bridge")
	fs.BoolVar(&cmd.Opts.MPLS, "M", false, "Set protocol family to mpls")
	fs.BoolVar(&cmd.Opts.Link, "0", false, "Set protocol family to link")
	fs.BoolVar(&cmd.Opts.Details, "d", false, "Display details")
	fs.BoolVar(&cmd.Opts.Details, "details", false, "Display details")
	fs.BoolVar(&cmd.Opts.Stats, "s", false, "Display statistics")
	fs.BoolVar(&cmd.Opts.Stats, "statistics", false, "Display statistics")
	fs.IntVar(&cmd.Opts.Loops, "l", 0, "Set maximum number of attempts to flush all addresses")
	fs.IntVar(&cmd.Opts.Loops, "loops", 1, "Set maximum number of attempts to flush all addresses")
	fs.BoolVar(&cmd.Opts.HumanReadable, "h", false, "Display timings and sizes in human readable format")
	fs.BoolVar(&cmd.Opts.HumanReadable, "humanreadable", false, "Display timings and sizes in human-readable format")
	fs.BoolVar(&cmd.Opts.Iec, "iec", false, "Use 1024-based block sizes for human-readable sizes")
	fs.BoolVar(&cmd.Opts.Brief, "br", false, "Brief output")
	fs.BoolVar(&cmd.Opts.Brief, "brief", false, "Brief output")
	fs.BoolVar(&cmd.Opts.JSON, "j", false, "Output in JSON format")
	fs.BoolVar(&cmd.Opts.JSON, "json", false, "Output in JSON format")
	fs.BoolVar(&cmd.Opts.Prettify, "p", false, "Make JSON output pretty")
	fs.BoolVar(&cmd.Opts.Prettify, "pretty", false, "Make JSON output pretty")
	fs.StringVar(&cmd.Opts.Color, "c", "", "Use color output")
	fs.StringVar(&cmd.Opts.Color, "color", "", "Use color output")
	fs.StringVar(&cmd.Opts.RcvBuf, "rc", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.StringVar(&cmd.Opts.RcvBuf, "rcvbuf", "", "Set the netlink socket receive buffer size, defaults to 1MB")
	fs.BoolVar(&cmd.Opts.TimeStamp, "t", false, "Display time stamps")
	fs.BoolVar(&cmd.Opts.TimeStamp, "timestamp", false, "Display time stamps")
	fs.BoolVar(&cmd.Opts.TimeStampShort, "ts", false, "Display short time stamps")
	fs.BoolVar(&cmd.Opts.TimeStampShort, "tshort", false, "Display short time stamps")
	fs.BoolVar(&cmd.Opts.All, "a", false, "Display all information")
	fs.BoolVar(&cmd.Opts.All, "all", false, "Display all information")
	fs.BoolVar(&cmd.Opts.Numeric, "N", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.BoolVar(&cmd.Opts.Numeric, "numeric", false, "Print the number of protocol, scope, dsfield, etc directly instead of converting it to human readable name.")
	fs.StringVar(&cmd.Opts.Batch, "b", "", "Read commands from a file")
	fs.StringVar(&cmd.Opts.Batch, "batch", "", "Read commands from a file")
	fs.BoolVar(&cmd.Opts.Force, "force", false, "Don't terminate ip on errors in batch mode.  If there were any errors during execution of the commands, the application return code will be non zero.")
	fs.BoolVar(&cmd.Opts.Oneline, "o", false, "Output each record on a single line")
	fs.BoolVar(&cmd.Opts.Oneline, "oneline", false, "Output each record on a single line")
	fs.StringVar(&cmd.Opts.Netns, "n", "", "Switch to network namespace")
	fs.StringVar(&cmd.Opts.Netns, "netns", "", "Switch to network namespace")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", ipHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))
	cmd.Args = fs.Args()

	cmd.Family = netlink.FAMILY_ALL

	if cmd.Opts.Inet4 {
		cmd.Family = netlink.FAMILY_V4
	}
	if cmd.Opts.Inet6 {
		cmd.Family = netlink.FAMILY_V6
	}

	if cmd.Opts.MPLS {
		return cmd, fmt.Errorf("protocol family MPLS is not yet supported")
	}

	if cmd.Opts.Bridge {
		return cmd, fmt.Errorf("protocol family bridge is not yet supported")
	}

	if cmd.Opts.Link {
		cmd.Family = netlink.FAMILY_ALL
	}

	if cmd.Opts.Family != "" {
		switch cmd.Opts.Family {
		case "inet":
			cmd.Family = netlink.FAMILY_V4
		case "inet6":
			cmd.Family = netlink.FAMILY_V6
		default:
			return cmd, fmt.Errorf("invalid family %q", cmd.Opts.Family)
		}
	}

	if cmd.Opts.Resolve {
		return cmd, fmt.Errorf("resolving DNS names is unsupported")
	}

	if cmd.Opts.Color != "" {
		return cmd, fmt.Errorf("color output is unsupported")
	}

	var (
		err    error
		handle *netlink.Handle
	)

	if cmd.Opts.Netns != "" {
		nsHandle, err := netns.GetFromName(cmd.Opts.Netns)
		if err != nil {
			return cmd, fmt.Errorf("failed to find network namespace %q: %w", cmd.Opts.Netns, err)
		}
		defer nsHandle.Close()

		handle, err = netlink.NewHandleAt(nsHandle, unix.NETLINK_ROUTE)
		if err != nil {
			return cmd, fmt.Errorf("failed to create netlink handle in network namespace %q: %w", cmd.Opts.Netns, err)
		}
	} else {
		handle, err = netlink.NewHandle(unix.NETLINK_ROUTE)
		if err != nil {
			return cmd, fmt.Errorf("failed to create netlink handle: %w", err)
		}
	}

	if cmd.Opts.RcvBuf != "" {
		bufSize, err := strconv.Atoi(cmd.Opts.RcvBuf)
		if err != nil {
			return cmd, fmt.Errorf("failed to parse rcvbuf flag: %w", err)
		}

		handle.SetSocketReceiveBufferSize(bufSize, true)
	}

	cmd.handle = handle

	return cmd, nil
}

type cmd struct {
	// Output writer
	Out io.Writer
	// Netlink handle for all netlink ops
	handle *netlink.Handle
	// Cursor is our next token pointer
	Cursor int
	// Options
	Opts flags
	// Arguments
	Args []string
	// Expected values for current token placement
	ExpectedValues []string
	// Selected protocol Family
	Family int
}

func (cmd *cmd) run() error {
	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", cmd.Args, cmd.Cursor, cmd.ExpectedValues)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				log.Fatalf("ip: args: %v, I got to arg %v, expected %v after that", cmd.Args, cmd.Cursor, cmd.ExpectedValues)
			}
			log.Fatalf("ip: %v", err)
		default:
			log.Fatalf("ip: unexpected panic value: %T(%v)", err, err)
		}
	}()

	if cmd.Opts.Batch != "" {
		return cmd.batchCmds()
	}

	return cmd.runSubCommand()
}

func (cmd *cmd) batchCmds() error {
	file, err := os.Open(cmd.Opts.Batch)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd.Args = strings.Fields(line)

		if len(cmd.Args) == 0 { // Skip empty lines
			continue
		}

		err := cmd.runSubCommand()
		if err != nil {
			if cmd.Opts.Force {
				log.Printf("Error (force mode on, continuing): Failed to run command '%s': %v", line, err)
			} else {
				return fmt.Errorf("failed to run command '%s': %w", line, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading batch file: %v", err)
	}

	return nil
}

func (cmd *cmd) runSubCommand() error {
	cmd.Cursor = -1

	if !cmd.tokenRemains() {
		fmt.Fprint(cmd.Out, ipHelp)
	}

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
		fmt.Fprint(cmd.Out, ipHelp)

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
