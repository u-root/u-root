// +build amd64

// Package hwapi provides access to low level hardware
package hwapi

import "github.com/intel-go/cpuid"

func cpuidLow(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32) // implemented in cpuidlow_amd64.s

//VersionString returns the vendor ID
func (t TxtAPI) VersionString() string {
	return cpuid.VendorIdentificatorString
}

//HasSMX returns true if SMX is supported
func (t TxtAPI) HasSMX() bool {
	return cpuid.HasFeature(cpuid.SMX)
}

//HasVMX returns true if VMX is supported
func (t TxtAPI) HasVMX() bool {
	return cpuid.HasFeature(cpuid.VMX)
}

//HasMTRR returns true if MTRR are supported
func (t TxtAPI) HasMTRR() bool {
	return cpuid.HasFeature(cpuid.MTRR) || cpuid.HasExtraFeature(cpuid.MTRR_2)
}

//ProcessorBrandName returns the CPU brand name
func (t TxtAPI) ProcessorBrandName() string {
	return cpuid.ProcessorBrandString
}

//CPUSignature returns CPUID=1 eax
func (t TxtAPI) CPUSignature() uint32 {
	eax, _, _, _ := cpuidLow(1, 0)
	return eax
}

//CPULogCount returns number of logical CPU cores
func (t TxtAPI) CPULogCount() uint32 {
	return uint32(cpuid.MaxLogocalCPUId)
}
