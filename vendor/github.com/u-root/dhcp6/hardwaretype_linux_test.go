// +build linux

package dhcp6

import (
	"net"
	"strings"
	"testing"
)

// TestHardwareTypeLinux verifies that HardwareType detects the expected
// hardware type for a given interface on a Linux-based machine.
func TestHardwareTypeLinux(t *testing.T) {
	// Use eth0, since it is most likely the most common Linux interface
	// with a standard hardware type number
	eth0, err := net.InterfaceByName("eth0")
	if err != nil {
		if strings.Contains(err.Error(), "no such network interface") {
			t.Skip("skipping, system does not have interface eth0")
		}

		t.Fatal(err)
	}

	var tests = []struct {
		desc  string
		ifi   *net.Interface
		htype uint16
		err   error
	}{
		{
			desc: "fake interface foo0",
			ifi: &net.Interface{
				Index: 0,
				Name:  "foo0",
			},
			htype: 0,
			err:   ErrParseHardwareType,
		},
		{
			desc:  "real interface eth0 with htype 1",
			ifi:   eth0,
			htype: 1,
		},
	}

	for i, tt := range tests {
		htype, err := HardwareType(tt.ifi)
		if err != nil {
			if want, got := tt.err, err; want != got {
				t.Fatalf("[%02d] test %q, unexpected error for HardwareType(%v): %v != %v",
					i, tt.desc, tt.ifi, want, got)
			}

			continue
		}

		if want, got := tt.htype, htype; want != got {
			t.Fatalf("[%02d] test %q, unexpected error for HardwareType(%v): %v != %v",
				i, tt.desc, tt.ifi, want, got)
		}
	}
}
