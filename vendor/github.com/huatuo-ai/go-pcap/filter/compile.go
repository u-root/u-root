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
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/net/bpf"
)

var (
	ip4MaskFull = net.CIDRMask(32, 32)
	ip6MaskFull = net.CIDRMask(128, 128)
	returnDrop  = bpf.RetConstant{Val: 0}
	returnKeep  = bpf.RetConstant{Val: 0x40000}
)

// =============================================================================
// Builder methods on *prog — emit helpers with label targets
// =============================================================================

func (b *prog) loadEtherKind() {
	b.emit(b.layout.genLinkProbe()...)
}

func (b *prog) compareProtocolIP4(onMatch, onMiss labelID) {
	j := b.layout.genLinkType(etherTypeIPv4, 0, 0).(bpf.JumpIf)
	b.emitJumpIf(j.Cond, j.Val, onMatch, onMiss)
}

func (b *prog) compareProtocolIP6(onMatch, onMiss labelID) {
	j := b.layout.genLinkType(etherTypeIPv6, 0, 0).(bpf.JumpIf)
	b.emitJumpIf(j.Cond, j.Val, onMatch, onMiss)
}

func (b *prog) compareProtocolArp(onMatch, onMiss labelID) {
	j := b.layout.genLinkType(etherTypeArp, 0, 0).(bpf.JumpIf)
	b.emitJumpIf(j.Cond, j.Val, onMatch, onMiss)
}

func (b *prog) compareProtocolRarp(onMatch, onMiss labelID) {
	j := b.layout.genLinkType(etherTypeRarp, 0, 0).(bpf.JumpIf)
	b.emitJumpIf(j.Cond, j.Val, onMatch, onMiss)
}

func (b *prog) compareSubProtocolTCP(onMatch, onMiss labelID) {
	b.emitJumpIf(bpf.JumpEqual, ipProtocolTCP, onMatch, onMiss)
}

func (b *prog) compareSubProtocolUDP(onMatch, onMiss labelID) {
	b.emitJumpIf(bpf.JumpEqual, ipProtocolUDP, onMatch, onMiss)
}

func (b *prog) compareSubProtocolSctp(onMatch, onMiss labelID) {
	b.emitJumpIf(bpf.JumpEqual, ipProtocolSctp, onMatch, onMiss)
}

func (b *prog) compareProtocolSTP(onMatch, onMiss labelID) {
	if !b.layout.hasL2Protocols() {
		b.emitJump(onMiss)
		return
	}
	lengthOK := b.newLabel()
	b.loadEtherKind()
	b.emitJumpIf(bpf.JumpGreaterThan, 1500, onMiss, lengthOK)
	b.bind(lengthOK)
	b.emit(bpf.LoadAbsolute{Off: b.layout.l3Off(), Size: lengthWord})
	b.emit(bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: 0xffffff00})
	b.emitJumpIf(bpf.JumpEqual, 0x42420300, onMatch, onMiss)
}

func (b *prog) compareIPv4Protocol(proto uint32, onMatch, onMiss labelID) {
	b.emit(loadIPv4Protocol(b.layout))
	b.emitJumpIf(bpf.JumpEqual, proto, onMatch, onMiss)
}

func (b *prog) compareIPv6Protocol(proto uint32, onMatch, onMiss labelID) {
	mid := b.newLabel()
	b.emit(loadIPv6Protocol(b.layout))
	b.emitJumpIf(bpf.JumpEqual, proto, onMatch, mid)
	b.bind(mid)
	restart := b.newLabel()
	b.emitJumpIf(bpf.JumpEqual, ip6ContinuationPacket, restart, onMiss)
	b.bind(restart)
	b.emit(loadIPv6ContinuationProtocol(b.layout))
	b.emitJumpIf(bpf.JumpEqual, proto, onMatch, onMiss)
}

func (b *prog) loadIPv4HeaderOffset(onMiss labelID) {
	cont := b.newLabel()
	b.emit(bpf.LoadAbsolute{Off: b.layout.l3Off() + intraIP4HeaderFlags, Size: lengthHalf})
	b.emitJumpIf(bpf.JumpBitsSet, jumpMask, onMiss, cont)
	b.bind(cont)
	b.emit(bpf.LoadMemShift{Off: b.layout.l3Off() + intraIP4HeaderSize})
}

