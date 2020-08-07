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
 *  Date: Apr  2, 2019  Time: 03:18:00PM
 */

#ifndef    _TPM_FAIL_FP_H_
#define    _TPM_FAIL_FP_H_

//*** SetForceFailureMode()
// This function is called by the simulator to enable failure mode testing.
#if SIMULATION
LIB_EXPORT void
SetForceFailureMode(
    void
);
#endif

//*** TpmLogFailure()
// This function saves the failure values when the code will continue to operate. It
// if similar to TpmFail() but returns to the caller. The assumption is that the
// caller will propagate a failure back up the stack.
void
TpmLogFailure(
#if FAIL_TRACE
    const char      *function,
    int              line,
#endif
    int              code
);

//*** TpmFail()
// This function is called by TPM.lib when a failure occurs. It will set up the
// failure values to be returned on TPM2_GetTestResult().
NORETURN void
TpmFail(
#if FAIL_TRACE
    const char      *function,
    int              line,
#endif
    int              code
);

//*** TpmFailureMode(
// This function is called by the interface code when the platform is in failure
// mode.
void
TpmFailureMode(
    unsigned int     inRequestSize,     // IN: command buffer size
    unsigned char   *inRequest,         // IN: command buffer
    unsigned int    *outResponseSize,   // OUT: response buffer size
    unsigned char   **outResponse       // OUT: response buffer
);

//*** UnmarshalFail()
// This is a stub that is used to catch an attempt to unmarshal an entry
// that is not defined. Don't ever expect this to be called but...
void
UnmarshalFail(
    void            *type,
    BYTE            **buffer,
    INT32           *size
);

#endif  // _TPM_FAIL_FP_H_
