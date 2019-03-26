// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

type SDT struct {
	Header
	Tables []uint64
}

func unmarshalSDT(t Tabler) (Tabler, error) {
	s := &SDT{
		Header: *GetHeader(t),
	}
	// Now the fun. In 1999, 64-bit micros had been out for about 10 years.
	// Intel had announced the ia64 years earlier. In 2000 the ACPI committee
	// chose 32-bit pointers anyway, then had to backfill a bunch of table
	// types to do 64 bits.
	return s, nil
}
