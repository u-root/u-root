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
// This file contains the structure definitions used for ECC in the LibTomCrypt
// version of the code. These definitions would change, based on the library.
// The ECC-related structures that cross the TPM interface are defined
// in TpmTypes.h
//

#ifndef MATH_LIB_DEFINED
#define MATH_LIB_DEFINED

#define MATH_LIB_WOLF

#if ALG_ECC
#define HAVE_ECC
#endif

#include <wolfssl/wolfcrypt/tfm.h>
#include <wolfssl/wolfcrypt/ecc.h>

#define MP_VAR(name)                      \
    mp_int          _##name;                                   \
    mp_int          *name = MpInitialize(&_##name);

// Allocate a mp_int and initialize with the values in a mp_int* initializer
#define MP_INITIALIZED(name, initializer)                      \
    MP_VAR(name);                                              \
    BnToWolf(name, initializer);

#define POINT_CREATE(name, initializer)                   \
    ecc_point       *name = EcPointInitialized(initializer);

#define POINT_DELETE(name)                                \
    wc_ecc_del_point(name);                               \
    name = NULL;

typedef ECC_CURVE_DATA bnCurve_t;

typedef bnCurve_t  *bigCurve;

#define AccessCurveData(E)  (E)

#define CURVE_INITIALIZED(name, initializer)                        \
    bnCurve_t      *name = (ECC_CURVE_DATA *)GetCurveData(initializer)

#define CURVE_FREE(E)

#include "TpmToWolfSupport_fp.h"

#define WOLF_ENTER()

#define WOLF_LEAVE()

// This definition would change if there were something to report
#define MathLibSimulationEnd()

#endif // MATH_LIB_DEFINED
