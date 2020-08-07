/* Microsoft Reference Implementation for TPM 2.0
 *
 *  The copyright in this software is being made available under the BSD
 * License, included below. This software may be subject to other third party
 * and contributor rights, including patent rights, and no such rights are
 * granted under this license.
 *
 *  Copyright (c) Microsoft Corporation
 *
 *  All rights reserved.
 *
 *  BSD License
 *
 *  Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 *  Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS ""AS
 * IS"" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
 * THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
 * OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
 * OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
 * ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
//** Description
//
// This file contains the routines that are used by the simulator to mimic
// a hardware clock on a TPM.
//
// In this implementation, all the time values are measured in millisecond.
// However, the precision of the clock functions may be implementation
// dependent.

#include <time.h>

#include "PlatformData.h"
#include "Platform_fp.h"

unsigned int s_adjustRate;
bool s_timerReset;

clock64_t s_realTimePrevious;
clock64_t s_tpmTime;
clock64_t s_lastSystemTime;
clock64_t s_lastReportedTime;

void _plat__TimerReset() {
  s_lastSystemTime = 0;
  s_tpmTime = 0;
  s_adjustRate = CLOCK_NOMINAL;
  s_timerReset = true;
  return;
}

static uint64_t _plat__RealTime() {
  struct timespec systime;
  clock_gettime(CLOCK_MONOTONIC, &systime);
  return (clock64_t)systime.tv_sec * 1000 + (systime.tv_nsec / 1000000);
}

uint64_t _plat__TimerRead() {
  clock64_t timeDiff;
  clock64_t adjustedTimeDiff;
  clock64_t timeNow;
  clock64_t readjustedTimeDiff;

  // This produces a timeNow that is basically locked to the system clock.
  timeNow = _plat__RealTime();

  // if this hasn't been initialized, initialize it
  if (s_lastSystemTime == 0) {
    s_lastSystemTime = timeNow;
    s_lastReportedTime = 0;
    s_realTimePrevious = 0;
  }
  // The system time can bounce around and that's OK as long as we don't allow
  // time to go backwards. When the time does appear to go backwards, set
  // lastSystemTime to be the new value and then update the reported time.
  if (timeNow < s_lastReportedTime) s_lastSystemTime = timeNow;
  s_lastReportedTime = s_lastReportedTime + timeNow - s_lastSystemTime;
  s_lastSystemTime = timeNow;
  timeNow = s_lastReportedTime;

  // The code above produces a timeNow that is similar to the value returned
  // by Clock(). The difference is that timeNow does not max out, and it is
  // at a ms. rate rather than at a CLOCKS_PER_SEC rate. The code below
  // uses that value and does the rate adjustment on the time value.
  // If there is no difference in time, then skip all the computations
  if (s_realTimePrevious >= timeNow) return s_tpmTime;
  // Compute the amount of time since the last update of the system clock
  timeDiff = timeNow - s_realTimePrevious;

  // Do the time rate adjustment and conversion from CLOCKS_PER_SEC to mSec
  adjustedTimeDiff = (timeDiff * CLOCK_NOMINAL) / ((uint64_t)s_adjustRate);

  // update the TPM time with the adjusted timeDiff
  s_tpmTime += (clock64_t)adjustedTimeDiff;

  // Might have some rounding error that would loose CLOCKS. See what is not
  // being used. As mentioned above, this could result in putting back more than
  // is taken out. Here, we are trying to recreate timeDiff.
  readjustedTimeDiff =
      (adjustedTimeDiff * (uint64_t)s_adjustRate) / CLOCK_NOMINAL;

  // adjusted is now converted back to being the amount we should advance the
  // previous sampled time. It should always be less than or equal to timeDiff.
  // That is, we could not have use more time than we started with.
  s_realTimePrevious = s_realTimePrevious + readjustedTimeDiff;

  return s_tpmTime;
}

bool _plat__TimerWasReset() {
  bool retVal = s_timerReset;
  s_timerReset = false;
  return retVal;
}

void _plat__ClockAdjustRate(int adjust) {
  // We expect the caller should only use a fixed set of constant values to
  // adjust the rate
  switch (adjust) {
    case CLOCK_ADJUST_COARSE:
      s_adjustRate += CLOCK_ADJUST_COARSE;
      break;
    case -CLOCK_ADJUST_COARSE:
      s_adjustRate -= CLOCK_ADJUST_COARSE;
      break;
    case CLOCK_ADJUST_MEDIUM:
      s_adjustRate += CLOCK_ADJUST_MEDIUM;
      break;
    case -CLOCK_ADJUST_MEDIUM:
      s_adjustRate -= CLOCK_ADJUST_MEDIUM;
      break;
    case CLOCK_ADJUST_FINE:
      s_adjustRate += CLOCK_ADJUST_FINE;
      break;
    case -CLOCK_ADJUST_FINE:
      s_adjustRate -= CLOCK_ADJUST_FINE;
      break;
    default:
      // ignore any other values;
      break;
  }

  if (s_adjustRate > (CLOCK_NOMINAL + CLOCK_ADJUST_LIMIT))
    s_adjustRate = CLOCK_NOMINAL + CLOCK_ADJUST_LIMIT;
  if (s_adjustRate < (CLOCK_NOMINAL - CLOCK_ADJUST_LIMIT))
    s_adjustRate = CLOCK_NOMINAL - CLOCK_ADJUST_LIMIT;

  return;
}
