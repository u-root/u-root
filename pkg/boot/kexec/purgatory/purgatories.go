// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

type asm struct {
	name string
	cc   []string
	ld   []string
	code string
}

var asms = []asm{
	{
		name: "to32bit_3000",
		cc:   []string{"x86_64-linux-gnu-gcc", "-c", "-nostdlib", "-nostdinc", "-static"},
		ld:   []string{"ld", "-N", "-e entry64", "-Ttext=0x3000"},
		code: `
# This code does its best to be PIC. As in, PIC by code, not by linker.
# assumptions:
# we have 8 bytes of stack to do ONE call 
# we are in shared space -- code and data in dram
# we are in long mode
# in case of emergency, break glass.
// 1: jmp 1b

jmp 1f
# save space for two arguments and stack
.align 8
entry: .long 0
params: .long 0
# This can be useful for being sure things are where we think they are,
# and ensuring we are storing correctly.
.long 0xdeadcafe, 0x12345678
.long 0xcafe5555, 0xabcdefaa
.align 128
1:
call 1f
1:
popq %rax
andl $0xfffffff0, %eax
movl %eax, %esp
movl $gdt, %eax
jmp 1f
# recursive gdt. 
	.balign 16
gdt:	/* 0x00 unusable segment 
	 * 0x08 unused
	 * so use them as the gdt ptr
	 */
	.word	gdt_end - gdt - 1
	.quad	gdt
	.word	0, 0, 0
			
	/* Documented linux kernel segments */
	/* 0x10 4GB flat code segment */
	.word	0xFFFF, 0x0000, 0x9A00, 0x00CF
	/* 0x18 4GB flat data segment */
	.word	0xFFFF, 0x0000, 0x9200, 0x00CF
gdt_end:
# NOTE: NOT PIC
to32indir:
	.long to32
	.long 0x10

1: 
	lgdt	gdt(%rip)
	ljmp	*to32indir(%rip)
to32:
	.code32

	.equ	CR0_PG,        0x80000000
	/* Disable paging */
	movl	%cr0, %eax
	andl	$~CR0_PG, %eax
	movl	%eax, %cr0

	/* Disable long mode */
	.equ	MSR_K6_EFER,   0xC0000080
	movl	$MSR_K6_EFER, %ecx
	rdmsr
	.equ	EFER_LME,      0x00000100
	andl	$~EFER_LME, %eax
	wrmsr

	/* Disable PAE */
	xorl	%eax, %eax
	movl	%eax, %cr4

	/* load the data segments */
	movl	$0x18, %eax	/* data segment */
	movl	%eax, %ds
	movl	%eax, %es
	movl	%eax, %ss
	movl	%eax, %fs
	movl	%eax, %gs
	movl 	entry, %eax
	movl 	params, %esi
	# in case of emergency, uncomment this. 1: jmp 1b
	jmp *%eax
	

1: jmp 1b
.data
.globl entry32_regs
entry32_regs: .long 0
.globl cmdline_end
cmdline_end: .long 0
`,
	},
	// This is the default purgatory, a simple from-64 to-64 trampoline
	{
		name: "default",
		cc:   []string{"x86_64-linux-gnu-gcc", "-c", "-nostdlib", "-nostdinc", "-static"},
		ld:   []string{"ld", "-N", "-e entry64", "-Ttext=0x3000"},
		code: `
jmp 1f
.align 8
// Known to Go.
entry: .long 0x87654321
.long 0xcafecafe
params: .long 0xab5aab5a
.long 0xbeeffeed
// end Known to Go
# This can be useful for being sure things are where we think they are,
# and ensuring we are storing correctly.
.long 0xdeadcafe, 0x12345678
.long 0xcafe5555, 0xabcdefaa
# save space for two arguments and stack
.align 128
1:
call 1f
1:
popq %rax
andl $0xfffffff0, %eax
movl %eax, %esp
/* Crudely reset a VGA card to text mode 3, by writing plausible default    */
/* values into its registers.                                               */
/* Tim Deegan (tjd21 at cl.cam.ac.uk), March 2003                              */
#define inb(p) movw $p, %dx; inb %dx, %al
#define outb(v, p) movw $p, %dx; movb $v, %al; outb %al, %dx
#define outw(v, p) movw $p, %dx; movw $v, %ax; outw %ax, %dx

	/* Hello */
	inb(0x3da)
	outb(0, 0x3c0)

	/* Sequencer registers */
	outw(0x0300, 0x3c4)
	outw(0x0001, 0x3c4)
	outw(0x0302, 0x3c4)
	outw(0x0003, 0x3c4)
	outw(0x0204, 0x3c4)

	/* Ensure CRTC regs 0-7 are unlocked by clearing bit 7 of CRTC[17] */
	outw(0x0e11, 0x3d4)
	/* CRTC registers */
	outw(0x5f00, 0x3d4)
	outw(0x4f01, 0x3d4)
	outw(0x5002, 0x3d4)
	outw(0x8203, 0x3d4)
	outw(0x5504, 0x3d4)
	outw(0x8105, 0x3d4)
	outw(0xbf06, 0x3d4)
	outw(0x1f07, 0x3d4)
	outw(0x0008, 0x3d4)
	outw(0x4f09, 0x3d4)
	outw(0x200a, 0x3d4)
	outw(0x0e0b, 0x3d4)
	outw(0x000c, 0x3d4)
	outw(0x000d, 0x3d4)
	outw(0x010e, 0x3d4)
	outw(0xe00f, 0x3d4)
	outw(0x9c10, 0x3d4)
	outw(0x8e11, 0x3d4)
	outw(0x8f12, 0x3d4)
	outw(0x2813, 0x3d4)
	outw(0x1f14, 0x3d4)
	outw(0x9615, 0x3d4)
	outw(0xb916, 0x3d4)
	outw(0xa317, 0x3d4)
	outw(0xff18, 0x3d4)

	/* Graphic registers */
	outw(0x0000, 0x3ce)
	outw(0x0001, 0x3ce)
	outw(0x0002, 0x3ce)
	outw(0x0003, 0x3ce)
	outw(0x0004, 0x3ce)
	outw(0x1005, 0x3ce)
	outw(0x0e06, 0x3ce)
	outw(0x0007, 0x3ce)
	outw(0xff08, 0x3ce)

	/* Attribute registers */
	inb(0x3da)
	outb(0x00, 0x3c0)
	outb(0x00, 0x3c0)

	inb(0x3da)
	outb(0x01, 0x3c0)
	outb(0x01, 0x3c0)

	inb(0x3da)
	outb(0x02, 0x3c0)
	outb(0x02, 0x3c0)

	inb(0x3da)
	outb(0x03, 0x3c0)
	outb(0x03, 0x3c0)

	inb(0x3da)
	outb(0x04, 0x3c0)
	outb(0x04, 0x3c0)

	inb(0x3da)
	outb(0x05, 0x3c0)
	outb(0x05, 0x3c0)

	inb(0x3da)
	outb(0x06, 0x3c0)
	outb(0x14, 0x3c0)

	inb(0x3da)
	outb(0x07, 0x3c0)
	outb(0x07, 0x3c0)

	inb(0x3da)
	outb(0x08, 0x3c0)
	outb(0x38, 0x3c0)

	inb(0x3da)
	outb(0x09, 0x3c0)
	outb(0x39, 0x3c0)

	inb(0x3da)
	outb(0x0a, 0x3c0)
	outb(0x3a, 0x3c0)

	inb(0x3da)
	outb(0x0b, 0x3c0)
	outb(0x3b, 0x3c0)

	inb(0x3da)
	outb(0x0c, 0x3c0)
	outb(0x3c, 0x3c0)

	inb(0x3da)
	outb(0x0d, 0x3c0)
	outb(0x3d, 0x3c0)

	inb(0x3da)
	outb(0x0e, 0x3c0)
	outb(0x3e, 0x3c0)

	inb(0x3da)
	outb(0x0f, 0x3c0)
	outb(0x3f, 0x3c0)

	inb(0x3da)
	outb(0x10, 0x3c0)
	outb(0x0c, 0x3c0)

	inb(0x3da)
	outb(0x11, 0x3c0)
	outb(0x00, 0x3c0)

	inb(0x3da)
	outb(0x12, 0x3c0)
	outb(0x0f, 0x3c0)

	inb(0x3da)
	outb(0x13, 0x3c0)
	outb(0x08, 0x3c0)

	inb(0x3da)
	outb(0x14, 0x3c0)
	outb(0x00, 0x3c0)

	/* Goodbye */
	inb(0x3da)
	outb(0x20, 0x3c0)
	movq 	entry, %rax
	movq 	params, %rsi
	# in case of emergency, uncomment this.
	# 1: jmp 1b
	jmp *%rax
	

1: jmp 1b
.data
.globl entry32_regs
entry32_regs: .long 0
.globl cmdline_end
cmdline_end: .long 0
`,
	},
	{
		name: "loop_3000",
		cc:   []string{"x86_64-linux-gnu-gcc", "-c", "-nostdlib", "-nostdinc", "-static"},
		ld:   []string{"ld", "-N", "-e entry64", "-Ttext=0x3000"},
		code: `
1: jmp 1b
`,
	},
}
