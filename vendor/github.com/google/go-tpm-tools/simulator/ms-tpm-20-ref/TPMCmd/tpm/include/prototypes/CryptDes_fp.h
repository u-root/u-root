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

#ifndef    _CRYPT_DES_FP_H_
#define    _CRYPT_DES_FP_H_

#if ALG_TDES

//*** CryptSetOddByteParity()
// This function sets the per byte parity of a 64-bit value. The least-significant
// bit is of each byte is replaced with the odd parity of the other 7 bits in the
// byte. With odd parity, no byte will ever be 0x00.
UINT64
CryptSetOddByteParity(
    UINT64          k
);

//*** CryptDesValidateKey()
// Function to check to see if the input key is a valid DES key where the definition
// of valid is that none of the elements are on the list of weak, semi-weak, or
// possibly weak keys; and that for two keys, K1!=K2, and for three keys that
// K1!=K2 and K2!=K3.
BOOL
CryptDesValidateKey(
    TPM2B_SYM_KEY       *desKey     // IN: key to validate
);

//*** CryptGenerateKeyDes()
// This function is used to create a DES key of the appropriate size. The key will
// have odd parity in the bytes.
TPM_RC
CryptGenerateKeyDes(
    TPMT_PUBLIC             *publicArea,        // IN/OUT: The public area template
                                                //     for the new key.
    TPMT_SENSITIVE          *sensitive,         // OUT: sensitive area
    RAND_STATE              *rand               // IN: the "entropy" source for
);
#endif

#endif  // _CRYPT_DES_FP_H_
