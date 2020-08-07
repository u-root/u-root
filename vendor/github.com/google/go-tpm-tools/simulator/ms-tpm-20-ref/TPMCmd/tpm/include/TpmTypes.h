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

#ifndef _TPM_TYPES_H_
#define _TPM_TYPES_H_

// Table 1:2 - Definition of TPM_ALG_ID Constants
typedef UINT16                          TPM_ALG_ID;
#define TYPE_OF_TPM_ALG_ID              UINT16
#define     ALG_ERROR_VALUE             0x0000
#define TPM_ALG_ERROR                   (TPM_ALG_ID)(ALG_ERROR_VALUE)
#define     ALG_RSA_VALUE               0x0001
#define TPM_ALG_RSA                     (TPM_ALG_ID)(ALG_RSA_VALUE)
#define     ALG_TDES_VALUE              0x0003
#define TPM_ALG_TDES                    (TPM_ALG_ID)(ALG_TDES_VALUE)
#define     ALG_SHA_VALUE               0x0004
#define TPM_ALG_SHA                     (TPM_ALG_ID)(ALG_SHA_VALUE)
#define     ALG_SHA1_VALUE              0x0004
#define TPM_ALG_SHA1                    (TPM_ALG_ID)(ALG_SHA1_VALUE)
#define     ALG_HMAC_VALUE              0x0005
#define TPM_ALG_HMAC                    (TPM_ALG_ID)(ALG_HMAC_VALUE)
#define     ALG_AES_VALUE               0x0006
#define TPM_ALG_AES                     (TPM_ALG_ID)(ALG_AES_VALUE)
#define     ALG_MGF1_VALUE              0x0007
#define TPM_ALG_MGF1                    (TPM_ALG_ID)(ALG_MGF1_VALUE)
#define     ALG_KEYEDHASH_VALUE         0x0008
#define TPM_ALG_KEYEDHASH               (TPM_ALG_ID)(ALG_KEYEDHASH_VALUE)
#define     ALG_XOR_VALUE               0x000A
#define TPM_ALG_XOR                     (TPM_ALG_ID)(ALG_XOR_VALUE)
#define     ALG_SHA256_VALUE            0x000B
#define TPM_ALG_SHA256                  (TPM_ALG_ID)(ALG_SHA256_VALUE)
#define     ALG_SHA384_VALUE            0x000C
#define TPM_ALG_SHA384                  (TPM_ALG_ID)(ALG_SHA384_VALUE)
#define     ALG_SHA512_VALUE            0x000D
#define TPM_ALG_SHA512                  (TPM_ALG_ID)(ALG_SHA512_VALUE)
#define     ALG_NULL_VALUE              0x0010
#define TPM_ALG_NULL                    (TPM_ALG_ID)(ALG_NULL_VALUE)
#define     ALG_SM3_256_VALUE           0x0012
#define TPM_ALG_SM3_256                 (TPM_ALG_ID)(ALG_SM3_256_VALUE)
#define     ALG_SM4_VALUE               0x0013
#define TPM_ALG_SM4                     (TPM_ALG_ID)(ALG_SM4_VALUE)
#define     ALG_RSASSA_VALUE            0x0014
#define TPM_ALG_RSASSA                  (TPM_ALG_ID)(ALG_RSASSA_VALUE)
#define     ALG_RSAES_VALUE             0x0015
#define TPM_ALG_RSAES                   (TPM_ALG_ID)(ALG_RSAES_VALUE)
#define     ALG_RSAPSS_VALUE            0x0016
#define TPM_ALG_RSAPSS                  (TPM_ALG_ID)(ALG_RSAPSS_VALUE)
#define     ALG_OAEP_VALUE              0x0017
#define TPM_ALG_OAEP                    (TPM_ALG_ID)(ALG_OAEP_VALUE)
#define     ALG_ECDSA_VALUE             0x0018
#define TPM_ALG_ECDSA                   (TPM_ALG_ID)(ALG_ECDSA_VALUE)
#define     ALG_ECDH_VALUE              0x0019
#define TPM_ALG_ECDH                    (TPM_ALG_ID)(ALG_ECDH_VALUE)
#define     ALG_ECDAA_VALUE             0x001A
#define TPM_ALG_ECDAA                   (TPM_ALG_ID)(ALG_ECDAA_VALUE)
#define     ALG_SM2_VALUE               0x001B
#define TPM_ALG_SM2                     (TPM_ALG_ID)(ALG_SM2_VALUE)
#define     ALG_ECSCHNORR_VALUE         0x001C
#define TPM_ALG_ECSCHNORR               (TPM_ALG_ID)(ALG_ECSCHNORR_VALUE)
#define     ALG_ECMQV_VALUE             0x001D
#define TPM_ALG_ECMQV                   (TPM_ALG_ID)(ALG_ECMQV_VALUE)
#define     ALG_KDF1_SP800_56A_VALUE    0x0020
#define TPM_ALG_KDF1_SP800_56A          (TPM_ALG_ID)(ALG_KDF1_SP800_56A_VALUE)
#define     ALG_KDF2_VALUE              0x0021
#define TPM_ALG_KDF2                    (TPM_ALG_ID)(ALG_KDF2_VALUE)
#define     ALG_KDF1_SP800_108_VALUE    0x0022
#define TPM_ALG_KDF1_SP800_108          (TPM_ALG_ID)(ALG_KDF1_SP800_108_VALUE)
#define     ALG_ECC_VALUE               0x0023
#define TPM_ALG_ECC                     (TPM_ALG_ID)(ALG_ECC_VALUE)
#define     ALG_SYMCIPHER_VALUE         0x0025
#define TPM_ALG_SYMCIPHER               (TPM_ALG_ID)(ALG_SYMCIPHER_VALUE)
#define     ALG_CAMELLIA_VALUE          0x0026
#define TPM_ALG_CAMELLIA                (TPM_ALG_ID)(ALG_CAMELLIA_VALUE)
#define     ALG_SHA3_256_VALUE          0x0027
#define TPM_ALG_SHA3_256                (TPM_ALG_ID)(ALG_SHA3_256_VALUE)
#define     ALG_SHA3_384_VALUE          0x0028
#define TPM_ALG_SHA3_384                (TPM_ALG_ID)(ALG_SHA3_384_VALUE)
#define     ALG_SHA3_512_VALUE          0x0029
#define TPM_ALG_SHA3_512                (TPM_ALG_ID)(ALG_SHA3_512_VALUE)
#define     ALG_CMAC_VALUE              0x003F
#define TPM_ALG_CMAC                    (TPM_ALG_ID)(ALG_CMAC_VALUE)
#define     ALG_CTR_VALUE               0x0040
#define TPM_ALG_CTR                     (TPM_ALG_ID)(ALG_CTR_VALUE)
#define     ALG_OFB_VALUE               0x0041
#define TPM_ALG_OFB                     (TPM_ALG_ID)(ALG_OFB_VALUE)
#define     ALG_CBC_VALUE               0x0042
#define TPM_ALG_CBC                     (TPM_ALG_ID)(ALG_CBC_VALUE)
#define     ALG_CFB_VALUE               0x0043
#define TPM_ALG_CFB                     (TPM_ALG_ID)(ALG_CFB_VALUE)
#define     ALG_ECB_VALUE               0x0044
#define TPM_ALG_ECB                     (TPM_ALG_ID)(ALG_ECB_VALUE)
// Values derived from Table 1:2
#define     ALG_FIRST_VALUE             0x0001
#define TPM_ALG_FIRST                   (TPM_ALG_ID)(ALG_FIRST_VALUE)
#define     ALG_LAST_VALUE              0x0044
#define TPM_ALG_LAST                    (TPM_ALG_ID)(ALG_LAST_VALUE)

// Table 1:3 - Definition of TPM_ECC_CURVE Constants
typedef UINT16              TPM_ECC_CURVE;
#define TYPE_OF_TPM_ECC_CURVE   UINT16
#define TPM_ECC_NONE        (TPM_ECC_CURVE)(0x0000)
#define TPM_ECC_NIST_P192   (TPM_ECC_CURVE)(0x0001)
#define TPM_ECC_NIST_P224   (TPM_ECC_CURVE)(0x0002)
#define TPM_ECC_NIST_P256   (TPM_ECC_CURVE)(0x0003)
#define TPM_ECC_NIST_P384   (TPM_ECC_CURVE)(0x0004)
#define TPM_ECC_NIST_P521   (TPM_ECC_CURVE)(0x0005)
#define TPM_ECC_BN_P256     (TPM_ECC_CURVE)(0x0010)
#define TPM_ECC_BN_P638     (TPM_ECC_CURVE)(0x0011)
#define TPM_ECC_SM2_P256    (TPM_ECC_CURVE)(0x0020)

// Table 2:12 - Definition of TPM_CC Constants
typedef UINT32                              TPM_CC;
#define TYPE_OF_TPM_CC                      UINT32
#define TPM_CC_NV_UndefineSpaceSpecial      (TPM_CC)(0x0000011F)
#define TPM_CC_EvictControl                 (TPM_CC)(0x00000120)
#define TPM_CC_HierarchyControl             (TPM_CC)(0x00000121)
#define TPM_CC_NV_UndefineSpace             (TPM_CC)(0x00000122)
#define TPM_CC_ChangeEPS                    (TPM_CC)(0x00000124)
#define TPM_CC_ChangePPS                    (TPM_CC)(0x00000125)
#define TPM_CC_Clear                        (TPM_CC)(0x00000126)
#define TPM_CC_ClearControl                 (TPM_CC)(0x00000127)
#define TPM_CC_ClockSet                     (TPM_CC)(0x00000128)
#define TPM_CC_HierarchyChangeAuth          (TPM_CC)(0x00000129)
#define TPM_CC_NV_DefineSpace               (TPM_CC)(0x0000012A)
#define TPM_CC_PCR_Allocate                 (TPM_CC)(0x0000012B)
#define TPM_CC_PCR_SetAuthPolicy            (TPM_CC)(0x0000012C)
#define TPM_CC_PP_Commands                  (TPM_CC)(0x0000012D)
#define TPM_CC_SetPrimaryPolicy             (TPM_CC)(0x0000012E)
#define TPM_CC_FieldUpgradeStart            (TPM_CC)(0x0000012F)
#define TPM_CC_ClockRateAdjust              (TPM_CC)(0x00000130)
#define TPM_CC_CreatePrimary                (TPM_CC)(0x00000131)
#define TPM_CC_NV_GlobalWriteLock           (TPM_CC)(0x00000132)
#define TPM_CC_GetCommandAuditDigest        (TPM_CC)(0x00000133)
#define TPM_CC_NV_Increment                 (TPM_CC)(0x00000134)
#define TPM_CC_NV_SetBits                   (TPM_CC)(0x00000135)
#define TPM_CC_NV_Extend                    (TPM_CC)(0x00000136)
#define TPM_CC_NV_Write                     (TPM_CC)(0x00000137)
#define TPM_CC_NV_WriteLock                 (TPM_CC)(0x00000138)
#define TPM_CC_DictionaryAttackLockReset    (TPM_CC)(0x00000139)
#define TPM_CC_DictionaryAttackParameters   (TPM_CC)(0x0000013A)
#define TPM_CC_NV_ChangeAuth                (TPM_CC)(0x0000013B)
#define TPM_CC_PCR_Event                    (TPM_CC)(0x0000013C)
#define TPM_CC_PCR_Reset                    (TPM_CC)(0x0000013D)
#define TPM_CC_SequenceComplete             (TPM_CC)(0x0000013E)
#define TPM_CC_SetAlgorithmSet              (TPM_CC)(0x0000013F)
#define TPM_CC_SetCommandCodeAuditStatus    (TPM_CC)(0x00000140)
#define TPM_CC_FieldUpgradeData             (TPM_CC)(0x00000141)
#define TPM_CC_IncrementalSelfTest          (TPM_CC)(0x00000142)
#define TPM_CC_SelfTest                     (TPM_CC)(0x00000143)
#define TPM_CC_Startup                      (TPM_CC)(0x00000144)
#define TPM_CC_Shutdown                     (TPM_CC)(0x00000145)
#define TPM_CC_StirRandom                   (TPM_CC)(0x00000146)
#define TPM_CC_ActivateCredential           (TPM_CC)(0x00000147)
#define TPM_CC_Certify                      (TPM_CC)(0x00000148)
#define TPM_CC_PolicyNV                     (TPM_CC)(0x00000149)
#define TPM_CC_CertifyCreation              (TPM_CC)(0x0000014A)
#define TPM_CC_Duplicate                    (TPM_CC)(0x0000014B)
#define TPM_CC_GetTime                      (TPM_CC)(0x0000014C)
#define TPM_CC_GetSessionAuditDigest        (TPM_CC)(0x0000014D)
#define TPM_CC_NV_Read                      (TPM_CC)(0x0000014E)
#define TPM_CC_NV_ReadLock                  (TPM_CC)(0x0000014F)
#define TPM_CC_ObjectChangeAuth             (TPM_CC)(0x00000150)
#define TPM_CC_PolicySecret                 (TPM_CC)(0x00000151)
#define TPM_CC_Rewrap                       (TPM_CC)(0x00000152)
#define TPM_CC_Create                       (TPM_CC)(0x00000153)
#define TPM_CC_ECDH_ZGen                    (TPM_CC)(0x00000154)
#define TPM_CC_HMAC                         (TPM_CC)(0x00000155)
#define TPM_CC_MAC                          (TPM_CC)(0x00000155)
#define TPM_CC_Import                       (TPM_CC)(0x00000156)
#define TPM_CC_Load                         (TPM_CC)(0x00000157)
#define TPM_CC_Quote                        (TPM_CC)(0x00000158)
#define TPM_CC_RSA_Decrypt                  (TPM_CC)(0x00000159)
#define TPM_CC_HMAC_Start                   (TPM_CC)(0x0000015B)
#define TPM_CC_MAC_Start                    (TPM_CC)(0x0000015B)
#define TPM_CC_SequenceUpdate               (TPM_CC)(0x0000015C)
#define TPM_CC_Sign                         (TPM_CC)(0x0000015D)
#define TPM_CC_Unseal                       (TPM_CC)(0x0000015E)
#define TPM_CC_PolicySigned                 (TPM_CC)(0x00000160)
#define TPM_CC_ContextLoad                  (TPM_CC)(0x00000161)
#define TPM_CC_ContextSave                  (TPM_CC)(0x00000162)
#define TPM_CC_ECDH_KeyGen                  (TPM_CC)(0x00000163)
#define TPM_CC_EncryptDecrypt               (TPM_CC)(0x00000164)
#define TPM_CC_FlushContext                 (TPM_CC)(0x00000165)
#define TPM_CC_LoadExternal                 (TPM_CC)(0x00000167)
#define TPM_CC_MakeCredential               (TPM_CC)(0x00000168)
#define TPM_CC_NV_ReadPublic                (TPM_CC)(0x00000169)
#define TPM_CC_PolicyAuthorize              (TPM_CC)(0x0000016A)
#define TPM_CC_PolicyAuthValue              (TPM_CC)(0x0000016B)
#define TPM_CC_PolicyCommandCode            (TPM_CC)(0x0000016C)
#define TPM_CC_PolicyCounterTimer           (TPM_CC)(0x0000016D)
#define TPM_CC_PolicyCpHash                 (TPM_CC)(0x0000016E)
#define TPM_CC_PolicyLocality               (TPM_CC)(0x0000016F)
#define TPM_CC_PolicyNameHash               (TPM_CC)(0x00000170)
#define TPM_CC_PolicyOR                     (TPM_CC)(0x00000171)
#define TPM_CC_PolicyTicket                 (TPM_CC)(0x00000172)
#define TPM_CC_ReadPublic                   (TPM_CC)(0x00000173)
#define TPM_CC_RSA_Encrypt                  (TPM_CC)(0x00000174)
#define TPM_CC_StartAuthSession             (TPM_CC)(0x00000176)
#define TPM_CC_VerifySignature              (TPM_CC)(0x00000177)
#define TPM_CC_ECC_Parameters               (TPM_CC)(0x00000178)
#define TPM_CC_FirmwareRead                 (TPM_CC)(0x00000179)
#define TPM_CC_GetCapability                (TPM_CC)(0x0000017A)
#define TPM_CC_GetRandom                    (TPM_CC)(0x0000017B)
#define TPM_CC_GetTestResult                (TPM_CC)(0x0000017C)
#define TPM_CC_Hash                         (TPM_CC)(0x0000017D)
#define TPM_CC_PCR_Read                     (TPM_CC)(0x0000017E)
#define TPM_CC_PolicyPCR                    (TPM_CC)(0x0000017F)
#define TPM_CC_PolicyRestart                (TPM_CC)(0x00000180)
#define TPM_CC_ReadClock                    (TPM_CC)(0x00000181)
#define TPM_CC_PCR_Extend                   (TPM_CC)(0x00000182)
#define TPM_CC_PCR_SetAuthValue             (TPM_CC)(0x00000183)
#define TPM_CC_NV_Certify                   (TPM_CC)(0x00000184)
#define TPM_CC_EventSequenceComplete        (TPM_CC)(0x00000185)
#define TPM_CC_HashSequenceStart            (TPM_CC)(0x00000186)
#define TPM_CC_PolicyPhysicalPresence       (TPM_CC)(0x00000187)
#define TPM_CC_PolicyDuplicationSelect      (TPM_CC)(0x00000188)
#define TPM_CC_PolicyGetDigest              (TPM_CC)(0x00000189)
#define TPM_CC_TestParms                    (TPM_CC)(0x0000018A)
#define TPM_CC_Commit                       (TPM_CC)(0x0000018B)
#define TPM_CC_PolicyPassword               (TPM_CC)(0x0000018C)
#define TPM_CC_ZGen_2Phase                  (TPM_CC)(0x0000018D)
#define TPM_CC_EC_Ephemeral                 (TPM_CC)(0x0000018E)
#define TPM_CC_PolicyNvWritten              (TPM_CC)(0x0000018F)
#define TPM_CC_PolicyTemplate               (TPM_CC)(0x00000190)
#define TPM_CC_CreateLoaded                 (TPM_CC)(0x00000191)
#define TPM_CC_PolicyAuthorizeNV            (TPM_CC)(0x00000192)
#define TPM_CC_EncryptDecrypt2              (TPM_CC)(0x00000193)
#define TPM_CC_AC_GetCapability             (TPM_CC)(0x00000194)
#define TPM_CC_AC_Send                      (TPM_CC)(0x00000195)
#define TPM_CC_Policy_AC_SendSelect         (TPM_CC)(0x00000196)
#define TPM_CC_CertifyX509                  (TPM_CC)(0x00000197)
#define CC_VEND                             0x20000000
#define TPM_CC_Vendor_TCG_Test              (TPM_CC)(0x20000000)

