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

#ifndef    _TIME_FP_H_
#define    _TIME_FP_H_

//*** TimePowerOn()
// This function initialize time info at _TPM_Init().
//
// This function is called at _TPM_Init() so that the TPM time can start counting
// as soon as the TPM comes out of reset and doesn't have to wait until
// TPM2_Startup() in order to begin the new time epoch. This could be significant
// for systems that could get powered up but not run any TPM commands for some
// period of time.
//
void
TimePowerOn(
    void
);

//*** TimeStartup()
// This function updates the resetCount and restartCount components of
// TPMS_CLOCK_INFO structure at TPM2_Startup().
//
// This function will deal with the deferred creation of a new epoch.
// TimeUpdateToCurrent() will not start a new epoch even if one is due when
// TPM_Startup() has not been run. This is because the state of NV is not known
// until startup completes. When Startup is done, then it will create the epoch
// nonce to complete the initializations by calling this function.
BOOL
TimeStartup(
    STARTUP_TYPE     type           // IN: start up type
);

//*** TimeClockUpdate()
// This function updates go.clock. If 'newTime' requires an update of NV, then
// NV is checked for availability. If it is not available or is rate limiting, then
// go.clock is not updated and the function returns an error. If 'newTime' would
// not cause an NV write, then go.clock is updated. If an NV write occurs, then
// go.safe is SET.
void
TimeClockUpdate(
    UINT64           newTime    // IN: New time value in mS.
);

//*** TimeUpdate()
// This function is used to update the time and clock values. If the TPM
// has run TPM2_Startup(), this function is called at the start of each command.
// If the TPM has not run TPM2_Startup(), this is called from TPM2_Startup() to
// get the clock values initialized. It is not called on command entry because, in
// this implementation, the go structure is not read from NV until TPM2_Startup().
// The reason for this is that the initialization code (_TPM_Init()) may run before
// NV is accessible.
void
TimeUpdate(
    void
);

//*** TimeUpdateToCurrent()
// This function updates the 'Time' and 'Clock' in the global
// TPMS_TIME_INFO structure.
//
// In this implementation, 'Time' and 'Clock' are updated at the beginning
// of each command and the values are unchanged for the duration of the
// command.
//
// Because 'Clock' updates may require a write to NV memory, 'Time' and 'Clock'
// are not allowed to advance if NV is not available. When clock is not advancing,
// any function that uses 'Clock' will fail and return TPM_RC_NV_UNAVAILABLE or
// TPM_RC_NV_RATE.
//
// This implementation does not do rate limiting. If the implementation does do
// rate limiting, then the 'Clock' update should not be inhibited even when doing
// rate limiting.
void
TimeUpdateToCurrent(
    void
);

//*** TimeSetAdjustRate()
// This function is used to perform rate adjustment on 'Time' and 'Clock'.
void
TimeSetAdjustRate(
    TPM_CLOCK_ADJUST     adjust         // IN: adjust constant
);

//*** TimeGetMarshaled()
// This function is used to access TPMS_TIME_INFO in canonical form.
// The function collects the time information and marshals it into 'dataBuffer'
// and returns the marshaled size
UINT16
TimeGetMarshaled(
    TIME_INFO       *dataBuffer     // OUT: result buffer
);

//*** TimeFillInfo
// This function gathers information to fill in a TPMS_CLOCK_INFO structure.
void
TimeFillInfo(
    TPMS_CLOCK_INFO     *clockInfo
);

#endif  // _TIME_FP_H_
