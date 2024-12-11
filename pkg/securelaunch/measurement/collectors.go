// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package measurement provides different collectors to hash files, disks, dmi info and cpuid info.
package measurement

import (
	"encoding/json"
	"fmt"
)

// PCR number to extend all measurements taken by the securelaunch package.
var pcr = uint32(22)

// All collectors (e.g., cpuid, dmi, etc.) should satisfy this interface.
// Collectors collect data and extend its hash into a PCR.
type Collector interface {
	Collect() error
}

var supportedCollectors = map[string]func([]byte) (Collector, error){
	"storage": NewStorageCollector,
	"dmi":     NewDmiCollector,
	"files":   NewFileCollector,
	"cpuid":   NewCPUIDCollector,
}

// GetCollector calls the appropriate init handlers for a particular
// collector JSON object argument and returns a new Collector Interface.
//
// An error is returned if unmarshalling fails or an unsupported collector is
// passed as an argument.
func GetCollector(config []byte) (Collector, error) {
	var header struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(config, &header)
	if err != nil {
		fmt.Printf("Measurement: Unmarshal error\n")
		return nil, err
	}

	if init, ok := supportedCollectors[header.Type]; ok {
		return init(config)
	}

	return nil, fmt.Errorf("unsupported collector %s", header.Type)
}

// SetPCR sets the PCR value to use for measurements.
func SetPCR(value uint32) {
	pcr = value
}
