// Copyright 2026 The HuaTuo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/bpf"
)

// primitive implements Filter and Element
type primitive struct {
	kind        filterKind
	direction   filterDirection
	protocol    filterProtocol
	subProtocol filterSubProtocol
	negator     bool
	id          string
	resolved    []net.IP
}

func (p primitive) IsPrimitive() bool { return true }
func (p primitive) Type() ElementType { return Primitive }
func (p primitive) Distill() Filter   { return p }

func (p primitive) Combine(o *primitive) *primitive {
	if p.kind == filterKindVLAN || p.kind == filterKindMPLS ||
		o.kind == filterKindVLAN || o.kind == filterKindMPLS {
		return nil
	}
	if p.kind == filterKindMulticast || o.kind == filterKindMulticast {
		return nil
	}
	if !p.kindCanCombineSubProtocol() && o.subProtocol != filterSubProtocolUnset {
		return nil
	}
	if !o.kindCanCombineSubProtocol() && p.subProtocol != filterSubProtocolUnset {
		return nil
	}
	if p.Equal(o) {
		return &p
	}
	c := primitive{}
	switch {
	case p.kind == o.kind || o.kind == filterKindUnset:
		c.kind = p.kind
	case p.kind == filterKindUnset:
		c.kind = o.kind
	default:
		return nil
	}
	switch {
	case p.direction == o.direction || o.direction == filterDirectionUnset:
		c.direction = p.direction
	case p.direction == filterDirectionUnset:
		c.direction = o.direction
	default:
		return nil
	}
	switch {
	case p.protocol == o.protocol || o.protocol == filterProtocolUnset:
		c.protocol = p.protocol
	case p.protocol == filterProtocolUnset:
		c.protocol = o.protocol
	default:
		return nil
	}
	switch {
	case p.subProtocol == o.subProtocol || o.subProtocol == filterSubProtocolUnset:
		c.subProtocol = p.subProtocol
	case p.subProtocol == filterSubProtocolUnset:
		c.subProtocol = o.subProtocol
	default:
		return nil
	}
	switch {
	case p.id == o.id || o.id == "":
		c.id = p.id
	case p.id == "":
		c.id = o.id
	default:
		return nil
	}
	switch p.negator {
	case o.negator:
		c.negator = p.negator
	default:
		return nil
	}
	return &c
}

func (p primitive) kindCanCombineSubProtocol() bool {
	switch p.kind {
	case filterKindUnset, filterKindPort, filterKindPortRange:
		return true
	default:
		return false
	}
}

// Compile validates the primitive and delegates to the two-pass assembler.
func (p primitive) Compile(layout linkLayout) ([]bpf.Instruction, error) {
	return compileFilter(p, layout)
}

// Size compiles and returns the instruction count. Guarantees len(Compile) == Size.
func (p primitive) Size(layout linkLayout) uint8 {
	if p.isAlwaysReject(layout) {
		return 2
	}
	if p.isAlwaysAccept(layout) {
		return 2
	}
	insns, err := p.Compile(layout)
	if err != nil || len(insns) == 0 {
		return 2
	}
	return uint8(len(insns))
}

// isAlwaysReject returns true for L2-only primitives on non-L2 link types.
func (p primitive) isAlwaysReject(layout linkLayout) bool {
	if layout.hasL2Protocols() {
		return false
	}
	if p.kind == filterKindVLAN || p.kind == filterKindMPLS {
		return true
	}
	switch p.kind {
	case filterKindHost:
		switch p.protocol {
		case filterProtocolEther, filterProtocolArp, filterProtocolRarp:
			return true
		}
	case filterKindNet:
		switch p.protocol {
		case filterProtocolArp, filterProtocolRarp:
			return true
		}
	case filterKindMulticast:
		switch p.protocol {
		case filterProtocolUnset, filterProtocolEther:
			return true
		}
	case filterKindUnset:
		if p.subProtocol == filterSubProtocolSTP {
			return true
		}
		switch p.protocol {
		case filterProtocolArp, filterProtocolRarp:
			return true
		}
	}
	return false
}

// isAlwaysAccept always returns false.
func (p primitive) isAlwaysAccept(layout linkLayout) bool { return false }

