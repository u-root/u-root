// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/orangecms/go-framebuffer/framebuffer"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

func TestReadRSDPAt(t *testing.T) {
	memory := kexec.MemoryMap{{Range: kexec.Range{Start: 32, Size: 4096}, Type: kexec.RangeReserved}}
	data := bytes.Repeat([]byte{0xa5}, 128)
	got, err := readRSDPAt(bytes.NewReader(data), memory, 32)
	if err != nil || !bytes.Equal(got, data[32:68]) {
		t.Fatalf("readRSDPAt() = %x, %v; want %x, nil", got, err, data[32:68])
	}
	if _, err := readRSDPAt(bytes.NewReader(nil), memory, 32); !errors.Is(err, io.EOF) {
		t.Errorf("readRSDPAt(short memory) = %v; want EOF", err)
	}
	for name, address := range map[string]uint64{
		"unsigned overflow": math.MaxUint64,
		"negative offset":   1 << 63,
		"end exceeds int64": math.MaxInt64 - 35,
		"not preserved":     4096,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := readRSDPAt(bytes.NewReader(data), memory, address); err == nil {
				t.Errorf("readRSDPAt(%#x) = nil; want error", address)
			}
		})
	}
	// A malicious high firmware range must not make an unsigned address
	// acceptable as a signed /dev/mem offset (the CodeQL regression).
	memory[0].Start = 1 << 63
	if _, err := readRSDPAt(bytes.NewReader(data), memory, 1<<63); err == nil {
		t.Error("readRSDPAt(high preserved address) = nil; want error")
	}
}

func TestPreserved(t *testing.T) {
	memory := testFirmware(t, 3).memory
	for _, tt := range []struct {
		address, size uint64
		want          bool
	}{
		{0xf0000, 4096, true},
		{0xf0000, 0, false},
		{0xeffff, 36, false},
		{0xf0000, 4097, false},
		{0x100000, 36, false},
		{math.MaxUint64 - 1, 36, false},
	} {
		if got := preserved(memory, tt.address, tt.size); got != tt.want {
			t.Errorf("preserved(%#x, %d) = %v; want %v", tt.address, tt.size, got, tt.want)
		}
	}
}

func TestRSDPRevisions(t *testing.T) {
	legacy := make([]byte, 20)
	copy(legacy, "RSD PTR ")
	legacy[8] = -byteSum(legacy)
	if err := validateRSDP(legacy); err != nil {
		t.Fatalf("validateRSDP(ACPI 1.0) = %v; want nil", err)
	}
	extended := append(append([]byte(nil), legacy...), make([]byte, 16)...)
	extended[15] = 2
	binary.LittleEndian.PutUint32(extended[20:], uint32(len(extended)))
	extended[8] = 0
	extended[8] = -byteSum(extended[:20])
	extended[32] = -byteSum(extended)
	for name, change := range map[string]func([]byte) []byte{
		"short legacy":      func(b []byte) []byte { return b[:19] },
		"short extended":    func(b []byte) []byte { return b[:35] },
		"signature":         func(b []byte) []byte { b[0] ^= 1; return b },
		"legacy checksum":   func(b []byte) []byte { b[8] ^= 1; return b },
		"length":            func(b []byte) []byte { b[20]++; return b },
		"extended checksum": func(b []byte) []byte { b[32] ^= 1; return b },
	} {
		t.Run(name, func(t *testing.T) {
			if err := validateRSDP(change(append([]byte(nil), extended...))); err == nil {
				t.Error("validateRSDP(invalid table) = nil; want error")
			}
		})
	}
}

func TestFramebufferRejectedFormats(t *testing.T) {
	var valid framebuffer.Framebuffer
	valid.Finfo.Smem_start = 0xe0000000
	valid.Finfo.Smem_len = 4096 * 768
	valid.Finfo.Visual = 2
	valid.Finfo.Line_length = 4096
	valid.Vinfo.Xres, valid.Vinfo.Yres, valid.Vinfo.Bits_per_pixel = 1024, 768, 32
	valid.Vinfo.Red.Offset, valid.Vinfo.Red.Length = 16, 8
	valid.Vinfo.Green.Offset, valid.Vinfo.Green.Length = 8, 8
	valid.Vinfo.Blue.Length = 8
	// XRGB is valid too: the reserved channel need not have any bits.
	if _, err := framebufferRecord(&valid); err != nil {
		t.Fatalf("framebufferRecord(XRGB) = %v; want nil", err)
	}
	for name, change := range map[string]func(*framebuffer.Framebuffer){
		"nonlinear":        func(f *framebuffer.Framebuffer) { f.Finfo.Type = 1 },
		"pseudocolor":      func(f *framebuffer.Framebuffer) { f.Finfo.Visual = 3 },
		"24bpp":            func(f *framebuffer.Framebuffer) { f.Vinfo.Bits_per_pixel = 24 },
		"panned":           func(f *framebuffer.Framebuffer) { f.Vinfo.Xoffset = 1 },
		"empty":            func(f *framebuffer.Framebuffer) { f.Vinfo.Xres = 0 },
		"unaligned stride": func(f *framebuffer.Framebuffer) { f.Finfo.Line_length++ },
		"short stride":     func(f *framebuffer.Framebuffer) { f.Finfo.Line_length = 4 },
		"short buffer":     func(f *framebuffer.Framebuffer) { f.Finfo.Smem_len-- },
		"unmapped address": func(f *framebuffer.Framebuffer) { f.Finfo.Smem_start = 64 << 30 },
		"unmapped end":     func(f *framebuffer.Framebuffer) { f.Finfo.Smem_start = (64 << 30) - 1 },
		"reversed channel": func(f *framebuffer.Framebuffer) { f.Vinfo.Red.Right = 1 },
		"wide channel":     func(f *framebuffer.Framebuffer) { f.Vinfo.Red.Length = 9 },
		"channel overflow": func(f *framebuffer.Framebuffer) { f.Vinfo.Red.Offset = 32 },
		"missing blue":     func(f *framebuffer.Framebuffer) { f.Vinfo.Blue.Length = 0 },
	} {
		t.Run(name, func(t *testing.T) {
			fb := valid
			change(&fb)
			if _, err := framebufferRecord(&fb); err == nil {
				t.Error("framebufferRecord(invalid format) = nil; want error")
			}
		})
	}
}
