// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

	// From go runtime compiler source:
	//TODO: INB DX, AL                      // ec
	//TODO: INW DX, AX                      // 66ed
	//TODO: INL DX, AX                      // ed

TEXT ·archInb(SB),$0-5
	MOVW    arg+0(FP), DX
	BYTE	$0xec //INB	DX, AL
	MOVB    AX, ret+4(FP)
	RET

TEXT ·archInw(SB),$0-6
	MOVW    arg+0(FP), DX
	BYTE	$0x66 // Do the next instruction (INL) in 16-bit mode
	BYTE	$0xec //INW	DX, AL
	MOVW    AX, ret+4(FP)
	RET


TEXT ·archInl(SB),$0-8
	MOVW    arg+0(FP), DX
	BYTE	$0xed //INL	DX, AL
	MOVL    AX, ret+4(FP)
	RET

	// From go runtime compiler source:
	//TODO: OUTB AL, DX                     // ee
	//TODO: OUTW AX, DX                     // 66ef
	//TODO: OUTL AX, DX                     // ef

TEXT ·archOutb(SB),$0-3
	MOVW    arg+0(FP), DX
	MOVB	arg1+2(FP), AX
	BYTE	$0xee //OUTB	DX, AL
	RET

TEXT ·archOutw(SB),$0-4
	MOVW    arg+0(FP), DX
	MOVW	arg1+2(FP), AX
	BYTE	$0x66 // Do the next instruction (OUTL) in 16-bit mode
	BYTE	$0xef //OUTW	DX, AL
	RET


TEXT ·archOutl(SB),$0-8
	MOVW    arg+0(FP), DX
	MOVL	arg1+4(FP), AX
	BYTE	$0xef //OUTL	DX, AL
	RET
