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
// This file implements a DRBG with a behavior according to SP800-90A using
// a block cypher. This is also compliant to ISO/IEC 18031:2011(E) C.3.2.
//
// A state structure is created for use by TPM.lib and functions
// within the CryptoEngine my use their own state structures when they need to have
// deterministic values.
//
// A debug mode is available that allows the random numbers generated for TPM.lib
// to be repeated during runs of the simulator. The switch for it is in 
// TpmBuildSwitches.h. It is USE_DEBUG_RNG.
//
//
// This is the implementation layer of CTR DRGB mechanism as defined in SP800-90A 
// and the functions are organized as closely as practical to the organization in 
// SP800-90A. It is intended to be compiled as a separate module that is linked 
// with a secure application so that both reside inside the same boundary 
// [SP 800-90A 8.5]. The secure application in particular manages the accesses
// protected storage for the state of the DRBG instantiations, and supplies the
// implementation functions here with a valid pointer to the working state of the
// given instantiations (as a DRBG_STATE structure).
//
// This DRBG mechanism implementation does not support prediction resistance. Thus
// 'prediction_resistance_flag' is omitted from Instantiate_function(),
// Reseed_function(), Generate_function() argument lists [SP 800-90A 9.1, 9.2,
// 9.3], as well as from the working state data structure DRBG_STATE [SP 800-90A
// 9.1].
//
// This DRBG mechanism implementation always uses the highest security strength of
// available in the block ciphers. Thus 'requested_security_strength' parameter is
// omitted from Instantiate_function() and Generate_function() argument lists
// [SP 800-90A 9.1, 9.2, 9.3], as well as from the working state data structure
// DRBG_STATE [SP 800-90A 9.1].
//
// Internal functions (ones without Crypt prefix) expect validated arguments and
// therefore use assertions instead of runtime parameter checks and mostly return
// void instead of a status value.

#include "Tpm.h"

// Pull in the test vector definitions and define the space
#include    "PRNG_TestVectors.h"

const BYTE DRBG_NistTestVector_Entropy[] = {DRBG_TEST_INITIATE_ENTROPY};
const BYTE DRBG_NistTestVector_GeneratedInterm[] = 
                                {DRBG_TEST_GENERATED_INTERM};

const BYTE DRBG_NistTestVector_EntropyReseed[] = 
                                {DRBG_TEST_RESEED_ENTROPY};
const BYTE DRBG_NistTestVector_Generated[] = {DRBG_TEST_GENERATED};

//** Derivation Functions
//*** Description
// The functions in this section are used to reduce the personalization input values 
// to make them usable as input for reseeding and instantiation. The overall
// behavior is intended to produce the same results as described in SP800-90A,
// section 10.4.2 "Derivation Function Using a Block Cipher Algorithm
// (Block_Cipher_df)." The code is broken into several subroutines to deal with the
// fact that the data used for personalization may come in several separate blocks
// such as a Template hash and a proof value and a primary seed.

//*** Derivation Function Defines and Structures

#define     DF_COUNT (DRBG_KEY_SIZE_WORDS / DRBG_IV_SIZE_WORDS + 1)
#if DRBG_KEY_SIZE_BITS != 128 && DRBG_KEY_SIZE_BITS != 256
#   error "CryptRand.c only written for AES with 128- or 256-bit keys."
#endif

typedef struct
{
    DRBG_KEY_SCHEDULE   keySchedule;
    DRBG_IV             iv[DF_COUNT];
    DRBG_IV             out1;
    DRBG_IV             buf;
    int                 contents;
} DF_STATE, *PDF_STATE;

//*** DfCompute()
// This function does the incremental update of the derivation function state. It
// encrypts the 'iv' value and XOR's the results into each of the blocks of the
// output. This is equivalent to processing all of input data for each output block.
static void
DfCompute(
    PDF_STATE        dfState
    )
{
    int              i;
    int              iv;
    crypt_uword_t   *pIv;
    crypt_uword_t    temp[DRBG_IV_SIZE_WORDS] = {0};
//
    for(iv = 0; iv < DF_COUNT; iv++)
    {
        pIv = (crypt_uword_t *)&dfState->iv[iv].words[0];
        for(i = 0; i < DRBG_IV_SIZE_WORDS; i++)
        {
            temp[i] ^= pIv[i] ^ dfState->buf.words[i];
        }
        DRBG_ENCRYPT(&dfState->keySchedule, &temp, pIv);
    }
    for(i = 0; i < DRBG_IV_SIZE_WORDS; i++)
        dfState->buf.words[i] = 0;
    dfState->contents = 0;
}

