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
#include "Object_spt_fp.h"

//** Local Functions

//*** GetIV2BSize()
// Get the size of TPM2B_IV in canonical form that will be append to the start of
// the sensitive data.  It includes both size of size field and size of iv data
static UINT16
GetIV2BSize(
    OBJECT              *protector          // IN: the protector handle
    )
{
    TPM_ALG_ID          symAlg;
    UINT16              keyBits;

    // Determine the symmetric algorithm and size of key
    if(protector == NULL)
    {
        // Use the context encryption algorithm and key size
        symAlg = CONTEXT_ENCRYPT_ALG;
        keyBits = CONTEXT_ENCRYPT_KEY_BITS;
    }
    else
    {
        symAlg = protector->publicArea.parameters.asymDetail.symmetric.algorithm;
        keyBits = protector->publicArea.parameters.asymDetail.symmetric.keyBits.sym;
    }

    // The IV size is a UINT16 size field plus the block size of the symmetric
    // algorithm
    return sizeof(UINT16) + CryptGetSymmetricBlockSize(symAlg, keyBits);
}

//*** ComputeProtectionKeyParms()
// This function retrieves the symmetric protection key parameters for
// the sensitive data
// The parameters retrieved from this function include encryption algorithm,
// key size in bit, and a TPM2B_SYM_KEY containing the key material as well as
// the key size in bytes
// This function is used for any action that requires encrypting or decrypting of
// the sensitive area of an object or a credential blob
//
/*(See part 1 specification)
    KDF for generating the protection key material:
    KDFa(hashAlg, seed, "STORAGE", Name, NULL , bits)
where
    hashAlg     for a Primary Object, an algorithm chosen by the TPM vendor
                for derivations from Primary Seeds. For all other objects,
                the nameAlg of the object's parent.
    seed        for a Primary Object in the Platform Hierarchy, the PPS.
                For Primary Objects in either Storage or Endorsement Hierarchy,
                the SPS. For Temporary Objects, the context encryption seed.
                For all other objects, the symmetric seed value in the
                sensitive area of the object's parent.
    STORAGE     label to differentiate use of KDFa() (see 4.7)
    Name        the Name of the object being encrypted
    bits        the number of bits required for a  symmetric key and IV
*/
//  Return Type: void
static void
ComputeProtectionKeyParms(
    OBJECT          *protector,         // IN: the protector object
    TPM_ALG_ID       hashAlg,           // IN: hash algorithm for KDFa
    TPM2B           *name,              // IN: name of the object
    TPM2B           *seedIn,            // IN: optional seed for duplication blob.
                                        //     For non duplication blob, this
                                        //     parameter should be NULL
    TPM_ALG_ID      *symAlg,            // OUT: the symmetric algorithm
    UINT16          *keyBits,           // OUT: the symmetric key size in bits
    TPM2B_SYM_KEY   *symKey             // OUT: the symmetric key
    )
{
    const TPM2B         *seed = seedIn;

    // Determine the algorithms for the KDF and the encryption/decryption
    // For TPM_RH_NULL, using context settings
    if(protector == NULL)
    {
        // Use the context encryption algorithm and key size
        *symAlg = CONTEXT_ENCRYPT_ALG;
        symKey->t.size = CONTEXT_ENCRYPT_KEY_BYTES;
        *keyBits = CONTEXT_ENCRYPT_KEY_BITS;
    }
    else
    {
        TPMT_SYM_DEF_OBJECT *symDef;
        symDef = &protector->publicArea.parameters.asymDetail.symmetric;
        *symAlg = symDef->algorithm;
        *keyBits = symDef->keyBits.sym;
        symKey->t.size = (*keyBits + 7) / 8;
    }
    // Get seed for KDF
    if(seed == NULL)
        seed = GetSeedForKDF(protector);
    // KDFa to generate symmetric key and IV value
    CryptKDFa(hashAlg, seed, STORAGE_KEY, name, NULL,
              symKey->t.size * 8, symKey->t.buffer, NULL, FALSE);
    return;
}

//*** ComputeOuterIntegrity()
// The sensitive area parameter is a buffer that holds a space for
// the integrity value and the marshaled sensitive area. The caller should
// skip over the area set aside for the integrity value
// and compute the hash of the remainder of the object.
// The size field of sensitive is in unmarshaled form and the
// sensitive area contents is an array of bytes.
/*(See part 1 specification)
    KDFa(hashAlg, seed, "INTEGRITY", NULL, NULL , bits)   (38)
where
    hashAlg     for a Primary Object, the nameAlg of the object. For all other
                objects the nameAlg of the object's parent.
    seed        for a Primary Object in the Platform Hierarchy, the PPS. For
                Primary Objects in either Storage or Endorsement Hierarchy,
                the SPS. For a Temporary Object, the context encryption key.
                For all other objects, the symmetric seed value in the sensitive
                area of the object's parent.
    "INTEGRITY" a value used to differentiate the uses of the KDF.
    bits        the number of bits in the digest produced by hashAlg.
Key is then used in the integrity computation.
    HMACnameAlg(HMACkey, encSensitive || Name )
where
    HMACnameAlg()   the HMAC function using nameAlg of the object's parent
    HMACkey         value derived from the parent symmetric protection value
    encSensitive    symmetrically encrypted sensitive area
    Name            the Name of the object being protected
*/
//  Return Type: void
static void
ComputeOuterIntegrity(
    TPM2B           *name,              // IN: the name of the object
    OBJECT          *protector,         // IN: the object that
                                        //     provides protection. For an object,
                                        //     it is a parent. For a credential, it
                                        //     is the encrypt object. For
                                        //     a Temporary Object, it is NULL
    TPMI_ALG_HASH    hashAlg,           // IN: algorithm to use for integrity
    TPM2B           *seedIn,            // IN: an external seed may be provided for
                                        //     duplication blob. For non duplication
                                        //     blob, this parameter should be NULL
    UINT32           sensitiveSize,     // IN: size of the marshaled sensitive data
    BYTE            *sensitiveData,     // IN: sensitive area
    TPM2B_DIGEST    *integrity          // OUT: integrity
    )
{
    HMAC_STATE       hmacState;
    TPM2B_DIGEST     hmacKey;
    const TPM2B     *seed = seedIn;
//
    // Get seed for KDF
    if(seed == NULL)
        seed = GetSeedForKDF(protector);
    // Determine the HMAC key bits
    hmacKey.t.size = CryptHashGetDigestSize(hashAlg);

    // KDFa to generate HMAC key
    CryptKDFa(hashAlg, seed, INTEGRITY_KEY, NULL, NULL,
              hmacKey.t.size * 8, hmacKey.t.buffer, NULL, FALSE);
    // Start HMAC and get the size of the digest which will become the integrity
    integrity->t.size = CryptHmacStart2B(&hmacState, hashAlg, &hmacKey.b);

    // Adding the marshaled sensitive area to the integrity value
    CryptDigestUpdate(&hmacState.hashState, sensitiveSize, sensitiveData);

    // Adding name
    CryptDigestUpdate2B(&hmacState.hashState, name);

    // Compute HMAC
    CryptHmacEnd2B(&hmacState, &integrity->b);

    return;
}

//*** ComputeInnerIntegrity()
// This function computes the integrity of an inner wrap
static void
ComputeInnerIntegrity(
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm for inner wrap
    TPM2B           *name,          // IN: the name of the object
    UINT16           dataSize,      // IN: the size of sensitive data
    BYTE            *sensitiveData, // IN: sensitive data
    TPM2B_DIGEST    *integrity      // OUT: inner integrity
    )
{
    HASH_STATE      hashState;
//
    // Start hash and get the size of the digest which will become the integrity
    integrity->t.size = CryptHashStart(&hashState, hashAlg);

    // Adding the marshaled sensitive area to the integrity value
    CryptDigestUpdate(&hashState, dataSize, sensitiveData);

    // Adding name
    CryptDigestUpdate2B(&hashState, name);

    // Compute hash
    CryptHashEnd2B(&hashState, &integrity->b);

    return;
}

