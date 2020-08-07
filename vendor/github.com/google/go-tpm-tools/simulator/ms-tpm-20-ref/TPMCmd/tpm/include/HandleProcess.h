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
 *  Created by TpmDispatch; Version 4.0 July 8,2017
 *  Date: Oct  9, 2018  Time: 07:25:19PM
 */
#if CC_Startup
case TPM_CC_Startup:
    break;
#endif     // CC_Startup
#if CC_Shutdown
case TPM_CC_Shutdown:
    break;
#endif     // CC_Shutdown
#if CC_SelfTest
case TPM_CC_SelfTest:
    break;
#endif     // CC_SelfTest
#if CC_IncrementalSelfTest
case TPM_CC_IncrementalSelfTest:
    break;
#endif     // CC_IncrementalSelfTest
#if CC_GetTestResult
case TPM_CC_GetTestResult:
    break;
#endif     // CC_GetTestResult
#if CC_StartAuthSession
case TPM_CC_StartAuthSession:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_ENTITY_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_StartAuthSession
#if CC_PolicyRestart
case TPM_CC_PolicyRestart:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyRestart
#if CC_Create
case TPM_CC_Create:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Create
#if CC_Load
case TPM_CC_Load:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Load
#if CC_LoadExternal
case TPM_CC_LoadExternal:
    break;
#endif     // CC_LoadExternal
#if CC_ReadPublic
case TPM_CC_ReadPublic:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ReadPublic
#if CC_ActivateCredential
case TPM_CC_ActivateCredential:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_ActivateCredential
#if CC_MakeCredential
case TPM_CC_MakeCredential:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_MakeCredential
#if CC_Unseal
case TPM_CC_Unseal:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Unseal
#if CC_ObjectChangeAuth
case TPM_CC_ObjectChangeAuth:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_ObjectChangeAuth
#if CC_CreateLoaded
case TPM_CC_CreateLoaded:
    *handleCount = 1;
    result = TPMI_DH_PARENT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_CreateLoaded
#if CC_Duplicate
case TPM_CC_Duplicate:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_Duplicate
#if CC_Rewrap
case TPM_CC_Rewrap:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_Rewrap
#if CC_Import
case TPM_CC_Import:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Import
#if CC_RSA_Encrypt
case TPM_CC_RSA_Encrypt:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_RSA_Encrypt
#if CC_RSA_Decrypt
case TPM_CC_RSA_Decrypt:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_RSA_Decrypt
#if CC_ECDH_KeyGen
case TPM_CC_ECDH_KeyGen:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ECDH_KeyGen
#if CC_ECDH_ZGen
case TPM_CC_ECDH_ZGen:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ECDH_ZGen
#if CC_ECC_Parameters
case TPM_CC_ECC_Parameters:
    break;
#endif     // CC_ECC_Parameters
#if CC_ZGen_2Phase
case TPM_CC_ZGen_2Phase:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ZGen_2Phase
#if CC_EncryptDecrypt
case TPM_CC_EncryptDecrypt:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_EncryptDecrypt
#if CC_EncryptDecrypt2
case TPM_CC_EncryptDecrypt2:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_EncryptDecrypt2
#if CC_Hash
case TPM_CC_Hash:
    break;
#endif     // CC_Hash
#if CC_HMAC
case TPM_CC_HMAC:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_HMAC
#if CC_MAC
case TPM_CC_MAC:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_MAC
#if CC_GetRandom
case TPM_CC_GetRandom:
    break;
#endif     // CC_GetRandom
#if CC_StirRandom
case TPM_CC_StirRandom:
    break;
#endif     // CC_StirRandom
#if CC_HMAC_Start
case TPM_CC_HMAC_Start:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_HMAC_Start
#if CC_MAC_Start
case TPM_CC_MAC_Start:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_MAC_Start
#if CC_HashSequenceStart
case TPM_CC_HashSequenceStart:
    break;
