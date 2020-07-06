// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"fmt"
)

// EfiDppHwSubType is the dpp subtype for hardware.
type EfiDppHwSubType EfiDevPathProtoSubType

const (
	DppHTypePCI EfiDppHwSubType = iota + 1
	DppHTypePCCARD
	DppHTypeMMap
	DppHTypeVendor
	DppHTypeCtrl
	DppHTypeBMC
)

var efiDppHwSubTypeStrings = map[EfiDppHwSubType]string{
	DppHTypePCI:    "PCI",
	DppHTypePCCARD: "PCCARD",
	DppHTypeMMap:   "MMap",
	DppHTypeVendor: "Vendor",
	DppHTypeCtrl:   "Control",
	DppHTypeBMC:    "BMC",
}

func (e EfiDppHwSubType) String() string {
	if s, ok := efiDppHwSubTypeStrings[e]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN-0x%x", uint8(e))
}

// DppHwPci is the struct in EfiDevicePathProtocol for DppHTypePCI
type DppHwPci struct {
	Hdr              EfiDevicePathProtocolHdr
	Function, Device uint8
}

var _ EfiDevicePathProtocol = (*DppHwPci)(nil)

// Parses input into a DppHwPci struct.
func ParseDppHwPci(h EfiDevicePathProtocolHdr, b []byte) (*DppHwPci, error) {
	if len(b) != 2 {
		return nil, ErrParse
	}
	return &DppHwPci{
		Hdr:      h,
		Function: b[0],
		Device:   b[1],
	}, nil
}

func (e *DppHwPci) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppHwPci) ProtoSubTypeStr() string {
	return EfiDppHwSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppHwPci) String() string {
	return fmt.Sprintf("PCI(0x%x,0x%x)", e.Function, e.Device)
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppHwPci) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}
