// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

// tokenRemains returns true if there are more tokens to parse.
func (cmd *cmd) tokenRemains() bool {
	return cmd.Cursor < len(cmd.Args)-1
}

// currentToken returns the current token.
func (cmd *cmd) currentToken() string {
	return cmd.Args[cmd.Cursor]
}

// nextToken returns the next token and sets its expected values.
func (cmd *cmd) nextToken(expectedValues ...string) string {
	cmd.ExpectedValues = expectedValues
	cmd.Cursor++
	return cmd.Args[cmd.Cursor]
}

// lastToken returns the last token and sets its expected values.
func (cmd *cmd) lastToken(expectedValues ...string) string {
	cmd.ExpectedValues = expectedValues
	cmd.Cursor--

	return cmd.Args[cmd.Cursor]
}

// peekToken returns the next token without moving the cursor.
func (cmd *cmd) peekToken(expectedValues ...string) string {
	cmd.ExpectedValues = expectedValues

	return cmd.Args[cmd.Cursor+1]
}

// findPrefix returns the prefix of the next token.
// If the prefix is not found, an empty string is returned.
func (cmd *cmd) findPrefix(expectedValue ...string) string {
	token := cmd.nextToken(expectedValue...)
	var x, n int

	for i, v := range cmd.ExpectedValues {
		if strings.HasPrefix(v, token) {
			n++
			x = i
		}
	}

	if n == 1 {
		return cmd.ExpectedValues[x]
	}

	return ""
}

var ErrNotFound = fmt.Errorf("not found")

// in the ip command, turns out 'dev' is a noisy word.
// The BNF it shows is not right in that case.
// Some commands require 'dev' to be present, some don't.
// If 'dev' is present, it is skipped.
// If no device name is present, an error ErrNotFound is returned.
// The mandatory flag will make sure that the program will panic if the device name is not found.
func (cmd *cmd) parseDeviceName(mandatory bool) (netlink.Link, error) {
	switch mandatory {
	case true:
		if cmd.nextToken("dev", "device-name") == "dev" {
			cmd.Cursor++
		}

		cmd.ExpectedValues = []string{"device-name"}
		return netlink.LinkByName(cmd.currentToken())
	default:
		if !cmd.tokenRemains() {
			return nil, ErrNotFound
		}

		if cmd.nextToken("dev", "device-name") == "dev" {
			cmd.Cursor++
		}

		cmd.ExpectedValues = []string{"device-name"}
		return netlink.LinkByName(cmd.currentToken())
	}
}

func (cmd *cmd) parseAddress() (net.IP, error) {
	token := cmd.nextToken("address", "PREFIX")
	if token == "address" {
		token = cmd.nextToken("PREFIX")
	}

	ip := net.ParseIP(token)

	if ip == nil {
		return nil, fmt.Errorf("failed to parse address: %v", token)
	}

	return ip, nil
}

func (cmd *cmd) parseIPNet() (*net.IPNet, error) {
	token := cmd.nextToken("CIDR")
	_, ipNet, err := net.ParseCIDR(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", token)
	}

	return ipNet, nil
}

func (cmd *cmd) parseAddressorCIDR() (net.IP, *net.IPNet, error) {
	addrStr := cmd.nextToken("PREFIX")

	// Check if it's a CIDR notation
	if strings.Contains(addrStr, "/") {
		ip, ipNet, err := net.ParseCIDR(addrStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse address: %s", addrStr)
		}
		return ip, ipNet, nil
	}

	// Regular IP address
	ip := net.ParseIP(addrStr)
	if ip == nil {
		return nil, nil, fmt.Errorf("failed to parse address: %s", addrStr)
	}
	return ip, nil, nil
}

func (cmd *cmd) parseHardwareAddress() (net.HardwareAddr, error) {
	return net.ParseMAC(cmd.nextToken("<MAC-ADDR>"))
}

func (cmd *cmd) parseByte(expected ...string) ([]byte, error) {
	return hex.DecodeString(cmd.nextToken(expected...))
}

func (cmd *cmd) parseInt(expected ...string) (int, error) {
	val, err := strconv.ParseInt(cmd.nextToken(expected...), 10, 0)

	return int(val), err
}

func (cmd *cmd) parseUint8(expected ...string) (uint8, error) {
	val, err := strconv.ParseInt(cmd.nextToken(expected...), 10, 8)

	return uint8(val), err
}

func (cmd *cmd) parseUint16(expected ...string) (uint16, error) {
	val, err := strconv.ParseInt(cmd.nextToken(expected...), 10, 16)

	return uint16(val), err
}

func (cmd *cmd) parseUint32(expected ...string) (uint32, error) {
	val, err := strconv.ParseInt(cmd.nextToken(expected...), 10, 32)

	return uint32(val), err
}

func (cmd *cmd) parseUint64(expected ...string) (uint64, error) {
	val, err := strconv.ParseInt(cmd.nextToken(expected...), 10, 64)

	return uint64(val), err
}

// parseBool parses a boolean value from the cmd.argument list.
// expectedTrue and expectedFalse are the strings that represent true and false.
func (cmd *cmd) parseBool(expectedTrue, expectedFalse string) (bool, error) {
	switch cmd.nextToken([]string{expectedTrue, expectedFalse}...) {
	case expectedTrue:
		return true, nil
	case expectedFalse:
		return false, nil
	}

	return false, fmt.Errorf("invalid bool value: %v", cmd.Args[cmd.Cursor])
}

func (cmd *cmd) parseName() string {
	nextToken := cmd.nextToken("name", "device-name")
	if nextToken == "name" {
		nextToken = cmd.nextToken("device-name")
	}

	return nextToken
}

func (cmd *cmd) parseNextHop() (string, net.IP, error) {
	nh := cmd.nextToken("via")
	if nh != "via" {
		return "", nil, cmd.usage()
	}

	addr := net.ParseIP(cmd.nextToken("Gateway CIDR"))
	if addr == nil {
		return "", nil, fmt.Errorf("failed to parse gateway IP: %v", cmd.Args[cmd.Cursor])
	}

	return nh, addr, nil
}