#endif     // CC_HashSequenceStart
#if CC_SequenceUpdate
case TPM_CC_SequenceUpdate:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_SequenceUpdate
#if CC_SequenceComplete
case TPM_CC_SequenceComplete:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_SequenceComplete
#if CC_EventSequenceComplete
case TPM_CC_EventSequenceComplete:
    *handleCount = 2;
    result = TPMI_DH_PCR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_EventSequenceComplete
#if CC_Certify
case TPM_CC_Certify:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_Certify
#if CC_CertifyCreation
case TPM_CC_CertifyCreation:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_CertifyCreation
#if CC_Quote
case TPM_CC_Quote:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Quote
#if CC_GetSessionAuditDigest
case TPM_CC_GetSessionAuditDigest:
    *handleCount = 3;
    result = TPMI_RH_ENDORSEMENT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    result = TPMI_SH_HMAC_Unmarshal(&handles[2], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_3;
    break;
#endif     // CC_GetSessionAuditDigest
#if CC_GetCommandAuditDigest
case TPM_CC_GetCommandAuditDigest:
    *handleCount = 2;
    result = TPMI_RH_ENDORSEMENT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_GetCommandAuditDigest
#if CC_GetTime
case TPM_CC_GetTime:
    *handleCount = 2;
    result = TPMI_RH_ENDORSEMENT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_GetTime
#if CC_CertifyX509
case TPM_CC_CertifyX509:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_CertifyX509
#if CC_Commit
case TPM_CC_Commit:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Commit
#if CC_EC_Ephemeral
case TPM_CC_EC_Ephemeral:
    break;
#endif     // CC_EC_Ephemeral
#if CC_VerifySignature
case TPM_CC_VerifySignature:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_VerifySignature
#if CC_Sign
case TPM_CC_Sign:
    *handleCount = 1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Sign
#if CC_SetCommandCodeAuditStatus
case TPM_CC_SetCommandCodeAuditStatus:
    *handleCount = 1;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_SetCommandCodeAuditStatus
#if CC_PCR_Extend
case TPM_CC_PCR_Extend:
    *handleCount = 1;
    result = TPMI_DH_PCR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_Extend
#if CC_PCR_Event
case TPM_CC_PCR_Event:
    *handleCount = 1;
    result = TPMI_DH_PCR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_Event
#if CC_PCR_Read
case TPM_CC_PCR_Read:
    break;
#endif     // CC_PCR_Read
#if CC_PCR_Allocate
case TPM_CC_PCR_Allocate:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_Allocate
#if CC_PCR_SetAuthPolicy
case TPM_CC_PCR_SetAuthPolicy:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_SetAuthPolicy
#if CC_PCR_SetAuthValue
case TPM_CC_PCR_SetAuthValue:
    *handleCount = 1;
    result = TPMI_DH_PCR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_SetAuthValue
#if CC_PCR_Reset
case TPM_CC_PCR_Reset:
    *handleCount = 1;
    result = TPMI_DH_PCR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PCR_Reset
#if CC_PolicySigned
case TPM_CC_PolicySigned:
    *handleCount = 2;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_PolicySigned
#if CC_PolicySecret
case TPM_CC_PolicySecret:
    *handleCount = 2;
    result = TPMI_DH_ENTITY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_PolicySecret
#if CC_PolicyTicket
case TPM_CC_PolicyTicket:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyTicket
#if CC_PolicyOR
case TPM_CC_PolicyOR:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyOR
#if CC_PolicyPCR
case TPM_CC_PolicyPCR:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyPCR
#if CC_PolicyLocality
case TPM_CC_PolicyLocality:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyLocality
#if CC_PolicyNV
case TPM_CC_PolicyNV:
    *handleCount = 3;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    result = TPMI_SH_POLICY_Unmarshal(&handles[2], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_3;
    break;
#endif     // CC_PolicyNV
#if CC_PolicyCounterTimer
case TPM_CC_PolicyCounterTimer:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyCounterTimer
#if CC_PolicyCommandCode
case TPM_CC_PolicyCommandCode:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyCommandCode
#if CC_PolicyPhysicalPresence
case TPM_CC_PolicyPhysicalPresence:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyPhysicalPresence
#if CC_PolicyCpHash
case TPM_CC_PolicyCpHash:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyCpHash
#if CC_PolicyNameHash
case TPM_CC_PolicyNameHash:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyNameHash
#if CC_PolicyDuplicationSelect
case TPM_CC_PolicyDuplicationSelect:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyDuplicationSelect
#if CC_PolicyAuthorize
case TPM_CC_PolicyAuthorize:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyAuthorize
#if CC_PolicyAuthValue
case TPM_CC_PolicyAuthValue:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyAuthValue
#if CC_PolicyPassword
case TPM_CC_PolicyPassword:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyPassword
#if CC_PolicyGetDigest
case TPM_CC_PolicyGetDigest:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyGetDigest
#if CC_PolicyNvWritten
case TPM_CC_PolicyNvWritten:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyNvWritten
#if CC_PolicyTemplate
case TPM_CC_PolicyTemplate:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PolicyTemplate
#if CC_PolicyAuthorizeNV
case TPM_CC_PolicyAuthorizeNV:
    *handleCount = 3;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    result = TPMI_SH_POLICY_Unmarshal(&handles[2], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_3;
    break;
#endif     // CC_PolicyAuthorizeNV
#if CC_CreatePrimary
case TPM_CC_CreatePrimary:
    *handleCount = 1;
    result = TPMI_RH_HIERARCHY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_CreatePrimary
#if CC_HierarchyControl
case TPM_CC_HierarchyControl:
    *handleCount = 1;
    result = TPMI_RH_HIERARCHY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_HierarchyControl
#if CC_SetPrimaryPolicy
case TPM_CC_SetPrimaryPolicy:
    *handleCount = 1;
    result = TPMI_RH_HIERARCHY_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_SetPrimaryPolicy
#if CC_ChangePPS
case TPM_CC_ChangePPS:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ChangePPS
#if CC_ChangeEPS
case TPM_CC_ChangeEPS:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ChangeEPS
#if CC_Clear
case TPM_CC_Clear:
    *handleCount = 1;
    result = TPMI_RH_CLEAR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Clear
#if CC_ClearControl
case TPM_CC_ClearControl:
    *handleCount = 1;
    result = TPMI_RH_CLEAR_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ClearControl
#if CC_HierarchyChangeAuth
case TPM_CC_HierarchyChangeAuth:
    *handleCount = 1;
    result = TPMI_RH_HIERARCHY_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_HierarchyChangeAuth
#if CC_DictionaryAttackLockReset
case TPM_CC_DictionaryAttackLockReset:
    *handleCount = 1;
    result = TPMI_RH_LOCKOUT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_DictionaryAttackLockReset
#if CC_DictionaryAttackParameters
case TPM_CC_DictionaryAttackParameters:
    *handleCount = 1;
    result = TPMI_RH_LOCKOUT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_DictionaryAttackParameters
#if CC_PP_Commands
case TPM_CC_PP_Commands:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_PP_Commands
#if CC_SetAlgorithmSet
case TPM_CC_SetAlgorithmSet:
    *handleCount = 1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_SetAlgorithmSet
#if CC_FieldUpgradeStart
case TPM_CC_FieldUpgradeStart:
    *handleCount = 2;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_FieldUpgradeStart
#if CC_FieldUpgradeData
case TPM_CC_FieldUpgradeData:
    break;
#endif     // CC_FieldUpgradeData
#if CC_FirmwareRead
case TPM_CC_FirmwareRead:
    break;
#endif     // CC_FirmwareRead
#if CC_ContextSave
case TPM_CC_ContextSave:
    *handleCount = 1;
    result = TPMI_DH_CONTEXT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ContextSave
#if CC_ContextLoad
case TPM_CC_ContextLoad:
    break;
#endif     // CC_ContextLoad
#if CC_FlushContext
case TPM_CC_FlushContext:
    break;
#endif     // CC_FlushContext
#if CC_EvictControl
case TPM_CC_EvictControl:
    *handleCount = 2;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_EvictControl
#if CC_ReadClock
case TPM_CC_ReadClock:
    break;
#endif     // CC_ReadClock
#if CC_ClockSet
case TPM_CC_ClockSet:
    *handleCount = 1;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ClockSet
#if CC_ClockRateAdjust
case TPM_CC_ClockRateAdjust:
    *handleCount = 1;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_ClockRateAdjust
#if CC_GetCapability
case TPM_CC_GetCapability:
    break;
#endif     // CC_GetCapability
#if CC_TestParms
case TPM_CC_TestParms:
    break;
#endif     // CC_TestParms
#if CC_NV_DefineSpace
case TPM_CC_NV_DefineSpace:
    *handleCount = 1;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_NV_DefineSpace
#if CC_NV_UndefineSpace
case TPM_CC_NV_UndefineSpace:
    *handleCount = 2;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_UndefineSpace
#if CC_NV_UndefineSpaceSpecial
case TPM_CC_NV_UndefineSpaceSpecial:
    *handleCount = 2;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_PLATFORM_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_UndefineSpaceSpecial
#if CC_NV_ReadPublic
case TPM_CC_NV_ReadPublic:
    *handleCount = 1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_NV_ReadPublic
#if CC_NV_Write
case TPM_CC_NV_Write:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_Write
#if CC_NV_Increment
case TPM_CC_NV_Increment:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_Increment
#if CC_NV_Extend
case TPM_CC_NV_Extend:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_Extend
#if CC_NV_SetBits
case TPM_CC_NV_SetBits:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_SetBits
#if CC_NV_WriteLock
case TPM_CC_NV_WriteLock:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_WriteLock
#if CC_NV_GlobalWriteLock
case TPM_CC_NV_GlobalWriteLock:
    *handleCount = 1;
    result = TPMI_RH_PROVISION_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_NV_GlobalWriteLock
#if CC_NV_Read
case TPM_CC_NV_Read:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_Read
#if CC_NV_ReadLock
case TPM_CC_NV_ReadLock:
    *handleCount = 2;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    break;
#endif     // CC_NV_ReadLock
#if CC_NV_ChangeAuth
case TPM_CC_NV_ChangeAuth:
    *handleCount = 1;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_NV_ChangeAuth
#if CC_NV_Certify
case TPM_CC_NV_Certify:
    *handleCount = 3;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, TRUE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    result = TPMI_RH_NV_INDEX_Unmarshal(&handles[2], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_3;
    break;
#endif     // CC_NV_Certify
#if CC_AC_GetCapability
case TPM_CC_AC_GetCapability:
    *handleCount = 1;
    result = TPMI_RH_AC_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_AC_GetCapability
#if CC_AC_Send
case TPM_CC_AC_Send:
    *handleCount = 3;
    result = TPMI_DH_OBJECT_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize, FALSE);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    result = TPMI_RH_NV_AUTH_Unmarshal(&handles[1], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_2;
    result = TPMI_RH_AC_Unmarshal(&handles[2], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_3;
    break;
#endif     // CC_AC_Send
#if CC_Policy_AC_SendSelect
case TPM_CC_Policy_AC_SendSelect:
    *handleCount = 1;
    result = TPMI_SH_POLICY_Unmarshal(&handles[0], handleBufferStart, 
                                              bufferRemainingSize);
    if(TPM_RC_SUCCESS != result) return result + TPM_RC_H + TPM_RC_1;
    break;
#endif     // CC_Policy_AC_SendSelect
#if CC_Vendor_TCG_Test
case TPM_CC_Vendor_TCG_Test:
    break;
#endif     // CC_Vendor_TCG_Test
