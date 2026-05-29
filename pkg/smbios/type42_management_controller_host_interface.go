// Copyright 2016-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

// TableTypeManagementControllerHostInterface is the SMBIOS Type 42
const TableTypeManagementControllerHostInterface TableType = 42

// IP is a 16-byte SMBIOS representation of an IP address
type IP [16]byte

// IPAddressFormat represents the format of an IP address.
type IPAddressFormat uint8

const (
	IPAddressFormatUnknown IPAddressFormat = 0
	IPAddressFormatIPv4    IPAddressFormat = 1
	IPAddressFormatIPv6    IPAddressFormat = 2
)

// IPAssignmentType represents the assignment method of an IP address.
type IPAssignmentType uint8

const (
	IPAssignmentTypeUnknown       IPAssignmentType = 0
	IPAssignmentTypeStatic        IPAssignmentType = 1
	IPAssignmentTypeDHCP          IPAssignmentType = 2
	IPAssignmentTypeAutoConfigure IPAssignmentType = 3
	IPAssignmentTypeHostSelected  IPAssignmentType = 4
)

// BmcIPDiscoveryType represents the BMC IP address discovery method.
type BmcIPDiscoveryType uint8

const (
	BmcIPDiscoveryTypeUnknown       BmcIPDiscoveryType = 0
	BmcIPDiscoveryTypeStatic        BmcIPDiscoveryType = 1
	BmcIPDiscoveryTypeDHCP          BmcIPDiscoveryType = 2
	BmcIPDiscoveryTypeAutoConfigure BmcIPDiscoveryType = 3
	BmcIPDiscoveryTypeHostSelected  BmcIPDiscoveryType = 4
)

// NetIP converts the SMBIOS IP representation to Go's standard net.IP
func (ip IP) NetIP() net.IP {
	return net.IP(ip[:])
}

// ParseField parses an IP field within a table.
func (ip *IP) ParseField(t *Table, off int) (int, error) {
	b, err := t.GetBytesAt(off, 16)
	if err != nil {
		return off, err
	}
	copy(ip[:], b)
	return off + 16, nil
}

// RedfishHostInterfaceConfig holds all parameters needed to configure a Redfish Host Interface over USB(Type 42).
type RedfishHostInterfaceConfig struct {
	InterfaceType        uint8
	USBDeviceType        uint8
	VendorID             uint16
	ProductID            uint16
	UsbLength            uint8
	UsbDescriptorType    uint8
	SerialNumber         string
	HostIPAssignmentType uint8
	BmcIPDiscoveryType   uint8
	ServiceUUID          UUID
	HostIP               net.IP
	HostMask             net.IPMask
	BmcIP                net.IP
	BmcMask              net.IPMask
	RedfishPort          uint16
	VlanID               uint32
	Hostname             string
}

// ManagementControllerHostInterface is defined in DSP0134 7.43.
// This layout specifically models the USB physical interface carrying the Redfish-over-IP protocol.
type ManagementControllerHostInterface struct {
	Header
	InterfaceType                   uint8  // 04h physical interface type (e.g. 40h Network Host)
	InterfaceTypeSpecificDataLength uint8  // 05h size of interface specific data (8 bytes)
	DeviceType                      uint8  // 06h USB Network Interface device type (e.g. 02h)
	VendorID                        uint16 // 07h USB Vendor ID (idVendor)
	ProductID                       uint16 // 09h USB Product ID (idProduct)
	UsbLength                       uint8  // 0Bh USB descriptor bLength (2)
	UsbDescriptorType               uint8  // 0Ch USB descriptor bDescriptorType (3)
	SerialNumber                    string // 0Dh USB SerialNumber string index (resolves to trailing string)
	ProtocolCount                   uint8  // 0Eh number of supported protocols (1)
	ProtocolIdentifier              uint8  // 0Fh protocol identifier (e.g. 04h Redfish over IP)
	ProtocolLength                  uint8  // 10h size of Redfish protocol data block
	ServiceUUID                     UUID   // 11h Redfish service UUID
	HostIPAssignmentType            uint8  // 21h
	HostIPAddressFormat             uint8  // 22h
	HostIPAddress                   IP     // 23h
	HostIPMask                      IP     // 33h
	BmcIPDiscoveryType              uint8  // 43h
	BmcIPAddressFormat              uint8  // 44h
	BmcIPAddress                    IP     // 45h
	BmcIPMask                       IP     // 55h
	RedfishPort                     uint16 // 65h
	VlanID                          uint32 // 67h
	Hostname                        string // 6Bh
}

// ParseManagementControllerHostInterface parses a generic Table into ManagementControllerHostInterface.
func ParseManagementControllerHostInterface(t *Table) (*ManagementControllerHostInterface, error) {
	return parseManagementControllerHostInterface(parseStruct, t)
}

