// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"io/ioutil"
	"path/filepath"
)

// The RawTable is used until such time as we've implemented Table
// or a standin for everything.

// Table contains a []byte for a single ACPI table.
// The Name is provided for convenience.
type RawTable struct {
	Name string
	Data []byte
}

// RawTables returns an array of RawTable
func RawTables() ([]*RawTable, error) {
	n, err := filepath.Glob("/sys/firmware/acpi/tables/[A-Z]*")
	if err != nil {
		return nil, err
	}

	var tabs []*RawTable
	for _, t := range n {
		b, err := ioutil.ReadFile(t)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, &RawTable{Name: t, Data: b})
	}
	return tabs, nil
}

// TablesData returns all ACPI tables as a single []byte
func RawTablesData() ([]byte, error) {
	t, err := RawTables()
	if err != nil {
		return nil, err
	}
	var b []byte

	for _, tab := range t {
		b = append(b, tab.Data...)
	}
	return b, nil
}
