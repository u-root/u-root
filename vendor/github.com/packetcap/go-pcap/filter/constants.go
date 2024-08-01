package filter

const (
	lengthByte                 int    = 1
	lengthHalf                 int    = 2
	lengthWord                 int    = 4
	bitsPerWord                int    = 32
	etherTypeIPv4              uint32 = 0x0800
	etherTypeIPv6              uint32 = 0x86dd
	etherTypeArp               uint32 = 0x806
	etherTypeRarp              uint32 = 0x8035
	jumpMask                   uint32 = 0x1fff
	ipProtocolTCP              uint32 = 0x06
	ipProtocolUDP              uint32 = 0x11
	ipProtocolSctp             uint32 = 0x84
	ip6SourcePort              uint32 = 54
	ip6DestinationPort         uint32 = 56
	ip4SourcePort              uint32 = 14
	ip4DestinationPort         uint32 = 16
	ip4HeaderSize              uint32 = 14
	ip4HeaderFlags             uint32 = 20
	ip6SourceAddressStart      uint32 = 22
	ip6DestinationAddressStart uint32 = 38
	ip6ContinuationPacket      uint32 = 0x2c
)

type filterKind int

const (
	filterKindUnset filterKind = iota
	filterKindHost
	filterKindNet
	filterKindPort
	filterKindPortRange
)

//nolint:unused
var kinds = map[string]filterKind{
	"host":      filterKindHost,
	"net":       filterKindNet,
	"port":      filterKindPort,
	"portrange": filterKindPortRange,
}
var kinds2 = map[ExpressionToken]filterKind{
	tokenHost:      filterKindHost,
	tokenNet:       filterKindNet,
	tokenPort:      filterKindPort,
	tokenPortRange: filterKindPortRange,
}

type filterDirection int

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

type filterProtocol int

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
	"decnett": filterProtocolDecnet,
}

type filterSubProtocol int

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
	filterSubProtocolStp
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
	"stp":     filterSubProtocolStp,
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
}
