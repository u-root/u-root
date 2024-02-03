// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esxi

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/uio/uio"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		file string
		want options
	}{
		{
			file: "testdata/kernel_cmdline_mods.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
				modules: []module{
					{
						path:    "testdata/b.b00",
						cmdline: "b.b00 blabla",
					},
					{
						path:    "testdata/k.b00",
						cmdline: "k.b00",
					},
					{
						path:    "testdata/m.m00",
						cmdline: "m.m00 marg marg2",
					},
				},
			},
		},
		{
			file: "testdata/kernelopt_first.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
			},
		},
		{
			file: "testdata/empty_mods.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
			},
		},
		{
			file: "testdata/no_mods.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
			},
		},
		{
			file: "testdata/no_cmdline.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 ",
			},
		},
		{
			file: "testdata/empty_cmdline.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 ",
			},
		},
		{
			file: "testdata/empty_updated.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
				// Explicitly stating this as the wanted value.
				updated: 0,
			},
		},
		{
			file: "testdata/updated_twice.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
				// Explicitly stating this as the wanted value.
				updated: 0,
			},
		},
		{
			file: "testdata/updated.cfg",
			want: options{
				title:   "VMware ESXi",
				kernel:  "testdata/b.b00",
				args:    "b.b00 zee",
				updated: 4,
			},
		},
		{
			file: "testdata/empty_bootstate.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
				// Explicitly stating this as the wanted value.
				bootstate: bootValid,
			},
		},
		{
			file: "testdata/bootstate_twice.cfg",
			want: options{
				title:  "VMware ESXi",
				kernel: "testdata/b.b00",
				args:   "b.b00 zee",
				// Explicitly stating this as the wanted value.
				bootstate: bootValid,
			},
		},
		{
			file: "testdata/bootstate.cfg",
			want: options{
				title:     "VMware ESXi",
				kernel:    "testdata/b.b00",
				args:      "b.b00 zee",
				bootstate: bootDirty,
			},
		},
		{
			file: "testdata/bootstate_invalid.cfg",
			want: options{
				title:     "VMware ESXi",
				kernel:    "testdata/b.b00",
				args:      "b.b00 zee",
				bootstate: bootInvalid,
			},
		},
		{
			file: "testdata/no_bootstate.cfg",
			want: options{
				title:     "VMware ESXi",
				kernel:    "testdata/b.b00",
				args:      "b.b00 zee",
				bootstate: bootInvalid,
			},
		},
	} {
		got, err := parse(tt.file)
		if err != nil {
			t.Fatalf("cannot parse config at %s: %v", tt.file, err)
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("LoadConfig(%s) = %#v want %#v", tt.file, got, tt.want)
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

// Poor man's equal.
//
// the Kernel and Modules fields will be full of uio.NewLazyFiles. We just want
// them to be pointing to the same file name; we can't compare the function
// pointers obviously. Lazy files will always print their name.
func multibootEqual(a, b []*boot.MultibootImage) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func TestDev5Valid(t *testing.T) {
	want := []*boot.MultibootImage{
		{
			Name:    "VMware ESXi from testdata/dev5",
			Kernel:  uio.NewLazyFile("testdata/k"),
			Cmdline: fmt.Sprintf(" bootUUID=%s", uuid5),
			Modules: []multiboot.Module{},
		},
	}

	opts5 := &options{
		title:     "VMware ESXi",
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}

	// No opts6 at all.
	imgs, _ := getImages(device, opts5, nil)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opts5, nil, imgs, want)
	}

	// Invalid opts6. Higher updated, but invalid state.
	invalidOpts6 := &options{
		title:     "VMware ESXi",
		kernel:    "foobar",
		updated:   2,
		bootstate: bootInvalid,
	}
	imgs, _ = getImages(device, opts5, invalidOpts6)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opts5, invalidOpts6, imgs, want)
	}
}

func TestDev6Valid(t *testing.T) {
	want := []*boot.MultibootImage{
		{
			Name:    "VMware ESXi from testdata/dev6",
			Kernel:  uio.NewLazyFile("testdata/k"),
			Cmdline: fmt.Sprintf(" bootUUID=%s", uuid6),
			Modules: []multiboot.Module{},
		},
	}

	opts6 := &options{
		title:     "VMware ESXi",
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}

	// No opts5 at all.
	imgs, err := getImages(device, nil, opts6)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v (err %v)", device, nil, opts6, imgs, want, err)
	}

	// Invalid opts5. Higher updated, but invalid state.
	invalidOpts5 := &options{
		title:     "VMware ESXi",
		kernel:    "foobar",
		updated:   2,
		bootstate: bootInvalid,
	}
	imgs, _ = getImages(device, invalidOpts5, opts6)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, invalidOpts5, opts6, imgs, want)
	}
}

func TestImageOrder(t *testing.T) {
	prevGetBlockSize := getBlockSize
	defer func() {
		getBlockSize = prevGetBlockSize
	}()
	getBlockSize = func(dev string) (int, error) {
		return 512, nil
	}

	opt5 := &options{
		title:     "VMware ESXi",
		kernel:    "foobar",
		updated:   2,
		bootstate: bootValid,
	}
	want5 := &boot.MultibootImage{
		Name:    "VMware ESXi from testdata/dev5",
		Kernel:  uio.NewLazyFile("foobar"),
		Cmdline: fmt.Sprintf(" bootUUID=%s", uuid5),
		Modules: []multiboot.Module{},
	}

	opt6 := &options{
		title:     "VMware ESXi",
		kernel:    "testdata/k",
		updated:   1,
		bootstate: bootValid,
	}
	want6 := &boot.MultibootImage{
		Name:    "VMware ESXi from testdata/dev6",
		Kernel:  uio.NewLazyFile("testdata/k"),
		Cmdline: fmt.Sprintf(" bootUUID=%s", uuid6),
		Modules: []multiboot.Module{},
	}

	// Way 1.
	want := []*boot.MultibootImage{want5, want6}
	imgs, _ := getImages(device, opt5, opt6)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opt5, opt6, imgs, want)
	}

	opt5.updated = 1
	opt6.updated = 2
	// Vice versa priority.
	want = []*boot.MultibootImage{want6, want5}
	imgs, _ = getImages(device, opt5, opt6)
	if !multibootEqual(imgs, want) {
		t.Fatalf("getImages(%s, %v, %v) = %v, want %v", device, opt5, opt6, imgs, want)
	}
}

func FuzzParse(f *testing.F) {
	seeds, err := filepath.Glob("testdata/*.cfg")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}
	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			return
		}

		parse(string(data))
	})
}
