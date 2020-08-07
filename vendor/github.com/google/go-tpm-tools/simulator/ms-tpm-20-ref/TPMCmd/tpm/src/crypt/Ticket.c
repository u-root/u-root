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
//** Introduction
/*
    This clause contains the functions used for ticket computations.
*/

//** Includes
#include "Tpm.h"

//** Functions

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
    )
{
    TPM_GENERATED   valueToCompare = TPM_GENERATED_VALUE;
    BYTE            bufferToCompare[sizeof(valueToCompare)];
    BYTE            *marshalBuffer;
//
    // If the buffer size is less than the size of TPM_GENERATED_VALUE, assume
    // it is not safe to generate a ticket
    if(buffer->size < sizeof(valueToCompare))
        return FALSE;
    marshalBuffer = bufferToCompare;
    TPM_GENERATED_Marshal(&valueToCompare, &marshalBuffer, NULL);
    if(MemoryEqual(buffer->buffer, bufferToCompare, sizeof(valueToCompare)))
        return FALSE;
    else
        return TRUE;
}

//*** TicketComputeVerified()
// This function creates a TPMT_TK_VERIFIED ticket.
/*(See part 2 specification)
//  The ticket is computed as:
//      HMAC(proof, (TPM_ST_VERIFIED | digest | keyName))
//  Where:
//      HMAC()              an HMAC using the hash of proof
//      proof               a TPM secret value associated with the hierarchy
//                          associated with keyName
//      TPM_ST_VERIFIED     a value to differentiate the tickets
//      digest              the signed digest
//      keyName             the Name of the key that signed digest
*/
void
TicketComputeVerified(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy constant for ticket
    TPM2B_DIGEST        *digest,        // IN: digest
    TPM2B_NAME          *keyName,       // IN: name of key that signed the values
    TPMT_TK_VERIFIED    *ticket         // OUT: verified ticket
    )
{
    TPM2B_PROOF         *proof;
    HMAC_STATE           hmacState;
//
    // Fill in ticket fields
    ticket->tag = TPM_ST_VERIFIED;
    ticket->hierarchy = hierarchy;
    proof = HierarchyGetProof(hierarchy);

    // Start HMAC using the proof value of the hierarchy as the HMAC key
    ticket->digest.t.size = CryptHmacStart2B(&hmacState, CONTEXT_INTEGRITY_HASH_ALG,
                                             &proof->b);
        //  TPM_ST_VERIFIED
    CryptDigestUpdateInt(&hmacState, sizeof(TPM_ST), ticket->tag);
        //  digest
    CryptDigestUpdate2B(&hmacState.hashState, &digest->b);
        // key name
    CryptDigestUpdate2B(&hmacState.hashState, &keyName->b);
        // done
    CryptHmacEnd2B(&hmacState, &ticket->digest.b);

    return;
}

//*** TicketComputeAuth()
// This function creates a TPMT_TK_AUTH ticket.
/*(See part 2 specification)
//  The ticket is computed as:
//      HMAC(proof, (type || timeout || timeEpoch || cpHash 
//                        || policyRef || keyName))
//  where:
//      HMAC()      an HMAC using the hash of proof
//      proof       a TPM secret value associated with the hierarchy of the key
//                  associated with keyName.
//      type        a value to differentiate the tickets.  It could be either
//                  TPM_ST_AUTH_SECRET or TPM_ST_AUTH_SIGNED
//      timeout     TPM-specific value indicating when the authorization expires
//      timeEpoch   TPM-specific value indicating the epoch for the timeout
//      cpHash      optional hash (digest only) of the authorized command
//      policyRef   optional reference to a policy value
//      keyName name of the key that signed the authorization
*/
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
    )
{
    TPM2B_PROOF         *proof;
    HMAC_STATE           hmacState;
//
    // Get proper proof
    proof = HierarchyGetProof(hierarchy);

    // Fill in ticket fields
    ticket->tag = type;
    ticket->hierarchy = hierarchy;

    // Start HMAC with hierarchy proof as the HMAC key
    ticket->digest.t.size = CryptHmacStart2B(&hmacState, CONTEXT_INTEGRITY_HASH_ALG,
                                             &proof->b);
        //  TPM_ST_AUTH_SECRET or TPM_ST_AUTH_SIGNED,
    CryptDigestUpdateInt(&hmacState, sizeof(UINT16), ticket->tag);
    // cpHash
    CryptDigestUpdate2B(&hmacState.hashState, &cpHashA->b);
        //  policyRef
    CryptDigestUpdate2B(&hmacState.hashState, &policyRef->b);
        //  keyName
    CryptDigestUpdate2B(&hmacState.hashState, &entityName->b);
        //  timeout
    CryptDigestUpdateInt(&hmacState, sizeof(timeout), timeout);
    if(timeout != 0)
    {
            //  epoch
        CryptDigestUpdateInt(&hmacState.hashState, sizeof(CLOCK_NONCE), 
                             g_timeEpoch);
            // reset count
        if(expiresOnReset)
            CryptDigestUpdateInt(&hmacState.hashState, sizeof(gp.totalResetCount),
                                 gp.totalResetCount);
    }
        // done
    CryptHmacEnd2B(&hmacState, &ticket->digest.b);

    return;
}

