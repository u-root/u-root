package filter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/net/bpf"
)

// Compile take a filter string compatible with tcpdump at
// https://www.tcpdump.org/manpages/pcap-filter.7.html and return
// bpf instructions

var (
	ip4MaskFull                  = net.CIDRMask(32, 32)   //[]byte{0xff, 0xff, 0xff, 0xff}
	ip6MaskFull                  = net.CIDRMask(128, 128) //[]byte{0xff, 0xff, 0xff, 0xff,0xff, 0xff, 0xff, 0xff,0xff, 0xff, 0xff, 0xff,0xff, 0xff, 0xff, 0xff}
	returnDrop                   = bpf.RetConstant{Val: 0}
	returnKeep                   = bpf.RetConstant{Val: 0x40000}
	loadIPv4SourcePort           = bpf.LoadIndirect{Off: ip4SourcePort, Size: lengthHalf}
	loadIPv4DestinationPort      = bpf.LoadIndirect{Off: ip4DestinationPort, Size: lengthHalf}
	loadIPv6SourcePort           = bpf.LoadAbsolute{Off: ip6SourcePort, Size: lengthHalf}
	loadIPv6DestinationPort      = bpf.LoadAbsolute{Off: ip6DestinationPort, Size: lengthHalf}
	loadEtherKind                = bpf.LoadAbsolute{Off: 12, Size: lengthHalf}
	loadIPv4SourceAddress        = bpf.LoadAbsolute{Off: 26, Size: lengthWord}
	loadIPv4DestinationAddress   = bpf.LoadAbsolute{Off: 30, Size: lengthWord}
	loadArpSenderAddress         = bpf.LoadAbsolute{Off: 28, Size: lengthWord}
	loadArpTargetAddress         = bpf.LoadAbsolute{Off: 38, Size: lengthWord}
	loadIPv4Protocol             = bpf.LoadAbsolute{Off: 23, Size: lengthByte}
	loadIPv6Protocol             = bpf.LoadAbsolute{Off: 20, Size: lengthByte}
	loadIPv6ContinuationProtocol = bpf.LoadAbsolute{Off: 54, Size: lengthByte}
	loadEthernetSourceFirst      = bpf.LoadAbsolute{Off: 6, Size: lengthHalf}
	loadEthernetSourceLast       = bpf.LoadAbsolute{Off: 8, Size: lengthWord}
	loadEthernetDestinationFirst = bpf.LoadAbsolute{Off: 0, Size: lengthHalf}
	loadEthernetDestinationLast  = bpf.LoadAbsolute{Off: 2, Size: lengthWord}
)

func loadIPv4HeaderOffset(skipFail uint8) []bpf.Instruction {
	return []bpf.Instruction{
		bpf.LoadAbsolute{Off: ip4HeaderFlags, Size: lengthHalf},                  // flags+fragment offset, since we need to calc where the src/dst port is
		bpf.JumpIf{Cond: bpf.JumpBitsSet, Val: jumpMask, SkipTrue: skipFail - 1}, // do we have an L4 header?
		bpf.LoadMemShift{Off: ip4HeaderSize},                                     // calculate size of IP header
	}
}

