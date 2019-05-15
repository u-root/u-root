// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bytes"
	"fmt"
	"os/exec"
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

func compress(n string) error {
	c := exec.Command("xz", "-f", n)
	o, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v", string(o), err)
	}
	return nil
}

// Not sure how to do what needs doing here.
func TestXZModules(t *testing.T) {
	t.Skip("TODO: figure out how to make this work")
	// Compress the modules using the basic command, i.e. xz
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
	}
	for k := range m {
		if err := compress(k); err != nil {
			t.Errorf("Can't compress %v: %v", k, err)
		}
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