// emit generates instructions for this primitive, branching to onMatch when
// the filter matches and onMiss when it does not. Negation is handled
// externally by the negated wrapper which swaps onMatch/onMiss before calling emit.
func (p primitive) emit(b *prog, onMatch, onMiss labelID) {
	switch p.kind {
	case filterKindHost:
		p.emitHost(b, onMatch, onMiss)
	case filterKindPort:
		p.emitPort(b, onMatch, onMiss)
	case filterKindPortRange:
		p.emitPortRange(b, onMatch, onMiss)
	case filterKindNet:
		p.emitNet(b, onMatch, onMiss)
	case filterKindMulticast:
		p.emitMulticast(b, onMatch, onMiss)
	case filterKindVLAN:
		p.emitVLAN(b, onMatch, onMiss)
	case filterKindMPLS:
		p.emitMPLS(b, onMatch, onMiss)
	case filterKindUnset:
		p.emitUnset(b, onMatch, onMiss)
	}
}

func (p primitive) emitVLAN(b *prog, onMatch, onMiss labelID) {
	if !b.layout.hasL2Protocols() {
		b.emitJump(onMiss)
		return
	}
	b.loadEtherKind()
	matchedType := b.newLabel()
	tryQinQ := b.newLabel()
	tryVLAN9100 := b.newLabel()
	b.emitJumpIf(bpf.JumpEqual, etherTypeVLAN, matchedType, tryQinQ)
	b.bind(tryQinQ)
	b.emitJumpIf(bpf.JumpEqual, etherTypeQinQ, matchedType, tryVLAN9100)
	b.bind(tryVLAN9100)
	b.emitJumpIf(bpf.JumpEqual, etherTypeVLAN9100, matchedType, onMiss)
	b.bind(matchedType)
	if p.id == "" {
		b.emitJump(onMatch)
		return
	}
	vlanID, err := strconv.ParseUint(p.id, 10, 12)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	b.emit(bpf.LoadAbsolute{Off: b.layout.l3Off(), Size: lengthHalf})
	b.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0x0fff})
	b.emitJumpIf(bpf.JumpEqual, uint32(vlanID), onMatch, onMiss)
}

func (p primitive) emitMPLS(b *prog, onMatch, onMiss labelID) {
	layout, ok := b.layout.(cursorLayout)
	if !ok || layout.cursor.probe != probeEtherType {
		if ok && layout.cursor.mplsDepth > 0 {
			p.emitNextMPLSLabel(b, layout.cursor, onMatch, onMiss)
			return
		}
		b.emitJump(onMiss)
		return
	}
	b.loadEtherKind()
	matchedType := b.newLabel()
	tryMulticast := b.newLabel()
	b.emitJumpIf(bpf.JumpEqual, etherTypeMPLSUnicast, matchedType, tryMulticast)
	b.bind(tryMulticast)
	b.emitJumpIf(bpf.JumpEqual, etherTypeMPLSMulticast, matchedType, onMiss)
	b.bind(matchedType)
	if p.id == "" {
		b.emitJump(onMatch)
		return
	}
	p.emitMPLSLabelValue(b, layout.cursor.l3Offset, onMatch, onMiss)
}

func (p primitive) emitNextMPLSLabel(
	b *prog,
	cursor packetCursor,
	onMatch labelID,
	onMiss labelID,
) {
	nextLabel := b.newLabel()
	b.emit(bpf.LoadAbsolute{Off: cursor.l3Offset - 2, Size: lengthByte})
	b.emitJumpIf(bpf.JumpBitsSet, 0x01, onMiss, nextLabel)
	b.bind(nextLabel)
	if p.id == "" {
		b.emitJump(onMatch)
		return
	}
	p.emitMPLSLabelValue(b, cursor.l3Offset, onMatch, onMiss)
}

func (p primitive) emitMPLSLabelValue(b *prog, offset uint32, onMatch, onMiss labelID) {
	label, err := strconv.ParseUint(p.id, 10, 20)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	b.emit(bpf.LoadAbsolute{Off: offset, Size: lengthWord})
	b.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0xfffff000})
	b.emitJumpIf(bpf.JumpEqual, uint32(label)<<12, onMatch, onMiss)
}

