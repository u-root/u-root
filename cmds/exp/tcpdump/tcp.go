// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gopacket/gopacket/layers"
)

// tcpData returns a string representation of the TCP layer.
func tcpData(layer *layers.TCP, length int, verbose, quiet bool) string {
	flags := tcpFlags(*layer)
	opts := tcpOptions(layer.Options)

	switch {
	case quiet:
		return fmt.Sprintf("TCP, length %d", length)

	case verbose:
		return fmt.Sprintf("Flags [%s], cksum 0x%x, seq %d, ack %d, win %d, options [%s], length %d",
			flags, layer.Checksum, layer.Seq, layer.Ack, layer.Window, opts, length)

	default:
		return fmt.Sprintf("Flags [%s], seq %d, ack %d, win %d, options [%s], length %d",
			flags, layer.Seq, layer.Ack, layer.Window, opts, length)
	}
}

// tcpFlags returns a string representation of the TCP flags.
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

// tcpOptions returns a string representation of the TCP options.
func tcpOptions(options []layers.TCPOption) string {
	var opts string

	for _, opt := range options {
		opts += tcpOptionToString(opt) + ","
	}

	return strings.TrimRight(opts, ",")
}

// tcpOptionToString returns a string representation of the TCP option.
func tcpOptionToString(opt layers.TCPOption) string {
	if opt.OptionType == layers.TCPOptionKindMSS && len(opt.OptionData) == 2 {
		return fmt.Sprintf("%s val %v",
			opt.OptionType,
			binary.BigEndian.Uint16(opt.OptionData))
	}

	if opt.OptionType == layers.TCPOptionKindTimestamps && len(opt.OptionData) == 8 {
		return fmt.Sprintf("%s val %v",
			opt.OptionType,
			binary.BigEndian.Uint32(opt.OptionData[:4]))
	}

	return opt.OptionType.String()
}
