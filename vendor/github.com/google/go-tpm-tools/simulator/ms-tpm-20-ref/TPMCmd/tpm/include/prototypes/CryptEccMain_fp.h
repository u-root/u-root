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

#ifndef    _CRYPT_ECC_MAIN_FP_H_
#define    _CRYPT_ECC_MAIN_FP_H_

#if ALG_ECC

//** Functions
#if SIMULATION
void
EccSimulationEnd(
    void
);
#endif // SIMULATION

//*** CryptEccInit()
// This function is called at _TPM_Init
BOOL
CryptEccInit(
    void
);

//*** CryptEccStartup()
// This function is called at TPM2_Startup().
BOOL
CryptEccStartup(
    void
);

//*** ClearPoint2B(generic)
// Initialize the size values of a TPMS_ECC_POINT structure.
void
ClearPoint2B(
    TPMS_ECC_POINT      *p          // IN: the point
);

//*** CryptEccGetParametersByCurveId()
// This function returns a pointer to the curve data that is associated with
// the indicated curveId.
// If there is no curve with the indicated ID, the function returns NULL. This
// function is in this module so that it can be called by GetCurve data.
//  Return Type: const ECC_CURVE_DATA
//      NULL            curve with the indicated TPM_ECC_CURVE is not implemented
//      != NULL         pointer to the curve data
LIB_EXPORT const ECC_CURVE *
CryptEccGetParametersByCurveId(
    TPM_ECC_CURVE       curveId     // IN: the curveID
);

//*** CryptEccGetKeySizeForCurve()
// This function returns the key size in bits of the indicated curve.
LIB_EXPORT UINT16
CryptEccGetKeySizeForCurve(
    TPM_ECC_CURVE            curveId    // IN: the curve
);

//*** GetCurveData()
// This function returns the a pointer for the parameter data
// associated with a curve.
const ECC_CURVE_DATA *
GetCurveData(
    TPM_ECC_CURVE        curveId     // IN: the curveID
);

//***CryptEccGetOID()
const BYTE *
CryptEccGetOID(
    TPM_ECC_CURVE       curveId
);

//*** CryptEccGetCurveByIndex()
// This function returns the number of the 'i'-th implemented curve. The normal
// use would be to call this function with 'i' starting at 0. When the 'i' is greater
// than or equal to the number of implemented curves, TPM_ECC_NONE is returned.
LIB_EXPORT TPM_ECC_CURVE
CryptEccGetCurveByIndex(
    UINT16               i
);

//*** CryptEccGetParameter()
// This function returns an ECC curve parameter. The parameter is
// selected by a single character designator from the set of ""PNABXYH"".
//  Return Type: BOOL
//      TRUE(1)         curve exists and parameter returned
//      FALSE(0)        curve does not exist or parameter selector
LIB_EXPORT BOOL
CryptEccGetParameter(
    TPM2B_ECC_PARAMETER     *out,       // OUT: place to put parameter
    char                     p,         // IN: the parameter selector
    TPM_ECC_CURVE            curveId    // IN: the curve id
);

//*** CryptCapGetECCCurve()
// This function returns the list of implemented ECC curves.
//  Return Type: TPMI_YES_NO
//      YES             if no more ECC curve is available
//      NO              if there are more ECC curves not reported
TPMI_YES_NO
CryptCapGetECCCurve(
    TPM_ECC_CURVE    curveID,       // IN: the starting ECC curve
    UINT32           maxCount,      // IN: count of returned curves
    TPML_ECC_CURVE  *curveList      // OUT: ECC curve list
);

//*** CryptGetCurveSignScheme()
// This function will return a pointer to the scheme of the curve.
const TPMT_ECC_SCHEME *
CryptGetCurveSignScheme(
    TPM_ECC_CURVE    curveId        // IN: The curve selector
);

//*** CryptGenerateR()
// This function computes the commit random value for a split signing scheme.
//
// If 'c' is NULL, it indicates that 'r' is being generated
// for TPM2_Commit.
// If 'c' is not NULL, the TPM will validate that the 'gr.commitArray'
// bit associated with the input value of 'c' is SET. If not, the TPM
// returns FALSE and no 'r' value is generated.
//  Return Type: BOOL
//      TRUE(1)         r value computed
//      FALSE(0)        no r value computed
BOOL
CryptGenerateR(
    TPM2B_ECC_PARAMETER     *r,             // OUT: the generated random value
    UINT16                  *c,             // IN/OUT: count value.
    TPMI_ECC_CURVE           curveID,       // IN: the curve for the value
    TPM2B_NAME              *name           // IN: optional name of a key to
                                            //     associate with 'r'
);

