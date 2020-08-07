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
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _AC_SPT_FP_H_
#define    _AC_SPT_FP_H_

//*** AcToCapabilities()
// This function returns a pointer to a list of AC capabilities.
TPML_AC_CAPABILITIES *
AcToCapabilities(
    TPMI_RH_AC      component       // IN: component
);

//*** AcIsAccessible()
// Function to determine if an AC handle references an actual AC
//  Return Type: BOOL
BOOL
AcIsAccessible(
    TPM_HANDLE          acHandle
);

//*** AcCapabilitiesGet()
// This function returns a list of capabilities associated with an AC
//  Return Type: TPMI_YES_NO
//      YES         if there are more handles available
//      NO          all the available handles has been returned
TPMI_YES_NO
AcCapabilitiesGet(
    TPMI_RH_AC               component,     // IN: the component
    TPM_AT                   type,          // IN: start capability type
    TPML_AC_CAPABILITIES    *capabilityList // OUT: list of handle
);

//*** AcSendObject()
// Stub to handle sending of an AC object
//  Return Type: TPM_RC
TPM_RC
AcSendObject(
    TPM_HANDLE           acHandle,      // IN: Handle of AC receiving object
    OBJECT              *object,        // IN: object structure to send
    TPMS_AC_OUTPUT      *acDataOut      // OUT: results of operation
);

#endif  // _AC_SPT_FP_H_
