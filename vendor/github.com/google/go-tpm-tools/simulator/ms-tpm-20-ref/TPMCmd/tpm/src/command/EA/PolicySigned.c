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
#include "Policy_spt_fp.h"
#include "PolicySigned_fp.h"

#if CC_PolicySigned  // Conditional expansion of this file

/*(See part 3 specification)
// Include an asymmetrically signed authorization to the policy evaluation
*/
//  Return Type: TPM_RC
//      TPM_RC_CPHASH           cpHash was previously set to a different value
//      TPM_RC_EXPIRED          'expiration' indicates a time in the past or
//                              'expiration' is non-zero but no nonceTPM is present
//      TPM_RC_NONCE            'nonceTPM' is not the nonce associated with the
//                              'policySession'
//      TPM_RC_SCHEME           the signing scheme of 'auth' is not supported by the
//                              TPM
//      TPM_RC_SIGNATURE        the signature is not genuine
//      TPM_RC_SIZE             input cpHash has wrong size
TPM_RC
TPM2_PolicySigned(
    PolicySigned_In     *in,            // IN: input parameter list
    PolicySigned_Out    *out            // OUT: output parameter list
    )
{
    TPM_RC                   result = TPM_RC_SUCCESS;
    SESSION                 *session;
    TPM2B_NAME               entityName;
    TPM2B_DIGEST             authHash;
    HASH_STATE               hashState;
    UINT64                   authTimeout = 0;
// Input Validation
    // Set up local pointers
    session = SessionGet(in->policySession);    // the session structure

    // Only do input validation if this is not a trial policy session
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        authTimeout = ComputeAuthTimeout(session, in->expiration, &in->nonceTPM);

        result = PolicyParameterChecks(session, authTimeout,
                                       &in->cpHashA, &in->nonceTPM,
                                       RC_PolicySigned_nonceTPM,
                                       RC_PolicySigned_cpHashA,
                                       RC_PolicySigned_expiration);
        if(result != TPM_RC_SUCCESS)
            return result;
        // Re-compute the digest being signed
        /*(See part 3 specification)
        // The digest is computed as:
        //     aHash := hash ( nonceTPM | expiration | cpHashA | policyRef)
        //  where:
        //      hash()      the hash associated with the signed authorization
        //      nonceTPM    the nonceTPM value from the TPM2_StartAuthSession .
        //                  response If the authorization is not limited to this
        //                  session, the size of this value is zero.
        //      expiration  time limit on authorization set by authorizing object.
        //                  This 32-bit value is set to zero if the expiration
        //                  time is not being set.
        //      cpHashA     hash of the command parameters for the command being
        //                  approved using the hash algorithm of the PSAP session.
        //                  Set to NULLauth if the authorization is not limited
        //                  to a specific command.
        //      policyRef   hash of an opaque value determined by the authorizing
        //                  object.  Set to the NULLdigest if no hash is present.
        */
        // Start hash
        authHash.t.size = CryptHashStart(&hashState,
                                         CryptGetSignHashAlg(&in->auth));
        // If there is no digest size, then we don't have a verification function
        // for this algorithm (e.g. TPM_ALG_ECDAA) so indicate that it is a
        // bad scheme.
        if(authHash.t.size == 0)
            return TPM_RCS_SCHEME + RC_PolicySigned_auth;

        //  nonceTPM
        CryptDigestUpdate2B(&hashState, &in->nonceTPM.b);

        //  expiration
        CryptDigestUpdateInt(&hashState, sizeof(UINT32), in->expiration);

        //  cpHashA
        CryptDigestUpdate2B(&hashState, &in->cpHashA.b);

        //  policyRef
        CryptDigestUpdate2B(&hashState, &in->policyRef.b);

        //  Complete digest
        CryptHashEnd2B(&hashState, &authHash.b);

        // Validate Signature.  A TPM_RC_SCHEME, TPM_RC_HANDLE or TPM_RC_SIGNATURE
        // error may be returned at this point
        result = CryptValidateSignature(in->authObject, &authHash, &in->auth);
        if(result != TPM_RC_SUCCESS)
            return RcSafeAddToResult(result, RC_PolicySigned_auth);
    }
// Internal Data Update
    // Update policy with input policyRef and name of authorization key
    // These values are updated even if the session is a trial session
    PolicyContextUpdate(TPM_CC_PolicySigned,
                        EntityGetName(in->authObject, &entityName),
                        &in->policyRef,
                        &in->cpHashA, authTimeout, session);
// Command Output
    // Create ticket and timeout buffer if in->expiration < 0 and this is not
    // a trial session.
    // NOTE: PolicyParameterChecks() makes sure that nonceTPM is present
    // when expiration is non-zero.
    if(in->expiration < 0
       && session->attributes.isTrialPolicy == CLEAR)
    {
        BOOL        expiresOnReset = (in->nonceTPM.t.size == 0);
        // Compute policy ticket
        authTimeout &= ~EXPIRATION_BIT;

        TicketComputeAuth(TPM_ST_AUTH_SIGNED, EntityGetHierarchy(in->authObject),
                          authTimeout, expiresOnReset, &in->cpHashA, &in->policyRef,
                          &entityName, &out->policyTicket);
        // Generate timeout buffer.  The format of output timeout buffer is
        // TPM-specific.
        // Note: In this implementation, the timeout buffer value is computed after 
        // the ticket is produced so, when the ticket is checked, the expiration
        // flag needs to be extracted before the ticket is checked.
        // In the Windows compatible version, the least-significant bit of the
        // timeout value is used as a flag to indicate if the authorization expires
        // on reset. The flag is the MSb.
        out->timeout.t.size = sizeof(authTimeout);
        if(expiresOnReset)
            authTimeout |= EXPIRATION_BIT;
        UINT64_TO_BYTE_ARRAY(authTimeout, out->timeout.t.buffer);
    }
    else
    {
        // Generate a null ticket.
        // timeout buffer is null
        out->timeout.t.size = 0;

        // authorization ticket is null
        out->policyTicket.tag = TPM_ST_AUTH_SIGNED;
        out->policyTicket.hierarchy = TPM_RH_NULL;
        out->policyTicket.digest.t.size = 0;
    }
    return TPM_RC_SUCCESS;
}

#endif // CC_PolicySigned