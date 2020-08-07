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
// library (yet). These math functions will call the wolfcrypt library to execute
// the operations. There is a difference between the internal format and the
// wolfcrypt format. To call the wolfcrypt function, a mp_int structure is created
// for each passed variable. We define USE_FAST_MATH wolfcrypt option, which allocates
// mp_int on the stack. We must copy each word to the new structure, and set the used
// size. 
//
// Not using USE_FAST_MATH would allow for a simple pointer swap for the big integer
// buffer 'd', however wolfcrypt expects to manage this memory, and will swap out
// the pointer to and from temporary variables and free the reference underneath us.
// Using USE_FAST_MATH also instructs wolfcrypt to use the stack for all these 
// intermediate variables


//** Includes and Defines
#include "Tpm.h"

#ifdef MATH_LIB_WOLF
#include "BnConvert_fp.h"
#include "TpmToWolfMath_fp.h"

#define WOLF_HALF_RADIX     (RADIX_BITS == 64 && !defined(FP_64BIT))

//** Functions

//*** BnFromWolf()
// This function converts a wolfcrypt mp_int to a TPM bignum. In this implementation
// it is assumed that wolfcrypt used the same format for a big number as does the
// TPM -- an array of native-endian words in little-endian order.
void
BnFromWolf(
    bigNum          bn,
    mp_int          *wolfBn
    )
{
    if(bn != NULL)
    {
        int         i;
#if WOLF_HALF_RADIX
        pAssert((unsigned)wolfBn->used <= 2 * BnGetAllocated(bn));
#else
        pAssert((unsigned)wolfBn->used <= BnGetAllocated(bn));
#endif
        for (i = 0; i < wolfBn->used; i++)
        {
#if WOLF_HALF_RADIX
            if (i & 1)
                bn->d[i/2] |= (crypt_uword_t)wolfBn->dp[i] << 32;
            else
                bn->d[i/2] = wolfBn->dp[i];
#else
            bn->d[i] = wolfBn->dp[i];
#endif
        }

#if WOLF_HALF_RADIX
        BnSetTop(bn, (wolfBn->used + 1)/2);
#else
        BnSetTop(bn, wolfBn->used);
#endif
    }
}

//*** BnToWolf()
// This function converts a TPM bignum to a wolfcrypt mp_init, and has the same
// assumptions as made by BnFromWolf()
void
BnToWolf(
    mp_int              *toInit,
    bigConst            initializer
    )
{
    uint32_t         i;
    if (toInit != NULL && initializer != NULL)
    {
        for (i = 0; i < initializer->size; i++)
        {
#if WOLF_HALF_RADIX
            toInit->dp[2 * i] = (fp_digit)initializer->d[i];
            toInit->dp[2 * i + 1] = (fp_digit)(initializer->d[i] >> 32);
#else
            toInit->dp[i] = initializer->d[i];
#endif
        }

#if WOLF_HALF_RADIX
        toInit->used = (int)initializer->size * 2;
        if (toInit->dp[toInit->used - 1] == 0 && toInit->dp[toInit->used - 2] != 0)
            --toInit->used;
#else
        toInit->used = (int)initializer->size;
#endif
        toInit->sign = 0;
    }
}

//*** MpInitialize()
// This function initializes an wolfcrypt mp_int.
mp_int *
MpInitialize(
    mp_int              *toInit
)
{
    mp_init( toInit );
    return toInit;
}

