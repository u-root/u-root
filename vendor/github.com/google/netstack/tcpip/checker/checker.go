// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package checker provides helper functions to check networking packets for
// validity.
package checker

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
)

// NetworkChecker is a function to check a property of a network packet.
type NetworkChecker func(*testing.T, []header.Network)

// TransportChecker is a function to check a property of a transport packet.
type TransportChecker func(*testing.T, header.Transport)

// IPv4 checks the validity and properties of the given IPv4 packet. It is
// expected to be used in conjunction with other network checkers for specific
// properties. For example, to check the source and destination address, one
// would call:
//
// checker.IPv4(t, b, checker.SrcAddr(x), checker.DstAddr(y))
func IPv4(t *testing.T, b []byte, checkers ...NetworkChecker) {
	ipv4 := header.IPv4(b)

	if !ipv4.IsValid(len(b)) {
		t.Fatalf("Not a valid IPv4 packet")
	}

	xsum := ipv4.CalculateChecksum()
	if xsum != 0 && xsum != 0xffff {
		t.Fatalf("Bad checksum: 0x%x, checksum in packet: 0x%x", xsum, ipv4.Checksum())
	}

	for _, f := range checkers {
		f(t, []header.Network{ipv4})
	}
}

// IPv6 checks the validity and properties of the given IPv6 packet. The usage
// is similar to IPv4.
func IPv6(t *testing.T, b []byte, checkers ...NetworkChecker) {
	ipv6 := header.IPv6(b)
	if !ipv6.IsValid(len(b)) {
		t.Fatalf("Not a valid IPv6 packet")
	}

	for _, f := range checkers {
		f(t, []header.Network{ipv6})
	}
}

// SrcAddr creates a checker that checks the source address.
func SrcAddr(addr tcpip.Address) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if a := h[0].SourceAddress(); a != addr {
			t.Fatalf("Bad source address, got %v, want %v", a, addr)
		}
	}
}

// DstAddr creates a checker that checks the destination address.
func DstAddr(addr tcpip.Address) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if a := h[0].DestinationAddress(); a != addr {
			t.Fatalf("Bad destination address, got %v, want %v", a, addr)
		}
	}
}

// PayloadLen creates a checker that checks the payload length.
func PayloadLen(plen int) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if l := len(h[0].Payload()); l != plen {
			t.Fatalf("Bad payload length, got %v, want %v", l, plen)
		}
	}
}

// FragmentOffset creates a checker that checks the FragmentOffset field.
func FragmentOffset(offset uint16) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		// We only do this of IPv4 for now.
		switch ip := h[0].(type) {
		case header.IPv4:
			if v := ip.FragmentOffset(); v != offset {
				t.Fatalf("Bad fragment offset, got %v, want %v", v, offset)
			}
		}
	}
}

// FragmentFlags creates a checker that checks the fragment flags field.
func FragmentFlags(flags uint8) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		// We only do this of IPv4 for now.
		switch ip := h[0].(type) {
		case header.IPv4:
			if v := ip.Flags(); v != flags {
				t.Fatalf("Bad fragment offset, got %v, want %v", v, flags)
			}
		}
	}
}

// TOS creates a checker that checks the TOS field.
func TOS(tos uint8, label uint32) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if v, l := h[0].TOS(); v != tos || l != label {
			t.Fatalf("Bad TOS, got (%v, %v), want (%v,%v)", v, l, tos, label)
		}
	}
}

// Raw creates a checker that checks the bytes of payload.
// The checker always checks the payload of the last network header.
// For instance, in case of IPv6 fragments, the payload that will be checked
// is the one containing the actual data that the packet is carrying, without
// the bytes added by the IPv6 fragmentation.
func Raw(want []byte) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if got := h[len(h)-1].Payload(); !reflect.DeepEqual(got, want) {
			t.Fatalf("Wrong payload, got %v, want %v", got, want)
		}
	}
}

// IPv6Fragment creates a checker that validates an IPv6 fragment.
func IPv6Fragment(checkers ...NetworkChecker) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		if p := h[0].TransportProtocol(); p != header.IPv6FragmentHeader {
			t.Fatalf("Bad protocol, got %v, want %v", p, header.UDPProtocolNumber)
		}

		ipv6Frag := header.IPv6Fragment(h[0].Payload())
		if !ipv6Frag.IsValid() {
			t.Fatalf("Not a valid IPv6 fragment")
		}

		for _, f := range checkers {
			f(t, []header.Network{h[0], ipv6Frag})
		}
	}
}

