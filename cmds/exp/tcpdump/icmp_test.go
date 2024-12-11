// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"testing"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

func TestParseICMP(t *testing.T) {
	tests := []struct {
		name           string
		packet         gopacket.Packet
		expectedOutput string
	}{
		{
			name: "ICMPv4 Echo Request",
			packet: gopacket.NewPacket(
				[]byte{
					0x08, 0x00, 0x4d, 0x3d, 0x00, 0x46, 0x00, 0x01, 0x61, 0x62, 0x63, 0x64,
				},
				layers.LayerTypeICMPv4,
				gopacket.Default,
			),
			expectedOutput: "ICMP EchoRequest, id 70, seq 1, length 12",
		},
		{
			name: "ICMPv6 Echo Request",
			packet: gopacket.NewPacket(
				[]byte{
					0x80, 0x00, 0x4d, 0x3d, 0x1c, 0x46, 0x00, 0x01, 0x61, 0x62, 0x63, 0x64,
				},
				layers.LayerTypeICMPv6,
				gopacket.Default,
			),
			expectedOutput: "ICMP6 EchoRequest, length 12",
		},
		{
			name: "Non-ICMP Packet",
			packet: gopacket.NewPacket(
				[]byte{
					0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
					0xc0, 0xa8, 0x00, 0x01,
				},
				layers.LayerTypeIPv4,
				gopacket.Default,
			),
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseICMP(tt.packet)
			if result != tt.expectedOutput {
				t.Errorf("parseICMP() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
