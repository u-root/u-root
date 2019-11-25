// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	sysfsPath = "/sys/firmware/dmi/tables"
)

func getSMBIOSData() ([]byte, []byte, error) {
	entry, err := ioutil.ReadFile(filepath.Join(sysfsPath, "smbios_entry_point"))
	if err != nil {
		return nil, nil, fmt.Errorf("error reading DMI data: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(sysfsPath, "DMI"))
	if err != nil {
		return nil, nil, fmt.Errorf("error reading DMI data: %v", err)
	}
	return entry, data, nil
}
