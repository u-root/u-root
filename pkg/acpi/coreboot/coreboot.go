// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package coreboot transforms ACPI tables in various ways, for use in coreboot.
// The most common transformation is code generation for use in mainboards
// or chipsets.
package coreboot

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/acpi"
)

type newcb func(acpi.Table) (Corebooter, error)

// Corebooter writes C code for coreboot.
// N.B. proper coreboot spelling is all lower case; we accept
// Go rules here and call it Corebooter, but not CoreBooter!
type Corebooter interface {
	Coreboot(w io.Writer) error
}

var (
	// corebooters is the map of names to a corebooter.
	corebooters = map[string]newcb{}

	// Debug enables various debug prints. External code can set it to, e.g., log.Printf
	Debug = func(string, ...interface{}) {}
)

// NewCorebooters returns a Corebooter for a Table or an error
func NewCorebooter(t acpi.Table) (Corebooter, error) {
	cb, ok := corebooters[t.Sig()]
	if !ok {
		return nil, fmt.Errorf("No corebooter for %q", t.Sig())
	}
	return cb(t)
}
