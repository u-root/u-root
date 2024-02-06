// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"reflect"
	"testing"
)

func TestLinuxModifiers(t *testing.T) {
	for _, tt := range []struct {
		images    []OSImage
		modifiers []LinuxModifier
		want      []OSImage
	}{
		{
			images: []OSImage{},
			want:   []OSImage{},
		},
		{
			images: nil,
			want:   nil,
		},
		{
			images: []OSImage{
				&LinuxImage{
					Cmdline: "foobar",
				},
				&MultibootImage{
					Cmdline: "blabla",
				},
			},
			modifiers: []LinuxModifier{
				func(img *LinuxImage) {
					img.Cmdline += " andsoon"
				},
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "foobar andsoon",
				},
				&MultibootImage{
					Cmdline: "blabla",
				},
			},
		},
		{
			images: []OSImage{
				&LinuxImage{
					Cmdline: "foobar",
				},
			},
			modifiers: []LinuxModifier{
				AppendLinux("andsoon"),
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "foobar andsoon",
				},
			},
		},
		{
			images: []OSImage{
				&LinuxImage{
					Cmdline: "foobar",
				},
			},
			modifiers: []LinuxModifier{
				PrependLinux("andsoon"),
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "andsoon foobar",
				},
			},
		},
		{
			images: []OSImage{
				&LinuxImage{},
			},
			modifiers: []LinuxModifier{
				PrependLinux("andsoon"),
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "andsoon",
				},
			},
		},
		{
			images: []OSImage{
				&LinuxImage{},
			},
			modifiers: []LinuxModifier{
				AppendLinux("andsoon"),
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "andsoon",
				},
			},
		},
	} {
		ApplyLinuxModifiers(tt.images, tt.modifiers...)
		if got := tt.images; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ApplyLinuxModifiers = %v, want %v", got, tt.want)
		}
	}
}

func TestMultibootModifiers(t *testing.T) {
	for _, tt := range []struct {
		images    []OSImage
		modifiers []MultibootModifier
		want      []OSImage
	}{
		{
			images: []OSImage{},
			want:   []OSImage{},
		},
		{
			images: nil,
			want:   nil,
		},
		{
			images: []OSImage{
				&LinuxImage{
					Cmdline: "foobar",
				},
				&MultibootImage{
					Cmdline: "blabla",
				},
			},
			modifiers: []MultibootModifier{
				func(img *MultibootImage) {
					img.Cmdline += " andsoon"
				},
			},
			want: []OSImage{
				&LinuxImage{
					Cmdline: "foobar",
				},
				&MultibootImage{
					Cmdline: "blabla andsoon",
				},
			},
		},
	} {
		ApplyMultibootModifiers(tt.images, tt.modifiers...)
		if got := tt.images; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ApplyMultibootModifiers = %v, want %v", got, tt.want)
		}
	}
}
