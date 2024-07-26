// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/traceroute"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type Args struct {
	Host string
}

func run(args []string) error {
	flags := &traceroute.Flags{}
	trargs := &Args{}

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	// Short form flags - must be provided with a single dash (-)
	// If not provided with cmdline, traceroute resolves host to IPv4 and IPv6 addresses, but
	// ALWAYS uses IPv4 in this case.
	f.BoolVar(&flags.AF4, "4", false, "Explicitly  force  IPv4 tracerouting.")
	f.BoolVar(&flags.AF6, "6", false, "Explicitly  force  IPv6 tracerouting.")
	f.UintVar(&flags.DestPortSeq, "p", 0, "Destination port")
	f.StringVar(&flags.Module, "m", "udp4", "udp, tcp, icmp")

	// Long form flags - must be provided with two dashes (--)
	f.UintVar(&flags.DestPortSeq, "port", 0, "Destination port")
	f.StringVar(&flags.Module, "module", "", "udp, tcp, icmp")
	f.BoolVar(&flags.ICMP, "icmp", false, "Use ICMP method. Same as -m icmp")
	f.BoolVar(&flags.TCP, "tcp", false, "Use TCP method. Same as -m tcp")
	f.BoolVar(&flags.UDP, "udp", true, "Use UDP method. Same as -m udp")

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	leftoverArgs := f.Args()
	if len(leftoverArgs) < 1 {
		// Error, print help and exit
		f.Usage()
		return nil
	}

	trargs.Host = leftoverArgs[0]

	flags.Host = trargs.Host

	// Evaluate AF and Module
	af := "4"
	if !flags.AF4 && flags.AF6 {
		af = "6"
	}

	if (flags.TCP || flags.ICMP || flags.UDP) && flags.Module == "" {
		if flags.TCP {
			flags.Module = "tcp"
		} else if flags.ICMP {
			flags.Module = "icmp"
		} else if flags.UDP {
			flags.Module = "udp"
		}
	}

	modaf := strings.ToLower(fmt.Sprintf("%s%s", flags.Module, af))
	// Pass execution to pkg/traceroute.
	// Setup can be quite complexe with such amount of flags
	// and the different modules.
	return traceroute.RunTraceroute(modaf, flags)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
