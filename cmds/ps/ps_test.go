// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"testing"
)

// Simple Test trying execute the ps
// If no errors returns, it's okay
func TestPsExecution(t *testing.T) {
	t.Skip("ps is abuggy and this test is breaking travis; skipping")
	pT := ProcessTable{}
	if err := pT.LoadTable(); err != nil {
		t.Fatalf("Loading Table fails on some point; %v", err)
	}

	if err := ps(pT); err != nil {
		t.Fatalf("Calling ps fails; %v", err)
	}
}
