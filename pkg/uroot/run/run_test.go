// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run

import (
	"os"
	"testing"
)

func TestDefaultParams(t *testing.T) {
	p := DefaultParams()
	if p.Wd == "" {
		t.Error("got empty working directory")
	}
	_, err := os.Getwd()
	if err != nil {
		t.Errorf("error from os.Getwd: %v", err)
	}
}
