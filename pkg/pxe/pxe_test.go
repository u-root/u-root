// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pxe

import (
	"net"
	"reflect"
	"testing"
)

func TestProbeFiles(t *testing.T) {
	// Anyone got some ideas for other test cases?
	for _, tt := range []struct {
		mac   net.HardwareAddr
		ip    net.IP
		files []string
	}{
		{
			mac: []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			ip:  []byte{192, 168, 0, 1},
			files: []string{
				"01-aa-bb-cc-dd-ee-ff",
				"C0A80001",
				"C0A8000",
				"C0A800",
				"C0A80",
				"C0A8",
				"C0A",
				"C0",
				"C",
				"default",
			},
		},
		{
			mac: []byte{0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd},
			ip:  []byte{192, 168, 2, 91},
			files: []string{
				"01-88-99-aa-bb-cc-dd",
				"C0A8025B",
				"C0A8025",
				"C0A802",
				"C0A80",
				"C0A8",
				"C0A",
				"C0",
				"C",
				"default",
			},
		},
	} {
		got := probeFiles(tt.mac, tt.ip)
		if !reflect.DeepEqual(got, tt.files) {
			t.Errorf("probeFiles(%s, %s) = %v, want %v", tt.mac, tt.ip, got, tt.files)
		}
	}
}
