// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
)

// SDT represents either an RSDT or XSDT. It has a Generic header
// and Tables, which are pointers. In the RSDT they are 32 bits;
// in the XSDT, 64. We unmarshal to 64 bits, and when we marshal,
// we use the signature to determine whether the table is 32 or 64.
type SDT struct {
	Generic
	Tables []int64
	Base   int64
}

func init() {
	addUnMarshaler("RSDT", unmarshalSDT)
	addUnMarshaler("XSDT", unmarshalSDT)
}

func unmarshalSDT(t Tabler) (Tabler, error) {
	s := &SDT{
		Generic: Generic{
			Header: *GetHeader(t),
			data:   t.AllData(),
		},
	}

	sig := s.Sig()
	if sig != "RSDT" && sig != "XSDT" {
		return nil, fmt.Errorf("%v is not RSDT or XSDT", sig)
	}

	// Now the fun. In 1999, 64-bit micros had been out for about 10 years.
	// Intel had announced the ia64 years earlier. In 2000 the ACPI committee
	// chose 32-bit pointers anyway, then had to backfill a bunch of table
	// types to do 64 bits shortly thereafter (i.e. v2). Geez.
	esize := 4
	if sig == "XSDT" {
		esize = 8
	}
	d := t.TableData()

	for i := 0; i < len(d); i += esize {
		val := int64(0)
		if sig == "XSDT" {
			val = int64(binary.LittleEndian.Uint64(d[i : i+8]))
		} else {
			val = int64(binary.LittleEndian.Uint32(d[i : i+4]))
		}
		s.Tables = append(s.Tables, val)
	}
	return s, nil
}

// Marshal marshals an [RX]SDT. If it has tables, it marshals them too.
// Note that tables are just pointers in this case. Most users will likely
// remove the tables (s->Tables = nil) and add their own in the call to
// MarshalAll.
func (s *SDT) Marshal() ([]byte, error) {
	h, err := s.Generic.Header.Marshal()
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(h)
	x := s.Sig() == "XSDT"
	for _, p := range s.Tables {
		if x {
			w(b, p)
		} else {
			w(b, uint32(p))
		}
	}
	return b.Bytes(), nil
}

// MarshalAll marshals out an SDT, and all the tables, in a blob
// suitable for kexec. The most common use of this call would be to
// set s->Tables = nil and then pass in the desired Tables as
// parameters to this function.  For passed-in tables, all addresses
// are recomputed, as there may be more tables. Further, even if
// tables were scattered all over, we unify them into one segment.
// There is one potential problem, which we can fix if needed:
// it is possible the [XR]SDT is placed so close to the top of memory
// there is no room for the table. In the unlikely event that ever
// happens, we will just figure out how to place the tables in memory
// lower than the [XR]SDT.
func (s *SDT) MarshalAll(t ...Tabler) ([]byte, error) {
	var tabs [][]byte
	Debug("SDT MarshalAll has %d tables %d extra tables", len(s.Tables), len(t))

	// Serialize the tables from pointers in the SDT. Note that
	// this pointer can be nil. Depending on your kernel and its
	// config settings, you won't be able to read these anyway. In
	// the case of u-root kexec, this pointer is always nil as a
	// defensive measure.
	for i, addr := range s.Tables {
		t, err := ReadRaw(addr)
		if err != nil {
			return nil, err
		}
		Debug("SDT MarshalAll: processed table %d to %d bytes", i, len(t.AllData()))
		tabs = append(tabs, t.AllData())
	}

	// Serialize the extra tables.
	for i, tt := range t {
		b, err := Marshal(tt)
		if err != nil {
			return nil, err
		}
		Debug("SDT MarshalAll: processed extra table %d to %d bytes", i, len(b))
		tabs = append(tabs, b)
	}

	Debug("processed tables")
	// The length of the SDT is SSDTSize + len(s.Tables) * pointersize.
	// The number of tables will likely be different, so the
	// current value in the header is almost certainly wrong.
	// The easiest path here is to replace the
	// data with the new data, but first we have to compute the
	// pointers. So we do this as follows: truncate ssd to just
	// the header, serialize pointers, then get the size.
	s.Generic.data = s.Generic.data[:HeaderLength]
	var (
		addrs bytes.Buffer
		st    []byte
	)

	base := s.Base + HeaderLength // This is where the pointers start
	x := s.Sig() == "XSDT"
	if x {
		base += int64(len(tabs) * 8)
	} else {
		base += int64(len(tabs) * 4)
	}

	// We use base as a basic bump allocator.
	for i, t := range tabs {
		Debug("Table %d: len %d, base %#x", i, len(t), base)
		st = append(st, t...)
		if x {
			w(&addrs, uint64(base))
		} else {
			w(&addrs, uint32(base))
		}
		base += int64(len(t))
	}
	s.Generic.data = append(s.Generic.data, addrs.Bytes()...)
	h, err := s.Generic.Marshal()
	// If you get really desperate ...
	if false {
		Debug("marshalled sdt is ")
		d := hex.Dumper(os.Stdout)
		d.Write(h)
		d.Close()
	}
	if err != nil {
		return nil, err
	}

	// Append the tables. We have to do this after Marshaling the SDT
	// as the ACPI tables length should not be included in the SDT length.
	h = append(h, st...)
	return h, nil
}

// ReadSDT reads an SDT in from memory, using UnMarshalSDT, which uses
// the io package. This is increasingly unlikely to work over time.
func ReadSDT() (*SDT, error) {
	_, r, err := GetRSDP()
	if err != nil {
		return nil, err
	}
	s, err := UnMarshalSDT(r)
	return s, err
}

// NewSDT creates a new SDT, defaulting to XSDT.
func NewSDT(opt ...func(*SDT)) (*SDT, error) {
	var s = &SDT{
		Generic: Generic{
			Header: Header{
				Sig:             "XSDT",
				Length:          HeaderLength,
				Revision:        1,
				OEMID:           "GOOGLE",
				OEMTableID:      "ACPI=TOY",
				OEMRevision:     1,
				CreatorID:       1,
				CreatorRevision: 1,
			},
		},
	}
	for _, o := range opt {
		o(s)
	}
	// It may seem odd to check for a marshaling error
	// in something that does no I/O, but consider this
	// is a good place to see that the user did not set
	// something wrong.
	h, err := s.Marshal()
	if err != nil {
		return nil, err
	}
	s.data = h
	return s, nil
}
