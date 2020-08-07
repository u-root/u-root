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

#ifndef    _SESSION_PROCESS_FP_H_
#define    _SESSION_PROCESS_FP_H_

//*** IsDAExempted()
// This function indicates if a handle is exempted from DA logic.
// A handle is exempted if it is
//  1. a primary seed handle,
//  2. an object with noDA bit SET,
//  3. an NV Index with TPMA_NV_NO_DA bit SET, or
//  4. a PCR handle.
//
//  Return Type: BOOL
//      TRUE(1)         handle is exempted from DA logic
//      FALSE(0)        handle is not exempted from DA logic
BOOL
IsDAExempted(
    TPM_HANDLE       handle         // IN: entity handle
);

//*** ClearCpRpHashes()
void
ClearCpRpHashes(
    COMMAND         *command
);

//*** CompareNameHash()
// This function computes the name hash and compares it to the nameHash in the
// session data.
BOOL
CompareNameHash(
    COMMAND         *command,       // IN: main parsing structure
    SESSION         *session        // IN: session structure with nameHash
);

//*** ParseSessionBuffer()
// This function is the entry function for command session processing.
// It iterates sessions in session area and reports if the required authorization
// has been properly provided. It also processes audit session and passes the
// information of encryption sessions to parameter encryption module.
//
//  Return Type: TPM_RC
//        various           parsing failure or authorization failure
//
TPM_RC
ParseSessionBuffer(
    COMMAND         *command        // IN: the structure that contains
);

//*** CheckAuthNoSession()
// Function to process a command with no session associated.
// The function makes sure all the handles in the command require no authorization.
//
//  Return Type: TPM_RC
//      TPM_RC_AUTH_MISSING         failure - one or more handles require
//                                  authorization
TPM_RC
CheckAuthNoSession(
    COMMAND         *command        // IN: command parsing structure
);

//*** BuildResponseSession()
// Function to build Session buffer in a response. The authorization data is added
// to the end of command->responseBuffer. The size of the authorization area is
// accumulated in command->authSize.
// When this is called, command->responseBuffer is pointing at the next location
// in the response buffer to be filled. This is where the authorization sessions
// will go, if any. command->parameterSize is the number of bytes that have been
// marshaled as parameters in the output buffer.
void
BuildResponseSession(
    COMMAND         *command        // IN: structure that has relevant command
                                    //     information
);

//*** SessionRemoveAssociationToHandle()
// This function deals with the case where an entity associated with an authorization
// is deleted during command processing. The primary use of this is to support
// UndefineSpaceSpecial().
void
SessionRemoveAssociationToHandle(
    TPM_HANDLE       handle
);

#endif  // _SESSION_PROCESS_FP_H_
