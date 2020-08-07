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
 *  Created by TpmStructures; Version 4.4 Mar 26, 2019
 *  Date: Apr 10, 2019  Time: 03:21:33PM
 */

#ifndef _TPM_PROFILE_H_
#define _TPM_PROFILE_H_

// Table 2:4 - Defines for Logic Values
#undef TRUE
#define TRUE                1
#undef FALSE
#define FALSE               0
#undef YES
#define YES                 1
#undef NO
#define NO                  0
#undef SET
#define SET                 1
#undef CLEAR
#define CLEAR               0

// Table 0:1 - Defines for Processor Values
#ifndef BIG_ENDIAN_TPM
#define BIG_ENDIAN_TPM              NO
#endif
#ifndef LITTLE_ENDIAN_TPM
#define LITTLE_ENDIAN_TPM           !BIG_ENDIAN_TPM
#endif
#ifndef MOST_SIGNIFICANT_BIT_0
#define MOST_SIGNIFICANT_BIT_0      NO
#endif
#ifndef LEAST_SIGNIFICANT_BIT_0
#define LEAST_SIGNIFICANT_BIT_0     !MOST_SIGNIFICANT_BIT_0
#endif
#ifndef AUTO_ALIGN
#define AUTO_ALIGN                  NO
#endif

// Table 0:4 - Defines for Implemented Curves
#ifndef ECC_NIST_P192
#define ECC_NIST_P192                   NO
#endif
#ifndef ECC_NIST_P224
#define ECC_NIST_P224                   NO
#endif
#ifndef ECC_NIST_P256
#define ECC_NIST_P256                   YES
#endif
#ifndef ECC_NIST_P384
#define ECC_NIST_P384                   YES
#endif
#ifndef ECC_NIST_P521
#define ECC_NIST_P521                   NO
#endif
#ifndef ECC_BN_P256
#define ECC_BN_P256                     YES
#endif
#ifndef ECC_BN_P638
#define ECC_BN_P638                     NO
#endif
#ifndef ECC_SM2_P256
#define ECC_SM2_P256                    NO
#endif

