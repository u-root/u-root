// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// SGDisk is the Linux SCSI Generic interface to SCSI/SATA devices.
// Control is achieved by ioctls on an fd.
// SG is extremely low level, requiring the assembly of Command and Data Blocks,
// and occasionaly the disassembly of Status Blocks.
//
// SG can operate with any version of SCSI or ATA, starting from ATA1 to the present.
// ATA packets became "16-bits wide and 64-bit aware in ATA6 standard in 2003.
// Block addresses in this standard are 48 bits.
// In our usage of SG on Linux, we only use ATA16 and LBA48.
//
// We figure that a standard defined in 2003 is
// fairly universal now, and we don't care about
// hardware older than that.
//
// In Linux, requests to SG are defined by a packet header, used by the kernel;
// a Command and Data Block (cdb), a serialized opcode header for the disk;
// and a block, always 512 bytes, containing data.
//
// We show the serialized format of an ATA16 packet below.
// In this layout, following common nomenclature,
// lob is low order byte, hob is high order byte.
// Why is it done this way? In earlier ATA packets,
// serialized over an 8 bit bus, the LBA was 3 bytes.
// It seems when they doubled the bus, and doubled other
// fields, they put the extra bytes "beside" the existing bytes,
// with the result shown below.
//
// The first 3 bytes of the CDB are information about the request,
// and the last 13 bytes are generic information.
// Command and Data Block layout:
// cdb[ 3] = hob_feat
// cdb[ 4] = feat
// cdb[ 5] = hob_nsect
// cdb[ 6] = nsect
// cdb[ 7] = hob_lbal
// cdb[ 8] = lbal
// cdb[ 9] = hob_lbam
// cdb[10] = lbam
// cdb[11] = hob_lbah
// cdb[12] = lbah
// cdb[13] = device
// cdb[14] = command
// Further, there is a direction which can be to, from, or none.

// packetHeader is the Linux SCSI Generic driver version 3 header structure, or
// sg_io_hdr_t in some code bases.
//
// A pointer to this struct must be passed to the SG_IO ioctl.
//
// Note that some pointers are not word-aligned, i.e. the
// compiler will insert padding; this struct is larger than
// the sum of its parts. This struct has some information
// also contained in the Command and Data Block.
type packetHeader struct {
	interfaceID       int32
	direction         direction
	cmdLen            uint8
	maxStatusBlockLen uint8
	iovCount          uint16
	dataLen           uint32
	data              uintptr
	cdb               uintptr
	sb                uintptr
	timeout           uint32
	flags             uint32
	packID            uint32
	usrPtr            uintptr
	status            uint8
	maskedStatus      uint8
	msgStatus         uint8
	sbLen             uint8
	hostStatus        uint16
	driverStatus      uint16
	resID             int32
	duration          uint32
	info              uint32
}

// packet contains the packetHeader and other information.
type packet struct {
	packetHeader

	// This is additional, per-request-type information
	// needed to create a command and data block.
	// It is assembled from both the Disk and the request type.
	ataType  uint8 // almost always lba48
	transfer uint8
	category uint8
	protocol uint8
	features uint16
	cmd      Cmd
	dev      uint8
	lba      uint64
	nsect    uint16
	dma      bool

	// There are pointers in the packetHeader to this data.
	//
	// We maintain them here to ensure they don't
	// get garbage collected, as the packetHeader only
	// contains uintptrs to refer to them.
	command commandDataBlock
	status  statusBlock
	block   dataBlock
	word    wordBlock
}

type diskFile interface {
	Close() error
	Name() string
	Fd() uintptr
}

// SGDisk implements a Disk using the Linux SG device
type SGDisk struct {
	f      diskFile
	dev    uint8
	packID uint32

	// Timeuut is the timeout on a disk operation.
	Timeout time.Duration
}

