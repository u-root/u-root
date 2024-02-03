// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ibft defines the iSCSI Boot Firmware Table.
//
// An iBFT is placed in memory by a bootloader to tell an OS which iSCSI target
// it was booted from and what additional iSCSI targets it shall connect to.
//
// An iBFT is typically placed in one of two places:
//
//	(1) an ACPI table named "iBFT", or
//
//	(2) in the 512K-1M physical memory range identified by its first 4 bytes.
//
// However, this package doesn't concern itself with the placement, just the
// marshaling of the table's bytes.
package ibft

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/u-root/uio/uio"
)

var (
	signature  = [4]byte{'i', 'B', 'F', 'T'}
	oemID      = [6]byte{'G', 'o', 'o', 'g', 'l', 'e'}
	oemTableID = [8]byte{'G', 'o', 'o', 'g', 'I', 'B', 'F', 'T'}
)

type structureID uint8

const (
	reserved structureID = iota
	ibftControlID
	ibftInitiatorID
	ibftNICID
	ibftTargetID
	ibftExtensionsID
)

const (
	ibftControlLen   uint16 = 18
	ibftInitiatorLen uint16 = 74
	ibftNICLen       uint16 = 102
	ibftTargetLen    uint16 = 54
)

type acpiHeader struct {
	Signature       [4]byte
	Length          uint32
	Revision        uint8
	Checksum        uint8
	OEMID           [6]byte
	OEMTableID      [8]byte
	OEMRevision     uint64
	CreatorID       uint64
	CreatorRevision uint64
}

func flags(s ...bool) uint8 {
	var i uint8
	for bit, f := range s {
		if f {
			i |= 1 << uint8(bit)
		} else {
			i &= ^(1 << uint8(bit))
		}
	}
	return i
}

func writeIP6(l *uio.Lexer, ip net.IP) {
	if ip == nil {
		l.WriteBytes(net.IPv6zero)
	} else {
		l.WriteBytes(ip.To16())
	}
}

// ibftStructHeader is the header of each iBFT structure, as defined by iBFT Spec 1.4.2.
type ibftStructHeader struct {
	StructureID structureID
	Version     uint8
	Length      uint16
	Index       uint8
	Flags       uint8
}

// Initiator defines an iSCSI initiator.
type Initiator struct {
	// Name is the name of the initiator.
	//
	// Some servers apparently use initiator names for permissions.
	Name string

	Valid bool

	// Boot indicates that this initiator was used to boot.
	Boot bool

	// I have no clue what this stuff is.
	SNSServer             net.IP
	SLPServer             net.IP
	PrimaryRadiusServer   net.IP
	SecondaryRadiusServer net.IP
}

func (i *Initiator) marshal(h *heapTable) {
	header := ibftStructHeader{
		StructureID: ibftInitiatorID,
		Version:     1,
		Length:      ibftInitiatorLen,
		Index:       0,
		Flags:       flags(i.Valid, i.Boot),
	}
	h.Table.WriteData(&header)

	writeIP6(h.Table, i.SNSServer)
	writeIP6(h.Table, i.SLPServer)
	writeIP6(h.Table, i.PrimaryRadiusServer)
	writeIP6(h.Table, i.SecondaryRadiusServer)
	h.writeHeap([]byte(i.Name))
}

// BDF is a Bus/Device/Function Identifier of a PCI device.
type BDF struct {
	Bus      uint8
	Device   uint8
	Function uint8
}

func (b BDF) marshal(h *heapTable) {
	h.Table.Write16(uint16(uint16(b.Bus)<<8 | uint16(b.Device)<<3 | uint16(b.Function)))
}

// Origin is the source of network configuration; for example, DHCP or static
// configuration.
//
// The spec links to a Microsoft.com 404 page. We got this info from
// iPXE code, but most likely it refers to the NL_PREFIX_ORIGIN enumeration in
// nldef.h, which can currently be found here
// https://docs.microsoft.com/en-us/windows/win32/api/nldef/ne-nldef-nl_prefix_origin
type Origin uint8