//*** DfStart()
// This initializes the output blocks with an encrypted counter value and
// initializes the key schedule.
static void
DfStart(
    PDF_STATE        dfState,
    uint32_t         inputLength
    )
{
    BYTE            init[8];
    int             i;
    UINT32          drbgSeedSize = sizeof(DRBG_SEED);

    const BYTE dfKey[DRBG_KEY_SIZE_BYTES] = {
        0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
        0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f
    #if DRBG_KEY_SIZE_BYTES > 16
        ,0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
        0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f
    #endif
    };
    memset(dfState, 0, sizeof(DF_STATE));
    DRBG_ENCRYPT_SETUP(&dfKey[0], DRBG_KEY_SIZE_BITS, &dfState->keySchedule);
    // Create the first chaining values
    for(i = 0; i < DF_COUNT; i++)
        ((BYTE *)&dfState->iv[i])[3] = (BYTE)i;
    DfCompute(dfState);
    // initialize the first 64 bits of the IV in a way that doesn't depend
    // on the size of the words used.
    UINT32_TO_BYTE_ARRAY(inputLength, init);
    UINT32_TO_BYTE_ARRAY(drbgSeedSize, &init[4]);
    memcpy(&dfState->iv[0], init, 8);
    dfState->contents = 4;
}

//*** DfUpdate()
// This updates the state with the input data. A byte at a time is moved into the
// state buffer until it is full and then that block is encrypted by DfCompute().
static void
DfUpdate(
    PDF_STATE        dfState,
    int              size,
    const BYTE      *data
    )
{
    while(size > 0)
    {
        int         toFill = DRBG_IV_SIZE_BYTES - dfState->contents;
        if(size < toFill)
            toFill = size;
        // Copy as many bytes as there are or until the state buffer is full
        memcpy(&dfState->buf.bytes[dfState->contents], data, toFill);
        // Reduce the size left by the amount copied
        size -= toFill;
        // Advance the data pointer by the amount copied
        data += toFill;
        // increase the buffer contents count by the amount copied
        dfState->contents += toFill;
        pAssert(dfState->contents <= DRBG_IV_SIZE_BYTES);
        // If we have a full buffer, do a computation pass.
        if(dfState->contents == DRBG_IV_SIZE_BYTES)
            DfCompute(dfState);
    }
}

//*** DfEnd()
// This function is called to get the result of the derivation function computation.
// If the buffer is not full, it is padded with zeros. The output buffer is
// structured to be the same as a DRBG_SEED value so that the function can return
// a pointer to the DRBG_SEED value in the DF_STATE structure.
static DRBG_SEED *
DfEnd(
    PDF_STATE        dfState
    )
{
    // Since DfCompute is always called when a buffer is full, there is always
    // space in the buffer for the terminator
    dfState->buf.bytes[dfState->contents++] = 0x80;
    // If the buffer is not full, pad with zeros
    while(dfState->contents < DRBG_IV_SIZE_BYTES)
        dfState->buf.bytes[dfState->contents++] = 0;
    // Do a final state update
    DfCompute(dfState);
    return (DRBG_SEED *)&dfState->iv;
}

//*** DfBuffer()
// Function to take an input buffer and do the derivation function to produce a
// DRBG_SEED value that can be used in DRBG_Reseed();
static DRBG_SEED *
DfBuffer(
    DRBG_SEED       *output,        // OUT: receives the result
    int              size,          // IN: size of the buffer to add
    BYTE            *buf            // IN: address of the buffer
    )
{
    DF_STATE        dfState;
    if(size == 0 || buf == NULL)
        return NULL;
    // Initialize the derivation function
    DfStart(&dfState, size);
    DfUpdate(&dfState, size, buf);
    DfEnd(&dfState);
    memcpy(output, &dfState.iv[0], sizeof(DRBG_SEED));
    return output;
}

//*** DRBG_GetEntropy()
// Even though this implementation never fails, it may get blocked
// indefinitely long in the call to get entropy from the platform
// (DRBG_GetEntropy32()).
// This function is only used during instantiation of the DRBG for
// manufacturing and on each start-up after an non-orderly shutdown.
//  Return Type: BOOL
//      TRUE(1)         requested entropy returned
//      FALSE(0)        entropy Failure
BOOL
DRBG_GetEntropy(
    UINT32           requiredEntropy,   // IN: requested number of bytes of full
                                        //     entropy
    BYTE            *entropy            // OUT: buffer to return collected entropy
    )
{
#if !USE_DEBUG_RNG

    UINT32       obtainedEntropy;
    INT32        returnedEntropy;

// If in debug mode, always use the self-test values for initialization
    if(IsSelfTest())
    {
#endif
        // If doing simulated DRBG, then check to see if the
        // entropyFailure condition is being tested
        if(!IsEntropyBad())
        {
            // In self-test, the caller should be asking for exactly the seed
            // size of entropy.
            pAssert(requiredEntropy == sizeof(DRBG_NistTestVector_Entropy));
            memcpy(entropy, DRBG_NistTestVector_Entropy,
                   sizeof(DRBG_NistTestVector_Entropy));
        }
#if !USE_DEBUG_RNG
    }
    else if(!IsEntropyBad())
    {
        // Collect entropy
        // Note: In debug mode, the only "entropy" value ever returned
        // is the value of the self-test vector.
        for(returnedEntropy = 1, obtainedEntropy = 0;
            obtainedEntropy < requiredEntropy && !IsEntropyBad();
            obtainedEntropy += returnedEntropy)
        {
            returnedEntropy = _plat__GetEntropy(&entropy[obtainedEntropy],
                                                requiredEntropy - obtainedEntropy);
            if(returnedEntropy <= 0)
                SetEntropyBad();
        }
    }
#endif
    return !IsEntropyBad();
}

//*** IncrementIv()
// This function increments the IV value by 1. It is used by EncryptDRBG().
void
IncrementIv(
    DRBG_IV         *iv
    )
{
    BYTE      *ivP = ((BYTE *)iv) + DRBG_IV_SIZE_BYTES;
    while((--ivP >= (BYTE *)iv) && ((*ivP = ((*ivP + 1) & 0xFF)) == 0));
}

//*** EncryptDRBG()
// This does the encryption operation for the DRBG. It will encrypt
// the input state counter (IV) using the state key. Into the output
// buffer for as many times as it takes to generate the required
// number of bytes.
static BOOL
EncryptDRBG(
    BYTE                *dOut,
    UINT32               dOutBytes,
    DRBG_KEY_SCHEDULE   *keySchedule,
    DRBG_IV             *iv,
    UINT32              *lastValue      // Points to the last output value
    )
{
#if FIPS_COMPLIANT
// For FIPS compliance, the DRBG has to do a continuous self-test to make sure that
// no two consecutive values are the same. This overhead is not incurred if the TPM
// is not required to be FIPS compliant
//
    UINT32           temp[DRBG_IV_SIZE_BYTES / sizeof(UINT32)];
    int              i;
    BYTE            *p;

    for(; dOutBytes > 0;)
    {
        // Increment the IV before each encryption (this is what makes this
        // different from normal counter-mode encryption
        IncrementIv(iv);
        DRBG_ENCRYPT(keySchedule, iv, temp);
// Expect a 16 byte block
#if DRBG_IV_SIZE_BITS != 128
#error  "Unsuppored IV size in DRBG"
#endif
        if((lastValue[0] == temp[0])
            && (lastValue[1] == temp[1])
            && (lastValue[2] == temp[2])
            && (lastValue[3] == temp[3])
            )
        {
            LOG_FAILURE(FATAL_ERROR_ENTROPY);
            return FALSE;
        }
        lastValue[0] = temp[0];
        lastValue[1] = temp[1];
        lastValue[2] = temp[2];
        lastValue[3] = temp[3];
        i = MIN(dOutBytes, DRBG_IV_SIZE_BYTES);
        dOutBytes -= i;
        for(p = (BYTE *)temp; i > 0; i--)
            *dOut++ = *p++;
    }
#else // version without continuous self-test
    NOT_REFERENCED(lastValue);
    for(; dOutBytes >= DRBG_IV_SIZE_BYTES;
    dOut = &dOut[DRBG_IV_SIZE_BYTES], dOutBytes -= DRBG_IV_SIZE_BYTES)
    {
        // Increment the IV
        IncrementIv(iv);
        DRBG_ENCRYPT(keySchedule, iv, dOut);
    }
    // If there is a partial, generate into a block-sized
    // temp buffer and copy to the output.
    if(dOutBytes != 0)
    {
        BYTE        temp[DRBG_IV_SIZE_BYTES];
        // Increment the IV
        IncrementIv(iv);
        DRBG_ENCRYPT(keySchedule, iv, temp);
        memcpy(dOut, temp, dOutBytes);
    }
#endif
    return TRUE;
}

//*** DRBG_Update()
// This function performs the state update function.
// According to SP800-90A, a temp value is created by doing CTR mode
// encryption of 'providedData' and replacing the key and IV with
// these values. The one difference is that, with counter mode, the
// IV is incremented after each block is encrypted and in this
// operation, the counter is incremented before each block is
// encrypted. This function implements an 'optimized' version
// of the algorithm in that it does the update of the drbgState->seed
// in place and then 'providedData' is XORed into drbgState->seed
// to complete the encryption of 'providedData'. This works because
// the IV is the last thing that gets encrypted.
//
static BOOL 
DRBG_Update(
    DRBG_STATE          *drbgState,     // IN:OUT state to update
    DRBG_KEY_SCHEDULE   *keySchedule,   // IN: the key schedule (optional)
    DRBG_SEED           *providedData   // IN: additional data
    )
{
    UINT32               i;
    BYTE                *temp = (BYTE *)&drbgState->seed;
    DRBG_KEY            *key = pDRBG_KEY(&drbgState->seed);
    DRBG_IV             *iv = pDRBG_IV(&drbgState->seed);
    DRBG_KEY_SCHEDULE    localKeySchedule;
//
    pAssert(drbgState->magic == DRBG_MAGIC);

    // If an key schedule was not provided, make one
    if(keySchedule == NULL)
    {
        if(DRBG_ENCRYPT_SETUP((BYTE *)key,
            DRBG_KEY_SIZE_BITS, &localKeySchedule) != 0)
        {
            LOG_FAILURE(FATAL_ERROR_INTERNAL);
            return FALSE;
        }
        keySchedule = &localKeySchedule;
    }
    // Encrypt the temp value

    EncryptDRBG(temp, sizeof(DRBG_SEED), keySchedule, iv,
                drbgState->lastValue);
    if(providedData != NULL)
    {
        BYTE        *pP = (BYTE *)providedData;
        for(i = DRBG_SEED_SIZE_BYTES; i != 0; i--)
            *temp++ ^= *pP++;
    }
    // Since temp points to the input key and IV, we are done and
    // don't need to copy the resulting 'temp' to drbgState->seed
    return TRUE;
}

//*** DRBG_Reseed()
// This function is used when reseeding of the DRBG is required. If
// entropy is provided, it is used in lieu of using hardware entropy.
// Note: the provided entropy must be the required size.
//  Return Type: BOOL
//      TRUE(1)         reseed succeeded
//      FALSE(0)        reseed failed, probably due to the entropy generation
BOOL
DRBG_Reseed(
    DRBG_STATE          *drbgState,         // IN: the state to update
    DRBG_SEED           *providedEntropy,   // IN: entropy
    DRBG_SEED           *additionalData     // IN:
    )
{
    DRBG_SEED            seed;

    pAssert((drbgState != NULL) && (drbgState->magic == DRBG_MAGIC));

    if(providedEntropy == NULL)
    {
        providedEntropy = &seed;
        if(!DRBG_GetEntropy(sizeof(DRBG_SEED), (BYTE *)providedEntropy))
            return FALSE;
    }
    if(additionalData != NULL)
    {
        unsigned int          i;

        // XOR the provided data into the provided entropy
        for(i = 0; i < sizeof(DRBG_SEED); i++)
            ((BYTE *)providedEntropy)[i] ^= ((BYTE *)additionalData)[i];
    }
    DRBG_Update(drbgState, NULL, providedEntropy);

    drbgState->reseedCounter = 1;

    return TRUE;
}

