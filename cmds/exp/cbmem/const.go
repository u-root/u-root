// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package main

// Constants from coreboot. We will not change these names.
const (
	LB_TAG_UNUSED                = 0x0000
	LB_TAG_MEMORY                = 0x0001
	LB_TAG_HWRPB                 = 0x0002
	LB_TAG_MAINBOARD             = 0x0003
	LB_TAG_VERSION               = 0x0004
	LB_TAG_EXTRA_VERSION         = 0x0005
	LB_TAG_BUILD                 = 0x0006
	LB_TAG_COMPILE_TIME          = 0x0007
	LB_TAG_COMPILE_BY            = 0x0008
	LB_TAG_COMPILE_HOST          = 0x0009
	LB_TAG_COMPILE_DOMAIN        = 0x000a
	LB_TAG_COMPILER              = 0x000b
	LB_TAG_LINKER                = 0x000c
	LB_TAG_ASSEMBLER             = 0x000d
	LB_TAG_SERIAL                = 0x000f
	LB_TAG_CONSOLE               = 0x0010
	LB_TAG_FORWARD               = 0x0011
	LB_TAG_FRAMEBUFFER           = 0x0012
	LB_TAG_GPIO                  = 0x0013
	LB_TAG_TIMESTAMPS            = 0x0016
	LB_TAG_CBMEM_CONSOLE         = 0x0017
	LB_TAG_MRC_CACHE             = 0x0018
	LB_TAG_VBNV                  = 0x0019
	LB_TAG_VBOOT_HANDOFF         = 0x0020 // deprecated
	LB_TAG_X86_ROM_MTRR          = 0x0021
	LB_TAG_DMA                   = 0x0022
	LB_TAG_RAM_OOPS              = 0x0023
	LB_TAG_ACPI_GNVS             = 0x0024
	LB_TAG_BOARD_ID              = 0x0025 // deprecated
	LB_TAG_VERSION_TIMESTAMP     = 0x0026
	LB_TAG_WIFI_CALIBRATION      = 0x0027
	LB_TAG_RAM_CODE              = 0x0028 // deprecated
	LB_TAG_SPI_FLASH             = 0x0029
	LB_TAG_SERIALNO              = 0x002a
	LB_TAG_MTC                   = 0x002b
	LB_TAG_VPD                   = 0x002c
	LB_TAG_SKU_ID                = 0x002d // deprecated
	LB_TAG_BOOT_MEDIA_PARAMS     = 0x0030
	LB_TAG_CBMEM_ENTRY           = 0x0031
	LB_TAG_TSC_INFO              = 0x0032
	LB_TAG_MAC_ADDRS             = 0x0033
	LB_TAG_VBOOT_WORKBUF         = 0x0034
	LB_TAG_MMC_INFO              = 0x0035
	LB_TAG_TCPA_LOG              = 0x0036
	LB_TAG_FMAP                  = 0x0037
	LB_TAG_PLATFORM_BLOB_VERSION = 0x0038
	LB_TAG_SMMSTOREV2            = 0x0039
	LB_TAG_TPM_PPI_HANDOFF       = 0x003a
	LB_TAG_BOARD_CONFIG          = 0x0040
	/* The following options are CMOS-related */
	LB_TAG_CMOS_OPTION_TABLE = 0x00c8
	LB_TAG_OPTION            = 0x00c9
	LB_TAG_OPTION_ENUM       = 0x00ca
	LB_TAG_OPTION_DEFAULTS   = 0x00cb
	LB_TAG_OPTION_CHECKSUM   = 0x00cc
)

const (
	LB_MEM_RAM                   = 1
	LB_MEM_RESERVED              = 2  // Don't use this memory region
	LB_MEM_ACPI                  = 3  // ACPI Tables
	LB_MEM_NVS                   = 4  // ACPI NVS Memory
	LB_MEM_UNUSABLE              = 5  // Unusable address space
	LB_MEM_VENDOR_RSVD           = 6  // Vendor Reserved
	LB_MEM_TABLE                 = 16 // Ram configuration tables are kept in
	LB_SERIAL_TYPE_IO_MAPPED     = 1
	LB_SERIAL_TYPE_MEMORY_MAPPED = 2
	ACTIVE_LOW                   = 0
	ACTIVE_HIGH                  = 1
	GPIO_MAX_NAME_LENGTH         = 16
	LB_TAB_VBOOT_HANDOFF         = 0x0020
	LB_TAB_DMA                   = 0x0022
	MAX_SERIALNO_LENGTH          = 32
	CMOS_MAX_NAME_LENGTH         = 32
	CMOS_MAX_TEXT_LENGTH         = 32
	CMOS_IMAGE_BUFFER_SIZE       = 256
	CHECKSUM_NONE                = 0
	CHECKSUM_PCBIOS              = 1

	CBMC_CURSOR_MASK = ((1 << 28) - 1)
	CBMC_OVERFLOW    = (1 << 31)
)

