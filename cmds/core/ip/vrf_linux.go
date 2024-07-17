// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/vishvananda/netlink"
)

const (
	vrfHelp = `Usage:	ip vrf show [NAME] ...
	ip vrf exec [NAME] cmd ...
	ip vrf identify [PID]
	ip vrf pids [NAME]
`
)

func vrf(w io.Writer) error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return vrfShow(w)
	}

	whatIWant = []string{"show", "help"}
	var c string

	switch c = findPrefix(arg[cursor], whatIWant); c {
	case "show":
		return vrfShow(w)
	case "help":
		fmt.Fprint(w, vrfHelp)

		return nil
	}
	return usage()
}

func vrfShow(w io.Writer) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	// Print header
	fmt.Fprintln(w, "Name              Table")
	fmt.Fprintln(w, "-----------------------")

	for _, link := range links {
		vrf, ok := link.(*netlink.Vrf)
		if !ok {
			continue
		}

		// Adjusted to print both the VRF name and its table ID in the specified format
		fmt.Fprintf(w, "%-17s %d\n", vrf.Name, vrf.Table)
	}
	return nil
}