//*** ProduceInnerIntegrity()
// This function produces an inner integrity for regular private, credential or
// duplication blob
// It requires the sensitive data being marshaled to the innerBuffer, with the
// leading bytes reserved for integrity hash.  It assume the sensitive data
// starts at address (innerBuffer + integrity size).
// This function integrity at the beginning of the inner buffer
// It returns the total size of buffer with the inner wrap
static UINT16
ProduceInnerIntegrity(
    TPM2B           *name,          // IN: the name of the object
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm for inner wrap
    UINT16           dataSize,      // IN: the size of sensitive data, excluding the
                                    //     leading integrity buffer size
    BYTE            *innerBuffer    // IN/OUT: inner buffer with sensitive data in
                                    //     it.  At input, the leading bytes of this
                                    //     buffer is reserved for integrity
    )
{
    BYTE            *sensitiveData; // pointer to the sensitive data
    TPM2B_DIGEST     integrity;
    UINT16           integritySize;
    BYTE            *buffer;        // Auxiliary buffer pointer
//
    // sensitiveData points to the beginning of sensitive data in innerBuffer
    integritySize = sizeof(UINT16) + CryptHashGetDigestSize(hashAlg);
    sensitiveData = innerBuffer + integritySize;

    ComputeInnerIntegrity(hashAlg, name, dataSize, sensitiveData, &integrity);

    // Add integrity at the beginning of inner buffer
    buffer = innerBuffer;
    TPM2B_DIGEST_Marshal(&integrity, &buffer, NULL);

    return dataSize + integritySize;
}

//*** CheckInnerIntegrity()
// This function check integrity of inner blob
//  Return Type: TPM_RC
//      TPM_RC_INTEGRITY        if the outer blob integrity is bad
//      unmarshal errors        unmarshal errors while unmarshaling integrity
static TPM_RC
CheckInnerIntegrity(
    TPM2B           *name,          // IN: the name of the object
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm for inner wrap
    UINT16           dataSize,      // IN: the size of sensitive data, including the
                                    //     leading integrity buffer size
    BYTE            *innerBuffer    // IN/OUT: inner buffer with sensitive data in
                                    //     it
    )
{
    TPM_RC          result;
    TPM2B_DIGEST    integrity;
    TPM2B_DIGEST    integrityToCompare;
    BYTE            *buffer;                // Auxiliary buffer pointer
    INT32           size;
//
    // Unmarshal integrity
    buffer = innerBuffer;
    size = (INT32)dataSize;
    result = TPM2B_DIGEST_Unmarshal(&integrity, &buffer, &size);
    if(result == TPM_RC_SUCCESS)
    {
        // Compute integrity to compare
        ComputeInnerIntegrity(hashAlg, name, (UINT16)size, buffer,
                              &integrityToCompare);
        // Compare outer blob integrity
        if(!MemoryEqual2B(&integrity.b, &integrityToCompare.b))
            result = TPM_RC_INTEGRITY;
    }
    return result;
}

//** Public Functions

//*** AdjustAuthSize()
// This function will validate that the input authValue is no larger than the
// digestSize for the nameAlg. It will then pad with zeros to the size of the
// digest.
BOOL
AdjustAuthSize(
    TPM2B_AUTH          *auth,          // IN/OUT: value to adjust
    TPMI_ALG_HASH        nameAlg        // IN:
    )
{
    UINT16               digestSize;
//
    // If there is no nameAlg, then this is a LoadExternal and the authVale can
    // be any size up to the maximum allowed by the 
    digestSize = (nameAlg == TPM_ALG_NULL) ? sizeof(TPMU_HA) 
        : CryptHashGetDigestSize(nameAlg);
    if(digestSize < MemoryRemoveTrailingZeros(auth))
        return FALSE;
    else if(digestSize > auth->t.size)
        MemoryPad2B(&auth->b, digestSize);
    auth->t.size = digestSize;

    return TRUE;
}

//*** AreAttributesForParent()
// This function is called by create, load, and import functions.
// Note: The 'isParent' attribute is SET when an object is loaded and it has
// attributes that are suitable for a parent object.
//  Return Type: BOOL
//      TRUE(1)         properties are those of a parent
//      FALSE(0)        properties are not those of a parent
BOOL
ObjectIsParent(
    OBJECT          *parentObject   // IN: parent handle
    )
{
    return parentObject->attributes.isParent;
}

//*** CreateChecks()
// Attribute checks that are unique to creation.
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   sensitiveDataOrigin is not consistent with the
//                          object type
//  other                   returns from PublicAttributesValidation()
TPM_RC
CreateChecks(
    OBJECT              *parentObject,
    TPMT_PUBLIC         *publicArea,
    UINT16               sensitiveDataSize
    )
{
    TPMA_OBJECT          attributes = publicArea->objectAttributes;
    TPM_RC               result = TPM_RC_SUCCESS;
//
    // If the caller indicates that they have provided the data, then make sure that
    // they have provided some data.
    if((!IS_ATTRIBUTE(attributes, TPMA_OBJECT, sensitiveDataOrigin))
       && (sensitiveDataSize == 0))
        return TPM_RCS_ATTRIBUTES;
    // For an ordinary object, data can only be provided when sensitiveDataOrigin
    // is CLEAR
    if((parentObject != NULL)
       && (IS_ATTRIBUTE(attributes, TPMA_OBJECT, sensitiveDataOrigin))
       && (sensitiveDataSize != 0))
        return TPM_RCS_ATTRIBUTES;
    switch(publicArea->type)
    {
        case ALG_KEYEDHASH_VALUE:
            // if this is a data object (sign == decrypt == CLEAR) then the
            // TPM cannot be the data source.
            if(!IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign) 
               && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt)
               && IS_ATTRIBUTE(attributes, TPMA_OBJECT, sensitiveDataOrigin))
                result = TPM_RC_ATTRIBUTES;
            // comment out the next line in order to prevent a fixedTPM derivation
            // parent
