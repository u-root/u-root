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
#include "TpmASN1_fp.h"
#include "X509_spt_fp.h"
#include "CryptHash_fp.h"
#include "CryptRsa_fp.h"

//** Functions

#if ALG_RSA

//*** X509AddSigningAlgorithmRSA()
// This creates the singing algorithm data.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddSigningAlgorithmRSA(
    OBJECT              *signKey,
    TPMT_SIG_SCHEME     *scheme,
    ASN1MarshalContext  *ctx
)
{
    TPM_ALG_ID           hashAlg = scheme->details.any.hashAlg;
    PHASH_DEF            hashDef = CryptGetHashDef(hashAlg);
//
    NOT_REFERENCED(signKey);
    // return failure if hash isn't implemented
    if(hashDef->hashAlg != hashAlg)
        return 0;
    switch(scheme->scheme)
    {
        case ALG_RSASSA_VALUE:
        {
            // if the hash is implemented but there is no PKCS1 OID defined
            // then this is not a valid signing combination.
            if(hashDef->PKCS1[0] != ASN1_OBJECT_IDENTIFIER)
                break;
            if(ctx == NULL)
                return 1;
            ASN1StartMarshalContext(ctx);
            ASN1PushOID(ctx, hashDef->PKCS1);
            return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
        }
        case ALG_RSAPSS_VALUE:
            // leave if this is just an implementation check
            if(ctx == NULL)
                return 1;
            // In the case of SHA1, everything is default and RFC4055 says that 
            // implementations that do signature generation MUST omit the parameter
            // when defaults are used. )-:
            if(hashDef->hashAlg == ALG_SHA1_VALUE)
            {
                return X509PushAlgorithmIdentifierSequence(ctx, OID_RSAPSS);
            }
            else
            {
                // Going to build something that looks like:
                //  SEQUENCE (2 elem)
                //     OBJECT IDENTIFIER 1.2.840.113549.1.1.10 rsaPSS (PKCS #1)
                //     SEQUENCE (3 elem)
                //       [0] (1 elem)
                //         SEQUENCE (2 elem)
                //           OBJECT IDENTIFIER 2.16.840.1.101.3.4.2.1 sha-256 
                //           NULL
                //       [1] (1 elem)
                //         SEQUENCE (2 elem)
                //           OBJECT IDENTIFIER 1.2.840.113549.1.1.8 pkcs1-MGF
                //           SEQUENCE (2 elem)
                //             OBJECT IDENTIFIER 2.16.840.1.101.3.4.2.1 sha-256
                //             NULL
                //       [2] (1 elem)  salt length
                //         INTEGER 32

                // The indentation is just to keep track of where we are in the 
                // structure
                ASN1StartMarshalContext(ctx); // SEQUENCE (2 elements)
                {
                    ASN1StartMarshalContext(ctx);   // SEQUENCE (3 elements)
                    {
                        // [2] (1 elem)  salt length
                        //    INTEGER 32
                        ASN1StartMarshalContext(ctx);
                        {
                            INT16       saltSize =
                                CryptRsaPssSaltSize((INT16)hashDef->digestSize,
                                (INT16)signKey->publicArea.unique.rsa.t.size);
                            ASN1PushUINT(ctx, saltSize);
                        }
                        ASN1EndEncapsulation(ctx, ASN1_APPLICAIION_SPECIFIC + 2);

                        // Add the mask generation algorithm
                        // [1] (1 elem)
                        //    SEQUENCE (2 elem) 1st
                        //      OBJECT IDENTIFIER 1.2.840.113549.1.1.8 pkcs1-MGF
                        //      SEQUENCE (2 elem) 2nd  
                        //        OBJECT IDENTIFIER 2.16.840.1.101.3.4.2.1 sha-256
                        //        NULL
                        ASN1StartMarshalContext(ctx);   // mask context [1] (1 elem)
                        {
                            ASN1StartMarshalContext(ctx);   // SEQUENCE (2 elem) 1st
                            // Handle the 2nd Sequence (sequence (object, null))
                            {
                                X509PushAlgorithmIdentifierSequence(ctx,
                                    hashDef->OID);
                                // add the pkcs1-MGF OID
                                ASN1PushOID(ctx, OID_MGF1);
                            }
                            // End outer sequence
                            ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
                        }
                        // End the [1] 
                        ASN1EndEncapsulation(ctx, ASN1_APPLICAIION_SPECIFIC + 1);

                        // Add the hash algorithm
                        // [0] (1 elem)
                        //   SEQUENCE (2 elem) (done by 
                        //              X509PushAlgorithmIdentifierSequence)
                        //     OBJECT IDENTIFIER 2.16.840.1.101.3.4.2.1 sha-256 (NIST)
                        //     NULL
                        ASN1StartMarshalContext(ctx); // [0] (1 elem)
                        {
                            X509PushAlgorithmIdentifierSequence(ctx, hashDef->OID);
                        }
                        ASN1EndEncapsulation(ctx, (ASN1_APPLICAIION_SPECIFIC + 0));
                    }
                    //  SEQUENCE (3 elements) end
                    ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);

                    // RSA PSS OID
                    // OBJECT IDENTIFIER 1.2.840.113549.1.1.10 rsaPSS (PKCS #1)
                    ASN1PushOID(ctx, OID_RSAPSS);
                }
                // End Sequence (2 elements)
                return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
            }
        default:
            break;
    }
    return 0;
}

//*** X509AddPublicRSA()
// This function will add the publicKey description to the DER data. If fillPtr is 
// NULL, then no data is transferred and this function will indicate if the TPM
// has the values for DER-encoding of the public key.
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509AddPublicRSA(
    OBJECT                  *object,
    ASN1MarshalContext    *ctx
)
{
    UINT32          exp = object->publicArea.parameters.rsaDetail.exponent;
//
/*
    SEQUENCE (2 elem) 1st
      SEQUENCE (2 elem) 2nd
        OBJECT IDENTIFIER 1.2.840.113549.1.1.1 rsaEncryption (PKCS #1)
        NULL
      BIT STRING (1 elem)
        SEQUENCE (2 elem) 3rd
          INTEGER (2048 bit) 2197304513741227955725834199357401…
          INTEGER 65537
*/
    // If this is a check to see if the key can be encoded, it can. 
    // Need to mark the end sequence
    if(ctx == NULL)
        return 1;
    ASN1StartMarshalContext(ctx); // SEQUENCE (2 elem) 1st
    ASN1StartMarshalContext(ctx); // BIT STRING
    ASN1StartMarshalContext(ctx); // SEQUENCE *(2 elem) 3rd

    // Get public exponent in big-endian byte order.
    if(exp == 0)
        exp = RSA_DEFAULT_PUBLIC_EXPONENT;

    // Push a 4 byte integer. This might get reduced if there are leading zeros or
    // extended if the high order byte is negative.
    ASN1PushUINT(ctx, exp);
    // Push the public key as an integer
    ASN1PushInteger(ctx, object->publicArea.unique.rsa.t.size,
                             object->publicArea.unique.rsa.t.buffer);
    // Embed this in a SEQUENCE tag and length in for the key, exponent sequence
    ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE); // SEQUENCE (3rd)
                                 
    // Embed this in a BIT STRING
    ASN1EndEncapsulation(ctx, ASN1_BITSTRING);

    // Now add the formatted SEQUENCE for the RSA public key OID. This is a
    // fully constructed value so it doesn't need to have a context started
    X509PushAlgorithmIdentifierSequence(ctx, OID_PKCS1_PUB);

    return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
}

#endif // ALG_RSA