// TCP creates a checker that checks that the transport protocol is TCP and
// potentially additional transport header fields.
func TCP(checkers ...TransportChecker) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		first := h[0]
		last := h[len(h)-1]

		if p := last.TransportProtocol(); p != header.TCPProtocolNumber {
			t.Fatalf("Bad protocol, got %v, want %v", p, header.TCPProtocolNumber)
		}

		// Verify the checksum.
		tcp := header.TCP(last.Payload())
		l := uint16(len(tcp))

		xsum := header.Checksum([]byte(first.SourceAddress()), 0)
		xsum = header.Checksum([]byte(first.DestinationAddress()), xsum)
		xsum = header.Checksum([]byte{0, byte(last.TransportProtocol())}, xsum)
		xsum = header.Checksum([]byte{byte(l >> 8), byte(l)}, xsum)
		xsum = header.Checksum(tcp, xsum)

		if xsum != 0 && xsum != 0xffff {
			t.Fatalf("Bad checksum: 0x%x, checksum in segment: 0x%x", xsum, tcp.Checksum())
		}

		// Run the transport checkers.
		for _, f := range checkers {
			f(t, tcp)
		}
	}
}

// UDP creates a checker that checks that the transport protocol is UDP and
// potentially additional transport header fields.
func UDP(checkers ...TransportChecker) NetworkChecker {
	return func(t *testing.T, h []header.Network) {
		last := h[len(h)-1]

		if p := last.TransportProtocol(); p != header.UDPProtocolNumber {
			t.Fatalf("Bad protocol, got %v, want %v", p, header.UDPProtocolNumber)
		}

		udp := header.UDP(last.Payload())
		for _, f := range checkers {
			f(t, udp)
		}
	}
}

// SrcPort creates a checker that checks the source port.
func SrcPort(port uint16) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		if p := h.SourcePort(); p != port {
			t.Fatalf("Bad source port, got %v, want %v", p, port)
		}
	}
}

// DstPort creates a checker that checks the destination port.
func DstPort(port uint16) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		if p := h.DestinationPort(); p != port {
			t.Fatalf("Bad destination port, got %v, want %v", p, port)
		}
	}
}

// SeqNum creates a checker that checks the sequence number.
func SeqNum(seq uint32) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}

		if s := tcp.SequenceNumber(); s != seq {
			t.Fatalf("Bad sequence number, got %v, want %v", s, seq)
		}
	}
}

// AckNum creates a checker that checks the ack number.
func AckNum(seq uint32) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}

		if s := tcp.AckNumber(); s != seq {
			t.Fatalf("Bad ack number, got %v, want %v", s, seq)
		}
	}
}

// Window creates a checker that checks the tcp window.
func Window(window uint16) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}

		if w := tcp.WindowSize(); w != window {
			t.Fatalf("Bad window, got 0x%x, want 0x%x", w, window)
		}
	}
}

// TCPFlags creates a checker that checks the tcp flags.
func TCPFlags(flags uint8) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}

		if f := tcp.Flags(); f != flags {
			t.Fatalf("Bad flags, got 0x%x, want 0x%x", f, flags)
		}
	}
}

// TCPFlagsMatch creates a checker that checks that the tcp flags, masked by the
// given mask, match the supplied flags.
func TCPFlagsMatch(flags, mask uint8) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}

		if f := tcp.Flags(); (f & mask) != (flags & mask) {
			t.Fatalf("Bad masked flags, got 0x%x, want 0x%x, mask 0x%x", f, flags, mask)
		}
	}
}

