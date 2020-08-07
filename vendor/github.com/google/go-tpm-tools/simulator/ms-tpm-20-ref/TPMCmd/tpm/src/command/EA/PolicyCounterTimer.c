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
#include "PolicyCounterTimer_fp.h"

#if CC_PolicyCounterTimer  // Conditional expansion of this file

#include "Policy_spt_fp.h"

/*(See part 3 specification)
// Add a conditional gating of a policy based on the contents of the
// TPMS_TIME_INFO structure.
*/
//  Return Type: TPM_RC
//      TPM_RC_POLICY           the comparison of the selected portion of the
//                              TPMS_TIME_INFO with 'operandB' failed
//      TPM_RC_RANGE            'offset' + 'size' exceed size of TPMS_TIME_INFO
//                              structure
TPM_RC
TPM2_PolicyCounterTimer(
    PolicyCounterTimer_In   *in             // IN: input parameter list
    )
{
    SESSION             *session;
    TIME_INFO            infoData;          // data buffer of  TPMS_TIME_INFO
    BYTE                *pInfoData = (BYTE *)&infoData;
    UINT16               infoDataSize;
    TPM_CC               commandCode = TPM_CC_PolicyCounterTimer;
    HASH_STATE           hashState;
    TPM2B_DIGEST         argHash;

// Input Validation
    // Get a marshaled time structure
    infoDataSize = TimeGetMarshaled(&infoData);
    // Make sure that the referenced stays within the bounds of the structure.
    // NOTE: the offset checks are made even for a trial policy because the policy
    // will not make any sense if the references are out of bounds of the timer
    // structure.
    if(in->offset > infoDataSize)
        return TPM_RCS_VALUE + RC_PolicyCounterTimer_offset;
    if((UINT32)in->offset + (UINT32)in->operandB.t.size > infoDataSize)
        return TPM_RCS_RANGE;
    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    //If this is a trial policy, skip the check to see if the condition is met.
    if(session->attributes.isTrialPolicy == CLEAR)
    {
        // If the command is going to use any part of the counter or timer, need
        // to verify that time is advancing.
        // The time and clock vales are the first two 64-bit values in the clock
        if(in->offset < sizeof(UINT64) + sizeof(UINT64))
        {
            // Using Clock or Time so see if clock is running. Clock doesn't 
            // run while NV is unavailable.
            // TPM_RC_NV_UNAVAILABLE or TPM_RC_NV_RATE error may be returned here.
            RETURN_IF_NV_IS_NOT_AVAILABLE;
        }
        // offset to the starting position
        pInfoData = (BYTE *)infoData;
        // Check to see if the condition is valid
        if(!PolicySptCheckCondition(in->operation, pInfoData + in->offset,
                                    in->operandB.t.buffer, in->operandB.t.size))
            return TPM_RC_POLICY;
    }
// Internal Data Update
    // Start argument list hash
    argHash.t.size = CryptHashStart(&hashState, session->authHashAlg);
    //  add operandB
    CryptDigestUpdate2B(&hashState, &in->operandB.b);
    //  add offset
    CryptDigestUpdateInt(&hashState, sizeof(UINT16), in->offset);
    //  add operation
    CryptDigestUpdateInt(&hashState, sizeof(TPM_EO), in->operation);
    //  complete argument hash
    CryptHashEnd2B(&hashState, &argHash.b);

    // update policyDigest
    //  start hash
    CryptHashStart(&hashState, session->authHashAlg);

    //  add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    //  add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    //  add argument digest
    CryptDigestUpdate2B(&hashState, &argHash.b);

    // complete the digest
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyCounterTimer