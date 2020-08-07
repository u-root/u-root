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
#include "PolicyTicket_fp.h"

#if CC_PolicyTicket  // Conditional expansion of this file

#include "Policy_spt_fp.h"

/*(See part 3 specification)
// Include ticket to the policy evaluation
*/
//  Return Type: TPM_RC
//      TPM_RC_CPHASH           policy's cpHash was previously set to a different
//                              value
//      TPM_RC_EXPIRED          'timeout' value in the ticket is in the past and the
//                              ticket has expired
//      TPM_RC_SIZE             'timeout' or 'cpHash' has invalid size for the
//      TPM_RC_TICKET           'ticket' is not valid
TPM_RC
TPM2_PolicyTicket(
    PolicyTicket_In     *in             // IN: input parameter list
    )
{
    TPM_RC                   result;
    SESSION                 *session;
    UINT64                   authTimeout;
    TPMT_TK_AUTH             ticketToCompare;
    TPM_CC                   commandCode = TPM_CC_PolicySecret;
    BOOL                     expiresOnReset;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // NOTE: A trial policy session is not allowed to use this command.
    // A ticket is used in place of a previously given authorization. Since
    // a trial policy doesn't actually authenticate, the validated
    // ticket is not necessary and, in place of using a ticket, one
    // should use the intended authorization for which the ticket
    // would be a substitute.
    if(session->attributes.isTrialPolicy)
        return TPM_RCS_ATTRIBUTES + RC_PolicyTicket_policySession;
    // Restore timeout data.  The format of timeout buffer is TPM-specific.
    // In this implementation, the most significant bit of the timeout value is
    // used as the flag to indicate that the ticket expires on TPM Reset or 
    // TPM Restart. The flag has to be removed before the parameters and ticket
    // are checked.
    if(in->timeout.t.size != sizeof(UINT64))
        return TPM_RCS_SIZE + RC_PolicyTicket_timeout;
    authTimeout = BYTE_ARRAY_TO_UINT64(in->timeout.t.buffer);

    // extract the flag
    expiresOnReset = (authTimeout & EXPIRATION_BIT) != 0;
    authTimeout &= ~EXPIRATION_BIT;

    // Do the normal checks on the cpHashA and timeout values
    result = PolicyParameterChecks(session, authTimeout,
                                   &in->cpHashA, 
                                   NULL,                    // no nonce
                                   0,                       // no bad nonce return
                                   RC_PolicyTicket_cpHashA,
                                   RC_PolicyTicket_timeout);
    if(result != TPM_RC_SUCCESS)
        return result;
    // Validate Ticket
    // Re-generate policy ticket by input parameters
    TicketComputeAuth(in->ticket.tag, in->ticket.hierarchy, 
                      authTimeout, expiresOnReset, &in->cpHashA, &in->policyRef, 
                      &in->authName, &ticketToCompare);
    // Compare generated digest with input ticket digest
    if(!MemoryEqual2B(&in->ticket.digest.b, &ticketToCompare.digest.b))
        return TPM_RCS_TICKET + RC_PolicyTicket_ticket;

// Internal Data Update

    // Is this ticket to take the place of a TPM2_PolicySigned() or
    // a TPM2_PolicySecret()?
    if(in->ticket.tag == TPM_ST_AUTH_SIGNED)
        commandCode = TPM_CC_PolicySigned;
    else if(in->ticket.tag == TPM_ST_AUTH_SECRET)
        commandCode = TPM_CC_PolicySecret;
    else
        // There could only be two possible tag values.  Any other value should
        // be caught by the ticket validation process.
        FAIL(FATAL_ERROR_INTERNAL);

    // Update policy context
    PolicyContextUpdate(commandCode, &in->authName, &in->policyRef,
                        &in->cpHashA, authTimeout, session);

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyTicket