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
#include "PP_Commands_fp.h"

#if CC_PP_Commands  // Conditional expansion of this file

/*(See part 3 specification)
// This command is used to determine which commands require assertion of
// Physical Presence in addition to platformAuth/platformPolicy.
*/
TPM_RC
TPM2_PP_Commands(
    PP_Commands_In  *in             // IN: input parameter list
    )
{
    UINT32          i;

    // The command needs NV update.  Check if NV is available.
    // A TPM_RC_NV_UNAVAILABLE or TPM_RC_NV_RATE error may be returned at
    // this point
    RETURN_IF_NV_IS_NOT_AVAILABLE;

// Internal Data Update

    // Process set list
    for(i = 0; i < in->setList.count; i++)
        // If command is implemented, set it as PP required.  If the input
        // command is not a PP command, it will be ignored at
        // PhysicalPresenceCommandSet().
        // Note: PhysicalPresenceCommandSet() checks if the command is implemented.
        PhysicalPresenceCommandSet(in->setList.commandCodes[i]);

    // Process clear list
    for(i = 0; i < in->clearList.count; i++)
        // If command is implemented, clear it as PP required.  If the input
        // command is not a PP command, it will be ignored at
        // PhysicalPresenceCommandClear().  If the input command is
        // TPM2_PP_Commands, it will be ignored as well
        PhysicalPresenceCommandClear(in->clearList.commandCodes[i]);

    // Save the change of PP list
    NV_SYNC_PERSISTENT(ppList);

    return TPM_RC_SUCCESS;
}

#endif // CC_PP_Commands