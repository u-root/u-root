// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"testing"

	"github.com/gopacket/gopacket/layers"
)

func TestTcpFlags(t *testing.T) {
	tests := []struct {
		name     string
		layer    layers.TCP
		expected string
	}{
		{
			name:     "No flags set",
			layer:    layers.TCP{},
			expected: "",
		},
		{
			name:     "PSH flag set",
			layer:    layers.TCP{PSH: true},
			expected: "P",
		},
		{
			name:     "FIN flag set",
			layer:    layers.TCP{FIN: true},
			expected: "F",
		},
		{
			name:     "SYN flag set",
			layer:    layers.TCP{SYN: true},
			expected: "S",
		},
		{
			name:     "RST flag set",
			layer:    layers.TCP{RST: true},
			expected: "R",
		},
		{
			name:     "URG flag set",
			layer:    layers.TCP{URG: true},
			expected: "U",
		},
		{
			name:     "ECE flag set",
			layer:    layers.TCP{ECE: true},
			expected: "E",
		},
		{
			name:     "CWR flag set",
			layer:    layers.TCP{CWR: true},
			expected: "C",
		},
		{
			name:     "NS flag set",
			layer:    layers.TCP{NS: true},
			expected: "N",
		},
		{
			name:     "ACK flag set",
			layer:    layers.TCP{ACK: true},
			expected: ".",
		},
		{
			name:     "Multiple flags set",
			layer:    layers.TCP{PSH: true, SYN: true, ACK: true},
			expected: "PS.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tcpFlags(tt.layer)
			if result != tt.expected {
				t.Errorf("tcpFlags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTcpOptionToString(t *testing.T) {
	tests := []struct {
		name     string
		opt      layers.TCPOption
		expected string
	}{
		{
			name: "MSS option with valid data",
			opt: layers.TCPOption{
				OptionType: layers.TCPOptionKindMSS,
				OptionData: []byte{0x00, 0x64},
			},
			expected: "MSS val 100",
		},
		{
			name: "MSS option with invalid data",
			opt: layers.TCPOption{
				OptionType: layers.TCPOptionKindMSS,
				OptionData: []byte{0x00},
			},
			expected: "MSS",
		},
		{
			name: "Timestamps option with valid data",
			opt: layers.TCPOption{
				OptionType: layers.TCPOptionKindTimestamps,
				OptionData: []byte{0x00, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x65},
			},
			expected: "Timestamps val 100",
		},
		{
			name: "Timestamps option with invalid data",
			opt: layers.TCPOption{
				OptionType: layers.TCPOptionKindTimestamps,
				OptionData: []byte{0x00, 0x00, 0x00, 0x64},
			},
			expected: "Timestamps",
		},
		{
			name: "Unknown option type",
			opt: layers.TCPOption{
				OptionType: 0xFF,
				OptionData: []byte{},
			},
			expected: "Unknown(255)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tcpOptionToString(tt.opt)
			if result != tt.expected {
				t.Errorf("tcpOptionToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTcpOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []layers.TCPOption
		expected string
	}{
		{
			name:     "No options",
			options:  []layers.TCPOption{},
			expected: "",
		},
		{
			name: "Single MSS option",
			options: []layers.TCPOption{
				{
					OptionType: layers.TCPOptionKindMSS,
					OptionData: []byte{0x00, 0x64},
				},
			},
			expected: "MSS val 100",
		},
		{
			name: "Multiple options",
			options: []layers.TCPOption{
				{
					OptionType: layers.TCPOptionKindMSS,
					OptionData: []byte{0x00, 0x64},
				},
				{
					OptionType: layers.TCPOptionKindTimestamps,
					OptionData: []byte{0x00, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x65},
				},
			},
			expected: "MSS val 100,Timestamps val 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tcpOptions(tt.options)
			if result != tt.expected {
				t.Errorf("tcpOptions() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTcpData(t *testing.T) {
	tests := []struct {
		name     string
		layer    *layers.TCP
		length   int
		verbose  bool
		quiet    bool
		expected string
	}{
		{
			name: "Quiet mode",
			layer: &layers.TCP{
				Seq:      1000,
				Ack:      2000,
				Window:   3000,
				Checksum: 0x1234,
				Options:  []layers.TCPOption{},
			},
			length:   100,
			verbose:  false,
			quiet:    true,
			expected: "TCP, length 100",
		},
		{
			name: "Verbose mode",
			layer: &layers.TCP{
				Seq:      1000,
				Ack:      2000,
				Window:   3000,
				Checksum: 0x1234,
				Options: []layers.TCPOption{
					{
						OptionType: layers.TCPOptionKindMSS,
						OptionData: []byte{0x00, 0x64},
					},
				},
			},
			length:   100,
			verbose:  true,
			quiet:    false,
			expected: "Flags [], cksum 0x1234, seq 1000, ack 2000, win 3000, options [MSS val 100], length 100",
		},
		{
			name: "Default mode",
			layer: &layers.TCP{
				Seq:      1000,
				Ack:      2000,
				Window:   3000,
				Checksum: 0x1234,
				Options: []layers.TCPOption{
					{
						OptionType: layers.TCPOptionKindMSS,
						OptionData: []byte{0x00, 0x64},
					},
				},
			},
			length:   100,
			verbose:  false,
			quiet:    false,
			expected: "Flags [], seq 1000, ack 2000, win 3000, options [MSS val 100], length 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tcpData(tt.layer, tt.length, tt.verbose, tt.quiet)
			if result != tt.expected {
				t.Errorf("tcpData() = %v, want %v", result, tt.expected)
			}
		})
	}
}
