// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/traceroute"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

func run(out io.Writer, args []string) error {
	var err error
	flags := &traceroute.Flags{}
	trargs := &traceroute.Args{}

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	// Short form flags - must be provided with a single dash (-)
	// If not provided with cmdline, traceroute resolves host to IPv4 and IPv6 addresses, but
	// ALWAYS uses IPv4 in this case.
	f.BoolVar(&flags.AF4, "4", false, "Explicitly  force  IPv4 tracerouting.")
	f.BoolVar(&flags.AF6, "6", false, "Explicitly  force  IPv6 tracerouting.")
	f.UintVar(&flags.DestPortSeq, "p", 0, "Destination port")
	//f.StringVar(&flags.Module, "m", "default", traceroute.SupportedModules())

	// Long form flags - must be provided with two dashes (--)
	f.UintVar(&flags.DestPortSeq, "port", 0, "Destination port")
	//f.StringVar(&flags.Module, "module", "default", traceroute.SupportedModules())

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	leftoverArgs := f.Args()
	if len(leftoverArgs) < 1 {
		// Error, print help and exit
		f.Usage()
		return nil
	}

	trargs.Host = leftoverArgs[0]

	if len(leftoverArgs) > 1 {
		trargs.PacketLen, err = strconv.Atoi(leftoverArgs[1])
		if err != nil {
			return err
		}
	}

	// Pass execution to pkg/traceroute.
	// Setup can be quite complexe with such amount of flags
	// and the different modules.
	return traceroute.RunTraceroute(trargs.Host, false)
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