func (b *prog) checkEtherAddresses(direction filterDirection, addr string, onMatch, onMiss labelID) {
	hwAddr, err := net.ParseMAC(addr)
	if err != nil || len(hwAddr) != 6 {
		b.emitJump(onMiss)
		return
	}
	lastFour := binary.BigEndian.Uint32(hwAddr[len(hwAddr)-4:])
	firstTwo := uint32(binary.BigEndian.Uint16(hwAddr[len(hwAddr)-6 : len(hwAddr)-4]))

	switch direction {
	case filterDirectionSrc:
		cont := b.newLabel()
		b.emit(loadEthernetSourceLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont, onMiss)
		b.bind(cont)
		b.emit(loadEthernetSourceFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, onMatch, onMiss)
	case filterDirectionDst:
		cont := b.newLabel()
		b.emit(loadEthernetDestinationLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont, onMiss)
		b.bind(cont)
		b.emit(loadEthernetDestinationFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, onMatch, onMiss)
	case filterDirectionSrcOrDst:
		tryDst := b.newLabel()
		cont1 := b.newLabel()
		b.emit(loadEthernetSourceLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont1, tryDst)
		b.bind(cont1)
		b.emit(loadEthernetSourceFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, onMatch, tryDst)
		b.bind(tryDst)
		cont2 := b.newLabel()
		b.emit(loadEthernetDestinationLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont2, onMiss)
		b.bind(cont2)
		b.emit(loadEthernetDestinationFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, onMatch, onMiss)
	case filterDirectionSrcAndDst:
		cont1 := b.newLabel()
		b.emit(loadEthernetSourceLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont1, onMiss)
		b.bind(cont1)
		cont2 := b.newLabel()
		b.emit(loadEthernetSourceFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, cont2, onMiss)
		b.bind(cont2)
		cont3 := b.newLabel()
		b.emit(loadEthernetDestinationLast)
		b.emitJumpIf(bpf.JumpEqual, lastFour, cont3, onMiss)
		b.bind(cont3)
		b.emit(loadEthernetDestinationFirst)
		b.emitJumpIf(bpf.JumpEqual, firstTwo, onMatch, onMiss)
	}
}

func (b *prog) checkIP4HostAddresses(direction filterDirection, addr net.IP, onMatch, onMiss labelID) {
	b.checkIP4Addresses(direction, addr, nil, onMatch, onMiss,
		loadIPv4SourceAddress(b.layout), loadIPv4DestinationAddress(b.layout))
}

func (b *prog) checkIP4HostAddressList(
	direction filterDirection,
	addrs []net.IP,
	onMatch labelID,
	onMiss labelID,
) {
	b.checkIP4AddressList(direction, addrs, onMatch, onMiss, b.checkIP4HostAddresses)
}

func (b *prog) checkIP4ArpAddressList(
	direction filterDirection,
	addrs []net.IP,
	onMatch labelID,
	onMiss labelID,
) {
	b.checkIP4AddressList(direction, addrs, onMatch, onMiss, b.checkIP4ArpAddresses)
}

func (b *prog) checkIP4AddressList(
	direction filterDirection,
	addrs []net.IP,
	onMatch labelID,
	onMiss labelID,
	check func(filterDirection, net.IP, labelID, labelID),
) {
	for i, addr := range addrs {
		miss := onMiss
		if i < len(addrs)-1 {
			miss = b.newLabel()
		}
		check(direction, addr, onMatch, miss)
		if miss != onMiss {
			b.bind(miss)
		}
	}
}

func (b *prog) checkIP4ArpAddresses(direction filterDirection, addr net.IP, onMatch, onMiss labelID) {
	b.checkIP4Addresses(direction, addr, nil, onMatch, onMiss,
		loadArpSenderAddress(b.layout), loadArpTargetAddress(b.layout))
}

func (b *prog) checkIP4NetHostAddresses(direction filterDirection, addr string, onMatch, onMiss labelID) {
	b.checkIP4NetAddresses(direction, addr, true, onMatch, onMiss)
}

