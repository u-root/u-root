package tap

import (
	"fmt"
	"syscall"
	"unsafe"
)

const ifNameSize = 0x10

type Tap struct {
	fd int
}

type ifReq struct {
	Name  [ifNameSize]byte
	Flags uint16
	_     [0x28 - ifNameSize - 2]byte
}

func ioctl(fd, op, arg uintptr) (uintptr, error) {
	res, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, fd, op, arg)
	if errno != 0 {
		return res, errno
	}

	return res, nil
}

func fcntl(fd, op, arg uintptr) (uintptr, error) {
	res, _, errno := syscall.Syscall(
		syscall.SYS_FCNTL, fd, op, arg)

	if errno != 0 {
		return res, errno
	}

	return res, nil
}

func New(name string) (*Tap, error) {
	var err error

	t := &Tap{}

	if t.fd, err = syscall.Open("/dev/net/tun", syscall.O_RDWR, 0); err != nil {
		return t, fmt.Errorf("/dev/net/tun: %w", err)
	}

	ifr := ifReq{
		Name:  [ifNameSize]byte{},
		Flags: syscall.IFF_TAP | syscall.IFF_NO_PI,
	}
	copy(ifr.Name[:ifNameSize-1], name)

	ifrPtr := uintptr(unsafe.Pointer(&ifr))
	if _, err = ioctl(uintptr(t.fd), syscall.TUNSETIFF, ifrPtr); err != nil {
		return t, fmt.Errorf("TUN TUNSETIFF: %w", err)
	}

	// issue SIGIO if this tap interface receive packets
	if _, err = fcntl(uintptr(t.fd), syscall.F_SETSIG, 0); err != nil {
		return t, fmt.Errorf("tun SETSIG: %w", err)
	}

	var flags uintptr

	// enable non-blocking IO for tap interface
	if flags, err = fcntl(uintptr(t.fd), syscall.F_GETFL, 0); err != nil {
		return t, fmt.Errorf("TUN GETFL: %w", err)
	}

	flags |= syscall.O_NONBLOCK | syscall.O_ASYNC
	if _, err = fcntl(uintptr(t.fd), syscall.F_SETFL, flags); err != nil {
		return t, fmt.Errorf("TUN SETFL NONBLOCK|ASYNC: %w", err)
	}

	return t, nil
}

func (t *Tap) Close() error {
	return syscall.Close(t.fd)
}

func (t Tap) Write(buf []byte) (n int, err error) {
	return syscall.Write(t.fd, buf)
}

func (t Tap) Read(buf []byte) (n int, err error) {
	return syscall.Read(t.fd, buf)
}
