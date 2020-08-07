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
// This header file is used to 'splice' the OpenSSL library into the TPM code.
//
// The support required of a library are a hash module, a block cipher module and
// portions of a big number library.

#ifndef SYM_LIB_DEFINED
#define SYM_LIB_DEFINED

#define SYM_LIB_OSSL

#include <openssl/aes.h>
#include <openssl/des.h>
#include <openssl/bn.h>
#include <openssl/ossl_typ.h>

//***************************************************************
//** Links to the OpenSSL AES code
//***************************************************************

#if ALG_SM4
#error "SM4 is not available"
#endif

#if ALG_CAMELLIA
#error "Camellia is not available"
#endif

// Define the order of parameters to the library functions that do block encryption
// and decryption.
typedef void(*TpmCryptSetSymKeyCall_t)(
    const BYTE  *in,
    BYTE        *out,
    void *keySchedule
    );

// The Crypt functions that call the block encryption function use the parameters 
// in the order:
//  1) keySchedule
//  2) in buffer
//  3) out buffer
// Since open SSL uses the order in encryptoCall_t above, need to swizzle the
// values to the order required by the library.
#define SWIZZLE(keySchedule, in, out)                                   \
    (const BYTE *)(in), (BYTE *)(out), (void *)(keySchedule)

// Macros to set up the encryption/decryption key schedules
//
// AES:
#define TpmCryptSetEncryptKeyAES(key, keySizeInBits, schedule)            \
    AES_set_encrypt_key((key), (keySizeInBits), (tpmKeyScheduleAES *)(schedule))
#define TpmCryptSetDecryptKeyAES(key, keySizeInBits, schedule)            \
    AES_set_decrypt_key((key), (keySizeInBits), (tpmKeyScheduleAES *)(schedule))

// TDES:
#define TpmCryptSetEncryptKeyTDES(key, keySizeInBits, schedule)            \
    TDES_set_encrypt_key((key), (keySizeInBits), (tpmKeyScheduleTDES *)(schedule))
#define TpmCryptSetDecryptKeyTDES(key, keySizeInBits, schedule)            \
    TDES_set_encrypt_key((key), (keySizeInBits), (tpmKeyScheduleTDES *)(schedule))

// Macros to alias encryption calls to specific algorithms. This should be used
// sparingly. Currently, only used by CryptRand.c
// 
// When using these calls, to call the AES block encryption code, the caller 
// should use:
//      TpmCryptEncryptAES(SWIZZLE(keySchedule, in, out));
#define TpmCryptEncryptAES          AES_encrypt
#define TpmCryptDecryptAES          AES_decrypt
#define tpmKeyScheduleAES           AES_KEY


#define TpmCryptEncryptTDES         TDES_encrypt
#define TpmCryptDecryptTDES         TDES_decrypt 
#define tpmKeyScheduleTDES          DES_key_schedule

typedef union tpmCryptKeySchedule_t tpmCryptKeySchedule_t;

#if ALG_TDES
#include "TpmToOsslDesSupport_fp.h"
#endif

// This definition would change if there were something to report
#define SymLibSimulationEnd()

#endif // SYM_LIB_DEFINED
