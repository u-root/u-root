// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ibft

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

var spewConfig = &spew.ConfigState{
	Indent:                  "  ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

// Header is described in iBFT Spec 1.4.3.
func emptyHeader() []byte {
	return append([]byte("iBFT\x00\x00\x00\x00\x01\x00GoogleGoogIBFT"), bytes.Repeat([]byte{0}, 24)...)
}

func defaultControl() []byte {
	return []byte{
		0x01 /* Control ID */, 0x01 /* version */, 18, 0 /* length */, 0 /* index */, 0 /* flags */, 0, 0, /* extensions */
		0x48, 0 /* Initiator */, 0x98, 0x0 /* NIC0 */, 0x0, 0x01 /* Target0 */, 0x0, 0x0 /* NIC1 */, 0x0, 0x0, /* Target1 */
		0, 0, 0, 0, 0, 0, /* padding */
	}
}

// Offsets are from iBFT Spec Section 1.4.5.
func initiator(flags byte, nameLen uint16, nameOffset uint16) []byte {
	b := append(append(
		[]byte{0x02 /* Initiator ID */, 0x01 /* version */, 74, 0 /* length */, 0 /* index */, flags /* flags */},
		bytes.Repeat([]byte{0}, 16*4)...),
		0, 0 /* initiator name length */, 0, 0, /* initiator name offset */
		0, 0, 0, 0, 0, 0, /* padding */
	)
	binary.LittleEndian.PutUint16(b[70:], nameLen)
	binary.LittleEndian.PutUint16(b[72:], nameOffset)
	return b
}

func emptyNIC() []byte {
	return nic(0, nil, 0, nil, nil, nil, nil, 0)
}

// Offsets are from iBFT Spec Section 1.4.6.
func nic(flags uint8, ip net.IP, subnetMaskPrefix uint8, gateway, dns, dhcp net.IP, mac net.HardwareAddr, bdf uint16) []byte {
	empty := append(
		append(
			[]byte{0x03 /* NIC ID */, 0x01 /* version */, 102, 0 /* length */, 0 /* index */, flags /* flags */},
			bytes.Repeat([]byte{0}, 16+1+1+16*4+2+6+2+2+2)..., /* all fields */
		),
		0, 0, /* padding for alignment */
	)
	if ip != nil {
		copy(empty[6:], ip.To16())
	}
	empty[22] = subnetMaskPrefix
	if gateway != nil {
		copy(empty[24:], gateway.To16())
	}
	if dns != nil {
		copy(empty[40:], dns.To16())
	}
	if dhcp != nil {
		copy(empty[72:], dhcp.To16())
	}
	if mac != nil {
		copy(empty[90:], mac)
	}
	binary.LittleEndian.PutUint16(empty[96:], bdf)
	return empty
}

func emptyTarget() []byte {
	return target(0, nil, 0, 0, 0)
}

// Offsets are from iBFT Spec Section 1.4.7.
func target(flags uint8, ip net.IP, port uint16, nameLen uint16, nameOffset uint16) []byte {
	empty := append(
		[]byte{0x04 /* Target ID */, 0x01 /* version */, 54, 0 /* length */, 0 /* index */, flags /* flags */},
		bytes.Repeat([]byte{0}, 49+1)..., /* all fields + 1 bytes padding for alignment */
	)
	if ip != nil {
		copy(empty[6:], ip.To16())
	}
	binary.LittleEndian.PutUint16(empty[22:], port)
	binary.LittleEndian.PutUint16(empty[34:], nameLen)
	binary.LittleEndian.PutUint16(empty[36:], nameOffset)
	return empty
}

func heap(s []string) []byte {
	var h []byte
	for _, t := range s {
		h = append(h, []byte(t)...)
		h = append(h, 0)
	}
	return h
}

func join(b ...[]byte) []byte {
	var r []byte
	for _, bb := range b {
		r = append(r, bb...)
	}
	return r
}

func TestMarshal(t *testing.T) {
	for _, tt := range []struct {
		desc string
		i    *IBFT
		want []byte
	}{
		{
			desc: "empty iBFT",
			i:    &IBFT{},
			want: fixACPIHeader(join(emptyHeader(), defaultControl(), initiator(0, 0, 0), emptyNIC(), emptyTarget())),
		},
		{
			desc: "IBFT with relevant stuff.",
			i: &IBFT{
				Initiator: Initiator{
					Valid: true,
					Boot:  true,
					Name:  "NERF",
				},
				NIC0: NIC{
					Valid:  true,
					Boot:   true,
					Global: true,
					IPNet: &net.IPNet{
						IP:   net.IP{192, 168, 1, 15},
						Mask: net.IPv4Mask(255, 255, 255, 0),
					},
					Gateway:    net.IP{192, 168, 1, 1},
					PrimaryDNS: net.IP{8, 8, 8, 8},
					DHCPServer: net.IP{192, 168, 1, 1},
					MACAddress: net.HardwareAddr{52, 54, 0o0, 12, 34, 56},
				},
				Target0: Target{
					Valid: true,
					Boot:  true,
					Target: &net.TCPAddr{
						IP:   net.IP{192, 168, 1, 1},
						Port: 3260,
					},
					TargetName: "iqn.2016-01.com.example:foo",
				},
			},
			want: fixACPIHeader(join(
				emptyHeader(),
				defaultControl(),
				initiator(1<<1|1, 4 /* len */, 0x138),
				nic(1<<2|1<<1|1, net.IP{192, 168, 1, 15} /* ip */, 24 /* netmask */, net.IP{192, 168, 1, 1} /* gateway */, net.IP{8, 8, 8, 8} /* dns */, net.IP{192, 168, 1, 1} /* dhcp serv */, net.HardwareAddr{52, 54, 0o0, 12, 34, 56} /* mac */, 0 /* BDF */),
				target(1<<1|1, net.IP{192, 168, 1, 1}, 3260, 27, 0x138+5),
				heap([]string{"NERF", "iqn.2016-01.com.example:foo"})),
			),
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			got := tt.i.Marshal()

			t.Logf("got:\n%s", hex.Dump(got))
			t.Logf("want:\n%s", hex.Dump(tt.want))

			if !cmp.Equal(got, tt.want) {
				t.Errorf("IBFT(%s).Marshal() differences: %s", spewConfig.Sdump(tt.i), cmp.Diff(got, tt.want))
			}
		})
	}
}