// Table 2:5 - Definition of Types for Documentation Clarity
typedef UINT32              TPM_ALGORITHM_ID;
#define TYPE_OF_TPM_ALGORITHM_ID    UINT32
typedef UINT32              TPM_MODIFIER_INDICATOR;
#define TYPE_OF_TPM_MODIFIER_INDICATOR  UINT32
typedef UINT32              TPM_AUTHORIZATION_SIZE;
#define TYPE_OF_TPM_AUTHORIZATION_SIZE  UINT32
typedef UINT32              TPM_PARAMETER_SIZE;
#define TYPE_OF_TPM_PARAMETER_SIZE  UINT32
typedef UINT16              TPM_KEY_SIZE;
#define TYPE_OF_TPM_KEY_SIZE    UINT16
typedef UINT16              TPM_KEY_BITS;
#define TYPE_OF_TPM_KEY_BITS    UINT16

// Table 2:6 - Definition of TPM_SPEC Constants
typedef UINT32                  TPM_SPEC;
#define TYPE_OF_TPM_SPEC        UINT32
#define SPEC_FAMILY             0x322E3000
#define TPM_SPEC_FAMILY         (TPM_SPEC)(SPEC_FAMILY)
#define SPEC_LEVEL              00
#define TPM_SPEC_LEVEL          (TPM_SPEC)(SPEC_LEVEL)
#define SPEC_VERSION            154
#define TPM_SPEC_VERSION        (TPM_SPEC)(SPEC_VERSION)
#define SPEC_YEAR               2019
#define TPM_SPEC_YEAR           (TPM_SPEC)(SPEC_YEAR)
#define SPEC_DAY_OF_YEAR        81
#define TPM_SPEC_DAY_OF_YEAR    (TPM_SPEC)(SPEC_DAY_OF_YEAR)

// Table 2:7 - Definition of TPM_GENERATED Constants
typedef UINT32                  TPM_GENERATED;
#define TYPE_OF_TPM_GENERATED   UINT32
#define TPM_GENERATED_VALUE     (TPM_GENERATED)(0xFF544347)

// Table 2:16 - Definition of TPM_RC Constants
typedef UINT32                      TPM_RC;
#define TYPE_OF_TPM_RC              UINT32
#define TPM_RC_SUCCESS              (TPM_RC)(0x000)
#define TPM_RC_BAD_TAG              (TPM_RC)(0x01E)
#define RC_VER1                     (TPM_RC)(0x100)
#define TPM_RC_INITIALIZE           (TPM_RC)(RC_VER1+0x000)
#define TPM_RC_FAILURE              (TPM_RC)(RC_VER1+0x001)
#define TPM_RC_SEQUENCE             (TPM_RC)(RC_VER1+0x003)
#define TPM_RC_PRIVATE              (TPM_RC)(RC_VER1+0x00B)
#define TPM_RC_HMAC                 (TPM_RC)(RC_VER1+0x019)
#define TPM_RC_DISABLED             (TPM_RC)(RC_VER1+0x020)
#define TPM_RC_EXCLUSIVE            (TPM_RC)(RC_VER1+0x021)
#define TPM_RC_AUTH_TYPE            (TPM_RC)(RC_VER1+0x024)
#define TPM_RC_AUTH_MISSING         (TPM_RC)(RC_VER1+0x025)
#define TPM_RC_POLICY               (TPM_RC)(RC_VER1+0x026)
#define TPM_RC_PCR                  (TPM_RC)(RC_VER1+0x027)
#define TPM_RC_PCR_CHANGED          (TPM_RC)(RC_VER1+0x028)
#define TPM_RC_UPGRADE              (TPM_RC)(RC_VER1+0x02D)
#define TPM_RC_TOO_MANY_CONTEXTS    (TPM_RC)(RC_VER1+0x02E)
#define TPM_RC_AUTH_UNAVAILABLE     (TPM_RC)(RC_VER1+0x02F)
#define TPM_RC_REBOOT               (TPM_RC)(RC_VER1+0x030)
#define TPM_RC_UNBALANCED           (TPM_RC)(RC_VER1+0x031)
#define TPM_RC_COMMAND_SIZE         (TPM_RC)(RC_VER1+0x042)
#define TPM_RC_COMMAND_CODE         (TPM_RC)(RC_VER1+0x043)
#define TPM_RC_AUTHSIZE             (TPM_RC)(RC_VER1+0x044)
#define TPM_RC_AUTH_CONTEXT         (TPM_RC)(RC_VER1+0x045)
#define TPM_RC_NV_RANGE             (TPM_RC)(RC_VER1+0x046)
#define TPM_RC_NV_SIZE              (TPM_RC)(RC_VER1+0x047)
#define TPM_RC_NV_LOCKED            (TPM_RC)(RC_VER1+0x048)
#define TPM_RC_NV_AUTHORIZATION     (TPM_RC)(RC_VER1+0x049)
#define TPM_RC_NV_UNINITIALIZED     (TPM_RC)(RC_VER1+0x04A)
#define TPM_RC_NV_SPACE             (TPM_RC)(RC_VER1+0x04B)
#define TPM_RC_NV_DEFINED           (TPM_RC)(RC_VER1+0x04C)
#define TPM_RC_BAD_CONTEXT          (TPM_RC)(RC_VER1+0x050)
#define TPM_RC_CPHASH               (TPM_RC)(RC_VER1+0x051)
#define TPM_RC_PARENT               (TPM_RC)(RC_VER1+0x052)
#define TPM_RC_NEEDS_TEST           (TPM_RC)(RC_VER1+0x053)
#define TPM_RC_NO_RESULT            (TPM_RC)(RC_VER1+0x054)
#define TPM_RC_SENSITIVE            (TPM_RC)(RC_VER1+0x055)
#define RC_MAX_FM0                  (TPM_RC)(RC_VER1+0x07F)
#define RC_FMT1                     (TPM_RC)(0x080)
#define TPM_RC_ASYMMETRIC           (TPM_RC)(RC_FMT1+0x001)
#define TPM_RCS_ASYMMETRIC          (TPM_RC)(RC_FMT1+0x001)
#define TPM_RC_ATTRIBUTES           (TPM_RC)(RC_FMT1+0x002)
#define TPM_RCS_ATTRIBUTES          (TPM_RC)(RC_FMT1+0x002)
#define TPM_RC_HASH                 (TPM_RC)(RC_FMT1+0x003)
#define TPM_RCS_HASH                (TPM_RC)(RC_FMT1+0x003)
#define TPM_RC_VALUE                (TPM_RC)(RC_FMT1+0x004)
#define TPM_RCS_VALUE               (TPM_RC)(RC_FMT1+0x004)
#define TPM_RC_HIERARCHY            (TPM_RC)(RC_FMT1+0x005)
#define TPM_RCS_HIERARCHY           (TPM_RC)(RC_FMT1+0x005)
#define TPM_RC_KEY_SIZE             (TPM_RC)(RC_FMT1+0x007)
#define TPM_RCS_KEY_SIZE            (TPM_RC)(RC_FMT1+0x007)
#define TPM_RC_MGF                  (TPM_RC)(RC_FMT1+0x008)
#define TPM_RCS_MGF                 (TPM_RC)(RC_FMT1+0x008)
#define TPM_RC_MODE                 (TPM_RC)(RC_FMT1+0x009)
#define TPM_RCS_MODE                (TPM_RC)(RC_FMT1+0x009)
#define TPM_RC_TYPE                 (TPM_RC)(RC_FMT1+0x00A)
#define TPM_RCS_TYPE                (TPM_RC)(RC_FMT1+0x00A)
#define TPM_RC_HANDLE               (TPM_RC)(RC_FMT1+0x00B)
#define TPM_RCS_HANDLE              (TPM_RC)(RC_FMT1+0x00B)
#define TPM_RC_KDF                  (TPM_RC)(RC_FMT1+0x00C)
#define TPM_RCS_KDF                 (TPM_RC)(RC_FMT1+0x00C)
#define TPM_RC_RANGE                (TPM_RC)(RC_FMT1+0x00D)
#define TPM_RCS_RANGE               (TPM_RC)(RC_FMT1+0x00D)
#define TPM_RC_AUTH_FAIL            (TPM_RC)(RC_FMT1+0x00E)
#define TPM_RCS_AUTH_FAIL           (TPM_RC)(RC_FMT1+0x00E)
#define TPM_RC_NONCE                (TPM_RC)(RC_FMT1+0x00F)
#define TPM_RCS_NONCE               (TPM_RC)(RC_FMT1+0x00F)
#define TPM_RC_PP                   (TPM_RC)(RC_FMT1+0x010)
#define TPM_RCS_PP                  (TPM_RC)(RC_FMT1+0x010)
#define TPM_RC_SCHEME               (TPM_RC)(RC_FMT1+0x012)
#define TPM_RCS_SCHEME              (TPM_RC)(RC_FMT1+0x012)
#define TPM_RC_SIZE                 (TPM_RC)(RC_FMT1+0x015)
#define TPM_RCS_SIZE                (TPM_RC)(RC_FMT1+0x015)
#define TPM_RC_SYMMETRIC            (TPM_RC)(RC_FMT1+0x016)
#define TPM_RCS_SYMMETRIC           (TPM_RC)(RC_FMT1+0x016)
#define TPM_RC_TAG                  (TPM_RC)(RC_FMT1+0x017)
#define TPM_RCS_TAG                 (TPM_RC)(RC_FMT1+0x017)
#define TPM_RC_SELECTOR             (TPM_RC)(RC_FMT1+0x018)
#define TPM_RCS_SELECTOR            (TPM_RC)(RC_FMT1+0x018)
#define TPM_RC_INSUFFICIENT         (TPM_RC)(RC_FMT1+0x01A)
#define TPM_RCS_INSUFFICIENT        (TPM_RC)(RC_FMT1+0x01A)
#define TPM_RC_SIGNATURE            (TPM_RC)(RC_FMT1+0x01B)
#define TPM_RCS_SIGNATURE           (TPM_RC)(RC_FMT1+0x01B)
#define TPM_RC_KEY                  (TPM_RC)(RC_FMT1+0x01C)
#define TPM_RCS_KEY                 (TPM_RC)(RC_FMT1+0x01C)
#define TPM_RC_POLICY_FAIL          (TPM_RC)(RC_FMT1+0x01D)
#define TPM_RCS_POLICY_FAIL         (TPM_RC)(RC_FMT1+0x01D)
#define TPM_RC_INTEGRITY            (TPM_RC)(RC_FMT1+0x01F)
#define TPM_RCS_INTEGRITY           (TPM_RC)(RC_FMT1+0x01F)
#define TPM_RC_TICKET               (TPM_RC)(RC_FMT1+0x020)
#define TPM_RCS_TICKET              (TPM_RC)(RC_FMT1+0x020)
#define TPM_RC_RESERVED_BITS        (TPM_RC)(RC_FMT1+0x021)
#define TPM_RCS_RESERVED_BITS       (TPM_RC)(RC_FMT1+0x021)
#define TPM_RC_BAD_AUTH             (TPM_RC)(RC_FMT1+0x022)
#define TPM_RCS_BAD_AUTH            (TPM_RC)(RC_FMT1+0x022)
#define TPM_RC_EXPIRED              (TPM_RC)(RC_FMT1+0x023)
#define TPM_RCS_EXPIRED             (TPM_RC)(RC_FMT1+0x023)
#define TPM_RC_POLICY_CC            (TPM_RC)(RC_FMT1+0x024)
#define TPM_RCS_POLICY_CC           (TPM_RC)(RC_FMT1+0x024)
#define TPM_RC_BINDING              (TPM_RC)(RC_FMT1+0x025)
#define TPM_RCS_BINDING             (TPM_RC)(RC_FMT1+0x025)
#define TPM_RC_CURVE                (TPM_RC)(RC_FMT1+0x026)
#define TPM_RCS_CURVE               (TPM_RC)(RC_FMT1+0x026)
#define TPM_RC_ECC_POINT            (TPM_RC)(RC_FMT1+0x027)
#define TPM_RCS_ECC_POINT           (TPM_RC)(RC_FMT1+0x027)
#define RC_WARN                     (TPM_RC)(0x900)
#define TPM_RC_CONTEXT_GAP          (TPM_RC)(RC_WARN+0x001)
#define TPM_RC_OBJECT_MEMORY        (TPM_RC)(RC_WARN+0x002)
#define TPM_RC_SESSION_MEMORY       (TPM_RC)(RC_WARN+0x003)
#define TPM_RC_MEMORY               (TPM_RC)(RC_WARN+0x004)
#define TPM_RC_SESSION_HANDLES      (TPM_RC)(RC_WARN+0x005)
#define TPM_RC_OBJECT_HANDLES       (TPM_RC)(RC_WARN+0x006)
#define TPM_RC_LOCALITY             (TPM_RC)(RC_WARN+0x007)
#define TPM_RC_YIELDED              (TPM_RC)(RC_WARN+0x008)
#define TPM_RC_CANCELED             (TPM_RC)(RC_WARN+0x009)
#define TPM_RC_TESTING              (TPM_RC)(RC_WARN+0x00A)
#define TPM_RC_REFERENCE_H0         (TPM_RC)(RC_WARN+0x010)
#define TPM_RC_REFERENCE_H1         (TPM_RC)(RC_WARN+0x011)
#define TPM_RC_REFERENCE_H2         (TPM_RC)(RC_WARN+0x012)
#define TPM_RC_REFERENCE_H3         (TPM_RC)(RC_WARN+0x013)
#define TPM_RC_REFERENCE_H4         (TPM_RC)(RC_WARN+0x014)
#define TPM_RC_REFERENCE_H5         (TPM_RC)(RC_WARN+0x015)
#define TPM_RC_REFERENCE_H6         (TPM_RC)(RC_WARN+0x016)
#define TPM_RC_REFERENCE_S0         (TPM_RC)(RC_WARN+0x018)
#define TPM_RC_REFERENCE_S1         (TPM_RC)(RC_WARN+0x019)
#define TPM_RC_REFERENCE_S2         (TPM_RC)(RC_WARN+0x01A)
#define TPM_RC_REFERENCE_S3         (TPM_RC)(RC_WARN+0x01B)
#define TPM_RC_REFERENCE_S4         (TPM_RC)(RC_WARN+0x01C)
#define TPM_RC_REFERENCE_S5         (TPM_RC)(RC_WARN+0x01D)
#define TPM_RC_REFERENCE_S6         (TPM_RC)(RC_WARN+0x01E)
#define TPM_RC_NV_RATE              (TPM_RC)(RC_WARN+0x020)
#define TPM_RC_LOCKOUT              (TPM_RC)(RC_WARN+0x021)
#define TPM_RC_RETRY                (TPM_RC)(RC_WARN+0x022)
#define TPM_RC_NV_UNAVAILABLE       (TPM_RC)(RC_WARN+0x023)
#define TPM_RC_NOT_USED             (TPM_RC)(RC_WARN+0x7F)
#define TPM_RC_H                    (TPM_RC)(0x000)
#define TPM_RC_P                    (TPM_RC)(0x040)
#define TPM_RC_S                    (TPM_RC)(0x800)
#define TPM_RC_1                    (TPM_RC)(0x100)
#define TPM_RC_2                    (TPM_RC)(0x200)
#define TPM_RC_3                    (TPM_RC)(0x300)
#define TPM_RC_4                    (TPM_RC)(0x400)
#define TPM_RC_5                    (TPM_RC)(0x500)
#define TPM_RC_6                    (TPM_RC)(0x600)
#define TPM_RC_7                    (TPM_RC)(0x700)
#define TPM_RC_8                    (TPM_RC)(0x800)
#define TPM_RC_9                    (TPM_RC)(0x900)
#define TPM_RC_A                    (TPM_RC)(0xA00)
#define TPM_RC_B                    (TPM_RC)(0xB00)
#define TPM_RC_C                    (TPM_RC)(0xC00)
#define TPM_RC_D                    (TPM_RC)(0xD00)
#define TPM_RC_E                    (TPM_RC)(0xE00)
#define TPM_RC_F                    (TPM_RC)(0xF00)
#define TPM_RC_N_MASK               (TPM_RC)(0xF00)

