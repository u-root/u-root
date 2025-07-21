// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fbpt reads Firmware Basic Performance Table within ACPI FPDT Table.
package fbpt

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/acpi/fpdt"
	"github.com/u-root/u-root/pkg/uefivars"
)

const (
	fbptStructureSig = "FBPT"
	memDevice        = "/dev/mem"

	// definition can be found at: https://uefi.org/sites/default/files/resources/ACPI%206_2_A_Sept29.pdf (page 208/page 212)
	efiAcpi5_0FpdtPerformanceRecordHeaderSize = 4
	efiAcpi5_0FbptHeaderSize                  = 8

	// maximum number of FBPTPerfRecords to return in 'FindAllFBPTRecords'
	maxNumberOfFBPTPerfRecords = 2000

	fpdtDynamicStringEventRecordIdentifier = 0x1011
	// see ACPI Table Spec: https://uefi.org/sites/default/files/resources/ACPI%206_2_A_Sept29.pdf (page 213)
	efiAcpi6_5FpdtFirmwareBasicBootRecordIdentifier = 0x2

	moduleStartID          = 0x01
	moduleEndID            = 0x02
	moduleLoadImageStartID = 0x03
	moduleLoadImageEndID   = 0x04
	moduleDBStartID        = 0x05
	moduleDBEndID          = 0x06
	moduleDBSupportStartID = 0x07
	moduleDBSupportEndID   = 0x08
	moduleDBStopStartID    = 0x09
	moduleDBStopEndID      = 0x0A

	perfEventSignalStartID = 0x10
	perfEventSignalEndID   = 0x11
	perfCallbackStartID    = 0x20
	perfCallbackEndID      = 0x21
	perfFunctionStartID    = 0x30
	perfFunctionEndID      = 0x31
	perfInmoduleStartID    = 0x40
	perfInmoduleEndID      = 0x41
	perfCrossmoduleStartID = 0x50
	perfCrossmoduleEndID   = 0x51
)

var eventTypeMap = map[uint16]string{
	moduleStartID:          "MODULE_START_ID",
	moduleEndID:            "MODULE_END_ID",
	moduleLoadImageStartID: "MODULE_LOADIMAGE_START_ID",
	moduleLoadImageEndID:   "MODULE_LOADIMAGE_END_ID",
	moduleDBStartID:        "MODULE_DB_START_ID",
	moduleDBEndID:          "MODULE_DB_END_ID",
	moduleDBSupportStartID: "MODULE_DB_SUPPORT_START_ID",
	moduleDBSupportEndID:   "MODULE_DB_SUPPORT_END_ID",
	moduleDBStopStartID:    "MODULE_DB_STOP_START_ID",
	moduleDBStopEndID:      "MODULE_DB_STOP_END_ID",

	perfEventSignalStartID: "PERF_EVENTSIGNAL_START_ID",
	perfEventSignalEndID:   "PERF_EVENTSIGNAL_END_ID",
	perfCallbackStartID:    "PERF_CALLBACK_START_ID",
	perfCallbackEndID:      "PERF_CALLBACK_END_ID",
	perfFunctionStartID:    "PERF_FUNCTION_START_ID",
	perfFunctionEndID:      "PERF_FUNCTION_END_ID",
	perfInmoduleStartID:    "PERF_INMODULE_START_ID",
	perfInmoduleEndID:      "PERF_INMODULE_END_ID",
	perfCrossmoduleStartID: "PERF_CROSSMODULE_START_ID",
	perfCrossmoduleEndID:   "PERF_CROSSMODULE_END_ID",
}

// based on struct definition found in edk2: /MdePkg/Include/IndustryStandard/Acpi50.h
type efiAcpi5_0FpdtPerformanceRecordHeader struct {
	Type     uint16
	Length   uint8
	Revision uint8
}

// EfiAcpi6_5FpdtFirmwareBasicBootRecord represents the
// the Firmware Basic Boot Record always found within FBPT
// based on struct definition found in edk2: /MdeModulePkg/Include/Guid/ExtendedFirmwarePerformance.h
type EfiAcpi6_5FpdtFirmwareBasicBootRecord struct {
	PerformanceRecordHeader efiAcpi5_0FpdtPerformanceRecordHeader
	ResetEnd                uint64
	OSLoaderLoadImageStart  uint64
	OSLoaderStartImageStart uint64
	ExitBootServicesEntry   uint64
	ExitBootServicesExit    uint64
}

// MeasurementRecord represents
// all the different Measurement Entries within FBPT
type MeasurementRecord struct {
	HookType            string
	ProcessorIdentifier uint32
	Timestamp           uint64
	GUID                uefivars.MixedGUID
	Description         string
}

func verifyFBPTSignature(mem io.ReadSeeker, fbptAddr uint64) (uint32, error) {
	// Read & confirm FBPT struct signature
	if _, err := mem.Seek(int64(fbptAddr), io.SeekStart); err != nil {
		return 0, err
	}
	// Read as slices
	var fbptSig [4]byte
	if _, err := io.ReadFull(mem, fbptSig[:]); err != nil {
		return 0, err
	}

	if string(fbptSig[:]) != fbptStructureSig {
		return 0, errors.New("FBPT structure signature check failed. Expected: FBPT, Got: " + string(fbptSig[:]))
	}

	var fbptLength [4]byte
	if _, err := io.ReadFull(mem, fbptLength[:]); err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(fbptLength[:]), nil
}