//*** DRBG_SelfTest()
// This is run when the DRBG is instantiated and at startup
//  Return Type: BOOL
//      TRUE(1)         test OK
//      FALSE(0)        test failed
BOOL
DRBG_SelfTest(
    void
    )
{
    BYTE             buf[sizeof(DRBG_NistTestVector_Generated)];
    DRBG_SEED        seed;
    UINT32           i;
    BYTE            *p;
    DRBG_STATE       testState;
//
    pAssert(!IsSelfTest());

    SetSelfTest();
    SetDrbgTested();
    // Do an instantiate
    if(!DRBG_Instantiate(&testState, 0, NULL))
        return FALSE;
#if DRBG_DEBUG_PRINT
    dbgDumpMemBlock(pDRBG_KEY(&testState), DRBG_KEY_SIZE_BYTES,
                    "Key after Instantiate");
    dbgDumpMemBlock(pDRBG_IV(&testState), DRBG_IV_SIZE_BYTES,
                    "Value after Instantiate");
#endif
    if(DRBG_Generate((RAND_STATE *)&testState, buf, sizeof(buf)) == 0)
        return FALSE;
#if DRBG_DEBUG_PRINT
    dbgDumpMemBlock(pDRBG_KEY(&testState.seed), DRBG_KEY_SIZE_BYTES,
                    "Key after 1st Generate");
    dbgDumpMemBlock(pDRBG_IV(&testState.seed), DRBG_IV_SIZE_BYTES,
                    "Value after 1st Generate");
#endif
    if(memcmp(buf, DRBG_NistTestVector_GeneratedInterm, sizeof(buf)) != 0)
        return FALSE;
    memcpy(seed.bytes, DRBG_NistTestVector_EntropyReseed, sizeof(seed));
    DRBG_Reseed(&testState, &seed, NULL);
#if DRBG_DEBUG_PRINT
    dbgDumpMemBlock((BYTE *)pDRBG_KEY(&testState.seed), DRBG_KEY_SIZE_BYTES,
                    "Key after 2nd Generate");
    dbgDumpMemBlock((BYTE *)pDRBG_IV(&testState.seed), DRBG_IV_SIZE_BYTES,
                    "Value after 2nd Generate");
    dbgDumpMemBlock(buf, sizeof(buf), "2nd Generated");
#endif
    if(DRBG_Generate((RAND_STATE *)&testState, buf, sizeof(buf)) == 0)
        return FALSE;
    if(memcmp(buf, DRBG_NistTestVector_Generated, sizeof(buf)) != 0)
        return FALSE;
    ClearSelfTest();

    DRBG_Uninstantiate(&testState);
    for(p = (BYTE *)&testState, i = 0; i < sizeof(DRBG_STATE); i++)
    {
        if(*p++)
            return FALSE;
    }
    // Simulate hardware failure to make sure that we get an error when
    // trying to instantiate
    SetEntropyBad();
    if(DRBG_Instantiate(&testState, 0, NULL))
       return FALSE;
    ClearEntropyBad();

    return TRUE;
}

//** Public Interface
//*** Description
// The functions in this section are the interface to the RNG. These
// are the functions that are used by TPM.lib. 

//*** CryptRandomStir()
// This function is used to cause a reseed. A DRBG_SEED amount of entropy is
// collected from the hardware and then additional data is added.
//  Return Type: TPM_RC
//      TPM_RC_NO_RESULT        failure of the entropy generator
LIB_EXPORT TPM_RC
CryptRandomStir(
    UINT16           additionalDataSize,
    BYTE            *additionalData
    )
{
#if !USE_DEBUG_RNG 
    DRBG_SEED        tmpBuf;
    DRBG_SEED        dfResult;
//
    // All reseed with outside data starts with a buffer full of entropy
    if(!DRBG_GetEntropy(sizeof(tmpBuf), (BYTE *)&tmpBuf))
        return TPM_RC_NO_RESULT;

    DRBG_Reseed(&drbgDefault, &tmpBuf,
                DfBuffer(&dfResult, additionalDataSize, additionalData));
    drbgDefault.reseedCounter = 1;

    return TPM_RC_SUCCESS;

#else 
    // If doing debug, use the input data as the initial setting for the RNG state
    // so that the test can be reset at any time.
    // Note: If this is called with a data size of 0 or less, nothing happens. The
    // presumption is that, in a debug environment, the caller will have specific
    // values for initialization, so this check is just a simple way to prevent
    // inadvertent programming errors from screwing things up. This doesn't use an
    // pAssert() because the non-debug version of this function will accept these
    // parameters as meaning that there is no additionalData and only hardware
    // entropy is used. 
    if((additionalDataSize > 0) && (additionalData != NULL))
    {
        memset(drbgDefault.seed.bytes, 0, sizeof(drbgDefault.seed.bytes));
        memcpy(drbgDefault.seed.bytes, additionalData, 
               MIN(additionalDataSize, sizeof(drbgDefault.seed.bytes)));
    }
    drbgDefault.reseedCounter = 1;

    return TPM_RC_SUCCESS;
#endif
}