func (p primitive) emitHost(b *prog, onMatch, onMiss labelID) {
	if p.subProtocol != filterSubProtocolUnset {
		hostOK := b.newLabel()
		p.emitUnsetProtocol(b, hostOK, onMiss)
		b.bind(hostOK)
		p.subProtocol = filterSubProtocolUnset
		p.emitHost(b, onMatch, onMiss)
		return
	}

	switch p.protocol {
	case filterProtocolEther:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		// No loadEtherKind: ether address check is independent of ethertype.
		b.checkEtherAddresses(p.direction, p.id, onMatch, onMiss)

	case filterProtocolIP6:
		b.loadEtherKind()
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		_, a6, err := p.getAddrs(net.DefaultResolver)
		if err != nil {
			b.emitJump(onMiss)
			return
		}
		b.checkIP6HostAddressList(p.direction, a6, onMatch, onMiss)

	case filterProtocolIP:
		b.loadEtherKind()
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		a4, _, err := p.getAddrs(net.DefaultResolver)
		if err != nil {
			b.emitJump(onMiss)
			return
		}
		b.checkIP4HostAddressList(p.direction, a4, onMatch, onMiss)

	case filterProtocolArp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.loadEtherKind()
		arpOk := b.newLabel()
		b.compareProtocolArp(arpOk, onMiss)
		b.bind(arpOk)
		a4, _, err := p.getAddrs(net.DefaultResolver)
		if err != nil {
			b.emitJump(onMiss)
			return
		}
		b.checkIP4ArpAddressList(p.direction, a4, onMatch, onMiss)

	case filterProtocolRarp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.loadEtherKind()
		rarpOk := b.newLabel()
		b.compareProtocolRarp(rarpOk, onMiss)
		b.bind(rarpOk)
		a4, _, err := p.getAddrs(net.DefaultResolver)
		if err != nil {
			b.emitJump(onMiss)
			return
		}
		b.checkIP4ArpAddressList(p.direction, a4, onMatch, onMiss)

	case filterProtocolUnset:
		p.emitHostUnset(b, onMatch, onMiss)
	}
}

func (p primitive) emitHostUnset(b *prog, onMatch, onMiss labelID) {
	a4, a6, err := p.getAddrs(net.DefaultResolver)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	b.loadEtherKind()
	hasL2 := b.layout.hasL2Protocols()

	if len(a4) > 0 {
		ip4Ok := b.newLabel()
		// Determine where the IPv4 protocol-miss branch goes.
		// With L2: always to afterIP4 (to check ARP/RARP).
		// Without L2: directly to onMiss (no ARP/RARP on non-L2 links).
		if hasL2 {
			afterIP4 := b.newLabel()
			b.compareProtocolIP4(ip4Ok, afterIP4)
			b.bind(ip4Ok)
			b.checkIP4HostAddressList(p.direction, a4, onMatch, onMiss)
			b.bind(afterIP4)
			// ARP and RARP share the same address-check label.
			arpRarpCheck := b.newLabel()
			notARP := b.newLabel()
			b.compareProtocolArp(arpRarpCheck, notARP)
			b.bind(notARP)
			rarpMiss := onMiss
			if len(a6) > 0 {
				rarpMiss = b.newLabel()
			}
			b.compareProtocolRarp(arpRarpCheck, rarpMiss)
			b.bind(arpRarpCheck)
			b.checkIP4ArpAddressList(p.direction, a4, onMatch, onMiss)
			if len(a6) > 0 {
				b.bind(rarpMiss)
			}
		} else {
			ip4Miss := onMiss
			if len(a6) > 0 {
				ip4Miss = b.newLabel()
			}
			b.compareProtocolIP4(ip4Ok, ip4Miss)
			b.bind(ip4Ok)
			b.checkIP4HostAddressList(p.direction, a4, onMatch, onMiss)
			if ip4Miss != onMiss {
				b.bind(ip4Miss)
			}
		}
	}

	if len(a6) > 0 {
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		b.checkIP6HostAddressList(p.direction, a6, onMatch, onMiss)
	}
}

func (p primitive) emitPort(b *prog, onMatch, onMiss labelID) {
	portInt, err := findPort(p.id)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	port := uint32(portInt)
	b.loadEtherKind()

	switch p.protocol {
	case filterProtocolIP6:
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		b.emit(loadIPv6Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, true)
		b.checkPorts(p.direction, port, onMatch, onMiss, true)

	case filterProtocolIP:
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.emit(loadIPv4Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, false)
		b.checkPorts(p.direction, port, onMatch, onMiss, false)

	case filterProtocolUnset:
		// Dual-stack: try IPv6 first, then IPv4.
		tryIP4 := b.newLabel()
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, tryIP4)
		b.bind(ip6Ok)
		b.emit(loadIPv6Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, true)
		b.checkPorts(p.direction, port, onMatch, onMiss, true)

		b.bind(tryIP4)
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.emit(loadIPv4Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, false)
		b.checkPorts(p.direction, port, onMatch, onMiss, false)
	}
}

