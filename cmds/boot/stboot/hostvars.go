// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const hostvarsFile = "hostvars.json"

//go:generate jsonenums -type=bootmode
type bootmode int

// bootmodes values defines where to load a bootball from.
const (
	NetworkStatic bootmode = iota
	NetworkDHCP
	LocalStorage
)

func (b bootmode) string() string {
	return []string{"NetworkStatic", "NetworkDHCP", "LocalStorage"}[b]
}

// Hostvars contains platform-specific data.
type Hostvars struct {
	// MinimalSignaturesMatch is the min number of signatures that must pass validation.
	MinimalSignaturesMatch int `json:"minimal_signatures_match"`
	// Fingerprints are used to validate the root certificate insinde the bootball.
	Fingerprints []string `json:"fingerprints"`
	// Timestamp is the UNIX build time of the bootloader
	Timestamp int `json:"build_timestamp"`
	//BootMode
	BootMode bootmode `json:"boot_mode"`
}

// loadHostVars parses hostvars.json file.
// It is expected to be in /etc.
func loadHostvars(path string) (*Hostvars, error) {
	var vars Hostvars
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s due to: %v", path, err)
	}
	if err = json.Unmarshal(data, &vars); err != nil {
		return nil, fmt.Errorf("cannot parse data - invalid hostvars in %s:  %v", path, err)
	}
	return &vars, nil
}
