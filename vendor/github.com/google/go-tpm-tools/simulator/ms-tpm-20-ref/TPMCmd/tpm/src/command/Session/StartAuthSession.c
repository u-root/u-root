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
#include "StartAuthSession_fp.h"

#if CC_StartAuthSession  // Conditional expansion of this file

/*(See part 3 specification)
// Start an authorization session
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       'tpmKey' does not reference a decrypt key
//      TPM_RC_CONTEXT_GAP      the difference between the most recently created
//                              active context and the oldest active context is at
//                              the limits of the TPM
//      TPM_RC_HANDLE           input decrypt key handle only has public portion
//                              loaded
//      TPM_RC_MODE             'symmetric' specifies a block cipher but the mode
//                              is not TPM_ALG_CFB.
//      TPM_RC_SESSION_HANDLES  no session handle is available
//      TPM_RC_SESSION_MEMORY   no more slots for loading a session
//      TPM_RC_SIZE             nonce less than 16 octets or greater than the size
//                              of the digest produced by 'authHash'
//      TPM_RC_VALUE            secret size does not match decrypt key type; or the
//                              recovered secret is larger than the digest size of
//                              the nameAlg of 'tpmKey'; or, for an RSA decrypt key,
//                              if 'encryptedSecret' is greater than the
//                              public modulus of 'tpmKey'.
TPM_RC
TPM2_StartAuthSession(
    StartAuthSession_In     *in,            // IN: input parameter buffer
    StartAuthSession_Out    *out            // OUT: output parameter buffer
    )
{
    TPM_RC                   result = TPM_RC_SUCCESS;
    OBJECT                  *tpmKey;                // TPM key for decrypt salt
    TPM2B_DATA               salt;

// Input Validation

    // Check input nonce size.  IT should be at least 16 bytes but not larger
    // than the digest size of session hash.
    if(in->nonceCaller.t.size < 16
       || in->nonceCaller.t.size > CryptHashGetDigestSize(in->authHash))
        return TPM_RCS_SIZE + RC_StartAuthSession_nonceCaller;

    // If an decrypt key is passed in, check its validation
    if(in->tpmKey != TPM_RH_NULL)
    {
        // Get pointer to loaded decrypt key
        tpmKey = HandleToObject(in->tpmKey);

        // key must be asymmetric with its sensitive area loaded. Since this
        // command does not require authorization, the presence of the sensitive
        // area was not already checked as it is with most other commands that
        // use the sensitive are so check it here
        if(!CryptIsAsymAlgorithm(tpmKey->publicArea.type))
            return TPM_RCS_KEY + RC_StartAuthSession_tpmKey;
        // secret size cannot be 0
        if(in->encryptedSalt.t.size == 0)
            return TPM_RCS_VALUE + RC_StartAuthSession_encryptedSalt;
        // Decrypting salt requires accessing the private portion of a key.
        // Therefore, tmpKey can not be a key with only public portion loaded
        if(tpmKey->attributes.publicOnly)
            return TPM_RCS_HANDLE + RC_StartAuthSession_tpmKey;
        // HMAC session input handle check.
        // tpmKey should be a decryption key
        if(!IS_ATTRIBUTE(tpmKey->publicArea.objectAttributes, TPMA_OBJECT, decrypt))
            return TPM_RCS_ATTRIBUTES + RC_StartAuthSession_tpmKey;
        // Secret Decryption.  A TPM_RC_VALUE, TPM_RC_KEY or Unmarshal errors
        // may be returned at this point
        result = CryptSecretDecrypt(tpmKey, &in->nonceCaller, SECRET_KEY,
                                    &in->encryptedSalt, &salt);
        if(result != TPM_RC_SUCCESS)
            return TPM_RCS_VALUE + RC_StartAuthSession_encryptedSalt;
    }
    else
    {
        // secret size must be 0
        if(in->encryptedSalt.t.size != 0)
            return TPM_RCS_VALUE + RC_StartAuthSession_encryptedSalt;
        salt.t.size = 0;
    }
    switch(HandleGetType(in->bind))
    {
        case TPM_HT_TRANSIENT:
        {
            OBJECT      *object = HandleToObject(in->bind);
            // If the bind handle references a transient object, make sure that we
            // can get to the authorization value. Also, make sure that the object
            // has a proper Name (nameAlg != TPM_ALG_NULL). If it doesn't, then
            // it might be possible to bind to an object where the authValue is
            // known. This does not create a real issue in that, if you know the
            // authorization value, you can actually bind to the object. However,
            // there is a potential 
            if(object->attributes.publicOnly == SET) 
                return TPM_RCS_HANDLE + RC_StartAuthSession_bind;
            break;
        }
        case TPM_HT_NV_INDEX:
        // a PIN index can't be a bind object
        {
            NV_INDEX       *nvIndex = NvGetIndexInfo(in->bind, NULL);
            if(IsNvPinPassIndex(nvIndex->publicArea.attributes)
               || IsNvPinFailIndex(nvIndex->publicArea.attributes))
                return TPM_RCS_HANDLE + RC_StartAuthSession_bind;
            break;
        }
        default:
            break;
    }
    // If 'symmetric' is a symmetric block cipher (not TPM_ALG_NULL or TPM_ALG_XOR)
    // then the mode must be CFB.
    if(in->symmetric.algorithm != TPM_ALG_NULL
       && in->symmetric.algorithm != TPM_ALG_XOR
       && in->symmetric.mode.sym != TPM_ALG_CFB)
        return TPM_RCS_MODE + RC_StartAuthSession_symmetric;

// Internal Data Update and command output

    // Create internal session structure.  TPM_RC_CONTEXT_GAP, TPM_RC_NO_HANDLES
    // or TPM_RC_SESSION_MEMORY errors may be returned at this point.
    //
    // The detailed actions for creating the session context are not shown here
    // as the details are implementation dependent
    // SessionCreate sets the output handle and nonceTPM
    result = SessionCreate(in->sessionType, in->authHash, &in->nonceCaller,
                           &in->symmetric, in->bind, &salt, &out->sessionHandle,
                           &out->nonceTPM);
    return result;
}

#endif // CC_StartAuthSession