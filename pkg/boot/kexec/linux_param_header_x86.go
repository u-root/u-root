package kexec

import (
	"github.com/u-root/u-root/pkg/boot/bzimage"
)

const (
	CL_MAGIC_VALUE = 0xA33F

	RAMDISK_IMAGE_START_MASK = 0x07FF
	RAMDISK_PROMPT_FLAG      = 0x8000
	RAMDISK_LOAD_FLAG        = 0x4000

	LOADER_TYPE_LOADLIN         = 1
	LOADER_TYPE_BOOTSECT_LOADER = 2
	LOADER_TYPE_SYSLINUX        = 3
	LOADER_TYPE_ETHERBOOT       = 4
	LOADER_TYPE_KEXEC           = 0x0D
	LOADER_TYPE_UNKNOWN         = 0xFF

	COMMAND_LINE_SIZE = 2048
)

/*
 * EDD stuff
 */
const (
	EDD_MBR_SIG_MAX = 16
	EDDMAXNR        = 6 /* number of EddInfo structs starting at EDDBUF  */

	EDD_EXT_FIXED_DISK_ACCESS           = (1 << 0)
	EDD_EXT_DEVICE_LOCKING_AND_EJECTING = (1 << 1)
	EDD_EXT_ENHANCED_DISK_DRIVE_SUPPORT = (1 << 2)
	EDD_EXT_64BIT_EXTENSIONS            = (1 << 3)

	EDD_DEVICE_PARAM_SIZE = 74
)

type EddInfo struct {
	device                   uint8
	version                  uint8
	interface_support        uint16
	legacy_max_cylinder      uint16
	legacy_max_head          uint8
	legacy_sectors_per_track uint8
	edd_device_params        [EDD_DEVICE_PARAM_SIZE]uint8
}

/*
 * e820 ram stuff.
 */
const E820MAX = 128 /* number of entries in E820MAP */

const (
	E820_RAM      = 1
	E820_RESERVED = 2
	E820_ACPI     = 3 /* usable as RAM once ACPI tables have been read */
	E820_NVS      = 4
	E820_PMEM     = 7
	E820_PRAM     = 12
)

type E820entry struct {
	Addr uint64 /* start of memory segment */
	Size uint64 /* size of memory segment */
	Typ  uint32 /* type of memory segment */
}

type ApmBIOSInfo struct {
	Version   uint16 `offset:"0x40"`
	Cseg      uint16 `offset:"0x42"`
	Offset    uint32 `offset:"0x44"`
	Cseg16    uint16 `offset:"0x48"`
	Dseg      uint16 `offset:"0x4a"`
	Flags     uint16 `offset:"0x4c"`
	CsegLen   uint16 `offset:"0x4e"`
	Cseg16Len uint16 `offset:"0x50"`
	DsegLen   uint16 `offset:"0x52"`
}

type DriveInfo struct {
	Dummy [32]uint8
}

type SysDescTable struct {
	Length uint16
	Table  [30]uint8
}

