package ztpv4

import (
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/require"
)

func TestParseV4VendorClass(t *testing.T) {
	tt := []struct {
		name         string
		vc, hostname string
		want         *VendorData
		fail         bool
	}{
		{name: "empty", fail: true},
		{name: "unknownVendor", vc: "VendorX;BFR10K;XX12345", fail: true},
		{name: "truncatedVendor", vc: "Arista;1234", fail: true},
		{
			name: "arista",
			vc:   "Arista;DCS-7050S-64;01.23;JPE12345678",
			want: &VendorData{VendorName: "Arista", Model: "DCS-7050S-64", Serial: "JPE12345678"},
		},
		{
			name: "juniper",
			vc:   "Juniper-ptx1000-DD123",
			want: &VendorData{VendorName: "Juniper", Model: "ptx1000", Serial: "DD123"},
		},
		{
			name: "juniperModelDash",
			vc:   "Juniper-qfx10002-36q-DN817",
			want: &VendorData{VendorName: "Juniper", Model: "qfx10002-36q", Serial: "DN817"},
		},
		{
			name:     "juniperHostnameSerial",
			vc:       "Juniper-qfx10008",
			hostname: "DE123",
			want:     &VendorData{VendorName: "Juniper", Model: "qfx10008", Serial: "DE123"},
		},
		{name: "juniperNoSerial", vc: "Juniper-qfx10008", fail: true},
		{
			name: "zpe",
			vc:   "ZPESystems:NSC:001234567",
			want: &VendorData{VendorName: "ZPESystems", Model: "NSC", Serial: "001234567"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := dhcpv4.New()
			if err != nil {
				t.Fatalf("failed to creat dhcpv4 packet object: %v", err)
			}

			if tc.vc != "" {
				packet.UpdateOption(dhcpv4.OptClassIdentifier(tc.vc))
			}
			if tc.hostname != "" {
				packet.UpdateOption(dhcpv4.OptHostName(tc.hostname))
			}

			vd, err := ParseVendorData(packet)
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, vd)
			}
		})
	}
}
