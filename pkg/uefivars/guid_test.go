// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	for _, td := range []struct {
		name, want string
		in         MixedGUID
	}{
		{
			name: "1",
			in:   MixedGUID{0xCD, 0x5C, 0x63, 0x81, 0x4F, 0x1B, 0x3F, 0x4D, 0xB7, 0xB7, 0xF7, 0x8A, 0x5B, 0x02, 0x9F, 0x35},
			want: "81635ccd-1b4f-4d3f-b7b7-f78a5b029f35",
		}, {
			name: "2",
			in:   MixedGUID{0xa2, 0xd1, 0x1b, 0x1d, 0xd9, 0x0f, 0xe9, 0x41, 0xbb, 0xb5, 0xa9, 0x8b, 0xac, 0x57, 0x0b, 0x2a},
			want: "1d1bd1a2-0fd9-41e9-bbb5-a98bac570b2a",
		}, {
			name: "3",
			in:   MixedGUID{0x3e, 0x14, 0xbe, 0xcf, 0x9e, 0x5e, 0x25, 0x46, 0xa5, 0x00, 0xc3, 0xf0, 0x36, 0x20, 0x04, 0x11},
			want: "cfbe143e-5e9e-4625-a500-c3f036200411",
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			mstr := td.in.String()
			if mstr != td.want {
				t.Errorf("mismatch\nwant %s\n got %s", td.want, mstr)
			}
			std := td.in.ToStdEnc()
			sstr := std.String()
			if sstr != td.want {
				t.Errorf("mismatch\nwant %s\n got %s", td.want, sstr)
			}
			guid := std.ToMixedGUID()
			if !bytes.Equal(guid[:], td.in[:]) {
				t.Errorf("mismatch\nwant %x\n got %x", td.in, guid)
			}
		})
	}
}
