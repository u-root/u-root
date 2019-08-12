// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"errors"
	"fmt"
	"strings"
)

// Much of this is auto-generated. If adding a new type, see README for instructions.

// ProcessorInformation is defined in DSP0134 x.x.
type ProcessorInformation struct {
	Table
	SocketDesignation string                   // 04h
	Type              ProcessorType            // 05h
	Family            uint8                    // 06h
	Manufacturer      string                   // 07h
	ID                uint64                   // 08h
	Version           string                   // 10h
	Voltage           uint8                    // 11h
	ExternalClock     uint16                   // 12h
	MaxSpeed          uint16                   // 14h
	CurrentSpeed      uint16                   // 16h
	Status            ProcessorStatus          // 18h
	Upgrade           ProcessorUpgrade         // 19h
	L1CacheHandle     uint16                   // 1Ah
	L2CacheHandle     uint16                   // 1Ch
	L3CacheHandle     uint16                   // 1Eh
	SerialNumber      string                   // 20h
	AssetTag          string                   // 21h
	PartNumber        string                   // 22h
	CoreCount         uint8                    // 23h
	CoreEnabled       uint8                    // 24h
	ThreadCount       uint8                    // 25h
	Characteristics   ProcessorCharacteristics // 26h
	Family2           ProcessorFamily          // 28h
	CoreCount2        uint16                   // 2Ah
	CoreEnabled2      uint16                   // 2Ch
	ThreadCount2      uint16                   // 2Eh
}