//            break;  
        case ALG_SYMCIPHER_VALUE:
            // A restricted key symmetric key (SYMCIPHER and KEYEDHASH)
            // must have sensitiveDataOrigin SET unless it has fixedParent and 
            // fixedTPM CLEAR.
            if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted))
                if(!IS_ATTRIBUTE(attributes, TPMA_OBJECT, sensitiveDataOrigin))
                    if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedParent)
                       || IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM))
                        result = TPM_RCS_ATTRIBUTES;
            break;
        default: // Asymmetric keys cannot have the sensitive portion provided
            if(!IS_ATTRIBUTE(attributes, TPMA_OBJECT, sensitiveDataOrigin))
                result = TPM_RCS_ATTRIBUTES;
            break;
    }
    if(TPM_RC_SUCCESS == result)
    {
        result = PublicAttributesValidation(parentObject, publicArea);
    }
    return result;
}
//*** SchemeChecks
// This function is called by TPM2_LoadExternal() and PublicAttributesValidation().
// This function validates the schemes in the public area of an object.
//  Return Type: TPM_RC
//      TPM_RC_HASH         non-duplicable storage key and its parent have different
//                          name algorithm
//      TPM_RC_KDF          incorrect KDF specified for decrypting keyed hash object
//      TPM_RC_KEY          invalid key size values in an asymmetric key public area
//      TPM_RCS_SCHEME       inconsistent attributes 'decrypt', 'sign', 'restricted'
//                          and key's scheme ID; or hash algorithm is inconsistent
//                          with the scheme ID for keyed hash object
//      TPM_RC_SYMMETRIC    a storage key with no symmetric algorithm specified; or
//                          non-storage key with symmetric algorithm different from
// ALG_NULL
TPM_RC
SchemeChecks(
    OBJECT          *parentObject,  // IN: parent (null if primary seed)
    TPMT_PUBLIC     *publicArea     // IN: public area of the object
    )
{
    TPMT_SYM_DEF_OBJECT     *symAlgs = NULL;
    TPM_ALG_ID               scheme = TPM_ALG_NULL;
    TPMA_OBJECT              attributes = publicArea->objectAttributes;
    TPMU_PUBLIC_PARMS        *parms = &publicArea->parameters;
//
    switch(publicArea->type)
    {
        case ALG_SYMCIPHER_VALUE:
            symAlgs = &parms->symDetail.sym;
            // If this is a decrypt key, then only the block cipher modes (not
            // SMAC) are valid. TPM_ALG_NULL is OK too. If this is a 'sign' key,
            // then any mode that got through the unmarshaling is OK.
            if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt)
               && !CryptSymModeIsValid(symAlgs->mode.sym, TRUE))
                return TPM_RCS_SCHEME;
            break;
        case ALG_KEYEDHASH_VALUE:
            scheme = parms->keyedHashDetail.scheme.scheme;
            // if both sign and decrypt
            if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign) 
               == IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
            {
                // if both sign and decrypt are set or clear, then need
                // ALG_NULL as scheme
                if(scheme != TPM_ALG_NULL)
                    return TPM_RCS_SCHEME;
            }
            else if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign) 
                    && scheme != TPM_ALG_HMAC)
                return TPM_RCS_SCHEME;
            else if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
            {
                if(scheme != TPM_ALG_XOR)
                    return TPM_RCS_SCHEME;
                // If this is a derivation parent, then the KDF needs to be
                // SP800-108 for this implementation. This is the only derivation
                // supported by this implementation. Other implementations could
                // support additional schemes. There is no default.
                if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted))
                {
                    if(parms->keyedHashDetail.scheme.details.xor.kdf 
                       != TPM_ALG_KDF1_SP800_108)
                        return TPM_RCS_SCHEME;
                    // Must select a digest.
                    if(CryptHashGetDigestSize(
                        parms->keyedHashDetail.scheme.details.xor.hashAlg) == 0)
                        return TPM_RCS_HASH;
                }
            }
            break;
        default: // handling for asymmetric
            scheme = parms->asymDetail.scheme.scheme;
            symAlgs = &parms->asymDetail.symmetric;
            // if the key is both sign and decrypt, then the scheme must be
            // ALG_NULL because there is no way to specify both a sign and a
            // decrypt scheme in the key.
            if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign) 
               == IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
            {
                // scheme must be TPM_ALG_NULL
                if(scheme != TPM_ALG_NULL)
                    return TPM_RCS_SCHEME;
            }
            else if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign))
            {
                // If this is a signing key, see if it has a signing scheme
                if(CryptIsAsymSignScheme(publicArea->type, scheme))
                {
                    // if proper signing scheme then it needs a proper hash
                    if(parms->asymDetail.scheme.details.anySig.hashAlg 
                       == TPM_ALG_NULL)
                        return TPM_RCS_SCHEME;
                }
                else
                {
                    // signing key that does not have a proper signing scheme.
                    // This is OK if the key is not restricted and its scheme
                    // is TPM_ALG_NULL
                    if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted) 
                       || scheme != TPM_ALG_NULL)
                        return TPM_RCS_SCHEME;
                }
            }
            else if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
            {
                if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted))
                {
                    // for a restricted decryption key (a parent), scheme 
                    // is required to be TPM_ALG_NULL
                    if(scheme != TPM_ALG_NULL)
                        return TPM_RCS_SCHEME;
                }
                else
                {
                    // For an unrestricted decryption key, the scheme has to
                    // be a valid scheme or TPM_ALG_NULL
                    if(scheme != TPM_ALG_NULL &&
                       !CryptIsAsymDecryptScheme(publicArea->type, scheme))
                        return TPM_RCS_SCHEME;
                }
            }
            if(!IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted) 
               || !IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
            {
                // For an asymmetric key that is not a parent, the symmetric
                // algorithms must be TPM_ALG_NULL
                if(symAlgs->algorithm != TPM_ALG_NULL)
                    return TPM_RCS_SYMMETRIC;
            }
            // Special checks for an ECC key
#if ALG_ECC
            if(publicArea->type == TPM_ALG_ECC)
            {
                TPM_ECC_CURVE            curveID;
                const TPMT_ECC_SCHEME   *curveScheme;

                curveID = publicArea->parameters.eccDetail.curveID;
                curveScheme = CryptGetCurveSignScheme(curveID);
                // The curveId must be valid or the unmarshaling is busted.
                pAssert(curveScheme != NULL);

                // If the curveID requires a specific scheme, then the key must
                // select the same scheme
                if(curveScheme->scheme != TPM_ALG_NULL)
                {
                    TPMS_ECC_PARMS      *ecc = &publicArea->parameters.eccDetail;
                    if(scheme != curveScheme->scheme)
                        return TPM_RCS_SCHEME;
                    // The scheme can allow any hash, or not...
                    if(curveScheme->details.anySig.hashAlg != TPM_ALG_NULL
                       && (ecc->scheme.details.anySig.hashAlg
                           != curveScheme->details.anySig.hashAlg))
                        return TPM_RCS_SCHEME;
                }
                // For now, the KDF must be TPM_ALG_NULL
                if(publicArea->parameters.eccDetail.kdf.scheme != TPM_ALG_NULL)
                    return TPM_RCS_KDF;
            }
#endif
            break;
    }
    // If this is a restricted decryption key with symmetric algorithms, then it 
    // is an ordinary parent (not a derivation parent). It needs to specific
    // symmetric algorithms other than TPM_ALG_NULL
    if(symAlgs != NULL 
       && IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted) 
       && IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
    {
        if(symAlgs->algorithm == TPM_ALG_NULL)
            return TPM_RCS_SYMMETRIC;
#if 0       //??
// This next check is under investigation. Need to see if it will break Windows 
// before it is enabled. If it does not, then it should be default because a 
// the mode used with a parent is always CFB and Part 2 indicates as much.
        if(symAlgs->mode.sym != TPM_ALG_CFB)
            return TPM_RCS_MODE;
#endif
        // If this parent is not duplicable, then the symmetric algorithms 
        // (encryption and hash) must match those of its parent
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedParent) 
           && (parentObject != NULL))
        {
            if(publicArea->nameAlg != parentObject->publicArea.nameAlg)
                return TPM_RCS_HASH;
            if(!MemoryEqual(symAlgs, &parentObject->publicArea.parameters,
                            sizeof(TPMT_SYM_DEF_OBJECT)))
                return TPM_RCS_SYMMETRIC;
        }
    }
    return TPM_RC_SUCCESS;
}

