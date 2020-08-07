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
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Mar 28, 2019  Time: 08:25:19PM
 */

#ifndef    _TICKET_FP_H_
#define    _TICKET_FP_H_

//*** TicketIsSafe()
// This function indicates if producing a ticket is safe.
// It checks if the leading bytes of an input buffer is TPM_GENERATED_VALUE
// or its substring of canonical form.  If so, it is not safe to produce ticket
// for an input buffer claiming to be TPM generated buffer
//  Return Type: BOOL
//      TRUE(1)         safe to produce ticket
//      FALSE(0)        not safe to produce ticket
BOOL
TicketIsSafe(
    TPM2B           *buffer
);

//*** TicketComputeVerified()
// This function creates a TPMT_TK_VERIFIED ticket.
void
TicketComputeVerified(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy constant for ticket
    TPM2B_DIGEST        *digest,        // IN: digest
    TPM2B_NAME          *keyName,       // IN: name of key that signed the values
    TPMT_TK_VERIFIED    *ticket         // OUT: verified ticket
);

//*** TicketComputeAuth()
// This function creates a TPMT_TK_AUTH ticket.
void
TicketComputeAuth(
    TPM_ST               type,          // IN: the type of ticket.
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy constant for ticket
    UINT64               timeout,       // IN: timeout
    BOOL                 expiresOnReset,// IN: flag to indicate if ticket expires on
                                        //      TPM Reset
    TPM2B_DIGEST        *cpHashA,       // IN: input cpHashA
    TPM2B_NONCE         *policyRef,     // IN: input policyRef
    TPM2B_NAME          *entityName,    // IN: name of entity
    TPMT_TK_AUTH        *ticket         // OUT: Created ticket
);

//*** TicketComputeHashCheck()
// This function creates a TPMT_TK_HASHCHECK ticket.
void
TicketComputeHashCheck(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy constant for ticket
    TPM_ALG_ID           hashAlg,       // IN: the hash algorithm for 'digest'
    TPM2B_DIGEST        *digest,        // IN: input digest
    TPMT_TK_HASHCHECK   *ticket         // OUT: Created ticket
);

//*** TicketComputeCreation()
// This function creates a TPMT_TK_CREATION ticket.
void
TicketComputeCreation(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy for ticket
    TPM2B_NAME          *name,          // IN: object name
    TPM2B_DIGEST        *creation,      // IN: creation hash
    TPMT_TK_CREATION    *ticket         // OUT: created ticket
);

#endif  // _TICKET_FP_H_
