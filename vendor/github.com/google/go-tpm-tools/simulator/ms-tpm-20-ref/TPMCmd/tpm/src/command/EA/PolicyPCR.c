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
#include "PolicyPCR_fp.h"

#if CC_PolicyPCR  // Conditional expansion of this file

/*(See part 3 specification)
// Add a PCR gate for a policy session
*/
//  Return Type: TPM_RC
//      TPM_RC_VALUE          if provided, 'pcrDigest' does not match the
//                            current PCR settings
//      TPM_RC_PCR_CHANGED    a previous TPM2_PolicyPCR() set
//                            pcrCounter and it has changed
TPM_RC
TPM2_PolicyPCR(
    PolicyPCR_In    *in             // IN: input parameter list
    )
{
    SESSION         *session;
    TPM2B_DIGEST     pcrDigest;
    BYTE             pcrs[sizeof(TPML_PCR_SELECTION)];
    UINT32           pcrSize;
    BYTE            *buffer;
    TPM_CC           commandCode = TPM_CC_PolicyPCR;
    HASH_STATE       hashState;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // Compute current PCR digest
    PCRComputeCurrentDigest(session->authHashAlg, &in->pcrs, &pcrDigest);
    
    // Do validation for non trial session
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        // Make sure that this is not going to invalidate a previous PCR check
        if(session->pcrCounter != 0 && session->pcrCounter != gr.pcrCounter)
            return TPM_RC_PCR_CHANGED;

        // If the caller specified the PCR digest and it does not
        // match the current PCR settings, return an error..
        if(in->pcrDigest.t.size != 0)
        {
            if(!MemoryEqual2B(&in->pcrDigest.b, &pcrDigest.b))
                return TPM_RCS_VALUE + RC_PolicyPCR_pcrDigest;
        }
    }
    else
    {
        // For trial session, just use the input PCR digest if one provided
        // Note: It can't be too big because it is a TPM2B_DIGEST and the size 
        // would have been checked during unmarshaling
        if(in->pcrDigest.t.size != 0)
            pcrDigest = in->pcrDigest;
    }
// Internal Data Update
    // Update policy hash
    // policyDigestnew = hash(   policyDigestold || TPM_CC_PolicyPCR
    //                      || PCRS || pcrDigest)
    //  Start hash
    CryptHashStart(&hashState, session->authHashAlg);

    //  add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    //  add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    //  add PCRS
    buffer = pcrs;
    pcrSize = TPML_PCR_SELECTION_Marshal(&in->pcrs, &buffer, NULL);
    CryptDigestUpdate(&hashState, pcrSize, pcrs);

    //  add PCR digest
    CryptDigestUpdate2B(&hashState, &pcrDigest.b);

    //  complete the hash and get the results
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    //  update pcrCounter in session context for non trial session
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        session->pcrCounter = gr.pcrCounter;
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyPCR