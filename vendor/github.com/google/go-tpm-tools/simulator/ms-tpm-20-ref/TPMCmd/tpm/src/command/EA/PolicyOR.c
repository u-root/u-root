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
#include "PolicyOR_fp.h"

#if CC_PolicyOR  // Conditional expansion of this file

#include "Policy_spt_fp.h"

/*(See part 3 specification)
// PolicyOR command
*/
//  Return Type: TPM_RC
//      TPM_RC_VALUE            no digest in 'pHashList' matched the current
//                              value of policyDigest for 'policySession'
TPM_RC
TPM2_PolicyOR(
    PolicyOR_In     *in             // IN: input parameter list
    )
{
    SESSION     *session;
    UINT32       i;

// Input Validation and Update

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // Compare and Update Internal Session policy if match
    for(i = 0; i < in->pHashList.count; i++)
    {
        if(session->attributes.isTrialPolicy == SET
           || (MemoryEqual2B(&session->u2.policyDigest.b,
                             &in->pHashList.digests[i].b)))
        {
            // Found a match
            HASH_STATE      hashState;
            TPM_CC          commandCode = TPM_CC_PolicyOR;

            // Start hash
            session->u2.policyDigest.t.size
                = CryptHashStart(&hashState, session->authHashAlg);
            // Set policyDigest to 0 string and add it to hash
            MemorySet(session->u2.policyDigest.t.buffer, 0,
                      session->u2.policyDigest.t.size);
            CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

            // add command code
            CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

            // Add each of the hashes in the list
            for(i = 0; i < in->pHashList.count; i++)
            {
                // Extend policyDigest
                CryptDigestUpdate2B(&hashState, &in->pHashList.digests[i].b);
            }
            // Complete digest
            CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

            return TPM_RC_SUCCESS;
        }
    }
    // None of the values in the list matched the current policyDigest
    return TPM_RCS_VALUE + RC_PolicyOR_pHashList;
}

#endif // CC_PolicyOR