// Table 0:7 - Defines for Implementation Values
#ifndef FIELD_UPGRADE_IMPLEMENTED
#define FIELD_UPGRADE_IMPLEMENTED       NO
#endif
#ifndef HASH_ALIGNMENT
#define HASH_ALIGNMENT                  4
#endif
#ifndef SYMMETRIC_ALIGNMENT
#define SYMMETRIC_ALIGNMENT             4
#endif
#ifndef HASH_LIB
#define HASH_LIB                        Ossl
#endif
#ifndef SYM_LIB
#define SYM_LIB                         Ossl
#endif
#ifndef MATH_LIB
#define MATH_LIB                        Ossl
#endif
#ifndef BSIZE
#define BSIZE                           UINT16
#endif
#ifndef IMPLEMENTATION_PCR
#define IMPLEMENTATION_PCR              24
#endif
#ifndef PCR_SELECT_MAX
#define PCR_SELECT_MAX                  ((IMPLEMENTATION_PCR+7)/8)
#endif
#ifndef PLATFORM_PCR
#define PLATFORM_PCR                    24
#endif
#ifndef PCR_SELECT_MIN
#define PCR_SELECT_MIN                  ((PLATFORM_PCR+7)/8)
#endif
#ifndef DRTM_PCR
#define DRTM_PCR                        17
#endif
#ifndef HCRTM_PCR
#define HCRTM_PCR                       0
#endif
#ifndef NUM_LOCALITIES
#define NUM_LOCALITIES                  5
#endif
#ifndef MAX_HANDLE_NUM
#define MAX_HANDLE_NUM                  3
#endif
#ifndef MAX_ACTIVE_SESSIONS
#define MAX_ACTIVE_SESSIONS             64
#endif
#ifndef CONTEXT_SLOT
#define CONTEXT_SLOT                    UINT16
#endif
#ifndef CONTEXT_COUNTER
#define CONTEXT_COUNTER                 UINT64
#endif
#ifndef MAX_LOADED_SESSIONS
#define MAX_LOADED_SESSIONS             3
#endif
#ifndef MAX_SESSION_NUM
#define MAX_SESSION_NUM                 3
#endif
#ifndef MAX_LOADED_OBJECTS
#define MAX_LOADED_OBJECTS              3
#endif
#ifndef MIN_EVICT_OBJECTS
#define MIN_EVICT_OBJECTS               2
#endif
#ifndef NUM_POLICY_PCR_GROUP
#define NUM_POLICY_PCR_GROUP            1
#endif
#ifndef NUM_AUTHVALUE_PCR_GROUP
#define NUM_AUTHVALUE_PCR_GROUP         1
#endif
#ifndef MAX_CONTEXT_SIZE
#define MAX_CONTEXT_SIZE                1264
#endif
#ifndef MAX_DIGEST_BUFFER
#define MAX_DIGEST_BUFFER               1024
#endif
#ifndef MAX_NV_INDEX_SIZE
#define MAX_NV_INDEX_SIZE               2048
#endif
#ifndef MAX_NV_BUFFER_SIZE
#define MAX_NV_BUFFER_SIZE              1024
#endif
#ifndef MAX_CAP_BUFFER
#define MAX_CAP_BUFFER                  1024
#endif
#ifndef NV_MEMORY_SIZE
#define NV_MEMORY_SIZE                  16384
#endif
#ifndef MIN_COUNTER_INDICES
#define MIN_COUNTER_INDICES             8
#endif
#ifndef NUM_STATIC_PCR
#define NUM_STATIC_PCR                  16
#endif
#ifndef MAX_ALG_LIST_SIZE
#define MAX_ALG_LIST_SIZE               64
#endif
#ifndef PRIMARY_SEED_SIZE
#define PRIMARY_SEED_SIZE               32
#endif
#ifndef CONTEXT_ENCRYPT_ALGORITHM
#define CONTEXT_ENCRYPT_ALGORITHM       AES
#endif
#ifndef NV_CLOCK_UPDATE_INTERVAL
#define NV_CLOCK_UPDATE_INTERVAL        12
#endif
#ifndef NUM_POLICY_PCR
#define NUM_POLICY_PCR                  1
#endif
#ifndef MAX_COMMAND_SIZE
#define MAX_COMMAND_SIZE                4096
#endif
#ifndef MAX_RESPONSE_SIZE
#define MAX_RESPONSE_SIZE               4096
#endif
#ifndef ORDERLY_BITS
#define ORDERLY_BITS                    8
#endif
#ifndef MAX_SYM_DATA
#define MAX_SYM_DATA                    128
#endif
#ifndef MAX_RNG_ENTROPY_SIZE
#define MAX_RNG_ENTROPY_SIZE            64
#endif
#ifndef RAM_INDEX_SPACE
#define RAM_INDEX_SPACE                 512
#endif
#ifndef RSA_DEFAULT_PUBLIC_EXPONENT
#define RSA_DEFAULT_PUBLIC_EXPONENT     0x00010001
#endif
#ifndef ENABLE_PCR_NO_INCREMENT
#define ENABLE_PCR_NO_INCREMENT         YES
#endif
#ifndef CRT_FORMAT_RSA
#define CRT_FORMAT_RSA                  YES
#endif
#ifndef VENDOR_COMMAND_COUNT
#define VENDOR_COMMAND_COUNT            0
#endif
#ifndef MAX_VENDOR_BUFFER_SIZE
#define MAX_VENDOR_BUFFER_SIZE          1024
#endif
#ifndef TPM_MAX_DERIVATION_BITS
#define TPM_MAX_DERIVATION_BITS         8192
#endif
#ifndef RSA_MAX_PRIME
#define RSA_MAX_PRIME                   (MAX_RSA_KEY_BYTES/2)
#endif
#ifndef RSA_PRIVATE_SIZE
#define RSA_PRIVATE_SIZE                (RSA_MAX_PRIME*5)
#endif
#ifndef SIZE_OF_X509_SERIAL_NUMBER
#define SIZE_OF_X509_SERIAL_NUMBER      20
#endif
#ifndef PRIVATE_VENDOR_SPECIFIC_BYTES
#define PRIVATE_VENDOR_SPECIFIC_BYTES   RSA_PRIVATE_SIZE
#endif