// Table 2:17 - Definition of TPM_CLOCK_ADJUST Constants
typedef INT8                        TPM_CLOCK_ADJUST;
#define TYPE_OF_TPM_CLOCK_ADJUST    UINT8
#define TPM_CLOCK_COARSE_SLOWER     (TPM_CLOCK_ADJUST)(-3)
#define TPM_CLOCK_MEDIUM_SLOWER     (TPM_CLOCK_ADJUST)(-2)
#define TPM_CLOCK_FINE_SLOWER       (TPM_CLOCK_ADJUST)(-1)
#define TPM_CLOCK_NO_CHANGE         (TPM_CLOCK_ADJUST)(0)
#define TPM_CLOCK_FINE_FASTER       (TPM_CLOCK_ADJUST)(1)
#define TPM_CLOCK_MEDIUM_FASTER     (TPM_CLOCK_ADJUST)(2)
#define TPM_CLOCK_COARSE_FASTER     (TPM_CLOCK_ADJUST)(3)

// Table 2:18 - Definition of TPM_EO Constants
typedef UINT16              TPM_EO;
#define TYPE_OF_TPM_EO      UINT16
#define TPM_EO_EQ           (TPM_EO)(0x0000)
#define TPM_EO_NEQ          (TPM_EO)(0x0001)
#define TPM_EO_SIGNED_GT    (TPM_EO)(0x0002)
#define TPM_EO_UNSIGNED_GT  (TPM_EO)(0x0003)
#define TPM_EO_SIGNED_LT    (TPM_EO)(0x0004)
#define TPM_EO_UNSIGNED_LT  (TPM_EO)(0x0005)
#define TPM_EO_SIGNED_GE    (TPM_EO)(0x0006)
#define TPM_EO_UNSIGNED_GE  (TPM_EO)(0x0007)
#define TPM_EO_SIGNED_LE    (TPM_EO)(0x0008)
#define TPM_EO_UNSIGNED_LE  (TPM_EO)(0x0009)
#define TPM_EO_BITSET       (TPM_EO)(0x000A)
#define TPM_EO_BITCLEAR     (TPM_EO)(0x000B)

// Table 2:19 - Definition of TPM_ST Constants
typedef UINT16                          TPM_ST;
#define TYPE_OF_TPM_ST                  UINT16
#define TPM_ST_RSP_COMMAND              (TPM_ST)(0x00C4)
#define TPM_ST_NULL                     (TPM_ST)(0x8000)
#define TPM_ST_NO_SESSIONS              (TPM_ST)(0x8001)
#define TPM_ST_SESSIONS                 (TPM_ST)(0x8002)
#define TPM_ST_ATTEST_NV                (TPM_ST)(0x8014)
#define TPM_ST_ATTEST_COMMAND_AUDIT     (TPM_ST)(0x8015)
#define TPM_ST_ATTEST_SESSION_AUDIT     (TPM_ST)(0x8016)
#define TPM_ST_ATTEST_CERTIFY           (TPM_ST)(0x8017)
#define TPM_ST_ATTEST_QUOTE             (TPM_ST)(0x8018)
#define TPM_ST_ATTEST_TIME              (TPM_ST)(0x8019)
#define TPM_ST_ATTEST_CREATION          (TPM_ST)(0x801A)
#define TPM_ST_ATTEST_NV_DIGEST         (TPM_ST)(0x801C)
#define TPM_ST_CREATION                 (TPM_ST)(0x8021)
#define TPM_ST_VERIFIED                 (TPM_ST)(0x8022)
#define TPM_ST_AUTH_SECRET              (TPM_ST)(0x8023)
#define TPM_ST_HASHCHECK                (TPM_ST)(0x8024)
#define TPM_ST_AUTH_SIGNED              (TPM_ST)(0x8025)
#define TPM_ST_FU_MANIFEST              (TPM_ST)(0x8029)

// Table 2:20 - Definition of TPM_SU Constants
typedef UINT16              TPM_SU;
#define TYPE_OF_TPM_SU      UINT16
#define TPM_SU_CLEAR        (TPM_SU)(0x0000)
#define TPM_SU_STATE        (TPM_SU)(0x0001)

// Table 2:21 - Definition of TPM_SE Constants
typedef UINT8               TPM_SE;
#define TYPE_OF_TPM_SE      UINT8
#define TPM_SE_HMAC         (TPM_SE)(0x00)
#define TPM_SE_POLICY       (TPM_SE)(0x01)
#define TPM_SE_TRIAL        (TPM_SE)(0x03)

// Table 2:22 - Definition of TPM_CAP Constants
typedef UINT32                      TPM_CAP;
#define TYPE_OF_TPM_CAP             UINT32
#define TPM_CAP_FIRST               (TPM_CAP)(0x00000000)
#define TPM_CAP_ALGS                (TPM_CAP)(0x00000000)
#define TPM_CAP_HANDLES             (TPM_CAP)(0x00000001)
#define TPM_CAP_COMMANDS            (TPM_CAP)(0x00000002)
#define TPM_CAP_PP_COMMANDS         (TPM_CAP)(0x00000003)
#define TPM_CAP_AUDIT_COMMANDS      (TPM_CAP)(0x00000004)
#define TPM_CAP_PCRS                (TPM_CAP)(0x00000005)
#define TPM_CAP_TPM_PROPERTIES      (TPM_CAP)(0x00000006)
#define TPM_CAP_PCR_PROPERTIES      (TPM_CAP)(0x00000007)
#define TPM_CAP_ECC_CURVES          (TPM_CAP)(0x00000008)
#define TPM_CAP_AUTH_POLICIES       (TPM_CAP)(0x00000009)
#define TPM_CAP_LAST                (TPM_CAP)(0x00000009)
#define TPM_CAP_VENDOR_PROPERTY     (TPM_CAP)(0x00000100)

// Table 2:23 - Definition of TPM_PT Constants
typedef UINT32                      TPM_PT;
#define TYPE_OF_TPM_PT              UINT32
#define TPM_PT_NONE                 (TPM_PT)(0x00000000)
#define PT_GROUP                    (TPM_PT)(0x00000100)
#define PT_FIXED                    (TPM_PT)(PT_GROUP*1)
#define TPM_PT_FAMILY_INDICATOR     (TPM_PT)(PT_FIXED+0)
#define TPM_PT_LEVEL                (TPM_PT)(PT_FIXED+1)
#define TPM_PT_REVISION             (TPM_PT)(PT_FIXED+2)
#define TPM_PT_DAY_OF_YEAR          (TPM_PT)(PT_FIXED+3)
#define TPM_PT_YEAR                 (TPM_PT)(PT_FIXED+4)
#define TPM_PT_MANUFACTURER         (TPM_PT)(PT_FIXED+5)
#define TPM_PT_VENDOR_STRING_1      (TPM_PT)(PT_FIXED+6)
#define TPM_PT_VENDOR_STRING_2      (TPM_PT)(PT_FIXED+7)
#define TPM_PT_VENDOR_STRING_3      (TPM_PT)(PT_FIXED+8)
#define TPM_PT_VENDOR_STRING_4      (TPM_PT)(PT_FIXED+9)
#define TPM_PT_VENDOR_TPM_TYPE      (TPM_PT)(PT_FIXED+10)
#define TPM_PT_FIRMWARE_VERSION_1   (TPM_PT)(PT_FIXED+11)
#define TPM_PT_FIRMWARE_VERSION_2   (TPM_PT)(PT_FIXED+12)
#define TPM_PT_INPUT_BUFFER         (TPM_PT)(PT_FIXED+13)
#define TPM_PT_HR_TRANSIENT_MIN     (TPM_PT)(PT_FIXED+14)
#define TPM_PT_HR_PERSISTENT_MIN    (TPM_PT)(PT_FIXED+15)
#define TPM_PT_HR_LOADED_MIN        (TPM_PT)(PT_FIXED+16)
#define TPM_PT_ACTIVE_SESSIONS_MAX  (TPM_PT)(PT_FIXED+17)
#define TPM_PT_PCR_COUNT            (TPM_PT)(PT_FIXED+18)
#define TPM_PT_PCR_SELECT_MIN       (TPM_PT)(PT_FIXED+19)
#define TPM_PT_CONTEXT_GAP_MAX      (TPM_PT)(PT_FIXED+20)
#define TPM_PT_NV_COUNTERS_MAX      (TPM_PT)(PT_FIXED+22)
#define TPM_PT_NV_INDEX_MAX         (TPM_PT)(PT_FIXED+23)
#define TPM_PT_MEMORY               (TPM_PT)(PT_FIXED+24)
#define TPM_PT_CLOCK_UPDATE         (TPM_PT)(PT_FIXED+25)
#define TPM_PT_CONTEXT_HASH         (TPM_PT)(PT_FIXED+26)
#define TPM_PT_CONTEXT_SYM          (TPM_PT)(PT_FIXED+27)
#define TPM_PT_CONTEXT_SYM_SIZE     (TPM_PT)(PT_FIXED+28)
#define TPM_PT_ORDERLY_COUNT        (TPM_PT)(PT_FIXED+29)
#define TPM_PT_MAX_COMMAND_SIZE     (TPM_PT)(PT_FIXED+30)
#define TPM_PT_MAX_RESPONSE_SIZE    (TPM_PT)(PT_FIXED+31)
#define TPM_PT_MAX_DIGEST           (TPM_PT)(PT_FIXED+32)
#define TPM_PT_MAX_OBJECT_CONTEXT   (TPM_PT)(PT_FIXED+33)
#define TPM_PT_MAX_SESSION_CONTEXT  (TPM_PT)(PT_FIXED+34)
#define TPM_PT_PS_FAMILY_INDICATOR  (TPM_PT)(PT_FIXED+35)
#define TPM_PT_PS_LEVEL             (TPM_PT)(PT_FIXED+36)
#define TPM_PT_PS_REVISION          (TPM_PT)(PT_FIXED+37)
#define TPM_PT_PS_DAY_OF_YEAR       (TPM_PT)(PT_FIXED+38)
#define TPM_PT_PS_YEAR              (TPM_PT)(PT_FIXED+39)
#define TPM_PT_SPLIT_MAX            (TPM_PT)(PT_FIXED+40)
#define TPM_PT_TOTAL_COMMANDS       (TPM_PT)(PT_FIXED+41)
#define TPM_PT_LIBRARY_COMMANDS     (TPM_PT)(PT_FIXED+42)
#define TPM_PT_VENDOR_COMMANDS      (TPM_PT)(PT_FIXED+43)
#define TPM_PT_NV_BUFFER_MAX        (TPM_PT)(PT_FIXED+44)
#define TPM_PT_MODES                (TPM_PT)(PT_FIXED+45)
#define TPM_PT_MAX_CAP_BUFFER       (TPM_PT)(PT_FIXED+46)
#define PT_VAR                      (TPM_PT)(PT_GROUP*2)
#define TPM_PT_PERMANENT            (TPM_PT)(PT_VAR+0)
#define TPM_PT_STARTUP_CLEAR        (TPM_PT)(PT_VAR+1)
#define TPM_PT_HR_NV_INDEX          (TPM_PT)(PT_VAR+2)
#define TPM_PT_HR_LOADED            (TPM_PT)(PT_VAR+3)
#define TPM_PT_HR_LOADED_AVAIL      (TPM_PT)(PT_VAR+4)
#define TPM_PT_HR_ACTIVE            (TPM_PT)(PT_VAR+5)
#define TPM_PT_HR_ACTIVE_AVAIL      (TPM_PT)(PT_VAR+6)
#define TPM_PT_HR_TRANSIENT_AVAIL   (TPM_PT)(PT_VAR+7)
#define TPM_PT_HR_PERSISTENT        (TPM_PT)(PT_VAR+8)
#define TPM_PT_HR_PERSISTENT_AVAIL  (TPM_PT)(PT_VAR+9)
#define TPM_PT_NV_COUNTERS          (TPM_PT)(PT_VAR+10)
#define TPM_PT_NV_COUNTERS_AVAIL    (TPM_PT)(PT_VAR+11)
#define TPM_PT_ALGORITHM_SET        (TPM_PT)(PT_VAR+12)
#define TPM_PT_LOADED_CURVES        (TPM_PT)(PT_VAR+13)
#define TPM_PT_LOCKOUT_COUNTER      (TPM_PT)(PT_VAR+14)
#define TPM_PT_MAX_AUTH_FAIL        (TPM_PT)(PT_VAR+15)
#define TPM_PT_LOCKOUT_INTERVAL     (TPM_PT)(PT_VAR+16)
#define TPM_PT_LOCKOUT_RECOVERY     (TPM_PT)(PT_VAR+17)
#define TPM_PT_NV_WRITE_RECOVERY    (TPM_PT)(PT_VAR+18)
#define TPM_PT_AUDIT_COUNTER_0      (TPM_PT)(PT_VAR+19)
#define TPM_PT_AUDIT_COUNTER_1      (TPM_PT)(PT_VAR+20)

// Table 2:24 - Definition of TPM_PT_PCR Constants
typedef UINT32                      TPM_PT_PCR;
#define TYPE_OF_TPM_PT_PCR          UINT32
#define TPM_PT_PCR_FIRST            (TPM_PT_PCR)(0x00000000)
#define TPM_PT_PCR_SAVE             (TPM_PT_PCR)(0x00000000)
#define TPM_PT_PCR_EXTEND_L0        (TPM_PT_PCR)(0x00000001)
#define TPM_PT_PCR_RESET_L0         (TPM_PT_PCR)(0x00000002)
#define TPM_PT_PCR_EXTEND_L1        (TPM_PT_PCR)(0x00000003)
#define TPM_PT_PCR_RESET_L1         (TPM_PT_PCR)(0x00000004)
#define TPM_PT_PCR_EXTEND_L2        (TPM_PT_PCR)(0x00000005)
#define TPM_PT_PCR_RESET_L2         (TPM_PT_PCR)(0x00000006)
#define TPM_PT_PCR_EXTEND_L3        (TPM_PT_PCR)(0x00000007)
#define TPM_PT_PCR_RESET_L3         (TPM_PT_PCR)(0x00000008)
#define TPM_PT_PCR_EXTEND_L4        (TPM_PT_PCR)(0x00000009)
#define TPM_PT_PCR_RESET_L4         (TPM_PT_PCR)(0x0000000A)
#define TPM_PT_PCR_NO_INCREMENT     (TPM_PT_PCR)(0x00000011)
#define TPM_PT_PCR_DRTM_RESET       (TPM_PT_PCR)(0x00000012)
#define TPM_PT_PCR_POLICY           (TPM_PT_PCR)(0x00000013)
#define TPM_PT_PCR_AUTH             (TPM_PT_PCR)(0x00000014)
#define TPM_PT_PCR_LAST             (TPM_PT_PCR)(0x00000014)

// Table 2:25 - Definition of TPM_PS Constants
typedef UINT32                  TPM_PS;
#define TYPE_OF_TPM_PS          UINT32
#define TPM_PS_MAIN             (TPM_PS)(0x00000000)
#define TPM_PS_PC               (TPM_PS)(0x00000001)
#define TPM_PS_PDA              (TPM_PS)(0x00000002)
#define TPM_PS_CELL_PHONE       (TPM_PS)(0x00000003)
#define TPM_PS_SERVER           (TPM_PS)(0x00000004)
#define TPM_PS_PERIPHERAL       (TPM_PS)(0x00000005)
#define TPM_PS_TSS              (TPM_PS)(0x00000006)
#define TPM_PS_STORAGE          (TPM_PS)(0x00000007)
#define TPM_PS_AUTHENTICATION   (TPM_PS)(0x00000008)
#define TPM_PS_EMBEDDED         (TPM_PS)(0x00000009)
#define TPM_PS_HARDCOPY         (TPM_PS)(0x0000000A)
#define TPM_PS_INFRASTRUCTURE   (TPM_PS)(0x0000000B)
#define TPM_PS_VIRTUALIZATION   (TPM_PS)(0x0000000C)
#define TPM_PS_TNC              (TPM_PS)(0x0000000D)
#define TPM_PS_MULTI_TENANT     (TPM_PS)(0x0000000E)
#define TPM_PS_TC               (TPM_PS)(0x0000000F)

