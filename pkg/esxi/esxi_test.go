// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esxi

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/boot"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		file string
		want *boot.MultibootImage
	}{
		{
			file: "testdata/kernel_cmdline_mods.cfg",
			want: &boot.MultibootImage{
				Path:    "testdata/b.b00",
				Cmdline: "zee",
				Modules: []string{
					"testdata/b.b00 blabla",
					"testdata/k.b00",
					"testdata/m.m00 marg marg2",
				},
			},
		},
		{
			file: "testdata/empty_mods.cfg",
			want: &boot.MultibootImage{
				Path:    "testdata/b.b00",
				Cmdline: "zee",
			},
		},
		{
			file: "testdata/no_mods.cfg",
			want: &boot.MultibootImage{
				Path:    "testdata/b.b00",
				Cmdline: "zee",
			},
		},
		{
			file: "testdata/no_cmdline.cfg",
			want: &boot.MultibootImage{
				Path: "testdata/b.b00",
			},
		},
	} {
		got, err := LoadConfig(tt.file)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(got, tt.want) {
			t.Errorf("LoadConfig(./testdata/boot.cfg) = %s want %s:\ndifferences %s", got, tt.want, cmp.Diff(got, tt.want))
		}
	}
}