// These values and descriptions were taken from the NL_PREFIX_ORIGIN
// enumeration documentation for Windows.
//
// They can also be found in iPXE source at src/include/ipxe/ibft.h.
const (
	// OriginOther means the IP prefix was provided by a source other than
	// those in this enumeration.
	OriginOther Origin = 0

	// OriginManual means the IP address prefix was manually specified.
	OriginManual Origin = 1

	// OriginWellKnown means the IP address prefix is from a well-known
	// source.
	OriginWellKnown Origin = 2

	// OriginDHCP means the IP addrss prefix was provided by DHCP settings.
	OriginDHCP Origin = 3

	// OriginRA means the IP address prefix was obtained through a router
	// advertisement (RA).
	OriginRA Origin = 4

	// OriginUnchanged menas the IP address prefix should be unchanged.
	// This value is used when setting the properties for a unicast IP
	// interface when the value for the IP prefix origin should be left
	// unchanged.
	OriginUnchanged Origin = 0xf
)

// NIC defines NIC network configuration such as IP, gateway, DNS.
type NIC struct {
	// Valid NIC.
	Valid bool

	// Boot indicates this NIC was used to boot from.
	Boot bool

	// Global indicates a globally reachable IP (as opposed to a link-local IP).
	Global bool

	// IPNet is the IP and subnet mask for this network interface.
	IPNet *net.IPNet

	// Origin is the source of the IP address prefix.
	//
	// It's hard to say what loaded operating systems do with this field.
	//
	// The Windows docs for NL_PREFIX_ORIGIN consider it the origin of only
	// the IP address prefix. iPXE only ever sets the Manual and DHCP
	// constants, even in all IPv6 cases where it came from RAs.
	Origin Origin

	// Gateway is the network gateway. The gateway must be within the IPNet.
	Gateway net.IP

	// PrimaryDNS is the primary DNS server.
	PrimaryDNS net.IP

	// SecondaryDNS is a secondary DNS server.
	SecondaryDNS net.IP

	// DHCPServer is the address for the DHCP server, if one was used.
	DHCPServer net.IP

	// VLAN is the VLAN for this network interface.
	VLAN uint16

	// MACAddress is the MAC of the network interface.
	MACAddress net.HardwareAddr

	// PCIBDF is the Bus/Device/Function identifier of the PCI NIC device.
	PCIBDF BDF

	// HostName is the host name of the machine.
	HostName string
}

func (n *NIC) marshal(h *heapTable) {
	header := ibftStructHeader{
		StructureID: ibftNICID,
		Version:     1,
		Length:      ibftNICLen,
		Index:       0,
		Flags:       flags(n.Valid, n.Boot, n.Global),
	}
	h.Table.WriteData(&header)

	// IP + subnet mask prefix.
	if n.IPNet == nil {
		writeIP6(h.Table, net.IPv6zero)
		h.Table.Write8(0)
	} else {
		writeIP6(h.Table, n.IPNet.IP)
		ones, _ := n.IPNet.Mask.Size()
		h.Table.Write8(uint8(ones))
	}

	h.Table.Write8(uint8(n.Origin))

	writeIP6(h.Table, n.Gateway)
	writeIP6(h.Table, n.PrimaryDNS)
	writeIP6(h.Table, n.SecondaryDNS)
	writeIP6(h.Table, n.DHCPServer)

	h.Table.Write16(n.VLAN)
	copy(h.Table.WriteN(6), n.MACAddress)
	n.PCIBDF.marshal(h)
	h.writeHeap([]byte(n.HostName))
}

// Target carries info about an iSCSI target server.
type Target struct {
	Valid bool
	Boot  bool

	CHAP  bool
	RCHAP bool

	// Target is the server to connect to.
	Target *net.TCPAddr

	// BootLUN is the LUN to connect to.
	BootLUN uint64

	CHAPType uint8

	// NICAssociation is the Index of the NIC.
	NICAssociation uint8

	// TargetName is the name of the iSCSI target.
	//
	// The target name must be a valid iSCSI Qualifier Name or EUI.
	TargetName string

	CHAPName          string
	CHAPSecret        string
	ReverseCHAPName   string
	ReverseCHAPSecret string
}

func (t *Target) marshal(h *heapTable) {
	header := ibftStructHeader{
		StructureID: ibftTargetID,
		Version:     1,
		Length:      ibftTargetLen,
		Index:       0,
		Flags:       flags(t.Valid, t.Boot, t.CHAP, t.RCHAP),
	}
	h.Table.WriteData(&header)

	if t.Target == nil {
		writeIP6(h.Table, net.IPv6zero)
		h.Table.Write16(0)
	} else {
		writeIP6(h.Table, t.Target.IP)
		h.Table.Write16(uint16(t.Target.Port))
	}
	h.Table.Write64(t.BootLUN)
	h.Table.Write8(t.CHAPType)
	h.Table.Write8(t.NICAssociation)

	h.writeHeap([]byte(t.TargetName))
	h.writeHeap([]byte(t.CHAPName))
	h.writeHeap([]byte(t.CHAPSecret))
	h.writeHeap([]byte(t.ReverseCHAPName))
	h.writeHeap([]byte(t.ReverseCHAPSecret))
}