// Table 2:26 - Definition of Types for Handles
typedef UINT32              TPM_HANDLE;
#define TYPE_OF_TPM_HANDLE  UINT32

// Table 2:27 - Definition of TPM_HT Constants
typedef UINT8                   TPM_HT;
#define TYPE_OF_TPM_HT          UINT8
#define TPM_HT_PCR              (TPM_HT)(0x00)
#define TPM_HT_NV_INDEX         (TPM_HT)(0x01)
#define TPM_HT_HMAC_SESSION     (TPM_HT)(0x02)
#define TPM_HT_LOADED_SESSION   (TPM_HT)(0x02)
#define TPM_HT_POLICY_SESSION   (TPM_HT)(0x03)
#define TPM_HT_SAVED_SESSION    (TPM_HT)(0x03)
#define TPM_HT_PERMANENT        (TPM_HT)(0x40)
#define TPM_HT_TRANSIENT        (TPM_HT)(0x80)
#define TPM_HT_PERSISTENT       (TPM_HT)(0x81)
#define TPM_HT_AC               (TPM_HT)(0x90)

// Table 2:28 - Definition of TPM_RH Constants
typedef TPM_HANDLE          TPM_RH;
#define TPM_RH_FIRST        (TPM_RH)(0x40000000)
#define TPM_RH_SRK          (TPM_RH)(0x40000000)
#define TPM_RH_OWNER        (TPM_RH)(0x40000001)
#define TPM_RH_REVOKE       (TPM_RH)(0x40000002)
#define TPM_RH_TRANSPORT    (TPM_RH)(0x40000003)
#define TPM_RH_OPERATOR     (TPM_RH)(0x40000004)
#define TPM_RH_ADMIN        (TPM_RH)(0x40000005)
#define TPM_RH_EK           (TPM_RH)(0x40000006)
#define TPM_RH_NULL         (TPM_RH)(0x40000007)
#define TPM_RH_UNASSIGNED   (TPM_RH)(0x40000008)
#define TPM_RS_PW           (TPM_RH)(0x40000009)
#define TPM_RH_LOCKOUT      (TPM_RH)(0x4000000A)
#define TPM_RH_ENDORSEMENT  (TPM_RH)(0x4000000B)
#define TPM_RH_PLATFORM     (TPM_RH)(0x4000000C)
#define TPM_RH_PLATFORM_NV  (TPM_RH)(0x4000000D)
#define TPM_RH_AUTH_00      (TPM_RH)(0x40000010)
#define TPM_RH_AUTH_FF      (TPM_RH)(0x4000010F)
#define TPM_RH_LAST         (TPM_RH)(0x4000010F)

// Table 2:29 - Definition of TPM_HC Constants
typedef TPM_HANDLE              TPM_HC;
#define HR_HANDLE_MASK          (TPM_HC)(0x00FFFFFF)
#define HR_RANGE_MASK           (TPM_HC)(0xFF000000)
#define HR_SHIFT                (TPM_HC)(24)
#define HR_PCR                  (TPM_HC)((TPM_HT_PCR<<HR_SHIFT))
#define HR_HMAC_SESSION         (TPM_HC)((TPM_HT_HMAC_SESSION<<HR_SHIFT))
#define HR_POLICY_SESSION       (TPM_HC)((TPM_HT_POLICY_SESSION<<HR_SHIFT))
#define HR_TRANSIENT            (TPM_HC)((TPM_HT_TRANSIENT<<HR_SHIFT))
#define HR_PERSISTENT           (TPM_HC)((TPM_HT_PERSISTENT<<HR_SHIFT))
#define HR_NV_INDEX             (TPM_HC)((TPM_HT_NV_INDEX<<HR_SHIFT))
#define HR_PERMANENT            (TPM_HC)((TPM_HT_PERMANENT<<HR_SHIFT))
#define PCR_FIRST               (TPM_HC)((HR_PCR+0))
#define PCR_LAST                (TPM_HC)((PCR_FIRST+IMPLEMENTATION_PCR-1))
#define HMAC_SESSION_FIRST      (TPM_HC)((HR_HMAC_SESSION+0))
#define HMAC_SESSION_LAST       (TPM_HC)((HMAC_SESSION_FIRST+MAX_ACTIVE_SESSIONS-1))
#define LOADED_SESSION_FIRST    (TPM_HC)(HMAC_SESSION_FIRST)
#define LOADED_SESSION_LAST     (TPM_HC)(HMAC_SESSION_LAST)
#define POLICY_SESSION_FIRST    (TPM_HC)((HR_POLICY_SESSION+0))
#define POLICY_SESSION_LAST     \
            (TPM_HC)((POLICY_SESSION_FIRST+MAX_ACTIVE_SESSIONS-1))
#define TRANSIENT_FIRST         (TPM_HC)((HR_TRANSIENT+0))
#define ACTIVE_SESSION_FIRST    (TPM_HC)(POLICY_SESSION_FIRST)
#define ACTIVE_SESSION_LAST     (TPM_HC)(POLICY_SESSION_LAST)
#define TRANSIENT_LAST          (TPM_HC)((TRANSIENT_FIRST+MAX_LOADED_OBJECTS-1))
#define PERSISTENT_FIRST        (TPM_HC)((HR_PERSISTENT+0))
#define PERSISTENT_LAST         (TPM_HC)((PERSISTENT_FIRST+0x00FFFFFF))
#define PLATFORM_PERSISTENT     (TPM_HC)((PERSISTENT_FIRST+0x00800000))
#define NV_INDEX_FIRST          (TPM_HC)((HR_NV_INDEX+0))
#define NV_INDEX_LAST           (TPM_HC)((NV_INDEX_FIRST+0x00FFFFFF))
#define PERMANENT_FIRST         (TPM_HC)(TPM_RH_FIRST)
#define PERMANENT_LAST          (TPM_HC)(TPM_RH_LAST)
#define HR_NV_AC                (TPM_HC)(((TPM_HT_NV_INDEX<<HR_SHIFT)+0xD00000))
#define NV_AC_FIRST             (TPM_HC)((HR_NV_AC+0))
#define NV_AC_LAST              (TPM_HC)((HR_NV_AC+0x0000FFFF))
#define HR_AC                   (TPM_HC)((TPM_HT_AC<<HR_SHIFT))
#define AC_FIRST                (TPM_HC)((HR_AC+0))
#define AC_LAST                 (TPM_HC)((HR_AC+0x0000FFFF))

