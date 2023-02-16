// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/vfile"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

const (
	// Number of configs in the fitimage.itb
	fbcCnt = 2
	// Size in bytes the of content for each image in the fitimage.itb
	// This arbitrary value is defined in testdata/README.md for data creation.
	imageSize = 100
	// Fill patterns for the image content
	kernel0Fill  = "k0"
	kernel1Fill  = "k1"
	initram0Fill = "i0"
)

func TestLoadConfig(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	kn, rn, err := i.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("kernel name: %s", kn)
	t.Logf("ramdisk name: %s", rn)
	if kn != "kernel@0" {
		t.Errorf("Expected kernel %s, got %s", "kernel@0", kn)
	}
	if rn != "ramdisk@0" {
		t.Errorf("Expected ramdisk %s, got %s", "ramdisk@0", rn)
	}
}

func TestLoadConfigMiss(t *testing.T) {
	i, err := New("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	i.ConfigOverride = "MagicNonExistentConfig"
	kn, rn, err := i.LoadConfig()

	if kn != "" {
		t.Errorf("Kernel %s returned on expected config miss", kn)
	}

	if rn != "" {
		t.Errorf("Initramfs %s returned on expected config miss", rn)
	}

	if err == nil {
		t.Fatal("Expected error message for miss on FIT config, got nil")
	}
}

func TestLoad(t *testing.T) {
	keyFiles := []string{"key0", "key1"}

	var keys []*openpgp.Entity
	for _, k := range keyFiles {
		b, err := os.ReadFile(filepath.Join("testdata", k))
		if err != nil {
			t.Fatal(err)
		}
		key, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(b)))
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}

	for _, tt := range []struct {
		desc           string
		kernel         string
		initram        string
		keyring        openpgp.KeyRing
		want           error
		wantKernelFill string
		wantInitFill   string
	}{
		{
			desc:           "Successful kernel0/init0 read with key0",
			keyring:        openpgp.EntityList{keys[0]},
			kernel:         "kernel@0",
			initram:        "ramdisk@0",
			want:           nil,
			wantKernelFill: kernel0Fill,
			wantInitFill:   initram0Fill,
		},
		{
			desc:           "Successful unsigned kernel1/init",
			keyring:        nil,
			kernel:         "kernel@1",
			initram:        "ramdisk@0",
			want:           nil,
			wantKernelFill: kernel1Fill,
			wantInitFill:   initram0Fill,
		},
		{
			desc:           "Successful unsigned kernel1",
			keyring:        nil,
			kernel:         "kernel@1",
			initram:        "",
			want:           nil,
			wantKernelFill: kernel1Fill,
			wantInitFill:   "",
		},
		{
			desc:           "bad kernel0 good init0 read with key1",
			keyring:        openpgp.EntityList{keys[1]},
			kernel:         "kernel@0",
			initram:        "ramdisk@0",
			want:           vfile.ErrUnsigned{},
			wantKernelFill: "",
			wantInitFill:   "",
		},
		{
			desc:           "missing kernel",
			keyring:        nil,
			kernel:         "",
			initram:        "ramdisk@0",
			want:           fmt.Errorf(""),
			wantKernelFill: "",
			wantInitFill:   "",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			i, err := New("testdata/fitimage.itb")
			if err != nil {
				t.Fatal(err)
			}

			i.Kernel, i.InitRAMFS, i.KeyRing = tt.kernel, tt.initram, tt.keyring

			defer func(old func(i *boot.LinuxImage, opts ...boot.LoadOption) error) { loadImage = old }(loadImage)

			loadImage = func(i *boot.LinuxImage, opts ...boot.LoadOption) error {
				if i == nil {
					t.Errorf("Load() of kernel:%s, init:%s, keys: %v - passed nil to loadImage", tt.kernel, tt.initram, tt.keyring)
					return nil
				}

				if i.Kernel == nil && tt.wantKernelFill != "" {
					t.Errorf("loadImage gave nil kernel: want pattern '%s'", tt.wantKernelFill)
				}
				if i.Kernel != nil {
					compareImage(t, "kernel", tt.wantKernelFill, i.Kernel)
				}

				if i.Initrd == nil && tt.wantInitFill != "" {
					t.Errorf("loadImage gave nil initram: want pattern '%s'", tt.wantInitFill)
				}
				if i.Initrd != nil {
					compareImage(t, "initram", tt.wantInitFill, i.Initrd)
				}
				return nil
			}

			gotErr := i.Load()

			// Shallow check to verify the correct error
			if (tt.want == nil && gotErr != nil) || reflect.TypeOf(gotErr) != reflect.TypeOf(tt.want) {
				t.Errorf("Load() of kernel:%s, init:%s, keys: %v - got %T, want %T", tt.kernel, tt.initram, tt.keyring, gotErr, tt.want)
			}
		})
	}
}

