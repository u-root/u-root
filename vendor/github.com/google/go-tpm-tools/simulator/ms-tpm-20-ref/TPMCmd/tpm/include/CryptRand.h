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
// This file contains constant definition shared by CryptUtil and the parts
// of the Crypto Engine.
//

#ifndef _CRYPT_RAND_H
#define _CRYPT_RAND_H


//** DRBG Structures and Defines

// Values and structures for the random number generator. These values are defined
// in this header file so that the size of the RNG state can be known to TPM.lib.
// This allows the allocation of some space in NV memory for the state to
// be stored on an orderly shutdown.

// The DRBG based on a symmetric block cipher is defined by three values,
// 1) the key size
// 2) the block size (the IV size)
// 3) the symmetric algorithm

#define DRBG_KEY_SIZE_BITS      AES_MAX_KEY_SIZE_BITS
#define DRBG_IV_SIZE_BITS       (AES_MAX_BLOCK_SIZE * 8)
#define DRBG_ALGORITHM          TPM_ALG_AES


typedef tpmKeyScheduleAES     DRBG_KEY_SCHEDULE;
#define DRBG_ENCRYPT_SETUP(key, keySizeInBits, schedule)        \
            TpmCryptSetEncryptKeyAES(key, keySizeInBits, schedule)
#define DRBG_ENCRYPT(keySchedule, in, out)                      \
            TpmCryptEncryptAES(SWIZZLE(keySchedule, in, out))

#if     ((DRBG_KEY_SIZE_BITS % RADIX_BITS) != 0) \
    || ((DRBG_IV_SIZE_BITS % RADIX_BITS) != 0)
#error "Key size and IV for DRBG must be even multiples of the radix"
#endif
#if (DRBG_KEY_SIZE_BITS % DRBG_IV_SIZE_BITS) != 0
#error "Key size for DRBG must be even multiple of the cypher block size"
#endif

// Derived values
#define DRBG_MAX_REQUESTS_PER_RESEED (1 << 48)
#define DRBG_MAX_REQEST_SIZE (1 << 32)

#define pDRBG_KEY(seed)    ((DRBG_KEY *)&(((BYTE *)(seed))[0]))
#define pDRBG_IV(seed)     ((DRBG_IV *)&(((BYTE *)(seed))[DRBG_KEY_SIZE_BYTES]))

#define DRBG_KEY_SIZE_WORDS     (BITS_TO_CRYPT_WORDS(DRBG_KEY_SIZE_BITS))
#define DRBG_KEY_SIZE_BYTES     (DRBG_KEY_SIZE_WORDS * RADIX_BYTES)

#define DRBG_IV_SIZE_WORDS      (BITS_TO_CRYPT_WORDS(DRBG_IV_SIZE_BITS))
#define DRBG_IV_SIZE_BYTES      (DRBG_IV_SIZE_WORDS * RADIX_BYTES)

#define DRBG_SEED_SIZE_WORDS    (DRBG_KEY_SIZE_WORDS + DRBG_IV_SIZE_WORDS)
#define DRBG_SEED_SIZE_BYTES    (DRBG_KEY_SIZE_BYTES + DRBG_IV_SIZE_BYTES)


typedef union
{
    BYTE            bytes[DRBG_KEY_SIZE_BYTES];
    crypt_uword_t   words[DRBG_KEY_SIZE_WORDS];
} DRBG_KEY;

typedef union
{
    BYTE            bytes[DRBG_IV_SIZE_BYTES];
    crypt_uword_t   words[DRBG_IV_SIZE_WORDS];
} DRBG_IV;

typedef union
{
    BYTE            bytes[DRBG_SEED_SIZE_BYTES];
    crypt_uword_t   words[DRBG_SEED_SIZE_WORDS];
} DRBG_SEED;

#define CTR_DRBG_MAX_REQUESTS_PER_RESEED        ((UINT64)1 << 20)
#define CTR_DRBG_MAX_BYTES_PER_REQUEST          (1 << 16)

#   define CTR_DRBG_MIN_ENTROPY_INPUT_LENGTH    DRBG_SEED_SIZE_BYTES
#   define CTR_DRBG_MAX_ENTROPY_INPUT_LENGTH    DRBG_SEED_SIZE_BYTES
#   define CTR_DRBG_MAX_ADDITIONAL_INPUT_LENGTH DRBG_SEED_SIZE_BYTES

#define     TESTING         (1 << 0)
#define     ENTROPY         (1 << 1)
#define     TESTED          (1 << 2)

#define     IsTestStateSet(BIT)    ((g_cryptoSelfTestState.rng & BIT) != 0)
#define     SetTestStateBit(BIT)   (g_cryptoSelfTestState.rng |= BIT)
#define     ClearTestStateBit(BIT) (g_cryptoSelfTestState.rng &= ~BIT)

#define     IsSelfTest()    IsTestStateSet(TESTING)
#define     SetSelfTest()   SetTestStateBit(TESTING)
#define     ClearSelfTest() ClearTestStateBit(TESTING)

#define     IsEntropyBad()      IsTestStateSet(ENTROPY)
#define     SetEntropyBad()     SetTestStateBit(ENTROPY)
#define     ClearEntropyBad()   ClearTestStateBit(ENTROPY)

#define     IsDrbgTested()      IsTestStateSet(TESTED)
#define     SetDrbgTested()     SetTestStateBit(TESTED)
#define     ClearDrbgTested()   ClearTestStateBit(TESTED)

typedef struct
{
    UINT64      reseedCounter;
    UINT32      magic;
    DRBG_SEED   seed; // contains the key and IV for the counter mode DRBG
    UINT32      lastValue[4];   // used when the TPM does continuous self-test
                                // for FIPS compliance of DRBG
} DRBG_STATE, *pDRBG_STATE;
#define DRBG_MAGIC   ((UINT32) 0x47425244) // "DRBG" backwards so that it displays

typedef struct
{
    UINT64               counter;
    UINT32               magic;
    UINT32               limit;
    TPM2B               *seed;
    const TPM2B         *label;
    TPM2B               *context;
    TPM_ALG_ID           hash;
    TPM_ALG_ID           kdf;
    UINT16               digestSize;
    TPM2B_DIGEST         residual;
} KDF_STATE, *pKDR_STATE;
#define KDF_MAGIC    ((UINT32) 0x4048444a) // "KDF " backwards

// Make sure that any other structures added to this union start with a 64-bit
// counter and a 32-bit magic number
typedef union
{
    DRBG_STATE      drbg;
    KDF_STATE       kdf;
} RAND_STATE;

// This is the state used when the library uses a random number generator.
// A special function is installed for the library to call. That function
// picks up the state from this location and uses it for the generation
// of the random number.
extern RAND_STATE           *s_random;

// When instrumenting RSA key sieve
#if  RSA_INSTRUMENT
#define PRIME_INDEX(x)  ((x) == 512 ? 0 : (x) == 1024 ? 1 : 2)
#   define INSTRUMENT_SET(a, b) ((a) = (b))
#   define INSTRUMENT_ADD(a, b) (a) = (a) + (b)
#   define INSTRUMENT_INC(a)    (a) = (a) + 1

extern UINT32  PrimeIndex;
extern UINT32  failedAtIteration[10];
extern UINT32  PrimeCounts[3];
extern UINT32  MillerRabinTrials[3];
extern UINT32  totalFieldsSieved[3];
extern UINT32  bitsInFieldAfterSieve[3];
extern UINT32  emptyFieldsSieved[3];
extern UINT32  noPrimeFields[3];
extern UINT32  primesChecked[3];
extern UINT16  lastSievePrime;
#else
#   define INSTRUMENT_SET(a, b)
#   define INSTRUMENT_ADD(a, b)
#   define INSTRUMENT_INC(a)
#endif

#endif // _CRYPT_RAND_H
