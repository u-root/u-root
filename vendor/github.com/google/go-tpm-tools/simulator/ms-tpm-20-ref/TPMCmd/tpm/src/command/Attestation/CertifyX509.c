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
#include "CertifyX509_fp.h"
#include "X509.h"
#include "TpmASN1_fp.h"
#include "X509_spt_fp.h"
#include "Attest_spt_fp.h"

#if CC_CertifyX509 // Conditional expansion of this file

/*(See part 3 specification)
// Certify 
*/
// return type: TPM_RC
//      TPM_RC_ATTRIBUTES       the attributes of 'objectHandle' are not compatible
//                              with the KeyUsage or TPMA_OBJECT values in the 
//                              extensions fields
//      TPM_RC_BINDING          the public and private portions of the key are not
//                              properly bound.
//      TPM_RC_HASH             the hash algorithm in the scheme is not supported
//      TPM_RC_KEY              'signHandle' does not reference a signing key;
//      TPM_RC_SCHEME           the scheme is not compatible with sign key type,
//                              or input scheme is not compatible with default
//                              scheme, or the chosen scheme is not a valid
//                              sign scheme
//      TPM_RC_VALUE            most likely a problem with the format of 
//                              'partialCertificate'
TPM_RC
TPM2_CertifyX509(
    CertifyX509_In          *in,          // IN: input parameter list
    CertifyX509_Out         *out            // OUT: output parameter list
)
{
    TPM_RC                   result;
    OBJECT                  *signKey = HandleToObject(in->signHandle);
    OBJECT                  *object = HandleToObject(in->objectHandle);
    HASH_STATE               hash;
    INT16                    length;        // length for a tagged element
    ASN1UnmarshalContext     ctx;
    ASN1MarshalContext       ctxOut;
    // certTBS holds an array of pointers and lengths. Each entry references the 
    // corresponding value in a TBSCertificate structure. For example, the 1th 
    // element references the version number
    stringRef                certTBS[REF_COUNT] = {{0}};
#define ALLOWED_SEQUENCES   (SUBJECT_PUBLIC_KEY_REF - SIGNATURE_REF) 
    stringRef                partial[ALLOWED_SEQUENCES] = {{0}};
    INT16                    countOfSequences = 0;
    INT16                    i;
    //
#if CERTIFYX509_DEBUG
    DebugFileOpen();
    DebugDumpBuffer(in->partialCertificate.t.size, in->partialCertificate.t.buffer,
        "partialCertificate");
#endif

    // Input Validation
    // signing key must be able to sign
    if(!IsSigningObject(signKey))
        return TPM_RCS_KEY + RC_CertifyX509_signHandle;
    // Pick a scheme for sign.  If the input sign scheme is not compatible with
    // the default scheme, return an error.
    if(!CryptSelectSignScheme(signKey, &in->inScheme))
        return TPM_RCS_SCHEME + RC_CertifyX509_inScheme;
    // Make sure that the public Key encoding is known
    if(X509AddPublicKey(NULL, object) == 0)
        return TPM_RCS_ASYMMETRIC + RC_CertifyX509_objectHandle;
    // Unbundle 'partialCertificate'.
        // Initialize the unmarshaling context
    if(!ASN1UnmarshalContextInitialize(&ctx, in->partialCertificate.t.size,
        in->partialCertificate.t.buffer))
        return TPM_RCS_VALUE + RC_CertifyX509_partialCertificate;
    // Make sure that this is a constructed SEQUENCE
    length = ASN1NextTag(&ctx);
    // Must be a constructed SEQUENCE that uses all of the input parameter
    if((ctx.tag != (ASN1_CONSTRUCTED_SEQUENCE))
        || ((ctx.offset + length) != in->partialCertificate.t.size))
        return TPM_RCS_SIZE + RC_CertifyX509_partialCertificate;

    // This scans through the contents of the outermost SEQUENCE. This would be the
    // 'issuer', 'validity', 'subject', 'issuerUniqueID' (optional), 
    // 'subjectUniqueID' (optional), and 'extensions.'
    while(ctx.offset < ctx.size)
    {
        INT16           startOfElement = ctx.offset;
        //
            // Read the next tag and length field. 
        length = ASN1NextTag(&ctx);
        if(length < 0)
            break;
        if(ctx.tag == ASN1_CONSTRUCTED_SEQUENCE)
        {
            partial[countOfSequences].buf = &ctx.buffer[startOfElement];
            ctx.offset += length;
            partial[countOfSequences].len = (INT16)ctx.offset - startOfElement;
            if(++countOfSequences > ALLOWED_SEQUENCES)
                break;
        }
        else if(ctx.tag  == X509_EXTENSIONS)
        {
            if(certTBS[EXTENSIONS_REF].len != 0)
                return TPM_RCS_VALUE + RC_CertifyX509_partialCertificate;
            certTBS[EXTENSIONS_REF].buf = &ctx.buffer[startOfElement];
            ctx.offset += length;
            certTBS[EXTENSIONS_REF].len =
                (INT16)ctx.offset - startOfElement;
        }
        else
            return TPM_RCS_VALUE + RC_CertifyX509_partialCertificate;
    }
    // Make sure that we used all of the data and found at least the required
    // number of elements. 
    if((ctx.offset != ctx.size) || (countOfSequences < 3)
        || (countOfSequences > 4)
        || (certTBS[EXTENSIONS_REF].buf == NULL))
        return TPM_RCS_VALUE + RC_CertifyX509_partialCertificate;
    // Now that we know how many sequences there were, we can put them where they
    // belong
    for(i = 0; i < countOfSequences; i++)
        certTBS[SUBJECT_KEY_REF - i] = partial[countOfSequences - 1 - i];

    // If only three SEQUENCES, then the TPM needs to produce the signature algorithm.
    // See if it can
    if((countOfSequences == 3) && 
        (X509AddSigningAlgorithm(NULL, signKey, &in->inScheme) == 0))
            return TPM_RCS_SCHEME + RC_CertifyX509_signHandle;

    // Process the extensions
    result = X509ProcessExtensions(object, &certTBS[EXTENSIONS_REF]);
    if(result != TPM_RC_SUCCESS)
        // If the extension has the TPMA_OBJECT extension and the attributes don't 
        // match, then the error code will be TPM_RCS_ATTRIBUTES. Otherwise, the error
        // indicates a malformed partialCertificate.
        return result + ((result == TPM_RCS_ATTRIBUTES)
                         ? RC_CertifyX509_objectHandle
                         : RC_CertifyX509_partialCertificate);
// Command Output
// Create the addedToCertificate values

    // Build the addedToCertificate from the bottom up.
    // Initialize the context structure
    ASN1InitialializeMarshalContext(&ctxOut, sizeof(out->addedToCertificate.t.buffer),
                                    out->addedToCertificate.t.buffer);
    // Place a marker for the overall context
    ASN1StartMarshalContext(&ctxOut);  // SEQUENCE for addedToCertificate

    // Add the subject public key descriptor
    certTBS[SUBJECT_PUBLIC_KEY_REF].len = X509AddPublicKey(&ctxOut, object);
    certTBS[SUBJECT_PUBLIC_KEY_REF].buf = ctxOut.buffer + ctxOut.offset;
    // If the caller didn't provide the algorithm identifier, create it
    if(certTBS[SIGNATURE_REF].len == 0)
    {
        certTBS[SIGNATURE_REF].len = X509AddSigningAlgorithm(&ctxOut, signKey,
            &in->inScheme);
        certTBS[SIGNATURE_REF].buf = ctxOut.buffer + ctxOut.offset;
    }
    // Create the serial number value. Use the out->tbsDigest as scratch.
    {
        TPM2B                   *digest = &out->tbsDigest.b;
        //
        digest->size = (INT16)CryptHashStart(&hash, signKey->publicArea.nameAlg);
        pAssert(digest->size != 0);

        // The serial number size is the smaller of the digest and the vendor-defined
        // value
        digest->size = MIN(digest->size, SIZE_OF_X509_SERIAL_NUMBER);
        // Add all the parts of the certificate other than the serial number 
        // and version number
        for(i = SIGNATURE_REF; i < REF_COUNT; i++)
            CryptDigestUpdate(&hash, certTBS[i].len, certTBS[i].buf);
        // throw in the Name of the signing key...
        CryptDigestUpdate2B(&hash, &signKey->name.b);
        // ...and the Name of the signed key.
        CryptDigestUpdate2B(&hash, &object->name.b);
        // Done
        CryptHashEnd2B(&hash, digest);
    }

    // Add the serial number
    certTBS[SERIAL_NUMBER_REF].len = 
        ASN1PushInteger(&ctxOut, out->tbsDigest.t.size, out->tbsDigest.t.buffer);
    certTBS[SERIAL_NUMBER_REF].buf = ctxOut.buffer + ctxOut.offset;

    // Add the static version number
    ASN1StartMarshalContext(&ctxOut);
    ASN1PushUINT(&ctxOut, 2);
    certTBS[VERSION_REF].len = 
        ASN1EndEncapsulation(&ctxOut, ASN1_APPLICAIION_SPECIFIC);
    certTBS[VERSION_REF].buf = ctxOut.buffer + ctxOut.offset;

    // Create a fake tag and length for the TBS in the space used for 
    // 'addedToCertificate'
    {
        for(length = 0, i = 0; i < REF_COUNT; i++) 
            length += certTBS[i].len;
        // Put a fake tag and length into the buffer for use in the tbsDigest
        certTBS[ENCODED_SIZE_REF].len =
            ASN1PushTagAndLength(&ctxOut, ASN1_CONSTRUCTED_SEQUENCE, length);
        certTBS[ENCODED_SIZE_REF].buf = ctxOut.buffer + ctxOut.offset;
        // Restore the buffer pointer to add back the number of octets used for the
        // tag and length
        ctxOut.offset += certTBS[ENCODED_SIZE_REF].len;
    }
    // sanity check
    if(ctxOut.offset < 0)
        return TPM_RC_FAILURE;
    // Create the tbsDigest to sign
    out->tbsDigest.t.size = CryptHashStart(&hash, in->inScheme.details.any.hashAlg);
    for(i = 0; i < REF_COUNT; i++)
        CryptDigestUpdate(&hash, certTBS[i].len, certTBS[i].buf);
    CryptHashEnd2B(&hash, &out->tbsDigest.b);

#if CERTIFYX509_DEBUG
    {
        BYTE                 fullTBS[4096];
        BYTE                *fill = fullTBS;
        int         		 j;
        for (j = 0; j < REF_COUNT; j++)
        {
            MemoryCopy(fill, certTBS[j].buf, certTBS[j].len);
            fill += certTBS[j].len;
        }
        DebugDumpBuffer((int)(fill - &fullTBS[0]), fullTBS, "\nfull TBS");
    }
#endif

// Finish up the processing of addedToCertificate
    // Create the actual tag and length for the addedToCertificate structure
    out->addedToCertificate.t.size =
        ASN1EndEncapsulation(&ctxOut, ASN1_CONSTRUCTED_SEQUENCE);
    // Now move all the addedToContext to the start of the buffer
    MemoryCopy(out->addedToCertificate.t.buffer, ctxOut.buffer + ctxOut.offset,
               out->addedToCertificate.t.size);
#if CERTIFYX509_DEBUG
    DebugDumpBuffer(out->addedToCertificate.t.size, out->addedToCertificate.t.buffer,
                    "\naddedToCertificate");
#endif
    // only thing missing is the signature
    result = CryptSign(signKey, &in->inScheme, &out->tbsDigest, &out->signature);

    return result;
}

#endif // CC_CertifyX509
