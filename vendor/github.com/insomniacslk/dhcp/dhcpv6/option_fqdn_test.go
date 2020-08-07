package dhcpv6

import (
	"bytes"
	"testing"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/require"
)

func TestParseOptFQDN(t *testing.T) {
	data := []byte{
		0, // Flags
		4, 'c', 'n', 'o', 's', 9, 'l', 'o', 'c', 'a', 'l',
		'h', 'o', 's', 't', 0,
	}
	opt, err := ParseOptFQDN(data)

	require.NoError(t, err)
	require.Equal(t, OptionFQDN, opt.Code())
	require.Equal(t, uint8(0), opt.Flags)
	require.Equal(t, "cnos.localhost", opt.DomainName.Labels[0])
	require.Equal(t, "OptFQDN{flags=0, domainname=[cnos.localhost]}", opt.String())
}

func TestOptFQDNToBytes(t *testing.T) {
	opt := OptFQDN{
		Flags:      0,
		DomainName: &rfc1035label.Labels{
			Labels: []string{"cnos.localhost"},
		},
	}
	want := []byte{
		0, // Flags
		4, 'c', 'n', 'o', 's', 9, 'l', 'o', 'c', 'a', 'l',
		'h', 'o', 's', 't', 0,
	}
	b := opt.ToBytes()
	if !bytes.Equal(b, want) {
		t.Fatalf("opt.ToBytes()=%v, want %v", b, want)
	}
}
