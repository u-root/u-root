// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ip manipulates network addresses, interfaces, routing, and other config.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
)

var inet6 bool

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
	return fmt.Errorf("this was fine: '%v', and this was left, '%v', and this was not understood, '%v'; only options are '%v'",
		arg[0:cursor], arg[cursor:], arg[cursor], whatIWant)
}

func run(out io.Writer) error {
	// When this is embedded in busybox we need to reinit some things.
	whatIWant = []string{"address", "route", "link", "neigh"}
	cursor = 0

	defer func() error {
		switch err := recover().(type) {
		case nil:
		case error:
			if strings.Contains(err.Error(), "index out of range") {
				return fmt.Errorf("args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			} else if strings.Contains(err.Error(), "slice bounds out of range") {
				return fmt.Errorf("args: %v, I got to arg %v, I wanted %v after that", arg, cursor, whatIWant)
			}
			return fmt.Errorf("bummer: %v", err)
		default:
			return fmt.Errorf("unexpected panic value: %T(%v)", err, err)
		}
		return nil
	}()

	// The ip command doesn't actually follow the BNF it prints on error.
	// There are lots of handy shortcuts that people will expect.
	var err error
	switch findPrefix(arg[cursor], whatIWant) {
	case "address":
		err = address(out)
	case "link":
		err = link(out)
	case "route":
		err = route(out)
	case "neigh":
		err = neigh(out)
	default:
		err = usage()
	}
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.BoolVar(&inet6, "6", false, "use inet6")
	flag.Parse()
	arg = flag.Args()
	if err := run(os.Stdout); err != nil {
		log.Fatalf("ip: %v", err)
	}
}
