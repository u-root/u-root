// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"net"
	"reflect"
	"testing"
)

func mockCmd(cursor int, args ...string) *cmd {
	return &cmd{
		Args:   args,
		Cursor: cursor,
	}
}

func TestTokenRemains(t *testing.T) {
	c := mockCmd(0, "arg1", "arg2", "arg3")
	if !c.tokenRemains() {
		t.Errorf("Expected tokenRemains to return true, got false")
	}

	c.Cursor = 2
	if c.tokenRemains() {
		t.Errorf("Expected tokenRemains to return false, got true")
	}
}

func TestCurrentToken(t *testing.T) {
	c := mockCmd(0, "arg1", "arg2", "arg3")
	expected := "arg1"
	if token := c.currentToken(); token != expected {
		t.Errorf("Expected currentToken to return %s, got %s", expected, token)
	}
}

func TestNextToken(t *testing.T) {
	c := mockCmd(0, "arg1", "arg2", "arg3")
	expected := "arg2"
	if token := c.nextToken(); token != expected {
		t.Errorf("Expected nextToken to return %s, got %s", expected, token)
	}

	expectedValues := []string{"val1", "val2"}
	c.nextToken(expectedValues...)
	if len(c.ExpectedValues) != len(expectedValues) {
		t.Errorf("Expected nextToken to assign expectedValues %v, got %v", expectedValues, c.ExpectedValues)
	}
}

func TestLastToken(t *testing.T) {
	c := mockCmd(2, "arg1", "arg2", "arg3")

	expected := "arg2"
	if token := c.lastToken(); token != expected {
		t.Errorf("Expected lastToken to return %s, got %s", expected, token)
	}

	expectedValues := []string{"val1", "val2"}
	c.lastToken(expectedValues...)
	if len(c.ExpectedValues) != len(expectedValues) {
		t.Errorf("Expected lastToken to assign expectedValues %v, got %v", expectedValues, c.ExpectedValues)
	}
}

