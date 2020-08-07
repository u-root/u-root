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
//** Description
//
// This file contains implementation of cryptographic functions for hashing.
//
//** Includes, Defines, and Types

#define _CRYPT_HASH_C_
#include "Tpm.h"
#include "CryptHash_fp.h"
#include "CryptHash.h"
#include "OIDs.h"

#define HASH_TABLE_SIZE     (HASH_COUNT + 1)


#if     ALG_SHA1
HASH_DEF_TEMPLATE(SHA1, Sha1);
#endif
#if     ALG_SHA256
HASH_DEF_TEMPLATE(SHA256, Sha256);
#endif
#if     ALG_SHA384
HASH_DEF_TEMPLATE(SHA384, Sha384);
#endif
#if     ALG_SHA512
HASH_DEF_TEMPLATE(SHA512, Sha512);
#endif
#if ALG_SM3_256
HASH_DEF_TEMPLATE(SM3_256, Sm3_256);
#endif
HASH_DEF NULL_Def = {{0}};

PHASH_DEF       HashDefArray[] = {
#if ALG_SHA1
    &Sha1_Def,
#endif
#if ALG_SHA256
    &Sha256_Def,
#endif
#if ALG_SHA384
    &Sha384_Def,
#endif
#if ALG_SHA512
    &Sha512_Def,
#endif
#if ALG_SM3_256
    &Sm3_256_Def,
#endif
    &NULL_Def
};


//** Obligatory Initialization Functions

//*** CryptHashInit()
// This function is called by _TPM_Init do perform the initialization operations for
// the library.
BOOL
CryptHashInit(
    void
    )
{
    LibHashInit();
    return TRUE;
}

//*** CryptHashStartup()
// This function is called by TPM2_Startup(). It checks that the size of the 
// HashDefArray is consistent with the HASH_COUNT.
BOOL
CryptHashStartup(
    void
    )
{
    int         i = sizeof(HashDefArray) / sizeof(PHASH_DEF) - 1;
    return (i == HASH_COUNT);
}

//** Hash Information Access Functions
//*** Introduction
// These functions provide access to the hash algorithm description information.

//*** CryptGetHashDef()
// This function accesses the hash descriptor associated with a hash a
// algorithm. The function returns a pointer to a 'null' descriptor if hashAlg is 
// TPM_ALG_NULL or not a defined algorithm. 
PHASH_DEF
CryptGetHashDef(
    TPM_ALG_ID       hashAlg
    )
{
    size_t           i;
#define HASHES  (sizeof(HashDefArray) / sizeof(PHASH_DEF))
    for(i = 0; i < HASHES; i++)
    {
        PHASH_DEF p = HashDefArray[i];
        if(p->hashAlg == hashAlg)
            return p;
    }
    return &NULL_Def;
}

//*** CryptHashIsValidAlg()
// This function tests to see if an algorithm ID is a valid hash algorithm. If
// flag is true, then TPM_ALG_NULL is a valid hash.
//  Return Type: BOOL
//      TRUE(1)         hashAlg is a valid, implemented hash on this TPM
//      FALSE(0)        hashAlg is not valid for this TPM
BOOL
CryptHashIsValidAlg(
    TPM_ALG_ID       hashAlg,           // IN: the algorithm to check
    BOOL             flag               // IN: TRUE if TPM_ALG_NULL is to be treated
                                        //     as a valid hash
    )
{
    if(hashAlg == TPM_ALG_NULL)
        return flag;
    return CryptGetHashDef(hashAlg) != &NULL_Def;
}

//*** CryptHashGetAlgByIndex()
// This function is used to iterate through the hashes. TPM_ALG_NULL
// is returned for all indexes that are not valid hashes.
// If the TPM implements 3 hashes, then an 'index' value of 0 will
// return the first implemented hash and an 'index' of 2 will return the
// last. All other index values will return TPM_ALG_NULL.
//
//  Return Type: TPM_ALG_ID
// TPM_ALG_xxx         a hash algorithm
// TPM_ALG_NULL        this can be used as a stop value
LIB_EXPORT TPM_ALG_ID
CryptHashGetAlgByIndex(
    UINT32           index          // IN: the index
    )
{
    TPM_ALG_ID       hashAlg;
    if(index >= HASH_COUNT)
        hashAlg = TPM_ALG_NULL;
    else
        hashAlg = HashDefArray[index]->hashAlg;
    return hashAlg;
}

