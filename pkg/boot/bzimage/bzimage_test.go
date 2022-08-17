// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

import (
	"bytes"
	"encoding/binary"
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
	for _, tc := range testImages {
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
		}, {
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
			// Copy the image to ensure that the test does not change the original image.
			newImage := make([]byte, len(baseImage))
			copy(newImage, baseImage)

			// Write the desired version, Little-Endian style, into the image.
			var b bytes.Buffer // satisfies the io.Writer interface used by binary.Write.
			if err := binary.Write(&b, binary.LittleEndian, tc.version); err != nil {
				t.Fatalf("failed to convert version to LittleEndian: %v", err)
			}
			copy(newImage[0x0206:], b.Bytes())

			// Try to unmarshal the image with the modified version.
			if gotErr := ((&BzImage{}).UnmarshalBinary(newImage) != nil); gotErr != tc.wantErr {
				t.Fatalf("got error: %v, expected error: %t", gotErr, tc.wantErr)
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

	_, err = b.MarshalBinary()
	if err == nil {
		t.Logf("Overflow test, want %v, got nil", "Marshal: compressed KernelCode too big: was 986532, now 1422388")
		t.Fatal(err)
	}

	b.KernelCode = k[:len(k)-len(k)/2]

	_, err = b.MarshalBinary()
	if err != nil {
		t.Logf("shrink test, want nil, got %v", err)
		t.Fatal(err)
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
