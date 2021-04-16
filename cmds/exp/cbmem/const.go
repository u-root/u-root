// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package main

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
	LB_TAG_PLATFORM_BLOB_VERSION     = 0x0038
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
	CBMC_CURSOR_MASK   = ((1 << 28) - 1)
	CBMC_OVERFLOW      = (1 << 31)
)

var (
	TimeStampNames = map[int]string{
		0:                    "1st timestamp",
		TS_START_ROMSTAGE:    "start of rom stage",
		TS_BEFORE_INITRAM:    "before ram initialization",
		TS_AFTER_INITRAM:     "after ram initialization",
		TS_END_ROMSTAGE:      "end of romstage",
		TS_START_VBOOT:       "start of verified boot",
		TS_END_VBOOT:         "end of verified boot",
		TS_START_COPYRAM:     "starting to load ramstage",
		TS_END_COPYRAM:       "finished loading ramstage",
		TS_START_RAMSTAGE:    "start of ramstage",
		TS_START_BOOTBLOCK:   "start of bootblock",
		TS_END_BOOTBLOCK:     "end of bootblock",
		TS_START_COPYROM:     "starting to load romstage",
		TS_END_COPYROM:       "finished loading romstage",
		TS_START_ULZMA:       "starting LZMA decompress (ignore for x86)",
		TS_END_ULZMA:         "finished LZMA decompress (ignore for x86)",
		TS_DEVICE_ENUMERATE:  "device enumeration",
		TS_DEVICE_CONFIGURE:  "device configuration",
		TS_DEVICE_ENABLE:     "device enable",
		TS_DEVICE_INITIALIZE: "device initialization",
		TS_DEVICE_DONE:       "device setup done",
		TS_CBMEM_POST:        "cbmem post",
		TS_WRITE_TABLES:      "write tables",
		TS_LOAD_PAYLOAD:      "load payload",
		TS_ACPI_WAKE_JUMP:    "ACPI wake jump",
		TS_SELFBOOT_JUMP:     "selfboot jump",

		TS_START_COPYVER:     "starting to load verstage",
		TS_END_COPYVER:       "finished loading verstage",
		TS_START_TPMINIT:     "starting to initialize TPM",
		TS_END_TPMINIT:       "finished TPM initialization",
		TS_START_VERIFY_SLOT: "starting to verify keyblock/preamble (RSA)",
		TS_END_VERIFY_SLOT:   "finished verifying keyblock/preamble (RSA)",
		TS_START_HASH_BODY:   "starting to verify body (load+SHA2+RSA) ",
		TS_DONE_LOADING:      "finished loading body (ignore for x86)",
		TS_DONE_HASHING:      "finished calculating body hash (SHA2)",
		TS_END_HASH_BODY:     "finished verifying body signature (RSA)",

		TS_DC_START:                     "depthcharge start",
		TS_RO_PARAMS_INIT:               "RO parameter init",
		TS_RO_VB_INIT:                   "RO vboot init",
		TS_RO_VB_SELECT_FIRMWARE:        "RO vboot select firmware",
		TS_RO_VB_SELECT_AND_LOAD_KERNEL: "RO vboot select&load kernel",
		TS_RW_VB_SELECT_AND_LOAD_KERNEL: "RW vboot select&load kernel",
		TS_VB_SELECT_AND_LOAD_KERNEL:    "vboot select&load kernel",
		TS_VB_EC_VBOOT_DONE:             "finished EC verification",
		TS_CROSSYSTEM_DATA:              "crossystem data",
		TS_START_KERNEL:                 "start kernel",

		// FSP related timestamps
		TS_FSP_MEMORY_INIT_START:   "calling FspMemoryInit",
		TS_FSP_MEMORY_INIT_END:     "returning from FspMemoryInit",
		TS_FSP_TEMP_RAM_EXIT_START: "calling FspTempRamExit",
		TS_FSP_TEMP_RAM_EXIT_END:   "returning from FspTempRamExit",
		TS_FSP_SILICON_INIT_START:  "calling FspSiliconInit",
		TS_FSP_SILICON_INIT_END:    "returning from FspSiliconInit",
		TS_FSP_BEFORE_ENUMERATE:    "calling FspNotify(AfterPciEnumeration)",
		TS_FSP_AFTER_ENUMERATE:     "returning from FspNotify(AfterPciEnumeration)",
		TS_FSP_BEFORE_FINALIZE:     "calling FspNotify(ReadyToBoot)",
		TS_FSP_AFTER_FINALIZE:      "returning from FspNotify(ReadyToBoot)",
	}

	memTags = map[uint32]string{
		LB_MEM_RAM:         "LB_MEM_RAM",
		LB_MEM_RESERVED:    "LB_MEM_RESERVED",
		LB_MEM_ACPI:        "LB_MEM_ACPI",
		LB_MEM_NVS:         "LB_MEM_NVS",
		LB_MEM_UNUSABLE:    "LB_MEM_UNUSABLE",
		LB_MEM_VENDOR_RSVD: "LB_MEM_VENDOR_RSVD",
		LB_MEM_TABLE:       "LB_MEM_TABLE",
	}
	serialNames = map[uint32]string{
		LB_SERIAL_TYPE_IO_MAPPED:     "IO_MAPPED",
		LB_SERIAL_TYPE_MEMORY_MAPPED: "MEMORY_MAPPED",
	}
	consoleNames = map[uint32]string{
		LB_TAG_CONSOLE_SERIAL8250:    "SERIAL8250",
		LB_TAG_CONSOLE_VGA:           "VGA",
		LB_TAG_CONSOLE_BTEXT:         "BTEXT",
		LB_TAG_CONSOLE_LOGBUF:        "LOGBUF",
		LB_TAG_CONSOLE_SROM:          "SROM",
		LB_TAG_CONSOLE_EHCI:          "EHCI",
		LB_TAG_CONSOLE_SERIAL8250MEM: "SERIAL8250MEM",
	}
	tagNames = map[uint32]string{
		LB_TAG_UNUSED:                "LB_TAG_UNUSED",
		LB_TAG_MEMORY:                "LB_TAG_MEMORY",
		LB_TAG_HWRPB:                 "LB_TAG_HWRPB",
		LB_TAG_MAINBOARD:             "LB_TAG_MAINBOARD",
		LB_TAG_VERSION:               "LB_TAG_VERSION",
		LB_TAG_EXTRA_VERSION:         "LB_TAG_EXTRA_VERSION",
		LB_TAG_BUILD:                 "LB_TAG_BUILD",
		LB_TAG_COMPILE_TIME:          "LB_TAG_COMPILE_TIME",
		LB_TAG_COMPILE_BY:            "LB_TAG_COMPILE_BY",
		LB_TAG_COMPILE_HOST:          "LB_TAG_COMPILE_HOST",
		LB_TAG_COMPILE_DOMAIN:        "LB_TAG_COMPILE_DOMAIN",
		LB_TAG_COMPILER:              "LB_TAG_COMPILER",
		LB_TAG_LINKER:                "LB_TAG_LINKER",
		LB_TAG_ASSEMBLER:             "LB_TAG_ASSEMBLER",
		LB_TAG_VERSION_TIMESTAMP:     "LB_TAG_VERSION_TIMESTAMP",
		LB_TAG_SERIAL:                "LB_TAG_SERIAL",
		LB_TAG_CONSOLE:               "LB_TAG_CONSOLE",
		LB_TAG_FORWARD:               "LB_TAG_FORWARD",
		LB_TAG_FRAMEBUFFER:           "LB_TAG_FRAMEBUFFER",
		LB_TAG_GPIO:                  "LB_TAG_GPIO",
		LB_TAG_VDAT:                  "LB_TAG_VDAT",
		LB_TAG_VBNV:                  "LB_TAG_VBNV",
		LB_TAB_VBOOT_HANDOFF:         "LB_TAB_VBOOT_HANDOFF",
		LB_TAB_DMA:                   "LB_TAB_DMA",
		LB_TAG_RAM_OOPS:              "LB_TAG_RAM_OOPS",
		LB_TAG_MTC:                   "LB_TAG_MTC",
		LB_TAG_TIMESTAMPS:            "LB_TAG_TIMESTAMPS",
		LB_TAG_CBMEM_CONSOLE:         "LB_TAG_CBMEM_CONSOLE",
		LB_TAG_MRC_CACHE:             "LB_TAG_MRC_CACHE",
		LB_TAG_ACPI_GNVS:             "LB_TAG_ACPI_GNVS",
		LB_TAG_WIFI_CALIBRATION:      "LB_TAG_WIFI_CALIBRATION",
		LB_TAG_X86_ROM_MTRR:          "LB_TAG_X86_ROM_MTRR",
		LB_TAG_BOARD_ID:              "LB_TAG_BOARD_ID",
		LB_TAG_MAC_ADDRS:             "LB_TAG_MAC_ADDRS",
		LB_TAG_RAM_CODE:              "LB_TAG_RAM_CODE",
		LB_TAG_SPI_FLASH:             "LB_TAG_SPI_FLASH",
		LB_TAG_BOOT_MEDIA_PARAMS:     "LB_TAG_BOOT_MEDIA_PARAMS",
		LB_TAG_CBMEM_ENTRY:           "LB_TAG_CBMEM_ENTRY",
		LB_TAG_SERIALNO:              "LB_TAG_SERIALNO",
		LB_TAG_CMOS_OPTION_TABLE:     "LB_TAG_CMOS_OPTION_TABLE",
		LB_TAG_OPTION:                "LB_TAG_OPTION",
		LB_TAG_OPTION_ENUM:           "LB_TAG_OPTION_ENUM",
		LB_TAG_OPTION_DEFAULTS:       "LB_TAG_OPTION_DEFAULTS",
		LB_TAG_OPTION_CHECKSUM:       "LB_TAG_OPTION_CHECKSUM",
		LB_TAG_PLATFORM_BLOB_VERSION: "LB_TAG_PLATFORM_BLOB_VERSION",
	}
	tsNames = map[uint32]string{
		TS_DC_START:                     "TS_DC_START",
		TS_RO_PARAMS_INIT:               "TS_RO_PARAMS_INIT",
		TS_RO_VB_INIT:                   "TS_RO_VB_INIT",
		TS_RO_VB_SELECT_FIRMWARE:        "TS_RO_VB_SELECT_FIRMWARE",
		TS_RO_VB_SELECT_AND_LOAD_KERNEL: "TS_RO_VB_SELECT_AND_LOAD_KERNEL",
		TS_RW_VB_SELECT_AND_LOAD_KERNEL: "TS_RW_VB_SELECT_AND_LOAD_KERNEL",
		TS_VB_SELECT_AND_LOAD_KERNEL:    "TS_VB_SELECT_AND_LOAD_KERNEL",
		TS_VB_EC_VBOOT_DONE:             "TS_VB_EC_VBOOT_DONE",
		TS_CROSSYSTEM_DATA:              "TS_CROSSYSTEM_DATA",
		TS_START_KERNEL:                 "TS_START_KERNEL",
	}
)
