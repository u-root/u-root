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
#include "Rewrap_fp.h"

#if CC_Rewrap  // Conditional expansion of this file

#include "Object_spt_fp.h"

/*(See part 3 specification)
// This command allows the TPM to serve in the role as an MA.
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES       'newParent' is not a decryption key
//      TPM_RC_HANDLE           'oldParent' does not consistent with inSymSeed
//      TPM_RC_INTEGRITY        the integrity check of 'inDuplicate' failed
//      TPM_RC_KEY              for an ECC key, the public key is not on the curve
//                              of the curve ID
//      TPM_RC_KEY_SIZE         the decrypted input symmetric key size
//                              does not matches the symmetric algorithm
//                              key size of 'oldParent'
//      TPM_RC_TYPE             'oldParent' is not a storage key, or 'newParent
//                              is not a storage key
//      TPM_RC_VALUE            for an 'oldParent; RSA key, the data to be decrypted
//                              is greater than the public exponent
//      Unmarshal errors        errors during unmarshaling the input
//                              encrypted buffer to a ECC public key, or
//                              unmarshal the private buffer to sensitive
TPM_RC
TPM2_Rewrap(
    Rewrap_In       *in,            // IN: input parameter list
    Rewrap_Out      *out            // OUT: output parameter list
    )
{
    TPM_RC                  result = TPM_RC_SUCCESS;
    TPM2B_DATA              data;               // symmetric key
    UINT16                  hashSize = 0;
    TPM2B_PRIVATE           privateBlob;        // A temporary private blob
                                                // to transit between old
                                                // and new wrappers
// Input Validation
    if((in->inSymSeed.t.size == 0 && in->oldParent != TPM_RH_NULL)
       || (in->inSymSeed.t.size != 0 && in->oldParent == TPM_RH_NULL))
        return TPM_RCS_HANDLE + RC_Rewrap_oldParent;
    if(in->oldParent != TPM_RH_NULL)
    {
        OBJECT              *oldParent = HandleToObject(in->oldParent);

        // old parent key must be a storage object
        if(!ObjectIsStorage(in->oldParent))
            return TPM_RCS_TYPE + RC_Rewrap_oldParent;
        // Decrypt input secret data via asymmetric decryption.  A
        // TPM_RC_VALUE, TPM_RC_KEY or unmarshal errors may be returned at this
        // point
        result = CryptSecretDecrypt(oldParent, NULL, DUPLICATE_STRING, 
                                    &in->inSymSeed, &data);
        if(result != TPM_RC_SUCCESS)
            return TPM_RCS_VALUE + RC_Rewrap_inSymSeed;
        // Unwrap Outer
        result = UnwrapOuter(oldParent, &in->name.b,
                             oldParent->publicArea.nameAlg, &data.b,
                             FALSE,
                             in->inDuplicate.t.size, in->inDuplicate.t.buffer);
        if(result != TPM_RC_SUCCESS)
            return RcSafeAddToResult(result, RC_Rewrap_inDuplicate);
        // Copy unwrapped data to temporary variable, remove the integrity field
        hashSize = sizeof(UINT16) +
            CryptHashGetDigestSize(oldParent->publicArea.nameAlg);
        privateBlob.t.size = in->inDuplicate.t.size - hashSize;
        pAssert(privateBlob.t.size <= sizeof(privateBlob.t.buffer));
        MemoryCopy(privateBlob.t.buffer, in->inDuplicate.t.buffer + hashSize,
                     privateBlob.t.size);
    }
    else
    {
        // No outer wrap from input blob.  Direct copy.
        privateBlob = in->inDuplicate;
    }
    if(in->newParent != TPM_RH_NULL)
    {
        OBJECT          *newParent;
        newParent = HandleToObject(in->newParent);

        // New parent must be a storage object
        if(!ObjectIsStorage(in->newParent))
            return TPM_RCS_TYPE + RC_Rewrap_newParent;
        // Make new encrypt key and its associated secret structure.  A
        // TPM_RC_VALUE error may be returned at this point if RSA algorithm is
        // enabled in TPM
        out->outSymSeed.t.size = sizeof(out->outSymSeed.t.secret);
        result = CryptSecretEncrypt(newParent, DUPLICATE_STRING, &data, 
                                    &out->outSymSeed);
        if(result != TPM_RC_SUCCESS)
            return result;
        // Copy temporary variable to output, reserve the space for integrity
        hashSize = sizeof(UINT16) +
            CryptHashGetDigestSize(newParent->publicArea.nameAlg);
 // Make sure that everything fits into the output buffer
 // Note: this is mostly only an issue if there was no outer wrapper on
 // 'inDuplicate'. It could be as large as a TPM2B_PRIVATE buffer. If we add
 // a digest for an outer wrapper, it won't fit anymore.
        if((privateBlob.t.size + hashSize) > sizeof(out->outDuplicate.t.buffer))
            return TPM_RCS_VALUE + RC_Rewrap_inDuplicate;
// Command output
        out->outDuplicate.t.size = privateBlob.t.size;
        pAssert(privateBlob.t.size 
                <= sizeof(out->outDuplicate.t.buffer) - hashSize);
        MemoryCopy(out->outDuplicate.t.buffer + hashSize, privateBlob.t.buffer, 
                   privateBlob.t.size);
        // Produce outer wrapper for output
        out->outDuplicate.t.size = ProduceOuterWrap(newParent, &in->name.b,
                                                    newParent->publicArea.nameAlg,
                                                    &data.b,
                                                    FALSE,
                                                    out->outDuplicate.t.size,
                                                    out->outDuplicate.t.buffer);
    }
    else  // New parent is a null key so there is no seed
    {
        out->outSymSeed.t.size = 0;

        // Copy privateBlob directly
        out->outDuplicate = privateBlob;
    }
    return TPM_RC_SUCCESS;
}

#endif // CC_Rewrap