//*** CryptHashGetDigestSize()
// Returns the size of the digest produced by the hash. If 'hashAlg' is not a hash
// algorithm, the TPM will FAIL.
//  Return Type: UINT16
//   0       TPM_ALG_NULL
//   > 0     the digest size
//
LIB_EXPORT UINT16
CryptHashGetDigestSize(
    TPM_ALG_ID       hashAlg        // IN: hash algorithm to look up
    )
{
    return CryptGetHashDef(hashAlg)->digestSize;
}

//*** CryptHashGetBlockSize()
// Returns the size of the block used by the hash. If 'hashAlg' is not a hash
// algorithm, the TPM will FAIL.
//  Return Type: UINT16
//   0       TPM_ALG_NULL
//   > 0     the digest size
//
LIB_EXPORT UINT16
CryptHashGetBlockSize(
    TPM_ALG_ID       hashAlg        // IN: hash algorithm to look up
    )
{
    return CryptGetHashDef(hashAlg)->blockSize;
}

//*** CryptHashGetOid()
// This function returns a pointer to DER=encoded OID for a hash algorithm. All OIDs
// are full OID values including the Tag (0x06) and length byte. 
LIB_EXPORT const BYTE *
CryptHashGetOid(
    TPM_ALG_ID      hashAlg
)
{
    return CryptGetHashDef(hashAlg)->OID;
}

//***  CryptHashGetContextAlg()
// This function returns the hash algorithm associated with a hash context.
TPM_ALG_ID
CryptHashGetContextAlg(
    PHASH_STATE      state          // IN: the context to check
    )
{
    return state->hashAlg;
}

//** State Import and Export

//*** CryptHashCopyState
// This function is used to clone a HASH_STATE.
LIB_EXPORT void
CryptHashCopyState(
    HASH_STATE          *out,           // OUT: destination of the state
    const HASH_STATE    *in             // IN: source of the state
    )
{
    pAssert(out->type == in->type);
    out->hashAlg = in->hashAlg;
    out->def = in->def;
    if(in->hashAlg != TPM_ALG_NULL)
    {
         HASH_STATE_COPY(out, in);
    }
    if(in->type == HASH_STATE_HMAC)
    {
        const HMAC_STATE    *hIn = (HMAC_STATE *)in;
        HMAC_STATE          *hOut = (HMAC_STATE *)out;
        hOut->hmacKey = hIn->hmacKey;
    }
    return;
}

//*** CryptHashExportState()
// This function is used to export a hash or HMAC hash state. This function
// would be called when preparing to context save a sequence object.
void
CryptHashExportState(
    PCHASH_STATE         internalFmt,   // IN: the hash state formatted for use by
                                        //     library
    PEXPORT_HASH_STATE   externalFmt    // OUT: the exported hash state
    )
{
    BYTE                    *outBuf = (BYTE *)externalFmt;
//
    cAssert(sizeof(HASH_STATE) <= sizeof(EXPORT_HASH_STATE));
    // the following #define is used to move data from an aligned internal data
    // structure to a byte buffer (external format data.
#define CopyToOffset(value)                                                     \
        memcpy(&outBuf[offsetof(HASH_STATE,value)], &internalFmt->value,        \
                sizeof(internalFmt->value))
    // Copy the hashAlg
    CopyToOffset(hashAlg);
    CopyToOffset(type);
#ifdef HASH_STATE_SMAC
    if(internalFmt->type == HASH_STATE_SMAC)
    {
        memcpy(outBuf, internalFmt, sizeof(HASH_STATE));
        return;

    }
#endif
    if(internalFmt->type == HASH_STATE_HMAC)
    {
        HMAC_STATE              *from = (HMAC_STATE *)internalFmt;
        memcpy(&outBuf[offsetof(HMAC_STATE, hmacKey)], &from->hmacKey,
               sizeof(from->hmacKey));
    }
    if(internalFmt->hashAlg != TPM_ALG_NULL)
         HASH_STATE_EXPORT(externalFmt, internalFmt);
}

