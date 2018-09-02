// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uefi contents data types for the components found in UEFI and an
// Parse function for reading an image.
package uefi

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

// ROMAttributes is used to hold global variables that apply across the whole image.
// We have to do this to avoid passing too many things down each time.
type ROMAttributes struct {
	ErasePolarity byte // Either 0xFF or 0
}

// Attributes holds the global attributes
var Attributes ROMAttributes

// Firmware is an interface to describe generic firmware types. When the
// firmware is parsed, all the Firmware objects are laid out in a tree (similar
// to an AST). This interface represents one node in said tree. The
// implementations (e.g. Flash image, or FirmwareVolume) must implement this
// interface.
type Firmware interface {
	Validate() []error

	Buf() []byte
	SetBuf(buf []byte)

	// Apply a visitor to the Firmware.
	Apply(v Visitor) error

	// Apply a visitor to all the direct children of the Firmware
	// (excluding the Firmware itself).
	ApplyChildren(v Visitor) error
}

// TypedFirmware includes the Firmware interface's type when exporting it to
// JSON. The type is required when unmarshalling.
type TypedFirmware struct {
	Type  string
	Value Firmware
}

// MakeTyped takes a Firmware interface and makes a (type, value) pair.
func MakeTyped(f Firmware) *TypedFirmware {
	return &TypedFirmware{
		Type:  reflect.TypeOf(f).String(),
		Value: f,
	}
}

// UnmarshalJSON unmarshals a TypedFirmware struct and correctly deduces the
// type of the interface.
func (f *TypedFirmware) UnmarshalJSON(b []byte) error {
	var getType struct {
		Type  string
		Value json.RawMessage
	}
	if err := json.Unmarshal(b, &getType); err != nil {
		return err
	}
	factory, ok := firmwareTypes[getType.Type]
	if !ok {
		return fmt.Errorf("unknown TypedFirmware type '%s', unable to unmarshal", getType.Type)
	}
	f.Type = getType.Type
	f.Value = factory()
	return json.Unmarshal(getType.Value, &f.Value)
}

// This should never be exposed, it is only used for marshalling different types to json.
type marshalFirmware struct {
	FType           string
	FirmwareElement json.RawMessage
}

var firmwareTypes = map[string]func() Firmware{
	"*uefi.BIOSRegion":      func() Firmware { return &BIOSRegion{} },
	"*uefi.BIOSPadding":     func() Firmware { return &BIOSPadding{} },
	"*uefi.File":            func() Firmware { return &File{} },
	"*uefi.FirmwareVolume":  func() Firmware { return &FirmwareVolume{} },
	"*uefi.FlashDescriptor": func() Firmware { return &FlashDescriptor{} },
	"*uefi.FlashImage":      func() Firmware { return &FlashImage{} },
	"*uefi.GBERegion":       func() Firmware { return &GBERegion{} },
	"*uefi.MERegion":        func() Firmware { return &MERegion{} },
	"*uefi.PDRegion":        func() Firmware { return &PDRegion{} },
	"*uefi.Section":         func() Firmware { return &Section{} },
}

// MarshalFirmware marshals the firmware element to JSON, including the type information at the top.
func MarshalFirmware(f Firmware) ([]byte, error) {
	b, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		return nil, err
	}

	m := marshalFirmware{FType: reflect.TypeOf(f).String(), FirmwareElement: json.RawMessage(b)}
	return json.MarshalIndent(m, "", "    ")
}

// UnmarshalFirmware unmarshals the firmware element from JSON, using the type information at the top.
func UnmarshalFirmware(b []byte) (Firmware, error) {
	var m marshalFirmware
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	factory, ok := firmwareTypes[m.FType]
	if !ok {
		return nil, fmt.Errorf("unknown Firmware type '%s', unable to unmarshal", m.FType)
	}
	f := factory()
	err := json.Unmarshal(m.FirmwareElement, &f)
	return f, err
}

// Parse exposes a high-level parser for generic firmware types. It does not
// implement any parser itself, but it calls known parsers that implement the
// Firmware interface.
func Parse(buf []byte) (Firmware, error) {
	if _, err := FindSignature(buf); err == nil {
		// Intel rom.
		return NewFlashImage(buf)
	}
	// Non intel image such as edk2's OVMF
	// We don't know how to parse this header, so treat it as a large BIOSRegion
	return NewBIOSRegion(buf, nil)
}

// ExtractBinary simply dumps the binary to a specified directory and filename.
// It creates the directory if it doesn't already exist, and dumps the buffer to it.
// It returns the filepath of the binary, and an error if it exists.
// This is meant as a helper function for other Extract functions.
func ExtractBinary(buf []byte, dirPath string, filename string) (string, error) {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	// Dump the binary.
	fp := filepath.Join(dirPath, filename)
	if err := ioutil.WriteFile(fp, buf, 0666); err != nil {
		// Make sure we return "" since we don't want an invalid path to be serialized out.
		return "", err
	}
	return fp, nil
}

// Checksum8 does a 8 bit checksum of the slice passed in.
func Checksum8(buf []byte) uint8 {
	var sum uint8
	for _, val := range buf {
		sum += val
	}
	return sum
}

// Checksum16 does a 16 bit checksum of the byte slice passed in.
func Checksum16(buf []byte) (uint16, error) {
	r := bytes.NewReader(buf)
	buflen := len(buf)
	if buflen%2 != 0 {
		return 0, fmt.Errorf("byte slice does not have even length, not able to do 16 bit checksum. Length was %v",
			buflen)
	}
	var temp, sum uint16
	for i := 0; i < buflen; i += 2 {
		if err := binary.Read(r, binary.LittleEndian, &temp); err != nil {
			return 0, err
		}
		sum += temp
	}
	return sum, nil
}

// Read3Size reads a 3-byte size and returns it as a uint64
func Read3Size(size [3]uint8) uint64 {
	return uint64(size[2])<<16 |
		uint64(size[1])<<8 | uint64(size[0])
}

// Write3Size writes a size into a 3-byte array
func Write3Size(size uint64) [3]uint8 {
	if size >= 0xFFFFFF {
		return [3]uint8{0xFF, 0xFF, 0xFF}
	}
	b := [3]uint8{uint8(size), uint8(size >> 8), uint8(size >> 16)}
	return b
}

// Align aligns an address
func Align(val uint64, base uint64) uint64 {
	return (val + base - 1) & ^(base - 1)
}

// Align4 aligns an address to 4 bytes
func Align4(val uint64) uint64 {
	return Align(val, 4)
}

// Align8 aligns an address to 8 bytes
func Align8(val uint64) uint64 {
	return Align(val, 8)
}

// Erase sets the buffer to be ErasePolarity
func Erase(buf []byte, polarity byte) {
	for j, blen := 0, len(buf); j < blen; j++ {
		buf[j] = Attributes.ErasePolarity
	}
}
