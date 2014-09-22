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
)

var (
	// Cursor is out next token pointer.
	// The language of this command doesn't require much more.
	cursor    int
	arg       []string
	whatIWant = "addr|route"
	l         = log.New(os.Stdout, "ip: ", 0)
)

// the pattern:
// at each level parse off arg[0]. If it matches, continue. If it does not, all error with how far you got, what arg you saw,
// and why it did not work out.

func usage() {
	log.Fatalf("This was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], whatIWant)
}

func adddelip(op, ip, dev string) error {
	addr, network, err := net.ParseCIDR(ip)
	if err != nil {
		l.Fatalf("%v is not in CIDR format: %v", ip, err)
	}
	iface, err := net.InterfaceByName(dev)
	if err != nil {
		l.Fatalf("%v not found: %v", dev, err)
		return err
	}

	switch op {
	case "add":
		if err := NetworkLinkAddIp(iface, addr, network); err != nil {
			l.Fatalf("Adding %v to %v failed: %v", ip, dev, err)
		}
	default:
		l.Fatalf("%v is not supported yet", op)
	}
	return nil

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
	whatIWant = "add"
	cmd := arg[cursor]

	switch cmd {
	case "add":
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
		if err := NetworkLinkAddIp(iface, addr, network); err != nil {
			l.Fatalf("Adding %v to %v failed: %v", arg[1], arg[2], err)
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
	whatIWant="up|down"
	switch arg[cursor] {
		case "up":
		if err := NetworkLinkUp(iface); err != nil {
			l.Fatalf("%v can't make it up: %v", dev, err)
		}
		case "down":
		if err := NetworkLinkDown(iface); err != nil {
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
	switch {
	case arg[cursor] == "addr":
		addrip()

	case arg[cursor] == "link":
		link()
	case len(arg) == 1 && arg[0] == "route":
		if b, err := ioutil.ReadFile("/proc/net/route"); err == nil {
			l.Printf("%s", string(b))
		} else {
			l.Fatalf("Route failed: %v", err)
		}
	// oh, barf.
	case len(arg) == 5 && arg[0] == "link" && arg[1] == "set" && arg[2] == "dev" && arg[4] == "up":
		dev := arg[3]
		iface, err := net.InterfaceByName(dev)
		if err != nil {
			l.Fatalf("%v not found", dev)
		}
		if err = NetworkLinkUp(iface); err != nil {
			l.Fatalf("%v can't make it up: %v", dev, err)
		}
	case len(arg) == 7 && arg[0] == "route" && arg[1] == "add" && arg[2] == "default" && arg[3] == "via" && arg[5] == "dev":
		AddDefaultGw(arg[4], arg[6])
	default:
		l.Fatalf("We don't do this: %v; try addr or link or route", arg)
	}
}