func compareImage(t *testing.T, name string, wantFill string, image io.ReaderAt) {
	// Make sure we can only read imageSize by expecting an EOF instead of
	// the imageSize+1 byte.
	gotImage := make([]byte, imageSize+1)
	gotSize, err := image.ReadAt(gotImage, 0)
	if !errors.Is(err, io.EOF) {
		t.Errorf("failed reading %s passed to loadImage: ReadAt(%d)/expected %d: got %v, want %v", name, imageSize+1, imageSize, err, io.EOF)
	}
	if gotSize != imageSize {
		t.Errorf("failed reading %s passed to loadImage: bytes got %v bytes, want %v of pattern %s.\nReturned image:\t%v", name, gotSize, imageSize, wantFill, gotImage)
	} else {
		fill := []byte(wantFill)
		for i := 0; i < gotSize; i++ {
			if gotImage[i] != fill[i%len(fill)] {
				t.Errorf("loadImage gave %s: %v, want pattern: \"%v...\". Mismatch at index %d: want %d, got %d", name, gotImage, wantFill, i, gotImage[i], fill[i%len(fill)])
			}
		}
	}
}

func TestReadSignedImage(t *testing.T) {
	keyFiles := []string{"key0", "key1"}

	var keys []*openpgp.Entity
	for _, k := range keyFiles {
		b, err := os.ReadFile(filepath.Join("testdata", k))
		if err != nil {
			t.Fatal(err)
		}
		key, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(b)))
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key)
	}

	for _, tt := range []struct {
		desc             string
		keyring          openpgp.KeyRing
		image            string
		want             error
		isSignatureValid bool
		wantContentFill  string
	}{
		{
			desc:             "Successful kernel read with key0",
			keyring:          openpgp.EntityList{keys[0]},
			image:            "kernel@0",
			want:             nil,
			wantContentFill:  kernel0Fill,
			isSignatureValid: true,
		},
		{
			desc:             "Successful initram read with key0",
			keyring:          openpgp.EntityList{keys[0]},
			image:            "ramdisk@0",
			want:             nil,
			wantContentFill:  initram0Fill,
			isSignatureValid: true,
		},
		{
			desc:             "Successful initram read with key1",
			keyring:          openpgp.EntityList{keys[1]},
			image:            "ramdisk@0",
			want:             nil,
			wantContentFill:  initram0Fill,
			isSignatureValid: true,
		},
		{
			desc:             "Read unsigned kernel1",
			keyring:          openpgp.EntityList{keys[0], keys[1]},
			image:            "kernel@1",
			want:             vfile.ErrUnsigned{},
			wantContentFill:  kernel1Fill,
			isSignatureValid: false,
		},
		{
			desc:             "Read signed kernel0 with wrong key",
			keyring:          openpgp.EntityList{keys[1]},
			image:            "kernel@0",
			want:             vfile.ErrUnsigned{},
			wantContentFill:  kernel0Fill,
			isSignatureValid: false,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			i, err := New("testdata/fitimage.itb")
			if err != nil {
				t.Fatal(err)
			}

			b, gotErr := i.ReadSignedImage(tt.image, tt.keyring)

			// Shallow check to verify we're claiming it's signed or unsigned
			if (tt.want == nil && gotErr != nil) || reflect.TypeOf(gotErr) != reflect.TypeOf(tt.want) {
				t.Errorf("ReadSignedImage(%s, %v) = %T, want %T", tt.image, tt.keyring, gotErr, tt.want)
			}

			if isSignatureValid := (gotErr == nil); isSignatureValid != tt.isSignatureValid {
				t.Errorf("isSignatureValid(%v) = %v, want %v", gotErr, isSignatureValid, tt.isSignatureValid)
			}

			if b != nil {
				compareImage(t, tt.image, tt.wantContentFill, b)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	f, err := os.Open("testdata/fitimage.itb")
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := ParseConfig(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(imgs) != fbcCnt {
		t.Fatalf("Expected 2 images from ParseConfig, got %x", len(imgs))
	}

	cs := [fbcCnt]string{"conf@1", "conf_bz@1"}
	for c, i := range imgs {
		t.Logf("config name: %s", i.ConfigOverride)
		t.Logf("kernel name: %s", i.Kernel)
		t.Logf("ramdisk name: %s", i.InitRAMFS)
		if i.ConfigOverride != cs[c] {
			t.Errorf("Expected config %s, got %s", cs[c], i.ConfigOverride)
		}
		if i.Kernel != "kernel@0" {
			t.Errorf("Expected kernel %s, got %s", "kernel@0", i.Kernel)
		}
		if i.InitRAMFS != "ramdisk@0" {
			t.Errorf("Expected ramdisk %s, got %s", "ramdisk@0", i.InitRAMFS)
		}
	}
}

func TestParseConfigMiss(t *testing.T) {
	f, err := os.Open("testdata/fitimage.its")
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := ParseConfig(f)

	if imgs != nil {
		t.Errorf("Expected nil on bad image parse, got %#v", imgs)
	}

	if err == nil {
		t.Fatal("Expected error on failed ParseConfig, got nil")
	}
}

func TestLabel(t *testing.T) {
	n, kn, rn := "conf_bz@1", "kernel@0", "ramdisk@0"
	img := &Image{name: n, Kernel: kn, InitRAMFS: rn}
	l := img.Label()
	if !strings.Contains(l, n) {
		t.Fatalf("Expected Image label to contain name %s, got %s", n, l)
	}
}

func TestRank(t *testing.T) {
	testRank := 2
	img := &Image{BootRank: testRank}
	l := img.Rank()
	if l != testRank {
		t.Fatalf("Expected Image rank %d, got %d", testRank, l)
	}
}
