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
//** Includes
#include "Tpm.h"
#include "Policy_spt_fp.h"
#include "PolicySigned_fp.h"
#include "PolicySecret_fp.h"
#include "PolicyTicket_fp.h"

//** Functions
//*** PolicyParameterChecks()
// This function validates the common parameters of TPM2_PolicySiged()
// and TPM2_PolicySecret(). The common parameters are 'nonceTPM',
// 'expiration', and 'cpHashA'.
TPM_RC
PolicyParameterChecks(
    SESSION         *session,
    UINT64           authTimeout,
    TPM2B_DIGEST    *cpHashA,
    TPM2B_NONCE     *nonce,
    TPM_RC           blameNonce,
    TPM_RC           blameCpHash,
    TPM_RC           blameExpiration
    )
{
    // Validate that input nonceTPM is correct if present
    if(nonce != NULL && nonce->t.size != 0)
    {
        if(!MemoryEqual2B(&nonce->b, &session->nonceTPM.b))
            return TPM_RCS_NONCE + blameNonce;
    }
    // If authTimeout is set (expiration != 0...
    if(authTimeout != 0)
    {
        // Validate input expiration.
        // Cannot compare time if clock stop advancing.  A TPM_RC_NV_UNAVAILABLE
        // or TPM_RC_NV_RATE error may be returned here.
        RETURN_IF_NV_IS_NOT_AVAILABLE;

        // if the time has already passed or the time epoch has changed then the
        // time value is no longer good. 
        if((authTimeout < g_time) 
           || (session->epoch != g_timeEpoch)) 
            return TPM_RCS_EXPIRED + blameExpiration;
    }
    // If the cpHash is present, then check it
    if(cpHashA != NULL && cpHashA->t.size != 0)
    {
        // The cpHash input has to have the correct size
        if(cpHashA->t.size != session->u2.policyDigest.t.size)
            return TPM_RCS_SIZE + blameCpHash;

        // If the cpHash has already been set, then this input value
        // must match the current value.
        if(session->u1.cpHash.b.size != 0
           && !MemoryEqual2B(&cpHashA->b, &session->u1.cpHash.b))
            return TPM_RC_CPHASH;
    }
    return TPM_RC_SUCCESS;
}

//*** PolicyContextUpdate()
// Update policy hash
//      Update the policyDigest in policy session by extending policyRef and
//      objectName to it. This will also update the cpHash if it is present.
//  Return Type: void
void
PolicyContextUpdate(
    TPM_CC           commandCode,   // IN: command code
    TPM2B_NAME      *name,          // IN: name of entity
    TPM2B_NONCE     *ref,           // IN: the reference data
    TPM2B_DIGEST    *cpHash,        // IN: the cpHash (optional)
    UINT64           policyTimeout, // IN: the timeout value for the policy
    SESSION         *session        // IN/OUT: policy session to be updated
    )
{
    HASH_STATE           hashState;

    // Start hash
   CryptHashStart(&hashState, session->authHashAlg);
    

    // policyDigest size should always be the digest size of session hash algorithm.
    pAssert(session->u2.policyDigest.t.size 
            == CryptHashGetDigestSize(session->authHashAlg));

    // add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    // add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(commandCode), commandCode);

    // add name if applicable
    if(name != NULL)
        CryptDigestUpdate2B(&hashState, &name->b);

    // Complete the digest and get the results
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    // If the policy reference is not null, do a second update to the digest.
    if(ref != NULL)
    {

        // Start second hash computation
        CryptHashStart(&hashState, session->authHashAlg);

        // add policyDigest
        CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

        // add policyRef
        CryptDigestUpdate2B(&hashState, &ref->b);

        // Complete second digest
        CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);
    }
    // Deal with the cpHash. If the cpHash value is present
    // then it would have already been checked to make sure that
    // it is compatible with the current value so all we need
    // to do here is copy it and set the isCpHashDefined attribute
    if(cpHash != NULL && cpHash->t.size != 0)
    {
        session->u1.cpHash = *cpHash;
        session->attributes.isCpHashDefined = SET;
    }

    // update the timeout if it is specified
    if(policyTimeout != 0)
    {
        // If the timeout has not been set, then set it to the new value
        // than the current timeout then set it to the new value
        if(session->timeout == 0 || session->timeout > policyTimeout)
            session->timeout = policyTimeout;
    }
    return;
}
//*** ComputeAuthTimeout()
// This function is used to determine what the authorization timeout value for
// the session should be.
UINT64
ComputeAuthTimeout(
    SESSION         *session,               // IN: the session containing the time
                                            //     values
    INT32            expiration,            // IN: either the number of seconds from
                                            //     the start of the session or the
                                            //     time in g_timer;
    TPM2B_NONCE     *nonce                  // IN: indicator of the time base
    )
{
    UINT64           policyTime;
    // If no expiration, policy time is 0
    if(expiration == 0)
        policyTime = 0;
    else
    {
        if(expiration < 0)
            expiration = -expiration;
        if(nonce->t.size == 0)
            // The input time is absolute Time (not Clock), but it is expressed
            // in seconds. To make sure that we don't time out too early, take the
            // current value of milliseconds in g_time and add that to the input
            // seconds value.
            policyTime = (((UINT64)expiration) * 1000) + g_time % 1000;
        else
            // The policy timeout is the absolute value of the expiration in seconds
            // added to the start time of the policy.
            policyTime = session->startTime + (((UINT64)expiration) * 1000);

    }
    return policyTime;
}

