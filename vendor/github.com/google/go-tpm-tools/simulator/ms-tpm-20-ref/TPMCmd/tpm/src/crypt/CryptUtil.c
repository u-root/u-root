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
//** Introduction
//
//  This module contains the interfaces to the CryptoEngine and provides
//  miscellaneous cryptographic functions in support of the TPM.
//

//** Includes
#include "Tpm.h"

//****************************************************************************/
//**     Hash/HMAC Functions
//****************************************************************************/

//*** CryptHmacSign()
// Sign a digest using an HMAC key. This an HMAC of a digest, not an HMAC of a 
// message.
//  Return Type: TPM_RC
//      TPM_RC_HASH         not a valid hash
static TPM_RC
CryptHmacSign(
    TPMT_SIGNATURE      *signature,     // OUT: signature
    OBJECT              *signKey,       // IN: HMAC key sign the hash
    TPM2B_DIGEST        *hashData       // IN: hash to be signed
    )
{
    HMAC_STATE       hmacState;
    UINT32           digestSize;

    digestSize = CryptHmacStart2B(&hmacState, signature->signature.any.hashAlg,
                                  &signKey->sensitive.sensitive.bits.b);
    CryptDigestUpdate2B(&hmacState.hashState, &hashData->b);
    CryptHmacEnd(&hmacState, digestSize,
                 (BYTE *)&signature->signature.hmac.digest);
    return TPM_RC_SUCCESS;
}

//*** CryptHMACVerifySignature()
// This function will verify a signature signed by a HMAC key.
// Note that a caller needs to prepare 'signature' with the signature algorithm
// (TPM_ALG_HMAC) and the hash algorithm to use. This function then builds a
// signature of that type.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME           not the proper scheme for this key type
//      TPM_RC_SIGNATURE        if invalid input or signature is not genuine
static TPM_RC
CryptHMACVerifySignature(
    OBJECT              *signKey,       // IN: HMAC key signed the hash
    TPM2B_DIGEST        *hashData,      // IN: digest being verified
    TPMT_SIGNATURE      *signature      // IN: signature to be verified
    )
{
    TPMT_SIGNATURE           test;
    TPMT_KEYEDHASH_SCHEME   *keyScheme = 
                         &signKey->publicArea.parameters.keyedHashDetail.scheme;
//
    if((signature->sigAlg != ALG_HMAC_VALUE)
       || (signature->signature.hmac.hashAlg == ALG_NULL_VALUE))
        return TPM_RC_SCHEME;
    // This check is not really needed for verification purposes. However, it does
    // prevent someone from trying to validate a signature using a weaker hash 
    // algorithm than otherwise allowed by the key. That is, a key with a scheme
    // other than TMP_ALG_NULL can only be used to validate signatures that have
    // a matching scheme.
    if((keyScheme->scheme != ALG_NULL_VALUE)
       && ((keyScheme->scheme != signature->sigAlg)
           || (keyScheme->details.hmac.hashAlg 
               != signature->signature.any.hashAlg)))
        return TPM_RC_SIGNATURE;
    test.sigAlg = signature->sigAlg;
    test.signature.hmac.hashAlg = signature->signature.hmac.hashAlg;

    CryptHmacSign(&test, signKey, hashData);

    // Compare digest
    if(!MemoryEqual(&test.signature.hmac.digest,
                    &signature->signature.hmac.digest,
                    CryptHashGetDigestSize(signature->signature.any.hashAlg)))
        return TPM_RC_SIGNATURE;

    return TPM_RC_SUCCESS;
}

//*** CryptGenerateKeyedHash()
// This function creates a keyedHash object.
// Return type: TPM_RC
//      TPM_RC_NO_RESULT    cannot get values from random number generator
//      TPM_RC_SIZE         sensitive data size is larger than allowed for
//                          the scheme
static TPM_RC
CryptGenerateKeyedHash(
    TPMT_PUBLIC             *publicArea,        // IN/OUT: the public area template
                                                //     for the new key.
    TPMT_SENSITIVE          *sensitive,         // OUT: sensitive area
    TPMS_SENSITIVE_CREATE   *sensitiveCreate,   // IN: sensitive creation data
    RAND_STATE              *rand               // IN: "entropy" source
    )
{
    TPMT_KEYEDHASH_SCHEME   *scheme;
    TPM_ALG_ID               hashAlg;
    UINT16                   hashBlockSize;
    UINT16                   digestSize;

    scheme = &publicArea->parameters.keyedHashDetail.scheme;

    if(publicArea->type != ALG_KEYEDHASH_VALUE)
        return TPM_RC_FAILURE;

    // Pick the limiting hash algorithm
    if(scheme->scheme == ALG_NULL_VALUE)
        hashAlg = publicArea->nameAlg;
    else if(scheme->scheme == ALG_XOR_VALUE)
        hashAlg = scheme->details.xor.hashAlg;
    else
        hashAlg = scheme->details.hmac.hashAlg;
    hashBlockSize = CryptHashGetBlockSize(hashAlg);
    digestSize = CryptHashGetDigestSize(hashAlg);

    // if this is a signing or a decryption key, then the limit
    // for the data size is the block size of the hash. This limit
    // is set because larger values have lower entropy because of the
    // HMAC function. The lower limit is 1/2 the size of the digest
    //
    //If the user provided the key, check that it is a proper size
    if(sensitiveCreate->data.t.size != 0)
    {
        if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, decrypt)
           || IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, sign))
        {
            if(sensitiveCreate->data.t.size > hashBlockSize)
                return TPM_RC_SIZE;
#if 0   // May make this a FIPS-mode requirement
            if(sensitiveCreate->data.t.size < (digestSize / 2))
                return TPM_RC_SIZE;
#endif
        }
        // If this is a data blob, then anything that will get past the unmarshaling
        // is OK
        MemoryCopy2B(&sensitive->sensitive.bits.b, &sensitiveCreate->data.b,
                     sizeof(sensitive->sensitive.bits.t.buffer));
    }
    else
    {
        // The TPM is going to generate the data so set the size to be the
        // size of the digest of the algorithm
        sensitive->sensitive.bits.t.size =
            DRBG_Generate(rand, sensitive->sensitive.bits.t.buffer, digestSize);
        if(sensitive->sensitive.bits.t.size == 0)
            return (g_inFailureMode) ? TPM_RC_FAILURE : TPM_RC_NO_RESULT;
    }
    return TPM_RC_SUCCESS;
}

//*** CryptIsSchemeAnonymous()
// This function is used to test a scheme to see if it is an anonymous scheme
// The only anonymous scheme is ECDAA. ECDAA can be used to do things
// like U-Prove.
BOOL
CryptIsSchemeAnonymous(
    TPM_ALG_ID       scheme         // IN: the scheme algorithm to test
    )
{
    return scheme == ALG_ECDAA_VALUE;
}

//**** ************************************************************************
//** Symmetric Functions
//**** ************************************************************************

//*** ParmDecryptSym()
//  This function performs parameter decryption using symmetric block cipher.
/*(See Part 1 specification)
// Symmetric parameter decryption
//      When parameter decryption uses a symmetric block cipher, a decryption
//      key and IV will be generated from:
//      KDFa(hash, sessionAuth, "CFB", nonceNewer, nonceOlder, bits)    (24)
//      Where:
//      hash            the hash function associated with the session
//      sessionAuth     the sessionAuth associated with the session
//      nonceNewer      nonceCaller for a command
//      nonceOlder      nonceTPM for a command
//      bits            the number of bits required for the symmetric key
//                      plus an IV
*/
void
ParmDecryptSym(
    TPM_ALG_ID       symAlg,        // IN: the symmetric algorithm
    TPM_ALG_ID       hash,          // IN: hash algorithm for KDFa
    UINT16           keySizeInBits, // IN: the key size in bits
    TPM2B           *key,           // IN: KDF HMAC key
    TPM2B           *nonceCaller,   // IN: nonce caller
    TPM2B           *nonceTpm,      // IN: nonce TPM
    UINT32           dataSize,      // IN: size of parameter buffer
    BYTE            *data           // OUT: buffer to be decrypted
    )
{
    // KDF output buffer
    // It contains parameters for the CFB encryption
    // From MSB to LSB, they are the key and iv
    BYTE             symParmString[MAX_SYM_KEY_BYTES + MAX_SYM_BLOCK_SIZE];
    // Symmetric key size in byte
    UINT16           keySize = (keySizeInBits + 7) / 8;
    TPM2B_IV         iv;

    iv.t.size = CryptGetSymmetricBlockSize(symAlg, keySizeInBits);
    // If there is decryption to do...
    if(iv.t.size > 0)
    {
        // Generate key and iv
        CryptKDFa(hash, key, CFB_KEY, nonceCaller, nonceTpm,
                  keySizeInBits + (iv.t.size * 8), symParmString, NULL, FALSE);
        MemoryCopy(iv.t.buffer, &symParmString[keySize], iv.t.size);

        CryptSymmetricDecrypt(data, symAlg, keySizeInBits, symParmString,
                              &iv, ALG_CFB_VALUE, dataSize, data);
    }
    return;
}

