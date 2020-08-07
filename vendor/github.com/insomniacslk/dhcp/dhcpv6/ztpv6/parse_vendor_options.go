package ztpv6

import (
	"errors"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

var (
	errVendorOptionMalformed = errors.New("malformed vendor option")
)

// VendorData contains fields extracted from Option 17 data
type VendorData struct {
	VendorName, Model, Serial string
}

// ParseVendorData will try to parse dhcp6 Vendor Specific Information options
// ( 16 and 17) data looking for more specific vendor data (like model, serial
// number, etc). If the options are missing we will just return nil
func ParseVendorData(packet dhcpv6.DHCPv6) (*VendorData, error) {
	// check for both options 16 and 17 if both are present will use opt 17
	opt16 := packet.GetOneOption(dhcpv6.OptionVendorClass)
	opt17 := packet.GetOneOption(dhcpv6.OptionVendorOpts)
	if (opt16 == nil) && (opt17 == nil) {
		return nil, errors.New("no vendor options or vendor class found")
	}

	vd := VendorData{}
	vData := []string{}

	if opt17 != nil {
		vo := opt17.(*dhcpv6.OptVendorOpts).VendorOpts
		for _, opt := range vo {
			vData = append(vData, string(opt.(*dhcpv6.OptionGeneric).OptionData))
		}
	} else {
		data := opt16.(*dhcpv6.OptVendorClass).Data
		for _, d := range data {
			vData = append(vData, string(d))
		}
	}
	for _, d := range vData {
		switch {
		// Arista;DCS-0000;00.00;ZZZ00000000
		case strings.HasPrefix(d, "Arista;"):
			p := strings.Split(d, ";")
			if len(p) < 4 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[3]
			return &vd, nil

		// ZPESystems:NSC:000000000
		case strings.HasPrefix(d, "ZPESystems:"):
			p := strings.Split(d, ":")
			if len(p) < 3 {
				return nil, errVendorOptionMalformed
			}

			vd.VendorName = p[0]
			vd.Model = p[1]
			vd.Serial = p[2]
			return &vd, nil
		}
	}
	return nil, errors.New("failed to parse vendor option data")
}
