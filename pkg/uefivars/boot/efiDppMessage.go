// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"encoding/binary"
	"fmt"
)

type EfiDppMsgSubType EfiDevPathProtoSubType

const (
	DppMsgTypeATAPI      EfiDppMsgSubType = iota + 1
	DppMsgTypeSCSI                        // 2
	DppMsgTypeFibreCh                     // 3
	DppMsgTypeFirewire                    // 4
	DppMsgTypeUSB                         // 5
	DppMsgTypeIIO                         // 6
	_                                     // 7
	_                                     // 8
	DppMsgTypeInfiniband                  // 9
	DppMsgTypeVendor                      // 10 //uart flow control, sas are subtypes
	DppMsgTypeMAC                         // 11
	DppMsgTypeIP4                         // 12
	DppMsgTypeIP6                         // 13
	DppMsgTypeUART                        // 14
	DppMsgTypeUSBClass                    // 15
	DppMsgTypeUSBWWID                     // 16
	DppMsgTypeDevLU                       // 17
	DppMsgTypeSATA                        // 18
	DppMsgTypeISCSI                       // 19
	DppMsgTypeVLAN                        // 20
	_                                     // 21
	DppMsgTypeSASEx                       // 22
	DppMsgTypeNVME                        // 23
	DppMsgTypeURI                         // 24
	DppMsgTypeUFS                         // 25
	DppMsgTypeSD                          // 26
	DppMsgTypeBT                          // 27
	DppMsgTypeWiFi                        // 28
	DppMsgTypeeMMC                        // 29
	DppMsgTypeBLE                         // 30
	DppMsgTypeDNS                         // 31
	DppMsgTypeNVDIMM                      // 32
	DppMsgTypeRest                        // documented as 32, likely 33
)

var efiDppMsgSubTypeStrings = map[EfiDppMsgSubType]string{
	DppMsgTypeATAPI:      "ATAPI",
	DppMsgTypeSCSI:       "SCSI",
	DppMsgTypeFibreCh:    "Fibre Channel",
	DppMsgTypeFirewire:   "1394",
	DppMsgTypeUSB:        "USB",
	DppMsgTypeIIO:        "I2O",
	DppMsgTypeInfiniband: "Infiniband",
	DppMsgTypeVendor:     "Vendor",
	DppMsgTypeMAC:        "MAC",
	DppMsgTypeIP4:        "IPv4",
	DppMsgTypeIP6:        "IPv6",
	DppMsgTypeUART:       "UART",
	DppMsgTypeUSBClass:   "USB Class",
	DppMsgTypeUSBWWID:    "USB WWID",
	DppMsgTypeDevLU:      "Device Logical Unit",
	DppMsgTypeSATA:       "SATA",
	DppMsgTypeISCSI:      "iSCSI",
	DppMsgTypeVLAN:       "VLAN",
	DppMsgTypeSASEx:      "SAS Ex",
	DppMsgTypeNVME:       "NVME",
	DppMsgTypeURI:        "URI",
	DppMsgTypeUFS:        "UFS",
	DppMsgTypeSD:         "SD",
	DppMsgTypeBT:         "Bluetooth",
	DppMsgTypeWiFi:       "WiFi",
	DppMsgTypeeMMC:       "eMMC",
	DppMsgTypeBLE:        "BLE",
	DppMsgTypeDNS:        "DNS",
	DppMsgTypeNVDIMM:     "NVDIMM",
	DppMsgTypeRest:       "REST",
}

func (e EfiDppMsgSubType) String() string {
	if s, ok := efiDppMsgSubTypeStrings[e]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN-0x%x", uint8(e))
}

// DppMsgATAPI is a struct describing an atapi dpp message.
// pg 293
type DppMsgATAPI struct {
	Hdr             EfiDevicePathProtocolHdr
	Primary, Master bool
	LUN             uint16
}

var _ EfiDevicePathProtocol = (*DppMsgATAPI)(nil)

// ParseDppMsgATAPI parses input into a DppMsgATAPI.
func ParseDppMsgATAPI(h EfiDevicePathProtocolHdr, b []byte) (*DppMsgATAPI, error) {
	if h.Length != 8 {
		return nil, ErrParse
	}
	msg := &DppMsgATAPI{
		Hdr:     h,
		Primary: b[0] == 0,
		Master:  b[1] == 0,
		LUN:     binary.LittleEndian.Uint16(b[2:4]),
	}
	return msg, nil
}

func (e *DppMsgATAPI) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMsgATAPI) ProtoSubTypeStr() string {
	return EfiDppMsgSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMsgATAPI) String() string {
	return fmt.Sprintf("ATAPI(pri=%t,master=%t,lun=%d)", e.Primary, e.Master, e.LUN)
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppMsgATAPI) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}

// DppMsgMAC contains a MAC address.
// pg 300
type DppMsgMAC struct {
	Hdr    EfiDevicePathProtocolHdr
	Mac    [32]byte // 0-padded
	IfType uint8    // RFC3232; seems ethernet is 6
}

// ParseDppMsgMAC parses input into a DppMsgMAC.
func ParseDppMsgMAC(h EfiDevicePathProtocolHdr, b []byte) (*DppMsgMAC, error) {
	if h.Length != 37 {
		return nil, ErrParse
	}
	mac := &DppMsgMAC{
		Hdr: h,
		// Mac:    b[:32],
		IfType: b[32],
	}
	copy(mac.Mac[:], b[:32])
	return mac, nil
}

func (e *DppMsgMAC) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMsgMAC) ProtoSubTypeStr() string {
	return EfiDppMsgSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMsgMAC) String() string {
	switch e.IfType {
	case 1:
		return fmt.Sprintf("MAC(%x)", e.Mac[:6])
	default:
		return fmt.Sprintf("MAC(mac=%08x,iftype=0x%x)", e.Mac, e.IfType)
	}
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppMsgMAC) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}
