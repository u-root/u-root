/*
 * kexec: Linux boots Linux
 *
 * Copyright (C) 2003,2004  Eric Biederman (ebiederm@xmission.com)
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation (version 2 of the License).
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.
 */

#include "arch/debug.h"
		
.text
.code64
	.balign 16	
	.globl entry64, entry64_regs
entry64:
	/* Don't worry about special registers... */

	/* Setup a gdt that should be preserved */
	lgdt	gdt(%rip)

	/* load the data segments */
	movl	$0x18, %eax	/* data segment */
	movl	%eax, %ds
	movl	%eax, %es
	movl	%eax, %ss
	movl	%eax, %fs
	movl	%eax, %gs

	leaq	stack_init(%rip), %rsp
	pushq	$0x10 /* CS */
	leaq	new_cs_exit(%rip), %rax
	pushq	%rax
	lretq
new_cs_exit:

	/* Load the registers */
	movq	rax(%rip), %rax
	movq	rbx(%rip), %rbx
	movq	rcx(%rip), %rcx
	movq	rdx(%rip), %rdx
	movq	rsi(%rip), %rsi
	movq	rdi(%rip), %rdi
	movq	rsp(%rip), %rsp
	movq	rbp(%rip), %rbp
	movq	r8(%rip), %r8
	movq	r9(%rip), %r9
	movq	r10(%rip), %r10
	movq	r11(%rip), %r11
	movq	r12(%rip), %r12
	movq	r13(%rip), %r13
	movq	r14(%rip), %r14
	movq	r15(%rip), %r15

	/* Jump to the new code... */
	jmpq	*rip(%rip)
		
	.section ".rodata"       
	.balign 4
entry64_regs:
rax:	.quad 0x00000000
rbx:	.quad 0x00000000
rcx:	.quad 0x00000000
rdx:	.quad 0x00000000
rsi:	.quad 0x00000000
rdi:	.quad 0x00000000
rsp:	.quad 0x00000000
rbp:	.quad 0x00000000
r8:	.quad 0x00000000
r9:	.quad 0x00000000
r10:	.quad 0x00000000
r11:	.quad 0x00000000
r12:	.quad 0x00000000
r13:	.quad 0x00000000
r14:	.quad 0x00000000
r15:	.quad 0x00000000
rip:	.quad entry32
	.size entry64_regs, . - entry64_regs		

	.section ".rodata"       
	.balign 16
gdt:	/* 0x00 unusable segment 
	 * 0x08 unused
	 * so use them as the gdt ptr
	 */
	.word	gdt_end - gdt - 1
	.quad	gdt
	.word	0, 0, 0
			
	/* 0x10 4GB flat code segment */
	.word	0xFFFF, 0x0000, 0x9A00, 0x00AF
	
	/* 0x18 4GB flat data segment */
	.word	0xFFFF, 0x0000, 0x9200, 0x00CF
gdt_end:
stack:	.quad	0, 0
stack_init:
