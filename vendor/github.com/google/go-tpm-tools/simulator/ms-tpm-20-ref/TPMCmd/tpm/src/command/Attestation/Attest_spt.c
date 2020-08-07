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
//** Includes
#include "Tpm.h"
#include "Attest_spt_fp.h"

//** Functions

//***FillInAttestInfo()
// Fill in common fields of TPMS_ATTEST structure.
void
FillInAttestInfo(
    TPMI_DH_OBJECT       signHandle,    // IN: handle of signing object
    TPMT_SIG_SCHEME     *scheme,        // IN/OUT: scheme to be used for signing
    TPM2B_DATA          *data,          // IN: qualifying data
    TPMS_ATTEST         *attest         // OUT: attest structure
    )
{
    OBJECT              *signObject = HandleToObject(signHandle);

        // Magic number
    attest->magic = TPM_GENERATED_VALUE;

    if(signObject == NULL)
    {
        // The name for a null handle is TPM_RH_NULL
        // This is defined because UINT32_TO_BYTE_ARRAY does a cast. If the
        // size of the cast is smaller than a constant, the compiler warns
        // about the truncation of a constant value.
        TPM_HANDLE      nullHandle = TPM_RH_NULL;
        attest->qualifiedSigner.t.size = sizeof(TPM_HANDLE);
        UINT32_TO_BYTE_ARRAY(nullHandle, attest->qualifiedSigner.t.name);
    }
    else
    {
        // Certifying object qualified name
        // if the scheme is anonymous, this is an empty buffer
        if(CryptIsSchemeAnonymous(scheme->scheme))
            attest->qualifiedSigner.t.size = 0;
        else
            attest->qualifiedSigner = signObject->qualifiedName;
    }
    // current clock in plain text
    TimeFillInfo(&attest->clockInfo);

    // Firmware version in plain text
    attest->firmwareVersion = ((UINT64)gp.firmwareV1 << (sizeof(UINT32) * 8));
    attest->firmwareVersion += gp.firmwareV2;

    // Check the hierarchy of sign object.  For NULL sign handle, the hierarchy
    // will be TPM_RH_NULL
    if((signObject == NULL)
       || (!signObject->attributes.epsHierarchy
           && !signObject->attributes.ppsHierarchy))
    {
        // For signing key that is not in platform or endorsement hierarchy,
        // obfuscate the reset, restart and firmware version information
        UINT64          obfuscation[2];
        CryptKDFa(CONTEXT_INTEGRITY_HASH_ALG, &gp.shProof.b, OBFUSCATE_STRING,
                  &attest->qualifiedSigner.b, NULL, 128,
                  (BYTE *)&obfuscation[0], NULL, FALSE);
        // Obfuscate data
        attest->firmwareVersion += obfuscation[0];
        attest->clockInfo.resetCount += (UINT32)(obfuscation[1] >> 32);
        attest->clockInfo.restartCount += (UINT32)obfuscation[1];
    }
    // External data
    if(CryptIsSchemeAnonymous(scheme->scheme))
        attest->extraData.t.size = 0;
    else
    {
        // If we move the data to the attestation structure, then it is not
        // used in the signing operation except as part of the signed data
        attest->extraData = *data;
        data->t.size = 0;
    }
}

//***SignAttestInfo()
// Sign a TPMS_ATTEST structure. If signHandle is TPM_RH_NULL, a null signature
// is returned.
//
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   'signHandle' references not a signing key
//      TPM_RC_SCHEME       'scheme' is not compatible with 'signHandle' type
//      TPM_RC_VALUE        digest generated for the given 'scheme' is greater than
//                          the modulus of 'signHandle' (for an RSA key);
//                          invalid commit status or failed to generate "r" value
//                          (for an ECC key)
TPM_RC
SignAttestInfo(
    OBJECT              *signKey,           // IN: sign object
    TPMT_SIG_SCHEME     *scheme,            // IN: sign scheme
    TPMS_ATTEST         *certifyInfo,       // IN: the data to be signed
    TPM2B_DATA          *qualifyingData,    // IN: extra data for the signing
                                            //     process
    TPM2B_ATTEST        *attest,            // OUT: marshaled attest blob to be
                                            //     signed
    TPMT_SIGNATURE      *signature          // OUT: signature
    )
{
    BYTE                    *buffer;
    HASH_STATE              hashState;
    TPM2B_DIGEST            digest;
    TPM_RC                  result;

    // Marshal TPMS_ATTEST structure for hash
    buffer = attest->t.attestationData;
    attest->t.size = TPMS_ATTEST_Marshal(certifyInfo, &buffer, NULL);

    if(signKey == NULL)
    {
        signature->sigAlg = TPM_ALG_NULL;
        result = TPM_RC_SUCCESS;
    }
    else
    {
        TPMI_ALG_HASH           hashAlg;
        // Compute hash
        hashAlg = scheme->details.any.hashAlg;
        // need to set the receive buffer to get something put in it
        digest.t.size = sizeof(digest.t.buffer);
        digest.t.size = CryptHashBlock(hashAlg, attest->t.size,
                                       attest->t.attestationData,
                                       digest.t.size, digest.t.buffer);
        // If there is qualifying data, need to rehash the data
        // hash(qualifyingData || hash(attestationData))
        if(qualifyingData->t.size != 0)
        {
            CryptHashStart(&hashState, hashAlg);
            CryptDigestUpdate2B(&hashState, &qualifyingData->b);
            CryptDigestUpdate2B(&hashState, &digest.b);
            CryptHashEnd2B(&hashState, &digest.b);
        }
        // Sign the hash. A TPM_RC_VALUE, TPM_RC_SCHEME, or
        // TPM_RC_ATTRIBUTES error may be returned at this point
        result = CryptSign(signKey, scheme, &digest, signature);

        // Since the clock is used in an attestation, the state in NV is no longer
        // "orderly" with respect to the data in RAM if the signature is valid
        if(result == TPM_RC_SUCCESS)
        {
            // Command uses the clock so need to clear the orderly state if it is
            // set.
            result = NvClearOrderly();
        }
    }
    return result;
}

//*** IsSigningObject()
// Checks to see if the object is OK for signing. This is here rather than in
// Object_spt.c because all the attestation commands use this file but not
// Object_spt.c.
//  Return Type: BOOL
//      TRUE(1)         object may sign
//      FALSE(0)        object may not sign
BOOL
IsSigningObject(
    OBJECT          *object         // IN:
    )
{
    return ((object == NULL) 
            || ((IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, sign)
               && object->publicArea.type != TPM_ALG_SYMCIPHER)));
}