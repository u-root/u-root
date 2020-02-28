// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stboot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// HostVars contains contains platform-specific data
type HostVars struct {
	//IP configuration. If empty DHCP will be used.
	HostIP         string `json:"host_ip"`
	HostNetmask    string `json:"netmask"`
	DefaultGateway string `json:"gateway"`
	DNSServer      string `json:"dns"`
	// MinimalSignaturesMatch is the min number of signatures that must pass validation.
	MinimalSignaturesMatch int `json:"minimal_signatures_match"`
	// Fingerprints are used to validate the root certificate insinde the bootball.
	Fingerprints []string `json:"fingerprints"`
	// Timestamp is the UNIX build time of the bootloader
	Timestamp int `json:"build_timestamp"`
}

// FindHostVars parses hostvars.json file.
// It is expected to be in /etc.
func FindHostVars() (HostVars, error) {
	var vars HostVars
	file := filepath.Join("etc/", HostVarsName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return vars, fmt.Errorf("%s not found: %v", file, err)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return vars, fmt.Errorf("cant open %s: %v", file, err)
	}
	if err = json.Unmarshal(data, &vars); err != nil {
		return vars, fmt.Errorf("cant parse data from %s", file)
	}
	return vars, nil
}
