// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build riscv64 && !tinygo

#include "textflag.h"

// Linux machine_kexec enters this snippet of code with the MMU off,
// register a0 set to the boot hart id, and register a1 set to
// whichever kexec segment looked like an FDT first.
//
// Each address is materialized with AUIPC+ADDI, which is PC-relative.
// Register a0 is left unchanged. Copying these instructions together
// with the data slots, keeps the references valid wherever kexec
// places the stub.
TEXT trampoline_start(SB),NOSPLIT,$0
	PCALIGN $2048
	MOV	$stack_top(SB), T0
	MOV	(T0), X2

	MOV	$hob_addr(SB), T0
	MOV	(T0), A1

	MOV	$entry_point(SB), T0
	MOV	(T0), T1
	JMP	(T1)

TEXT stack_top(SB),$0
	PCALIGN $32
	WORD	$0
TEXT hob_addr(SB),$0
	PCALIGN $32
	WORD	$0
TEXT entry_point(SB),$0
	PCALIGN $32
	WORD	$0

// func addrOfStart() uintptr
TEXT ·addrOfStart(SB),NOSPLIT,$0-8
	MOV	$trampoline_start(SB), T0
	MOV	T0, ret+0(FP)
	RET

// func addrOfStackTop() uintptr
TEXT ·addrOfStackTop(SB),NOSPLIT,$0-8
	MOV	$stack_top(SB), T0
	MOV	T0, ret+0(FP)
	RET

// func addrOfHobAddr() uintptr
TEXT ·addrOfHobAddr(SB),NOSPLIT,$0-8
	MOV	$hob_addr(SB), T0
	MOV	T0, ret+0(FP)
	RET

// func addrOfEntry() uintptr
TEXT ·addrOfEntry(SB),NOSPLIT,$0-8
	MOV	$entry_point(SB), T0
	MOV	T0, ret+0(FP)
	RET
