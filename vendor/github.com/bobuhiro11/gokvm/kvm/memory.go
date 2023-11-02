package kvm

import (
	"errors"
	"syscall"
	"unsafe"
)

// UserSpaceMemoryRegion defines Memory Regions.
type UserspaceMemoryRegion struct {
	Slot          uint32
	Flags         uint32
	GuestPhysAddr uint64
	MemorySize    uint64
	UserspaceAddr uint64
}

// SetMemLogDirtyPages sets region flags to log dirty pages.
// This is useful in many situations, including migration.
func (r *UserspaceMemoryRegion) SetMemLogDirtyPages() {
	r.Flags |= 1 << 0
}

// SetMemReadonly marks a region as read only.
func (r *UserspaceMemoryRegion) SetMemReadonly() {
	r.Flags |= 1 << 1
}

// SetUserMemoryRegion adds a memory region to a vm -- not a vcpu, a vm.
func SetUserMemoryRegion(vmFd uintptr, region *UserspaceMemoryRegion) error {
	_, err := Ioctl(vmFd, IIOW(kvmSetUserMemoryRegion, unsafe.Sizeof(UserspaceMemoryRegion{})),
		uintptr(unsafe.Pointer(region)))

	return err
}

// SetTSSAddr sets the Task Segment Selector for a vm.
func SetTSSAddr(vmFd uintptr, addr uint32) error {
	_, err := Ioctl(vmFd, IIO(kvmSetTSSAddr), uintptr(addr))

	return err
}

// SetIdentityMapAddr sets the address of a 4k-sized-page for a vm.
func SetIdentityMapAddr(vmFd uintptr, addr uint32) error {
	_, err := Ioctl(vmFd, IIOW(kvmSetIdentityMapAddr, 8), uintptr(unsafe.Pointer(&addr)))

	return err
}

type DirtyLog struct {
	Slot   uint32
	_      uint32
	BitMap uint64
}

// GetDirtyLog provides a memory slot and return a bitmap containing any pages dirtied since the
// last call to this ioctl. Bit 0 is the first page in the memory slot.
// Ensure the entire structure is cleared to avoid padding issues.
func GetDirtyLog(vmFd uintptr, dirtlog *DirtyLog) error {
	_, err := Ioctl(vmFd,
		IIOW(kvmGetDirtyLog, unsafe.Sizeof(DirtyLog{})),
		uintptr(unsafe.Pointer(dirtlog)))

	if errors.Is(err, syscall.ENOENT) {
		return nil
	}

	return err
}

func SetNrMMUPages(vmFd uintptr, shadowMem uint64) error {
	_, err := Ioctl(vmFd,
		IIO(kvmSetNrMMUPages),
		uintptr(shadowMem))

	return err
}

func GetNrMMUPages(vmFd uintptr, shadowMem *uint64) error {
	_, err := Ioctl(vmFd,
		IIO(kvmGetNrMMUPages),
		uintptr(unsafe.Pointer(shadowMem)))

	return err
}

type coalescedMMIOZone struct {
	Addr   uint64
	Size   uint32
	PadPio uint32
}

// RegisterCoalescedMMIO registers a address space for Coalesced MMIO.
func RegisterCoalescedMMIO(vmFd uintptr, addr uint64, size uint32) error {
	zone := &coalescedMMIOZone{
		Addr:   addr,
		Size:   size,
		PadPio: 0,
	}

	_, err := Ioctl(vmFd,
		IIOW(kvmResgisterCoalescedMMIO, unsafe.Sizeof(coalescedMMIOZone{})),
		uintptr(unsafe.Pointer(zone)))

	return err
}

// UNregisterCoaloescedMMIO unregister a address space from Coalesced MMIO.
func UnregisterCoalescedMMIO(vmFd uintptr, addr uint64, size uint32) error {
	zone := &coalescedMMIOZone{
		Addr:   addr,
		Size:   size,
		PadPio: 0,
	}

	_, err := Ioctl(vmFd,
		IIOW(kvmUnResgisterCoalescedMMIO, unsafe.Sizeof(coalescedMMIOZone{})),
		uintptr(unsafe.Pointer(zone)))

	return err
}
