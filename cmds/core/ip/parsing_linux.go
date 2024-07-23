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

// tokenRemains returns true if there are more tokens to parse.
func (cmd *cmd) tokenRemains() bool {
	return cmd.cursor < len(cmd.args)-1
}

// currentToken returns the current token.
func (cmd *cmd) currentToken() string {
	return cmd.args[cmd.cursor]
}

// nextToken returns the next token and sets its expected values.
func (cmd *cmd) nextToken(expectedValues ...string) string {
	cmd.expectedValues = expectedValues
	cmd.cursor++

	return cmd.args[cmd.cursor]
}

// lastToken returns the last token and sets its expected values.
func (cmd *cmd) lastToken(expectedValues ...string) string {
	cmd.expectedValues = expectedValues
	cmd.cursor--

	return cmd.args[cmd.cursor]
}

// findPrefix returns the prefix of the next token.
// If the prefix is not found, an empty string is returned.
func (cmd *cmd) findPrefix(expectedValue ...string) string {
	cmd.expectedValues = expectedValue
	cmd.cursor++
	var x, n int

	for i, v := range cmd.expectedValues {
		if strings.HasPrefix(v, cmd.currentToken()) {
			n++
			x = i
		}
	}

	if n == 1 {
		return cmd.expectedValues[x]
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
func (cmd cmd) parseDeviceName(mandatory bool) (netlink.Link, error) {
	switch mandatory {
	case true:
		cmd.cursor++
		cmd.expectedValues = []string{"dev", "device name"}

		if cmd.args[cmd.cursor] == "dev" {
			cmd.cursor++
		}

		cmd.expectedValues = []string{"device name"}
		return netlink.LinkByName(cmd.args[cmd.cursor])
	case false:
		if cmd.cursor == len(cmd.args)-1 {
			return nil, ErrNotFound
		}

		cmd.cursor++
		cmd.expectedValues = []string{"dev", "device name"}

		if cmd.cursor > len(cmd.args)-1 {
			return nil, ErrNotFound
		}

		if cmd.args[cmd.cursor] == "dev" {
			cmd.cursor++

			if cmd.cursor > len(cmd.args)-1 {
				return nil, ErrNotFound
			}

		}

		cmd.expectedValues = []string{"device name"}
		return netlink.LinkByName(cmd.args[cmd.cursor])
	}

	return nil, ErrNotFound
}

// parseType parses the type of the command.
// The type is the next argument after the 'type' keyword.
// The type is optional in some commands, hence an `ErrNotFound` is returned if the type is not found.
func (cmd cmd) parseType() (string, error) {
	if cmd.cursor == len(cmd.args)-1 {
		return "", ErrNotFound
	}

	cmd.cursor++
	cmd.expectedValues = []string{"type"}

	if cmd.cursor > len(cmd.args)-1 {
		return "", ErrNotFound
	}

	if cmd.args[cmd.cursor] != "type" {
		return "", ErrNotFound
	}

	cmd.cursor++

	cmd.expectedValues = []string{"type name"}
	return cmd.args[cmd.cursor], nil
}

func (cmd cmd) parseAddress() (net.IP, error) {
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

func (cmd cmd) parseIPNet() (*net.IPNet, error) {
	token := cmd.nextToken("CIDR")
	_, ipNet, err := net.ParseCIDR(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", token)
	}

	if ipNet == nil {
		return nil, fmt.Errorf("failed to parse CIDR: %v", token)
	}

	return ipNet, nil
}

func (cmd cmd) parseHardwareAddress() (net.HardwareAddr, error) {
	return net.ParseMAC(cmd.nextToken("<MAC-ADDR>"))
}

type Integer interface {
	string | uint8 | uint16 | uint32 | uint64 | int | []byte
}

// parseValue parses a value from the cmd.argument list.
// expected is the string that represents the expected value.
// allowed types are string, []byte, int, uint16, uint32, uint64, uint8.
func parseValue[T Integer](cmd cmd, expected string) (val T, err error) {
	token := cmd.nextToken(expected)

	var value interface{}

	switch any(val).(type) {
	case string:
		return any(token).(T), nil
	case []byte:
		value, err = hex.DecodeString(token)
		if err != nil {
			return val, fmt.Errorf("failed to parse hex: %v", err)
		}

		return any(val).(T), nil
	case int:
		value, err = strconv.Atoi(token)
		if err != nil {
			return val, fmt.Errorf("failed to parse integer: %v", err)
		}

		return any(val).(T), nil
	case uint16:
		value, err = strconv.ParseUint(token, 10, 16)
		if err != nil {
			return val, fmt.Errorf("failed to parse integer: %v", err)
		}
	case uint32:
		value, err = strconv.ParseUint(token, 10, 32)
		if err != nil {
			return val, fmt.Errorf("failed to parse integer: %v", err)
		}
	case uint64:
		value, err = strconv.ParseUint(token, 10, 64)
		if err != nil {
			return val, fmt.Errorf("failed to parse integer: %v", err)
		}
	default:
		value, err = strconv.ParseUint(token, 10, 8)
		if err != nil {
			return val, fmt.Errorf("failed to parse integer: %v", err)
		}
	}

	return any(value).(T), nil
}

// parseBool parses a boolean value from the cmd.argument list.
// expectedTrue and expectedFalse are the strings that represent true and false.
func (cmd cmd) parseBool(expectedTrue, expectedFalse string) (bool, error) {
	switch cmd.nextToken([]string{expectedTrue, expectedFalse}...) {
	case expectedTrue:
		return true, nil
	case expectedFalse:
		return false, nil
	}

	return false, fmt.Errorf("invalid bool value: %v", cmd.args[cmd.cursor])
}

func (cmd cmd) parseName() string {
	nextToken := cmd.nextToken("name", "device name")

	if nextToken == "name" {
		nextToken = cmd.nextToken("device name")
	}

	return nextToken
}

func (cmd cmd) parseNodeSpec() string {
	return cmd.nextToken("default", "CIDR")
}

func (cmd cmd) parseNextHop() (string, net.IP, error) {
	nh := cmd.nextToken("via")
	if nh != "via" {
		return "", nil, cmd.usage()
	}

	addr := net.ParseIP(cmd.nextToken("Gateway CIDR"))
	if addr == nil {
		return "", nil, fmt.Errorf("failed to parse gateway IP: %v", cmd.args[cmd.cursor])
	}

	return nh, addr, nil
}
