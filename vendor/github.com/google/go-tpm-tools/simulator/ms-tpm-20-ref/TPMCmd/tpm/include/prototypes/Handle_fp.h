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

#ifndef    _HANDLE_FP_H_
#define    _HANDLE_FP_H_

//*** HandleGetType()
// This function returns the type of a handle which is the MSO of the handle.
TPM_HT
HandleGetType(
    TPM_HANDLE       handle         // IN: a handle to be checked
);

//*** NextPermanentHandle()
// This function returns the permanent handle that is equal to the input value or
// is the next higher value. If there is no handle with the input value and there
// is no next higher value, it returns 0:
TPM_HANDLE
NextPermanentHandle(
    TPM_HANDLE       inHandle       // IN: the handle to check
);

//*** PermanentCapGetHandles()
// This function returns a list of the permanent handles of PCR, started from
// 'handle'. If 'handle' is larger than the largest permanent handle, an empty list
// will be returned with 'more' set to NO.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
PermanentCapGetHandles(
    TPM_HANDLE       handle,        // IN: start handle
    UINT32           count,         // IN: count of returned handles
    TPML_HANDLE     *handleList     // OUT: list of handle
);

//*** PermanentHandleGetPolicy()
// This function returns a list of the permanent handles of PCR, started from
// 'handle'. If 'handle' is larger than the largest permanent handle, an empty list
// will be returned with 'more' set to NO.
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
PermanentHandleGetPolicy(
    TPM_HANDLE           handle,        // IN: start handle
    UINT32               count,         // IN: max count of returned handles
    TPML_TAGGED_POLICY  *policyList     // OUT: list of handle
);

#endif  // _HANDLE_FP_H_