func compareProtocolIP4(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: etherTypeIPv4, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareProtocolIP6(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: etherTypeIPv6, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareProtocolArp(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: etherTypeArp, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareProtocolRarp(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: etherTypeRarp, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareSubProtocolTCP(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: ipProtocolTCP, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareSubProtocolUDP(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: ipProtocolUDP, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareSubProtocolSctp(skipTrue, skipFalse uint8) bpf.Instruction {
	return bpf.JumpIf{Cond: bpf.JumpEqual, Val: ipProtocolSctp, SkipFalse: skipFalse, SkipTrue: skipTrue}
}

func compareIPv6Protocol(proto uint32, skipTrue, skipFalse uint8) []bpf.Instruction {
	st, sf := skipTrue, skipFalse
	if st == 0 {
		st = 4
	}
	if sf == 0 {
		sf = 4
	}
	return []bpf.Instruction{
		loadIPv6Protocol,
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: proto, SkipFalse: 0, SkipTrue: st - 1},
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: ip6ContinuationPacket, SkipFalse: sf - 2},
		loadIPv6ContinuationProtocol,
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: proto, SkipFalse: sf - 4, SkipTrue: st - 4},
	}
}

func compareIPv4Protocol(proto uint32, skipTrue, skipFalse uint8) []bpf.Instruction {
	st, sf := skipTrue, skipFalse
	if st == 0 {
		st = 1
	}
	if sf == 0 {
		sf = 1
	}
	return []bpf.Instruction{
		loadIPv4Protocol,
		bpf.JumpIf{Cond: bpf.JumpEqual, Val: proto, SkipFalse: sf - 1, SkipTrue: st - 1},
	}
}

// checkEtherAddresses add steps to check Ethernet addresses
// fail and succeed are the number of steps to skip the succeed or fail instructions.
// For example, if the next one is succeed, then succeed will be 0
func checkEtherAddresses(direction filterDirection, addr string, fail, succeed uint8) []bpf.Instruction {
	inst := make([]bpf.Instruction, 0)
	// ignore errors as we already validated
	hwAddr, _ := net.ParseMAC(addr)
	if hwAddr == nil {
		return nil
	}
	// need last 4 bytes and first 2 bytes separately
	lastFour := binary.BigEndian.Uint32(hwAddr[len(hwAddr)-4:])
	firstTwo := uint32(binary.BigEndian.Uint16(hwAddr[len(hwAddr)-6 : len(hwAddr)-4]))

	switch direction {
	case filterDirectionSrc:
		inst = append(inst, loadEthernetSourceLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: fail - 1})
		inst = append(inst, loadEthernetSourceFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipTrue: succeed - 3, SkipFalse: fail - 3})
	case filterDirectionDst:
		inst = append(inst, loadEthernetDestinationLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: fail - 1})
		inst = append(inst, loadEthernetDestinationFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipTrue: succeed - 3, SkipFalse: fail - 3})
	case filterDirectionSrcOrDst:
		inst = append(inst, loadEthernetSourceLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: 2})
		inst = append(inst, loadEthernetSourceFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipTrue: succeed - 3})
		inst = append(inst, loadEthernetDestinationLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: fail - 5})
		inst = append(inst, loadEthernetDestinationFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipTrue: succeed - 7, SkipFalse: fail - 7})
	case filterDirectionSrcAndDst:
		inst = append(inst, loadEthernetSourceLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: fail - 1})
		inst = append(inst, loadEthernetSourceFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipFalse: fail - 3})
		inst = append(inst, loadEthernetDestinationLast)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: lastFour, SkipFalse: fail - 5})
		inst = append(inst, loadEthernetDestinationFirst)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: firstTwo, SkipFalse: fail - 7})
	}
	return inst
}

// checkIP4HostAddresses check for host addresses
func checkIP4HostAddresses(direction filterDirection, addr net.IP, fail, succeed uint8) []bpf.Instruction {
	return checkIP4Addresses(direction, addr, nil, fail, succeed, loadIPv4SourceAddress, loadIPv4DestinationAddress)
}

// checkIP4ArpAddresses check for arp addresses
func checkIP4ArpAddresses(direction filterDirection, addr net.IP, fail, succeed uint8) []bpf.Instruction {
	return checkIP4Addresses(direction, addr, nil, fail, succeed, loadArpSenderAddress, loadArpTargetAddress)
}

func checkIP4NetAddresses(direction filterDirection, addr string, ip bool, fail, succeed uint8) []bpf.Instruction {
	// maskCheck is used for networks where a CIDR is supplied, so we need to check if the mask is valid
	// ignore error since it already was validated
	addrBytes, network, _ := getNetAndMask(addr)
	if addrBytes == nil {
		return nil
	}
	var maskCheck *bpf.ALUOpConstant
	if !bytes.Equal(network.Mask, ip4MaskFull) {
		maskCheck = &bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: binary.BigEndian.Uint32(network.Mask)}
	}
	loadSource, loadDestination := loadIPv4SourceAddress, loadIPv4DestinationAddress
	if !ip {
		loadSource, loadDestination = loadArpSenderAddress, loadArpTargetAddress
	}
	return checkIP4Addresses(direction, addrBytes, maskCheck, fail, succeed, loadSource, loadDestination)
}

func checkIP4NetHostAddresses(direction filterDirection, addr string, fail, succeed uint8) []bpf.Instruction {
	return checkIP4NetAddresses(direction, addr, true, fail, succeed)
}
func checkIP4NetArpAddresses(direction filterDirection, addr string, fail, succeed uint8) []bpf.Instruction {
	return checkIP4NetAddresses(direction, addr, false, fail, succeed)
}