//*** CryptCommit()
// This function is called when the count value is committed. The 'gr.commitArray'
// value associated with the current count value is SET and g_commitCounter is
// incremented. The low-order 16 bits of old value of the counter is returned.
UINT16
CryptCommit(
    void
);

//*** CryptEndCommit()
// This function is called when the signing operation using the committed value
// is completed. It clears the gr.commitArray bit associated with the count
// value so that it can't be used again.
void
CryptEndCommit(
    UINT16           c              // IN: the counter value of the commitment
);

//*** CryptEccGetParameters()
// This function returns the ECC parameter details of the given curve.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        unsupported ECC curve ID
BOOL
CryptEccGetParameters(
    TPM_ECC_CURVE                curveId,       // IN: ECC curve ID
    TPMS_ALGORITHM_DETAIL_ECC   *parameters     // OUT: ECC parameters
);

//*** BnGetCurvePrime()
// This function is used to get just the prime modulus associated with a curve.
const bignum_t *
BnGetCurvePrime(
    TPM_ECC_CURVE            curveId
);

//*** BnGetCurveOrder()
// This function is used to get just the curve order
const bignum_t *
BnGetCurveOrder(
    TPM_ECC_CURVE            curveId
);

//*** BnIsOnCurve()
// This function checks if a point is on the curve.
BOOL
BnIsOnCurve(
    pointConst                   Q,
    const ECC_CURVE_DATA        *C
);

//*** BnIsValidPrivateEcc()
// Checks that 0 < 'x' < 'q'
BOOL
BnIsValidPrivateEcc(
    bigConst                 x,         // IN: private key to check
    bigCurve                 E          // IN: the curve to check
);

LIB_EXPORT BOOL
CryptEccIsValidPrivateKey(
    TPM2B_ECC_PARAMETER     *d,
    TPM_ECC_CURVE            curveId
);

//*** BnPointMul()
// This function does a point multiply of the form 'R' = ['d']'S' + ['u']'Q' where the
// parameters are bigNum values. If 'S' is NULL and d is not NULL, then it computes
// 'R' = ['d']'G' + ['u']'Q'  or just 'R' = ['d']'G' if 'u' and 'Q' are NULL.
// If 'skipChecks' is TRUE, then the function will not verify that the inputs are
// correct for the domain. This would be the case when the values were created by the
// CryptoEngine code.
// It will return TPM_RC_NO_RESULT if the resulting point is the point at infinity.
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        result of multiplication is a point at infinity
//      TPM_RC_ECC_POINT        'S' or 'Q' is not on the curve
//      TPM_RC_VALUE            'd' or 'u' is not < n
TPM_RC
BnPointMult(
    bigPoint             R,         // OUT: computed point
    pointConst           S,         // IN: optional point to multiply by 'd'
    bigConst             d,         // IN: scalar for [d]S or [d]G
    pointConst           Q,         // IN: optional second point
    bigConst             u,         // IN: optional second scalar
    bigCurve             E          // IN: curve parameters
);

//***BnEccGetPrivate()
// This function gets random values that are the size of the key plus 64 bits. The
// value is reduced (mod ('q' - 1)) and incremented by 1 ('q' is the order of the
// curve. This produces a value ('d') such that 1 <= 'd' < 'q'. This is the method
// of FIPS 186-4 Section B.4.1 ""Key Pair Generation Using Extra Random Bits"".
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure generating private key
BOOL
BnEccGetPrivate(
    bigNum                   dOut,      // OUT: the qualified random value
    const ECC_CURVE_DATA    *C,         // IN: curve for which the private key
                                        //     needs to be appropriate
    RAND_STATE              *rand       // IN: state for DRBG
);

//*** BnEccGenerateKeyPair()
// This function gets a private scalar from the source of random bits and does
// the point multiply to get the public key.
BOOL
BnEccGenerateKeyPair(
    bigNum               bnD,            // OUT: private scalar
    bn_point_t          *ecQ,            // OUT: public point
    bigCurve             E,              // IN: curve for the point
    RAND_STATE          *rand            // IN: DRBG state to use
);