#define TYPE_OF_TPMA_ALGORITHM  UINT32
#define TPMA_ALGORITHM_TO_UINT32(a)  (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_ALGORITHM(a)  (*((TPMA_ALGORITHM *)&(a)))
#define TPMA_ALGORITHM_TO_BYTE_ARRAY(i, a)                                         \
            UINT32_TO_BYTE_ARRAY((TPMA_ALGORITHM_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_ALGORITHM(i, a)                                         \
            {UINT32 x = BYTE_ARRAY_TO_UINT32(a);                                   \
             i = UINT32_TO_TPMA_ALGORITHM(x);                                      \
             }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_ALGORITHM {                     // Table 2:30
    unsigned    asymmetric           : 1;
    unsigned    symmetric            : 1;
    unsigned    hash                 : 1;
    unsigned    object               : 1;
    unsigned    Reserved_bits_at_4   : 4;
    unsigned    signing              : 1;
    unsigned    encrypting           : 1;
    unsigned    method               : 1;
    unsigned    Reserved_bits_at_11  : 21;
} TPMA_ALGORITHM;                                   /* Bits */
// This is the initializer for a TPMA_ALGORITHM structure
#define TPMA_ALGORITHM_INITIALIZER(                                                \
             asymmetric, symmetric,  hash,       object,     bits_at_4,            \
             signing,    encrypting, method,     bits_at_11)                       \
            {asymmetric, symmetric,  hash,       object,     bits_at_4,            \
             signing,    encrypting, method,     bits_at_11}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:30 TPMA_ALGORITHM using bit masking
typedef UINT32                      TPMA_ALGORITHM;
#define TYPE_OF_TPMA_ALGORITHM      UINT32
#define TPMA_ALGORITHM_asymmetric   ((TPMA_ALGORITHM)1 << 0)
#define TPMA_ALGORITHM_symmetric    ((TPMA_ALGORITHM)1 << 1)
#define TPMA_ALGORITHM_hash         ((TPMA_ALGORITHM)1 << 2)
#define TPMA_ALGORITHM_object       ((TPMA_ALGORITHM)1 << 3)
#define TPMA_ALGORITHM_signing      ((TPMA_ALGORITHM)1 << 8)
#define TPMA_ALGORITHM_encrypting   ((TPMA_ALGORITHM)1 << 9)
#define TPMA_ALGORITHM_method       ((TPMA_ALGORITHM)1 << 10)
//  This is the initializer for a TPMA_ALGORITHM bit array.
#define TPMA_ALGORITHM_INITIALIZER(                                                \
             asymmetric, symmetric,  hash,       object,     bits_at_4,            \
             signing,    encrypting, method,     bits_at_11)                       \
            {(asymmetric << 0) + (symmetric << 1)  + (hash << 2)       +           \
             (object << 3)     + (signing << 8)    + (encrypting << 9) +           \
             (method << 10)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_OBJECT UINT32
#define TPMA_OBJECT_TO_UINT32(a)     (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_OBJECT(a)     (*((TPMA_OBJECT *)&(a)))
#define TPMA_OBJECT_TO_BYTE_ARRAY(i, a)                                            \
            UINT32_TO_BYTE_ARRAY((TPMA_OBJECT_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_OBJECT(i, a)                                            \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPMA_OBJECT(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_OBJECT {                        // Table 2:31
    unsigned    Reserved_bit_at_0        : 1;
    unsigned    fixedTPM                 : 1;
    unsigned    stClear                  : 1;
    unsigned    Reserved_bit_at_3        : 1;
    unsigned    fixedParent              : 1;
    unsigned    sensitiveDataOrigin      : 1;
    unsigned    userWithAuth             : 1;
    unsigned    adminWithPolicy          : 1;
    unsigned    Reserved_bits_at_8       : 2;
    unsigned    noDA                     : 1;
    unsigned    encryptedDuplication     : 1;
    unsigned    Reserved_bits_at_12      : 4;
    unsigned    restricted               : 1;
    unsigned    decrypt                  : 1;
    unsigned    sign                     : 1;
    unsigned    x509sign                 : 1;
    unsigned    Reserved_bits_at_20      : 12;
} TPMA_OBJECT;                                      /* Bits */
// This is the initializer for a TPMA_OBJECT structure
#define TPMA_OBJECT_INITIALIZER(                                                   \
             bit_at_0,             fixedtpm,             stclear,                  \
             bit_at_3,             fixedparent,          sensitivedataorigin,      \
             userwithauth,         adminwithpolicy,      bits_at_8,                \
             noda,                 encryptedduplication, bits_at_12,               \
             restricted,           decrypt,              sign,                     \
             x509sign,             bits_at_20)                                     \
            {bit_at_0,             fixedtpm,             stclear,                  \
             bit_at_3,             fixedparent,          sensitivedataorigin,      \
             userwithauth,         adminwithpolicy,      bits_at_8,                \
             noda,                 encryptedduplication, bits_at_12,               \
             restricted,           decrypt,              sign,                     \
             x509sign,             bits_at_20}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:31 TPMA_OBJECT using bit masking
typedef UINT32                              TPMA_OBJECT;
#define TYPE_OF_TPMA_OBJECT                 UINT32
#define TPMA_OBJECT_fixedTPM                ((TPMA_OBJECT)1 << 1)
#define TPMA_OBJECT_stClear                 ((TPMA_OBJECT)1 << 2)
#define TPMA_OBJECT_fixedParent             ((TPMA_OBJECT)1 << 4)
#define TPMA_OBJECT_sensitiveDataOrigin     ((TPMA_OBJECT)1 << 5)
#define TPMA_OBJECT_userWithAuth            ((TPMA_OBJECT)1 << 6)
#define TPMA_OBJECT_adminWithPolicy         ((TPMA_OBJECT)1 << 7)
#define TPMA_OBJECT_noDA                    ((TPMA_OBJECT)1 << 10)
#define TPMA_OBJECT_encryptedDuplication    ((TPMA_OBJECT)1 << 11)
#define TPMA_OBJECT_restricted              ((TPMA_OBJECT)1 << 16)
#define TPMA_OBJECT_decrypt                 ((TPMA_OBJECT)1 << 17)
#define TPMA_OBJECT_sign                    ((TPMA_OBJECT)1 << 18)
#define TPMA_OBJECT_x509sign                ((TPMA_OBJECT)1 << 19)
//  This is the initializer for a TPMA_OBJECT bit array.
#define TPMA_OBJECT_INITIALIZER(                                                   \
             bit_at_0,             fixedtpm,             stclear,                  \
             bit_at_3,             fixedparent,          sensitivedataorigin,      \
             userwithauth,         adminwithpolicy,      bits_at_8,                \
             noda,                 encryptedduplication, bits_at_12,               \
             restricted,           decrypt,              sign,                     \
             x509sign,             bits_at_20)                                     \
            {(fixedtpm << 1)              + (stclear << 2)               +         \
             (fixedparent << 4)           + (sensitivedataorigin << 5)   +         \
             (userwithauth << 6)          + (adminwithpolicy << 7)       +         \
             (noda << 10)                 + (encryptedduplication << 11) +         \
             (restricted << 16)           + (decrypt << 17)              +         \
             (sign << 18)                 + (x509sign << 19)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_SESSION    UINT8
#define TPMA_SESSION_TO_UINT8(a)     (*((UINT8 *)&(a)))
#define UINT8_TO_TPMA_SESSION(a)     (*((TPMA_SESSION *)&(a)))
#define TPMA_SESSION_TO_BYTE_ARRAY(i, a)                                           \
            UINT8_TO_BYTE_ARRAY((TPMA_SESSION_TO_UINT8(i)), (a))
#define BYTE_ARRAY_TO_TPMA_SESSION(i, a)                                           \
            { UINT8 x = BYTE_ARRAY_TO_UINT8(a); i = UINT8_TO_TPMA_SESSION(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_SESSION {                       // Table 2:32
    unsigned    continueSession      : 1;
    unsigned    auditExclusive       : 1;
    unsigned    auditReset           : 1;
    unsigned    Reserved_bits_at_3   : 2;
    unsigned    decrypt              : 1;
    unsigned    encrypt              : 1;
    unsigned    audit                : 1;
} TPMA_SESSION;                                     /* Bits */
// This is the initializer for a TPMA_SESSION structure
#define TPMA_SESSION_INITIALIZER(                                                  \
             continuesession, auditexclusive,  auditreset,      bits_at_3,         \
             decrypt,         encrypt,         audit)                              \
            {continuesession, auditexclusive,  auditreset,      bits_at_3,         \
             decrypt,         encrypt,         audit}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:32 TPMA_SESSION using bit masking
typedef UINT8                           TPMA_SESSION;
#define TYPE_OF_TPMA_SESSION            UINT8
#define TPMA_SESSION_continueSession    ((TPMA_SESSION)1 << 0)
#define TPMA_SESSION_auditExclusive     ((TPMA_SESSION)1 << 1)
#define TPMA_SESSION_auditReset         ((TPMA_SESSION)1 << 2)
#define TPMA_SESSION_decrypt            ((TPMA_SESSION)1 << 5)
#define TPMA_SESSION_encrypt            ((TPMA_SESSION)1 << 6)
#define TPMA_SESSION_audit              ((TPMA_SESSION)1 << 7)
//  This is the initializer for a TPMA_SESSION bit array.
#define TPMA_SESSION_INITIALIZER(                                                  \
             continuesession, auditexclusive,  auditreset,      bits_at_3,         \
             decrypt,         encrypt,         audit)                              \
            {(continuesession << 0) + (auditexclusive << 1)  +                     \
             (auditreset << 2)      + (decrypt << 5)         +                     \
             (encrypt << 6)         + (audit << 7)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_LOCALITY   UINT8
#define TPMA_LOCALITY_TO_UINT8(a)    (*((UINT8 *)&(a)))
#define UINT8_TO_TPMA_LOCALITY(a)    (*((TPMA_LOCALITY *)&(a)))
#define TPMA_LOCALITY_TO_BYTE_ARRAY(i, a)                                          \
            UINT8_TO_BYTE_ARRAY((TPMA_LOCALITY_TO_UINT8(i)), (a))
#define BYTE_ARRAY_TO_TPMA_LOCALITY(i, a)                                          \
            { UINT8 x = BYTE_ARRAY_TO_UINT8(a); i = UINT8_TO_TPMA_LOCALITY(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_LOCALITY {                      // Table 2:33
    unsigned    TPM_LOC_ZERO         : 1;
    unsigned    TPM_LOC_ONE          : 1;
    unsigned    TPM_LOC_TWO          : 1;
    unsigned    TPM_LOC_THREE        : 1;
    unsigned    TPM_LOC_FOUR         : 1;
    unsigned    Extended             : 3;
} TPMA_LOCALITY;                                    /* Bits */
// This is the initializer for a TPMA_LOCALITY structure
#define TPMA_LOCALITY_INITIALIZER(                                                 \
             tpm_loc_zero,  tpm_loc_one,   tpm_loc_two,   tpm_loc_three,           \
             tpm_loc_four,  extended)                                              \
            {tpm_loc_zero,  tpm_loc_one,   tpm_loc_two,   tpm_loc_three,           \
             tpm_loc_four,  extended}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:33 TPMA_LOCALITY using bit masking
typedef UINT8                           TPMA_LOCALITY;
#define TYPE_OF_TPMA_LOCALITY           UINT8
#define TPMA_LOCALITY_TPM_LOC_ZERO      ((TPMA_LOCALITY)1 << 0)
#define TPMA_LOCALITY_TPM_LOC_ONE       ((TPMA_LOCALITY)1 << 1)
#define TPMA_LOCALITY_TPM_LOC_TWO       ((TPMA_LOCALITY)1 << 2)
#define TPMA_LOCALITY_TPM_LOC_THREE     ((TPMA_LOCALITY)1 << 3)
#define TPMA_LOCALITY_TPM_LOC_FOUR      ((TPMA_LOCALITY)1 << 4)
#define TPMA_LOCALITY_Extended_SHIFT    5
#define TPMA_LOCALITY_Extended          ((TPMA_LOCALITY)0x7 << 5)
//  This is the initializer for a TPMA_LOCALITY bit array.
#define TPMA_LOCALITY_INITIALIZER(                                                 \
             tpm_loc_zero,  tpm_loc_one,   tpm_loc_two,   tpm_loc_three,           \
             tpm_loc_four,  extended)                                              \
            {(tpm_loc_zero << 0)  + (tpm_loc_one << 1)   + (tpm_loc_two << 2)   +  \
             (tpm_loc_three << 3) + (tpm_loc_four << 4)  + (extended << 5)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_PERMANENT  UINT32
#define TPMA_PERMANENT_TO_UINT32(a)  (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_PERMANENT(a)  (*((TPMA_PERMANENT *)&(a)))
#define TPMA_PERMANENT_TO_BYTE_ARRAY(i, a)                                         \
            UINT32_TO_BYTE_ARRAY((TPMA_PERMANENT_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_PERMANENT(i, a)                                         \
            {UINT32 x = BYTE_ARRAY_TO_UINT32(a);                                   \
             i = UINT32_TO_TPMA_PERMANENT(x);                                      \
             }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_PERMANENT {                     // Table 2:34
    unsigned    ownerAuthSet         : 1;
    unsigned    endorsementAuthSet   : 1;
    unsigned    lockoutAuthSet       : 1;
    unsigned    Reserved_bits_at_3   : 5;
    unsigned    disableClear         : 1;
    unsigned    inLockout            : 1;
    unsigned    tpmGeneratedEPS      : 1;
    unsigned    Reserved_bits_at_11  : 21;
} TPMA_PERMANENT;                                   /* Bits */
// This is the initializer for a TPMA_PERMANENT structure
#define TPMA_PERMANENT_INITIALIZER(                                                \
             ownerauthset,       endorsementauthset, lockoutauthset,               \
             bits_at_3,          disableclear,       inlockout,                    \
             tpmgeneratedeps,    bits_at_11)                                       \
            {ownerauthset,       endorsementauthset, lockoutauthset,               \
             bits_at_3,          disableclear,       inlockout,                    \
             tpmgeneratedeps,    bits_at_11}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:34 TPMA_PERMANENT using bit masking
typedef UINT32                              TPMA_PERMANENT;
#define TYPE_OF_TPMA_PERMANENT              UINT32
#define TPMA_PERMANENT_ownerAuthSet         ((TPMA_PERMANENT)1 << 0)
#define TPMA_PERMANENT_endorsementAuthSet   ((TPMA_PERMANENT)1 << 1)
#define TPMA_PERMANENT_lockoutAuthSet       ((TPMA_PERMANENT)1 << 2)
#define TPMA_PERMANENT_disableClear         ((TPMA_PERMANENT)1 << 8)
#define TPMA_PERMANENT_inLockout            ((TPMA_PERMANENT)1 << 9)
#define TPMA_PERMANENT_tpmGeneratedEPS      ((TPMA_PERMANENT)1 << 10)
//  This is the initializer for a TPMA_PERMANENT bit array.
#define TPMA_PERMANENT_INITIALIZER(                                                \
             ownerauthset,       endorsementauthset, lockoutauthset,               \
             bits_at_3,          disableclear,       inlockout,                    \
             tpmgeneratedeps,    bits_at_11)                                       \
            {(ownerauthset << 0)       + (endorsementauthset << 1) +               \
             (lockoutauthset << 2)     + (disableclear << 8)       +               \
             (inlockout << 9)          + (tpmgeneratedeps << 10)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_STARTUP_CLEAR  UINT32
#define TPMA_STARTUP_CLEAR_TO_UINT32(a)  (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_STARTUP_CLEAR(a)  (*((TPMA_STARTUP_CLEAR *)&(a)))
#define TPMA_STARTUP_CLEAR_TO_BYTE_ARRAY(i, a)                                     \
            UINT32_TO_BYTE_ARRAY((TPMA_STARTUP_CLEAR_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_STARTUP_CLEAR(i, a)                                     \
            {UINT32 x = BYTE_ARRAY_TO_UINT32(a);                                   \
             i = UINT32_TO_TPMA_STARTUP_CLEAR(x);                                  \
             }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_STARTUP_CLEAR {                 // Table 2:35
    unsigned    phEnable             : 1;
    unsigned    shEnable             : 1;
    unsigned    ehEnable             : 1;
    unsigned    phEnableNV           : 1;
    unsigned    Reserved_bits_at_4   : 27;
    unsigned    orderly              : 1;
} TPMA_STARTUP_CLEAR;                               /* Bits */
// This is the initializer for a TPMA_STARTUP_CLEAR structure
#define TPMA_STARTUP_CLEAR_INITIALIZER(                                            \
             phenable, shenable, ehenable, phenablenv, bits_at_4, orderly)         \
            {phenable, shenable, ehenable, phenablenv, bits_at_4, orderly}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:35 TPMA_STARTUP_CLEAR using bit masking
typedef UINT32                          TPMA_STARTUP_CLEAR;
#define TYPE_OF_TPMA_STARTUP_CLEAR      UINT32
#define TPMA_STARTUP_CLEAR_phEnable     ((TPMA_STARTUP_CLEAR)1 << 0)
#define TPMA_STARTUP_CLEAR_shEnable     ((TPMA_STARTUP_CLEAR)1 << 1)
#define TPMA_STARTUP_CLEAR_ehEnable     ((TPMA_STARTUP_CLEAR)1 << 2)
#define TPMA_STARTUP_CLEAR_phEnableNV   ((TPMA_STARTUP_CLEAR)1 << 3)
#define TPMA_STARTUP_CLEAR_orderly      ((TPMA_STARTUP_CLEAR)1 << 31)
//  This is the initializer for a TPMA_STARTUP_CLEAR bit array.
#define TPMA_STARTUP_CLEAR_INITIALIZER(                                            \
             phenable, shenable, ehenable, phenablenv, bits_at_4, orderly)         \
            {(phenable << 0)   + (shenable << 1)   + (ehenable << 2)   +           \
             (phenablenv << 3) + (orderly << 31)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_MEMORY UINT32
#define TPMA_MEMORY_TO_UINT32(a)     (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_MEMORY(a)     (*((TPMA_MEMORY *)&(a)))
#define TPMA_MEMORY_TO_BYTE_ARRAY(i, a)                                            \
            UINT32_TO_BYTE_ARRAY((TPMA_MEMORY_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_MEMORY(i, a)                                            \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPMA_MEMORY(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_MEMORY {                        // Table 2:36
    unsigned    sharedRAM            : 1;
    unsigned    sharedNV             : 1;
    unsigned    objectCopiedToRam    : 1;
    unsigned    Reserved_bits_at_3   : 29;
} TPMA_MEMORY;                                      /* Bits */
// This is the initializer for a TPMA_MEMORY structure
#define TPMA_MEMORY_INITIALIZER(                                                   \
             sharedram, sharednv, objectcopiedtoram, bits_at_3)                    \
            {sharedram, sharednv, objectcopiedtoram, bits_at_3}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:36 TPMA_MEMORY using bit masking
typedef UINT32                          TPMA_MEMORY;
#define TYPE_OF_TPMA_MEMORY             UINT32
#define TPMA_MEMORY_sharedRAM           ((TPMA_MEMORY)1 << 0)
#define TPMA_MEMORY_sharedNV            ((TPMA_MEMORY)1 << 1)
#define TPMA_MEMORY_objectCopiedToRam   ((TPMA_MEMORY)1 << 2)
//  This is the initializer for a TPMA_MEMORY bit array.
#define TPMA_MEMORY_INITIALIZER(                                                   \
             sharedram, sharednv, objectcopiedtoram, bits_at_3)                    \
            {(sharedram << 0) + (sharednv << 1) + (objectcopiedtoram << 2)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_CC     UINT32
#define TPMA_CC_TO_UINT32(a)     (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_CC(a)     (*((TPMA_CC *)&(a)))
#define TPMA_CC_TO_BYTE_ARRAY(i, a)                                                \
            UINT32_TO_BYTE_ARRAY((TPMA_CC_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_CC(i, a)                                                \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPMA_CC(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_CC {                            // Table 2:37
    unsigned    commandIndex         : 16;
    unsigned    Reserved_bits_at_16  : 6;
    unsigned    nv                   : 1;
    unsigned    extensive            : 1;
    unsigned    flushed              : 1;
    unsigned    cHandles             : 3;
    unsigned    rHandle              : 1;
    unsigned    V                    : 1;
    unsigned    Reserved_bits_at_30  : 2;
} TPMA_CC;                                          /* Bits */
// This is the initializer for a TPMA_CC structure
#define TPMA_CC_INITIALIZER(                                                       \
             commandindex, bits_at_16,   nv,           extensive,    flushed,      \
             chandles,     rhandle,      v,            bits_at_30)                 \
            {commandindex, bits_at_16,   nv,           extensive,    flushed,      \
             chandles,     rhandle,      v,            bits_at_30}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:37 TPMA_CC using bit masking
typedef UINT32                      TPMA_CC;
#define TYPE_OF_TPMA_CC             UINT32
#define TPMA_CC_commandIndex_SHIFT  0
#define TPMA_CC_commandIndex        ((TPMA_CC)0xffff << 0)
#define TPMA_CC_nv                  ((TPMA_CC)1 << 22)
#define TPMA_CC_extensive           ((TPMA_CC)1 << 23)
#define TPMA_CC_flushed             ((TPMA_CC)1 << 24)
#define TPMA_CC_cHandles_SHIFT      25
#define TPMA_CC_cHandles            ((TPMA_CC)0x7 << 25)
#define TPMA_CC_rHandle             ((TPMA_CC)1 << 28)
#define TPMA_CC_V                   ((TPMA_CC)1 << 29)
//  This is the initializer for a TPMA_CC bit array.
#define TPMA_CC_INITIALIZER(                                                       \
             commandindex, bits_at_16,   nv,           extensive,    flushed,      \
             chandles,     rhandle,      v,            bits_at_30)                 \
            {(commandindex << 0) + (nv << 22)          + (extensive << 23)   +     \
             (flushed << 24)     + (chandles << 25)    + (rhandle << 28)     +     \
             (v << 29)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_MODES  UINT32
#define TPMA_MODES_TO_UINT32(a)  (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_MODES(a)  (*((TPMA_MODES *)&(a)))
#define TPMA_MODES_TO_BYTE_ARRAY(i, a)                                             \
            UINT32_TO_BYTE_ARRAY((TPMA_MODES_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_MODES(i, a)                                             \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPMA_MODES(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_MODES {                         // Table 2:38
    unsigned    FIPS_140_2           : 1;
    unsigned    Reserved_bits_at_1   : 31;
} TPMA_MODES;                                       /* Bits */
// This is the initializer for a TPMA_MODES structure
#define TPMA_MODES_INITIALIZER(fips_140_2, bits_at_1) {fips_140_2, bits_at_1}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:38 TPMA_MODES using bit masking
typedef UINT32                  TPMA_MODES;
#define TYPE_OF_TPMA_MODES      UINT32
#define TPMA_MODES_FIPS_140_2   ((TPMA_MODES)1 << 0)
//  This is the initializer for a TPMA_MODES bit array.
#define TPMA_MODES_INITIALIZER(fips_140_2, bits_at_1) {(fips_140_2 << 0)}
#endif // USE_BIT_FIELD_STRUCTURES

#define TYPE_OF_TPMA_X509_KEY_USAGE UINT32
#define TPMA_X509_KEY_USAGE_TO_UINT32(a)     (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_X509_KEY_USAGE(a)     (*((TPMA_X509_KEY_USAGE *)&(a)))
#define TPMA_X509_KEY_USAGE_TO_BYTE_ARRAY(i, a)                                    \
            UINT32_TO_BYTE_ARRAY((TPMA_X509_KEY_USAGE_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_X509_KEY_USAGE(i, a)                                    \
            {UINT32 x = BYTE_ARRAY_TO_UINT32(a);                                   \
             i = UINT32_TO_TPMA_X509_KEY_USAGE(x);                                 \
             }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_X509_KEY_USAGE {                // Table 2:39
    unsigned    digitalSignature     : 1;
    unsigned    nonrepudiation       : 1;
    unsigned    keyEncipherment      : 1;
    unsigned    dataEncipherment     : 1;
    unsigned    keyAgreement         : 1;
    unsigned    keyCertSign          : 1;
    unsigned    crlSign              : 1;
    unsigned    encipherOnly         : 1;
    unsigned    decipherOnly         : 1;
    unsigned    Reserved_bits_at_9   : 23;
} TPMA_X509_KEY_USAGE;                              /* Bits */
// This is the initializer for a TPMA_X509_KEY_USAGE structure
#define TPMA_X509_KEY_USAGE_INITIALIZER(                                           \
             digitalsignature, nonrepudiation,   keyencipherment,                  \
             dataencipherment, keyagreement,     keycertsign,                      \
             crlsign,          encipheronly,     decipheronly,                     \
             bits_at_9)                                                            \
            {digitalsignature, nonrepudiation,   keyencipherment,                  \
             dataencipherment, keyagreement,     keycertsign,                      \
             crlsign,          encipheronly,     decipheronly,                     \
             bits_at_9}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:39 TPMA_X509_KEY_USAGE using bit masking
typedef UINT32                                  TPMA_X509_KEY_USAGE;
#define TYPE_OF_TPMA_X509_KEY_USAGE             UINT32
#define TPMA_X509_KEY_USAGE_digitalSignature    ((TPMA_X509_KEY_USAGE)1 << 0)
#define TPMA_X509_KEY_USAGE_nonrepudiation      ((TPMA_X509_KEY_USAGE)1 << 1)
#define TPMA_X509_KEY_USAGE_keyEncipherment     ((TPMA_X509_KEY_USAGE)1 << 2)
#define TPMA_X509_KEY_USAGE_dataEncipherment    ((TPMA_X509_KEY_USAGE)1 << 3)
#define TPMA_X509_KEY_USAGE_keyAgreement        ((TPMA_X509_KEY_USAGE)1 << 4)
#define TPMA_X509_KEY_USAGE_keyCertSign         ((TPMA_X509_KEY_USAGE)1 << 5)
#define TPMA_X509_KEY_USAGE_crlSign             ((TPMA_X509_KEY_USAGE)1 << 6)
#define TPMA_X509_KEY_USAGE_encipherOnly        ((TPMA_X509_KEY_USAGE)1 << 7)
#define TPMA_X509_KEY_USAGE_decipherOnly        ((TPMA_X509_KEY_USAGE)1 << 8)
//  This is the initializer for a TPMA_X509_KEY_USAGE bit array.
#define TPMA_X509_KEY_USAGE_INITIALIZER(                                           \
             digitalsignature, nonrepudiation,   keyencipherment,                  \
             dataencipherment, keyagreement,     keycertsign,                      \
             crlsign,          encipheronly,     decipheronly,                     \
             bits_at_9)                                                            \
            {(digitalsignature << 0) + (nonrepudiation << 1)   +                   \
             (keyencipherment << 2)  + (dataencipherment << 3) +                   \
             (keyagreement << 4)     + (keycertsign << 5)      +                   \
             (crlsign << 6)          + (encipheronly << 7)     +                   \
             (decipheronly << 8)}
#endif // USE_BIT_FIELD_STRUCTURES

typedef BYTE                TPMI_YES_NO;            // Table 2:40  /* Interface */

typedef TPM_HANDLE          TPMI_DH_OBJECT;         // Table 2:41  /* Interface */

typedef TPM_HANDLE          TPMI_DH_PARENT;         // Table 2:42  /* Interface */

typedef TPM_HANDLE          TPMI_DH_PERSISTENT;     // Table 2:43  /* Interface */

typedef TPM_HANDLE          TPMI_DH_ENTITY;         // Table 2:44  /* Interface */

typedef TPM_HANDLE          TPMI_DH_PCR;            // Table 2:45  /* Interface */

typedef TPM_HANDLE          TPMI_SH_AUTH_SESSION;   // Table 2:46  /* Interface */

typedef TPM_HANDLE          TPMI_SH_HMAC;           // Table 2:47  /* Interface */

typedef TPM_HANDLE          TPMI_SH_POLICY;         // Table 2:48  /* Interface */

typedef TPM_HANDLE          TPMI_DH_CONTEXT;        // Table 2:49  /* Interface */

typedef TPM_HANDLE          TPMI_DH_SAVED;          // Table 2:50  /* Interface */

typedef TPM_HANDLE          TPMI_RH_HIERARCHY;      // Table 2:51  /* Interface */

typedef TPM_HANDLE          TPMI_RH_ENABLES;        // Table 2:52  /* Interface */

typedef TPM_HANDLE          TPMI_RH_HIERARCHY_AUTH; // Table 2:53  /* Interface */

typedef TPM_HANDLE          TPMI_RH_PLATFORM;       // Table 2:54  /* Interface */

typedef TPM_HANDLE          TPMI_RH_OWNER;          // Table 2:55  /* Interface */

typedef TPM_HANDLE          TPMI_RH_ENDORSEMENT;    // Table 2:56  /* Interface */

typedef TPM_HANDLE          TPMI_RH_PROVISION;      // Table 2:57  /* Interface */

typedef TPM_HANDLE          TPMI_RH_CLEAR;          // Table 2:58  /* Interface */

typedef TPM_HANDLE          TPMI_RH_NV_AUTH;        // Table 2:59  /* Interface */

typedef TPM_HANDLE          TPMI_RH_LOCKOUT;        // Table 2:60  /* Interface */

typedef TPM_HANDLE          TPMI_RH_NV_INDEX;       // Table 2:61  /* Interface */

typedef TPM_HANDLE          TPMI_RH_AC;             // Table 2:62  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_HASH;          // Table 2:63  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_ASYM;          // Table 2:64  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_SYM;           // Table 2:65  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_SYM_OBJECT;    // Table 2:66  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_SYM_MODE;      // Table 2:67  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_KDF;           // Table 2:68  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_SIG_SCHEME;    // Table 2:69  /* Interface */

typedef TPM_ALG_ID          TPMI_ECC_KEY_EXCHANGE;  // Table 2:70  /* Interface */

typedef TPM_ST              TPMI_ST_COMMAND_TAG;    // Table 2:71  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_MAC_SCHEME;    // Table 2:72  /* Interface */

typedef TPM_ALG_ID          TPMI_ALG_CIPHER_MODE;   // Table 2:73  /* Interface */

typedef BYTE                TPMS_EMPTY;             // Table 2:74

typedef struct {                                    // Table 2:75
    TPM_ALG_ID              alg;
    TPMA_ALGORITHM          attributes;
} TPMS_ALGORITHM_DESCRIPTION;                       /* Structure */

typedef union {                                     // Table 2:76
#if ALG_SHA1
    BYTE                    sha1[SHA1_DIGEST_SIZE];
#endif // ALG_SHA1
#if ALG_SHA256
    BYTE                    sha256[SHA256_DIGEST_SIZE];
#endif // ALG_SHA256
#if ALG_SHA384
    BYTE                    sha384[SHA384_DIGEST_SIZE];
#endif // ALG_SHA384
#if ALG_SHA512
    BYTE                    sha512[SHA512_DIGEST_SIZE];
#endif // ALG_SHA512
#if ALG_SM3_256
    BYTE                    sm3_256[SM3_256_DIGEST_SIZE];
#endif // ALG_SM3_256
#if ALG_SHA3_256
    BYTE                    sha3_256[SHA3_256_DIGEST_SIZE];
#endif // ALG_SHA3_256
#if ALG_SHA3_384
    BYTE                    sha3_384[SHA3_384_DIGEST_SIZE];
#endif // ALG_SHA3_384
#if ALG_SHA3_512
    BYTE                    sha3_512[SHA3_512_DIGEST_SIZE];
#endif // ALG_SHA3_512
} TPMU_HA;                                          /* Structure */

typedef struct {                                    // Table 2:77
    TPMI_ALG_HASH           hashAlg;
    TPMU_HA                 digest;
} TPMT_HA;                                          /* Structure */

typedef union {                                     // Table 2:78
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMU_HA)];
    }            t;
    TPM2B        b;
} TPM2B_DIGEST;                                     /* Structure */

typedef union {                                     // Table 2:79
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMT_HA)];
    }            t;
    TPM2B        b;
} TPM2B_DATA;                                       /* Structure */

// Table 2:80 - Definition of Types for TPM2B_NONCE
typedef TPM2B_DIGEST        TPM2B_NONCE;

// Table 2:81 - Definition of Types for TPM2B_AUTH
typedef TPM2B_DIGEST        TPM2B_AUTH;

// Table 2:82 - Definition of Types for TPM2B_OPERAND
typedef TPM2B_DIGEST        TPM2B_OPERAND;

typedef union {                                     // Table 2:83
    struct {
        UINT16              size;
        BYTE                buffer[1024];
    }            t;
    TPM2B        b;
} TPM2B_EVENT;                                      /* Structure */

typedef union {                                     // Table 2:84
    struct {
        UINT16              size;
        BYTE                buffer[MAX_DIGEST_BUFFER];
    }            t;
    TPM2B        b;
} TPM2B_MAX_BUFFER;                                 /* Structure */

typedef union {                                     // Table 2:85
    struct {
        UINT16              size;
        BYTE                buffer[MAX_NV_BUFFER_SIZE];
    }            t;
    TPM2B        b;
} TPM2B_MAX_NV_BUFFER;                              /* Structure */

typedef union {                                     // Table 2:86
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(UINT64)];
    }            t;
    TPM2B        b;
} TPM2B_TIMEOUT;                                    /* Structure */

typedef union {                                     // Table 2:87
    struct {
        UINT16              size;
        BYTE                buffer[MAX_SYM_BLOCK_SIZE];
    }            t;
    TPM2B        b;
} TPM2B_IV;                                         /* Structure */

typedef union {                                     // Table 2:88
    TPMT_HA                 digest;
    TPM_HANDLE              handle;
} TPMU_NAME;                                        /* Structure */

typedef union {                                     // Table 2:89
    struct {
        UINT16              size;
        BYTE                name[sizeof(TPMU_NAME)];
    }            t;
    TPM2B        b;
} TPM2B_NAME;                                       /* Structure */

typedef struct {                                    // Table 2:90
    UINT8                   sizeofSelect;
    BYTE                    pcrSelect[PCR_SELECT_MAX];
} TPMS_PCR_SELECT;                                  /* Structure */

typedef struct {                                    // Table 2:91
    TPMI_ALG_HASH           hash;
    UINT8                   sizeofSelect;
    BYTE                    pcrSelect[PCR_SELECT_MAX];
} TPMS_PCR_SELECTION;                               /* Structure */

typedef struct {                                    // Table 2:94
    TPM_ST                  tag;
    TPMI_RH_HIERARCHY       hierarchy;
    TPM2B_DIGEST            digest;
} TPMT_TK_CREATION;                                 /* Structure */

typedef struct {                                    // Table 2:95
    TPM_ST                  tag;
    TPMI_RH_HIERARCHY       hierarchy;
    TPM2B_DIGEST            digest;
} TPMT_TK_VERIFIED;                                 /* Structure */

typedef struct {                                    // Table 2:96
    TPM_ST                  tag;
    TPMI_RH_HIERARCHY       hierarchy;
    TPM2B_DIGEST            digest;
} TPMT_TK_AUTH;                                     /* Structure */

typedef struct {                                    // Table 2:97
    TPM_ST                  tag;
    TPMI_RH_HIERARCHY       hierarchy;
    TPM2B_DIGEST            digest;
} TPMT_TK_HASHCHECK;                                /* Structure */

typedef struct {                                    // Table 2:98
    TPM_ALG_ID              alg;
    TPMA_ALGORITHM          algProperties;
} TPMS_ALG_PROPERTY;                                /* Structure */

typedef struct {                                    // Table 2:99
    TPM_PT                  property;
    UINT32                  value;
} TPMS_TAGGED_PROPERTY;                             /* Structure */

typedef struct {                                    // Table 2:100
    TPM_PT_PCR              tag;
    UINT8                   sizeofSelect;
    BYTE                    pcrSelect[PCR_SELECT_MAX];
} TPMS_TAGGED_PCR_SELECT;                           /* Structure */

typedef struct {                                    // Table 2:101
    TPM_HANDLE              handle;
    TPMT_HA                 policyHash;
} TPMS_TAGGED_POLICY;                               /* Structure */

typedef struct {                                    // Table 2:102
    UINT32                  count;
    TPM_CC                  commandCodes[MAX_CAP_CC];
} TPML_CC;                                          /* Structure */

typedef struct {                                    // Table 2:103
    UINT32                  count;
    TPMA_CC                 commandAttributes[MAX_CAP_CC];
} TPML_CCA;                                         /* Structure */

typedef struct {                                    // Table 2:104
    UINT32                  count;
    TPM_ALG_ID              algorithms[MAX_ALG_LIST_SIZE];
} TPML_ALG;                                         /* Structure */

typedef struct {                                    // Table 2:105
    UINT32                  count;
    TPM_HANDLE              handle[MAX_CAP_HANDLES];
} TPML_HANDLE;                                      /* Structure */

typedef struct {                                    // Table 2:106
    UINT32                  count;
    TPM2B_DIGEST            digests[8];
} TPML_DIGEST;                                      /* Structure */

typedef struct {                                    // Table 2:107
    UINT32                  count;
    TPMT_HA                 digests[HASH_COUNT];
} TPML_DIGEST_VALUES;                               /* Structure */

typedef struct {                                    // Table 2:108
    UINT32                  count;
    TPMS_PCR_SELECTION      pcrSelections[HASH_COUNT];
} TPML_PCR_SELECTION;                               /* Structure */

typedef struct {                                    // Table 2:109
    UINT32                  count;
    TPMS_ALG_PROPERTY       algProperties[MAX_CAP_ALGS];
} TPML_ALG_PROPERTY;                                /* Structure */

typedef struct {                                    // Table 2:110
    UINT32                      count;
    TPMS_TAGGED_PROPERTY        tpmProperty[MAX_TPM_PROPERTIES];
} TPML_TAGGED_TPM_PROPERTY;                         /* Structure */

typedef struct {                                    // Table 2:111
    UINT32                      count;
    TPMS_TAGGED_PCR_SELECT      pcrProperty[MAX_PCR_PROPERTIES];
} TPML_TAGGED_PCR_PROPERTY;                         /* Structure */

typedef struct {                                    // Table 2:112
    UINT32                  count;
    TPM_ECC_CURVE           eccCurves[MAX_ECC_CURVES];
} TPML_ECC_CURVE;                                   /* Structure */

typedef struct {                                    // Table 2:113
    UINT32                  count;
    TPMS_TAGGED_POLICY      policies[MAX_TAGGED_POLICIES];
} TPML_TAGGED_POLICY;                               /* Structure */

typedef union {                                     // Table 2:114
    TPML_ALG_PROPERTY               algorithms;
    TPML_HANDLE                     handles;
    TPML_CCA                        command;
    TPML_CC                         ppCommands;
    TPML_CC                         auditCommands;
    TPML_PCR_SELECTION              assignedPCR;
    TPML_TAGGED_TPM_PROPERTY        tpmProperties;
    TPML_TAGGED_PCR_PROPERTY        pcrProperties;
#if ALG_ECC
    TPML_ECC_CURVE                  eccCurves;
#endif // ALG_ECC
    TPML_TAGGED_POLICY              authPolicies;
} TPMU_CAPABILITIES;                                /* Structure */

typedef struct {                                    // Table 2:115
    TPM_CAP                 capability;
    TPMU_CAPABILITIES       data;
} TPMS_CAPABILITY_DATA;                             /* Structure */

typedef struct {                                    // Table 2:116
    UINT64                  clock;
    UINT32                  resetCount;
    UINT32                  restartCount;
    TPMI_YES_NO             safe;
} TPMS_CLOCK_INFO;                                  /* Structure */

typedef struct {                                    // Table 2:117
    UINT64                  time;
    TPMS_CLOCK_INFO         clockInfo;
} TPMS_TIME_INFO;                                   /* Structure */

typedef struct {                                    // Table 2:118
    TPMS_TIME_INFO          time;
    UINT64                  firmwareVersion;
} TPMS_TIME_ATTEST_INFO;                            /* Structure */

typedef struct {                                    // Table 2:119
    TPM2B_NAME              name;
    TPM2B_NAME              qualifiedName;
} TPMS_CERTIFY_INFO;                                /* Structure */

typedef struct {                                    // Table 2:120
    TPML_PCR_SELECTION      pcrSelect;
    TPM2B_DIGEST            pcrDigest;
} TPMS_QUOTE_INFO;                                  /* Structure */

typedef struct {                                    // Table 2:121
    UINT64                  auditCounter;
    TPM_ALG_ID              digestAlg;
    TPM2B_DIGEST            auditDigest;
    TPM2B_DIGEST            commandDigest;
} TPMS_COMMAND_AUDIT_INFO;                          /* Structure */

typedef struct {                                    // Table 2:122
    TPMI_YES_NO             exclusiveSession;
    TPM2B_DIGEST            sessionDigest;
} TPMS_SESSION_AUDIT_INFO;                          /* Structure */

typedef struct {                                    // Table 2:123
    TPM2B_NAME              objectName;
    TPM2B_DIGEST            creationHash;
} TPMS_CREATION_INFO;                               /* Structure */

typedef struct {                                    // Table 2:124
    TPM2B_NAME                  indexName;
    UINT16                      offset;
    TPM2B_MAX_NV_BUFFER         nvContents;
} TPMS_NV_CERTIFY_INFO;                             /* Structure */

typedef struct {                                    // Table 2:125
    TPM2B_NAME              indexName;
    TPM2B_DIGEST            nvDigest;
} TPMS_NV_DIGEST_CERTIFY_INFO;                      /* Structure */

typedef TPM_ST              TPMI_ST_ATTEST;         // Table 2:126  /* Interface */

typedef union {                                             // Table 2:127
    TPMS_CERTIFY_INFO                   certify;
    TPMS_CREATION_INFO                  creation;
    TPMS_QUOTE_INFO                     quote;
    TPMS_COMMAND_AUDIT_INFO             commandAudit;
    TPMS_SESSION_AUDIT_INFO             sessionAudit;
    TPMS_TIME_ATTEST_INFO               time;
    TPMS_NV_CERTIFY_INFO                nv;
    TPMS_NV_DIGEST_CERTIFY_INFO         nvDigest;
} TPMU_ATTEST;                                              /* Structure */

typedef struct {                                    // Table 2:128
    TPM_GENERATED           magic;
    TPMI_ST_ATTEST          type;
    TPM2B_NAME              qualifiedSigner;
    TPM2B_DATA              extraData;
    TPMS_CLOCK_INFO         clockInfo;
    UINT64                  firmwareVersion;
    TPMU_ATTEST             attested;
} TPMS_ATTEST;                                      /* Structure */

typedef union {                                     // Table 2:129
    struct {
        UINT16              size;
        BYTE                attestationData[sizeof(TPMS_ATTEST)];
    }            t;
    TPM2B        b;
} TPM2B_ATTEST;                                     /* Structure */

typedef struct {                                    // Table 2:130
    TPMI_SH_AUTH_SESSION        sessionHandle;
    TPM2B_NONCE                 nonce;
    TPMA_SESSION                sessionAttributes;
    TPM2B_AUTH                  hmac;
} TPMS_AUTH_COMMAND;                                /* Structure */

typedef struct {                                    // Table 2:131
    TPM2B_NONCE             nonce;
    TPMA_SESSION            sessionAttributes;
    TPM2B_AUTH              hmac;
} TPMS_AUTH_RESPONSE;                               /* Structure */

typedef TPM_KEY_BITS        TPMI_TDES_KEY_BITS;     // Table 2:132  /* Interface */

typedef TPM_KEY_BITS        TPMI_AES_KEY_BITS;      // Table 2:132  /* Interface */

typedef TPM_KEY_BITS        TPMI_SM4_KEY_BITS;      // Table 2:132  /* Interface */

typedef TPM_KEY_BITS        TPMI_CAMELLIA_KEY_BITS; // Table 2:132  /* Interface */

typedef union {                                     // Table 2:133
#if ALG_TDES
    TPMI_TDES_KEY_BITS          tdes;
#endif // ALG_TDES
#if ALG_AES
    TPMI_AES_KEY_BITS           aes;
#endif // ALG_AES
#if ALG_SM4
    TPMI_SM4_KEY_BITS           sm4;
#endif // ALG_SM4
#if ALG_CAMELLIA
    TPMI_CAMELLIA_KEY_BITS      camellia;
#endif // ALG_CAMELLIA
    TPM_KEY_BITS                sym;
#if ALG_XOR
    TPMI_ALG_HASH               xor;
#endif // ALG_XOR
} TPMU_SYM_KEY_BITS;                                /* Structure */

typedef union {                                     // Table 2:134
#if ALG_TDES
    TPMI_ALG_SYM_MODE       tdes;
#endif // ALG_TDES
#if ALG_AES
    TPMI_ALG_SYM_MODE       aes;
#endif // ALG_AES
#if ALG_SM4
    TPMI_ALG_SYM_MODE       sm4;
#endif // ALG_SM4
#if ALG_CAMELLIA
    TPMI_ALG_SYM_MODE       camellia;
#endif // ALG_CAMELLIA
    TPMI_ALG_SYM_MODE       sym;
} TPMU_SYM_MODE;                                    /* Structure */

typedef struct {                                    // Table 2:136
    TPMI_ALG_SYM            algorithm;
    TPMU_SYM_KEY_BITS       keyBits;
    TPMU_SYM_MODE           mode;
} TPMT_SYM_DEF;                                     /* Structure */

typedef struct {                                    // Table 2:137
    TPMI_ALG_SYM_OBJECT         algorithm;
    TPMU_SYM_KEY_BITS           keyBits;
    TPMU_SYM_MODE               mode;
} TPMT_SYM_DEF_OBJECT;                              /* Structure */

typedef union {                                     // Table 2:138
    struct {
        UINT16              size;
        BYTE                buffer[MAX_SYM_KEY_BYTES];
    }            t;
    TPM2B        b;
} TPM2B_SYM_KEY;                                    /* Structure */

typedef struct {                                    // Table 2:139
    TPMT_SYM_DEF_OBJECT         sym;
} TPMS_SYMCIPHER_PARMS;                             /* Structure */

typedef union {                                     // Table 2:140
    struct {
        UINT16              size;
        BYTE                buffer[LABEL_MAX_BUFFER];
    }            t;
    TPM2B        b;
} TPM2B_LABEL;                                      /* Structure */

typedef struct {                                    // Table 2:141
    TPM2B_LABEL             label;
    TPM2B_LABEL             context;
} TPMS_DERIVE;                                      /* Structure */

typedef union {                                     // Table 2:142
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMS_DERIVE)];
    }            t;
    TPM2B        b;
} TPM2B_DERIVE;                                     /* Structure */

typedef union {                                     // Table 2:143
    BYTE                    create[MAX_SYM_DATA];
    TPMS_DERIVE             derive;
} TPMU_SENSITIVE_CREATE;                            /* Structure */

typedef union {                                     // Table 2:144
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMU_SENSITIVE_CREATE)];
    }            t;
    TPM2B        b;
} TPM2B_SENSITIVE_DATA;                             /* Structure */

