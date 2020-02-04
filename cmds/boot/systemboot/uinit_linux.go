// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/u-root/u-root/pkg/smbios"
)

const (
	sysfsPath = "/sys/firmware/dmi/tables"
)

func getSMBIOSInfo() (*smbios.Info, error) {
	entry, err := ioutil.ReadFile(filepath.Join(sysfsPath, "smbios_entry_point"))
	if err != nil {
		return nil, fmt.Errorf("error reading DMI data: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(sysfsPath, "DMI"))
	if err != nil {
		return nil, fmt.Errorf("error reading DMI data: %v", err)
	}
	si, err := smbios.ParseInfo(entry, data)
	if err != nil {
		return nil, err
	}
	return si, nil
}
