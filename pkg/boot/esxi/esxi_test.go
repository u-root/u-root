// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esxi

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		file string
		want options
	}{
		{
			file: "testdata/kernel_cmdline_mods.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
				modules: []string{
					"testdata/b.b00 blabla",
					"testdata/k.b00",
					"testdata/m.m00 marg marg2",
				},
			},
		},
		{
			file: "testdata/empty_mods.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
			},
		},
		{
			file: "testdata/no_mods.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
			},
		},
		{
			file: "testdata/no_cmdline.cfg",
			want: options{
				kernel: "testdata/b.b00",
			},
		},
		{
			file: "testdata/empty_cmdline.cfg",
			want: options{
				kernel: "testdata/b.b00",
			},
		},
		{
			file: "testdata/empty_updated.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
				// Explicitly stating this as the wanted value.
				updated: 0,
			},
		},
		{
			file: "testdata/updated_twice.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
				// Explicitly stating this as the wanted value.
				updated: 0,
			},
		},
		{
			file: "testdata/updated.cfg",
			want: options{
				kernel:  "testdata/b.b00",
				args:    "zee",
				updated: 4,
			},
		},
		{
			file: "testdata/empty_bootstate.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
				// Explicitly stating this as the wanted value.
				bootstate: bootValid,
			},
		},
		{
			file: "testdata/bootstate_twice.cfg",
			want: options{
				kernel: "testdata/b.b00",
				args:   "zee",
				// Explicitly stating this as the wanted value.
				bootstate: bootValid,
			},
		},
		{
			file: "testdata/bootstate.cfg",
			want: options{
				kernel:    "testdata/b.b00",
				args:      "zee",
				bootstate: bootDirty,
			},
		},
		{
			file: "testdata/bootstate_invalid.cfg",
			want: options{
				kernel:    "testdata/b.b00",
				args:      "zee",
				bootstate: bootInvalid,
			},
		},
		{
			file: "testdata/no_bootstate.cfg",
			want: options{
				kernel:    "testdata/b.b00",
				args:      "zee",
				bootstate: bootInvalid,
			},
		},
	} {
		got, err := parse(tt.file)
		if err != nil {
			t.Fatalf("cannot parse config at %s: %v", tt.file, err)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("LoadConfig(%s) = %v want %v", tt.file, got, tt.want)
		}
	}
}

// This is in the second block of testdata/dev5 and testdata/dev6.
var (
	dev5GUID = "aabbccddeeff0011"
	dev6GUID = "00112233445566aa"
	uuid5    = hex.EncodeToString([]byte(dev5GUID))
	uuid6    = hex.EncodeToString([]byte(dev6GUID))
	device   = "testdata/dev"
)

func TestDev5Valid(t *testing.T) {
	want := []*boot.MultibootImage{
		{
			Path:    "testdata/k",
			Cmdline: fmt.Sprintf(" bootUUID=%s", uuid5),
		},
	}

	opts5 := &options{
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}

	// No opts6 at all.
	imgs, _ := getImages(device, opts5, nil)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opts5, nil, imgs, want)
	}

	// Invalid opts6. Higher updated, but invalid state.
	invalidOpts6 := &options{
		kernel:    "foobar",
		updated:   2,
		bootstate: bootInvalid,
	}
	imgs, _ = getImages(device, opts5, invalidOpts6)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opts5, invalidOpts6, imgs, want)
	}
}

func TestDev6Valid(t *testing.T) {
	want := []*boot.MultibootImage{
		{
			Path:    "testdata/k",
			Cmdline: fmt.Sprintf(" bootUUID=%s", uuid6),
		},
	}

	opts6 := &options{
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}

	// No opts5 at all.
	imgs, _ := getImages(device, nil, opts6)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, nil, opts6, imgs, want)
	}

	// Invalid opts5. Higher updated, but invalid state.
	invalidOpts5 := &options{
		kernel:    "foobar",
		updated:   2,
		bootstate: bootInvalid,
	}
	imgs, _ = getImages(device, invalidOpts5, opts6)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, invalidOpts5, opts6, imgs, want)
	}
}

func TestImageOrder(t *testing.T) {
	getBlockSize = func(dev string) (int, error) {
		return 512, nil
	}

	opt5 := &options{
		kernel:    "foobar",
		updated:   2,
		bootstate: bootValid,
	}
	want5 := &boot.MultibootImage{
		Path:    "foobar",
		Cmdline: fmt.Sprintf(" bootUUID=%s", uuid5),
	}

	opt6 := &options{
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}
	want6 := &boot.MultibootImage{
		Path:    "testdata/k",
		Cmdline: fmt.Sprintf(" bootUUID=%s", uuid6),
	}

	// Way 1.
	want := []*boot.MultibootImage{want5, want6}
	imgs, _ := getImages(device, opt5, opt6)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opt5, opt6, imgs, want)
	}

	opt5.updated = 1
	opt6.updated = 2
	// Vice versa priority.
	want = []*boot.MultibootImage{want6, want5}
	imgs, _ = getImages(device, opt5, opt6)
	if !reflect.DeepEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opt5, opt6, imgs, want)
	}
}
