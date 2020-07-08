// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package uefivars

import (
	"encoding/binary"
	"fmt"
)

// MixedGUID is a mixed-endianness guid, as used by MS and UEFI.
type MixedGUID [16]byte

// UUID uses the normal ordering, compatible with github.com/google/uuid. Use
// a package such as that if you need to generate UUIDs, check types, etc.
type UUID [16]byte

// ToStdEnc converts MixedGuid to a UUID.
func (m MixedGUID) ToStdEnc() (u UUID) {
	u[0], u[1], u[2], u[3] = m[3], m[2], m[1], m[0]
	u[4], u[5] = m[5], m[4]
	u[6], u[7] = m[7], m[6]
	copy(u[8:], m[8:])
	return
}

// String converts a MixedGUID to string.
func (m MixedGUID) String() string {
	le := binary.LittleEndian
	data1 := le.Uint32(m[:4])
	data2 := le.Uint16(m[4:6])
	data3 := le.Uint16(m[6:8])
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", data1, data2, data3, m[8:10], m[10:])
}

// ToMixedGuid converts UUID to MixedGuid.
func (u UUID) ToMixedGUID() (m MixedGUID) {
	m[0], m[1], m[2], m[3] = u[3], u[2], u[1], u[0]
	m[4], m[5] = u[5], u[4]
	m[6], m[7] = u[7], u[6]
	copy(m[8:], u[8:])
	return
}

// String converts a UUID to string.
func (u UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}