//*** PublicAttributesValidation()
// This function validates the values in the public area of an object.
// This function is used in the processing of TPM2_Create, TPM2_CreatePrimary, 
// TPM2_CreateLoaded(), TPM2_Load(),  TPM2_Import(), and TPM2_LoadExternal().
// For TPM2_Import() this is only used if the new parent has fixedTPM SET. For
// TPM2_LoadExternal(), this is not used for a public-only key
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   'fixedTPM', 'fixedParent', or 'encryptedDuplication'
//                          attributes are inconsistent between themselves or with
//                          those of the parent object;
//                          inconsistent 'restricted', 'decrypt' and 'sign'
//                          attributes;
//                          attempt to inject sensitive data for an asymmetric key;
//                          attempt to create a symmetric cipher key that is not
//                          a decryption key
//      TPM_RC_HASH         nameAlg is TPM_ALG_NULL
//      TPM_RC_SIZE         'authPolicy' size does not match digest size of the name
//                          algorithm in 'publicArea'
//   other                  returns from SchemeChecks()
TPM_RC
PublicAttributesValidation(
    OBJECT          *parentObject,  // IN: input parent object
    TPMT_PUBLIC     *publicArea     // IN: public area of the object
    )
{
    TPMA_OBJECT      attributes = publicArea->objectAttributes;
    TPMA_OBJECT      parentAttributes = TPMA_ZERO_INITIALIZER();
//
    if(parentObject != NULL)
        parentAttributes = parentObject->publicArea.objectAttributes;
    if(publicArea->nameAlg == TPM_ALG_NULL)
        return TPM_RCS_HASH;
    // If there is an authPolicy, it needs to be the size of the digest produced
    // by the nameAlg of the object
    if((publicArea->authPolicy.t.size != 0
        && (publicArea->authPolicy.t.size
            != CryptHashGetDigestSize(publicArea->nameAlg))))
        return TPM_RCS_SIZE;
    // If the parent is fixedTPM (including a Primary Object) the object must have
    // the same value for fixedTPM and fixedParent
    if(parentObject == NULL 
       || IS_ATTRIBUTE(parentAttributes, TPMA_OBJECT, fixedTPM))
    {
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedParent) 
           != IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM))
            return TPM_RCS_ATTRIBUTES;
    }
    else
    {
        // The parent is not fixedTPM so the object can't be fixedTPM
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM))
            return  TPM_RCS_ATTRIBUTES;
    }
    // See if sign and decrypt are the same
    if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign) 
       == IS_ATTRIBUTE(attributes, TPMA_OBJECT, decrypt))
    {
        // a restricted key cannot have both SET or both CLEAR
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, restricted))
            return TPM_RC_ATTRIBUTES;
        // only a data object may have both sign and decrypt CLEAR
        // BTW, since we know that decrypt==sign, no need to check both
        if(publicArea->type != TPM_ALG_KEYEDHASH 
           && !IS_ATTRIBUTE(attributes, TPMA_OBJECT, sign))
            return TPM_RC_ATTRIBUTES;
    }
    // If the object can't be duplicated (directly or indirectly) then there
    // is no justification for having encryptedDuplication SET
    if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM) 
       && IS_ATTRIBUTE(attributes, TPMA_OBJECT, encryptedDuplication))
        return TPM_RCS_ATTRIBUTES;
    // If a parent object has fixedTPM CLEAR, the child must have the
    // same encryptedDuplication value as its parent.
    // Primary objects are considered to have a fixedTPM parent (the seeds).
    if(parentObject != NULL 
       && !IS_ATTRIBUTE(parentAttributes, TPMA_OBJECT, fixedTPM))
    {
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, encryptedDuplication) 
           != IS_ATTRIBUTE(parentAttributes, TPMA_OBJECT, encryptedDuplication))
            return TPM_RCS_ATTRIBUTES;
    }
    // Special checks for derived objects
    if((parentObject != NULL) && (parentObject->attributes.derivation == SET))
    {
        // A derived object has the same settings for fixedTPM as its parent
        if(IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedTPM) 
           != IS_ATTRIBUTE(parentAttributes, TPMA_OBJECT, fixedTPM))
            return TPM_RCS_ATTRIBUTES;
        // A derived object is required to be fixedParent
        if(!IS_ATTRIBUTE(attributes, TPMA_OBJECT, fixedParent))
            return TPM_RCS_ATTRIBUTES;
    }
    return SchemeChecks(parentObject, publicArea);
}

//*** FillInCreationData()
// Fill in creation data for an object.
//  Return Type: void
void
FillInCreationData(
    TPMI_DH_OBJECT           parentHandle,  // IN: handle of parent
    TPMI_ALG_HASH            nameHashAlg,   // IN: name hash algorithm
    TPML_PCR_SELECTION      *creationPCR,   // IN: PCR selection
    TPM2B_DATA              *outsideData,   // IN: outside data
    TPM2B_CREATION_DATA     *outCreation,   // OUT: creation data for output
    TPM2B_DIGEST            *creationDigest // OUT: creation digest
    )
{
    BYTE                 creationBuffer[sizeof(TPMS_CREATION_DATA)];
    BYTE                *buffer;
    HASH_STATE           hashState;
//
    // Fill in TPMS_CREATION_DATA in outCreation

    // Compute PCR digest
    PCRComputeCurrentDigest(nameHashAlg, creationPCR,
                            &outCreation->creationData.pcrDigest);

    // Put back PCR selection list
    outCreation->creationData.pcrSelect = *creationPCR;

    // Get locality
    outCreation->creationData.locality
        = LocalityGetAttributes(_plat__LocalityGet());
    outCreation->creationData.parentNameAlg = TPM_ALG_NULL;

    // If the parent is either a primary seed or TPM_ALG_NULL, then  the Name
    // and QN of the parent are the parent's handle.
    if(HandleGetType(parentHandle) == TPM_HT_PERMANENT)
    {
        buffer = &outCreation->creationData.parentName.t.name[0];
        outCreation->creationData.parentName.t.size =
            TPM_HANDLE_Marshal(&parentHandle, &buffer, NULL);
        // For a primary or temporary object, the parent name (a handle) and the
        // parent's QN are the same
        outCreation->creationData.parentQualifiedName
            = outCreation->creationData.parentName;
    }
    else         // Regular object
    {
        OBJECT          *parentObject = HandleToObject(parentHandle);
//
        // Set name algorithm
        outCreation->creationData.parentNameAlg = parentObject->publicArea.nameAlg;

        // Copy parent name
        outCreation->creationData.parentName = parentObject->name;

        // Copy parent qualified name
        outCreation->creationData.parentQualifiedName = parentObject->qualifiedName;
    }
    // Copy outside information
    outCreation->creationData.outsideInfo = *outsideData;

    // Marshal creation data to canonical form
    buffer = creationBuffer;
    outCreation->size = TPMS_CREATION_DATA_Marshal(&outCreation->creationData,
                                                   &buffer, NULL);
    // Compute hash for creation field in public template
    creationDigest->t.size = CryptHashStart(&hashState, nameHashAlg);
    CryptDigestUpdate(&hashState, outCreation->size, creationBuffer);
    CryptHashEnd2B(&hashState, &creationDigest->b);

    return;
}

//*** GetSeedForKDF()
// Get a seed for KDF.  The KDF for encryption and HMAC key use the same seed.
const TPM2B *
GetSeedForKDF(
    OBJECT          *protector         // IN: the protector handle
    )
{
    // Get seed for encryption key.  Use input seed if provided.
    // Otherwise, using protector object's seedValue.  TPM_RH_NULL is the only
    // exception that we may not have a loaded object as protector.  In such a
    // case, use nullProof as seed.
    if(protector == NULL)
        return &gr.nullProof.b;
    else
        return &protector->sensitive.seedValue.b;
}

