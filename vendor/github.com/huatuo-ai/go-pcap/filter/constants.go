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

const (
	lengthByte  int    = 1
	lengthHalf  int    = 2
	lengthWord  int    = 4
	bitsPerWord int    = 32
	jumpMask    uint32 = 0x1fff

	// Protocol numbers (network layer).
	etherTypeIPv4          uint32 = 0x0800
	etherTypeIPv6          uint32 = 0x86dd
	etherTypeArp           uint32 = 0x806
	etherTypeRarp          uint32 = 0x8035
	etherTypeVLAN          uint32 = 0x8100
	etherTypeQinQ          uint32 = 0x88a8
	etherTypeVLAN9100      uint32 = 0x9100
	etherTypeMPLSUnicast   uint32 = 0x8847
	etherTypeMPLSMulticast uint32 = 0x8848

	vlanHeaderLength uint32 = 4
	mplsLabelLength  uint32 = 4

	// Linux socket BPF extensions are encoded as SKF_AD_OFF + SKF_AD_*.
	// They are accepted by SO_ATTACH_FILTER but are not portable packet loads.
	linuxBPFExtVLANTag        uint32 = 0xfffff02c
	linuxBPFExtVLANTagPresent uint32 = 0xfffff030

	ipProtocolTCP   uint32 = 0x06
	ipProtocolUDP   uint32 = 0x11
	ipProtocolSctp  uint32 = 0x84
	ipProtocolIcmp  uint32 = 0x01
	ipProtocolIgmp  uint32 = 0x02
	ipProtocolIcmp6 uint32 = 0x3a
	ipProtocolEsp   uint32 = 0x32
	ipProtocolAh    uint32 = 0x33
	ipProtocolPim   uint32 = 0x67
	ipProtocolVrrp  uint32 = 0x70

	ip6ContinuationPacket uint32 = 0x2c

	// Intra-IP offsets (byte offset from the start of the IP header).
	// Absolute packet offset = layout.l3Off() + intraIPOffset.
	intraIP4HeaderFlags uint32 = 6
	intraIP4Protocol    uint32 = 9
	intraIP4SrcAddr     uint32 = 12
	intraIP4DstAddr     uint32 = 16
	intraIP4SrcPort     uint32 = 0  // relative to end of IP header (LoadIndirect)
	intraIP4DstPort     uint32 = 2  // relative to end of IP header (LoadIndirect)
	intraIP4HeaderSize  uint32 = 0  // offset of IHL nibble within IP header (LoadMemShift)
	intraArpSenderAddr  uint32 = 14 // offset within ARP packet
	intraArpTargetAddr  uint32 = 24 // offset within ARP packet

	intraIP6NextHeader   uint32 = 6
	intraIP6SrcAddrStart uint32 = 8
	intraIP6DstAddrStart uint32 = 24
	intraIP6SrcPort      uint32 = 40 // offset within IPv6 header
	intraIP6DstPort      uint32 = 42 // offset within IPv6 header
	intraIP6ContHdrProto uint32 = 40 // offset of protocol in continuation header
)

type filterKind uint8

const (
	filterKindUnset filterKind = iota
	filterKindHost
	filterKindNet
	filterKindPort
	filterKindPortRange
	filterKindMulticast
	filterKindGateway
	filterKindVLAN
	filterKindMPLS
)

//nolint:unused
var kinds = map[string]filterKind{
	"host":      filterKindHost,
	"net":       filterKindNet,
	"port":      filterKindPort,
	"portrange": filterKindPortRange,
	"multicast": filterKindMulticast,
}

var kinds2 = map[ExpressionToken]filterKind{
	tokenHost:      filterKindHost,
	tokenNet:       filterKindNet,
	tokenPort:      filterKindPort,
	tokenPortRange: filterKindPortRange,
	tokenMulticast: filterKindMulticast,
}

type filterDirection uint8

