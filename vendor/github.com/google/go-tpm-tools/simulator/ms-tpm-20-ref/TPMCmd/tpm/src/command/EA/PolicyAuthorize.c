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
#include "PolicyAuthorize_fp.h"

#if CC_PolicyAuthorize  // Conditional expansion of this file

#include "Policy_spt_fp.h"

/*(See part 3 specification)
// Change policy by a signature from authority
*/
//  Return Type: TPM_RC
//      TPM_RC_HASH         hash algorithm in 'keyName' is not supported
//      TPM_RC_SIZE         'keyName' is not the correct size for its hash algorithm
//      TPM_RC_VALUE        the current policyDigest of 'policySession' does not
//                          match 'approvedPolicy'; or 'checkTicket' doesn't match
//                          the provided values
TPM_RC
TPM2_PolicyAuthorize(
    PolicyAuthorize_In  *in             // IN: input parameter list
    )
{
    SESSION                 *session;
    TPM2B_DIGEST             authHash;
    HASH_STATE               hashState;
    TPMT_TK_VERIFIED         ticket;
    TPM_ALG_ID               hashAlg;
    UINT16                   digestSize;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // Extract from the Name of the key, the algorithm used to compute it's Name
    hashAlg = BYTE_ARRAY_TO_UINT16(in->keySign.t.name);

    // 'keySign' parameter needs to use a supported hash algorithm, otherwise
    // can't tell how large the digest should be
    if(!CryptHashIsValidAlg(hashAlg, FALSE))
        return TPM_RCS_HASH + RC_PolicyAuthorize_keySign;

    digestSize = CryptHashGetDigestSize(hashAlg);
    if(digestSize != (in->keySign.t.size - 2))
        return TPM_RCS_SIZE + RC_PolicyAuthorize_keySign;

    //If this is a trial policy, skip all validations
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        // Check that "approvedPolicy" matches the current value of the
        // policyDigest in policy session
        if(!MemoryEqual2B(&session->u2.policyDigest.b,
                          &in->approvedPolicy.b))
            return TPM_RCS_VALUE + RC_PolicyAuthorize_approvedPolicy;

        // Validate ticket TPMT_TK_VERIFIED
        // Compute aHash.  The authorizing object sign a digest
        //  aHash := hash(approvedPolicy || policyRef).
        // Start hash
        authHash.t.size = CryptHashStart(&hashState, hashAlg);

        // add approvedPolicy
        CryptDigestUpdate2B(&hashState, &in->approvedPolicy.b);

        // add policyRef
        CryptDigestUpdate2B(&hashState, &in->policyRef.b);

        // complete hash
        CryptHashEnd2B(&hashState, &authHash.b);

        // re-compute TPMT_TK_VERIFIED
        TicketComputeVerified(in->checkTicket.hierarchy, &authHash,
                              &in->keySign, &ticket);

        // Compare ticket digest.  If not match, return error
        if(!MemoryEqual2B(&in->checkTicket.digest.b, &ticket.digest.b))
            return TPM_RCS_VALUE + RC_PolicyAuthorize_checkTicket;
    }

// Internal Data Update

    // Set policyDigest to zero digest
    PolicyDigestClear(session);

    // Update policyDigest
    PolicyContextUpdate(TPM_CC_PolicyAuthorize, &in->keySign, &in->policyRef,
                        NULL, 0, session);

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyAuthorize