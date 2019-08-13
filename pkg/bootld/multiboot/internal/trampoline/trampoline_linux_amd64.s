// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// We want all the trampoline's assembly code to be located
// in a contiguous byte range in a compiled binary.
// Go compiler does not guarantee that. Still, current version
// of compiler puts all pieces together.

#include "textflag.h"

#define MSR_EFER	0xC0000080
#define EFER_LME	0xFFFFFEFF
#define CR0_PG		0x0FFFFFFF

#define DATA_SEGMENT	0x00CF92000000FFFF
#define CODE_SEGMENT	0x00CF9A000000FFFF

#define MAGIC	0x2BADB002

TEXT begin(SB),NOSPLIT,$0
	// u-root-trampoline-begin
	BYTE $'u'; BYTE $'-'; BYTE $'r'; BYTE $'o'; BYTE $'o';
	BYTE $'t'; BYTE $'-'; BYTE $'t'; BYTE $'r'; BYTE $'a';
	BYTE $'m'; BYTE $'p'; BYTE $'o'; BYTE $'l'; BYTE $'i';
	BYTE $'n'; BYTE $'e'; BYTE $'-'; BYTE $'b'; BYTE $'e';
	BYTE $'g'; BYTE $'i'; BYTE $'n';

TEXT ·start(SB),NOSPLIT,$0
	// Create GDT pointer on stack.
	LEAQ	gdt(SB), CX
	SHLQ	$16, CX
	ORQ	$(4*8 - 1), CX
	PUSHQ	CX

	LGDT	(SP)

	// Store value of multiboot info addr in BX.
	// Don't modify BX.
	MOVL	info(SB), BX

	// Far return doesn't work on QEMU in 64-bit mode,
	// let's do far jump.
	//
	// In a regular plan9 assembly we can do something like:
	//	BYTE	$0xFF; BYTE $0x2D
	//	LONG	$bootaddr(SB)
	// TEXT bootaddr(SB),NOSPLIT,$0
	//	LONG	$boot(SB)
	//	LONG	$0x8
	//
	// Go compiler doesn't let us do it.
	//
	// Setup offset to make a far jump from boot(SB)
	// to a final kernel in a 32-bit mode.
	MOVL	entry(SB), AX
	MOVL	AX, farjump32+1(SB)

	// Setup offset to make a far jump to boot(SB)
	// to switch from 64-bit mode to 32-bit mode.
	LEAQ	boot(SB), CX
	MOVL	CX, farjump64+6(SB)
	JMP	farjump64(SB)


TEXT boot(SB),NOSPLIT,$0
	// We are in 32-bit mode now.
	//
	// Be careful editing this code!!! Go compiler
	// interprets all commands as 64-bit commands.

	// Disable paging.
	MOVL	CR0, AX
	ANDL	$CR0_PG, AX
	MOVL	AX, CR0

	// Disable long mode.
	MOVL	$MSR_EFER, CX
	RDMSR
	ANDL	$EFER_LME, AX
	WRMSR

	// Disable PAE.
	XORL	AX, AX
	MOVL	AX, CR4

	// Load data segments.
	MOVL	$0x10, AX // GDT 0x10 data segment
	BYTE	$0x8e; BYTE $0xd8 // MOVL AX, DS
	BYTE	$0x8e; BYTE $0xc0 // MOVL AX, ES
	BYTE	$0x8e; BYTE $0xd0 // MOVL AX, SS
	BYTE	$0x8e; BYTE $0xe0 // MOVL AX, FS
	BYTE	$0x8e; BYTE $0xe8 // MOVL AX, GS

	MOVL	$MAGIC, AX
	JMP	farjump32(SB)

	// Unreachable code.
	// Need reference text labels for compiler to
	// include them to a binary.
	JMP	begin(SB)
	JMP	infotext(SB)
	JMP	entrytext(SB)
	JMP	end(SB)

TEXT farjump64(SB),NOSPLIT,$0
	BYTE	$0xFF; BYTE $0x2D; LONG $0x0 // ljmp *(ip)

	LONG	$0x0 // farjump64+6(SB)
	LONG	$0x8 // code segment

TEXT farjump32(SB),NOSPLIT,$0
	// ljmp $0x18, offset
	BYTE	$0xEA
	LONG	$0x0 // farjump32+1(SB)
	WORD	$0x18 // code segment

TEXT gdt(SB),NOSPLIT,$0
	QUAD	$0x0		// 0x0 null entry
	QUAD	$CODE_SEGMENT	// 0x8
	QUAD	$DATA_SEGMENT	// 0x10
	QUAD	$CODE_SEGMENT	// 0x18

TEXT infotext(SB),NOSPLIT,$0
	// u-root-info-long
	BYTE $'u'; BYTE $'-'; BYTE $'r'; BYTE $'o'; BYTE $'o';
	BYTE $'t'; BYTE $'-'; BYTE $'i'; BYTE $'n'; BYTE $'f';
	BYTE $'o'; BYTE $'-'; BYTE $'l'; BYTE $'o'; BYTE $'n';
	BYTE $'g';
TEXT info(SB),NOSPLIT,$0
	LONG	$0x0

TEXT entrytext(SB),NOSPLIT,$0
	// u-root-entry-long
	BYTE $'u'; BYTE $'-'; BYTE $'r'; BYTE $'o'; BYTE $'o';
	BYTE $'t'; BYTE $'-'; BYTE $'e'; BYTE $'n'; BYTE $'t';
	BYTE $'r'; BYTE $'y'; BYTE $'-'; BYTE $'l'; BYTE $'o';
	BYTE $'n'; BYTE $'g';
TEXT entry(SB),NOSPLIT,$0
	LONG	$0x0

TEXT end(SB),NOSPLIT,$0
	// u-root-trampoline-end
	BYTE $'u'; BYTE $'-'; BYTE $'r'; BYTE $'o'; BYTE $'o';
	BYTE $'t'; BYTE $'-'; BYTE $'t'; BYTE $'r'; BYTE $'a';
	BYTE $'m'; BYTE $'p'; BYTE $'o'; BYTE $'l'; BYTE $'i';
	BYTE $'n'; BYTE $'e'; BYTE $'-'; BYTE $'e'; BYTE $'n';
	BYTE $'d';
