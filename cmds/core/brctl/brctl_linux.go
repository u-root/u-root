// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/brctl"
)

// cli
const usage = `
Usage: brctl [commands]
commands:

INSTANCES:
brctl addbr <name>		creates a new instance of the ethernet bridge
brctl delbr <name>		deletes the instance <name> of the ethernet bridge
brctl show			show current instance(s) of the ethernet bridge

PORTS:
brctl addif <brname> <ifname>	will make the interface <ifname> a port of the bridge <brname>
brctl delif <brname> <ifname>	will detach the interface <ifname> from the bridge <brname>
brctl show <brname>		will show some information on the bridge and its attached ports

AGEING:
brctl showmacs <brname> 		shows a list of learned MAC addresses for this bridge
brctl setageingtime <brname> <time>	sets the ethernet (MAC) address ageing time, in seconds [OPT]

SPANNING TREE PROTOCOL (IEEE 802.1d):
brctl stp <bridge> <state>			controls this bridge instance's participation in the spanning tree protocol.
brctl showstp <bridge>				shows stp related information of <bridge>.
brctl setbridgeprio <bridge> <priority>		sets the bridge's priority to <priority>
brctl setfd <bridge> <time>			sets the bridge's 'bridge forward delay' to <time> seconds
brctl sethello <bridge> <time>			sets the bridge's 'bridge hello time' to <time> seconds
brctl setmaxage <bridge> <time>			sets the bridge's 'maximum message age' to <time> seconds.
brctl setpathcost <bridge> <port> <cost>	sets the port cost of the port <port> to <cost>. This is a dimensionless metric
brctl setportprio <bridge> <port> <priority>	sets the port <port>'s priority to <priority>
brctl hairpin <bridge> <port> <state>		enable/disable hairpin mode on the port <port> of the bridge <bridge>
`

var (
	errFewArgs    = errors.New("too few args")
	errInvalidCmd = errors.New("unknown command")
)

func run(out io.Writer, argv []string) error {
	var err error
	command := argv[0]
	args := argv[1:]

	switch command {
	case "addbr":
		if len(args) != 1 {
			return errFewArgs
		}
		err = brctl.Addbr(args[0])

	case "delbr":
		if len(args) != 1 {
			return errFewArgs
		}
		err = brctl.Delbr(args[0])

	case "addif":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Addif(args[0], args[1])

	case "delif":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Delif(args[0], args[1])

	case "show":
		err = brctl.Show(out, args...)

	case "showmacs":
		if len(args) != 1 {
			return errFewArgs
		}
		err = brctl.Showmacs(args[0], out)

	case "setageingtime":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Setageingtime(args[0], args[1])

	case "stp":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Stp(args[0], args[1])

	case "showstp":
		if len(args) != 1 {
			return errFewArgs
		}
		err = brctl.ShowStp(out, args[0])

	case "setbridgeprio":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Setbridgeprio(args[0], args[1])

	case "setfd":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Setfd(args[0], args[1])

	case "sethello":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Sethello(args[0], args[1])

	case "setmaxage":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.Setmaxage(args[0], args[1])

	case "setpathcost":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.Setpathcost(args[0], args[1], args[2])

	case "setportprio":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.Setportprio(args[0], args[1], args[2])

	case "hairpin":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.Hairpin(args[0], args[1], args[2])

	default:
		return fmt.Errorf("%w: %s", errInvalidCmd, command)
	}

	return err
}

func main() {
	argv := os.Args

	if len(argv) < 2 {
		log.Fatal(usage)
	}

	if err := run(os.Stdout, argv[1:]); err != nil {
		log.Fatalf("brctl: %v", err)
	}
}
