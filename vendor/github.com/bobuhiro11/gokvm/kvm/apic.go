package kvm

import (
	"unsafe"
)

type TRPAccessCtl struct {
	Enable uint32
	Flags  uint32
	_      [8]uint32
}

// TRPAccessReporting sets status of Trap Access Reporting for APIC.
func TRPAccessReporting(vcpuFd uintptr, ctl *TRPAccessCtl) error {
	_, err := Ioctl(vcpuFd,
		IIOWR(kvmTRPAccessReporting, unsafe.Sizeof(TRPAccessCtl{})),
		uintptr(unsafe.Pointer(ctl)))

	return err
}
