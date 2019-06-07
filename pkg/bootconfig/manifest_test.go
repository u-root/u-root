// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManifest(t *testing.T) {
	m := NewManifest()
	require.NotNil(t, m)
	require.Equal(t, m.Version, CurrentManifestVersion)
}

func TestManifestFromBytes(t *testing.T) {
	data := []byte(`{
	"version": 1,
	"configs": [
		{
			"name": "some_boot_config",
			"kernel": "/path/to/kernel",
			"initramfs": "/path/to/initramfs",
			"kernel_args": "init=/bin/bash",
			"devicetree": "some data here"
		}
	]
}`)
	m, err := ManifestFromBytes(data)
	require.NoError(t, err)
	require.Equal(t, 1, len(m.Configs))
}

func TestManifestFromBytesInvalid(t *testing.T) {
	data := []byte(`{
		"nonexisting": "baaah",
		"configs": {
			"broken": true
		}
}`)
	_, err := ManifestFromBytes(data)
	require.Error(t, err)
}

func TestManifestGetBootConfig(t *testing.T) {
	data := []byte(`{
	"version": 1,
	"configs": [
		{
			"name": "some_boot_config",
			"kernel": "/path/to/kernel"
		}
	]
}`)
	m, err := ManifestFromBytes(data)
	require.NoError(t, err)
	config, err := m.GetBootConfig(0)
	require.NoError(t, err)
	assert.Equal(t, "some_boot_config", config.Name)
	assert.Equal(t, "/path/to/kernel", config.Kernel)
}

func TestManifestGetBootConfigMissing(t *testing.T) {
	data := []byte(`{
	"version": 1,
	"configs": [
		{
			"name": "some_boot_config",
			"kernel": "/path/to/kernel"
		}
	]
}`)
	m, err := ManifestFromBytes(data)
	require.NoError(t, err)
	_, err = m.GetBootConfig(1)
	require.Error(t, err)
}
