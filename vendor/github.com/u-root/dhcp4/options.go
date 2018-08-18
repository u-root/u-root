// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp4

import (
	"encoding"
	"io"
	"math"
	"sort"

	"github.com/u-root/dhcp4/internal/buffer"
)

// Options is a map of OptionCode keys with a slice of byte values.
//
// Its methods can be used to easily check for additional information from a
// packet. Get should be used to access data from Options.
type Options map[OptionCode][]byte

// Add adds a new OptionCode key and BinaryMarshaler's bytes to the Options
// map.
func (o Options) Add(key OptionCode, value encoding.BinaryMarshaler) error {
	if value == nil {
		o.AddRaw(key, []byte{})
		return nil
	}

	b, err := value.MarshalBinary()
	if err != nil {
		return err
	}

	o.AddRaw(key, b)
	return nil
}

// AddRaw adds a new OptionCode key and raw value byte slice to the Options
// map.
func (o Options) AddRaw(key OptionCode, value []byte) {
	o[key] = append(o[key], value...)
}

// Get attempts to retrieve the value specified by an OptionCode key.
//
// If a value is found, get returns a non-nil byte slice. If it is not found,
// Get returns nil.
func (o Options) Get(key OptionCode) []byte {
	// Check for value by key.
	v, ok := o[key]
	if !ok {
		return nil
	}

	// Some options can actually have zero length (OptionRapidCommit), so
	// just return an empty byte slice if this is the case.
	if len(v) == 0 {
		return []byte{}
	}
	return v
}

// Unmarshal fills opts with option codes and corresponding values from an
// input byte slice.
//
// It is used with various different types to enable parsing of both top-level
// options. If options data is malformed, it returns ErrInvalidOptions or
// io.ErrUnexpectedEOF.
func (o *Options) Unmarshal(buf *buffer.Buffer) error {
	*o = make(Options)

	var end bool
	for buf.Len() >= 1 {
		// 1 byte: option code
		// 1 byte: option length n
		// n bytes: data
		code := OptionCode(buf.Read8())

		if code == Pad {
			continue
		} else if code == End {
			end = true
			break
		}
		if !buf.Has(1) {
			return io.ErrUnexpectedEOF
		}

		length := int(buf.Read8())
		if length == 0 {
			continue
		}

		if !buf.Has(length) {
			return io.ErrUnexpectedEOF
		}

		// N bytes: option data
		data := buf.Consume(length)
		if data == nil {
			return io.ErrUnexpectedEOF
		}
		data = data[:length:length]

		// RFC 3396: Just concatenate the data if the option code was
		// specified multiple times.
		o.AddRaw(code, data)
	}

	if !end {
		return io.ErrUnexpectedEOF
	}

	// Any bytes left must be padding.
	for buf.Len() >= 1 {
		if OptionCode(buf.Read8()) != Pad {
			return ErrInvalidOptions
		}
	}
	return nil
}

// Marshal writes options into the provided Buffer sorted by option codes.
func (o Options) Marshal(b *buffer.Buffer) {
	for _, c := range o.sortedKeys() {
		code := OptionCode(c)
		data := o[code]

		// RFC 3396: If more than 256 bytes of data are given, the
		// option is simply listed multiple times.
		for len(data) > 0 {
			// 1 byte: option code
			b.Write8(uint8(code))

			// Some DHCPv4 options have fixed length and do not put
			// length on the wire.
			if code == End || code == Pad {
				continue
			}

			n := len(data)
			if n > math.MaxUint8 {
				n = math.MaxUint8
			}

			// 1 byte: option length
			b.Write8(uint8(n))

			// N bytes: option data
			b.WriteBytes(data[:n])
			data = data[n:]
		}
	}

	// If "End" option is not in map, marshal it manually.
	if _, ok := o[End]; !ok {
		b.Write8(uint8(End))
	}
}

// sortedKeys returns an ordered slice of option keys from the Options map, for
// use in serializing options to binary.
func (o Options) sortedKeys() []int {
	// Send all values for a given key
	var codes []int
	for k := range o {
		codes = append(codes, int(k))
	}

	sort.Sort(sort.IntSlice(codes))
	return codes
}
