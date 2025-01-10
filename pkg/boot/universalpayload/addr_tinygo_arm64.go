// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && tinygo

package universalpayload

/*
#include "trampoline_tinygo_arm64.h"
*/
import "C"

func addrOfStart() uintptr {
	return C.addrOfStartU()
}

func addrOfStackTop() uintptr {
	return C.addrOfStackTopU()
}

func addrOfHobAddr() uintptr {
	return C.addrOfHobAddrU()
}
