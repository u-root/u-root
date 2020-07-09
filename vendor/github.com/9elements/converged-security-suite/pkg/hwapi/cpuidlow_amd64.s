// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
#include "textflag.h"

// func cpuidLow(arg1, arg2 uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuidLow(SB),NOSPLIT,$0-24
    MOVL    arg1+0(FP), AX
    MOVL    arg2+4(FP), CX
    CPUID
    MOVL AX, eax+8(FP)
    MOVL BX, ebx+12(FP)
    MOVL CX, ecx+16(FP)
    MOVL DX, edx+20(FP)
    RET
