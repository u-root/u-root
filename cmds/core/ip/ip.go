// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"strconv"

	"encoding/hex"

	flag "github.com/spf13/pflag"

	"github.com/vishvananda/netlink"
)

var inet6 = flag.BoolP("6", "6", false, "use ipv6")

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
	cursor    int
	arg       []string
	whatIWant []string
	flog       = log.New(os.Stdout, "ip: ", 0)

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
	return fmt.Errorf("This was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], whatIWant)
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
// The BNF it shows is not right in that case.
// Always make 'dev' optional.
func dev() (netlink.Link, error) {
	cursor++
	whatIWant = []string{"dev", "device name"}
	if arg[cursor] == "dev" {
		cursor++
	}
	whatIWant = []string{"device name"}
	return netlink.LinkByName(arg[cursor])
}

func addrip() error {
	var err error
	var addr *netlink.Addr
	if len(arg) == 1 {
		return showLinks(os.Stdout, true)
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
			return err
		}
	default:
		return usage()
	}
	iface, err := dev()
	if err != nil {
		return err
	}
	switch c {
	case "add":
		if err := netlink.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("Adding %v to %v failed: %v", arg[1], arg[2], err)
		}
	case "del":
		if err := netlink.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("Deleting %v from %v failed: %v", arg[1], arg[2], err)
		}
	default:
		return fmt.Errorf("devip: arg[0] changed: can't happen")
	}
	return nil
}

func neigh() error {
	if len(arg) != 1 {
		return errors.New("neigh subcommands not supported yet")
	}
	return showNeighbours(os.Stdout, true)
}

func linkshow() error {
	cursor++
	whatIWant = []string{"<nothing>", "<device name>"}
	if len(arg[cursor:]) == 0 {
		return showLinks(os.Stdout, false)
	}
	return nil
}

func setHardwareAddress(iface netlink.Link) error {
	cursor++
	hwAddr, err := net.ParseMAC(arg[cursor])
	if err != nil {
		return fmt.Errorf("%v cant parse mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
	}
	err = netlink.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return fmt.Errorf("%v cant set mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
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
		return setHardwareAddress(iface)
	case "up":
		if err := netlink.LinkSetUp(iface); err != nil {
			return fmt.Errorf("%v can't make it up: %v", iface.Attrs().Name, err)
		}
	case "down":
		if err := netlink.LinkSetDown(iface); err != nil {
			return fmt.Errorf("%v can't make it down: %v", iface.Attrs().Name, err)
		}
	default:
		return usage()
	}
	return nil
}

func link() error {
	if len(arg) == 1 {
		return linkshow()
	}

	cursor++
	whatIWant = []string{"show", "set"}
	cmd := arg[cursor]

	switch one(cmd, whatIWant) {
	case "show":
		return linkshow()
	case "set":
		return linkset()
	}
	return usage()
}

func routeshow() error {
	path := "/proc/net/route"
	if *inet6 {
		path = "/proc/net/ipv6_route"
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Route show failed: %v", err)
	}
	// Need to formation this better.
	if !(*inet6) {
	/*
	ip: Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT
  ens33	00000000	02C210AC	0003	0	0	100	00000000	0	0	0
  ens33	0000FEA9	00000000	0001	0	0	1000	0000FFFF	0	0	0
  ens33	00C210AC	00000000	0001	0	0	100	00FFFFFF	0	0	0

	want to turn it into something like this

	eoinokane@ubuntu:~/.gvm/pkgsets/go1.12/global/src/github.com/u-root/u-root/cmds/core/ip$ ip route
	default via 172.16.194.2 dev ens33 proto dhcp metric 100
	169.254.0.0/16 dev ens33 scope link metric 1000
	172.16.194.0/24 dev ens33 proto kernel scope link src 172.16.194.129 metric 100

	*/
		// Split the string by new line (Assumed format above)
		rows := strings.Split(string(b),"\n")

		// May be over kill but check for the format of the headerpattern for ubunto (is this this same
		// for all linux implementations??? No clue, need to check with Andrea)
		ubunto_headerpattern,_  := regexp.Compile("Iface	Destination	Gateway 	Flags	RefCnt	Use	Metric	Mask		MTU	Window	IRTT")
		ubunto_bodypattern,_ := regexp.Compile(("([a-z])\\w+	([0-9A-F]{8})	([0-9A-F]{8})	(000[1-3])	(\\d){1}	(\\d){1}	(\\d){3,4}	([0-9A-F]{8})	(\\d){1}	(\\d){1}	(\\d){1}"))
		ubunto_eofpattern,_ := regexp.Compile("^$")

		// check for the ubunto to line for ipv4
		if matched:= ubunto_headerpattern.MatchString(rows[0]); matched {

			// Skip the top line, (already verified with regex)
			for v :=1; v < len(rows); v++ {

				if matched = ubunto_eofpattern.MatchString(rows[v]); matched  {
						// Got to the EOF, exit cleanly
						return nil
				} else if matched = ubunto_bodypattern.MatchString(rows[v]); matched {

					// Get the cols, can make some assumptions based on the regex match used
					cols := strings.Split(strings.Trim(rows[v]," "), "\t")

					// Create a simple array to hold the output
					var o []string
					// Get the source ip address, it's in hex format.
					if src_address, err := hextoipaddress(cols[1]); err!=nil {
						flog.Printf("Cannot decode the source ip address %v",err)
						return err
					} else {
						o = append(o, src_address)
						if gateway_address, err := hextoipaddress(cols[2]); err!=nil {
							flog.Printf("Cannot decode the gateway ip address %v",err)
							return err
						}	else {
							o = append(o, gateway_address)
							if mask_address,err := hextoipaddress(cols[7]); err!=nil {
								flog.Printf("Cannot decode the mask ip address %v",err)
								return err
							} else {
								mask_address := net.IPMask(net.ParseIP(mask_address).To4())
								mask_size, _ := mask_address.Size()
								if mask_size != 0 {
										o[0] += "/"
										o[0] += strconv.Itoa(mask_size)
								}
								o = append(o, cols[0])
								o = append(o, cols[6])
							}
						}
					}
					flog.Printf("%s via %s ?dev? %s ?proto? ?dhcp? metric %s",o[0],o[1],o[2],o[3])
					/*
					Flags := cols[3]
					RefCnt := cols[4]
					Use := cols[5]
					MTU := cols[8]
					Window := cols[9]
					IRTT := cols[10]
					*/
					/*
					0  ip: col content 'ens33'
					1  ip: col content '00C210AC'
					2  ip: col content '00000000'
					3  ip: col content '0001'
					4  ip: col content '0'
					5  ip: col content '0'
					6  ip: col content '100'
					7  ip: col content '00FFFFFF'
					8  ip: col content '0'
					9  ip: col content '0'
					10 ip: col content '0'
					*/
				} else {
						flog.Printf("error matching body %t %v", matched, err)
						return nil
				}
			}
		} else {
				flog.Printf("error matching header %t %v", matched, err)
				return nil
		}
	} else {
		flog.Printf("%s", string(b))
	}
	return nil
}

func hextoipaddress(src_address_hex string) (string, error) {
	if a, err := hex.DecodeString(src_address_hex); err!=nil {
		flog.Printf("hextoipaddress(%s) Cannot decode the source ip address %v",src_address_hex, err)
		return string(""), err
	} else {
		/*
		for i:=0; i < len(a); i++ {
			flog.Printf("hexttoipaddress(%s) a[%d]='%s'",src_address_hex,i,strconv.Itoa(int(a[i])))
		}
		*/
		src_address_string := fmt.Sprintf("%v.%v.%v.%v",strconv.Itoa(int(a[3])), strconv.Itoa(int(a[2])), strconv.Itoa(int(a[1])), strconv.Itoa(int(a[0])))

		// Check if it's the default address, am I being overly clever
		if ip := net.ParseIP(src_address_string); ip.Equal(net.IPv4zero) {
			return "default", nil
		} else {
			return src_address_string, nil
		}
	}
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
		return "", nil, fmt.Errorf("Gateway CIDR: %v", err)
	}
	return nh, addr, nil
}

