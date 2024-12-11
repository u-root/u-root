// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"os"
	"path/filepath"
)

// FromSysfs parses SMBIOS info from sysfs tables.
func FromSysfs() (*Info, error) {
	return fromSysfs("/sys/firmware/dmi/tables")
}

func fromSysfs(sysfsPath string) (*Info, error) {
	entry, err := os.ReadFile(filepath.Join(sysfsPath, "smbios_entry_point"))
	if err != nil {
		return nil, fmt.Errorf("error reading SMBIOS entry data: %w", err)
	}
	data, err := os.ReadFile(filepath.Join(sysfsPath, "DMI"))
	if err != nil {
		return nil, fmt.Errorf("error reading DMI data: %w", err)
	}
	return ParseInfo(entry, data)
}