//*** ParmEncryptSym()
//  This function performs parameter encryption using symmetric block cipher.
/*(See part 1 specification)
//      When parameter decryption uses a symmetric block cipher, an encryption
//      key and IV will be generated from:
//      KDFa(hash, sessionAuth, "CFB", nonceNewer, nonceOlder, bits)    (24)
//      Where:
//      hash            the hash function associated with the session
//      sessionAuth     the sessionAuth associated with the session
//      nonceNewer      nonceTPM for a response
//      nonceOlder      nonceCaller for a response
//      bits            the number of bits required for the symmetric key
//                      plus an IV
*/
void
ParmEncryptSym(
    TPM_ALG_ID       symAlg,        // IN: symmetric algorithm
    TPM_ALG_ID       hash,          // IN: hash algorithm for KDFa
    UINT16           keySizeInBits, // IN: symmetric key size in bits
    TPM2B           *key,           // IN: KDF HMAC key
    TPM2B           *nonceCaller,   // IN: nonce caller
    TPM2B           *nonceTpm,      // IN: nonce TPM
    UINT32           dataSize,      // IN: size of parameter buffer
    BYTE            *data           // OUT: buffer to be encrypted
    )
{
    // KDF output buffer
    // It contains parameters for the CFB encryption
    BYTE             symParmString[MAX_SYM_KEY_BYTES + MAX_SYM_BLOCK_SIZE];

    // Symmetric key size in bytes
    UINT16           keySize = (keySizeInBits + 7) / 8;

    TPM2B_IV         iv;

    iv.t.size = CryptGetSymmetricBlockSize(symAlg, keySizeInBits);
    // See if there is any encryption to do
    if(iv.t.size > 0)
    {
        // Generate key and iv
        CryptKDFa(hash, key, CFB_KEY, nonceTpm, nonceCaller,
                  keySizeInBits + (iv.t.size * 8), symParmString, NULL, FALSE);
        MemoryCopy(iv.t.buffer, &symParmString[keySize], iv.t.size);

        CryptSymmetricEncrypt(data, symAlg, keySizeInBits, symParmString, &iv,
                              ALG_CFB_VALUE, dataSize, data);
    }
    return;
}

//*** CryptGenerateKeySymmetric()
// This function generates a symmetric cipher key. The derivation process is
// determined by the type of the provided 'rand'
// Return type: TPM_RC
//      TPM_RC_NO_RESULT    cannot get a random value      
//      TPM_RC_KEY_SIZE     key size in the public area does not match the size
//                          in the sensitive creation area
//      TPM_RC_KEY          provided key value is not allowed
static TPM_RC
CryptGenerateKeySymmetric(
    TPMT_PUBLIC             *publicArea,        // IN/OUT: The public area template
                                                //     for the new key.
    TPMT_SENSITIVE          *sensitive,         // OUT: sensitive area
    TPMS_SENSITIVE_CREATE   *sensitiveCreate,   // IN: sensitive creation data
    RAND_STATE              *rand               // IN: the "entropy" source for
    )
{
    UINT16           keyBits = publicArea->parameters.symDetail.sym.keyBits.sym;
    TPM_RC           result;
//
    // only do multiples of RADIX_BITS
    if((keyBits % RADIX_BITS) != 0)
        return TPM_RC_KEY_SIZE;
    // If this is not a new key, then the provided key data must be the right size
    if(sensitiveCreate->data.t.size != 0)
    {
        result = CryptSymKeyValidate(&publicArea->parameters.symDetail.sym,
                                     (TPM2B_SYM_KEY *)&sensitiveCreate->data);
        if(result == TPM_RC_SUCCESS)
            MemoryCopy2B(&sensitive->sensitive.sym.b, &sensitiveCreate->data.b,
                         sizeof(sensitive->sensitive.sym.t.buffer));
    } 
#if ALG_TDES
    else if(publicArea->parameters.symDetail.sym.algorithm == ALG_TDES_VALUE)
    {
        result = CryptGenerateKeyDes(publicArea, sensitive, rand);
    }
#endif
    else
    {
        sensitive->sensitive.sym.t.size = 
            DRBG_Generate(rand, sensitive->sensitive.sym.t.buffer, 
                          BITS_TO_BYTES(keyBits));
        if(g_inFailureMode)
            result = TPM_RC_FAILURE;
        else if(sensitive->sensitive.sym.t.size == 0)
            result = TPM_RC_NO_RESULT;
        else
            result = TPM_RC_SUCCESS;
    }
    return result;
}

//*** CryptXORObfuscation()
// This function implements XOR obfuscation. It should not be called if the
// hash algorithm is not implemented. The only return value from this function
// is TPM_RC_SUCCESS.
void
CryptXORObfuscation(
    TPM_ALG_ID       hash,          // IN: hash algorithm for KDF
    TPM2B           *key,           // IN: KDF key
    TPM2B           *contextU,      // IN: contextU
    TPM2B           *contextV,      // IN: contextV
    UINT32           dataSize,      // IN: size of data buffer
    BYTE            *data           // IN/OUT: data to be XORed in place
    )
{
    BYTE             mask[MAX_DIGEST_SIZE]; // Allocate a digest sized buffer
    BYTE            *pm;
    UINT32           i;
    UINT32           counter = 0;
    UINT16           hLen = CryptHashGetDigestSize(hash);
    UINT32           requestSize = dataSize * 8;
    INT32            remainBytes = (INT32)dataSize;

    pAssert((key != NULL) && (data != NULL) && (hLen != 0));

    // Call KDFa to generate XOR mask
    for(; remainBytes > 0; remainBytes -= hLen)
    {
        // Make a call to KDFa to get next iteration
        CryptKDFa(hash, key, XOR_KEY, contextU, contextV,
                  requestSize, mask, &counter, TRUE);

    // XOR next piece of the data
        pm = mask;
        for(i = hLen < remainBytes ? hLen : remainBytes; i > 0; i--)
            *data++ ^= *pm++;
    }
    return;
}

//****************************************************************************
//** Initialization and shut down
//****************************************************************************

//*** CryptInit()
// This function is called when the TPM receives a _TPM_Init indication.
//
// NOTE: The hash algorithms do not have to be tested, they just need to be
// available. They have to be tested before the TPM can accept HMAC authorization
// or return any result that relies on a hash algorithm.
//  Return Type: BOOL
//      TRUE(1)         initializations succeeded
//      FALSE(0)        initialization failed and caller should place the TPM into
//                      Failure Mode
BOOL
CryptInit(
    void
    )
{
    BOOL         ok;
    // Initialize the vector of implemented algorithms
    AlgorithmGetImplementedVector(&g_implementedAlgorithms);

    // Indicate that all test are necessary
    CryptInitializeToTest();

    // Do any library initializations that are necessary. If any fails,
    // the caller should go into failure mode;
    ok = SupportLibInit();
    ok = ok && CryptSymInit();
    ok = ok && CryptRandInit();
    ok = ok && CryptHashInit();
#if ALG_RSA
    ok = ok && CryptRsaInit();
#endif // ALG_RSA
#if ALG_ECC
    ok = ok && CryptEccInit();
#endif // ALG_ECC
    return ok;
}

//*** CryptStartup()
// This function is called by TPM2_Startup() to initialize the functions in
// this cryptographic library and in the provided CryptoLibrary. This function
// and CryptUtilInit() are both provided so that the implementation may move the
// initialization around to get the best interaction.
//  Return Type: BOOL
//      TRUE(1)         startup succeeded
//      FALSE(0)        startup failed and caller should place the TPM into
//                      Failure Mode
BOOL
CryptStartup(
    STARTUP_TYPE     type           // IN: the startup type
    )
{
    BOOL            OK;
    NOT_REFERENCED(type);

    OK = CryptSymStartup() && CryptRandStartup() && CryptHashStartup()
#if ALG_RSA
        &&  CryptRsaStartup()
#endif // ALG_RSA
#if ALG_ECC
        &&  CryptEccStartup()
#endif // ALG_ECC
        ;
#if ALG_ECC
    // Don't directly check for SU_RESET because that is the default
    if(OK && (type != SU_RESTART) && (type != SU_RESUME))
    {
    // If the shutdown was orderly, then the values recovered from NV will
    // be OK to use.
        // Get a new  random commit nonce
        gr.commitNonce.t.size = sizeof(gr.commitNonce.t.buffer);
        CryptRandomGenerate(gr.commitNonce.t.size, gr.commitNonce.t.buffer);
        // Reset the counter and commit array
        gr.commitCounter = 0;
        MemorySet(gr.commitArray, 0, sizeof(gr.commitArray));
    }
#endif // ALG_ECC
    return OK;
}

