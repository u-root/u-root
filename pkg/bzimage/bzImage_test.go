// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

import (
	"io/ioutil"
	"testing"
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
	d, err := b.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != len(image) {
		t.Fatalf("Marshal: want %d as output len, got %d", len(image), len(d))
	}
	if err := Equal(image, d); err != nil {
		t.Fatalf("Check if images are the same: want nil, got %v", err)
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
	Debug = t.Logf
	initramfsimage, err := ioutil.ReadFile("testdata/bzImage")
	if err != nil {
		t.Fatal(err)
	}
	var b BzImage
	if err := b.UnmarshalBinary(initramfsimage); err != nil {
		t.Fatal(err)
	}
	b.AddInitRAMFS("testdata/init.cpio")
	d, err := b.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	// For testing, you can enable this write, and then:
	// qemu-system-x86_64 -serial stdio -kernel /tmp/x
	// I mainly left this here as a memo.
	if false {
		if err := ioutil.WriteFile("/tmp/x", d, 0644); err != nil {
			t.Fatal(err)
		}
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
