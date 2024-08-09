// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package scuzz

import (
	"testing"
	"unsafe"
)

// check checks the packetHeader, cdb, and sb.
// The packetHeader check is a bit complex as it contains
// pointers. The pointer values for the original executiion of hdparm
// will not match the values we get for this test. We therefore skip any uintptrs.
func check(t *testing.T, got *packet, want *packet) {
	if got.interfaceID != want.interfaceID {
		t.Errorf("interfaceID: got %v, want %v", got.interfaceID, want.interfaceID)
	}
	if got.direction != want.direction {
		t.Errorf("direction: got %v, want %v", got.direction, want.direction)
	}
	if got.cmdLen != want.cmdLen {
		t.Errorf("cmdLen: got %v, want %v", got.cmdLen, want.cmdLen)
	}
	if got.maxStatusBlockLen != want.maxStatusBlockLen {
		t.Errorf("maxStatusBlockLen: got %v, want %v", got.maxStatusBlockLen, want.maxStatusBlockLen)
	}
	if got.iovCount != want.iovCount {
		t.Errorf("iovCount: got %v, want %v", got.iovCount, want.iovCount)
	}
	if got.dataLen != want.dataLen {
		t.Errorf("dataLen: got %v, want %v", got.dataLen, want.dataLen)
	}
	if got.timeout != want.timeout {
		t.Errorf("timeout: got %v, want %v", got.timeout, want.timeout)
	}
	if got.flags != want.flags {
		t.Errorf("flags: got %v, want %v", got.flags, want.flags)
	}
	if got.packID != want.packID {
		t.Errorf("packID: got %v, want %v", got.packID, want.packID)
	}
	if got.status != want.status {
		t.Errorf("status: got %v, want %v", got.status, want.status)
	}
	if got.maskedStatus != want.maskedStatus {
		t.Errorf("maskedStatus: got %v, want %v", got.maskedStatus, want.maskedStatus)
	}
	if got.msgStatus != want.msgStatus {
		t.Errorf("msgStatus: got %v, want %v", got.msgStatus, want.msgStatus)
	}
	if got.sbLen != want.sbLen {
		t.Errorf("sbLen: got %v, want %v", got.sbLen, want.sbLen)
	}
	if got.hostStatus != want.hostStatus {
		t.Errorf("hostStatus: got %v, want %v", got.hostStatus, want.hostStatus)
	}
	if got.driverStatus != want.driverStatus {
		t.Errorf("driverStatus: got %v, want %v", got.driverStatus, want.driverStatus)
	}
	if got.resID != want.resID {
		t.Errorf("resID: got %v, want %v", got.resID, want.resID)
	}
	if got.duration != want.duration {
		t.Errorf("duration: got %v, want %v", got.duration, want.duration)
	}
	if got.info != want.info {
		t.Errorf("info: got %v, want %v", got.info, want.info)
	}

	for i := range got.command {
		if got.command[i] != want.command[i] {
			t.Errorf("command[%d]: got %#02x, want %#02x", i, got.command[i], want.command[i])
		}
	}

	for i := range got.block {
		if got.block[i] != want.block[i] {
			t.Errorf("cblock[%d]: got %#02x, want %#02x", i, got.block[i], want.block[i])
		}
	}
}

// TestSizes makes sure that everything marshals to the right size.
// The sizes are magic numbers from Linux. Thanks to the compatibility
// guarantee, we know they don't change.
func TestSizes(t *testing.T) {
	hs := unsafe.Sizeof(packetHeader{})
	if hs != hdrSize {
		t.Errorf("PacketHeader.Marshal(): got %d, want %d", hs, hdrSize)
	}
	l := len(&commandDataBlock{})
	if l != cdbSize {
		t.Errorf("commandDataBlock.Marshal(): got %d, want %d", l, cdbSize)
	}
	l = len(&statusBlock{})
	if l != maxStatusBlockLen {
		t.Errorf("sbsize: got %d, want %d", l, maxStatusBlockLen)
	}
}

func TestUnlock(t *testing.T) {
	Debug = t.Logf
	// This command: ./hdparm --security-unlock 12345678901234567890123456789012 /dev/null
	// yields this header and data to ioctl(fd, SECURITY_UNLOCK, ...)
	// The 'want' data is derived from a modified version of hdparm (github.com/rminnich/hdparmm)
	// which prints the ioctl parameters as initialized go structs.

	want := &packet{
		packetHeader: packetHeader{
			interfaceID:       'S',
			direction:         -2,
			cmdLen:            16,
			maxStatusBlockLen: 32,
			iovCount:          0,
			dataLen:           512,
			data:              0,
			cdb:               0,
			sb:                0,
			timeout:           15000,
			flags:             0,
			packID:            0,
			usrPtr:            0,
			status:            0,
			maskedStatus:      0,
			msgStatus:         0,
			sbLen:             0,
			hostStatus:        0,
			driverStatus:      0,
			resID:             0,
			duration:          0,
			info:              0,
		},
		command: commandDataBlock{0x85, 0xb, 0x6, 0o0, 0o0, 0o0, 0x1, 0o0, 0o0, 0o0, 0o0, 0o0, 0o0, 0x40, 0xf2, 0o0},
		block: dataBlock{
			0x00, 0x01, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34,
			0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30,
			0x31, 0x32,
		},
	}
	// sb statusBlock

	d := &SGDisk{dev: 0x40, Timeout: DefaultTimeout}
	p := d.unlockPacket("12345678901234567890123456789012", true)
	check(t, p, want)
}

func TestIdentify(t *testing.T) {
	Debug = t.Logf
	// The 'want' data is derived from a modified version of hdparm (github.com/rminnich/hdparmm)
	// which prints the ioctl parameters as initialized go structs.

	want := &packet{
		packetHeader: packetHeader{
			interfaceID:       'S',
			direction:         -3,
			cmdLen:            16,
			maxStatusBlockLen: 32,
			iovCount:          0,
			dataLen:           512,
			data:              0,
			cdb:               0,
			sb:                0,
			timeout:           15000,
			flags:             0,
			packID:            0,
			usrPtr:            0,
			status:            0,
			maskedStatus:      0,
			msgStatus:         0,
			sbLen:             0,
			hostStatus:        0,
			driverStatus:      0,
			resID:             0,
			duration:          0,
			info:              0,
		},
		command: commandDataBlock{0x85, 0x08, 0x0e, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0xec, 0x00},
	}
	// TODO: check status block. Requires a qemu device that supports these operations.
	// sb = statusBlock{0x70, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	p := (&SGDisk{dev: 0x40, Timeout: DefaultTimeout}).identifyPacket()
	check(t, p, want)
}