//*** CryptHashImportState()
// This function is used to import the hash state. This function
// would be called to import a hash state when the context of a sequence object
// was being loaded.
void
CryptHashImportState(
    PHASH_STATE          internalFmt,   // OUT: the hash state formatted for use by
                                        //     the library
    PCEXPORT_HASH_STATE  externalFmt    // IN: the exported hash state
    )
{
    BYTE                    *inBuf = (BYTE *)externalFmt;
//
#define CopyFromOffset(value)                                                   \
    memcpy(&internalFmt->value, &inBuf[offsetof(HASH_STATE,value)],             \
           sizeof(internalFmt->value))

    // Copy the hashAlg of the byte-aligned input structure to the structure-aligned
    // internal structure.
    CopyFromOffset(hashAlg);
    CopyFromOffset(type);
    if(internalFmt->hashAlg != TPM_ALG_NULL)
    {
#ifdef HASH_STATE_SMAC
        if(internalFmt->type == HASH_STATE_SMAC)
        {
            memcpy(internalFmt, inBuf, sizeof(HASH_STATE));
            return;
        }
#endif
        internalFmt->def = CryptGetHashDef(internalFmt->hashAlg);
        HASH_STATE_IMPORT(internalFmt, inBuf);
        if(internalFmt->type == HASH_STATE_HMAC)
        {
            HMAC_STATE              *to = (HMAC_STATE *)internalFmt;
            memcpy(&to->hmacKey, &inBuf[offsetof(HMAC_STATE, hmacKey)],
                   sizeof(to->hmacKey));
        }
    }
}

//** State Modification Functions

//***HashEnd()
// Local function to complete a hash that uses the hashDef instead of an algorithm
// ID. This function is used to complete the hash and only return a partial digest.
// The return value is the size of the data copied.
static UINT16
HashEnd(
    PHASH_STATE      hashState,     // IN: the hash state
    UINT32           dOutSize,      // IN: the size of receive buffer
    PBYTE            dOut           // OUT: the receive buffer
    )
{
    BYTE                temp[MAX_DIGEST_SIZE];
    if((hashState->hashAlg == TPM_ALG_NULL) 
       || (hashState->type != HASH_STATE_HASH))
        dOutSize = 0;
    if(dOutSize > 0)
    {
        hashState->def = CryptGetHashDef(hashState->hashAlg);
        // Set the final size
        dOutSize = MIN(dOutSize, hashState->def->digestSize);
        // Complete into the temp buffer and then copy
        HASH_END(hashState, temp);
        // Don't want any other functions calling the HASH_END method
        // directly.
#undef HASH_END
        memcpy(dOut, &temp, dOutSize);
    }
    hashState->type = HASH_STATE_EMPTY;
    return (UINT16)dOutSize;
}

//*** CryptHashStart()
// Functions starts a hash stack
// Start a hash stack and returns the digest size. As a side effect, the
// value of 'stateSize' in hashState is updated to indicate the number of bytes
// of state that were saved. This function calls GetHashServer() and that function
// will put the TPM into failure mode if the hash algorithm is not supported.
//
// This function does not use the sequence parameter. If it is necessary to import
// or export context, this will start the sequence in a local state
// and export the state to the input buffer. Will need to add a flag to the state
// structure to indicate that it needs to be imported before it can be used.
// (BLEH).
//  Return Type: UINT16
//  0           hash is TPM_ALG_NULL
// >0           digest size
LIB_EXPORT UINT16
CryptHashStart(
    PHASH_STATE      hashState,     // OUT: the running hash state
    TPM_ALG_ID       hashAlg        // IN: hash algorithm
    )
{
    UINT16               retVal;

    TEST(hashAlg);

    hashState->hashAlg = hashAlg;
    if(hashAlg == TPM_ALG_NULL)
    {
        retVal = 0;
    }
    else
    {
        hashState->def = CryptGetHashDef(hashAlg);
        HASH_START(hashState);
        retVal = hashState->def->digestSize;
    }
#undef HASH_START
    hashState->type = HASH_STATE_HASH;
    return retVal;
}

