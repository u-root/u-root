package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptInterfaceId(t *testing.T) {
	expected := []byte("DSLAM01 eth2/1/01/21")
	opt, err := ParseOptInterfaceId(expected)
	if err != nil {
		t.Fatal(err)
	}
	if url := opt.InterfaceID(); !bytes.Equal(url, expected) {
		t.Fatalf("Invalid Interface ID. Expected %v, got %v", expected, url)
	}
}

func TestOptInterfaceIdToBytes(t *testing.T) {
	interfaceId := []byte("DSLAM01 eth2/1/01/21")
	opt := OptInterfaceId{}
	opt.SetInterfaceID(interfaceId)
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, interfaceId) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", interfaceId, toBytes)
	}
}

func TestOptInterfaceIdString(t *testing.T) {
	interfaceId := []byte("DSLAM01 eth2/1/01/21")
	opt := OptInterfaceId{}
	opt.SetInterfaceID(interfaceId)
	require.Contains(
		t,
		opt.String(),
		"68 83 76 65 77 48 49 32 101 116 104 50 47 49 47 48 49 47 50 49",
		"String() should return the interfaceId as bytes",
	)
}
