// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package spidev wraps the Linux spidev driver and performs low-level SPI
// operations.
package spidev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"unsafe"

	"github.com/u-root/u-root/pkg/flash/chips"
	"github.com/u-root/u-root/pkg/flash/op"
	"golang.org/x/sys/unix"
)

// See Linux "include/uapi/linux/spi/spidev.h" and
// "Documentation/spi/spidev.rst"

// Various ioctl numbers.
const (
	iocRdMode        = 0x80016b01
	iocWrMode        = 0x40016b01
	iocRdLSBFirst    = 0x80016b02
	iocWrLSBFirst    = 0x40016b02
	iocRdBitsPerWord = 0x80016b03
	iocWrBitsPerWord = 0x40016b03
	iocRdMaxSpeedHz  = 0x80046b04
	iocWrMaxSpeedHz  = 0x40046b04
	iocRdMode32      = 0x80046b05
	iocWrMode32      = 0x40046b05
)

// Constants used by iocMessage function.
const (
	// iocMessage0 is the length of a message of 0 length. Use the
	// iocMessage(n) function for an iocMessage of length n.
	iocMessage0 = 0x40006b00
	sizeBits    = 14
	sizeShift   = 16
	sizeMask    = ((1 << sizeBits) - 1) << sizeShift
)

// maxTransferSize is the maximum size of a transfer. This is limited by the
// `length uint32` in the iocTransfer struct.
var maxTransferSize = math.MaxInt32

// iocMessage is an ioctl number for n Transfers. Since the ioctl number
// contains the size of the message, it is not a constant.
func iocMessage(n int) uint32 {
	size := uint32(n * binary.Size(iocTransfer{}))
	if n < 0 || size > (1<<sizeBits) {
		return iocMessage(0)
	}
	return iocMessage0 | (size << sizeShift)
}

// Mode is a bitset of the current SPI mode bits.
type Mode uint32

const (
	// CPHA determines clock phase sampling (1=trailing edge).
	CPHA Mode = 1 << iota
	// CPOL determines clock polarity (1=idle high).
	CPOL
	// CS_HIGH determines chip select polarity (1=idle high).
	CS_HIGH
	// LSB_FIRST determines whether least significant bit comes first in a
	// word (1=LSB first).
	LSB_FIRST
	// THREE_WIRE determines whether the bus operates in three wire mode
	// (1=three wire).
	THREE_WIRE
	// LOOP determines whether the device should loop (1=loop enabled).
	LOOP
	// NO_CS determines whether to disable chip select (1=no chip select).
	NO_CS
	// READY determins ready mode bit.
	READY
	// TX_DUAL determines whether to transmit in dual mode.
	TX_DUAL
	// TX_QUAD determines whether to transmit in quad mode.
	TX_QUAD
	// RX_DUAL determines whether to receive in dual mode.
	RX_DUAL
	// RX_QUAD determines whether to receive in quad mode.
	RX_QUAD
)

// iocTransfer is the data type used by the iocMessage ioctl. Multiple such
// transfers may be chained together in a single ioctl call.
type iocTransfer struct {
	// txBuf contains the userspace address of data to send. If this is 0,
	// then zeros are shifted onto the SPI bus.
	txBuf uint64
	// rxBuf contains the userspace address of data to receive. If this is
	// 0, data received is ignored.
	rxBuf uint64
	// length is the length of max{transfer, txBuf, rxBuf} in bytes.
	length uint32

	speedHz        uint32
	delayUSecs     uint16
	bitsPerWord    uint8
	csChange       uint8
	txNBits        uint8
	rxNBits        uint8
	wordDelayUSecs uint8
	pad            uint8
}

// Transfer is the data and options for a single SPI transfer. Note that a SPI
// transfer is full-duplex meaning data is shifted out of Tx and shifted into
// Rx on the same clock cycle.
type Transfer struct {
	// Tx contains a slice sent over the SPI bus.
	Tx []byte
	// Rx contains a slice received from the SPI bus.
	Rx []byte

	// The following temporarily override the SPI mode. They only apply to
	// the current transfer.

	// SpeedHz sets speed in Hz (optional).
	SpeedHz uint32
	// DelayUSecs is how long to delay before the next transfer (optional).
	DelayUSecs uint16
	// BitsPerWord is the device's wordsize (optional).
	BitsPerWord uint8
	// CSChange controls whether the CS (Chip Select) signal will be
	// de-asserted at the end of the transfer.
	CSChange bool
	// TxNbits controls single, dual or quad mode (optional).
	TxNBits uint8
	// RxNbits controls single, dual or quad mode (optional).
	RxNBits uint8
	// WordDelayUSecs is the delay between words (optional).
	WordDelayUSecs uint8
}

