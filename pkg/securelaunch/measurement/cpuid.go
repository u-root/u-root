// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package measurement

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	cpuid "github.com/u-root/cpuid"
	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/tpm"
)

const defaultCPUIDFile = "cpuid.txt"

// CPUIDCollector describes the "cpuid" portion of the policy file.
type CPUIDCollector struct {
	Type     string `json:"type"`
	Location string `json:"location"`
}

// NewCPUIDCollector extracts the "cpuid" portion from the policy file and
// initializes a new CPUIDCollector structure.
//
// An error is returned if unmarshalling of CPUIDCollector fails.
func NewCPUIDCollector(config []byte) (Collector, error) {
	slaunch.Debug("New CPUID Collector initialized\n")
	fc := new(CPUIDCollector)
	err := json.Unmarshal(config, &fc)
	if err != nil {
		return nil, err
	}
	return fc, nil
}

// getCPUIDInfo uses a string builder to store data obtained from the
// intel-go/cpuid package.
//
// It returns a byte slice of the string built.
func getCPUIDInfo() []byte {
	var w strings.Builder
	fmt.Fprintf(&w, "VendorString:           %s\n", cpuid.VendorIdentificatorString)
	fmt.Fprintf(&w, "ProcessorBrandString:   %s\n", cpuid.ProcessorBrandString)
	fmt.Fprintf(&w, "SteppingId:     %d\n", cpuid.SteppingId)
	fmt.Fprintf(&w, "ProcessorType:  %d\n", cpuid.ProcessorType)
	fmt.Fprintf(&w, "DisplayFamily:  %d\n", cpuid.DisplayFamily)
	fmt.Fprintf(&w, "DisplayModel:   %d\n", cpuid.DisplayModel)
	fmt.Fprintf(&w, "CacheLineSize:  %d\n", cpuid.CacheLineSize)
	fmt.Fprintf(&w, "MaxLogocalCPUId:%d\n", cpuid.MaxLogicalCPUId)
	fmt.Fprintf(&w, "InitialAPICId:  %d\n", cpuid.InitialAPICId)
	fmt.Fprintf(&w, "Smallest monitor-line size in bytes:  %d\n", cpuid.MonLineSizeMin)
	fmt.Fprintf(&w, "Largest monitor-line size in bytes:   %d\n", cpuid.MonLineSizeMax)
	fmt.Fprintf(&w, "Monitor Interrupt break-event is supported:  %v\n", cpuid.MonitorIBE)
	fmt.Fprintf(&w, "MONITOR/MWAIT extensions are supported:      %v\n", cpuid.MonitorEMX)
	fmt.Fprintf(&w, "AVX state:     %v\n", cpuid.EnabledAVX)
	fmt.Fprintf(&w, "AVX-512 state: %v\n", cpuid.EnabledAVX512)
	fmt.Fprintf(&w, "Interrupt thresholds in digital thermal sensor: %v\n", cpuid.ThermalSensorInterruptThresholds)

	fmt.Fprintf(&w, "Features: ")
	for i := uint64(0); i < 64; i++ {
		if cpuid.HasFeature(1 << i) {
			fmt.Fprintf(&w, "%s ", cpuid.FeatureNames[1<<i])
		}
	}
	fmt.Fprintf(&w, "\n")

	fmt.Fprintf(&w, "ExtendedFeatures: ")
	for i := uint64(0); i < 64; i++ {
		if cpuid.HasExtendedFeature(1 << i) {
			fmt.Fprintf(&w, "%s ", cpuid.ExtendedFeatureNames[1<<i])
		}
	}
	fmt.Fprintf(&w, "\n")

	fmt.Fprintf(&w, "ExtraFeatures: ")
	for i := uint64(0); i < 64; i++ {
		if cpuid.HasExtraFeature(1 << i) {
			fmt.Fprintf(&w, "%s ", cpuid.ExtraFeatureNames[1<<i])
		}
	}
	fmt.Fprintf(&w, "\n")

	fmt.Fprintf(&w, "ThermalAndPowerFeatures: ")
	for i := uint32(0); i < 64; i++ {
		if cpuid.HasThermalAndPowerFeature(1 << i) {
			if name, found := cpuid.ThermalAndPowerFeatureNames[1<<i]; found {
				fmt.Fprintf(&w, "%s ", name)
			}
		}
	}
	fmt.Fprintf(&w, "\n")

	for _, cacheDescription := range cpuid.CacheDescriptors {
		fmt.Fprintf(&w, "CacheDescriptor: %v\n", cacheDescription)
	}

	return []byte(w.String())
}

// measureCPUIDFile extends the CPUIDInfo obtained from the cpuid package into
// a TPM PCR.
func measureCPUIDFile() ([]byte, error) {
	d := getCPUIDInfo() // return strings builder
	eventDesc := "CPUID Collector: Measured CPUID Info"
	if e := tpm.ExtendPCRDebug(pcr, bytes.NewReader(d), eventDesc); e != nil {
		return nil, e
	}

	return d, nil
}

// Collect gets the cpuid and extends it into the TPM and stores a copy on disk.
//
// It satisfies the Collector interface.
func (s *CPUIDCollector) Collect() error {
	d, err := measureCPUIDFile()
	if err != nil {
		log.Printf("CPUID Collector: err = %v", err)
		return err
	}

	return slaunch.AddToPersistQueue("CPUID Collector", d, s.Location, defaultCPUIDFile)
}
