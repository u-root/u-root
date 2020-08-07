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
#include "Tpm.h"
#include "Attest_spt_fp.h"
#include "Certify_fp.h"

#if CC_Certify  // Conditional expansion of this file

/*(See part 3 specification)
// prove an object with a specific Name is loaded in the TPM
*/
//  Return Type: TPM_RC
//      TPM_RC_KEY          key referenced by 'signHandle' is not a signing key
//      TPM_RC_SCHEME       'inScheme' is not compatible with 'signHandle'
//      TPM_RC_VALUE        digest generated for 'inScheme' is greater or has larger
//                          size than the modulus of 'signHandle', or the buffer for
//                          the result in 'signature' is too small (for an RSA key);
//                          invalid commit status (for an ECC key with a split scheme)
TPM_RC
TPM2_Certify(
    Certify_In      *in,            // IN: input parameter list
    Certify_Out     *out            // OUT: output parameter list
    )
{
    TPMS_ATTEST             certifyInfo;
    OBJECT                  *signObject = HandleToObject(in->signHandle);
    OBJECT                  *certifiedObject = HandleToObject(in->objectHandle);
// Input validation
    if(!IsSigningObject(signObject))
        return TPM_RCS_KEY + RC_Certify_signHandle;
    if(!CryptSelectSignScheme(signObject, &in->inScheme))
        return TPM_RCS_SCHEME + RC_Certify_inScheme;

// Command Output
    // Filling in attest information
    // Common fields
    FillInAttestInfo(in->signHandle, &in->inScheme, &in->qualifyingData,
                     &certifyInfo);

    // Certify specific fields
    certifyInfo.type = TPM_ST_ATTEST_CERTIFY;
    // NOTE: the certified object is not allowed to be TPM_ALG_NULL so
    // 'certifiedObject' will never be NULL
    certifyInfo.attested.certify.name = certifiedObject->name;

    // When using an anonymous signing scheme, need to set the qualified Name to the
    // empty buffer to avoid correlation between keys
    if(CryptIsSchemeAnonymous(in->inScheme.scheme))
        certifyInfo.attested.certify.qualifiedName.t.size = 0;
    else
        certifyInfo.attested.certify.qualifiedName = certifiedObject->qualifiedName;


    // Sign attestation structure.  A NULL signature will be returned if
    // signHandle is TPM_RH_NULL.  A TPM_RC_NV_UNAVAILABLE, TPM_RC_NV_RATE,
    // TPM_RC_VALUE, TPM_RC_SCHEME or TPM_RC_ATTRIBUTES error may be returned
    // by SignAttestInfo()
    return SignAttestInfo(signObject, &in->inScheme, &certifyInfo,
                          &in->qualifyingData, &out->certifyInfo, &out->signature);
}

#endif // CC_Certify