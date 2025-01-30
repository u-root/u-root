// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

const (
	vrfHelp = `Usage:	ip vrf show [NAME] ...`
)

func (cmd *cmd) vrf() error {
	if !cmd.tokenRemains() {
		return cmd.vrfShow()
	}

	switch cmd.findPrefix("show", "help") {
	case "show":
		return cmd.vrfShow()
	case "help":
		fmt.Fprint(cmd.Out, vrfHelp)

		return nil
	}
	return cmd.usage()
}

// VrfJSON represents a VRF entry for JSON output format.
type VrfJSON struct {
	Name  string `json:"name"`
	Table uint32 `json:"table"`
}

func (cmd *cmd) vrfShow() error {
	links, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	return cmd.printVrf(links)
}

func (cmd *cmd) printVrf(links []netlink.Link) error {
	if cmd.Opts.JSON {
		vrfs := make([]VrfJSON, 0, len(links))

		for _, link := range links {
			vrf, ok := link.(*netlink.Vrf)
			if !ok {
				continue
			}

			vrfs = append(vrfs, VrfJSON{
				Name:  vrf.Name,
				Table: vrf.Table,
			})
		}

		return printJSON(*cmd, vrfs)
	}

	// Print header
	fmt.Fprintln(cmd.Out, "Name              Table")
	fmt.Fprintln(cmd.Out, "-----------------------")

	for _, link := range links {
		vrf, ok := link.(*netlink.Vrf)
		if !ok {
			continue
		}

		// Adjusted to print both the VRF name and its table ID in the specified format
		fmt.Fprintf(cmd.Out, "%-17s %d\n", vrf.Name, vrf.Table)
	}

	return nil
}
