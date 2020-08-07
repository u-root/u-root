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
// This file contains the functions relating to the TPM's time functions including
// the interface to the implementation-specific time functions.
//
//** Includes
#include "Tpm.h"
#include "PlatformData.h"

//** Functions

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
    )
{
    g_time = _plat__TimerRead();
}

//*** TimeNewEpoch()
// This function does the processing to generate a new time epoch nonce and
// set NV for update. This function is only called when NV is known to be available
// and the clock is running. The epoch is updated to persistent data.
static void
TimeNewEpoch(
    void
    )
{
#if CLOCK_STOPS
    CryptRandomGenerate(sizeof(CLOCK_NONCE), (BYTE *)&g_timeEpoch);
#else
    // if the epoch is kept in NV, update it.
    gp.timeEpoch++;
    NV_SYNC_PERSISTENT(timeEpoch);
#endif
    // Clean out any lingering state
    _plat__TimerWasStopped();
}

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
    )
{
    NOT_REFERENCED(type);
    // If the previous cycle is orderly shut down, the value of the safe bit
    // the same as previously saved.  Otherwise, it is not safe.
    if(!NV_IS_ORDERLY)
        go.clockSafe = NO;
    return TRUE;
}

//*** TimeClockUpdate()
// This function updates go.clock. If 'newTime' requires an update of NV, then
// NV is checked for availability. If it is not available or is rate limiting, then
// go.clock is not updated and the function returns an error. If 'newTime' would
// not cause an NV write, then go.clock is updated. If an NV write occurs, then
// go.safe is SET.
void
TimeClockUpdate(
    UINT64           newTime    // IN: New time value in mS.
    )
{
#define CLOCK_UPDATE_MASK  ((1ULL << NV_CLOCK_UPDATE_INTERVAL)- 1)

    // Check to see if the update will cause a need for an nvClock update
    if((newTime | CLOCK_UPDATE_MASK) > (go.clock | CLOCK_UPDATE_MASK))
    {
        pAssert(g_NvStatus == TPM_RC_SUCCESS);

        // Going to update the NV time state so SET the safe flag
        go.clockSafe = YES;

        // update the time
        go.clock = newTime;

        NvWrite(NV_ORDERLY_DATA, sizeof(go), &go);
    }
    else
        // No NV update needed so just update
        go.clock = newTime;

}

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
    )
{
    UINT64          elapsed;
//
    // Make sure that we consume the current _plat__TimerWasStopped() state.
   if(_plat__TimerWasStopped())
    {
        TimeNewEpoch();
    }
    // Get the difference between this call and the last time we updated the tick
    // timer.
    elapsed = _plat__TimerRead() - g_time;
    // Don't read +
    g_time += elapsed;

    // Don't need to check the result because it has to be success because have
    // already checked that NV is available.
    TimeClockUpdate(go.clock + elapsed);

    // Call self healing logic for dictionary attack parameters
    DASelfHeal();
}

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
)
{
    // Can't update time during the dark interval or when rate limiting so don't
    // make any modifications to the internal clock value. Also, defer any clock
    // processing until TPM has run TPM2_Startup()
    if(!NV_IS_AVAILABLE || !TPMIsStarted())
        return;

    TimeUpdate();
}


//*** TimeSetAdjustRate()
// This function is used to perform rate adjustment on 'Time' and 'Clock'.
void
TimeSetAdjustRate(
    TPM_CLOCK_ADJUST     adjust         // IN: adjust constant
    )
{
    switch(adjust)
    {
        case TPM_CLOCK_COARSE_SLOWER:
            _plat__ClockAdjustRate(CLOCK_ADJUST_COARSE);
            break;
        case TPM_CLOCK_COARSE_FASTER:
            _plat__ClockAdjustRate(-CLOCK_ADJUST_COARSE);
            break;
        case TPM_CLOCK_MEDIUM_SLOWER:
            _plat__ClockAdjustRate(CLOCK_ADJUST_MEDIUM);
            break;
        case TPM_CLOCK_MEDIUM_FASTER:
            _plat__ClockAdjustRate(-CLOCK_ADJUST_MEDIUM);
            break;
        case TPM_CLOCK_FINE_SLOWER:
            _plat__ClockAdjustRate(CLOCK_ADJUST_FINE);
            break;
        case TPM_CLOCK_FINE_FASTER:
            _plat__ClockAdjustRate(-CLOCK_ADJUST_FINE);
            break;
        case TPM_CLOCK_NO_CHANGE:
            break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }

    return;
}

//*** TimeGetMarshaled()
// This function is used to access TPMS_TIME_INFO in canonical form.
// The function collects the time information and marshals it into 'dataBuffer'
// and returns the marshaled size
UINT16
TimeGetMarshaled(
    TIME_INFO       *dataBuffer     // OUT: result buffer
    )
{
    TPMS_TIME_INFO      timeInfo;

    // Fill TPMS_TIME_INFO structure
    timeInfo.time = g_time;
    TimeFillInfo(&timeInfo.clockInfo);

    // Marshal TPMS_TIME_INFO to canonical form
    return TPMS_TIME_INFO_Marshal(&timeInfo, (BYTE **)&dataBuffer, NULL);
}

//*** TimeFillInfo
// This function gathers information to fill in a TPMS_CLOCK_INFO structure.
void
TimeFillInfo(
    TPMS_CLOCK_INFO     *clockInfo
    )
{
    clockInfo->clock = go.clock;
    clockInfo->resetCount = gp.resetCount;
    clockInfo->restartCount = gr.restartCount;

    // If NV is not available, clock stopped advancing and the value reported is
    // not "safe".
    if(NV_IS_AVAILABLE)
        clockInfo->safe = go.clockSafe;
    else
        clockInfo->safe = NO;

    return;
}