package test

import (
	"fmt"

	"github.com/9elements/converged-security-suite/pkg/hwapi"
)

func notImplemented(txtAPI hwapi.APIInterfaces) (bool, error, error) {
	return false, nil, fmt.Errorf("Not implemented")
}

var (
	testRSDPChecksum = Test{
		Name:                    "ACPI RSDP has valid checksum",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 1",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testRSDTChecksum = Test{
		Name:                    "ACPI RSDT present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 2",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testRSDTValid = Test{
		Name:                    "ACPI RSDT is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 3",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testDMARPresent = Test{
		Name:                    "ACPI DMAR is present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 4",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testDMARValid = Test{
		Name:                    "ACPI DMAR is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 5",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testMADTPresent = Test{
		Name:                    "ACPI MADT is present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 16",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testMADTValid = Test{
		Name:                    "ACPI MADT is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 7",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testRSDPValid = Test{
		Name:                    "ACPI RSDP is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 8",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testXSDTValid = Test{
		Name:                    "ACPI XSDT is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 9",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}

	testTXTHeapSizeFitsMADTCopy = Test{
		Name:                    "ACPI MADT copy fits into TXT heap",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 9 Major 7 Minor 1",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testTXTHeapSizeFitsDynamicMadt = Test{
		Name:                    "Dynamic ACPI MADT fits into TXT heap",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 9 Major 7 Minor 2",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testTXTHeapSizeFitsDMARCopy = Test{
		Name:                    "ACPI DMAR copy fits into TXT heap",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 9 Major 7 Minor 3",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIRSDPInOSToSINITData = Test{
		Name:                    "ACPI RSDP in 'OS to SINIT data' points to address below 4 GiB",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 9 Major 0xc",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARValidHPET = Test{
		Name:                    "ACPI DMAR table has valid HPET configuration",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 1",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARValidBus = Test{
		Name:                    "ACPI DMAR table has valid BUS configuration",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 2",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARValidAzalia = Test{
		Name:                    "ACPI DMAR table Azalia device scope is valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 3",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDeviceScopePresent = Test{
		Name:                    "ACPI DMAR table device scope is present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 4",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARHPETScopeDuplicated = Test{
		Name:                    "ACPI DMAR table has no duplicated HPET scope",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 5",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdVtdDevice = Test{
		Name:                    "ACPI DMAR table DRHD device ",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 6",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdVtdScope = Test{
		Name:                    "ACPI DMAR table DRHD device scope",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 7",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdPchApic = Test{
		Name:                    "ACPI DMAR table DRHD PCH APIC present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 8",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdBaseaddressBelowFourGiB = Test{
		Name:                    "ACPI DMAR table DRHD base address below 4 GiB",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 9",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdTopaddressBelowFourGiB = Test{
		Name:                    "ACPI DMAR table DRHD top address below 4 GiB",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 0xa",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdBadDevicescopeEntry = Test{
		Name:                    "ACPI DMAR table DRHD device scope entries are valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 0xb",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}
	testACPIDMARDrhdBadDevicescopeLength = Test{
		Name:                    "ACPI DMAR table DRHD device scope length are valid",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xA Major 3 Minor 0xc",
		SpecificiationTitle:     ServerGrantleyPlatformSpecificationTitle,
		SpecificationDocumentID: ServerGrantleyPlatformDocumentID,
	}

	testACPIPWRMBarBelowFourGib = Test{
		Name:                    "ACPI PWRM BAR is below 4 GiB",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0x35 Major 4",
		SpecificiationTitle:     CBtGTXTPlatformSpecificationTitle,
		SpecificationDocumentID: CBtGTXTPlatformDocumentID,
	}
	testMCFGPresent = Test{
		Name:                    "ACPI MCFG is present",
		Required:                true,
		function:                notImplemented,
		Status:                  NotImplemented,
		SpecificationChapter:    "SINIT Class 0xC Major 0xa",
		SpecificiationTitle:     CBtGTXTPlatformSpecificationTitle,
		SpecificationDocumentID: CBtGTXTPlatformDocumentID,
	}
)
