// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package txtlog

// BIOSLogID is the legacy eventlog type
type BIOSLogID uint32

const (
	// EvPrebootCert see [2] specification in tcpa_log.go
	EvPrebootCert BIOSLogID = 0x0
	// EvPostCode see [2] specification in tcpa_log.go
	EvPostCode BIOSLogID = 0x1
	// EvUnused see [2] specification in tcpa_log.go
	EvUnused BIOSLogID = 0x2
	// EvNoAction see [2] specification in tcpa_log.go
	EvNoAction BIOSLogID = 0x3
	// EvSeparator see [2] specification in tcpa_log.go
	EvSeparator BIOSLogID = 0x4
	// EvAction see [2] specification in tcpa_log.go
	EvAction BIOSLogID = 0x5
	// EvEventTag see [2] specification in tcpa_log.go
	EvEventTag BIOSLogID = 0x6
	// EvSCRTMContents see [2] specification in tcpa_log.go
	EvSCRTMContents BIOSLogID = 0x7
	// EvSCRTMVersion see [2] specification in tcpa_log.go
	EvSCRTMVersion BIOSLogID = 0x8
	// EvCPUMicrocode see [2] specification in tcpa_log.go
	EvCPUMicrocode BIOSLogID = 0x9
	// EvPlatformConfigFlags see [2] specification in tcpa_log.go
	EvPlatformConfigFlags BIOSLogID = 0xA
	// EvTableOfServices see [2] specification in tcpa_log.go
	EvTableOfServices BIOSLogID = 0xB
	// EvCompactHash see [2] specification in tcpa_log.go
	EvCompactHash BIOSLogID = 0xC
	// EvIPL see [2] specification in tcpa_log.go
	EvIPL BIOSLogID = 0xD
	// EvIPLPartitionData see [2] specification in tcpa_log.go
	EvIPLPartitionData BIOSLogID = 0xE
	// EvNonHostCode see [2] specification in tcpa_log.go
	EvNonHostCode BIOSLogID = 0xF
	// EvNonHostConfig see [2] specification in tcpa_log.go
	EvNonHostConfig BIOSLogID = 0x10
	// EvNonHostInfo see [2] specification in tcpa_log.go
	EvNonHostInfo BIOSLogID = 0x11
	// EvOmitBootDeviceEvents see [2] specification in tcpa_log.go
	EvOmitBootDeviceEvents BIOSLogID = 0x12
)

// BIOSLogTypes are the BIOS eventlog types
var BIOSLogTypes = map[BIOSLogID]string{
	EvPrebootCert:          "EV_PREBOOT_CERT",
	EvPostCode:             "EV_POST_CODE",
	EvUnused:               "EV_UNUSED",
	EvNoAction:             "EV_NO_ACTION",
	EvSeparator:            "EV_SEPARATOR",
	EvAction:               "EV_ACTION",
	EvEventTag:             "EV_EVENT_TAG",
	EvSCRTMContents:        "EV_S_CRTM_CONTENTS",
	EvSCRTMVersion:         "EV_S_CRTM_VERSION",
	EvCPUMicrocode:         "EV_CPU_MICROCODE",
	EvPlatformConfigFlags:  "EV_PLATFORM_CONFIG_FLAGS",
	EvTableOfServices:      "EV_TABLE_OF_DEVICES",
	EvCompactHash:          "EV_COMPACT_HASH",
	EvIPL:                  "EV_IPL",
	EvIPLPartitionData:     "EV_IPL_PARTITION_DATA",
	EvNonHostCode:          "EV_NONHOST_CODE",
	EvNonHostConfig:        "EV_NONHOST_CONFIG",
	EvNonHostInfo:          "EV_NONHOST_INFO",
	EvOmitBootDeviceEvents: "EV_OMIT_BOOT_DEVICE_EVENTS",
}

// EFILogID is the EFI eventlog type
type EFILogID uint32

