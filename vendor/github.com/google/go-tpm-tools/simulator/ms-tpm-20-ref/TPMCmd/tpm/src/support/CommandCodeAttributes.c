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
// This file contains the functions for testing various command properties.

//** Includes and Defines

#include "Tpm.h"
#include "CommandCodeAttributes_fp.h"

// Set the default value for CC_VEND if not already set
#ifndef CC_VEND
#define     CC_VEND     (TPM_CC)(0x20000000)
#endif

typedef UINT16          ATTRIBUTE_TYPE;

// The following file is produced from the command tables in part 3 of the
// specification. It defines the attributes for each of the commands.
// NOTE: This file is currently produced by an automated process. Files
// produced from Part 2 or Part 3 tables through automated processes are not
// included in the specification so that their is no ambiguity about the
// table containing the information being the normative definition.
#define _COMMAND_CODE_ATTRIBUTES_
#include    "CommandAttributeData.h"

//** Command Attribute Functions

//*** NextImplementedIndex()
// This function is used when the lists are not compressed. In a compressed list,
// only the implemented commands are present. So, a search might find a value
// but that value may not be implemented. This function checks to see if the input
// commandIndex points to an implemented command and, if not, it searches upwards
// until it finds one. When the list is compressed, this function gets defined
// as a no-op.
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     command is not implemented
//  other                           index of the command
#if !COMPRESSED_LISTS
static COMMAND_INDEX
NextImplementedIndex(
    COMMAND_INDEX       commandIndex
    )
{
    for(;commandIndex < COMMAND_COUNT; commandIndex++)
    {
        if(s_commandAttributes[commandIndex] & IS_IMPLEMENTED)
            return commandIndex;
    }
    return UNIMPLEMENTED_COMMAND_INDEX;
}
#else
#define NextImplementedIndex(x) (x)
#endif

