// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

// These are the semi-documented things that define a bzImage
// Thanks to coreboot for documenting the basic layout.

const (
	Ram      e820type = 1
	Reserved e820type = 2
	ACPI     e820type = 3
	NVS      e820type = 4
)

const (
	NotSet    boottype = 0
	LoadLin   boottype = 1
	BootSect  boottype = 2
	SysLinux  boottype = 3
	EtherBoot boottype = 4
	Kernel    boottype = 5
)

const (
	RamdiskStartMask = 0x07FF
	Prompt           = 0x8000
	Load             = 0x4000

	CommandLineMagic = 0x7ff
	CommandLineSize  = 256

	DefaultInitrdAddrMax  = 0x37FFFFFF
	DefaultBzimageAddrMax = 0x37FFFFFF

	E820Max = 128
	E820Map = 0x2d0
	E820NR  = 0x1e8
)

// what's an EDD? No idea.
/*
 * EDD stuff
 */

const (
	EDDMBRSigMax       = 16
	EDDMaxNR           = 6 /* number of edd_info structs starting at EDDBUF  */
	EDDDeviceParamSize = 74
)

const (
	EDDExtFixedDiskAccess = 1 << iota
	EDDExtDeviceLockingAndEjecting
	EDDExtEnhancedDiskDriveSupport
	EDDExt64BitExtensions
)

type EDDInfo struct {
	Device                uint8
	Version               uint8
	InterfaceSupport      uint16
	LegacyMaxCylinder     uint16
	LegacyMaxHead         uint8
	LegacySectorsPerTrace uint8
	EDDDeviceParams       [EDDDeviceParamSize]uint8
}

type e820type uint32
type boottype uint8
type E820Entry struct {
	Addr    uint64
	Size    uint64
	MemType e820type
}

// The header of Linux/i386 kernel
type LinuxHeader struct {
	MBRCode         [0xc0]uint8         `offset:"0x000"`
	ExtRamdiskImage uint32              `offset:"0xc0"`
	ExtRamdiskSize  uint32              `offset:"0xc4"`
	ExtCmdlinePtr   uint32              `offset:"0xc8"`
	O               [0x1f1 - 0xcc]uint8 `offset:"0xcc"`
	SetupSects      uint8               `offset:"0x1f1"`
	RootFlags       uint16              `offset:"0x1f2"`
	Syssize         uint32              `offset:"0x1f4"` //(2.04+)
	RamSize         uint16              `offset:"0x1f8"`
	Vidmode         uint16              `offset:"0x1fa"`
	RootDev         uint16              `offset:"0x1fc"`
	Bootsectormagic uint16              `offset:"0x1fe"`
	// 0.00+
	Jump            uint16   `offset:"0x200"`
	HeaderMagic     [4]uint8 `offset:"0x202"`
	Protocolversion uint16   `offset:"0x206"`
	RealModeSwitch  uint32   `offset:"0x208"`
	StartSys        uint16   `offset:"0x20c"`
	Kveraddr        uint16   `offset:"0x20e"`
	TypeOfLoader    uint8    `offset:"0x210"`
	Loadflags       uint8    `offset:"0x211"`
	Setupmovesize   uint16   `offset:"0x212"`
	Code32Start     uint32   `offset:"0x214"`
	RamDiskImage    uint32   `offset:"0x218"`
	RamDiskSize     uint32   `offset:"0x21c"`
	BootSectKludge  [4]uint8 `offset:"0x220"`
	// 2.01+
	Heapendptr    uint16 `offset:"0x224"`
	ExtLoaderVer  uint8  `offset:"0x226"`
	ExtLoaderType uint8  `offset:"0x227"`
	// 2.02+
	Cmdlineptr uint32 `offset:"0x228"`
	// 2.03+
	InitrdAddrMax uint32 `offset:"0x22c"`
	// 2.05+
	Kernelalignment   uint32 `offset:"0x230"`
	RelocatableKernel uint8  `offset:"0x234"`
	MinAlignment      uint8  `offset:"0x235"` //(2.10+)
	XLoadFlags        uint16 `offset:"0x236"`
	// 2.06+
	CmdLineSize uint32 `offset:"0x238"`
	// 2.07+
	HardwareSubArch     uint32 `offset:"0x23c"`
	HardwareSubArchData uint64 `offset:"0x240"`
	// 2.08+
	PayloadOffset uint32 `offset:"0x248"`
	PayloadSize   uint32 `offset:"0x24c"`
	// 2.09+
	SetupData uint64 `offset:"0x250"`
	// 2.10+
	PrefAddress    uint64 `offset:"0x258"`
	InitSize       uint32 `offset:"0x260"`
	HandoverOffset uint32 `offset:"0x264"`
}

