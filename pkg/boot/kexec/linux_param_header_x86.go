package kexec

import (
	"unsafe"

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

// LinuxParamHeader is the x86 linux params header.
type LinuxParamHeader struct {
	OrigX           uint8  `offset:"0x00"`
	OrigY           uint8  `offset:"0x01"`
	ExtMemK         uint16 `offset:"0x02"` /* EXT_MEM_K sits here */
	OrigVideoPage   uint16 `offset:"0x04"`
	OrigVideoMode   uint8  `offset:"0x06"`
	OrigVideoCols   uint8  `offset:"0x07"`
	unused2         uint16 `offset:"0x08"`
	OrigVideoEgaBx  uint16 `offset:"0x0a"`
	unused3         uint16 `offset:"0x0c"`
	OrigVideoLines  uint8  `offset:"0x0e"`
	OrigVideoIsVGA  uint8  `offset:"0x0f"`
	OrigVideoPoints uint16 `offset:"0x10"`

	/* VESA graphic mode -- linear frame buffer */
	LfbWidth  uint16 `offset:"0x12"`
	LfbHeight uint16 `offset:"0x14"`
	LfbDepth  uint16 `offset:"0x16"`
	LfbBase   uint32 `offset:"0x18"`
	LfbSize   uint32 `offset:"0x1c"`
	ClMagic   uint16 `offset:"0x20"`

	/* CL_MAGIC_VALUE 0xA33F */
	ClOffset       uint16   `offset:"0x22"`
	LfbLinelength  uint16   `offset:"0x24"`
	RedSize        uint8    `offset:"0x26"`
	RedPos         uint8    `offset:"0x27"`
	GreenSize      uint8    `offset:"0x28"`
	GreenPos       uint8    `offset:"0x29"`
	BlueSize       uint8    `offset:"0x2a"`
	BluePos        uint8    `offset:"0x2b"`
	RsvdSize       uint8    `offset:"0x2c"`
	RsvdPos        uint8    `offset:"0x2d"`
	VesapmSeg      uint16   `offset:"0x2e"`
	VesapmOff      uint16   `offset:"0x30"`
	Pages          uint16   `offset:"0x32"`
	VesaAttributes uint16   `offset:"0x34"`
	Capabilities   uint32   `offset:"0x36"`
	ExtLfbBase     uint32   `offset:"0x3a"`
	reserved4      [2]uint8 `offset:"0x3e"` /* 0x3e -- 0x3f reserved for future expansion */

	ApmBIOSInfo         ApmBIOSInfo         `offset:"0x40"`
	reserved4_1         [28]uint8           `offset:"0x54"` /* 0x54 */
	ACPIRsdpAddr        uint64              `offset:"0x70"`
	reserved4_2         [8]uint8            `offset:"0x78"`
	DriveInfo           DriveInfo           `offset:"0x80"`
	SysDescTable        SysDescTable        `offset:"0xa0"`
	ExtRamdiskImage     uint32              `offset:"0xc0"`
	ExtRamdiskSize      uint32              `offset:"0xc4"`
	ExtCmdLinePtr       uint32              `offset:"0xc8"`
	reserved4_3         [0x1c0 - 0xcc]uint8 `offset:"0xe4"`
	EfiInfo             [32]uint8           `offset:"0x1c0"`
	AltMemK             uint32              `offset:"0x1e0"`
	reserved5           [4]uint8            `offset:"0x1e4"`
	E820MapNr           uint8               `offset:"0x1e8"`
	EddbufEntries       uint8               `offset:"0x1e9"`
	EddMBRSigBufEntries uint8               `offset:"0x1ea"`
	reserved6           [6]uint8            `offset:"0x1eb"`
	SetupSects          uint8               `offset:"0x1f1"`
	MountRootRdonly     uint16              `offset:"0x1f2"`
	Syssize             uint16              `offset:"0x1f4"`
	Swapdev             uint16              `offset:"0x1f6"`
	RamdiskFlags        uint16              `offset:"0x1f8"`
	// RAMDISK_IMAGE_START_MASK	0x07FF
	// RAMDISK_PROMPT_FLAG		0x8000
	// RAMDISK_LOAD_FLAG		0x4000
	VidMode       uint16      `offset:"0x1fa"`
	RootDev       uint16      `offset:"0x1fc"`
	Reserved9     [1][1]uint8 `offset:"0x1fe"`
	AuxDeviceInfo uint8       `offset:"0x1ff"`
	/* 2.00+ */
	Reserved10      [2]uint8 `offset:"0x200"`
	HeaderMagic     [4]uint8 `offset:"0x202"`
	ProtocolVersion uint16   `offset:"0x206"`
	RmodeSwitchIP   uint16   `offset:"0x208"`
	RmodeSwitchCs   uint16   `offset:"0x20a"`
	Reserved11      [4]uint8 `offset:"0x20c"`
	LoaderType      uint8    `offset:"0x210"`
	// LOADER_TYPE_LOADLIN         1
	// LOADER_TYPE_BOOTSECT_LOADER 2
	// LOADER_TYPE_SYSLINUX        3
	// LOADER_TYPE_ETHERBOOT       4
	// LOADER_TYPE_KEXEC           0x0D
	// LOADER_TYPE_UNKNOWN         0xFF
	LoaderFlagzaAq uint8    `offset:"0x211"`
	Reserved12     [2]uint8 `offset:"0x212"`
	KernelStart    uint32   `offset:"0x214"`
	InitrdStart    uint32   `offset:"0x218"`
	InitrdSize     uint32   `offset:"0x21c"`
	Reserved13     [4]uint8 `offset:"0x220"`
	/* 2.01+ */
	HeapEndPtr uint16   `offset:"0x224"`
	Reserved14 [2]uint8 `offset:"0x226"`
	/* 2.02+ */
	CmdLinePtr uint32 `offset:"0x228"`
	/* 2.03+ */
	InitrdAddrMax uint32 `offset:"0x22c"`

	// TENATIVE = 0
	//
	// Code that is tenatively correct but hasn't yet been officially accepted
	//
	/* 2.04+ */
	// uint16_t entry32_off;			/* 0x230 */
	// uint16_t internal_cmdline_off;		/* 0x232 */
	// uint32_t low_base;			/* 0x234 */
	// uint32_t low_memsz;			/* 0x238 */
	// uint32_t low_filesz;			/* 0x23c */
	// uint32_t real_base;			/* 0x240 */
	// uint32_t real_memsz;			/* 0x244 */
	// uint32_t real_filesz;			/* 0x248 */
	// uint32_t high_base;			/* 0x24C */
	// uint32_t high_memsz;			/* 0x250 */
	// uint32_t high_filesz;			/* 0x254 */
	// uint8_t  reserved15[0x2d0 - 0x258];	/* 0x258 */

	/* 2.04+ */
	KernelAlignment     uint32                  `offset:"0x230"`
	RelocatableKernel   uint8                   `offset:"0x234"`
	MinAlignment        uint8                   `offset:"0x235"`
	XloadFlags          uint16                  `offset:"0x236"`
	CmdlineSize         uint32                  `offset:"0x238"`
	HardwareSubarch     uint32                  `offset:"0x23C"`
	HardwareSubarchData uint64                  `offset:"0x240"`
	PayloadOffset       uint32                  `offset:"0x248"`
	PayloadLength       uint32                  `offset:"0x24C"`
	SetupData           uint64                  `offset:"0x250"`
	PrefAddress         uint64                  `offset:"0x258"`
	InitSize            uint32                  `offset:"0x260"`
	HandoverOffset      uint32                  `offset:"0x264"`
	reserved16          [0x290 - 0x268]uint8    `offset:"0x268"`
	EddMBRSigBuffer     [EDD_MBR_SIG_MAX]uint32 `offset:"0x290"`

	E820_map [E820MAX]E820entry `offset:"0x2d0"`
	_pad8    [48]uint8          `offset:"0xcd0"`
	Eddbuf   [EDDMAXNR]EddInfo  `offset:"0xd00"`
	/* 0xeec */
}

func lowerByteFromUint16(v uint16) byte {
	return uint8(v & uint16(^uint8(0)))
}

func higherByteFromUint16(v uint16) byte {
	return uint8(v >> 16)
}

// bytesFromUint16 retrieves the 2 bytes from a given uint16.
//
// Assume little endian for now.
// TODO(10000TB): check endianess and return accordingly.
func bytesFromUint16(v uint16) [2]byte {
	return [2]byte{lowerByteFromUint16(v), higherByteFromUint16(v)}
}

func lowerUint16FromUint32(v uint32) uint16 {
	return uint16(v & uint32(^uint16(0)))
}

func higherUint16FromUint32(v uint32) uint16 {
	return uint16(v >> 32)
}

// uint16sFromUint32 retrieves the 2 uint16 from a given uint32.
func uint16sFromUint32(v uint32) [2]uint16 {
	return [2]uint16{lowerUint16FromUint32(v), higherUint16FromUint32(v)}
}

// ToBytes serializes current struct to a bytes slice with field values put in their respective offsets.
func (lph *LinuxParamHeader) ToBytes() []byte {
	sz := unsafe.Sizeof(*lph)
	buf := make([]byte, sz)

	// TODO(10000TB): complete this func impl.

	// Basic impl to make it 32bit entry to work somwhow first.
	//
	// TODO(10000TB): better way to "reflect" and programatically do these.
	// glob ...
	buf[0x1f1] = lph.SetupSects
	vals := bytesFromUint16(lph.MountRootRdonly)
	buf[0x1f2], buf[0x1f2+1] = vals[0], vals[1]
	vals = bytesFromUint16(lph.Syssize)
	buf[0x1f4], buf[0x1f4+1] = vals[0], vals[1]
	vals = bytesFromUint16(lph.Swapdev)
	buf[0x1f6], buf[0x1f6+1] = vals[0], vals[1]
	vals = bytesFromUint16(lph.RamdiskFlags)
	buf[0x1f8], buf[0x1f8+1] = vals[0], vals[1]
	vals = bytesFromUint16(lph.VidMode)
	buf[0x1fa], buf[0x1fa+1] = vals[0], vals[1]
	vals = bytesFromUint16(lph.RootDev)
	buf[0x1fc], buf[0x1fc+1] = vals[0], vals[1]

	return buf
}

// TODO(10000TB): bzimage pkg defines a similar struct as LinuxParamHeader
// but there are differences evaluate both, and maybe merge them, and move
// logic here as part of bzimage.

// getLinuxParamHeader returns a linux param header from the given bzimage.
//
// Current max offset to consder is 0x7f.
func getLinuxParamHeader(b *bzimage.BzImage) *LinuxParamHeader {
	h := LinuxParamHeader{}
	if b == nil {
		return nil
	}
	h.SetupSects = b.Header.SetupSects
	h.MountRootRdonly = b.Header.RootFlags

	// (2.04+)
	// b.Header.Syssize uint32: [0x1f4, 0x1f8)
	// h.syssize uint16: [0x1f4, 0x1f6)
	// h.swapdev uint16: [0x1f6, 0x1f8)
	h.Syssize = uint16(b.Header.Syssize & uint32(^uint16(0))) // 0x1f4.
	h.Swapdev = uint16(b.Header.Syssize >> 32)                // 0x1f6.

	h.RamdiskFlags = b.Header.RAMSize // 0x1f8.
	h.VidMode = b.Header.Vidmode      // 0x1fa.
	h.RootDev = b.Header.RootDev      // 0x1fc.

	// TODO(10000TB): 0x1fe pointed to reserved9 in x86 linux params
	// header. Is it right when we simply copy over lower uint8 bits
	// of Bootsectormagic, an uint16 from bzimage?
	h.Reserved9[0][0] = bytesFromUint16(b.Header.Bootsectormagic)[0] // 0x1fe.
	h.AuxDeviceInfo = bytesFromUint16(b.Header.Bootsectormagic)[1]   // 0x1ff.

	// (2.00+)
	// Reserved10, [2]uint8, offset at 0x200 corresponds to Jump, uint16
	// from bzimage's linux header.
	h.Reserved10[0] = lowerByteFromUint16(b.Header.Jump)  // 0x200.
	h.Reserved10[1] = higherByteFromUint16(b.Header.Jump) // 0x201.
	h.HeaderMagic = b.Header.HeaderMagic                  // 0x202.

	h.ProtocolVersion = b.Header.Protocolversion                     // 0x206.
	h.RmodeSwitchIP = lowerUint16FromUint32(b.Header.RealModeSwitch) // 0x208.
	h.RmodeSwitchCs = lowerUint16FromUint32(b.Header.RealModeSwitch) // 0x20a.

	h.Reserved11[0] = lowerByteFromUint16(b.Header.StartSys)  // 0x20c.
	h.Reserved11[1] = higherByteFromUint16(b.Header.StartSys) // 0x20d.
	h.Reserved11[2] = lowerByteFromUint16(b.Header.Kveraddr)  // 0x20e.
	h.Reserved11[3] = higherByteFromUint16(b.Header.Kveraddr) // 0x20f.

	h.LoaderType = b.Header.TypeOfLoader  // 0x210.
	h.LoaderFlagzaAq = b.Header.Loadflags // 0x211.

	h.Reserved12[0] = lowerByteFromUint16(b.Header.Setupmovesize)  // 0x212.
	h.Reserved12[1] = higherByteFromUint16(b.Header.Setupmovesize) // 0x213.

	h.KernelStart = b.Header.Code32Start  // 0x214.
	h.InitrdStart = b.Header.RamdiskImage // 0x218.
	h.InitrdSize = b.Header.RamdiskSize   // 0x21c.

	h.Reserved13 = b.Header.BootSectKludge // 0x220.

	// (2.01+)
	h.HeapEndPtr = b.Header.Heapendptr // 0x224.

	h.Reserved14[0] = b.Header.ExtLoaderVer  // 0x226.
	h.Reserved14[1] = b.Header.ExtLoaderType // 0x227.

	// (2.02+)
	h.CmdLinePtr = b.Header.Cmdlineptr // 0x228.

	// (2.03+)
	h.InitrdAddrMax = b.Header.InitrdAddrMax // 0x22c.

	// (2.05+)
	h.KernelAlignment = b.Header.Kernelalignment     // 0x230.
	h.RelocatableKernel = b.Header.RelocatableKernel // 0x234.
	h.MinAlignment = b.Header.MinAlignment           // 0x235 (2.10+).
	h.XloadFlags = b.Header.XLoadFlags               // 0x236.
	// (2.06+)
	h.CmdlineSize = b.Header.CmdLineSize // 0x238.
	// (2.07+)
	h.HardwareSubarch = b.Header.HardwareSubArch         // 0x23c.
	h.HardwareSubarchData = b.Header.HardwareSubArchData // 0x240.
	// (2.08+)
	h.PayloadOffset = b.Header.PayloadOffset // 0x248.
	h.PayloadLength = b.Header.PayloadSize   // 0x24c.
	// (2.09+)
	h.SetupData = b.Header.SetupData // 0x250.
	// (2.10+)
	h.PrefAddress = b.Header.PrefAddress       // 0x258.
	h.InitSize = b.Header.InitSize             // 0x260.
	h.HandoverOffset = b.Header.HandoverOffset // 0x264.
	return &h
}
