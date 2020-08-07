/* Microsoft Reference Implementation for TPM 2.0
 *
 *  The copyright in this software is being made available under the BSD License,
 *  included below. This software may be subject to other third party and
 *  contributor rights, including patent rights, and no such rights are granted
 *  under this license.
 *
 *  Copyright (c) Microsoft Corporation
 *
 *  All rights reserved.
 *
 *  BSD License
 *
 *  Redistribution and use in source and binary forms, with or without modification,
 *  are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list
 *  of conditions and the following disclaimer.
 *
 *  Redistributions in binary form must reproduce the above copyright notice, this
 *  list of conditions and the following disclaimer in the documentation and/or
 *  other materials provided with the distribution.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ""AS IS""
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 *  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 *  ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 *  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 *  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 *  ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 *  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 *  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
//** Introduction
// This file contains bit manipulation routines.  They operate on bit arrays.
//
// The 0th bit in the array is the right-most bit in the 0th octet in
// the array.
//
// NOTE: If pAssert() is defined, the functions will assert if the indicated bit
// number is outside of the range of 'bArray'. How the assert is handled is
// implementation dependent.

//** Includes

#include "Tpm.h"

//** Functions

//*** TestBit()
// This function is used to check the setting of a bit in an array of bits.
//  Return Type: BOOL
//      TRUE(1)         bit is set
//      FALSE(0)        bit is not set
BOOL
TestBit(
    unsigned int     bitNum,        // IN: number of the bit in 'bArray'
    BYTE            *bArray,        // IN: array containing the bits
    unsigned int     bytesInArray   // IN: size in bytes of 'bArray'
    )
{
    pAssert(bytesInArray > (bitNum >> 3));
    return((bArray[bitNum >> 3] & (1 << (bitNum & 7))) != 0);
}

//*** SetBit()
// This function will set the indicated bit in 'bArray'.
void
SetBit(
    unsigned int     bitNum,        // IN: number of the bit in 'bArray'
    BYTE            *bArray,        // IN: array containing the bits
    unsigned int     bytesInArray   // IN: size in bytes of 'bArray'
    )
{
    pAssert(bytesInArray > (bitNum >> 3));
    bArray[bitNum >> 3] |= (1 << (bitNum & 7));
}

//*** ClearBit()
// This function will clear the indicated bit in 'bArray'.
void
ClearBit(
    unsigned int     bitNum,        // IN: number of the bit in 'bArray'.
    BYTE            *bArray,        // IN: array containing the bits
    unsigned int     bytesInArray   // IN: size in bytes of 'bArray'
    )
{
    pAssert(bytesInArray > (bitNum >> 3));
    bArray[bitNum >> 3] &= ~(1 << (bitNum & 7));
}

