package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadHostvars(t *testing.T) {
	h, err := loadHostvars("testdata/hostvars.json")
	require.NoError(t, err)
	require.Equal(t, LocalStorage, h.BootMode)
}

func TestLoadHostvarsInvalid(t *testing.T) {
	_, err := loadHostvars("testdata/hostvars_invalid.json")
	require.Error(t, err)
}