func (t *Transfer) String() string {
	var x [8]byte
	n, _ := io.ReadAtLeast(bytes.NewBuffer(t.Tx), x[:], len(x))
	return fmt.Sprintf("%#02x...[:%d](%s)", x[:n], len(t.Tx), op.OpCode(x[0]).String())
}

// ErrTxOverflow is returned if the Transfer buffer is too large.
type ErrTxOverflow struct {
	TxLen, TxMax int
}

// Error implements the error interface.
func (e ErrTxOverflow) Error() string {
	return fmt.Sprintf("SPI tx buffer overflow (%d > %d)", e.TxLen, e.TxMax)
}

// ErrRxOverflow is returned if the Transfer buffer is too large.
type ErrRxOverflow struct {
	RxLen, RxMax int
}

// Error implements the error interface.
func (e ErrRxOverflow) Error() string {
	return fmt.Sprintf("SPI rx buffer overflow (%d > %d)", e.RxLen, e.RxMax)
}

// ErrBufferMismatch is returned if the rx and tx buffers do not have equal
// length.
type ErrBufferMismatch struct {
	TxLen, RxLen int
}

// Error implements the error interface.
func (e ErrBufferMismatch) Error() string {
	return fmt.Sprintf("SPI tx and rx buffers of unequal length (%d != %d)", e.TxLen, e.RxLen)
}

// SPI wraps the Linux device driver and provides low-level SPI operations.
type SPI struct {
	f *os.File
	// Used for mocking.
	syscall func(trap, a1, a2 uintptr, a3 unsafe.Pointer) (r1, r2 uintptr, err unix.Errno)
	// logger allows logging
	logger func(string, ...any)
}

type opt func(s *SPI) error

// WithLogger returns an opt which can be used in Open to add
// a logger. A common usage would be:
// spidev.Open("/dev/spidev0.0", WithLogger(log.Printf))
func WithLogger(l func(string, ...any)) opt {
	return func(s *SPI) error {
		s.logger = l
		return nil
	}
}

// safe tries to set "safe" settings for initial SPI operation.
// However, settings may not succeed, for $REASONS$.
// Hardware is highly variable.
// If there is an error, log it, and continue.
func (s *SPI) safe() {
	if err := s.SetSpeedHz(500000); err != nil {
		s.logger("warning only: set speed to %d HZ err %v", 500000, err)
	}
}

// Open opens a new SPI device. dev is a filename such as "/dev/spidev0.0".
// Remember to call Close() once done.
func Open(dev string, opts ...opt) (*SPI, error) {
	f, err := os.OpenFile(dev, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	s := &SPI{
		f:      f,
		logger: func(string, ...any) {}, // log.Printf,
		// logger: log.Printf,
		// a3 must be an unsafe.Pointer instead of a uintptr, otherwise
		// we cannot mock out in the test without creating a race
		// condition. See `go doc unsafe.Pointer`.
		syscall: func(trap, a1, a2 uintptr, a3 unsafe.Pointer) (r1, r2 uintptr, err unix.Errno) {
			return unix.Syscall(trap, a1, a2, uintptr(a3))
		},
	}

	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	s.safe()

	return s, nil
}

// Close closes the SPI device.
func (s *SPI) Close() error {
	return s.f.Close()
}

// Transfer performs multiple SPI reads and/or writes in a single function.
// See the Transfer struct for details.
func (s *SPI) Transfer(transfers []Transfer) error {
	// Convert []Transfer to []iocTransfer.
	it := make([]iocTransfer, len(transfers))
	for i, t := range transfers {
		s.logger("%d:%s", i, t.String())
		it[i] = iocTransfer{
			speedHz:        t.SpeedHz,
			delayUSecs:     t.DelayUSecs,
			bitsPerWord:    t.BitsPerWord,
			txNBits:        t.TxNBits,
			rxNBits:        t.RxNBits,
			wordDelayUSecs: t.WordDelayUSecs,
		}
		if len(t.Tx) > maxTransferSize {
			return ErrTxOverflow{len(t.Tx), maxTransferSize}
		}
		if len(t.Rx) > maxTransferSize {
			return ErrRxOverflow{len(t.Rx), maxTransferSize}
		}
		if len(t.Tx) != 0 && len(t.Rx) != 0 && len(t.Tx) != len(t.Rx) {
			return ErrBufferMismatch{len(t.Tx), len(t.Rx)}
		}
		if len(t.Tx) != 0 {
			txBuf := &transfers[i].Tx[0]
			it[i].txBuf = uint64(uintptr(unsafe.Pointer(txBuf)))
			it[i].length = uint32(len(t.Tx))
			// Defer garbage collection until after the syscall.
			defer runtime.KeepAlive(txBuf)
		}
		if len(t.Rx) != 0 {
			rxBuf := &transfers[i].Rx[0]
			it[i].rxBuf = uint64(uintptr(unsafe.Pointer(rxBuf)))
			it[i].length = uint32(len(t.Rx))
			// Defer garbage collection until after the syscall.
			defer runtime.KeepAlive(rxBuf)
		}
		if transfers[i].CSChange {
			it[i].csChange = 1
		}
	}

	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), uintptr(iocMessage(len(transfers))), unsafe.Pointer(&it[0])); err != 0 {
		return os.NewSyscallError("ioctl(SPI_IOC_MESSAGE)", err)
	}
	return nil
}

