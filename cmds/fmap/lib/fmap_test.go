// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmap

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

// Flash map is stored in little-endian.
var fmapName = []byte("Fake flash" + strings.Repeat("\x00", 32-10))
var area0Name = []byte("Area Number 1\x00\x00\x00Hello" + strings.Repeat("\x00", 32-21))
var area1Name = []byte("Area Number 2xxxxxxxxxxxxxxxxxxx")
var fakeFlash = bytes.Join([][]byte{
	// Arbitrary data
	bytes.Repeat([]byte{0x53, 0x11, 0x34, 0x22}, 94387),

	// Signature
	[]byte("__FMAP__"),
	// VerMajor, VerMinor
	{1, 0},
	// Base
	{0xef, 0xbe, 0xad, 0xde, 0xbe, 0xba, 0xfe, 0xca},
	// Size
	{0x11, 0x22, 0x33, 0x44},
	// Name (32 bytes)
	fmapName,
	// NAreas
	{0x02, 0x00},

	// Areas[0].Offset
	{0xef, 0xbe, 0xad, 0xde},
	// Areas[0].Size
	{0x11, 0x11, 0x11, 0x11},
	// Areas[0].Name (32 bytes)
	area0Name,
	// Areas[0].Flags
	{0x13, 0x10},

	// Areas[1].Offset
	{0xbe, 0xba, 0xfe, 0xca},
	// Areas[1].Size
	{0x22, 0x22, 0x22, 0x22},
	// Areas[1].Name (32 bytes)
	area1Name,
	// Areas[1].Flags
	{0x00, 0x00},
}, []byte{})

func TestReadFMap(t *testing.T) {
	r := bytes.NewReader(fakeFlash)
	fmap, _, err := Read(r)
	if err != nil {
		t.Fatal(err)
	}
	expected := FMap{
		Header: Header{
			VerMajor: 1,
			VerMinor: 0,
			Base:     0xcafebabedeadbeef,
			Size:     0x44332211,
			NAreas:   2,
		},
		Areas: []Area{
			{
				Offset: 0xdeadbeef,
				Size:   0x11111111,
				Flags:  0x1013,
			}, {
				Offset: 0xcafebabe,
				Size:   0x22222222,
				Flags:  0x0000,
			},
		},
	}
	copy(expected.Signature[:], []byte("__FMAP__"))
	copy(expected.Name.Value[:], fmapName)
	copy(expected.Areas[0].Name.Value[:], area0Name)
	copy(expected.Areas[1].Name.Value[:], area1Name)
	if !reflect.DeepEqual(*fmap, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v", expected, *fmap)
	}
}

func TestReadMetadata(t *testing.T) {
	r := bytes.NewReader(fakeFlash)
	_, metadata, err := Read(r)
	if err != nil {
		t.Fatal(err)
	}
	expected := Metadata{
		Start: 4 * 94387,
	}
	if !reflect.DeepEqual(*metadata, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v", expected, *metadata)
	}
}

func TestFieldNames(t *testing.T) {
	r := bytes.NewReader(fakeFlash)
	fmap, _, err := Read(r)
	if err != nil {
		t.Fatal(err)
	}
	for i, expected := range []string{"STATIC|COMPRESSED|0x1010", "0x0"} {
		got := FlagNames(fmap.Areas[i].Flags)
		if got != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
		}
	}
}

func TestNoSignature(t *testing.T) {
	fakeFlash := bytes.Repeat([]byte{0x53, 0x11, 0x34, 0x22}, 94387)
	r := bytes.NewReader(fakeFlash)
	_, _, err := Read(r)
	expected := "Cannot find fmap signature"
	got := err.Error()
	if expected != got {
		t.Errorf("expected: %s; got: %s", expected, got)
	}
}

func TestTwoSignatures(t *testing.T) {
	fakeFlash := bytes.Repeat(fakeFlash, 2)
	r := bytes.NewReader(fakeFlash)
	_, _, err := Read(r)
	expected := "Found multiple signatures"
	got := err.Error()
	if expected != got {
		t.Errorf("expected: %s; got: %s", expected, got)
	}
}

func TestTruncatedFmap(t *testing.T) {
	r := bytes.NewReader(fakeFlash[:len(fakeFlash)-2])
	_, _, err := Read(r)
	expected := "Unexpected EOF while parsing fmap"
	got := err.Error()
	if expected != got {
		t.Errorf("expected: %s; got: %s", expected, got)
	}
}

func TestReadArea(t *testing.T) {
	fmap := FMap{
		Header: Header{
			NAreas: 3,
		},
		Areas: []Area{
			{
				Offset: 0x0,
				Size:   0x10,
			}, {
				Offset: 0x10,
				Size:   0x20,
			}, {
				Offset: 0x30,
				Size:   0x40,
			},
		},
	}
	fakeFlash := bytes.Repeat([]byte{0x53, 0x11, 0x34, 0x22}, 0x70)
	r := bytes.NewReader(fakeFlash)
	area, err := fmap.ReadArea(r, 1)
	if err != nil {
		t.Fatal(err)
	}
	expected := fakeFlash[0x10:0x30]
	got, err := ioutil.ReadAll(area)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, got) {
		t.Errorf("expected: %v; got: %v", expected, got)
	}
}

func TestChecksum(t *testing.T) {
	fmap := FMap{
		Header: Header{
			NAreas: 3,
		},
		Areas: []Area{
			{
				Offset: 0x00,
				Size:   0x03,
				Flags:  FmapAreaStatic,
			}, {
				Offset: 0x03,
				Size:   0x20,
				Flags:  0x00,
			}, {
				Offset: 0x23,
				Size:   0x04,
				Flags:  FmapAreaStatic | FmapAreaCompressed,
			},
		},
	}
	fakeFlash := bytes.Repeat([]byte("abcd"), 0x70)
	r := bytes.NewReader(fakeFlash)
	checksum, err := fmap.Checksum(r, sha256.New())
	if err != nil {
		t.Fatal(err)
	}
	// $ echo -n abcdabc | sha256sum
	want := "8a50a4422d673f463f8e4141d8c4b68c4f001ba16f83ad77b8a31bde53ee7273"
	got := fmt.Sprintf("%x", checksum)
	if want != got {
		t.Errorf("want: %v; got: %v", want, got)
	}
}
