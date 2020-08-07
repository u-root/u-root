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
//
// This file contains the implementation of the symmetric block cipher modes
// allowed for a TPM. These functions only use the single block encryption functions
// of the selected symmetric cryptographic library.

//** Includes, Defines, and Typedefs
#ifndef CRYPT_SYM_H
#define CRYPT_SYM_H

typedef union tpmCryptKeySchedule_t {
#if ALG_AES
    tpmKeyScheduleAES           AES;
#endif
#if ALG_SM4
    tpmKeyScheduleSM4           SM4;
#endif
#if ALG_CAMELLIA
    tpmKeyScheduleCAMELLIA      CAMELLIA;
#endif
 
#if ALG_TDES
    tpmKeyScheduleTDES          TDES[3];
#endif
#if SYMMETRIC_ALIGNMENT == 8
    uint64_t            alignment;
#else
    uint32_t            alignment;
#endif
} tpmCryptKeySchedule_t;


// Each block cipher within a library is expected to conform to the same calling 
// conventions with three parameters ('keySchedule', 'in', and 'out') in the same
// order. That means that all algorithms would use the same order of the same 
// parameters. The code is written assuming the ('keySchedule', 'in', and 'out') 
// order. However, if the library uses a different order, the order can be changed
// with a SWIZZLE macro that puts the parameters in the correct order.
// Note that all algorithms have to use the same order and number of parameters
// because the code to build the calling list is common for each call to encrypt 
// or decrypt with the algorithm chosen by setting a function pointer to select
// the algorithm that is used.

#   define ENCRYPT(keySchedule, in, out)                                            \
      encrypt(SWIZZLE(keySchedule, in, out))

#   define DECRYPT(keySchedule, in, out)                                            \
      decrypt(SWIZZLE(keySchedule, in, out))


// Note that the macros rely on 'encrypt' as local values in the
// functions that use these macros. Those parameters are set by the macro that 
// set the key schedule to be used for the call.


#define ENCRYPT_CASE(ALG)                                                   \
    case TPM_ALG_##ALG:                                                     \
        TpmCryptSetEncryptKey##ALG(key, keySizeInBits, &keySchedule.ALG);   \
        encrypt = (TpmCryptSetSymKeyCall_t)TpmCryptEncrypt##ALG;            \
        break;
#define DECRYPT_CASE(ALG)                                                   \
    case TPM_ALG_##ALG:                                                     \
        TpmCryptSetDecryptKey##ALG(key, keySizeInBits, &keySchedule.ALG);   \
        decrypt = (TpmCryptSetSymKeyCall_t)TpmCryptDecrypt##ALG;            \
        break;

#if ALG_AES
#define ENCRYPT_CASE_AES    ENCRYPT_CASE(AES)
#define DECRYPT_CASE_AES    DECRYPT_CASE(AES)
#else
#define ENCRYPT_CASE_AES
#define DECRYPT_CASE_AES
#endif
#if ALG_SM4
#define ENCRYPT_CASE_SM4    ENCRYPT_CASE(SM4)
#define DECRYPT_CASE_SM4    DECRYPT_CASE(SM4)
#else
#define ENCRYPT_CASE_SM4
#define DECRYPT_CASE_SM4
#endif
#if ALG_CAMELLIA
#define ENCRYPT_CASE_CAMELLIA    ENCRYPT_CASE(CAMELLIA)
#define DECRYPT_CASE_CAMELLIA    DECRYPT_CASE(CAMELLIA)
#else
#define ENCRYPT_CASE_CAMELLIA
#define DECRYPT_CASE_CAMELLIA
#endif
#if ALG_TDES
#define ENCRYPT_CASE_TDES    ENCRYPT_CASE(TDES)
#define DECRYPT_CASE_TDES    DECRYPT_CASE(TDES)
#else
#define ENCRYPT_CASE_TDES
#define DECRYPT_CASE_TDES
#endif

// For each algorithm the case will either be defined or null.
#define     SELECT(direction)                               \
    switch(algorithm)                                       \
    {                                                       \
        direction##_CASE_AES                                \
        direction##_CASE_SM4                                \
        direction##_CASE_CAMELLIA                           \
        direction##_CASE_TDES                               \
        default:                                            \
            FAIL(FATAL_ERROR_INTERNAL);                     \
    }


#endif // CRYPT_SYM_H