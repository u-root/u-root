// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
)

const monitorHelp = `Usage: ip monitor [ all | OBJECTS ] [ label ]
OBJECTS :=  address | link | mroute | neigh | netconf |
            nexthop | nsid | prefix | route | rule
`

var (
	addressLabel string
	linkLabel    string
	mrouteLabel  string
	neighLabel   string
	netconfLabel string
	nexthopLabel string
	nsidLabel    string
	prefixLabel  string
	routeLabel   string
)

func (cmd cmd) monitor() error {
	addrUpdates := make(chan netlink.AddrUpdate)
	linkUpdates := make(chan netlink.LinkUpdate)
	neighUpdates := make(chan netlink.NeighUpdate)
	routeUpdates := make(chan netlink.RouteUpdate)
	done := make(chan struct{})
	defer close(done)

	// catch signals to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	var singleOptionSelected bool

	for cmd.tokenRemains() {
		token := cmd.nextToken("all", "address", "link", "mroute", "neigh", "netconf", "nexthop", "nsid", "prefix", "route", "rule", "label", "all-nsid")

		switch token {
		case "all":
			if singleOptionSelected {
				return fmt.Errorf("all option can't be used with other options")
			}
		case "address":
			singleOptionSelected = true
			if err := netlink.AddrSubscribe(addrUpdates, done); err != nil {
				return fmt.Errorf("failed to subscribe to address updates: %v", err)
			}
		case "link":
			singleOptionSelected = true
			if err := netlink.LinkSubscribe(linkUpdates, done); err != nil {
				return fmt.Errorf("failed to subscribe to link updates: %v", err)
			}
		case "neigh":
			singleOptionSelected = true
			if err := netlink.NeighSubscribe(neighUpdates, done); err != nil {
				return fmt.Errorf("failed to subscribe to neigh updates: %v", err)
			}
		case "route":
			singleOptionSelected = true
			if err := netlink.RouteSubscribe(routeUpdates, done); err != nil {
				return fmt.Errorf("failed to subscribe to route updates: %v", err)
			}
		case "label":
			addressLabel = "[ADDR]"
			linkLabel = "[LINK]"
			mrouteLabel = "[MROUTE]"
			neighLabel = "[NEIGH]"
			netconfLabel = "[NETCONF]"
			nexthopLabel = "[NEXTHOP]"
			nsidLabel = "[NSID]"
			prefixLabel = "[PREFIX]"
			routeLabel = "[ROUTE]"
		case "mroute", "netconf", "nexthop", "nsid", "prefix", "rule":
			return fmt.Errorf("monitoring %s is not yet supported", cmd.currentToken())
		default:

			return cmd.usage()
		}

	}

	// if either the all option was selected or no option was selected, subscribe to all
	if !singleOptionSelected {
		if err := netlink.AddrSubscribe(addrUpdates, done); err != nil {
			return fmt.Errorf("failed to subscribe to address updates: %v", err)
		}
		if err := netlink.LinkSubscribe(linkUpdates, done); err != nil {
			return fmt.Errorf("failed to subscribe to link updates: %v", err)
		}
		if err := netlink.NeighSubscribe(neighUpdates, done); err != nil {
			return fmt.Errorf("failed to subscribe to neigh updates: %v", err)
		}
		if err := netlink.RouteSubscribe(routeUpdates, done); err != nil {
			return fmt.Errorf("failed to subscribe to route updates: %v", err)
		}
	}

	timestamp := ""

	for {

		if cmd.opts.timeStamp {
			currentTime := time.Now()
			timestamp = currentTime.Format("Timestamp: Mon Jan 2 15:04:05 2006") + fmt.Sprintf(" %06d usec", currentTime.Nanosecond()/1000) + "\n"
		} else if cmd.opts.timeStampShort {
			currentTime := time.Now()
			timestamp = currentTime.Format("[2006-01-02T15:04:05.000000]")
		}

		select {
		case update := <-addrUpdates:

			link, err := netlink.LinkByIndex(update.LinkIndex)
			if err != nil {
				return fmt.Errorf("failed to get link by index %d: %v", update.LinkIndex, err)
			}

			var action string
			if !update.NewAddr {
				action = "Deleted"
			}

			fmt.Fprintf(cmd.out, "%s%s%s %d: %s    %v %v scope %d %v\n", timestamp, addressLabel, action, update.LinkIndex, link.Attrs().Name, ipFamily(update.LinkAddress.IP), update.LinkAddress.String(), update.Scope, link.Attrs().Name)

			validLft := fmt.Sprintf("%v", update.ValidLft)
			preferedLft := fmt.Sprintf("%v", update.PreferedLft)

			if update.ValidLft >= math.MaxInt32 {
				validLft = "forever"
			}

			if update.PreferedLft >= math.MaxInt32 {
				preferedLft = "forever"
			}

			fmt.Fprintf(cmd.out, "    valid_lft %s preferred_lft %s\n", validLft, preferedLft)

		case update := <-neighUpdates:
			var action string

			if update.Type == syscall.RTM_DELNEIGH {
				action = "Deleted "
			}

			link, err := netlink.LinkByIndex(update.Neigh.LinkIndex)
			if err != nil {
				return fmt.Errorf("failed to get link by index %d: %v", update.Neigh.LinkIndex, err)
			}

			fmt.Fprintf(cmd.out, "%s%s%s%s dev %v lladdr %s %v\n", timestamp, neighLabel, action, update.Neigh.IP, link.Attrs().Name, update.Neigh.HardwareAddr.String(), neighStateToString(update.Neigh.State))

		case update := <-routeUpdates:
			var action string
			switch update.Type {
			case syscall.RTM_NEWROUTE:
				action = "Added"
			case syscall.RTM_DELROUTE:
				action = "Deleted"
			}

			link, err := netlink.LinkByIndex(update.Route.LinkIndex)
			if err != nil {
				return fmt.Errorf("failed to get link by index %d: %v", update.Route.LinkIndex, err)
			}

			fmt.Fprintf(cmd.out, "%s%s%s %s dev %s table %d proto %s scope %s src %s\n", timestamp, routeLabel, action, update.Route.Dst, link.Attrs().Name, update.Route.Table, update.Route.Protocol.String(), update.Route.Scope.String(), update.Route.Src)
		case update := <-linkUpdates:
			fmt.Fprintf(cmd.out, "%s%s%d: %s: <%s>\n", timestamp, linkLabel, update.Link.Attrs().Index, update.Link.Attrs().Name, strings.Replace(strings.ToUpper(net.Flags(update.Flags).String()), "|", ",", -1))
			fmt.Fprintf(cmd.out, "    link/%v\n", update.Link.Attrs().EncapType)
		case <-sig:
			return nil
		case <-done:
			return nil
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func neighStateToString(state int) string {
	stateMap := map[int]string{
		0x01: "INCOMPLETE",
		0x02: "REACHABLE",
		0x04: "STALE",
		0x08: "DELAY",
		0x10: "PROBE",
		0x20: "FAILED",
		0x40: "NOARP",
		0x80: "PERMANENT",
	}

	if stateStr, exists := stateMap[state]; exists {
		return stateStr
	}
	return "UNKNOWN"
}

func ipFamily(ip net.IP) string {
	if ip.To4() != nil {
		return "inet"
	}
	return "inet6"
}
