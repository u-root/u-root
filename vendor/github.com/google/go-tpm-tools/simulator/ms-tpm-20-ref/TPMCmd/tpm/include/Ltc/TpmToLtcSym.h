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
// This header file is used to "splice" the TPM to the LTC symmetric cipher code.

#ifndef SYM_LIB_DEFINED
#define SYM_LIB_DEFINED

#define SYM_LIB_LTC

// Avoid pulling in the MPA math if not doing asymmetric with LTC
#if !(defined MATH_LIB_LTC)
#  define LTC_NO_ASYMMETRIC
#endif

#include "LtcSettings.h"

//***************************************************************
//******** Linking to the TomCrypt AES code *********************
//***************************************************************

#if ALG_SM4
#error "SM4 is not available"
#endif

#if ALG_CAMELLIA
#error "Camellia is not available"
#endif

// Define the order of parameters to the functions that do block encryption and
// decryption.
typedef void(*TpmCryptSetSymKeyCall_t)(
    const void      *in,
    void            *out,
    void            *keySchedule
    );

// Macro to put the parameters in the order required by the library
#define SWIZZLE(keySchedule, in, out)                                               \
    (const void *)(in), (void *)(out), (void *)(keySchedule)

// Macros to set up the encryption/decryption key schedules
//
// AES:
# define TpmCryptSetEncryptKeyAES(key, keySizeInBits, schedule)                     \
    aes_setup((key), BITS_TO_BYTES(keySizeInBits), 0, (symmetric_key *)(schedule))
# define TpmCryptSetDecryptKeyAES(key, keySizeInBits, schedule)                     \
    aes_setup((key), BITS_TO_BYTES(keySizeInBits), 0, (symmetric_key *)(schedule))

// TDES:
# define TpmCryptSetEncryptKeyTDES(key, keySizeInBits, schedule)                    \
    TDES_setup((key), (keySizeInBits), (symmetric_key *)(schedule))
# define TpmCryptSetDecryptKeyTDES(key, keySizeInBits, schedule)                    \
    TDES_setup((key), (keySizeInBits), (symmetric_key *)(schedule))


// Macros to alias encrypt and decrypt function calls to library-specific values
// sparingly. These should be used sparingly. Currently, they are only used by
// CryptRand.c in the AES version of the DRBG.
#define TpmCryptEncryptAES      aes_ecb_encrypt
#define TpmCryptDecryptAES      aes_ecb_decrypt
#define tpmKeyScheduleAES       struct rijndael_key
//
#define TpmCryptEncryptTDES     des3_ecb_encrypt
#define TpmCryptDecryptTDES     des3_ecb_decrypt
#define tpmKeyScheduleTDES      struct des3_key

typedef union tpmCryptKeySchedule_t tpmCryptKeySchedule_t;

#include "TpmToLtcDesSupport_fp.h"

// This is used to trigger printing of simulation statistics

#define SymLibSimulationEnd()

#endif // SYM_LIB_DEFINED
