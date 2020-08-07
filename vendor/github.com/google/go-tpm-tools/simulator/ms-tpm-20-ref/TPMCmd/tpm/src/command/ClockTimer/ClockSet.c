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
#include "Tpm.h"
#include "ClockSet_fp.h"

#if CC_ClockSet  // Conditional expansion of this file

// Read the current TPMS_TIMER_INFO structure settings
//  Return Type: TPM_RC
//      TPM_RC_NV_RATE              NV is unavailable because of rate limit
//      TPM_RC_NV_UNAVAILABLE       NV is inaccessible
//      TPM_RC_VALUE                invalid new clock

TPM_RC
TPM2_ClockSet(
    ClockSet_In     *in             // IN: input parameter list
    )
{
// Input Validation
    // new time can not be bigger than 0xFFFF000000000000 or smaller than
    // current clock
    if(in->newTime > 0xFFFF000000000000ULL
       || in->newTime < go.clock)
        return TPM_RCS_VALUE + RC_ClockSet_newTime;

// Internal Data Update
    // Can't modify the clock if NV is not available.
    RETURN_IF_NV_IS_NOT_AVAILABLE;

    TimeClockUpdate(in->newTime);
    return TPM_RC_SUCCESS;
}

#endif // CC_ClockSet