func (b *prog) checkIP4NetArpAddresses(direction filterDirection, addr string, onMatch, onMiss labelID) {
	b.checkIP4NetAddresses(direction, addr, false, onMatch, onMiss)
}

func (b *prog) checkIP4NetAddresses(direction filterDirection, addr string, ip bool, onMatch, onMiss labelID) {
	addrBytes, network, err := getNetAndMask(addr)
	if err != nil || addrBytes == nil {
		b.emitJump(onMiss)
		return
	}
	var maskCheck *bpf.ALUOpConstant
	if !bytes.Equal(network.Mask, ip4MaskFull) {
		maskCheck = &bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: binary.BigEndian.Uint32(network.Mask)}
	}
	if ip {
		b.checkIP4Addresses(direction, addrBytes, maskCheck, onMatch, onMiss,
			loadIPv4SourceAddress(b.layout), loadIPv4DestinationAddress(b.layout))
	} else {
		b.checkIP4Addresses(direction, addrBytes, maskCheck, onMatch, onMiss,
			loadArpSenderAddress(b.layout), loadArpTargetAddress(b.layout))
	}
}

func (b *prog) checkIP4Addresses(
	direction filterDirection,
	addr []byte,
	maskCheck *bpf.ALUOpConstant,
	onMatch labelID,
	onMiss labelID,
	loadSrcDst ...bpf.Instruction,
) {
	if addr == nil {
		return
	}
	addrVal := binary.BigEndian.Uint32(addr[len(addr)-4:])

	switch direction {
	case filterDirectionSrc:
		b.emit(loadSrcDst[0])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, onMatch, onMiss)
	case filterDirectionDst:
		b.emit(loadSrcDst[1])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, onMatch, onMiss)
	case filterDirectionSrcOrDst:
		cont := b.newLabel()
		b.emit(loadSrcDst[0])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, onMatch, cont)
		b.bind(cont)
		b.emit(loadSrcDst[1])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, onMatch, onMiss)
	case filterDirectionSrcAndDst:
		cont := b.newLabel()
		b.emit(loadSrcDst[0])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, cont, onMiss)
		b.bind(cont)
		b.emit(loadSrcDst[1])
		if maskCheck != nil {
			b.emit(*maskCheck)
		}
		b.emitJumpIf(bpf.JumpEqual, addrVal, onMatch, onMiss)
	}
}

func (b *prog) checkIP6HostAddresses(direction filterDirection, addr net.IP, onMatch, onMiss labelID) {
	b.checkIP6Addresses(direction, addr, nil, onMatch, onMiss)
}

func (b *prog) checkIP6HostAddressList(
	direction filterDirection,
	addrs []net.IP,
	onMatch labelID,
	onMiss labelID,
) {
	for i, addr := range addrs {
		miss := onMiss
		if i < len(addrs)-1 {
			miss = b.newLabel()
		}
		b.checkIP6HostAddresses(direction, addr, onMatch, miss)
		if miss != onMiss {
			b.bind(miss)
		}
	}
}

func (b *prog) checkIP6NetAddresses(direction filterDirection, addr net.IP, mask net.IPMask, onMatch, onMiss labelID) {
	b.checkIP6Addresses(direction, addr, mask, onMatch, onMiss)
}

func (b *prog) checkIP6Addresses(direction filterDirection, addr []byte, mask net.IPMask, onMatch, onMiss labelID) {
	if len(addr) < net.IPv6len {
		b.emitJump(onMiss)
		return
	}
	addrArray := [4]uint32{
		binary.BigEndian.Uint32(addr[:4]),
		binary.BigEndian.Uint32(addr[4:8]),
		binary.BigEndian.Uint32(addr[8:12]),
		binary.BigEndian.Uint32(addr[12:16]),
	}
	switch direction {
	case filterDirectionSrc:
		b.loadAndCompareIPv6Address(addrArray, mask, true, onMatch, onMiss)
	case filterDirectionDst:
		b.loadAndCompareIPv6Address(addrArray, mask, false, onMatch, onMiss)
	case filterDirectionSrcOrDst:
		cont := b.newLabel()
		b.loadAndCompareIPv6Address(addrArray, mask, true, onMatch, cont)
		b.bind(cont)
		b.loadAndCompareIPv6Address(addrArray, mask, false, onMatch, onMiss)
	case filterDirectionSrcAndDst:
		cont := b.newLabel()
		b.loadAndCompareIPv6Address(addrArray, mask, true, cont, onMiss)
		b.bind(cont)
		b.loadAndCompareIPv6Address(addrArray, mask, false, onMatch, onMiss)
	}
}

