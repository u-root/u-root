// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"
)

func TestGetFamily(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want uint16
	}{
		{
			name: "Type not 0xfe",
			val: ProcessorInfo{
				Family:  0xff,
				Family2: 0xaa,
			},
			want: 0xff,
		},
		{
			name: "Type 0xfe wrong length",
			val: ProcessorInfo{
				Family:  0xfe,
				Family2: 0xaa,
			},
			want: 0xfe,
		},
		{
			name: "Type 0xfe",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				Family:  0xfe,
				Family2: 0xaa,
			},
			want: 0xaa,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.GetFamily()
			if result != ProcessorFamily(tt.want) {
				t.Errorf("GeFamily(): '%v', want '%v'", result, tt.want)
			}
		})
	}
}

func TestGetVoltage(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want float32
	}{
		{
			name: "5V",
			val: ProcessorInfo{
				Voltage: 0x1,
			},
			want: 5,
		},
		{
			name: "3.3V",
			val: ProcessorInfo{
				Voltage: 0x2,
			},
			want: 3.3,
		},
		{
			name: "2.9V",
			val: ProcessorInfo{
				Voltage: 0x4,
			},
			want: 2.9,
		},
		{
			name: "Unknown",
			val: ProcessorInfo{
				Voltage: 0xFF,
			},
			want: 12.7,
		},
		{
			name: "Zero",
			val: ProcessorInfo{
				Voltage: 0x10,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.GetVoltage()

			if result != tt.want {
				t.Errorf("GetVoltage(): '%v', want '%v'", result, tt.want)
			}
		})
	}
}

func TestGetCoreCount(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want int
	}{
		{
			name: "",
			val: ProcessorInfo{
				CoreCount:  4,
				CoreCount2: 8,
			},
			want: 4,
		},
		{
			name: "",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				CoreCount:  4,
				CoreCount2: 8,
			},
			want: 4,
		},
		{
			name: "",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				CoreCount:  0xff,
				CoreCount2: 8,
			},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.GetCoreCount()

			if result != tt.want {
				t.Errorf("GetCorCount(): '%v', want '%v'", result, tt.want)
			}
		})
	}
}

func TestGetCoreEnabled(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want int
	}{
		{
			name: "",
			val: ProcessorInfo{
				CoreEnabled:  4,
				CoreEnabled2: 8,
			},
			want: 4,
		},
		{
			name: "",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				CoreEnabled:  4,
				CoreEnabled2: 8,
			},
			want: 4,
		},
		{
			name: "",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				CoreEnabled:  0xff,
				CoreEnabled2: 0x8ff,
			},
			want: 0x8ff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.GetCoreEnabled()

			if result != tt.want {
				t.Errorf("GetCoreEnabled() '%v', want '%v'", result, tt.want)
			}
		})
	}
}

func TestGetThreadCount(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want int
	}{
		{
			name: "Count 4",
			val: ProcessorInfo{
				ThreadCount:  4,
				ThreadCount2: 8,
			},
			want: 4,
		},
		{
			name: "Count 4",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				ThreadCount:  4,
				ThreadCount2: 8,
			},
			want: 4,
		},
		{
			name: "Count 8ff",
			val: ProcessorInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				ThreadCount:  0xff,
				ThreadCount2: 0x8ff,
			},
			want: 0x8ff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.GetThreadCount()

			if result != tt.want {
				t.Errorf("GetThreadCount(): '%v', want '%v'", result, tt.want)
			}
		})
	}
}

func TestProcessorTypeString(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorType
		want string
	}{
		{
			name: "Other",
			val:  ProcessorType(1),
			want: "Other",
		},
		{
			name: "Unknown",
			val:  ProcessorType(2),
			want: "Unknown",
		},
		{
			name: "Central Processor",
			val:  ProcessorType(3),
			want: "Central Processor",
		},
		{
			name: "Math Processor",
			val:  ProcessorType(4),
			want: "Math Processor",
		},
		{
			name: "DSP Processor",
			val:  ProcessorType(5),
			want: "DSP Processor",
		},
		{
			name: "Video Processor",
			val:  ProcessorType(6),
			want: "Video Processor",
		},
		{
			name: "Unknown Processor",
			val:  ProcessorType(8),
			want: "0x8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()

			if result != tt.want {
				t.Errorf("ProcessorType().String(): '%s', want '%s'", result, tt.want)
			}
		})
	}
}

func TestProcessorUpgradeString(t *testing.T) {
	resultStrings := []string{
		"Other",
		"Unknown",
		"Daughter Board",
		"ZIF Socket",
		"Replaceable Piggy Back",
		"None",
		"LIF Socket",
		"Slot 1",
		"Slot 2",
		"370-pin Socket",
		"Slot A",
		"Slot M",
		"Socket 423",
		"Socket A (Socket 462)",
		"Socket 478",
		"Socket 754",
		"Socket 940",
		"Socket 939",
		"Socket mPGA604",
		"Socket LGA771",
		"Socket LGA775",
		"Socket S1",
		"Socket AM2",
		"Socket F (1207)",
		"Socket LGA1366",
		"Socket G34",
		"Socket AM3",
		"Socket C32",
		"Socket LGA1156",
		"Socket LGA1567",
		"Socket PGA988A",
		"Socket BGA1288",
		"Socket rPGA988B",
		"Socket BGA1023",
		"Socket BGA1224",
		"Socket BGA1155",
		"Socket LGA1356",
		"Socket LGA2011",
		"Socket FS1",
		"Socket FS2",
		"Socket FM1",
		"Socket FM2",
		"Socket LGA2011-3",
		"Socket LGA1356-3",
		"Socket LGA1150",
		"Socket BGA1168",
		"Socket BGA1234",
		"Socket BGA1364",
		"Socket AM4",
		"Socket LGA1151",
		"Socket BGA1356",
		"Socket BGA1440",
		"Socket BGA1515",
		"Socket LGA3647-1",
		"Socket SP3",
		"Socket SP3r2",
		"Socket LGA2066",
		"Socket BGA1392",
		"Socket BGA1510",
		"Socket BGA1528",
	}

	for id := range resultStrings {
		ProcessorUpgr := ProcessorUpgrade(id + 1)
		if ProcessorUpgr.String() != resultStrings[id] {
			t.Errorf("ProcessorUpgrade().String(): '%s', want '%s'", ProcessorUpgr.String(), resultStrings[id])
		}
	}
}