func (p primitive) emitPortRange(b *prog, onMatch, onMiss labelID) {
	minimum, maximum, err := findPortRange(p.id)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	b.loadEtherKind()

	switch p.protocol {
	case filterProtocolIP6:
		ip6OK := b.newLabel()
		b.compareProtocolIP6(ip6OK, onMiss)
		b.bind(ip6OK)
		b.emit(loadIPv6Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, true)
		b.checkPortRange(p.direction, minimum, maximum, onMatch, onMiss, true)
	case filterProtocolIP:
		ip4OK := b.newLabel()
		b.compareProtocolIP4(ip4OK, onMiss)
		b.bind(ip4OK)
		b.emit(loadIPv4Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, false)
		b.checkPortRange(p.direction, minimum, maximum, onMatch, onMiss, false)
	case filterProtocolUnset:
		tryIP4 := b.newLabel()
		ip6OK := b.newLabel()
		b.compareProtocolIP6(ip6OK, tryIP4)
		b.bind(ip6OK)
		b.emit(loadIPv6Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, true)
		b.checkPortRange(p.direction, minimum, maximum, onMatch, onMiss, true)

		b.bind(tryIP4)
		ip4OK := b.newLabel()
		b.compareProtocolIP4(ip4OK, onMiss)
		b.bind(ip4OK)
		b.emit(loadIPv4Protocol(b.layout))
		p.emitSubProtocolCompare(b, onMatch, onMiss, false)
		b.checkPortRange(p.direction, minimum, maximum, onMatch, onMiss, false)
	}
}

func (p primitive) emitNet(b *prog, onMatch, onMiss labelID) {
	if p.subProtocol != filterSubProtocolUnset {
		netOK := b.newLabel()
		p.emitUnsetProtocol(b, netOK, onMiss)
		b.bind(netOK)
		p.subProtocol = filterSubProtocolUnset
		p.emitNet(b, onMatch, onMiss)
		return
	}

	switch p.protocol {
	case filterProtocolIP6:
		b.loadEtherKind()
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		addr, network, err := getNetAndMask(p.id)
		if err != nil {
			b.emitJump(onMiss)
			return
		}
		b.checkIP6NetAddresses(p.direction, addr, network.Mask, onMatch, onMiss)

	case filterProtocolIP:
		b.loadEtherKind()
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.checkIP4NetHostAddresses(p.direction, p.id, onMatch, onMiss)

	case filterProtocolArp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.loadEtherKind()
		arpOk := b.newLabel()
		b.compareProtocolArp(arpOk, onMiss)
		b.bind(arpOk)
		b.checkIP4NetArpAddresses(p.direction, p.id, onMatch, onMiss)

	case filterProtocolRarp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.loadEtherKind()
		rarpOk := b.newLabel()
		b.compareProtocolRarp(rarpOk, onMiss)
		b.bind(rarpOk)
		b.checkIP4NetArpAddresses(p.direction, p.id, onMatch, onMiss)

	case filterProtocolUnset:
		p.emitNetUnset(b, onMatch, onMiss)
	}
}

func (p primitive) emitNetUnset(b *prog, onMatch, onMiss labelID) {
	b.loadEtherKind()
	addr, network, err := getNetAndMask(p.id)
	if err != nil {
		b.emitJump(onMiss)
		return
	}
	hasL2 := b.layout.hasL2Protocols()

	if addr.To4() != nil {
		ip4Ok := b.newLabel()
		var arpStart labelID

		if hasL2 {
			arpStart = b.newLabel()
			b.compareProtocolIP4(ip4Ok, arpStart)
		} else {
			b.compareProtocolIP4(ip4Ok, onMiss)
		}
		b.bind(ip4Ok)
		b.checkIP4NetHostAddresses(p.direction, p.id, onMatch, onMiss)

		if hasL2 {
			b.bind(arpStart)
			arpOk := b.newLabel()
			tryRARP := b.newLabel()
			b.compareProtocolArp(arpOk, tryRARP)
			b.bind(arpOk)

			arpAddrDo := b.newLabel()
			b.emitJump(arpAddrDo)

			b.bind(tryRARP)
			rarpOk := b.newLabel()
			b.compareProtocolRarp(rarpOk, onMiss)
			b.bind(rarpOk)

			b.bind(arpAddrDo)
			b.checkIP4NetArpAddresses(p.direction, p.id, onMatch, onMiss)
		}
	} else {
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		b.checkIP6NetAddresses(p.direction, addr, network.Mask, onMatch, onMiss)
	}
}