//****************************************************************************
//** Algorithm-Independent Functions
//****************************************************************************
//*** Introduction
// These functions are used generically when a function of a general type
// (e.g., symmetric encryption) is required.  The functions will modify the
// parameters as required to interface to the indicated algorithms.
//
//*** CryptIsAsymAlgorithm()
// This function indicates if an algorithm is an asymmetric algorithm.
//  Return Type: BOOL
//      TRUE(1)         if it is an asymmetric algorithm
//      FALSE(0)        if it is not an asymmetric algorithm
BOOL
CryptIsAsymAlgorithm(
    TPM_ALG_ID       algID          // IN: algorithm ID
    )
{
    switch(algID)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
#endif
#if ALG_ECC
        case ALG_ECC_VALUE:
#endif
            return TRUE;
            break;
        default:
            break;
    }
    return FALSE;
}

//*** CryptSecretEncrypt()
// This function creates a secret value and its associated secret structure using
// an asymmetric algorithm.
//
// This function is used by TPM2_Rewrap() TPM2_MakeCredential(),
// and TPM2_Duplicate().
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   'keyHandle' does not reference a valid decryption key
//      TPM_RC_KEY          invalid ECC key (public point is not on the curve)
//      TPM_RC_SCHEME       RSA key with an unsupported padding scheme
//      TPM_RC_VALUE        numeric value of the data to be decrypted is greater
//                          than the RSA key modulus
TPM_RC
CryptSecretEncrypt(
    OBJECT                  *encryptKey,    // IN: encryption key object
    const TPM2B             *label,         // IN: a null-terminated string as L
    TPM2B_DATA              *data,          // OUT: secret value
    TPM2B_ENCRYPTED_SECRET  *secret         // OUT: secret structure
    )
{
    TPMT_RSA_DECRYPT         scheme;
    TPM_RC                   result = TPM_RC_SUCCESS;
//
    if(data == NULL || secret == NULL)
        return TPM_RC_FAILURE;

    // The output secret value has the size of the digest produced by the nameAlg.
    data->t.size = CryptHashGetDigestSize(encryptKey->publicArea.nameAlg);
    // The encryption scheme is OAEP using the nameAlg of the encrypt key.
    scheme.scheme = ALG_OAEP_VALUE;
    scheme.details.anySig.hashAlg = encryptKey->publicArea.nameAlg;

    if(!IS_ATTRIBUTE(encryptKey->publicArea.objectAttributes, TPMA_OBJECT, decrypt))
        return TPM_RC_ATTRIBUTES;
    switch(encryptKey->publicArea.type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
        {
            // Create secret data from RNG
            CryptRandomGenerate(data->t.size, data->t.buffer);

            // Encrypt the data by RSA OAEP into encrypted secret
            result = CryptRsaEncrypt((TPM2B_PUBLIC_KEY_RSA *)secret, &data->b,
                                     encryptKey, &scheme, label, NULL);
        }
        break;
#endif // ALG_RSA

#if ALG_ECC
        case ALG_ECC_VALUE:
        {
            TPMS_ECC_POINT      eccPublic;
            TPM2B_ECC_PARAMETER eccPrivate;
            TPMS_ECC_POINT      eccSecret;
            BYTE                *buffer = secret->t.secret;

            // Need to make sure that the public point of the key is on the
            // curve defined by the key.
            if(!CryptEccIsPointOnCurve(
                encryptKey->publicArea.parameters.eccDetail.curveID,
                &encryptKey->publicArea.unique.ecc))
                result = TPM_RC_KEY;
            else
            {
                // Call crypto engine to create an auxiliary ECC key
                // We assume crypt engine initialization should always success.
                // Otherwise, TPM should go to failure mode.

                CryptEccNewKeyPair(&eccPublic, &eccPrivate,
                               encryptKey->publicArea.parameters.eccDetail.curveID);
                // Marshal ECC public to secret structure. This will be used by the
                // recipient to decrypt the secret with their private key.
                secret->t.size = TPMS_ECC_POINT_Marshal(&eccPublic, &buffer, NULL);

                // Compute ECDH shared secret which is R = [d]Q where d is the
                // private part of the ephemeral key and Q is the public part of a
                // TPM key. TPM_RC_KEY error return from CryptComputeECDHSecret
                // because the auxiliary ECC key is just created according to the
                // parameters of input ECC encrypt key.
                if(CryptEccPointMultiply(&eccSecret,
                             encryptKey->publicArea.parameters.eccDetail.curveID,
                             &encryptKey->publicArea.unique.ecc, &eccPrivate,
                             NULL, NULL) 
                   != TPM_RC_SUCCESS)
                    result = TPM_RC_KEY;
                else
                {
                    // The secret value is computed from Z using KDFe as:
                    // secret := KDFe(HashID, Z, Use, PartyUInfo, PartyVInfo, bits)
                    // Where:
                    //  HashID  the nameAlg of the decrypt key
                    //  Z   the x coordinate (Px) of the product (P) of the point
                    //      (Q) of the secret and the private x coordinate (de,V)
                    //      of the decryption key
                    //  Use a null-terminated string containing "SECRET"
                    //  PartyUInfo  the x coordinate of the point in the secret
                    //              (Qe,U )
                    //  PartyVInfo  the x coordinate of the public key (Qs,V )
                    //  bits    the number of bits in the digest of HashID
                    // Retrieve seed from KDFe
                    CryptKDFe(encryptKey->publicArea.nameAlg, &eccSecret.x.b,
                              label, &eccPublic.x.b,
                              &encryptKey->publicArea.unique.ecc.x.b,
                              data->t.size * 8, data->t.buffer);
                }
            }
        }
        break;
#endif // ALG_ECC
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    return result;
}

//*** CryptSecretDecrypt()
// Decrypt a secret value by asymmetric (or symmetric) algorithm
// This function is used for ActivateCredential and Import for asymmetric
// decryption, and StartAuthSession for both asymmetric and symmetric
// decryption process
//
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES        RSA key is not a decryption key
//      TPM_RC_BINDING           Invalid RSA key (public and private parts are not
//                               cryptographically bound.
//      TPM_RC_ECC_POINT         ECC point in the secret is not on the curve
//      TPM_RC_INSUFFICIENT      failed to retrieve ECC point from the secret
//      TPM_RC_NO_RESULT         multiplication resulted in ECC point at infinity
//      TPM_RC_SIZE              data to decrypt is not of the same size as RSA key
//      TPM_RC_VALUE             For RSA key, numeric value of the encrypted data is
//                               greater than the modulus, or the recovered data is
//                               larger than the output buffer.
//                               For keyedHash or symmetric key, the secret is
//                               larger than the size of the digest produced by
//                               the name algorithm.
//      TPM_RC_FAILURE           internal error
TPM_RC
CryptSecretDecrypt(
    OBJECT                  *decryptKey,    // IN: decrypt key
    TPM2B_NONCE             *nonceCaller,   // IN: nonceCaller.  It is needed for
                                            //     symmetric decryption.  For
                                            //     asymmetric decryption, this
                                            //     parameter is NULL
    const TPM2B             *label,         // IN: a value for L
    TPM2B_ENCRYPTED_SECRET  *secret,        // IN: input secret
    TPM2B_DATA              *data           // OUT: decrypted secret value
    )
{
    TPM_RC      result = TPM_RC_SUCCESS;

    // Decryption for secret
    switch(decryptKey->publicArea.type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
        {
            TPMT_RSA_DECRYPT        scheme;
            TPMT_RSA_SCHEME         *keyScheme
                = &decryptKey->publicArea.parameters.rsaDetail.scheme;
            UINT16                   digestSize;

            scheme = *(TPMT_RSA_DECRYPT *)keyScheme;
            // If the key scheme is ALG_NULL_VALUE, set the scheme to OAEP and
            // set the algorithm to the name algorithm.
            if(scheme.scheme == ALG_NULL_VALUE)
            {
                // Use OAEP scheme
                scheme.scheme = ALG_OAEP_VALUE;
                scheme.details.oaep.hashAlg = decryptKey->publicArea.nameAlg;
            }
            // use the digestSize as an indicator of whether or not the scheme
            // is using a supported hash algorithm.
            // Note: depending on the scheme used for encryption, a hashAlg might
            // not be needed. However, the return value has to have some upper
            // limit on the size. In this case, it is the size of the digest of the
            // hash algorithm. It is checked after the decryption is done but, there
            // is no point in doing the decryption if the size is going to be
            // 'wrong' anyway.
            digestSize = CryptHashGetDigestSize(scheme.details.oaep.hashAlg);
            if(scheme.scheme != ALG_OAEP_VALUE || digestSize == 0)
                return TPM_RC_SCHEME;

            // Set the output buffer capacity
            data->t.size = sizeof(data->t.buffer);

            // Decrypt seed by RSA OAEP
            result = CryptRsaDecrypt(&data->b, &secret->b,
                                     decryptKey, &scheme, label);
            if((result == TPM_RC_SUCCESS) && (data->t.size > digestSize))
                result = TPM_RC_VALUE;
        }
        break;
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
        {
            TPMS_ECC_POINT       eccPublic;
            TPMS_ECC_POINT       eccSecret;
            BYTE                *buffer = secret->t.secret;
            INT32                size = secret->t.size;

            // Retrieve ECC point from secret buffer
            result = TPMS_ECC_POINT_Unmarshal(&eccPublic, &buffer, &size);
            if(result == TPM_RC_SUCCESS)
            {
                result = CryptEccPointMultiply(&eccSecret,
                              decryptKey->publicArea.parameters.eccDetail.curveID,
                              &eccPublic, &decryptKey->sensitive.sensitive.ecc,
                              NULL, NULL);
                if(result == TPM_RC_SUCCESS)
                {
                    // Set the size of the "recovered" secret value to be the size
                    // of the digest produced by the nameAlg.
                    data->t.size =
                        CryptHashGetDigestSize(decryptKey->publicArea.nameAlg);

                // The secret value is computed from Z using KDFe as:
                // secret := KDFe(HashID, Z, Use, PartyUInfo, PartyVInfo, bits)
                // Where:
                //  HashID -- the nameAlg of the decrypt key
                //  Z --  the x coordinate (Px) of the product (P) of the point
                //        (Q) of the secret and the private x coordinate (de,V)
                //        of the decryption key
                //  Use -- a null-terminated string containing "SECRET"
                //  PartyUInfo -- the x coordinate of the point in the secret
                //              (Qe,U )
                //  PartyVInfo -- the x coordinate of the public key (Qs,V )
                //  bits -- the number of bits in the digest of HashID
                // Retrieve seed from KDFe
                    CryptKDFe(decryptKey->publicArea.nameAlg, &eccSecret.x.b, label,
                              &eccPublic.x.b,
                              &decryptKey->publicArea.unique.ecc.x.b,
                              data->t.size * 8, data->t.buffer);
                }
            }
        }
        break;
#endif // ALG_ECC
#if     !ALG_KEYEDHASH
#   error   "KEYEDHASH support is required"
#endif
        case ALG_KEYEDHASH_VALUE:
            // The seed size can not be bigger than the digest size of nameAlg
            if(secret->t.size >
               CryptHashGetDigestSize(decryptKey->publicArea.nameAlg))
                result = TPM_RC_VALUE;
            else
            {
                // Retrieve seed by XOR Obfuscation:
                //    seed = XOR(secret, hash, key, nonceCaller, nullNonce)
                //    where:
                //    secret  the secret parameter from the TPM2_StartAuthHMAC
                //            command that contains the seed value
                //    hash    nameAlg  of tpmKey
                //    key     the key or data value in the object referenced by
                //            entityHandle in the TPM2_StartAuthHMAC command
                //    nonceCaller the parameter from the TPM2_StartAuthHMAC command
                //    nullNonce   a zero-length nonce
                // XOR Obfuscation in place
                CryptXORObfuscation(decryptKey->publicArea.nameAlg,
                                    &decryptKey->sensitive.sensitive.bits.b,
                                    &nonceCaller->b, NULL,
                                    secret->t.size, secret->t.secret);
                // Copy decrypted seed
                MemoryCopy2B(&data->b, &secret->b, sizeof(data->t.buffer));
            }
            break;
        case ALG_SYMCIPHER_VALUE:
        {
            TPM2B_IV                iv = {{0}};
            TPMT_SYM_DEF_OBJECT     *symDef;
            // The seed size can not be bigger than the digest size of nameAlg
            if(secret->t.size >
               CryptHashGetDigestSize(decryptKey->publicArea.nameAlg))
                result = TPM_RC_VALUE;
            else
            {
                symDef = &decryptKey->publicArea.parameters.symDetail.sym;
                iv.t.size = CryptGetSymmetricBlockSize(symDef->algorithm,
                                                       symDef->keyBits.sym);
                if(iv.t.size == 0)
                    return TPM_RC_FAILURE;
                if(nonceCaller->t.size >= iv.t.size)
                {
                    MemoryCopy(iv.t.buffer, nonceCaller->t.buffer, iv.t.size);
                }
                else
                {
                    if(nonceCaller->t.size > sizeof(iv.t.buffer))
                        return TPM_RC_FAILURE;
                    MemoryCopy(iv.b.buffer, nonceCaller->t.buffer,
                                 nonceCaller->t.size);
                }
                // make sure secret will fit
                if(secret->t.size > data->t.size)
                    return TPM_RC_FAILURE;
                data->t.size = secret->t.size;
                // CFB decrypt, using nonceCaller as iv
                CryptSymmetricDecrypt(data->t.buffer, symDef->algorithm,
                                      symDef->keyBits.sym,
                                      decryptKey->sensitive.sensitive.sym.t.buffer,
                                      &iv, ALG_CFB_VALUE, secret->t.size,
                                      secret->t.secret);
            }
        }
        break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    return result;
}

//*** CryptParameterEncryption()
// This function does in-place encryption of a response parameter.
void
CryptParameterEncryption(
    TPM_HANDLE       handle,            // IN: encrypt session handle
    TPM2B           *nonceCaller,       // IN: nonce caller
    UINT16           leadingSizeInByte, // IN: the size of the leading size field in
                                        //     bytes
    TPM2B_AUTH      *extraKey,          // IN: additional key material other than
                                        //     sessionAuth
    BYTE            *buffer             // IN/OUT: parameter buffer to be encrypted
    )
{
    SESSION     *session = SessionGet(handle);  // encrypt session
    TPM2B_TYPE(TEMP_KEY, (sizeof(extraKey->t.buffer)
                          + sizeof(session->sessionKey.t.buffer)));
    TPM2B_TEMP_KEY        key;               // encryption key
    UINT32               cipherSize = 0;    // size of cipher text
//
    // Retrieve encrypted data size.
    if(leadingSizeInByte == 2)
    {
        // Extract the first two bytes as the size field as the data size
        // encrypt
        cipherSize = (UINT32)BYTE_ARRAY_TO_UINT16(buffer);
        // advance the buffer
        buffer = &buffer[2];
    }
#ifdef      TPM4B
    else if(leadingSizeInByte == 4)
    {
        // use the first four bytes to indicate the number of bytes to encrypt
        cipherSize = BYTE_ARRAY_TO_UINT32(buffer);
        //advance pointer
        buffer = &buffer[4];
    }
#endif
    else
    {
        FAIL(FATAL_ERROR_INTERNAL);
    }

    // Compute encryption key by concatenating sessionKey with extra key
    MemoryCopy2B(&key.b, &session->sessionKey.b, sizeof(key.t.buffer));
    MemoryConcat2B(&key.b, &extraKey->b, sizeof(key.t.buffer));

    if(session->symmetric.algorithm == ALG_XOR_VALUE)

        // XOR parameter encryption formulation:
        //    XOR(parameter, hash, sessionAuth, nonceNewer, nonceOlder)
        CryptXORObfuscation(session->authHashAlg, &(key.b),
                            &(session->nonceTPM.b),
                            nonceCaller, cipherSize, buffer);
    else
        ParmEncryptSym(session->symmetric.algorithm, session->authHashAlg,
                       session->symmetric.keyBits.aes, &(key.b),
                       nonceCaller, &(session->nonceTPM.b),
                       cipherSize, buffer);
    return;
}

//*** CryptParameterDecryption()
// This function does in-place decryption of a command parameter.
//  Return Type: TPM_RC
//      TPM_RC_SIZE             The number of bytes in the input buffer is less than
//                              the number of bytes to be decrypted.
TPM_RC
CryptParameterDecryption(
    TPM_HANDLE       handle,            // IN: encrypted session handle
    TPM2B           *nonceCaller,       // IN: nonce caller
    UINT32           bufferSize,        // IN: size of parameter buffer
    UINT16           leadingSizeInByte, // IN: the size of the leading size field in
                                        //     byte
    TPM2B_AUTH      *extraKey,          // IN: the authValue
    BYTE            *buffer             // IN/OUT: parameter buffer to be decrypted
    )
{
    SESSION         *session = SessionGet(handle);  // encrypt session
    // The HMAC key is going to be the concatenation of the session key and any
    // additional key material (like the authValue). The size of both of these
    // is the size of the buffer which can contain a TPMT_HA.
    TPM2B_TYPE(HMAC_KEY, (sizeof(extraKey->t.buffer)
                          + sizeof(session->sessionKey.t.buffer)));
    TPM2B_HMAC_KEY          key;            // decryption key
    UINT32                  cipherSize = 0; // size of cipher text
//
    // Retrieve encrypted data size.
    if(leadingSizeInByte == 2)
    {
        // The first two bytes of the buffer are the size of the
        // data to be decrypted
        cipherSize = (UINT32)BYTE_ARRAY_TO_UINT16(buffer);
        buffer = &buffer[2];   // advance the buffer
    }
#ifdef  TPM4B
    else if(leadingSizeInByte == 4)
    {
        // the leading size is four bytes so get the four byte size field
        cipherSize = BYTE_ARRAY_TO_UINT32(buffer);
        buffer = &buffer[4];   //advance pointer
    }
#endif
    else
    {
        FAIL(FATAL_ERROR_INTERNAL);
    }
    if(cipherSize > bufferSize)
        return TPM_RC_SIZE;

    // Compute decryption key by concatenating sessionAuth with extra input key
    MemoryCopy2B(&key.b, &session->sessionKey.b, sizeof(key.t.buffer));
    MemoryConcat2B(&key.b, &extraKey->b, sizeof(key.t.buffer));

    if(session->symmetric.algorithm == ALG_XOR_VALUE)
        // XOR parameter decryption formulation:
        //    XOR(parameter, hash, sessionAuth, nonceNewer, nonceOlder)
        // Call XOR obfuscation function
        CryptXORObfuscation(session->authHashAlg, &key.b, nonceCaller,
                            &(session->nonceTPM.b), cipherSize, buffer);
    else
        // Assume that it is one of the symmetric block ciphers.
        ParmDecryptSym(session->symmetric.algorithm, session->authHashAlg,
                       session->symmetric.keyBits.sym,
                       &key.b, nonceCaller, &session->nonceTPM.b,
                       cipherSize, buffer);

    return TPM_RC_SUCCESS;
}

//*** CryptComputeSymmetricUnique()
// This function computes the unique field in public area for symmetric objects.
void
CryptComputeSymmetricUnique(
    TPMT_PUBLIC     *publicArea,    // IN: the object's public area
    TPMT_SENSITIVE  *sensitive,     // IN: the associated sensitive area
    TPM2B_DIGEST    *unique         // OUT: unique buffer
    )
{
    // For parents (symmetric and derivation), use an HMAC to compute
    // the 'unique' field
    if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, restricted)
       && IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, decrypt))
    {
        // Unique field is HMAC(sensitive->seedValue, sensitive->sensitive)
        HMAC_STATE      hmacState;
        unique->b.size = CryptHmacStart2B(&hmacState, publicArea->nameAlg,
                                          &sensitive->seedValue.b);
        CryptDigestUpdate2B(&hmacState.hashState, 
                            &sensitive->sensitive.any.b);
        CryptHmacEnd2B(&hmacState, &unique->b);
    }
    else
    {
        HASH_STATE  hashState;
        // Unique := Hash(sensitive->seedValue || sensitive->sensitive)
        unique->t.size = CryptHashStart(&hashState, publicArea->nameAlg);
        CryptDigestUpdate2B(&hashState, &sensitive->seedValue.b);
        CryptDigestUpdate2B(&hashState, &sensitive->sensitive.any.b);
        CryptHashEnd2B(&hashState, &unique->b);
    }
    return;
}