//*** TicketComputeHashCheck()
// This function creates a TPMT_TK_HASHCHECK ticket.
/*(See part 2 specification)
//  The ticket is computed as:
//      HMAC(proof, (TPM_ST_HASHCHECK || digest ))
//  where:
//      HMAC()  an HMAC using the hash of proof
//      proof   a TPM secret value associated with the hierarchy
//      TPM_ST_HASHCHECK
//              a value to differentiate the tickets
//      digest  the digest of the data
*/
void
TicketComputeHashCheck(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy constant for ticket
    TPM_ALG_ID           hashAlg,       // IN: the hash algorithm for 'digest'
    TPM2B_DIGEST        *digest,        // IN: input digest
    TPMT_TK_HASHCHECK   *ticket         // OUT: Created ticket
    )
{
    TPM2B_PROOF         *proof;
    HMAC_STATE           hmacState;
//
    // Get proper proof
    proof = HierarchyGetProof(hierarchy);

    // Fill in ticket fields
    ticket->tag = TPM_ST_HASHCHECK;
    ticket->hierarchy = hierarchy;

    // Start HMAC using hierarchy proof as HMAC key
    ticket->digest.t.size = CryptHmacStart2B(&hmacState, CONTEXT_INTEGRITY_HASH_ALG,
                                             &proof->b);
        //  TPM_ST_HASHCHECK
    CryptDigestUpdateInt(&hmacState, sizeof(TPM_ST), ticket->tag);
        //  hash algorithm
    CryptDigestUpdateInt(&hmacState, sizeof(hashAlg), hashAlg);
        //  digest
    CryptDigestUpdate2B(&hmacState.hashState, &digest->b);
        // done
    CryptHmacEnd2B(&hmacState, &ticket->digest.b);

    return;
}

//*** TicketComputeCreation()
// This function creates a TPMT_TK_CREATION ticket.
/*(See part 2 specification)
// The ticket is computed as:
//      HMAC(proof, (TPM_ST_CREATION || Name || hash(TPMS_CREATION_DATA)))
//  Where:
//  HMAC()  an HMAC using the hash of proof
//  proof   a TPM secret value associated with the hierarchy associated with Name
//  TPM_ST_VERIFIED     a value to differentiate the tickets
//  Name    the Name of the object to which the creation data is to be associated
//  TPMS_CREATION_DATA  the creation data structure associated with Name
*/
void
TicketComputeCreation(
    TPMI_RH_HIERARCHY    hierarchy,     // IN: hierarchy for ticket
    TPM2B_NAME          *name,          // IN: object name
    TPM2B_DIGEST        *creation,      // IN: creation hash
    TPMT_TK_CREATION    *ticket         // OUT: created ticket
    )
{
    TPM2B_PROOF         *proof;
    HMAC_STATE           hmacState;

    // Get proper proof
    proof = HierarchyGetProof(hierarchy);

    // Fill in ticket fields
    ticket->tag = TPM_ST_CREATION;
    ticket->hierarchy = hierarchy;

    // Start HMAC using hierarchy proof as HMAC key
    ticket->digest.t.size = CryptHmacStart2B(&hmacState, CONTEXT_INTEGRITY_HASH_ALG,
                                             &proof->b);
        //  TPM_ST_CREATION
    CryptDigestUpdateInt(&hmacState, sizeof(TPM_ST), ticket->tag);
        //  name if provided
    if(name != NULL)
        CryptDigestUpdate2B(&hmacState.hashState, &name->b);
        //  creation hash
    CryptDigestUpdate2B(&hmacState.hashState, &creation->b);
        // Done
    CryptHmacEnd2B(&hmacState, &ticket->digest.b);

    return;
}