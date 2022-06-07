// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package align

import "os"

var pageSize = uint(os.Getpagesize())

func AlignUpBySize(v uint, alignSize uint) uint {
	// Align everything to at least page size.
	if alignSize < pageSize {
		alignSize = pageSize
	}
	mask := alignSize - 1
	return (v + mask) &^ mask
}

func AlignUpPageSize(p uint) uint {
	return AlignUpBySize(p, pageSize)
}

func AlignUpPageSizePtr(p uintptr) uintptr {
	return uintptr(AlignUpPageSize(uint(p)))
}