//*** CryptCreateObject()
// This function creates an object.
// For an asymmetric key, it will create a key pair and, for a parent key, a seed
// value for child protections.
//
// For an symmetric object, (TPM_ALG_SYMCIPHER or TPM_ALG_KEYEDHASH), it will
// create a secret key if the caller did not provide one. It will create a random
// secret seed value that is hashed with the secret value to create the public
// unique value.
//
// 'publicArea', 'sensitive', and 'sensitiveCreate' are the only required parameters
// and are the only ones that are used by TPM2_Create(). The other parameters
// are optional and are used when the generated Object needs to be deterministic.
// This is the case for both Primary Objects and Derived Objects.
//
// When a seed value is provided, a RAND_STATE will be populated and used for
// all operations in the object generation that require a random number. In the
// simplest case, TPM2_CreatePrimary() will use 'seed', 'label' and 'context' with
// context being the hash of the template. If the Primary Object is in
// the Endorsement hierarchy, it will also populate 'proof' with ehProof.
//
// For derived keys, 'seed' will be the secret value from the parent, 'label' and
// 'context' will be set according to the parameters of TPM2_CreateLoaded() and 
// 'hashAlg' will be set which causes the RAND_STATE to be a KDF generator.
//
//  Return Type: TPM_RC
//      TPM_RC_KEY          a provided key is not an allowed value
//      TPM_RC_KEY_SIZE     key size in the public area does not match the size
//                          in the sensitive creation area for a symmetric key
//      TPM_RC_NO_RESULT    unable to get random values (only in derivation)
//      TPM_RC_RANGE        for an RSA key, the exponent is not supported
//      TPM_RC_SIZE         sensitive data size is larger than allowed for the
//                          scheme for a keyed hash object
//      TPM_RC_VALUE        exponent is not prime or could not find a prime using
//                          the provided parameters for an RSA key;
//                          unsupported name algorithm for an ECC key
TPM_RC
CryptCreateObject(
    OBJECT                  *object,            // IN: new object structure pointer
    TPMS_SENSITIVE_CREATE   *sensitiveCreate,   // IN: sensitive creation
    RAND_STATE              *rand               // IN: the random number generator
                                                //      to use
    )
{
    TPMT_PUBLIC             *publicArea = &object->publicArea;
    TPMT_SENSITIVE          *sensitive = &object->sensitive;
    TPM_RC                   result = TPM_RC_SUCCESS;
//
    // Set the sensitive type for the object
    sensitive->sensitiveType = publicArea->type;

    // For all objects, copy the initial authorization data
    sensitive->authValue = sensitiveCreate->userAuth;

    // If the TPM is the source of the data, set the size of the provided data to
    // zero so that there's no confusion about what to do.
    if(IS_ATTRIBUTE(publicArea->objectAttributes, 
                    TPMA_OBJECT, sensitiveDataOrigin))
        sensitiveCreate->data.t.size = 0;

    // Generate the key and unique fields for the asymmetric keys and just the
    // sensitive value for symmetric object
    switch(publicArea->type)
    {
#if ALG_RSA
        // Create RSA key
        case ALG_RSA_VALUE:
            // RSA uses full object so that it has a place to put the private
            // exponent
            result = CryptRsaGenerateKey(publicArea, sensitive, rand);
            break;
#endif // ALG_RSA

#if ALG_ECC
        // Create ECC key
        case ALG_ECC_VALUE:
            result = CryptEccGenerateKey(publicArea, sensitive, rand);
            break;
#endif // ALG_ECC
        case ALG_SYMCIPHER_VALUE:
            result = CryptGenerateKeySymmetric(publicArea, sensitive,
                                               sensitiveCreate, rand);
            break;
        case ALG_KEYEDHASH_VALUE:
            result = CryptGenerateKeyedHash(publicArea, sensitive,
                                            sensitiveCreate, rand);
            break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    if(result != TPM_RC_SUCCESS)
        return result;
// Create the sensitive seed value
    // If this is a primary key in the endorsement hierarchy, stir the DRBG state
    // This implementation uses both shProof and ehProof to make sure that there
    // is no leakage of either.
    if(object->attributes.primary && object->attributes.epsHierarchy)
    {
        DRBG_AdditionalData((DRBG_STATE *)rand, &gp.shProof.b);
        DRBG_AdditionalData((DRBG_STATE *)rand, &gp.ehProof.b);
    }
    // Generate a seedValue that is the size of the digest produced by nameAlg
    sensitive->seedValue.t.size =
        DRBG_Generate(rand, sensitive->seedValue.t.buffer, 
                      CryptHashGetDigestSize(publicArea->nameAlg));
    if(g_inFailureMode)
        return TPM_RC_FAILURE;
    else if(sensitive->seedValue.t.size == 0)
        return TPM_RC_NO_RESULT;
    // For symmetric objects, need to compute the unique value for the public area
    if(publicArea->type == ALG_SYMCIPHER_VALUE
       || publicArea->type == ALG_KEYEDHASH_VALUE)
    {
        CryptComputeSymmetricUnique(publicArea, sensitive, &publicArea->unique.sym);
    }
    else
    {
        // if this is an asymmetric key and it isn't a parent, then
        // get rid of the seed.
        if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, sign)
           || !IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, restricted))
            memset(&sensitive->seedValue, 0, sizeof(sensitive->seedValue));
    }
    // Compute the name
    PublicMarshalAndComputeName(publicArea, &object->name);
    return result;
}

