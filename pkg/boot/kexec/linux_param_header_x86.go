package kexec

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
