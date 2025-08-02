// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msr

// MSR defines an MSR address.
type MSR uint32

// MSRVal defines an MSR Value to be used in TestAndSet, Test, and other
// operations.
type MSRVal struct {
	// Addr is the address of the MSR.
	Addr MSR
	// Name is a printable name.
	Name string
	// Clear are bits to clear in TestAndSet and Test.
	Clear uint64
	// Set are bits to set in TestAndSet and Test; or the value to write
	// if the MSR is write-only.
	Set uint64
	// WriteOnly indicates an MSR is write-only.
	WriteOnly bool
}

// String implements String() for MSRVal
func (m MSRVal) String() string {
	return m.Name
}

// Debug can be set for debug prints on MSR operations.
// It can be set to, e.g., log.Printf.
// It's default action is to do nothing.
var Debug = func(string, ...any) {}
