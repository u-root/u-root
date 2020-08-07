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
//
// This file contains implementation of the math functions that are performed
// with canonical integers in byte buffers. The canonical integer is
// big-endian bytes.
//
#include "Tpm.h"

//** Functions

//*** UnsignedCmpB
// This function compare two unsigned values. The values are byte-aligned,
// big-endian numbers (e.g, a hash).
//  Return Type: int
//      1          if (a > b)
//      0          if (a = b)
//      -1         if (a < b)
LIB_EXPORT int
UnsignedCompareB(
    UINT32           aSize,         // IN: size of a
    const BYTE      *a,             // IN: a
    UINT32           bSize,         // IN: size of b
    const BYTE      *b              // IN: b
    )
{
    UINT32             i;
    if(aSize > bSize)
        return 1;
    else if(aSize < bSize)
        return -1;
    else
    {
        for(i = 0; i < aSize; i++)
        {
            if(a[i] != b[i])
                return (a[i] > b[i]) ? 1 : -1;
        }
    }
    return 0;
}

//***SignedCompareB()
// Compare two signed integers:
//  Return Type: int
//      1         if a > b
//      0         if a = b
//      -1        if a < b
int
SignedCompareB(
    const UINT32     aSize,         // IN: size of a
    const BYTE      *a,             // IN: a buffer
    const UINT32     bSize,         // IN: size of b
    const BYTE      *b              // IN: b buffer
    )
{
    int      signA, signB;       // sign of a and b

    // For positive or 0, sign_a is 1
    // for negative, sign_a is 0
    signA = ((a[0] & 0x80) == 0) ? 1 : 0;

    // For positive or 0, sign_b is 1
    // for negative, sign_b is 0
    signB = ((b[0] & 0x80) == 0) ? 1 : 0;

    if(signA != signB)
    {
        return signA - signB;
    }
    if(signA == 1)
        // do unsigned compare function
        return UnsignedCompareB(aSize, a, bSize, b);
    else
        // do unsigned compare the other way
        return 0 - UnsignedCompareB(aSize, a, bSize, b);
}

//*** ModExpB
// This function is used to do modular exponentiation in support of RSA.
// The most typical uses are: 'c' = 'm'^'e' mod 'n' (RSA encrypt) and
// 'm' = 'c'^'d' mod 'n' (RSA decrypt).  When doing decryption, the 'e' parameter
// of the function will contain the private exponent 'd' instead of the public
// exponent 'e'.
//
// If the results will not fit in the provided buffer,
// an error is returned (CRYPT_ERROR_UNDERFLOW). If the results is smaller
// than the buffer, the results is de-normalized.
//
// This version is intended for use with RSA and requires that 'm' be
// less than 'n'.
//
//  Return Type: TPM_RC
//      TPM_RC_SIZE         number to exponentiate is larger than the modulus
//      TPM_RC_NO_RESULT    result will not fit into the provided buffer
//
TPM_RC
ModExpB(
    UINT32           cSize,         // IN: the size of the output buffer. It will
                                    //     need to be the same size as the modulus
    BYTE            *c,             // OUT: the buffer to receive the results
                                    //     (c->size must be set to the maximum size
                                    //     for the returned value)
    const UINT32     mSize,
    const BYTE      *m,             // IN: number to exponentiate
    const UINT32     eSize,
    const BYTE      *e,             // IN: power
    const UINT32     nSize,
    const BYTE      *n              // IN: modulus
    )
{
    BN_MAX(bnC);
    BN_MAX(bnM);
    BN_MAX(bnE);
    BN_MAX(bnN);
    NUMBYTES         tSize = (NUMBYTES)nSize;
    TPM_RC           retVal = TPM_RC_SUCCESS;

    // Convert input parameters
    BnFromBytes(bnM, m, (NUMBYTES)mSize);
    BnFromBytes(bnE, e, (NUMBYTES)eSize);
    BnFromBytes(bnN, n, (NUMBYTES)nSize);


    // Make sure that the output is big enough to hold the result
    // and that 'm' is less than 'n' (the modulus)
    if(cSize < nSize)
        ERROR_RETURN(TPM_RC_NO_RESULT);
    if(BnUnsignedCmp(bnM, bnN) >= 0)
        ERROR_RETURN(TPM_RC_SIZE);
    BnModExp(bnC, bnM, bnE, bnN);
    BnToBytes(bnC, c, &tSize);
Exit:
    return retVal;
}

//*** DivideB()
// Divide an integer ('n') by an integer ('d') producing a quotient ('q') and
// a remainder ('r'). If 'q' or 'r' is not needed, then the pointer to them
// may be set to NULL.
//
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT         'q' or 'r' is too small to receive the result
//
LIB_EXPORT TPM_RC
DivideB(
    const TPM2B     *n,             // IN: numerator
    const TPM2B     *d,             // IN: denominator
    TPM2B           *q,             // OUT: quotient
    TPM2B           *r              // OUT: remainder
    )
{
    BN_MAX_INITIALIZED(bnN, n);
    BN_MAX_INITIALIZED(bnD, d);
    BN_MAX(bnQ);
    BN_MAX(bnR);
//
    // Do divide with converted values
    BnDiv(bnQ, bnR, bnN, bnD);

    // Convert the BIGNUM result back to 2B format using the size of the original
    // number
    if(q != NULL)
        if(!BnTo2B(bnQ, q, q->size))
            return TPM_RC_NO_RESULT;
    if(r != NULL)
        if(!BnTo2B(bnR, r, r->size))
            return TPM_RC_NO_RESULT;
    return TPM_RC_SUCCESS;
}

//*** AdjustNumberB()
// Remove/add leading zeros from a number in a TPM2B. Will try to make the number 
// by adding or removing leading zeros. If the number is larger than the requested 
// size, it will make the number as small as possible. Setting 'requestedSize' to 
// zero is equivalent to requesting that the number be normalized.
UINT16
AdjustNumberB(
    TPM2B           *num,
    UINT16           requestedSize
    )
{
    BYTE            *from;
    UINT16           i;
    // See if number is already the requested size
    if(num->size == requestedSize)
        return requestedSize;
    from = num->buffer;
    if (num->size > requestedSize)
    {
    // This is a request to shift the number to the left (remove leading zeros)
        // Find the first non-zero byte. Don't look past the point where removing
        // more zeros would make the number smaller than requested, and don't throw
        // away any significant digits.
        for(i = num->size; *from == 0 && i > requestedSize; from++, i--);
        if(i < num->size)
        {
            num->size = i;
            MemoryCopy(num->buffer, from, i);
        }
    }
    // This is a request to shift the number to the right (add leading zeros)
    else 
    {
        MemoryCopy(&num->buffer[requestedSize - num->size], num->buffer, num->size);
        MemorySet(num->buffer, 0, requestedSize- num->size);
        num->size = requestedSize;
    }
    return num->size;
}

//*** ShiftLeft()
// This function shifts a byte buffer (a TPM2B) one byte to the left. That is, 
// the most significant bit of the most significant byte is lost.
TPM2B *
ShiftLeft(
    TPM2B       *value          // IN/OUT: value to shift and shifted value out
)
{
    UINT16       count = value->size;
    BYTE        *buffer = value->buffer;
    if(count > 0)
    {
        for(count -= 1; count > 0; buffer++, count--)
        {
            buffer[0] = (buffer[0] << 1) + ((buffer[1] & 0x80) ? 1 : 0);
        }
        *buffer <<= 1;
    }
    return value;
}

