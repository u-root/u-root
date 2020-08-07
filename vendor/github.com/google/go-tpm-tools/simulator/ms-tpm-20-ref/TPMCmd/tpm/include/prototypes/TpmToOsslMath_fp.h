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
 *  Date: Apr  2, 2019  Time: 03:18:00PM
 */

#ifndef    _TPM_TO_OSSL_MATH_FP_H_
#define    _TPM_TO_OSSL_MATH_FP_H_

#ifdef MATH_LIB_OSSL

//*** OsslToTpmBn()
// This function converts an OpenSSL BIGNUM to a TPM bignum. In this implementation
// it is assumed that OpenSSL uses a different control structure but the same data
// layout -- an array of native-endian words in little-endian order.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure because value will not fit or OpenSSL variable doesn't
//                      exist
BOOL
OsslToTpmBn(
    bigNum          bn,
    BIGNUM          *osslBn
);

//*** BigInitialized()
// This function initializes an OSSL BIGNUM from a TPM bigConst. Do not use this for
// values that are passed to OpenSLL when they are not declared as const in the
// function prototype. Instead, use BnNewVariable().
BIGNUM *
BigInitialized(
    BIGNUM             *toInit,
    bigConst            initializer
);
#if LIBRARY_COMPATIBILITY_CHECK

//*** MathLibraryCompatibilityCheck()
void
MathLibraryCompatibilityCheck(
    void
);
#endif

//*** BnModMult()
// This function does a modular multiply. It first does a multiply and then a divide
// and returns the remainder of the divide.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnModMult(
    bigNum              result,
    bigConst            op1,
    bigConst            op2,
    bigConst            modulus
);

//*** BnMult()
// Multiplies two numbers
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnMult(
    bigNum               result,
    bigConst             multiplicand,
    bigConst             multiplier
);

//*** BnDiv()
// This function divides two bigNum values. The function returns FALSE if
// there is an error in the operation.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnDiv(
    bigNum               quotient,
    bigNum               remainder,
    bigConst             dividend,
    bigConst             divisor
);

#if ALG_RSA
//*** BnGcd()
// Get the greatest common divisor of two numbers
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnGcd(
    bigNum      gcd,            // OUT: the common divisor
    bigConst    number1,        // IN:
    bigConst    number2         // IN:
);

//***BnModExp()
// Do modular exponentiation using bigNum values. The conversion from a bignum_t to
// a bigNum is trivial as they are based on the same structure
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnModExp(
    bigNum               result,         // OUT: the result
    bigConst             number,         // IN: number to exponentiate
    bigConst             exponent,       // IN:
    bigConst             modulus         // IN:
);

//*** BnModInverse()
// Modular multiplicative inverse
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation
LIB_EXPORT BOOL
BnModInverse(
    bigNum               result,
    bigConst             number,
    bigConst             modulus
);
#endif // ALG_RSA
#if ALG_ECC

//*** BnCurveInitialize()
// This function initializes the OpenSSL curve information structure. This
// structure points to the TPM-defined values for the curve, to the context for the
// number values in the frame, and to the OpenSSL-defined group values.
//  Return Type: bigCurve *
//      NULL        the TPM_ECC_CURVE is not valid or there was a problem in
//                  in initializing the curve data
//      non-NULL    points to 'E'
LIB_EXPORT bigCurve
BnCurveInitialize(
    bigCurve          E,           // IN: curve structure to initialize
    TPM_ECC_CURVE     curveId      // IN: curve identifier
);

//*** BnCurveFree()
// This function will free the allocated components of the curve and end the
// frame in which the curve data exists
LIB_EXPORT void
BnCurveFree(
    bigCurve                    E
);

//*** BnEccModMult()
// This function does a point multiply of the form R = [d]S
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: point to multiply by 'd' (optional)
    bigConst             d,         // IN: scalar for [d]S
    bigCurve             E
);

//*** BnEccModMult2()
// This function does a point multiply of the form R = [d]G + [u]Q
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccModMult2(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: optional point
    bigConst             d,         // IN: scalar for [d]S or [d]G
    pointConst           Q,         // IN: second point
    bigConst             u,         // IN: second scalar
    bigCurve             E          // IN: curve
);

//** BnEccAdd()
// This function does addition of two points.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure in operation; treat as result being point at infinity
LIB_EXPORT BOOL
BnEccAdd(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: point to multiply by 'd'
    pointConst           Q,         // IN: second point
    bigCurve             E          // IN: curve
);
#endif // ALG_ECC
#endif // MATHLIB OSSL

#endif  // _TPM_TO_OSSL_MATH_FP_H_