//*** CryptGetSignHashAlg()
// Get the hash algorithm of signature from a TPMT_SIGNATURE structure.
// It assumes the signature is not NULL
//  This is a function for easy access
TPMI_ALG_HASH
CryptGetSignHashAlg(
    TPMT_SIGNATURE  *auth           // IN: signature
    )
{
    if(auth->sigAlg == ALG_NULL_VALUE)
        FAIL(FATAL_ERROR_INTERNAL);

    // Get authHash algorithm based on signing scheme
    switch(auth->sigAlg)
    {
#if ALG_RSA
    // If RSA is supported, both RSASSA and RSAPSS are required
#   if !defined ALG_RSASSA_VALUE || !defined ALG_RSAPSS_VALUE
#       error "RSASSA and RSAPSS are required for RSA"
#   endif
        case ALG_RSASSA_VALUE:
            return auth->signature.rsassa.hash;
        case ALG_RSAPSS_VALUE:
            return auth->signature.rsapss.hash;
#endif // ALG_RSA

#if ALG_ECC
    // If ECC is defined, ECDSA is mandatory
#   if !ALG_ECDSA
#       error "ECDSA is requried for ECC"
#   endif
        case ALG_ECDSA_VALUE:
    // SM2 and ECSCHNORR are optional

#   if      ALG_SM2
        case ALG_SM2_VALUE:
#   endif
#   if  ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
#   endif
        //all ECC signatures look the same
            return auth->signature.ecdsa.hash;

#   if  ALG_ECDAA
        // Don't know how to verify an ECDAA signature
        case ALG_ECDAA_VALUE:
            break;
#   endif

#endif // ALG_ECC

        case ALG_HMAC_VALUE:
            return auth->signature.hmac.hashAlg;

        default:
            break;
    }
    return ALG_NULL_VALUE;
}

