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


/* TPM specific preprocessor flags for wolfcrypt */


#ifndef WOLF_CRYPT_USER_SETTINGS_H
#define WOLF_CRYPT_USER_SETTINGS_H

/* Remove the automatic setting of the default I/O functions EmbedSend()
    and EmbedReceive(). */
#define WOLFSSL_USER_IO

/* Avoid naming conflicts */
#define NO_OLD_WC_NAMES

/* Use stack based fast math for all big integer math */
#define USE_FAST_MATH
#define TFM_TIMING_RESISTANT

/* Expose direct encryption functions */
#define WOLFSSL_AES_DIRECT

/* Enable/Disable algorithm support based on TPM implementation header */
#if ALG_SHA256
    #define WOLFSSL_SHA256
#endif
#if ALG_SHA384 || ALG_SHA512
    #define WOLFSSL_SHA384
    #define WOLFSSL_SHA512
#endif
#if ALG_TDES
    #define WOLFSSL_DES_ECB
#endif
#if ALG_RSA
    /* Turn on RSA key generation functionality */
    #define WOLFSSL_KEY_GEN
#endif
#if ALG_ECC || defined(WOLFSSL_LIB)
    #define HAVE_ECC

    /* Expose additional ECC primitives */
    #define WOLFSSL_PUBLIC_ECC_ADD_DBL 
    #define ECC_TIMING_RESISTANT

    /* Enables Shamir calc method */
    #define ECC_SHAMIR

    /* The TPM only needs low level ECC crypto */
    #define NO_ECC_SIGN
    #define NO_ECC_VERIFY
    #define NO_ECC_SECP

    #undef ECC_BN_P256
    #undef ECC_SM2_P256
    #undef ECC_BN_P638
    #define ECC_BN_P256     NO
    #define ECC_SM2_P256    NO
    #define ECC_BN_P638     NO

#endif

/* Disable explicit RSA. The TPM support for RSA is dependent only on TFM */
#define NO_RSA
#define NO_RC4
#define NO_ASN

/* Enable debug wolf library check */
//#define LIBRARY_COMPATIBILITY_CHECK

#define WOLFSSL_

#endif // WOLF_CRYPT_USER_SETTINGS_H
