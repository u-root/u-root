// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import "fmt"

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

	if s.Sig() != "RSDT" && s.Sig() != "XSDT" {
		return nil, fmt.Errorf("%v is not RSDT or XSDT", s.Sig())
	}

	// Now the fun. In 1999, 64-bit micros had been out for about 10 years.
	// Intel had announced the ia64 years earlier. In 2000 the ACPI committee
	// chose 32-bit pointers anyway, then had to backfill a bunch of table
	// types to do 64 bits.
	return s, nil
}