//*** CryptIsSplitSign()
// This function us used to determine if the signing operation is a split
// signing operation that required a TPM2_Commit().
//
BOOL
CryptIsSplitSign(
    TPM_ALG_ID       scheme         // IN: the algorithm selector
    )
{
    switch(scheme)
    {
#   if  ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TRUE;
            break;
#   endif   // ALG_ECDAA
        default:
            return FALSE;
            break;
    }
}

//*** CryptIsAsymSignScheme()
// This function indicates if a scheme algorithm is a sign algorithm.
BOOL
CryptIsAsymSignScheme(
    TPMI_ALG_PUBLIC          publicType,        // IN: Type of the object
    TPMI_ALG_ASYM_SCHEME     scheme             // IN: the scheme
    )
{
    BOOL            isSignScheme = TRUE;

    switch(publicType)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            switch(scheme)
            {
#   if !ALG_RSASSA  || !ALG_RSAPSS
#       error "RSASSA and PSAPSS required if RSA used."
#   endif
                case ALG_RSASSA_VALUE:
                case ALG_RSAPSS_VALUE:
                    break;
                default:
                    isSignScheme = FALSE;
                    break;
            }
            break;
#endif // ALG_RSA

#if ALG_ECC
        // If ECC is implemented ECDSA is required
        case ALG_ECC_VALUE:
            switch(scheme)
            {
                // Support for ECDSA is required for ECC
                case ALG_ECDSA_VALUE:
#if ALG_ECDAA // ECDAA is optional
                case ALG_ECDAA_VALUE:
#endif
#if ALG_ECSCHNORR // Schnorr is also optional
                case ALG_ECSCHNORR_VALUE:
#endif
#if ALG_SM2 // SM2 is optional
                case ALG_SM2_VALUE:
#endif
                    break;
                default:
                    isSignScheme = FALSE;
                    break;
            }
            break;
#endif // ALG_ECC
        default:
            isSignScheme = FALSE;
            break;
    }
    return isSignScheme;
}

//*** CryptIsAsymDecryptScheme()
// This function indicate if a scheme algorithm is a decrypt algorithm.
BOOL
CryptIsAsymDecryptScheme(
    TPMI_ALG_PUBLIC          publicType,        // IN: Type of the object
    TPMI_ALG_ASYM_SCHEME     scheme             // IN: the scheme
    )
{
    BOOL        isDecryptScheme = TRUE;

    switch(publicType)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            switch(scheme)
            {
                case ALG_RSAES_VALUE:
                case ALG_OAEP_VALUE:
                    break;
                default:
                    isDecryptScheme = FALSE;
                    break;
            }
            break;
#endif // ALG_RSA

#if ALG_ECC
        // If ECC is implemented ECDH is required
        case ALG_ECC_VALUE:
            switch(scheme)
            {
#if  !ALG_ECDH
#   error "ECDH is required for ECC"
#endif
                case ALG_ECDH_VALUE:
#if ALG_SM2
                case ALG_SM2_VALUE:
#endif
#if ALG_ECMQV
                case ALG_ECMQV_VALUE:
#endif
                    break;
                default:
                    isDecryptScheme = FALSE;
                    break;
            }
            break;
#endif // ALG_ECC
        default:
            isDecryptScheme = FALSE;
            break;
    }
    return isDecryptScheme;
}

//*** CryptSelectSignScheme()
// This function is used by the attestation and signing commands.  It implements
// the rules for selecting the signature scheme to use in signing. This function
// requires that the signing key either be TPM_RH_NULL or be loaded.
//
// If a default scheme is defined in object, the default scheme should be chosen,
// otherwise, the input scheme should be chosen.
// In the case that  both object and input scheme has a non-NULL scheme
// algorithm, if the schemes are compatible, the input scheme will be chosen.
//
// This function should not be called if 'signObject->publicArea.type' ==
// ALG_SYMCIPHER.
//
//  Return Type: BOOL
//      TRUE(1)         scheme selected
//      FALSE(0)        both 'scheme' and key's default scheme are empty; or
//                      'scheme' is empty while key's default scheme requires
//                      explicit input scheme (split signing); or
//                      non-empty default key scheme differs from 'scheme'
BOOL
CryptSelectSignScheme(
    OBJECT              *signObject,    // IN: signing key
    TPMT_SIG_SCHEME     *scheme         // IN/OUT: signing scheme
    )
{
    TPMT_SIG_SCHEME     *objectScheme;
    TPMT_PUBLIC         *publicArea;
    BOOL                 OK;

    // If the signHandle is TPM_RH_NULL, then the NULL scheme is used, regardless
    // of the setting of scheme
    if(signObject == NULL)
    {
        OK = TRUE;
        scheme->scheme = ALG_NULL_VALUE;
        scheme->details.any.hashAlg = ALG_NULL_VALUE;
    }
    else
    {
        // assignment to save typing.
        publicArea = &signObject->publicArea;

        // A symmetric cipher can be used to encrypt and decrypt but it can't
        // be used for signing
        if(publicArea->type == ALG_SYMCIPHER_VALUE)
            return FALSE;
        // Point to the scheme object
        if(CryptIsAsymAlgorithm(publicArea->type))
            objectScheme =
            (TPMT_SIG_SCHEME *)&publicArea->parameters.asymDetail.scheme;
        else
            objectScheme =
            (TPMT_SIG_SCHEME *)&publicArea->parameters.keyedHashDetail.scheme;

    // If the object doesn't have a default scheme, then use the
    // input scheme.
        if(objectScheme->scheme == ALG_NULL_VALUE)
        {
            // Input and default can't both be NULL
            OK = (scheme->scheme != ALG_NULL_VALUE);
            // Assume that the scheme is compatible with the key. If not,
            // an error will be generated in the signing operation.
        }
        else if(scheme->scheme == ALG_NULL_VALUE)
        {
            // input scheme is NULL so use default

            // First, check to see if the default requires that the caller
            // provided scheme data
            OK = !CryptIsSplitSign(objectScheme->scheme);
            if(OK)
            {
                // The object has a scheme and the input is TPM_ALG_NULL so copy
                // the object scheme as the final scheme. It is better to use a
                // structure copy than a copy of the individual fields.
                *scheme = *objectScheme;
            }
        }
        else
        {
            // Both input and object have scheme selectors
            // If the scheme and the hash are not the same then...
            // NOTE: the reason that there is no copy here is that the input
            // might contain extra data for a split signing scheme and that
            // data is not in the object so, it has to be preserved.
            OK = (objectScheme->scheme == scheme->scheme)
                && (objectScheme->details.any.hashAlg
                    == scheme->details.any.hashAlg);
        }
    }
    return OK;
}

