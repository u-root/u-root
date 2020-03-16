// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// FromSysfs parses SMBIOS info from sysfs tables.
func FromSysfs() (*Info, error) {
	return fromSysfs("/sys/firmware/dmi/tables")
}

func fromSysfs(sysfsPath string) (*Info, error) {
	entry, err := ioutil.ReadFile(filepath.Join(sysfsPath, "smbios_entry_point"))
	if err != nil {
		return nil, fmt.Errorf("error reading SMBIOS entry data: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(sysfsPath, "DMI"))
	if err != nil {
		return nil, fmt.Errorf("error reading DMI data: %v", err)
	}
	return ParseInfo(entry, data)
}