//*** ProduceOuterWrap()
// This function produce outer wrap for a buffer containing the sensitive data.
// It requires the sensitive data being marshaled to the outerBuffer, with the
// leading bytes reserved for integrity hash.  If iv is used, iv space should
// be reserved at the beginning of the buffer.  It assumes the sensitive data
// starts at address (outerBuffer + integrity size {+ iv size}).
// This function performs:
//  1. Add IV before sensitive area if required
//  2. encrypt sensitive data, if iv is required, encrypt by iv.  otherwise,
//     encrypted by a NULL iv
//  3. add HMAC integrity at the beginning of the buffer
// It returns the total size of blob with outer wrap
UINT16
ProduceOuterWrap(
    OBJECT          *protector,     // IN: The handle of the object that provides
                                    //     protection.  For object, it is parent
                                    //     handle. For credential, it is the handle
                                    //     of encrypt object.
    TPM2B           *name,          // IN: the name of the object
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm for outer wrap
    TPM2B           *seed,          // IN: an external seed may be provided for
                                    //     duplication blob. For non duplication
                                    //     blob, this parameter should be NULL
    BOOL             useIV,         // IN: indicate if an IV is used
    UINT16           dataSize,      // IN: the size of sensitive data, excluding the
                                    //     leading integrity buffer size or the
                                    //     optional iv size
    BYTE            *outerBuffer    // IN/OUT: outer buffer with sensitive data in
                                    //     it
    )
{
    TPM_ALG_ID      symAlg;
    UINT16          keyBits;
    TPM2B_SYM_KEY   symKey;
    TPM2B_IV        ivRNG;          // IV from RNG
    TPM2B_IV        *iv = NULL;
    UINT16          ivSize = 0;     // size of iv area, including the size field
    BYTE            *sensitiveData; // pointer to the sensitive data
    TPM2B_DIGEST    integrity;
    UINT16          integritySize;
    BYTE            *buffer;        // Auxiliary buffer pointer
//
    // Compute the beginning of sensitive data.  The outer integrity should
    // always exist if this function is called to make an outer wrap
    integritySize = sizeof(UINT16) + CryptHashGetDigestSize(hashAlg);
    sensitiveData = outerBuffer + integritySize;

    // If iv is used, adjust the pointer of sensitive data and add iv before it
    if(useIV)
    {
        ivSize = GetIV2BSize(protector);

        // Generate IV from RNG.  The iv data size should be the total IV area
        // size minus the size of size field
        ivRNG.t.size = ivSize - sizeof(UINT16);
        CryptRandomGenerate(ivRNG.t.size, ivRNG.t.buffer);

        // Marshal IV to buffer
        buffer = sensitiveData;
        TPM2B_IV_Marshal(&ivRNG, &buffer, NULL);

        // adjust sensitive data starting after IV area
        sensitiveData += ivSize;

        // Use iv for encryption
        iv = &ivRNG;
    }
    // Compute symmetric key parameters for outer buffer encryption
    ComputeProtectionKeyParms(protector, hashAlg, name, seed,
                              &symAlg, &keyBits, &symKey);
    // Encrypt inner buffer in place
    CryptSymmetricEncrypt(sensitiveData, symAlg, keyBits,
                          symKey.t.buffer, iv, TPM_ALG_CFB, dataSize,
                          sensitiveData);
    // Compute outer integrity.  Integrity computation includes the optional IV
    // area
    ComputeOuterIntegrity(name, protector, hashAlg, seed, dataSize + ivSize,
                          outerBuffer + integritySize, &integrity);
    // Add integrity at the beginning of outer buffer
    buffer = outerBuffer;
    TPM2B_DIGEST_Marshal(&integrity, &buffer, NULL);

    // return the total size in outer wrap
    return dataSize + integritySize + ivSize;
}

//*** UnwrapOuter()
// This function remove the outer wrap of a blob containing sensitive data
// This function performs:
//  1. check integrity of outer blob
//  2. decrypt outer blob
//
//  Return Type: TPM_RC
//      TPM_RCS_INSUFFICIENT     error during sensitive data unmarshaling
//      TPM_RCS_INTEGRITY        sensitive data integrity is broken
//      TPM_RCS_SIZE             error during sensitive data unmarshaling
//      TPM_RCS_VALUE            IV size for CFB does not match the encryption
//                               algorithm block size
TPM_RC
UnwrapOuter(
    OBJECT          *protector,     // IN: The object that provides
                                    //     protection.  For object, it is parent
                                    //     handle. For credential, it is the
                                    //     encrypt object.
    TPM2B           *name,          // IN: the name of the object
    TPM_ALG_ID       hashAlg,       // IN: hash algorithm for outer wrap
    TPM2B           *seed,          // IN: an external seed may be provided for
                                    //     duplication blob. For non duplication
                                    //     blob, this parameter should be NULL.
    BOOL             useIV,         // IN: indicates if an IV is used
    UINT16           dataSize,      // IN: size of sensitive data in outerBuffer,
                                    //     including the leading integrity buffer
                                    //     size, and an optional iv area
    BYTE            *outerBuffer    // IN/OUT: sensitive data
    )
{
    TPM_RC          result;
    TPM_ALG_ID      symAlg = TPM_ALG_NULL;
    TPM2B_SYM_KEY   symKey;
    UINT16          keyBits = 0;
    TPM2B_IV        ivIn;               // input IV retrieved from input buffer
    TPM2B_IV        *iv = NULL;
    BYTE            *sensitiveData;     // pointer to the sensitive data
    TPM2B_DIGEST    integrityToCompare;
    TPM2B_DIGEST    integrity;
    INT32           size;
//
    // Unmarshal integrity
    sensitiveData = outerBuffer;
    size = (INT32)dataSize;
    result = TPM2B_DIGEST_Unmarshal(&integrity, &sensitiveData, &size);
    if(result == TPM_RC_SUCCESS)
    {
        // Compute integrity to compare
        ComputeOuterIntegrity(name, protector, hashAlg, seed,
                              (UINT16)size, sensitiveData,
                              &integrityToCompare);
        // Compare outer blob integrity
        if(!MemoryEqual2B(&integrity.b, &integrityToCompare.b))
            return TPM_RCS_INTEGRITY;
        // Get the symmetric algorithm parameters used for encryption
        ComputeProtectionKeyParms(protector, hashAlg, name, seed,
                                  &symAlg, &keyBits, &symKey);
        // Retrieve IV if it is used
        if(useIV)
        {
            result = TPM2B_IV_Unmarshal(&ivIn, &sensitiveData, &size);
            if(result == TPM_RC_SUCCESS)
            {
                // The input iv size for CFB must match the encryption algorithm
                // block size
                if(ivIn.t.size != CryptGetSymmetricBlockSize(symAlg, keyBits))
                    result = TPM_RC_VALUE;
                else
                    iv = &ivIn;
            }
        }
    }
    // If no errors, decrypt private in place. Since this function uses CFB, 
    // CryptSymmetricDecrypt() will not return any errors. It may fail but it will
    // not return an error.
    if(result == TPM_RC_SUCCESS)
        CryptSymmetricDecrypt(sensitiveData, symAlg, keyBits,
                              symKey.t.buffer, iv, TPM_ALG_CFB,
                              (UINT16)size, sensitiveData);
    return result;
}

//*** MarshalSensitive()
// This function is used to marshal a sensitive area. Among other things, it
// adjusts the size of the authValue to be no smaller than the digest of
// 'nameAlg'. It will also make sure that the RSA sensitive contains the right number
// of values.
// Returns the size of the marshaled area.
static UINT16
MarshalSensitive(
    OBJECT              *parent,            // IN: the object parent (optional)
    BYTE                *buffer,            // OUT: receiving buffer
    TPMT_SENSITIVE      *sensitive,         // IN: the sensitive area to marshal
    TPMI_ALG_HASH        nameAlg            // IN:
    )
{
    BYTE                *sizeField = buffer;    // saved so that size can be 
                                                // marshaled after it is known
    UINT16               retVal;
//
    // Pad the authValue if needed
    MemoryPad2B(&sensitive->authValue.b, CryptHashGetDigestSize(nameAlg));
    buffer += 2;

    // Marshal the structure
#if ALG_RSA
    // If the sensitive size is the special case for a prime in the type 
    if((sensitive->sensitive.rsa.t.size & RSA_prime_flag) > 0)
    {
        UINT16               sizeSave = sensitive->sensitive.rsa.t.size;
    //
        // Turn off the flag that indicates that the sensitive->sensitive contains
        // the CRT form of the exponent.
        sensitive->sensitive.rsa.t.size &= ~(RSA_prime_flag);
        // If the parent isn't fixedTPM, then truncate the sensitive data to be
        // the size of the prime. Otherwise, leave it at the current size which 
        // is the full CRT size.
        if(parent == NULL
           || !IS_ATTRIBUTE(parent->publicArea.objectAttributes,
                            TPMA_OBJECT, fixedTPM))
            sensitive->sensitive.rsa.t.size /= 5;
        retVal = TPMT_SENSITIVE_Marshal(sensitive, &buffer, NULL);
        // Restore the flag and the size.
        sensitive->sensitive.rsa.t.size = sizeSave;
    }
    else
#endif
    retVal = TPMT_SENSITIVE_Marshal(sensitive, &buffer, NULL);

    // Marshal the size
    retVal = (UINT16)(retVal + UINT16_Marshal(&retVal, &sizeField, NULL));

    return retVal;
}

