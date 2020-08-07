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
// This header contains the hash structure definitions used in the TPM code
// to define the amount of space to be reserved for the hash state. This allows
// the TPM code to not have to import all of the symbols used by the hash
// computations. This lets the build environment of the TPM code not to have
// include the header files associated with the CryptoEngine code.

#ifndef _CRYPT_HASH_H
#define _CRYPT_HASH_H

//** Hash-related Structures

union SMAC_STATES;

// These definitions add the high-level methods for processing state that may be
// an SMAC
typedef void(* SMAC_DATA_METHOD)(
    union SMAC_STATES       *state,
    UINT32                   size,
    const BYTE              *buffer
    );

typedef UINT16(* SMAC_END_METHOD)(
    union SMAC_STATES       *state,
    UINT32                   size,
    BYTE                    *buffer
    );

typedef struct sequenceMethods {
    SMAC_DATA_METHOD          data;
    SMAC_END_METHOD           end;
} SMAC_METHODS;

#define SMAC_IMPLEMENTED (CC_MAC || CC_MAC_Start)

// These definitions are here because the SMAC state is in the union of hash states.
typedef struct tpmCmacState {
    TPM_ALG_ID              symAlg;
    UINT16                  keySizeBits;
    INT16                   bcount; // current count of bytes accumulated in IV
    TPM2B_IV                iv;     // IV buffer
    TPM2B_SYM_KEY           symKey;
} tpmCmacState_t;

typedef union SMAC_STATES {
#if ALG_CMAC
    tpmCmacState_t          cmac;
#endif
    UINT64                  pad;
} SMAC_STATES;

typedef struct SMAC_STATE {
    SMAC_METHODS            smacMethods;
    SMAC_STATES             state;
} SMAC_STATE;


typedef union
{
#if ALG_SHA1
    tpmHashStateSHA1_t         Sha1;
#endif
#if ALG_SHA256
    tpmHashStateSHA256_t       Sha256;
#endif
#if ALG_SHA384
    tpmHashStateSHA384_t       Sha384;
#endif
#if ALG_SHA512
    tpmHashStateSHA512_t       Sha512;
#endif

// Additions for symmetric block cipher MAC
#if SMAC_IMPLEMENTED
    SMAC_STATE                 smac;
#endif
    // to force structure alignment to be no worse than HASH_ALIGNMENT
#if HASH_ALIGNMENT == 4
    uint32_t             align;
#else
    uint64_t             align;
#endif
} ANY_HASH_STATE;

typedef ANY_HASH_STATE *PANY_HASH_STATE;
typedef const ANY_HASH_STATE    *PCANY_HASH_STATE;

#define ALIGNED_SIZE(x, b) ((((x) + (b) - 1) / (b)) * (b))
// MAX_HASH_STATE_SIZE will change with each implementation. It is assumed that
// a hash state will not be larger than twice the block size plus some
// overhead (in this case, 16 bytes). The overall size needs to be as
// large as any of the hash contexts. The structure needs to start on an
// alignment boundary and be an even multiple of the alignment
#define MAX_HASH_STATE_SIZE ((2 * MAX_HASH_BLOCK_SIZE) + 16)
#define MAX_HASH_STATE_SIZE_ALIGNED                                             \
                    ALIGNED_SIZE(MAX_HASH_STATE_SIZE, HASH_ALIGNMENT)

// This is an aligned byte array that will hold any of the hash contexts.
typedef  ANY_HASH_STATE ALIGNED_HASH_STATE;

// The header associated with the hash library is expected to define the methods
// which include the calling sequence. When not compiling CryptHash.c, the methods
// are not defined so we need placeholder functions for the structures

#ifndef HASH_START_METHOD_DEF
#   define HASH_START_METHOD_DEF    void (HASH_START_METHOD)(void)
#endif
#ifndef HASH_DATA_METHOD_DEF
#   define HASH_DATA_METHOD_DEF     void (HASH_DATA_METHOD)(void)
#endif
#ifndef HASH_END_METHOD_DEF
#   define HASH_END_METHOD_DEF      void (HASH_END_METHOD)(void)
#endif
#ifndef HASH_STATE_COPY_METHOD_DEF
#   define HASH_STATE_COPY_METHOD_DEF     void (HASH_STATE_COPY_METHOD)(void)
#endif
#ifndef  HASH_STATE_EXPORT_METHOD_DEF
#   define  HASH_STATE_EXPORT_METHOD_DEF   void (HASH_STATE_EXPORT_METHOD)(void)
#endif
#ifndef  HASH_STATE_IMPORT_METHOD_DEF
#   define  HASH_STATE_IMPORT_METHOD_DEF   void (HASH_STATE_IMPORT_METHOD)(void)
#endif

// Define the prototypical function call for each of the methods. This defines the
// order in which the parameters are passed to the underlying function.
typedef HASH_START_METHOD_DEF;
typedef HASH_DATA_METHOD_DEF;
typedef HASH_END_METHOD_DEF;
typedef HASH_STATE_COPY_METHOD_DEF;
typedef HASH_STATE_EXPORT_METHOD_DEF;
typedef HASH_STATE_IMPORT_METHOD_DEF;


typedef struct _HASH_METHODS
{
    HASH_START_METHOD           *start;
    HASH_DATA_METHOD            *data;
    HASH_END_METHOD             *end;
    HASH_STATE_COPY_METHOD      *copy;      // Copy a hash block
    HASH_STATE_EXPORT_METHOD    *copyOut;   // Copy a hash block from a hash
                                            // context
    HASH_STATE_IMPORT_METHOD    *copyIn;    // Copy a hash block to a proper hash
                                            // context
} HASH_METHODS, *PHASH_METHODS;

