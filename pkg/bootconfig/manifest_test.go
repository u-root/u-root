package bootconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManifest(t *testing.T) {
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
	m, err := NewManifest(data)
	require.NoError(t, err)
	require.Equal(t, 1, len(m.Configs))
}

func TestNewManifestInvalid(t *testing.T) {
	data := []byte(`{
		"nonexisting": "baaah",
		"configs": {
			"broken": true
		}
}`)
	_, err := NewManifest(data)
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
	m, err := NewManifest(data)
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
	m, err := NewManifest(data)
	require.NoError(t, err)
	_, err = m.GetBootConfig(1)
	require.Error(t, err)
}
