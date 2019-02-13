// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

type (
	// marshalers marshal ACPI tables into a head and a heap.
	marshaler func(head, heap *bytes.Buffer, i interface{}) error
)

const (
	// CSUMOffset is the offset of the single byte checksum in *most* ACPI tables
	CSUMOffset = 9
)

var (
	// Debug implements fmt.Sprintf and can be used for debug printing
	Debug      = func(string, ...interface{}) {}
	marshalers = map[reflect.Kind]marshaler{}
)

// Flags takes 0 or more flags and produces a uint32 value.
// Each argument represents one bit position, with the first flag
// being bit 0. For each flag with a value of "1", the bit for that
// flag will be set. For each flag with a value of "0", that bit in
// the flag will be cleared.
// This is mainly for consistency but it would allow
// us in future to pass in an initial value and set or clear bits in it
// depending on flags.
// Current allowed values are "0" and "1". In future, if the flags are
// not contiguous, we can allow "ignore" in future to ignore a flag.
func flags(s ...flag) (uint8, error) {
	var i, bit uint8
	for _, f := range s {
		switch f {
		case "1":
			i |= 1 << bit
		case "0":
			i &= ^(1 << bit)
		default:
			return 0, fmt.Errorf("%s is not a valid value: only 0 or 1 are valied", f)
		}
		bit++
	}
	return i, nil
}

// w writes 0 or more values to a bytes.Buffer, in LittleEndian order,
func w(b *bytes.Buffer, val ...interface{}) {
	for _, v := range val {
		binary.Write(b, binary.LittleEndian, v)
		Debug("\t %T %v b is %d bytes", v, v, b.Len())
	}
	Debug("w: done: b is %d bytes", b.Len())
}

// uw writes strings as unsigned words to a bytes.Buffer.
// Currently it only supports 16, 32, and 64 bit writes.
func uw(b *bytes.Buffer, s string, bits int) error {
	// convenience case: if they don't set it, it comes in as "",
	// take that to mean 0.
	var v uint64
	if s != "" {
		var err error
		if v, err = strconv.ParseUint(string(s), 0, bits); err != nil {
			return err
		}
	}
	switch bits {
	case 8:
		w(b, uint8(v))
	case 16:
		w(b, uint16(v))
	case 32:
		w(b, uint32(v))
	case 64:
		w(b, uint64(v))
	}
	return nil

}

// Marshal marshals support ACPI tables into a byte slice.
func Marshal(i interface{}) ([]byte, error) {
	var head, heap bytes.Buffer
	Debug("Marshall %T", i)
	m, ok := marshalers[reflect.TypeOf(i).Kind()]
	if !ok {
		return nil, fmt.Errorf("No marshaler for %T, have %v", i, marshalers)
	}
	Debug("Marshaler is %v", m)
	// We pass in both a head and a heap. For most ACPI tables,
	// only the head is written. For some tables, the heap is used
	// as well. The top level handler in marshal is required to return
	// with the heap reset and the head containing any tables.
	if err := m(&head, &heap, i); err != nil {
		return nil, err
	}

	if heap.Len() != 0 {
		return nil, fmt.Errorf("%T top level returned with non-empty heap", i)
	}
	s := head.Bytes()
	binary.LittleEndian.PutUint32(s[4:], uint32(len(s)))
	c := gencsum(head.Bytes())
	Debug("CSUM is %#x", c)
	s[9] = c

	return s, nil
}