//*** CryptDigestUpdate()
// Add data to a hash or HMAC, SMAC stack.
//
void
CryptDigestUpdate(
    PHASH_STATE      hashState,     // IN: the hash context information
    UINT32           dataSize,      // IN: the size of data to be added
    const BYTE      *data           // IN: data to be hashed
    )
{
    if(hashState->hashAlg != TPM_ALG_NULL)
    {
        if((hashState->type == HASH_STATE_HASH)
           || (hashState->type == HASH_STATE_HMAC))
            HASH_DATA(hashState, dataSize, (BYTE *)data);
#if SMAC_IMPLEMENTED
        else if(hashState->type == HASH_STATE_SMAC)
            (hashState->state.smac.smacMethods.data)(&hashState->state.smac.state, 
                                                     dataSize, data);
#endif // SMAC_IMPLEMENTED
        else
            FAIL(FATAL_ERROR_INTERNAL);
    }
    return;
}

//*** CryptHashEnd()
// Complete a hash or HMAC computation. This function will place the smaller of
// 'digestSize' or the size of the digest in 'dOut'. The number of bytes in the
// placed in the buffer is returned. If there is a failure, the returned value
// is <= 0.
//  Return Type: UINT16
//       0      no data returned
//      > 0     the number of bytes in the digest or dOutSize, whichever is smaller
LIB_EXPORT UINT16
CryptHashEnd(
    PHASH_STATE      hashState,     // IN: the state of hash stack
    UINT32           dOutSize,      // IN: size of digest buffer
    BYTE            *dOut           // OUT: hash digest
    )
{
    pAssert(hashState->type == HASH_STATE_HASH);
    return HashEnd(hashState, dOutSize, dOut);
}

//*** CryptHashBlock()
// Start a hash, hash a single block, update 'digest' and return the size of
// the results.
//
// The 'digestSize' parameter can be smaller than the digest. If so, only the more
// significant bytes are returned.
//  Return Type: UINT16
//  >= 0        number of bytes placed in 'dOut'
LIB_EXPORT UINT16
CryptHashBlock(
    TPM_ALG_ID       hashAlg,       // IN: The hash algorithm
    UINT32           dataSize,      // IN: size of buffer to hash
    const BYTE      *data,          // IN: the buffer to hash
    UINT32           dOutSize,      // IN: size of the digest buffer
    BYTE            *dOut           // OUT: digest buffer
    )
{
    HASH_STATE          state;
    CryptHashStart(&state, hashAlg);
    CryptDigestUpdate(&state, dataSize, data);
    return HashEnd(&state, dOutSize, dOut);
}

//*** CryptDigestUpdate2B()
// This function updates a digest (hash or HMAC) with a TPM2B.
//
// This function can be used for both HMAC and hash functions so the
// 'digestState' is void so that either state type can be passed.
LIB_EXPORT void
CryptDigestUpdate2B(
    PHASH_STATE      state,         // IN: the digest state
    const TPM2B     *bIn            // IN: 2B containing the data
    )
{
    // Only compute the digest if a pointer to the 2B is provided.
    // In CryptDigestUpdate(), if size is zero or buffer is NULL, then no change
    // to the digest occurs. This function should not provide a buffer if bIn is
    // not provided.
    pAssert(bIn != NULL);
    CryptDigestUpdate(state, bIn->size, bIn->buffer);
    return;
}

//*** CryptHashEnd2B()
// This function is the same as CryptCompleteHash() but the digest is
// placed in a TPM2B. This is the most common use and this is provided
// for specification clarity. 'digest.size' should be set to indicate the number of
// bytes to place in the buffer
//  Return Type: UINT16
//      >=0     the number of bytes placed in 'digest.buffer'
LIB_EXPORT UINT16
CryptHashEnd2B(
    PHASH_STATE      state,         // IN: the hash state
    P2B              digest         // IN: the size of the buffer Out: requested
                                    //     number of bytes
    )
{
    return CryptHashEnd(state, digest->size, digest->buffer);
}

