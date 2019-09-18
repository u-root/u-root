// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
)

// getData returns SMBIOS entry point and DMI table data.
// If dumpFile is non-empty, it is read from that file, otherwise it is read
// from sysfsPath (smbios_entry_point and DMI files respectively).
func getData(textOut io.Writer, dumpFile, sysfsPath string) ([]byte, []byte, error) {
	var err error
	var entry, data []byte
	if dumpFile != "" {
		fmt.Fprintf(textOut, "Reading SMBIOS/DMI data from file %s.\n", dumpFile)
		data, err = ioutil.ReadFile(dumpFile)
		if err != nil {
			return nil, nil, fmt.Errorf("error reading dump: %v", err)
		}
		if len(data) < 36 {
			return nil, nil, errors.New("dump is too short")
		}
		entry = data[:32]
		data = data[32:]
	} else {
		fmt.Fprintf(textOut, "Reading SMBIOS/DMI data from sysfs.\n")
		entry, err = ioutil.ReadFile(filepath.Join(sysfsPath, "smbios_entry_point"))
		if err != nil {
			return nil, nil, fmt.Errorf("error reading DMI data: %v", err)
		}
		data, err = ioutil.ReadFile(filepath.Join(sysfsPath, "DMI"))
		if err != nil {
			return nil, nil, fmt.Errorf("error reading DMI data: %v", err)
		}
	}
	return entry, data, nil
}
