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
// This file contains constant definitions used for self-test.

#ifndef _CRYPT_TEST_H
#define _CRYPT_TEST_H

// This is the definition of a bit array with one bit per algorithm. 
// NOTE: Since bit numbering starts at zero, when ALG_LAST_VALUE is a multiple of 8, 
// ALGORITHM_VECTOR will need to have byte for the single bit in the last byte. So, 
// for example, when ALG_LAST_VECTOR is 8, ALGORITHM_VECTOR will need 2 bytes.
#define ALGORITHM_VECTOR_BYTES  ((ALG_LAST_VALUE + 8) / 8) 
typedef BYTE    ALGORITHM_VECTOR[ALGORITHM_VECTOR_BYTES];

#ifdef  TEST_SELF_TEST
LIB_EXPORT    extern  ALGORITHM_VECTOR    LibToTest;
#endif

// This structure is used to contain self-test tracking information for the 
// cryptographic modules. Each of the major modules is given a 32-bit value in 
// which it may maintain its own self test information. The convention for this 
// state is that when all of the bits in this structure are 0, all functions need 
// to be tested.
typedef struct
{
    UINT32      rng;
    UINT32      hash;
    UINT32      sym;
#if ALG_RSA
    UINT32      rsa;
#endif
#if ALG_ECC
    UINT32      ecc;
#endif
} CRYPTO_SELF_TEST_STATE;


#endif // _CRYPT_TEST_H