//*** CryptDigestUpdateInt()
// This function is used to include an integer value to a hash stack. The function
// marshals the integer into its canonical form before calling CryptDigestUpdate().
LIB_EXPORT void
CryptDigestUpdateInt(
    void            *state,         // IN: the state of hash stack
    UINT32           intSize,       // IN: the size of 'intValue' in bytes
    UINT64           intValue       // IN: integer value to be hashed
    )
{
#if LITTLE_ENDIAN_TPM
    intValue = REVERSE_ENDIAN_64(intValue);
#endif
    CryptDigestUpdate(state, intSize, &((BYTE *)&intValue)[8 - intSize]);
}

//** HMAC Functions

//*** CryptHmacStart()
// This function is used to start an HMAC using a temp
// hash context. The function does the initialization
// of the hash with the HMAC key XOR iPad and updates the
// HMAC key XOR oPad.
//
// The function returns the number of bytes in a digest produced by 'hashAlg'.
//  Return Type: UINT16
//  >= 0        number of bytes in digest produced by 'hashAlg' (may be zero)
//
LIB_EXPORT UINT16
CryptHmacStart(
    PHMAC_STATE      state,         // IN/OUT: the state buffer
    TPM_ALG_ID       hashAlg,       // IN: the algorithm to use
    UINT16           keySize,       // IN: the size of the HMAC key
    const BYTE      *key            // IN: the HMAC key
    )
{
    PHASH_DEF                hashDef;
    BYTE *                   pb;
    UINT32                   i;
//
    hashDef = CryptGetHashDef(hashAlg);
    if(hashDef->digestSize != 0)
    {
    // If the HMAC key is larger than the hash block size, it has to be reduced
    // to fit. The reduction is a digest of the hashKey.
        if(keySize > hashDef->blockSize)
        {
            // if the key is too big, reduce it to a digest of itself
            state->hmacKey.t.size = CryptHashBlock(hashAlg, keySize, key,
                                                   hashDef->digestSize, 
                                                   state->hmacKey.t.buffer);
        }
        else
        {
            memcpy(state->hmacKey.t.buffer, key, keySize);
            state->hmacKey.t.size = keySize;
        }
        // XOR the key with iPad (0x36)
        pb = state->hmacKey.t.buffer;
        for(i = state->hmacKey.t.size; i > 0; i--)
            *pb++ ^= 0x36;

        // if the keySize is smaller than a block, fill the rest with 0x36
        for(i = hashDef->blockSize - state->hmacKey.t.size; i > 0; i--)
            *pb++ = 0x36;

        // Increase the oPadSize to a full block
        state->hmacKey.t.size = hashDef->blockSize;

        // Start a new hash with the HMAC key
        // This will go in the caller's state structure and may be a sequence or not
        CryptHashStart((PHASH_STATE)state, hashAlg);
        CryptDigestUpdate((PHASH_STATE)state, state->hmacKey.t.size,
                          state->hmacKey.t.buffer);
        // XOR the key block with 0x5c ^ 0x36
        for(pb = state->hmacKey.t.buffer, i = hashDef->blockSize; i > 0; i--)
            *pb++ ^= (0x5c ^ 0x36);
    }
    // Set the hash algorithm
    state->hashState.hashAlg = hashAlg;
    // Set the hash state type
    state->hashState.type = HASH_STATE_HMAC;

    return hashDef->digestSize;
}

//*** CryptHmacEnd()
// This function is called to complete an HMAC. It will finish the current
// digest, and start a new digest. It will then add the oPadKey and the
// completed digest and return the results in dOut. It will not return more
// than dOutSize bytes.
//  Return Type: UINT16
//  >= 0        number of bytes in 'dOut' (may be zero)
LIB_EXPORT UINT16
CryptHmacEnd(
    PHMAC_STATE      state,         // IN: the hash state buffer
    UINT32           dOutSize,      // IN: size of digest buffer
    BYTE            *dOut           // OUT: hash digest
    )
{
    BYTE                 temp[MAX_DIGEST_SIZE];
    PHASH_STATE          hState = (PHASH_STATE)&state->hashState;

#if SMAC_IMPLEMENTED
    if(hState->type == HASH_STATE_SMAC)
        return (state->hashState.state.smac.smacMethods.end)
                    (&state->hashState.state.smac.state,
                     dOutSize,
                     dOut);
#endif
    pAssert(hState->type == HASH_STATE_HMAC);
    hState->def = CryptGetHashDef(hState->hashAlg);
    // Change the state type for completion processing
    hState->type = HASH_STATE_HASH;
    if(hState->hashAlg == TPM_ALG_NULL)
        dOutSize = 0;
    else
    {
        
    // Complete the current hash
        HashEnd(hState, hState->def->digestSize, temp);
        // Do another hash starting with the oPad
        CryptHashStart(hState, hState->hashAlg);
        CryptDigestUpdate(hState, state->hmacKey.t.size, state->hmacKey.t.buffer);
        CryptDigestUpdate(hState, hState->def->digestSize, temp);
    }
    return HashEnd(hState, dOutSize, dOut);
}