func (p primitive) emitMulticast(b *prog, onMatch, onMiss labelID) {
	switch p.protocol {
	case filterProtocolUnset, filterProtocolEther:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.genEtherMulticast(onMatch, onMiss)

	case filterProtocolIP:
		b.loadEtherKind()
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.genIPv4Multicast(onMatch, onMiss)

	case filterProtocolIP6:
		b.loadEtherKind()
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		b.genIPv6Multicast(onMatch, onMiss)
	}
}

func (p primitive) emitUnset(b *prog, onMatch, onMiss labelID) {
	b.loadEtherKind()

	switch p.protocol {
	case filterProtocolIP:
		// Bare `ip` matches any IPv4 packet by EtherType alone, like tcpdump.
		// Only narrow by protocol number when a sub-protocol was given,
		// e.g. `ip proto tcp`.
		if p.subProtocol == filterSubProtocolUnset {
			b.compareProtocolIP4(onMatch, onMiss)
			return
		}
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.compareIPv4Protocol(p.protoNum(), onMatch, onMiss)

	case filterProtocolIP6:
		// Bare `ip6` matches any IPv6 packet by EtherType alone, like tcpdump.
		// Only narrow by next-header when a sub-protocol was given,
		// e.g. `ip6 proto udp`.
		if p.subProtocol == filterSubProtocolUnset {
			b.compareProtocolIP6(onMatch, onMiss)
			return
		}
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		p.emitIPv6SubProtocol(b, onMatch, onMiss)

	case filterProtocolArp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		arpOk := b.newLabel()
		b.compareProtocolArp(arpOk, onMiss)
		b.bind(arpOk)
		b.emitJump(onMatch)

	case filterProtocolRarp:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		rarpOk := b.newLabel()
		b.compareProtocolRarp(rarpOk, onMiss)
		b.bind(rarpOk)
		b.emitJump(onMatch)

	case filterProtocolEther:
		switch p.subProtocol {
		case filterSubProtocolIP:
			b.compareProtocolIP4(onMatch, onMiss)
		case filterSubProtocolIP6:
			b.compareProtocolIP6(onMatch, onMiss)
		case filterSubProtocolArp:
			b.compareProtocolArp(onMatch, onMiss)
		case filterSubProtocolRarp:
			b.compareProtocolRarp(onMatch, onMiss)
		case filterSubProtocolSTP:
			b.compareProtocolSTP(onMatch, onMiss)
		}

	case filterProtocolUnset:
		p.emitUnsetProtocol(b, onMatch, onMiss)
	}
}

// emitIPv6SubProtocol checks for a sub-protocol in an IPv6 packet.
func (p primitive) emitIPv6SubProtocol(b *prog, onMatch, onMiss labelID) {
	switch p.subProtocol {
	case filterSubProtocolSCTP:
		subOk := b.newLabel()
		b.compareSubProtocolSctp(subOk, onMiss)
		b.bind(subOk)
		b.emitJump(onMatch)
	default:
		b.compareIPv6Protocol(p.protoNum(), onMatch, onMiss)
	}
}

// emitUnsetProtocol emits the dual-stack sub-protocol checks for
// filterProtocolUnset with filterKindUnset.
func (p primitive) emitUnsetProtocol(b *prog, onMatch, onMiss labelID) {
	switch p.subProtocol {
	// IPv4-only sub-protocols.
	case filterSubProtocolIcmp, filterSubProtocolIgmp:
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.compareIPv4Protocol(p.protoNum(), onMatch, onMiss)

	// IPv6-only sub-protocols.
	case filterSubProtocolIcmp6:
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, onMiss)
		b.bind(ip6Ok)
		b.compareIPv6Protocol(p.protoNum(), onMatch, onMiss)

	// Dual-stack sub-protocols: try IPv6 first, then IPv4.
	case filterSubProtocolUDP, filterSubProtocolTCP, filterSubProtocolSCTP,
		filterSubProtocolPim, filterSubProtocolEsp, filterSubProtocolAh,
		filterSubProtocolVrrp:
		tryIP4 := b.newLabel()
		ip6Ok := b.newLabel()
		b.compareProtocolIP6(ip6Ok, tryIP4)
		b.bind(ip6Ok)
		p.emitIPv6SubProtocol(b, onMatch, onMiss)

		b.bind(tryIP4)
		ip4Ok := b.newLabel()
		b.compareProtocolIP4(ip4Ok, onMiss)
		b.bind(ip4Ok)
		b.compareIPv4Protocol(p.protoNum(), onMatch, onMiss)

	case filterSubProtocolSTP:
		if !b.layout.hasL2Protocols() {
			b.emitJump(onMiss)
			return
		}
		b.compareProtocolSTP(onMatch, onMiss)
	}
}

