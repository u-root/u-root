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
#include "ActivateCredential_fp.h"

#if CC_ActivateCredential  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// Activate Credential with an object
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       'keyHandle' does not reference a decryption key
//      TPM_RC_ECC_POINT        'secret' is invalid (when 'keyHandle' is an ECC key)
//      TPM_RC_INSUFFICIENT     'secret' is invalid (when 'keyHandle' is an ECC key)
//      TPM_RC_INTEGRITY        'credentialBlob' fails integrity test
//      TPM_RC_NO_RESULT        'secret' is invalid (when 'keyHandle' is an ECC key)
//      TPM_RC_SIZE             'secret' size is invalid or the 'credentialBlob'
//                              does not unmarshal correctly
//      TPM_RC_TYPE             'keyHandle' does not reference an asymmetric key.
//      TPM_RC_VALUE            'secret' is invalid (when 'keyHandle' is an RSA key)
TPM_RC
TPM2_ActivateCredential(
    ActivateCredential_In   *in,            // IN: input parameter list
    ActivateCredential_Out  *out            // OUT: output parameter list
    )
{
    TPM_RC                   result = TPM_RC_SUCCESS;
    OBJECT                  *object;            // decrypt key
    OBJECT                  *activateObject;    // key associated with credential
    TPM2B_DATA               data;          // credential data

// Input Validation

    // Get decrypt key pointer
    object = HandleToObject(in->keyHandle);

    // Get certificated object pointer
    activateObject = HandleToObject(in->activateHandle);

    // input decrypt key must be an asymmetric, restricted decryption key
    if(!CryptIsAsymAlgorithm(object->publicArea.type)
       || !IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, decrypt)
       || !IS_ATTRIBUTE(object->publicArea.objectAttributes, 
                        TPMA_OBJECT, restricted))
        return TPM_RCS_TYPE + RC_ActivateCredential_keyHandle;

// Command output

    // Decrypt input credential data via asymmetric decryption.  A
    // TPM_RC_VALUE, TPM_RC_KEY or unmarshal errors may be returned at this
    // point
    result = CryptSecretDecrypt(object, NULL, IDENTITY_STRING, &in->secret, &data);
    if(result != TPM_RC_SUCCESS)
    {
        if(result == TPM_RC_KEY)
            return TPM_RC_FAILURE;
        return RcSafeAddToResult(result, RC_ActivateCredential_secret);
    }

    // Retrieve secret data.  A TPM_RC_INTEGRITY error or unmarshal
    // errors may be returned at this point
    result = CredentialToSecret(&in->credentialBlob.b,
                                &activateObject->name.b,
                                &data.b,
                                object,
                                &out->certInfo);
    if(result != TPM_RC_SUCCESS)
        return RcSafeAddToResult(result, RC_ActivateCredential_credentialBlob);

    return TPM_RC_SUCCESS;
}

#endif // CC_ActivateCredential