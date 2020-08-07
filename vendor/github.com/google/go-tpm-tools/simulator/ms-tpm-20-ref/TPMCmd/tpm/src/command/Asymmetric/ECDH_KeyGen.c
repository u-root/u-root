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
#include "ECDH_KeyGen_fp.h"

#if CC_ECDH_KeyGen  // Conditional expansion of this file

/*(See part 3 specification)
// This command uses the TPM to generate an ephemeral public key and the product
// of the ephemeral private key and the public portion of an ECC key.
*/
//  Return Type: TPM_RC
//      TPM_RC_KEY              'keyHandle' does not reference an ECC key
TPM_RC
TPM2_ECDH_KeyGen(
    ECDH_KeyGen_In      *in,            // IN: input parameter list
    ECDH_KeyGen_Out     *out            // OUT: output parameter list
    )
{
    OBJECT                  *eccKey;
    TPM2B_ECC_PARAMETER      sensitive;
    TPM_RC                   result;

// Input Validation

    eccKey = HandleToObject(in->keyHandle);

    // Referenced key must be an ECC key
    if(eccKey->publicArea.type != TPM_ALG_ECC)
        return TPM_RCS_KEY + RC_ECDH_KeyGen_keyHandle;

// Command Output
    do
    {
        TPMT_PUBLIC         *keyPublic = &eccKey->publicArea;
        // Create ephemeral ECC key
        result = CryptEccNewKeyPair(&out->pubPoint.point, &sensitive,
                                    keyPublic->parameters.eccDetail.curveID);
        if(result == TPM_RC_SUCCESS)
        {
            // Compute Z
            result = CryptEccPointMultiply(&out->zPoint.point,
                                           keyPublic->parameters.eccDetail.curveID,
                                           &keyPublic->unique.ecc, 
                                           &sensitive,
                                           NULL, NULL);
                    // The point in the key is not on the curve. Indicate
                    // that the key is bad.
            if(result == TPM_RC_ECC_POINT)
                return TPM_RCS_KEY + RC_ECDH_KeyGen_keyHandle;
             // The other possible error from CryptEccPointMultiply is
             // TPM_RC_NO_RESULT indicating that the multiplication resulted in
             // the point at infinity, so get a new random key and start over
             // BTW, this never happens.
        }
    } while(result == TPM_RC_NO_RESULT);
    return result;
}

#endif // CC_ECDH_KeyGen