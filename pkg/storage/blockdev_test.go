// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFindMountPointNotExists checks that non existent
// entry is checked and nil returned
func TestFindMountPointNotExists(t *testing.T) {
	LinuxMountsPath = "testdata/mounts"
	_, err := GetMountpointByDevice("/dev/mapper/sys-oldxxxxxx")
	require.Error(t, err)
}

// TestFindMountPointValid check for valid output of
// test mountpoint.
func TestFindMountPointValid(t *testing.T) {
	LinuxMountsPath = "testdata/mounts"
	mountpoint, err := GetMountpointByDevice("/dev/mapper/sys-old")
	require.NoError(t, err)
	require.Equal(t, *mountpoint, "/media/usb")
}

func TestBlockStatFromBytes15Fields(t *testing.T) {
	// dummy values, don't judge me
	input := []byte("       0        1        2        3        4        5        6        7        8        9        10        11        12        13        14\n")
	bs, err := BlockStatFromBytes(input)
	require.NoError(t, err)
	require.Equal(t, uint64(5), bs.WriteMerges)
	require.Equal(t, uint64(14), bs.DiscardTicks)
}

func TestBlockStatFromBytes11Fields(t *testing.T) {
	// dummy values, don't judge me
	input := []byte("       0        1        2        3        4        5        6        7        8        9        10\n")
	bs, err := BlockStatFromBytes(input)
	require.NoError(t, err)
	require.Equal(t, uint64(5), bs.WriteMerges)
	require.Equal(t, uint64(0), bs.DiscardTicks)
}

func TestBlockStatFromBytesNotEnoughFields(t *testing.T) {
	// dummy values, don't judge me
	input := []byte("       0        1        2        3        4        5        6        7        8\n")
	_, err := BlockStatFromBytes(input)
	require.Error(t, err)
}
