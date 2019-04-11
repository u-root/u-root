// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func withNetbootInfo(bootFileName, serverHostName string) dhcpv4.Modifier {
	return func(m *dhcpv4.DHCPv4) {
		m.BootFileName = bootFileName
		m.ServerHostName = serverHostName
	}
}

func mustNew(t *testing.T, modifiers ...dhcpv4.Modifier) *dhcpv4.DHCPv4 {
	m, err := dhcpv4.New(modifiers...)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	return m
}

func TestBoot(t *testing.T) {
	for i, tt := range []struct {
		message *dhcpv4.DHCPv4
		want    *url.URL
		err     error
	}{
		{
			message: mustNew(t),
			err:     ErrNoBootFile,
		},
		{
			message: mustNew(t,
				withNetbootInfo("pxelinux.0", "10.0.0.1"),
			),
			want: &url.URL{
				Scheme: "tftp",
				Host:   "10.0.0.1",
				Path:   "pxelinux.0",
			},
		},
		{
			message: mustNew(t,
				withNetbootInfo("pxelinux.0", ""),
			),
			err: ErrNoServerHostName,
		},
		{
			message: mustNew(t,
				withNetbootInfo("pxelinux.0", ""),
				dhcpv4.WithServerIP(net.IP{10, 0, 0, 2}),
			),
			want: &url.URL{
				Scheme: "tftp",
				Host:   "10.0.0.2",
				Path:   "pxelinux.0",
			},
		},
		{
			message: mustNew(t,
				withNetbootInfo("pxelinux.0", ""),
				dhcpv4.WithServerIP(net.IP{10, 0, 0, 2}),
				dhcpv4.WithOption(dhcpv4.OptServerIdentifier(net.IP{10, 0, 0, 3})),
			),
			want: &url.URL{
				Scheme: "tftp",
				Host:   "10.0.0.3",
				Path:   "pxelinux.0",
			},
		},
		{
			message: mustNew(t,
				withNetbootInfo("pxelinux.0", "10.0.0.1"),
				dhcpv4.WithServerIP(net.IP{10, 0, 0, 2}),
				dhcpv4.WithOption(dhcpv4.OptServerIdentifier(net.IP{10, 0, 0, 3})),
			),
			want: &url.URL{
				Scheme: "tftp",
				Host:   "10.0.0.1",
				Path:   "pxelinux.0",
			},
		},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			p := NewPacket4(nil, tt.message)
			got, err := p.Boot()
			if err != tt.err {
				t.Errorf("Boot() = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Boot() = %s, want %s", got, tt.want)
			}
		})
	}
}