// getLinuxParamHeader returns a linux param header from the given bzimage.
//
// Current max offset to consder is 0x7f.
func getLinuxParamHeader(b *bzimage.BzImage) *bzimage.LinuxParams {
	var h = &bzimage.LinuxParams{}
	if b == nil {
		return h
	}
	// h.SetupSects = b.Header.SetupSects
	// h.MountRootRdonly = b.Header.RootFlags

	// // (2.04+)
	// // b.Header.Syssize uint32: [0x1f4, 0x1f8)
	// // h.syssize uint16: [0x1f4, 0x1f6)
	// // h.swapdev uint16: [0x1f6, 0x1f8)
	// h.Syssize = uint16(b.Header.Syssize & uint32(^uint16(0))) // 0x1f4.
	// h.Swapdev = uint16(b.Header.Syssize >> 32)                // 0x1f6.

	// h.RamdiskFlags = b.Header.RAMSize // 0x1f8.
	// h.VidMode = b.Header.Vidmode      // 0x1fa.
	// h.RootDev = b.Header.RootDev      // 0x1fc.

	// // TODO(10000TB): 0x1fe pointed to reserved9 in x86 linux params
	// // header. Is it right when we simply copy over lower uint8 bits
	// // of Bootsectormagic, an uint16 from bzimage?
	// h.Reserved9[0][0] = bytesFromUint16(b.Header.Bootsectormagic)[0] // 0x1fe.
	// h.AuxDeviceInfo = bytesFromUint16(b.Header.Bootsectormagic)[1]   // 0x1ff.

	// // (2.00+)
	// // Reserved10, [2]uint8, offset at 0x200 corresponds to Jump, uint16
	// // from bzimage's linux header.
	// h.Reserved10[0] = lowerByteFromUint16(b.Header.Jump)  // 0x200.
	// h.Reserved10[1] = higherByteFromUint16(b.Header.Jump) // 0x201.
	// h.HeaderMagic = b.Header.HeaderMagic                  // 0x202.

	// h.ProtocolVersion = b.Header.Protocolversion                     // 0x206.
	// h.RmodeSwitchIP = lowerUint16FromUint32(b.Header.RealModeSwitch) // 0x208.
	// h.RmodeSwitchCs = lowerUint16FromUint32(b.Header.RealModeSwitch) // 0x20a.

	// h.Reserved11[0] = lowerByteFromUint16(b.Header.StartSys)  // 0x20c.
	// h.Reserved11[1] = higherByteFromUint16(b.Header.StartSys) // 0x20d.
	// h.Reserved11[2] = lowerByteFromUint16(b.Header.Kveraddr)  // 0x20e.
	// h.Reserved11[3] = higherByteFromUint16(b.Header.Kveraddr) // 0x20f.

	// h.LoaderType = b.Header.TypeOfLoader  // 0x210.
	// h.LoaderFlagzaAq = b.Header.Loadflags // 0x211.

	// h.Reserved12[0] = lowerByteFromUint16(b.Header.Setupmovesize)  // 0x212.
	// h.Reserved12[1] = higherByteFromUint16(b.Header.Setupmovesize) // 0x213.

	// h.KernelStart = b.Header.Code32Start  // 0x214.
	// h.InitrdStart = b.Header.RamdiskImage // 0x218.
	// h.InitrdSize = b.Header.RamdiskSize   // 0x21c.

	// h.Reserved13 = b.Header.BootSectKludge // 0x220.

	// // (2.01+)
	// h.HeapEndPtr = b.Header.Heapendptr // 0x224.

	// h.Reserved14[0] = b.Header.ExtLoaderVer  // 0x226.
	// h.Reserved14[1] = b.Header.ExtLoaderType // 0x227.

	// // (2.02+)
	// h.CmdLinePtr = b.Header.Cmdlineptr // 0x228.

	// // (2.03+)
	// h.InitrdAddrMax = b.Header.InitrdAddrMax // 0x22c.

	// // (2.05+)
	// h.KernelAlignment = b.Header.Kernelalignment     // 0x230.
	// h.RelocatableKernel = b.Header.RelocatableKernel // 0x234.
	// h.MinAlignment = b.Header.MinAlignment           // 0x235 (2.10+).
	// h.XloadFlags = b.Header.XLoadFlags               // 0x236.
	// // (2.06+)
	// h.CmdlineSize = b.Header.CmdLineSize // 0x238.
	// // (2.07+)
	// h.HardwareSubarch = b.Header.HardwareSubArch         // 0x23c.
	// h.HardwareSubarchData = b.Header.HardwareSubArchData // 0x240.
	// // (2.08+)
	// h.PayloadOffset = b.Header.PayloadOffset // 0x248.
	// h.PayloadLength = b.Header.PayloadSize   // 0x24c.
	// // (2.09+)
	// h.SetupData = b.Header.SetupData // 0x250.
	// // (2.10+)
	// h.PrefAddress = b.Header.PrefAddress       // 0x258.
	// h.InitSize = b.Header.InitSize             // 0x260.
	// h.HandoverOffset = b.Header.HandoverOffset // 0x264.
	return h
}
