// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netstat

import (
	"strings"
)

type CFMode uint8

const (
	FIB   CFMode = 0x00
	CACHE CFMode = 0xFF
)

const (
	RTFUP        uint32 = 0x0001
	RTFGATEWAY   uint32 = 0x0002
	RTFHOST      uint32 = 0x0004
	RTFREINSTATE uint32 = 0x0008
	RTFDYNAMIC   uint32 = 0x0010
	RTFMODIFIED  uint32 = 0x0020
	RTFMTU       uint32 = 0x0040
	RTFMSS       uint32 = RTFMTU
	RTFWINDOW    uint32 = 0x0080
	RTFIRTT      uint32 = 0x0100
	RTFREJECT    uint32 = 0x0200
	RTFNOTCACHED uint32 = 0x0400

	RTFDEFAULT   uint32 = 0x00010000
	RTFALLONLINK uint32 = 0x00020000
	RTFADDRCONF  uint32 = 0x00040000
	RTFNONEXTHOP uint32 = 0x00200000
	RTFEXPIRES   uint32 = 0x00400000
	RTFCACHE     uint32 = 0x01000000
	RTFFLOW      uint32 = 0x02000000
	RTFPOLICY    uint32 = 0x04000000
	RTFLOCAL     uint32 = 0x80000000
)

func convertFlagData(flag uint32) string {
	var s strings.Builder

	flags := []struct {
		flag uint32
		n    string
	}{
		{RTFUP, "U"},
		{RTFGATEWAY, "G"},
		{RTFREJECT, "!"},
		{RTFHOST, "H"},
		{RTFREINSTATE, "R"},
		{RTFDYNAMIC, "D"},
		{RTFMODIFIED, "M"},
		{RTFDEFAULT, "d"},
		{RTFALLONLINK, "a"},
		{RTFADDRCONF, "c"},
		{RTFNONEXTHOP, "o"},
		{RTFEXPIRES, "e"},
		{RTFCACHE, "c"},
		{RTFFLOW, "f"},
		{RTFPOLICY, "p"},
		{RTFLOCAL, "l"},
		{RTFMTU, "u"},
		{RTFWINDOW, "w"},
		{RTFIRTT, "i"},
		{RTFNOTCACHED, "n"},
	}

	for _, f := range flags {
		if (f.flag & flag) > 0 {
			s.WriteString(f.n)
		}
	}

	return s.String()
}
