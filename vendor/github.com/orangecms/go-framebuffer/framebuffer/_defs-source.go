// Copyright 2013 Konstantin Kulikov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package framebuffer

/*
#include <linux/fb.h>
#include <sys/mman.h>
*/
import "C"

type fixedScreenInfo C.struct_fb_fix_screeninfo
type variableScreenInfo C.struct_fb_var_screeninfo
type bitField C.struct_fb_bitfield

const (
	getFixedScreenInfo    uintptr = C.FBIOGET_FSCREENINFO
	getVariableScreenInfo uintptr = C.FBIOGET_VSCREENINFO
)

const (
	protocolRead  int = C.PROT_READ
	protocolWrite int = C.PROT_WRITE
	mapShared     int = C.MAP_SHARED
)
