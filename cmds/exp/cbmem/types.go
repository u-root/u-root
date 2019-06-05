// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const (
	TableSize = 48
)

// TODO: we don't have a platform yet that creates these.
type TS struct {
	entry_id    uint32
	entry_stamp uint64
}

type TSTable struct {
	base_time     uint64
	max_entries   uint16
	tick_freq_mhz uint16
	num_entries   uint32
	entries       []TS
}

const (
	TS_START_ROMSTAGE            int = 1
	TS_BEFORE_INITRAM            int = 2
	TS_AFTER_INITRAM             int = 3
	TS_END_ROMSTAGE              int = 4
	TS_START_VBOOT               int = 5
	TS_END_VBOOT                 int = 6
	TS_START_COPYRAM             int = 8
	TS_END_COPYRAM               int = 9
	TS_START_RAMSTAGE            int = 10
	TS_START_BOOTBLOCK           int = 11
	TS_END_BOOTBLOCK             int = 12
	TS_START_COPYROM             int = 13
	TS_END_COPYROM               int = 14
	TS_START_ULZMA               int = 15
	TS_END_ULZMA                 int = 16
	TS_DEVICE_ENUMERATE          int = 30
	TS_DEVICE_CONFIGURE          int = 40
	TS_DEVICE_ENABLE             int = 50
	TS_DEVICE_INITIALIZE         int = 60
	TS_DEVICE_DONE               int = 70
	TS_CBMEM_POST                int = 75
	TS_WRITE_TABLES              int = 80
	TS_LOAD_PAYLOAD              int = 90
	TS_ACPI_WAKE_JUMP            int = 98
	TS_SELFBOOT_JUMP             int = 99
	TS_START_COPYVER             int = 501
	TS_END_COPYVER               int = 502
	TS_START_TPMINIT             int = 503
	TS_END_TPMINIT               int = 504
	TS_START_VERIFY_SLOT         int = 505
	TS_END_VERIFY_SLOT           int = 506
	TS_START_HASH_BODY           int = 507
	TS_DONE_LOADING              int = 508
	TS_DONE_HASHING              int = 509
	TS_END_HASH_BODY             int = 510
	TS_FSP_MEMORY_INIT_START     int = 950
	TS_FSP_MEMORY_INIT_END       int = 951
	TS_FSP_TEMP_RAM_EXIT_START   int = 952
	TS_FSP_TEMP_RAM_EXIT_END     int = 953
	TS_FSP_SILICON_INIT_START    int = 954
	TS_FSP_SILICON_INIT_END      int = 955
	TS_FSP_BEFORE_ENUMERATE      int = 956
	TS_FSP_AFTER_ENUMERATE       int = 957
	TS_FSP_BEFORE_FINALIZE       int = 958
	TS_FSP_AFTER_FINALIZE        int = 959
	LB_TAG_UNUSED                    = 0x0000
	LB_TAG_MEMORY                    = 0x0001
	LB_MEM_RAM                       = 1
	LB_MEM_RESERVED                  = 2  // Don't use this memory region
	LB_MEM_ACPI                      = 3  // ACPI Tables
	LB_MEM_NVS                       = 4  // ACPI NVS Memory
	LB_MEM_UNUSABLE                  = 5  // Unusable address space
	LB_MEM_VENDOR_RSVD               = 6  // Vendor Reserved
	LB_MEM_TABLE                     = 16 // Ram configuration tables are kept in
	LB_TAG_HWRPB                     = 0x0002
	LB_TAG_MAINBOARD                 = 0x0003
	LB_TAG_VERSION                   = 0x0004
	LB_TAG_EXTRA_VERSION             = 0x0005
	LB_TAG_BUILD                     = 0x0006
	LB_TAG_COMPILE_TIME              = 0x0007
	LB_TAG_COMPILE_BY                = 0x0008
	LB_TAG_COMPILE_HOST              = 0x0009
	LB_TAG_COMPILE_DOMAIN            = 0x000a
	LB_TAG_COMPILER                  = 0x000b
	LB_TAG_LINKER                    = 0x000c
	LB_TAG_ASSEMBLER                 = 0x000d
	LB_TAG_VERSION_TIMESTAMP         = 0x0026
	LB_TAG_SERIAL                    = 0x000f
	LB_SERIAL_TYPE_IO_MAPPED         = 1
	LB_SERIAL_TYPE_MEMORY_MAPPED     = 2
	LB_TAG_CONSOLE                   = 0x0010
	LB_TAG_CONSOLE_SERIAL8250        = 0
	LB_TAG_CONSOLE_VGA               = 1 // OBSOLETE
	LB_TAG_CONSOLE_BTEXT             = 2 // OBSOLETE
	LB_TAG_CONSOLE_LOGBUF            = 3 // OBSOLETE
	LB_TAG_CONSOLE_SROM              = 4 // OBSOLETE
	LB_TAG_CONSOLE_EHCI              = 5
	LB_TAG_CONSOLE_SERIAL8250MEM     = 6
	LB_TAG_FORWARD                   = 0x0011
	LB_TAG_FRAMEBUFFER               = 0x0012
	LB_TAG_GPIO                      = 0x0013
	ACTIVE_LOW                       = 0
	ACTIVE_HIGH                      = 1
	GPIO_MAX_NAME_LENGTH             = 16
	LB_TAG_VDAT                      = 0x0015
	LB_TAG_VBNV                      = 0x0019
	LB_TAB_VBOOT_HANDOFF             = 0x0020
	LB_TAB_DMA                       = 0x0022
	LB_TAG_RAM_OOPS                  = 0x0023
	LB_TAG_MTC                       = 0x002b
	LB_TAG_TIMESTAMPS                = 0x0016
	LB_TAG_CBMEM_CONSOLE             = 0x0017
	LB_TAG_MRC_CACHE                 = 0x0018
	LB_TAG_ACPI_GNVS                 = 0x0024
	LB_TAG_WIFI_CALIBRATION          = 0x0027
	LB_TAG_X86_ROM_MTRR              = 0x0021
	LB_TAG_BOARD_ID                  = 0x0025
	LB_TAG_RAM_CODE                  = 0x0028
	LB_TAG_SPI_FLASH                 = 0x0029
	LB_TAG_BOOT_MEDIA_PARAMS         = 0x0030
	LB_TAG_CBMEM_ENTRY               = 0x0031
	LB_TAG_SERIALNO                  = 0x002a
	LB_TAG_MAC_ADDRS                 = 0x0033
	MAX_SERIALNO_LENGTH              = 32
	LB_TAG_CMOS_OPTION_TABLE         = 200
	LB_TAG_OPTION                    = 201
	CMOS_MAX_NAME_LENGTH             = 32
	LB_TAG_OPTION_ENUM               = 202
	CMOS_MAX_TEXT_LENGTH             = 32
	LB_TAG_OPTION_DEFAULTS           = 203
	CMOS_IMAGE_BUFFER_SIZE           = 256
	LB_TAG_OPTION_CHECKSUM           = 204
	CHECKSUM_NONE                    = 0
	CHECKSUM_PCBIOS                  = 1

	// Depthcharge entry IDs start at 1000.
	TS_DC_START                     = 1000
	TS_RO_PARAMS_INIT               = 1001
	TS_RO_VB_INIT                   = 1002
	TS_RO_VB_SELECT_FIRMWARE        = 1003
	TS_RO_VB_SELECT_AND_LOAD_KERNEL = 1004

	TS_RW_VB_SELECT_AND_LOAD_KERNEL = 1010

	TS_VB_SELECT_AND_LOAD_KERNEL = 1020

	TS_VB_EC_VBOOT_DONE = 1030

	TS_CROSSYSTEM_DATA = 1100
	TS_START_KERNEL    = 1101
)