//*** CryptSign()
// Sign a digest with asymmetric key or HMAC.
// This function is called by attestation commands and the generic TPM2_Sign
// command.
// This function checks the key scheme and digest size.  It does not
// check if the sign operation is allowed for restricted key.  It should be
// checked before the function is called.
// The function will assert if the key is not a signing key.
//
//  Return Type: TPM_RC
//      TPM_RC_SCHEME      'signScheme' is not compatible with the signing key type
//      TPM_RC_VALUE       'digest' value is greater than the modulus of
//                         'signHandle' or size of 'hashData' does not match hash
//                         algorithm in'signScheme' (for an RSA key);
//                         invalid commit status or failed to generate "r" value
//                         (for an ECC key)
TPM_RC
CryptSign(
    OBJECT              *signKey,       // IN: signing key
    TPMT_SIG_SCHEME     *signScheme,    // IN: sign scheme.
    TPM2B_DIGEST        *digest,        // IN: The digest being signed
    TPMT_SIGNATURE      *signature      // OUT: signature
    )
{
    TPM_RC               result = TPM_RC_SCHEME;

    // Initialize signature scheme
    signature->sigAlg = signScheme->scheme;

    // If the signature algorithm is TPM_ALG_NULL or the signing key is NULL,
    // then we are done
    if((signature->sigAlg == ALG_NULL_VALUE) || (signKey == NULL))
        return TPM_RC_SUCCESS;

    // Initialize signature hash
    // Note: need to do the check for TPM_ALG_NULL first because the null scheme
    // doesn't have a hashAlg member.
    signature->signature.any.hashAlg = signScheme->details.any.hashAlg;

    // perform sign operation based on different key type
    switch(signKey->publicArea.type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            result = CryptRsaSign(signature, signKey, digest, NULL);
            break;
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            // The reason that signScheme is passed to CryptEccSign but not to the
            // other signing methods is that the signing for ECC may be split and
            // need the 'r' value that is in the scheme but not in the signature.
            result = CryptEccSign(signature, signKey, digest,
                                  (TPMT_ECC_SCHEME *)signScheme, NULL);
            break;
#endif // ALG_ECC
        case ALG_KEYEDHASH_VALUE:
            result = CryptHmacSign(signature, signKey, digest);
            break;
        default:
            FAIL(FATAL_ERROR_INTERNAL);
            break;
    }
    return result;
}

//*** CryptValidateSignature()
// This function is used to verify a signature.  It is called by
// TPM2_VerifySignature() and TPM2_PolicySigned.
//
// Since this operation only requires use of a public key, no consistency
// checks are necessary for the key to signature type because a caller can load
// any public key that they like with any scheme that they like. This routine
// simply makes sure that the signature is correct, whatever the type.
//
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE            the signature is not genuine
//      TPM_RC_SCHEME               the scheme is not supported
//      TPM_RC_HANDLE               an HMAC key was selected but the
//                                  private part of the key is not loaded
TPM_RC
CryptValidateSignature(
    TPMI_DH_OBJECT   keyHandle,     // IN: The handle of sign key
    TPM2B_DIGEST    *digest,        // IN: The digest being validated
    TPMT_SIGNATURE  *signature      // IN: signature
    )
{
    // NOTE: HandleToObject will either return a pointer to a loaded object or
    // will assert. It will never return a non-valid value. This makes it save
    // to initialize 'publicArea' with the return value from HandleToObject() 
    // without checking it first.
    OBJECT              *signObject = HandleToObject(keyHandle);
    TPMT_PUBLIC         *publicArea = &signObject->publicArea;
    TPM_RC               result = TPM_RC_SCHEME;

    // The input unmarshaling should prevent any input signature from being
    // a NULL signature, but just in case
    if(signature->sigAlg == ALG_NULL_VALUE)
        return TPM_RC_SIGNATURE;

    switch(publicArea->type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
        {
    //
            // Call RSA code to verify signature
            result = CryptRsaValidateSignature(signature, signObject, digest);
            break;
        }
#endif // ALG_RSA

#if ALG_ECC
        case ALG_ECC_VALUE:
            result = CryptEccValidateSignature(signature, signObject, digest);
            break;
#endif // ALG_ECC

        case ALG_KEYEDHASH_VALUE:
            if(signObject->attributes.publicOnly)
                result = TPM_RCS_HANDLE;
            else
                result = CryptHMACVerifySignature(signObject, digest, signature);
            break;
        default:
            break;
    }
    return result;
}

//*** CryptGetTestResult
// This function returns the results of a self-test function.
// Note: the behavior in this function is NOT the correct behavior for a real
// TPM implementation.  An artificial behavior is placed here due to the
// limitation of a software simulation environment.  For the correct behavior,
// consult the part 3 specification for TPM2_GetTestResult().
TPM_RC
CryptGetTestResult(
    TPM2B_MAX_BUFFER    *outData        // OUT: test result data
    )
{
    outData->t.size = 0;
    return TPM_RC_SUCCESS;
}

