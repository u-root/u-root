// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msr

const (
	// This is Intel's name. It makes no sense: this register is present
	// in 64-bit CPUs.
	IntelIA32FeatureControl  MSR = 0x3A  // MSR_IA32_FEATURE_CONTROL
	IntelPkgCstConfigControl MSR = 0xE2  // MSR_PKG_CST_CONFIG_CONTROL
	IntelFeatureConfig       MSR = 0x13c // MSR_FEATURE_CONFIG
	IntelDramPowerLimit      MSR = 0x618 // MSR_DRAM_POWER_LIMIT
	IntelConfigTDPControl    MSR = 0x64B // MSR_CONFIG_TDP_CONTROL
	IntelIA32DebugInterface  MSR = 0xC80 // IA32_DEBUG_INTERFACE
)

var LockIntel = []MSRVal{
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
