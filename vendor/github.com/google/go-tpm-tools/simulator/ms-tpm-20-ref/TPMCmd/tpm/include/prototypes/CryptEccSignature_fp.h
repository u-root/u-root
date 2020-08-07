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

#ifndef    _CRYPT_ECC_SIGNATURE_FP_H_
#define    _CRYPT_ECC_SIGNATURE_FP_H_

#if ALG_ECC

//*** BnSignEcdsa()
// This function implements the ECDSA signing algorithm. The method is described
// in the comments below.
TPM_RC
BnSignEcdsa(
    bigNum                   bnR,           // OUT: 'r' component of the signature
    bigNum                   bnS,           // OUT: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in the signature
                                            //     process
    bigNum                   bnD,           // IN: private signing key
    const TPM2B_DIGEST      *digest,        // IN: the digest to sign
    RAND_STATE              *rand           // IN: used in debug of signing
);

//*** CryptEccSign()
// This function is the dispatch function for the various ECC-based
// signing schemes.
// There is a bit of ugliness to the parameter passing. In order to test this,
// we sometime would like to use a deterministic RNG so that we can get the same
// signatures during testing. The easiest way to do this for most schemes is to
// pass in a deterministic RNG and let it return canned values during testing.
// There is a competing need for a canned parameter to use in ECDAA. To accommodate
// both needs with minimal fuss, a special type of RAND_STATE is defined to carry
// the address of the commit value. The setup and handling of this is not very
// different for the caller than what was in previous versions of the code.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME            'scheme' is not supported
LIB_EXPORT TPM_RC
CryptEccSign(
    TPMT_SIGNATURE          *signature,     // OUT: signature
    OBJECT                  *signKey,       // IN: ECC key to sign the hash
    const TPM2B_DIGEST      *digest,        // IN: digest to sign
    TPMT_ECC_SCHEME         *scheme,        // IN: signing scheme
    RAND_STATE              *rand
);
#if ALG_ECDSA

//*** BnValidateSignatureEcdsa()
// This function validates an ECDSA signature. rIn and sIn should have been checked
// to make sure that they are in the range 0 < 'v' < 'n'
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE           signature not valid
TPM_RC
BnValidateSignatureEcdsa(
    bigNum                   bnR,           // IN: 'r' component of the signature
    bigNum                   bnS,           // IN: 's' component of the signature
    bigCurve                 E,             // IN: the curve used in the signature
                                            //     process
    bn_point_t              *ecQ,           // IN: the public point of the key
    const TPM2B_DIGEST      *digest         // IN: the digest that was signed
);
#endif      // ALG_ECDSA

//*** CryptEccValidateSignature()
// This function validates an EcDsa or EcSchnorr signature.
// The point 'Qin' needs to have been validated to be on the curve of 'curveId'.
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE            not a valid signature
LIB_EXPORT TPM_RC
CryptEccValidateSignature(
    TPMT_SIGNATURE          *signature,     // IN: signature to be verified
    OBJECT                  *signKey,       // IN: ECC key signed the hash
    const TPM2B_DIGEST      *digest         // IN: digest that was signed
);

//***CryptEccCommitCompute()
// This function performs the point multiply operations required by TPM2_Commit.
//
// If 'B' or 'M' is provided, they must be on the curve defined by 'curveId'. This
// routine does not check that they are on the curve and results are unpredictable
// if they are not.
//
// It is a fatal error if 'r' is NULL. If 'B' is not NULL, then it is a
// fatal error if 'd' is NULL or if 'K' and 'L' are both NULL.
// If 'M' is not NULL, then it is a fatal error if 'E' is NULL.
//
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        if 'K', 'L' or 'E' was computed to be the point
//                              at infinity
//      TPM_RC_CANCELED         a cancel indication was asserted during this
//                              function
LIB_EXPORT TPM_RC
CryptEccCommitCompute(
    TPMS_ECC_POINT          *K,             // OUT: [d]B or [r]Q
    TPMS_ECC_POINT          *L,             // OUT: [r]B
    TPMS_ECC_POINT          *E,             // OUT: [r]M
    TPM_ECC_CURVE            curveId,       // IN: the curve for the computations
    TPMS_ECC_POINT          *M,             // IN: M (optional)
    TPMS_ECC_POINT          *B,             // IN: B (optional)
    TPM2B_ECC_PARAMETER     *d,             // IN: d (optional)
    TPM2B_ECC_PARAMETER     *r              // IN: the computed r value (required)
);
#endif  // ALG_ECC

#endif  // _CRYPT_ECC_SIGNATURE_FP_H_
