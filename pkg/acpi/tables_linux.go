// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"io/ioutil"
	"path/filepath"
)

// Table contains a []byte for a single ACPI table.
// The Name is provided for convenience.
type Table struct {
	Name string
	Data []byte
}

// Tables returns an array of Table
func Tables() ([]*Table, error) {
	n, err := filepath.Glob("/sys/firmware/acpi/tables/[A-Z]*")
	if err != nil {
		return nil, err
	}

	var tabs []*Table
	for _, t := range n {
		b, err := ioutil.ReadFile(t)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, &Table{Name: t, Data: b})
	}
	return tabs, nil
}

// TablesData returns all ACPI tables as a single []byte
func TablesData() ([]byte, error) {
	t, err := Tables()
	if err != nil {
		return nil, err
	}
	var b []byte

	for _, tab := range t {
		b = append(b, tab.Data...)
	}
	return b, nil
}