func (b *prog) loadAndCompareIPv6Address(addr [4]uint32, mask net.IPMask, source bool, onMatch, onMiss labelID) {
	maskSize := 128
	var maskInst bpf.Instruction
	start := b.layout.l3Off() + intraIP6SrcAddrStart
	if !source {
		start = b.layout.l3Off() + intraIP6DstAddrStart
	}
	if mask != nil {
		var maskBits int
		maskSize, maskBits = mask.Size()
		if maskSize < 0 || maskBits != net.IPv6len*8 {
			b.emitJump(onMiss)
			return
		}
		partWords := maskSize % bitsPerWord
		if partWords != 0 {
			maskStartOff := (maskSize / bitsPerWord) * 4
			maskTerm := binary.BigEndian.Uint32(mask[maskStartOff : maskStartOff+4])
			if maskTerm != 0xffffffff {
				maskInst = bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: maskTerm}
			}
		}
	}

	bitsUsed := 0
	for i, a := range addr {
		b.emit(bpf.LoadAbsolute{Off: start + uint32(i*4), Size: 4})
		bitsUsed += bitsPerWord
		if bitsUsed > maskSize && maskInst != nil {
			b.emit(maskInst)
		}
		if bitsUsed >= maskSize {
			b.emitJumpIf(bpf.JumpEqual, a, onMatch, onMiss)
			return
		}
		cont := b.newLabel()
		b.emitJumpIf(bpf.JumpEqual, a, cont, onMiss)
		b.bind(cont)
	}
}

func (b *prog) checkPorts(direction filterDirection, port uint32, onMatch, onMiss labelID, ip6 bool) {
	var loadSource, loadDestination bpf.Instruction
	if ip6 {
		loadSource = loadIPv6SourcePort(b.layout)
		loadDestination = loadIPv6DestinationPort(b.layout)
	} else {
		loadSource = loadIPv4SourcePort(b.layout)
		loadDestination = loadIPv4DestinationPort(b.layout)
		b.loadIPv4HeaderOffset(onMiss)
	}

	switch direction {
	case filterDirectionSrc:
		b.emit(loadSource)
		b.emitJumpIf(bpf.JumpEqual, port, onMatch, onMiss)
	case filterDirectionDst:
		b.emit(loadDestination)
		b.emitJumpIf(bpf.JumpEqual, port, onMatch, onMiss)
	case filterDirectionSrcOrDst:
		cont := b.newLabel()
		b.emit(loadSource)
		b.emitJumpIf(bpf.JumpEqual, port, onMatch, cont)
		b.bind(cont)
		b.emit(loadDestination)
		b.emitJumpIf(bpf.JumpEqual, port, onMatch, onMiss)
	case filterDirectionSrcAndDst:
		cont := b.newLabel()
		b.emit(loadSource)
		b.emitJumpIf(bpf.JumpEqual, port, cont, onMiss)
		b.bind(cont)
		b.emit(loadDestination)
		b.emitJumpIf(bpf.JumpEqual, port, onMatch, onMiss)
	}
}

func (b *prog) checkPortRange(
	direction filterDirection,
	minimum uint32,
	maximum uint32,
	onMatch labelID,
	onMiss labelID,
	ip6 bool,
) {
	var loadSource, loadDestination bpf.Instruction
	if ip6 {
		loadSource = loadIPv6SourcePort(b.layout)
		loadDestination = loadIPv6DestinationPort(b.layout)
	} else {
		loadSource = loadIPv4SourcePort(b.layout)
		loadDestination = loadIPv4DestinationPort(b.layout)
		b.loadIPv4HeaderOffset(onMiss)
	}

	check := func(load bpf.Instruction, match, miss labelID) {
		upper := b.newLabel()
		b.emit(load)
		b.emitJumpIf(bpf.JumpGreaterOrEqual, minimum, upper, miss)
		b.bind(upper)
		b.emitJumpIf(bpf.JumpGreaterThan, maximum, miss, match)
	}

	switch direction {
	case filterDirectionSrc:
		check(loadSource, onMatch, onMiss)
	case filterDirectionDst:
		check(loadDestination, onMatch, onMiss)
	case filterDirectionSrcOrDst:
		tryDestination := b.newLabel()
		check(loadSource, onMatch, tryDestination)
		b.bind(tryDestination)
		check(loadDestination, onMatch, onMiss)
	case filterDirectionSrcAndDst:
		checkDestination := b.newLabel()
		check(loadSource, checkDestination, onMiss)
		b.bind(checkDestination)
		check(loadDestination, onMatch, onMiss)
	}
}

