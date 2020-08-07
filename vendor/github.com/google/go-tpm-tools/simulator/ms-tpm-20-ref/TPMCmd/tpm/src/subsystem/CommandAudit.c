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
// This file contains the functions that support command audit.

//** Includes
#include "Tpm.h"

//** Functions

//*** CommandAuditPreInstall_Init()
// This function initializes the command audit list. This function simulates
// the behavior of manufacturing. A function is used instead of a structure
// definition because this is easier than figuring out the initialization value
// for a bit array.
//
// This function would not be implemented outside of a manufacturing or
// simulation environment.
void
CommandAuditPreInstall_Init(
    void
    )
{
    // Clear all the audit commands
    MemorySet(gp.auditCommands, 0x00, sizeof(gp.auditCommands));

    // TPM_CC_SetCommandCodeAuditStatus always being audited
    CommandAuditSet(TPM_CC_SetCommandCodeAuditStatus);

    // Set initial command audit hash algorithm to be context integrity hash
    // algorithm
    gp.auditHashAlg = CONTEXT_INTEGRITY_HASH_ALG;

    // Set up audit counter to be 0
    gp.auditCounter = 0;

    // Write command audit persistent data to NV
    NV_SYNC_PERSISTENT(auditCommands);
    NV_SYNC_PERSISTENT(auditHashAlg);
    NV_SYNC_PERSISTENT(auditCounter);

    return;
}

//*** CommandAuditStartup()
// This function clears the command audit digest on a TPM Reset.
BOOL
CommandAuditStartup(
    STARTUP_TYPE     type           // IN: start up type
    )
{
    if((type != SU_RESTART) && (type != SU_RESUME))
    {
        // Reset the digest size to initialize the digest
        gr.commandAuditDigest.t.size = 0;
    }
    return TRUE;
}

//*** CommandAuditSet()
// This function will SET the audit flag for a command. This function
// will not SET the audit flag for a command that is not implemented. This
// ensures that the audit status is not SET when TPM2_GetCapability() is
// used to read the list of audited commands.
//
// This function is only used by TPM2_SetCommandCodeAuditStatus().
//
// The actions in TPM2_SetCommandCodeAuditStatus() are expected to cause the
// changes to be saved to NV after it is setting and clearing bits.
//  Return Type: BOOL
//      TRUE(1)         command code audit status was changed
//      FALSE(0)        command code audit status was not changed
BOOL
CommandAuditSet(
    TPM_CC           commandCode    // IN: command code
    )
{
    COMMAND_INDEX        commandIndex = CommandCodeToCommandIndex(commandCode);

    // Only SET a bit if the corresponding command is implemented
    if(commandIndex != UNIMPLEMENTED_COMMAND_INDEX)
    {
        // Can't audit shutdown
        if(commandCode != TPM_CC_Shutdown)
        {
            if(!TEST_BIT(commandIndex, gp.auditCommands))
            {
                // Set bit
                SET_BIT(commandIndex, gp.auditCommands);
                return TRUE;
            }
        }
    }
    // No change
    return FALSE;
}

//*** CommandAuditClear()
// This function will CLEAR the audit flag for a command. It will not CLEAR the
// audit flag for TPM_CC_SetCommandCodeAuditStatus().
//
// This function is only used by TPM2_SetCommandCodeAuditStatus().
//
// The actions in TPM2_SetCommandCodeAuditStatus() are expected to cause the
// changes to be saved to NV after it is setting and clearing bits.
//  Return Type: BOOL
//      TRUE(1)         command code audit status was changed
//      FALSE(0)        command code audit status was not changed
BOOL
CommandAuditClear(
    TPM_CC           commandCode    // IN: command code
    )
{
    COMMAND_INDEX       commandIndex = CommandCodeToCommandIndex(commandCode);

    // Do nothing if the command is not implemented
    if(commandIndex != UNIMPLEMENTED_COMMAND_INDEX)
    {
        // The bit associated with TPM_CC_SetCommandCodeAuditStatus() cannot be
        // cleared
        if(commandCode != TPM_CC_SetCommandCodeAuditStatus)
        {
            if(TEST_BIT(commandIndex, gp.auditCommands))
            {
                // Clear bit
                CLEAR_BIT(commandIndex, gp.auditCommands);
                return TRUE;
            }
        }
    }
    // No change
    return FALSE;
}

//*** CommandAuditIsRequired()
// This function indicates if the audit flag is SET for a command.
//  Return Type: BOOL
//      TRUE(1)         command is audited
//      FALSE(0)        command is not audited
BOOL
CommandAuditIsRequired(
    COMMAND_INDEX    commandIndex   // IN: command index
    )
{
    // Check the bit map.  If the bit is SET, command audit is required
    return(TEST_BIT(commandIndex, gp.auditCommands));
}

//*** CommandAuditCapGetCCList()
// This function returns a list of commands that have their audit bit SET.
//
// The list starts at the input commandCode.
//  Return Type: TPMI_YES_NO
//      YES         if there are more command code available
//      NO          all the available command code has been returned
TPMI_YES_NO
CommandAuditCapGetCCList(
    TPM_CC           commandCode,   // IN: start command code
    UINT32           count,         // IN: count of returned TPM_CC
    TPML_CC         *commandList    // OUT: list of TPM_CC
    )
{
    TPMI_YES_NO     more = NO;
    COMMAND_INDEX   commandIndex;

    // Initialize output handle list
    commandList->count = 0;

    // The maximum count of command we may return is MAX_CAP_CC
    if(count > MAX_CAP_CC) count = MAX_CAP_CC;

    // Find the implemented command that has a command code that is the same or
    // higher than the input
    // Collect audit commands
    for(commandIndex = GetClosestCommandIndex(commandCode);
    commandIndex != UNIMPLEMENTED_COMMAND_INDEX;
        commandIndex = GetNextCommandIndex(commandIndex))
    {
        if(CommandAuditIsRequired(commandIndex))
        {
            if(commandList->count < count)
            {
                // If we have not filled up the return list, add this command
                // code to its
                TPM_CC      cc = GET_ATTRIBUTE(s_ccAttr[commandIndex], 
                                               TPMA_CC, commandIndex);
                if(IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V))   
                    cc += (1 << 29);
                commandList->commandCodes[commandList->count] = cc;
                commandList->count++;
            }
            else
            {
                // If the return list is full but we still have command
                // available, report this and stop iterating
                more = YES;
                break;
            }
        }
    }

    return more;
}

//*** CommandAuditGetDigest
// This command is used to create a digest of the commands being audited. The
// commands are processed in ascending numeric order with a list of TPM_CC being
// added to a hash. This operates as if all the audited command codes were
// concatenated and then hashed.
void
CommandAuditGetDigest(
    TPM2B_DIGEST    *digest         // OUT: command digest
    )
{
    TPM_CC                       commandCode;
    COMMAND_INDEX                commandIndex;
    HASH_STATE                   hashState;

    // Start hash
    digest->t.size = CryptHashStart(&hashState, gp.auditHashAlg);

    // Add command code
    for(commandIndex = 0; commandIndex < COMMAND_COUNT; commandIndex++)
    {
        if(CommandAuditIsRequired(commandIndex))
        {
            commandCode = GetCommandCode(commandIndex);
            CryptDigestUpdateInt(&hashState, sizeof(commandCode), commandCode);
        }
    }

    // Complete hash
    CryptHashEnd2B(&hashState, &digest->b);

    return;
}