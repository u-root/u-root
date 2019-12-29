package rtnl

import (
	"net"
	"testing"
)

func TestParseAddrs(t *testing.T) {

	tests := []struct {
		name  string
		ipstr string
		ipnet net.IPNet
		err   string
	}{
		{
			name:  "ipv6 subnet address",
			ipstr: "ff00::/64",
			err:   "address ff00::: attempted to parse a subnet address into a host address",
		},
		{
			name:  "ipv4 subnet address",
			ipstr: "10.0.0.0/8",
			err:   "address 10.0.0.0: attempted to parse a subnet address into a host address",
		},
		{
			name:  "ipv6 host address",
			ipstr: "ff00::1/64",
			ipnet: net.IPNet{
				IP:   net.IP{0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
				Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
		},
		{
			name:  "ipv4 host address",
			ipstr: "10.0.0.1/8",
			ipnet: net.IPNet{
				IP:   net.IP{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xa, 0x0, 0x0, 0x1},
				Mask: net.IPMask{0xff, 0x0, 0x0, 0x0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseAddr(tt.ipstr)

			if err != nil {
				if want, got := tt.err, err.Error(); want != got {
					t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
				}
				return
			}
			if tt.err != "" {
				t.Fatalf("expected error:\n  %s\nbut got nothing.. :(", tt.err)
			}

			if want, got := tt.ipnet, res; !want.IP.Equal(got.IP) {
				t.Fatalf("unexpected IP:\n- want: %+#v\n-  got: %+#v", want, got)
			}
		})
	}
}
