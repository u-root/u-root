// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

const (
	vrfHelp = `Usage:	ip vrf show [NAME] ...
	ip vrf exec [NAME] cmd ...
	ip vrf identify [PID]
	ip vrf pids [NAME]
`
)

func (cmd cmd) vrf() error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return cmd.vrfShow()
	}

	expectedValues = []string{"show", "help"}
	var c string

	switch c = findPrefix(arg[cursor], expectedValues); c {
	case "show":
		return cmd.vrfShow()
	case "help":
		fmt.Fprint(cmd.out, vrfHelp)

		return nil
	}
	return usage()
}

type Vrf struct {
	Name  string `json:"name"`
	Table uint32 `json:"table"`
}

func (cmd cmd) vrfShow() error {
	links, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	if f.json {
		vrfs := make([]Vrf, 0, len(links))

		for _, link := range links {
			vrf, ok := link.(*netlink.Vrf)
			if !ok {
				continue
			}

			vrfs = append(vrfs, Vrf{
				Name:  vrf.Name,
				Table: vrf.Table,
			})
		}

		return printJSON(cmd.out, vrfs)
	}

	// Print header
	fmt.Fprintln(cmd.out, "Name              Table")
	fmt.Fprintln(cmd.out, "-----------------------")

	for _, link := range links {
		vrf, ok := link.(*netlink.Vrf)
		if !ok {
			continue
		}

		// Adjusted to print both the VRF name and its table ID in the specified format
		fmt.Fprintf(cmd.out, "%-17s %d\n", vrf.Name, vrf.Table)
	}
	return nil
}
