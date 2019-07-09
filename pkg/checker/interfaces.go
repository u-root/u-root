// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/insomniacslk/dhcp/netboot"
	"github.com/safchain/ethtool"
)

func init() {
	registerCheckFun(InterfaceExists)
	registerCheckFun(LinkSpeed)
	registerCheckFun(LinkAutoneg)
	registerCheckFun(InterfaceCanDoDHCPv6)
	registerCheckFun(InterfaceHasLinkLocalAddress)
	registerCheckFun(InterfaceHasGlobalAddresses)
	registerCheckFun(InterfaceRemediate)
}

// InterfaceExists verifies if an interface is present on the system
func InterfaceExists(args CheckArgs) error {
	_, err := net.InterfaceByName(args["ifname"])
	return err
}

func ethStats(ifname string) (*ethtool.EthtoolCmd, error) {
	cmd := ethtool.EthtoolCmd{}
	_, err := cmd.CmdGet(ifname)
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

// LinkSpeed checks the link speed, and complains if smaller than `min`
// megabit/s.
func LinkSpeed(args CheckArgs) error {
	minSpeed, err := strconv.Atoi(args["minSpeed"])
	if err != nil || minSpeed <= 0 {
		return fmt.Errorf("Invalid value for 'minSpeed' argument: %v", args["minSpeed"])
	}

	eth, err := ethStats(args["ifname"])
	if err != nil {
		return err
	}
	if int(eth.Speed) < minSpeed {
		return fmt.Errorf("link speed %d < %d", eth.Speed, minSpeed)
	}
	return nil
}

// LinkAutoneg checks if the link auto-negotiation state, and return an error if
// it's not the expected state.
func LinkAutoneg(args CheckArgs) error {
	expected, err := strconv.ParseBool(args["expected"])
	if err != nil {
		return err
	}
	eth, err := ethStats(args["ifname"])
	if err != nil {
		return err
	}
	var want uint8
	if expected {
		want = 1
	}
	if eth.Autoneg != want {
		return fmt.Errorf("link autoneg %d; want %d", eth.Autoneg, want)
	}
	return nil
}

func addresses(ifname string) ([]net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	iplist := make([]net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			return nil, errors.New("not a net.IPNet")
		}
		iplist = append(iplist, ipnet.IP)
	}
	return iplist, nil
}

// InterfaceHasLinkLocalAddress verifies if an interface has a configured
// link-local address.
func InterfaceHasLinkLocalAddress(args CheckArgs) error {
	addrs, err := addresses(args["ifname"])
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if addr.IsLinkLocalUnicast() {
			return nil
		}
	}
	return fmt.Errorf("no link local addresses for interface %s", args["ifname"])
}

// InterfaceHasGlobalAddresses returns a Checker that verifies if an interface has
// at least one global address.
func InterfaceHasGlobalAddresses(args CheckArgs) error {
	addrs, err := addresses(args["ifname"])
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if addr.IsGlobalUnicast() {
			return nil
		}
	}
	return fmt.Errorf("no unicast global addresses for interface %s", args["ifname"])
}

// InterfaceRemediate returns a Remediator that tries to fix a missing
// interface issue.
func InterfaceRemediate(args CheckArgs) error {
	// TODO implement driver checking logic
	dmesg, err := getDmesg()
	if err != nil {
		return fmt.Errorf("cannot read dmesg to look for NIC driver information: %v", err)
	}
	lines := grep(dmesg, args["ifname"])
	if len(lines) == 0 {
		return fmt.Errorf("no trace of %s in dmesg", args["ifname"])
	}
	// TODO should this be returned as a string to the caller?
	fmt.Printf("  found %d references to %s in dmesg\n", len(lines), args["ifname"])
	return nil
}

// InterfaceCanDoDHCPv6 checks whether DHCPv6 succeeds on an interface, and if
// it has a valid netboot URL.
func InterfaceCanDoDHCPv6(args CheckArgs) error {
	conv, err := netboot.RequestNetbootv6(args["ifname"], 10*time.Second, 2)
	if err != nil {
		return err
	}
	_, _, err = netboot.ConversationToNetconf(conv)
	return err
}
