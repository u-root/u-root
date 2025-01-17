// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// func addrOfStart() *byte
TEXT ·addrOfStart(SB), $0-8
	PCALIGN $32
	MOVQ	$trampoline_start(SB), AX
	MOVQ	AX, ret+0(FP)
	RET

// func addrOfStackTop() *byte
TEXT ·addrOfStackTop(SB), $0-8
	MOVQ	$stack_top(SB), AX
	MOVQ	AX, ret+0(FP)
	RET

// func addrOfHobAddr() *byte
TEXT ·addrOfHobAddr(SB), $0-8
	MOVQ	$hob_addr(SB), AX
	MOVQ	AX, ret+0(FP)
	RET
