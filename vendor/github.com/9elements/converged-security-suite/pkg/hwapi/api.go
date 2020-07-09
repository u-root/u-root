package hwapi

import (
	"github.com/9elements/go-tss"
)

//APIInterfaces provides methods to access hardware found on modern x86_64 platforms
type APIInterfaces interface {
	// cpu_whitelist.go - cpu_blacklist.go
	CPUBlacklistTXTSupport() bool
	CPUWhitelistTXTSupport() bool

	// cpuid.go
	VersionString() string
	HasSMX() bool
	HasVMX() bool
	HasMTRR() bool
	ProcessorBrandName() string
	CPUSignature() uint32
	CPULogCount() uint32

	// e820.go
	IsReservedInE820(start uint64, end uint64) (bool, error)

	// iommu.go
	LookupIOAddress(addr uint64, regs VTdRegisters) ([]uint64, error)
	AddressRangesIsDMAProtected(first, end uint64) (bool, error)

	// msr.go
	HasSMRR() (bool, error)
	GetSMRRInfo() (SMRR, error)
	IA32FeatureControlIsLocked() (bool, error)
	IA32PlatformID() (uint64, error)
	AllowsVMXInSMX() (bool, error)
	TXTLeavesAreEnabled() (bool, error)
	IA32DebugInterfaceEnabledOrLocked() (bool, bool, bool, error)

	// pci.go
	PCIReadConfigSpace(bus int, device int, devFn int, off int, buf interface{}) error
	PCIReadConfig16(bus int, device int, devFn int, off int) (uint16, error)
	PCIReadConfig32(bus int, device int, devFn int, off int) (uint32, error)
	PCIReadVendorID(bus int, device int, devFn int) (uint16, error)
	PCIReadDeviceID(bus int, device int, devFn int) (uint16, error)
	ReadHostBridgeTseg() (uint32, uint32, error)
	ReadHostBridgeDPR() (DMAProtectedRange, error)

	// phys.go
	ReadPhys(addr int64, data UintN) error
	ReadPhysBuf(addr int64, buf []byte) error
	WritePhys(addr int64, data UintN) error

	// tpm.go
	NewTPM() (*tss.TPM, error)
	NVLocked(tpmCon *tss.TPM) (bool, error)
	ReadNVPublic(tpmCon *tss.TPM, index uint32) ([]byte, error)
	NVReadValue(tpmCon *tss.TPM, index uint32, password string, size, offhandle uint32) ([]byte, error)
	ReadPCR(tpmCon *tss.TPM, pcr uint32) ([]byte, error)
}

//TxtAPI The context object for TXT Api
type TxtAPI struct{}

//GetAPI Returns an initialized TxtApi object
func GetAPI() APIInterfaces {
	return TxtAPI{}
}
