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
#include "MakeCredential_fp.h"

#if CC_MakeCredential  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// Make Credential with an object
*/
//  Return Type: TPM_RC
//      TPM_RC_KEY              'handle' referenced an ECC key that has a unique
//                              field that is not a point on the curve of the key
//      TPM_RC_SIZE             'credential' is larger than the digest size of
//                              Name algorithm of 'handle'
//      TPM_RC_TYPE             'handle' does not reference an asymmetric
//                              decryption key
TPM_RC
TPM2_MakeCredential(
    MakeCredential_In   *in,            // IN: input parameter list
    MakeCredential_Out  *out            // OUT: output parameter list
    )
{
    TPM_RC               result = TPM_RC_SUCCESS;

    OBJECT              *object;
    TPM2B_DATA           data;

// Input Validation

    // Get object pointer
    object = HandleToObject(in->handle);

    // input key must be an asymmetric, restricted decryption key
    // NOTE: Needs to be restricted to have a symmetric value.
    if(!CryptIsAsymAlgorithm(object->publicArea.type)
       || !IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, decrypt)
       || !IS_ATTRIBUTE(object->publicArea.objectAttributes, 
                        TPMA_OBJECT, restricted))
        return TPM_RCS_TYPE + RC_MakeCredential_handle;

    // The credential information may not be larger than the digest size used for
    // the Name of the key associated with handle.
    if(in->credential.t.size > CryptHashGetDigestSize(object->publicArea.nameAlg))
        return TPM_RCS_SIZE + RC_MakeCredential_credential;

// Command Output

    // Make encrypt key and its associated secret structure.
    out->secret.t.size = sizeof(out->secret.t.secret);
    result = CryptSecretEncrypt(object, IDENTITY_STRING, &data, &out->secret);
    if(result != TPM_RC_SUCCESS)
        return result;

    // Prepare output credential data from secret
    SecretToCredential(&in->credential, &in->objectName.b, &data.b,
                       object, &out->credentialBlob);

    return TPM_RC_SUCCESS;
}

#endif // CC_MakeCredential