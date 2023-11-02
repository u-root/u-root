package kvm

import (
	"unsafe"
)

type X86MCE struct {
	Status    uint64
	Addr      uint64
	Misc      uint64
	MCGStatus uint64
	Bank      uint8
	_         [7]uint8
	_         [3]uint64
}

// X86SetupMCE initializes MCE support for the given vcpu.
func X86SetupMCE(vcpuFd uintptr, mceCap *uint64) error {
	_, err := Ioctl(vcpuFd,
		IIOW(kvmX86SetupMCE, unsafe.Sizeof(mceCap)),
		uintptr(unsafe.Pointer(mceCap)))

	return err
}

// X86GetMCECapSupported returns supported MCE capabilities.
func X86GetMCECapSupported(kvmFd uintptr, mceCap *uint64) error {
	_, err := Ioctl(kvmFd,
		IIOR(kvmX86GetMCECapSupported, unsafe.Sizeof(mceCap)),
		uintptr(unsafe.Pointer(mceCap)))

	return err
}
