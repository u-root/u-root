// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

const (
	e820max   = 32
	paramBase = 0x90000
	paramSize = 0x1000
)

type e820type uint32

const (
	ram e820type = 1
	reserved
	acpi
	nvs
)

type e820 struct {
	addr uint64
	size uint64
	typ  e820type
}

///* The header of Linux/i386 kernel */
//struct linux_header {
//	u8 reserved1[0x1f1];	/* 0x000 */
//	u8 setup_sects;		/* 0x1f1 */
//	u16 root_flags;		/* 0x1f2 */
//	u32 syssize;		/* 0x1f4 (2.04+) */
//	u8 reserved2[2];	/* 0x1f8 */
//	u16 vid_mode;		/* 0x1fa */
//	u16 root_dev;		/* 0x1fc */
//	u16 boot_sector_magic;	/* 0x1fe */
//	/* 2.00+ */
//	u8 reserved3[2];	/* 0x200 */
//	u8 header_magic[4];	/* 0x202 */
//	u16 protocol_version;	/* 0x206 */
//	u32 realmode_swtch;	/* 0x208 */
//	u16 start_sys;		/* 0x20c */
//	u16 kver_addr;		/* 0x20e */
//	u8 type_of_loader;	/* 0x210 */
//	u8 loadflags;		/* 0x211 */
//	u16 setup_move_size;	/* 0x212 */
//	u32 code32_start;	/* 0x214 */
//	u32 ramdisk_image;	/* 0x218 */
//	u32 ramdisk_size;	/* 0x21c */
//	u8 reserved4[4];	/* 0x220 */
//	/* 2.01+ */
//	u16 heap_end_ptr;	/* 0x224 */
//	u8 reserved5[2];	/* 0x226 */
//	/* 2.02+ */
//	u32 cmd_line_ptr;	/* 0x228 */
//	/* 2.03+ */
//	u32 initrd_addr_max;	/* 0x22c */
//	/* 2.05+ */
//	u32 kernel_alignment;	/* 0x230 */
//	u8 relocatable_kernel;	/* 0x234 */
//	u8 min_alignment;	/* 0x235 (2.10+) */
//	u8 reserved6[2];	/* 0x236 */
//	/* 2.06+ */
//	u32 cmdline_size;	/* 0x238 */
//	/* 2.07+ */
//	u32 hardware_subarch;	/* 0x23c */
//	u64 hardware_subarch_data;// 0x240 */
//	// 2.08+ */
//	u32 payload_offset;	// 0x248 */
//	u32 payload_length;	// 0x24c */
//	// 2.09+ */
//	u64 setup_data;		// 0x250 */
//	// 2.10+ */
//	u64 pref_address;	// 0x258 */
//	u32 init_size;		// 0x260 */
//} __packed;

const (
	CLMagic = 0xA33F
)

type ramDiskFlags uint16

const (
	startMask ramDiskFlags = 0x7ff
	prompt                 = 0x8000
	load                   = 0x4000
)

type loaderType uint8

const (
	loadLin   loaderType = 1
	bootSect             = 2
	sysLinux             = 3
	etherBoot            = 4
	kernel               = 5
)

const commandLineSize = 256

