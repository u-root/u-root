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
	"log"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/uefivars"
)

type EfiDppMediaSubType EfiDevPathProtoSubType

const (
	// DppTypeMedia, pg 319 +
	DppMTypeHdd      EfiDppMediaSubType = iota + 1 // 0x01
	DppMTypeCd                                     // 0x02
	DppMTypeVendor                                 // 0x03
	DppMTypeFilePath                               // 0x04 //p321
	DppMTypeMedia                                  // 0x05 //media protocol i.e. filesystem format??
	DppMTypePIWGFF                                 // 0x06
	DppMTypePIWGFV                                 // 0x07
	DppMTypeRelOff                                 // 0x08
	DppMTypeRAM                                    // 0x09
)

var efiDppMediaSubTypeStrings = map[EfiDppMediaSubType]string{
	DppMTypeHdd:      "HDD",
	DppMTypeCd:       "CD",
	DppMTypeVendor:   "Vendor",
	DppMTypeFilePath: "FilePath",
	DppMTypeMedia:    "Media",
	DppMTypePIWGFF:   "PIWG Firmware File",
	DppMTypePIWGFV:   "PIWG Firmware Volume",
	DppMTypeRelOff:   "Relative Offset",
	DppMTypeRAM:      "RAMDisk",
}

func (e EfiDppMediaSubType) String() string {
	if s, ok := efiDppMediaSubTypeStrings[e]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN-0x%x", uint8(e))
}

// DppMediaHDD is the struct in EfiDevicePathProtocol for DppMTypeHdd
type DppMediaHDD struct {
	Hdr EfiDevicePathProtocolHdr

	PartNum   uint32             // index into partition table for MBR or GPT; 0 indicates entire disk
	PartStart uint64             // starting LBA. only used for MBR?
	PartSize  uint64             // size in LB's. only used for MBR?
	PartSig   uefivars.MixedGUID // format determined by SigType below. unused bytes must be 0x0.
	PartFmt   uint8              // 0x01 for MBR, 0x02 for GPT
	SigType   uint8              // 0x00 - none; 0x01 - 32bit MBR sig (@ 0x1b8); 0x02 - GUID
}

var _ EfiDevicePathProtocol = (*DppMediaHDD)(nil)

// ParseDppMediaHdd parses input into a DppMediaHDD struct.
func ParseDppMediaHdd(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaHDD, error) {
	if len(b) < 38 {
		return nil, ErrParse
	}
	hdd := &DppMediaHDD{
		Hdr:       h,
		PartNum:   binary.LittleEndian.Uint32(b[:4]),
		PartStart: binary.LittleEndian.Uint64(b[4:12]),
		PartSize:  binary.LittleEndian.Uint64(b[12:20]),
		// PartSig:   b[20:36], //cannot assign slice to array
		PartFmt: b[36],
		SigType: b[37],
	}
	copy(hdd.PartSig[:], b[20:36])
	return hdd, nil
}

func (e *DppMediaHDD) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMediaHDD) ProtoSubTypeStr() string {
	return EfiDppMediaSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMediaHDD) String() string {
	//             (part#,pttype,guid,begin,length)
	return fmt.Sprintf("HD(%d,%s,%s,0x%x,0x%x)", e.PartNum, e.pttype(), e.sig(), e.PartStart, e.PartSize)
}

// Resolver returns an EfiPathSegmentResolver which can find and mount the
// partition described by DppMediaHDD.
func (e *DppMediaHDD) Resolver() (EfiPathSegmentResolver, error) {
	allBlocks, err := block.GetBlockDevices()
	if err != nil {
		return nil, err
	}
	var blocks block.BlockDevices
	if e.SigType == 2 {
		guid := e.PartSig.ToStdEnc().String()
		blocks = allBlocks.FilterPartID(guid)
	} else if e.SigType == 1 {
		// use MBR ID (e.PartSig[:4]) and partition number (e.PartNum)
		// see PartSig comments and UEFI documentation for location of MBR ID
		log.Printf("Sig Type 1: unimplemented/cannot identify")
		return nil, ErrNotFound
	} else {
		// SigType==0: no sig, would need to compare partition #/start/len
		log.Printf("Sig Type %d: unimplemented/cannot identify", e.SigType)
		return nil, ErrNotFound
	}
	if len(blocks) != 1 {
		log.Printf("blocks: %#v", blocks)
		return nil, ErrNotFound
	}
	return &HddResolver{BlockDev: blocks[0]}, nil
}

