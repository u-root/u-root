// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msr

const (
	// This is Intel's name. It makes no sense: this register is present
	// in 64-bit CPUs.
	IntelIA32FeatureControl  MSR = 0x3A
	IntelPkgCstConfigControl MSR = 0xE2
	IntelFeatureConfig       MSR = 0x13
	IntelDramPowerLimit      MSR = 0x61
	IntelConfigTDPControl    MSR = 0x64B
	IntelIA32DebugInterface  MSR = 0xC80
)

var Intel = []MSRVal{
	{
		// Architectural MSR. All systems.
		// Enables features like VMX.
		Addr: IntelIA32FeatureControl,
		Name: "IntelIA32FeatureControl",
		Set:  1 << 0, // Locks this register.
	},
	{
		// Silvermont, Airmont, Nehalem...
		// Controls Processor C States.
		Addr: IntelPkgCstConfigControl,
		Name: "IntelPkgCstConfigControl",
		Set:  1 << 15, // Locks this register.
	},
	{
		// Westmere onwards.
		// Note that this turns on AES instructions, however
		// 3 will turn off AES until reset.
		Addr: IntelFeatureConfig,
		Name: "IntelFeatureConfig",
		Set:  1 << 0,
	},
	{
		// Goldmont, SandyBridge
		// Controls DRAM power limits. See Intel SDM
		Addr: IntelDramPowerLimit,
		Name: "IntelDramPowerLimit",
		Set:  1 << 31, // Locks this register.
	},
	{
		// IvyBridge Onwards.
		// Not much information in the SDM, seems to control power limits
		Addr: IntelConfigTDPControl,
		Name: "IntelConfigTDPControl",
		Set:  1 << 31, // Locks this register.
	},
	{
		// Architectural MSR. All systems.
		// This is the actual spelling of the MSR in the manual.
		// Controls availability of silicon debug interfaces
		Addr: IntelIA32DebugInterface,
		Name: "IntelIA32DebugInterface",
		Set:  1 << 30, // Locks this register.
	},
}
