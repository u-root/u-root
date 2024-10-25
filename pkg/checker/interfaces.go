// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/netboot"
	"github.com/safchain/ethtool"
)

// InterfaceExists returns a Checker that verifies if an interface is present on
// the system
func InterfaceExists(ifname string) Checker {
	return func() error {
		_, err := net.InterfaceByName(ifname)
		return err
	}
}

// LinkSpeed checks the link speed, and complains if smaller than `min`
// megabit/s.
func LinkSpeed(ifname string, minSpeed uint32) Checker {
	return func() error {
		ec := ethtool.EthtoolCmd{}
		speed, err := ec.CmdGet(ifname)
		if err != nil {
			return err
		}
		if speed < minSpeed {
			return fmt.Errorf("link speed %d < %d", speed, minSpeed)
		}
		return nil
	}
}

// LinkAutoneg checks if the link auto-negotiation state, and return an error if
// it's not the expected state.
func LinkAutoneg(ifname string, expected bool) Checker {
	return func() error {
		ec := ethtool.EthtoolCmd{}
		_, err := ec.CmdGet(ifname)
		if err != nil {
			return err
		}
		var want uint8
		if expected {
			want = 1
		}
		if ec.Autoneg != want {
			return fmt.Errorf("link autoneg %d; want %d", ec.Autoneg, want)
		}
		return nil
	}
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

// InterfaceHasLinkLocalAddress returns a Checker that verifies if an interface
// has a configured link-local address.
func InterfaceHasLinkLocalAddress(ifname string) Checker {
	return func() error {
		addrs, err := addresses(ifname)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if addr.IsLinkLocalUnicast() {
				return nil
			}
		}
		return fmt.Errorf("no link local addresses for interface %s", ifname)
	}
}

// InterfaceHasGlobalAddresses returns a Checker that verifies if an interface has
// at least one global address.
func InterfaceHasGlobalAddresses(ifname string) Checker {
	return func() error {
		addrs, err := addresses(ifname)
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if addr.IsGlobalUnicast() {
				return nil
			}
		}
		return fmt.Errorf("no unicast global addresses for interface %s", ifname)
	}
}

// InterfaceRemediate returns a Remediator that tries to fix a missing
// interface issue.
func InterfaceRemediate(ifname string) Remediator {
	return func() error {
		// TODO implement driver loading logic
		dmesg, err := getDmesg()
		if err != nil {
			return fmt.Errorf("cannot read dmesg to look for NIC driver information: %w", err)
		}
		lines := grep(dmesg, ifname)
		if len(lines) == 0 {
			return fmt.Errorf("no trace of %s in dmesg", ifname)
		}
		// TODO should this be returned as a string to the caller?
		fmt.Printf("  found %d references to %s in dmesg\n", len(lines), ifname)
		return nil
	}
}

// InterfaceCanDoDHCPv6 checks whether DHCPv6 succeeds on an interface, and if
// it has a valid netboot URL.
func InterfaceCanDoDHCPv6(ifname string) Checker {
	return func() error {
		conv, err := netboot.RequestNetbootv6(ifname, 10*time.Second, 2)
		if err != nil {
			return err
		}
		_, err = netboot.ConversationToNetconf(conv)
		return err
	}
}
