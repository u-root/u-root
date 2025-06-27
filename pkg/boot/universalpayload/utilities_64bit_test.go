// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 || arm64

package universalpayload

import (
	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

// TestSkipReservedRange64Bit tests 64-bit specific cases for skipReservedRange
func TestSkipReservedRange64Bit(t *testing.T) {
	// Create a test memory map with 64-bit reserved regions
	testMemoryMap := kexec.MemoryMap{
		kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeReserved},
		kexec.TypedRange{Range: kexec.Range{Start: 0x800000, Size: 0x50000}, Type: kexec.RangeReserved},
		kexec.TypedRange{Range: kexec.Range{Start: 0x100000000, Size: 0x1000000}, Type: kexec.RangeReserved}, // 64-bit address
	}

	tests := []struct {
		name           string
		memoryMap      kexec.MemoryMap
		base           uintptr
		attr           uint64
		expectedResult bool
		description    string
	}{
		{
			name:           "Base in 64-bit reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x100000000, // Start of 64-bit reserved region
			attr:           0x40200,     // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address at start of 64-bit reserved memory region should be skipped",
		},
		{
			name:           "Base in middle of 64-bit reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x100800000, // Middle of 64-bit reserved region (0x100000000 + 0x800000)
			attr:           0x40200,     // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address in middle of 64-bit reserved memory region should be skipped",
		},
		{
			name:           "Base at end of 64-bit reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x100FFFFFF, // End of 64-bit reserved region (0x100000000 + 0x1000000 - 1)
			attr:           0x40200,     // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address at end of 64-bit reserved memory region should be skipped",
		},
		{
			name:           "Base outside 64-bit reserved memory region should not skip",
			memoryMap:      testMemoryMap,
			base:           0x101000000, // Just outside 64-bit reserved region
			attr:           0x40200,     // PCIMMIO32Attr
			expectedResult: false,
			description:    "Base address outside 64-bit reserved memory region should not be skipped",
		},
		{
			name:           "Base in 64-bit reserved region with 64-bit MMIO attribute should skip",
			memoryMap:      testMemoryMap,
			base:           0x100000000, // Start of 64-bit reserved region
			attr:           0x140204,    // PCIMMIO64Attr
			expectedResult: true,
			description:    "64-bit MMIO in 64-bit reserved memory region should be skipped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipReservedRange(tt.memoryMap, tt.base, tt.attr)

			if result != tt.expectedResult {
				t.Errorf("skipReservedRange:\n%v, (base: 0x%x, attr: 0x%x) = %v, want %v\nDescription: %s",
					tt.memoryMap, tt.base, tt.attr, result, tt.expectedResult, tt.description)
			}
		})
	}
}
