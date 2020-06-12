// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package systembooter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBooterForNetBooter(t *testing.T) {
	validConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`),
	}
	booter := GetBooterFor(validConfig)
	require.NotNil(t, booter)
	require.Equal(t, booter.TypeName(), "netboot")
	require.NotNil(t, booter.(*NetBooter))
}

func TestGetBooterForNullBooter(t *testing.T) {
	validConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "null"}`),
	}
	booter := GetBooterFor(validConfig)
	require.NotNil(t, booter)
	require.Equal(t, booter.TypeName(), "null")
	require.NotNil(t, booter.(*NullBooter))
	require.Nil(t, booter.Boot())
}

func TestGetBooterForInvalidBooter(t *testing.T) {
	invalidConfig := BootEntry{
		Name:   "Boot0000",
		Config: []byte(`{"type": "invalid"`),
	}
	booter := GetBooterFor(invalidConfig)
	require.NotNil(t, booter)
	// an invalid config returns always a NullBooter
	require.Equal(t, booter.TypeName(), "null")
	require.NotNil(t, booter.(*NullBooter))
	require.Nil(t, booter.Boot())
}

func TestGetBootEntries(t *testing.T) {
	var (
		bootConfig0000 = []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`)
		bootConfig0001 = []byte(`{"type": "localboot", "uuid": "blah-bleh", "kernel": "/path/to/kernel"}`)
	)
	// Override the package-level variable Get so it will use our dummy getter
	// instead of VPD
	Get = func(key string, readOnly bool) ([]byte, error) {
		switch key {
		case "Boot0000":
			return bootConfig0000, nil
		case "Boot0001":
			return bootConfig0001, nil
		default:
			return nil, errors.New("No such key")
		}
	}
	entries := GetBootEntries()
	require.Equal(t, len(entries), 2)
	require.Equal(t, "Boot0000", entries[0].Name)
	require.Equal(t, bootConfig0000, entries[0].Config)
	require.Equal(t, "Boot0001", entries[1].Name)
	require.Equal(t, bootConfig0001, entries[1].Config)
}

func TestGetBootEntriesOnlyRO(t *testing.T) {
	// Override the package-level variable Get so it will use our dummy getter
	// instead of VPD
	Get = func(key string, readOnly bool) ([]byte, error) {
		if !readOnly || key != "Boot0000" {
			return nil, errors.New("No such key")
		}
		return []byte(`{"type": "netboot", "method": "dhcpv6", "mac": "aa:bb:cc:dd:ee:ff"}`), nil
	}
	entries := GetBootEntries()
	require.Equal(t, len(entries), 1)
}
