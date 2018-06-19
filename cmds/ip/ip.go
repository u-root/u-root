// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	l "log"
	"math"
	"net"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
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
	whatIWant []string
	log       = l.New(os.Stdout, "ip: ", 0)

	addrScopes = map[netlink.Scope]string{
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
	return errors.New(fmt.Sprintf("This was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], whatIWant))
}

func one(cmd string, cmds []string) string {
	var x, n int
	for i, v := range cmds {
		if strings.HasPrefix(v, cmd) {
			n++
			x = i
		}
	}
	if n == 1 {
		return cmds[x]
	}
	return ""
}

// in the ip command, turns out 'dev' is a noise word.
// The BNF is not right there either.
// Always make it optional.
func dev() (netlink.Link, error) {
	cursor++
	whatIWant = []string{"dev", "device name"}
	if arg[cursor] == "dev" {
		cursor++
	}
	whatIWant = []string{"device name"}
	iface, err := netlink.LinkByName(arg[cursor])
	if err != nil {
		//usage()
		return nil, err
	}
	return iface, nil
}

func showLinks(w io.Writer, withAddresses bool) error {
	ifaces, err := netlink.LinkList()
	if err != nil {
		return errors.New(fmt.Sprintf("Can't enumerate interfaces? %v", err))
	}

	for _, v := range ifaces {
		l := v.Attrs()

		fmt.Fprintf(w, "%d: %s: <%s> mtu %d state %s\n", l.Index, l.Name,
			strings.Replace(strings.ToUpper(fmt.Sprintf("%s", l.Flags)), "|", ",", -1),
			l.MTU, strings.ToUpper(l.OperState.String()))

		fmt.Fprintf(w, "    link/%s %s\n", l.EncapType, l.HardwareAddr)

		if withAddresses {
			err := showLinkAddresses(w, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func showLinkAddresses(w io.Writer, link netlink.Link) error {
	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		log.Printf("Can't enumerate addresses")
	}

	for _, addr := range addrs {

		var inet string
		switch len(addr.IPNet.IP) {
		case 4:
			inet = "inet"
		case 16:
			inet = "inet6"
		default:
			return errors.New("Can't figure out IP protocol version")
		}

		fmt.Fprintf(w, "    %s %s", inet, addr.Peer)
		if addr.Broadcast != nil {
			fmt.Fprintf(w, " brd %s", addr.Broadcast)
		}
		fmt.Fprintf(w, " scope %s %s\n", addrScopes[netlink.Scope(addr.Scope)], addr.Label)

		var validLft, preferredLft string
		// TODO: fix vishnavanda/netlink. *Lft should be uint32, not int.
		if uint32(addr.PreferedLft) == math.MaxUint32 {
			preferredLft = "forever"
		} else {
			preferredLft = fmt.Sprintf("%dsec", addr.PreferedLft)
		}
		if uint32(addr.ValidLft) == math.MaxUint32 {
			validLft = "forever"
		} else {
			validLft = fmt.Sprintf("%dsec", addr.ValidLft)
		}
		fmt.Fprintf(w, "       valid_lft %s preferred_lft %s\n", validLft, preferredLft)
	}
	return nil
}

func addrip() error {
	var err error
	var addr *netlink.Addr
	if len(arg) == 1 {
		err := showLinks(os.Stdout, true)
		if err != nil {
			return err
		}
		return nil
	}
	cursor++
	whatIWant = []string{"add", "del"}
	cmd := arg[cursor]

	c := one(cmd, whatIWant)
	switch c {
	case "add", "del":
		cursor++
		whatIWant = []string{"CIDR format address"}
		addr, err = netlink.ParseAddr(arg[cursor])
		if err != nil {
			return usage()
		}
	default:
		return usage()
	}
	iface, err := dev()
	if err != nil {
		return nil
	}
	switch c {
	case "add":
		if err := netlink.AddrAdd(iface, addr); err != nil {
			return errors.New(fmt.Sprintf("Adding %v to %v failed: %v", arg[1], arg[2], err))
		}
	case "del":
		if err := netlink.AddrDel(iface, addr); err != nil {
			return errors.New(fmt.Sprintf("Deleting %v from %v failed: %v", arg[1], arg[2], err))
		}
	default:
		return errors.New("devip: arg[0] changed: can't happen")
	}
	return nil
}

func linkshow() error {
	cursor++
	whatIWant = []string{"<nothing>", "<device name>"}
	if len(arg[cursor:]) == 0 {
		err := showLinks(os.Stdout, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func setHardwareAddress(iface netlink.Link) error {
	cursor++
	hwAddr, err := net.ParseMAC(arg[cursor])
	if err != nil {
		return errors.New(fmt.Sprintf("%v cant parse mac addr %v: %v", iface, hwAddr, err))
	}
	err = netlink.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return errors.New(fmt.Sprintf("%v cant set mac addr %v: %v", iface, hwAddr, err))
	}
	return nil
}

func linkset() error {
	iface, err := dev()
	if err != nil {
		return err
	}
	cursor++
	whatIWant = []string{"address", "up", "down"}
	switch one(arg[cursor], whatIWant) {
	case "address":
		if err := setHardwareAddress(iface); err != nil {
			return nil
		}
	case "up":
		if err := netlink.LinkSetUp(iface); err != nil {
			return errors.New(fmt.Sprintf("%v can't make it up: %v", iface, err))
		}
	case "down":
		if err := netlink.LinkSetDown(iface); err != nil {
			return errors.New(fmt.Sprintf("%v can't make it down: %v", iface, err))
		}
	default:
		return usage()
	}
	return nil
}

func link() error {
	if len(arg) == 1 {
		err := linkshow()
		if err != nil {
			return err
		}
		return nil
	}

	cursor++
	whatIWant = []string{"show", "set"}
	cmd := arg[cursor]

	switch one(cmd, whatIWant) {
	case "show":
		err := linkshow()
		if err != nil {
			return err
		}
	case "set":
		err := linkset()
		if err != nil {
			return err
		}
	default:
		return usage()
	}
	return nil
}

func routeshow() error {
	if b, err := ioutil.ReadFile("/proc/net/route"); err == nil {
		log.Printf("%s", string(b))
	} else {
		return errors.New(fmt.Sprintf("Route show failed: %v", err))
	}
	return nil
}

func nodespec() string {
	cursor++
	whatIWant = []string{"default", "CIDR"}
	return arg[cursor]
}

func nexthop() (string, *netlink.Addr, error) {
	cursor++
	whatIWant = []string{"via"}
	if arg[cursor] != "via" {
		return "", nil, usage()
	}
	nh := arg[cursor]
	cursor++
	whatIWant = []string{"Gateway CIDR"}
	addr, err := netlink.ParseAddr(arg[cursor])
	if err != nil {
		return "", nil, errors.New(fmt.Sprintf("Gateway CIDR: %v", err))
	}
	return nh, addr, nil
}

func routeadddefault() error {
	nh, nhval, err := nexthop()
	if err != nil {
		return nil
	}
	// TODO: NHFLAGS.
	l, err := dev()
	if err != nil {
		return nil
	}
	switch nh {
	case "via":
		log.Printf("Add default route %v via %v", nhval, l)
		r := &netlink.Route{LinkIndex: l.Attrs().Index, Gw: nhval.IPNet.IP}
		if err := netlink.RouteAdd(r); err != nil {
			return errors.New(fmt.Sprintf("Add default route: %v", err))
		}

	default:
		return usage()
	}
	return nil
}

func routeadd() error {
	ns := nodespec()
	switch ns {
	case "default":
		err := routeadddefault()
		if err != nil {
			return err
		}
	default:
		return usage()
	}
	return nil
}

func route() error {
	cursor++
	if len(arg[cursor:]) == 0 {
		err := routeshow()
		if err != nil {
			return err
		}
		return nil
	}

	whatIWant = []string{"show", "add"}
	switch one(arg[cursor], whatIWant) {
	case "show":
		err := routeshow()
		if err != nil {
			return err
		}
	case "add":
		err := routeadd()
		if err != nil {
			return err
		}
	default:
		return usage()
	}
	return nil
}

func main() {
	// When this is embedded in busybox we need to reinit some things.
	whatIWant = []string{"addr", "route", "link"}
	cursor = 0
	flag.Parse()
	arg = flag.Args()

	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				log.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				log.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			}
			log.Fatalf("Bummer: %v", err)
		default:
			log.Fatalf("unexpected panic value: %T(%v)", err, err)
		}
	}()

	// The ip command doesn't actually follow the BNF it prints on error.
	// There are lots of handy shortcuts that people will expect.
	switch one(arg[cursor], whatIWant) {
	case "addr":
		err := addrip()
		if err != nil {
			log.Fatalf("%v", err)
		}
	case "link":
		err := link()
		if err != nil {
			log.Fatalf("%v", err)
		}
	case "route":
		err := route()
		if err != nil {
			log.Fatalf("%v", err)
		}
	default:
		log.Fatalf("%v", usage())
	}
}