#if LIBRARY_COMPATIBILITY_CHECK
//** MathLibraryCompatibililtyCheck()
// This function is only used during development to make sure that the library
// that is being referenced is using the same size of data structures as the TPM.
void
MathLibraryCompatibilityCheck(
    void 
    )
{
    BN_VAR(tpmTemp, 64 * 8); // allocate some space for a test value
    crypt_uword_t           i;
    TPM2B_TYPE(TEST, 16);
    TPM2B_TEST              test = {{16, {0x0F, 0x0E, 0x0D, 0x0C, 
                                          0x0B, 0x0A, 0x09, 0x08, 
                                          0x07, 0x06, 0x05, 0x04, 
                                          0x03, 0x02, 0x01, 0x00}}};
    // Convert the test TPM2B to a bigNum
    BnFrom2B(tpmTemp, &test.b);
    MP_INITIALIZED(wolfTemp, tpmTemp);
    (wolfTemp); // compiler warning
    // Make sure the values are consistent
    cAssert(wolfTemp->used == (int)tpmTemp->size);
    for(i = 0; i < tpmTemp->size; i++)
        cAssert(wolfTemp->dp[i] == tpmTemp->d[i]);
}
#endif

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
    WOLF_ENTER();
    BOOL                OK;
    MP_INITIALIZED(bnOp1, op1);
    MP_INITIALIZED(bnOp2, op2);
    MP_INITIALIZED(bnTemp, NULL);
    BN_VAR(temp, LARGEST_NUMBER_BITS * 2);

    pAssert(BnGetAllocated(result) >= BnGetSize(modulus));

    OK = (mp_mul( bnOp1, bnOp2, bnTemp ) == MP_OKAY);
    if(OK)
    {
        BnFromWolf(temp, bnTemp);
        OK = BnDiv(NULL, result, temp, modulus);
    }

    WOLF_LEAVE();
    return OK;
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
    WOLF_ENTER();
    BOOL                OK;
    MP_INITIALIZED(bnTemp, NULL);
    MP_INITIALIZED(bnA, multiplicand);
    MP_INITIALIZED(bnB, multiplier);

    pAssert(result->allocated >=
            (BITS_TO_CRYPT_WORDS(BnSizeInBits(multiplicand)
                                 + BnSizeInBits(multiplier))));

    OK = (mp_mul( bnA, bnB, bnTemp ) == MP_OKAY);
    if(OK)
    {
        BnFromWolf(result, bnTemp);
    }

    WOLF_LEAVE();
    return OK;
}

//*** BnDiv()
// This function divides two bigNum values. The function returns FALSE if
// there is an error in the operation.
LIB_EXPORT BOOL
BnDiv(
    bigNum               quotient,
    bigNum               remainder,
    bigConst             dividend,
    bigConst             divisor
    )
{
    WOLF_ENTER();
    BOOL        OK;
    MP_INITIALIZED(bnQ, quotient);
    MP_INITIALIZED(bnR, remainder);
    MP_INITIALIZED(bnDend, dividend);
    MP_INITIALIZED(bnSor, divisor);
    pAssert(!BnEqualZero(divisor));
    if(BnGetSize(dividend) < BnGetSize(divisor))
    {
        if(quotient)
            BnSetWord(quotient, 0);
        if(remainder)
            BnCopy(remainder, dividend);
        OK = TRUE;
    }
    else
    {
        pAssert((quotient == NULL)
                || (quotient->allocated >= (unsigned)(dividend->size 
                                                      - divisor->size)));
        pAssert((remainder == NULL)
                || (remainder->allocated >= divisor->size));
        OK = (mp_div(bnDend , bnSor, bnQ, bnR) == MP_OKAY);
        if(OK)
        {
            BnFromWolf(quotient, bnQ);
            BnFromWolf(remainder, bnR);
        }
    }

    WOLF_LEAVE();
    return OK;
}

#if ALG_RSA
//*** BnGcd()
// Get the greatest common divisor of two numbers
LIB_EXPORT BOOL
BnGcd(
    bigNum      gcd,            // OUT: the common divisor
    bigConst    number1,        // IN:
    bigConst    number2         // IN:
    )
{
    WOLF_ENTER();
    BOOL            OK;
    MP_INITIALIZED(bnGcd, gcd);
    MP_INITIALIZED(bn1, number1);
    MP_INITIALIZED(bn2, number2);
    pAssert(gcd != NULL);
    OK = (mp_gcd( bn1, bn2, bnGcd ) == MP_OKAY);
    if(OK)
    {
        BnFromWolf(gcd, bnGcd);
    }
    WOLF_LEAVE();
    return OK;
}

//***BnModExp()
// Do modular exponentiation using bigNum values. The conversion from a mp_int to
// a bigNum is trivial as they are based on the same structure
LIB_EXPORT BOOL
BnModExp(
    bigNum               result,         // OUT: the result
    bigConst             number,         // IN: number to exponentiate
    bigConst             exponent,       // IN:
    bigConst             modulus         // IN:
    )
{
    WOLF_ENTER();
    BOOL            OK;
    MP_INITIALIZED(bnResult, result);
    MP_INITIALIZED(bnN, number);
    MP_INITIALIZED(bnE, exponent);
    MP_INITIALIZED(bnM, modulus);
    OK = (mp_exptmod( bnN, bnE, bnM, bnResult ) == MP_OKAY);
    if(OK)
    {
        BnFromWolf(result, bnResult);
    }

    WOLF_LEAVE();
    return OK;
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
    WOLF_ENTER();
    BOOL            OK;
    MP_INITIALIZED(bnResult, result);
    MP_INITIALIZED(bnN, number);
    MP_INITIALIZED(bnM, modulus);

    OK = (mp_invmod(bnN, bnM, bnResult) == MP_OKAY);
    if(OK)
    {
        BnFromWolf(result, bnResult);
    }

    WOLF_LEAVE();
    return OK;
}
#endif // TPM_ALG_RSA

