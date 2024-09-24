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

func TestEthernetInfo(t *testing.T) {
	tests := []struct {
		name           string
		cmd            cmd
		ethernetLayer  gopacket.LinkLayer
		networkLayer   gopacket.NetworkLayer
		expectedOutput string
	}{
		{
			name: "Ether option disabled",
			cmd: cmd{
				Opts: flags{
					Ether:  false,
					Device: "eth0",
				},
			},
			ethernetLayer: &layers.Ethernet{
				SrcMAC:       []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				DstMAC:       []byte{0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB},
				EthernetType: layers.EthernetTypeIPv4,
			},
			networkLayer: &layers.IPv4{
				SrcIP:    []byte{192, 168, 0, 1},
				DstIP:    []byte{192, 168, 0, 2},
				Protocol: layers.IPProtocolTCP,
			},
			expectedOutput: "IPv4 eth0",
		},
		{
			name: "Ether option enabled with broadcast",
			cmd: cmd{
				Opts: flags{
					Ether:  true,
					Device: "eth0",
				},
			},
			ethernetLayer: &layers.Ethernet{
				SrcMAC:       []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				DstMAC:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
				EthernetType: layers.EthernetTypeIPv4,
				Length:       14,
				BaseLayer: layers.BaseLayer{
					Contents: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0x08, 0x00},
				},
			},
			networkLayer: &layers.IPv4{
				SrcIP:    []byte{192, 168, 0, 1},
				DstIP:    []byte{192, 168, 0, 2},
				Protocol: layers.IPProtocolTCP,
			},
			expectedOutput: "00:11:22:33:44:55 > Broadcast, ethertype IPv4, length 14:",
		},
		{
			name: "Ether option enabled with unicast",
			cmd: cmd{
				Opts: flags{
					Ether:  true,
					Device: "eth0",
				},
			},
			ethernetLayer: &layers.Ethernet{
				SrcMAC:       []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				DstMAC:       []byte{0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB},
				EthernetType: layers.EthernetTypeIPv4,
				Length:       14,
				BaseLayer: layers.BaseLayer{
					Contents: []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0x08, 0x00},
				},
			},
			networkLayer: &layers.IPv4{
				SrcIP:    []byte{192, 168, 0, 1},
				DstIP:    []byte{192, 168, 0, 2},
				Protocol: layers.IPProtocolTCP,
			},
			expectedOutput: "00:11:22:33:44:55 > 66:77:88:99:aa:bb, ethertype IPv4, length 14:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.ethernetInfo(tt.ethernetLayer, tt.networkLayer)
			if result != tt.expectedOutput {
				t.Errorf("ethernetInfo() = %v, want %v", result, tt.expectedOutput)
			}
		})
	}
}