#if ALG_SHA1
    TPM2B_TYPE(SHA1_DIGEST, SHA1_DIGEST_SIZE);
#endif
#if ALG_SHA256
    TPM2B_TYPE(SHA256_DIGEST, SHA256_DIGEST_SIZE);
#endif
#if ALG_SHA384
    TPM2B_TYPE(SHA384_DIGEST, SHA384_DIGEST_SIZE);
#endif
#if ALG_SHA512
    TPM2B_TYPE(SHA512_DIGEST, SHA512_DIGEST_SIZE);
#endif
#if ALG_SM3_256
    TPM2B_TYPE(SM3_256_DIGEST, SM3_256_DIGEST_SIZE);
#endif

// When the TPM implements RSA, the hash-dependent OID pointers are part of the
// HASH_DEF. These macros conditionally add the OID reference to the HASH_DEF and the
// HASH_DEF_TEMPLATE.
#if ALG_RSA
#define PKCS1_HASH_REF   const BYTE  *PKCS1;
#define PKCS1_OID(NAME)  , OID_PKCS1_##NAME
#else
#define PKCS1_HASH_REF
#define PKCS1_OID(NAME)
#endif

// When the TPM implements ECC, the hash-dependent OID pointers are part of the
// HASH_DEF. These macros conditionally add the OID reference to the HASH_DEF and the
// HASH_DEF_TEMPLATE.
#if ALG_ECDSA 
#define ECDSA_HASH_REF    const BYTE  *ECDSA;
#define ECDSA_OID(NAME)  , OID_ECDSA_##NAME
#else
#define ECDSA_HASH_REF
#define ECDSA_OID(NAME)
#endif

typedef const struct HASH_DEF
{
    HASH_METHODS         method;
    uint16_t             blockSize;
    uint16_t             digestSize;
    uint16_t             contextSize;
    uint16_t             hashAlg;
    const BYTE          *OID;
    PKCS1_HASH_REF      // PKCS1 OID
    ECDSA_HASH_REF      // ECDSA OID
} HASH_DEF, *PHASH_DEF;

// Macro to fill in the HASH_DEF for an algorithm. For SHA1, the instance would be:
//  HASH_DEF_TEMPLATE(Sha1, SHA1)
// This handles the difference in capitalization for the various pieces.
#define HASH_DEF_TEMPLATE(HASH, Hash)                                               \
    HASH_DEF    Hash##_Def= {                                                       \
                        {(HASH_START_METHOD *)&tpmHashStart_##HASH,                 \
                         (HASH_DATA_METHOD *)&tpmHashData_##HASH,                   \
                         (HASH_END_METHOD *)&tpmHashEnd_##HASH,                     \
                         (HASH_STATE_COPY_METHOD *)&tpmHashStateCopy_##HASH,        \
                         (HASH_STATE_EXPORT_METHOD *)&tpmHashStateExport_##HASH,    \
                         (HASH_STATE_IMPORT_METHOD *)&tpmHashStateImport_##HASH,    \
                        },                                                          \
                        HASH##_BLOCK_SIZE,     /*block size */                      \
                        HASH##_DIGEST_SIZE,    /*data size */                       \
                        sizeof(tpmHashState##HASH##_t),                             \
                        TPM_ALG_##HASH, OID_##HASH                                  \
                        PKCS1_OID(HASH) ECDSA_OID(HASH)};

// These definitions are for the types that can be in a hash state structure.
// These types are used in the cryptographic utilities. This is a define rather than
// an enum so that the size of this field can be explicit.
typedef BYTE    HASH_STATE_TYPE;
#define HASH_STATE_EMPTY        ((HASH_STATE_TYPE) 0)
#define HASH_STATE_HASH         ((HASH_STATE_TYPE) 1)
#define HASH_STATE_HMAC         ((HASH_STATE_TYPE) 2)
#if CC_MAC || CC_MAC_Start
#define HASH_STATE_SMAC         ((HASH_STATE_TYPE) 3)
#endif


// This is the structure that is used for passing a context into the hashing
// functions. It should be the same size as the function context used within
// the hashing functions. This is checked when the hash function is initialized.
// This version uses a new layout for the contexts and a different definition. The
// state buffer is an array of HASH_UNIT values so that a decent compiler will put
// the structure on a HASH_UNIT boundary. If the structure is not properly aligned,
// the code that manipulates the structure will copy to a properly aligned
// structure before it is used and copy the result back. This just makes things
// slower.
// NOTE: This version of the state had the pointer to the update method in the
// state. This is to allow the SMAC functions to use the same structure without 
// having to replicate the entire HASH_DEF structure.
typedef struct _HASH_STATE
{
    HASH_STATE_TYPE          type;               // type of the context
    TPM_ALG_ID               hashAlg;
    PHASH_DEF                def;
    ANY_HASH_STATE           state;
} HASH_STATE, *PHASH_STATE;
typedef const HASH_STATE *PCHASH_STATE;


//** HMAC State Structures

// An HMAC_STATE structure contains an opaque HMAC stack state. A caller would
// use this structure when performing incremental HMAC operations. This structure
// contains a hash state and an HMAC key and allows slightly better stack
// optimization than adding an HMAC key to each hash state.
typedef struct hmacState
{
    HASH_STATE           hashState;          // the hash state
    TPM2B_HASH_BLOCK     hmacKey;            // the HMAC key
} HMAC_STATE, *PHMAC_STATE;

// This is for the external hash state. This implementation assumes that the size
// of the exported hash state is no larger than the internal hash state.
typedef struct
{
    BYTE                     buffer[sizeof(HASH_STATE)];
} EXPORT_HASH_STATE, *PEXPORT_HASH_STATE;

typedef const EXPORT_HASH_STATE *PCEXPORT_HASH_STATE;

#endif //  _CRYPT_HASH_H
