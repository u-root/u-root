// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cbmem prints out coreboot mem table information in JSON by default,
// and also implements the basic cbmem -list and -console commands.
// TODO: checksum tables.
package main

import (
	"fmt"
	"os"
	"sort"
	"unsafe"
)

const (
	TS_START_ROMSTAGE    uint32 = 1
	TS_BEFORE_INITRAM    uint32 = 2
	TS_AFTER_INITRAM     uint32 = 3
	TS_END_ROMSTAGE      uint32 = 4
	TS_START_VBOOT       uint32 = 5
	TS_END_VBOOT         uint32 = 6
	TS_START_COPYRAM     uint32 = 8
	TS_END_COPYRAM       uint32 = 9
	TS_START_RAMSTAGE    uint32 = 10
	TS_START_BOOTBLOCK   uint32 = 11
	TS_END_BOOTBLOCK     uint32 = 12
	TS_START_COPYROM     uint32 = 13
	TS_END_COPYROM       uint32 = 14
	TS_START_ULZMA       uint32 = 15
	TS_END_ULZMA         uint32 = 16
	TS_START_ULZ4F       uint32 = 17
	TS_END_ULZ4F         uint32 = 18
	TS_DEVICE_ENUMERATE  uint32 = 30
	TS_DEVICE_CONFIGURE  uint32 = 40
	TS_DEVICE_ENABLE     uint32 = 50
	TS_DEVICE_INITIALIZE uint32 = 60
	TS_OPROM_INITIALIZE  uint32 = 65
	TS_OPROM_COPY_END    uint32 = 66
	TS_OPROM_END         uint32 = 67
	TS_DEVICE_DONE       uint32 = 70
	TS_CBMEM_POST        uint32 = 75
	TS_WRITE_TABLES      uint32 = 80
	TS_FINALIZE_CHIPS    uint32 = 85
	TS_LOAD_PAYLOAD      uint32 = 90
	TS_ACPI_WAKE_JUMP    uint32 = 98
	TS_SELFBOOT_JUMP     uint32 = 99
	TS_START_POSTCAR     uint32 = 100
	TS_END_POSTCAR       uint32 = 101
	TS_DELAY_START       uint32 = 110
	TS_DELAY_END         uint32 = 111

	/* 500+ reserved for vendorcode extensions (500-600: google/chromeos) */
	TS_START_COPYVER       uint32 = 501
	TS_END_COPYVER         uint32 = 502
	TS_START_TPMINIT       uint32 = 503
	TS_END_TPMINIT         uint32 = 504
	TS_START_VERIFY_SLOT   uint32 = 505
	TS_END_VERIFY_SLOT     uint32 = 506
	TS_START_HASH_BODY     uint32 = 507
	TS_DONE_LOADING        uint32 = 508
	TS_DONE_HASHING        uint32 = 509
	TS_END_HASH_BODY       uint32 = 510
	TS_START_TPMPCR        uint32 = 511
	TS_END_TPMPCR          uint32 = 512
	TS_START_TPMLOCK       uint32 = 513
	TS_END_TPMLOCK         uint32 = 514
	TS_START_EC_SYNC       uint32 = 515
	TS_EC_HASH_READY       uint32 = 516
	TS_EC_POWER_LIMIT_WAIT uint32 = 517
	TS_END_EC_SYNC         uint32 = 518
	TS_START_COPYVPD       uint32 = 550
	TS_END_COPYVPD_RO      uint32 = 551
	TS_END_COPYVPD_RW      uint32 = 552

	/* 900-920 reserved for vendorcode extensions (900-940: AMD AGESA) */
	TS_AGESA_INIT_RESET_START  uint32 = 900
	TS_AGESA_INIT_RESET_DONE   uint32 = 901
	TS_AGESA_INIT_EARLY_START  uint32 = 902
	TS_AGESA_INIT_EARLY_DONE   uint32 = 903
	TS_AGESA_INIT_POST_START   uint32 = 904
	TS_AGESA_INIT_POST_DONE    uint32 = 905
	TS_AGESA_INIT_ENV_START    uint32 = 906
	TS_AGESA_INIT_ENV_DONE     uint32 = 907
	TS_AGESA_INIT_MID_START    uint32 = 908
	TS_AGESA_INIT_MID_DONE     uint32 = 909
	TS_AGESA_INIT_LATE_START   uint32 = 910
	TS_AGESA_INIT_LATE_DONE    uint32 = 911
	TS_AGESA_INIT_RTB_START    uint32 = 912
	TS_AGESA_INIT_RTB_DONE     uint32 = 913
	TS_AGESA_INIT_RESUME_START uint32 = 914
	TS_AGESA_INIT_RESUME_DONE  uint32 = 915
	TS_AGESA_S3_LATE_START     uint32 = 916
	TS_AGESA_S3_LATE_DONE      uint32 = 917
	TS_AGESA_S3_FINAL_START    uint32 = 918
	TS_AGESA_S3_FINAL_DONE     uint32 = 919

	/* 940-950 reserved for vendorcode extensions (940-950: Intel ME) */
	TS_ME_INFORM_DRAM_WAIT uint32 = 940
	TS_ME_INFORM_DRAM_DONE uint32 = 941

	/* 950+ reserved for vendorcode extensions (950-999: intel/fsp) */
	TS_FSP_MEMORY_INIT_START         uint32 = 950
	TS_FSP_MEMORY_INIT_END           uint32 = 951
	TS_FSP_TEMP_RAM_EXIT_START       uint32 = 952
	TS_FSP_TEMP_RAM_EXIT_END         uint32 = 953
	TS_FSP_SILICON_INIT_START        uint32 = 954
	TS_FSP_SILICON_INIT_END          uint32 = 955
	TS_FSP_BEFORE_ENUMERATE          uint32 = 956
	TS_FSP_AFTER_ENUMERATE           uint32 = 957
	TS_FSP_BEFORE_FINALIZE           uint32 = 958
	TS_FSP_AFTER_FINALIZE            uint32 = 959
	TS_FSP_BEFORE_END_OF_FIRMWARE    uint32 = 960
	TS_FSP_AFTER_END_OF_FIRMWARE     uint32 = 961
	TS_FSP_MULTI_PHASE_SI_INIT_START uint32 = 962
	TS_FSP_MULTI_PHASE_SI_INIT_END   uint32 = 963

	/* 1000+ reserved for payloads (1000-1200: ChromeOS depthcharge) */

	/* Depthcharge entry IDs start at 1000 */
	TS_DC_START uint32 = 1000

	TS_RO_PARAMS_INIT               uint32 = 1001
	TS_RO_VB_INIT                   uint32 = 1002
	TS_RO_VB_SELECT_FIRMWARE        uint32 = 1003
	TS_RO_VB_SELECT_AND_LOAD_KERNEL uint32 = 1004

	TS_RW_VB_SELECT_AND_LOAD_KERNEL uint32 = 1010

	TS_VB_SELECT_AND_LOAD_KERNEL uint32 = 1020
	TS_VB_EC_VBOOT_DONE          uint32 = 1030
	TS_VB_STORAGE_INIT_DONE      uint32 = 1040
	TS_VB_READ_KERNEL_DONE       uint32 = 1050
	TS_VB_VBOOT_DONE             uint32 = 1100

	TS_START_KERNEL         uint32 = 1101
	TS_KERNEL_DECOMPRESSION uint32 = 1102
)

