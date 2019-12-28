// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tss provides TPM 1.2/2.0 core functionality and
// abstraction layer for high-level functions
package tss

import (
	"testing"
)

func TestReadPCRs(t *testing.T) {
	tss, err := OpenTPM()
	if err != nil {
		t.Error()
	}
	_, err = tss.ReadPCRs(HashSHA1)
	if err != nil {
		t.Error()
	}
}
