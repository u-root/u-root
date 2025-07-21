// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fbptcat dumps the contents of FBPT Table
// within the ACPI FPDT Table
// FBPT stands for Firmware Basic Performance Table

package main

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/acpi/fpdt"
	"github.com/u-root/u-root/pkg/acpi/fpdt/fbpt"
)

func main() {
	// Get FPDT table from ACPI
	var acpiFPDT acpi.Table
	var err error
	if acpiFPDT, err = fpdt.ReadACPIFPDTTable(); err != nil {
		log.Fatal(err)
	}

	// Get FBPT Pointer from FPDT Table
	var FBPTAddr uint64
	if FBPTAddr, err = fpdt.FindFBPTTableAdrr(acpiFPDT); err != nil {
		log.Fatal(err)
	}

	var basicBootRecord fbpt.EfiAcpi6_5FpdtFirmwareBasicBootRecord
	var measurementRecords []fbpt.MeasurementRecord
	if _, measurementRecords, basicBootRecord, err = fbpt.FindAllFBPTRecords(FBPTAddr); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ResetEnd: %d, OSLoaderLoadImageStart: %d, OsLoaderStartImageStart: %d, ExitBootServicesEntry: %d, ExitBootServicesExit: %d \n", basicBootRecord.ResetEnd, basicBootRecord.OSLoaderLoadImageStart, basicBootRecord.OSLoaderStartImageStart, basicBootRecord.ExitBootServicesEntry, basicBootRecord.ExitBootServicesExit)

	for i, measurementRecord := range measurementRecords {
		if measurementRecord.Timestamp == 0 && len(measurementRecord.HookType) == 0 && len(measurementRecord.Description) == 0 {
			continue
		}
		fmt.Printf("Index: %d,Hook Type: %s, Processor Identifier/APIC ID: %d, Timestamp: %d, Guid: %s, Description: %s\n", i, measurementRecord.HookType, measurementRecord.ProcessorIdentifier, measurementRecord.Timestamp, measurementRecord.GUID.String(), measurementRecord.Description)
	}
}
