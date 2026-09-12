// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"encoding/binary"
	"errors"
	"os"
	"testing"
)

// riscvHeader is the first 64 bytes of a real riscv64 kernel Image.
func riscvHeader(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/riscv64_header.bin")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestParseRISCVFromBytes(t *testing.T) {
	got, err := ParseRISCVFromBytes(riscvHeader(t))
	if err != nil {
		t.Fatalf("ParseRISCVFromBytes = %v, want nil", err)
	}

	want := RISCVHeader{
		Code0:      0x1280006f,
		Code1:      0x0,
		TextOffset: 0x200000,
		ImageSize:  0xb429b8,
		Flags:      0x0,
		Version:    0x2,
		Res1:       0x0,
		Res2:       0x0,
		Magic:      0x5643534952,
		Magic2:     RISCVMagic2,
		Res3:       0x0,
	}
	if got.Header != want {
		t.Errorf("header: got %+v, want %+v", got.Header, want)
	}

	if major, minor := got.Header.MajorVersion(), got.Header.MinorVersion(); major != 0 || minor != 2 {
		t.Errorf("version: got v%d.%d, want v0.2", major, minor)
	}
}

// The riscv64 and arm64 headers are laid out differently, so each parser must
// reject the other architecture's Image rather than silently misreading it.
func TestParseRISCVRejectsArm64(t *testing.T) {
	arm64, err := os.ReadFile("testdata/Image")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseRISCVFromBytes(arm64); !errors.Is(err, ErrBadMagic) {
		t.Errorf("ParseRISCVFromBytes(arm64 Image) = %v, want %v", err, ErrBadMagic)
	}
}

func TestParseArm64RejectsRISCV(t *testing.T) {
	if _, err := ParseFromBytes(riscvHeader(t)); !errors.Is(err, ErrBadMagic) {
		t.Errorf("ParseFromBytes(riscv64 Image) = %v, want %v", err, ErrBadMagic)
	}
}

func TestParseRISCVErrors(t *testing.T) {
	for _, tt := range []struct {
		name string
		mod  func([]byte)
		want error
	}{
		{
			name: "bad magic2",
			mod:  func(b []byte) { binary.LittleEndian.PutUint32(b[0x38:], 0xdeadbeef) },
			want: ErrBadMagic,
		},
		{
			// magic is deprecated as of header v0.2, so a wrong one
			// must not by itself make the Image unparseable.
			name: "deprecated magic ignored",
			mod:  func(b []byte) { binary.LittleEndian.PutUint64(b[0x30:], 0) },
			want: nil,
		},
		{
			name: "zero image size",
			mod:  func(b []byte) { binary.LittleEndian.PutUint64(b[0x10:], 0) },
			want: ErrNoImageSize,
		},
		{
			name: "big endian kernel",
			mod:  func(b []byte) { binary.LittleEndian.PutUint64(b[0x18:], 1) },
			want: ErrBadEndianness,
		},
		{
			name: "truncated",
			mod:  func(b []byte) {},
			want: nil, // replaced below
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b := riscvHeader(t)
			if tt.name == "truncated" {
				if _, err := ParseRISCVFromBytes(b[:32]); err == nil {
					t.Error("ParseRISCVFromBytes(short) = nil, want error")
				}
				return
			}
			tt.mod(b)
			_, err := ParseRISCVFromBytes(b)
			if !errors.Is(err, tt.want) {
				t.Errorf("ParseRISCVFromBytes = %v, want %v", err, tt.want)
			}
		})
	}
}
