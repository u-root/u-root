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
 *  list of conditions and the following disclaimer in the documentation and/or other
 *  materials provided with the distribution.
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
// This file contains the structure definitions used for linking from the TPM
// code to the MPA and LTC math libraries.

#ifndef MATH_LIB_DEFINED
#define MATH_LIB_DEFINED

#define MATH_LIB_LTC

_REDUCE_WARNING_LEVEL_(2)
#include "LtcSettings.h"
#include "mpalib.h"
#include "mpa.h"
#include "tomcrypt_mpa.h"
_NORMAL_WARNING_LEVEL_


#if RADIX_BITS != 32
#error "The mpa library used with LibTomCrypt only works for 32-bit words"
#endif

// These macros handle entering and leaving a scope
// from which an MPA or LibTomCrypt function may be called.
// Many of these functions require a scratch pool from which
// they will allocate scratch variables (rather than using their
// own stack).
extern mpa_scratch_mem external_mem_pool;

#define MPA_ENTER(vars, bits)                                       \
    mpa_word_t           POOL_ [                                    \
                         mpa_scratch_mem_size_in_U32(vars, bits)];  \
    mpa_scratch_mem      pool_save = external_mem_pool;             \
    mpa_scratch_mem      POOL = LtcPoolInit(POOL_, vars, bits)

#define MPA_LEAVE()     init_mpa_tomcrypt(pool_save)

typedef ECC_CURVE_DATA bnCurve_t;

typedef bnCurve_t  *bigCurve;

#define AccessCurveData(E)  (E)

// Include the support functions for the routines that are used by LTC thunk.
#include "TpmToLtcSupport_fp.h"

#define CURVE_INITIALIZED(name, initializer)                        \
    bnCurve_t      *name = (ECC_CURVE_DATA *)GetCurveData(initializer)

#define CURVE_FREE(E)

// This definition would change if there were something to report
#define MathLibSimulationEnd()

#endif // MATH_LIB_DEFINED