// TCPSynOptions creates a checker that checks the presence of TCP options in
// SYN segments.
//
// If wndscale is negative, the window scale option must not be present.
func TCPSynOptions(wantOpts header.TCPSynOptions) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}
		opts := tcp.Options()
		limit := len(opts)
		foundMSS := false
		foundWS := false
		foundTS := false
		tsVal := uint32(0)
		tsEcr := uint32(0)
		for i := 0; i < limit; {
			switch opts[i] {
			case header.TCPOptionEOL:
				i = limit
			case header.TCPOptionNOP:
				i++
			case header.TCPOptionMSS:
				v := uint16(opts[i+2])<<8 | uint16(opts[i+3])
				if wantOpts.MSS != v {
					t.Fatalf("Bad MSS: got %v, want %v", v, wantOpts.MSS)
				}
				foundMSS = true
				i += 4
			case header.TCPOptionWS:
				if wantOpts.WS < 0 {
					t.Fatalf("WS present when it shouldn't be")
				}
				v := int(opts[i+2])
				if v != wantOpts.WS {
					t.Fatalf("Bad WS: got %v, want %v", v, wantOpts.WS)
				}
				foundWS = true
				i += 3
			case header.TCPOptionTS:
				if i+10 > limit || opts[i+1] != 10 {
					t.Fatalf("bad length %d for TS option, limit: %d", opts[i+1], limit)
				}
				tsVal = binary.BigEndian.Uint32(opts[i+2:])
				tsEcr = uint32(0)
				if tcp.Flags()&header.TCPFlagAck != 0 {
					// If the syn is an SYN-ACK then read
					// the tsEcr value as well.
					tsEcr = binary.BigEndian.Uint32(opts[i+6:])
				}
				foundTS = true
				i += 10
			default:
				i += int(opts[i+1])
			}
		}

		if !foundMSS {
			t.Fatalf("MSS option not found. Options: %x", opts)
		}

		if !foundWS && wantOpts.WS >= 0 {
			t.Fatalf("WS option not found. Options: %x", opts)
		}
		if wantOpts.TS && !foundTS {
			t.Fatalf("TS option not found. Options: %x", opts)
		}
		if foundTS && tsVal == 0 {
			t.Fatalf("TS option specified but the timestamp value is zero")
		}
		if foundTS && tsEcr == 0 && wantOpts.TSEcr != 0 {
			t.Fatalf("TS option specified but TSEcr is incorrect: got %d, want: %d", tsEcr, wantOpts.TSEcr)
		}
	}
}

// TCPTimestampChecker creates a checker that validates that a TCP segment has a
// TCP Timestamp option if wantTS is true, it also compares the wantTSVal and
// wantTSEcr values with those in the TCP segment (if present).
//
// If wantTSVal or wantTSEcr is zero then the corresponding comparison is
// skipped.
func TCPTimestampChecker(wantTS bool, wantTSVal uint32, wantTSEcr uint32) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		tcp, ok := h.(header.TCP)
		if !ok {
			return
		}
		opts := []byte(tcp.Options())
		limit := len(opts)
		foundTS := false
		tsVal := uint32(0)
		tsEcr := uint32(0)
		for i := 0; i < limit; {
			switch opts[i] {
			case header.TCPOptionEOL:
				i = limit
			case header.TCPOptionNOP:
				i++
			case header.TCPOptionTS:
				if i+10 > limit {
					t.Fatalf("TS option found, but option is truncated, option length: %d, want 10 bytes", limit-i)
				}
				if opts[i+1] != 10 {
					t.Fatalf("TS option found, but bad length specified: %d, want: 10", opts[i+1])
				}
				tsVal = binary.BigEndian.Uint32(opts[i+2:])
				tsEcr = binary.BigEndian.Uint32(opts[i+6:])
				foundTS = true
				i += 10
			default:
				// We don't recognize this option, just skip over it.
				if i+2 > limit {
					return
				}
				l := int(opts[i+1])
				if i < 2 || i+l > limit {
					return
				}
				i += l
			}
		}

		if wantTS != foundTS {
			t.Fatalf("TS Option mismatch: got TS= %v, want TS= %v", foundTS, wantTS)
		}
		if wantTS && wantTSVal != 0 && wantTSVal != tsVal {
			t.Fatalf("Timestamp value is incorrect: got: %d, want: %d", tsVal, wantTSVal)
		}
		if wantTS && wantTSEcr != 0 && tsEcr != wantTSEcr {
			t.Fatalf("Timestamp Echo Reply is incorrect: got: %d, want: %d", tsEcr, wantTSEcr)
		}
	}
}

// Payload creates a checker that checks the payload.
func Payload(want []byte) TransportChecker {
	return func(t *testing.T, h header.Transport) {
		if got := h.Payload(); !reflect.DeepEqual(got, want) {
			t.Fatalf("Wrong payload, got %v, want %v", got, want)
		}
	}
}
