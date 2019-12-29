package interfaces

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func fakeIface(idx int, name string, loopback bool) net.Interface {
	var flags net.Flags
	if loopback {
		flags |= net.FlagLoopback
	}
	return net.Interface{
		Index:        idx,
		MTU:          1500,
		Name:         name,
		HardwareAddr: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Flags:        flags,
	}
}

func TestGetLoopbackInterfaces(t *testing.T) {
	interfaceGetter = func() ([]net.Interface, error) {
		return []net.Interface{
			fakeIface(0, "lo", true),
			fakeIface(1, "eth0", false),
			fakeIface(2, "eth1", false),
		}, nil
	}
	ifaces, err := GetLoopbackInterfaces()
	// this has to be reassigned before any require.* call
	interfaceGetter = net.Interfaces

	require.NoError(t, err)
	require.Equal(t, 1, len(ifaces))
}

func TestGetLoopbackInterfacesError(t *testing.T) {
	interfaceGetter = func() ([]net.Interface, error) {
		return nil, errors.New("expected error")

	}
	_, err := GetLoopbackInterfaces()
	// this has to be reassigned before any require.* call
	interfaceGetter = net.Interfaces

	require.Error(t, err)
}

func TestGetNonLoopbackInterfaces(t *testing.T) {
	interfaceGetter = func() ([]net.Interface, error) {
		return []net.Interface{
			fakeIface(0, "lo", true),
			fakeIface(1, "eth0", false),
			fakeIface(2, "eth1", false),
		}, nil
	}
	ifaces, err := GetNonLoopbackInterfaces()
	// this has to be reassigned before any require.* call
	interfaceGetter = net.Interfaces

	require.NoError(t, err)
	require.Equal(t, 2, len(ifaces))
}

func TestGetNonLoopbackInterfacesError(t *testing.T) {
	interfaceGetter = func() ([]net.Interface, error) {
		return nil, errors.New("expected error")

	}
	_, err := GetNonLoopbackInterfaces()
	// this has to be reassigned before any require.* call
	interfaceGetter = net.Interfaces

	require.Error(t, err)
}