type Header struct {
	Signature    [4]uint8
	HeaderSz     uint32
	HeaderCSUM   uint32
	TableSz      uint32
	TableCSUM    uint32
	TableEntries uint32
}
type Record struct {
	Tag  uint32
	Size uint32
}

type memoryRange struct {
	Start uint64
	Size  uint64
	Mtype uint32
}
type memoryEntry struct {
	Record
	Maps []memoryRange
}
type hwrpbEntry struct {
	Record
	HwrPB uint64
}

type consoleEntry struct {
	Record
	Console uint32
}

type mainboardEntry struct {
	Record
	Vendor     string
	PartNumber string
}

type stringEntry struct {
	Record
	String []uint8
}

type timestampEntry struct {
	Record
	TimeStamp uint32
}
type serialEntry struct {
	Record
	Type     uint32
	BaseAddr uint32
	Baud     uint32
	RegWidth uint32
}
type memconsoleEntry struct {
	Record
	CSize  uint32
	Cursor uint32
	Data   []byte
}
type forwardEntry struct {
	Record
	Forward uint64
}
type framebufferEntry struct {
	Record
	PhysicalAddress  uint64
	XResolution      uint32
	YRresolution     uint32
	BytesPerLine     uint32
	BitsPerPixel     uint8
	RedMaskPos       uint8
	RedMaskSize      uint8
	GreenMaskPos     uint8
	GreenMaskSize    uint8
	BlueMaskPos      uint8
	BlueMaskSize     uint8
	ReservedMaskPos  uint8
	ReservedMaskSize uint8
}

