package dhcpv4

import (
	"net"
	"reflect"
	"testing"
)

func mustParseIPNet(s string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ipnet
}

func TestParseRoutes(t *testing.T) {
	for _, tt := range []struct {
		p       []byte
		want    Routes
		wantErr bool
	}{
		{
			p: []byte{32, 10, 2, 3, 4, 0, 0, 0, 0},
			want: Routes{
				&Route{
					Dest:   mustParseIPNet("10.2.3.4/32"),
					Router: net.IP{0, 0, 0, 0},
				},
			},
		},
		{
			p: []byte{0, 0, 0, 0, 0},
			want: Routes{
				&Route{
					Dest:   mustParseIPNet("0.0.0.0/0"),
					Router: net.IP{0, 0, 0, 0},
				},
			},
		},
		{
			p: []byte{32, 10, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want: Routes{
				&Route{
					Dest:   mustParseIPNet("10.2.3.4/32"),
					Router: net.IP{0, 0, 0, 0},
				},
				&Route{
					Dest:   mustParseIPNet("0.0.0.0/0"),
					Router: net.IP{0, 0, 0, 0},
				},
			},
		},
		{
			p:       []byte{64, 10, 2, 3, 4},
			wantErr: true, // Mask length 64 > 32
		},
	} {
		var r Routes
		if err := r.FromBytes(tt.p); (err != nil) != tt.wantErr {
			t.Errorf("FromBytes(%v) Unexpected error state: %v", tt.p, err)
		}

		if !reflect.DeepEqual(r, tt.want) {
			t.Errorf("FromBytes(%v) = %v, want %v", tt.p, r, tt.want)
		}
	}
}
