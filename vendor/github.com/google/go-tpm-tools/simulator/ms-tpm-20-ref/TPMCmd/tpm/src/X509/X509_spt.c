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
#include "TpmASN1.h"
#include "TpmASN1_fp.h"
#define _X509_SPT_
#include "X509.h"
#include "X509_spt_fp.h"
#if ALG_RSA
#   include "X509_RSA_fp.h"
#endif // ALG_RSA
#if ALG_ECC
#   include "X509_ECC_fp.h"
#endif // ALG_ECC
#if ALG_SM2
//#   include "X509_SM2_fp.h"
#endif // ALG_RSA



//** Unmarshaling Functions

//*** X509FindExtensionOID()
// This will search a list of X508 extensions to find an extension with the
// requested OID. If the extension is found, the output context ('ctx') is set up
// to point to the OID in the extension.
//  Return Type: BOOL
//      TRUE(1)         success
//      FALSE(0)        failure (could be catastrophic)
BOOL
X509FindExtensionByOID(
    ASN1UnmarshalContext    *ctxIn,         // IN: the context to search
    ASN1UnmarshalContext    *ctx,           // OUT: the extension context
    const BYTE              *OID            // IN: oid to search for
)
{
    INT16                length;
//
    pAssert(ctxIn != NULL);
    // Make the search non-destructive of the input if ctx provided. Otherwise, use
    // the provided context.
    if (ctx == NULL)
        ctx = ctxIn;
    else if(ctx != ctxIn)
        *ctx = *ctxIn;
    for(;ctx->size > ctx->offset; ctx->offset += length)
    {
        VERIFY((length = ASN1NextTag(ctx)) >= 0);
        // If this is not a constructed sequence, then it doesn't belong
        // in the extensions.
        VERIFY(ctx->tag == ASN1_CONSTRUCTED_SEQUENCE);
        // Make sure that this entry could hold the OID
        if (length >= OID_SIZE(OID))
        {
            // See if this is a match for the provided object identifier. 
            if (MemoryEqual(OID, &(ctx->buffer[ctx->offset]), OID_SIZE(OID)))
            {
                // Return with ' ctx' set to point to the start of the OID with the size
                // set to be the size of the SEQUENCE
                ctx->buffer += ctx->offset;
                ctx->offset = 0;
                ctx->size = length;
                return TRUE;
            }
        }
    }
    VERIFY(ctx->offset == ctx->size);
    return FALSE;
Error:
    ctxIn->size = -1;
    ctx->size = -1;
    return FALSE;
}

//*** X509GetExtensionBits()
// This function will extract a bit field from an extension. If the extension doesn't
// contain a bit string, it will fail.
// Return Type: BOOL
//  TRUE(1)         success
//  FALSE(0)        failure
UINT32
X509GetExtensionBits(
    ASN1UnmarshalContext            *ctx,
    UINT32                          *value
)
{
    INT16                length;
//
    while (((length = ASN1NextTag(ctx)) > 0) && (ctx->size > ctx->offset))
    {
        // Since this is an extension, the extension value will be in an OCTET STRING
        if (ctx->tag == ASN1_OCTET_STRING)
        {
            return ASN1GetBitStringValue(ctx, value);
        }
        ctx->offset += length;
    }
    ctx->size = -1;
    return FALSE;
}