// NewSGDisk returns a Disk that uses the Linux SCSI Generic Device.
// It also does an Identify to verify that the target name is a true
// lba48 device.
func NewSGDisk(n string, opt ...SGDiskOpt) (*SGDisk, error) {
	f, err := os.OpenFile(n, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return NewSGDiskFromFile(f, opt...)
}

func NewSGDiskFromFile(f diskFile, opt ...SGDiskOpt) (*SGDisk, error) {
	s := &SGDisk{f: f, Timeout: DefaultTimeout}
	if _, err := s.Identify(); err != nil {
		return nil, err
	}
	for _, o := range opt {
		o(s)
	}
	return s, nil
}

// Close closes any open FDs.
func (s *SGDisk) Close() error {
	return s.f.Close()
}

// genCommandDataBlock creates a Command and Data Block used by
// Linux SGDisk.
func (p *packet) genCommandDataBlock() {
	p.command[0] = ata16
	p.command[1] = p.ataType | p.transfer | p.category | p.protocol
	switch {
	case p.dma && p.dataLen != 0:
		p.command[1] |= protoDMA
	case p.dma && p.dataLen == 0:
		p.command[1] |= nonData
	case !p.dma && p.dataLen != 0:
		if p.direction == _SG_DXFER_TO_DEV {
			p.command[1] |= pioOut
		} else {
			p.command[1] |= pioIn
		}
	case !p.dma && p.dataLen == 0:
		p.command[1] |= nonData
	}
	// libata/AHCI workaround: don't demand sense data for IDENTIFY commands
	// We learned this from hdparm.
	if p.dataLen != 0 {
		p.command[2] |= tlenNsect | tlenSectors
		if p.direction == _SG_DXFER_TO_DEV {
			p.command[2] |= tdirTo
		} else {
			p.command[2] |= tdirFrom
		}
	} else {
		p.command[2] = checkCond
	}
	p.command[3] = uint8(p.features >> 8)
	p.command[4] = uint8(p.features)
	p.command[5] = uint8(p.nsect >> 8)
	p.command[6] = uint8(p.nsect)
	p.command[7] = uint8(p.lba >> 8)
	p.command[8] = uint8(p.lba)
	p.command[9] = uint8(p.lba >> 24)
	p.command[10] = uint8(p.lba >> 16)
	p.command[11] = uint8(p.lba >> 40)
	p.command[12] = uint8(p.lba >> 32)
	p.command[13] = p.dev
	p.command[14] = uint8(p.cmd)
}

func (s *SGDisk) newPacket(cmd Cmd, direction direction, ataType uint8) *packet {
	p := &packet{}
	// These are invariant across all uses of SGDisk.
	p.interfaceID = 'S'
	p.cmdLen = uint8(len(p.command))
	p.data = uintptr(unsafe.Pointer(&p.block[0]))
	p.sb = uintptr(unsafe.Pointer(&p.status[0]))
	p.cdb = uintptr(unsafe.Pointer(&p.command[0]))

	// These are determined by the request.
	p.cmd = cmd
	p.dev = s.dev
	p.packID = uint32(s.packID)
	p.direction = direction
	// Go 1.12 appears not to have Milliseconds.
	p.timeout = uint32(s.Timeout.Seconds() * 1000)

	// These are reasonable defaults which the caller
	// can override.
	p.maxStatusBlockLen = maxStatusBlockLen
	p.iovCount = 0 // is this ever non-zero?
	p.dataLen = uint32(oldSchoolBlockLen)
	p.nsect = 1
	p.ataType = ataType

	return p
}

func (s *SGDisk) unlockPacket(password string, admin bool) *packet {
	p := s.newPacket(unix.WIN_SECURITY_UNLOCK, _SG_DXFER_TO_DEV, lba48)
	p.genCommandDataBlock()
	if admin {
		p.block[1] = 1
	}
	copy(p.block[2:], []byte(password))
	return p
}

// Unlock performs unlock requests for Linux SCSI Generic Disks
func (s *SGDisk) Unlock(password string, admin bool) error {
	p := s.unlockPacket(password, admin)
	if err := s.operate(p); err != nil {
		return err
	}
	return nil
}

func (s *SGDisk) identifyPacket() *packet {
	p := s.newPacket(unix.WIN_IDENTIFY, _SG_DXFER_FROM_DEV, 0)
	p.genCommandDataBlock()
	return p
}

// Identify returns identifying information for Linux SCSI Generic Disks.
func (s *SGDisk) Identify() (*Info, error) {
	p := s.identifyPacket()
	if err := s.operate(p); err != nil {
		return nil, err
	}
	return unpackIdentify(p.status, p.block, p.word), nil
}

// _SG_IO is the ioctl request number for SCSI operations.
const _SG_IO = 0x2285

func (s *SGDisk) operate(p *packet) error {
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(s.f.Fd()), _SG_IO, uintptr(unsafe.Pointer(&p.packetHeader)))
	sb := p.status[0]
	if errno != 0 || sb != 0 {
		return &os.PathError{
			Op:   "ioctl SG_IO",
			Path: s.f.Name(),
			Err:  fmt.Errorf("SCSI generic error %w and drive error status %#02x", errno, sb),
		}
	}
	w, err := p.block.toWordBlock()
	if err != nil {
		return &os.PathError{
			Op:   "converting SG_IO block output",
			Path: s.f.Name(),
			Err:  err,
		}
	}
	// the drive must be ata48. The status should show that even if we did not issue an ata48 command.
	/*if err := w.mustLBA(); err != nil {
		return err
	}*/
	p.word = w
	return nil
}

// SGDiskOpt allows callers of NewSGDisk to set values
type SGDiskOpt func(*SGDisk)

// WithTimeout returns an SGDiskOpt that allows setting a non-default TimeOut
func WithTimeout(timeout time.Duration) SGDiskOpt {
	return func(s *SGDisk) {
		s.Timeout = timeout
	}
}
