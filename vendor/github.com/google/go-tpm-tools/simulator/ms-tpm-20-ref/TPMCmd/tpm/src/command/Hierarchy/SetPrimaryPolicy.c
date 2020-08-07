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
#include "SetPrimaryPolicy_fp.h"

#if CC_SetPrimaryPolicy  // Conditional expansion of this file

/*(See part 3 specification)
// Set a hierarchy policy
*/
//  Return Type: TPM_RC
//      TPM_RC_SIZE           size of input authPolicy is not consistent with
//                            input hash algorithm
TPM_RC
TPM2_SetPrimaryPolicy(
    SetPrimaryPolicy_In     *in             // IN: input parameter list
    )
{
// Input Validation

    // Check the authPolicy consistent with hash algorithm. If the policy size is
    // zero, then the algorithm is required to be TPM_ALG_NULL
    if(in->authPolicy.t.size != CryptHashGetDigestSize(in->hashAlg))
        return TPM_RCS_SIZE + RC_SetPrimaryPolicy_authPolicy;

    // The command need NV update for OWNER and ENDORSEMENT hierarchy, and
    // might need orderlyState update for PLATFROM hierarchy.
    // Check if NV is available.  A TPM_RC_NV_UNAVAILABLE or TPM_RC_NV_RATE
    // error may be returned at this point
    RETURN_IF_NV_IS_NOT_AVAILABLE;

// Internal Data Update

    // Set hierarchy policy
    switch(in->authHandle)
    {
        case TPM_RH_OWNER:
            gp.ownerAlg = in->hashAlg;
            gp.ownerPolicy = in->authPolicy;
            NV_SYNC_PERSISTENT(ownerAlg);
            NV_SYNC_PERSISTENT(ownerPolicy);
            break;
        case TPM_RH_ENDORSEMENT:
            gp.endorsementAlg = in->hashAlg;
            gp.endorsementPolicy = in->authPolicy;
            NV_SYNC_PERSISTENT(endorsementAlg);
            NV_SYNC_PERSISTENT(endorsementPolicy);
            break;
        case TPM_RH_PLATFORM:
            gc.platformAlg = in->hashAlg;
            gc.platformPolicy = in->authPolicy;
            // need to update orderly state
            g_clearOrderly = TRUE;
            break;
        case TPM_RH_LOCKOUT:
            gp.lockoutAlg = in->hashAlg;
            gp.lockoutPolicy = in->authPolicy;
            NV_SYNC_PERSISTENT(lockoutAlg);
            NV_SYNC_PERSISTENT(lockoutPolicy);
            break;

        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_SetPrimaryPolicy