//***CryptEccNewKeyPair(***)
// This function creates an ephemeral ECC. It is ephemeral in that
// is expected that the private part of the key will be discarded
LIB_EXPORT TPM_RC
CryptEccNewKeyPair(
    TPMS_ECC_POINT          *Qout,      // OUT: the public point
    TPM2B_ECC_PARAMETER     *dOut,      // OUT: the private scalar
    TPM_ECC_CURVE            curveId    // IN: the curve for the key
);

//*** CryptEccPointMultiply()
// This function computes 'R' := ['dIn']'G' + ['uIn']'QIn'. Where 'dIn' and
// 'uIn' are scalars, 'G' and 'QIn' are points on the specified curve and 'G' is the
// default generator of the curve.
//
// The 'xOut' and 'yOut' parameters are optional and may be set to NULL if not
// used.
//
// It is not necessary to provide 'uIn' if 'QIn' is specified but one of 'uIn' and
// 'dIn' must be provided. If 'dIn' and 'QIn' are specified but 'uIn' is not
// provided, then 'R' = ['dIn']'QIn'.
//
// If the multiply produces the point at infinity, the TPM_RC_NO_RESULT is returned.
//
// The sizes of 'xOut' and yOut' will be set to be the size of the degree of
// the curve
//
// It is a fatal error if 'dIn' and 'uIn' are both unspecified (NULL) or if 'Qin'
// or 'Rout' is unspecified.
//
//  Return Type: TPM_RC
//      TPM_RC_ECC_POINT         the point 'Pin' or 'Qin' is not on the curve
//      TPM_RC_NO_RESULT         the product point is at infinity
//      TPM_RC_CURVE             bad curve
//      TPM_RC_VALUE             'dIn' or 'uIn' out of range
//
LIB_EXPORT TPM_RC
CryptEccPointMultiply(
    TPMS_ECC_POINT      *Rout,              // OUT: the product point R
    TPM_ECC_CURVE        curveId,           // IN: the curve to use
    TPMS_ECC_POINT      *Pin,               // IN: first point (can be null)
    TPM2B_ECC_PARAMETER *dIn,               // IN: scalar value for [dIn]Qin
                                            //     the Pin
    TPMS_ECC_POINT      *Qin,               // IN: point Q
    TPM2B_ECC_PARAMETER *uIn                // IN: scalar value for the multiplier
                                            //     of Q
);

//*** CryptEccIsPointOnCurve()
// This function is used to test if a point is on a defined curve. It does this
// by checking that 'y'^2 mod 'p' = 'x'^3 + 'a'*'x' + 'b' mod 'p'.
//
// It is a fatal error if 'Q' is not specified (is NULL).
//  Return Type: BOOL
//      TRUE(1)         point is on curve
//      FALSE(0)        point is not on curve or curve is not supported
LIB_EXPORT BOOL
CryptEccIsPointOnCurve(
    TPM_ECC_CURVE            curveId,       // IN: the curve selector
    TPMS_ECC_POINT          *Qin            // IN: the point.
);

//*** CryptEccGenerateKey()
// This function generates an ECC key pair based on the input parameters.
// This routine uses KDFa to produce candidate numbers. The method is according
// to FIPS 186-3, section B.1.2 "Key Pair Generation by Testing Candidates."
// According to the method in FIPS 186-3, the resulting private value 'd' should be
// 1 <= 'd' < 'n' where 'n' is the order of the base point.
//
// It is a fatal error if 'Qout', 'dOut', is not provided (is NULL).
//
// If the curve is not supported
// If 'seed' is not provided, then a random number will be used for the key
//  Return Type: TPM_RC
//      TPM_RC_CURVE            curve is not supported
//      TPM_RC_NO_RESULT        could not verify key with signature (FIPS only)
LIB_EXPORT TPM_RC
CryptEccGenerateKey(
    TPMT_PUBLIC         *publicArea,        // IN/OUT: The public area template for
                                            //      the new key. The public key
                                            //      area will be replaced computed
                                            //      ECC public key
    TPMT_SENSITIVE      *sensitive,         // OUT: the sensitive area will be
                                            //      updated to contain the private
                                            //      ECC key and the symmetric
                                            //      encryption key
    RAND_STATE          *rand               // IN: if not NULL, the deterministic
                                            //     RNG state
);
#endif  // ALG_ECC

#endif  // _CRYPT_ECC_MAIN_FP_H_
