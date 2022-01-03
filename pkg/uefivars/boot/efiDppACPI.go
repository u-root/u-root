// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"bytes"
	"fmt"
)

// EfiDppACPISubType is the dpp subtype for ACPI.
type EfiDppACPISubType EfiDevPathProtoSubType

const (
	DppAcpiTypeDevPath EfiDppACPISubType = iota + 1
	DppAcpiTypeExpandedDevPath
	DppAcpiTypeADR
	DppAcpiTypeNVDIMM
)

var efiDppACPISubTypeStrings = map[EfiDppACPISubType]string{
	DppAcpiTypeDevPath:         "Device Path",
	DppAcpiTypeExpandedDevPath: "Expanded Device Path",
	DppAcpiTypeADR:             "_ADR",
	DppAcpiTypeNVDIMM:          "NVDIMM",
}

func (e EfiDppACPISubType) String() string {
	if s, ok := efiDppACPISubTypeStrings[e]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN-0x%x", uint8(e))
}

// DppAcpiDevPath is an acpi device path.
type DppAcpiDevPath struct {
	Hdr      EfiDevicePathProtocolHdr
	HID, UID []byte // both length 4; not sure of endianness
}

var _ EfiDevicePathProtocol = (*DppAcpiDevPath)(nil)

// ParseDppAcpiDevPath parses input into a DppAcpiDevPath.
func ParseDppAcpiDevPath(h EfiDevicePathProtocolHdr, b []byte) (*DppAcpiDevPath, error) {
	if h.Length != 12 {
		return nil, ErrParse
	}
	return &DppAcpiDevPath{
		Hdr: h,
		HID: b[:4],
		UID: b[4:8],
	}, nil
}

func (e *DppAcpiDevPath) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppAcpiDevPath) ProtoSubTypeStr() string {
	return EfiDppACPISubType(e.Hdr.ProtoSubType).String()
}

func (e *DppAcpiDevPath) String() string { return fmt.Sprintf("ACPI(0x%x,0x%x)", e.HID, e.UID) }

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppAcpiDevPath) Resolver() (EfiPathSegmentResolver, error) { return nil, ErrUnimpl }

// DppAcpiExDevPath is an expanded dpp acpi device path.
type DppAcpiExDevPath struct {
	Hdr                    EfiDevicePathProtocolHdr
	HID, UID, CID          []byte // all length 4; not sure of endianness
	HIDSTR, UIDSTR, CIDSTR string
}

var _ EfiDevicePathProtocol = (*DppAcpiExDevPath)(nil)

// ParseDppAcpiExDevPath parses input into a DppAcpiExDevPath.
func ParseDppAcpiExDevPath(h EfiDevicePathProtocolHdr, b []byte) (*DppAcpiExDevPath, error) {
	if h.Length < 19 {
		return nil, ErrParse
	}
	ex := &DppAcpiExDevPath{
		Hdr: h,
		HID: b[:4],
		UID: b[4:8],
		CID: b[8:12],
	}
	b = b[12:]
	var err error
	ex.HIDSTR, err = readToNull(b)
	if err != nil {
		return nil, err
	}
	b = b[len(ex.HIDSTR)+1:]
	ex.UIDSTR, err = readToNull(b)
	if err != nil {
		return nil, err
	}
	b = b[len(ex.UIDSTR)+1:]
	ex.CIDSTR, err = readToNull(b)
	if err != nil {
		return nil, err
	}
	return ex, nil
}

func (e *DppAcpiExDevPath) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppAcpiExDevPath) ProtoSubTypeStr() string {
	return EfiDppACPISubType(e.Hdr.ProtoSubType).String()
}

func (e *DppAcpiExDevPath) String() string {
	return fmt.Sprintf("ACPI_EX(0x%x,0x%x,0x%x,%s,%s,%s)", e.HID, e.UID, e.CID, e.HIDSTR, e.UIDSTR, e.CIDSTR)
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppAcpiExDevPath) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}

func readToNull(b []byte) (string, error) {
	i := bytes.IndexRune(b, 0)
	if i < 0 {
		return "", ErrParse
	}
	return string(b[:i]), nil
}
