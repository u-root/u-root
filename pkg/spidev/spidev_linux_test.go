// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spidev

import (
	"encoding/binary"
	"errors"
	"os"
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/u-root/u-root/pkg/flash/chips"
	"github.com/u-root/u-root/pkg/flash/op"
	"golang.org/x/sys/unix"
)

// mockSpidev simulates the ioctls for spidev.
type mockSpidev struct {
	// forceErrno when set will always return the given error from syscall.
	forceErrno unix.Errno

	mode        Mode
	bitsPerWord uint8
	speedHz     uint32
	transfers   []iocTransfer
}

func (s *mockSpidev) syscall(trap, a1, a2 uintptr, a3 unsafe.Pointer) (r1, r2 uintptr, err unix.Errno) {
	if s.forceErrno != 0 {
		return 0, 0, s.forceErrno
	}

	if trap != unix.SYS_IOCTL {
		return 0, 0, unix.EINVAL
	}

	switch a2 {
	case iocRdBitsPerWord:
		*(*uint8)(a3) = uint8(s.bitsPerWord)
	case iocWrBitsPerWord:
		s.bitsPerWord = *(*uint8)(a3)
	case iocRdMaxSpeedHz:
		*(*uint32)(a3) = uint32(s.speedHz)
	case iocWrMaxSpeedHz:
		s.speedHz = *(*uint32)(a3)
	case iocRdMode32:
		*(*uint32)(a3) = uint32(s.mode)
	case iocWrMode32:
		s.mode = Mode(*(*uint32)(a3))
	default:
		if uint32(a2&^sizeMask) != iocMessage(0) {
			return 0, 0, unix.EINVAL
		}

		// Parse multiple transfer structs.
		size := int((a2 & sizeMask) >> sizeShift)
		if size%binary.Size(iocTransfer{}) != 0 {
			return 0, 0, unix.EINVAL
		}

		// Re-create the slice from the pointer.
		s.transfers = unsafe.Slice((*iocTransfer)(a3), size/binary.Size(iocTransfer{}))

		// Make sure the original pointer is not freed up until this point.
		runtime.KeepAlive(a3)

		// Replace all the non-zero address with 0xdeadbeef because the
		// pointer addresses might change during the test.
		for i := range s.transfers {
			t := &s.transfers[i]
			if t.txBuf != 0 {
				t.txBuf = 0xdeadbeef
			}
			if t.rxBuf != 0 {
				t.rxBuf = 0xdeadbeef
			}
		}
	}

	return 0, 0, 0
}

// TestOpenError tests when Open returns an error like file does not exist.
func TestOpenError(t *testing.T) {
	if _, err := Open("/dev/blahblahblahblah"); !os.IsNotExist(err) {
		t.Fatalf(`Open("/dev/blahblahblahblah got %v; want %v`, err, os.ErrNotExist)
	}
}

// TestGetters tests the functions which return values like Mode, SpeedHz, ...
func TestGetters(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Could not create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Could not open spidev: %v", err)
	}
	defer s.Close()

	m := &mockSpidev{
		// You wouldn't use these values in practice, but it is good
		// for a unit test.
		mode:        0x1234,
		bitsPerWord: 10,
		speedHz:     12345,
	}
	s.syscall = m.syscall

	// Test syscall with and without error.
	for _, tt := range []struct {
		name       string
		forceErrno unix.Errno
		wantErr    error
	}{
		{"", 0, nil},
		{"WithErrno", unix.EAGAIN, unix.EAGAIN},
	} {
		m.forceErrno = tt.forceErrno

		t.Run("ID"+tt.name, func(t *testing.T) {
			m, err := s.ID()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Mode() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			want := chips.ID(0)
			if m != want {
				t.Errorf("ID() = %#v; want %#v", m, want)
			}
		})

		t.Run("Mode"+tt.name, func(t *testing.T) {
			m, err := s.Mode()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Mode() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			want := Mode(0x1234)
			if m != want {
				t.Errorf("Mode() = %#v; want %#v", m, want)
			}
		})

		t.Run("BitsPerWord"+tt.name, func(t *testing.T) {
			bpw, err := s.BitsPerWord()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("BitsPerWord() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			want := uint8(10)
			if bpw != want {
				t.Errorf("BitsPerWord() = %d; want %d", bpw, want)
			}
		})

		t.Run("SpeedHz"+tt.name, func(t *testing.T) {
			hz, err := s.SpeedHz()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SpeedHz() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			want := uint32(12345)
			if hz != want {
				t.Errorf("SpeedHz() = %d; want %d", hz, want)
			}
		})
	}
}

