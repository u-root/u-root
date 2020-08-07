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
#include "HMAC_fp.h"

#if CC_HMAC  // Conditional expansion of this file

/*(See part 3 specification)
// Compute HMAC on a data buffer
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       key referenced by 'handle' is a restricted key
//      TPM_RC_KEY              'handle' does not reference a signing key
//      TPM_RC_TYPE             key referenced by 'handle' is not an HMAC key
//      TPM_RC_VALUE           'hashAlg' is not compatible with the hash algorithm
//                              of the scheme of the object referenced by 'handle'
TPM_RC
TPM2_HMAC(
    HMAC_In         *in,            // IN: input parameter list
    HMAC_Out        *out            // OUT: output parameter list
    )
{
    HMAC_STATE               hmacState;
    OBJECT                  *hmacObject;
    TPMI_ALG_HASH            hashAlg;
    TPMT_PUBLIC             *publicArea;

// Input Validation

    // Get HMAC key object and public area pointers
    hmacObject = HandleToObject(in->handle);
    publicArea = &hmacObject->publicArea;
    // Make sure that the key is an HMAC key
    if(publicArea->type != TPM_ALG_KEYEDHASH)
        return TPM_RCS_TYPE + RC_HMAC_handle;

    // and that it is unrestricted
    if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, restricted))
        return TPM_RCS_ATTRIBUTES + RC_HMAC_handle;

    // and that it is a signing key
    if(!IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, sign))
        return TPM_RCS_KEY + RC_HMAC_handle;

    // See if the key has a default
    if(publicArea->parameters.keyedHashDetail.scheme.scheme == TPM_ALG_NULL)
        // it doesn't so use the input value
        hashAlg = in->hashAlg;
    else
    {
        // key has a default so use it
        hashAlg
            = publicArea->parameters.keyedHashDetail.scheme.details.hmac.hashAlg;
        // and verify that the input was either the  TPM_ALG_NULL or the default
        if(in->hashAlg != TPM_ALG_NULL && in->hashAlg != hashAlg)
            hashAlg = TPM_ALG_NULL;
    }
    // if we ended up without a hash algorithm then return an error
    if(hashAlg == TPM_ALG_NULL)
        return TPM_RCS_VALUE + RC_HMAC_hashAlg;

// Command Output

    // Start HMAC stack
    out->outHMAC.t.size = CryptHmacStart2B(&hmacState, hashAlg,
                                           &hmacObject->sensitive.sensitive.bits.b);
    // Adding HMAC data
    CryptDigestUpdate2B(&hmacState.hashState, &in->buffer.b);

    // Complete HMAC
    CryptHmacEnd2B(&hmacState, &out->outHMAC.b);

    return TPM_RC_SUCCESS;
}

#endif // CC_HMAC