typedef struct {                                    // Table 2:145
    TPM2B_AUTH                  userAuth;
    TPM2B_SENSITIVE_DATA        data;
} TPMS_SENSITIVE_CREATE;                            /* Structure */

typedef struct {                                    // Table 2:146
    UINT16                      size;
    TPMS_SENSITIVE_CREATE       sensitive;
} TPM2B_SENSITIVE_CREATE;                           /* Structure */

typedef struct {                                    // Table 2:147
    TPMI_ALG_HASH           hashAlg;
} TPMS_SCHEME_HASH;                                 /* Structure */

typedef struct {                                    // Table 2:148
    TPMI_ALG_HASH           hashAlg;
    UINT16                  count;
} TPMS_SCHEME_ECDAA;                                /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_KEYEDHASH_SCHEME;

// Table 2:150 - Definition of Types for HMAC_SIG_SCHEME
typedef TPMS_SCHEME_HASH    TPMS_SCHEME_HMAC;

typedef struct {                                    // Table 2:151
    TPMI_ALG_HASH           hashAlg;
    TPMI_ALG_KDF            kdf;
} TPMS_SCHEME_XOR;                                  /* Structure */

typedef union {                                     // Table 2:152
#if ALG_HMAC
    TPMS_SCHEME_HMAC        hmac;
#endif // ALG_HMAC
#if ALG_XOR
    TPMS_SCHEME_XOR         xor;
#endif // ALG_XOR
} TPMU_SCHEME_KEYEDHASH;                            /* Structure */

