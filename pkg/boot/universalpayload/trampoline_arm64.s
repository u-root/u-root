// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"


TEXT trampoline_start(SB), NOSPLIT, $0
	PCALIGN $2048
	MOVD $entry_point(SB), R4
	MOVD (R4), R4

	MOVD $hob_addr(SB), R0
	MOVD (R0), R0

	MOVD $stack_top(SB), R2
	MOVD (R2), R2
	MOVD R2, RSP

	JMP  (R4)

TEXT stack_top(SB),$0
	PCALIGN $32
	NOP
TEXT hob_addr(SB),$0
	PCALIGN $32
	NOP
TEXT entry_point(SB),$0
	PCALIGN $32
	NOP

// func addrOfStart() uintptr
TEXT ·addrOfStart(SB), $0-8
	PCALIGN $32
	MOVD	$trampoline_start(SB), R0
	MOVD	R0, ret+0(FP)
	RET

// func addrOfStackTop() uintptr
TEXT ·addrOfStackTop(SB), $0-8
	MOVD	$stack_top(SB), R0
	MOVD	R0, ret+0(FP)
	RET

// func addrOfHobAddr() uintptr
TEXT ·addrOfHobAddr(SB), $0-8
	MOVD	$hob_addr(SB), R0
	MOVD	R0, ret+0(FP)
	RET