const (
	// EvEFIEventBase is the base value for all EFI platform
	EvEFIEventBase EFILogID = 0x80000000
	// EvEFIVariableDriverConfig see [1] specification in tcpa_log.go
	EvEFIVariableDriverConfig EFILogID = 0x80000001
	// EvEFIVariableBoot see [1] specification in tcpa_log.go
	EvEFIVariableBoot EFILogID = 0x80000002
	// EvEFIBootServicesApplication see [1] specification in tcpa_log.go
	EvEFIBootServicesApplication EFILogID = 0x80000003
	// EvEFIBootServicesDriver see [1] specification in tcpa_log.go
	EvEFIBootServicesDriver EFILogID = 0x80000004
	// EvEFIRuntimeServicesDriver see [1] specification in tcpa_log.go
	EvEFIRuntimeServicesDriver EFILogID = 0x80000005
	// EvEFIGPTEvent see [1] specification in tcpa_log.go
	EvEFIGPTEvent EFILogID = 0x80000006
	// EvEFIAction see [1] specification in tcpa_log.go
	EvEFIAction EFILogID = 0x80000007
	// EvEFIPlatformFirmwareBlob see [1] specification in tcpa_log.go
	EvEFIPlatformFirmwareBlob EFILogID = 0x80000008
	// EvEFIHandoffTables see [1] specification in tcpa_log.go
	EvEFIHandoffTables EFILogID = 0x80000009
	// EvEFIHCRTMEvent see [1] specification in tcpa_log.go
	EvEFIHCRTMEvent EFILogID = 0x80000010
	// EvEFIVariableAuthority see [1] specification in tcpa_log.go
	EvEFIVariableAuthority EFILogID = 0x800000E0
)

// EFILogTypes are the EFI eventlog types
var EFILogTypes = map[EFILogID]string{
	EvEFIEventBase:               "EV_EFI_EVENT_BASE",
	EvEFIVariableDriverConfig:    "EV_EFI_VARIABLE_DRIVER_CONFIG",
	EvEFIVariableBoot:            "EV_EFI_VARIABLE_BOOT",
	EvEFIBootServicesApplication: "EV_EFI_BOOT_SERVICES_APPLICATION",
	EvEFIBootServicesDriver:      "EV_EFI_BOOT_SERVICES_DRIVER",
	EvEFIRuntimeServicesDriver:   "EV_EFI_RUNTIME_SERVICES_DRIVER",
	EvEFIGPTEvent:                "EV_EFI_GPT_EVENT",
	EvEFIAction:                  "EV_EFI_ACTION",
	EvEFIPlatformFirmwareBlob:    "EV_EFI_PLATFORM_FIRMWARE_BLOB",
	EvEFIHandoffTables:           "EV_EFI_HANDOFF_TABLES",
	EvEFIHCRTMEvent:              "EV_EFI_HCRTM_EVENT",
	EvEFIVariableAuthority:       "EV_EFI_VARIABLE_AUTHORITY",
}

// TCGAgileEventFormatID is the agile eventlog identifier for EV_NO_ACTION events
const TCGAgileEventFormatID string = "Spec ID Event03"

// TCGOldEfiFormatID is the legacy eventlog identifier for EV_NO_ACTION events
const TCGOldEfiFormatID string = "Spec ID Event02"

// HCRTM string for event type EV_EFI_HCRTM_EVENT
const HCRTM string = "HCRTM"

// FirmwareType (BIOS)
type FirmwareType string

const (
	// Uefi is an Open Source UEFI implementation, www.tianocore.org
	Uefi FirmwareType = "UEFI"
	// Coreboot is an Open Source firmware, www.coreboot.org
	Coreboot FirmwareType = "coreboot"
	// UBoot is an Open Source firmware, www.denx.de/wiki/U-Boot
	UBoot FirmwareType = "U-Boot"
	// LinuxBoot is an Open Source firmware based on UEFI and a Linux runtime,
	// www.linuxboot.org
	LinuxBoot FirmwareType = "LinuxBoot"
	// Bios is the legacy BIOS
	Bios FirmwareType = "BIOS"
	// TXT is Intel TXT launch
	Txt FirmwareType = "TXT"
)