func TestFindPrefix(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		cursor         int
		expectedValues []string
		want           string
	}{
		{
			name:           "No match",
			args:           []string{"cmd", "option1", "value1"},
			cursor:         0,
			expectedValues: []string{"opt", "option2"},
			want:           "",
		},
		{
			name:           "Single match",
			args:           []string{"cmd", "option1", "value1"},
			cursor:         0,
			expectedValues: []string{"opt", "option1"},
			want:           "option1",
		},
		{
			name:           "Multiple matches",
			args:           []string{"cmd", "option", "value1"},
			cursor:         0,
			expectedValues: []string{"option1", "option2"},
			want:           "",
		},
		{
			name:           "Exact match",
			args:           []string{"cmd", "option1", "value1"},
			cursor:         1,
			expectedValues: []string{"value1", "value2"},
			want:           "value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mockCmd(tt.cursor, tt.args...)
			if got := c.findPrefix(tt.expectedValues...); got != tt.want {
				t.Errorf("findPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    net.IP
		wantErr bool
	}{
		{"Invalid address", cmd{Args: []string{"cmd", "address", "invalid"}}, nil, true},
		{"Valid IPv4", cmd{Args: []string{"cmd", "address", "192.168.1.1"}}, net.ParseIP("192.168.1.1"), false},
		{"Valid IPv6", cmd{Args: []string{"cmd", "address", "fe80::1"}}, net.ParseIP("fe80::1"), false},
		{"Valid without address", cmd{Args: []string{"cmd", "fe80::1"}}, net.ParseIP("fe80::1"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIPNet(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    *net.IPNet
		wantErr bool
	}{
		{"Invalid CIDR", cmd{Args: []string{"cmd", "invalid"}}, nil, true},
		{"Valid IPv4 CIDR", cmd{Args: []string{"cmd", "192.168.1.0/24"}}, mustParseCIDR("192.168.1.0/24"), false},
		{"Valid IPv6 CIDR", cmd{Args: []string{"cmd", "fe80::/64"}}, mustParseCIDR("fe80::/64"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseIPNet()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIPNet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIPNet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAddressorCIDR(t *testing.T) {
	tests := []struct {
		name      string
		cmd       cmd
		wantIP    net.IP
		wantIPNet *net.IPNet
		wantErr   bool
	}{
		{
			name:      "Valid IPv4",
			cmd:       cmd{Args: []string{"cmd", "192.168.1.1"}},
			wantIP:    net.ParseIP("192.168.1.1"),
			wantIPNet: nil,
			wantErr:   false,
		},
		{
			name:      "Valid IPv6",
			cmd:       cmd{Args: []string{"cmd", "fe80::1"}},
			wantIP:    net.ParseIP("fe80::1"),
			wantIPNet: nil,
			wantErr:   false,
		},
		{
			name:      "Valid IPv4 CIDR",
			cmd:       cmd{Args: []string{"cmd", "192.168.1.0/24"}},
			wantIP:    net.ParseIP("192.168.1.0"),
			wantIPNet: mustParseCIDR("192.168.1.0/24"),
			wantErr:   false,
		},
		{
			name:      "Valid IPv6 CIDR",
			cmd:       cmd{Args: []string{"cmd", "fe80::/64"}},
			wantIP:    net.ParseIP("fe80::"),
			wantIPNet: mustParseCIDR("fe80::/64"),
			wantErr:   false,
		},
		{
			name:      "Invalid address",
			cmd:       cmd{Args: []string{"cmd", "invalid"}},
			wantIP:    nil,
			wantIPNet: nil,
			wantErr:   true,
		},
		{
			name:      "Invalid CIDR",
			cmd:       cmd{Args: []string{"cmd", "192.168.1.1/invalid"}},
			wantIP:    nil,
			wantIPNet: nil,
			wantErr:   true,
		},
		{
			name:      "Invalid IPv6 CIDR",
			cmd:       cmd{Args: []string{"cmd", "fe80::/invalid"}},
			wantIP:    nil,
			wantIPNet: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIP, gotIPNet, err := tt.cmd.parseAddressorCIDR()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddressorCIDR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIP, tt.wantIP) {
				t.Errorf("parseAddressorCIDR() gotIP = %v, want %v", gotIP, tt.wantIP)
			}
			if !reflect.DeepEqual(gotIPNet, tt.wantIPNet) {
				t.Errorf("parseAddressorCIDR() gotIPNet = %v, want %v", gotIPNet, tt.wantIPNet)
			}
		})
	}
}

// Helper function to parse CIDR for test expectations.
func mustParseCIDR(cidrStr string) *net.IPNet {
	_, cidr, _ := net.ParseCIDR(cidrStr)
	return cidr
}

// Test for parseHardwareAddress
func TestParseHardwareAddress(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    net.HardwareAddr
		wantErr bool
	}{
		{"Valid MAC", cmd{Args: []string{"cmd", "01:23:45:67:89:ab"}}, net.HardwareAddr{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}, false},
		{"Invalid MAC", cmd{Args: []string{"cmd", "invalid"}}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseHardwareAddress()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHardwareAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHardwareAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test for parseByte
func TestParseByte(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    []byte
		wantErr bool
	}{
		{"Valid hex", cmd{Args: []string{"cmd", "deadbeef"}}, []byte{0xde, 0xad, 0xbe, 0xef}, false},
		{"Invalid hex", cmd{Args: []string{"cmd", "xyz"}}, []byte{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    int
		wantErr bool
	}{
		{"Valid integer", cmd{Args: []string{"cmd", "123"}}, 123, false},
		{"Negative integer", cmd{Args: []string{"cmd", "-123"}}, -123, false},
		{"Invalid integer", cmd{Args: []string{"cmd", "abc"}}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUint8(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    uint8
		wantErr bool
	}{
		{"Valid uint8", cmd{Args: []string{"cmd", "123"}}, 123, false},
		{"Invalid uint8", cmd{Args: []string{"cmd", "abc"}}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseUint8()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUint16(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    uint16
		wantErr bool
	}{
		{"Valid uint16", cmd{Args: []string{"cmd", "123"}}, 123, false},
		{"Invalid uint16", cmd{Args: []string{"cmd", "abc"}}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseUint16()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUint32(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    uint32
		wantErr bool
	}{
		{"Valid uint32", cmd{Args: []string{"cmd", "123"}}, 123, false},
		{"Invalid uint32", cmd{Args: []string{"cmd", "abc"}}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseUint32()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUint64(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    uint64
		wantErr bool
	}{
		{"Valid uint64", cmd{Args: []string{"cmd", "123"}}, 123, false},
		{"Invalid uint64", cmd{Args: []string{"cmd", "abc"}}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseUint64()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name          string
		cmd           cmd
		expectedTrue  string
		expectedFalse string
		want          bool
		wantErr       bool
	}{
		{"Valid True", cmd{Args: []string{"cmd", "true"}}, "true", "false", true, false},
		{"Valid False", cmd{Args: []string{"cmd", "false"}}, "true", "false", false, false},
		{"Invalid Value", cmd{Args: []string{"cmd", "maybe"}}, "true", "false", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseBool(tt.expectedTrue, tt.expectedFalse)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseName(t *testing.T) {
	tests := []struct {
		name string
		cmd  cmd
		want string
	}{
		{"Valid Name", cmd{Args: []string{"cmd", "name", "eth0"}}, "eth0"},
		{"Valid Device Name", cmd{Args: []string{"cmd", "eth0"}}, "eth0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cmd.parseName()
			if got != tt.want {
				t.Errorf("parseName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNextHop(t *testing.T) {
	tests := []struct {
		name    string
		cmdArgs []string
		want    string
		wantIP  net.IP
		wantErr bool
	}{
		{"Valid Next Hop", []string{"cmd", "via", "192.168.1.1"}, "via", net.ParseIP("192.168.1.1"), false},
		{"Valid Next Hop", []string{"cmd", "via", "aaa"}, "via", nil, true},

		{"Invalid Next Hop", []string{"cmd", "192.168.1.0/24"}, "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{Args: tt.cmdArgs}
			got, gotIP, err := cmd.parseNextHop()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNextHop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got != tt.want {
					t.Errorf("parseNextHop() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(gotIP, tt.wantIP) {
					t.Errorf("parseNextHop() gotIP = %v, want %v", gotIP, tt.wantIP)
				}
			}
		})
	}
}
