// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mei implements a wrapper on top of Linux's MEI (Intel ME Interface,
// formerly known as HECI). This module requires Linux, and the `mei_me` driver.
// Once loaded, this driver will expose a `/dev/mei0` device, that can be
// accessed through this library.
package mei

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/vtolstov/go-ioctl"
)

// DefaultMEIDevicePath is the path of the default MEI device. This file will be
// present if the "mei_me" kernel module is loaded.
var DefaultMEIDevicePath = "/dev/mei0"

// HECIGuids maps the known HECI GUIDs to their values. The MEI interface wants
// little-endian. See all the GUIDs at
// https://github.com/intel/lms/blob/master/MEIClient/Include/HECI_if.h
var (
	// "8e6a6715-9abc-4043-88ef-9e39c6f63e0f"
	MKHIGuid = ClientGUID{0x15, 0x67, 0x6a, 0x8e, 0xbc, 0x9a, 0x43, 0x40, 0x88, 0xef, 0x9e, 0x39, 0xc6, 0xf6, 0x3e, 0xf}
)

// see include/uapi/linux/mei.h
var (
	IoctlMEIConnectClient = ioctl.IOWR('H', 0x01, uintptr(len(ClientGUID{})))
)

// ClientGUID is the data buffer to pass to `ioctl` to connect to
// MEI. See include/uapi/linux/mei.h .
type ClientGUID [16]byte

// ClientProperties is the data buffer returned by `ioctl` after connecting to
// MEI. See include/uapi/linux/mei.h .
type ClientProperties [6]byte

// MaxMsgLength is the maximum size of a message for this client.
func (c ClientProperties) MaxMsgLength() uint32 {
	return binary.LittleEndian.Uint32(c[:4])
}

// ProtocolVersion is this client's protocol version.
func (c ClientProperties) ProtocolVersion() uint8 {
	return c[4]
}

// MEI represents an Intel ME Interface object.
type MEI struct {
	fd               *int
	ClientProperties ClientProperties
}

// OpenMEI opens the specified MEI device, using the client type defined by GUID.
// See `HECIGuids` in this package.
func OpenMEI(meiPath string, guid ClientGUID) (*MEI, error) {
	var m MEI
	fd, err := syscall.Open(meiPath, os.O_RDWR, 0o755)
	if err != nil {
		return nil, err
	}
	data := [16]byte(guid)
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), IoctlMEIConnectClient, uintptr(unsafe.Pointer(&data))); err != 0 {
		return nil, fmt.Errorf("ioctl IOCTL_MEI_CONNECT_CLIENT failed: %w", err)
	}
	// can be racy, unless protected by a mutex
	m.fd = &fd
	copy(m.ClientProperties[:], data[:])
	return &m, nil
}

// Close closes the MEI device, if open, and does nothing otherwise.
func (m *MEI) Close() error {
	if m.fd != nil {
		err := syscall.Close(*m.fd)
		m.fd = nil
		return err
	}
	return nil
}

// Write writes to the MEI file descriptor.
func (m *MEI) Write(p []byte) (int, error) {
	// use syscall.Write instead of m.fd.Write to avoid epoll
	return syscall.Write(*m.fd, p)
}

// Read reads from the MEI file descriptor.
func (m *MEI) Read(p []byte) (int, error) {
	// use syscall.Read instead of m.fd.Read to avoid epoll
	return syscall.Read(*m.fd, p)
}
