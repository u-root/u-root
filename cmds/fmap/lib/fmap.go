// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

var signature = []byte("__FMAP__")

const (
	FmapAreaStatic = 1 << iota
	FmapAreaCompressed
	FmapAreaReadOnly
)

type fMapInternal struct {
	Signature [8]uint8
	VerMajor  uint8
	VerMinor  uint8
	Base      uint64
	Size      uint32
	Name      [32]uint8
}

type FMap struct {
	fMapInternal
	Start uint64
	Areas []FMapArea
}

type FMapArea struct {
	Offset uint32
	Size   uint32
	Name   [32]uint8
	Flags  uint16
}

// FlagNames returns human readable representation of the flags.
func FlagNames(flags uint16) string {
	names := []string{}
	m := []struct {
		val  uint16
		name string
	}{
		{FmapAreaStatic, "STATIC"},
		{FmapAreaCompressed, "COMPRESSED"},
		{FmapAreaReadOnly, "READ_ONLY"},
	}
	for _, v := range m {
		if v.val&flags != 0 {
			names = append(names, v.name)
			flags -= v.val
		}
	}
	// Write a hex value for unknown flags.
	if flags != 0 || len(names) == 0 {
		names = append(names, fmt.Sprintf("%#x", flags))
	}
	return strings.Join(names, "|")
}

func readField(r io.Reader, data interface{}) {
	// The endianness might depend on your machine or it might not.
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		log.Fatal("Unexpected EOF while parsing fmap")
	}
}

// Read an FMap into the data structure.
func ReadFMap(f io.Reader) *FMap {
	// Read flash into memory.
	// TODO: it is possible to parse fmap without reading entire file into memory
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	// Check for too many fmaps.
	if bytes.Count(data, signature) >= 2 {
		log.Print("Warning: Found multiple signatures, using first")
	}

	// Check for too few fmaps.
	start := bytes.Index(data, signature)
	if start == -1 {
		log.Fatal("Cannot find fmap signature")
	}

	// Reader anchored to the start of the fmap
	r := bytes.NewReader(data[start:])

	// Read fields.
	fmap := FMap{Start: uint64(start)}
	readField(r, &fmap.fMapInternal)
	var nAreas uint16
	readField(r, &nAreas)
	fmap.Areas = make([]FMapArea, nAreas)
	readField(r, &fmap.Areas)

	return &fmap
}
