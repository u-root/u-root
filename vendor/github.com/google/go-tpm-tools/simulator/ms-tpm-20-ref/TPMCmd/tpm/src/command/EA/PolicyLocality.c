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
#include "PolicyLocality_fp.h"


#if CC_PolicyLocality  // Conditional expansion of this file

//  Return Type: TPM_RC
//      TPM_RC_RANGE          all the locality values selected by
//                            'locality' have been disabled
//                            by previous TPM2_PolicyLocality() calls.
TPM_RC
TPM2_PolicyLocality(
    PolicyLocality_In   *in             // IN: input parameter list
    )
{
    SESSION     *session;
    BYTE         marshalBuffer[sizeof(TPMA_LOCALITY)];
    BYTE         prevSetting[sizeof(TPMA_LOCALITY)];
    UINT32       marshalSize;
    BYTE        *buffer;
    TPM_CC       commandCode = TPM_CC_PolicyLocality;
    HASH_STATE   hashState;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // Get new locality setting in canonical form
    marshalBuffer[0] = 0;   // Code analysis says that this is not initialized
    buffer = marshalBuffer;
    marshalSize = TPMA_LOCALITY_Marshal(&in->locality, &buffer, NULL);

    // Its an error if the locality parameter is zero
    if(marshalBuffer[0] == 0)
        return TPM_RCS_RANGE + RC_PolicyLocality_locality;

    // Get existing locality setting in canonical form
    prevSetting[0] = 0;     // Code analysis says that this is not initialized
    buffer = prevSetting;
    TPMA_LOCALITY_Marshal(&session->commandLocality, &buffer, NULL);

    // If the locality has previously been set
    if(prevSetting[0] != 0
        // then the current locality setting and the requested have to be the same
        // type (that is, either both normal or both extended
        && ((prevSetting[0] < 32) != (marshalBuffer[0] < 32)))
        return TPM_RCS_RANGE + RC_PolicyLocality_locality;

    // See if the input is a regular or extended locality
    if(marshalBuffer[0] < 32)
    {
        // if there was no previous setting, start with all normal localities
        // enabled
        if(prevSetting[0] == 0)
            prevSetting[0] = 0x1F;

        // AND the new setting with the previous setting and store it in prevSetting
        prevSetting[0] &= marshalBuffer[0];

        // The result setting can not be 0
        if(prevSetting[0] == 0)
            return TPM_RCS_RANGE + RC_PolicyLocality_locality;
    }
    else
    {
        // for extended locality
        // if the locality has already been set, then it must match the
        if(prevSetting[0] != 0 && prevSetting[0] != marshalBuffer[0])
            return TPM_RCS_RANGE + RC_PolicyLocality_locality;

        // Setting is OK
        prevSetting[0] = marshalBuffer[0];
    }

// Internal Data Update

    // Update policy hash
    // policyDigestnew = hash(policyDigestold || TPM_CC_PolicyLocality || locality)
    // Start hash
    CryptHashStart(&hashState, session->authHashAlg);

    // add old digest
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    // add commandCode
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    // add input locality
    CryptDigestUpdate(&hashState, marshalSize, marshalBuffer);

    // complete the digest
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    // update session locality by unmarshal function.  The function must succeed
    // because both input and existing locality setting have been validated.
    buffer = prevSetting;
    TPMA_LOCALITY_Unmarshal(&session->commandLocality, &buffer,
                            (INT32 *)&marshalSize);

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyLocality