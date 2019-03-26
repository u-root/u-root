// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/binary"
	"fmt"
)

type SDT struct {
	Generic
	Tables []uint64
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
	// types to do 64 bits. Geez.
	esize := 4
	if sig == "XSDT" {
		esize = 8
	}
	d := t.TableData()

	for i := 0; i < len(d); i += esize {
		val := uint64(0)
		if sig == "XSDT" {
			val = binary.LittleEndian.Uint64(d[i : i+8])
		} else {
			val = uint64(binary.LittleEndian.Uint32(d[i : i+4]))
		}
		s.Tables = append(s.Tables, val)
	}
	return s, nil
}
