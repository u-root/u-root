package dhcp4

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/d2g/dhcp4"
)

func TestPacketMarshalBinary(t *testing.T) {
	for i, tt := range []struct {
		packet  *Packet
		wantD2G func() dhcp4.Packet
		err     error
	}{
		{
			packet: NewPacket(BootRequest),
			wantD2G: func() dhcp4.Packet {
				return dhcp4.NewPacket(dhcp4.BootRequest)
			},
			err: nil,
		},
		{
			packet: &Packet{
				Op:            BootRequest,
				HType:         1,
				Hops:          2,
				TransactionID: [4]byte{0xa, 0xb, 0xc, 0xd},
				Broadcast:     true,
				CIAddr:        net.IP{0xff, 0xee, 0xdd, 0xcc},
				YIAddr:        net.IP{192, 168, 0, 1},
				SIAddr:        net.IP{2, 3, 4, 5},
				CHAddr:        net.HardwareAddr{0xfe, 0xab, 0x67},
			},
			wantD2G: func() dhcp4.Packet {
				p := dhcp4.NewPacket(dhcp4.BootRequest)
				p.SetHops(2)
				p.SetXId([]byte{0xa, 0xb, 0xc, 0xd})
				p.SetBroadcast(true)
				p.SetCIAddr(net.IP{0xff, 0xee, 0xdd, 0xcc})
				p.SetYIAddr(net.IP{192, 168, 0, 1})
				p.SetSIAddr(net.IP{2, 3, 4, 5})
				p.SetGIAddr(nil)
				p.SetCHAddr(net.HardwareAddr{0xfe, 0xab, 0x67})
				return p
			},
			err: nil,
		},
	} {
		t.Run(fmt.Sprintf("Test [%02d]", i), func(t *testing.T) {
			got, err := tt.packet.MarshalBinary()
			if err != tt.err {
				t.Fatalf("Marshal(%#v) got error %v, want %v", tt.packet, err, tt.err)
			} else if err != nil {
				return
			}

			want := tt.wantD2G()
			if len(got) != len(want) {
				t.Errorf("Marshal(%#v) got length %d, want length %d", tt.packet, len(got), len(want))
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("Marshal(%#v) got %v, want %v", tt.packet, got, want)
			}
		})
	}
}

func TestPacketUnmarshalBinary(t *testing.T) {
	for i, tt := range []struct {
		packet func() dhcp4.Packet
		want   *Packet
		err    error
	}{
		{
			packet: func() dhcp4.Packet {
				p := dhcp4.NewPacket(dhcp4.BootRequest)
				p.SetBroadcast(true)
				return p
			},
			want: &Packet{
				Op:        BootRequest,
				HType:     1,
				Broadcast: true,
				CIAddr:    net.IP{0, 0, 0, 0},
				YIAddr:    net.IP{0, 0, 0, 0},
				SIAddr:    net.IP{0, 0, 0, 0},
				GIAddr:    net.IP{0, 0, 0, 0},
				CHAddr:    net.HardwareAddr{},
				Options:   Options{},
			},
		},
		{
			packet: func() dhcp4.Packet {
				p := dhcp4.NewPacket(dhcp4.BootRequest)
				p.SetHops(2)
				p.SetXId([]byte{0xa, 0xb, 0xc, 0xd})
				p.SetBroadcast(true)
				p.SetCIAddr(net.IP{0xff, 0xee, 0xdd, 0xcc})
				p.SetYIAddr(net.IP{192, 168, 0, 1})
				p.SetSIAddr(net.IP{2, 3, 4, 5})
				p.SetGIAddr(nil)
				p.SetCHAddr(net.HardwareAddr{0xfe, 0xab, 0x67})
				return p
			},
			want: &Packet{
				Op:            BootRequest,
				HType:         1,
				Hops:          2,
				TransactionID: [4]byte{0xa, 0xb, 0xc, 0xd},
				Broadcast:     true,
				CIAddr:        net.IP{0xff, 0xee, 0xdd, 0xcc},
				YIAddr:        net.IP{192, 168, 0, 1},
				SIAddr:        net.IP{2, 3, 4, 5},
				GIAddr:        net.IP{0, 0, 0, 0},
				CHAddr:        net.HardwareAddr{0xfe, 0xab, 0x67},
				Options:       Options{},
			},
		},
	} {
		t.Run(fmt.Sprintf("Test %02d", i), func(t *testing.T) {
			input := tt.packet()
			got := new(Packet)
			if err := got.UnmarshalBinary(input); err != tt.err {
				t.Fatalf("got %v, want %v", err, tt.err)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