//*** PolicyDigestClear()
// Function to reset the policyDigest of a session
void
PolicyDigestClear(
    SESSION         *session
    )
{
    session->u2.policyDigest.t.size = CryptHashGetDigestSize(session->authHashAlg);
    MemorySet(session->u2.policyDigest.t.buffer, 0, 
              session->u2.policyDigest.t.size);
}

BOOL
PolicySptCheckCondition(
    TPM_EO          operation,
    BYTE            *opA,
    BYTE            *opB,
    UINT16           size
    )
{
    // Arithmetic Comparison
    switch(operation)
    {
        case TPM_EO_EQ:
            // compare A = B
            return (UnsignedCompareB(size, opA, size, opB) == 0);
            break;
        case TPM_EO_NEQ:
            // compare A != B
            return (UnsignedCompareB(size, opA, size, opB) != 0);
            break;
        case TPM_EO_SIGNED_GT:
            // compare A > B signed
            return (SignedCompareB(size, opA, size, opB) > 0);
            break;
        case TPM_EO_UNSIGNED_GT:
            // compare A > B unsigned
            return (UnsignedCompareB(size, opA, size, opB) > 0);
            break;
        case TPM_EO_SIGNED_LT:
            // compare A < B signed
            return (SignedCompareB(size, opA, size, opB) < 0);
            break;
        case TPM_EO_UNSIGNED_LT:
            // compare A < B unsigned
            return (UnsignedCompareB(size, opA, size, opB) < 0);
            break;
        case TPM_EO_SIGNED_GE:
            // compare A >= B signed
            return (SignedCompareB(size, opA, size, opB) >= 0);
            break;
        case TPM_EO_UNSIGNED_GE:
            // compare A >= B unsigned
            return (UnsignedCompareB(size, opA, size, opB) >= 0);
            break;
        case TPM_EO_SIGNED_LE:
            // compare A <= B signed
            return (SignedCompareB(size, opA, size, opB) <= 0);
            break;
        case TPM_EO_UNSIGNED_LE:
            // compare A <= B unsigned
            return (UnsignedCompareB(size, opA, size, opB) <= 0);
            break;
        case TPM_EO_BITSET:
            // All bits SET in B are SET in A. ((A&B)=B)
        {
            UINT32 i;
            for(i = 0; i < size; i++)
                if((opA[i] & opB[i]) != opB[i])
                    return FALSE;
        }
        break;
        case TPM_EO_BITCLEAR:
            // All bits SET in B are CLEAR in A. ((A&B)=0)
        {
            UINT32 i;
            for(i = 0; i < size; i++)
                if((opA[i] & opB[i]) != 0)
                    return FALSE;
        }
        break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    return TRUE;
}
