// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

// These are the semi-documented things that define a bzImage
// Thanks to coreboot for documenting the basic layout.

// E820 types.
const (
	RAM      e820type = 1
	Reserved e820type = 2
	ACPI     e820type = 3
	NVS      e820type = 4
)

// Boot types.
const (
	NotSet    boottype = 0
	LoadLin   boottype = 1
	BootSect  boottype = 2
	SysLinux  boottype = 3
	EtherBoot boottype = 4
	Kernel    boottype = 5
)

//Offsets and magic values.
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

// EDD consts.
const (
	EDDMBRSigMax       = 16
	EDDMaxNR           = 6 /* number of edd_info structs starting at EDDBUF  */
	EDDDeviceParamSize = 74
)

// EDDExt consts.
const (
	EDDExtFixedDiskAccess = 1 << iota
	EDDExtDeviceLockingAndEjecting
	EDDExtEnhancedDiskDriveSupport
	EDDExt64BitExtensions
)

// EDDInfo struct.
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

// E820Entry is one e820 entry.
type E820Entry struct {
	Addr    uint64
	Size    uint64
	MemType e820type
}

// LinuxHeader is the header of Linux/i386 kernel
type LinuxHeader struct {
	MBRCode         [0xc0]uint8         `offset:"0x000"`
	ExtRamdiskImage uint32              `offset:"0xc0"`
	ExtRamdiskSize  uint32              `offset:"0xc4"`
	ExtCmdlinePtr   uint32              `offset:"0xc8"`
	O               [0x1f1 - 0xcc]uint8 `offset:"0xcc"`
	SetupSects      uint8               `offset:"0x1f1"`
	RootFlags       uint16              `offset:"0x1f2"`
	Syssize         uint32              `offset:"0x1f4"` //(2.04+)
	RAMSize         uint16              `offset:"0x1f8"`
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
	RamdiskImage    uint32   `offset:"0x218"`
	RamdiskSize     uint32   `offset:"0x21c"`
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

var (
	// LoaderType contains strings describing boot types.
	LoaderType = map[boottype]string{
		NotSet:    "Not set",
		LoadLin:   "loadlin",
		BootSect:  "bootsector",
		SysLinux:  "syslinux",
		EtherBoot: "etherboot",
		Kernel:    "kernel (kexec)",
	}
	//E820 contains strings describing e820types.
	E820 = map[e820type]string{
		RAM:      "RAM",
		Reserved: "Reserved",
		ACPI:     "ACPI",
		NVS:      "NVS",
	}
	// HeaderMagic is kernel header magic bytes.
	HeaderMagic = [4]uint8{'H', 'd', 'r', 'S'}
)

// BzImage represents sections extracted from a kernel.
type BzImage struct {
	Header       LinuxHeader
	BootCode     []byte
	HeadCode     []byte
	KernelCode   []byte
	TailCode     []byte
	KernelBase   uintptr
	KernelOffset uintptr
	compressed   []byte
	// Some operations don't need the decompressed code; this speeds them up significantly.
	NoDecompress bool
}