func TestProcessorCharacteristicsString(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorCharacteristics
		want string
	}{
		{
			name: "Reserved Characteristics",
			val:  ProcessorCharacteristics(0x01),
			want: "		Reserved",
		},
		{
			name: "Unknown Characteristics",
			val:  ProcessorCharacteristics(0x02),
			want: "		Unknown",
		},
		{
			name: "64-bit capable Characteristics",
			val:  ProcessorCharacteristics(0x04),
			want: "		64-bit capable",
		},
		{
			name: "Multi-Core Characteristics",
			val:  ProcessorCharacteristics(0x08),
			want: "		Multi-Core",
		},
		{
			name: "Hardware Thread Characteristics",
			val:  ProcessorCharacteristics(0x10),
			want: "		Hardware Thread",
		},
		{
			name: "Execute Protection Characteristics",
			val:  ProcessorCharacteristics(0x20),
			want: "		Execute Protection",
		},
		{
			name: "Enhanced Virtualization Characteristics",
			val:  ProcessorCharacteristics(0x40),
			want: "		Enhanced Virtualization",
		},
		{
			name: "Power/Performance Control Characteristics",
			val:  ProcessorCharacteristics(0x80),
			want: "		Power/Performance Control",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val.String() != tt.want {
				t.Errorf("ProcessorCharacteristics().String(): '%s', want '%s'", tt.val.String(), tt.want)
			}
		})
	}
}

func TestParseProcessorInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   *ProcessorInfo
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  &ProcessorInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeBIOSInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: fmt.Errorf("invalid table type 0"),
		},
		{
			name: "Required fields are missing",
			val:  &ProcessorInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeProcessorInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  &ProcessorInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeProcessorInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: fmt.Errorf("error parsing structure"),
		},
		{
			name: "Parse valid ProcessorInfo",
			val: &ProcessorInfo{
				Table: Table{
					Header: Header{
						Type: TableTypeProcessorInfo,
					},
					data: []byte{
						0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a, 0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a,
					},
				},
			},
			table: Table{
				Header: Header{
					Type: TableTypeProcessorInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return len(tt.val.data), tt.want
			}
			_, err := parseProcessorInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseProcessorInfo(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func TestProcessorInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  ProcessorInfo
		want string
	}{
		{
			name: "Empty Struct",
			val:  ProcessorInfo{},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Type: 0x0
	Family: 0x0
	Manufacturer: 
	ID: 00 00 00 00 00 00 00 00
	Version: 
	Voltage: 0.0 V
	External Clock: Unknown
	Max Speed: Unknown
	Current Speed: Unknown
	Status: Unpopulated
	Upgrade: 0x0`,
		},
		{
			name: "Full Struct 0x10 Fam",
			val: ProcessorInfo{
				Type:          ProcessorTypeCentralProcessor,
				Family:        0x10,
				Manufacturer:  "Google",
				ID:            0x1234567812345678,
				L1CacheHandle: 0x1337,
				L2CacheHandle: 0xDEAD,
				L3CacheHandle: 0xBEEF,
				Table: Table{
					data: []byte{
						0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a, 0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a,
					},
				},
				ThreadCount: 8,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Type: Central Processor
	Family: Pentium II Xeon
	Manufacturer: Google
	ID: 78 56 34 12 78 56 34 12
	Signature: Type 1, Family 41, Model 71, Stepping 8
	Flags:
		PSE (Page size extension)
		TSC (Time stamp counter)
		MSR (Model specific registers)
		PAE (Physical address extension)
		APIC (On-chip APIC hardware supported)
		MTRR (Memory type range registers)
		MCA (Machine check architecture)
		PSN (Processor serial number present and enabled)
		DS (Debug store)
		SSE (Streaming SIMD extensions)
		HTT (Multi-threading)
	Version: 
	Voltage: 0.0 V
	External Clock: Unknown
	Max Speed: Unknown
	Current Speed: Unknown
	Status: Unpopulated
	Upgrade: 0x0
	L1 Cache Handle: 0x1337
	L2 Cache Handle: 0xDEAD
	L3 Cache Handle: 0xBEEF
	Serial Number: 
	Asset Tag: 
	Part Number: 
	Core Count: 0
	Core Enabled: 0
	Thread Count: 8
	Characteristics:
		`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val.String() != tt.want {
				t.Errorf("ProcessorInfo().String(): '%s', want '%s'", tt.val.String(), tt.want)
			}
		})
	}
}
