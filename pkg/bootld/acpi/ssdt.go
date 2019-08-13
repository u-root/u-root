// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import "bytes"

// genssdt generates an ssdt header for an existing []byte.
// We are unlikely to ever need this, and may remove it. It is from an
// early iteration of this package.
// It is only useful if the slice contains AML, not a table, so we don't actually
// use it currently. It turns out that it's easisest just to do it programatically
// with a function call, rather than create a magic struct and serializing it out, due to
// the need to create a checksum once it's all assembled.
// From the Big Bad Book of ACPI:
// SSDT
// Signature 4 0 ‘SSDT’ Signature for the Secondary System Description Table.
// Length 4 4 Length, in bytes, of the entire SSDT (including the header).
// Revision 1 8 2
// Checksum 1 9 Entire table must sum to zero.
// OEMID 6 10 OEM ID
// OEM Table ID 8 16 The manufacture model ID.
// OEM Revision 4 24 OEM revision of DSDT for supplied OEM Table ID.
// Creator ID 4 28 Vendor ID for the ASL Compiler.
// Creator Revision 4 32 Revision number of the ASL Compiler.
func genssdt(b []byte) []byte {
	var (
		ssdt = &bytes.Buffer{}
		csum uint8
	)

	l := uint32(HeaderLength + len(b))
	w(ssdt, 1, []byte("SSDT"), l, uint8(0), csum, []byte("ACPIXX"), []byte("GOXR00LZ"), uint32(0), []byte("VEND"), uint32(0xdecafbad), b)
	csum = gencsum(ssdt.Bytes())
	Debug("CSUM is %#x", csum)
	s := ssdt.Bytes()
	s[CSUMOffset] = csum
	return s
}
