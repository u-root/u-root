package tpm

// TPMMaxPCRListSize is the maximum number of PCRs for a TPM
const TPMMaxPCRListSize = 24

// IAlgHash is the TPM hash algorithm
type IAlgHash uint16

// We only define TPM hash algorithms here we use
const (
	// TPMAlgError is an algorithm error
	TPMAlgError IAlgHash = 0x0000
	// TPMAlgSha
	TPMAlgSha     IAlgHash = 0x0004
	TPMAlgSha256  IAlgHash = 0x000B
	TPMAlgSha384  IAlgHash = 0x000C
	TPMAlgSha512  IAlgHash = 0x000D
	TPMAlgSm3s256 IAlgHash = 0x0012
)

// IAlgHashSize is the TPM hash algorithm length
type IAlgHashSize uint8

const (
	// TPMAlgShaSize SHA hash size
	TPMAlgShaSize IAlgHashSize = 20
	// TPMAlgSha256Size SHA256 hash size
	TPMAlgSha256Size IAlgHashSize = 32
	// TPMAlgSha384Size SHA384 hash size
	TPMAlgSha384Size IAlgHashSize = 48
	// TPMAlgSha512Size SHA512 hash size
	TPMAlgSha512Size IAlgHashSize = 64
	// TPMAlgSm3s256Size SM3-256 hash size
	TPMAlgSm3s256Size IAlgHashSize = 32
)

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