// Parameters passed to 32-bit part of Linux
type LinuxParams struct {
	Origx           uint8  `offset:"0x00"`
	Origy           uint8  `offset:"0x01"`
	ExtMemK         uint16 `offset:"0x02"` //-- EXTMEMK sits here
	OrigVideoPage   uint16 `offset:"0x04"`
	OrigVideoMode   uint8  `offset:"0x06"`
	OrigVideoCols   uint8  `offset:"0x07"`
	_               uint16 `offset:"0x08"`
	OrigVideoeGabx  uint16 `offset:"0x0a"`
	_               uint16 `offset:"0x0c"`
	OrigVideoLines  uint8  `offset:"0x0e"`
	OrigVideoIsVGA  uint8  `offset:"0x0f"`
	OrigVideoPoints uint16 `offset:"0x10"`

	// VESA graphic mode -- linear frame buffer
	Lfbwidth      uint16    `offset:"0x12"`
	Lfbheight     uint16    `offset:"0x14"`
	Lfbdepth      uint16    `offset:"0x16"`
	Lfbbase       uint32    `offset:"0x18"`
	Lfbsize       uint32    `offset:"0x1c"`
	CLMagic       uint16    `offset:"0x20"` // DON'T USE
	CLOffset      uint16    `offset:"0x22"` // DON'T USE
	Lfblinelength uint16    `offset:"0x24"`
	Redsize       uint8     `offset:"0x26"`
	Redpos        uint8     `offset:"0x27"`
	Greensize     uint8     `offset:"0x28"`
	Greenpos      uint8     `offset:"0x29"`
	Bluesize      uint8     `offset:"0x2a"`
	Bluepos       uint8     `offset:"0x2b"`
	Rsvdsize      uint8     `offset:"0x2c"`
	Rsvdpos       uint8     `offset:"0x2d"`
	Vesapmseg     uint16    `offset:"0x2e"`
	Vesapmoff     uint16    `offset:"0x30"`
	Pages         uint16    `offset:"0x32"`
	_             [12]uint8 `offset:"0x34"` //-- 0x3f reserved for future expansion

	//struct apmbiosinfo apmbiosinfo;
	Apmbiosinfo [0x40]uint8 `offset:"0x40"`
	//struct driveinfostruct driveinfo;
	Driveinfo [0x20]uint8 `offset:"0x80"`
	//struct sysdesctable sysdesctable;
	Sysdesctable        [0x140]uint8 `offset:"0xa0"`
	Altmemk             uint32       `offset:"0x1e0"`
	_                   [4]uint8     `offset:"0x1e4"`
	E820MapNr           uint8        `offset:"0x1e8"`
	_                   [9]uint8     `offset:"0x1e9"`
	MountRootReadonly   uint16       `offset:"0x1f2"`
	_                   [4]uint8     `offset:"0x1f4"`
	Ramdiskflags        uint16       `offset:"0x1f8"`
	_                   [2]uint8     `offset:"0x1fa"`
	OrigRootDev         uint16       `offset:"0x1fc"`
	_                   [1]uint8     `offset:"0x1fe"`
	Auxdeviceinfo       uint8        `offset:"0x1ff"`
	_                   [2]uint8     `offset:"0x200"`
	Paramblocksignature [4]uint8     `offset:"0x202"`
	Paramblockversion   uint16       `offset:"0x206"`
	_                   [8]uint8     `offset:"0x208"`
	LoaderType          uint8        `offset:"0x210"`
	Loaderflags         uint8        `offset:"0x211"`
	_                   [2]uint8     `offset:"0x212"`
	KernelStart         uint32       `offset:"0x214"`
	Initrdstart         uint32       `offset:"0x218"`
	Initrdsize          uint32       `offset:"0x21c"`
	_                   [8]uint8     `offset:"0x220"`
	CLPtr               uint32       `offset:"0x228"` // USE THIS.
	InitrdAddrMax       uint32       `offset:"0x22c"`
	/* 2.04+ */
	KernelAlignment     uint32               `offset:"0x230"`
	RelocatableKernel   uint8                `offset:"0x234"`
	MinAlignment        uint8                `offset:"0x235"`
	XLoadFlags          uint16               `offset:"0x236"`
	CmdLineSize         uint32               `offset:"0x238"`
	HardwareSubarch     uint32               `offset:"0x23C"`
	HardwareSubarchData uint64               `offset:"0x240"`
	Payload_Ofset       uint32               `offset:"0x248"`
	PayloadLength       uint32               `offset:"0x24C"`
	SetupData           uint64               `offset:"0x250"`
	PrefAddress         uint64               `offset:"0x258"`
	InitSize            uint32               `offset:"0x260"`
	HandoverOffset      uint32               `offset:"0x264"`
	_                   [0x290 - 0x268]uint8 `offset:"0x268"`
	EDDMBRSigBuffer     [EDDMBRSigMax]uint32 `offset:"0x290"`
	// e820map is another cockup from the usual suspects.
	// Go rounds the size to something reasonable. Oh well. No checking for you.
	// So the next two offsets are bogus, sorry about that.
	E820Map [E820Max]E820Entry `offset:"0x2d0"`
	// we lie.
	_      [48]uint8         `offset:"0xed0"` // `offset:"0xcd0"`
	EDDBuf [EDDMaxNR]EDDInfo `offset:"0xf00"` // `offset:"0xd00"`

}

var (
	LoaderType = map[boottype]string{
		NotSet:    "Not set",
		LoadLin:   "loadlin",
		BootSect:  "bootsector",
		SysLinux:  "syslinux",
		EtherBoot: "etherboot",
		Kernel:    "kernel (kexec)",
	}
	E820 = map[e820type]string{
		Ram:      "Ram",
		Reserved: "Reserved",
		ACPI:     "ACPI",
		NVS:      "NVS",
	}
	HeaderMagic = [4]uint8{'H', 'd', 'r', 'S'}
)

type BzImage struct {
	Header       LinuxHeader
	BootCode     []byte
	HeadCode     []byte
	KernelCode   []byte
	TailCode     []byte
	KernelBase   uintptr
	KernelOffset uintptr
	compressed   []byte
}
