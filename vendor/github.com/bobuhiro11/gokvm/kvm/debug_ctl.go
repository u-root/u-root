package kvm

import "unsafe"

// debugControl controls guest debug.
type debugControl struct { // nolint:unused
	Control  uint32
	_        uint32
	DebugReg [8]uint64
}

// SingleStep enables single stepping on all vCPUS in a VM. At present, it seems
// not to work. It is based on working code from the linuxboot Voodoo project.
func SingleStep(vmFd uintptr, onoff bool) error {
	const (
		// Enable enables debug options in the guest
		Enable = 1
		// SingleStep enables single step.
		SingleStep = 2
	)

	var (
		debug         [unsafe.Sizeof(debugControl{})]byte
		setGuestDebug = IIOW(0x9b, unsafe.Sizeof(debugControl{}))
	)

	if onoff {
		debug[2] = 0x0002 // 0000
		debug[0] = Enable | SingleStep
	}

	// this is not very nice, but it is easy.
	// And TBH, the tricks the Linux kernel people
	// play are a lot nastier.
	_, err := Ioctl(vmFd, setGuestDebug, uintptr(unsafe.Pointer(&debug[0])))

	return err
}
