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
#include "RSA_Encrypt_fp.h"

#if CC_RSA_Encrypt  // Conditional expansion of this file

/*(See part 3 specification)
// This command performs the padding and encryption of a data block
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES           'decrypt' attribute is not SET in key referenced
//                                  by 'keyHandle'
//      TPM_RC_KEY                  'keyHandle' does not reference an RSA key
//      TPM_RC_SCHEME               incorrect input scheme, or the chosen
//                                  scheme is not a valid RSA decrypt scheme
//      TPM_RC_VALUE                the numeric value of 'message' is greater than
//                                  the public modulus of the key referenced by
//                                  'keyHandle', or 'label' is not a null-terminated
//                                  string
TPM_RC
TPM2_RSA_Encrypt(
    RSA_Encrypt_In      *in,            // IN: input parameter list
    RSA_Encrypt_Out     *out            // OUT: output parameter list
    )
{
    TPM_RC                  result;
    OBJECT                  *rsaKey;
    TPMT_RSA_DECRYPT        *scheme;
// Input Validation
    rsaKey = HandleToObject(in->keyHandle);

    // selected key must be an RSA key
    if(rsaKey->publicArea.type != TPM_ALG_RSA)
        return TPM_RCS_KEY + RC_RSA_Encrypt_keyHandle;
    // selected key must have the decryption attribute
    if(!IS_ATTRIBUTE(rsaKey->publicArea.objectAttributes, TPMA_OBJECT, decrypt))
        return TPM_RCS_ATTRIBUTES + RC_RSA_Encrypt_keyHandle;

    // Is there a label?
    if(!IsLabelProperlyFormatted(&in->label.b))
        return TPM_RCS_VALUE + RC_RSA_Encrypt_label;
// Command Output
    // Select a scheme for encryption
    scheme = CryptRsaSelectScheme(in->keyHandle, &in->inScheme);
    if(scheme == NULL)
        return TPM_RCS_SCHEME + RC_RSA_Encrypt_inScheme;

    // Encryption.  TPM_RC_VALUE, or TPM_RC_SCHEME errors my be returned buy
    // CryptEncyptRSA.
    out->outData.t.size = sizeof(out->outData.t.buffer);

    result = CryptRsaEncrypt(&out->outData, &in->message.b, rsaKey, scheme,
                             &in->label.b, NULL);
    return result;
}

#endif // CC_RSA_Encrypt