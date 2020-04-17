// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

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
	require.NotContains(t, bc.ID(), "/", "\\")
}

func TestFiles(t *testing.T) {
	bc := BootConfig{
		Kernel:    "path/to/kernel",
		Initramfs: "path/to/initramfs",
		Multiboot: "path/to/multibootkernel",
		Modules: []string{
			"path/to/mod1",
			"path/to/mod2 -arg1 -arg2 -arg3",
			"path/to/mod3 arg1",
		},
	}
	names := bc.Files()
	require.Equal(t, []string{"path/to/kernel", "path/to/initramfs", "path/to/multibootkernel", "path/to/mod1", "path/to/mod2", "path/to/mod3"}, names)
}

func TestChangeFilePaths(t *testing.T) {
	bc := BootConfig{
		Kernel:    "path/to/kernel",
		Initramfs: "path/to/initramfs",
		Multiboot: "path/to/multibootkernel",
		Modules: []string{
			"path/to/mod1",
			"path/to/mod2 -arg1 -arg2 -arg3",
			"path/to/mod3 arg1",
		},
	}
	bc.ChangeFilePaths("new/path/with/same/base/")
	require.Equal(t, "new/path/with/same/base/kernel", bc.Kernel)
	require.Equal(t, "new/path/with/same/base/initramfs", bc.Initramfs)
	require.Equal(t, "new/path/with/same/base/multibootkernel", bc.Multiboot)
	require.Equal(t, "new/path/with/same/base/mod1", bc.Modules[0])
	require.Equal(t, "new/path/with/same/base/mod2 -arg1 -arg2 -arg3", bc.Modules[1])
	require.Equal(t, "new/path/with/same/base/mod3 arg1", bc.Modules[2])
}

func TestSetFilePathsPrefix(t *testing.T) {
	bc := BootConfig{
		Kernel:    "path/to/kernel",
		Initramfs: "path/to/initramfs",
		Multiboot: "path/to/multibootkernel",
		Modules: []string{
			"path/to/mod1",
			"path/to/mod2 -arg1 -arg2 -arg3",
			"path/to/mod3 arg1",
		},
	}
	bc.SetFilePathsPrefix("prefix")
	require.Equal(t, "prefix/path/to/kernel", bc.Kernel)
	require.Equal(t, "prefix/path/to/initramfs", bc.Initramfs)
	require.Equal(t, "prefix/path/to/multibootkernel", bc.Multiboot)
	require.Equal(t, "prefix/path/to/mod1", bc.Modules[0])
	require.Equal(t, "prefix/path/to/mod2 -arg1 -arg2 -arg3", bc.Modules[1])
	require.Equal(t, "prefix/path/to/mod3 arg1", bc.Modules[2])
}
