// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo && linux && (amd64 || 386)

package memio

/*
#include "ioport_linux_tinygo.h"
*/
import "C"

func archInl(port uint16) uint32 {
	return C.archInl(port)
}

func archInw(port uint16) uint16 {
	return C.archInw(port)
}

func archInb(port uint16) uint8 {
	return C.archInb(port)
}

func archOutl(port uint16, val uint32) {
	C.archOutl(port, val)
}

func archOutw(port uint16, val uint16) {
	C.archOutw(port, val)
}

func archOutb(port uint16, val uint8) {
	C.archOutb(port, val)
}