//*** CryptHmacStart2B()
// This function starts an HMAC and returns the size of the digest
// that will be produced.
//
// This function is provided to support the most common use of starting an HMAC
// with a TPM2B key.
//
// The caller must provide a block of memory in which the hash sequence state
// is kept.  The caller should not alter the contents of this buffer until the
// hash sequence is completed or abandoned.
//
//  Return Type: UINT16
//      > 0     the digest size of the algorithm
//      = 0     the hashAlg was TPM_ALG_NULL
LIB_EXPORT UINT16
CryptHmacStart2B(
    PHMAC_STATE      hmacState,     // OUT: the state of HMAC stack. It will be used
                                    //     in HMAC update and completion
    TPMI_ALG_HASH    hashAlg,       // IN: hash algorithm
    P2B              key            // IN: HMAC key
    )
{
    return CryptHmacStart(hmacState, hashAlg, key->size, key->buffer);
}

//*** CryptHmacEnd2B()
//   This function is the same as CryptHmacEnd() but the HMAC result
//   is returned in a TPM2B which is the most common use.
//  Return Type: UINT16
//      >=0     the number of bytes placed in 'digest'
LIB_EXPORT UINT16
CryptHmacEnd2B(
    PHMAC_STATE      hmacState,     // IN: the state of HMAC stack
    P2B              digest         // OUT: HMAC
    )
{
    return CryptHmacEnd(hmacState, digest->size, digest->buffer);
}

//** Mask and Key Generation Functions
//*** CryptMGF1()
// This function performs MGF1 using the selected hash. MGF1 is
// T(n) = T(n-1) || H(seed || counter).
// This function returns the length of the mask produced which
// could be zero if the digest algorithm is not supported
//  Return Type: UINT16
//      0       hash algorithm was TPM_ALG_NULL
//    > 0       should be the same as 'mSize'
LIB_EXPORT UINT16
CryptMGF1(
    UINT32           mSize,         // IN: length of the mask to be produced
    BYTE            *mask,          // OUT: buffer to receive the mask
    TPM_ALG_ID       hashAlg,       // IN: hash to use
    UINT32           seedSize,      // IN: size of the seed
    BYTE            *seed           // IN: seed size
    )
{
    HASH_STATE           hashState;
    PHASH_DEF            hDef = CryptGetHashDef(hashAlg);
    UINT32               remaining;
    UINT32               counter = 0;
    BYTE                 swappedCounter[4];

    // If there is no digest to compute return
    if((hashAlg == TPM_ALG_NULL) || (mSize == 0))
        return 0;

    for(remaining = mSize; ; remaining -= hDef->digestSize)
    {
        // Because the system may be either Endian...
        UINT32_TO_BYTE_ARRAY(counter, swappedCounter);

        // Start the hash and include the seed and counter
        CryptHashStart(&hashState, hashAlg);
        CryptDigestUpdate(&hashState, seedSize, seed);
        CryptDigestUpdate(&hashState, 4, swappedCounter);

        // Handling the completion depends on how much space remains in the mask
        // buffer. If it can hold the entire digest, put it there. If not
        // put the digest in a temp buffer and only copy the amount that
        // will fit into the mask buffer.
        HashEnd(&hashState, remaining, mask);
        if(remaining <= hDef->digestSize)
            break;
        mask = &mask[hDef->digestSize];
        counter++;
    }
    return (UINT16)mSize;
}

