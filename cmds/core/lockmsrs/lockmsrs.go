// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lock-msrs locks a bunch of intel MSRs, all specified in the Intel
// Software developer's manual. This seems like a good set of bits to
// lock down when booting through NERF/LINUXBOOT to some other OS.
// When locked, these MSRs generally prevent further modifications until
// reset.
package main

import (
	"log"

	"github.com/intel-go/cpuid"
	"github.com/u-root/u-root/pkg/msr"
)

const MSR_IA32_FEATURE_CONTROL uint32 = 0x3A
const MSR_PKG_CST_CONFIG_CONTROL uint32 = 0xE2
const MSR_FEATURE_CONFIG uint32 = 0x13C
const MSR_DRAM_POWER_LIMIT uint32 = 0x618
const MSR_CONFIG_TDP_CONTROL uint32 = 0x64B
const IA32_DEBUG_INTERFACE uint32 = 0xC80

var msrList = []struct {
	msrAdd    uint32
	clearMask uint64
	setMask   uint64
}{
	{
		// Architectural MSR. All systems.
		// Enables features like VMX.
		msrAdd:  MSR_IA32_FEATURE_CONTROL,
		setMask: 1 << 0, // Locks this register.
	},
	{
		// Silvermont, Airmont, Nehalem...
		// Controls Processor C States.
		msrAdd:  MSR_PKG_CST_CONFIG_CONTROL,
		setMask: 1 << 15, // Locks this register.
	},
	{
		// Westmere onwards.
		// Note that this turns on AES instructions, however
		// 3 will turn off AES until reset.
		msrAdd:  MSR_FEATURE_CONFIG,
		setMask: 1 << 0,
	},
	{
		// Goldmont, SandyBridge
		// Controls DRAM power limits. See Intel SDM
		msrAdd:  MSR_DRAM_POWER_LIMIT,
		setMask: 1 << 31, // Locks this register.
	},
	{
		// IvyBridge Onwards.
		// Not much information in the SDM, seems to control power limits
		msrAdd:  MSR_CONFIG_TDP_CONTROL,
		setMask: 1 << 31, // Locks this register.
	},
	{
		// Architectural MSR. All systems.
		// This is the actual spelling of the MSR in the manual.
		// Controls availability of silicon debug interfaces
		msrAdd:  IA32_DEBUG_INTERFACE,
		setMask: 1 << 30, // Locks this register.
	},
}

func main() {
	if cpuid.VendorIdentificatorString != "GenuineIntel" {
		log.Fatalf("Only Intel CPUs supported, expected GenuineIntel, got %v", cpuid.VendorIdentificatorString)
	}
	p := msr.Paths("*")
	var errs []error

	for _, m := range msrList {
		errs = msr.MaskBits(p, m.msrAdd, m.clearMask, m.setMask)
	}
	for i, e := range errs {
		if e != nil {
			log.Printf("%v: %v\n", p[i], e)
		}
	}
}
