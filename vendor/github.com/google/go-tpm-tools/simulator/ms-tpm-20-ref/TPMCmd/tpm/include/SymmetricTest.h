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

// This file contains the structures and data definitions for the symmetric tests.
// This file references the header file that contains the actual test vectors. This
// organization was chosen so that the program that is used to generate the test
// vector values does not have to also re-generate this data.
#ifndef     SELF_TEST_DATA
#error  "This file may only be included in AlgorithmTests.c"
#endif

#ifndef     _SYMMETRIC_TEST_H
#define     _SYMMETRIC_TEST_H
#include    "SymmetricTestData.h"


//** Symmetric Test Structures

const SYMMETRIC_TEST_VECTOR   c_symTestValues[NUM_SYMS + 1] = {
#if ALG_AES && AES_128
    {ALG_AES_VALUE, 128, key_AES128, 16, sizeof(dataIn_AES128), dataIn_AES128,
    {dataOut_AES128_CTR, dataOut_AES128_OFB, dataOut_AES128_CBC, 
     dataOut_AES128_CFB, dataOut_AES128_ECB}},
#endif
#if ALG_AES && AES_192
    {ALG_AES_VALUE, 192, key_AES192, 16, sizeof(dataIn_AES192), dataIn_AES192,
    {dataOut_AES192_CTR, dataOut_AES192_OFB, dataOut_AES192_CBC, 
     dataOut_AES192_CFB, dataOut_AES192_ECB}},
#endif
#if ALG_AES && AES_256
    {ALG_AES_VALUE, 256, key_AES256, 16, sizeof(dataIn_AES256), dataIn_AES256,
    {dataOut_AES256_CTR, dataOut_AES256_OFB, dataOut_AES256_CBC,
    dataOut_AES256_CFB, dataOut_AES256_ECB}},
#endif
#if ALG_SM4 && SM4_128
    {ALG_SM4_VALUE, 128, key_SM4128, 16, sizeof(dataIn_SM4128), dataIn_SM4128,
    {dataOut_SM4128_CTR, dataOut_SM4128_OFB, dataOut_SM4128_CBC, 
     dataOut_SM4128_CFB, dataOut_AES128_ECB}},
#endif
    {0}
};

#endif  // _SYMMETRIC_TEST_H