const (
	filterDirectionUnset filterDirection = iota
	filterDirectionSrcAndDst
	filterDirectionSrcOrDst
	filterDirectionSrc
	filterDirectionDst
	filterDirectionRa
	filterDirectionTa
	filterDirectionAddr1
	filterDirectionAddr2
	filterDirectionAddr3
	filterDirectionAddr4
)

//nolint:unused
var directions = map[string]filterDirection{
	"src":         filterDirectionSrc,
	"dst":         filterDirectionDst,
	"src and dst": filterDirectionSrcAndDst,
	"src or dst":  filterDirectionSrcOrDst,
	"ra":          filterDirectionRa,
	"ta":          filterDirectionTa,
	"addr1":       filterDirectionAddr1,
	"addr2":       filterDirectionAddr2,
	"addr3":       filterDirectionAddr3,
	"addr4":       filterDirectionAddr4,
}

type filterProtocol uint8

const (
	filterProtocolUnset filterProtocol = iota
	filterProtocolEther
	filterProtocolFddi
	filterProtocolTr
	filterProtocolWlan
	filterProtocolIP
	filterProtocolIP6
	filterProtocolArp
	filterProtocolRarp
	filterProtocolDecnet
)

var protocols = map[string]filterProtocol{
	"ether":   filterProtocolEther,
	"fddi":    filterProtocolFddi,
	"tr":      filterProtocolTr,
	"wlan":    filterProtocolWlan,
	"ip":      filterProtocolIP,
	"ip6":     filterProtocolIP6,
	"arp":     filterProtocolArp,
	"rarp":    filterProtocolRarp,
	"decnet":  filterProtocolDecnet,
	"decnett": filterProtocolDecnet,
}

type filterSubProtocol uint8

const (
	filterSubProtocolUnset filterSubProtocol = iota
	filterSubProtocolIP
	filterSubProtocolIP6
	filterSubProtocolArp
	filterSubProtocolRarp
	filterSubProtocolAtalk
	filterSubProtocolAarp
	filterSubProtocolDecnet
	filterSubProtocolSca
	filterSubProtocolLat
	filterSubProtocolMopdl
	filterSubProtocolMoprc
	filterSubProtocolIso
	filterSubProtocolSTP
	filterSubProtocolIPx
	filterSubProtocolNetbeui
	filterSubProtocolIcmp
	filterSubProtocolIcmp6
	filterSubProtocolIgmp
	filterSubProtocolIgrp
	filterSubProtocolPim
	filterSubProtocolAh
	filterSubProtocolEsp
	filterSubProtocolVrrp
	filterSubProtocolUDP
	filterSubProtocolTCP
	filterSubProtocolSCTP
	filterSubProtocolUnknown
)

var subProtocols = map[string]filterSubProtocol{
	"ip":      filterSubProtocolIP,
	"ip6":     filterSubProtocolIP6,
	"arp":     filterSubProtocolArp,
	"rarp":    filterSubProtocolRarp,
	"atalk":   filterSubProtocolAtalk,
	"aarp":    filterSubProtocolAarp,
	"decnet":  filterSubProtocolDecnet,
	"sca":     filterSubProtocolSca,
	"lat":     filterSubProtocolLat,
	"modpl":   filterSubProtocolMopdl,
	"morpc":   filterSubProtocolMoprc,
	"iso":     filterSubProtocolIso,
	"stp":     filterSubProtocolSTP,
	"ipx":     filterSubProtocolIPx,
	"netbeui": filterSubProtocolNetbeui,
	"icmp":    filterSubProtocolIcmp,
	"icmp6":   filterSubProtocolIcmp6,
	"igmp":    filterSubProtocolIgmp,
	"igrp":    filterSubProtocolIgrp,
	"pim":     filterSubProtocolPim,
	"ah":      filterSubProtocolAh,
	"esp":     filterSubProtocolEsp,
	"vrrp":    filterSubProtocolVrrp,
	"udp":     filterSubProtocolUDP,
	"tcp":     filterSubProtocolTCP,
	"sctp":    filterSubProtocolSCTP,
}
