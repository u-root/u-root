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
// This file contains the implementation of the message authentication codes based
// on a symmetric block cipher. These functions only use the single block 
// encryption functions of the selected symmetric cryptographic library.

//** Includes, Defines, and Typedefs
#define _CRYPT_HASH_C_
#include "Tpm.h"
#include "CryptSym.h"

#if ALG_CMAC

//** Functions

//*** CryptCmacStart()
// This is the function to start the CMAC sequence operation. It initializes the
// dispatch functions for the data and end operations for CMAC and initializes the
// parameters that are used for the processing of data, including the key, key size
// and block cipher algorithm.
UINT16
CryptCmacStart(
    SMAC_STATE          *state,
    TPMU_PUBLIC_PARMS   *keyParms,
    TPM_ALG_ID           macAlg,
    TPM2B               *key
)
{
    tpmCmacState_t      *cState = &state->state.cmac;
    TPMT_SYM_DEF_OBJECT *def = &keyParms->symDetail.sym;
//
    if(macAlg != TPM_ALG_CMAC)
        return 0;
    // set up the encryption algorithm and parameters
    cState->symAlg = def->algorithm;
    cState->keySizeBits = def->keyBits.sym;
    cState->iv.t.size = CryptGetSymmetricBlockSize(def->algorithm, 
                                                   def->keyBits.sym);
    MemoryCopy2B(&cState->symKey.b, key, sizeof(cState->symKey.t.buffer));

    // Set up the dispatch methods for the CMAC
    state->smacMethods.data = CryptCmacData;
    state->smacMethods.end = CryptCmacEnd;
    return cState->iv.t.size;
}


//*** CryptCmacData()
// This function is used to add data to the CMAC sequence computation. The function
// will XOR new data into the IV. If the buffer is full, and there is additional
// input data, the data is encrypted into the IV buffer, the new data is then
// XOR into the IV. When the data runs out, the function returns without encrypting
// even if the buffer is full. The last data block of a sequence will not be
// encrypted until the call to CryptCmacEnd(). This is to allow the proper subkey
// to be computed and applied before the last block is encrypted.
void
CryptCmacData(
    SMAC_STATES         *state,
    UINT32               size,
    const BYTE          *buffer
)
{
    tpmCmacState_t          *cmacState = &state->cmac;
    TPM_ALG_ID               algorithm = cmacState->symAlg;
    BYTE                    *key = cmacState->symKey.t.buffer;
    UINT16                   keySizeInBits = cmacState->keySizeBits;
    tpmCryptKeySchedule_t    keySchedule;
    TpmCryptSetSymKeyCall_t  encrypt;
//
    SELECT(ENCRYPT);
    while(size > 0)
    {
        if(cmacState->bcount == cmacState->iv.t.size)
        {
            ENCRYPT(&keySchedule, cmacState->iv.t.buffer, cmacState->iv.t.buffer);
            cmacState->bcount = 0;
        }
        for(;(size > 0) && (cmacState->bcount < cmacState->iv.t.size);
            size--, cmacState->bcount++)
        {
            cmacState->iv.t.buffer[cmacState->bcount] ^= *buffer++;
        }
    }
}

//*** CryptCmacEnd()
// This is the completion function for the CMAC. It does padding, if needed, and
// selects the subkey to be applied before the last block is encrypted.
UINT16
CryptCmacEnd(
    SMAC_STATES             *state,
    UINT32                   outSize,
    BYTE                    *outBuffer
)
{
    tpmCmacState_t          *cState = &state->cmac;
    // Need to set algorithm, key, and keySizeInBits in the local context so that  
    // the SELECT and ENCRYPT macros will work here
    TPM_ALG_ID               algorithm = cState->symAlg;
    BYTE                    *key = cState->symKey.t.buffer;
    UINT16                   keySizeInBits = cState->keySizeBits;
    tpmCryptKeySchedule_t    keySchedule;
    TpmCryptSetSymKeyCall_t  encrypt;
    TPM2B_IV                 subkey = {{0, {0}}};
    BOOL                     xorVal;
    UINT16                   i;

    subkey.t.size = cState->iv.t.size;
    // Encrypt a block of zero
    SELECT(ENCRYPT);
    ENCRYPT(&keySchedule, subkey.t.buffer, subkey.t.buffer);

    // shift left by 1 and XOR with 0x0...87 if the MSb was 0
    xorVal = ((subkey.t.buffer[0] & 0x80) == 0) ? 0 : 0x87;
    ShiftLeft(&subkey.b);
    subkey.t.buffer[subkey.t.size - 1] ^= xorVal;
    // this is a sanity check to make sure that the algorithm is working properly.
    // remove this check when debug is done
    pAssert(cState->bcount <= cState->iv.t.size);
    // If the buffer is full then no need to compute subkey 2.
    if(cState->bcount < cState->iv.t.size)
    {
        //Pad the data
        cState->iv.t.buffer[cState->bcount++] ^= 0x80;
        // The rest of the data is a pad of zero which would simply be XORed
        // with the iv value so nothing to do...
        // Now compute K2
        xorVal = ((subkey.t.buffer[0] & 0x80) == 0) ? 0 : 0x87;
        ShiftLeft(&subkey.b);
        subkey.t.buffer[subkey.t.size - 1] ^= xorVal;
    }
    // XOR the subkey into the IV
    for(i = 0; i < subkey.t.size; i++)
        cState->iv.t.buffer[i] ^= subkey.t.buffer[i];
    ENCRYPT(&keySchedule, cState->iv.t.buffer, cState->iv.t.buffer);
    i = (UINT16)MIN(cState->iv.t.size, outSize);
    MemoryCopy(outBuffer, cState->iv.t.buffer, i);

    return i;
}
#endif