//*** GetClosestCommandIndex()
// This function returns the command index for the command with a value that is
// equal to or greater than the input value
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     command is not implemented
//  other                           index of a command
COMMAND_INDEX
GetClosestCommandIndex(
    TPM_CC           commandCode    // IN: the command code to start at
    )
{
    BOOL                vendor = (commandCode & CC_VEND) != 0;
    COMMAND_INDEX       searchIndex = (COMMAND_INDEX)commandCode;

    // The commandCode is a UINT32 and the search index is UINT16. We are going to
    // search for a match but need to make sure that the commandCode value is not
    // out of range. To do this, need to clear the vendor bit of the commandCode
    // (if set) and compare the result to the 16-bit searchIndex value. If it is 
    // out of range, indicate that the command is not implemented
    if((commandCode & ~CC_VEND) != searchIndex)
        return UNIMPLEMENTED_COMMAND_INDEX;

    // if there is at least one vendor command, the last entry in the array will
    // have the v bit set. If the input commandCode is larger than the last
    // vendor-command, then it is out of range.
    if(vendor)
    {
#if VENDOR_COMMAND_ARRAY_SIZE > 0
        COMMAND_INDEX       commandIndex;
        COMMAND_INDEX       min;
        COMMAND_INDEX       max;
        int                 diff;
#if LIBRARY_COMMAND_ARRAY_SIZE == COMMAND_COUNT
#error "Constants are not consistent."
#endif
        // Check to see if the value is equal to or below the minimum
        // entry.
        // Note: Put this check first so that the typical case of only one vendor-
        // specific command doesn't waste any more time.
        if(GET_ATTRIBUTE(s_ccAttr[LIBRARY_COMMAND_ARRAY_SIZE], TPMA_CC, 
                         commandIndex) >= searchIndex)
        {
            // the vendor array is always assumed to be packed so there is
            // no need to check to see if the command is implemented
            return LIBRARY_COMMAND_ARRAY_SIZE;
        }
        // See if this is out of range on the top
        if(GET_ATTRIBUTE(s_ccAttr[COMMAND_COUNT - 1], TPMA_CC, commandIndex) 
           < searchIndex)
        {
            return UNIMPLEMENTED_COMMAND_INDEX;
        }
        commandIndex = UNIMPLEMENTED_COMMAND_INDEX; // Needs initialization to keep
                                                    // compiler happy
        min = LIBRARY_COMMAND_ARRAY_SIZE;       // first vendor command
        max = COMMAND_COUNT - 1;                // last vendor command
        diff = 1;                               // needs initialization to keep
                                                // compiler happy
        while(min <= max)
        {
            commandIndex = (min + max + 1) / 2;
            diff = GET_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, commandIndex)
                - searchIndex;
            if(diff == 0)
                return commandIndex;
            if(diff > 0)
                max = commandIndex - 1;
            else
                min = commandIndex + 1;
        }
        // didn't find and exact match. commandIndex will be pointing at the last
        // item tested. If 'diff' is positive, then the last item tested was
        // larger index of the command code so it is the smallest value
        // larger than the requested value.
        if(diff > 0)
            return commandIndex;
        // if 'diff' is negative, then the value tested was smaller than
        // the commandCode index and the next higher value is the correct one.
        // Note: this will necessarily be in range because of the earlier check
        // that the index was within range.
        return commandIndex + 1;
#else
    // If there are no vendor commands so anything with the vendor bit set is out
    // of range
        return UNIMPLEMENTED_COMMAND_INDEX;
#endif
    }
    // Get here if the V-Bit was not set in 'commandCode'

    if(GET_ATTRIBUTE(s_ccAttr[LIBRARY_COMMAND_ARRAY_SIZE - 1], TPMA_CC, 
                     commandIndex) < searchIndex)
    {
        // requested index is out of the range to the top
#if VENDOR_COMMAND_ARRAY_SIZE > 0
        // If there are vendor commands, then the first vendor command
        // is the next value greater than the commandCode.
        // NOTE: we got here if the starting index did not have the V bit but we
        // reached the end of the array of library commands (non-vendor). Since
        // there is at least one vendor command, and vendor commands are always
        // in a compressed list that starts after the library list, the next
        // index value contains a valid vendor command.
        return LIBRARY_COMMAND_ARRAY_SIZE;
#else
        // if there are no vendor commands, then this is out of range
        return UNIMPLEMENTED_COMMAND_INDEX;
#endif
    }
    // If the request is lower than any value in the array, then return
    // the lowest value (needs to be an index for an implemented command
    if(GET_ATTRIBUTE(s_ccAttr[0], TPMA_CC, commandIndex) >= searchIndex)
    {
        return NextImplementedIndex(0);
    }
    else
    {
#if COMPRESSED_LISTS
        COMMAND_INDEX       commandIndex = UNIMPLEMENTED_COMMAND_INDEX;
        COMMAND_INDEX       min = 0;
        COMMAND_INDEX       max = LIBRARY_COMMAND_ARRAY_SIZE - 1;
        int                 diff = 1;
#if LIBRARY_COMMAND_ARRAY_SIZE == 0
#error  "Something is terribly wrong"
#endif
        // The s_ccAttr array contains an extra entry at the end (a zero value).
        // Don't count this as an array entry. This means that max should start
        // out pointing to the last valid entry in the array which is - 2
        pAssert(max == (sizeof(s_ccAttr) / sizeof(TPMA_CC)
                        - VENDOR_COMMAND_ARRAY_SIZE - 2));
        while(min <= max)
        {
            commandIndex = (min + max + 1) / 2;
            diff = GET_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, 
                                 commandIndex) - searchIndex;
            if(diff == 0)
                return commandIndex;
            if(diff > 0)
                max = commandIndex - 1;
            else
                min = commandIndex + 1;
        }
        // didn't find and exact match. commandIndex will be pointing at the
        // last item tested. If diff is positive, then the last item tested was
        // larger index of the command code so it is the smallest value
        // larger than the requested value.
        if(diff > 0)
            return commandIndex;
        // if diff is negative, then the value tested was smaller than
        // the commandCode index and the next higher value is the correct one.
        // Note: this will necessarily be in range because of the earlier check
        // that the index was within range.
        return commandIndex + 1;
#else
        // The list is not compressed so offset into the array by the command
        // code value of the first entry in the list. Then go find the first
        // implemented command.
        return NextImplementedIndex(searchIndex
                                    - (COMMAND_INDEX)s_ccAttr[0].commandIndex);
#endif
    }
}

