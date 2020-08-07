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
 *  Date: Mar 28, 2019  Time: 08:25:19PM
 */

#ifndef    _CRYPT_UTIL_FP_H_
#define    _CRYPT_UTIL_FP_H_

//*** CryptIsSchemeAnonymous()
// This function is used to test a scheme to see if it is an anonymous scheme
// The only anonymous scheme is ECDAA. ECDAA can be used to do things
// like U-Prove.
BOOL
CryptIsSchemeAnonymous(
    TPM_ALG_ID       scheme         // IN: the scheme algorithm to test
);

//*** ParmDecryptSym()
//  This function performs parameter decryption using symmetric block cipher.
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
);

//*** ParmEncryptSym()
//  This function performs parameter encryption using symmetric block cipher.
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
);

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
);

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
);

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
);

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
);

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
);

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
);

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
);

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
);

//*** CryptComputeSymmetricUnique()
// This function computes the unique field in public area for symmetric objects.
void
CryptComputeSymmetricUnique(
    TPMT_PUBLIC     *publicArea,    // IN: the object's public area
    TPMT_SENSITIVE  *sensitive,     // IN: the associated sensitive area
    TPM2B_DIGEST    *unique         // OUT: unique buffer
);

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
);

//*** CryptGetSignHashAlg()
// Get the hash algorithm of signature from a TPMT_SIGNATURE structure.
// It assumes the signature is not NULL
//  This is a function for easy access
TPMI_ALG_HASH
CryptGetSignHashAlg(
    TPMT_SIGNATURE  *auth           // IN: signature
);

//*** CryptIsSplitSign()
// This function us used to determine if the signing operation is a split
// signing operation that required a TPM2_Commit().
//
BOOL
CryptIsSplitSign(
    TPM_ALG_ID       scheme         // IN: the algorithm selector
);

//*** CryptIsAsymSignScheme()
// This function indicates if a scheme algorithm is a sign algorithm.
BOOL
CryptIsAsymSignScheme(
    TPMI_ALG_PUBLIC          publicType,        // IN: Type of the object
    TPMI_ALG_ASYM_SCHEME     scheme             // IN: the scheme
);

//*** CryptIsAsymDecryptScheme()
// This function indicate if a scheme algorithm is a decrypt algorithm.
BOOL
CryptIsAsymDecryptScheme(
    TPMI_ALG_PUBLIC          publicType,        // IN: Type of the object
    TPMI_ALG_ASYM_SCHEME     scheme             // IN: the scheme
);

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
);

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
);

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
);

//*** CryptGetTestResult
// This function returns the results of a self-test function.
// Note: the behavior in this function is NOT the correct behavior for a real
// TPM implementation.  An artificial behavior is placed here due to the
// limitation of a software simulation environment.  For the correct behavior,
// consult the part 3 specification for TPM2_GetTestResult().
TPM_RC
CryptGetTestResult(
    TPM2B_MAX_BUFFER    *outData        // OUT: test result data
);

//*** CryptIsUniqueSizeValid()
// This function validates that the unique values are consistent.
// NOTE: This is not a comprehensive test of the public key.
//  Return Type: BOOL
//      TRUE(1)         sizes are consistent
//      FALSE(0)        sizes are not consistent
BOOL
CryptIsUniqueSizeValid(
    TPMT_PUBLIC     *publicArea     // IN: the public area to check
);

//*** CryptIsSensitiveSizeValid()
// This function is used by TPM2_LoadExternal() to validate that the sensitive area
// contains a 'sensitive' value that is consistent with the values in the public
// area.
BOOL
CryptIsSensitiveSizeValid(
    TPMT_PUBLIC             *publicArea,        // IN: the object's public part
    TPMT_SENSITIVE          *sensitiveArea      // IN: the object's sensitive part
);

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
);

//*** CryptAlgSetImplemented()
// This function initializes the bit vector with one bit for each implemented
// algorithm. This function is called from _TPM_Init(). The vector of implemented
// algorithms should be generated by the part 2 parser so that the
// 'g_implementedAlgorithms' vector can be a constant. That's not how it is now
void
CryptAlgsSetImplemented(
    void
);

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
);

//*** CryptMacIsValidForKey()
// Check to see if the key type is compatible with the mac type
BOOL
CryptMacIsValidForKey(
    TPM_ALG_ID          keyType,
    TPM_ALG_ID          macAlg,
    BOOL                flag
);

//*** CryptSmacIsValidAlg()
// This function is used to test if an algorithm is a supported SMAC algorithm. It
// needs to be updated as new algorithms are added.
BOOL
CryptSmacIsValidAlg(
    TPM_ALG_ID      alg,
    BOOL            FLAG        // IN: Indicates if TPM_ALG_NULL is valid
);

//*** CryptSymModeIsValid()
// Function checks to see if an algorithm ID is a valid, symmetric block cipher
// mode for the TPM. If 'flag' is SET, them TPM_ALG_NULL is a valid mode.
// not include the modes used for SMAC
BOOL
CryptSymModeIsValid(
    TPM_ALG_ID          mode,
    BOOL                flag
);

#endif  // _CRYPT_UTIL_FP_H_