typedef struct {                                    // Table 2:153
    TPMI_ALG_KEYEDHASH_SCHEME       scheme;
    TPMU_SCHEME_KEYEDHASH           details;
} TPMT_KEYEDHASH_SCHEME;                            /* Structure */

// Table 2:154 - Definition of Types for RSA Signature Schemes
typedef TPMS_SCHEME_HASH    TPMS_SIG_SCHEME_RSASSA;
typedef TPMS_SCHEME_HASH    TPMS_SIG_SCHEME_RSAPSS;

// Table 2:155 - Definition of Types for ECC Signature Schemes
typedef TPMS_SCHEME_HASH    TPMS_SIG_SCHEME_ECDSA;
typedef TPMS_SCHEME_HASH    TPMS_SIG_SCHEME_SM2;
typedef TPMS_SCHEME_HASH    TPMS_SIG_SCHEME_ECSCHNORR;
typedef TPMS_SCHEME_ECDAA   TPMS_SIG_SCHEME_ECDAA;

typedef union {                                     // Table 2:156
#if ALG_ECC
    TPMS_SIG_SCHEME_ECDAA           ecdaa;
#endif // ALG_ECC
#if ALG_RSASSA
    TPMS_SIG_SCHEME_RSASSA          rsassa;
#endif // ALG_RSASSA
#if ALG_RSAPSS
    TPMS_SIG_SCHEME_RSAPSS          rsapss;
#endif // ALG_RSAPSS
#if ALG_ECDSA
    TPMS_SIG_SCHEME_ECDSA           ecdsa;
#endif // ALG_ECDSA
#if ALG_SM2
    TPMS_SIG_SCHEME_SM2             sm2;
#endif // ALG_SM2
#if ALG_ECSCHNORR
    TPMS_SIG_SCHEME_ECSCHNORR       ecschnorr;
#endif // ALG_ECSCHNORR
#if ALG_HMAC
    TPMS_SCHEME_HMAC                hmac;
#endif // ALG_HMAC
    TPMS_SCHEME_HASH                any;
} TPMU_SIG_SCHEME;                                  /* Structure */

typedef struct {                                    // Table 2:157
    TPMI_ALG_SIG_SCHEME         scheme;
    TPMU_SIG_SCHEME             details;
} TPMT_SIG_SCHEME;                                  /* Structure */

// Table 2:158 - Definition of Types for Encryption Schemes
typedef TPMS_SCHEME_HASH    TPMS_ENC_SCHEME_OAEP;
typedef TPMS_EMPTY          TPMS_ENC_SCHEME_RSAES;

// Table 2:159 - Definition of Types for ECC Key Exchange
typedef TPMS_SCHEME_HASH    TPMS_KEY_SCHEME_ECDH;
typedef TPMS_SCHEME_HASH    TPMS_KEY_SCHEME_ECMQV;

// Table 2:160 - Definition of Types for KDF Schemes
typedef TPMS_SCHEME_HASH    TPMS_SCHEME_MGF1;
typedef TPMS_SCHEME_HASH    TPMS_SCHEME_KDF1_SP800_56A;
typedef TPMS_SCHEME_HASH    TPMS_SCHEME_KDF2;
typedef TPMS_SCHEME_HASH    TPMS_SCHEME_KDF1_SP800_108;

typedef union {                                     // Table 2:161
#if ALG_MGF1
    TPMS_SCHEME_MGF1                mgf1;
#endif // ALG_MGF1
#if ALG_KDF1_SP800_56A
    TPMS_SCHEME_KDF1_SP800_56A      kdf1_sp800_56a;
#endif // ALG_KDF1_SP800_56A
#if ALG_KDF2
    TPMS_SCHEME_KDF2                kdf2;
#endif // ALG_KDF2
#if ALG_KDF1_SP800_108
    TPMS_SCHEME_KDF1_SP800_108      kdf1_sp800_108;
#endif // ALG_KDF1_SP800_108
} TPMU_KDF_SCHEME;                                  /* Structure */

typedef struct {                                    // Table 2:162
    TPMI_ALG_KDF            scheme;
    TPMU_KDF_SCHEME         details;
} TPMT_KDF_SCHEME;                                  /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_ASYM_SCHEME;   // Table 2:163  /* Interface */

typedef union {                                     // Table 2:164
#if ALG_ECDH
    TPMS_KEY_SCHEME_ECDH            ecdh;
#endif // ALG_ECDH
#if ALG_ECMQV
    TPMS_KEY_SCHEME_ECMQV           ecmqv;
#endif // ALG_ECMQV
#if ALG_ECC
    TPMS_SIG_SCHEME_ECDAA           ecdaa;
#endif // ALG_ECC
#if ALG_RSASSA
    TPMS_SIG_SCHEME_RSASSA          rsassa;
#endif // ALG_RSASSA
#if ALG_RSAPSS
    TPMS_SIG_SCHEME_RSAPSS          rsapss;
#endif // ALG_RSAPSS
#if ALG_ECDSA
    TPMS_SIG_SCHEME_ECDSA           ecdsa;
#endif // ALG_ECDSA
#if ALG_SM2
    TPMS_SIG_SCHEME_SM2             sm2;
#endif // ALG_SM2
#if ALG_ECSCHNORR
    TPMS_SIG_SCHEME_ECSCHNORR       ecschnorr;
#endif // ALG_ECSCHNORR
#if ALG_RSAES
    TPMS_ENC_SCHEME_RSAES           rsaes;
#endif // ALG_RSAES
#if ALG_OAEP
    TPMS_ENC_SCHEME_OAEP            oaep;
#endif // ALG_OAEP
    TPMS_SCHEME_HASH                anySig;
} TPMU_ASYM_SCHEME;                                 /* Structure */

typedef struct {                                    // Table 2:165
    TPMI_ALG_ASYM_SCHEME        scheme;
    TPMU_ASYM_SCHEME            details;
} TPMT_ASYM_SCHEME;                                 /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_RSA_SCHEME;    // Table 2:166  /* Interface */

typedef struct {                                    // Table 2:167
    TPMI_ALG_RSA_SCHEME         scheme;
    TPMU_ASYM_SCHEME            details;
} TPMT_RSA_SCHEME;                                  /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_RSA_DECRYPT;   // Table 2:168  /* Interface */

typedef struct {                                    // Table 2:169
    TPMI_ALG_RSA_DECRYPT        scheme;
    TPMU_ASYM_SCHEME            details;
} TPMT_RSA_DECRYPT;                                 /* Structure */

typedef union {                                     // Table 2:170
    struct {
        UINT16              size;
        BYTE                buffer[MAX_RSA_KEY_BYTES];
    }            t;
    TPM2B        b;
} TPM2B_PUBLIC_KEY_RSA;                             /* Structure */

typedef TPM_KEY_BITS        TPMI_RSA_KEY_BITS;      // Table 2:171  /* Interface */

typedef union {                                     // Table 2:172
    struct {
        UINT16              size;
        BYTE                buffer[RSA_PRIVATE_SIZE];
    }            t;
    TPM2B        b;
} TPM2B_PRIVATE_KEY_RSA;                            /* Structure */

typedef union {                                     // Table 2:173
    struct {
        UINT16              size;
        BYTE                buffer[MAX_ECC_KEY_BYTES];
    }            t;
    TPM2B        b;
} TPM2B_ECC_PARAMETER;                              /* Structure */

typedef struct {                                    // Table 2:174
    TPM2B_ECC_PARAMETER         x;
    TPM2B_ECC_PARAMETER         y;
} TPMS_ECC_POINT;                                   /* Structure */

typedef struct {                                    // Table 2:175
    UINT16                  size;
    TPMS_ECC_POINT          point;
} TPM2B_ECC_POINT;                                  /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_ECC_SCHEME;    // Table 2:176  /* Interface */

typedef TPM_ECC_CURVE       TPMI_ECC_CURVE;         // Table 2:177  /* Interface */

typedef struct {                                    // Table 2:178
    TPMI_ALG_ECC_SCHEME         scheme;
    TPMU_ASYM_SCHEME            details;
} TPMT_ECC_SCHEME;                                  /* Structure */

typedef struct {                                    // Table 2:179
    TPM_ECC_CURVE               curveID;
    UINT16                      keySize;
    TPMT_KDF_SCHEME             kdf;
    TPMT_ECC_SCHEME             sign;
    TPM2B_ECC_PARAMETER         p;
    TPM2B_ECC_PARAMETER         a;
    TPM2B_ECC_PARAMETER         b;
    TPM2B_ECC_PARAMETER         gX;
    TPM2B_ECC_PARAMETER         gY;
    TPM2B_ECC_PARAMETER         n;
    TPM2B_ECC_PARAMETER         h;
} TPMS_ALGORITHM_DETAIL_ECC;                        /* Structure */

typedef struct {                                    // Table 2:180
    TPMI_ALG_HASH               hash;
    TPM2B_PUBLIC_KEY_RSA        sig;
} TPMS_SIGNATURE_RSA;                               /* Structure */

// Table 2:181 - Definition of Types for Signature
typedef TPMS_SIGNATURE_RSA  TPMS_SIGNATURE_RSASSA;
typedef TPMS_SIGNATURE_RSA  TPMS_SIGNATURE_RSAPSS;

typedef struct {                                    // Table 2:182
    TPMI_ALG_HASH               hash;
    TPM2B_ECC_PARAMETER         signatureR;
    TPM2B_ECC_PARAMETER         signatureS;
} TPMS_SIGNATURE_ECC;                               /* Structure */

// Table 2:183 - Definition of Types for TPMS_SIGNATURE_ECC
typedef TPMS_SIGNATURE_ECC  TPMS_SIGNATURE_ECDAA;
typedef TPMS_SIGNATURE_ECC  TPMS_SIGNATURE_ECDSA;
typedef TPMS_SIGNATURE_ECC  TPMS_SIGNATURE_SM2;
typedef TPMS_SIGNATURE_ECC  TPMS_SIGNATURE_ECSCHNORR;

typedef union {                                     // Table 2:184
#if ALG_ECC
    TPMS_SIGNATURE_ECDAA            ecdaa;
#endif // ALG_ECC
#if ALG_RSA
    TPMS_SIGNATURE_RSASSA           rsassa;
#endif // ALG_RSA
#if ALG_RSA
    TPMS_SIGNATURE_RSAPSS           rsapss;
#endif // ALG_RSA
#if ALG_ECC
    TPMS_SIGNATURE_ECDSA            ecdsa;
#endif // ALG_ECC
#if ALG_ECC
    TPMS_SIGNATURE_SM2              sm2;
#endif // ALG_ECC
#if ALG_ECC
    TPMS_SIGNATURE_ECSCHNORR        ecschnorr;
#endif // ALG_ECC
#if ALG_HMAC
    TPMT_HA                         hmac;
#endif // ALG_HMAC
    TPMS_SCHEME_HASH                any;
} TPMU_SIGNATURE;                                   /* Structure */

typedef struct {                                    // Table 2:185
    TPMI_ALG_SIG_SCHEME         sigAlg;
    TPMU_SIGNATURE              signature;
} TPMT_SIGNATURE;                                   /* Structure */

typedef union {                                     // Table 2:186
#if ALG_ECC
    BYTE                    ecc[sizeof(TPMS_ECC_POINT)];
#endif // ALG_ECC
#if ALG_RSA
    BYTE                    rsa[MAX_RSA_KEY_BYTES];
#endif // ALG_RSA
#if ALG_SYMCIPHER
    BYTE                    symmetric[sizeof(TPM2B_DIGEST)];
#endif // ALG_SYMCIPHER
#if ALG_KEYEDHASH
    BYTE                    keyedHash[sizeof(TPM2B_DIGEST)];
#endif // ALG_KEYEDHASH
} TPMU_ENCRYPTED_SECRET;                            /* Structure */

typedef union {                                     // Table 2:187
    struct {
        UINT16              size;
        BYTE                secret[sizeof(TPMU_ENCRYPTED_SECRET)];
    }            t;
    TPM2B        b;
} TPM2B_ENCRYPTED_SECRET;                           /* Structure */

typedef TPM_ALG_ID          TPMI_ALG_PUBLIC;        // Table 2:188  /* Interface */

typedef union {                                     // Table 2:189
#if ALG_KEYEDHASH
    TPM2B_DIGEST                keyedHash;
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
    TPM2B_DIGEST                sym;
#endif // ALG_SYMCIPHER
#if ALG_RSA
    TPM2B_PUBLIC_KEY_RSA        rsa;
#endif // ALG_RSA
#if ALG_ECC
    TPMS_ECC_POINT              ecc;
#endif // ALG_ECC
    TPMS_DERIVE                 derive;
} TPMU_PUBLIC_ID;                                   /* Structure */

