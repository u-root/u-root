package ioctl

import "syscall"

// Generic ioctl constants
const (
	IOC_NONE  = 0x0
	IOC_WRITE = 0x1
	IOC_READ  = 0x2

	IOC_NRBITS   = 8
	IOC_TYPEBITS = 8

	IOC_SIZEBITS = 14
	IOC_DIRBITS  = 2

	IOC_NRSHIFT   = 0
	IOC_TYPESHIFT = IOC_NRSHIFT + IOC_NRBITS
	IOC_SIZESHIFT = IOC_TYPESHIFT + IOC_TYPEBITS
	IOC_DIRSHIFT  = IOC_SIZESHIFT + IOC_SIZEBITS

	IOC_NRMASK   = ((1 << IOC_NRBITS) - 1)
	IOC_TYPEMASK = ((1 << IOC_TYPEBITS) - 1)
	IOC_SIZEMASK = ((1 << IOC_SIZEBITS) - 1)
	IOC_DIRMASK  = ((1 << IOC_DIRBITS) - 1)
)

// Some useful additional ioctl constanst
const (
	IOC_IN        = IOC_WRITE << IOC_DIRSHIFT
	IOC_OUT       = IOC_READ << IOC_DIRSHIFT
	IOC_INOUT     = (IOC_WRITE | IOC_READ) << IOC_DIRSHIFT
	IOCSIZE_MASK  = IOC_SIZEMASK << IOC_SIZESHIFT
	IOCSIZE_SHIFT = IOC_SIZESHIFT
)

// IOC generate IOC
func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IOC_DIRSHIFT) | (t << IOC_TYPESHIFT) |
		(nr << IOC_NRSHIFT) | (size << IOC_SIZESHIFT)
}

// IOR generate IOR
func IOR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ, t, nr, size)
}

// IOW generate IOW
func IOW(t, nr, size uintptr) uintptr {
	return IOC(IOC_WRITE, t, nr, size)
}

// IOWR generate IOWR
func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ|IOC_WRITE, t, nr, size)
}

// IO generate IO
func IO(t, nr uintptr) uintptr {
	return IOC(IOC_NONE, t, nr, 0)
}

// IOCTL send ioctl
func IOCTL(fd, name, data uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, name, data)
	if err != 0 {
		return syscall.Errno(err)
	}
	return nil
}