// Mode returns the Mode bitset.
func (s *SPI) Mode() (Mode, error) {
	var m Mode
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocRdMode32, unsafe.Pointer(&m)); err != 0 {
		return 0, os.NewSyscallError("ioctl(SPI_IOC_RD_MODE32)", err)
	}
	return m, nil
}

// SetMode sets the Mode bitset.
func (s *SPI) SetMode(m Mode) error {
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocWrMode32, unsafe.Pointer(&m)); err != 0 {
		return os.NewSyscallError("ioctl(SPI_IOC_WR_MODE32)", err)
	}
	return nil
}

// BitsPerWord sets the number of bits per word. Myy understanding is this is
// only useful if there is a word delay.
func (s *SPI) BitsPerWord() (uint8, error) {
	var bpw uint8
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocRdBitsPerWord, unsafe.Pointer(&bpw)); err != 0 {
		return bpw, os.NewSyscallError("ioctl(SPI_IOC_RD_BITS_PER_WORD)", err)
	}
	return bpw, nil
}

// SetBitsPerWord sets the number of bits per word.
func (s *SPI) SetBitsPerWord(bpw uint8) error {
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocWrBitsPerWord, unsafe.Pointer(&bpw)); err != 0 {
		return os.NewSyscallError("ioctl(SPI_IOC_WR_BITS_PER_WORD)", err)
	}
	return nil
}

// SpeedHz gets the transfer speed.
func (s *SPI) SpeedHz() (uint32, error) {
	var hz uint32
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocRdMaxSpeedHz, unsafe.Pointer(&hz)); err != 0 {
		return hz, os.NewSyscallError("ioctl(SPI_IOC_RD_MAX_SPEED_HZ)", err)
	}
	return hz, nil
}

// SetSpeedHz sets the transfer speed.
func (s *SPI) SetSpeedHz(hz uint32) error {
	if _, _, err := s.syscall(unix.SYS_IOCTL, s.f.Fd(), iocWrMaxSpeedHz, unsafe.Pointer(&hz)); err != 0 {
		return os.NewSyscallError("ioctl(SPI_IOC_WR_MAX_SPEED_HZ)", err)
	}
	return nil
}

// ID gets ID.
func (s *SPI) ID() (chips.ID, error) {
	// Wake it up, then get the id.
	// PRDRES is not universally handled on all devices, but that's ok.
	// but CE MUST drop, so we structure this as two separate
	// transfers to ensure that happens.
	var id [4]byte
	transfers := []Transfer{
		{
			Tx:       []byte{byte(op.PRDRES)},
			Rx:       make([]byte, 1),
			CSChange: true,
		},
		{
			Tx: []byte{byte(op.ReadJEDECID), 0, 0, 0},
			Rx: id[:],
		},
	}

	if err := s.Transfer(transfers); err != nil {
		return -1, err
	}

	id[0] = 0
	return chips.ID(binary.BigEndian.Uint32(id[:])), nil
}

// Status gets the 8 bit status register.
// While this is similar to ID, there is a good chance
// there will be special cases for each opcode type,
// so they should probably remain separate.
// SPI is always full of surprises.
func (s *SPI) Status() (op.Status, error) {
	var status [2]byte
	transfers := []Transfer{
		{
			Tx:       []byte{byte(op.PRDRES)},
			Rx:       make([]byte, 1),
			CSChange: true,
		},
		{
			Tx: []byte{byte(op.ReadStatus), 0},
			Rx: status[:],
		},
	}

	// in the event of an error, return all 1s,
	// making the chip look busy.
	if err := s.Transfer(transfers); err != nil {
		return op.Status(0xff), err
	}

	return op.Status(status[1]), nil
}