// emitSubProtocolCompare emits the sub-protocol comparison for Port primitives.
// On match it falls through to the port check that follows. On miss it branches
// to onMiss. This is shared by the IPv4 and IPv6 port paths.
func (p primitive) emitSubProtocolCompare(b *prog, onMatch, onMiss labelID, ip6 bool) {
	switch p.subProtocol {
	case filterSubProtocolTCP:
		subOk := b.newLabel()
		b.compareSubProtocolTCP(subOk, onMiss)
		b.bind(subOk)

	case filterSubProtocolUDP:
		subOk := b.newLabel()
		b.compareSubProtocolUDP(subOk, onMiss)
		b.bind(subOk)

	case filterSubProtocolSCTP:
		subOk := b.newLabel()
		b.compareSubProtocolSctp(subOk, onMiss)
		b.bind(subOk)

	case filterSubProtocolUnset:
		// Check SCTP then TCP then UDP. Fall through on any match.
		afterAll := b.newLabel()

		sctpOk := b.newLabel()
		tryTCP := b.newLabel()
		b.compareSubProtocolSctp(sctpOk, tryTCP)
		b.bind(sctpOk)
		b.emitJump(afterAll)

		b.bind(tryTCP)
		tcpOk := b.newLabel()
		tryUDP := b.newLabel()
		b.compareSubProtocolTCP(tcpOk, tryUDP)
		b.bind(tcpOk)
		b.emitJump(afterAll)

		b.bind(tryUDP)
		udpOk := b.newLabel()
		b.compareSubProtocolUDP(udpOk, onMiss)
		b.bind(udpOk)

		b.bind(afterAll)
	}
}

// protoNum maps the primitive's subProtocol to its IP protocol number.
func (p primitive) protoNum() uint32 {
	switch p.subProtocol {
	case filterSubProtocolTCP:
		return ipProtocolTCP
	case filterSubProtocolUDP:
		return ipProtocolUDP
	case filterSubProtocolSCTP:
		return ipProtocolSctp
	case filterSubProtocolIcmp:
		return ipProtocolIcmp
	case filterSubProtocolIgmp:
		return ipProtocolIgmp
	case filterSubProtocolIcmp6:
		return ipProtocolIcmp6
	case filterSubProtocolPim:
		return ipProtocolPim
	case filterSubProtocolEsp:
		return ipProtocolEsp
	case filterSubProtocolAh:
		return ipProtocolAh
	case filterSubProtocolVrrp:
		return ipProtocolVrrp
	}
	return 0
}

func (p primitive) Equal(f Filter) bool {
	if f == nil {
		return false
	}
	o, ok := f.(primitive)
	if !ok {
		return false
	}
	return p.kind == o.kind &&
		p.direction == o.direction &&
		p.protocol == o.protocol &&
		p.subProtocol == o.subProtocol &&
		p.negator == o.negator &&
		p.id == o.id
}