//*** CryptRandomGenerate()
// Generate a 'randomSize' number or random bytes.
LIB_EXPORT UINT16
CryptRandomGenerate(
    UINT16           randomSize,
    BYTE            *buffer
    )
{
    return DRBG_Generate((RAND_STATE *)&drbgDefault, buffer, randomSize);
}



//*** DRBG_InstantiateSeededKdf()
// This function is used to instantiate a KDF-based RNG. This is used for derivations.
// This function always returns TRUE.
LIB_EXPORT BOOL
DRBG_InstantiateSeededKdf(
    KDF_STATE       *state,         // OUT: buffer to hold the state
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm
    TPM_ALG_ID       kdf,           // IN: the KDF to use
    TPM2B           *seed,          // IN: the seed to use
    const TPM2B     *label,         // IN: a label for the generation process.
    TPM2B           *context,       // IN: the context value
    UINT32           limit          // IN: Maximum number of bits from the KDF
    )
{
    state->magic = KDF_MAGIC;
    state->limit = limit;
    state->seed = seed;
    state->hash = hashAlg;
    state->kdf = kdf;
    state->label = label;
    state->context = context;
    state->digestSize = CryptHashGetDigestSize(hashAlg);
    state->counter = 0;
    state->residual.t.size = 0;
    return TRUE;
}

//*** DRBG_AdditionalData()
// Function to reseed the DRBG with additional entropy. This is normally called
// before computing the protection value of a primary key in the Endorsement 
// hierarchy.
LIB_EXPORT void
DRBG_AdditionalData(
    DRBG_STATE      *drbgState,     // IN:OUT state to update
    TPM2B           *additionalData // IN: value to incorporate
    )
{
    DRBG_SEED        dfResult;
    if(drbgState->magic == DRBG_MAGIC)
    {
        DfBuffer(&dfResult, additionalData->size, additionalData->buffer);
        DRBG_Reseed(drbgState, &dfResult, NULL);
    }
}


//*** DRBG_InstantiateSeeded()
// This function is used to instantiate a random number generator from seed values.
// The nominal use of this generator is to create sequences of pseudo-random
// numbers from a seed value.
// Return Type: TPM_RC
//  TPM_RC_FAILURE      DRBG self-test failure
LIB_EXPORT TPM_RC
DRBG_InstantiateSeeded(
    DRBG_STATE      *drbgState,     // IN/OUT: buffer to hold the state
    const TPM2B     *seed,          // IN: the seed to use
    const TPM2B     *purpose,       // IN: a label for the generation process.
    const TPM2B     *name,          // IN: name of the object
    const TPM2B     *additional     // IN: additional data
    )
{
    DF_STATE         dfState;
    int              totalInputSize;
    // DRBG should have been tested, but...
    if(!IsDrbgTested() && !DRBG_SelfTest())
    {
        LOG_FAILURE(FATAL_ERROR_SELF_TEST);
        return TPM_RC_FAILURE;
    }
    // Initialize the DRBG state
    memset(drbgState, 0, sizeof(DRBG_STATE));
    drbgState->magic = DRBG_MAGIC;

    // Size all of the values
    totalInputSize = (seed != NULL) ? seed->size : 0;
    totalInputSize += (purpose != NULL) ? purpose->size : 0;
    totalInputSize += (name != NULL) ? name->size : 0;
    totalInputSize += (additional != NULL) ? additional->size : 0;

    // Initialize the derivation
    DfStart(&dfState, totalInputSize);

    // Run all the input strings through the derivation function
    if(seed != NULL)
        DfUpdate(&dfState, seed->size, seed->buffer);
    if(purpose != NULL)
        DfUpdate(&dfState, purpose->size, purpose->buffer);
    if(name != NULL)
        DfUpdate(&dfState, name->size, name->buffer);
    if(additional != NULL)
        DfUpdate(&dfState, additional->size, additional->buffer);

    // Used the derivation function output as the "entropy" input. This is not
    // how it is described in SP800-90A but this is the equivalent function
    DRBG_Reseed(((DRBG_STATE *)drbgState), DfEnd(&dfState), NULL);

    return TPM_RC_SUCCESS;
}

