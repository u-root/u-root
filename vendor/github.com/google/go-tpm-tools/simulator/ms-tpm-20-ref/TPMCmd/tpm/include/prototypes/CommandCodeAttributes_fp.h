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
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Mar 28, 2019  Time: 08:25:19PM
 */

#ifndef    _COMMAND_CODE_ATTRIBUTES_FP_H_
#define    _COMMAND_CODE_ATTRIBUTES_FP_H_

//*** GetClosestCommandIndex()
// This function returns the command index for the command with a value that is
// equal to or greater than the input value
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     command is not implemented
//  other                           index of a command
COMMAND_INDEX
GetClosestCommandIndex(
    TPM_CC           commandCode    // IN: the command code to start at
);

//*** CommandCodeToComandIndex()
// This function returns the index in the various attributes arrays of the
// command.
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     command is not implemented
//  other                           index of the command
COMMAND_INDEX
CommandCodeToCommandIndex(
    TPM_CC           commandCode    // IN: the command code to look up
);

//*** GetNextCommandIndex()
// This function returns the index of the next implemented command.
//  Return Type: COMMAND_INDEX
//  UNIMPLEMENTED_COMMAND_INDEX     no more implemented commands
//  other                           the index of the next implemented command
COMMAND_INDEX
GetNextCommandIndex(
    COMMAND_INDEX    commandIndex   // IN: the starting index
);

//*** GetCommandCode()
// This function returns the commandCode associated with the command index
TPM_CC
GetCommandCode(
    COMMAND_INDEX    commandIndex   // IN: the command index
);

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
);

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
);

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
);

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
);

//*** IsHandleInResponse()
// This function determines if a command has a handle in the response
BOOL
IsHandleInResponse(
    COMMAND_INDEX    commandIndex
);

//*** IsWriteOperation()
// Checks to see if an operation will write to an NV Index and is subject to being
// blocked by read-lock
BOOL
IsWriteOperation(
    COMMAND_INDEX    commandIndex   // IN: Command to check
);

//*** IsReadOperation()
// Checks to see if an operation will write to an NV Index and is
// subject to being blocked by write-lock.
BOOL
IsReadOperation(
    COMMAND_INDEX    commandIndex   // IN: Command to check
);

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
);

//*** IsVendorCommand()
// Function indicates if a command index references a vendor command.
//  Return Type: BOOL
//      TRUE(1)         command is a vendor command
//      FALSE(0)        command is not a vendor command
BOOL
IsVendorCommand(
    COMMAND_INDEX    commandIndex   // IN: command index to check
);

#endif  // _COMMAND_CODE_ATTRIBUTES_FP_H_
