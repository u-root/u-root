package ztpv4

import (
	"errors"
	"strings"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

// VendorData is optional data a particular vendor may or may not include
// in the Vendor Class options.
type VendorData struct {
	VendorName, Model, Serial string
}

var errVendorOptionMalformed = errors.New("malformed vendor option")

// ParseVendorData will try to parse dhcp4 options looking for more
// specific vendor data (like model, serial number, etc).
func ParseVendorData(packet *dhcpv4.DHCPv4) (*VendorData, error) {
	vc := packet.ClassIdentifier()
	if len(vc) == 0 {
		return nil, errors.New("vendor options not found")
	}
	vd := &VendorData{}

	switch {
	// Arista;DCS-7050S-64;01.23;JPE12221671
	case strings.HasPrefix(vc, "Arista;"):
		p := strings.Split(vc, ";")
		if len(p) < 4 {
			return nil, errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[3]
		return vd, nil

	// ZPESystems:NSC:002251623
	case strings.HasPrefix(vc, "ZPESystems:"):
		p := strings.Split(vc, ":")
		if len(p) < 3 {
			return nil, errVendorOptionMalformed
		}

		vd.VendorName = p[0]
		vd.Model = p[1]
		vd.Serial = p[2]
		return vd, nil

	// Juniper option 60 parsing is a bit more nuanced.  The following are all
	// "valid" identifying stings for Juniper:
	//    Juniper-ptx1000-DD576      <vendor>-<model>-<serial
	//    Juniper-qfx10008           <vendor>-<model> (serial in hostname option)
	//    Juniper-qfx10002-361-DN817 <vendor>-<model>-<serial> (model has a dash in it!)
	case strings.HasPrefix(vc, "Juniper-"):
		p := strings.Split(vc, "-")
		if len(p) < 3 {
			vd.Model = p[1]
			vd.Serial = packet.HostName()
			if len(vd.Serial) == 0 {
				return nil, errors.New("host name option is missing")
			}
		} else {
			vd.Model = strings.Join(p[1:len(p)-1], "-")
			vd.Serial = p[len(p)-1]
		}

		vd.VendorName = p[0]
		return vd, nil
	}

	// We didn't match anything.
	return nil, errors.New("no known ZTP vendor found")
}
