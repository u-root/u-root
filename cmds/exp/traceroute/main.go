// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/traceroute"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errFlags = errors.New("invalid flag/argument usage")

func parseFlags(args []string) (*traceroute.Flags, error) {
	flags := &traceroute.Flags{}
	trargs := &traceroute.Args{}

	var af4, af6 bool

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	// Short form flags - must be provided with a single dash (-)
	// If not provided with cmdline, traceroute resolves host to IPv4 and IPv6 addresses, but
	// ALWAYS uses IPv4 in this case.
	f.BoolVar(&af4, "4", false, "Explicitly force IPv4 tracerouting.")
	f.BoolVar(&af6, "6", false, "Explicitly force IPv6 tracerouting.")
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

	if len(leftoverArgs) != 1 {
		// Error, print help and exit
		f.Usage()
		return nil, errFlags
	}

	trargs.Host = leftoverArgs[0]

	flags.Host = trargs.Host

	// Evaluate AF and Module
	af := "4"
	if !af4 && af6 {
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

	flags.Proto = strings.ToLower(fmt.Sprintf("%s%s", flags.Module, af))

	return flags, nil
}

func run(args []string) error {
	flags, err := parseFlags(args)
	if err != nil {
		return err
	}
	// Pass execution to pkg/traceroute.
	// Setup can be quite complex with such amount of flags
	// and the different modules.
	return traceroute.RunTraceroute(flags)
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