//*** CryptRandStartup()
// This function is called when TPM_Startup is executed. This function always returns
// TRUE.
LIB_EXPORT BOOL
CryptRandStartup(
    void
    )
{
#if ! _DRBG_STATE_SAVE
    // If not saved in NV, re-instantiate on each startup
    DRBG_Instantiate(&drbgDefault, 0, NULL);
#else
    // If the running state is saved in NV, NV has to be loaded before it can
    // be updated
    if(go.drbgState.magic == DRBG_MAGIC)
        DRBG_Reseed(&go.drbgState, NULL, NULL);
    else
        DRBG_Instantiate(&go.drbgState, 0, NULL);
#endif
    return TRUE;
}

//**** CryptRandInit()
// This function is called when _TPM_Init is being processed.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure
LIB_EXPORT BOOL
CryptRandInit(
    void
    )
{
#if !USE_DEBUG_RNG
    _plat__GetEntropy(NULL, 0);
#endif
    return DRBG_SelfTest();
}

//*** DRBG_Generate()
// This function generates a random sequence according SP800-90A.
// If 'random' is not NULL, then 'randomSize' bytes of random values are generated.
// If 'random' is NULL or 'randomSize' is zero, then the function returns
// zero without generating any bits or updating the reseed counter.
// This function returns the number of bytes produced which could be less than the 
// number requested if the request is too large ("too large" is implementation
// dependent.)
LIB_EXPORT UINT16
DRBG_Generate(
    RAND_STATE      *state,
    BYTE            *random,        // OUT: buffer to receive the random values
    UINT16           randomSize     // IN: the number of bytes to generate
    )
{
    if(state == NULL)
        state = (RAND_STATE *)&drbgDefault;
    if(random == NULL)
        return 0;

    // If the caller used a KDF state, generate a sequence from the KDF not to 
    // exceed the limit.
    if(state->kdf.magic == KDF_MAGIC)
    {
        KDF_STATE       *kdf = (KDF_STATE *)state;
        UINT32           counter = (UINT32)kdf->counter;
        INT32            bytesLeft = randomSize;
//
        // If the number of bytes to be returned would put the generator 
        // over the limit, then return 0
        if((((kdf->counter * kdf->digestSize) + randomSize) * 8) > kdf->limit)
            return 0;
        // Process partial and full blocks until all requested bytes provided
        while(bytesLeft > 0)
        {
            // If there is any residual data in the buffer, copy it to the output
            // buffer
            if(kdf->residual.t.size > 0)
            {
                INT32      size;
//
                // Don't use more of the residual than will fit or more than are
                // available
                size = MIN(kdf->residual.t.size, bytesLeft);
                
                // Copy some or all of the residual to the output. The residual is
                // at the end of the buffer. The residual might be a full buffer.
                MemoryCopy(random,
                           &kdf->residual.t.buffer
                           [kdf->digestSize - kdf->residual.t.size], size);
                
                // Advance the buffer pointer
                random += size;

                // Reduce the number of bytes left to get
                bytesLeft -= size;

                // And reduce the residual size appropriately
                kdf->residual.t.size -= (UINT16)size;
            } 
            else
            {
                UINT16           blocks = (UINT16)(bytesLeft / kdf->digestSize);
// 
                // Get the number of required full blocks
                if(blocks > 0)
                {
                    UINT16      size = blocks * kdf->digestSize;
// Get some number of full blocks and put them in the return buffer
                    CryptKDFa(kdf->hash, kdf->seed, kdf->label, kdf->context, NULL,
                              kdf->limit, random, &counter, blocks);

                    // reduce the size remaining to be moved and advance the pointer
                    bytesLeft -= size;
                    random += size;
                }
                else
                {
                    // Fill the residual buffer with a full block and then loop to
                    // top to get part of it copied to the output.
                    kdf->residual.t.size = CryptKDFa(kdf->hash, kdf->seed,
                                                     kdf->label, kdf->context, NULL,
                                                     kdf->limit,
                                                     kdf->residual.t.buffer,
                                                     &counter, 1);
                }
            }
        }
        kdf->counter = counter;
        return randomSize;
    }
    else if(state->drbg.magic == DRBG_MAGIC)
    {
        DRBG_STATE          *drbgState = (DRBG_STATE *)state;
        DRBG_KEY_SCHEDULE    keySchedule;
        DRBG_SEED           *seed = &drbgState->seed;

        if(drbgState->reseedCounter >= CTR_DRBG_MAX_REQUESTS_PER_RESEED)
        {
            if(drbgState == &drbgDefault)
            {
                DRBG_Reseed(drbgState, NULL, NULL);
                if(IsEntropyBad() && !IsSelfTest())
                    return 0;
            }
            else
            {
                // If this is a PRNG then the only way to get
                // here is if the SW has run away.
                LOG_FAILURE(FATAL_ERROR_INTERNAL);
                return 0;
            }
        }
        // if the allowed number of bytes in a request is larger than the
        // less than the number of bytes that can be requested, then check
#if UINT16_MAX >=  CTR_DRBG_MAX_BYTES_PER_REQUEST
        if(randomSize > CTR_DRBG_MAX_BYTES_PER_REQUEST)
            randomSize = CTR_DRBG_MAX_BYTES_PER_REQUEST;
#endif
        // Create  encryption schedule
        if(DRBG_ENCRYPT_SETUP((BYTE *)pDRBG_KEY(seed),
                              DRBG_KEY_SIZE_BITS, &keySchedule) != 0)
        {
            LOG_FAILURE(FATAL_ERROR_INTERNAL);
            return 0;
        }
        // Generate the random data
        EncryptDRBG(random, randomSize, &keySchedule, pDRBG_IV(seed),
                    drbgState->lastValue);
        // Do a key update
        DRBG_Update(drbgState, &keySchedule, NULL);

        // Increment the reseed counter
        drbgState->reseedCounter += 1;
    }
    else
    {
        LOG_FAILURE(FATAL_ERROR_INTERNAL);
        return FALSE;
    }
    return randomSize;
}

