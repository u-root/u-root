// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
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

var ErrNotFound = fmt.Errorf("not found")

// in the ip command, turns out 'dev' is a noisy word.
// The BNF it shows is not right in that case.
// Always make 'dev' optional.
func parseDeviceName(mandatory bool) (netlink.Link, error) {
	switch mandatory {
	case true:
		cursor++
		whatIWant = []string{"dev", "device name"}

		if arg[cursor] == "dev" {
			cursor++
		}

		whatIWant = []string{"device name"}
		return netlink.LinkByName(arg[cursor])
	case false:
		if cursor == len(arg)-1 {
			return nil, ErrNotFound
		}

		cursor++
		whatIWant = []string{"dev", "device name"}

		if cursor > len(arg)-1 {
			return nil, ErrNotFound
		}

		if arg[cursor] == "dev" {
			cursor++

			if cursor > len(arg)-1 {
				return nil, ErrNotFound
			}

		}

		whatIWant = []string{"device name"}
		return netlink.LinkByName(arg[cursor])
	}

	return nil, ErrNotFound
}

func parseType() (string, error) {
	if cursor == len(arg)-1 {
		return "", ErrNotFound
	}

	cursor++
	whatIWant = []string{"type"}

	if cursor > len(arg)-1 {
		return "", ErrNotFound
	}

	if arg[cursor] != "type" {
		return "", ErrNotFound
	}

	cursor++

	whatIWant = []string{"type name"}
	return arg[cursor], nil
}

func parseAddress() (net.IP, error) {
	cursor++
	whatIWant = []string{"[address] PREFIX"}

	if arg[cursor] == "address" {
		cursor++
	}

	whatIWant = []string{"PREFIX"}
	ip := net.ParseIP(arg[cursor])

	if ip == nil {
		return nil, fmt.Errorf("failed to parse adress: %v", arg[cursor])
	}
	return net.ParseIP(arg[cursor]), nil
}

func parseIPNet() (*net.IPNet, error) {
	cursor++
	whatIWant = []string{"ADDR/PLEN"}
	_, ipNet, err := net.ParseCIDR(arg[cursor])
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", arg[cursor])
	}

	if ipNet == nil {
		return nil, fmt.Errorf("failed to parse adress: %v", arg[cursor])
	}

	return ipNet, nil
}

func parseHardwareAddress() (net.HardwareAddr, error) {
	cursor++
	whatIWant = []string{"<mac address>"}

	return net.ParseMAC(arg[cursor])
}

func parseString(expected string) string {
	cursor++
	whatIWant = []string{expected}

	return arg[cursor]
}

func parseByte(expected string) ([]byte, error) {
	cursor++
	whatIWant = []string{expected}

	return hex.DecodeString(arg[cursor])
}

func parseInt(expected string) (int, error) {
	cursor++
	whatIWant = []string{expected}

	return strconv.Atoi(arg[cursor])
}

func parseUint8(expected string) (uint8, error) {
	cursor++
	whatIWant = []string{expected}

	val, err := strconv.ParseUint(arg[cursor], 10, 8)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint8: %v", err)
	}

	return uint8(val), nil
}

func parseUint16(expected string) (uint16, error) {
	cursor++
	whatIWant = []string{expected}

	val, err := strconv.ParseUint(arg[cursor], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint16: %v", err)
	}

	return uint16(val), nil
}

func parseUint32(expected string) (uint32, error) {
	cursor++
	whatIWant = []string{expected}

	val, err := strconv.ParseUint(arg[cursor], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint32: %v", err)
	}

	return uint32(val), nil
}

func parseUint64(expected string) (uint64, error) {
	cursor++
	whatIWant = []string{expected}

	return strconv.ParseUint(arg[cursor], 10, 64)
}

func parseBool() (bool, error) {
	cursor++
	whatIWant = []string{"true", "false"}

	switch arg[cursor] {
	case "on":
		return true, nil
	case "off":
		return false, nil
	}

	return false, fmt.Errorf("invalid bool value: %v", arg[cursor])
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
