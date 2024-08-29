// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/u-root/u-root/pkg/dt"
)

func fdtReader(t *testing.T, fdt *dt.FDT) io.ReaderAt {
	t.Helper()
	var b bytes.Buffer
	fdt.Header.Magic = dt.Magic
	fdt.Header.Version = 17
	if _, err := fdt.Write(&b); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b.Bytes())
}

func TestGetFdtInfo(t *testing.T) {
	for _, tt := range []struct {
		// Inputs
		name string
		fdt  io.ReaderAt

		// Results
		fdtLoad *FdtLoad
		err     error
	}{
		// CASE 1: normal case, a FdtLoad object is returned with expected values
		{
			name: "testdata/upl.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("entry-start", 0x00805ac3),
							dt.PropertyU64("load", 0x00800000),
						)),
					)),
				)),
			}),
			fdtLoad: &FdtLoad{
				Load:       uint64(0x0000800000),
				EntryStart: uint64(0x0000805ac3),
			},
			err: nil,
		},
		// CASE 2: dtb file not found
		{
			name:    "testdata/not_exist_file.dtb",
			fdt:     nil,
			fdtLoad: nil,
			err:     ErrFailToReadFdtFile,
		},
		// CASE 3: missing first level node: /images
		{
			name: "testdata/missing_first_node_images.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("description", dt.WithProperty(
						dt.PropertyString("arch", "x86_64"),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeImagesNotFound,
		},
		// CASE 4: missing second level node: /images/tianocore
		{
			name: "testdata/missing_second_node_images.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithProperty(
						dt.PropertyString("arch", "x86_64"),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeTianocoreNotFound,
		},
		// CASE 5: failed to get /images/tianocore/load property
		{
			name: "testdata/missing_property_load.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("entry-start", 0x00805ac3),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeLoadNotFound,
		},
		// CASE 6: failed to convert /images/tianocore/load property (type error)
		{
			name: "testdata/missing_property_load.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyString("load", "0x00800000"),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrFailToConvertLoad,
		},
		// CASE 7: failed to get /images/tianocore/entry-start property
		{
			name: "testdata/missing_property_entry_start.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("load", 0x00800000),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeEntryStartNotFound,
		},
		// CASE 8: failed to convert /images/tianocore/entry-start property (type error)
		{
			name: "testdata/fail_convert_property_entry_start.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("load", 0x00800000),
							dt.PropertyString("entry-start", "0x00800000"),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrFailToConvertEntryStart,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFdtInfo(tt.name, tt.fdt)
			if tt.err != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.err)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("Unxpected error %q, want = %q", err.Error(), tt.err)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			if tt.fdtLoad != nil && got == nil {
				t.Fatalf("getFdtInfo fdtLoad = nil, want = %v", tt.fdtLoad)
			}

			if tt.fdtLoad != nil {
				if tt.fdtLoad.Load != got.Load {
					t.Fatalf("getFdtInfo fdtLoad.Load = %d, want = %v", got.Load, tt.fdtLoad.Load)
				}
				if tt.fdtLoad.EntryStart != got.EntryStart {
					t.Fatalf("getFdtInfo fdtLoad.EntryStart = %v, want = %v", got.EntryStart, tt.fdtLoad.EntryStart)
				}
			}
		})
	}
}

func mockCPUTempInfoFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tempFile, err := os.CreateTemp(tmpDir, "cpuinfo")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	sysfsCPUInfoPath = tempFile.Name()

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tempFile.Close()
	return tempFile.Name()
}

func TestGetPhysicalAddressSizes(t *testing.T) {
	tests := []struct {
		name           string
		cpuInfoContent string
		expectedBits   uint8
		expectedErr    error
	}{
		{
			name: "Valid Address Size",
			cpuInfoContent: `
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 142
model name	: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
stepping	: 10
microcode	: 0xea
cpu MHz		: 1992.000
cache size	: 8192 KB
physical id	: 0
siblings	: 4
core id		: 0
cpu cores	: 2
apicid		: 0
initial apicid	: 0
address sizes	: 39 bits physical, 48 bits virtual
`,
			expectedBits: 39,
			expectedErr:  nil,
		},
		{
			name: "No Address Size Info",
			cpuInfoContent: `
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 142value out of range
model name	: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
`,
			expectedBits: 0,
			expectedErr:  ErrCPUAddressNotFound,
		},
		{
			name: "Invalid Address Size",
			// number value out of range
			cpuInfoContent: `
address sizes	: 1000 bits physical, 48 bits virtual
`,
			expectedBits: 0,
			expectedErr:  strconv.ErrRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := mockCPUTempInfoFile(t, tt.cpuInfoContent)
			defer os.Remove(tempFile)

			physicalBits, err := getPhysicalAddressSizes()

			if tt.expectedErr == nil {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if physicalBits != tt.expectedBits {
					t.Errorf("Unexpected physical address size %d, want = %d", physicalBits, tt.expectedBits)
				}
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Unxpected error %+v, want = %q", err, tt.expectedErr)
				}
			}
		})
	}
}

func TestAlignHOBLength(t *testing.T) {
	tests := []struct {
		name        string
		expectLen   uint64
		bufLen      int
		expectedErr error
	}{
		{
			name:      "Exact Length",
			expectLen: 5,
			bufLen:    5,
		},
		{
			name:      "Padding Required",
			expectLen: 10,
			bufLen:    5,
		},
		{
			name:        "Negative Padding",
			expectLen:   3,
			bufLen:      5,
			expectedErr: ErrAlignPadRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("12345")
			err := alignHOBLength(tt.expectLen, tt.bufLen, buf)

			if tt.expectedErr == nil {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if uint64(buf.Len()) != tt.expectLen {
					t.Fatalf("alignHOBLength() got = %d, want = %d", buf.Len(), tt.expectLen)
				}
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Unxpected error %+v, want = %q", err, tt.expectedErr)
				}
			}
		})
	}
}