//*** SensitiveToPrivate()
// This function prepare the private blob for off the chip storage
// The operations in this function:
//  1. marshal TPM2B_SENSITIVE structure into the buffer of TPM2B_PRIVATE
//  2. apply encryption to the sensitive area.
//  3. apply outer integrity computation.
void
SensitiveToPrivate(
    TPMT_SENSITIVE  *sensitive,     // IN: sensitive structure
    TPM2B_NAME      *name,          // IN: the name of the object
    OBJECT          *parent,        // IN: The parent object
    TPM_ALG_ID       nameAlg,       // IN: hash algorithm in public area.  This
                                    //     parameter is used when parentHandle is
                                    //     NULL, in which case the object is
                                    //     temporary.
    TPM2B_PRIVATE   *outPrivate     // OUT: output private structure
    )
{
    BYTE                *sensitiveData;     // pointer to the sensitive data
    UINT16              dataSize;           // data blob size
    TPMI_ALG_HASH       hashAlg;            // hash algorithm for integrity
    UINT16              integritySize;
    UINT16              ivSize;
//
    pAssert(name != NULL && name->t.size != 0);

    // Find the hash algorithm for integrity computation
    if(parent == NULL)
    {
        // For Temporary Object, using self name algorithm
        hashAlg = nameAlg;
    }
    else
    {
        // Otherwise, using parent's name algorithm
        hashAlg = parent->publicArea.nameAlg;
    }
    // Starting of sensitive data without wrappers
    sensitiveData = outPrivate->t.buffer;

    // Compute the integrity size
    integritySize = sizeof(UINT16) + CryptHashGetDigestSize(hashAlg);

    // Reserve space for integrity
    sensitiveData += integritySize;

    // Get iv size
    ivSize = GetIV2BSize(parent);

    // Reserve space for iv
    sensitiveData += ivSize;

    // Marshal the sensitive area including authValue size adjustments.
    dataSize = MarshalSensitive(parent, sensitiveData, sensitive, nameAlg);

    //Produce outer wrap, including encryption and HMAC
    outPrivate->t.size = ProduceOuterWrap(parent, &name->b, hashAlg, NULL,
                                          TRUE, dataSize, outPrivate->t.buffer);
    return;
}

//*** PrivateToSensitive()
// Unwrap a input private area.  Check the integrity, decrypt and retrieve data
// to a sensitive structure.
// The operations in this function:
//  1. check the integrity HMAC of the input private area
//  2. decrypt the private buffer
//  3. unmarshal TPMT_SENSITIVE structure into the buffer of TPMT_SENSITIVE
//  Return Type: TPM_RC
//      TPM_RCS_INTEGRITY       if the private area integrity is bad
//      TPM_RC_SENSITIVE        unmarshal errors while unmarshaling TPMS_ENCRYPT
//                              from input private
//      TPM_RCS_SIZE            error during sensitive data unmarshaling
//      TPM_RCS_VALUE           outer wrapper does not have an iV of the correct
//                              size
TPM_RC
PrivateToSensitive(
    TPM2B           *inPrivate,     // IN: input private structure
    TPM2B           *name,          // IN: the name of the object
    OBJECT          *parent,        // IN: parent object
    TPM_ALG_ID       nameAlg,       // IN: hash algorithm in public area.  It is
                                    //     passed separately because we only pass
                                    //     name, rather than the whole public area
                                    //     of the object.  This parameter is used in
                                    //     the following two cases: 1. primary
                                    //     objects. 2. duplication blob with inner
                                    //     wrap.  In other cases, this parameter
                                    //     will be ignored
    TPMT_SENSITIVE  *sensitive      // OUT: sensitive structure
    )
{
    TPM_RC          result;
    BYTE            *buffer;
    INT32           size;
    BYTE            *sensitiveData; // pointer to the sensitive data
    UINT16          dataSize;
    UINT16          dataSizeInput;
    TPMI_ALG_HASH   hashAlg;        // hash algorithm for integrity
    UINT16          integritySize;
    UINT16          ivSize;
//
    // Make sure that name is provided
    pAssert(name != NULL && name->size != 0);

    // Find the hash algorithm for integrity computation
    // For Temporary Object (parent == NULL) use self name algorithm;
    // Otherwise, using parent's name algorithm
    hashAlg = (parent == NULL) ? nameAlg : parent->publicArea.nameAlg;

    // unwrap outer
    result = UnwrapOuter(parent, name, hashAlg, NULL, TRUE,
                         inPrivate->size, inPrivate->buffer);
    if(result != TPM_RC_SUCCESS)
        return result;
    // Compute the inner integrity size.
    integritySize = sizeof(UINT16) + CryptHashGetDigestSize(hashAlg);

    // Get iv size
    ivSize = GetIV2BSize(parent);

    // The starting of sensitive data and data size without outer wrapper
    sensitiveData = inPrivate->buffer + integritySize + ivSize;
    dataSize = inPrivate->size - integritySize - ivSize;

    // Unmarshal input data size
    buffer = sensitiveData;
    size = (INT32)dataSize;
    result = UINT16_Unmarshal(&dataSizeInput, &buffer, &size);
    if(result == TPM_RC_SUCCESS)
    {
        if((dataSizeInput + sizeof(UINT16)) != dataSize)
            result = TPM_RC_SENSITIVE;
        else
        {
            // Unmarshal sensitive buffer to sensitive structure
            result = TPMT_SENSITIVE_Unmarshal(sensitive, &buffer, &size);
            if(result != TPM_RC_SUCCESS || size != 0)
            {
                result = TPM_RC_SENSITIVE;
            }
        }
    }
    return result;
}

