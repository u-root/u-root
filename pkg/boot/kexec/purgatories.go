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
.align 8
dat:
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
	movl 	dat, %eax
	1: jmp 1b
	jmp *%eax
	

1: jmp 1b
.data
.globl entry32_regs
entry32_regs: .long 0
.globl cmdline_end
cmdline_end: .long 0
`,
	},
	{
		name: "to64",
		cc:   []string{"x86_64-linux-gnu-gcc", "-c", "-nostdlib", "-nostdinc", "-static"},
		ld:   []string{"ld", "-N", "-e entry64", "-Ttext=0x3000"},
		code: `
jmp 1f
.align 8
dat:
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
	movq 	dat, %rax
	1: jmp 1b
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
