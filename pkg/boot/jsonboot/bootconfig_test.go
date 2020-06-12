// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonboot

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.Equal(t, "some_conf", c.Name)
	require.Equal(t, "/path/to/kernel", c.Kernel)
	require.Equal(t, "/path/to/initramfs", c.Initramfs)
	require.Equal(t, "init=/bin/bash", c.KernelArgs)
	require.Equal(t, "some data here", c.DeviceTree)
	require.Equal(t, true, c.IsValid())
}

func TestNewBootConfigInvalidJSON(t *testing.T) {
	data := []byte(`{
	"name": "broken
}`)
	_, err := NewBootConfig(data)
	require.Error(t, err)
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
	require.NoError(t, err)
	require.Equal(t, false, c.IsValid())
}

func TestID(t *testing.T) {
	bc := BootConfig{
		Name: "Slash and space should not \\ appear /here",
	}
	id := bc.ID()
	t.Log(id)
}