//*** CommandCodeToComandIndex()
// This function returns the index in the various attributes arrays of the
// command.
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     command is not implemented
//  other                           index of the command
COMMAND_INDEX
CommandCodeToCommandIndex(
    TPM_CC           commandCode    // IN: the command code to look up
    )
{
    // Extract the low 16-bits of the command code to get the starting search index
    COMMAND_INDEX       searchIndex = (COMMAND_INDEX)commandCode;
    BOOL                vendor = (commandCode & CC_VEND) != 0;
    COMMAND_INDEX       commandIndex;
#if !COMPRESSED_LISTS
    if(!vendor)
    {
        commandIndex = searchIndex - (COMMAND_INDEX)s_ccAttr[0].commandIndex;
        // Check for out of range or unimplemented.
        // Note, since a COMMAND_INDEX is unsigned, if searchIndex is smaller than
        // the lowest value of command, it will become a 'negative' number making
        // it look like a large unsigned number, this will cause it to fail
        // the unsigned check below.
        if(commandIndex >= LIBRARY_COMMAND_ARRAY_SIZE
           || (s_commandAttributes[commandIndex] & IS_IMPLEMENTED) == 0)
            return UNIMPLEMENTED_COMMAND_INDEX;
        return commandIndex;
    }
#endif
    // Need this code for any vendor code lookup or for compressed lists
    commandIndex = GetClosestCommandIndex(commandCode);

    // Look at the returned value from get closest. If it isn't the one that was
    // requested, then the command is not implemented.
    if(commandIndex != UNIMPLEMENTED_COMMAND_INDEX)
    {
        if((GET_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, commandIndex) 
            != searchIndex)
           || (IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V)) != vendor)
            commandIndex = UNIMPLEMENTED_COMMAND_INDEX;
    }
    return commandIndex;
}

//*** GetNextCommandIndex()
// This function returns the index of the next implemented command.
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     no more implemented commands
//  other                           the index of the next implemented command
COMMAND_INDEX
GetNextCommandIndex(
    COMMAND_INDEX    commandIndex   // IN: the starting index
    )
{
    while(++commandIndex < COMMAND_COUNT)
    {
#if !COMPRESSED_LISTS
        if(s_commandAttributes[commandIndex] & IS_IMPLEMENTED)
#endif
            return commandIndex;
    }
    return UNIMPLEMENTED_COMMAND_INDEX;
}

//*** GetCommandCode()
// This function returns the commandCode associated with the command index
TPM_CC
GetCommandCode(
    COMMAND_INDEX    commandIndex   // IN: the command index
    )
{
    TPM_CC           commandCode = GET_ATTRIBUTE(s_ccAttr[commandIndex],
                                                 TPMA_CC, commandIndex);
    if(IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V))
        commandCode += CC_VEND;
    return commandCode;
}

//*** CommandAuthRole()
//
//  This function returns the authorization role required of a handle.
//
//  Return Type: AUTH_ROLE
//  AUTH_NONE       no authorization is required
//  AUTH_USER       user role authorization is required
//  AUTH_ADMIN      admin role authorization is required
//  AUTH_DUP        duplication role authorization is required
AUTH_ROLE
CommandAuthRole(
    COMMAND_INDEX    commandIndex,  // IN: command index
    UINT32           handleIndex    // IN: handle index (zero based)
    )
{
    if(0 == handleIndex)
    {
        // Any authorization role set?
        COMMAND_ATTRIBUTES  properties = s_commandAttributes[commandIndex];

        if(properties & HANDLE_1_USER)
            return AUTH_USER;
        if(properties & HANDLE_1_ADMIN)
            return AUTH_ADMIN;
        if(properties & HANDLE_1_DUP)
            return AUTH_DUP;
    }
    else if(1 == handleIndex)
    {
        if(s_commandAttributes[commandIndex] & HANDLE_2_USER)
            return AUTH_USER;
    }
    return AUTH_NONE;
}

//*** EncryptSize()
// This function returns the size of the decrypt size field. This function returns
// 0 if encryption is not allowed
//  Return Type: int
//  0       encryption not allowed
//  2       size field is two bytes
//  4       size field is four bytes
int
EncryptSize(
    COMMAND_INDEX    commandIndex   // IN: command index
    )
{
    return ((s_commandAttributes[commandIndex] & ENCRYPT_2) ? 2 : 
            (s_commandAttributes[commandIndex] & ENCRYPT_4) ? 4 : 0);
}

//*** DecryptSize()
// This function returns the size of the decrypt size field. This function returns
// 0 if decryption is not allowed
//  Return Type: int
//  0       encryption not allowed
//  2       size field is two bytes
//  4       size field is four bytes
int
DecryptSize(
    COMMAND_INDEX    commandIndex   // IN: command index
    )
{
    return ((s_commandAttributes[commandIndex] & DECRYPT_2) ? 2 : 
            (s_commandAttributes[commandIndex] & DECRYPT_4) ? 4 : 0);
}

