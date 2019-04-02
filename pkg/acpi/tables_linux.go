// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package acpi

import (
	"path/filepath"
)

// RawTables returns an array of Raw, for all ACPI tables
// available in /sys
func RawTables() ([]Tabler, error) {
	n, err := filepath.Glob("/sys/firmware/acpi/tables/[A-Z]*")
	if err != nil {
		return nil, err
	}

	var tabs []Tabler
	for _, t := range n {
		r, err := RawFromFile(t)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, r)
	}
	return tabs, nil
}
