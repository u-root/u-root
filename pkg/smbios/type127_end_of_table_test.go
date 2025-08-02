// Copyright 2016-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import "testing"

var validEndOfTableData = []byte{127, 4, 255, 254}

func validEndOfTableRaw(t *testing.T) []byte {
	return joinBytesT(t,
		validEndOfTableData,
		0, // String terminator
		0, // Table terminator
	)
}