func (p primitive) validate() error {
	switch {
	case p.subProtocol == filterSubProtocolUnknown:
		if _, err := strconv.ParseUint(p.id, 0, 8); err == nil {
			return unsupportedFeature("numeric protocol")
		}
		return fmt.Errorf("unknown protocol %s", p.id)
	case p.kind == filterKindGateway:
		return unsupportedFeature("gateway")
	case p.kind == filterKindVLAN:
		if p.id == "" {
			return nil
		}
		vlanID, err := strconv.ParseUint(p.id, 10, 12)
		if err != nil || vlanID > 4095 {
			return fmt.Errorf("invalid vlan id: %s", p.id)
		}
		return nil
	case p.kind == filterKindMPLS:
		if p.id == "" {
			return nil
		}
		if _, err := strconv.ParseUint(p.id, 10, 20); err != nil {
			return fmt.Errorf("invalid mpls label: %s", p.id)
		}
		return nil
	case p.protocol == filterProtocolFddi:
		return unsupportedFeature("fddi")
	case p.protocol == filterProtocolTr:
		return unsupportedFeature("tr")
	case p.protocol == filterProtocolWlan:
		return unsupportedFeature("wlan")
	case p.protocol == filterProtocolDecnet:
		return unsupportedFeature("decnet")
	case isUnsupportedSubProtocol(p.subProtocol):
		return unsupportedFeature(subProtocolName(p.subProtocol))
	case !isSupportedCombination(p):
		return unsupportedFeature("qualifier combination")
	case p.kind == filterKindUnset && p.id != "":
		return fmt.Errorf("parse error")
	case p.kind == filterKindUnset && p.direction != filterDirectionSrcOrDst:
		return fmt.Errorf("direction qualifier requires host, net, port, or portrange")
	case p.kind == filterKindHost:
		switch p.protocol {
		case filterProtocolIP, filterProtocolIP6, filterProtocolArp, filterProtocolRarp, filterProtocolUnset:
			if p.id == "" {
				return fmt.Errorf("blank host")
			}
			addr, network, parseErr := getNetAndMask(p.id)
			var maskFull net.IPMask
			if parseErr == nil && addr != nil && network != nil {
				if addr.To4() != nil {
					maskFull = ip4MaskFull
				} else {
					maskFull = ip6MaskFull
				}
				if !bytes.Equal(network.Mask, maskFull) {
					return fmt.Errorf("invalid host address with CIDR: %s", p.id)
				}
			}
		case filterProtocolEther:
			hardwareAddr, err := net.ParseMAC(p.id)
			if err != nil || len(hardwareAddr) != 6 {
				return fmt.Errorf("invalid ethernet address: %s", p.id)
			}
		}
	case p.kind == filterKindMulticast:
		switch p.protocol {
		case filterProtocolUnset, filterProtocolEther, filterProtocolIP, filterProtocolIP6:
		default:
			return fmt.Errorf("multicast not supported for protocol")
		}
	case p.kind == filterKindUnset && p.protocol == filterProtocolUnset && p.subProtocol == filterSubProtocolUnset:
		return fmt.Errorf("parse error")
	case p.kind == filterKindPort:
		if _, err := findPort(p.id); err != nil {
			return err
		}
		if !isPortSubProtocol(p.subProtocol) {
			return unsupportedFeature("port protocol " + subProtocolName(p.subProtocol))
		}
	case p.kind == filterKindPortRange:
		if _, _, err := findPortRange(p.id); err != nil {
			return err
		}
		if !isPortSubProtocol(p.subProtocol) {
			return unsupportedFeature("portrange protocol " + subProtocolName(p.subProtocol))
		}
	case p.kind == filterKindNet:
		addr, network, err := getNetAndMask(p.id)
		if err != nil {
			return err
		}
		isIPv4 := addr.To4() != nil
		if p.protocol == filterProtocolIP6 && isIPv4 {
			return fmt.Errorf("ipv4 network used with ip6 qualifier: %s", p.id)
		}
		needsIPv4 := p.protocol == filterProtocolIP || p.protocol == filterProtocolArp ||
			p.protocol == filterProtocolRarp
		if needsIPv4 && !isIPv4 {
			return fmt.Errorf("ipv6 network used with ipv4 qualifier: %s", p.id)
		}
		masked := addr.Mask(network.Mask)
		if !addr.Equal(masked) {
			return fmt.Errorf("invalid network, network bits extend past mask bits: %s", p.id)
		}
	case p.kind == filterKindUnset && p.protocol == filterProtocolEther && p.subProtocol == filterSubProtocolUnset:
		return fmt.Errorf("parse error")
	}
	return nil
}

func isSupportedCombination(p primitive) bool {
	switch p.kind {
	case filterKindHost:
		if p.protocol != filterProtocolUnset && p.protocol != filterProtocolEther &&
			p.protocol != filterProtocolIP && p.protocol != filterProtocolIP6 &&
			p.protocol != filterProtocolArp && p.protocol != filterProtocolRarp {
			return false
		}
		return p.subProtocol != filterSubProtocolSTP && (p.subProtocol == filterSubProtocolUnset ||
			isSupportedProtocolPrimitive(p.protocol, p.subProtocol))
	case filterKindNet:
		if p.protocol != filterProtocolUnset && p.protocol != filterProtocolIP &&
			p.protocol != filterProtocolIP6 && p.protocol != filterProtocolArp &&
			p.protocol != filterProtocolRarp {
			return false
		}
		return p.subProtocol != filterSubProtocolSTP && (p.subProtocol == filterSubProtocolUnset ||
			isSupportedProtocolPrimitive(p.protocol, p.subProtocol))
	case filterKindPort, filterKindPortRange:
		return p.protocol == filterProtocolUnset || p.protocol == filterProtocolIP ||
			p.protocol == filterProtocolIP6
	case filterKindMulticast:
		return p.subProtocol == filterSubProtocolUnset
	case filterKindUnset:
		return isSupportedProtocolPrimitive(p.protocol, p.subProtocol)
	default:
		return true
	}
}