var (
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
	consoleNames = map[uint32]string{}
	tagNames     = map[uint32]string{
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
		LB_TAG_SERIAL:                "LB_TAG_SERIAL",
		LB_TAG_CONSOLE:               "LB_TAG_CONSOLE",
		LB_TAG_FORWARD:               "LB_TAG_FORWARD",
		LB_TAG_FRAMEBUFFER:           "LB_TAG_FRAMEBUFFER",
		LB_TAG_GPIO:                  "LB_TAG_GPIO",
		LB_TAG_TIMESTAMPS:            "LB_TAG_TIMESTAMPS",
		LB_TAG_CBMEM_CONSOLE:         "LB_TAG_CBMEM_CONSOLE",
		LB_TAG_MRC_CACHE:             "LB_TAG_MRC_CACHE",
		LB_TAG_VBNV:                  "LB_TAG_VBNV",
		LB_TAG_VBOOT_HANDOFF:         "LB_TAG_VBOOT_HANDOFF",
		LB_TAG_X86_ROM_MTRR:          "LB_TAG_X86_ROM_MTRR",
		LB_TAG_DMA:                   "LB_TAG_DMA",
		LB_TAG_RAM_OOPS:              "LB_TAG_RAM_OOPS",
		LB_TAG_ACPI_GNVS:             "LB_TAG_ACPI_GNVS",
		LB_TAG_BOARD_ID:              "LB_TAG_BOARD_ID",
		LB_TAG_VERSION_TIMESTAMP:     "LB_TAG_VERSION_TIMESTAMP",
		LB_TAG_WIFI_CALIBRATION:      "LB_TAG_WIFI_CALIBRATION",
		LB_TAG_RAM_CODE:              "LB_TAG_RAM_CODE",
		LB_TAG_SPI_FLASH:             "LB_TAG_SPI_FLASH",
		LB_TAG_SERIALNO:              "LB_TAG_SERIALNO",
		LB_TAG_MTC:                   "LB_TAG_MTC",
		LB_TAG_VPD:                   "LB_TAG_VPD",
		LB_TAG_SKU_ID:                "LB_TAG_SKU_ID",
		LB_TAG_BOOT_MEDIA_PARAMS:     "LB_TAG_BOOT_MEDIA_PARAMS",
		LB_TAG_CBMEM_ENTRY:           "LB_TAG_CBMEM_ENTRY",
		LB_TAG_TSC_INFO:              "LB_TAG_TSC_INFO",
		LB_TAG_MAC_ADDRS:             "LB_TAG_MAC_ADDRS",
		LB_TAG_VBOOT_WORKBUF:         "LB_TAG_VBOOT_WORKBUF",
		LB_TAG_MMC_INFO:              "LB_TAG_MMC_INFO",
		LB_TAG_TCPA_LOG:              "LB_TAG_TCPA_LOG",
		LB_TAG_FMAP:                  "LB_TAG_FMAP",
		LB_TAG_PLATFORM_BLOB_VERSION: "LB_TAG_PLATFORM_BLOB_VERSION",
		LB_TAG_SMMSTOREV2:            "LB_TAG_SMMSTOREV2",
		LB_TAG_TPM_PPI_HANDOFF:       "LB_TAG_TPM_PPI_HANDOFF",
		LB_TAG_BOARD_CONFIG:          "LB_TAG_BOARD_CONFIG",
		LB_TAG_CMOS_OPTION_TABLE:     "LB_TAG_CMOS_OPTION_TABLE",
		LB_TAG_OPTION:                "LB_TAG_OPTION",
		LB_TAG_OPTION_ENUM:           "LB_TAG_OPTION_ENUM",
		LB_TAG_OPTION_DEFAULTS:       "LB_TAG_OPTION_DEFAULTS",
		LB_TAG_OPTION_CHECKSUM:       "LB_TAG_OPTION_CHECKSUM",
	}
)