// Table 0:2 - Defines for Implemented Algorithms
#ifndef ALG_AES
#define ALG_AES                         ALG_YES
#endif
#ifndef ALG_CAMELLIA
#define ALG_CAMELLIA                    ALG_NO      /* Not specified by vendor */
#endif
#ifndef ALG_CBC
#define ALG_CBC                         ALG_YES
#endif
#ifndef ALG_CFB
#define ALG_CFB                         ALG_YES
#endif
#ifndef ALG_CMAC
#define ALG_CMAC                        ALG_YES
#endif
#ifndef ALG_CTR
#define ALG_CTR                         ALG_YES
#endif
#ifndef ALG_ECB
#define ALG_ECB                         ALG_YES
#endif
#ifndef ALG_ECC
#define ALG_ECC                         ALG_YES
#endif
#ifndef ALG_ECDAA
#define ALG_ECDAA                       (ALG_YES && ALG_ECC)
#endif
#ifndef ALG_ECDH
#define ALG_ECDH                        (ALG_YES && ALG_ECC)
#endif
#ifndef ALG_ECDSA
#define ALG_ECDSA                       (ALG_YES && ALG_ECC)
#endif
#ifndef ALG_ECMQV
#define ALG_ECMQV                       (ALG_NO && ALG_ECC)
#endif
#ifndef ALG_ECSCHNORR
#define ALG_ECSCHNORR                   (ALG_YES && ALG_ECC)
#endif
#ifndef ALG_HMAC
#define ALG_HMAC                        ALG_YES
#endif
#ifndef ALG_KDF1_SP800_108
#define ALG_KDF1_SP800_108              ALG_YES
#endif
#ifndef ALG_KDF1_SP800_56A
#define ALG_KDF1_SP800_56A              (ALG_YES && ALG_ECC)
#endif
#ifndef ALG_KDF2
#define ALG_KDF2                        ALG_NO
#endif
#ifndef ALG_KEYEDHASH
#define ALG_KEYEDHASH                   ALG_YES
#endif
#ifndef ALG_MGF1
#define ALG_MGF1                        ALG_YES
#endif
#ifndef ALG_OAEP
#define ALG_OAEP                        (ALG_YES && ALG_RSA)
#endif
#ifndef ALG_OFB
#define ALG_OFB                         ALG_YES
#endif
#ifndef ALG_RSA
#define ALG_RSA                         ALG_YES
#endif
#ifndef ALG_RSAES
#define ALG_RSAES                       (ALG_YES && ALG_RSA)
#endif
#ifndef ALG_RSAPSS
#define ALG_RSAPSS                      (ALG_YES && ALG_RSA)
#endif
#ifndef ALG_RSASSA
#define ALG_RSASSA                      (ALG_YES && ALG_RSA)
#endif
#ifndef ALG_SHA
#define ALG_SHA                         ALG_NO      /* Not specified by vendor */
#endif
#ifndef ALG_SHA1
#define ALG_SHA1                        ALG_YES
#endif
#ifndef ALG_SHA256
#define ALG_SHA256                      ALG_YES
#endif
#ifndef ALG_SHA384
#define ALG_SHA384                      ALG_YES
#endif
#ifndef ALG_SHA3_256
#define ALG_SHA3_256                    ALG_NO      /* Not specified by vendor */
#endif
#ifndef ALG_SHA3_384
#define ALG_SHA3_384                    ALG_NO      /* Not specified by vendor */
#endif
#ifndef ALG_SHA3_512
#define ALG_SHA3_512                    ALG_NO      /* Not specified by vendor */
#endif
#ifndef ALG_SHA512
#define ALG_SHA512                      ALG_NO
#endif
#ifndef ALG_SM2
#define ALG_SM2                         (ALG_NO && ALG_ECC)
#endif
#ifndef ALG_SM3_256
#define ALG_SM3_256                     ALG_NO
#endif
#ifndef ALG_SM4
#define ALG_SM4                         ALG_NO
#endif
#ifndef ALG_SYMCIPHER
#define ALG_SYMCIPHER                   ALG_YES
#endif
#ifndef ALG_TDES
#define ALG_TDES                        ALG_NO
#endif
#ifndef ALG_XOR
#define ALG_XOR                         ALG_YES
#endif

// Table 1:00 - Defines for RSA Asymmetric Cipher Algorithm Constants
#ifndef RSA_1024
#define RSA_1024                    (ALG_RSA & YES)
#endif
#ifndef RSA_2048
#define RSA_2048                    (ALG_RSA & YES)
#endif
#ifndef RSA_3072
#define RSA_3072                    (ALG_RSA & NO)
#endif
#ifndef RSA_4096
#define RSA_4096                    (ALG_RSA & NO)
#endif

// Table 1:17 - Defines for AES Symmetric Cipher Algorithm Constants
#ifndef AES_128
#define AES_128                     (ALG_AES & YES)
#endif
#ifndef AES_192
#define AES_192                     (ALG_AES & NO)
#endif
#ifndef AES_256
#define AES_256                     (ALG_AES & YES)
#endif

