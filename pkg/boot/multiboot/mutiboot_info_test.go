// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInfoMarshal(t *testing.T) {
	for _, tt := range []struct {
		name string
		mi   *esxBootInfoInfo
		want []byte
	}{
		{
			name: "no elements",
			mi: &esxBootInfoInfo{
				cmdline: 0xdeadbeef,
				elems:   nil,
			},
			want: []byte{
				// cmdline
				0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
				// 0 elements
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		{
			name: "one memrange element",
			mi: &esxBootInfoInfo{
				cmdline: 0xdeadbeef,
				elems: []elem{
					&esxBootInfoMemRange{
						startAddr: 0xbeefdead,
						length:    0xdeadbeef,
						memType:   2,
					},
				},
			},
			want: []byte{
				// cmdline
				0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
				// 1 element
				0x1, 0, 0, 0, 0, 0, 0, 0,

				// TLV -- type, length, value

				// type
				byte(ESXBOOTINFO_MEMRANGE_TYPE), 0, 0, 0,
				// length - 20 bytes + 8 for the length + 4 for the type
				32, 0, 0, 0, 0, 0, 0, 0,
				// values
				0xad, 0xde, 0xef, 0xbe, 0, 0, 0, 0,
				0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
				2, 0, 0, 0,
			},
		},
		{
			name: "one module element",
			mi: &esxBootInfoInfo{
				cmdline: 0xdeadbeef,
				elems: []elem{
					&esxBootInfoModule{
						cmdline:    0xbeefdead,
						moduleSize: 0x1000,
						ranges: []esxBootInfoModuleRange{
							{
								startPageNum: 0x100,
								numPages:     1,
							},
						},
					},
				},
			},
			want: []byte{
				// cmdline
				0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
				// 1 element
				0x1, 0, 0, 0, 0, 0, 0, 0,

				// TLV -- type, length, value

				// type
				byte(ESXBOOTINFO_MODULE_TYPE), 0, 0, 0,
				// length - 36 bytes + 8 for the length + 4 for the type
				48, 0, 0, 0, 0, 0, 0, 0,
				// values
				// cmdline
				0xad, 0xde, 0xef, 0xbe, 0, 0, 0, 0,
				// moduleSize
				0x00, 0x10, 0, 0, 0, 0, 0, 0,
				// numRanges
				1, 0, 0, 0,
				// range - startPageNum
				0x00, 0x01, 0, 0, 0, 0, 0, 0,
				// range - numPages
				1, 0, 0, 0,
				// padding
				0, 0, 0, 0,
			},
		},
		{
			name: "one zero-length module element",
			mi: &esxBootInfoInfo{
				cmdline: 0xdeadbeef,
				elems: []elem{
					&esxBootInfoModule{
						cmdline:    0xbeefdead,
						moduleSize: 0,
					},
				},
			},
			want: []byte{
				// cmdline
				0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
				// 1 element
				0x1, 0, 0, 0, 0, 0, 0, 0,

				// TLV -- type, length, value

				// type
				byte(ESXBOOTINFO_MODULE_TYPE), 0, 0, 0,
				// length - 20 bytes + 8 for the length + 4 for the type
				32, 0, 0, 0, 0, 0, 0, 0,
				// values
				// cmdline
				0xad, 0xde, 0xef, 0xbe, 0, 0, 0, 0,
				// moduleSize
				0, 0, 0, 0, 0, 0, 0, 0,
				// numRanges
				0, 0, 0, 0,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mi.marshal()
			if !cmp.Equal(got, tt.want) {
				t.Errorf("marshaled bytes not the same. diff (-want, +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