//*** DRBG_Instantiate()
// This is CTR_DRBG_Instantiate_algorithm() from [SP 800-90A 10.2.1.3.1].
// This is called when a the TPM DRBG is to be instantiated. This is
// called to instantiate a DRBG used by the TPM for normal
// operations.
//  Return Type: BOOL
//      TRUE(1)         instantiation succeeded
//      FALSE(0)        instantiation failed
LIB_EXPORT BOOL
DRBG_Instantiate(
    DRBG_STATE      *drbgState,         // OUT: the instantiated value
    UINT16           pSize,             // IN: Size of personalization string
    BYTE            *personalization    // IN: The personalization string
    )
{
    DRBG_SEED        seed;
    DRBG_SEED        dfResult;
//
    pAssert((pSize == 0) || (pSize <= sizeof(seed)) || (personalization != NULL));
    // If the DRBG has not been tested, test when doing an instantiation. Since
    // Instantiation is called during self test, make sure we don't get stuck in a
    // loop.
    if(!IsDrbgTested() && !IsSelfTest() && !DRBG_SelfTest())
        return FALSE;
    // If doing a self test, DRBG_GetEntropy will return the NIST
    // test vector value.
    if(!DRBG_GetEntropy(sizeof(seed), (BYTE *)&seed))
        return FALSE;
    // set everything to zero
    memset(drbgState, 0, sizeof(DRBG_STATE));
    drbgState->magic = DRBG_MAGIC;

    // Steps 1, 2, 3, 6, 7 of SP 800-90A 10.2.1.3.1 are exactly what
    // reseeding does. So, do a reduction on the personalization value (if any)
    // and do a reseed.
    DRBG_Reseed(drbgState, &seed, DfBuffer(&dfResult, pSize, personalization));

    return TRUE;
}

//*** DRBG_Uninstantiate()
// This is Uninstantiate_function() from [SP 800-90A 9.4].
//
//  Return Type: TPM_RC
//      TPM_RC_VALUE        not a valid state
LIB_EXPORT TPM_RC
DRBG_Uninstantiate(
    DRBG_STATE      *drbgState      // IN/OUT: working state to erase
    )
{
    if((drbgState == NULL) || (drbgState->magic != DRBG_MAGIC))
        return TPM_RC_VALUE;
    memset(drbgState, 0, sizeof(DRBG_STATE));
    return TPM_RC_SUCCESS;
}
