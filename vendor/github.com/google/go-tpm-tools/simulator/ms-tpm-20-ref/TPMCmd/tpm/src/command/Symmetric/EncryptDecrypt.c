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
#include "Tpm.h"
#include "EncryptDecrypt_fp.h"
#if CC_EncryptDecrypt2
#include  "EncryptDecrypt_spt_fp.h"
#endif

#if CC_EncryptDecrypt  // Conditional expansion of this file

/*(See part 3 specification)
// symmetric encryption or decryption
*/
//  Return Type: TPM_RC
//      TPM_RC_KEY          is not a symmetric decryption key with both
//                          public and private portions loaded
//      TPM_RC_SIZE         'IvIn' size is incompatible with the block cipher mode;
//                          or 'inData' size is not an even multiple of the block
//                          size for CBC or ECB mode
//      TPM_RC_VALUE        'keyHandle' is restricted and the argument 'mode' does
//                          not match the key's mode
TPM_RC
TPM2_EncryptDecrypt(
    EncryptDecrypt_In   *in,            // IN: input parameter list
    EncryptDecrypt_Out  *out            // OUT: output parameter list
    )
{
#if CC_EncryptDecrypt2
    return EncryptDecryptShared(in->keyHandle, in->decrypt, in->mode,
                                &in->ivIn, &in->inData, out);
#else
    OBJECT              *symKey;
    UINT16               keySize;
    UINT16               blockSize;
    BYTE                *key;
    TPM_ALG_ID           alg;
    TPM_ALG_ID           mode;
    TPM_RC               result;
    BOOL                 OK;
    TPMA_OBJECT          attributes;

// Input Validation
    symKey = HandleToObject(in->keyHandle);
    mode = symKey->publicArea.parameters.symDetail.sym.mode.sym;
    attributes = symKey->publicArea.objectAttributes;

    // The input key should be a symmetric key
    if(symKey->publicArea.type != TPM_ALG_SYMCIPHER)
        return TPM_RCS_KEY + RC_EncryptDecrypt_keyHandle;
    // The key must be unrestricted and allow the selected operation
    OK = IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted)
    if(YES == in->decrypt)
        OK = OK && IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt);
    else
        OK = OK && IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign);
    if(!OK)
        return TPM_RCS_ATTRIBUTES + RC_EncryptDecrypt_keyHandle;
    
    // If the key mode is not TPM_ALG_NULL...
    // or TPM_ALG_NULL
    if(mode != TPM_ALG_NULL)
    {
        // then the input mode has to be TPM_ALG_NULL or the same as the key
        if((in->mode != TPM_ALG_NULL) && (in->mode != mode))
            return TPM_RCS_MODE + RC_EncryptDecrypt_mode;
    }
    else
    {
        // if the key mode is null, then the input can't be null
        if(in->mode == TPM_ALG_NULL)
            return TPM_RCS_MODE + RC_EncryptDecrypt_mode;
        mode = in->mode;
    }
    // The input iv for ECB mode should be an Empty Buffer.  All the other modes
    // should have an iv size same as encryption block size
    keySize = symKey->publicArea.parameters.symDetail.sym.keyBits.sym;
    alg = symKey->publicArea.parameters.symDetail.sym.algorithm;
    blockSize = CryptGetSymmetricBlockSize(alg, keySize);
    
    // reverify the algorithm. This is mainly to keep static analysis tools happy
    if(blockSize == 0)
        return TPM_RCS_KEY + RC_EncryptDecrypt_keyHandle;

    // Note: When an algorithm is not supported by a TPM, the TPM_ALG_xxx for that
    // algorithm is not defined. However, it is assumed that the ALG_xxx_VALUE for
    // the algorithm is always defined. Both have the same numeric value.
    // ALG_xxx_VALUE is used here so that the code does not get cluttered with
    // #ifdef's. Having this check does not mean that the algorithm is supported.
    // If it was not supported the unmarshaling code would have rejected it before
    // this function were called. This means that, depending on the implementation,
    // the check could be redundant but it doesn't hurt.
    if(((mode == ALG_ECB_VALUE) && (in->ivIn.t.size != 0))
       || ((mode != ALG_ECB_VALUE) && (in->ivIn.t.size != blockSize)))
        return TPM_RCS_SIZE + RC_EncryptDecrypt_ivIn;

    // The input data size of CBC mode or ECB mode must be an even multiple of
    // the symmetric algorithm's block size
    if(((mode == ALG_CBC_VALUE) || (mode == ALG_ECB_VALUE))
       && ((in->inData.t.size % blockSize) != 0))
        return TPM_RCS_SIZE + RC_EncryptDecrypt_inData;

    // Copy IV
    // Note: This is copied here so that the calls to the encrypt/decrypt functions
    // will modify the output buffer, not the input buffer
    out->ivOut = in->ivIn;

// Command Output
    key = symKey->sensitive.sensitive.sym.t.buffer;
    // For symmetric encryption, the cipher data size is the same as plain data
    // size.
    out->outData.t.size = in->inData.t.size;
    if(in->decrypt == YES)
    {
        // Decrypt data to output
        result = CryptSymmetricDecrypt(out->outData.t.buffer, alg, keySize, key,
                                       &(out->ivOut), mode, in->inData.t.size,
                                       in->inData.t.buffer);
    }
    else
    {
        // Encrypt data to output
        result = CryptSymmetricEncrypt(out->outData.t.buffer, alg, keySize, key,
                                       &(out->ivOut), mode, in->inData.t.size,
                                       in->inData.t.buffer);
    }
    return result;
#endif // CC_EncryptDecrypt2

}

#endif // CC_EncryptDecrypt