//*** CryptValidateKeys()
// This function is used to verify that the key material of and object is valid.
// For a 'publicOnly' object, the key is verified for size and, if it is an ECC
// key, it is verified to be on the specified curve. For a key with a sensitive 
// area, the binding between the public and private parts of the key are verified.
// If the nameAlg of the key is TPM_ALG_NULL, then the size of the sensitive area
// is verified but the public portion is not verified, unless the key is an RSA key.
// For an RSA key, the reason for loading the sensitive area is to use it. The 
// only way to use a private RSA key is to compute the private exponent. To compute
// the private exponent, the public modulus is used.
//  Return Type: TPM_RC
//      TPM_RC_BINDING      the public and private parts are not cryptographically 
//                          bound
//      TPM_RC_HASH         cannot have a publicOnly key with nameAlg of TPM_ALG_NULL
//      TPM_RC_KEY          the public unique is not valid
//      TPM_RC_KEY_SIZE     the private area key is not valid
//      TPM_RC_TYPE         the types of the sensitive and private parts do not match
TPM_RC
CryptValidateKeys(
    TPMT_PUBLIC      *publicArea,
    TPMT_SENSITIVE   *sensitive,
    TPM_RC            blamePublic,
    TPM_RC            blameSensitive
    )
{
    TPM_RC               result;
    UINT16               keySizeInBytes;
    UINT16               digestSize = CryptHashGetDigestSize(publicArea->nameAlg);
    TPMU_PUBLIC_PARMS   *params = &publicArea->parameters;
    TPMU_PUBLIC_ID      *unique = &publicArea->unique;

    if(sensitive != NULL)
    {   
        // Make sure that the types of the public and sensitive are compatible
        if(publicArea->type != sensitive->sensitiveType)
            return TPM_RCS_TYPE + blameSensitive;
        // Make sure that the authValue is not bigger than allowed
        // If there is no name algorithm, then the size just needs to be less than
        // the maximum size of the buffer used for authorization. That size check
        // was made during unmarshaling of the sensitive area
        if((sensitive->authValue.t.size) > digestSize && (digestSize > 0))
            return TPM_RCS_SIZE + blameSensitive;
    }
    switch(publicArea->type)
    {
#if ALG_RSA
        case ALG_RSA_VALUE:
            keySizeInBytes = BITS_TO_BYTES(params->rsaDetail.keyBits);

            // Regardless of whether there is a sensitive area, the public modulus
            // needs to have the correct size. Otherwise, it can't be used for
            // any public key operation nor can it be used to compute the private
            // exponent.
            // NOTE: This implementation only supports key sizes that are multiples 
            // of 1024 bits which means that the MSb of the 0th byte will always be
            // SET in any prime and in the public modulus.
            if((unique->rsa.t.size != keySizeInBytes)
                || (unique->rsa.t.buffer[0] < 0x80))
                return TPM_RCS_KEY + blamePublic;
            if(params->rsaDetail.exponent != 0
               && params->rsaDetail.exponent < 7)
                return TPM_RCS_VALUE + blamePublic;
            if(sensitive != NULL)
            {
                // If there is a sensitive area, it has to be the correct size
                // including having the correct high order bit SET. 
                if(((sensitive->sensitive.rsa.t.size * 2) != keySizeInBytes)
                   || (sensitive->sensitive.rsa.t.buffer[0] < 0x80))
                    return TPM_RCS_KEY_SIZE + blameSensitive;
            }
            break;
#endif
#if ALG_ECC
        case ALG_ECC_VALUE:
        {
            TPMI_ECC_CURVE      curveId;
            curveId = params->eccDetail.curveID;
            keySizeInBytes = BITS_TO_BYTES(CryptEccGetKeySizeForCurve(curveId));
            if(sensitive == NULL)
            {
                // Validate the public key size
                if(unique->ecc.x.t.size != keySizeInBytes
                   || unique->ecc.y.t.size != keySizeInBytes)
                    return TPM_RCS_KEY + blamePublic;
                if(publicArea->nameAlg != ALG_NULL_VALUE)
                {
                    if(!CryptEccIsPointOnCurve(curveId, &unique->ecc))
                        return TPM_RCS_ECC_POINT + blamePublic;
                }
            }
            else
            {
                // If the nameAlg is TPM_ALG_NULL, then only verify that the 
                // private part of the key is OK.
                if(!CryptEccIsValidPrivateKey(&sensitive->sensitive.ecc,
                                                  curveId))
                        return TPM_RCS_KEY_SIZE;
                if(publicArea->nameAlg != ALG_NULL_VALUE)
                {
                // Full key load, verify that the public point belongs to the
                // private key.
                    TPMS_ECC_POINT          toCompare;
                    result = CryptEccPointMultiply(&toCompare, curveId, NULL,
                                                   &sensitive->sensitive.ecc,
                                                   NULL, NULL);
                    if(result != TPM_RC_SUCCESS)
                        return TPM_RCS_BINDING;
                    else
                    {
                    // Make sure that the private key generated the public key.
                    // The input values and the values produced by the point 
                    // multiply may not be the same size so adjust the computed 
                    // value to match the size of the input value by adding or 
                    // removing zeros.
                        AdjustNumberB(&toCompare.x.b, unique->ecc.x.t.size);
                        AdjustNumberB(&toCompare.y.b, unique->ecc.y.t.size);
                        if(!MemoryEqual2B(&unique->ecc.x.b, &toCompare.x.b)
                           || !MemoryEqual2B(&unique->ecc.y.b, &toCompare.y.b))
                            return TPM_RCS_BINDING;
                    }
                }
            }
            break;
        }
#endif
        default:
            // Checks for SYMCIPHER and KEYEDHASH are largely the same
            // If public area has a nameAlg, then validate the public area size
            // and if there is also a sensitive area, validate the binding

            // For consistency, if the object is public-only just make sure that 
            // the unique field is consistent with the name algorithm
            if(sensitive == NULL)
            {
                if(unique->sym.t.size != digestSize)
                    return TPM_RCS_KEY + blamePublic;
            }
            else
            {
                // Make sure that the key size in the sensitive area is consistent.
                if(publicArea->type == ALG_SYMCIPHER_VALUE)
                {
                    result = CryptSymKeyValidate(&params->symDetail.sym, 
                                                &sensitive->sensitive.sym);
                    if(result != TPM_RC_SUCCESS)
                        return result + blameSensitive;
                }
                else
                {
                    // For a keyed hash object, the key has to be less than the
                    // smaller of the block size of the hash used in the scheme or
                    // 128 bytes. The worst case value is limited by the 
                    // unmarshaling code so the only thing left to be checked is 
                    // that it does not exceed the block size of the hash.
                    // by the hash algorithm of the scheme.
                    TPMT_KEYEDHASH_SCHEME       *scheme;
                    UINT16                       maxSize;
                    scheme = &params->keyedHashDetail.scheme;
                    if(scheme->scheme == ALG_XOR_VALUE)
                    {
                        maxSize = CryptHashGetBlockSize(scheme->details.xor.hashAlg);
                    }
                    else if(scheme->scheme == ALG_HMAC_VALUE)
                    {
                        maxSize = CryptHashGetBlockSize(scheme->details.hmac.hashAlg);
                    }
                    else if(scheme->scheme == ALG_NULL_VALUE)
                    {
                        // Not signing or xor so must be a data block
                        maxSize = 128;
                    }
                    else
                        return TPM_RCS_SCHEME + blamePublic;
                    if(sensitive->sensitive.bits.t.size > maxSize)
                        return TPM_RCS_KEY_SIZE + blameSensitive;
                }
                // If there is a nameAlg, check the binding
                if(publicArea->nameAlg != ALG_NULL_VALUE)
                {
                    TPM2B_DIGEST            compare;
                    if(sensitive->seedValue.t.size != digestSize)
                        return TPM_RCS_KEY_SIZE + blameSensitive;

                    CryptComputeSymmetricUnique(publicArea, sensitive, &compare);
                    if(!MemoryEqual2B(&unique->sym.b, &compare.b))
                        return TPM_RC_BINDING;
                }
            }
            break;
    }
    // For a parent, need to check that the seedValue is the correct size for
    // protections. It should be at least half the size of the nameAlg
    if(IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, restricted)
       && IS_ATTRIBUTE(publicArea->objectAttributes, TPMA_OBJECT, decrypt)
       && sensitive != NULL
       && publicArea->nameAlg != ALG_NULL_VALUE)
    {
        if((sensitive->seedValue.t.size < (digestSize / 2))
           || (sensitive->seedValue.t.size > digestSize))
            return TPM_RCS_SIZE + blameSensitive;
    }
    return TPM_RC_SUCCESS;
}

//*** CryptSelectMac()
// This function is used to set the MAC scheme based on the key parameters and
// the input scheme.
//  Return Type: TPM_RC
//      TPM_RC_SCHEME       the scheme is not a valid mac scheme
//      TPM_RC_TYPE         the input key is not a type that supports a mac
//      TPM_RC_VALUE        the input scheme and the key scheme are not compatible
TPM_RC
CryptSelectMac(
    TPMT_PUBLIC             *publicArea,
    TPMI_ALG_MAC_SCHEME     *inMac
)
{
    TPM_ALG_ID              macAlg = ALG_NULL_VALUE;
    switch(publicArea->type)
    {
        case ALG_KEYEDHASH_VALUE:
        {
            // Local value to keep lines from getting too long
            TPMT_KEYEDHASH_SCHEME   *scheme;
            scheme = &publicArea->parameters.keyedHashDetail.scheme;
            // Expect that the scheme is either HMAC or NULL
            if(scheme->scheme != ALG_NULL_VALUE)
                macAlg = scheme->details.hmac.hashAlg;
            break;
        }
        case ALG_SYMCIPHER_VALUE:
        {
            TPMT_SYM_DEF_OBJECT     *scheme;
            scheme = &publicArea->parameters.symDetail.sym;
            // Expect that the scheme is either valid symmetric cipher or NULL
            if(scheme->algorithm != ALG_NULL_VALUE)
                macAlg = scheme->mode.sym;
            break;
        }
        default:
            return TPM_RCS_TYPE;
    }
    // If the input value is not TPM_ALG_NULL ...
    if(*inMac != ALG_NULL_VALUE) 
    {
        // ... then either the scheme in the key must be TPM_ALG_NULL or the input
        // value must match
        if((macAlg != ALG_NULL_VALUE) && (*inMac != macAlg))
            return TPM_RCS_VALUE;
    }
    else
    {
        // Since the input value is TPM_ALG_NULL, then the key value can't be
        // TPM_ALG_NULL
        if(macAlg == ALG_NULL_VALUE)
            return TPM_RCS_VALUE;
        *inMac = macAlg;
    }
    if(!CryptMacIsValidForKey(publicArea->type, *inMac, FALSE))
        return TPM_RCS_SCHEME;
    return TPM_RC_SUCCESS;
}

//*** CryptMacIsValidForKey()
// Check to see if the key type is compatible with the mac type
BOOL
CryptMacIsValidForKey(
    TPM_ALG_ID          keyType,
    TPM_ALG_ID          macAlg,
    BOOL                flag
)
{
    switch(keyType)
    {
        case ALG_KEYEDHASH_VALUE:
            return CryptHashIsValidAlg(macAlg, flag);
            break;
        case ALG_SYMCIPHER_VALUE:
            return CryptSmacIsValidAlg(macAlg, flag);
            break;
        default:
            break;
    }
    return FALSE;
}

//*** CryptSmacIsValidAlg()
// This function is used to test if an algorithm is a supported SMAC algorithm. It
// needs to be updated as new algorithms are added.
BOOL
CryptSmacIsValidAlg(
    TPM_ALG_ID      alg,
    BOOL            FLAG        // IN: Indicates if TPM_ALG_NULL is valid
)
{
    switch (alg)
    {
#if ALG_CMAC
        case ALG_CMAC_VALUE:
            return TRUE;
            break;
#endif
        case ALG_NULL_VALUE:
            return FLAG;
            break;
        default:
            return FALSE;
    }
}

//*** CryptSymModeIsValid()
// Function checks to see if an algorithm ID is a valid, symmetric block cipher 
// mode for the TPM. If 'flag' is SET, them TPM_ALG_NULL is a valid mode.
// not include the modes used for SMAC
BOOL
CryptSymModeIsValid(
    TPM_ALG_ID          mode,
    BOOL                flag
)
{
    switch(mode)
    {
#if ALG_CTR
        case ALG_CTR_VALUE:
#endif // ALG_CTR
#if ALG_OFB
        case ALG_OFB_VALUE:
#endif // ALG_OFB
#if ALG_CBC
        case ALG_CBC_VALUE:
#endif // ALG_CBC
#if ALG_CFB
        case ALG_CFB_VALUE:
#endif // ALG_CFB
#if ALG_ECB
        case ALG_ECB_VALUE:
#endif // ALG_ECB
            return TRUE;
        case ALG_NULL_VALUE:
            return flag;
            break;
        default:
            break;
    }
    return FALSE;
}



