// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lacp provides methods to setup and teardown linux LACP bonds.
package lacp

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

// Follow netbase's bonding settings found in gbuild/netbase/etc/kernel/%uname-r%/netbase/libbond.sh
// These parameters can be found at https://www.kernel.org/doc/Documentation/networking/bonding.txt
// Exception : lacp_rate is fast.
const (
	lacpMiimon         int                        = 100
	lacpUpDelay        int                        = 0
	lacpDownDelay      int                        = 0
	lacpUseCarrier     int                        = 1
	lacpXmitHashPolicy netlink.BondXmitHashPolicy = netlink.BOND_XMIT_HASH_POLICY_LAYER3_4
	lacpRate           netlink.BondLacpRate       = netlink.BOND_LACP_RATE_FAST
	lacpAdSelect       netlink.BondAdSelect       = netlink.BOND_AD_SELECT_BANDWIDTH
)

// RemoveExistingBonds removes all existing bonds, LACP or others.
func RemoveExistingBonds() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}
	for _, l := range links {
		if l.Type() == "bond" {
			if err := netlink.LinkDel(l); err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateLACPBond creates an LACP interface using bondName and adds links, return netlink.Link for the bonded interface.
func CreateLACPBond(links []netlink.Link, bondName string) (netlink.Link, error) {
	if len(links) == 0 {
		return nil, fmt.Errorf("error creating bond %v, no links provided", bondName)
	}

	// Create the bond.
	bl := netlink.NewLinkBond(netlink.NewLinkAttrs())
	bl.Attrs().Name = bondName
	bl.Attrs().HardwareAddr = links[0].Attrs().HardwareAddr // Use the first link's MAC as HWaddr for the bond.
	bl.Mode = netlink.BOND_MODE_802_3AD
	bl.Miimon = lacpMiimon
	bl.UpDelay = lacpUpDelay
	bl.DownDelay = lacpDownDelay
	bl.UseCarrier = lacpUseCarrier
	bl.XmitHashPolicy = lacpXmitHashPolicy
	bl.LacpRate = lacpRate
	bl.AdSelect = lacpAdSelect

	if err := netlink.LinkAdd(bl); err != nil {
		return nil, fmt.Errorf("error creating bond %v: %v", bondName, err)
	}
	// Bring up the bond before adding links.
	if err := netlink.LinkSetUp(bl); err != nil {
		return nil, fmt.Errorf("error bringing up bond %v: %v", bondName, err)
	}
	// Add each link to the bond.
	for _, l := range links {
		if err := netlink.LinkSetMaster(l, bl); err != nil {
			return nil, fmt.Errorf("error adding link %v to bond %v", l.Attrs().Name, bondName)
		}
		if err := netlink.LinkSetUp(l); err != nil {
			return nil, fmt.Errorf("error bringing up link %v: %v", l.Attrs().Name, err)
		}
	}
	return bl, nil
}
