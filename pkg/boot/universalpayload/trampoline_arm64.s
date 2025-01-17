// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// func addrOfStart() *byte
TEXT ·addrOfStart(SB), $0-8
	PCALIGN $32
	MOVD	$trampoline_start(SB), R0
	MOVD	R0, ret+0(FP)
	RET

// func addrOfStackTop() *byte
TEXT ·addrOfStackTop(SB), $0-8
	MOVD	$stack_top(SB), R0
	MOVD	R0, ret+0(FP)
	RET

// func addrOfHobAddr() *byte
TEXT ·addrOfHobAddr(SB), $0-8
	MOVD	$hob_addr(SB), R0
	MOVD	R0, ret+0(FP)
	RET
