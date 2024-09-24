// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"

	"github.com/gopacket/gopacket"
)

func (cmd cmd) ethernetInfo(ethernetLayer gopacket.LinkLayer, networkLayer gopacket.NetworkLayer) string {
	if !cmd.Opts.Ether {
		return fmt.Sprintf("%s %s", networkLayer.NetworkFlow().EndpointType(), cmd.Opts.Device)
	}

	src, dst := ethernetLayer.LinkFlow().Endpoints()
	dstHost := dst.String()

	if dstHost == "ff:ff:ff:ff:ff:ff" {
		dstHost = "Broadcast"
	}

	length := len(ethernetLayer.LayerContents()) + len(ethernetLayer.LayerPayload())

	return fmt.Sprintf("%s > %s, ethertype %s, length %d:", src, dstHost, networkLayer.NetworkFlow().EndpointType(), length)
}