// Table 1:18 - Defines for SM4 Symmetric Cipher Algorithm Constants
#ifndef SM4_128
#define SM4_128                     (ALG_SM4 & YES)
#endif

// Table 1:19 - Defines for CAMELLIA Symmetric Cipher Algorithm Constants
#ifndef CAMELLIA_128
#define CAMELLIA_128                    (ALG_CAMELLIA & YES)
#endif
#ifndef CAMELLIA_192
#define CAMELLIA_192                    (ALG_CAMELLIA & NO)
#endif
#ifndef CAMELLIA_256
#define CAMELLIA_256                    (ALG_CAMELLIA & NO)
#endif

// Table 1:17 - Defines for TDES Symmetric Cipher Algorithm Constants
#ifndef TDES_128
#define TDES_128                    (ALG_TDES & YES)
#endif
#ifndef TDES_192
#define TDES_192                    (ALG_TDES & YES)
#endif

// Table 0:5 - Defines for Implemented Commands
#ifndef CC_AC_GetCapability
#define CC_AC_GetCapability                 CC_YES
#endif
#ifndef CC_AC_Send
#define CC_AC_Send                          CC_YES
#endif
#ifndef CC_ActivateCredential
#define CC_ActivateCredential               CC_YES
#endif
#ifndef CC_Certify
#define CC_Certify                          CC_YES
#endif
#ifndef CC_CertifyCreation
#define CC_CertifyCreation                  CC_YES
#endif
#ifndef CC_CertifyX509
#define CC_CertifyX509                      CC_YES
#endif
#ifndef CC_ChangeEPS
#define CC_ChangeEPS                        CC_YES
#endif
#ifndef CC_ChangePPS
#define CC_ChangePPS                        CC_YES
#endif
#ifndef CC_Clear
#define CC_Clear                            CC_YES
#endif
#ifndef CC_ClearControl
#define CC_ClearControl                     CC_YES
#endif
#ifndef CC_ClockRateAdjust
#define CC_ClockRateAdjust                  CC_YES
#endif
#ifndef CC_ClockSet
#define CC_ClockSet                         CC_YES
#endif
#ifndef CC_Commit
#define CC_Commit                           (CC_YES && ALG_ECC)
#endif
#ifndef CC_ContextLoad
#define CC_ContextLoad                      CC_YES
#endif
#ifndef CC_ContextSave
#define CC_ContextSave                      CC_YES
#endif
#ifndef CC_Create
#define CC_Create                           CC_YES
#endif
#ifndef CC_CreateLoaded
#define CC_CreateLoaded                     CC_YES
#endif
#ifndef CC_CreatePrimary
#define CC_CreatePrimary                    CC_YES
#endif
#ifndef CC_DictionaryAttackLockReset
#define CC_DictionaryAttackLockReset        CC_YES
#endif
#ifndef CC_DictionaryAttackParameters
#define CC_DictionaryAttackParameters       CC_YES
#endif
#ifndef CC_Duplicate
#define CC_Duplicate                        CC_YES
#endif
#ifndef CC_ECC_Parameters
#define CC_ECC_Parameters                   (CC_YES && ALG_ECC)
#endif
#ifndef CC_ECDH_KeyGen
#define CC_ECDH_KeyGen                      (CC_YES && ALG_ECC)
#endif
#ifndef CC_ECDH_ZGen
#define CC_ECDH_ZGen                        (CC_YES && ALG_ECC)
#endif
#ifndef CC_EC_Ephemeral
#define CC_EC_Ephemeral                     (CC_YES && ALG_ECC)
#endif
#ifndef CC_EncryptDecrypt
#define CC_EncryptDecrypt                   CC_YES
#endif
#ifndef CC_EncryptDecrypt2
#define CC_EncryptDecrypt2                  CC_YES
#endif
#ifndef CC_EventSequenceComplete
#define CC_EventSequenceComplete            CC_YES
#endif
#ifndef CC_EvictControl
#define CC_EvictControl                     CC_YES
#endif
#ifndef CC_FieldUpgradeData
#define CC_FieldUpgradeData                 CC_NO
#endif
#ifndef CC_FieldUpgradeStart
#define CC_FieldUpgradeStart                CC_NO
#endif
#ifndef CC_FirmwareRead
#define CC_FirmwareRead                     CC_NO
#endif
#ifndef CC_FlushContext
#define CC_FlushContext                     CC_YES
#endif
#ifndef CC_GetCapability
#define CC_GetCapability                    CC_YES
#endif
#ifndef CC_GetCommandAuditDigest
#define CC_GetCommandAuditDigest            CC_YES
#endif
#ifndef CC_GetRandom
#define CC_GetRandom                        CC_YES
#endif
#ifndef CC_GetSessionAuditDigest
#define CC_GetSessionAuditDigest            CC_YES
#endif
#ifndef CC_GetTestResult
#define CC_GetTestResult                    CC_YES
#endif
#ifndef CC_GetTime
#define CC_GetTime                          CC_YES
#endif
#ifndef CC_HMAC
#define CC_HMAC                             (CC_YES && !ALG_CMAC)
#endif
#ifndef CC_HMAC_Start
#define CC_HMAC_Start                       (CC_YES && !ALG_CMAC)
#endif
#ifndef CC_Hash
#define CC_Hash                             CC_YES
#endif
#ifndef CC_HashSequenceStart
#define CC_HashSequenceStart                CC_YES
#endif
#ifndef CC_HierarchyChangeAuth
#define CC_HierarchyChangeAuth              CC_YES
#endif
#ifndef CC_HierarchyControl
#define CC_HierarchyControl                 CC_YES
#endif
#ifndef CC_Import
#define CC_Import                           CC_YES
#endif
#ifndef CC_IncrementalSelfTest
#define CC_IncrementalSelfTest              CC_YES
#endif
#ifndef CC_Load
#define CC_Load                             CC_YES
#endif
#ifndef CC_LoadExternal
#define CC_LoadExternal                     CC_YES
#endif
#ifndef CC_MAC
#define CC_MAC                              (CC_YES && ALG_CMAC)
#endif
#ifndef CC_MAC_Start
#define CC_MAC_Start                        (CC_YES && ALG_CMAC)
#endif
#ifndef CC_MakeCredential
#define CC_MakeCredential                   CC_YES
#endif
#ifndef CC_NV_Certify
#define CC_NV_Certify                       CC_YES
#endif
#ifndef CC_NV_ChangeAuth
#define CC_NV_ChangeAuth                    CC_YES
#endif
#ifndef CC_NV_DefineSpace
#define CC_NV_DefineSpace                   CC_YES
#endif
#ifndef CC_NV_Extend
#define CC_NV_Extend                        CC_YES
#endif
#ifndef CC_NV_GlobalWriteLock
#define CC_NV_GlobalWriteLock               CC_YES
#endif
#ifndef CC_NV_Increment
#define CC_NV_Increment                     CC_YES
#endif
#ifndef CC_NV_Read
#define CC_NV_Read                          CC_YES
#endif
#ifndef CC_NV_ReadLock
#define CC_NV_ReadLock                      CC_YES
#endif
#ifndef CC_NV_ReadPublic
#define CC_NV_ReadPublic                    CC_YES
#endif
#ifndef CC_NV_SetBits
#define CC_NV_SetBits                       CC_YES
#endif
#ifndef CC_NV_UndefineSpace
#define CC_NV_UndefineSpace                 CC_YES
#endif
#ifndef CC_NV_UndefineSpaceSpecial
#define CC_NV_UndefineSpaceSpecial          CC_YES
#endif
#ifndef CC_NV_Write
#define CC_NV_Write                         CC_YES
#endif
#ifndef CC_NV_WriteLock
#define CC_NV_WriteLock                     CC_YES
#endif
#ifndef CC_ObjectChangeAuth
#define CC_ObjectChangeAuth                 CC_YES
#endif
#ifndef CC_PCR_Allocate
#define CC_PCR_Allocate                     CC_YES
#endif
#ifndef CC_PCR_Event
#define CC_PCR_Event                        CC_YES
#endif
#ifndef CC_PCR_Extend
#define CC_PCR_Extend                       CC_YES
#endif
#ifndef CC_PCR_Read
#define CC_PCR_Read                         CC_YES
#endif
#ifndef CC_PCR_Reset
#define CC_PCR_Reset                        CC_YES
#endif
#ifndef CC_PCR_SetAuthPolicy
#define CC_PCR_SetAuthPolicy                CC_YES
#endif
#ifndef CC_PCR_SetAuthValue
#define CC_PCR_SetAuthValue                 CC_YES
#endif
#ifndef CC_PP_Commands
#define CC_PP_Commands                      CC_YES
#endif
#ifndef CC_PolicyAuthValue
#define CC_PolicyAuthValue                  CC_YES
#endif
#ifndef CC_PolicyAuthorize
#define CC_PolicyAuthorize                  CC_YES
#endif
#ifndef CC_PolicyAuthorizeNV
#define CC_PolicyAuthorizeNV                CC_YES
#endif
#ifndef CC_PolicyCommandCode
#define CC_PolicyCommandCode                CC_YES
#endif
#ifndef CC_PolicyCounterTimer
#define CC_PolicyCounterTimer               CC_YES
#endif
#ifndef CC_PolicyCpHash
#define CC_PolicyCpHash                     CC_YES
#endif
#ifndef CC_PolicyDuplicationSelect
#define CC_PolicyDuplicationSelect          CC_YES
#endif
#ifndef CC_PolicyGetDigest
#define CC_PolicyGetDigest                  CC_YES
#endif
#ifndef CC_PolicyLocality
#define CC_PolicyLocality                   CC_YES
#endif
#ifndef CC_PolicyNV
#define CC_PolicyNV                         CC_YES
#endif
#ifndef CC_PolicyNameHash
#define CC_PolicyNameHash                   CC_YES
#endif
#ifndef CC_PolicyNvWritten
#define CC_PolicyNvWritten                  CC_YES
#endif
#ifndef CC_PolicyOR
#define CC_PolicyOR                         CC_YES
#endif
#ifndef CC_PolicyPCR
#define CC_PolicyPCR                        CC_YES
#endif
#ifndef CC_PolicyPassword
#define CC_PolicyPassword                   CC_YES
#endif
#ifndef CC_PolicyPhysicalPresence
#define CC_PolicyPhysicalPresence           CC_YES
#endif
#ifndef CC_PolicyRestart
#define CC_PolicyRestart                    CC_YES
#endif
#ifndef CC_PolicySecret
#define CC_PolicySecret                     CC_YES
#endif
#ifndef CC_PolicySigned
#define CC_PolicySigned                     CC_YES
#endif
#ifndef CC_PolicyTemplate
#define CC_PolicyTemplate                   CC_YES
#endif
#ifndef CC_PolicyTicket
#define CC_PolicyTicket                     CC_YES
#endif
#ifndef CC_Policy_AC_SendSelect
#define CC_Policy_AC_SendSelect             CC_YES
#endif
#ifndef CC_Quote
#define CC_Quote                            CC_YES
#endif
#ifndef CC_RSA_Decrypt
#define CC_RSA_Decrypt                      (CC_YES && ALG_RSA)
#endif
#ifndef CC_RSA_Encrypt
#define CC_RSA_Encrypt                      (CC_YES && ALG_RSA)
#endif
#ifndef CC_ReadClock
#define CC_ReadClock                        CC_YES
#endif
#ifndef CC_ReadPublic
#define CC_ReadPublic                       CC_YES
#endif
#ifndef CC_Rewrap
#define CC_Rewrap                           CC_YES
#endif
#ifndef CC_SelfTest
#define CC_SelfTest                         CC_YES
#endif
#ifndef CC_SequenceComplete
#define CC_SequenceComplete                 CC_YES
#endif
#ifndef CC_SequenceUpdate
#define CC_SequenceUpdate                   CC_YES
#endif
#ifndef CC_SetAlgorithmSet
#define CC_SetAlgorithmSet                  CC_YES
#endif
#ifndef CC_SetCommandCodeAuditStatus
#define CC_SetCommandCodeAuditStatus        CC_YES
#endif
#ifndef CC_SetPrimaryPolicy
#define CC_SetPrimaryPolicy                 CC_YES
#endif
#ifndef CC_Shutdown
#define CC_Shutdown                         CC_YES
#endif
#ifndef CC_Sign
#define CC_Sign                             CC_YES
#endif
#ifndef CC_StartAuthSession
#define CC_StartAuthSession                 CC_YES
#endif
#ifndef CC_Startup
#define CC_Startup                          CC_YES
#endif
#ifndef CC_StirRandom
#define CC_StirRandom                       CC_YES
#endif
#ifndef CC_TestParms
#define CC_TestParms                        CC_YES
#endif
#ifndef CC_Unseal
#define CC_Unseal                           CC_YES
#endif
#ifndef CC_Vendor_TCG_Test
#define CC_Vendor_TCG_Test                  CC_YES
#endif
#ifndef CC_VerifySignature
#define CC_VerifySignature                  CC_YES
#endif
#ifndef CC_ZGen_2Phase
#define CC_ZGen_2Phase                      (CC_YES && ALG_ECC)
#endif


#endif // _TPM_PROFILE_H_
