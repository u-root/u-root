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
// This file contains the structure definitions for the self-test. It also contains
// macros for use when the self-test is implemented.
#ifndef                 _SELF_TEST_H_
#define                 _SELF_TEST_H_

//** Defines

// Was typing this a lot
#define SELF_TEST_FAILURE   FAIL(FATAL_ERROR_SELF_TEST)

// Use the definition of key sizes to set algorithm values for key size. 
#define AES_ENTRIES (AES_128 + AES_192 + AES_256)
#define SM4_ENTRIES (SM4_128)
#define CAMELLIA_ENTRIES (CAMELLIA_128 + CAMELLIA_192 + CAMELLIA_256)
#define TDES_ENTRIES (TDES_128 + TDES_192)

#define NUM_SYMS    (AES_ENTRIES + SM4_ENTRIES + CAMELLIA_ENTRIES + TDES_ENTRIES)

typedef UINT32      SYM_INDEX;

// These two defines deal with the fact that the TPM_ALG_ID table does not delimit
// the symmetric mode values with a TPM_SYM_MODE_FIRST and TPM_SYM_MODE_LAST
#define TPM_SYM_MODE_FIRST       ALG_CTR_VALUE
#define TPM_SYM_MODE_LAST        ALG_ECB_VALUE

#define NUM_SYM_MODES   (TPM_SYM_MODE_LAST - TPM_SYM_MODE_FIRST + 1)

// Define a type to hold a bit vector for the modes.
#if NUM_SYM_MODES <= 0
#error  "No symmetric modes implemented"
#elif NUM_SYM_MODES <= 8
typedef BYTE    SYM_MODES;
#elif NUM_SYM_MODES <= 16
typedef UINT16  SYM_MODES;
#elif NUM_SYM_MODES <= 32
typedef UINT32  SYM_MODES;
#else
#error "Too many symmetric modes"
#endif

typedef struct SYMMETRIC_TEST_VECTOR {
    const TPM_ALG_ID     alg;                   // the algorithm
    const UINT16         keyBits;               // bits in the key
    const BYTE          *key;                   // The test key
    const UINT32         ivSize;                // block size of the algorithm
    const UINT32         dataInOutSize;         // size  to encrypt/decrypt
    const BYTE          *dataIn;                // data to encrypt
    const BYTE          *dataOut[NUM_SYM_MODES];// data to decrypt
} SYMMETRIC_TEST_VECTOR;

#if ALG_SHA512
#       define  DEFAULT_TEST_HASH               ALG_SHA512_VALUE
#       define  DEFAULT_TEST_DIGEST_SIZE        SHA512_DIGEST_SIZE
#       define  DEFAULT_TEST_HASH_BLOCK_SIZE    SHA512_BLOCK_SIZE
#elif ALG_SHA384
#       define  DEFAULT_TEST_HASH               ALG_SHA384_VALUE
#       define  DEFAULT_TEST_DIGEST_SIZE        SHA384_DIGEST_SIZE
#       define  DEFAULT_TEST_HASH_BLOCK_SIZE    SHA384_BLOCK_SIZE
#elif ALG_SHA256
#       define  DEFAULT_TEST_HASH               ALG_SHA256_VALUE
#       define  DEFAULT_TEST_DIGEST_SIZE        SHA256_DIGEST_SIZE
#       define  DEFAULT_TEST_HASH_BLOCK_SIZE    SHA256_BLOCK_SIZE
#elif ALG_SHA1
#       define  DEFAULT_TEST_HASH               ALG_SHA1_VALUE
#       define  DEFAULT_TEST_DIGEST_SIZE        SHA1_DIGEST_SIZE
#       define  DEFAULT_TEST_HASH_BLOCK_SIZE    SHA1_BLOCK_SIZE
#endif


#endif      // _SELF_TEST_H_