//***X509ProcessExtensions()
// This function is used to process the TPMA_OBJECT and KeyUsage extensions. It is not
// in the CertifyX509.c code because it makes the code harder to follow.
// Return Type: TPM_RC
//      TPM_RCS_ATTRIBUTES      the attributes of object are not consistent with 
//                              the extension setting
//      TPM_RC_VALUE            problem parsing the extensions
TPM_RC
X509ProcessExtensions(
    OBJECT              *object,        // IN: The object with the attributes to
                                        //      check
    stringRef           *extension      // IN: The start and length of the extensions
)
{
    ASN1UnmarshalContext     ctx;
    ASN1UnmarshalContext     extensionCtx;
    INT16                    length;
    UINT32                   value;
//
    if(!ASN1UnmarshalContextInitialize(&ctx, extension->len, extension->buf)
       || ((length = ASN1NextTag(&ctx)) < 0)
       || (ctx.tag != X509_EXTENSIONS))
        return TPM_RCS_VALUE;
    if( ((length = ASN1NextTag(&ctx)) < 0)
       || (ctx.tag != (ASN1_CONSTRUCTED_SEQUENCE)))
        return TPM_RCS_VALUE;

    // Get the extension for the TPMA_OBJECT if there is one
    if(X509FindExtensionByOID(&ctx, &extensionCtx, OID_TCG_TPMA_OBJECT) &&
        X509GetExtensionBits(&extensionCtx, &value))
    {
        // If an keyAttributes extension was found, it must be exactly the same as the
        // attributes of the object.
        // This cast will work because we know that a TPMA_OBJECT is in a UINT32. 
        // Set RUNTIME_SIZE_CHECKS to YES to force a check to verify this assumption
        // during debug. Doing this is lot easier than having to revisit the code
        // any time a new attribute is added.
        // NOTE: MemoryEqual() is used to avoid type-punned pointer warning/error.
        if(!MemoryEqual(&value, &object->publicArea.objectAttributes, sizeof(value)))
            return TPM_RCS_ATTRIBUTES;
    }
    // Make sure the failure to find the value wasn't because of a fatal error 
    else if(extensionCtx.size < 0)
        return TPM_RCS_VALUE;

    // Get the keyUsage extension. This one is required
    if(X509FindExtensionByOID(&ctx, &extensionCtx, OID_KEY_USAGE_EXTENSTION) &&
        X509GetExtensionBits(&extensionCtx, &value))
    {
        x509KeyUsageUnion   keyUsage;
        TPMA_OBJECT         attributes = object->publicArea.objectAttributes;
    //
        keyUsage.integer = value;
        // For KeyUsage:
        //    the 'sign' attribute is SET if Key Usage includes signing
        if(   (   (keyUsageSign.integer & keyUsage.integer) != 0
               && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign))
           // OR the 'decrypt' attribute is Set if Key Usage includes decryption uses
           || (   (keyUsageDecrypt.integer & keyUsage.integer) != 0
               && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
           // OR that 'fixedTPM' is SET if Key Usage is non-repudiation
           || (   IS_ATTRIBUTE(keyUsage.x509, TPMA_X509_KEY_USAGE, nonrepudiation)
               && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM))
           // OR that 'restricted' is SET if Key Usage is key agreement
           || (   IS_ATTRIBUTE(keyUsage.x509, TPMA_X509_KEY_USAGE, keyAgreement)
               && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted))
           )
            return TPM_RCS_ATTRIBUTES;
    }
    else
        // The KeyUsage extension is required
        return TPM_RCS_VALUE;

    return TPM_RC_SUCCESS;
}

//** Marshaling Functions

//*** X509AddSigningAlgorithm()
// This creates the singing algorithm data.
// Return Type: INT16
//  > 0                 number of octets added
// <= 0                 failure
INT16
X509AddSigningAlgorithm(
    ASN1MarshalContext  *ctx,
    OBJECT              *signKey,
    TPMT_SIG_SCHEME     *scheme
)
{
    switch(signKey->publicArea.type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            return X509AddSigningAlgorithmRSA(signKey, scheme, ctx);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return X509AddSigningAlgorithmECC(signKey, scheme, ctx);
#endif // ALG_ECC
#if ALG_SM2
        case ALG_SM2:
            return X509AddSigningAlgorithmSM2(signKey, scheme,ctx);
#endif // ALG_SM2
        default:
            break;
    }
    return 0;
}

//*** X509AddPublicKey()
// This function will add the publicKey description to the DER data. If fillPtr is 
// NULL, then no data is transferred and this function will indicate if the TPM
// has the values for DER-encoding of the public key.
//  Return Type: INT16
//      > 0         number of octets added
//      == 0        failure
INT16
X509AddPublicKey(
    ASN1MarshalContext  *ctx,
    OBJECT              *object
)
{
    switch(object->publicArea.type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            return X509AddPublicRSA(object, ctx);
#endif
#if ALG_ECC
        case ALG_ECC_VALUE:
            return X509AddPublicECC(object, ctx);
#endif
#if ALG_SM2
        case ALG_SM2_VALUE:
            break;
#endif
        default:
            break;
    }
    return FALSE;
}


//*** X509PushAlgorithmIdentifierSequence()
//  Return Type: INT16
//      > 0         number of bytes added
//     == 0         failure
INT16
X509PushAlgorithmIdentifierSequence(
    ASN1MarshalContext          *ctx,
    const BYTE                  *OID
    )
{
    ASN1StartMarshalContext(ctx);   // hash algorithm
    ASN1PushNull(ctx);
    ASN1PushOID(ctx, OID);
    return ASN1EndEncapsulation(ctx, ASN1_CONSTRUCTED_SEQUENCE);
}


