// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestExtraMounts(t *testing.T) {
	for _, tt := range []struct {
		name    string
		extra   string
		wantErr string
	}{
		{
			name: "mountList == ''",
		},
		{
			name:    "len(bin) == 0",
			extra:   "::",
			wantErr: "[\"\" \"\" \"\"] is not in the form src:target",
		},
		{
			name:  "switch case 1",
			extra: "/tmp, /etc"},
		{
			name:  "switch case 2",
			extra: "/tmp:/tmp,/etc:/etc",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*extra = tt.extra
			if got := extraMounts(tt.extra); got != nil {
				if !strings.Contains(got.Error(), tt.wantErr) {
					t.Errorf("extraMounts() = %q, want: %q", got.Error(), tt.wantErr)
				}
			}
		})
	}
}