// GPIO is General Purpose IO
type GPIO struct {
	Port     uint32
	Polarity uint32
	Value    uint32
	Name     []uint8
}

type gpioEntry struct {
	Record
	Count uint32
	GPIOs []GPIO
}

type rangeEntry struct {
	Record
	RangeStart uint64
	RangeSize  uint32
}

type cbmemEntry struct {
	Record
	CbmemAddr uint64
}

type MTRREntry struct {
	Record
	Index uint32
}

type BoardIDEntry struct {
	Record
	BoardID uint32
}
type MACEntry struct {
	MACaddr []uint8
	pad     []uint8
}

type LBEntry struct {
	Record
	Count    uint32
	MACAddrs []MACEntry
}

type LBRAMEntry struct {
	Record
	RamCode uint32
}

type SPIFlashEntry struct {
	Record
	FlashSize  uint32
	SectorSize uint32
	EraseCmd   uint32
}
type bootMediaParamsEntry struct {
	Record
	FMAPOffset    uint64
	CBFSOffset    uint64
	CBFSSize      uint64
	BootMediaSize uint64
}
type CBMemEntry struct {
	Record
	Address   uint64
	EntrySize uint32
	ID        uint32
}
type CMOSTable struct {
	Record
	HeaderLength uint32
}
type CMOSEntry struct {
	Record
	Bit      uint32
	Length   uint32
	Config   uint32
	ConfigID uint32
	Name     []uint8
}
type CMOSEnums struct {
	Record
	ConfigID uint32
	Value    uint32
	Text     []uint8
}
type CMOSDefaults struct {
	Record
	NameLength uint32
	Name       []uint8
	DefaultSet []uint8
}
type CMOSChecksum struct {
	Record
	RangeStart uint32
	RangeEnd   uint32
	Location   uint32
	Type       uint32
}

type CBmem struct {
	Memory           *memoryEntry
	MemConsole       *memconsoleEntry
	Consoles         []string
	TimeStamps       []timestampEntry
	UART             []serialEntry
	MainBoard        mainboardEntry
	Hwrpb            hwrpbEntry
	CBMemory         []cbmemEntry
	BoardID          BoardIDEntry
	StringVars       map[string]string
	BootMediaParams  bootMediaParamsEntry
	VersionTimeStamp uint32
	Unknown          []uint32
	Ignored          []string
}
