// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align

import "os"

var pageMask = uint(os.Getpagesize() - 1)

func AlignUpPageSize(p uint) uint {
	return (p + pageMask) &^ pageMask
}

func AlignUpPageSizePtr(p uintptr) uintptr {
	return uintptr(AlignUpPageSize(uint(p)))
}
