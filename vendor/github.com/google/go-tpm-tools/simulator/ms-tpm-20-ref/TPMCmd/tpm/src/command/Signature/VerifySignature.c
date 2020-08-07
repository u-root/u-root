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
#include "VerifySignature_fp.h"

#if CC_VerifySignature  // Conditional expansion of this file

/*(See part 3 specification)
// This command uses loaded key to validate an asymmetric signature on a message
// with the message digest passed to the TPM.
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES         'keyHandle' does not reference a signing key
//      TPM_RC_SIGNATURE          signature is not genuine
//      TPM_RC_SCHEME             CryptValidateSignature()
//      TPM_RC_HANDLE             the input handle is references an HMAC key but
//                                the private portion is not loaded
TPM_RC
TPM2_VerifySignature(
    VerifySignature_In      *in,            // IN: input parameter list
    VerifySignature_Out     *out            // OUT: output parameter list
    )
{
    TPM_RC                   result;
    OBJECT                  *signObject = HandleToObject(in->keyHandle);
    TPMI_RH_HIERARCHY        hierarchy;

// Input Validation
    // The object to validate the signature must be a signing key.
    if(!IS_ATTRIBUTE(signObject->publicArea.objectAttributes, TPMA_OBJECT, sign))
        return TPM_RCS_ATTRIBUTES + RC_VerifySignature_keyHandle;

    // Validate Signature.  TPM_RC_SCHEME, TPM_RC_HANDLE or TPM_RC_SIGNATURE
    // error may be returned by CryptCVerifySignatrue()
    result = CryptValidateSignature(in->keyHandle, &in->digest, &in->signature);
    if(result != TPM_RC_SUCCESS)
        return RcSafeAddToResult(result, RC_VerifySignature_signature);

// Command Output

    hierarchy = GetHeriarchy(in->keyHandle);
    if(hierarchy == TPM_RH_NULL
       || signObject->publicArea.nameAlg == TPM_ALG_NULL)
    {
        // produce empty ticket if hierarchy is TPM_RH_NULL or nameAlg is
        // ALG_NULL
        out->validation.tag = TPM_ST_VERIFIED;
        out->validation.hierarchy = TPM_RH_NULL;
        out->validation.digest.t.size = 0;
    }
    else
    {
        // Compute ticket
        TicketComputeVerified(hierarchy, &in->digest, &signObject->name,
                              &out->validation);
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_VerifySignature