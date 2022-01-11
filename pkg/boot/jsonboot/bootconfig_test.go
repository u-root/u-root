// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonboot

import (
	"testing"
)

func TestNewBootConfig(t *testing.T) {
	data := []byte(`{
	"name": "some_conf",
	"kernel": "/path/to/kernel",
	"initramfs": "/path/to/initramfs",
	"kernel_args": "init=/bin/bash",
	"devicetree": "some data here"
}`)
	c, err := NewBootConfig(data)
	if c.Name != "some_conf" || err != nil {
		t.Errorf(`NewBootConfig(data).Name = %q, %v, want "some_conf", nil`, c.Name, err)
	}
	if c.Kernel != "/path/to/kernel" {
		t.Errorf(`NewBootConfig(data).Kernel = %q, %v, want "/path/to/kernel", nil`, c.Kernel, err)
	}
	if c.Initramfs != "/path/to/initramfs" {
		t.Errorf(`NewBootConfig(data).Initramfs = %q, %v, want "/path/to/initramfs", nil`, c.Initramfs, err)
	}
	if c.KernelArgs != "init=/bin/bash" {
		t.Errorf(`NewBootConfig(data).KernelArgs = %q, %v, want "init=/bin/bash", nil`, c.KernelArgs, err)
	}
	if c.DeviceTree != "some data here" {
		t.Errorf(`NewBootConfig(data).DeviceTree = %q, %v, want "some data here", nil`, c.DeviceTree, err)
	}
	if !c.IsValid() {
		t.Errorf(`NewBootConfig(data).IsValid() = %t, %v, want "true", nil`, c.IsValid(), err)
	}
}

func TestNewBootConfigInvalidJSON(t *testing.T) {
	data := []byte(`{
	"name": "broken
}`)
	c, err := NewBootConfig(data)
	if err == nil {
		t.Errorf(`NewBootConfig(data) = %q, %v, want "nil", error`, c, err)
	}
}

func TestNewBootConfigMissingKernel(t *testing.T) {
	data := []byte(`{
	"name": "some_conf",
	"kernel_is_missing": "/path/to/kernel",
	"initramfs": "/path/to/initramfs",
	"kernel_args": "init=/bin/bash",
	"devicetree": "some data here"
}`)
	c, err := NewBootConfig(data)
	if c.IsValid() != false || err != nil {
		t.Errorf(`NewBootConfig(data).IsValid() = %t, %v, want "false", nil`, c.IsValid(), err)
	}
}

func TestID(t *testing.T) {
	bc := BootConfig{
		Name: "Slash and space should not \\ appear /here",
	}
	id := bc.ID()
	t.Log(id)
}