// checkIP4Addresses add steps to check IPv4 addresses
// fail and succeed are the number of steps to skip the succeed or fail instructions.
// For example, if the next one is succeed, then succeed will be 0
func checkIP4Addresses(direction filterDirection, addr []byte, maskCheck *bpf.ALUOpConstant, fail, succeed uint8, loadSource, loadTarget bpf.Instruction) []bpf.Instruction {
	inst := make([]bpf.Instruction, 0)
	if addr == nil {
		return nil
	}

	// need last 4 bytes for ipv4
	addrVal := binary.BigEndian.Uint32(addr[len(addr)-4:])

	switch direction {
	case filterDirectionSrc:
		inst = append(inst, loadSource)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipTrue: succeed - uint8(len(inst)), SkipFalse: fail - uint8(len(inst))})
	case filterDirectionDst:
		inst = append(inst, loadTarget)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipTrue: succeed - uint8(len(inst)), SkipFalse: fail - uint8(len(inst))})
	case filterDirectionSrcOrDst:
		inst = append(inst, loadSource)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipTrue: succeed - uint8(len(inst))})
		inst = append(inst, loadTarget)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipTrue: succeed - uint8(len(inst)), SkipFalse: fail - uint8(len(inst))})
	case filterDirectionSrcAndDst:
		inst = append(inst, loadSource)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipFalse: fail - uint8(len(inst))})
		inst = append(inst, loadTarget)
		if maskCheck != nil {
			inst = append(inst, *maskCheck)
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: addrVal, SkipTrue: succeed - uint8(len(inst)), SkipFalse: fail - uint8(len(inst))})
	}
	return inst
}

// checkIP6HostAddresses check for host addresses
func checkIP6HostAddresses(direction filterDirection, addr net.IP, fail, succeed uint8) []bpf.Instruction {
	return checkIP6Addresses(direction, addr, nil, fail, succeed)
}

// checkIP6NetAddresses check for net addresses
func checkIP6NetAddresses(direction filterDirection, addr net.IP, mask net.IPMask, fail, succeed uint8) []bpf.Instruction {
	return checkIP6Addresses(direction, addr, mask, fail, succeed)
}

// checkIP6Addresses add steps to check IPv6 addresses
// fail and succeed are the number of steps to skip the succeed or fail instructions.
// For example, if the next one is succeed, then succeed will be 0
func checkIP6Addresses(direction filterDirection, addr []byte, mask net.IPMask, fail, succeed uint8) []bpf.Instruction {
	inst := make([]bpf.Instruction, 0)

	// need each chunk of 4 bytes
	addrArray := [4]uint32{binary.BigEndian.Uint32(addr[:4]), binary.BigEndian.Uint32(addr[4:8]), binary.BigEndian.Uint32(addr[8:12]), binary.BigEndian.Uint32(addr[12:16])}

	// add the netmask calculation, if it is provided

	switch direction {
	case filterDirectionSrc:
		inst = append(inst, loadAndCompareIPv6SourceAddress(addrArray, mask, succeed, fail)...)
	case filterDirectionDst:
		inst = append(inst, loadAndCompareIPv6DestinationAddress(addrArray, mask, succeed, fail)...)
	case filterDirectionSrcOrDst:
		inst = append(inst, loadAndCompareIPv6SourceAddress(addrArray, mask, succeed, 0)...)
		inst = append(inst, loadAndCompareIPv6DestinationAddress(addrArray, mask, succeed-uint8(len(inst)), fail-uint8(len(inst)))...)
	case filterDirectionSrcAndDst:
		inst = append(inst, loadAndCompareIPv6SourceAddress(addrArray, mask, 0, fail)...)
		inst = append(inst, loadAndCompareIPv6DestinationAddress(addrArray, mask, succeed-uint8(len(inst)), fail-uint8(len(inst)))...)
	}
	return inst
}

// fail and succeed are the number of steps to skip the succeed or fail instructions.
// For example, if the next one is succeed, then succeed will be 0
func checkPorts(direction filterDirection, port uint32, fail, succeed uint8, ip6 bool) []bpf.Instruction {
	inst := make([]bpf.Instruction, 0)

	var (
		loadSource, loadDestination bpf.Instruction
	)

	if ip6 {
		loadSource = loadIPv6SourcePort
		loadDestination = loadIPv6DestinationPort
	} else {
		loadSource = loadIPv4SourcePort
		loadDestination = loadIPv4DestinationPort
		preInst := len(inst)
		inst = append(inst, loadIPv4HeaderOffset(fail)...)
		postInst := len(inst)
		diff := uint8(postInst - preInst)
		//
		fail -= diff
		succeed -= diff
	}

	switch direction {
	case filterDirectionSrc:
		inst = append(inst, loadSource)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipTrue: succeed - 1, SkipFalse: fail - 1})
	case filterDirectionDst:
		inst = append(inst, loadDestination)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipTrue: succeed - 1, SkipFalse: fail - 1})
	case filterDirectionSrcOrDst:
		inst = append(inst, loadSource)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipTrue: succeed - 1})
		inst = append(inst, loadDestination)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipTrue: succeed - 3, SkipFalse: fail - 3})
	case filterDirectionSrcAndDst:
		inst = append(inst, loadSource)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipFalse: fail - 1})
		inst = append(inst, loadDestination)
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: port, SkipTrue: succeed - 3, SkipFalse: fail - 3})
	}
	return inst
}