#if ALG_ECC

//*** PointFromWolf()
// Function to copy the point result from a wolf ecc_point to a bigNum
void
PointFromWolf(
    bigPoint         pOut,      // OUT: resulting point
    ecc_point       *pIn       // IN: the point to return
    )
{
    BnFromWolf(pOut->x, pIn->x);
    BnFromWolf(pOut->y, pIn->y);
    BnFromWolf(pOut->z, pIn->z);
}

//*** PointToWolf()
// Function to copy the point result from a bigNum to a wolf ecc_point
void
PointToWolf(
    ecc_point      *pOut,      // OUT: resulting point
    pointConst      pIn       // IN: the point to return
    )
{
    BnToWolf(pOut->x, pIn->x);
    BnToWolf(pOut->y, pIn->y);
    BnToWolf(pOut->z, pIn->z);
}

//*** EcPointInitialized()
// Allocate and initialize a point.
static ecc_point *
EcPointInitialized(
    pointConst          initializer
    )
{
    ecc_point           *P;

    P = wc_ecc_new_point();
    pAssert(P != NULL);
    // mp_int x,y,z are stack allocated.
    // initializer is not required
    if (P != NULL && initializer != NULL)
    {
        PointToWolf( P, initializer );
    }

    return P;
}

//*** BnEccModMult()
// This function does a point multiply of the form R = [d]S
// return type: BOOL
//  FALSE       failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: point to multiply by 'd' (optional)
    bigConst             d,         // IN: scalar for [d]S
    bigCurve             E
    )
{
    WOLF_ENTER();
    BOOL                 OK;
    MP_INITIALIZED(bnD, d);
    MP_INITIALIZED(bnPrime, CurveGetPrime(E));
    POINT_CREATE(pS, NULL);
    POINT_CREATE(pR, NULL);

    if(S == NULL)
        S = CurveGetG(AccessCurveData(E));

    PointToWolf(pS, S);

    OK = (wc_ecc_mulmod(bnD, pS, pR, NULL, bnPrime, 1 ) == MP_OKAY);
    if(OK)
    {
        PointFromWolf(R, pR);
    }

    POINT_DELETE(pR);
    POINT_DELETE(pS);

    WOLF_LEAVE();
    return !BnEqualZero(R->z);
}

//*** BnEccModMult2()
// This function does a point multiply of the form R = [d]G + [u]Q
// return type: BOOL
//  FALSE       failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult2(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: optional point
    bigConst             d,         // IN: scalar for [d]S or [d]G
    pointConst           Q,         // IN: second point
    bigConst             u,         // IN: second scalar
    bigCurve             E          // IN: curve
    )
{
    WOLF_ENTER();
    BOOL                 OK;
    POINT_CREATE(pR, NULL);
    POINT_CREATE(pS, NULL);
    POINT_CREATE(pQ, Q);
    MP_INITIALIZED(bnD, d);
    MP_INITIALIZED(bnU, u);
    MP_INITIALIZED(bnPrime, CurveGetPrime(E));
    MP_INITIALIZED(bnA, CurveGet_a(E));

    if(S == NULL)
        S = CurveGetG(AccessCurveData(E));
    PointToWolf( pS, S );

    OK = (ecc_mul2add(pS, bnD, pQ, bnU, pR, bnA, bnPrime, NULL) == MP_OKAY);
    if(OK)
    {
        PointFromWolf(R, pR);
    }

    POINT_DELETE(pS);
    POINT_DELETE(pQ);
    POINT_DELETE(pR);

    WOLF_LEAVE();
    return !BnEqualZero(R->z);
}

//** BnEccAdd()
// This function does addition of two points.
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
    WOLF_ENTER();
    BOOL                 OK;
    mp_digit             mp;
    POINT_CREATE(pR, NULL);
    POINT_CREATE(pS, S);
    POINT_CREATE(pQ, Q);
    MP_INITIALIZED(bnA, CurveGet_a(E));
    MP_INITIALIZED(bnMod, CurveGetPrime(E));
//
    OK = (mp_montgomery_setup(bnMod, &mp) == MP_OKAY);
    OK = OK && (ecc_projective_add_point(pS, pQ, pR, bnA, bnMod, mp ) == MP_OKAY);
    if(OK)
    {
        PointFromWolf(R, pR);
    }

    POINT_DELETE(pS);
    POINT_DELETE(pQ);
    POINT_DELETE(pR);

    WOLF_LEAVE();
    return !BnEqualZero(R->z);
}

#endif // TPM_ALG_ECC

#endif // MATH_LIB_WOLF