//*** SensitiveToDuplicate()
// This function prepare the duplication blob from the sensitive area.
// The operations in this function:
//  1. marshal TPMT_SENSITIVE structure into the buffer of TPM2B_PRIVATE
//  2. apply inner wrap to the sensitive area if required
//  3. apply outer wrap if required
void
SensitiveToDuplicate(
    TPMT_SENSITIVE      *sensitive,     // IN: sensitive structure
    TPM2B               *name,          // IN: the name of the object
    OBJECT              *parent,        // IN: The new parent object
    TPM_ALG_ID           nameAlg,       // IN: hash algorithm in public area. It
                                        //     is passed separately because we
                                        //     only pass name, rather than the
                                        //     whole public area of the object.
    TPM2B               *seed,          // IN: the external seed. If external
                                        //     seed is provided with size of 0,
                                        //     no outer wrap should be applied
                                        //     to duplication blob.
    TPMT_SYM_DEF_OBJECT *symDef,        // IN: Symmetric key definition. If the
                                        //     symmetric key algorithm is NULL,
                                        //     no inner wrap should be applied.
    TPM2B_DATA          *innerSymKey,   // IN/OUT: a symmetric key may be
                                        //     provided to encrypt the inner
                                        //     wrap of a duplication blob. May
                                        //     be generated here if needed.
    TPM2B_PRIVATE       *outPrivate     // OUT: output private structure
    )
{
    BYTE            *sensitiveData; // pointer to the sensitive data
    TPMI_ALG_HASH   outerHash = TPM_ALG_NULL;// The hash algorithm for outer wrap
    TPMI_ALG_HASH   innerHash = TPM_ALG_NULL;// The hash algorithm for inner wrap
    UINT16          dataSize;       // data blob size
    BOOL            doInnerWrap = FALSE;
    BOOL            doOuterWrap = FALSE;
//
    // Make sure that name is provided
    pAssert(name != NULL && name->size != 0);

    // Make sure symDef and innerSymKey are not NULL
    pAssert(symDef != NULL && innerSymKey != NULL);

    // Starting of sensitive data without wrappers
    sensitiveData = outPrivate->t.buffer;

    // Find out if inner wrap is required
    if(symDef->algorithm != TPM_ALG_NULL)
    {
        doInnerWrap = TRUE;

        // Use self nameAlg as inner hash algorithm
        innerHash = nameAlg;

        // Adjust sensitive data pointer
        sensitiveData += sizeof(UINT16) + CryptHashGetDigestSize(innerHash);
    }
    // Find out if outer wrap is required
    if(seed->size != 0)
    {
        doOuterWrap = TRUE;

        // Use parent nameAlg as outer hash algorithm
        outerHash = parent->publicArea.nameAlg;

        // Adjust sensitive data pointer
        sensitiveData += sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
    }
    // Marshal sensitive area
    dataSize = MarshalSensitive(NULL, sensitiveData, sensitive, nameAlg);

    // Apply inner wrap for duplication blob.  It includes both integrity and
    // encryption
    if(doInnerWrap)
    {
        BYTE            *innerBuffer = NULL;
        BOOL            symKeyInput = TRUE;
        innerBuffer = outPrivate->t.buffer;
        // Skip outer integrity space
        if(doOuterWrap)
            innerBuffer += sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
        dataSize = ProduceInnerIntegrity(name, innerHash, dataSize,
                                         innerBuffer);
        // Generate inner encryption key if needed
        if(innerSymKey->t.size == 0)
        {
            innerSymKey->t.size = (symDef->keyBits.sym + 7) / 8;
            CryptRandomGenerate(innerSymKey->t.size, innerSymKey->t.buffer);

            // TPM generates symmetric encryption.  Set the flag to FALSE
            symKeyInput = FALSE;
        }
        else
        {
            // assume the input key size should matches the symmetric definition
            pAssert(innerSymKey->t.size == (symDef->keyBits.sym + 7) / 8);
        }

        // Encrypt inner buffer in place
        CryptSymmetricEncrypt(innerBuffer, symDef->algorithm,
                              symDef->keyBits.sym, innerSymKey->t.buffer, NULL,
                              TPM_ALG_CFB, dataSize, innerBuffer);

        // If the symmetric encryption key is imported, clear the buffer for
        // output
        if(symKeyInput)
            innerSymKey->t.size = 0;
    }
    // Apply outer wrap for duplication blob.  It includes both integrity and
    // encryption
    if(doOuterWrap)
    {
        dataSize = ProduceOuterWrap(parent, name, outerHash, seed, FALSE,
                                    dataSize, outPrivate->t.buffer);
    }
    // Data size for output
    outPrivate->t.size = dataSize;

    return;
}

//*** DuplicateToSensitive()
// Unwrap a duplication blob.  Check the integrity, decrypt and retrieve data
// to a sensitive structure.
// The operations in this function:
//  1. check the integrity HMAC of the input private area
//  2. decrypt the private buffer
//  3. unmarshal TPMT_SENSITIVE structure into the buffer of TPMT_SENSITIVE
//
//  Return Type: TPM_RC
//      TPM_RC_INSUFFICIENT      unmarshaling sensitive data from 'inPrivate' failed
//      TPM_RC_INTEGRITY         'inPrivate' data integrity is broken
//      TPM_RC_SIZE              unmarshaling sensitive data from 'inPrivate' failed
TPM_RC
DuplicateToSensitive(
    TPM2B               *inPrivate,     // IN: input private structure
    TPM2B               *name,          // IN: the name of the object
    OBJECT              *parent,        // IN: the parent
    TPM_ALG_ID           nameAlg,       // IN: hash algorithm in public area.
    TPM2B               *seed,          // IN: an external seed may be provided.
                                        //     If external seed is provided with
                                        //     size of 0, no outer wrap is
                                        //     applied
    TPMT_SYM_DEF_OBJECT *symDef,        // IN: Symmetric key definition. If the
                                        //     symmetric key algorithm is NULL,
                                        //     no inner wrap is applied
    TPM2B               *innerSymKey,   // IN: a symmetric key may be provided
                                        //     to decrypt the inner wrap of a
                                        //     duplication blob.
    TPMT_SENSITIVE      *sensitive      // OUT: sensitive structure
    )
{
    TPM_RC               result;
    BYTE                *buffer;
    INT32                size;
    BYTE                *sensitiveData; // pointer to the sensitive data
    UINT16               dataSize;
    UINT16               dataSizeInput;
//
    // Make sure that name is provided
    pAssert(name != NULL && name->size != 0);

    // Make sure symDef and innerSymKey are not NULL
    pAssert(symDef != NULL && innerSymKey != NULL);

    // Starting of sensitive data
    sensitiveData = inPrivate->buffer;
    dataSize = inPrivate->size;

    // Find out if outer wrap is applied
    if(seed->size != 0)
    {
        // Use parent nameAlg as outer hash algorithm
        TPMI_ALG_HASH   outerHash = parent->publicArea.nameAlg;

        result = UnwrapOuter(parent, name, outerHash, seed, FALSE,
                             dataSize, sensitiveData);
        if(result != TPM_RC_SUCCESS)
            return result;
        // Adjust sensitive data pointer and size
        sensitiveData += sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
        dataSize -= sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
    }
    // Find out if inner wrap is applied
    if(symDef->algorithm != TPM_ALG_NULL)
    {
        // assume the input key size matches the symmetric definition
        pAssert(innerSymKey->size == (symDef->keyBits.sym + 7) / 8);

        // Decrypt inner buffer in place
        CryptSymmetricDecrypt(sensitiveData, symDef->algorithm,
                              symDef->keyBits.sym, innerSymKey->buffer, NULL,
                              TPM_ALG_CFB, dataSize, sensitiveData);
        // Check inner integrity
        result = CheckInnerIntegrity(name, nameAlg, dataSize, sensitiveData);
        if(result != TPM_RC_SUCCESS)
            return result;
        // Adjust sensitive data pointer and size
        sensitiveData += sizeof(UINT16) + CryptHashGetDigestSize(nameAlg);
        dataSize -= sizeof(UINT16) + CryptHashGetDigestSize(nameAlg);
    }
    // Unmarshal input data size
    buffer = sensitiveData;
    size = (INT32)dataSize;
    result = UINT16_Unmarshal(&dataSizeInput, &buffer, &size);
    if(result == TPM_RC_SUCCESS)
    {
        if((dataSizeInput + sizeof(UINT16)) != dataSize)
            result = TPM_RC_SIZE;
        else
        {
            // Unmarshal sensitive buffer to sensitive structure
            result = TPMT_SENSITIVE_Unmarshal(sensitive, &buffer, &size);

            // if the results is OK make sure that all the data was unmarshaled
            if(result == TPM_RC_SUCCESS && size != 0)
                result = TPM_RC_SIZE;
        }
    }
    return result;
}

//*** SecretToCredential()
// This function prepare the credential blob from a secret (a TPM2B_DIGEST)
// The operations in this function:
//  1. marshal TPM2B_DIGEST structure into the buffer of TPM2B_ID_OBJECT
//  2. encrypt the private buffer, excluding the leading integrity HMAC area
//  3. compute integrity HMAC and append to the beginning of the buffer.
//  4. Set the total size of TPM2B_ID_OBJECT buffer
void
SecretToCredential(
    TPM2B_DIGEST        *secret,        // IN: secret information
    TPM2B               *name,          // IN: the name of the object
    TPM2B               *seed,          // IN: an external seed.
    OBJECT              *protector,     // IN: the protector
    TPM2B_ID_OBJECT     *outIDObject    // OUT: output credential
    )
{
    BYTE                *buffer;        // Auxiliary buffer pointer
    BYTE                *sensitiveData; // pointer to the sensitive data
    TPMI_ALG_HASH        outerHash;     // The hash algorithm for outer wrap
    UINT16               dataSize;      // data blob size
//
    pAssert(secret != NULL && outIDObject != NULL);

    // use protector's name algorithm as outer hash ????
    outerHash = protector->publicArea.nameAlg;

    // Marshal secret area to credential buffer, leave space for integrity
    sensitiveData = outIDObject->t.credential
        + sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
// Marshal secret area
    buffer = sensitiveData;
    dataSize = TPM2B_DIGEST_Marshal(secret, &buffer, NULL);

    // Apply outer wrap
    outIDObject->t.size = ProduceOuterWrap(protector, name, outerHash, seed, FALSE,
                                           dataSize, outIDObject->t.credential);
    return;
}

