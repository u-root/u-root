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
#include "PolicyNameHash_fp.h"

#if CC_PolicyNameHash  // Conditional expansion of this file

/*(See part 3 specification)
// Add a nameHash restriction to the policyDigest
*/
//  Return Type: TPM_RC
//      TPM_RC_CPHASH     'nameHash' has been previously set to a different value
//      TPM_RC_SIZE       'nameHash' is not the size of the digest produced by the
//                        hash algorithm associated with 'policySession'
TPM_RC
TPM2_PolicyNameHash(
    PolicyNameHash_In   *in             // IN: input parameter list
    )
{
    SESSION             *session;
    TPM_CC               commandCode = TPM_CC_PolicyNameHash;
    HASH_STATE           hashState;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // A valid nameHash must have the same size as session hash digest
    // Since the authHashAlg for a session cannot be TPM_ALG_NULL, the digest size
    // is always non-zero.
    if(in->nameHash.t.size != CryptHashGetDigestSize(session->authHashAlg))
        return TPM_RCS_SIZE + RC_PolicyNameHash_nameHash;

    // u1 in the policy session context cannot otherwise be occupied
    if(session->u1.cpHash.b.size != 0
       || session->attributes.isBound 
       || session->attributes.isCpHashDefined 
       || session->attributes.isTemplateSet)
        return TPM_RC_CPHASH;

// Internal Data Update

    // Update policy hash
    // policyDigestnew = hash(policyDigestold || TPM_CC_PolicyNameHash || nameHash)
    //  Start hash
    CryptHashStart(&hashState, session->authHashAlg);

    //  add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    //  add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    //  add nameHash
    CryptDigestUpdate2B(&hashState, &in->nameHash.b);

    //  complete the digest
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    // update nameHash in session context
    session->u1.cpHash = in->nameHash;

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyNameHash