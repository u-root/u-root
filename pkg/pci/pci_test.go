// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPCIReadConfigRegister(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name     string
		pci      PCI
		offset   int64
		size     int64
		valsWant uint64
		errWant  string
	}{
		{
			name: "read byte 1 from config file",
			pci: PCI{
				FullPath: dir,
			},
			offset:   0,
			size:     8,
			valsWant: 0x00,
		},
		{
			name: "read byte 1,2 from config file",
			pci: PCI{
				FullPath: dir,
			},
			offset:   0,
			size:     16,
			valsWant: 0x1100,
		},
		{
			name: "read byte 1,2,3,4 from config file",
			pci: PCI{
				FullPath: dir,
			},
			offset:   0,
			size:     32,
			valsWant: 0x33221100,
		},
		{
			name: "read byte 1,2,3,4,5,6,7,8 from config file",
			pci: PCI{
				FullPath: dir,
			},
			offset:   0,
			size:     64,
			valsWant: 0x7766554433221100,
		},
		{
			name: "read byte 1,2,3,4,5,6,7,8 from config file with error",
			pci: PCI{
				FullPath: dir,
			},
			offset:   2,
			size:     64,
			valsWant: 0x7766554433221100,
			errWant:  "EOF",
		},
		{
			name: "wrong size",
			pci: PCI{
				FullPath: dir,
			},
			offset:  0,
			size:    0,
			errWant: "only options are 8, 16, 32, 64",
		},
		{
			name: "config file does not exist",
			pci: PCI{
				FullPath: "d",
			},
			errWant: "no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if vals, got := tt.pci.ReadConfigRegister(tt.offset, tt.size); got != nil {
				if !strings.Contains(got.Error(), tt.errWant) {
					t.Errorf("ReadConfig() = %q, want to contain: %q", got, tt.errWant)
				}
			} else {
				if vals != tt.valsWant {
					t.Errorf("ReadConfig() = '%#x', want: '%#x'", vals, tt.valsWant)
				}
			}
		})
	}
}

func TestPCIWriteConfigRegister(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name    string
		pci     PCI
		offset  int64
		size    int64
		val     uint64
		want    string
		errWant string
	}{
		{
			name: "Writing 1 byte to config file",
			pci: PCI{
				FullPath: dir,
			},
			offset: 0,
			size:   8,
			val:    0x00,
			want:   "0011223344556677",
		},
		{
			name: "Writing 2 bytes to config file",
			pci: PCI{
				FullPath: dir,
			},
			offset: 0,
			size:   16,
			val:    0x0011,
			want:   "1100223344556677",
		},
		{
			name: "Writing 4 bytes to config file",
			pci: PCI{
				FullPath: dir,
			},
			offset: 0,
			size:   32,
			val:    0x00112233,
			want:   "3322110044556677",
		},
		{
			name: "Writing 8 bytes to config file",
			pci: PCI{
				FullPath: dir,
			},
			offset: 0,
			size:   64,
			val:    0x0011223344556677,
			want:   "7766554433221100",
		},
		{
			name: "Writing 8 bytes to config file with offset of 2 bytes",
			pci: PCI{
				FullPath: dir,
			},
			offset: 2,
			size:   64,
			val:    0x0011223344556677,
			want:   "77667766554433221100",
		},
		{
			name: "wrong size",
			pci: PCI{
				FullPath: dir,
			},
			offset:  0,
			size:    0,
			errWant: "only options are 8, 16, 32, 64",
		},
		{
			name: "config file does not exist",
			pci: PCI{
				FullPath: "d",
			},
			errWant: "no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pci.WriteConfigRegister(tt.offset, tt.size, tt.val); got != nil {
				if !strings.Contains(got.Error(), tt.errWant) {
					t.Errorf("ReadConfig() = %q, want to contain: %q", got, tt.errWant)
				}
			} else {
				got, err := os.ReadFile(filepath.Join(dir, "config"))
				if err != nil {
					t.Errorf("Failed to read file %v", err)
				}
				if hex.EncodeToString(got) != tt.want {
					t.Errorf("Config file contains = %q, want: %q", hex.EncodeToString(got), tt.want)
				}
			}
		})
	}
}