// NewProcessorInformation parses a generic Table into ProcessorInformation.
func NewProcessorInformation(t *Table) (*ProcessorInformation, error) {
	if t.Type != TableTypeProcessorInformation {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x1a {
		return nil, errors.New("required fields missing")
	}
	pi := &ProcessorInformation{Table: *t}
	_, err := parseStruct(t, 0 /* off */, false /* complete */, pi)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

// GetFamily returns the processor family, taken from the appropriate field.
func (pi *ProcessorInformation) GetFamily() ProcessorFamily {
	if pi.Family == 0xfe && pi.Len() >= 0x2a {
		return pi.Family2
	}
	return ProcessorFamily(pi.Family)
}

// GetVoltage returns the processor voltage, in volts.
func (pi *ProcessorInformation) GetVoltage() float32 {
	if pi.Voltage&0x80 == 0 {
		switch {
		case pi.Voltage&1 != 0:
			return 5000
		case pi.Voltage&2 != 0:
			return 3300
		case pi.Voltage&4 != 0:
			return 2900
		}
		return 0
	}
	return float32(pi.Voltage&0x7f) / 10.0
}

// GetCoreCount returns the number of cores detected by the BIOS for this processor socket.
func (pi *ProcessorInformation) GetCoreCount() int {
	if pi.Len() >= 0x2c && pi.CoreCount == 0xff {
		return int(pi.CoreCount2)
	}
	return int(pi.CoreCount)
}

// GetCoreEnabled returns the number of cores that are enabled by the BIOS and available for Operating System use.
func (pi *ProcessorInformation) GetCoreEnabled() int {
	if pi.Len() >= 0x2e && pi.CoreEnabled == 0xff {
		return int(pi.CoreEnabled2)
	}
	return int(pi.CoreEnabled)
}

// GetThreadCount returns the total number of threads detected by the BIOS for this processor socket.
func (pi *ProcessorInformation) GetThreadCount() int {
	if pi.Len() >= 0x30 && pi.ThreadCount == 0xff {
		return int(pi.ThreadCount2)
	}
	return int(pi.ThreadCount)
}

func (pi *ProcessorInformation) String() string {
	freqStr := func(v uint16) string {
		if v == 0 {
			return "Unknown"
		}
		return fmt.Sprintf("%d MHz", v)
	}
	cacheHandleStr := func(h uint16) string {
		if h == 0xffff {
			return "n/a"
		}
		return fmt.Sprintf("0x%04X", h)
	}
	f := pi.GetFamily()
	sig := ""
	haveFlags := false
	switch {
	// Intel
	case (f >= 0x0B && f <= 0x15) || /* Intel, Cyrix */
		(f >= 0x28 && f <= 0x2F) ||
		(f >= 0xA1 && f <= 0xB3) ||
		f == 0xB5 ||
		(f >= 0xB9 && f <= 0xC7) ||
		(f >= 0xCD && f <= 0xCF) ||
		(f >= 0xD2 && f <= 0xDB) || /* VIA, Intel */
		(f >= 0xDD && f <= 0xE0):
		eax := uint32(pi.ID & 0xffffffff)
		sig = fmt.Sprintf("Type %d, Family %d, Model %d, Stepping %d",
			(eax>>12)&0x3, ((eax>>20)&0xff)+((eax>>8)&0xf), ((eax>>12)&0xf0)+((eax>>4)&0xf), eax&0xf)
		haveFlags = true
	// AMD
	case (f >= 0x18 && f <= 0x1D) ||
		f == 0x1F ||
		(f >= 0x38 && f <= 0x3F) ||
		(f >= 0x46 && f <= 0x4F) ||
		(f >= 0x66 && f <= 0x6B) ||
		(f >= 0x83 && f <= 0x8F) ||
		(f >= 0xB6 && f <= 0xB7) ||
		(f >= 0xE4 && f <= 0xEF):
		eax := uint32(pi.ID & 0xffffffff)
		fam := (eax >> 8) & 0xf
		if ((eax >> 8) & 0xf) == 0xf {
			fam += (eax >> 20) & 0xff
		}
		mod := (eax >> 4) & 0xf
		if ((eax >> 8) & 0xf) == 0xf {
			mod += (eax >> 12) & 0xf0
		}
		sig = fmt.Sprintf("Family %d, Model %d, Stepping %d", fam, mod, eax&0xf)
		haveFlags = true
	// ARM
	case (f >= 0x100 && f <= 0x101) ||
		(f >= 0x118 && f <= 0x119):
		if midr := uint32(pi.ID & 0xffffffff); midr != 0 {
			sig = fmt.Sprintf("Implementor 0x%02x, Variant 0x%x, Architecture %d, Part 0x%03x, Revision %d",
				midr>>24, (midr>>20)&0xf, (midr>>16)&0xf, (midr>>4)&0xfff, midr&0xf)
		}
	}
	lines := []string{
		pi.Header.String(),
		fmt.Sprintf("Socket Designation: %s", pi.SocketDesignation),
		fmt.Sprintf("Type: %s", pi.Type),
		fmt.Sprintf("Family: %s", pi.GetFamily()),
		fmt.Sprintf("Manufacturer: %s", pi.Manufacturer),
		fmt.Sprintf("ID: %02X %02X %02X %02X %02X %02X %02X %02X",
			(pi.ID>>0)&0xff, (pi.ID>>8)&0xff, (pi.ID>>16)&0xff, (pi.ID>>24)&0xff,
			(pi.ID>>32)&0xff, (pi.ID>>40)&0xff, (pi.ID>>48)&0xff, (pi.ID>>56)&0xff,
		),
	}
	if sig != "" {
		lines = append(lines, fmt.Sprintf("Signature: %s", sig))
	}
	if haveFlags {
		lines = append(lines, "Flags:")
		edx := uint32(pi.ID >> 32)
		for n, s := range [32]string{
			"FPU (Floating-point unit on-chip)", /* 0 */
			"VME (Virtual mode extension)",
			"DE (Debugging extension)",
			"PSE (Page size extension)",
			"TSC (Time stamp counter)",
			"MSR (Model specific registers)",
			"PAE (Physical address extension)",
			"MCE (Machine check exception)",
			"CX8 (CMPXCHG8 instruction supported)",
			"APIC (On-chip APIC hardware supported)",
			"", /* 10 */
			"SEP (Fast system call)",
			"MTRR (Memory type range registers)",
			"PGE (Page global enable)",
			"MCA (Machine check architecture)",
			"CMOV (Conditional move instruction supported)",
			"PAT (Page attribute table)",
			"PSE-36 (36-bit page size extension)",
			"PSN (Processor serial number present and enabled)",
			"CLFSH (CLFLUSH instruction supported)",
			"", /* 20 */
			"DS (Debug store)",
			"ACPI (ACPI supported)",
			"MMX (MMX technology supported)",
			"FXSR (FXSAVE and FXSTOR instructions supported)",
			"SSE (Streaming SIMD extensions)",
			"SSE2 (Streaming SIMD extensions 2)",
			"SS (Self-snoop)",
			"HTT (Multi-threading)",
			"TM (Thermal monitor supported)",
			"",                            /* 30 */
			"PBE (Pending break enabled)", /* 31 */
		} {
			if edx&(1<<uint(n)) != 0 && s != "" {
				lines = append(lines, "\t"+s)
			}
		}
	}
	lines = append(lines,
		fmt.Sprintf("Version: %s", pi.Version),
		fmt.Sprintf("Voltage: %.1f V", pi.GetVoltage()),
		fmt.Sprintf("External Clock: %s", freqStr(pi.ExternalClock)),
		fmt.Sprintf("Max Speed: %s", freqStr(pi.MaxSpeed)),
		fmt.Sprintf("Current Speed: %s", freqStr(pi.CurrentSpeed)),
		fmt.Sprintf("Status: %s", pi.Status),
		fmt.Sprintf("Upgrade: %s", pi.Upgrade),
	)
	if pi.Len() > 0x1a {
		lines = append(lines,
			fmt.Sprintf("L1 Cache Handle: %s", cacheHandleStr(pi.L1CacheHandle)),
			fmt.Sprintf("L2 Cache Handle: %s", cacheHandleStr(pi.L2CacheHandle)),
			fmt.Sprintf("L3 Cache Handle: %s", cacheHandleStr(pi.L3CacheHandle)),
		)
	}
	if pi.Len() > 0x20 {
		lines = append(lines,
			fmt.Sprintf("Serial Number: %s", pi.SerialNumber),
			fmt.Sprintf("Asset Tag: %s", pi.AssetTag),
			fmt.Sprintf("Part Number: %s", pi.PartNumber),
		)
	}
	if pi.Len() > 0x23 {
		lines = append(lines,
			fmt.Sprintf("Core Count: %d", pi.GetCoreCount()),
			fmt.Sprintf("Core Enabled: %d", pi.GetCoreEnabled()),
			fmt.Sprintf("Thread Count: %d", pi.GetThreadCount()),
			fmt.Sprintf("Characteristics:\n%s", pi.Characteristics),
		)
	}
	return strings.Join(lines, "\n\t")
}

// ProcessorType is defined in DSP0134 7.5.1.
type ProcessorType uint8

// ProcessorType values are defined in DSP0134 7.5.1.
const (
	ProcessorTypeOther            ProcessorType = 0x01 // Other
	ProcessorTypeUnknown                        = 0x02 // Unknown
	ProcessorTypeCentralProcessor               = 0x03 // Central Processor
	ProcessorTypeMathProcessor                  = 0x04 // Math Processor
	ProcessorTypeDSPProcessor                   = 0x05 // DSP Processor
	ProcessorTypeVideoProcessor                 = 0x06 // Video Processor
)

func (v ProcessorType) String() string {
	switch v {
	case ProcessorTypeOther:
		return "Other"
	case ProcessorTypeUnknown:
		return "Unknown"
	case ProcessorTypeCentralProcessor:
		return "Central Processor"
	case ProcessorTypeMathProcessor:
		return "Math Processor"
	case ProcessorTypeDSPProcessor:
		return "DSP Processor"
	case ProcessorTypeVideoProcessor:
		return "Video Processor"
	}
	return fmt.Sprintf("%d", v)
}

// ProcessorFamily is defined in DSP0134 7.5.2.
type ProcessorFamily uint16

// ProcessorFamily values are defined in DSP0134 7.5.2.
const (
	ProcessorFamilyOther                        ProcessorFamily = 0x01  // Other
	ProcessorFamilyUnknown                                      = 0x02  // Unknown
	ProcessorFamily8086                                         = 0x03  // 8086
	ProcessorFamily80286                                        = 0x04  // 80286
	ProcessorFamily80386                                        = 0x05  // 80386
	ProcessorFamily80486                                        = 0x06  // 80486
	ProcessorFamily8087                                         = 0x07  // 8087
	ProcessorFamily80287                                        = 0x08  // 80287
	ProcessorFamily80387                                        = 0x09  // 80387
	ProcessorFamily80487                                        = 0x0a  // 80487
	ProcessorFamilyPentium                                      = 0x0b  // Pentium
	ProcessorFamilyPentiumPro                                   = 0x0c  // Pentium Pro
	ProcessorFamilyPentiumII                                    = 0x0d  // Pentium II
	ProcessorFamilyPentiumMMX                                   = 0x0e  // Pentium MMX
	ProcessorFamilyCeleron                                      = 0x0f  // Celeron
	ProcessorFamilyPentiumIIXeon                                = 0x10  // Pentium II Xeon
	ProcessorFamilyPentiumIII                                   = 0x11  // Pentium III
	ProcessorFamilyM1                                           = 0x12  // M1
	ProcessorFamilyM2                                           = 0x13  // M2
	ProcessorFamilyCeleronM                                     = 0x14  // Celeron M
	ProcessorFamilyPentium4HT                                   = 0x15  // Pentium 4 HT
	ProcessorFamilyDuron                                        = 0x18  // Duron
	ProcessorFamilyK5                                           = 0x19  // K5
	ProcessorFamilyK6                                           = 0x1a  // K6
	ProcessorFamilyK62                                          = 0x1b  // K6-2
	ProcessorFamilyK63                                          = 0x1c  // K6-3
	ProcessorFamilyAthlon                                       = 0x1d  // Athlon
	ProcessorFamilyAMD29000                                     = 0x1e  // AMD29000
	ProcessorFamilyK62Plus                                      = 0x1f  // K6-2+
	ProcessorFamilyPowerPC                                      = 0x20  // Power PC
	ProcessorFamilyPowerPC601                                   = 0x21  // Power PC 601
	ProcessorFamilyPowerPC603                                   = 0x22  // Power PC 603
	ProcessorFamilyPowerPC603Plus                               = 0x23  // Power PC 603+
	ProcessorFamilyPowerPC604                                   = 0x24  // Power PC 604
	ProcessorFamilyPowerPC620                                   = 0x25  // Power PC 620
	ProcessorFamilyPowerPCX704                                  = 0x26  // Power PC x704
	ProcessorFamilyPowerPC750                                   = 0x27  // Power PC 750
	ProcessorFamilyCoreDuo                                      = 0x28  // Core Duo
	ProcessorFamilyCoreDuoMobile                                = 0x29  // Core Duo Mobile
	ProcessorFamilyCoreSoloMobile                               = 0x2a  // Core Solo Mobile
	ProcessorFamilyAtom                                         = 0x2b  // Atom
	ProcessorFamilyCoreM                                        = 0x2c  // Core M
	ProcessorFamilyCoreM3                                       = 0x2d  // Core m3
	ProcessorFamilyCoreM5                                       = 0x2e  // Core m5
	ProcessorFamilyCoreM7                                       = 0x2f  // Core m7
	ProcessorFamilyAlpha                                        = 0x30  // Alpha
	ProcessorFamilyAlpha21064                                   = 0x31  // Alpha 21064
	ProcessorFamilyAlpha21066                                   = 0x32  // Alpha 21066
	ProcessorFamilyAlpha21164                                   = 0x33  // Alpha 21164
	ProcessorFamilyAlpha21164PC                                 = 0x34  // Alpha 21164PC
	ProcessorFamilyAlpha21164a                                  = 0x35  // Alpha 21164a
	ProcessorFamilyAlpha21264                                   = 0x36  // Alpha 21264
	ProcessorFamilyAlpha21364                                   = 0x37  // Alpha 21364
	ProcessorFamilyTurionIIUltraDualCoreMobileM                 = 0x38  // Turion II Ultra Dual-Core Mobile M
	ProcessorFamilyTurionIIDualCoreMobileM                      = 0x39  // Turion II Dual-Core Mobile M
	ProcessorFamilyAthlonIIDualCoreM                            = 0x3a  // Athlon II Dual-Core M
	ProcessorFamilyOpteron6100                                  = 0x3b  // Opteron 6100
	ProcessorFamilyOpteron4100                                  = 0x3c  // Opteron 4100
	ProcessorFamilyOpteron6200                                  = 0x3d  // Opteron 6200
	ProcessorFamilyOpteron4200                                  = 0x3e  // Opteron 4200
	ProcessorFamilyFX                                           = 0x3f  // FX
	ProcessorFamilyMIPS                                         = 0x40  // MIPS
	ProcessorFamilyMIPSR4000                                    = 0x41  // MIPS R4000
	ProcessorFamilyMIPSR4200                                    = 0x42  // MIPS R4200
	ProcessorFamilyMIPSR4400                                    = 0x43  // MIPS R4400
	ProcessorFamilyMIPSR4600                                    = 0x44  // MIPS R4600
	ProcessorFamilyMIPSR10000                                   = 0x45  // MIPS R10000
	ProcessorFamilyCSeries                                      = 0x46  // C-Series
	ProcessorFamilyESeries                                      = 0x47  // E-Series
	ProcessorFamilyASeries                                      = 0x48  // A-Series
	ProcessorFamilyGSeries                                      = 0x49  // G-Series
	ProcessorFamilyZSeries                                      = 0x4a  // Z-Series
	ProcessorFamilyRSeries                                      = 0x4b  // R-Series
	ProcessorFamilyOpteron4300                                  = 0x4c  // Opteron 4300
	ProcessorFamilyOpteron6300                                  = 0x4d  // Opteron 6300
	ProcessorFamilyOpteron3300                                  = 0x4e  // Opteron 3300
	ProcessorFamilyFirePro                                      = 0x4f  // FirePro
	ProcessorFamilySPARC                                        = 0x50  // SPARC
	ProcessorFamilySuperSPARC                                   = 0x51  // SuperSPARC
	ProcessorFamilyMicroSPARCII                                 = 0x52  // MicroSPARC II
	ProcessorFamilyMicroSPARCIIep                               = 0x53  // MicroSPARC IIep
	ProcessorFamilyUltraSPARC                                   = 0x54  // UltraSPARC
	ProcessorFamilyUltraSPARCII                                 = 0x55  // UltraSPARC II
	ProcessorFamilyUltraSPARCIIi                                = 0x56  // UltraSPARC IIi
	ProcessorFamilyUltraSPARCIII                                = 0x57  // UltraSPARC III
	ProcessorFamilyUltraSPARCIIIi                               = 0x58  // UltraSPARC IIIi
	ProcessorFamily68040                                        = 0x60  // 68040
	ProcessorFamily68xxx                                        = 0x61  // 68xxx
	ProcessorFamily68000                                        = 0x62  // 68000
	ProcessorFamily68010                                        = 0x63  // 68010
	ProcessorFamily68020                                        = 0x64  // 68020
	ProcessorFamily68030                                        = 0x65  // 68030
	ProcessorFamilyAthlonX4                                     = 0x66  // Athlon X4
	ProcessorFamilyOpteronX1000                                 = 0x67  // Opteron X1000
	ProcessorFamilyOpteronX2000                                 = 0x68  // Opteron X2000
	ProcessorFamilyOpteronASeries                               = 0x69  // Opteron A-Series
	ProcessorFamilyOpteronX3000                                 = 0x6a  // Opteron X3000
	ProcessorFamilyZen                                          = 0x6b  // Zen
	ProcessorFamilyHobbit                                       = 0x70  // Hobbit
	ProcessorFamilyCrusoeTM5000                                 = 0x78  // Crusoe TM5000
	ProcessorFamilyCrusoeTM3000                                 = 0x79  // Crusoe TM3000
	ProcessorFamilyEfficeonTM8000                               = 0x7a  // Efficeon TM8000
	ProcessorFamilyWeitek                                       = 0x80  // Weitek
	ProcessorFamilyItanium                                      = 0x82  // Itanium
	ProcessorFamilyAthlon64                                     = 0x83  // Athlon 64
	ProcessorFamilyOpteron                                      = 0x84  // Opteron
	ProcessorFamilySempron                                      = 0x85  // Sempron
	ProcessorFamilyTurion64                                     = 0x86  // Turion 64
	ProcessorFamilyDualCoreOpteron                              = 0x87  // Dual-Core Opteron
	ProcessorFamilyAthlon64X2                                   = 0x88  // Athlon 64 X2
	ProcessorFamilyTurion64X2                                   = 0x89  // Turion 64 X2
	ProcessorFamilyQuadCoreOpteron                              = 0x8a  // Quad-Core Opteron
	ProcessorFamilyThirdGenerationOpteron                       = 0x8b  // Third-Generation Opteron
	ProcessorFamilyPhenomFX                                     = 0x8c  // Phenom FX
	ProcessorFamilyPhenomX4                                     = 0x8d  // Phenom X4
	ProcessorFamilyPhenomX2                                     = 0x8e  // Phenom X2
	ProcessorFamilyAthlonX2                                     = 0x8f  // Athlon X2
	ProcessorFamilyPARISC                                       = 0x90  // PA-RISC
	ProcessorFamilyPARISC8500                                   = 0x91  // PA-RISC 8500
	ProcessorFamilyPARISC8000                                   = 0x92  // PA-RISC 8000
	ProcessorFamilyPARISC7300LC                                 = 0x93  // PA-RISC 7300LC
	ProcessorFamilyPARISC7200                                   = 0x94  // PA-RISC 7200
	ProcessorFamilyPARISC7100LC                                 = 0x95  // PA-RISC 7100LC
	ProcessorFamilyPARISC7100                                   = 0x96  // PA-RISC 7100
	ProcessorFamilyV30                                          = 0xa0  // V30
	ProcessorFamilyQuadCoreXeon3200                             = 0xa1  // Quad-Core Xeon 3200
	ProcessorFamilyDualCoreXeon3000                             = 0xa2  // Dual-Core Xeon 3000
	ProcessorFamilyQuadCoreXeon5300                             = 0xa3  // Quad-Core Xeon 5300
	ProcessorFamilyDualCoreXeon5100                             = 0xa4  // Dual-Core Xeon 5100
	ProcessorFamilyDualCoreXeon5000                             = 0xa5  // Dual-Core Xeon 5000
	ProcessorFamilyDualCoreXeonLV                               = 0xa6  // Dual-Core Xeon LV
	ProcessorFamilyDualCoreXeonULV                              = 0xa7  // Dual-Core Xeon ULV
	ProcessorFamilyDualCoreXeon7100                             = 0xa8  // Dual-Core Xeon 7100
	ProcessorFamilyQuadCoreXeon5400                             = 0xa9  // Quad-Core Xeon 5400
	ProcessorFamilyQuadCoreXeon                                 = 0xaa  // Quad-Core Xeon
	ProcessorFamilyDualCoreXeon5200                             = 0xab  // Dual-Core Xeon 5200
	ProcessorFamilyDualCoreXeon7200                             = 0xac  // Dual-Core Xeon 7200
	ProcessorFamilyQuadCoreXeon7300                             = 0xad  // Quad-Core Xeon 7300
	ProcessorFamilyQuadCoreXeon7400                             = 0xae  // Quad-Core Xeon 7400
	ProcessorFamilyMultiCoreXeon7400                            = 0xaf  // Multi-Core Xeon 7400
	ProcessorFamilyPentiumIIIXeon                               = 0xb0  // Pentium III Xeon
	ProcessorFamilyPentiumIIISpeedstep                          = 0xb1  // Pentium III Speedstep
	ProcessorFamilyPentium4                                     = 0xb2  // Pentium 4
	ProcessorFamilyXeon                                         = 0xb3  // Xeon
	ProcessorFamilyAS400                                        = 0xb4  // AS400
	ProcessorFamilyXeonMP                                       = 0xb5  // Xeon MP
	ProcessorFamilyAthlonXP                                     = 0xb6  // Athlon XP
	ProcessorFamilyAthlonMP                                     = 0xb7  // Athlon MP
	ProcessorFamilyItanium2                                     = 0xb8  // Itanium 2
	ProcessorFamilyPentiumM                                     = 0xb9  // Pentium M
	ProcessorFamilyCeleronD                                     = 0xba  // Celeron D
	ProcessorFamilyPentiumD                                     = 0xbb  // Pentium D
	ProcessorFamilyPentiumEE                                    = 0xbc  // Pentium EE
	ProcessorFamilyCoreSolo                                     = 0xbd  // Core Solo
	ProcessorFamilyHandledAsASpecialCase                        = 0xbe  // handled as a special case */
	ProcessorFamilyCore2Duo                                     = 0xbf  // Core 2 Duo
	ProcessorFamilyCore2Solo                                    = 0xc0  // Core 2 Solo
	ProcessorFamilyCore2Extreme                                 = 0xc1  // Core 2 Extreme
	ProcessorFamilyCore2Quad                                    = 0xc2  // Core 2 Quad
	ProcessorFamilyCore2ExtremeMobile                           = 0xc3  // Core 2 Extreme Mobile
	ProcessorFamilyCore2DuoMobile                               = 0xc4  // Core 2 Duo Mobile
	ProcessorFamilyCore2SoloMobile                              = 0xc5  // Core 2 Solo Mobile
	ProcessorFamilyCoreI7                                       = 0xc6  // Core i7
	ProcessorFamilyDualCoreCeleron                              = 0xc7  // Dual-Core Celeron
	ProcessorFamilyIBM390                                       = 0xc8  // IBM390
	ProcessorFamilyG4                                           = 0xc9  // G4
	ProcessorFamilyG5                                           = 0xca  // G5
	ProcessorFamilyESA390G6                                     = 0xcb  // ESA/390 G6
	ProcessorFamilyZarchitecture                                = 0xcc  // z/Architecture
	ProcessorFamilyCoreI5                                       = 0xcd  // Core i5
	ProcessorFamilyCoreI3                                       = 0xce  // Core i3
	ProcessorFamilyCoreI9                                       = 0xcf  // Core i9
	ProcessorFamilyC7M                                          = 0xd2  // C7-M
	ProcessorFamilyC7D                                          = 0xd3  // C7-D
	ProcessorFamilyC7                                           = 0xd4  // C7
	ProcessorFamilyEden                                         = 0xd5  // Eden
	ProcessorFamilyMultiCoreXeon                                = 0xd6  // Multi-Core Xeon
	ProcessorFamilyDualCoreXeon3xxx                             = 0xd7  // Dual-Core Xeon 3xxx
	ProcessorFamilyQuadCoreXeon3xxx                             = 0xd8  // Quad-Core Xeon 3xxx
	ProcessorFamilyNano                                         = 0xd9  // Nano
	ProcessorFamilyDualCoreXeon5xxx                             = 0xda  // Dual-Core Xeon 5xxx
	ProcessorFamilyQuadCoreXeon5xxx                             = 0xdb  // Quad-Core Xeon 5xxx
	ProcessorFamilyDualCoreXeon7xxx                             = 0xdd  // Dual-Core Xeon 7xxx
	ProcessorFamilyQuadCoreXeon7xxx                             = 0xde  // Quad-Core Xeon 7xxx
	ProcessorFamilyMultiCoreXeon7xxx                            = 0xdf  // Multi-Core Xeon 7xxx
	ProcessorFamilyMultiCoreXeon3400                            = 0xe0  // Multi-Core Xeon 3400
	ProcessorFamilyOpteron3000                                  = 0xe4  // Opteron 3000
	ProcessorFamilySempronII                                    = 0xe5  // Sempron II
	ProcessorFamilyEmbeddedOpteronQuadCore                      = 0xe6  // Embedded Opteron Quad-Core
	ProcessorFamilyPhenomTripleCore                             = 0xe7  // Phenom Triple-Core
	ProcessorFamilyTurionUltraDualCoreMobile                    = 0xe8  // Turion Ultra Dual-Core Mobile
	ProcessorFamilyTurionDualCoreMobile                         = 0xe9  // Turion Dual-Core Mobile
	ProcessorFamilyAthlonDualCore                               = 0xea  // Athlon Dual-Core
	ProcessorFamilySempronSI                                    = 0xeb  // Sempron SI
	ProcessorFamilyPhenomII                                     = 0xec  // Phenom II
	ProcessorFamilyAthlonII                                     = 0xed  // Athlon II
	ProcessorFamilySixCoreOpteron                               = 0xee  // Six-Core Opteron
	ProcessorFamilySempronM                                     = 0xef  // Sempron M
	ProcessorFamilyI860                                         = 0xfa  // i860
	ProcessorFamilyI960                                         = 0xfb  // i960
	ProcessorFamilyARMv7                                        = 0x100 // ARMv7
	ProcessorFamilyARMv8                                        = 0x101 // ARMv8
	ProcessorFamilySH3                                          = 0x104 // SH-3
	ProcessorFamilySH4                                          = 0x105 // SH-4
	ProcessorFamilyARM                                          = 0x118 // ARM
	ProcessorFamilyStrongARM                                    = 0x119 // StrongARM
	ProcessorFamily6x86                                         = 0x12c // 6x86
	ProcessorFamilyMediaGX                                      = 0x12d // MediaGX
	ProcessorFamilyMII                                          = 0x12e // MII
	ProcessorFamilyWinChip                                      = 0x140 // WinChip
	ProcessorFamilyDSP                                          = 0x15e // DSP
	ProcessorFamilyVideoProcessor                               = 0x1f4 // Video Processor
)

func (v ProcessorFamily) String() string {
	switch v {
	case ProcessorFamilyOther:
		return "Other"
	case ProcessorFamilyUnknown:
		return "Unknown"
	case ProcessorFamily8086:
		return "8086"
	case ProcessorFamily80286:
		return "80286"
	case ProcessorFamily80386:
		return "80386"
	case ProcessorFamily80486:
		return "80486"
	case ProcessorFamily8087:
		return "8087"
	case ProcessorFamily80287:
		return "80287"
	case ProcessorFamily80387:
		return "80387"
	case ProcessorFamily80487:
		return "80487"
	case ProcessorFamilyPentium:
		return "Pentium"
	case ProcessorFamilyPentiumPro:
		return "Pentium Pro"
	case ProcessorFamilyPentiumII:
		return "Pentium II"
	case ProcessorFamilyPentiumMMX:
		return "Pentium MMX"
	case ProcessorFamilyCeleron:
		return "Celeron"
	case ProcessorFamilyPentiumIIXeon:
		return "Pentium II Xeon"
	case ProcessorFamilyPentiumIII:
		return "Pentium III"
	case ProcessorFamilyM1:
		return "M1"
	case ProcessorFamilyM2:
		return "M2"
	case ProcessorFamilyCeleronM:
		return "Celeron M"
	case ProcessorFamilyPentium4HT:
		return "Pentium 4 HT"
	case ProcessorFamilyDuron:
		return "Duron"
	case ProcessorFamilyK5:
		return "K5"
	case ProcessorFamilyK6:
		return "K6"
	case ProcessorFamilyK62:
		return "K6-2"
	case ProcessorFamilyK63:
		return "K6-3"
	case ProcessorFamilyAthlon:
		return "Athlon"
	case ProcessorFamilyAMD29000:
		return "AMD29000"
	case ProcessorFamilyK62Plus:
		return "K6-2+"
	case ProcessorFamilyPowerPC:
		return "Power PC"
	case ProcessorFamilyPowerPC601:
		return "Power PC 601"
	case ProcessorFamilyPowerPC603:
		return "Power PC 603"
	case ProcessorFamilyPowerPC603Plus:
		return "Power PC 603+"
	case ProcessorFamilyPowerPC604:
		return "Power PC 604"
	case ProcessorFamilyPowerPC620:
		return "Power PC 620"
	case ProcessorFamilyPowerPCX704:
		return "Power PC x704"
	case ProcessorFamilyPowerPC750:
		return "Power PC 750"
	case ProcessorFamilyCoreDuo:
		return "Core Duo"
	case ProcessorFamilyCoreDuoMobile:
		return "Core Duo Mobile"
	case ProcessorFamilyCoreSoloMobile:
		return "Core Solo Mobile"
	case ProcessorFamilyAtom:
		return "Atom"
	case ProcessorFamilyCoreM:
		return "Core M"
	case ProcessorFamilyCoreM3:
		return "Core m3"
	case ProcessorFamilyCoreM5:
		return "Core m5"
	case ProcessorFamilyCoreM7:
		return "Core m7"
	case ProcessorFamilyAlpha:
		return "Alpha"
	case ProcessorFamilyAlpha21064:
		return "Alpha 21064"
	case ProcessorFamilyAlpha21066:
		return "Alpha 21066"
	case ProcessorFamilyAlpha21164:
		return "Alpha 21164"
	case ProcessorFamilyAlpha21164PC:
		return "Alpha 21164PC"
	case ProcessorFamilyAlpha21164a:
		return "Alpha 21164a"
	case ProcessorFamilyAlpha21264:
		return "Alpha 21264"
	case ProcessorFamilyAlpha21364:
		return "Alpha 21364"
	case ProcessorFamilyTurionIIUltraDualCoreMobileM:
		return "Turion II Ultra Dual-Core Mobile M"
	case ProcessorFamilyTurionIIDualCoreMobileM:
		return "Turion II Dual-Core Mobile M"
	case ProcessorFamilyAthlonIIDualCoreM:
		return "Athlon II Dual-Core M"
	case ProcessorFamilyOpteron6100:
		return "Opteron 6100"
	case ProcessorFamilyOpteron4100:
		return "Opteron 4100"
	case ProcessorFamilyOpteron6200:
		return "Opteron 6200"
	case ProcessorFamilyOpteron4200:
		return "Opteron 4200"
	case ProcessorFamilyFX:
		return "FX"
	case ProcessorFamilyMIPS:
		return "MIPS"
	case ProcessorFamilyMIPSR4000:
		return "MIPS R4000"
	case ProcessorFamilyMIPSR4200:
		return "MIPS R4200"
	case ProcessorFamilyMIPSR4400:
		return "MIPS R4400"
	case ProcessorFamilyMIPSR4600:
		return "MIPS R4600"
	case ProcessorFamilyMIPSR10000:
		return "MIPS R10000"
	case ProcessorFamilyCSeries:
		return "C-Series"
	case ProcessorFamilyESeries:
		return "E-Series"
	case ProcessorFamilyASeries:
		return "A-Series"
	case ProcessorFamilyGSeries:
		return "G-Series"
	case ProcessorFamilyZSeries:
		return "Z-Series"
	case ProcessorFamilyRSeries:
		return "R-Series"
	case ProcessorFamilyOpteron4300:
		return "Opteron 4300"
	case ProcessorFamilyOpteron6300:
		return "Opteron 6300"
	case ProcessorFamilyOpteron3300:
		return "Opteron 3300"
	case ProcessorFamilyFirePro:
		return "FirePro"
	case ProcessorFamilySPARC:
		return "SPARC"
	case ProcessorFamilySuperSPARC:
		return "SuperSPARC"
	case ProcessorFamilyMicroSPARCII:
		return "MicroSPARC II"
	case ProcessorFamilyMicroSPARCIIep:
		return "MicroSPARC IIep"
	case ProcessorFamilyUltraSPARC:
		return "UltraSPARC"
	case ProcessorFamilyUltraSPARCII:
		return "UltraSPARC II"
	case ProcessorFamilyUltraSPARCIIi:
		return "UltraSPARC IIi"
	case ProcessorFamilyUltraSPARCIII:
		return "UltraSPARC III"
	case ProcessorFamilyUltraSPARCIIIi:
		return "UltraSPARC IIIi"
	case ProcessorFamily68040:
		return "68040"
	case ProcessorFamily68xxx:
		return "68xxx"
	case ProcessorFamily68000:
		return "68000"
	case ProcessorFamily68010:
		return "68010"
	case ProcessorFamily68020:
		return "68020"
	case ProcessorFamily68030:
		return "68030"
	case ProcessorFamilyAthlonX4:
		return "Athlon X4"
	case ProcessorFamilyOpteronX1000:
		return "Opteron X1000"
	case ProcessorFamilyOpteronX2000:
		return "Opteron X2000"
	case ProcessorFamilyOpteronASeries:
		return "Opteron A-Series"
	case ProcessorFamilyOpteronX3000:
		return "Opteron X3000"
	case ProcessorFamilyZen:
		return "Zen"
	case ProcessorFamilyHobbit:
		return "Hobbit"
	case ProcessorFamilyCrusoeTM5000:
		return "Crusoe TM5000"
	case ProcessorFamilyCrusoeTM3000:
		return "Crusoe TM3000"
	case ProcessorFamilyEfficeonTM8000:
		return "Efficeon TM8000"
	case ProcessorFamilyWeitek:
		return "Weitek"
	case ProcessorFamilyItanium:
		return "Itanium"
	case ProcessorFamilyAthlon64:
		return "Athlon 64"
	case ProcessorFamilyOpteron:
		return "Opteron"
	case ProcessorFamilySempron:
		return "Sempron"
	case ProcessorFamilyTurion64:
		return "Turion 64"
	case ProcessorFamilyDualCoreOpteron:
		return "Dual-Core Opteron"
	case ProcessorFamilyAthlon64X2:
		return "Athlon 64 X2"
	case ProcessorFamilyTurion64X2:
		return "Turion 64 X2"
	case ProcessorFamilyQuadCoreOpteron:
		return "Quad-Core Opteron"
	case ProcessorFamilyThirdGenerationOpteron:
		return "Third-Generation Opteron"
	case ProcessorFamilyPhenomFX:
		return "Phenom FX"
	case ProcessorFamilyPhenomX4:
		return "Phenom X4"
	case ProcessorFamilyPhenomX2:
		return "Phenom X2"
	case ProcessorFamilyAthlonX2:
		return "Athlon X2"
	case ProcessorFamilyPARISC:
		return "PA-RISC"
	case ProcessorFamilyPARISC8500:
		return "PA-RISC 8500"
	case ProcessorFamilyPARISC8000:
		return "PA-RISC 8000"
	case ProcessorFamilyPARISC7300LC:
		return "PA-RISC 7300LC"
	case ProcessorFamilyPARISC7200:
		return "PA-RISC 7200"
	case ProcessorFamilyPARISC7100LC:
		return "PA-RISC 7100LC"
	case ProcessorFamilyPARISC7100:
		return "PA-RISC 7100"
	case ProcessorFamilyV30:
		return "V30"
	case ProcessorFamilyQuadCoreXeon3200:
		return "Quad-Core Xeon 3200"
	case ProcessorFamilyDualCoreXeon3000:
		return "Dual-Core Xeon 3000"
	case ProcessorFamilyQuadCoreXeon5300:
		return "Quad-Core Xeon 5300"
	case ProcessorFamilyDualCoreXeon5100:
		return "Dual-Core Xeon 5100"
	case ProcessorFamilyDualCoreXeon5000:
		return "Dual-Core Xeon 5000"
	case ProcessorFamilyDualCoreXeonLV:
		return "Dual-Core Xeon LV"
	case ProcessorFamilyDualCoreXeonULV:
		return "Dual-Core Xeon ULV"
	case ProcessorFamilyDualCoreXeon7100:
		return "Dual-Core Xeon 7100"
	case ProcessorFamilyQuadCoreXeon5400:
		return "Quad-Core Xeon 5400"
	case ProcessorFamilyQuadCoreXeon:
		return "Quad-Core Xeon"
	case ProcessorFamilyDualCoreXeon5200:
		return "Dual-Core Xeon 5200"
	case ProcessorFamilyDualCoreXeon7200:
		return "Dual-Core Xeon 7200"
	case ProcessorFamilyQuadCoreXeon7300:
		return "Quad-Core Xeon 7300"
	case ProcessorFamilyQuadCoreXeon7400:
		return "Quad-Core Xeon 7400"
	case ProcessorFamilyMultiCoreXeon7400:
		return "Multi-Core Xeon 7400"
	case ProcessorFamilyPentiumIIIXeon:
		return "Pentium III Xeon"
	case ProcessorFamilyPentiumIIISpeedstep:
		return "Pentium III Speedstep"
	case ProcessorFamilyPentium4:
		return "Pentium 4"
	case ProcessorFamilyXeon:
		return "Xeon"
	case ProcessorFamilyAS400:
		return "AS400"
	case ProcessorFamilyXeonMP:
		return "Xeon MP"
	case ProcessorFamilyAthlonXP:
		return "Athlon XP"
	case ProcessorFamilyAthlonMP:
		return "Athlon MP"
	case ProcessorFamilyItanium2:
		return "Itanium 2"
	case ProcessorFamilyPentiumM:
		return "Pentium M"
	case ProcessorFamilyCeleronD:
		return "Celeron D"
	case ProcessorFamilyPentiumD:
		return "Pentium D"
	case ProcessorFamilyPentiumEE:
		return "Pentium EE"
	case ProcessorFamilyCoreSolo:
		return "Core Solo"
	case ProcessorFamilyHandledAsASpecialCase:
		return "handled as a special case */"
	case ProcessorFamilyCore2Duo:
		return "Core 2 Duo"
	case ProcessorFamilyCore2Solo:
		return "Core 2 Solo"
	case ProcessorFamilyCore2Extreme:
		return "Core 2 Extreme"
	case ProcessorFamilyCore2Quad:
		return "Core 2 Quad"
	case ProcessorFamilyCore2ExtremeMobile:
		return "Core 2 Extreme Mobile"
	case ProcessorFamilyCore2DuoMobile:
		return "Core 2 Duo Mobile"
	case ProcessorFamilyCore2SoloMobile:
		return "Core 2 Solo Mobile"
	case ProcessorFamilyCoreI7:
		return "Core i7"
	case ProcessorFamilyDualCoreCeleron:
		return "Dual-Core Celeron"
	case ProcessorFamilyIBM390:
		return "IBM390"
	case ProcessorFamilyG4:
		return "G4"
	case ProcessorFamilyG5:
		return "G5"
	case ProcessorFamilyESA390G6:
		return "ESA/390 G6"
	case ProcessorFamilyZarchitecture:
		return "z/Architecture"
	case ProcessorFamilyCoreI5:
		return "Core i5"
	case ProcessorFamilyCoreI3:
		return "Core i3"
	case ProcessorFamilyCoreI9:
		return "Core i9"
	case ProcessorFamilyC7M:
		return "C7-M"
	case ProcessorFamilyC7D:
		return "C7-D"
	case ProcessorFamilyC7:
		return "C7"
	case ProcessorFamilyEden:
		return "Eden"
	case ProcessorFamilyMultiCoreXeon:
		return "Multi-Core Xeon"
	case ProcessorFamilyDualCoreXeon3xxx:
		return "Dual-Core Xeon 3xxx"
	case ProcessorFamilyQuadCoreXeon3xxx:
		return "Quad-Core Xeon 3xxx"
	case ProcessorFamilyNano:
		return "Nano"
	case ProcessorFamilyDualCoreXeon5xxx:
		return "Dual-Core Xeon 5xxx"
	case ProcessorFamilyQuadCoreXeon5xxx:
		return "Quad-Core Xeon 5xxx"
	case ProcessorFamilyDualCoreXeon7xxx:
		return "Dual-Core Xeon 7xxx"
	case ProcessorFamilyQuadCoreXeon7xxx:
		return "Quad-Core Xeon 7xxx"
	case ProcessorFamilyMultiCoreXeon7xxx:
		return "Multi-Core Xeon 7xxx"
	case ProcessorFamilyMultiCoreXeon3400:
		return "Multi-Core Xeon 3400"
	case ProcessorFamilyOpteron3000:
		return "Opteron 3000"
	case ProcessorFamilySempronII:
		return "Sempron II"
	case ProcessorFamilyEmbeddedOpteronQuadCore:
		return "Embedded Opteron Quad-Core"
	case ProcessorFamilyPhenomTripleCore:
		return "Phenom Triple-Core"
	case ProcessorFamilyTurionUltraDualCoreMobile:
		return "Turion Ultra Dual-Core Mobile"
	case ProcessorFamilyTurionDualCoreMobile:
		return "Turion Dual-Core Mobile"
	case ProcessorFamilyAthlonDualCore:
		return "Athlon Dual-Core"
	case ProcessorFamilySempronSI:
		return "Sempron SI"
	case ProcessorFamilyPhenomII:
		return "Phenom II"
	case ProcessorFamilyAthlonII:
		return "Athlon II"
	case ProcessorFamilySixCoreOpteron:
		return "Six-Core Opteron"
	case ProcessorFamilySempronM:
		return "Sempron M"
	case ProcessorFamilyI860:
		return "i860"
	case ProcessorFamilyI960:
		return "i960"
	case ProcessorFamilyARMv7:
		return "ARMv7"
	case ProcessorFamilyARMv8:
		return "ARMv8"
	case ProcessorFamilySH3:
		return "SH-3"
	case ProcessorFamilySH4:
		return "SH-4"
	case ProcessorFamilyARM:
		return "ARM"
	case ProcessorFamilyStrongARM:
		return "StrongARM"
	case ProcessorFamily6x86:
		return "6x86"
	case ProcessorFamilyMediaGX:
		return "MediaGX"
	case ProcessorFamilyMII:
		return "MII"
	case ProcessorFamilyWinChip:
		return "WinChip"
	case ProcessorFamilyDSP:
		return "DSP"
	case ProcessorFamilyVideoProcessor:
		return "Video Processor"
	}
	return fmt.Sprintf("%d", v)
}

// ProcessorStatus is defined in DSP0134 7.5.
type ProcessorStatus uint8

var processorStatusStr = []string{
	"Unknown", "Enabled", "Disabled By User", "Disabled By BIOS", "Idle", "Reserved5", "Reserved6", "Other",
}

func (v ProcessorStatus) String() string {
	if v&0x40 == 0 {
		return "Unpopulated"
	}
	return "Populated, " + processorStatusStr[v&7]
}

// ProcessorUpgrade is defined in DSP0134 7.5.5.
type ProcessorUpgrade uint8

// ProcessorUpgrade values are defined in DSP0134 7.5.5.
const (
	ProcessorUpgradeOther                ProcessorUpgrade = 0x01 // Other
	ProcessorUpgradeUnknown                               = 0x02 // Unknown
	ProcessorUpgradeDaughterBoard                         = 0x03 // Daughter Board
	ProcessorUpgradeZIFSocket                             = 0x04 // ZIF Socket
	ProcessorUpgradeReplaceablePiggyBack                  = 0x05 // Replaceable Piggy Back
	ProcessorUpgradeNone                                  = 0x06 // None
	ProcessorUpgradeLIFSocket                             = 0x07 // LIF Socket
	ProcessorUpgradeSlot1                                 = 0x08 // Slot 1
	ProcessorUpgradeSlot2                                 = 0x09 // Slot 2
	ProcessorUpgrade370pinSocket                          = 0x0a // 370-pin Socket
	ProcessorUpgradeSlotA                                 = 0x0b // Slot A
	ProcessorUpgradeSlotM                                 = 0x0c // Slot M
	ProcessorUpgradeSocket423                             = 0x0d // Socket 423
	ProcessorUpgradeSocketA                               = 0x0e // Socket A (Socket 462)
	ProcessorUpgradeSocket478                             = 0x0f // Socket 478
	ProcessorUpgradeSocket754                             = 0x10 // Socket 754
	ProcessorUpgradeSocket940                             = 0x11 // Socket 940
	ProcessorUpgradeSocket939                             = 0x12 // Socket 939
	ProcessorUpgradeSocketMpga604                         = 0x13 // Socket mPGA604
	ProcessorUpgradeSocketLGA771                          = 0x14 // Socket LGA771
	ProcessorUpgradeSocketLGA775                          = 0x15 // Socket LGA775
	ProcessorUpgradeSocketS1                              = 0x16 // Socket S1
	ProcessorUpgradeSocketAM2                             = 0x17 // Socket AM2
	ProcessorUpgradeSocketF1207                           = 0x18 // Socket F (1207)
	ProcessorUpgradeSocketLGA1366                         = 0x19 // Socket LGA1366
	ProcessorUpgradeSocketG34                             = 0x1a // Socket G34
	ProcessorUpgradeSocketAM3                             = 0x1b // Socket AM3
	ProcessorUpgradeSocketC32                             = 0x1c // Socket C32
	ProcessorUpgradeSocketLGA1156                         = 0x1d // Socket LGA1156
	ProcessorUpgradeSocketLGA1567                         = 0x1e // Socket LGA1567
	ProcessorUpgradeSocketPGA988A                         = 0x1f // Socket PGA988A
	ProcessorUpgradeSocketBGA1288                         = 0x20 // Socket BGA1288
	ProcessorUpgradeSocketRpga988b                        = 0x21 // Socket rPGA988B
	ProcessorUpgradeSocketBGA1023                         = 0x22 // Socket BGA1023
	ProcessorUpgradeSocketBGA1224                         = 0x23 // Socket BGA1224
	ProcessorUpgradeSocketBGA1155                         = 0x24 // Socket BGA1155
	ProcessorUpgradeSocketLGA1356                         = 0x25 // Socket LGA1356
	ProcessorUpgradeSocketLGA2011                         = 0x26 // Socket LGA2011
	ProcessorUpgradeSocketFS1                             = 0x27 // Socket FS1
	ProcessorUpgradeSocketFS2                             = 0x28 // Socket FS2
	ProcessorUpgradeSocketFM1                             = 0x29 // Socket FM1
	ProcessorUpgradeSocketFM2                             = 0x2a // Socket FM2
	ProcessorUpgradeSocketLGA20113                        = 0x2b // Socket LGA2011-3
	ProcessorUpgradeSocketLGA13563                        = 0x2c // Socket LGA1356-3
	ProcessorUpgradeSocketLGA1150                         = 0x2d // Socket LGA1150
	ProcessorUpgradeSocketBGA1168                         = 0x2e // Socket BGA1168
	ProcessorUpgradeSocketBGA1234                         = 0x2f // Socket BGA1234
	ProcessorUpgradeSocketBGA1364                         = 0x30 // Socket BGA1364
	ProcessorUpgradeSocketAM4                             = 0x31 // Socket AM4
	ProcessorUpgradeSocketLGA1151                         = 0x32 // Socket LGA1151
	ProcessorUpgradeSocketBGA1356                         = 0x33 // Socket BGA1356
	ProcessorUpgradeSocketBGA1440                         = 0x34 // Socket BGA1440
	ProcessorUpgradeSocketBGA1515                         = 0x35 // Socket BGA1515
	ProcessorUpgradeSocketLGA36471                        = 0x36 // Socket LGA3647-1
	ProcessorUpgradeSocketSP3                             = 0x37 // Socket SP3
	ProcessorUpgradeSocketSP3r2                           = 0x38 // Socket SP3r2
	ProcessorUpgradeSocketLGA2066                         = 0x39 // Socket LGA2066
	ProcessorUpgradeSocketBGA1392                         = 0x3a // Socket BGA1392
	ProcessorUpgradeSocketBGA1510                         = 0x3b // Socket BGA1510
	ProcessorUpgradeSocketBGA1528                         = 0x3c // Socket BGA1528
)

func (v ProcessorUpgrade) String() string {
	switch v {
	case ProcessorUpgradeOther:
		return "Other"
	case ProcessorUpgradeUnknown:
		return "Unknown"
	case ProcessorUpgradeDaughterBoard:
		return "Daughter Board"
	case ProcessorUpgradeZIFSocket:
		return "ZIF Socket"
	case ProcessorUpgradeReplaceablePiggyBack:
		return "Replaceable Piggy Back"
	case ProcessorUpgradeNone:
		return "None"
	case ProcessorUpgradeLIFSocket:
		return "LIF Socket"
	case ProcessorUpgradeSlot1:
		return "Slot 1"
	case ProcessorUpgradeSlot2:
		return "Slot 2"
	case ProcessorUpgrade370pinSocket:
		return "370-pin Socket"
	case ProcessorUpgradeSlotA:
		return "Slot A"
	case ProcessorUpgradeSlotM:
		return "Slot M"
	case ProcessorUpgradeSocket423:
		return "Socket 423"
	case ProcessorUpgradeSocketA:
		return "Socket A (Socket 462)"
	case ProcessorUpgradeSocket478:
		return "Socket 478"
	case ProcessorUpgradeSocket754:
		return "Socket 754"
	case ProcessorUpgradeSocket940:
		return "Socket 940"
	case ProcessorUpgradeSocket939:
		return "Socket 939"
	case ProcessorUpgradeSocketMpga604:
		return "Socket mPGA604"
	case ProcessorUpgradeSocketLGA771:
		return "Socket LGA771"
	case ProcessorUpgradeSocketLGA775:
		return "Socket LGA775"
	case ProcessorUpgradeSocketS1:
		return "Socket S1"
	case ProcessorUpgradeSocketAM2:
		return "Socket AM2"
	case ProcessorUpgradeSocketF1207:
		return "Socket F (1207)"
	case ProcessorUpgradeSocketLGA1366:
		return "Socket LGA1366"
	case ProcessorUpgradeSocketG34:
		return "Socket G34"
	case ProcessorUpgradeSocketAM3:
		return "Socket AM3"
	case ProcessorUpgradeSocketC32:
		return "Socket C32"
	case ProcessorUpgradeSocketLGA1156:
		return "Socket LGA1156"
	case ProcessorUpgradeSocketLGA1567:
		return "Socket LGA1567"
	case ProcessorUpgradeSocketPGA988A:
		return "Socket PGA988A"
	case ProcessorUpgradeSocketBGA1288:
		return "Socket BGA1288"
	case ProcessorUpgradeSocketRpga988b:
		return "Socket rPGA988B"
	case ProcessorUpgradeSocketBGA1023:
		return "Socket BGA1023"
	case ProcessorUpgradeSocketBGA1224:
		return "Socket BGA1224"
	case ProcessorUpgradeSocketBGA1155:
		return "Socket BGA1155"
	case ProcessorUpgradeSocketLGA1356:
		return "Socket LGA1356"
	case ProcessorUpgradeSocketLGA2011:
		return "Socket LGA2011"
	case ProcessorUpgradeSocketFS1:
		return "Socket FS1"
	case ProcessorUpgradeSocketFS2:
		return "Socket FS2"
	case ProcessorUpgradeSocketFM1:
		return "Socket FM1"
	case ProcessorUpgradeSocketFM2:
		return "Socket FM2"
	case ProcessorUpgradeSocketLGA20113:
		return "Socket LGA2011-3"
	case ProcessorUpgradeSocketLGA13563:
		return "Socket LGA1356-3"
	case ProcessorUpgradeSocketLGA1150:
		return "Socket LGA1150"
	case ProcessorUpgradeSocketBGA1168:
		return "Socket BGA1168"
	case ProcessorUpgradeSocketBGA1234:
		return "Socket BGA1234"
	case ProcessorUpgradeSocketBGA1364:
		return "Socket BGA1364"
	case ProcessorUpgradeSocketAM4:
		return "Socket AM4"
	case ProcessorUpgradeSocketLGA1151:
		return "Socket LGA1151"
	case ProcessorUpgradeSocketBGA1356:
		return "Socket BGA1356"
	case ProcessorUpgradeSocketBGA1440:
		return "Socket BGA1440"
	case ProcessorUpgradeSocketBGA1515:
		return "Socket BGA1515"
	case ProcessorUpgradeSocketLGA36471:
		return "Socket LGA3647-1"
	case ProcessorUpgradeSocketSP3:
		return "Socket SP3"
	case ProcessorUpgradeSocketSP3r2:
		return "Socket SP3r2"
	case ProcessorUpgradeSocketLGA2066:
		return "Socket LGA2066"
	case ProcessorUpgradeSocketBGA1392:
		return "Socket BGA1392"
	case ProcessorUpgradeSocketBGA1510:
		return "Socket BGA1510"
	case ProcessorUpgradeSocketBGA1528:
		return "Socket BGA1528"
	}
	return fmt.Sprintf("%d", v)
}

// ProcessorCharacteristics values are defined in DSP0134 7.5.9.
type ProcessorCharacteristics uint16

// ProcessorCharacteristics fields are defined in DSP0134 x.x.x
const (
	ProcessorCharacteristicsReserved                ProcessorCharacteristics = (1 << 0) // Reserved
	ProcessorCharacteristicsUnknown                                          = (1 << 1) // Unknown
	ProcessorCharacteristics64bitCapable                                     = (1 << 2) // 64-bit Capable
	ProcessorCharacteristicsMultiCore                                        = (1 << 3) // Multi-Core
	ProcessorCharacteristicsHardwareThread                                   = (1 << 4) // Hardware Thread
	ProcessorCharacteristicsExecuteProtection                                = (1 << 5) // Execute Protection
	ProcessorCharacteristicsEnhancedVirtualization                           = (1 << 6) // Enhanced Virtualization
	ProcessorCharacteristicsPowerPerformanceControl                          = (1 << 7) // Power/Performance Control
)

func (v ProcessorCharacteristics) String() string {
	var lines []string
	if v&ProcessorCharacteristicsReserved != 0 {
		lines = append(lines, "Reserved")
	}
	if v&ProcessorCharacteristicsUnknown != 0 {
		lines = append(lines, "Unknown")
	}
	if v&ProcessorCharacteristics64bitCapable != 0 {
		lines = append(lines, "64-bit capable")
	}
	if v&ProcessorCharacteristicsMultiCore != 0 {
		lines = append(lines, "Multi-Core")
	}
	if v&ProcessorCharacteristicsHardwareThread != 0 {
		lines = append(lines, "Hardware Thread")
	}
	if v&ProcessorCharacteristicsExecuteProtection != 0 {
		lines = append(lines, "Execute Protection")
	}
	if v&ProcessorCharacteristicsEnhancedVirtualization != 0 {
		lines = append(lines, "Enhanced Virtualization")
	}
	if v&ProcessorCharacteristicsPowerPerformanceControl != 0 {
		lines = append(lines, "Power/Performance Control")
	}
	return "\t\t" + strings.Join(lines, "\n\t\t")
}
