// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	Debug = t.Logf
	d := DevName + "0"
	if _, err := os.Stat(d); err != nil {
		t.Skip("No device to test")
	}
	m, err := NewChipInfoFromDev(d)
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Errorf("no ChipInfo found in sysfs")
	}
	t.Logf("Chip info: name %v string %v", m.Name(), m.String())
}
