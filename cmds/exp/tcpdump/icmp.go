// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"

	"github.com/gopacket/gopacket"

	"github.com/gopacket/gopacket/layers"
)

// parseICMP parses ICMP packets and returns a string representation of the packet.
func parseICMP(packet gopacket.Packet) string {
	icmpv4 := packet.Layer(layers.LayerTypeICMPv4)
	if icmpv4 != nil {
		layer := icmpv4.(*layers.ICMPv4)

		return fmt.Sprintf("ICMP %s, id %d, seq %d, length %d", layer.TypeCode.String(), layer.Id, layer.Seq, len(layer.Contents)+len(layer.Payload))
	}

	icmpv6 := packet.Layer(layers.LayerTypeICMPv6)
	if icmpv6 != nil {
		layer := icmpv6.(*layers.ICMPv6)

		return fmt.Sprintf("ICMP6 %s, length %d", layer.TypeCode.String(), len(layer.Contents)+len(layer.Payload))
	}

	return ""
}
