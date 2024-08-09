// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package acpi

import "fmt"

// GetTable uses all the Methods until one works.
func GetTable() (string, []Table, error) {
	for m, f := range Methods {
		t, err := f()
		if err == nil {
			return m, t, nil
		}
	}
	return "", nil, fmt.Errorf("could not get ACPI tables")
}