// TestSetters tests the functions which set values like SetMode, SetSpeedHz, ...
func TestSetters(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Could not create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	s, err := Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Could not open spidev: %v", err)
	}
	defer s.Close()

	m := &mockSpidev{}
	s.syscall = m.syscall

	// Test syscall with and without error.
	for _, tt := range []struct {
		name       string
		forceErrno unix.Errno
		wantErr    error
	}{
		{"", 0, nil},
		{"WithErrno", unix.EAGAIN, unix.EAGAIN},
	} {
		m.forceErrno = tt.forceErrno

		t.Run("SetMode"+tt.name, func(t *testing.T) {
			if err := s.SetMode(0x12345); !errors.Is(err, tt.wantErr) {
				t.Errorf("SetMode() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			const want = Mode(0x12345)
			if m.mode != want {
				t.Errorf("SetMode() = %#v; want %#v", m.mode, want)
			}
		})

		t.Run("SetBitsPerWord"+tt.name, func(t *testing.T) {
			if err := s.SetBitsPerWord(20); !errors.Is(err, tt.wantErr) {
				t.Errorf("SetBitsPerWord() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			const want = 20
			if m.bitsPerWord != want {
				t.Errorf("SetBitsPerWord() = %d; want %d", m.bitsPerWord, want)
			}
		})

		t.Run("SetSpeedHz"+tt.name, func(t *testing.T) {
			if err := s.SetSpeedHz(12345); !errors.Is(err, tt.wantErr) {
				t.Errorf("SetSpeedHz() got error %q; want error %q", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			const want = 12345
			if m.speedHz != want {
				t.Errorf("SetSpeedHz() = %d; want %d", m.speedHz, want)
			}
		})
	}
}

func TestTransferString(t *testing.T) {
	for _, tt := range []struct {
		n string
		t Transfer
		s string
	}{
		{
			n: "empty",
			t: Transfer{},
			s: "00...[:0](Unknown(00))",
		},
		{
			n: "Read with no data",
			t: Transfer{Tx: op.Read.Bytes()},
			s: "0x03...[:1](Read)",
		},
	} {
		t.Run(tt.n, func(t *testing.T) {
			s := tt.t.String()
			if s != tt.s {
				t.Fatalf("got %q, want %q", s, tt.s)
			}
		})
	}
}

// TestTransfer tests multiple scenarios involving the Transfer method.
func TestTransfer(t *testing.T) {
	// To avoid OOMing the CI, we set the maxTransferSize to a smaller
	// value temporarily for this test.
	defer func(x int) { maxTransferSize = x }(maxTransferSize)
	maxTransferSize = 0x100000

	for _, tt := range []struct {
		name          string
		transfers     []Transfer
		forceErrno    unix.Errno
		wantTransfers []iocTransfer
		wantErr       error
	}{
		{
			name: "ErrTxOverflow",
			transfers: []Transfer{
				{
					Tx: make([]uint8, maxTransferSize+1),
				},
			},
			wantErr: ErrTxOverflow{
				TxLen: maxTransferSize + 1,
				TxMax: maxTransferSize,
			},
		},
		{
			name: "ErrRxOverflow",
			transfers: []Transfer{
				{
					Rx: make([]uint8, maxTransferSize+1),
				},
			},
			wantErr: ErrRxOverflow{
				RxLen: maxTransferSize + 1,
				RxMax: maxTransferSize,
			},
		},
		{
			name: "ErrBufferMismatch",
			transfers: []Transfer{
				{
					Tx: make([]uint8, 10),
					Rx: make([]uint8, 20),
				},
			},
			wantErr: ErrBufferMismatch{
				TxLen: 10,
				RxLen: 20,
			},
		},
		{
			name:       "Errno",
			forceErrno: unix.EAGAIN,
			transfers: []Transfer{
				{
					Tx: make([]uint8, 10),
					Rx: make([]uint8, 10),
				},
			},
			wantErr: unix.EAGAIN,
		},
		{
			name: "TxZero",
			transfers: []Transfer{
				{
					Rx: make([]uint8, 10),
				},
			},
			wantTransfers: []iocTransfer{
				{
					rxBuf:  0xdeadbeef,
					length: 10,
				},
			},
		},
		{
			name: "RxZero",
			transfers: []Transfer{
				{
					Tx: make([]uint8, 10),
				},
			},
			wantTransfers: []iocTransfer{
				{
					txBuf:  0xdeadbeef,
					length: 10,
				},
			},
		},
		{
			name: "OneTransfer",
			transfers: []Transfer{
				{
					Tx:             []uint8{1, 2, 3},
					Rx:             []uint8{0, 0, 0},
					SpeedHz:        0x12345678,
					DelayUSecs:     0x1234,
					BitsPerWord:    0x8,
					CSChange:       true,
					TxNBits:        24,
					RxNBits:        24,
					WordDelayUSecs: 0x10,
				},
			},
			wantTransfers: []iocTransfer{
				{
					txBuf:          0xdeadbeef,
					rxBuf:          0xdeadbeef,
					length:         3,
					speedHz:        0x12345678,
					delayUSecs:     0x1234,
					bitsPerWord:    0x8,
					csChange:       1,
					txNBits:        24,
					rxNBits:        24,
					wordDelayUSecs: 0x10,
				},
			},
		},
		{
			name: "TwoTransfers",
			transfers: []Transfer{
				{
					Tx: []uint8{1, 2, 3},
					Rx: []uint8{0, 0, 0},
				},
				{
					Tx: []uint8{4, 5, 6, 7},
				},
			},
			wantTransfers: []iocTransfer{
				{
					txBuf:  0xdeadbeef,
					rxBuf:  0xdeadbeef,
					length: 3,
				},
				{
					txBuf:  0xdeadbeef,
					length: 4,
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf("Could not create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			s, err := Open(tmpFile.Name())
			if err != nil {
				t.Fatalf("Could not open spidev: %v", err)
			}
			defer s.Close()

			m := &mockSpidev{
				forceErrno: tt.forceErrno,
			}
			s.syscall = m.syscall

			gotErr := s.Transfer(tt.transfers)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Got Transfer err %q; want %q", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(m.transfers, tt.wantTransfers) {
				t.Errorf("Got Transfers %#v; want %#v", m.transfers, tt.wantTransfers)
			}
		})
	}
}
