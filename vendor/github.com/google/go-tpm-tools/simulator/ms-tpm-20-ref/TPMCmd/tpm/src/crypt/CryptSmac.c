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

#if SMAC_IMPLEMENTED

//*** CryptSmacStart()
// Function to start an SMAC.
UINT16
CryptSmacStart(
    HASH_STATE              *state,
    TPMU_PUBLIC_PARMS       *keyParameters,
    TPM_ALG_ID               macAlg,          // IN: the type of MAC
    TPM2B                   *key
)
{
    UINT16                  retVal = 0;
//
    // Make sure that the key size is correct. This should have been checked
    // at key load, but...
    if(BITS_TO_BYTES(keyParameters->symDetail.sym.keyBits.sym) == key->size)
    {
        switch(macAlg)
        {
#if ALG_CMAC
            case ALG_CMAC_VALUE:
                retVal = CryptCmacStart(&state->state.smac, keyParameters, 
                                        macAlg, key);
                break;
#endif
            default:
                break;
        }
    }
    state->type = (retVal != 0) ? HASH_STATE_SMAC : HASH_STATE_EMPTY;
    return retVal;
}

//*** CryptMacStart()
// Function to start either an HMAC or an SMAC. Cannot reuse the CryptHmacStart
// function because of the difference in number of parameters.
UINT16
CryptMacStart(
    HMAC_STATE              *state,
    TPMU_PUBLIC_PARMS       *keyParameters,
    TPM_ALG_ID               macAlg,          // IN: the type of MAC
    TPM2B                   *key
)
{
    MemorySet(state, 0, sizeof(HMAC_STATE));
    if(CryptHashIsValidAlg(macAlg, FALSE))
    {
        return CryptHmacStart(state, macAlg, key->size, key->buffer);
    }
    else if(CryptSmacIsValidAlg(macAlg, FALSE))
    {
        return CryptSmacStart(&state->hashState, keyParameters, macAlg, key);
    }
    else
        return 0;
}

//*** CryptMacEnd()
// Dispatch to the MAC end function using a size and buffer pointer.
UINT16
CryptMacEnd(
    HMAC_STATE          *state,
    UINT32               size,
    BYTE                *buffer
)
{
    UINT16              retVal = 0;
    if(state->hashState.type == HASH_STATE_SMAC)
        retVal = (state->hashState.state.smac.smacMethods.end)(
                    &state->hashState.state.smac.state, size, buffer);
    else if(state->hashState.type == HASH_STATE_HMAC)
        retVal = CryptHmacEnd(state, size, buffer);
    state->hashState.type = HASH_STATE_EMPTY;
    return retVal;
}

//*** CryptMacEnd2B()
// Dispatch to the MAC end function using a 2B.
UINT16
CryptMacEnd2B (
    HMAC_STATE          *state,
    TPM2B               *data
)
{
    return CryptMacEnd(state, data->size, data->buffer);
}
#endif // SMAC_IMPLEMENTED
