// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import "testing"

func TestGetSizeBytes(t *testing.T) {
	tests := []struct {
		name   string
		MemDev MemoryDevice
		want   uint64
	}{
		{
			name: "Size Zero",
			MemDev: MemoryDevice{
				Size:         0,
				ExtendedSize: 0x1,
			},
			want: 0,
		},
		{
			name: "ExtendedSize",
			MemDev: MemoryDevice{
				Size:         0x7fff,
				ExtendedSize: 0x8000,
			},
			want: 0x800000000,
		},
		{
			name: "Size 0x8001",
			MemDev: MemoryDevice{
				Size:         0x8001,
				ExtendedSize: 0x10,
			},
			want: 0x400,
		},
		{
			name: "Any other Size",
			MemDev: MemoryDevice{
				Size:         0x1000,
				ExtendedSize: 0x1,
			},
			want: 0x100000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDev.GetSizeBytes()
			if got != tt.want {
				t.Errorf("GetSizeBytes() = '%x', want '%x'", got, tt.want)
			}
		})
	}
}

func TestMemoryDeviceFormFactor(t *testing.T) {
	tests := []struct {
		name             string
		MemDevFormFactor MemoryDeviceFormFactor
		want             string
	}{
		{
			name:             "MemoryDeviceFormFactorDIMM",
			MemDevFormFactor: MemoryDeviceFormFactorDIMM,
			want:             "DIMM",
		},
		{
			name:             "MemoryDeviceFormFactorSODIMM",
			MemDevFormFactor: MemoryDeviceFormFactorSODIMM,
			want:             "SODIMM",
		},
		{
			name:             "Unknown",
			MemDevFormFactor: 0x80,
			want:             "0x80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDevFormFactor.String()
			if got != tt.want {
				t.Errorf("MemDevFormFactor.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMemoryDeviceType(t *testing.T) {
	tests := []struct {
		name       string
		MemDevType MemoryDeviceType
		want       string
	}{
		{
			name:       "MemoryDeviceTypeDRAM",
			MemDevType: MemoryDeviceTypeDRAM,
			want:       "DRAM",
		},
		{
			name:       "MemoryDeviceTypeDDR4",
			MemDevType: MemoryDeviceTypeDDR4,
			want:       "DDR4",
		},
		{
			name:       "Unknown",
			MemDevType: 0x80,
			want:       "0x80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDevType.String()
			if got != tt.want {
				t.Errorf("MemDevType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMemoryDeviceTechnology(t *testing.T) {
	tests := []struct {
		name             string
		MemDevTechnology MemoryDeviceTechnology
		want             string
	}{
		{
			name:             "MemoryDeviceTechnologyDRAM",
			MemDevTechnology: MemoryDeviceTechnologyDRAM,
			want:             "DRAM",
		},
		{
			name:             "MemoryDeviceTechnologyNVDIMMP",
			MemDevTechnology: MemoryDeviceTechnologyNVDIMMP,
			want:             "NVDIMM-P",
		},
		{
			name:             "Unknown",
			MemDevTechnology: 0x80,
			want:             "0x80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDevTechnology.String()
			if got != tt.want {
				t.Errorf("MemDevTechnology.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMemoryDeviceTypeDetail(t *testing.T) {
	tests := []struct {
		name             string
		MemDevTypeDetail MemoryDeviceTypeDetail
		want             string
	}{
		{
			name:             "None",
			MemDevTypeDetail: 0x0001,
			want:             "None",
		},
		{
			name:             "MemoryDeviceTypeDetailOther",
			MemDevTypeDetail: MemoryDeviceTypeDetailOther,
			want:             "Other",
		},
		{
			name:             "MemoryDeviceTypeDetailOther",
			MemDevTypeDetail: 0xfffe,
			want:             "Other Unknown Fast-paged Static column Pseudo-static RAMBUS Synchronous CMOS EDO Window DRAM Cache DRAM Non-volatile Registered (Buffered) Unbuffered (Unregistered) LRDIMM",
		},
		{
			name:             "MemoryDeviceTypeDetailWindowDRAM",
			MemDevTypeDetail: MemoryDeviceTypeDetailWindowDRAM,
			want:             "Window DRAM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDevTypeDetail.String()
			if got != tt.want {
				t.Errorf("MemDevTypeDetail.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMemoryDeviceOperatingModeCapability(t *testing.T) {
	tests := []struct {
		name          string
		MemDevOpsCaps MemoryDeviceOperatingModeCapability
		want          string
	}{
		{
			name:          "None",
			MemDevOpsCaps: 0x0001,
			want:          "None",
		},
		{
			name:          "MemoryDeviceOperatingModeCapabilityOther",
			MemDevOpsCaps: MemoryDeviceOperatingModeCapabilityOther,
			want:          "Other",
		},
		{
			name:          "MemoryDeviceTypeDetailOther",
			MemDevOpsCaps: 0xfffe,
			want:          "Other Unknown Volatile memory Byte-accessible persistent memory Block-accessible persistent memory",
		},
		{
			name:          "MemoryDeviceOperatingModeCapabilityUnknown",
			MemDevOpsCaps: MemoryDeviceOperatingModeCapabilityUnknown,
			want:          "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDevOpsCaps.String()
			if got != tt.want {
				t.Errorf("MemoryDeviceOperatingModeCapability.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMemoryDevice(t *testing.T) {
	tests := []struct {
		name   string
		MemDev MemoryDevice
		want   string
	}{
		{
			name:   "Empty Struct",
			MemDev: MemoryDevice{},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Array Handle: 0x0000
	Error Information Handle: 0x0000
	Total Width: Unknown
	Data Width: Unknown
	Size: No Module Installed
	Form Factor: 0x0
	Set: None
	Locator: 
	Bank Locator: 
	Type: 0x0
	Type Detail: None`,
		},
		{
			name: "DeviceSet and Size",
			MemDev: MemoryDevice{
				MemoryErrorInfoHandle: 0xffff,
				DeviceSet:             0xff,
				Size:                  0x10,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Array Handle: 0x0000
	Error Information Handle: No Error
	Total Width: Unknown
	Data Width: Unknown
	Size: 16 MB
	Form Factor: 0x0
	Set: Unknown
	Locator: 
	Bank Locator: 
	Type: 0x0
	Type Detail: None`,
		},
		{
			name: "Metadata",
			MemDev: MemoryDevice{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
						0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15,
					},
				},
				MemoryErrorInfoHandle: 0xfffe,
				DeviceSet:             0xfe,
				Size:                  0x0,
				Speed:                 0x0100,
				Manufacturer:          "Google",
				SerialNumber:          "DEADBEEF",
				AssetTag:              "#354432",
				PartNumber:            "1",
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Array Handle: 0x0000
	Error Information Handle: Not Provided
	Total Width: Unknown
	Data Width: Unknown
	Size: No Module Installed
	Form Factor: 0x0
	Set: 254
	Locator: 
	Bank Locator: 
	Type: 0x0
	Type Detail: None
	Speed: 256 MT/s
	Manufacturer: Google
	Serial Number: DEADBEEF
	Asset Tag: #354432
	Part Number: 1`,
		},
		{
			name: "Voltage and Speed",
			MemDev: MemoryDevice{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
						0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
						0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
						0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
					},
				},
				MemoryErrorInfoHandle: 0xfffe,
				DeviceSet:             0xfe,
				Size:                  0x0,
				Speed:                 0x0100,
				Manufacturer:          "Google",
				SerialNumber:          "DEADBEEF",
				AssetTag:              "#354432",
				PartNumber:            "1",
				Attributes:            0x0D,
				ConfiguredSpeed:       0x0010,
				MinimumVoltage:        0x0,
				MaximumVoltage:        0x0064,
				ConfiguredVoltage:     0x2000,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Array Handle: 0x0000
	Error Information Handle: Not Provided
	Total Width: Unknown
	Data Width: Unknown
	Size: No Module Installed
	Form Factor: 0x0
	Set: 254
	Locator: 
	Bank Locator: 
	Type: 0x0
	Type Detail: None
	Speed: 256 MT/s
	Manufacturer: Google
	Serial Number: DEADBEEF
	Asset Tag: #354432
	Part Number: 1
	Rank: 13
	Configured Memory Speed: 16 MT/s
	Minimum Voltage: Unknown
	Maximum Voltage: 0.1 V
	Configured Voltage: 8.192 V`,
		},
		{
			name: "Size > 0x28",
			MemDev: MemoryDevice{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
						0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13,
						0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
						0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
						0x28, 0x29,
					},
				},
				MemoryErrorInfoHandle:             0xfffe,
				DeviceSet:                         0xfe,
				Size:                              0x0,
				Speed:                             0x0100,
				Manufacturer:                      "Google",
				SerialNumber:                      "DEADBEEF",
				AssetTag:                          "#354432",
				PartNumber:                        "1",
				Attributes:                        0x0D,
				ConfiguredSpeed:                   0x0010,
				MinimumVoltage:                    0x0,
				MaximumVoltage:                    0x0064,
				ConfiguredVoltage:                 0x2000,
				Technology:                        0x1,
				OperatingModeCapability:           0x1,
				FirmwareVersion:                   "2.3",
				ModuleManufacturerID:              0x1234,
				ModuleProductID:                   0x1234,
				SubsystemControllerManufacturerID: 0x11,
				SubsystemControllerProductID:      0x22,
				NonvolatileSize:                   0x00100000,
				CacheSize:                         0x00001000,
				LogicalSize:                       0x00001000,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Array Handle: 0x0000
	Error Information Handle: Not Provided
	Total Width: Unknown
	Data Width: Unknown
	Size: No Module Installed
	Form Factor: 0x0
	Set: 254
	Locator: 
	Bank Locator: 
	Type: 0x0
	Type Detail: None
	Speed: 256 MT/s
	Manufacturer: Google
	Serial Number: DEADBEEF
	Asset Tag: #354432
	Part Number: 1
	Rank: 13
	Configured Memory Speed: 16 MT/s
	Minimum Voltage: Unknown
	Maximum Voltage: 0.1 V
	Configured Voltage: 8.192 V
	Memory Technology: Other
	Memory Operating Mode Capability: None
	Firmware Version: 2.3
	Module Manufacturer ID: Bank 53, Hex 0x12
	Module Product ID: 0x1234
	Memory Subsystem Controller Manufacturer ID: Bank 18, Hex 0x00
	Memory Subsystem Controller Product ID: 0x0022
	Non-Volatile Size: 1 MB
	Volatile Size: None
	Cache Size: 4 kB
	Logical Size: 4 kB`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.MemDev.String()
			if got != tt.want {
				t.Errorf("MemoryDevice.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
