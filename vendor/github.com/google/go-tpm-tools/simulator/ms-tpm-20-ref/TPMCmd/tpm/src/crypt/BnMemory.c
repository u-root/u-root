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
// This file contains the memory setup functions used by the bigNum functions
// in CryptoEngine

//** Includes
#include "Tpm.h"

//** Functions

//*** BnSetTop()
// This function is used when the size of a bignum_t is changed. It
// makes sure that the unused words are set to zero and that any significant
// words of zeros are eliminated from the used size indicator.
LIB_EXPORT bigNum
BnSetTop(
    bigNum           bn,        // IN/OUT: number to clean
    crypt_uword_t    top        // IN: the new top
    )
{
    if(bn != NULL)
    {
        pAssert(top <= bn->allocated);
        // If forcing the size to be decreased, make sure that the words being
        // discarded are being set to 0
        while(bn->size > top)
            bn->d[--bn->size] = 0;
        bn->size = top;
        // Now make sure that the words that are left are 'normalized' (no high-order
        // words of zero.
        while((bn->size > 0) && (bn->d[bn->size - 1] == 0))
            bn->size -= 1;
    }
    return bn;
}

//*** BnClearTop()
// This function will make sure that all unused words are zero.
LIB_EXPORT bigNum
BnClearTop(
    bigNum          bn
    )
{
    crypt_uword_t       i;
//
    if(bn != NULL)
    {
        for(i = bn->size; i < bn->allocated; i++)
            bn->d[i] = 0;
        while((bn->size > 0) && (bn->d[bn->size] == 0))
            bn->size -= 1;
    }
    return bn;
}

//*** BnInitializeWord()
// This function is used to initialize an allocated bigNum with a word value. The
// bigNum does not have to be allocated with a single word.
LIB_EXPORT bigNum
BnInitializeWord(
    bigNum          bn,         // IN:
    crypt_uword_t   allocated,  // IN:
    crypt_uword_t   word        // IN:
    )
{
    bn->allocated = allocated;
    bn->size = (word != 0);
    bn->d[0] = word;
    while(allocated > 1)
        bn->d[--allocated] = 0;
    return bn;
}

//*** BnInit()
// This function initializes a stack allocated bignum_t. It initializes
// 'allocated' and 'size' and zeros the words of 'd'.
LIB_EXPORT bigNum
BnInit(
    bigNum               bn,
    crypt_uword_t        allocated
    )
{
    if(bn != NULL)
    {
        bn->allocated = allocated;
        bn->size = 0;
        while(allocated != 0)
            bn->d[--allocated] = 0;
    }
    return bn;
}

//*** BnCopy()
// Function to copy a bignum_t. If the output is NULL, then
// nothing happens. If the input is NULL, the output is set
// to zero.
LIB_EXPORT BOOL
BnCopy(
    bigNum           out,
    bigConst         in
    )
{
    if(in == out)
        BnSetTop(out, BnGetSize(out));
    else if(out != NULL)
    {
        if(in != NULL)
        {
            unsigned int         i;
            pAssert(BnGetAllocated(out) >= BnGetSize(in));
            for(i = 0; i < BnGetSize(in); i++)
                out->d[i] = in->d[i];
            BnSetTop(out, BnGetSize(in));
        }
        else
            BnSetTop(out, 0);
    }
    return TRUE;
}

#if ALG_ECC

//*** BnPointCopy()
// Function to copy a bn point.
LIB_EXPORT BOOL
BnPointCopy(
    bigPoint                 pOut,
    pointConst               pIn
    )
{
    return BnCopy(pOut->x, pIn->x)
        && BnCopy(pOut->y, pIn->y)
        && BnCopy(pOut->z, pIn->z);
}

//*** BnInitializePoint()
// This function is used to initialize a point structure with the addresses
// of the coordinates.
LIB_EXPORT bn_point_t *
BnInitializePoint(
    bigPoint             p,     // OUT: structure to receive pointers
    bigNum               x,     // IN: x coordinate
    bigNum               y,     // IN: y coordinate
    bigNum               z      // IN: x coordinate
    )
{
    p->x = x;
    p->y = y;
    p->z = z;
    BnSetWord(z, 1);
    return p;
}

#endif // ALG_ECC