func isSupportedProtocolPrimitive(protocol filterProtocol, subProtocol filterSubProtocol) bool {
	if subProtocol == filterSubProtocolUnset {
		return true
	}
	switch protocol {
	case filterProtocolUnset:
		switch subProtocol {
		case filterSubProtocolTCP, filterSubProtocolUDP, filterSubProtocolSCTP,
			filterSubProtocolIcmp, filterSubProtocolIcmp6, filterSubProtocolIgmp,
			filterSubProtocolPim, filterSubProtocolEsp, filterSubProtocolAh,
			filterSubProtocolVrrp, filterSubProtocolSTP:
			return true
		}
	case filterProtocolIP:
		switch subProtocol {
		case filterSubProtocolTCP, filterSubProtocolUDP, filterSubProtocolSCTP,
			filterSubProtocolIcmp, filterSubProtocolIgmp, filterSubProtocolPim,
			filterSubProtocolEsp, filterSubProtocolAh, filterSubProtocolVrrp:
			return true
		}
	case filterProtocolIP6:
		switch subProtocol {
		case filterSubProtocolTCP, filterSubProtocolUDP, filterSubProtocolSCTP,
			filterSubProtocolIcmp6, filterSubProtocolPim, filterSubProtocolEsp,
			filterSubProtocolAh, filterSubProtocolVrrp:
			return true
		}
	case filterProtocolEther:
		switch subProtocol {
		case filterSubProtocolIP, filterSubProtocolIP6, filterSubProtocolArp,
			filterSubProtocolRarp, filterSubProtocolSTP:
			return true
		}
	}
	return false
}

func (p primitive) getAddrs(resolver Resolver) ([]net.IP, []net.IP, error) {
	addrs := p.resolved
	if len(addrs) == 0 {
		var err error
		addrs, err = p.resolve(resolver)
		if err != nil {
			return nil, nil, err
		}
	}
	a4 := make([]net.IP, 0, len(addrs))
	a6 := make([]net.IP, 0, len(addrs))
	for _, a := range addrs {
		if a.To4() != nil {
			a4 = append(a4, a)
		} else {
			a6 = append(a6, a)
		}
	}
	return a4, a6, nil
}

func findPort(portStr string) (int, error) {
	if port, err := strconv.Atoi(portStr); err == nil {
		if port >= 0 && port <= 65535 {
			return port, nil
		}
		return -1, fmt.Errorf("invalid port: %s", portStr)
	}
	if port, err := net.LookupPort("tcp", portStr); err == nil {
		return port, nil
	}
	if port, err := net.LookupPort("udp", portStr); err == nil {
		return port, nil
	}
	return -1, fmt.Errorf("invalid port: %s", portStr)
}

func findPortRange(value string) (uint32, uint32, error) {
	first, last, ok := strings.Cut(value, "-")
	if !ok || first == "" || last == "" || strings.Contains(last, "-") {
		return 0, 0, fmt.Errorf("invalid portrange: %s", value)
	}
	minimum, err := findPort(first)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid portrange: %s", value)
	}
	maximum, err := findPort(last)
	if err != nil || minimum > maximum {
		return 0, 0, fmt.Errorf("invalid portrange: %s", value)
	}
	return uint32(minimum), uint32(maximum), nil
}

func isPortSubProtocol(protocol filterSubProtocol) bool {
	switch protocol {
	case filterSubProtocolUnset, filterSubProtocolTCP, filterSubProtocolUDP, filterSubProtocolSCTP:
		return true
	default:
		return false
	}
}

func isUnsupportedSubProtocol(protocol filterSubProtocol) bool {
	switch protocol {
	case filterSubProtocolAtalk, filterSubProtocolAarp, filterSubProtocolDecnet,
		filterSubProtocolSca, filterSubProtocolLat, filterSubProtocolMopdl,
		filterSubProtocolMoprc, filterSubProtocolIso, filterSubProtocolIPx,
		filterSubProtocolNetbeui, filterSubProtocolIgrp:
		return true
	default:
		return false
	}
}

func subProtocolName(protocol filterSubProtocol) string {
	for name, candidate := range subProtocols {
		if candidate == protocol {
			return name
		}
	}
	return "unknown"
}
