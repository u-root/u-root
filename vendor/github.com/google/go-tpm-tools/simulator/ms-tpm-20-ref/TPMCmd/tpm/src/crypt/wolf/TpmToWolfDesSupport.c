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
// The functions in this file are used for initialization of the interface to the
// wolfcrypt library.

//** Defines and Includes

#include "Tpm.h"

#if (defined SYM_LIB_WOLF) && ALG_TDES

//**Functions
//** TDES_setup
// This function calls the wolfcrypt function to generate a TDES key schedule. If the
// If the key is two key (16 bytes), then the first DES key is replicated to the third
// key position.
int TDES_setup(
    const BYTE          *key, 
    UINT32               keyBits, 
    tpmKeyScheduleTDES       *skey,
    int dir
    )
{
    BYTE                 k[24];
    BYTE                *kp;

    // If this is two-key, make it three key by replicating K1
    if(keyBits == 128)
    {
        memcpy(k, key, 16);
        memcpy(&k[16], key, 8);
        kp = k;
    }
    else
        kp = (BYTE *)key;

    return wc_Des3_SetKey( skey, kp, 0, dir );
}

//** TDES_setup_encrypt_key
// This function calls into TDES_setup(), specifically for an encryption key.
int TDES_setup_encrypt_key(
    const BYTE          *key,
    UINT32               keyBits,
    tpmKeyScheduleTDES       *skey
)
{
    return TDES_setup( key, keyBits, skey, DES_ENCRYPTION );
}

//** TDES_setup_decrypt_key
// This function calls into TDES_setup(), specifically for an decryption key.
int TDES_setup_decrypt_key(
    const BYTE          *key,
    UINT32               keyBits,
    tpmKeyScheduleTDES       *skey
)
{
    return TDES_setup( key, keyBits, skey, DES_DECRYPTION );
}

//*** TDES_encyrpt()
void TDES_encrypt(
    const BYTE              *in, 
    BYTE                    *out,
    tpmKeyScheduleTDES      *ks
    )
{
    wc_Des3_EcbEncrypt( ks, out, in, DES_BLOCK_SIZE );
}

//*** TDES_decrypt()
void TDES_decrypt(
    const BYTE          *in,
    BYTE                *out,
    tpmKeyScheduleTDES   *ks
    )
{
    wc_Des3_EcbDecrypt( ks, out, in, DES_BLOCK_SIZE );
}

#endif // MATH_LIB_WOLF && ALG_TDES
