// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package cpuid

import (
	intelcpuid "github.com/intel-go/cpuid"
)

const (
	ManufacturerIDAMD   = "AuthenticAMD"
	ManufacturerIDIntel = "GenuineIntel"
)

// Get the CPU Identification String and return it.
func CPUManufacturerID() (string, error) {
	return intelcpuid.VendorIdentificatorString, nil
}