// TimeStampNames map timestamp ints to names.
var TimeStampNames = map[uint32]string{
	0:                    "1st timestamp",
	TS_START_ROMSTAGE:    "start of romstage",
	TS_BEFORE_INITRAM:    "before RAM initialization",
	TS_AFTER_INITRAM:     "after RAM initialization",
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
	TS_START_ULZ4F:       "starting LZ4 decompress (ignore for x86)",
	TS_END_ULZ4F:         "finished LZ4 decompress (ignore for x86)",
	TS_DEVICE_ENUMERATE:  "device enumeration",
	TS_DEVICE_CONFIGURE:  "device configuration",
	TS_DEVICE_ENABLE:     "device enable",
	TS_DEVICE_INITIALIZE: "device initialization",
	TS_OPROM_INITIALIZE:  "Option ROM initialization",
	TS_OPROM_COPY_END:    "Option ROM copy done",
	TS_OPROM_END:         "Option ROM run done",
	TS_DEVICE_DONE:       "device setup done",
	TS_CBMEM_POST:        "cbmem post",
	TS_WRITE_TABLES:      "write tables",
	TS_FINALIZE_CHIPS:    "finalize chips",
	TS_LOAD_PAYLOAD:      "load payload",
	TS_ACPI_WAKE_JUMP:    "ACPI wake jump",
	TS_SELFBOOT_JUMP:     "selfboot jump",
	TS_DELAY_START:       "Forced delay start",
	TS_DELAY_END:         "Forced delay end",

	TS_START_COPYVER:     "starting to load verstage",
	TS_END_COPYVER:       "finished loading verstage",
	TS_START_TPMINIT:     "starting to initialize TPM",
	TS_END_TPMINIT:       "finished TPM initialization",
	TS_START_VERIFY_SLOT: "starting to verify keyblock/preamble (RSA)",
	TS_END_VERIFY_SLOT:   "finished verifying keyblock/preamble (RSA)",
	TS_START_HASH_BODY:   "starting to verify body (load+SHA2+RSA) ",
	TS_DONE_LOADING:      "finished loading body",
	TS_DONE_HASHING:      "finished calculating body hash (SHA2)",
	TS_END_HASH_BODY:     "finished verifying body signature (RSA)",
	TS_START_TPMPCR:      "starting TPM PCR extend",
	TS_END_TPMPCR:        "finished TPM PCR extend",
	TS_START_TPMLOCK:     "starting locking TPM",
	TS_END_TPMLOCK:       "finished locking TPM",

	TS_START_COPYVPD:  "starting to load Chrome OS VPD",
	TS_END_COPYVPD_RO: "finished loading Chrome OS VPD (RO)",
	TS_END_COPYVPD_RW: "finished loading Chrome OS VPD (RW)",

	TS_START_EC_SYNC:       "starting EC software sync",
	TS_EC_HASH_READY:       "EC vboot hash ready",
	TS_EC_POWER_LIMIT_WAIT: "waiting for EC to allow higher power draw",
	TS_END_EC_SYNC:         "finished EC software sync",

	TS_DC_START:                     "depthcharge start",
	TS_RO_PARAMS_INIT:               "RO parameter init",
	TS_RO_VB_INIT:                   "RO vboot init",
	TS_RO_VB_SELECT_FIRMWARE:        "RO vboot select firmware",
	TS_RO_VB_SELECT_AND_LOAD_KERNEL: "RO vboot select&load kernel",
	TS_RW_VB_SELECT_AND_LOAD_KERNEL: "RW vboot select&load kernel",
	TS_VB_SELECT_AND_LOAD_KERNEL:    "vboot select&load kernel",
	TS_VB_EC_VBOOT_DONE:             "finished EC verification",
	TS_VB_STORAGE_INIT_DONE:         "finished storage device initialization",
	TS_VB_READ_KERNEL_DONE:          "finished reading kernel from disk",
	TS_VB_VBOOT_DONE:                "finished vboot kernel verification",
	TS_KERNEL_DECOMPRESSION:         "starting kernel decompression/relocation",
	TS_START_KERNEL:                 "jumping to kernel",

	/* AMD AGESA related timestamps */
	TS_AGESA_INIT_RESET_START:  "calling AmdInitReset",
	TS_AGESA_INIT_RESET_DONE:   "back from AmdInitReset",
	TS_AGESA_INIT_EARLY_START:  "calling AmdInitEarly",
	TS_AGESA_INIT_EARLY_DONE:   "back from AmdInitEarly",
	TS_AGESA_INIT_POST_START:   "calling AmdInitPost",
	TS_AGESA_INIT_POST_DONE:    "back from AmdInitPost",
	TS_AGESA_INIT_ENV_START:    "calling AmdInitEnv",
	TS_AGESA_INIT_ENV_DONE:     "back from AmdInitEnv",
	TS_AGESA_INIT_MID_START:    "calling AmdInitMid",
	TS_AGESA_INIT_MID_DONE:     "back from AmdInitMid",
	TS_AGESA_INIT_LATE_START:   "calling AmdInitLate",
	TS_AGESA_INIT_LATE_DONE:    "back from AmdInitLate",
	TS_AGESA_INIT_RTB_START:    "calling AmdInitRtb/AmdS3Save",
	TS_AGESA_INIT_RTB_DONE:     "back from AmdInitRtb/AmdS3Save",
	TS_AGESA_INIT_RESUME_START: "calling AmdInitResume",
	TS_AGESA_INIT_RESUME_DONE:  "back from AmdInitResume",
	TS_AGESA_S3_LATE_START:     "calling AmdS3LateRestore",
	TS_AGESA_S3_LATE_DONE:      "back from AmdS3LateRestore",
	TS_AGESA_S3_FINAL_START:    "calling AmdS3FinalRestore",
	TS_AGESA_S3_FINAL_DONE:     "back from AmdS3FinalRestore",

	/* Intel ME related timestamps */
	TS_ME_INFORM_DRAM_WAIT: "waiting for ME acknowledgement of raminit",
	TS_ME_INFORM_DRAM_DONE: "finished waiting for ME response",

	/* FSP related timestamps */
	TS_FSP_MEMORY_INIT_START:      "calling FspMemoryInit",
	TS_FSP_MEMORY_INIT_END:        "returning from FspMemoryInit",
	TS_FSP_TEMP_RAM_EXIT_START:    "calling FspTempRamExit",
	TS_FSP_TEMP_RAM_EXIT_END:      "returning from FspTempRamExit",
	TS_FSP_SILICON_INIT_START:     "calling FspSiliconInit",
	TS_FSP_SILICON_INIT_END:       "returning from FspSiliconInit",
	TS_FSP_BEFORE_ENUMERATE:       "calling FspNotify(AfterPciEnumeration)",
	TS_FSP_AFTER_ENUMERATE:        "returning from FspNotify(AfterPciEnumeration)",
	TS_FSP_BEFORE_FINALIZE:        "calling FspNotify(ReadyToBoot)",
	TS_FSP_AFTER_FINALIZE:         "returning from FspNotify(ReadyToBoot)",
	TS_FSP_BEFORE_END_OF_FIRMWARE: "calling FspNotify(EndOfFirmware)",
	TS_FSP_AFTER_END_OF_FIRMWARE:  "returning from FspNotify(EndOfFirmware)",
	TS_START_POSTCAR:              "start of postcar",
	TS_END_POSTCAR:                "end of postcar",
}

