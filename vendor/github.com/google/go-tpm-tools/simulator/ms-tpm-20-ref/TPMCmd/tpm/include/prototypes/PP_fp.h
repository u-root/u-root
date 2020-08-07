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

#ifndef    _PP_FP_H_
#define    _PP_FP_H_

//*** PhysicalPresencePreInstall_Init()
// This function is used to initialize the array of commands that always require
// confirmation with physical presence. The array is an array of bits that
// has a correspondence with the command code.
//
// This command should only ever be executable in a manufacturing setting or in
// a simulation.
//
// When set, these cannot be cleared.
//
void
PhysicalPresencePreInstall_Init(
    void
);

//*** PhysicalPresenceCommandSet()
// This function is used to set the indicator that a command requires
// PP confirmation.
void
PhysicalPresenceCommandSet(
    TPM_CC           commandCode    // IN: command code
);

//*** PhysicalPresenceCommandClear()
// This function is used to clear the indicator that a command requires PP
// confirmation.
void
PhysicalPresenceCommandClear(
    TPM_CC           commandCode    // IN: command code
);

//*** PhysicalPresenceIsRequired()
// This function indicates if PP confirmation is required for a command.
//  Return Type: BOOL
//      TRUE(1)         physical presence is required
//      FALSE(0)        physical presence is not required
BOOL
PhysicalPresenceIsRequired(
    COMMAND_INDEX    commandIndex   // IN: command index
);

//*** PhysicalPresenceCapGetCCList()
// This function returns a list of commands that require PP confirmation. The
// list starts from the first implemented command that has a command code that
// the same or greater than 'commandCode'.
//  Return Type: TPMI_YES_NO
//      YES         if there are more command codes available
//      NO          all the available command codes have been returned
TPMI_YES_NO
PhysicalPresenceCapGetCCList(
    TPM_CC           commandCode,   // IN: start command code
    UINT32           count,         // IN: count of returned TPM_CC
    TPML_CC         *commandList    // OUT: list of TPM_CC
);

#endif  // _PP_FP_H_
