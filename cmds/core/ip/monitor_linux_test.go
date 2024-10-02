// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"net"
	"testing"
)

func TestNeighStateToString(t *testing.T) {
	tests := []struct {
		state    int
		expected string
	}{
		{0x01, "INCOMPLETE"},
		{0x02, "REACHABLE"},
		{0x04, "STALE"},
		{0x08, "DELAY"},
		{0x10, "PROBE"},
		{0x20, "FAILED"},
		{0x40, "NOARP"},
		{0x80, "PERMANENT"},
		{0x00, "UNKNOWN"}, // Test for an undefined state
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := neighStateToString(tt.state)
			if result != tt.expected {
				t.Errorf("neighStateToString(%d) = %s; want %s", tt.state, result, tt.expected)
			}
		})
	}
}

func TestIpFamily(t *testing.T) {
	tests := []struct {
		ip       string
		expected string
	}{
		{"192.168.1.1", "inet"},
		{"10.0.0.1", "inet"},
		{"::1", "inet6"},
		{"fe80::1", "inet6"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			result := ipFamily(ip)
			if result != tt.expected {
				t.Errorf("ipFamily(%s) = %s; want %s", tt.ip, result, tt.expected)
			}
		})
	}
}
