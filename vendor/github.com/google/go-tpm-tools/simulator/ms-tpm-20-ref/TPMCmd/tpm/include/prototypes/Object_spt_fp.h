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

#ifndef    _OBJECT_SPT_FP_H_
#define    _OBJECT_SPT_FP_H_

//*** AdjustAuthSize()
// This function will validate that the input authValue is no larger than the
// digestSize for the nameAlg. It will then pad with zeros to the size of the
// digest.
BOOL
AdjustAuthSize(
    TPM2B_AUTH          *auth,          // IN/OUT: value to adjust
    TPMI_ALG_HASH        nameAlg        // IN:
);

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
);

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
);

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
);

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
);

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
);

//*** GetSeedForKDF()
// Get a seed for KDF.  The KDF for encryption and HMAC key use the same seed.
const TPM2B *
GetSeedForKDF(
    OBJECT          *protector         // IN: the protector handle
);

//*** ProduceOuterWrap()
// This function produce outer wrap for a buffer containing the sensitive data.
// It requires the sensitive data being marshaled to the outerBuffer, with the
// leading bytes reserved for integrity hash.  If iv is used, iv space should
// be reserved at the beginning of the buffer.  It assumes the sensitive data
// starts at address (outerBuffer + integrity size @).
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
);

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
);

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
);

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
);

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
);

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
);

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
);

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
);

//*** MemoryRemoveTrailingZeros()
// This function is used to adjust the length of an authorization value.
// It adjusts the size of the TPM2B so that it does not include octets
// at the end of the buffer that contain zero.
// The function returns the number of non-zero octets in the buffer.
UINT16
MemoryRemoveTrailingZeros(
    TPM2B_AUTH      *auth           // IN/OUT: value to adjust
);

//*** SetLabelAndContext()
// This function sets the label and context for a derived key. It is possible
// that 'label' or 'context' can end up being an Empty Buffer.
TPM_RC
SetLabelAndContext(
    TPMS_DERIVE             *labelContext,  // IN/OUT: the recovered label and
                                            //      context
    TPM2B_SENSITIVE_DATA    *sensitive      // IN: the sensitive data
);

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
);

//*** ObjectSetExternal()
// Set the external attributes for an object.
void
ObjectSetExternal(
    OBJECT      *object
);

#endif  // _OBJECT_SPT_FP_H_
