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
//
// The functions in this file are used for initialization of the interface to the
// LibTomCrypt and MpsLib libraries. This is not used if only the LTC hash and
// symmetric functions are used.

//** Defines and Includes

#include "Tpm.h"

#if defined(HASH_LIB_LTC) || defined(MATH_LIB_LTC) || defined(SYM_LIB_LTC)

// This state is used because there is no way to pass the random number state
// to LibTomCrypt. I do not think that this is currently an issue because...
// Heck, just put in an assert and see what happens.
static void             *s_randState;

//*** LtcRand()
// This is a stub function that is called from the LibTomCrypt or libmpa code
// to get a random number. In turn, this will call the random RandGenerate
// function that was passed in LibraryInit(). This function will pass the pointer
// to the current rand state along with the random byte request.
uint32_t     LtcRand(
    void            *buf,
    size_t           blen
    )
{
    pAssert(1);
    DRBG_Generate(s_randState, buf, (uint16_t)blen);
    return 0;
}

//*** SupportLibInit()
// This does any initialization required by the support library.
LIB_EXPORT int
SupportLibInit(
    void
    )
{
    mpa_set_random_generator(LtcRand);
    s_randState = NULL;
    external_mem_pool = NULL;
    return 1;
}

//*** LtcPoolInit()
// Function to initialize a pool. ****
LIB_EXPORT mpa_scratch_mem
LtcPoolInit(
    mpa_word_t      *poolAddress,
    int              vars,
    int              bits
    )
{
    mpa_scratch_mem     pool = (mpa_scratch_mem)poolAddress;
    mpa_init_scratch_mem(pool, vars, bits);
    init_mpa_tomcrypt(pool);
    return pool;
}

#endif // HASH_LIB_LTC || MATH_LIB_LTC || SYM_LIB_LTC
