// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

func findPrefix(cmd string, cmds []string) string {
	var x, n int
	for i, v := range cmds {
		if strings.HasPrefix(v, cmd) {
			n++
			x = i
		}
	}
	if n == 1 {
		return cmds[x]
	}
	return ""
}

// in the ip command, turns out 'dev' is a noise word.
// The BNF it shows is not right in that case.
// Always make 'dev' optional.
func parseDeviceName() (netlink.Link, error) {
	cursor++
	whatIWant = []string{"dev", "device name"}
	if arg[cursor] == "dev" {
		cursor++
	}
	whatIWant = []string{"device name"}
	return netlink.LinkByName(arg[cursor])
}

func parseName() (string, error) {
	cursor++
	whatIWant = []string{"name", "device name"}
	if arg[cursor] == "name" {
		cursor++
	}
	whatIWant = []string{"device name"}
	return arg[cursor], nil
}

func parseNodeSpec() string {
	cursor++
	whatIWant = []string{"default", "CIDR"}
	return arg[cursor]
}

func parseNextHop() (string, net.IP, error) {
	cursor++
	whatIWant = []string{"via"}
	if arg[cursor] != "via" {
		return "", nil, usage()
	}
	nh := arg[cursor]
	cursor++
	whatIWant = []string{"Gateway CIDR"}
	addr := net.ParseIP(arg[cursor])
	if addr == nil {
		return "", nil, fmt.Errorf("failed to parse gateway IP: %v", arg[cursor])
	}
	return nh, addr, nil
}
