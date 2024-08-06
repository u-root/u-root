// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gopacket/gopacket/layers"
)

func tcpData(layer *layers.TCP, length int) string {
	var data string

	flags := tcpFlags(*layer)
	opts := tcpOptions(layer.Options)

	data = fmt.Sprintf("Flags [%s], seq %d, ack %d, win %d, options [%s], length %d", flags, layer.Seq, layer.Ack, layer.Window, opts, length)

	return data
}

func tcpFlags(layer layers.TCP) string {
	var flags string
	if layer.PSH {
		flags += "P"
	}
	if layer.FIN {
		flags += "F"
	}
	if layer.SYN {
		flags += "S"
	}
	if layer.RST {
		flags += "R"
	}
	if layer.URG {
		flags += "U"
	}
	if layer.ECE {
		flags += "E"
	}
	if layer.CWR {
		flags += "C"
	}
	if layer.NS {
		flags += "N"
	}
	if layer.ACK {
		flags += "."
	}

	return flags
}

func tcpOptions(options []layers.TCPOption) string {
	var opts string

	for _, opt := range options {
		opts += tcpOptionToString(opt) + ","
	}

	return strings.TrimRight(opts, ",")
}

func tcpOptionToString(opt layers.TCPOption) string {
	switch opt.OptionType {
	case layers.TCPOptionKindMSS:
		if len(opt.OptionData) >= 2 {
			return fmt.Sprintf("%s val %v",
				opt.OptionType,
				binary.BigEndian.Uint16(opt.OptionData))
		}

	case layers.TCPOptionKindTimestamps:
		if len(opt.OptionData) == 8 {
			return fmt.Sprintf("%s val %v",
				opt.OptionType,
				binary.BigEndian.Uint32(opt.OptionData[:4]))
		}
	}

	return fmt.Sprintf("%s", opt.OptionType)
}
