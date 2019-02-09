// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bytes"
	"path"
	"testing"
)

var procModsMock = `hid_generic 16384 0 - Live 0x0000000000000000
usbhid 49152 0 - Live 0x0000000000000000
ccm 20480 6 - Live 0x0000000000000000
`

func TestGenLoadedMods(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
	}
	br := bytes.NewBufferString(procModsMock)
	err := genLoadedMods(br, m)
	if err != nil {
		t.Fatalf("fail to genLoadedMods: %v\n", err)
	}
	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}