typedef struct {                                    // Table 2:190
    TPMT_KEYEDHASH_SCHEME       scheme;
} TPMS_KEYEDHASH_PARMS;                             /* Structure */

typedef struct {                                    // Table 2:191
    TPMT_SYM_DEF_OBJECT         symmetric;
    TPMT_ASYM_SCHEME            scheme;
} TPMS_ASYM_PARMS;                                  /* Structure */

typedef struct {                                    // Table 2:192
    TPMT_SYM_DEF_OBJECT         symmetric;
    TPMT_RSA_SCHEME             scheme;
    TPMI_RSA_KEY_BITS           keyBits;
    UINT32                      exponent;
} TPMS_RSA_PARMS;                                   /* Structure */

typedef struct {                                    // Table 2:193
    TPMT_SYM_DEF_OBJECT         symmetric;
    TPMT_ECC_SCHEME             scheme;
    TPMI_ECC_CURVE              curveID;
    TPMT_KDF_SCHEME             kdf;
} TPMS_ECC_PARMS;                                   /* Structure */

typedef union {                                     // Table 2:194
#if ALG_KEYEDHASH
    TPMS_KEYEDHASH_PARMS        keyedHashDetail;
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
    TPMS_SYMCIPHER_PARMS        symDetail;
#endif // ALG_SYMCIPHER
#if ALG_RSA
    TPMS_RSA_PARMS              rsaDetail;
#endif // ALG_RSA
#if ALG_ECC
    TPMS_ECC_PARMS              eccDetail;
#endif // ALG_ECC
    TPMS_ASYM_PARMS             asymDetail;
} TPMU_PUBLIC_PARMS;                                /* Structure */

typedef struct {                                    // Table 2:195
    TPMI_ALG_PUBLIC         type;
    TPMU_PUBLIC_PARMS       parameters;
} TPMT_PUBLIC_PARMS;                                /* Structure */

typedef struct {                                    // Table 2:196
    TPMI_ALG_PUBLIC         type;
    TPMI_ALG_HASH           nameAlg;
    TPMA_OBJECT             objectAttributes;
    TPM2B_DIGEST            authPolicy;
    TPMU_PUBLIC_PARMS       parameters;
    TPMU_PUBLIC_ID          unique;
} TPMT_PUBLIC;                                      /* Structure */

typedef struct {                                    // Table 2:197
    UINT16                  size;
    TPMT_PUBLIC             publicArea;
} TPM2B_PUBLIC;                                     /* Structure */

typedef union {                                     // Table 2:198
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMT_PUBLIC)];
    }            t;
    TPM2B        b;
} TPM2B_TEMPLATE;                                   /* Structure */

typedef union {                                     // Table 2:199
    struct {
        UINT16              size;
        BYTE                buffer[PRIVATE_VENDOR_SPECIFIC_BYTES];
    }            t;
    TPM2B        b;
} TPM2B_PRIVATE_VENDOR_SPECIFIC;                    /* Structure */

typedef union {                                     // Table 2:200
#if ALG_RSA
    TPM2B_PRIVATE_KEY_RSA               rsa;
#endif // ALG_RSA
#if ALG_ECC
    TPM2B_ECC_PARAMETER                 ecc;
#endif // ALG_ECC
#if ALG_KEYEDHASH
    TPM2B_SENSITIVE_DATA                bits;
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
    TPM2B_SYM_KEY                       sym;
#endif // ALG_SYMCIPHER
    TPM2B_PRIVATE_VENDOR_SPECIFIC       any;
} TPMU_SENSITIVE_COMPOSITE;                         /* Structure */

typedef struct {                                    // Table 2:201
    TPMI_ALG_PUBLIC                 sensitiveType;
    TPM2B_AUTH                      authValue;
    TPM2B_DIGEST                    seedValue;
    TPMU_SENSITIVE_COMPOSITE        sensitive;
} TPMT_SENSITIVE;                                   /* Structure */

typedef struct {                                    // Table 2:202
    UINT16                  size;
    TPMT_SENSITIVE          sensitiveArea;
} TPM2B_SENSITIVE;                                  /* Structure */

typedef struct {                                    // Table 2:203
    TPM2B_DIGEST            integrityOuter;
    TPM2B_DIGEST            integrityInner;
    TPM2B_SENSITIVE         sensitive;
} _PRIVATE;                                         /* Structure */

typedef union {                                     // Table 2:204
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(_PRIVATE)];
    }            t;
    TPM2B        b;
} TPM2B_PRIVATE;                                    /* Structure */

typedef struct {                                    // Table 2:205
    TPM2B_DIGEST            integrityHMAC;
    TPM2B_DIGEST            encIdentity;
} TPMS_ID_OBJECT;                                   /* Structure */

typedef union {                                     // Table 2:206
    struct {
        UINT16              size;
        BYTE                credential[sizeof(TPMS_ID_OBJECT)];
    }            t;
    TPM2B        b;
} TPM2B_ID_OBJECT;                                  /* Structure */

#define TYPE_OF_TPM_NV_INDEX    UINT32
#define TPM_NV_INDEX_TO_UINT32(a)    (*((UINT32 *)&(a)))
#define UINT32_TO_TPM_NV_INDEX(a)    (*((TPM_NV_INDEX *)&(a)))
#define TPM_NV_INDEX_TO_BYTE_ARRAY(i, a)                                           \
            UINT32_TO_BYTE_ARRAY((TPM_NV_INDEX_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPM_NV_INDEX(i, a)                                           \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPM_NV_INDEX(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPM_NV_INDEX {                       // Table 2:207
    unsigned    index                : 24;
    unsigned    RH_NV                : 8;
} TPM_NV_INDEX;                                     /* Bits */
// This is the initializer for a TPM_NV_INDEX structure
#define TPM_NV_INDEX_INITIALIZER(index, rh_nv) {index, rh_nv}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:207 TPM_NV_INDEX using bit masking
typedef UINT32                      TPM_NV_INDEX;
#define TYPE_OF_TPM_NV_INDEX        UINT32
#define TPM_NV_INDEX_index_SHIFT    0
#define TPM_NV_INDEX_index          ((TPM_NV_INDEX)0xffffff << 0)
#define TPM_NV_INDEX_RH_NV_SHIFT    24
#define TPM_NV_INDEX_RH_NV          ((TPM_NV_INDEX)0xff << 24)
//  This is the initializer for a TPM_NV_INDEX bit array.
#define TPM_NV_INDEX_INITIALIZER(index, rh_nv) {(index << 0) + (rh_nv << 24)}
#endif // USE_BIT_FIELD_STRUCTURES

// Table 2:208 - Definition of TPM_NT Constants
typedef UINT32              TPM_NT;
#define TYPE_OF_TPM_NT      UINT32
#define TPM_NT_ORDINARY     (TPM_NT)(0x0)
#define TPM_NT_COUNTER      (TPM_NT)(0x1)
#define TPM_NT_BITS         (TPM_NT)(0x2)
#define TPM_NT_EXTEND       (TPM_NT)(0x4)
#define TPM_NT_PIN_FAIL     (TPM_NT)(0x8)
#define TPM_NT_PIN_PASS     (TPM_NT)(0x9)

typedef struct {                                    // Table 2:209
    UINT32                  pinCount;
    UINT32                  pinLimit;
} TPMS_NV_PIN_COUNTER_PARAMETERS;                   /* Structure */

#define TYPE_OF_TPMA_NV     UINT32
#define TPMA_NV_TO_UINT32(a)     (*((UINT32 *)&(a)))
#define UINT32_TO_TPMA_NV(a)     (*((TPMA_NV *)&(a)))
#define TPMA_NV_TO_BYTE_ARRAY(i, a)                                                \
            UINT32_TO_BYTE_ARRAY((TPMA_NV_TO_UINT32(i)), (a))
#define BYTE_ARRAY_TO_TPMA_NV(i, a)                                                \
            { UINT32 x = BYTE_ARRAY_TO_UINT32(a); i = UINT32_TO_TPMA_NV(x); }
#if USE_BIT_FIELD_STRUCTURES
typedef struct TPMA_NV {                            // Table 2:210
    unsigned    PPWRITE              : 1;
    unsigned    OWNERWRITE           : 1;
    unsigned    AUTHWRITE            : 1;
    unsigned    POLICYWRITE          : 1;
    unsigned    TPM_NT               : 4;
    unsigned    Reserved_bits_at_8   : 2;
    unsigned    POLICY_DELETE        : 1;
    unsigned    WRITELOCKED          : 1;
    unsigned    WRITEALL             : 1;
    unsigned    WRITEDEFINE          : 1;
    unsigned    WRITE_STCLEAR        : 1;
    unsigned    GLOBALLOCK           : 1;
    unsigned    PPREAD               : 1;
    unsigned    OWNERREAD            : 1;
    unsigned    AUTHREAD             : 1;
    unsigned    POLICYREAD           : 1;
    unsigned    Reserved_bits_at_20  : 5;
    unsigned    NO_DA                : 1;
    unsigned    ORDERLY              : 1;
    unsigned    CLEAR_STCLEAR        : 1;
    unsigned    READLOCKED           : 1;
    unsigned    WRITTEN              : 1;
    unsigned    PLATFORMCREATE       : 1;
    unsigned    READ_STCLEAR         : 1;
} TPMA_NV;                                          /* Bits */
// This is the initializer for a TPMA_NV structure
#define TPMA_NV_INITIALIZER(                                                       \
             ppwrite,        ownerwrite,     authwrite,      policywrite,          \
             tpm_nt,         bits_at_8,      policy_delete,  writelocked,          \
             writeall,       writedefine,    write_stclear,  globallock,           \
             ppread,         ownerread,      authread,       policyread,           \
             bits_at_20,     no_da,          orderly,        clear_stclear,        \
             readlocked,     written,        platformcreate, read_stclear)         \
            {ppwrite,        ownerwrite,     authwrite,      policywrite,          \
             tpm_nt,         bits_at_8,      policy_delete,  writelocked,          \
             writeall,       writedefine,    write_stclear,  globallock,           \
             ppread,         ownerread,      authread,       policyread,           \
             bits_at_20,     no_da,          orderly,        clear_stclear,        \
             readlocked,     written,        platformcreate, read_stclear}
#else // USE_BIT_FIELD_STRUCTURES
// This implements Table 2:210 TPMA_NV using bit masking
typedef UINT32                  TPMA_NV;
#define TYPE_OF_TPMA_NV         UINT32
#define TPMA_NV_PPWRITE         ((TPMA_NV)1 << 0)
#define TPMA_NV_OWNERWRITE      ((TPMA_NV)1 << 1)
#define TPMA_NV_AUTHWRITE       ((TPMA_NV)1 << 2)
#define TPMA_NV_POLICYWRITE     ((TPMA_NV)1 << 3)
#define TPMA_NV_TPM_NT_SHIFT    4
#define TPMA_NV_TPM_NT          ((TPMA_NV)0xf << 4)
#define TPMA_NV_POLICY_DELETE   ((TPMA_NV)1 << 10)
#define TPMA_NV_WRITELOCKED     ((TPMA_NV)1 << 11)
#define TPMA_NV_WRITEALL        ((TPMA_NV)1 << 12)
#define TPMA_NV_WRITEDEFINE     ((TPMA_NV)1 << 13)
#define TPMA_NV_WRITE_STCLEAR   ((TPMA_NV)1 << 14)
#define TPMA_NV_GLOBALLOCK      ((TPMA_NV)1 << 15)
#define TPMA_NV_PPREAD          ((TPMA_NV)1 << 16)
#define TPMA_NV_OWNERREAD       ((TPMA_NV)1 << 17)
#define TPMA_NV_AUTHREAD        ((TPMA_NV)1 << 18)
#define TPMA_NV_POLICYREAD      ((TPMA_NV)1 << 19)
#define TPMA_NV_NO_DA           ((TPMA_NV)1 << 25)
#define TPMA_NV_ORDERLY         ((TPMA_NV)1 << 26)
#define TPMA_NV_CLEAR_STCLEAR   ((TPMA_NV)1 << 27)
#define TPMA_NV_READLOCKED      ((TPMA_NV)1 << 28)
#define TPMA_NV_WRITTEN         ((TPMA_NV)1 << 29)
#define TPMA_NV_PLATFORMCREATE  ((TPMA_NV)1 << 30)
#define TPMA_NV_READ_STCLEAR    ((TPMA_NV)1 << 31)
//  This is the initializer for a TPMA_NV bit array.
#define TPMA_NV_INITIALIZER(                                                       \
             ppwrite,        ownerwrite,     authwrite,      policywrite,          \
             tpm_nt,         bits_at_8,      policy_delete,  writelocked,          \
             writeall,       writedefine,    write_stclear,  globallock,           \
             ppread,         ownerread,      authread,       policyread,           \
             bits_at_20,     no_da,          orderly,        clear_stclear,        \
             readlocked,     written,        platformcreate, read_stclear)         \
            {(ppwrite << 0)         + (ownerwrite << 1)      +                     \
             (authwrite << 2)       + (policywrite << 3)     +                     \
             (tpm_nt << 4)          + (policy_delete << 10)  +                     \
             (writelocked << 11)    + (writeall << 12)       +                     \
             (writedefine << 13)    + (write_stclear << 14)  +                     \
             (globallock << 15)     + (ppread << 16)         +                     \
             (ownerread << 17)      + (authread << 18)       +                     \
             (policyread << 19)     + (no_da << 25)          +                     \
             (orderly << 26)        + (clear_stclear << 27)  +                     \
             (readlocked << 28)     + (written << 29)        +                     \
             (platformcreate << 30) + (read_stclear << 31)}
#endif // USE_BIT_FIELD_STRUCTURES

typedef struct {                                    // Table 2:211
    TPMI_RH_NV_INDEX        nvIndex;
    TPMI_ALG_HASH           nameAlg;
    TPMA_NV                 attributes;
    TPM2B_DIGEST            authPolicy;
    UINT16                  dataSize;
} TPMS_NV_PUBLIC;                                   /* Structure */

typedef struct {                                    // Table 2:212
    UINT16                  size;
    TPMS_NV_PUBLIC          nvPublic;
} TPM2B_NV_PUBLIC;                                  /* Structure */

typedef union {                                     // Table 2:213
    struct {
        UINT16              size;
        BYTE                buffer[MAX_CONTEXT_SIZE];
    }            t;
    TPM2B        b;
} TPM2B_CONTEXT_SENSITIVE;                          /* Structure */

typedef struct {                                    // Table 2:214
    TPM2B_DIGEST                    integrity;
    TPM2B_CONTEXT_SENSITIVE         encrypted;
} TPMS_CONTEXT_DATA;                                /* Structure */

typedef union {                                     // Table 2:215
    struct {
        UINT16              size;
        BYTE                buffer[sizeof(TPMS_CONTEXT_DATA)];
    }            t;
    TPM2B        b;
} TPM2B_CONTEXT_DATA;                               /* Structure */

typedef struct {                                    // Table 2:216
    UINT64                  sequence;
    TPMI_DH_SAVED           savedHandle;
    TPMI_RH_HIERARCHY       hierarchy;
    TPM2B_CONTEXT_DATA      contextBlob;
} TPMS_CONTEXT;                                     /* Structure */

typedef struct {                                    // Table 2:218
    TPML_PCR_SELECTION      pcrSelect;
    TPM2B_DIGEST            pcrDigest;
    TPMA_LOCALITY           locality;
    TPM_ALG_ID              parentNameAlg;
    TPM2B_NAME              parentName;
    TPM2B_NAME              parentQualifiedName;
    TPM2B_DATA              outsideInfo;
} TPMS_CREATION_DATA;                               /* Structure */

typedef struct {                                    // Table 2:219
    UINT16                  size;
    TPMS_CREATION_DATA      creationData;
} TPM2B_CREATION_DATA;                              /* Structure */

// Table 2:220 - Definition of TPM_AT Constants
typedef UINT32              TPM_AT;
#define TYPE_OF_TPM_AT      UINT32
#define TPM_AT_ANY          (TPM_AT)(0x00000000)
#define TPM_AT_ERROR        (TPM_AT)(0x00000001)
#define TPM_AT_PV1          (TPM_AT)(0x00000002)
#define TPM_AT_VEND         (TPM_AT)(0x80000000)

// Table 2:221 - Definition of TPM_AE Constants
typedef UINT32              TPM_AE;
#define TYPE_OF_TPM_AE      UINT32
#define TPM_AE_NONE         (TPM_AE)(0x00000000)

typedef struct {                                    // Table 2:222
    TPM_AT                  tag;
    UINT32                  data;
} TPMS_AC_OUTPUT;                                   /* Structure */

typedef struct {                                    // Table 2:223
    UINT32                  count;
    TPMS_AC_OUTPUT          acCapabilities[MAX_AC_CAPABILITIES];
} TPML_AC_CAPABILITIES;                             /* Structure */



#endif // _TPM_TYPES_H_
