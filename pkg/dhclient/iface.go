// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/vishvananda/netlink"
)

// FilterBondedInterfaces takes a slice of links, checks to see if any of
// them are already bonded to a bond, and returns the ones that aren't bonded.
func FilterBondedInterfaces(ifs []netlink.Link, verbose bool) []netlink.Link {
	t := []netlink.Link{}
	for _, iface := range ifs {
		n := iface.Attrs().Name
		if _, err := os.Stat(filepath.Join("/sys/class/net", n, "master")); err != nil {
			// if the master symlink does not exist, it probably means link isn't bonded to a bond.
			t = append(t, iface)
		} else if verbose {
			log.Printf("skipping bonded interface %v", n)
		}
	}
	return t
}

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
		return nil, fmt.Errorf("can not get list of link names: %w", err)
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
