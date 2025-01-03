// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT trampoline_start(SB),NOSPLIT,$0
	PCALIGN $2048
	LEAQ stack_top(SB), AX
	MOVQ (AX), AX
	MOVQ AX, SP

	LEAQ hob_addr(SB), AX
	MOVQ (AX), AX
	MOVQ AX, CX

	LEAQ entry_point(SB), AX
	MOVQ (AX), AX
	JMP  AX

TEXT stack_top(SB),$0
	PCALIGN $32
	LONG	$0xdeadbeef
TEXT hob_addr(SB),$0
	PCALIGN $32
	LONG	$0xdeadbeef
TEXT entry_point(SB),$0
	PCALIGN $32
	LONG	$0xdeadbeef

// func addrOfStart() uintptr
TEXT ·addrOfStart(SB), $0-8
	PCALIGN $32
	MOVQ	$trampoline_start(SB), AX
	MOVQ	AX, ret+0(FP)
	RET

// func addrOfStackTop() uintptr
TEXT ·addrOfStackTop(SB), $0-8
	MOVQ	$stack_top(SB), AX
	MOVQ	AX, ret+0(FP)
	RET

// func addrOfHobAddr() uintptr
TEXT ·addrOfHobAddr(SB), $0-8
	MOVQ	$hob_addr(SB), AX
	MOVQ	AX, ret+0(FP)
	RET