func (b *prog) genEtherMulticast(onMatch, onMiss labelID) {
	b.emit(bpf.LoadAbsolute{Off: 0, Size: 1})
	b.emitJumpIf(bpf.JumpBitsSet, 0x01, onMatch, onMiss)
}

func (b *prog) genIPv4Multicast(onMatch, onMiss labelID) {
	b.emit(bpf.LoadAbsolute{Off: b.layout.l3Off() + intraIP4DstAddr, Size: 1})
	b.emitJumpIf(bpf.JumpGreaterOrEqual, 0xe0, onMatch, onMiss)
}

func (b *prog) genIPv6Multicast(onMatch, onMiss labelID) {
	b.emit(bpf.LoadAbsolute{Off: b.layout.l3Off() + intraIP6DstAddrStart, Size: 1})
	b.emitJumpIf(bpf.JumpEqual, 0xff, onMatch, onMiss)
}

// =============================================================================
// Instruction helpers — used by the *prog builder methods above
// =============================================================================

func loadIPv4SourceAddress(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP4SrcAddr, Size: lengthWord}
}

func loadIPv4DestinationAddress(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP4DstAddr, Size: lengthWord}
}

func loadArpSenderAddress(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraArpSenderAddr, Size: lengthWord}
}

func loadArpTargetAddress(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraArpTargetAddr, Size: lengthWord}
}

func loadIPv4Protocol(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP4Protocol, Size: lengthByte}
}

func loadIPv6Protocol(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP6NextHeader, Size: lengthByte}
}

func loadIPv6ContinuationProtocol(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP6ContHdrProto, Size: lengthByte}
}

func loadIPv4SourcePort(layout linkLayout) bpf.LoadIndirect {
	return bpf.LoadIndirect{Off: layout.l3Off() + intraIP4SrcPort, Size: lengthHalf}
}

func loadIPv4DestinationPort(layout linkLayout) bpf.LoadIndirect {
	return bpf.LoadIndirect{Off: layout.l3Off() + intraIP4DstPort, Size: lengthHalf}
}

func loadIPv6SourcePort(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP6SrcPort, Size: lengthHalf}
}

func loadIPv6DestinationPort(layout linkLayout) bpf.LoadAbsolute {
	return bpf.LoadAbsolute{Off: layout.l3Off() + intraIP6DstPort, Size: lengthHalf}
}

var (
	loadEthernetSourceFirst      = bpf.LoadAbsolute{Off: 6, Size: lengthHalf}
	loadEthernetSourceLast       = bpf.LoadAbsolute{Off: 8, Size: lengthWord}
	loadEthernetDestinationFirst = bpf.LoadAbsolute{Off: 0, Size: lengthHalf}
	loadEthernetDestinationLast  = bpf.LoadAbsolute{Off: 2, Size: lengthWord}
)

// =============================================================================
// Network/mask parsing
// =============================================================================

func getNetAndMask(id string) (net.IP, *net.IPNet, error) {
	var (
		addr    net.IP
		network *net.IPNet
		mask    net.IPMask
	)
	if addr := net.ParseIP(id); addr != nil {
		if addr.To4() != nil {
			mask = ip4MaskFull
		} else {
			mask = ip6MaskFull
		}
		network = &net.IPNet{IP: addr, Mask: mask}
		return addr, network, nil
	}
	addr, network, err := net.ParseCIDR(id)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid net: %s", id)
	}
	return addr, network, nil
}
