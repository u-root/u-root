// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"reflect"
)

// gencsum generates a uint8 checksum of a []uint8
func gencsum(b []uint8) uint8 {
	var csum uint8
	for _, bb := range b {
		csum += bb
	}
	Debug("csum %#x %#x across %d bytes", csum, ^csum, len(b))
	return ^csum + 1
}

// HeapTable is for ACPI tables that have a heap, i.e. the strings
// are not subtables, as in most ACPI, but are contained in an area
// at the end of the tables, after the other table elements. So far,
// we only know of one such table, the IBFT.
type HeapTable struct {
	Head *bytes.Buffer
	Heap *bytes.Buffer
}

// Marshal marshals basic types into HeapTable
func (h *HeapTable) Marshal(i interface{}) error {
	switch s := i.(type) {
	case sockaddr:
		Debug("addr")
		a, err := net.ResolveTCPAddr("tcp", string(s))
		if err != nil {
			return fmt.Errorf("addr %s: %v", s, err)
		}
		w(h.Head, a.IP.To16(), uint16(a.Port))
	case ipaddr:
		a, err := net.ResolveIPAddr("ip", string(s))
		if err != nil {
			return fmt.Errorf("addr %s: %v", s, err)
		}
		w(h.Head, a.IP.To16())
		Debug("net")
	case flag:
	case mac:
		hw, err := net.ParseMAC(string(s))
		if err != nil {
			return err
		}
		if len(hw) != 6 {
			return fmt.Errorf("%q is not an ethernet MAC", s)
		}
		w(h.Head, hw)
		Debug("mac")
	case bdf:
		if err := uw(h.Head, string(s), 16); err != nil {
			return err
		}
		Debug("bdf")
	case u8:
		if err := uw(h.Head, string(s), 8); err != nil {
			return err
		}

	case u16:
		if err := uw(h.Head, string(s), 16); err != nil {
			return err
		}

	case u64:
		if err := uw(h.Head, string(s), 64); err != nil {
			return err
		}
	case sheap:
		w(h.Head, uint16(len(s)), uint16(h.Heap.Len()))
		Debug("Write %q to heap", string(s))
		w(h.Heap, []byte(s))
	default:
		return fmt.Errorf("Don't know what to do with %T", s)
	}
	return nil
}

// Marshal marshals an ACPI Header into a []byte.
func (h *Header) Marshal() ([]byte, error) {
	nt := reflect.TypeOf(h).Elem()
	nv := reflect.ValueOf(h).Elem()
	var b = &bytes.Buffer{}
	for i := 0; i < nt.NumField(); i++ {
		f := nt.Field(i)
		ft := f.Type
		fv := nv.Field(i)

		Debug("Header Marshal Field %d: %d ml %v %T (%v, %v)", i, b.Len(), f, f, ft, fv)
		var err error
		switch s := fv.Interface().(type) {

		case u8:
			if err := uw(b, string(s), 8); err != nil {
				return nil, err
			}

		case u16:
			if err := uw(b, string(s), 16); err != nil {
				return nil, err
			}

		case u32:
			if err := uw(b, string(s), 32); err != nil {
				return nil, err
			}
		case u64:
			if err := uw(b, string(s), 64); err != nil {
				return nil, err
			}
		case sig:
			if err := binary.Write(b, binary.LittleEndian, []byte(s)); err != nil {
				return nil, err
			}
		case oem:
			if err := binary.Write(b, binary.LittleEndian, []byte(s)); err != nil {
				return nil, err
			}
		case tableid:
			if err := binary.Write(b, binary.LittleEndian, []byte(s)); err != nil {
				return nil, err
			}
		case uint32, uint8, uint16, uint64:
			err = binary.Write(b, binary.LittleEndian, s)

		default:
			log.Panicf("Don't know what to do with %T", s)
			return nil, fmt.Errorf("Don't know what to do with %T", s)
		}
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

// ShowTable converts a Table into string.
func ShowTable(t Tabler) string {
	return fmt.Sprintf("%s %d %d %#02x %s %s %#08x %#08x %#08x",
		t.Sig(),
		t.Len(),
		t.Revision(),
		t.CheckSum(),
		t.OEMID(),
		t.OEMTableID(),
		t.OEMRevision(),
		t.CreatorID(),
		t.CreatorRevision())

}
