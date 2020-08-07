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
 *  Date: Apr  2, 2019  Time: 04:23:27PM
 */

#ifndef    _DA_FP_H_
#define    _DA_FP_H_

//*** DAPreInstall_Init()
// This function initializes the DA parameters to their manufacturer-default
// values. The default values are determined by a platform-specific specification.
//
// This function should not be called outside of a manufacturing or simulation
// environment.
//
// The DA parameters will be restored to these initial values by TPM2_Clear().
void
DAPreInstall_Init(
    void
);

//*** DAStartup()
// This function is called  by TPM2_Startup() to initialize the DA parameters.
// In the case of Startup(CLEAR), use of lockoutAuth will be enabled if the
// lockout recovery time is 0. Otherwise, lockoutAuth will not be enabled until
// the TPM has been continuously powered for the lockoutRecovery time.
//
// This function requires that NV be available and not rate limiting.
BOOL
DAStartup(
    STARTUP_TYPE     type           // IN: startup type
);

//*** DARegisterFailure()
// This function is called when a authorization failure occurs on an entity
// that is subject to dictionary-attack protection. When a DA failure is
// triggered, register the failure by resetting the relevant self-healing
// timer to the current time.
void
DARegisterFailure(
    TPM_HANDLE       handle         // IN: handle for failure
);

//*** DASelfHeal()
// This function is called to check if sufficient time has passed to allow
// decrement of failedTries or to re-enable use of lockoutAuth.
//
// This function should be called when the time interval is updated.
void
DASelfHeal(
    void
);

#endif  // _DA_FP_H_
