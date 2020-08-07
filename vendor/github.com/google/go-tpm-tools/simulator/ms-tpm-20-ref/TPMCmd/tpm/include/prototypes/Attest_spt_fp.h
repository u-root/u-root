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

#ifndef    _ATTEST_SPT_FP_H_
#define    _ATTEST_SPT_FP_H_

//***FillInAttestInfo()
// Fill in common fields of TPMS_ATTEST structure.
void
FillInAttestInfo(
    TPMI_DH_OBJECT       signHandle,    // IN: handle of signing object
    TPMT_SIG_SCHEME     *scheme,        // IN/OUT: scheme to be used for signing
    TPM2B_DATA          *data,          // IN: qualifying data
    TPMS_ATTEST         *attest         // OUT: attest structure
);

//***SignAttestInfo()
// Sign a TPMS_ATTEST structure. If signHandle is TPM_RH_NULL, a null signature
// is returned.
//
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   'signHandle' references not a signing key
//      TPM_RC_SCHEME       'scheme' is not compatible with 'signHandle' type
//      TPM_RC_VALUE        digest generated for the given 'scheme' is greater than
//                          the modulus of 'signHandle' (for an RSA key);
//                          invalid commit status or failed to generate "r" value
//                          (for an ECC key)
TPM_RC
SignAttestInfo(
    OBJECT              *signKey,           // IN: sign object
    TPMT_SIG_SCHEME     *scheme,            // IN: sign scheme
    TPMS_ATTEST         *certifyInfo,       // IN: the data to be signed
    TPM2B_DATA          *qualifyingData,    // IN: extra data for the signing
                                            //     process
    TPM2B_ATTEST        *attest,            // OUT: marshaled attest blob to be
                                            //     signed
    TPMT_SIGNATURE      *signature          // OUT: signature
);

//*** IsSigningObject()
// Checks to see if the object is OK for signing. This is here rather than in
// Object_spt.c because all the attestation commands use this file but not
// Object_spt.c.
//  Return Type: BOOL
//      TRUE(1)         object may sign
//      FALSE(0)        object may not sign
BOOL
IsSigningObject(
    OBJECT          *object         // IN:
);

#endif  // _ATTEST_SPT_FP_H_