// ByTime implements sort.Interface for []TS based on
// the EntryStamp field.
type ByTime []TS

func (bt ByTime) Len() int           { return len(bt) }
func (bt ByTime) Swap(i, j int)      { bt[i], bt[j] = bt[j], bt[i] }
func (bt ByTime) Less(i, j int) bool { return int64(bt[i].EntryStamp) < int64(bt[j].EntryStamp) }

func (c *CBmem) readTimeStamps(f *os.File) (*TimeStamps, error) {
	if c.TimeStampsTable.Addr == 0 {
		return nil, fmt.Errorf("no time stamps")
	}
	var t TSHeader
	a := int64(c.TimeStampsTable.Addr)
	r, err := newOffsetReader(f, a, int(unsafe.Sizeof(t)))
	if err != nil {
		return nil, fmt.Errorf("creating TSHeader offsetReader @ %#x: %w", a, err)
	}
	if err := readOne(r, &t, a); err != nil {
		return nil, fmt.Errorf("failed to read TSTable: %w", err)
	}
	a += int64(unsafe.Sizeof(t))
	stamps := make([]TS, t.NumEntries)
	ts := &TimeStamps{TS: []TS{{EntryID: 0, EntryStamp: t.BaseTime}}, TSHeader: t}
	if r, err = newOffsetReader(f, a, len(stamps)*int(unsafe.Sizeof(stamps[0]))); err != nil {
		return nil, fmt.Errorf("newOffsetReader for %d timestamps: %w", t.NumEntries, err)
	}
	if err := readOne(r, stamps, a); err != nil {
		return nil, fmt.Errorf("failed to read %d stamps: %w", t.NumEntries, err)
	}
	// Timestamps are unsigned. But we're seeing on one machine the first several
	// timestamps have the high order 32 bits set (!). To adjust them, it seems the
	// thing that works is take the first one as the base, and treat it as int64, and
	// subtract that base from each value.
	base := int64(stamps[0].EntryStamp)
	for i := range stamps {
		stamps[i].EntryStamp = uint64(int64(stamps[i].EntryStamp)-base) + t.BaseTime
	}
	ts.TS = append(ts.TS, stamps...)
	sort.Sort(ByTime(ts.TS))
	return ts, nil
}
