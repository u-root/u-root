package dhcpv6

import (
	"testing"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/stretchr/testify/require"
)

func TestParseOptDomainSearchList(t *testing.T) {
	data := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt, err := parseOptDomainSearchList(data)
	require.NoError(t, err)
	require.Equal(t, OptionDomainSearchList, opt.Code())
	require.Equal(t, 2, len(opt.DomainSearchList.Labels))
	require.Equal(t, "example.com", opt.DomainSearchList.Labels[0])
	require.Equal(t, "subnet.example.org", opt.DomainSearchList.Labels[1])
	require.Contains(t, opt.String(), "example.com subnet.example.org", "String() should contain the correct domain search output")
}

func TestOptDomainSearchListToBytes(t *testing.T) {
	expected := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'o', 'r', 'g', 0,
	}
	opt := OptDomainSearchList(
		&rfc1035label.Labels{
			Labels: []string{
				"example.com",
				"subnet.example.org",
			},
		},
	)
	require.Equal(t, expected, opt.ToBytes())
}

func TestParseOptDomainSearchListInvalidLength(t *testing.T) {
	data := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0,
		6, 's', 'u', 'b', 'n', 'e', 't', 7, 'e', // truncated
	}
	_, err := parseOptDomainSearchList(data)
	require.Error(t, err, "A truncated OptDomainSearchList should return an error")
}
