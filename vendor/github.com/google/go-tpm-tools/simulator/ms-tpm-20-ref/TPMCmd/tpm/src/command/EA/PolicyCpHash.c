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
#include "PolicyCpHash_fp.h"

#if CC_PolicyCpHash  // Conditional expansion of this file

/*(See part 3 specification)
// Add a cpHash restriction to the policyDigest
*/
//  Return Type: TPM_RC
//      TPM_RC_CPHASH           cpHash of 'policySession' has previously been set
//                              to a different value
//      TPM_RC_SIZE             'cpHashA' is not the size of a digest produced
//                              by the hash algorithm associated with
//                              'policySession'
TPM_RC
TPM2_PolicyCpHash(
    PolicyCpHash_In     *in             // IN: input parameter list
    )
{
    SESSION     *session;
    TPM_CC      commandCode = TPM_CC_PolicyCpHash;
    HASH_STATE  hashState;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // A valid cpHash must have the same size as session hash digest
    // NOTE: the size of the digest can't be zero because TPM_ALG_NULL
    // can't be used for the authHashAlg.
    if(in->cpHashA.t.size != CryptHashGetDigestSize(session->authHashAlg))
        return TPM_RCS_SIZE + RC_PolicyCpHash_cpHashA;

    // error if the cpHash in session context is not empty and is not the same 
    // as the input or is not a cpHash
    if((session->u1.cpHash.t.size != 0)
       && (!session->attributes.isCpHashDefined
           || !MemoryEqual2B(&in->cpHashA.b, &session->u1.cpHash.b)))
        return TPM_RC_CPHASH;
    

// Internal Data Update

    // Update policy hash
    // policyDigestnew = hash(policyDigestold || TPM_CC_PolicyCpHash || cpHashA)
    //  Start hash
    CryptHashStart(&hashState, session->authHashAlg);

    //  add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    //  add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    //  add cpHashA
    CryptDigestUpdate2B(&hashState, &in->cpHashA.b);

    //  complete the digest and get the results
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    // update cpHash in session context
    session->u1.cpHash = in->cpHashA;
    session->attributes.isCpHashDefined = SET;

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyCpHash