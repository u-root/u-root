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
// This file is a collection of miscellaneous macros.

#ifndef GP_MACROS_H
#define GP_MACROS_H

#ifndef NULL
#define NULL 0
#endif

#include "swap.h"
#include "VendorString.h"


//** For Self-test
// These macros are used in CryptUtil to invoke the incremental self test.
#if SELF_TEST
#   define     TEST(alg) if(TEST_BIT(alg, g_toTest)) CryptTestAlgorithm(alg, NULL)

// Use of TPM_ALG_NULL is reserved for RSAEP/RSADP testing. If someone is wanting
// to test a hash with that value, don't do it.
#   define     TEST_HASH(alg)                                                       \
            if(TEST_BIT(alg, g_toTest)                                              \
                &&  (alg != ALG_NULL_VALUE))                                        \
                CryptTestAlgorithm(alg, NULL)
#else
#   define TEST(alg)
#   define TEST_HASH(alg)
#endif // SELF_TEST

//** For Failures
#if defined _POSIX_ 
#   define FUNCTION_NAME        0
#else
#   define FUNCTION_NAME        __FUNCTION__
#endif

#if !FAIL_TRACE
#   define FAIL(errorCode) (TpmFail(errorCode))
#   define LOG_FAILURE(errorCode) (TpmLogFailure(errorCode))
#else
#   define FAIL(errorCode)        TpmFail(FUNCTION_NAME, __LINE__, errorCode)
#   define LOG_FAILURE(errorCode) TpmLogFailure(FUNCTION_NAME, __LINE__, errorCode)
#endif

// If implementation is using longjmp, then the call to TpmFail() does not return
// and the compiler will complain about unreachable code that comes after. To allow
// for not having longjmp, TpmFail() will return and the subsequent code will be
// executed. This macro accounts for the difference.
#ifndef NO_LONGJMP
#   define FAIL_RETURN(returnCode)
#   define TPM_FAIL_RETURN     NORETURN void
#else
#   define FAIL_RETURN(returnCode) return (returnCode)
#   define TPM_FAIL_RETURN     void
#endif

// This macro tests that a condition is TRUE and puts the TPM into failure mode
// if it is not. If longjmp is being used, then the FAIL(FATAL_ERROR_) macro makes 
// a call from which there is no return. Otherwise, it returns and the function 
// will exit with the appropriate return code.
#define REQUIRE(condition, errorCode, returnCode)                                   \
    {                                                                               \
        if(!!(condition))                                                           \
        {                                                                           \
            FAIL(FATAL_ERROR_errorCode);                                            \
            FAIL_RETURN(returnCode);                                                \
        }                                                                           \
    }

#define PARAMETER_CHECK(condition, returnCode)                                      \
    REQUIRE((condition), PARAMETER, returnCode)

#if (defined EMPTY_ASSERT) && (EMPTY_ASSERT != NO)
#   define pAssert(a)  ((void)0)
#else
#   define pAssert(a) {if(!(a)) FAIL(FATAL_ERROR_PARAMETER);}
#endif

//** Derived from Vendor-specific values
// Values derived from vendor specific settings in TpmProfile.h
#define PCR_SELECT_MIN          ((PLATFORM_PCR+7)/8)
#define PCR_SELECT_MAX          ((IMPLEMENTATION_PCR+7)/8)
#define MAX_ORDERLY_COUNT       ((1 << ORDERLY_BITS) - 1)

//** Compile-time Checks
// In some cases, the relationship between two values may be dependent
// on things that change based on various selections like the chosen cryptographic
// libraries. It is possible that these selections will result in incompatible
// settings. These are often detectable by the compiler but it isn't always 
// possible to do the check in the preprocessor code. For example, when the
// check requires use of "sizeof" then the preprocessor can't do the comparison.
// For these cases, we include a special macro that, depending on the compiler
// will generate a warning to indicate if the check always passes or always fails
// because it involves fixed constants. To run these checks, define COMPILER_CHECKS
// in TpmBuildSwitches.h
#if COMPILER_CHECKS
#   define  cAssert     pAssert
#else
#   define cAssert(value)
#endif

// This is used commonly in the "Crypt" code as a way to keep listings from 
// getting too long. This is not to save paper but to allow one to see more
// useful stuff on the screen at any given time.
#define     ERROR_RETURN(returnCode)                                                \
    {                                                                               \
         retVal = returnCode;                                                       \
         goto Exit;                                                                 \
    }

#ifndef MAX
#  define MAX(a, b) ((a) > (b) ? (a) : (b))
#endif
#ifndef MIN
#  define MIN(a, b) ((a) < (b) ? (a) : (b))
#endif
#ifndef IsOdd
#  define IsOdd(a)        (((a) & 1) != 0)
#endif

#ifndef BITS_TO_BYTES
#  define BITS_TO_BYTES(bits) (((bits) + 7) >> 3)
#endif

// These are defined for use when the size of the vector being checked is known
// at compile time.
#define TEST_BIT(bit, vector)   TestBit((bit), (BYTE *)&(vector), sizeof(vector))
#define SET_BIT(bit, vector)    SetBit((bit), (BYTE *)&(vector), sizeof(vector))
#define CLEAR_BIT(bit, vector) ClearBit((bit), (BYTE *)&(vector), sizeof(vector))


// The following definitions are used if they have not already been defined. The
// defaults for these settings are compatible with ISO/IEC 9899:2011 (E)
#ifndef LIB_EXPORT
#   define LIB_EXPORT
#   define LIB_IMPORT
#endif
#ifndef NORETURN
#   define NORETURN _Noreturn
#endif
#ifndef NOT_REFERENCED
#   define NOT_REFERENCED(x = x)   ((void) (x))
#endif

#define STD_RESPONSE_HEADER (sizeof(TPM_ST) + sizeof(UINT32) + sizeof(TPM_RC))

#define JOIN(x, y) x##y
#define JOIN3(x, y, z) x##y##z
#define CONCAT(x, y) JOIN(x, y)
#define CONCAT3(x, y, z) JOIN3(x,y,z)

// If CONTEXT_INTEGRITY_HASH_ALG is defined, then the vendor is using the old style
// table. Otherwise, pick the "strongest" implemented hash algorithm as the context
// hash.
#ifndef CONTEXT_HASH_ALGORITHM
#   if defined ALG_SHA512 && ALG_SHA512 == YES
#       define CONTEXT_HASH_ALGORITHM    SHA512
#   elif defined ALG_SHA384 && ALG_SHA384 == YES
#       define CONTEXT_HASH_ALGORITHM    SHA384
#   elif defined ALG_SHA256 && ALG_SHA256 == YES
#       define CONTEXT_HASH_ALGORITHM    SHA256
#   elif defined ALG_SM3_256 && ALG_SM3_256 == YES
#       define CONTEXT_HASH_ALGORITHM    SM3_256
#   elif defined ALG_SHA1 && ALG_SHA1 == YES
#       define CONTEXT_HASH_ALGORITHM  SHA1  
#   endif
#   define CONTEXT_INTEGRITY_HASH_ALG  CONCAT(TPM_ALG_, CONTEXT_HASH_ALGORITHM)
#endif

#ifndef CONTEXT_INTEGRITY_HASH_SIZE
#define CONTEXT_INTEGRITY_HASH_SIZE CONCAT(CONTEXT_HASH_ALGORITHM, _DIGEST_SIZE)
#endif
#if     ALG_RSA
#define     RSA_SECURITY_STRENGTH (MAX_RSA_KEY_BITS >= 15360 ? 256 :                \
                                  (MAX_RSA_KEY_BITS >=  7680 ? 192 :                \
                                  (MAX_RSA_KEY_BITS >=  3072 ? 128 :                \
                                  (MAX_RSA_KEY_BITS >=  2048 ? 112 :                \
                                  (MAX_RSA_KEY_BITS >=  1024 ?  80 :  0)))))
#else
#define     RSA_SECURITY_STRENGTH   0
#endif // ALG_RSA

#if     ALG_ECC
#define     ECC_SECURITY_STRENGTH (MAX_ECC_KEY_BITS >= 521 ? 256 :                  \
                                  (MAX_ECC_KEY_BITS >= 384 ? 192 :                  \
                                  (MAX_ECC_KEY_BITS >= 256 ? 128 : 0)))
#else
#define     ECC_SECURITY_STRENGTH   0
#endif // ALG_ECC

#define     MAX_ASYM_SECURITY_STRENGTH                                              \
                        MAX(RSA_SECURITY_STRENGTH, ECC_SECURITY_STRENGTH)

#define     MAX_HASH_SECURITY_STRENGTH  ((CONTEXT_INTEGRITY_HASH_SIZE * 8) / 2)

// Unless some algorithm is broken...
#define     MAX_SYM_SECURITY_STRENGTH   MAX_SYM_KEY_BITS

#define MAX_SECURITY_STRENGTH_BITS                                                  \
                        MAX(MAX_ASYM_SECURITY_STRENGTH,                             \
                        MAX(MAX_SYM_SECURITY_STRENGTH,                              \
                            MAX_HASH_SECURITY_STRENGTH))

// This is the size that was used before the 1.38 errata requiring that P1.14.4 be
// followed
#define PROOF_SIZE      CONTEXT_INTEGRITY_HASH_SIZE

// As required by P1.14.4
#define COMPLIANT_PROOF_SIZE                                                        \
            (MAX(CONTEXT_INTEGRITY_HASH_SIZE, (2 * MAX_SYM_KEY_BYTES)))
      
// As required by P1.14.3.1
#define COMPLIANT_PRIMARY_SEED_SIZE                                                 \
    BITS_TO_BYTES(MAX_SECURITY_STRENGTH_BITS * 2)

// This is the pre-errata version
#ifndef PRIMARY_SEED_SIZE
#   define PRIMARY_SEED_SIZE    PROOF_SIZE
#endif

#if USE_SPEC_COMPLIANT_PROOFS
#   undef PROOF_SIZE
#   define PROOF_SIZE           COMPLIANT_PROOF_SIZE
#   undef PRIMARY_SEED_SIZE
#   define PRIMARY_SEED_SIZE    COMPLIANT_PRIMARY_SEED_SIZE
#endif  // USE_SPEC_COMPLIANT_PROOFS

#if !SKIP_PROOF_ERRORS 
#   if PROOF_SIZE < COMPLIANT_PROOF_SIZE
#       error "PROOF_SIZE is not compliant with TPM specification"
#   endif
#   if PRIMARY_SEED_SIZE < COMPLIANT_PRIMARY_SEED_SIZE
#       error Non-compliant PRIMARY_SEED_SIZE
#   endif
#endif // !SKIP_PROOF_ERRORS

// If CONTEXT_ENCRYPT_ALG is defined, then the vendor is using the old style table
#if defined CONTEXT_ENCRYPT_ALG 
#   undef CONTEXT_ENCRYPT_ALGORITHM
#   if CONTEXT_ENCRYPT_ALG == ALG_AES_VALUE
#       define CONTEXT_ENCRYPT_ALGORITHM  AES
#   elif CONTEXT_ENCRYPT_ALG == ALG_SM4_VALUE
#       define CONTEXT_ENCRYPT_ALGORITHM  SM4
#   elif CONTEXT_ENCRYPT_ALG == ALG_CAMELLIA_VALUE
#       define CONTEXT_ENCRYPT_ALGORITHM  CAMELLIA
#   elif CONTEXT_ENCRYPT_ALG == ALG_TDES_VALUE
#   error Are you kidding? 
#   else
#       error Unknown value for CONTEXT_ENCRYPT_ALG
#   endif // CONTEXT_ENCRYPT_ALG == ALG_AES_VALUE
#else
#   define CONTEXT_ENCRYPT_ALG                                                      \
            CONCAT3(ALG_, CONTEXT_ENCRYPT_ALGORITHM, _VALUE)
#endif  // CONTEXT_ENCRYPT_ALG 
#define CONTEXT_ENCRYPT_KEY_BITS                                                    \
                CONCAT(CONTEXT_ENCRYPT_ALGORITHM, _MAX_KEY_SIZE_BITS)
#define CONTEXT_ENCRYPT_KEY_BYTES       ((CONTEXT_ENCRYPT_KEY_BITS+7)/8)

// This is updated to follow the requirement of P2 that the label not be larger
// than 32 bytes.
#ifndef LABEL_MAX_BUFFER
#define LABEL_MAX_BUFFER MIN(32, MAX(MAX_ECC_KEY_BYTES, MAX_DIGEST_SIZE))
#endif

// This bit is used to indicate that an authorization ticket expires on TPM Reset
// and TPM Restart. It is added to the timeout value returned by TPM2_PoliySigned()
// and TPM2_PolicySecret() and used by TPM2_PolicyTicket(). The timeout value is 
// relative to Time (g_time). Time is reset whenever the TPM loses power and cannot
// be moved forward by the user (as can Clock). 'g_time' is a 64-bit value expressing 
// time in ms. Stealing the MSb for a flag means that the TPM needs to be reset
// at least once every 292,471,208 years rather than once every 584,942,417 years.
#define EXPIRATION_BIT ((UINT64)1 << 63)

// Check for consistency of the bit ordering of bit fields
#if BIG_ENDIAN_TPM && MOST_SIGNIFICANT_BIT_0 && USE_BIT_FIELD_STRUCTURES
#   error "Settings not consistent"
#endif

// These macros are used to handle the variation in handling of bit fields. If 
#if USE_BIT_FIELD_STRUCTURES // The default, old version, with bit fields
#   define IS_ATTRIBUTE(a, type, b)    ((a.b) != 0)
#   define SET_ATTRIBUTE(a, type, b)       (a.b = SET)
#   define CLEAR_ATTRIBUTE(a, type, b)     (a.b = CLEAR)
#   define GET_ATTRIBUTE(a, type, b)        (a.b)
#   define TPMA_ZERO_INITIALIZER()          {0}
#else
#   define IS_ATTRIBUTE(a, type, b)         ((a & type##_##b) != 0)
#   define SET_ATTRIBUTE(a, type, b)        (a |= type##_##b)
#   define CLEAR_ATTRIBUTE(a, type, b)      (a &= ~type##_##b)
#   define GET_ATTRIBUTE(a, type, b)        \
        (type)((a & type##_##b) >> type##_##b##_SHIFT)
#   define TPMA_ZERO_INITIALIZER()         (0)
#endif

#define VERIFY(_X) if(!(_X)) goto Error 

#endif // GP_MACROS_H