func parseManagementControllerHostInterface(parseFn parseStructure, t *Table) (*ManagementControllerHostInterface, error) {
	if t.Type != TableTypeManagementControllerHostInterface {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 11 {
		return nil, errors.New("required fields missing")
	}

	mc := &ManagementControllerHostInterface{Header: t.Header}
	if _, err := parseFn(t, 0 /* off */, false /* complete */, mc); err != nil {
		return nil, err
	}

	// If the ProtocolLength is larger than 91, it means the hostname
	// is stored in-line (DSP0270 V1.0.0) rather than as a trailing string index.
	if mc.ProtocolLength > 91 {
		hLen, err := t.GetByteAt(107) // X+90 is offset 107 (since X starts at 17)
		if err == nil && hLen > 0 {
			if hBytes, err := t.GetBytesAt(108, int(hLen)); err == nil {
				mc.Hostname = string(hBytes)
			}
		}
	}

	return mc, nil
}

// MarshalBinary encodes the ManagementControllerHostInterface content into binary format.
func (mc *ManagementControllerHostInterface) MarshalBinary() ([]byte, error) {
	t, err := mc.toTable()
	if err != nil {
		return nil, err
	}
	return t.MarshalBinary()
}

func (mc *ManagementControllerHostInterface) toTable() (*Table, error) {
	var payload []byte
	var tableStr []string
	idx := byte(1)

	payload = append(payload, mc.InterfaceType)
	payload = append(payload, mc.InterfaceTypeSpecificDataLength)
	payload = append(payload, mc.DeviceType)

	vid := make([]byte, 2)
	binary.LittleEndian.PutUint16(vid, mc.VendorID)
	payload = append(payload, vid...)

	pid := make([]byte, 2)
	binary.LittleEndian.PutUint16(pid, mc.ProductID)
	payload = append(payload, pid...)

	payload = append(payload, mc.UsbLength)
	payload = append(payload, mc.UsbDescriptorType)

	if mc.SerialNumber != "" {
		payload = append(payload, idx)
		idx++
		tableStr = append(tableStr, mc.SerialNumber)
	} else {
		payload = append(payload, 0)
	}

	if mc.Hostname != "" && mc.ProtocolLength > 91 {
		mc.ProtocolLength = uint8(90 + 1 + len(mc.Hostname))
	}

	payload = append(payload, mc.ProtocolCount)
	payload = append(payload, mc.ProtocolIdentifier)
	payload = append(payload, mc.ProtocolLength)
	payload = append(payload, mc.ServiceUUID[:]...)
	payload = append(payload, mc.HostIPAssignmentType)
	payload = append(payload, mc.HostIPAddressFormat)
	payload = append(payload, mc.HostIPAddress[:]...)
	payload = append(payload, mc.HostIPMask[:]...)
	payload = append(payload, mc.BmcIPDiscoveryType)
	payload = append(payload, mc.BmcIPAddressFormat)
	payload = append(payload, mc.BmcIPAddress[:]...)
	payload = append(payload, mc.BmcIPMask[:]...)

	port := make([]byte, 2)
	binary.LittleEndian.PutUint16(port, mc.RedfishPort)
	payload = append(payload, port...)

	vlan := make([]byte, 4)
	binary.LittleEndian.PutUint32(vlan, mc.VlanID)
	payload = append(payload, vlan...)

	if mc.Hostname != "" {
		if mc.ProtocolLength > 91 {
			payload = append(payload, byte(len(mc.Hostname)))
			payload = append(payload, []byte(mc.Hostname)...)
		} else {
			payload = append(payload, idx)
			tableStr = append(tableStr, mc.Hostname)
		}
	} else {
		payload = append(payload, 0)
	}

	// Update header length to match payload size + 4-byte header
	mc.Header.Length = uint8(len(payload) + 4)

	h, err := mc.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var d []byte
	d = append(d, h...)
	d = append(d, payload...)

	t := &Table{
		Header:  mc.Header,
		data:    d,
		strings: tableStr,
	}
	return t, nil
}

func (mc *ManagementControllerHostInterface) String() string {
	lines := []string{
		mc.Header.String(),
		fmt.Sprintf("Interface Type: 0x%02X", mc.InterfaceType),
		fmt.Sprintf("Interface Specific Data Length: %d", mc.InterfaceTypeSpecificDataLength),
		fmt.Sprintf("Device Type: 0x%02X", mc.DeviceType),
		fmt.Sprintf("Vendor ID: 0x%04X", mc.VendorID),
		fmt.Sprintf("Product ID: 0x%04X", mc.ProductID),
		fmt.Sprintf("USB bLength: %d", mc.UsbLength),
		fmt.Sprintf("USB bDescriptorType: %d", mc.UsbDescriptorType),
		fmt.Sprintf("Serial Number: %s", mc.SerialNumber),
		fmt.Sprintf("Protocol Count: %d", mc.ProtocolCount),
		fmt.Sprintf("Protocol Identifier: 0x%02X", mc.ProtocolIdentifier),
		fmt.Sprintf("Protocol Length: %d", mc.ProtocolLength),
		fmt.Sprintf("Service UUID: %s", mc.ServiceUUID),
		fmt.Sprintf("Host IP Assignment Type: 0x%02X", mc.HostIPAssignmentType),
		fmt.Sprintf("Host IP Address Format: 0x%02X", mc.HostIPAddressFormat),
		fmt.Sprintf("Host IP Address: %s", mc.HostIPAddress.NetIP()),
		fmt.Sprintf("Host IP Mask: %s", net.IP(mc.HostIPMask[:])),
		fmt.Sprintf("BMC IP Discovery Type: 0x%02X", mc.BmcIPDiscoveryType),
		fmt.Sprintf("BMC IP Address Format: 0x%02X", mc.BmcIPAddressFormat),
		fmt.Sprintf("BMC IP Address: %s", mc.BmcIPAddress.NetIP()),
		fmt.Sprintf("BMC IP Mask: %s", net.IP(mc.BmcIPMask[:])),
		fmt.Sprintf("Redfish Port: %d", mc.RedfishPort),
		fmt.Sprintf("VLAN ID: %d", mc.VlanID),
		fmt.Sprintf("Hostname: %s", mc.Hostname),
	}
	return strings.Join(lines, "\n\t")
}
