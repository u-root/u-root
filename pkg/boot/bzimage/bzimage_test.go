// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

import (
	"fmt"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

type testImage struct {
	name string
	path string
}

var testImages = []testImage{
	{
		name: "basic bzImage",
		path: "testdata/bzImage",
	},
	{
		name: "a little larger bzImage, 64k random generated image",
		path: "testdata/bzimage-64kurandominitramfs",
	},
}

var badmagic = []byte("hi there")

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading %q: %v", path, err)
	}
	return data
}

func TestUnmarshal(t *testing.T) {
	Debug = t.Logf

	compressedTests := []testImage{
		// These test files have been created using .circleci/images/test-image-amd6/config_linux5.10_x86_64.txt
		{name: "bzip2", path: "testdata/bzImage-linux5.10-x86_64-bzip2"},
		{name: "signed-debian", path: "testdata/bzImage-debian-signed-linux5.10.0-6-amd64_5.10.28-1_amd64"},
		{name: "signed-rocky", path: "testdata/bzImage-rockylinux9"},
		{name: "gzip", path: "testdata/bzImage-linux5.10-x86_64-gzip"},
		{name: "xz", path: "testdata/bzImage-linux5.10-x86_64-xz"},
		{name: "lz4", path: "testdata/bzImage-linux5.10-x86_64-lz4"},
		{name: "lzma", path: "testdata/bzImage-linux5.10-x86_64-lzma"},
		// This test does not pass because the CircleCI environment does not include the `lzop` command.
		// TODO: Fix the CircleCI environment or (preferably) find a Go package which provides this functionality.
		//		{name: "lzo", path: "testdata/bzImage-linux5.10-x86_64-lzo"},
		{name: "zstd", path: "testdata/bzImage-linux5.10-x86_64-zstd"},
	}

	for _, tc := range append(testImages, compressedTests...) {
		t.Run(tc.name, func(t *testing.T) {
			image := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(image); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSupportedVersions(t *testing.T) {
	Debug = t.Logf

	tests := []struct {
		version uint16
		wantErr bool
	}{
		{
			version: 0x0207,
			wantErr: true,
		},
		{
			version: 0x0208,
			wantErr: false,
		},
		{
			version: 0x0209,
			wantErr: false,
		},
	}

	baseImage := mustReadFile(t, "testdata/bzImage")

	// Ensure that the base image unmarshals successfully.
	if err := (&BzImage{}).UnmarshalBinary(baseImage); err != nil {
		t.Fatalf("unable to unmarshal image: %v", err)
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("0x%04x", tc.version), func(t *testing.T) {
			// Unmarshal the base image.
			var bzImage BzImage
			if err := bzImage.UnmarshalBinary(baseImage); err != nil {
				t.Fatalf("failed to unmarshal base image: %v", err)
			}

			bzImage.Header.Protocolversion = tc.version

			// Marshal the image with the test version.
			modifiedImage, err := bzImage.MarshalBinary()
			if err != nil {
				t.Fatalf("failed to marshal image with the new version: %v", err)
			}

			// Try to unmarshal the image with the test version.
			err = (&BzImage{}).UnmarshalBinary(modifiedImage)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("got error: %v, expected error: %t", err, tc.wantErr)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	Debug = t.Logf
	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			image := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(image); err != nil {
				t.Fatal(err)
			}
			t.Logf("b header is %s", b.Header.String())
			image, err := b.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			// now unmarshall back into ourselves.
			// We can't perfectly recreate the image the kernel built,
			// but we need to know we are stable.
			if err := b.UnmarshalBinary(image); err != nil {
				t.Fatal(err)
			}
			d, err := b.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			var n BzImage
			if err := n.UnmarshalBinary(d); err != nil {
				t.Fatalf("Unmarshalling marshaled image: want nil, got  %v", err)
			}

			t.Logf("DIFF: %v", b.Header.Diff(&n.Header))
			if d := b.Header.Diff(&n.Header); d != "" {
				t.Errorf("Headers differ: %s", d)
			}
			if len(d) != len(image) {
				t.Fatalf("Marshal: want %d as output len, got %d; diff is %s", len(image), len(d), b.Diff(&b))
			}

			if err := Equal(image, d); err != nil {
				t.Logf("Check if images are the same: want nil, got %v", err)
			}

			// Corrupt little bits of thing.
			x := d[0x203]
			d[0x203] = 1
			if err := Equal(image, d); err == nil {
				t.Fatalf("Corrupting marshaled image: got nil, want err")
			}
			d[0x203] = x
			image[0x203] = 1
			if err := Equal(image, d); err == nil {
				t.Fatalf("Corrupting original image: got nil, want err")
			}
			image[0x203] = x
			x = d[0x208]
			d[0x208] = x + 1
			if err := Equal(image, d); err == nil {
				t.Fatalf("Corrupting marshaled header: got nil, want err")
			}
			d[0x208] = x
			d[20000] = d[20000] + 1
			if err := Equal(image, d); err == nil {
				t.Fatalf("Corrupting marshaled kernel: got nil, want err")
			}
		})
	}
}

func TestBadMagic(t *testing.T) {
	var b BzImage
	Debug = t.Logf
	if err := b.UnmarshalBinary(badmagic); err == nil {
		t.Fatal("Want err, got nil")
	}
}

func TestAddInitRAMFS(t *testing.T) {
	t.Logf("TestAddInitRAMFS")
	Debug = t.Logf
	initramfsimage := mustReadFile(t, "testdata/bzimage-64kurandominitramfs")
	var b BzImage
	if err := b.UnmarshalBinary(initramfsimage); err != nil {
		t.Fatal(err)
	}
	if err := b.AddInitRAMFS("testdata/init.cpio"); err != nil {
		t.Fatal(err)
	}
	d, err := b.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	// Ensure that we can still unmarshal the image.
	if err := (&BzImage{}).UnmarshalBinary(d); err != nil {
		t.Fatalf("unable to unmarshal the marshal'd image: %v", err)
	}

	// For testing, you can enable this write, and then:
	// qemu-system-x86_64 -serial stdio -kernel /tmp/x
	// I mainly left this here as a memo.
	if true {
		if err := os.WriteFile("/tmp/x", d, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Make KernelCode much bigger and watch it fail.
	k := b.KernelCode
	b.KernelCode = append(b.KernelCode, k...)
	b.KernelCode = append(b.KernelCode, k...)
	b.KernelCode = append(b.KernelCode, k...)
	b.KernelCode = append(b.KernelCode, k...)

	if _, err = b.MarshalBinary(); err == nil {
		t.Logf("Overflow test, want %v, got nil", "Marshal: compressed KernelCode too big: was 986532, now 1422388")
		t.Fatal(err)
	}

	b.KernelCode = k[:len(k)-len(k)/2]

	if _, err = b.MarshalBinary(); err != nil {
		t.Logf("shrink test, want nil, got %v", err)
		t.Fatal(err)
	}
	// Ensure that we can still unmarshal the image.
	if err := (&BzImage{}).UnmarshalBinary(d); err != nil {
		t.Fatalf("unable to unmarshal the marshal'd image: %v", err)
	}
}

func TestHeaderString(t *testing.T) {
	Debug = t.Logf
	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			initramfsimage := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(initramfsimage); err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", b.Header.String())
		})
	}
}

func TestExtract(t *testing.T) {
	Debug = t.Logf
	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			initramfsimage := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(initramfsimage); err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", b.Header.String())
			// The simplest test: what is extracted should be a valid elf.
			e, err := b.ELF()
			if err != nil {
				t.Fatalf("Extracted bzImage data is an elf: want nil, got %v", err)
			}
			t.Logf("Header: %v", e.FileHeader)
			for i, p := range e.Progs {
				t.Logf("%d: %v", i, *p)
			}
		})
	}
}

func TestELF(t *testing.T) {
	Debug = t.Logf
	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			initramfsimage := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(initramfsimage); err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", b.Header.String())
			e, err := b.ELF()
			if err != nil {
				t.Fatalf("Extract: want nil, got %v", err)
			}
			t.Logf("%v", e.FileHeader)
		})
	}
}

func TestInitRAMFS(t *testing.T) {
	Debug = t.Logf
	cpio.Debug = t.Logf
	for _, tc := range testImages {
		t.Run(tc.name, func(t *testing.T) {
			initramfsimage := mustReadFile(t, tc.path)
			var b BzImage
			if err := b.UnmarshalBinary(initramfsimage); err != nil {
				t.Fatal(err)
			}
			s, e, err := b.InitRAMFS()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Found %d byte initramfs@%d:%d", e-s, s, e)
		})
	}
}