// getNetAndMask get the address and the network with mask for an IP address.
// If it is *not* CIDR, will return full mask, i.e. 0xffffffff
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
		network = &net.IPNet{
			IP:   addr,
			Mask: mask,
		}
		return addr, network, nil
	}
	addr, network, err := net.ParseCIDR(id)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid net: %s", id)
	}
	return addr, network, nil
}

func calculateIP6MaskSteps(mask net.IPMask) uint8 {
	var count uint8
	// it takes up to 8 steps to check the src or dst, depending on the netmask
	maskSize, _ := mask.Size()
	wholeWords := maskSize / bitsPerWord
	partWords := maskSize % bitsPerWord
	// if it does not split evenly, we need another word and a bitmask line
	if partWords > 0 {
		wholeWords++
	}
	count += 2 * uint8(wholeWords)
	return count
}

// loadAndCompareIPv6SourceAddress check the IP6 source address. skipTrue and skipFalse
// are the number of steps to skip to true or false. If 0, then it means immediately after the
// steps in this section, not absolute. Since the number of steps in this section can change,
// it is important to know if it is absolute (positive number) or just right after (0).
func loadAndCompareIPv6SourceAddress(addr [4]uint32, mask net.IPMask, skipTrue, skipFalse uint8) []bpf.Instruction {
	return loadAndCompareIPv6Address(addr, mask, true, skipTrue, skipFalse)
}

// loadAndCompareIPv6DestinationAddress check the IP6 destination address. skipTrue and skipFalse
// are the number of steps to skip to true or false. If 0, then it means immediately after the
// steps in this section, not absolute. Since the number of steps in this section can change,
// it is important to know if it is absolute (positive number) or just right after (0).
func loadAndCompareIPv6DestinationAddress(addr [4]uint32, mask net.IPMask, skipTrue, skipFalse uint8) []bpf.Instruction {
	return loadAndCompareIPv6Address(addr, mask, false, skipTrue, skipFalse)
}

// loadAndCompareIPv6Address check the IP6 address. skipTrue and skipFalse
// are the number of steps to skip to true or false. If 0, then it means immediately after the
// steps in this section, not absolute. Since the number of steps in this section can change,
// it is important to know if it is absolute (positive number) or just right after (0).
func loadAndCompareIPv6Address(addr [4]uint32, mask net.IPMask, source bool, skipTrue, skipFalse uint8) []bpf.Instruction {
	var (
		maskSize int = 128
		maskInst bpf.Instruction
		start    uint32 = ip6SourceAddressStart
		st, sf   uint8
		// how many steps do we expect?
		size uint8 = 8
	)
	if mask != nil {
		maskSize, _ = mask.Size()
		// every 32 bits = 4 bytes = 1 word
		wholeWords := maskSize / bitsPerWord
		// each whole word requires 2 instructions
		size = 2 * uint8(wholeWords)
		partWords := maskSize % bitsPerWord
		// only apply the mask if it does not end precisely on a word boundary
		if partWords != 0 {
			size += 2
			maskStart := wholeWords * 4
			maskTerm := binary.BigEndian.Uint32(mask[maskStart : maskStart+4])
			if maskTerm != 0xffffffff {
				maskInst = bpf.ALUOpConstant{Op: bpf.ALUOpAnd, Val: maskTerm}
				size++
			}
		}
	}

	if !source {
		start = ip6DestinationAddressStart
	}
	inst := []bpf.Instruction{}

	var bitsUsed = 0
	for i, a := range addr {
		inst = append(inst, bpf.LoadAbsolute{Off: start + uint32(i*4), Size: 4}) // ip6 first 4 bytes
		bitsUsed += bitsPerWord
		if bitsUsed > maskSize {
			inst = append(inst, maskInst)
		}
		st, sf = getSkippers(skipTrue, skipFalse, size, inst)
		if bitsUsed >= maskSize {
			inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: a, SkipTrue: st, SkipFalse: sf})
			return inst
		}
		if i != len(addr)-1 {
			st = 0
		}
		inst = append(inst, bpf.JumpIf{Cond: bpf.JumpEqual, Val: a, SkipTrue: st, SkipFalse: sf})
	}

	return inst
}

// getSkipper calculate how much to skip at a stage in IP addresses.
// At each stage we have SkipTrue and SkipFalse. Here is how we
// calculate it. The rules are the same for either skiptrue/skipfalse
//   - if skip == 0, then at any stage, should show the amount to the end
//     This can be calculated as size-len(inst)-1
//   - if skip != 0, then at any stage, should show (skip - amount used).
//     This can be calculated as skip-len(inst)
func getSkipper(a, size uint8, inst []bpf.Instruction) uint8 {
	l := uint8(len(inst))
	if a == 0 {
		return size - l - 1
	}
	return a - l
}

func getSkippers(a, b, size uint8, inst []bpf.Instruction) (uint8, uint8) {
	return getSkipper(a, size, inst), getSkipper(b, size, inst)
}
