// Copyright 2016-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"net"
	"testing"
)

func TestType42RedfishWithTrailingStringHostname(t *testing.T) {
	// Header (4 bytes) + physical & USB data (10 bytes)
	headerAndInterfaceBytes := []byte{
		42, 108, 0x01, 0x00, // Type=42, Length=108, Handle=0x0001
		0x40,       // InterfaceType Network Host
		8,          // InterfaceTypeSpecificDataLength
		2,          // DeviceType USB
		0x12, 0x34, // VendorID
		0x01, 0x00, // ProductID
		2, // UsbLength
		3, // UsbDescriptorType
		1, // SerialNumber string index 1
	}

	uuidVal := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	hostIP := net.ParseIP("fe80::211:32ff:fe54:8801").To16()
	hostMask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::").To16())
	bmcIP := net.ParseIP("fe80::1").To16()
	bmcMask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::").To16())

	protoData := []byte{}
	protoData = append(protoData, uuidVal...)
	protoData = append(protoData, 1) // HostIPAssignmentType Static
	protoData = append(protoData, 2) // HostIPAddressFormat IPv6
	protoData = append(protoData, hostIP...)
	protoData = append(protoData, hostMask...)
	protoData = append(protoData, 1) // BmcIPDiscoveryType Static
	protoData = append(protoData, 2) // BmcIPAddressFormat IPv6
	protoData = append(protoData, bmcIP...)
	protoData = append(protoData, bmcMask...)
	protoData = append(protoData, 0xB6, 0x27) // RedfishPort 10166
	protoData = append(protoData, 0, 0, 0, 0) // VlanID 0
	protoData = append(protoData, 2)          // Hostname string index 2

	payload := []byte{}
	payload = append(payload, 1)  // ProtocolCount 1
	payload = append(payload, 4)  // ProtocolIdentifier Redfish
	payload = append(payload, 91) // ProtocolLength 91 (includes 1-byte Hostname index)
	payload = append(payload, protoData...)

	rawTableData := append([]byte(nil), headerAndInterfaceBytes...)
	rawTableData = append(rawTableData, payload...)

	strings := []string{"SerialNumber123", "bmc-hostname"}

	tObj := &Table{
		Header: Header{
			Type:   42,
			Length: uint8(len(rawTableData)),
			Handle: 0x0001,
		},
		data:    rawTableData,
		strings: strings,
	}

	// 1. Parse and verify
	mc, err := ParseManagementControllerHostInterface(tObj)
	if err != nil {
		t.Fatalf("Failed to parse Type 42 IPv6 with string index hostname: %v", err)
	}

	if mc.Hostname != "bmc-hostname" {
		t.Errorf("Expected resolved hostname 'bmc-hostname', got '%s'", mc.Hostname)
	}

	// 2. Marshal and verify
	tableMarshaled, err := mc.toTable()
	if err != nil {
		t.Fatalf("Failed to marshal back: %v", err)
	}

	if !bytes.Equal(tableMarshaled.data, tObj.data) {
		t.Errorf("Marshaled data mismatch")
	}
}

func TestType42RedfishWithInLineHostname(t *testing.T) {
	headerAndInterfaceBytes := []byte{
		42, 123, 0x01, 0x00, // Type=42, Length=123, Handle=0x0001
		0x40,       // InterfaceType Network Host
		8,          // InterfaceTypeSpecificDataLength
		2,          // DeviceType USB
		0x12, 0x34, // VendorID
		0x01, 0x00, // ProductID
		2, // UsbLength
		3, // UsbDescriptorType
		1, // SerialNumber string index 1
	}

	uuidVal := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	hostIP := net.ParseIP("fe80::211:32ff:fe54:8801").To16()
	hostMask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::").To16())
	bmcIP := net.ParseIP("fe80::1").To16()
	bmcMask := net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::").To16())

	protoData := []byte{}
	protoData = append(protoData, uuidVal...)
	protoData = append(protoData, 1) // HostIPAssignmentType Static
	protoData = append(protoData, 2) // HostIPAddressFormat IPv6
	protoData = append(protoData, hostIP...)
	protoData = append(protoData, hostMask...)
	protoData = append(protoData, 1) // BmcIPDiscoveryType Static
	protoData = append(protoData, 2) // BmcIPAddressFormat IPv6
	protoData = append(protoData, bmcIP...)
	protoData = append(protoData, bmcMask...)
	protoData = append(protoData, 0xB6, 0x27) // RedfishPort 10166
	protoData = append(protoData, 0, 0, 0, 0) // VlanID 0
	protoData = append(protoData, 15)         // Hostname string Length (15)
	protoData = append(protoData, []byte("inline-bmc-host")...)

	payload := []byte{}
	payload = append(payload, 1)   // ProtocolCount 1
	payload = append(payload, 4)   // ProtocolIdentifier Redfish
	payload = append(payload, 106) // ProtocolLength 106 (90 + 1 + 15 in-line bytes)
	payload = append(payload, protoData...)

	rawTableData := append([]byte(nil), headerAndInterfaceBytes...)
	rawTableData = append(rawTableData, payload...)

	strings := []string{"SerialNumber123"} // SerialNumber string section

	tObj := &Table{
		Header: Header{
			Type:   42,
			Length: uint8(len(rawTableData)),
			Handle: 0x0001,
		},
		data:    rawTableData,
		strings: strings,
	}

	// 1. Parse and verify
	mc, err := ParseManagementControllerHostInterface(tObj)
	if err != nil {
		t.Fatalf("Failed to parse Type 42 with in-line hostname: %v", err)
	}

	if mc.Hostname != "inline-bmc-host" {
		t.Errorf("Expected parsed in-line hostname 'inline-bmc-host', got '%s'", mc.Hostname)
	}

	// 2. Marshal and verify
	tableMarshaled, err := mc.toTable()
	if err != nil {
		t.Fatalf("Failed to marshal back: %v", err)
	}

	if !bytes.Equal(tableMarshaled.data, tObj.data) {
		t.Errorf("Marshaled data mismatch")
	}
}

func TestType42ErrorCases(t *testing.T) {
	// 1. Invalid table type (e.g., Type 1 System Info)
	tObjInvalidType := &Table{
		Header: Header{
			Type:   TableTypeSystemInfo,
			Length: 8,
			Handle: 0x0001,
		},
		data: []byte{1, 8, 0x01, 0x00, 0, 0, 0, 0},
	}
	if _, err := ParseManagementControllerHostInterface(tObjInvalidType); err == nil {
		t.Error("Expected error parsing non-Type 42 table, got nil")
	}

	// 2. Truncated data (shorter than 11 bytes)
	tObjTruncated := &Table{
		Header: Header{
			Type:   42,
			Length: 8,
			Handle: 0x0001,
		},
		data: []byte{42, 8, 0x01, 0x00, 0x40, 8, 2, 0},
	}
	if _, err := ParseManagementControllerHostInterface(tObjTruncated); err == nil {
		t.Error("Expected error parsing truncated Type 42 table, got nil")
	}
}

func TestType42String(t *testing.T) {
	tObj := &Table{
		Header: Header{
			Type:   42,
			Length: 108,
			Handle: 0x0001,
		},
		data: make([]byte, 108),
	}
	tObj.data[0] = 42
	tObj.data[1] = 108
	tObj.data[4] = 0x40
	tObj.data[5] = 8

	mc, err := ParseManagementControllerHostInterface(tObj)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	str := mc.String()
	if str == "" {
		t.Error("String representation is empty")
	}
}
