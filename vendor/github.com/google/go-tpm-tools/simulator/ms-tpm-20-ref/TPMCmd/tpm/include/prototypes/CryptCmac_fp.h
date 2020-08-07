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
/*(Auto-generated)
 *  Created by TpmPrototypes; Version 3.0 July 18, 2017
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _CRYPT_CMAC_FP_H_
#define    _CRYPT_CMAC_FP_H_

#if ALG_CMAC

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
);

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
);

//*** CryptCmacEnd()
// This is the completion function for the CMAC. It does padding, if needed, and
// selects the subkey to be applied before the last block is encrypted.
UINT16
CryptCmacEnd(
    SMAC_STATES             *state,
    UINT32                   outSize,
    BYTE                    *outBuffer
);
#endif

#endif  // _CRYPT_CMAC_FP_H_
