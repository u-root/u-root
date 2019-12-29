package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptRemoteId(t *testing.T) {
	expected := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected = append(expected, remoteId...)
	opt, err := ParseOptRemoteId(expected)
	if err != nil {
		t.Fatal(err)
	}
	if en := opt.EnterpriseNumber(); en != 0xaabbccdd {
		t.Fatalf("Invalid Enterprise Number. Expected 0xaabbccdd, got %v", en)
	}
	if rid := opt.RemoteID(); !bytes.Equal(rid, remoteId) {
		t.Fatalf("Invalid Remote ID. Expected %v, got %v", expected, rid)
	}
}

func TestOptRemoteIdToBytes(t *testing.T) {
	remoteId := []byte("DSLAM01 eth2/1/01/21")
	expected := append([]byte{0, 0, 0, 0}, remoteId...)
	opt := OptRemoteId{
		remoteId: remoteId,
	}
	toBytes := opt.ToBytes()
	if !bytes.Equal(toBytes, expected) {
		t.Fatalf("Invalid ToBytes result. Expected %v, got %v", expected, toBytes)
	}
}

func TestOptRemoteIdSet(t *testing.T) {
	enterpriseNumber := uint32(12345)
	remoteID := []byte("DSLAM01 eth2/1/01/21")
	opt := OptRemoteId{}
	opt.SetEnterpriseNumber(enterpriseNumber)
	opt.SetRemoteID(remoteID)

	require.Equal(t, uint32(12345), opt.EnterpriseNumber())
	require.Equal(t, []byte("DSLAM01 eth2/1/01/21"), opt.RemoteID())
}

func TestOptRemoteIdParseOptRemoteIdTooShort(t *testing.T) {
	buf := []byte{0xaa, 0xbb, 0xcc}
	_, err := ParseOptRemoteId(buf)
	require.Error(t, err, "A short option should return an error")
}

func TestOptRemoteIdString(t *testing.T) {
	buf := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	remoteId := []byte("Test1234")
	buf = append(buf, remoteId...)

	opt, err := ParseOptRemoteId(buf)
	require.NoError(t, err)
	str := opt.String()
	require.Contains(
		t,
		str,
		"enterprisenum=2864434397",
		"String() should contain the enterprisenum",
	)
	require.Contains(
		t,
		str,
		"remoteid=[84 101 115 116 49 50 51 52]",
		"String() should contain the remoteid bytes",
	)
}
