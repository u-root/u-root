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
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _BN_MATH_FP_H_
#define    _BN_MATH_FP_H_

//*** BnAdd()
// This function adds two bigNum values. This function always returns TRUE.
LIB_EXPORT BOOL
BnAdd(
    bigNum           result,
    bigConst         op1,
    bigConst         op2
);

//*** BnAddWord()
// This function adds a word value to a bigNum. This function always returns TRUE.
LIB_EXPORT BOOL
BnAddWord(
    bigNum           result,
    bigConst         op,
    crypt_uword_t    word
);

//*** BnSub()
// This function does subtraction of two bigNum values and returns result = op1 - op2
// when op1 is greater than op2. If op2 is greater than op1, then a fault is
// generated. This function always returns TRUE.
LIB_EXPORT BOOL
BnSub(
    bigNum           result,
    bigConst         op1,
    bigConst         op2
);

//*** BnSubWord()
// This function subtracts a word value from a bigNum. This function always
// returns TRUE.
LIB_EXPORT BOOL
BnSubWord(
    bigNum           result,
    bigConst     op,
    crypt_uword_t    word
);

//*** BnUnsignedCmp()
// This function performs a comparison of op1 to op2. The compare is approximately
// constant time if the size of the values used in the compare is consistent
// across calls (from the same line in the calling code).
//  Return Type: int
//      < 0             op1 is less than op2
//      0               op1 is equal to op2
//      > 0             op1 is greater than op2
LIB_EXPORT int
BnUnsignedCmp(
    bigConst               op1,
    bigConst               op2
);

//*** BnUnsignedCmpWord()
// Compare a bigNum to a crypt_uword_t.
//  Return Type: int
//      -1              op1 is less that word
//      0               op1 is equal to word
//      1               op1 is greater than word
LIB_EXPORT int
BnUnsignedCmpWord(
    bigConst             op1,
    crypt_uword_t        word
);

//*** BnModWord()
// This function does modular division of a big number when the modulus is a
// word value.
LIB_EXPORT crypt_word_t
BnModWord(
    bigConst         numerator,
    crypt_word_t     modulus
);

//*** Msb()
// This function returns the bit number of the most significant bit of a
// crypt_uword_t. The number for the least significant bit of any bigNum value is 0.
// The maximum return value is RADIX_BITS - 1,
//  Return Type: int
//      -1              the word was zero
//      n               the bit number of the most significant bit in the word
LIB_EXPORT int
Msb(
    crypt_uword_t           word
);

//*** BnMsb()
// This function returns the number of the MSb of a bigNum value.
//  Return Type: int
//      -1              the word was zero or 'bn' was NULL
//      n               the bit number of the most significant bit in the word
LIB_EXPORT int
BnMsb(
    bigConst            bn
);

//*** BnSizeInBits()
// This function returns the number of bits required to hold a number. It is one
// greater than the Msb.
//
LIB_EXPORT unsigned
BnSizeInBits(
    bigConst                 n
);

//*** BnSetWord()
// Change the value of a bignum_t to a word value.
LIB_EXPORT bigNum
BnSetWord(
    bigNum               n,
    crypt_uword_t        w
);

//*** BnSetBit()
// This function will SET a bit in a bigNum. Bit 0 is the least-significant bit in
// the 0th digit_t. The function always return TRUE
LIB_EXPORT BOOL
BnSetBit(
    bigNum           bn,        // IN/OUT: big number to modify
    unsigned int     bitNum     // IN: Bit number to SET
);

//*** BnTestBit()
// This function is used to check to see if a bit is SET in a bignum_t. The 0th bit
// is the LSb of d[0].
//  Return Type: BOOL
//      TRUE(1)         the bit is set
//      FALSE(0)        the bit is not set or the number is out of range
LIB_EXPORT BOOL
BnTestBit(
    bigNum               bn,        // IN: number to check
    unsigned int         bitNum     // IN: bit to test
);

//***BnMaskBits()
// This function is used to mask off high order bits of a big number.
// The returned value will have no more than 'maskBit' bits
// set.
// Note: There is a requirement that unused words of a bignum_t are set to zero.
//  Return Type: BOOL
//      TRUE(1)         result masked
//      FALSE(0)        the input was not as large as the mask
LIB_EXPORT BOOL
BnMaskBits(
    bigNum           bn,        // IN/OUT: number to mask
    crypt_uword_t    maskBit    // IN: the bit number for the mask.
);

//*** BnShiftRight()
// This function will shift a bigNum to the right by the shiftAmount.
// This function always returns TRUE.
LIB_EXPORT BOOL
BnShiftRight(
    bigNum           result,
    bigConst         toShift,
    uint32_t         shiftAmount
);

//*** BnGetRandomBits()
// This function gets random bits for use in various places. To make sure that the
// number is generated in a portable format, it is created as a TPM2B and then
// converted to the internal format.
//
// One consequence of the generation scheme is that, if the number of bits requested
// is not a multiple of 8, then the high-order bits are set to zero. This would come
// into play when generating a 521-bit ECC key. A 66-byte (528-bit) value is
// generated an the high order 7 bits are masked off (CLEAR).
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure
LIB_EXPORT BOOL
BnGetRandomBits(
    bigNum           n,
    size_t           bits,
    RAND_STATE      *rand
);

//*** BnGenerateRandomInRange()
// This function is used to generate a random number r in the range 1 <= r < limit.
// The function gets a random number of bits that is the size of limit. There is some
// some probability that the returned number is going to be greater than or equal
// to the limit. If it is, try again. There is no more than 50% chance that the
// next number is also greater, so try again. We keep trying until we get a
// value that meets the criteria. Since limit is very often a number with a LOT of
// high order ones, this rarely would need a second try.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure ('limit' is too small)
LIB_EXPORT BOOL
BnGenerateRandomInRange(
    bigNum           dest,
    bigConst         limit,
    RAND_STATE      *rand
);

#endif  // _BN_MATH_FP_H_
