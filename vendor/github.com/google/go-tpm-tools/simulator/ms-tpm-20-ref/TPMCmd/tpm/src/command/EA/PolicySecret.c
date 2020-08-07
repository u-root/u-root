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
#include "PolicySecret_fp.h"

#if CC_PolicySecret  // Conditional expansion of this file

#include "Policy_spt_fp.h"
#include "NV_spt_fp.h"

/*(See part 3 specification)
// Add a secret-based authorization to the policy evaluation
*/
//  Return Type: TPM_RC
//      TPM_RC_CPHASH           cpHash for policy was previously set to a
//                              value that is not the same as 'cpHashA'
//      TPM_RC_EXPIRED          'expiration' indicates a time in the past
//      TPM_RC_NONCE            'nonceTPM' does not match the nonce associated
//                              with 'policySession'
//      TPM_RC_SIZE             'cpHashA' is not the size of a digest for the
//                              hash associated with 'policySession'
TPM_RC
TPM2_PolicySecret(
    PolicySecret_In     *in,            // IN: input parameter list
    PolicySecret_Out    *out            // OUT: output parameter list
    )
{
    TPM_RC                   result;
    SESSION                 *session;
    TPM2B_NAME               entityName;
    UINT64                   authTimeout = 0;
// Input Validation
    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    //Only do input validation if this is not a trial policy session
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        authTimeout = ComputeAuthTimeout(session, in->expiration, &in->nonceTPM);

        result = PolicyParameterChecks(session, authTimeout,
                                       &in->cpHashA, &in->nonceTPM,
                                       RC_PolicySecret_nonceTPM,
                                       RC_PolicySecret_cpHashA,
                                       RC_PolicySecret_expiration);
        if(result != TPM_RC_SUCCESS)
            return result;
    }
// Internal Data Update
    // Update policy context with input policyRef and name of authorizing key
    // This value is computed even for trial sessions. Possibly update the cpHash
    PolicyContextUpdate(TPM_CC_PolicySecret,
                        EntityGetName(in->authHandle, &entityName), &in->policyRef,
                        &in->cpHashA, authTimeout, session);
// Command Output
    // Create ticket and timeout buffer if in->expiration < 0 and this is not
    // a trial session.
    // NOTE: PolicyParameterChecks() makes sure that nonceTPM is present
    // when expiration is non-zero.
    if(in->expiration < 0
       && session->attributes.isTrialPolicy == CLEAR
       && !NvIsPinPassIndex(in->authHandle))
    {
        BOOL        expiresOnReset = (in->nonceTPM.t.size == 0);
        // Compute policy ticket
        authTimeout &= ~EXPIRATION_BIT;
        TicketComputeAuth(TPM_ST_AUTH_SECRET, EntityGetHierarchy(in->authHandle),
                          authTimeout, expiresOnReset, &in->cpHashA, &in->policyRef,
                          &entityName, &out->policyTicket);
        // Generate timeout buffer.  The format of output timeout buffer is
        // TPM-specific.
        // Note: In this implementation, the timeout buffer value is computed after 
        // the ticket is produced so, when the ticket is checked, the expiration
        // flag needs to be extracted before the ticket is checked.
        out->timeout.t.size = sizeof(authTimeout);
        // In the Windows compatible version, the least-significant bit of the
        // timeout value is used as a flag to indicate if the authorization expires
        // on reset. The flag is the MSb.
        if(expiresOnReset)
            authTimeout |= EXPIRATION_BIT;
        UINT64_TO_BYTE_ARRAY(authTimeout, out->timeout.t.buffer);
    }
    else
    {
        // timeout buffer is null
        out->timeout.t.size = 0;

        // authorization ticket is null
        out->policyTicket.tag = TPM_ST_AUTH_SECRET;
        out->policyTicket.hierarchy = TPM_RH_NULL;
        out->policyTicket.digest.t.size = 0;
    }
    return TPM_RC_SUCCESS;
}

#endif // CC_PolicySecret