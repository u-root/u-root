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
 *  list of conditions and the following disclaimer in the documentation and/or other
 *  materials provided with the distribution.
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
// This file contains the math functions that are not implemented in the BnMath
// library (yet). These math functions will call the ST MPA library or the
// LibTomCrypt library to execute the operations. Since the TPM internal big number
// format is identical to the MPA format, no reformatting is required.

//** Includes
#include "Tpm.h"

#ifdef MATH_LIB_LTC

#if defined ECC_NIST_P256 && ECC_NIST_P256 == YES && ECC_CURVE_COUNT > 1
#error "LibTomCrypt only supports P256"
#endif

//** Functions

//*** BnModMult()
// Does multiply and divide returning the remainder of the divide.
LIB_EXPORT BOOL
BnModMult(
    bigNum              result,
    bigConst            op1,
    bigConst            op2,
    bigConst            modulus
    )
{
    BN_VAR(temp, LARGEST_NUMBER_BITS * 2);
    // mpa_mul does not allocate from the pool if the result is not the same as
    // op1 or op2. since this is assured by the stack allocation of 'temp', the
    // pool pointer can be NULL
    pAssert(BnGetAllocated(result) >= BnGetSize(modulus));
    mpa_mul((mpanum)temp, (const mpanum)op1, (const mpanum)op2,
            NULL);
    return BnDiv(NULL, result, temp, modulus);
}

//*** BnMult()
// Multiplies two numbers
LIB_EXPORT BOOL
BnMult(
    bigNum               result,
    bigConst             multiplicand,
    bigConst             multiplier
    )
{
    // Make sure that the mpa_mul function does not allocate anything
    // from the POOL by eliminating the reason for doing it.
    BN_VAR(tempResult, LARGEST_NUMBER_BITS * 2);
    if(result != multiplicand && result != multiplier)
        tempResult = result;
    mpa_mul((mpanum)tempResult, (const mpanum)multiplicand,
            (const mpanum)multiplier,
            NULL);
    BnCopy(result, tempResult);
    return TRUE;
}

//*** BnDiv()
// This function divides two BIGNUM values. The function always returns TRUE.
LIB_EXPORT BOOL
BnDiv(
    bigNum               quotient,
    bigNum               remainder,
    bigConst             dividend,
    bigConst             divisor
    )
{
    MPA_ENTER(10, LARGEST_NUMBER_BITS);
    pAssert(!BnEqualZero(divisor));
    if(BnGetSize(dividend) < BnGetSize(divisor))
    {
        if(quotient)
            BnSetWord(quotient, 0);
        if(remainder)
            BnCopy(remainder, dividend);
    }
    else
    {
        pAssert((quotient == NULL)
                || (quotient->allocated >= 
                        (unsigned)(dividend->size - divisor->size)));
        pAssert((remainder == NULL)
                || (remainder->allocated >= divisor->size));
        mpa_div((mpanum)quotient, (mpanum)remainder,
                (const mpanum)dividend, (const mpanum)divisor, POOL);
    }
    MPA_LEAVE();
    return TRUE;
}

#ifdef TPM_ALG_RSA
//*** BnGcd()
// Get the greatest common divisor of two numbers
LIB_EXPORT BOOL
BnGcd(
    bigNum      gcd,            // OUT: the common divisor
    bigConst    number1,        // IN:
    bigConst    number2         // IN:
    )
{
    MPA_ENTER(20, LARGEST_NUMBER_BITS);
//
    mpa_gcd((mpanum)gcd, (mpanum)number1, (mpanum)number2, POOL);
    MPA_LEAVE();
    return TRUE;
}

