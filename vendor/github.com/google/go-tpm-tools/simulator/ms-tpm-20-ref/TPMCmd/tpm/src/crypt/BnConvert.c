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
// This file contains the basic conversion functions that will convert TPM2B
// to/from the internal format. The internal format is a bigNum,
//

//** Includes

#include "Tpm.h"

//** Functions

//*** BnFromBytes()
// This function will convert a big-endian byte array to the internal number
// format. If bn is NULL, then the output is NULL. If bytes is null or the 
// required size is 0, then the output is set to zero
LIB_EXPORT bigNum
BnFromBytes(
    bigNum           bn,
    const BYTE      *bytes,
    NUMBYTES         nBytes
    )
{
    const BYTE      *pFrom; // 'p' points to the least significant bytes of source
    BYTE            *pTo;   // points to least significant bytes of destination
    crypt_uword_t    size;
//

    size = (bytes != NULL) ? BYTES_TO_CRYPT_WORDS(nBytes) : 0;

    // If nothing in, nothing out
    if(bn == NULL)
        return NULL;

    // make sure things fit
    pAssert(BnGetAllocated(bn) >= size);

    if(size > 0)
    {
        // Clear the topmost word in case it is not filled with data
        bn->d[size - 1] = 0;
        // Moving the input bytes from the end of the list (LSB) end
        pFrom = bytes + nBytes - 1;
        // To the LS0 of the LSW of the bigNum.
        pTo = (BYTE *)bn->d;
        for(; nBytes != 0; nBytes--)
            *pTo++ = *pFrom--;
        // For a little-endian machine, the conversion is a straight byte
        // reversal. For a big-endian machine, we have to put the words in
        // big-endian byte order
#if BIG_ENDIAN_TPM
        {
            crypt_word_t   t;
            for(t = (crypt_word_t)size - 1; t >= 0; t--)
                bn->d[t] = SWAP_CRYPT_WORD(bn->d[t]);
        }
#endif
    }
    BnSetTop(bn, size);
    return bn;
}

//*** BnFrom2B()
// Convert an TPM2B to a BIG_NUM.
// If the input value does not exist, or the output does not exist, or the input
// will not fit into the output the function returns NULL
LIB_EXPORT bigNum
BnFrom2B(
    bigNum           bn,         // OUT:
    const TPM2B     *a2B         // IN: number to convert
    )
{
    if(a2B != NULL)
        return BnFromBytes(bn, a2B->buffer, a2B->size);
    // Make sure that the number has an initialized value rather than whatever
    // was there before
    BnSetTop(bn, 0);    // Function accepts NULL
    return NULL;
}

//*** BnFromHex()
// Convert a hex string into a bigNum. This is primarily used in debugging.
LIB_EXPORT bigNum
BnFromHex(
    bigNum          bn,         // OUT:
    const char      *hex        // IN:
    )
{
#define FromHex(a)  ((a) - (((a) > 'a') ? ('a' + 10)                               \
                                       : ((a) > 'A') ? ('A' - 10) : '0'))
    unsigned             i;
    unsigned             wordCount;
    const char          *p;
    BYTE                *d = (BYTE *)&(bn->d[0]);
//
    pAssert(bn && hex);
    i = (unsigned)strlen(hex);
    wordCount = BYTES_TO_CRYPT_WORDS((i + 1) / 2);
    if((i == 0) || (wordCount >= BnGetAllocated(bn)))
        BnSetWord(bn, 0);
    else
    {
        bn->d[wordCount - 1] = 0;
        p = hex + i - 1;
        for(;i > 1; i -= 2)
        {
            BYTE a;
            a = FromHex(*p);
            p--;
            *d++ = a + (FromHex(*p) << 4);
            p--;
        }
        if(i == 1)
            *d = FromHex(*p);
    }
#if !BIG_ENDIAN_TPM
    for(i = 0; i < wordCount; i++)
        bn->d[i] = SWAP_CRYPT_WORD(bn->d[i]);
#endif // BIG_ENDIAN_TPM
    BnSetTop(bn, wordCount);
    return bn;
}