// FindAllFBPTRecords returns all FBPTRecords of type FPDT Dynamic String
// Record Type and the Firmware Basic Boot Record within FPBT.
func FindAllFBPTRecords(FBPTAddr uint64) (int, []MeasurementRecord, EfiAcpi6_5FpdtFirmwareBasicBootRecord, error) {
	var f *os.File
	var basicBootRecord EfiAcpi6_5FpdtFirmwareBasicBootRecord
	var err error
	if f, err = os.OpenFile(memDevice, os.O_RDONLY, 0); err != nil {
		return 0, nil, basicBootRecord, err
	}
	defer f.Close()

	var tablelength uint32
	if tablelength, err = verifyFBPTSignature(f, FBPTAddr); err != nil {
		return 0, nil, basicBootRecord, err
	}

	// iterate through FBPT table
	measurementRecords := make([]MeasurementRecord, maxNumberOfFBPTPerfRecords)
	var index int
	var tableBytesRead uint32
	var HeaderInfo efiAcpi5_0FpdtPerformanceRecordHeader
	for tableBytesRead < (tablelength-efiAcpi5_0FbptHeaderSize) && index < maxNumberOfFBPTPerfRecords {
		if HeaderInfo.Type, HeaderInfo.Length, _, err = fpdt.ReadFPDTRecordHeader(f); err != nil {
			return index, nil, basicBootRecord, err
		}
		if HeaderInfo.Type == fpdtDynamicStringEventRecordIdentifier {
			if measurementRecords[index], err = readFirmwarePerformanceDataTableDynamicRecord(f, HeaderInfo.Length); err != nil {
				return index, nil, basicBootRecord, err
			}
			index++
		} else if HeaderInfo.Type == efiAcpi6_5FpdtFirmwareBasicBootRecordIdentifier {
			// skip reserved section before Firmware Basic Boot Performance Record
			if _, err := f.Seek(4, io.SeekCurrent); err != nil {
				return index, nil, basicBootRecord, err
			}

			if basicBootRecord, err = readFirmwareBasicBootPerformanceRecord(f); err != nil {
				return index, nil, basicBootRecord, err
			}
		} else {
			if _, err := f.Seek(int64(HeaderInfo.Length-efiAcpi5_0FpdtPerformanceRecordHeaderSize), io.SeekCurrent); err != nil {
				return index, nil, basicBootRecord, err
			}
		}
		tableBytesRead += uint32(HeaderInfo.Length)
	}

	return index, measurementRecords, basicBootRecord, nil
}

func readFirmwarePerformanceDataTableDynamicRecord(mem io.ReadSeeker, recordLength uint8) (MeasurementRecord, error) {
	var measurementRecord MeasurementRecord
	var HookType [2]byte
	if _, err := io.ReadFull(mem, HookType[:]); err != nil {
		return measurementRecord, err
	}

	var ProcessorIdentifier [4]byte
	if _, err := io.ReadFull(mem, ProcessorIdentifier[:]); err != nil {
		return measurementRecord, err
	}

	var Timestamp [8]byte
	if _, err := io.ReadFull(mem, Timestamp[:]); err != nil {
		return measurementRecord, err
	}

	var GUID [16]byte
	if _, err := io.ReadFull(mem, GUID[:]); err != nil {
		return measurementRecord, err
	}

	String := make([]byte, recordLength-34)
	if _, err := io.ReadFull(mem, String[:]); err != nil {
		return measurementRecord, err
	}

	measurementRecord.HookType = eventTypeMap[binary.LittleEndian.Uint16(HookType[:])]
	measurementRecord.ProcessorIdentifier = binary.LittleEndian.Uint32(ProcessorIdentifier[:])
	measurementRecord.Timestamp = binary.LittleEndian.Uint64(Timestamp[:])
	measurementRecord.GUID = uefivars.MixedGUID(GUID)
	measurementRecord.Description = string(String[:])

	return measurementRecord, nil
}

func readFirmwareBasicBootPerformanceRecord(mem io.ReadSeeker) (EfiAcpi6_5FpdtFirmwareBasicBootRecord, error) {
	var basicBootRecord EfiAcpi6_5FpdtFirmwareBasicBootRecord
	var ResetEnd [8]byte
	if _, err := io.ReadFull(mem, ResetEnd[:]); err != nil {
		return basicBootRecord, err
	}

	var OSLoaderLoadImageStart [8]byte
	if _, err := io.ReadFull(mem, OSLoaderLoadImageStart[:]); err != nil {
		return basicBootRecord, err
	}

	var OSLoaderStartImageStart [8]byte
	if _, err := io.ReadFull(mem, OSLoaderStartImageStart[:]); err != nil {
		return basicBootRecord, err
	}

	var ExitBootServicesEntry [8]byte
	if _, err := io.ReadFull(mem, ExitBootServicesEntry[:]); err != nil {
		return basicBootRecord, err
	}

	var ExitBootServicesExit [8]byte
	if _, err := io.ReadFull(mem, ExitBootServicesExit[:]); err != nil {
		return basicBootRecord, err
	}

	basicBootRecord.ResetEnd = binary.LittleEndian.Uint64(ResetEnd[:])
	basicBootRecord.OSLoaderLoadImageStart = binary.LittleEndian.Uint64(OSLoaderLoadImageStart[:])
	basicBootRecord.OSLoaderStartImageStart = binary.LittleEndian.Uint64(OSLoaderStartImageStart[:])
	basicBootRecord.ExitBootServicesEntry = binary.LittleEndian.Uint64(ExitBootServicesEntry[:])
	basicBootRecord.ExitBootServicesExit = binary.LittleEndian.Uint64(ExitBootServicesExit[:])

	return basicBootRecord, nil
}
