// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dhcp4 implements constants and marshaling for DHCPv4 packets as
// described by RFC 2131, RFC 2132, and RFC 3396.
package dhcp4

import (
	"errors"
)

var (
	// ErrInvalidOptions is returned when invalid options data is
	// encountered during parsing. The data could report an incorrect
	// length or have trailing bytes which are not part of the option.
	ErrInvalidOptions = errors.New("invalid options data")

	// ErrInvalidPacket is returned when a byte slice does not contain
	// enough data to create a valid Packet.
	ErrInvalidPacket = errors.New("not enough bytes for valid packet")

	// ErrOptionNotPresent is returned when a requested opcode is not in
	// the packet.
	ErrOptionNotPresent = errors.New("option code not present in packet")
)