// return the partition table type as a string
func (e *DppMediaHDD) pttype() string {
	switch e.PartFmt {
	case 1:
		return "MBR"
	case 2:
		return "GPT"
	default:
		return "UNKNOWN"
	}
}

// return the signature as a string
func (e *DppMediaHDD) sig() string {
	switch e.SigType {
	case 1: // 32-bit MBR sig
		return fmt.Sprintf("%x", binary.LittleEndian.Uint32(e.PartSig[:4]))
	case 2: // GUID
		return e.PartSig.ToStdEnc().String()
	default:
		return "(NO SIG)"
	}
}

// DppMediaFilePath is a struct in EfiDevicePathProtocol for DppMTypeFilePath.
//
// If multiple are included in a load option, the docs say to concatenate them.
type DppMediaFilePath struct {
	Hdr EfiDevicePathProtocolHdr

	PathNameDecoded string // stored as utf16
}

var _ EfiDevicePathProtocol = (*DppMediaFilePath)(nil)

func ParseDppMediaFilePath(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaFilePath, error) {
	if len(b) < int(h.Length)-4 {
		return nil, ErrParse
	}
	path, err := uefivars.DecodeUTF16(b[:h.Length-4])
	if err != nil {
		return nil, err
	}
	// remove null termination byte, replace windows slashes
	path = strings.TrimSuffix(path, "\000")
	path = strings.Replace(path, "\\", string(os.PathSeparator), -1)
	fp := &DppMediaFilePath{
		Hdr:             h,
		PathNameDecoded: path,
	}
	return fp, nil
}

func (e *DppMediaFilePath) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMediaFilePath) ProtoSubTypeStr() string {
	return EfiDppMediaSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMediaFilePath) String() string {
	return fmt.Sprintf("File(%s)", e.PathNameDecoded)
}

// Resolver returns an EfiPathSegmentResolver decoding the DppMediaFilePath.
func (e *DppMediaFilePath) Resolver() (EfiPathSegmentResolver, error) {
	fp.Clean(e.PathNameDecoded)
	pr := PathResolver(e.PathNameDecoded)
	return &pr, nil
}

// struct in EfiDevicePathProtocol for DppMTypePIWGFV
type DppMediaPIWGFV struct {
	Hdr EfiDevicePathProtocolHdr
	Fv  []byte
}

var _ EfiDevicePathProtocol = (*DppMediaPIWGFV)(nil)

// ParseDppMediaPIWGFV parses input into a DppMediaPIWGFV.
func ParseDppMediaPIWGFV(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaPIWGFV, error) {
	if h.Length != 20 {
		return nil, ErrParse
	}
	fv := &DppMediaPIWGFV{
		Hdr: h,
		Fv:  b,
	}
	return fv, nil
}
func (e *DppMediaPIWGFV) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMediaPIWGFV) ProtoSubTypeStr() string {
	return EfiDppMediaSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMediaPIWGFV) String() string {
	var g uefivars.MixedGUID
	copy(g[:], e.Fv)
	return fmt.Sprintf("Fv(%s)", g.ToStdEnc().String())
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppMediaPIWGFV) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}

// struct in EfiDevicePathProtocol for DppMTypePIWGFF
type DppMediaPIWGFF struct {
	Hdr EfiDevicePathProtocolHdr
	Ff  []byte
}

var _ EfiDevicePathProtocol = (*DppMediaPIWGFF)(nil)

// ParseDppMediaPIWGFF parses the input into a DppMediaPIWGFF.
func ParseDppMediaPIWGFF(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaPIWGFF, error) {
	if h.Length != 20 {
		return nil, ErrParse
	}
	fv := &DppMediaPIWGFF{
		Hdr: h,
		Ff:  b,
	}
	return fv, nil
}

func (e *DppMediaPIWGFF) Header() EfiDevicePathProtocolHdr { return e.Hdr }

// ProtoSubTypeStr returns the subtype as human readable.
func (e *DppMediaPIWGFF) ProtoSubTypeStr() string {
	return EfiDppMediaSubType(e.Hdr.ProtoSubType).String()
}

func (e *DppMediaPIWGFF) String() string {
	var g uefivars.MixedGUID
	copy(g[:], e.Ff)
	return fmt.Sprintf("FvFile(%s)", g.ToStdEnc().String())
}

// Resolver returns a nil EfiPathSegmentResolver and ErrUnimpl. See the comment
// associated with ErrUnimpl.
func (e *DppMediaPIWGFF) Resolver() (EfiPathSegmentResolver, error) {
	return nil, ErrUnimpl
}