//*** CryptKDFa()
// This function performs the key generation according to Part 1 of the
// TPM specification.
//
// This function returns the number of bytes generated which may be zero.
//
// The 'key' and 'keyStream' pointers are not allowed to be NULL. The other
// pointer values may be NULL. The value of 'sizeInBits' must be no larger
// than (2^18)-1 = 256K bits (32385 bytes).
//
// The 'once' parameter is set to allow incremental generation of a large
// value. If this flag is TRUE, 'sizeInBits' will be used in the HMAC computation
// but only one iteration of the KDF is performed. This would be used for
// XOR obfuscation so that the mask value can be generated in digest-sized
// chunks rather than having to be generated all at once in an arbitrarily
// large buffer and then XORed into the result. If 'once' is TRUE, then
// 'sizeInBits' must be a multiple of 8.
//
// Any error in the processing of this command is considered fatal.
//  Return Type: UINT16
//     0            hash algorithm is not supported or is TPM_ALG_NULL
//    > 0           the number of bytes in the 'keyStream' buffer
LIB_EXPORT UINT16
CryptKDFa(
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm used in HMAC
    const TPM2B     *key,           // IN: HMAC key
    const TPM2B     *label,         // IN: a label for the KDF
    const TPM2B     *contextU,      // IN: context U
    const TPM2B     *contextV,      // IN: context V
    UINT32           sizeInBits,    // IN: size of generated key in bits
    BYTE            *keyStream,     // OUT: key buffer
    UINT32          *counterInOut,  // IN/OUT: caller may provide the iteration
                                    //     counter for incremental operations to
                                    //     avoid large intermediate buffers.
    UINT16           blocks         // IN: If non-zero, this is the maximum number
                                    //     of blocks to be returned, regardless
                                    //     of sizeInBits
    )
{
    UINT32                   counter = 0;       // counter value
    INT16                    bytes;             // number of bytes to produce
    UINT16                   generated;         // number of bytes generated
    BYTE                    *stream = keyStream;
    HMAC_STATE               hState;
    UINT16                   digestSize = CryptHashGetDigestSize(hashAlg);

    pAssert(key != NULL && keyStream != NULL);
    
    TEST(TPM_ALG_KDF1_SP800_108);

    if(digestSize == 0)
        return 0;

    if(counterInOut != NULL)
        counter = *counterInOut;

    // If the size of the request is larger than the numbers will handle,
    // it is a fatal error.
    pAssert(((sizeInBits + 7) / 8) <= INT16_MAX);

    // The number of bytes to be generated is the smaller of the sizeInBits bytes or
    // the number of requested blocks. The number of blocks is the smaller of the
    // number requested or the number allowed by sizeInBits. A partial block is
    // a full block.
    bytes = (blocks > 0) ? blocks * digestSize : (UINT16)BITS_TO_BYTES(sizeInBits);
    generated = bytes;

    // Generate required bytes
    for(; bytes > 0; bytes -= digestSize)
    {
        counter++;
        // Start HMAC
        if(CryptHmacStart(&hState, hashAlg, key->size, key->buffer) == 0)
            return 0;
        // Adding counter
        CryptDigestUpdateInt(&hState.hashState, 4, counter);

        // Adding label
        if(label != NULL)
            HASH_DATA(&hState.hashState, label->size, (BYTE *)label->buffer);
        // Add a null. SP108 is not very clear about when the 0 is needed but to
        // make this like the previous version that did not add an 0x00 after
        // a null-terminated string, this version will only add a null byte
        // if the label parameter did not end in a null byte, or if no label
        // is present.
        if((label == NULL)
           || (label->size == 0)
           || (label->buffer[label->size - 1] != 0))
            CryptDigestUpdateInt(&hState.hashState, 1, 0);
        // Adding contextU
        if(contextU != NULL)
            HASH_DATA(&hState.hashState, contextU->size, contextU->buffer);
        // Adding contextV
        if(contextV != NULL)
            HASH_DATA(&hState.hashState, contextV->size, contextV->buffer);
        // Adding size in bits
        CryptDigestUpdateInt(&hState.hashState, 4, sizeInBits);

        // Complete and put the data in the buffer
        CryptHmacEnd(&hState, bytes, stream);
        stream = &stream[digestSize];
    }
    // Masking in the KDF is disabled. If the calling function wants something
    // less than even number of bytes, then the caller should do the masking 
    // because there is no universal way to do it here
    if(counterInOut != NULL)
        *counterInOut = counter;
    return generated;
}

