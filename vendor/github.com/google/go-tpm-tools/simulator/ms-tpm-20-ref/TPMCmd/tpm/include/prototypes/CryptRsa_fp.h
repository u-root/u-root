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
 *  Date: Apr  2, 2019  Time: 03:18:00PM
 */

#ifndef    _CRYPT_RSA_FP_H_
#define    _CRYPT_RSA_FP_H_

#if ALG_RSA

//*** CryptRsaInit()
// Function called at _TPM_Init().
BOOL
CryptRsaInit(
    void
);

//*** CryptRsaStartup()
// Function called at TPM2_Startup()
BOOL
CryptRsaStartup(
    void
);

//*** CryptRsaPssSaltSize()
// This function computes the salt size used in PSS. It is broken out so that
// the X509 code can get the same value that is used by the encoding function in this
// module.
INT16
CryptRsaPssSaltSize(
    INT16              hashSize,
    INT16               outSize
);

//*** MakeDerTag()
// Construct the DER value that is used in RSASSA
//  Return Type: INT16
//   > 0        size of value
//   <= 0       no hash exists
INT16
MakeDerTag(
    TPM_ALG_ID   hashAlg,
    INT16        sizeOfBuffer,
    BYTE        *buffer
);

//*** CryptRsaSelectScheme()
// This function is used by TPM2_RSA_Decrypt and TPM2_RSA_Encrypt.  It sets up
// the rules to select a scheme between input and object default.
// This function assume the RSA object is loaded.
// If a default scheme is defined in object, the default scheme should be chosen,
// otherwise, the input scheme should be chosen.
// In the case that both the object and 'scheme' are not TPM_ALG_NULL, then
// if the schemes are the same, the input scheme will be chosen.
// if the scheme are not compatible, a NULL pointer will be returned.
//
// The return pointer may point to a TPM_ALG_NULL scheme.
TPMT_RSA_DECRYPT*
CryptRsaSelectScheme(
    TPMI_DH_OBJECT       rsaHandle,     // IN: handle of an RSA key
    TPMT_RSA_DECRYPT    *scheme         // IN: a sign or decrypt scheme
);

//*** CryptRsaLoadPrivateExponent()
// This function is called to generate the private exponent of an RSA key.
//  Return Type: TPM_RC
//      TPM_RC_BINDING      public and private parts of 'rsaKey' are not matched
TPM_RC
CryptRsaLoadPrivateExponent(
    TPMT_PUBLIC             *publicArea,
    TPMT_SENSITIVE          *sensitive
);

//*** CryptRsaEncrypt()
// This is the entry point for encryption using RSA. Encryption is
// use of the public exponent. The padding parameter determines what
// padding will be used.
//
// The 'cOutSize' parameter must be at least as large as the size of the key.
//
// If the padding is RSA_PAD_NONE, 'dIn' is treated as a number. It must be
// lower in value than the key modulus.
// NOTE: If dIn has fewer bytes than cOut, then we don't add low-order zeros to
//       dIn to make it the size of the RSA key for the call to RSAEP. This is
//       because the high order bytes of dIn might have a numeric value that is
//       greater than the value of the key modulus. If this had low-order zeros
//       added, it would have a numeric value larger than the modulus even though
//       it started out with a lower numeric value.
//
//  Return Type: TPM_RC
//      TPM_RC_VALUE     'cOutSize' is too small (must be the size
//                        of the modulus)
//      TPM_RC_SCHEME    'padType' is not a supported scheme
//
LIB_EXPORT TPM_RC
CryptRsaEncrypt(
    TPM2B_PUBLIC_KEY_RSA        *cOut,          // OUT: the encrypted data
    TPM2B                       *dIn,           // IN: the data to encrypt
    OBJECT                      *key,           // IN: the key used for encryption
    TPMT_RSA_DECRYPT            *scheme,        // IN: the type of padding and hash
                                                //     if needed
    const TPM2B                 *label,         // IN: in case it is needed
    RAND_STATE                  *rand           // IN: random number generator
                                                //     state (mostly for testing)
);

//*** CryptRsaDecrypt()
// This is the entry point for decryption using RSA. Decryption is
// use of the private exponent. The 'padType' parameter determines what
// padding was used.
//
//  Return Type: TPM_RC
//      TPM_RC_SIZE        'cInSize' is not the same as the size of the public
//                          modulus of 'key'; or numeric value of the encrypted
//                          data is greater than the modulus
//      TPM_RC_VALUE       'dOutSize' is not large enough for the result
//      TPM_RC_SCHEME      'padType' is not supported
//
LIB_EXPORT TPM_RC
CryptRsaDecrypt(
    TPM2B               *dOut,          // OUT: the decrypted data
    TPM2B               *cIn,           // IN: the data to decrypt
    OBJECT              *key,           // IN: the key to use for decryption
    TPMT_RSA_DECRYPT    *scheme,        // IN: the padding scheme
    const TPM2B         *label          // IN: in case it is needed for the scheme
);

//*** CryptRsaSign()
// This function is used to generate an RSA signature of the type indicated in
// 'scheme'.
//
//  Return Type: TPM_RC
//      TPM_RC_SCHEME       'scheme' or 'hashAlg' are not supported
//      TPM_RC_VALUE        'hInSize' does not match 'hashAlg' (for RSASSA)
//
LIB_EXPORT TPM_RC
CryptRsaSign(
    TPMT_SIGNATURE      *sigOut,
    OBJECT              *key,           // IN: key to use
    TPM2B_DIGEST        *hIn,           // IN: the digest to sign
    RAND_STATE          *rand           // IN: the random number generator
                                        //      to use (mostly for testing)
);

//*** CryptRsaValidateSignature()
// This function is used to validate an RSA signature. If the signature is valid
// TPM_RC_SUCCESS is returned. If the signature is not valid, TPM_RC_SIGNATURE is
// returned. Other return codes indicate either parameter problems or fatal errors.
//
//  Return Type: TPM_RC
//      TPM_RC_SIGNATURE    the signature does not check
//      TPM_RC_SCHEME       unsupported scheme or hash algorithm
//
LIB_EXPORT TPM_RC
CryptRsaValidateSignature(
    TPMT_SIGNATURE  *sig,           // IN: signature
    OBJECT          *key,           // IN: public modulus
    TPM2B_DIGEST    *digest         // IN: The digest being validated
);

//*** CryptRsaGenerateKey()
// Generate an RSA key from a provided seed
//  Return Type: TPM_RC
//      TPM_RC_CANCELED     operation was canceled
//      TPM_RC_RANGE        public exponent is not supported
//      TPM_RC_VALUE        could not find a prime using the provided parameters
LIB_EXPORT TPM_RC
CryptRsaGenerateKey(
    TPMT_PUBLIC         *publicArea,
    TPMT_SENSITIVE      *sensitive,
    RAND_STATE          *rand               // IN: if not NULL, the deterministic
                                            //     RNG state
);
#endif // ALG_RSA

#endif  // _CRYPT_RSA_FP_H_