func routeadddefault() error {
	nh, nhval, err := nexthop()
	if err != nil {
		return err
	}
	// TODO: NHFLAGS.
	l, err := dev()
	if err != nil {
		return err
	}
	switch nh {
	case "via":
		flog.Printf("Add default route %v via %v", nhval, l.Attrs().Name)
		r := &netlink.Route{LinkIndex: l.Attrs().Index, Gw: nhval.IPNet.IP}
		if err := netlink.RouteAdd(r); err != nil {
			return fmt.Errorf("error adding default route to %v: %v", l.Attrs().Name, err)
		}
		return nil
	}
	return usage()
}

func routeadd() error {
	ns := nodespec()
	switch ns {
	case "default":
		return routeadddefault()
	default:
		addr, err := netlink.ParseAddr(arg[cursor])
		if err != nil {
			return usage()
		}
		d, err := dev()
		if err != nil {
			return usage()
		}
		r := &netlink.Route{LinkIndex: d.Attrs().Index, Dst: addr.IPNet}
		if err := netlink.RouteAdd(r); err != nil {
			return fmt.Errorf("error adding route %s -> %s: %v", addr, d.Attrs().Name, err)
		}
		return nil
	}
}

func route() error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return routeshow()
	}

	whatIWant = []string{"show", "add"}
	switch one(arg[cursor], whatIWant) {
	case "show":
		return routeshow()
	case "add":
		return routeadd()
	}
	return usage()
}

func main() {
	// When this is embedded in busybox we need to reinit some things.
	whatIWant = []string{"addr", "route", "link", "neigh"}
	cursor = 0
	flag.Parse()
	arg = flag.Args()

	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				flog.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				flog.Fatalf("Args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			}
			flog.Fatalf("Bummer: %v", err)
		default:
			flog.Fatalf("unexpected panic value: %T(%v)", err, err)
		}
	}()

	// The ip command doesn't actually follow the BNF it prints on error.
	// There are lots of handy shortcuts that people will expect.
	var err error
	switch one(arg[cursor], whatIWant) {
	case "addr":
		err = addrip()
	case "link":
		err = link()
	case "route":
		err = route()
	case "neigh":
		err = neigh()
	default:
		err = usage()
	}
	if err != nil {
		flog.Fatal(err)
	}
}
