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
#include "HierarchyChangeAuth_fp.h"

#if CC_HierarchyChangeAuth  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// Set a hierarchy authValue
*/
//  Return Type: TPM_RC
//      TPM_RC_SIZE        'newAuth' size is greater than that of integrity hash 
//                          digest
TPM_RC
TPM2_HierarchyChangeAuth(
    HierarchyChangeAuth_In  *in             // IN: input parameter list
    )
{
    // The command needs NV update.
    RETURN_IF_NV_IS_NOT_AVAILABLE;

    // Make sure that the authorization value is a reasonable size (not larger than
    // the size of the digest produced by the integrity hash. The integrity
    // hash is assumed to produce the longest digest of any hash implemented
    // on the TPM. This will also remove trailing zeros from the authValue.
    if(MemoryRemoveTrailingZeros(&in->newAuth) > CONTEXT_INTEGRITY_HASH_SIZE)
        return TPM_RCS_SIZE + RC_HierarchyChangeAuth_newAuth;

    // Set hierarchy authValue
    switch(in->authHandle)
    {
        case TPM_RH_OWNER:
            gp.ownerAuth = in->newAuth;
            NV_SYNC_PERSISTENT(ownerAuth);
            break;
        case TPM_RH_ENDORSEMENT:
            gp.endorsementAuth = in->newAuth;
            NV_SYNC_PERSISTENT(endorsementAuth);
            break;
        case TPM_RH_PLATFORM:
            gc.platformAuth = in->newAuth;
            // orderly state should be cleared
            g_clearOrderly = TRUE;
            break;
        case TPM_RH_LOCKOUT:
            gp.lockoutAuth = in->newAuth;
            NV_SYNC_PERSISTENT(lockoutAuth);
            break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_HierarchyChangeAuth