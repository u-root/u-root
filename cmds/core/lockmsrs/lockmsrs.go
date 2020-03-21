// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lockmsrs locks important intel MSRs.
//
// All MSRs are specified in the Intel Software developer's manual.
// This seems like a good set of bits to lock down when booting through NERF/LINUXBOOT
// to some other OS. When locked, these MSRs generally prevent
// further modifications until reset.
package main

import (
	"flag"
	"log"

	"github.com/intel-go/cpuid"
	"github.com/u-root/u-root/pkg/msr"
)

const (
	msrIA32FeatureControl  msr.MSR = 0x3A  // MSR_IA32_FEATURE_CONTROL
	msrPkgCstConfigControl msr.MSR = 0xE2  // MSR_PKG_CST_CONFIG_CONTROL
	msrFeatureConfig       msr.MSR = 0x13C // MSR_FEATURE_CONFIG
	msrDramPowerLimit      msr.MSR = 0x618 // MSR_DRAM_POWER_LIMIT
	msrConfigTDPControl    msr.MSR = 0x64B // MSR_CONFIG_TDP_CONTROL
	ia32DebugInterface     msr.MSR = 0xC80 // IA32_DEBUG_INTERFACE
)

var verbose = flag.Bool("v", false, "verbose mode")

var msrList = []struct {
	msrAdd    msr.MSR
	clearMask uint64
	setMask   uint64
}{
	{
		// Architectural MSR. All systems.
		// Enables features like VMX.
		msrAdd:  msrIA32FeatureControl,
		setMask: 1 << 0, // Locks this register.
	},
	{
		// Silvermont, Airmont, Nehalem...
		// Controls Processor C States.
		msrAdd:  msrPkgCstConfigControl,
		setMask: 1 << 15, // Locks this register.
	},
	{
		// Westmere onwards.
		// Note that this turns on AES instructions, however
		// 3 will turn off AES until reset.
		msrAdd:  msrFeatureConfig,
		setMask: 1 << 0,
	},
	{
		// Goldmont, SandyBridge
		// Controls DRAM power limits. See Intel SDM
		msrAdd:  msrDramPowerLimit,
		setMask: 1 << 31, // Locks this register.
	},
	{
		// IvyBridge Onwards.
		// Not much information in the SDM, seems to control power limits
		msrAdd:  msrConfigTDPControl,
		setMask: 1 << 31, // Locks this register.
	},
	{
		// Architectural MSR. All systems.
		// This is the actual spelling of the MSR in the manual.
		// Controls availability of silicon debug interfaces
		msrAdd:  ia32DebugInterface,
		setMask: 1 << 30, // Locks this register.
	},
}

func main() {
	if cpuid.VendorIdentificatorString != "GenuineIntel" {
		log.Fatalf("Only Intel CPUs supported, expected GenuineIntel, got %v", cpuid.VendorIdentificatorString)
	}

	cpus, err := msr.AllCPUs()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range msrList {
		if *verbose {
			log.Printf("Locking MSR %v on cpus %v, clearmask 0x%8x, setmask 0x%8x", m.msrAdd, cpus, m.clearMask, m.setMask)
		}
		errs := m.msrAdd.TestAndSet(cpus, m.clearMask, m.setMask)
		for i, e := range errs {
			if e != nil {
				// Hope no one ever modifies this slice.
				log.Printf("Error locking msr %v on cpu %v: %v\n", m.msrAdd.String(), cpus[i], e)
			}
		}
	}
}
