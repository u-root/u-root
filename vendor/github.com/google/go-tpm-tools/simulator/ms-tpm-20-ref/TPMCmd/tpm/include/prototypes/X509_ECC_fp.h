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
 *  Date: Apr  2, 2019  Time: 11:00:49AM
 */

#ifndef    _X509_ECC_FP_H_
#define    _X509_ECC_FP_H_

//*** X509PushPoint()
// This seems like it might be used more than once so...
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509PushPoint(
    ASN1MarshalContext      *ctx,
    TPMS_ECC_POINT          *p
);

//*** X509AddSigningAlgorithmECC()
// This creates the singing algorithm data.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddSigningAlgorithmECC(
    OBJECT              *signKey,
    TPMT_SIG_SCHEME     *scheme,
    ASN1MarshalContext  *ctx
);

//*** X509AddPublicECC()
// This function will add the publicKey description to the DER data. If ctx is
// NULL, then no data is transferred and this function will indicate if the TPM
// has the values for DER-encoding of the public key.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddPublicECC(
    OBJECT                *object,
    ASN1MarshalContext    *ctx
);

#endif  // _X509_ECC_FP_H_
