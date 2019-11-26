// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"regexp"

	"github.com/vishvananda/netlink"
)

// Interfaces takes an RE and returns a
// []netlink.Link that matches it, or an error. It is an error
// for the returned list to be empty.
func Interfaces(ifName string) ([]netlink.Link, error) {
	ifRE, err := regexp.CompilePOSIX(ifName)
	if err != nil {
		return nil, err
	}

	ifnames, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("can not get list of link names: %v", err)
	}

	var filteredIfs []netlink.Link
	for _, iface := range ifnames {
		if ifRE.MatchString(iface.Attrs().Name) {
			filteredIfs = append(filteredIfs, iface)
		}
	}

	if len(filteredIfs) == 0 {
		return nil, fmt.Errorf("no interfaces match %s", ifName)
	}
	return filteredIfs, nil
}
