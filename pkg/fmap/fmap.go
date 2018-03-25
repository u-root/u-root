// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmap

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var signature = []byte("__FMAP__")

// Flags which can be applied to Area.Flags.
const (
	FmapAreaStatic = 1 << iota
	FmapAreaCompressed
	FmapAreaReadOnly
)

// String wraps around byte array to give us more control over how strings are
// serialized.
type String struct {
	Value [32]uint8
}

func (s *String) String() string {
	return strings.TrimRight(string(s.Value[:]), "\x00")
}

// MarshalJSON implements json.Marshaler.
func (s *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	if len(str) > len(s.Value) {
		return fmt.Errorf("String %#v is longer than 32 bytes", str)
	}
	copy(s.Value[:], []byte(str))
	return nil
}

// FMap structure serializable using encoding.Binary.
type FMap struct {
	Header
	Areas []Area
}

// Header describes the flash part.
type Header struct {
	Signature [8]uint8
	VerMajor  uint8
	VerMinor  uint8
	Base      uint64
	Size      uint32
	Name      String
	NAreas    uint16
}

// Area describes each area.
type Area struct {
	Offset uint32
	Size   uint32
	Name   String
	Flags  uint16
}

// Metadata contains additional data not part of the FMap.
type Metadata struct {
	Start uint64
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

func readField(r io.Reader, data interface{}) error {
	// The endianness might depend on your machine or it might not.
	if err := binary.Read(r, binary.LittleEndian, data); err != nil {
		return errors.New("Unexpected EOF while parsing fmap")
	}
	return nil
}

// Read an FMap into the data structure.
func Read(f io.Reader) (*FMap, *Metadata, error) {
	// Read flash into memory.
	// TODO: it is possible to parse fmap without reading entire file into memory
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	// Check for too many fmaps.
	if bytes.Count(data, signature) >= 2 {
		return nil, nil, errors.New("Found multiple signatures")
	}

	// Check for too few fmaps.
	start := bytes.Index(data, signature)
	if start == -1 {
		return nil, nil, errors.New("Cannot find fmap signature")
	}

	// Reader anchored to the start of the fmap
	r := bytes.NewReader(data[start:])

	// Read fields.
	var fmap FMap
	if err := readField(r, &fmap.Header); err != nil {
		return nil, nil, err
	}
	fmap.Areas = make([]Area, fmap.NAreas)
	if err := readField(r, &fmap.Areas); err != nil {
		return nil, nil, err
	}

	// Return useful metadata
	fmapMetadata := Metadata{
		Start: uint64(start),
	}

	return &fmap, &fmapMetadata, nil
}

// Write overwrites the fmap in the flash file.
func Write(f io.WriteSeeker, fmap *FMap, m *Metadata) error {
	if _, err := f.Seek(int64(m.Start), io.SeekStart); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, fmap.Header); err != nil {
		return err
	}
	return binary.Write(f, binary.LittleEndian, fmap.Areas)
}

// ReadArea reads an area from the fmap as a binary stream.
func (f *FMap) ReadArea(r io.ReadSeeker, i int) (io.Reader, error) {
	if i < 0 || int(f.NAreas) <= i {
		return nil, errors.New("Area index out of range")
	}
	if _, err := r.Seek(int64(f.Areas[i].Offset), io.SeekStart); err != nil {
		return nil, err
	}
	return io.LimitReader(r, int64(f.Areas[i].Size)), nil
}

// Checksum performs a hash of the static areas.
func (f *FMap) Checksum(r io.ReadSeeker, h hash.Hash) ([]byte, error) {
	for i, v := range f.Areas {
		if v.Flags&FmapAreaStatic == 0 {
			continue
		}
		areaReader, err := f.ReadArea(r, i)
		if err != nil {
			return nil, err
		}
		_, err = bufio.NewReader(areaReader).WriteTo(h)
		if err != nil {
			return nil, err
		}
	}
	return h.Sum([]byte{}), nil
}
