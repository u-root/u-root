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
#include "SetCommandCodeAuditStatus_fp.h"

#if CC_SetCommandCodeAuditStatus  // Conditional expansion of this file

/*(See part 3 specification)
// change the audit status of a command or to set the hash algorithm used for
// the audit digest.
*/
TPM_RC
TPM2_SetCommandCodeAuditStatus(
    SetCommandCodeAuditStatus_In    *in             // IN: input parameter list
    )
{

    // The command needs NV update.  Check if NV is available.
    // A TPM_RC_NV_UNAVAILABLE or TPM_RC_NV_RATE error may be returned at
    // this point
    RETURN_IF_NV_IS_NOT_AVAILABLE;

// Internal Data Update

    // Update hash algorithm
    if(in->auditAlg != TPM_ALG_NULL && in->auditAlg != gp.auditHashAlg)
    {
        // Can't change the algorithm and command list at the same time
        if(in->setList.count != 0 || in->clearList.count != 0)
            return TPM_RCS_VALUE + RC_SetCommandCodeAuditStatus_auditAlg;

        // Change the hash algorithm for audit
        gp.auditHashAlg = in->auditAlg;

        // Set the digest size to a unique value that indicates that the digest
        // algorithm has been changed. The size will be cleared to zero in the
        // command audit processing on exit.
        gr.commandAuditDigest.t.size = 1;

        // Save the change of command audit data (this sets g_updateNV so that NV
        // will be updated on exit.)
        NV_SYNC_PERSISTENT(auditHashAlg);
    }
    else
    {
        UINT32          i;
        BOOL            changed = FALSE;

        // Process set list
        for(i = 0; i < in->setList.count; i++)

            // If change is made in CommandAuditSet, set changed flag
            if(CommandAuditSet(in->setList.commandCodes[i]))
                changed = TRUE;

        // Process clear list
        for(i = 0; i < in->clearList.count; i++)
            // If change is made in CommandAuditClear, set changed flag
            if(CommandAuditClear(in->clearList.commandCodes[i]))
                changed = TRUE;

        // if change was made to command list, update NV
        if(changed)
            // this sets g_updateNV so that NV will be updated on exit.
            NV_SYNC_PERSISTENT(auditCommands);
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_SetCommandCodeAuditStatus