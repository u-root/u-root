// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"strings"
	"testing"
	"time"
)

func TestWellKnownPorts(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		port     string
		expected string
	}{
		{
			name:     "HTTP port with name",
			cmd:      cmd{Opts: flags{Numerical: false}},
			port:     "80",
			expected: "http",
		},
		{
			name:     "HTTPS port with name",
			cmd:      cmd{Opts: flags{Numerical: false}},
			port:     "443",
			expected: "https",
		},
		{
			name:     "Unknown port with numerical",
			cmd:      cmd{Opts: flags{Numerical: true}},
			port:     "8080",
			expected: "8080",
		},
		{
			name:     "Known port with numerical",
			cmd:      cmd{Opts: flags{Numerical: true}},
			port:     "22",
			expected: "22",
		},
		{
			name:     "Unknown port without numerical",
			cmd:      cmd{Opts: flags{Numerical: false}},
			port:     "666",
			expected: "666",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.wellKnownPorts(tt.port)
			if result != tt.expected {
				t.Errorf("wellKnownPorts(%v) = %v, want %v", tt.port, result, tt.expected)
			}
		})
	}
}

func TestParseTimeStamp(t *testing.T) {
	tests := []struct {
		name             string
		cmd              cmd
		currentTimestamp time.Time
		lastTimeStamp    time.Time
		expected         string
	}{
		{
			name:             "Option -t",
			cmd:              cmd{Opts: flags{T: true}},
			currentTimestamp: time.Unix(1633072800, 0),
			lastTimeStamp:    time.Unix(1633072800, 0),
			expected:         "",
		},
		{
			name:             "Option -tt",
			cmd:              cmd{Opts: flags{TT: true}},
			currentTimestamp: time.Unix(1633072800, 0),
			lastTimeStamp:    time.Unix(1633072700, 0),
			expected:         "1633072800",
		},
		{
			name:             "Option -ttt with nano",
			cmd:              cmd{Opts: flags{TTT: true, TimeStampInNanoSeconds: true}},
			currentTimestamp: time.Unix(1633072800, 0),
			lastTimeStamp:    time.Unix(1633072700, 0),
			expected:         "00:01:40.000000000",
		},
		{
			name:             "Option -ttt with micro",
			cmd:              cmd{Opts: flags{TTT: true, TimeStampInNanoSeconds: false}},
			currentTimestamp: time.Unix(1633072800, 0),
			lastTimeStamp:    time.Unix(1633072700, 0),
			expected:         "00:01:40.000000",
		},
		{
			name:             "Option -tttt",
			cmd:              cmd{Opts: flags{TTTT: true}},
			currentTimestamp: time.Unix(1633072800, 0),
			lastTimeStamp:    time.Unix(1633072700, 0),
			expected:         "00:01:40",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.parseTimeStamp(tt.currentTimestamp, tt.lastTimeStamp)
			if result != tt.expected {
				t.Errorf("parseTimeStamp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatPacketData(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11}
	expected := `0x0000:  0001 0203 0405 0607 0809 0a0b 0c0d 0e0f 
0x0010:  1011 
`

	result := formatPacketData(data)
	trimmedExpected := strings.TrimSpace(expected)
	trimmedResult := strings.TrimSpace(result)

	if trimmedResult != trimmedExpected {
		t.Errorf("Expected:\n%s\nGot:\n%s", trimmedExpected, trimmedResult)
	}
}
