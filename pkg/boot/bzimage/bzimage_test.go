// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

import (
	"io/ioutil"
	"testing"

	"github.com/u-root/u-root/pkg/cpio"
)

var badmagic = []byte("hi there")

func TestUnmarshal(t *testing.T) {
	Debug = t.Logf
	image, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
	var b BzImage
	if err := b.UnmarshalBinary(image); err != nil {
		t.Fatal(err)
	}
}

func TestMarshal(t *testing.T) {
	Debug = t.Logf
	image, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
	var b BzImage
	if err := b.UnmarshalBinary(image); err != nil {
		t.Fatal(err)
	}
	t.Logf("b header is %s", b.Header.String())
	image, err = b.MarshalBinary()
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
	initramfsimage, err := ioutil.ReadFile("testdata/bzimage-64kurandominitramfs")
	if err != nil {
		t.Fatal(err)
	}
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
		if err := ioutil.WriteFile("/tmp/x", d, 0644); err != nil {
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
	initramfsimage, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
	var b BzImage
	if err := b.UnmarshalBinary(initramfsimage); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b.Header.String())
}
func TestExtract(t *testing.T) {
	Debug = t.Logf
	initramfsimage, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
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
}

func TestELF(t *testing.T) {
	Debug = t.Logf
	initramfsimage, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
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
}

func TestInitRAMFS(t *testing.T) {
	Debug = t.Logf
	cpio.Debug = t.Logf
	for _, bz := range []string{"testdata/bzImage", "testdata/bzimage-64kurandominitramfs"} {
		initramfsimage, err := ioutil.ReadFile(bz)
		if err != nil {
			t.Fatal(err)
		}
		var b BzImage
		if err := b.UnmarshalBinary(initramfsimage); err != nil {
			t.Fatal(err)
		}
		s, e, err := b.InitRAMFS()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Found %d byte initramfs@%d:%d", e-s, s, e)
	}

}
