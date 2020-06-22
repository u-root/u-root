// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import (
	"net"
	"reflect"
	"testing"
)

func TestParseURI(t *testing.T) {
	for _, tt := range []struct {
		uri    string
		target *net.TCPAddr
		volume string
		want   string
	}{
		{
			uri:    "iscsi:192.168.1.1:::1:iqn.com.oracle:boot",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.oracle:boot",
		},
		{
			uri:    "iscsi:@192.168.1.1::3260::iqn.com.oracle:boot",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.oracle:boot",
		},
		{
			uri:    "iscsi:[fe80::1]::3260::iqn.com.google:esxi-boot-image",
			target: &net.TCPAddr{IP: net.ParseIP("fe80::1"), Port: 3260},
			volume: "iqn.com.google:esxi-boot-image",
		},
		{
			uri:    "iscsi:[fe80::1]::3260::iqn.com.google:esxi-boot-]:image",
			target: &net.TCPAddr{IP: net.ParseIP("fe80::1"), Port: 3260},
			volume: "iqn.com.google:esxi-boot-]:image",
		},
		{
			uri:    "iscsi:192.168.1.1::3260::iqn.com.google:[fe80::1]",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.google:[fe80::1]",
		},
		{
			uri:    "iscsi:[fe80::1]::3260::iqn.com.google:[foobar]",
			target: &net.TCPAddr{IP: net.ParseIP("fe80::1"), Port: 3260},
			volume: "iqn.com.google:[foobar]",
		},
		{
			uri:    "iscsi:192.168.1.1::3260::iqn.com.google:esxi-boot-image",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.google:esxi-boot-image",
		},
		{
			uri:    "iscsi:192.168.1.1::3000::iqn.com.google:esxi-boot-image",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3000},
			volume: "iqn.com.google:esxi-boot-image",
		},
		{
			uri:    "iscsi:192.168.1.1::::iqn.com.google:esxi-boot-image",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.google:esxi-boot-image",
		},
		{
			uri:    "iscsi:192.168.1.1::::iqn.com.google::::",
			target: &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 3260},
			volume: "iqn.com.google::::",
		},
		{
			uri:  "iscsi:192.168.1.1::::",
			want: "iSCSI URI \"iscsi:192.168.1.1::::\" is missing a volume name",
		},
		{
			uri:  "iscsi:192.168.1.1:::",
			want: "iSCSI URI \"iscsi:192.168.1.1:::\" failed to parse: fields missing",
		},
		{
			uri:  "iscs:192.168.1.1::::",
			want: "iSCSI URI \"iscs:192.168.1.1::::\" is missing iscsi scheme prefix, have iscs",
		},
		{
			uri:  "",
			want: "iSCSI URI \"\" failed to parse: fields missing",
		},
		{
			uri:  "iscsi:192.168.1.1::foobar::volume",
			want: "iSCSI URI \"iscsi:192.168.1.1::foobar::volume\" has invalid port: strconv.Atoi: parsing \"foobar\": invalid syntax",
		},
		{
			uri:  "iscsi:[fe80::1::::",
			want: "iSCSI URI \"iscsi:[fe80::1::::\" failed to parse: invalid IPv6 address",
		},
	} {
		gtarget, gvolume, got := ParseISCSIURI(tt.uri)
		if (got != nil && got.Error() != tt.want) || (got == nil && len(tt.want) > 0) {
			t.Errorf("parseISCSIURI(%s) = %v, want %v", tt.uri, got, tt.want)
		}
		if gvolume != tt.volume {
			t.Errorf("parseISCSIURI(%s) = volume %s, want %s", tt.uri, gvolume, tt.volume)
		}
		if !reflect.DeepEqual(gtarget, tt.target) {
			t.Errorf("parseISCSIURI(%s) = target %s, want %s", tt.uri, gtarget, tt.target)
		}
	}
}
