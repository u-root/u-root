// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package policy locates and parses a JSON policy file.
package policy

import (
	"encoding/json"
	"fmt"
	"log"

	slaunch "github.com/u-root/u-root/pkg/securelaunch"
	"github.com/u-root/u-root/pkg/securelaunch/eventlog"
	"github.com/u-root/u-root/pkg/securelaunch/launcher"
	"github.com/u-root/u-root/pkg/securelaunch/measurement"
)

// Policy describes the policy used to drive the security engine.
//
// The policy is stored as a JSON file.
type Policy struct {
	Collectors []measurement.Collector
	Launcher   launcher.Launcher
	EventLog   eventlog.EventLog
}

// policyBytes is a byte slice to hold a copy of the policy file in memory.
var policyBytes []byte

// Load reads the specified policy file.
func Load(policyLocation string) error {
	var err error

	policyBytes, err = slaunch.ReadFile(policyLocation)
	if err != nil {
		return err
	}

	return nil
}

// Parse accepts a JSON file as input, parses it into a well defined Policy
// structure and returns a pointer to the Policy structure.
func Parse() (*Policy, error) {
	if len(policyBytes) == 0 {
		return nil, fmt.Errorf("policy file not yet loaded or empty")
	}

	policy := &Policy{}
	var parse struct {
		Collectors []json.RawMessage `json:"collectors"`
		Attestor   json.RawMessage   `json:"attestor"`
		Launcher   json.RawMessage   `json:"launcher"`
		EventLog   json.RawMessage   `json:"eventlog"`
	}

	if err := json.Unmarshal(policyBytes, &parse); err != nil {
		log.Printf("policy: Error unmarshalling policy file: %v", err)
		return nil, err
	}

	for _, collectors := range parse.Collectors {
		collector, err := measurement.GetCollector(collectors)
		if err != nil {
			log.Printf("policy: Error GetCollector err:c=%s, collector=%v", collectors, collector)
			return nil, err
		}

		policy.Collectors = append(policy.Collectors, collector)
	}

	if len(parse.Launcher) > 0 {
		if err := json.Unmarshal(parse.Launcher, &policy.Launcher); err != nil {
			log.Printf("policy: Error unmarshalling `launcher` section of policy file: %v", err)
			return nil, err
		}
	}

	if len(parse.EventLog) > 0 {
		if err := json.Unmarshal(parse.EventLog, &policy.EventLog); err != nil {
			log.Printf("policy: Error unmarshalling `eventlog` section of policy file: %v", err)
			return nil, err
		}
	}

	return policy, nil
}

// Measure measures the policy file.
func Measure() error {
	if len(policyBytes) == 0 {
		return fmt.Errorf("policy file not yet loaded or empty")
	}

	eventDesc := "File Collector: measured securelaunch policy file"
	if err := measurement.HashBytes(policyBytes, eventDesc); err != nil {
		log.Printf("policy: ERR: could not measure policy file: %v", err)
		return err
	}

	return nil
}