//*** BnToBytes()
// This function converts a BIG_NUM to a byte array. It converts the bigNum to a 
// big-endian byte string and sets 'size' to the normalized value. If  'size' is an 
// input 0, then the receiving buffer is guaranteed to be large enough for the result 
// and the size will be set to the size required for bigNum (leading zeros 
// suppressed).
//
// The conversion for a little-endian machine simply requires that all significant
// bytes of the bigNum be reversed. For a big-endian machine, rather than 
// unpack each word individually, the bigNum is converted to little-endian words,
// copied, and then converted back to big-endian.
LIB_EXPORT BOOL
BnToBytes(
    bigConst             bn,
    BYTE                *buffer,
    NUMBYTES            *size           // This the number of bytes that are
                                        // available in the buffer. The result
                                        // should be this big.
    )
{
    crypt_uword_t        requiredSize;
    BYTE                *pFrom;
    BYTE                *pTo;
    crypt_uword_t        count;
//
    // validate inputs
    pAssert(bn && buffer && size);

    requiredSize = (BnSizeInBits(bn) + 7) / 8;
    if(requiredSize == 0)
    {
        // If the input value is 0, return a byte of zero
        *size = 1;
        *buffer = 0;
    }
    else
    {
#if BIG_ENDIAN_TPM
        // Copy the constant input value into a modifiable value
        BN_VAR(bnL, LARGEST_NUMBER_BITS * 2);
        BnCopy(bnL, bn);
        // byte swap the words in the local value to make them little-endian
        for(count = 0; count < bnL->size; count++)
            bnL->d[count] = SWAP_CRYPT_WORD(bnL->d[count]);
        bn = (bigConst)bnL;
#endif
        if(*size == 0)
            *size = (NUMBYTES)requiredSize;
        pAssert(requiredSize <= *size);
        // Byte swap the number (not words but the whole value)
        count = *size;
        // Start from the least significant word and offset to the most significant
        // byte which is in some high word
        pFrom = (BYTE *)(&bn->d[0]) + requiredSize - 1;
        pTo = buffer;

        // If the number of output bytes is larger than the number bytes required
        // for the input number, pad with zeros
        for(count = *size; count > requiredSize; count--)
            *pTo++ = 0;
        // Move the most significant byte at the end of the BigNum to the next most
        // significant byte position of the 2B and repeat for all significant bytes.
        for(; requiredSize > 0; requiredSize--)
            *pTo++ = *pFrom--;
    }
    return TRUE;
}

//*** BnTo2B()
// Function to convert a BIG_NUM to TPM2B.
// The TPM2B size is set to the requested 'size' which may require padding.
// If 'size' is non-zero and less than required by the value in 'bn' then an error
// is returned. If 'size' is zero, then the TPM2B is assumed to be large enough
// for the data and a2b->size will be adjusted accordingly.
LIB_EXPORT BOOL
BnTo2B(
    bigConst         bn,                // IN:
    TPM2B           *a2B,               // OUT:
    NUMBYTES         size               // IN: the desired size
    )
{
    // Set the output size
    if(bn && a2B)
    {
        a2B->size = size;
        return BnToBytes(bn, a2B->buffer, &a2B->size);
    }
    return FALSE;
}

#if ALG_ECC

//*** BnPointFrom2B()
// Function to create a BIG_POINT structure from a 2B point.
// A point is going to be two ECC values in the same buffer. The values are going
// to be the size of the modulus.  They are in modular form.
LIB_EXPORT bn_point_t   *
BnPointFrom2B(
    bigPoint             ecP,         // OUT: the preallocated point structure
    TPMS_ECC_POINT      *p            // IN: the number to convert
    )
{
    if(p == NULL)
        return NULL;

    if(NULL != ecP)
    {
        BnFrom2B(ecP->x, &p->x.b);
        BnFrom2B(ecP->y, &p->y.b);
        BnSetWord(ecP->z, 1);
    }
    return ecP;
}

//*** BnPointTo2B()
// This function converts a BIG_POINT into a TPMS_ECC_POINT. A TPMS_ECC_POINT
// contains two TPM2B_ECC_PARAMETER values. The maximum size of the parameters
// is dependent on the maximum EC key size used in an implementation.
// The presumption is that the TPMS_ECC_POINT is large enough to hold 2 TPM2B
// values, each as large as a MAX_ECC_PARAMETER_BYTES
LIB_EXPORT BOOL
BnPointTo2B(
    TPMS_ECC_POINT  *p,             // OUT: the converted 2B structure
    bigPoint         ecP,           // IN: the values to be converted
    bigCurve         E              // IN: curve descriptor for the point
    )
{
    UINT16           size;
//
    pAssert(p && ecP && E);
    pAssert(BnEqualWord(ecP->z, 1));
    // BnMsb is the bit number of the MSB. This is one less than the number of bits
    size = (UINT16)BITS_TO_BYTES(BnSizeInBits(CurveGetOrder(AccessCurveData(E))));
    BnTo2B(ecP->x, &p->x.b, size);
    BnTo2B(ecP->y, &p->y.b, size);
    return TRUE;
}

#endif // ALG_ECC