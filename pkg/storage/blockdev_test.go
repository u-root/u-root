package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFindMountPointNotExists checks that non existent
// entry is checked and nil returned
func TestFindMountPointNotExists(t *testing.T) {
	LinuxMountsPath = "tests/mounts"
	_, err := GetMountpointByDevice("/dev/mapper/sys-oldxxxxxx")
	require.Error(t, err)
}

// TestFindMountPointValid check for valid output of
// test mountpoint.
func TestFindMountPointValid(t *testing.T) {
	LinuxMountsPath = "tests/mounts"
	mountpoint, err := GetMountpointByDevice("/dev/mapper/sys-old")
	require.NoError(t, err)
	require.Equal(t, *mountpoint, "/media/usb")
}