// BootParams are passed via kexec to the kernel.
// They are place at the mis-named"zero page",
// which is at 0x90000.
// What we're doing here is kinda hokey. But it's just proven worth
// it to try to line up struct members with the real thing, even if
// it gets a bit over the top (see the e820 bits below).
// Most of the members here are not exported; most we don't care about.
type LinuxBootParams struct {
	OrigX         uint8  `offset:"0x00"`
	OrigY         uint8  `offset:"0x01"`
	ExtMemK       uint16 // 0x02
	OrigVideoPage uint16 `offset:"0x04"`
	OrigVideoMode uint8  `offset:"0x06"`
	OrigVideoCols uint8  `offset:"0x07"`
	//_               uint16 `offset:"0x08"`
	OrigVideoEGAbx uint16 `offset:"0x0a"`
	//_               uint16 `offset:"0x0c"`
	OrigVideoLines  uint8  `offset:"0x0e"`
	OrigVideoIsVGA  uint8  `offset:"0x0f"`
	OrigVideoPoints uint16 `offset:"0x10"`

	// VESA graphic mode -- linear frame buffer
	LfbWidth  uint16 `offset:"0x12"`
	LfbHeight uint16 `offset:"0x14"`
	LfbDepth  uint16 `offset:"0x16"`
	LfbBase   uint32 `offset:"0x18"`
	LfbSize   uint32 `offset:"0x1c"`
	ClMagic   uint16 `offset:"0x20"`

	ClOffset      uint16 `offset:"0x22"`
	LfbLineLength uint16 `offset:"0x24"`
	RedSize       uint8  `offset:"0x26"`
	RedPos        uint8  `offset:"0x27"`
	GreenSize     uint8  `offset:"0x28"`
	GreenPos      uint8  `offset:"0x29"`
	BlueSize      uint8  `offset:"0x2a"`
	BluePos       uint8  `offset:"0x2b"`
	RsvdSize      uint8  `offset:"0x2c"`
	RsvdPos       uint8  `offset:"0x2d"`
	VesapmSeg     uint16 `offset:"0x2e"`
	VesapmOff     uint16 `offset:"0x30"`
	Pages         uint16 `offset:"0x32"`
	//_             [12]uint8 `offset:"0x34"` //-- 0x3f reserved for future expansion

	//struct apm_bios_info apm_bios_info;   `offset:"0x40"`
	//_ [0x40]uint8 // obsolete apm bios info
	//struct drive_info_struct drive_info;  `offset:"0x80"`
	//_ [0x20]uint8 // obsolete drive info
	//struct sys_desc_table sys_desc_table; `offset:"0xa0"`
	//_                   [0x140]uint8 // obsolete sys_desc_table
	AltMemK uint32 `offset:"0x1e0"`
	//_                   [4]uint8     `offset:"0x1e4"`
	E820MapNr uint8 `offset:"0x1e8"`
	//_                   [8]uint8     `offset:"0x1e9"`
	SetupHdr        uint8  `offset:"0x1f1"`
	MountRootRdonly uint16 `offset:"0x1f2"`
	//_                   [4]uint8     `offset:"0x1f4"`
	RamDiskFlags ramDiskFlags `offset:"0x1f8"`
	//_                   [2]uint8     `offset:"0x1fa"`
	OrigRootDev uint16 `offset:"0x1fc"`
	//_                   [1]uint8     `offset:"0x1fe"`
	AuxDeviceInfo uint8 `offset:"0x1ff"`
	//_                   [2]uint8     `offset:"0x200"`
	ParamBlockSignature [4]uint8 `offset:"0x202"`
	ParamBlockVersion   uint16   `offset:"0x206"`
	//_                   [8]uint8     `offset:"0x208"`
	LoaderType  loaderType `offset:"0x210"`
	LoaderFlags uint8      `offset:"0x211"`
	//_                   [2]uint8     `offset:"0x212"`
	KernelStart uint32 `offset:"0x214"`
	InitrdStart uint32 `offset:"0x218"`
	InitrdSize  uint32 `offset:"0x21c"`
	//_                   [8]uint8     `offset:"0x220"`
	CmdLinePtr        uint32 `offset:"0x228"`
	InitrdAddrMax     uint32 `offset:"0x22c"`
	KernelAlignment   uint32 `offset:"0x230"`
	RelocatableKernel uint8  `offset:"0x234"`
	//_                   [0x2b]uint8  `offset:"0x235"`
	InitSize uint32 `offset:"0x260"`
	//_                   [0x6c]uint8  `offset:"0x264"`
	// This fails as Go has to pad the struct since it contains
	// mixed 32 and 64 bit.
	E820Map [e820max]e820 `offset:"0x2d0"`
	//_       [688]uint8             `offset:"0x550"`
	CmdLine [commandLineSize]uint8 `offset:"0x800"`
	//_       [1792]uint8            `offset:"0x900"` // - 0x1000
}

// NewLinuxBootParams returns a LinuxBootParams with the bare minimum filled in.
// We have learned from coreboot with Linux payloads just how minimal that can
// be, and we'll start there.
func NewLinuxBootParams() *LinuxBootParams {
	var b = &LinuxBootParams{
		OrigVideoMode:   3,
		OrigVideoCols:   80,
		OrigVideoLines:  25,
		OrigVideoIsVGA:  1,
		OrigVideoPoints: 16,
		MountRootRdonly: 1,
		//OrigRootDev: xx,
		//InitSize: notsure,
		LoaderType: 0xff,
		CmdLinePtr: 0x90800,
	}
	return b
}

func (bp *LinuxBootParams) Marshal() ([]byte, error) {
	var b = bytes.NewBuffer(make([]byte, 0, 0x1000))
	rv := reflect.ValueOf(*bp)
	rt := reflect.TypeOf(*bp)
	for ix := 0; ix < rv.NumField(); ix++ {
		f := rt.Field(ix)
		n := f.Tag.Get("offset")
		if n == "" {
			continue
		}
		o, err := strconv.ParseUint(n, 0, 16)
		if err != nil {
			return nil, fmt.Errorf("Bad offset tag for %v: %v: %v", f.Name, n, err)
		}
		if b.Len() < int(o) {
			b.Write(make([]byte, int(o)-b.Len()))
		}
		if err := binary.Write(b, binary.LittleEndian, rv.Field(ix).Interface()); err != nil {
			return nil, fmt.Errorf("Writing %v: got %v, want nil", f.Name, err)
		}

	}
	b.Write(make([]byte, 0x1000-b.Len()))
	return b.Bytes(), nil
}
func (bp *LinuxBootParams) Segment() ([]Segment, error) {
	b, err := bp.Marshal()
	if err != nil {
		return nil, err
	}
	return []Segment{NewSegment(b, Range{Start: 0x90000, Size: uint(len(b))})}, nil
}
