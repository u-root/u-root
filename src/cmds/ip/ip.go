// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"netlink"
)

// you will notice that I suck at parsers. That said, here is the method to my madness.
// The language of ip is not super consistent and has lots of convenience shortcuts.
// The BNF it shows you doesn't show them.
// The inputs is just the set of args.
// It's very short.
// Each token is just a string and we need not produce terminals with them -- they can
// just be the terminals and we can switch on them.
// The cursor is always our current token pointer. We do a dumb recursive descent parser
// and accumulate information into a global set of variables. At any point we can see into the
// whole set of args and see where we are. We can indicate at each point what we're expecting so
// that in usage() or recover() we can tell the user exactly what we wanted, unlike IP,
// which just barfs a whole (incorrect) BNF at you when you do anything wrong.
// To handle errors in too few arguments, we just do a recover block. That lets us blindly
// reference the arg[] array without having to check the length everywhere.

// Note the plethora of globals. The reason is simple: we parse one command, do it, and quit.
// It doesn't make sense to write this otherwise.
var (
	// Cursor is out next token pointer.
	// The language of this command doesn't require much more.
	cursor    int
	arg       []string
	whatIWant = "addr|route|link"
	l         = log.New(os.Stdout, "ip: ", 0)
)

// the pattern:
// at each level parse off arg[0]. If it matches, continue. If it does not, all error with how far you got, what arg you saw,
// and why it did not work out.

func usage() {
	log.Fatalf("This was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], whatIWant)
}

// in the ip command, turns out 'dev' is a noise word.
// The BNF is not right there either.
// Always make it optional.
func dev() *net.Interface {
	cursor++
	whatIWant = "dev|device name"
	if arg[cursor] == "dev" {
		cursor++
	}
	whatIWant = "device name"
	iface, err := net.InterfaceByName(arg[cursor])
	if err != nil {
		usage()
	}
	return iface
}

func showips() {
	ifaces, err := net.Interfaces()
	if err != nil {
		l.Fatalf("Can't enumerate interfaces? %v", err)
	}
	for _, v := range ifaces {
		addrs, err := v.Addrs()
		if err != nil {
			l.Printf("Can't enumerate addresses")
		}
		l.Printf("%v: %v", v, addrs)
	}
}

func addrip() {
	var err error
	var addr net.IP
	var network *net.IPNet
	if len(arg) == 1 {
		showips()
		return
	}
	cursor++
	whatIWant = "add|del"
	cmd := arg[cursor]

	switch cmd {
	case "add", "del":
		cursor++
		whatIWant = "CIDR format address"
		addr, network, err = net.ParseCIDR(arg[cursor])
		if err != nil {
			usage()
		}
	default:
		usage()
	}
	iface := dev()
	switch cmd {
	case "add":
		if err := netlink.NetworkLinkAddIp(iface, addr, network); err != nil {
			l.Fatalf("Adding %v to %v failed: %v", arg[1], arg[2], err)
		}
	case "del":
		if err := netlink.NetworkLinkDelIp(iface, addr, network); err != nil {
			l.Fatalf("Deleting %v from %v failed: %v", arg[1], arg[2], err)
		}
	default:
		l.Fatalf("devip: arg[0] changed: can't happen")
	}
	return

}

func linkshow() {
	cursor++
	whatIWant = "<nothing>|<device name>"
	if len(arg[cursor:]) == 0 {
		showips()
	}
}

func linkset() {
	iface := dev()
	cursor++
	whatIWant = "up|down"
	switch arg[cursor] {
	case "up":
		if err := netlink.NetworkLinkUp(iface); err != nil {
			l.Fatalf("%v can't make it up: %v", dev, err)
		}
	case "down":
		if err := netlink.NetworkLinkDown(iface); err != nil {
			l.Fatalf("%v can't make it down: %v", dev, err)
		}
	default:
		usage()
	}
}

func link() {
	cursor++
	whatIWant = "show|set"
	cmd := arg[cursor]

	switch cmd {
	case "show":
		linkshow()
	case "set":
		linkset()
	default:
		usage()
	}
	return
}

func routeshow() {
	if b, err := ioutil.ReadFile("/proc/net/route"); err == nil {
		l.Printf("%s", string(b))
	} else {
		l.Fatalf("Route show failed: %v", err)
	}
}
func nodespec() string {
	cursor++
	whatIWant = "default|CIDR"
	return arg[cursor]
}

func nexthop() (string, string) {
	cursor++
	whatIWant = "via"
	if arg[cursor] != "via" {
		usage()
	}
	nh := arg[cursor]
	cursor++
	whatIWant = "Gateway CIDR"
	return nh, arg[cursor]
}

func routeadddefault() {
	nh, nhval := nexthop()
	// TODO: NHFLAGS.
	d := dev()
	switch nh {
	case "via":
		l.Printf("Add default route %v via %v", nhval, d)
		netlink.AddDefaultGw(nhval, d.Name)
	default:
		usage()
	}
}

func routeadd() {
	ns := nodespec()
	switch ns {
	case "default":
		routeadddefault()
	default:
		usage()
	}
}

func route() {
	cursor++
	if len(arg[cursor:]) == 0 {
		routeshow()
		return
	}

	whatIWant = "show|add"
	switch arg[cursor] {
	case "show":
		routeshow()
	case "add":
		routeadd()
	default:
		usage()
	}

}
func main() {
	flag.Parse()
	arg = flag.Args()
	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				l.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				l.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			} else {
				l.Fatalf("FUCK: %v", err.Error())
			}
			l.Fatalf("Bummer: %v", err)
		default:
			l.Fatalf("unexpected panic value: %T(%v)", err, err)
		}
	}()
	// The ip command doesn't actually follow the BNF it prints on error.
	// There are lots of handy shortcuts that people will expect.
	switch arg[cursor] {
	case "addr":
		addrip()
	case "link":
		link()
	case "route":
		route()
	default:
		usage()
	}
}
