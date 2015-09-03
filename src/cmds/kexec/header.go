// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Thanks to coreboot for documenting the basic layout.
package main

const (
	Ram e820type = 1
	Reserved = 2
	ACPI = 3
	NVS  = 4
)

const (
	E820MAX = 32 // number of entries in E820MAP
)

const (
	NotSet boottype = 0
	LoadLin = 1
	BootSect = 2
	SysLinux = 3
	EtherBoot = 4
	Kernel = 5
)

const (
	RamdiskStartMask = 0x07FF
	Prompt           = 0x8000
	Load             = 0x4000

	CommandLineMagic = 0x7ff
	CommandLineSize  = 256

	HeaderMagic = 0x53726448 // "HdrS" but little-endian constant
)

type e820type uint8
type boottype uint8
type e820entry struct {
	Addr    uint64
	Size    uint64
	MemType e820type
}

// The header of Linux/i386 kernel 
type LinuxHeader struct {
	_               [0x1f1]uint8 // 0x000 
	SetupSects      uint8        // 0x1f1 
	RootFlags       uint16       // 0x1f2 
	Syssize         uint32       // 0x1f4 (2.04+) 
	_               [2]uint8     // 0x1f8 
	Vidmode         uint16       // 0x1fa 
	RootDev         uint16       // 0x1fc 
	Bootsectormagic uint16       // 0x1fe 
	// 2.00+ 
	_               [2]uint8 // 0x200 
	HeaderMagic     uint32 // 0x202 
	Protocolversion uint16   // 0x206 
	RealModeSwitch   uint32   // 0x208 
	StartSys        uint16   // 0x20c 
	Kveraddr        uint16   // 0x20e 
	TypeOfLoader    uint8    // 0x210 
	Loadflags       uint8    // 0x211 
	Setupmovesize   uint16   // 0x212 
	Code32start     uint32   // 0x214 
	RamdiskImage    uint32   // 0x218 
	RamdiskSize     uint32   // 0x21c 
	_               [4]uint8 // 0x220 
	// 2.01+ 
	Heapendptr uint16   // 0x224 
	_          [2]uint8 // 0x226 
	// 2.02+ 
	Cmdlineptr uint32 // 0x228 
	// 2.03+ 
	Initrdaddrmax uint32 // 0x22c 
	// 2.05+ 
	Kernelalignment   uint32   // 0x230 
	Relocatablekernel uint8    // 0x234 
	Minalignment      uint8    // 0x235 (2.10+) 
	_                 [2]uint8 // 0x236 
	// 2.06+ 
	Cmdlinesize uint32 // 0x238 
	// 2.07+ 
	Hardwaresubarch     uint32 // 0x23c 
	Hardwaresubarchdata uint64 // 0x240 
	// 2.08+ 
	Payloadoffset uint32 // 0x248 
	Payloadlength uint32 // 0x24c 
	// 2.09+ 
	Setupdata uint64 // 0x250 
	// 2.10+ 
	Prefaddress uint64 // 0x258 
	Initsize    uint32 // 0x260 
}

// Paramters passed to 32-bit part of Linux
type LinuxParams struct {
	Origx           uint8  // 0x00 
	Origy           uint8  // 0x01 
	ExtMemK         uint16 // 0x02 -- EXTMEMK sits here 
	OrigVideoPage   uint16 // 0x04 
	OrigVideoMode   uint8  // 0x06 
	OrigVideoCols   uint8  // 0x07 
	_               uint16 // 0x08 
	OrigVideoeGabx  uint16 // 0x0a 
	_               uint16 // 0x0c 
	OrigVideoLines  uint8  // 0x0e 
	OrigVideoIsVGA  uint8  // 0x0f 
	OrigVideoPoints uint16 // 0x10 

	// VESA graphic mode -- linear frame buffer 
	Lfbwidth      uint16    // 0x12 
	Lfbheight     uint16    // 0x14 
	Lfbdepth      uint16    // 0x16 
	Lfbbase       uint32    // 0x18 
	Lfbsize       uint32    // 0x1c 
	CLMagic       uint16    // 0x20 
	CLOffset      uint16    // 0x22 
	Lfblinelength uint16    // 0x24 
	Redsize       uint8     // 0x26 
	Redpos        uint8     // 0x27 
	Greensize     uint8     // 0x28 
	Greenpos      uint8     // 0x29 
	Bluesize      uint8     // 0x2a 
	Bluepos       uint8     // 0x2b 
	Rsvdsize      uint8     // 0x2c 
	Rsvdpos       uint8     // 0x2d 
	Vesapmseg     uint16    // 0x2e 
	Vesapmoff     uint16    // 0x30 
	Pages         uint16    // 0x32 
	_             [12]uint8 // 0x34 -- 0x3f reserved for future expansion 

	//struct apmbiosinfo apmbiosinfo;   // 0x40 
	Apmbiosinfo [0x40]uint8
	//struct driveinfostruct driveinfo;  // 0x80 
	Driveinfo [0x20]uint8
	//struct sysdesctable sysdesctable; // 0xa0 
	Sysdesctable        [0x140]uint8
	Altmemk             uint32             // 0x1e0 
	_                   [4]uint8           // 0x1e4 
	E820mapnr           uint8              // 0x1e8 
	_                   [9]uint8           // 0x1e9 
	MountRootReadonly     uint16             // 0x1f2 
	_                   [4]uint8           // 0x1f4 
	Ramdiskflags        uint16             // 0x1f8 
	_                   [2]uint8           // 0x1fa 
	OrigRootDev         uint16             // 0x1fc 
	_                   [1]uint8           // 0x1fe 
	Auxdeviceinfo       uint8              // 0x1ff 
	_                   [2]uint8           // 0x200 
	Paramblocksignature [4]uint8           // 0x202 
	Paramblockversion   uint16             // 0x206 
	_                   [8]uint8           // 0x208 
	LoaderType          uint8              // 0x210 
	Loaderflags         uint8              // 0x211 
	_                   [2]uint8           // 0x212 
	KernelStart         uint32             // 0x214 
	Initrdstart         uint32             // 0x218 
	Initrdsize          uint32             // 0x21c 
	_                   [8]uint8           // 0x220 
	Cmdlineptr          uint32             // 0x228 
	Initrdaddrmax       uint32             // 0x22c 
	Kernelalignment     uint32             // 0x230 
	Relocatablekernel   uint8              // 0x234 
	_                   [155]uint8         // 0x22c 
	E820Map             [E820MAX]e820entry // 0x2d0 
	_                   [688]uint8         // 0x550 
	Commandline [CommandLineSize]uint8 // 0x800 
	_           [1792]uint8            // 0x900 - 0x1000 
}

var (
	LoaderType = map[boottype] string {
		NotSet: "Not set",
		LoadLin: "loadlin",
		BootSect: "bootsector",
		SysLinux: "syslinux",
		EtherBoot: "etherboot",
		Kernel: "kernel (kexec)",
	}
	E820 = map[e820type] string {
		Ram: "Ram",
		Reserved: "Reserved",
		ACPI: "ACPI",
		NVS: "NVS",
	}

)