//*** CryptKDFe()
// This function implements KDFe() as defined in TPM specification part 1.
//
// This function returns the number of bytes generated which may be zero.
//
// The 'Z' and 'keyStream' pointers are not allowed to be NULL. The other
// pointer values may be NULL. The value of 'sizeInBits' must be no larger
// than (2^18)-1 = 256K bits (32385 bytes).
// Any error in the processing of this command is considered fatal.
//  Return Type: UINT16
//     0            hash algorithm is not supported or is TPM_ALG_NULL
//    > 0           the number of bytes in the 'keyStream' buffer
//
LIB_EXPORT UINT16
CryptKDFe(
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm used in HMAC
    TPM2B           *Z,             // IN: Z
    const TPM2B     *label,         // IN: a label value for the KDF
    TPM2B           *partyUInfo,    // IN: PartyUInfo
    TPM2B           *partyVInfo,    // IN: PartyVInfo
    UINT32           sizeInBits,    // IN: size of generated key in bits
    BYTE            *keyStream      // OUT: key buffer
    )
{
    HASH_STATE       hashState;
    PHASH_DEF        hashDef = CryptGetHashDef(hashAlg);

    UINT32           counter = 0;       // counter value
    UINT16           hLen;
    BYTE            *stream = keyStream;
    INT16            bytes;             // number of bytes to generate

    pAssert(keyStream != NULL && Z != NULL && ((sizeInBits + 7) / 8) < INT16_MAX);
//
    hLen = hashDef->digestSize;
    bytes = (INT16)((sizeInBits + 7) / 8);
    if(hashAlg == TPM_ALG_NULL || bytes == 0)
        return 0;

    // Generate required bytes
    //The inner loop of that KDF uses:
    //  Hash[i] := H(counter | Z | OtherInfo) (5)
    // Where:
    //  Hash[i]         the hash generated on the i-th iteration of the loop.
    //  H()             an approved hash function
    //  counter         a 32-bit counter that is initialized to 1 and incremented
    //                  on each iteration
    //  Z               the X coordinate of the product of a public ECC key and a
    //                  different private ECC key.
    //  OtherInfo       a collection of qualifying data for the KDF defined below.
    //  In this specification, OtherInfo will be constructed by:
    //      OtherInfo := Use | PartyUInfo  | PartyVInfo
    for(; bytes > 0; stream = &stream[hLen], bytes = bytes - hLen)
    {
        if(bytes < hLen)
            hLen = bytes;
        counter++;
        // Do the hash
        CryptHashStart(&hashState, hashAlg);
        // Add counter
        CryptDigestUpdateInt(&hashState, 4, counter);

        // Add Z
        if(Z != NULL)
            CryptDigestUpdate2B(&hashState, Z);
        // Add label
        if(label != NULL)
            CryptDigestUpdate2B(&hashState, label);
        // Add a null. SP108 is not very clear about when the 0 is needed but to
        // make this like the previous version that did not add an 0x00 after
        // a null-terminated string, this version will only add a null byte
        // if the label parameter did not end in a null byte, or if no label
        // is present.
        if((label == NULL)
           || (label->size == 0)
           || (label->buffer[label->size - 1] != 0))
            CryptDigestUpdateInt(&hashState, 1, 0);
        // Add PartyUInfo
        if(partyUInfo != NULL)
            CryptDigestUpdate2B(&hashState, partyUInfo);

        // Add PartyVInfo
        if(partyVInfo != NULL)
            CryptDigestUpdate2B(&hashState, partyVInfo);

        // Compute Hash. hLen was changed to be the smaller of bytes or hLen
        // at the start of each iteration.
        CryptHashEnd(&hashState, hLen, stream);
    }

    // Mask off bits if the required bits is not a multiple of byte size
    if((sizeInBits % 8) != 0)
        keyStream[0] &= ((1 << (sizeInBits % 8)) - 1);

    return (UINT16)((sizeInBits + 7) / 8);
}