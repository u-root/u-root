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
#include "X509.h"
#include "OIDs.h"
#include "TpmASN1_fp.h"
#include "X509_spt_fp.h"
#include "CryptHash_fp.h"

//** Functions

//*** X509PushPoint()
// This seems like it might be used more than once so...
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509PushPoint(
    ASN1MarshalContext      *ctx,
    TPMS_ECC_POINT          *p
)
{
    // Push a bit string containing the public key. For now, push the x, and y
    // coordinates of the public point, bottom up
    ASN1StartMarshalContext(ctx); // BIT STRING
    {
        ASN1PushBytes(ctx, p->y.t.size, p->y.t.buffer);
        ASN1PushBytes(ctx, p->x.t.size, p->x.t.buffer);
        ASN1PushByte(ctx, 0x04);
    }
    return ASN1EndEncapsulation(ctx, ASN1_BITSTRING); // Ends BIT STRING
}

//*** X509AddSigningAlgorithmECC()
// This creates the singing algorithm data.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddSigningAlgorithmECC(
    OBJECT              *signKey,
    TPMT_SIG_SCHEME     *scheme,
    ASN1MarshalContext  *ctx
)
{
    PHASH_DEF            hashDef = CryptGetHashDef(scheme->details.any.hashAlg);
//
    NOT_REFERENCED(signKey);
    // If the desired hashAlg definition wasn't found...
    if(hashDef->hashAlg != scheme->details.any.hashAlg)
        return 0;

    switch(scheme->scheme)
    {
        case ALG_ECDSA_VALUE:
            // Make sure that we have an OID for this hash and ECC
            if((hashDef->ECDSA)[0] != ASN1_OBJECT_IDENTIFIER)
                break;
            // if this is just an implementation check, indicate that this 
            // combination is supported
            if(!ctx)
                return 1;
            ASN1StartMarshalContext(ctx);
            ASN1PushOID(ctx, hashDef->ECDSA);
            return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
        default:
            break;
    }
    return 0;
}


//*** X509AddPublicECC()
// This function will add the publicKey description to the DER data. If ctx is 
// NULL, then no data is transferred and this function will indicate if the TPM
// has the values for DER-encoding of the public key.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddPublicECC(
    OBJECT                *object,
    ASN1MarshalContext    *ctx
)
{
    const BYTE      *curveOid =
        CryptEccGetOID(object->publicArea.parameters.eccDetail.curveID);
    if((curveOid == NULL) || (*curveOid != ASN1_OBJECT_IDENTIFIER))
        return 0;
//
//
//  SEQUENCE (2 elem) 1st
//    SEQUENCE (2 elem) 2nd
//      OBJECT IDENTIFIER 1.2.840.10045.2.1 ecPublicKey (ANSI X9.62 public key type)
//      OBJECT IDENTIFIER 1.2.840.10045.3.1.7 prime256v1 (ANSI X9.62 named curve)
//    BIT STRING (520 bit) 000001001010000111010101010111001001101101000100000010...
//
    // If this is a check to see if the key can be encoded, it can. 
    // Need to mark the end sequence
    if(ctx == NULL)
        return 1;
    ASN1StartMarshalContext(ctx); // SEQUENCE (2 elem) 1st
    {
        X509PushPoint(ctx, &object->publicArea.unique.ecc); // BIT STRING 
        ASN1StartMarshalContext(ctx); // SEQUENCE (2 elem) 2nd
        {
            ASN1PushOID(ctx, curveOid); // curve dependent
            ASN1PushOID(ctx, OID_ECC_PUBLIC); // (1.2.840.10045.2.1)
        }
        ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE); // Ends SEQUENCE 2nd
    }
    return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE); // Ends SEQUENCE 1st
}