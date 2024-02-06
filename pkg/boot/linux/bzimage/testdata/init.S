#include <x86_64-linux-gnu/asm/unistd_64.h>
1:
	.string	"------------------------------------------ hi from bzimage test ---------------------\n"
.globl _start
_start:
	movl	$70, %edx
	movq	$1b, %rsi
	movl	$1, %edi
	movl	$__NR_write, %eax 
	syscall
1:	jmp 1b
