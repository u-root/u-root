package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptInterfaceID(t *testing.T) {
	expected := []byte("DSLAM01 eth2/1/01/21")
	opt, err := parseOptInterfaceID(expected)
	if err != nil {
		t.Fatal(err)
	}
	if url := opt.ID; !bytes.Equal(url, expected) {
		t.Fatalf("Invalid Interface ID. Expected %v, got %v", expected, url)
	}
}

func TestOptInterfaceID(t *testing.T) {
	want := []byte("DSLAM01 eth2/1/01/21")
	opt := OptInterfaceID(want)
	if got := opt.ToBytes(); !bytes.Equal(got, want) {
		t.Fatalf("%s.ToBytes() = %v, want %v", opt, got, want)
	}

	require.Contains(
		t,
		opt.String(),
		"68 83 76 65 77 48 49 32 101 116 104 50 47 49 47 48 49 47 50 49",
		"String() should return the interfaceId as bytes",
	)
}