// TXT TPM1.2 log container signature
const Txt12EvtLogSignature = "TXT Event Container\000"

// TXT TPM1.2 log versions
const (
	Txt12EvtLog_Cntnr_Major_Ver = 1
	Txt12EvtLog_Cntnr_Minor_Ver = 0
	Txt12EvtLog_Evt_Major_Ver   = 1
	Txt12EvtLog_Evt_Minor_Ver   = 0
)

type TxtLogID uint32

const (
	TxtEvTypeBase TxtLogID = iota + 0x400
	TxtEvTypePcrMapping
	TxtEvTypeHashStart
	TxtEvTypeCombinedHash
	TxtEvTypeMleHash
	TxtEvTypeBiosAcRegData TxtLogID = iota + 0x405
	TxtEvTypeCPUScrtmStat
	TxtEvTypeLcpControlHash
	TxtEvTypeElementsHash
	TxtEvTypeStmHash
	TxtEvTypeOsSinitDataCapHash
	TxtEvTypeSinitPubKeyHash
	TxtEvTypeLcpHash
	TxtEvTypeLcpDetailsHash
	TxtEvTypeLcpAuthoritiesHash
	TxtEvTypeNvInfoHash
	TxtEvTypeColdBootBiosHash
	TxtEvTypeKmHash
	TxtEvTypeBpmHash
	TxtEvTypeKmInfoHash
	TxtEvTypeBpmInfoHash
	TxtEvTypeBootPolHash
	TxtEvTypeRandValue TxtLogID = iota + 0x4e8
	TxtEvTypeCapValue
)

// TxtLogTypes are the Intel TXT eventlog types
var TxtLogTypes = map[TxtLogID]string{
	TxtEvTypeBase:               "EVTYPE_BASE",
	TxtEvTypePcrMapping:         "EVTYPE_PCR_MAPPING",
	TxtEvTypeHashStart:          "EVTYPE_HASH_START",
	TxtEvTypeCombinedHash:       "EVTYPE_COMBINED_HASH",
	TxtEvTypeMleHash:            "EVTYPE_MLE_HASH",
	TxtEvTypeBiosAcRegData:      "EVTYPE_BIOSAC_REG_DATA",
	TxtEvTypeCPUScrtmStat:       "EVTYPE_CPU_SCRTM_STAT",
	TxtEvTypeLcpControlHash:     "EVTYPE_LCP_CONTROL_HASH",
	TxtEvTypeElementsHash:       "EVTYPE_ELEMENTS_HASH",
	TxtEvTypeStmHash:            "EVTYPE_STM_HASH",
	TxtEvTypeOsSinitDataCapHash: "EVTYPE_OSSINITDATA_CAP_HASH",
	TxtEvTypeSinitPubKeyHash:    "EVTYPE_SINIT_PUBKEY_ HASH",
	TxtEvTypeLcpHash:            "EVTYPE_LCP_HASH",
	TxtEvTypeLcpDetailsHash:     "EVTYPE_LCP_DETAILS_HASH",
	TxtEvTypeLcpAuthoritiesHash: "EVTYPE_LCP_AUTHORITIES_HASH",
	TxtEvTypeNvInfoHash:         "EVTYPE_NV_INFO_HASH",
	TxtEvTypeColdBootBiosHash:   "EVTYPE_COLD_BOOT_BIOS_HASH",
	TxtEvTypeKmHash:             "EVTYPE_KM_HASH",
	TxtEvTypeBpmHash:            "EVTYPE_KM_HASH",
	TxtEvTypeKmInfoHash:         "EVTYPE_KM_INFO_HASH",
	TxtEvTypeBpmInfoHash:        "EVTYPE_BPM_INFO_HASH",
	TxtEvTypeBootPolHash:        "EVTYPE_BOOT_POL_HASH",
	TxtEvTypeRandValue:          "EVTYPE_RANDOM_VALUE",
	TxtEvTypeCapValue:           "EVTYPE_CAP_VALUE",
}
