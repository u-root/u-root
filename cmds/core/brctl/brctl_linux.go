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
brctl addbr <name>   creates a new instance of the ethernet bridge
brctl delbr <name>   deletes the instance <name> of the ethernet bridge
brctl show           show current instance(s) of the ethernet bridge

PORTS:
brctl addif <brname> <ifname>              will make the interface <ifname> a port of the bridge <brname>
brctl delif <brname> <ifname>              will detach the interface <ifname> from the bridge <brname>
brctl show <brname>                        will show some information on the bridge and its attached ports
brctl hairpin <bridge> <port> {on | off}   enable/disable hairpin mode on the port <port> of the bridge <bridge>

AGEING:
brctl showmacs <brname>           shows a list of learned MAC addresses for this bridge
brctl setageing <brname> <time>   sets the ethernet (MAC) address ageing to <time> seconds

SPANNING TREE PROTOCOL (IEEE 802.1d):
brctl stp <bridge> {on | off | yes | no}       controls the bridge participation in the spanning tree protocol
brctl showstp <bridge>                         shows stp related information of <bridge>
brctl setbridgeprio <bridge> <PRIORITY>        sets the bridge's priority to <priority>
brctl setfd <bridge> <time>                    sets the bridge's 'bridge forward delay' to <time> seconds
brctl sethello <bridge> <time>                 sets the bridge's 'bridge hello time' to <time> seconds
brctl setmaxage <bridge> <time>                sets the bridge's 'maximum message age' to <time> seconds
brctl setpathcost <bridge> <port> <COST>       sets the port cost of the port <port> to <cost>
brctl setportprio <bridge> <port> <PRIORITY>   sets the port's priority to <priority>

COST is a dimensionless metric (from 1 to 65535, default is 100).
PRIORITY is a number from 0 (min) to 63 (max), default is 32, and has no dimension.
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
		err = brctl.ShowMACs(args[0], out)

	case "setageing":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetAgeingTime(args[0], args[1])

	case "stp":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetSTP(args[0], args[1])

	case "showstp":
		if len(args) != 1 {
			return errFewArgs
		}
		err = brctl.ShowStp(out, args[0])

	case "setbridgeprio":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetBridgePrio(args[0], args[1])

	case "setfd":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetForwardDelay(args[0], args[1])

	case "sethello":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetHello(args[0], args[1])

	case "setmaxage":
		if len(args) != 2 {
			return errFewArgs
		}
		err = brctl.SetMaxAge(args[0], args[1])

	case "setpathcost":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.SetPathCost(args[0], args[1], args[2])

	case "setportprio":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.SetPortPrio(args[0], args[1], args[2])

	case "hairpin":
		if len(args) != 3 {
			return errFewArgs
		}
		err = brctl.Hairpin(args[0], args[1], args[2])

	case "help", "-h", "--help":
		fmt.Fprintf(out, "%s\n", usage)
		return nil

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
		log.Fatal(err)
	}
}