//***BnModExp()
// Do modular exponentiation using BIGNUM values. The conversion from a bignum_t
// to a BIGNUM is trivial as they are based on the same structure
LIB_EXPORT BOOL
BnModExp(
    bigNum               result,         // OUT: the result
    bigConst             number,         // IN: number to exponentiate
    bigConst             exponent,       // IN:
    bigConst             modulus         // IN:
    )
{
    MPA_ENTER(20, LARGEST_NUMBER_BITS);
    BN_VAR(bnR, MAX_RSA_KEY_BITS);
    BN_VAR(bnR2, MAX_RSA_KEY_BITS);
    mpa_word_t              n_inv;
    mpa_word_t              ffmCtx[mpa_fmm_context_size_in_U32(MAX_RSA_KEY_BITS)];
//
    mpa_init_static_fmm_context((mpa_fmm_context_base *)ffmCtx,
                                BYTES_TO_CRYPT_WORDS(sizeof(ffmCtx)));
       // Generate modular form
    if(mpa_compute_fmm_context((const mpanum)modulus, (mpanum)bnR,
                               (mpanum)bnR2, &n_inv, POOL) != 0)
        FAIL(FATAL_ERROR_INTERNAL);
    // Do exponentiation
    mpa_exp_mod((mpanum)result, (const mpanum)number, (const mpanum)exponent,
                (const mpanum)modulus, (const mpanum)bnR, (const mpanum)bnR2,
                n_inv, POOL);
    MPA_LEAVE();
    return TRUE;
}

//*** BnModInverse()
// Modular multiplicative inverse
LIB_EXPORT BOOL
BnModInverse(
    bigNum               result,
    bigConst             number,
    bigConst             modulus
    )
{
    BOOL            retVal;
    MPA_ENTER(10, LARGEST_NUMBER_BITS);
    retVal = (mpa_inv_mod((mpanum)result, (const mpanum)number,
                          (const mpanum)modulus, POOL) == 0);
    MPA_LEAVE();
    return retVal;
}
#endif // TPM_ALG_RSA

#ifdef TPM_ALG_ECC


//*** BnEccModMult()
// This function does a point multiply of the form R = [d]S
// return type: BOOL
//  FALSE       failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: point to multiply by 'd'
    bigConst             d,         // IN: scalar for [d]S
    bigCurve             E
    )
{
    MPA_ENTER(30, MAX_ECC_KEY_BITS * 2);
    // The point multiply in LTC seems to need a large reciprocal for
    // intermediate results
    POINT_VAR(result, MAX_ECC_KEY_BITS * 4);
    BOOL                 OK;
//
    (POOL);     // Avoid compiler warning
    if(S == NULL)
        S = CurveGetG(AccessCurveData(E));
    OK = (ltc_ecc_mulmod((mpanum)d, (ecc_point *)S,
                         (ecc_point *)result, (void *)CurveGetPrime(E), 1)
          == CRYPT_OK);
    OK = OK && !BnEqualZero(result->z);
    if(OK)
        BnPointCopy(R, result);

    MPA_LEAVE();
    return OK ? TPM_RC_SUCCESS : TPM_RC_NO_RESULT;
}

//*** BnEccModMult2()
// This function does a point multiply of the form R = [d]S + [u]Q
// return type: BOOL
//  FALSE       failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult2(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: first point (optional)
    bigConst             d,         // IN: scalar for [d]S or [d]G
    pointConst           Q,         // IN: second point
    bigConst             u,         // IN: second scalar
    bigCurve             E          // IN: curve
    )
{
    MPA_ENTER(80, MAX_ECC_KEY_BITS);
    BOOL                 OK;
    // The point multiply in LTC seems to need a large reciprocal for
    // intermediate results
    POINT_VAR(result, MAX_ECC_KEY_BITS * 4);
//
    (POOL);     // Avoid compiler warning
    if(S == NULL)
        S = CurveGetG(AccessCurveData(E));

    OK = (ltc_ecc_mul2add((ecc_point  *)S, (mpanum)d, (ecc_point  *)Q, (mpanum)u,
                          (ecc_point  *)result, (mpanum)CurveGetPrime(E))
          == CRYPT_OK);
    OK = OK && !BnEqualZero(result->z);

    if(OK)
        BnPointCopy(R, result);

    MPA_LEAVE();
    return OK ? TPM_RC_SUCCESS : TPM_RC_NO_RESULT;
}

//*** BnEccAdd()
// This function does addition of two points. Since this is not implemented
// in LibTomCrypt() will try to trick it by doing multiply with scalar of 1.
// I have no idea if this will work and it's not needed unless MQV or the SM2
// variant is enabled.
// return type: BOOL
//  FALSE       failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccAdd(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: point to multiply by 'd'
    pointConst           Q,         // IN: second point
    bigCurve             E          // IN: curve
    )
{
    BN_WORD_INITIALIZED(one, 1);
    return BnEccModMult2(R, S, one, Q, one, E);
}

#endif // TPM_ALG_ECC

#endif // MATH_LIB_LTC