//*** IsSessionAllowed()
//
// This function indicates if the command is allowed to have sessions.
//
// This function must not be called if the command is not known to be implemented.
//
//  Return Type: BOOL
//      TRUE(1)         session is allowed with this command
//      FALSE(0)        session is not allowed with this command
BOOL
IsSessionAllowed(
    COMMAND_INDEX    commandIndex   // IN: the command to be checked
    )
{
    return ((s_commandAttributes[commandIndex] & NO_SESSIONS) == 0);
}

//*** IsHandleInResponse()
// This function determines if a command has a handle in the response
BOOL
IsHandleInResponse(
    COMMAND_INDEX    commandIndex
    )
{
    return ((s_commandAttributes[commandIndex] & R_HANDLE) != 0);
}

//*** IsWriteOperation()
// Checks to see if an operation will write to an NV Index and is subject to being
// blocked by read-lock
BOOL
IsWriteOperation(
    COMMAND_INDEX    commandIndex   // IN: Command to check
    )
{
#ifdef  WRITE_LOCK
    return ((s_commandAttributes[commandIndex] & WRITE_LOCK) != 0);
#else
    if(!IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V))
    {
        switch(GET_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, commandIndex))
        {
            case TPM_CC_NV_Write:
#if CC_NV_Increment
            case TPM_CC_NV_Increment:
#endif
#if CC_NV_SetBits
            case TPM_CC_NV_SetBits:
#endif
#if CC_NV_Extend
            case TPM_CC_NV_Extend:
#endif
#if CC_AC_Send
            case TPM_CC_AC_Send:
#endif
            // NV write lock counts as a write operation for authorization purposes.
            // We check to see if the NV is write locked before we do the
            // authorization. If it is locked, we fail the command early.
            case TPM_CC_NV_WriteLock:
                return TRUE;
            default:
                break;
        }
    }
    return FALSE;
#endif
}

//*** IsReadOperation()
// Checks to see if an operation will write to an NV Index and is
// subject to being blocked by write-lock.
BOOL
IsReadOperation(
    COMMAND_INDEX    commandIndex   // IN: Command to check
    )
{
#ifdef  READ_LOCK
    return ((s_commandAttributes[commandIndex] & READ_LOCK) != 0);
#else

    if(!IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V))
    {
        switch(GET_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, commandIndex))
        {
            case TPM_CC_NV_Read:
            case TPM_CC_PolicyNV:
            case TPM_CC_NV_Certify:
            // NV read lock counts as a read operation for authorization purposes.
            // We check to see if the NV is read locked before we do the
            // authorization. If it is locked, we fail the command early.
            case TPM_CC_NV_ReadLock:
                return TRUE;
            default:
                break;
        }
    }
    return FALSE;
#endif
}

//*** CommandCapGetCCList()
// This function returns a list of implemented commands and command attributes
// starting from the command in 'commandCode'.
//  Return Type: TPMI_YES_NO
//      YES         more command attributes are available
//      NO          no more command attributes are available
TPMI_YES_NO
CommandCapGetCCList(
    TPM_CC           commandCode,   // IN: start command code
    UINT32           count,         // IN: maximum count for number of entries in
                                    //     'commandList'
    TPML_CCA        *commandList    // OUT: list of TPMA_CC
    )
{
    TPMI_YES_NO      more = NO;
    COMMAND_INDEX    commandIndex;

    // initialize output handle list count
    commandList->count = 0;

    for(commandIndex = GetClosestCommandIndex(commandCode);
    commandIndex != UNIMPLEMENTED_COMMAND_INDEX;
        commandIndex = GetNextCommandIndex(commandIndex))
    {
#if !COMPRESSED_LISTS
        // this check isn't needed for compressed lists.
        if(!(s_commandAttributes[commandIndex] & IS_IMPLEMENTED))
            continue;
#endif
        if(commandList->count < count)
        {
            // If the list is not full, add the attributes for this command.
            commandList->commandAttributes[commandList->count]
                = s_ccAttr[commandIndex];
            commandList->count++;
        }
        else
        {
            // If the list is full but there are more commands to report,
            // indicate this and return.
            more = YES;
            break;
        }
    }
    return more;
}

//*** IsVendorCommand()
// Function indicates if a command index references a vendor command.
//  Return Type: BOOL
//      TRUE(1)         command is a vendor command
//      FALSE(0)        command is not a vendor command
BOOL
IsVendorCommand(
    COMMAND_INDEX    commandIndex   // IN: command index to check
    )
{
    return (IS_ATTRIBUTE(s_ccAttr[commandIndex], TPMA_CC, V));
}
