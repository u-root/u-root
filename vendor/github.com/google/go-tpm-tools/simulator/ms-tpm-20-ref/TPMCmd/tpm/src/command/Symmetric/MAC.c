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
#include "MAC_fp.h"

#if CC_MAC  // Conditional expansion of this file

/*(See part 3 specification)
// Compute MAC on a data buffer
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       key referenced by 'handle' is a restricted key
//      TPM_RC_KEY              'handle' does not reference a signing key
//      TPM_RC_TYPE             key referenced by 'handle' is not an HMAC key
//      TPM_RC_VALUE           'hashAlg' is not compatible with the hash algorithm
//                              of the scheme of the object referenced by 'handle'
TPM_RC
TPM2_MAC(
    MAC_In         *in,            // IN: input parameter list
    MAC_Out        *out            // OUT: output parameter list
    )
{
    OBJECT                  *keyObject;
    HMAC_STATE               state;
    TPMT_PUBLIC             *publicArea;
    TPM_RC                   result;

// Input Validation
    // Get MAC key object and public area pointers
    keyObject = HandleToObject(in->handle);
    publicArea = &keyObject->publicArea;

    // If the key is not able to do a MAC, indicate that the handle selects an
    // object that can't do a MAC
    result = CryptSelectMac(publicArea, &in->inScheme);
    if(result == TPM_RCS_TYPE)
        return TPM_RCS_TYPE + RC_MAC_handle;
    // If there is another error type, indicate that the scheme and key are not
    // compatible
    if(result != TPM_RC_SUCCESS)
        return RcSafeAddToResult(result, RC_MAC_inScheme);
    // Make sure that the key is not restricted
    if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, restricted))
        return TPM_RCS_ATTRIBUTES + RC_MAC_handle;
    // and that it is a signing key
    if(!IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, sign))
        return TPM_RCS_KEY + RC_MAC_handle;
// Command Output
    out->outMAC.t.size = CryptMacStart(&state, &publicArea->parameters, 
                                       in->inScheme, 
                                       &keyObject->sensitive.sensitive.any.b);
    // If the mac can't start, treat it as a fatal error
    if(out->outMAC.t.size == 0)
        return TPM_RC_FAILURE;
    CryptDigestUpdate2B(&state.hashState, &in->buffer.b);
    // If the MAC result is not what was expected, it is a fatal error
    if(CryptHmacEnd2B(&state, &out->outMAC.b) != out->outMAC.t.size)
        return TPM_RC_FAILURE;
    return TPM_RC_SUCCESS;
}

#endif // CC_MAC