// IBFT defines the entire iSCSI boot firmware table.
type IBFT struct {
	SingleLoginMode bool

	// Initiator offset: 0x50
	Initiator Initiator

	// NIC offset: 0xa0
	NIC0 NIC

	// Target offset: 0x110
	Target0 Target
}

// String is a short summary of the iBFT contents.
func (i *IBFT) String() string {
	return fmt.Sprintf("iBFT(iSCSI target=%s, IP=%s)", i.Target0.Target, i.NIC0.IPNet)
}

type heapTable struct {
	heapOffset uint64

	Table *uio.Lexer
	Heap  *uio.Lexer
}

func (h *heapTable) writeHeap(item []byte) {
	if len(item) == 0 {
		// Length.
		h.Table.Write16(0)

		// Offset.
		h.Table.Write16(0)
	} else {
		offset := h.heapOffset + uint64(h.Heap.Len())

		// Length.
		h.Table.Write16(uint16(len(item)))

		// Offset from beginning of iBFT.
		h.Table.Write16(uint16(offset))

		// Write a null-terminated array item.
		//
		// iBFT Spec, Section 1.3.5: "All array items stored in the
		// Heap area will be followed by a separate NULL. This
		// terminating NULL is not counted as part of the array
		// [item's] length."
		h.Heap.WriteBytes(item)
		h.Heap.Write8(0)
	}
}

type ibftControl struct {
	SingleLoginMode bool
	Initiator       uint16
	NIC0            uint16
	Target0         uint16
	NIC1            uint16
	Target1         uint16
}

func (c ibftControl) marshal(h *heapTable) {
	header := ibftStructHeader{
		StructureID: ibftControlID,
		Version:     1,
		Length:      ibftControlLen,
		Index:       0,
		Flags:       flags(c.SingleLoginMode),
	}
	h.Table.WriteData(&header)

	// Extensions: none. EVER.
	h.Table.Write16(0)

	h.Table.Write16(c.Initiator)
	h.Table.Write16(c.NIC0)
	h.Table.Write16(c.Target0)
	h.Table.Write16(c.NIC1)
	h.Table.Write16(c.Target1)
}

func gencsum(b []byte) byte {
	var csum byte
	for _, bb := range b {
		csum += bb
	}
	return ^csum + 1
}

const (
	lengthOffset   = 4
	checksumOffset = 9
)

func fixACPIHeader(b []byte) []byte {
	binary.LittleEndian.PutUint16(b[lengthOffset:], uint16(len(b)))
	b[checksumOffset] = gencsum(b)
	return b
}

const (
	controlOffset   = 0x30
	initiatorOffset = 0x48
	nic0Offset      = 0x98
	target0Offset   = 0x100
	heapOffset      = 0x138
)

// Marshal returns a binary representation of the iBFT.
//
// Pointers within an iBFT is relative, so this can be placed anywhere
// necessary.
func (i *IBFT) Marshal() []byte {
	h := &heapTable{
		heapOffset: heapOffset,

		Table: uio.NewLittleEndianBuffer(nil),
		Heap:  uio.NewLittleEndianBuffer(nil),
	}

	header := &acpiHeader{
		Signature:       signature,
		Length:          0,
		Revision:        1,
		Checksum:        0,
		OEMID:           oemID,
		OEMTableID:      oemTableID,
		OEMRevision:     0,
		CreatorID:       0,
		CreatorRevision: 0,
	}
	h.Table.WriteData(header)

	// iBFT spec, Section 1.4.4.4 "Each structure must be aligned on an 8
	// byte boundary."

	// 0x30
	h.Table.Align(8)
	control := ibftControl{
		SingleLoginMode: i.SingleLoginMode,
		Initiator:       initiatorOffset,
		NIC0:            nic0Offset,
		Target0:         target0Offset,
	}
	control.marshal(h)

	// 0x48
	h.Table.Align(8)
	i.Initiator.marshal(h)

	// 0x98
	h.Table.Align(8)
	i.NIC0.marshal(h)

	// 0x100
	h.Table.Align(8)
	i.Target0.marshal(h)

	// 0x138
	h.Table.Align(8)
	h.Table.WriteBytes(h.Heap.Data())

	return fixACPIHeader(h.Table.Data())
}
