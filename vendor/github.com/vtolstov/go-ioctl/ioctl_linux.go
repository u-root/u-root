// +build linux

package ioctl

import "unsafe"

// BlkPart send blkpart ioctl to fd
func BlkRRPart(fd uintptr) error {
	return IOCTL(fd, IO(0x12, 95), uintptr(0))
}

// BlkPg send blkpg ioctl to fd
func BlkPg(fd, data uintptr) error {
	return IOCTL(fd, IO(0x12, 105), data)
	return nil
}

type FsTrimRange struct {
	Start     uint64
	Length    uint64
	MinLength uint64
}

// Fitrim send fitrim ioctl to fd
func Fitrim(fd, data uintptr) error {
	r := FsTrimRange{}
	return IOCTL(fd, IOWR('X', 121, uintptr(unsafe.Pointer(&r))), data)
}

// Fifreeze send fifreeze ioctl to fd
func Fifreeze(fd, data uintptr) error {
	return IOCTL(fd, IOWR('X', 119, uintptr(0)), data)
}

// Fithaw send fithaw ioctl to fd
func Fithaw(fd, data uintptr) error {
	return IOCTL(fd, IOWR('X', 120, uintptr(0)), data)
}