//*** CredentialToSecret()
// Unwrap a credential.  Check the integrity, decrypt and retrieve data
// to a TPM2B_DIGEST structure.
// The operations in this function:
//  1. check the integrity HMAC of the input credential area
//  2. decrypt the credential buffer
//  3. unmarshal TPM2B_DIGEST structure into the buffer of TPM2B_DIGEST
//
//  Return Type: TPM_RC
//      TPM_RC_INSUFFICIENT      error during credential unmarshaling
//      TPM_RC_INTEGRITY         credential integrity is broken
//      TPM_RC_SIZE              error during credential unmarshaling
//      TPM_RC_VALUE             IV size does not match the encryption algorithm
//                               block size
TPM_RC
CredentialToSecret(
    TPM2B               *inIDObject,    // IN: input credential blob
    TPM2B               *name,          // IN: the name of the object
    TPM2B               *seed,          // IN: an external seed.
    OBJECT              *protector,     // IN: the protector
    TPM2B_DIGEST        *secret         // OUT: secret information
    )
{
    TPM_RC                   result;
    BYTE                    *buffer;
    INT32                    size;
    TPMI_ALG_HASH            outerHash;     // The hash algorithm for outer wrap
    BYTE                    *sensitiveData; // pointer to the sensitive data
    UINT16                   dataSize;
//
    // use protector's name algorithm as outer hash
    outerHash = protector->publicArea.nameAlg;

    // Unwrap outer, a TPM_RC_INTEGRITY error may be returned at this point
    result = UnwrapOuter(protector, name, outerHash, seed, FALSE,
                         inIDObject->size, inIDObject->buffer);
    if(result == TPM_RC_SUCCESS)
    {
        // Compute the beginning of sensitive data
        sensitiveData = inIDObject->buffer
            + sizeof(UINT16) + CryptHashGetDigestSize(outerHash);
        dataSize = inIDObject->size
            - (sizeof(UINT16) + CryptHashGetDigestSize(outerHash));
        // Unmarshal secret buffer to TPM2B_DIGEST structure
        buffer = sensitiveData;
        size = (INT32)dataSize;
        result = TPM2B_DIGEST_Unmarshal(secret, &buffer, &size);

        // If there were no other unmarshaling errors, make sure that the
        // expected amount of data was recovered
        if(result == TPM_RC_SUCCESS && size != 0)
            return TPM_RC_SIZE;
    }
    return result;
}

//*** MemoryRemoveTrailingZeros()
// This function is used to adjust the length of an authorization value.
// It adjusts the size of the TPM2B so that it does not include octets
// at the end of the buffer that contain zero.
// The function returns the number of non-zero octets in the buffer.
UINT16
MemoryRemoveTrailingZeros(
    TPM2B_AUTH      *auth           // IN/OUT: value to adjust
    )
{
    while((auth->t.size > 0) && (auth->t.buffer[auth->t.size - 1] == 0))
        auth->t.size--;
    return auth->t.size;
}

//*** SetLabelAndContext()
// This function sets the label and context for a derived key. It is possible
// that 'label' or 'context' can end up being an Empty Buffer.
TPM_RC
SetLabelAndContext(
    TPMS_DERIVE             *labelContext,  // IN/OUT: the recovered label and 
                                            //      context
    TPM2B_SENSITIVE_DATA    *sensitive      // IN: the sensitive data
    )
{
    TPMS_DERIVE              sensitiveValue;
    TPM_RC                   result;
    INT32                    size;
    BYTE                    *buff;
//
    // Unmarshal a TPMS_DERIVE from the TPM2B_SENSITIVE_DATA buffer
    // If there is something to unmarshal...
    if(sensitive->t.size != 0)
    {
        size = sensitive->t.size;
        buff = sensitive->t.buffer;
        result = TPMS_DERIVE_Unmarshal(&sensitiveValue, &buff, &size);
        if(result != TPM_RC_SUCCESS)
            return result;
        // If there was a label in the public area leave it there, otherwise, copy
        // the new value
        if(labelContext->label.t.size == 0)
            MemoryCopy2B(&labelContext->label.b, &sensitiveValue.label.b,
                         sizeof(labelContext->label.t.buffer));
        // if there was a context string in publicArea, it overrides
        if(labelContext->context.t.size == 0)
            MemoryCopy2B(&labelContext->context.b, &sensitiveValue.context.b,
                         sizeof(labelContext->label.t.buffer));
    }
    return TPM_RC_SUCCESS;
}

//*** UnmarshalToPublic()
// Support function to unmarshal the template. This is used because the
// Input may be a TPMT_TEMPLATE and that structure does not have the same
// size as a TPMT_PUBLIC because of the difference between the 'unique' and
// 'seed' fields.
// If 'derive' is not NULL, then the 'seed' field is assumed to contain
// a 'label' and 'context' that are unmarshaled into 'derive'.
TPM_RC
UnmarshalToPublic(
    TPMT_PUBLIC         *tOut,       // OUT: output
    TPM2B_TEMPLATE      *tIn,        // IN:
    BOOL                 derivation, // IN: indicates if this is for a derivation
    TPMS_DERIVE         *labelContext// OUT: label and context if derivation
    )
{
    BYTE                *buffer = tIn->t.buffer;
    INT32                size = tIn->t.size;
    TPM_RC               result;
//
    // make sure that tOut is zeroed so that there are no remnants from previous
    // uses
    MemorySet(tOut, 0, sizeof(TPMT_PUBLIC));
    // Unmarshal the components of the TPMT_PUBLIC up to the unique field
    result = TPMI_ALG_PUBLIC_Unmarshal(&tOut->type, &buffer, &size);
    if(result != TPM_RC_SUCCESS)
        return result;
    result = TPMI_ALG_HASH_Unmarshal(&tOut->nameAlg, &buffer, &size, FALSE);
    if(result != TPM_RC_SUCCESS)
        return result;
    result = TPMA_OBJECT_Unmarshal(&tOut->objectAttributes, &buffer, &size);
    if(result != TPM_RC_SUCCESS)
        return result;
    result = TPM2B_DIGEST_Unmarshal(&tOut->authPolicy, &buffer, &size);
    if(result != TPM_RC_SUCCESS)
        return result;
    result = TPMU_PUBLIC_PARMS_Unmarshal(&tOut->parameters, &buffer, &size, 
                                         tOut->type);
    if(result != TPM_RC_SUCCESS)
        return result;
    // Now unmarshal a TPMS_DERIVE if this is for derivation
    if(derivation)
        result = TPMS_DERIVE_Unmarshal(labelContext, &buffer, &size);
    else
        // otherwise, unmarshal a TPMU_PUBLIC_ID
        result = TPMU_PUBLIC_ID_Unmarshal(&tOut->unique, &buffer, &size, 
                                          tOut->type);
    // Make sure the template was used up
    if((result == TPM_RC_SUCCESS) && (size != 0))
        result = TPM_RC_SIZE;
    return result;
}


//*** ObjectSetExternal()
// Set the external attributes for an object.
void
ObjectSetExternal(
    OBJECT      *object
    )
{
    object->attributes.external = SET;
}