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
 *  Date: Oct 27, 2018  Time: 06:49:39PM
 */

// This macro is added just so that the code is only excessively long.
#define EXIT_IF_ERROR_PLUS(x)         \
    if(TPM_RC_SUCCESS != result) { result += (x); goto Exit; }
#if CC_Startup
case TPM_CC_Startup: {
    Startup_In *in = (Startup_In *)
            MemoryGetInBuffer(sizeof(Startup_In));
    result = TPM_SU_Unmarshal(&in->startupType, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Startup_startupType);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Startup (in);
break; 
}
#endif     // CC_Startup
#if CC_Shutdown
case TPM_CC_Shutdown: {
    Shutdown_In *in = (Shutdown_In *)
            MemoryGetInBuffer(sizeof(Shutdown_In));
    result = TPM_SU_Unmarshal(&in->shutdownType, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Shutdown_shutdownType);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Shutdown (in);
break; 
}
#endif     // CC_Shutdown
#if CC_SelfTest
case TPM_CC_SelfTest: {
    SelfTest_In *in = (SelfTest_In *)
            MemoryGetInBuffer(sizeof(SelfTest_In));
    result = TPMI_YES_NO_Unmarshal(&in->fullTest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SelfTest_fullTest);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SelfTest (in);
break; 
}
#endif     // CC_SelfTest
#if CC_IncrementalSelfTest
case TPM_CC_IncrementalSelfTest: {
    IncrementalSelfTest_In *in = (IncrementalSelfTest_In *)
            MemoryGetInBuffer(sizeof(IncrementalSelfTest_In));
    IncrementalSelfTest_Out *out = (IncrementalSelfTest_Out *) 
            MemoryGetOutBuffer(sizeof(IncrementalSelfTest_Out));
    result = TPML_ALG_Unmarshal(&in->toTest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_IncrementalSelfTest_toTest);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_IncrementalSelfTest (in, out);
    rSize = sizeof(IncrementalSelfTest_Out);
    *respParmSize += TPML_ALG_Marshal(&out->toDoList, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_IncrementalSelfTest
#if CC_GetTestResult
case TPM_CC_GetTestResult: {
    GetTestResult_Out *out = (GetTestResult_Out *) 
            MemoryGetOutBuffer(sizeof(GetTestResult_Out));
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetTestResult (out);
    rSize = sizeof(GetTestResult_Out);
    *respParmSize += TPM2B_MAX_BUFFER_Marshal(&out->outData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM_RC_Marshal(&out->testResult, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetTestResult
#if CC_StartAuthSession
case TPM_CC_StartAuthSession: {
    StartAuthSession_In *in = (StartAuthSession_In *)
            MemoryGetInBuffer(sizeof(StartAuthSession_In));
    StartAuthSession_Out *out = (StartAuthSession_Out *) 
            MemoryGetOutBuffer(sizeof(StartAuthSession_Out));
    in->tpmKey = handles[0];
    in->bind = handles[1];
    result = TPM2B_NONCE_Unmarshal(&in->nonceCaller, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_StartAuthSession_nonceCaller);
    result = TPM2B_ENCRYPTED_SECRET_Unmarshal(&in->encryptedSalt, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_StartAuthSession_encryptedSalt);
    result = TPM_SE_Unmarshal(&in->sessionType, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_StartAuthSession_sessionType);
    result = TPMT_SYM_DEF_Unmarshal(&in->symmetric, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_StartAuthSession_symmetric);
    result = TPMI_ALG_HASH_Unmarshal(&in->authHash, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_StartAuthSession_authHash);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_StartAuthSession (in, out);
    rSize = sizeof(StartAuthSession_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->sessionHandle;
    *respParmSize += TPM2B_NONCE_Marshal(&out->nonceTPM, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_StartAuthSession
#if CC_PolicyRestart
case TPM_CC_PolicyRestart: {
    PolicyRestart_In *in = (PolicyRestart_In *)
            MemoryGetInBuffer(sizeof(PolicyRestart_In));
    in->sessionHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyRestart (in);
break; 
}
#endif     // CC_PolicyRestart
#if CC_Create
case TPM_CC_Create: {
    Create_In *in = (Create_In *)
            MemoryGetInBuffer(sizeof(Create_In));
    Create_Out *out = (Create_Out *) 
            MemoryGetOutBuffer(sizeof(Create_Out));
    in->parentHandle = handles[0];
    result = TPM2B_SENSITIVE_CREATE_Unmarshal(&in->inSensitive, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Create_inSensitive);
    result = TPM2B_PUBLIC_Unmarshal(&in->inPublic, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_Create_inPublic);
    result = TPM2B_DATA_Unmarshal(&in->outsideInfo, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Create_outsideInfo);
    result = TPML_PCR_SELECTION_Unmarshal(&in->creationPCR, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Create_creationPCR);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Create (in, out);
    rSize = sizeof(Create_Out);
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->outPrivate, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_PUBLIC_Marshal(&out->outPublic, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_CREATION_DATA_Marshal(&out->creationData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->creationHash, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_CREATION_Marshal(&out->creationTicket, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Create
#if CC_Load
case TPM_CC_Load: {
    Load_In *in = (Load_In *)
            MemoryGetInBuffer(sizeof(Load_In));
    Load_Out *out = (Load_Out *) 
            MemoryGetOutBuffer(sizeof(Load_Out));
    in->parentHandle = handles[0];
    result = TPM2B_PRIVATE_Unmarshal(&in->inPrivate, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Load_inPrivate);
    result = TPM2B_PUBLIC_Unmarshal(&in->inPublic, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_Load_inPublic);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Load (in, out);
    rSize = sizeof(Load_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->objectHandle;
    *respParmSize += TPM2B_NAME_Marshal(&out->name, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Load
#if CC_LoadExternal
case TPM_CC_LoadExternal: {
    LoadExternal_In *in = (LoadExternal_In *)
            MemoryGetInBuffer(sizeof(LoadExternal_In));
    LoadExternal_Out *out = (LoadExternal_Out *) 
            MemoryGetOutBuffer(sizeof(LoadExternal_Out));
    result = TPM2B_SENSITIVE_Unmarshal(&in->inPrivate, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_LoadExternal_inPrivate);
    result = TPM2B_PUBLIC_Unmarshal(&in->inPublic, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_LoadExternal_inPublic);
    result = TPMI_RH_HIERARCHY_Unmarshal(&in->hierarchy, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_LoadExternal_hierarchy);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_LoadExternal (in, out);
    rSize = sizeof(LoadExternal_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->objectHandle;
    *respParmSize += TPM2B_NAME_Marshal(&out->name, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_LoadExternal
#if CC_ReadPublic
case TPM_CC_ReadPublic: {
    ReadPublic_In *in = (ReadPublic_In *)
            MemoryGetInBuffer(sizeof(ReadPublic_In));
    ReadPublic_Out *out = (ReadPublic_Out *) 
            MemoryGetOutBuffer(sizeof(ReadPublic_Out));
    in->objectHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ReadPublic (in, out);
    rSize = sizeof(ReadPublic_Out);
    *respParmSize += TPM2B_PUBLIC_Marshal(&out->outPublic, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_NAME_Marshal(&out->name, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_NAME_Marshal(&out->qualifiedName, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ReadPublic
#if CC_ActivateCredential
case TPM_CC_ActivateCredential: {
    ActivateCredential_In *in = (ActivateCredential_In *)
            MemoryGetInBuffer(sizeof(ActivateCredential_In));
    ActivateCredential_Out *out = (ActivateCredential_Out *) 
            MemoryGetOutBuffer(sizeof(ActivateCredential_Out));
    in->activateHandle = handles[0];
    in->keyHandle = handles[1];
    result = TPM2B_ID_OBJECT_Unmarshal(&in->credentialBlob, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ActivateCredential_credentialBlob);
    result = TPM2B_ENCRYPTED_SECRET_Unmarshal(&in->secret, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ActivateCredential_secret);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ActivateCredential (in, out);
    rSize = sizeof(ActivateCredential_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->certInfo, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ActivateCredential
#if CC_MakeCredential
case TPM_CC_MakeCredential: {
    MakeCredential_In *in = (MakeCredential_In *)
            MemoryGetInBuffer(sizeof(MakeCredential_In));
    MakeCredential_Out *out = (MakeCredential_Out *) 
            MemoryGetOutBuffer(sizeof(MakeCredential_Out));
    in->handle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->credential, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_MakeCredential_credential);
    result = TPM2B_NAME_Unmarshal(&in->objectName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_MakeCredential_objectName);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_MakeCredential (in, out);
    rSize = sizeof(MakeCredential_Out);
    *respParmSize += TPM2B_ID_OBJECT_Marshal(&out->credentialBlob, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ENCRYPTED_SECRET_Marshal(&out->secret, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_MakeCredential
#if CC_Unseal
case TPM_CC_Unseal: {
    Unseal_In *in = (Unseal_In *)
            MemoryGetInBuffer(sizeof(Unseal_In));
    Unseal_Out *out = (Unseal_Out *) 
            MemoryGetOutBuffer(sizeof(Unseal_Out));
    in->itemHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Unseal (in, out);
    rSize = sizeof(Unseal_Out);
    *respParmSize += TPM2B_SENSITIVE_DATA_Marshal(&out->outData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Unseal
#if CC_ObjectChangeAuth
case TPM_CC_ObjectChangeAuth: {
    ObjectChangeAuth_In *in = (ObjectChangeAuth_In *)
            MemoryGetInBuffer(sizeof(ObjectChangeAuth_In));
    ObjectChangeAuth_Out *out = (ObjectChangeAuth_Out *) 
            MemoryGetOutBuffer(sizeof(ObjectChangeAuth_Out));
    in->objectHandle = handles[0];
    in->parentHandle = handles[1];
    result = TPM2B_AUTH_Unmarshal(&in->newAuth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ObjectChangeAuth_newAuth);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ObjectChangeAuth (in, out);
    rSize = sizeof(ObjectChangeAuth_Out);
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->outPrivate, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ObjectChangeAuth
#if CC_CreateLoaded
case TPM_CC_CreateLoaded: {
    CreateLoaded_In *in = (CreateLoaded_In *)
            MemoryGetInBuffer(sizeof(CreateLoaded_In));
    CreateLoaded_Out *out = (CreateLoaded_Out *) 
            MemoryGetOutBuffer(sizeof(CreateLoaded_Out));
    in->parentHandle = handles[0];
    result = TPM2B_SENSITIVE_CREATE_Unmarshal(&in->inSensitive, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CreateLoaded_inSensitive);
    result = TPM2B_TEMPLATE_Unmarshal(&in->inPublic, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CreateLoaded_inPublic);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_CreateLoaded (in, out);
    rSize = sizeof(CreateLoaded_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->objectHandle;
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->outPrivate, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_PUBLIC_Marshal(&out->outPublic, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_NAME_Marshal(&out->name, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_CreateLoaded
#if CC_Duplicate
case TPM_CC_Duplicate: {
    Duplicate_In *in = (Duplicate_In *)
            MemoryGetInBuffer(sizeof(Duplicate_In));
    Duplicate_Out *out = (Duplicate_Out *) 
            MemoryGetOutBuffer(sizeof(Duplicate_Out));
    in->objectHandle = handles[0];
    in->newParentHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->encryptionKeyIn, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Duplicate_encryptionKeyIn);
    result = TPMT_SYM_DEF_OBJECT_Unmarshal(&in->symmetricAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Duplicate_symmetricAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Duplicate (in, out);
    rSize = sizeof(Duplicate_Out);
    *respParmSize += TPM2B_DATA_Marshal(&out->encryptionKeyOut, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->duplicate, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ENCRYPTED_SECRET_Marshal(&out->outSymSeed, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Duplicate
#if CC_Rewrap
case TPM_CC_Rewrap: {
    Rewrap_In *in = (Rewrap_In *)
            MemoryGetInBuffer(sizeof(Rewrap_In));
    Rewrap_Out *out = (Rewrap_Out *) 
            MemoryGetOutBuffer(sizeof(Rewrap_Out));
    in->oldParent = handles[0];
    in->newParent = handles[1];
    result = TPM2B_PRIVATE_Unmarshal(&in->inDuplicate, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Rewrap_inDuplicate);
    result = TPM2B_NAME_Unmarshal(&in->name, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Rewrap_name);
    result = TPM2B_ENCRYPTED_SECRET_Unmarshal(&in->inSymSeed, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Rewrap_inSymSeed);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Rewrap (in, out);
    rSize = sizeof(Rewrap_Out);
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->outDuplicate, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ENCRYPTED_SECRET_Marshal(&out->outSymSeed, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Rewrap
#if CC_Import
case TPM_CC_Import: {
    Import_In *in = (Import_In *)
            MemoryGetInBuffer(sizeof(Import_In));
    Import_Out *out = (Import_Out *) 
            MemoryGetOutBuffer(sizeof(Import_Out));
    in->parentHandle = handles[0];
    result = TPM2B_DATA_Unmarshal(&in->encryptionKey, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Import_encryptionKey);
    result = TPM2B_PUBLIC_Unmarshal(&in->objectPublic, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_Import_objectPublic);
    result = TPM2B_PRIVATE_Unmarshal(&in->duplicate, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Import_duplicate);
    result = TPM2B_ENCRYPTED_SECRET_Unmarshal(&in->inSymSeed, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Import_inSymSeed);
    result = TPMT_SYM_DEF_OBJECT_Unmarshal(&in->symmetricAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Import_symmetricAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Import (in, out);
    rSize = sizeof(Import_Out);
    *respParmSize += TPM2B_PRIVATE_Marshal(&out->outPrivate, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Import
#if CC_RSA_Encrypt
case TPM_CC_RSA_Encrypt: {
    RSA_Encrypt_In *in = (RSA_Encrypt_In *)
            MemoryGetInBuffer(sizeof(RSA_Encrypt_In));
    RSA_Encrypt_Out *out = (RSA_Encrypt_Out *) 
            MemoryGetOutBuffer(sizeof(RSA_Encrypt_Out));
    in->keyHandle = handles[0];
    result = TPM2B_PUBLIC_KEY_RSA_Unmarshal(&in->message, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_RSA_Encrypt_message);
    result = TPMT_RSA_DECRYPT_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_RSA_Encrypt_inScheme);
    result = TPM2B_DATA_Unmarshal(&in->label, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_RSA_Encrypt_label);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_RSA_Encrypt (in, out);
    rSize = sizeof(RSA_Encrypt_Out);
    *respParmSize += TPM2B_PUBLIC_KEY_RSA_Marshal(&out->outData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_RSA_Encrypt
#if CC_RSA_Decrypt
case TPM_CC_RSA_Decrypt: {
    RSA_Decrypt_In *in = (RSA_Decrypt_In *)
            MemoryGetInBuffer(sizeof(RSA_Decrypt_In));
    RSA_Decrypt_Out *out = (RSA_Decrypt_Out *) 
            MemoryGetOutBuffer(sizeof(RSA_Decrypt_Out));
    in->keyHandle = handles[0];
    result = TPM2B_PUBLIC_KEY_RSA_Unmarshal(&in->cipherText, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_RSA_Decrypt_cipherText);
    result = TPMT_RSA_DECRYPT_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_RSA_Decrypt_inScheme);
    result = TPM2B_DATA_Unmarshal(&in->label, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_RSA_Decrypt_label);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_RSA_Decrypt (in, out);
    rSize = sizeof(RSA_Decrypt_Out);
    *respParmSize += TPM2B_PUBLIC_KEY_RSA_Marshal(&out->message, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_RSA_Decrypt
#if CC_ECDH_KeyGen
case TPM_CC_ECDH_KeyGen: {
    ECDH_KeyGen_In *in = (ECDH_KeyGen_In *)
            MemoryGetInBuffer(sizeof(ECDH_KeyGen_In));
    ECDH_KeyGen_Out *out = (ECDH_KeyGen_Out *) 
            MemoryGetOutBuffer(sizeof(ECDH_KeyGen_Out));
    in->keyHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ECDH_KeyGen (in, out);
    rSize = sizeof(ECDH_KeyGen_Out);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->zPoint, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->pubPoint, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ECDH_KeyGen
#if CC_ECDH_ZGen
case TPM_CC_ECDH_ZGen: {
    ECDH_ZGen_In *in = (ECDH_ZGen_In *)
            MemoryGetInBuffer(sizeof(ECDH_ZGen_In));
    ECDH_ZGen_Out *out = (ECDH_ZGen_Out *) 
            MemoryGetOutBuffer(sizeof(ECDH_ZGen_Out));
    in->keyHandle = handles[0];
    result = TPM2B_ECC_POINT_Unmarshal(&in->inPoint, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ECDH_ZGen_inPoint);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ECDH_ZGen (in, out);
    rSize = sizeof(ECDH_ZGen_Out);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->outPoint, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ECDH_ZGen
#if CC_ECC_Parameters
case TPM_CC_ECC_Parameters: {
    ECC_Parameters_In *in = (ECC_Parameters_In *)
            MemoryGetInBuffer(sizeof(ECC_Parameters_In));
    ECC_Parameters_Out *out = (ECC_Parameters_Out *) 
            MemoryGetOutBuffer(sizeof(ECC_Parameters_Out));
    result = TPMI_ECC_CURVE_Unmarshal(&in->curveID, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ECC_Parameters_curveID);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ECC_Parameters (in, out);
    rSize = sizeof(ECC_Parameters_Out);
    *respParmSize += TPMS_ALGORITHM_DETAIL_ECC_Marshal(&out->parameters, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ECC_Parameters
#if CC_ZGen_2Phase
case TPM_CC_ZGen_2Phase: {
    ZGen_2Phase_In *in = (ZGen_2Phase_In *)
            MemoryGetInBuffer(sizeof(ZGen_2Phase_In));
    ZGen_2Phase_Out *out = (ZGen_2Phase_Out *) 
            MemoryGetOutBuffer(sizeof(ZGen_2Phase_Out));
    in->keyA = handles[0];
    result = TPM2B_ECC_POINT_Unmarshal(&in->inQsB, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ZGen_2Phase_inQsB);
    result = TPM2B_ECC_POINT_Unmarshal(&in->inQeB, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ZGen_2Phase_inQeB);
    result = TPMI_ECC_KEY_EXCHANGE_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_ZGen_2Phase_inScheme);
    result = UINT16_Unmarshal(&in->counter, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ZGen_2Phase_counter);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ZGen_2Phase (in, out);
    rSize = sizeof(ZGen_2Phase_Out);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->outZ1, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->outZ2, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ZGen_2Phase
#if CC_EncryptDecrypt
case TPM_CC_EncryptDecrypt: {
    EncryptDecrypt_In *in = (EncryptDecrypt_In *)
            MemoryGetInBuffer(sizeof(EncryptDecrypt_In));
    EncryptDecrypt_Out *out = (EncryptDecrypt_Out *) 
            MemoryGetOutBuffer(sizeof(EncryptDecrypt_Out));
    in->keyHandle = handles[0];
    result = TPMI_YES_NO_Unmarshal(&in->decrypt, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt_decrypt);
    result = TPMI_ALG_CIPHER_MODE_Unmarshal(&in->mode, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt_mode);
    result = TPM2B_IV_Unmarshal(&in->ivIn, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt_ivIn);
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->inData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt_inData);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_EncryptDecrypt (in, out);
    rSize = sizeof(EncryptDecrypt_Out);
    *respParmSize += TPM2B_MAX_BUFFER_Marshal(&out->outData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_IV_Marshal(&out->ivOut, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_EncryptDecrypt
#if CC_EncryptDecrypt2
case TPM_CC_EncryptDecrypt2: {
    EncryptDecrypt2_In *in = (EncryptDecrypt2_In *)
            MemoryGetInBuffer(sizeof(EncryptDecrypt2_In));
    EncryptDecrypt2_Out *out = (EncryptDecrypt2_Out *) 
            MemoryGetOutBuffer(sizeof(EncryptDecrypt2_Out));
    in->keyHandle = handles[0];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->inData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt2_inData);
    result = TPMI_YES_NO_Unmarshal(&in->decrypt, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt2_decrypt);
    result = TPMI_ALG_CIPHER_MODE_Unmarshal(&in->mode, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt2_mode);
    result = TPM2B_IV_Unmarshal(&in->ivIn, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EncryptDecrypt2_ivIn);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_EncryptDecrypt2 (in, out);
    rSize = sizeof(EncryptDecrypt2_Out);
    *respParmSize += TPM2B_MAX_BUFFER_Marshal(&out->outData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_IV_Marshal(&out->ivOut, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_EncryptDecrypt2
#if CC_Hash
case TPM_CC_Hash: {
    Hash_In *in = (Hash_In *)
            MemoryGetInBuffer(sizeof(Hash_In));
    Hash_Out *out = (Hash_Out *) 
            MemoryGetOutBuffer(sizeof(Hash_Out));
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->data, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Hash_data);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_Hash_hashAlg);
    result = TPMI_RH_HIERARCHY_Unmarshal(&in->hierarchy, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Hash_hierarchy);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Hash (in, out);
    rSize = sizeof(Hash_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->outHash, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_HASHCHECK_Marshal(&out->validation, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Hash
#if CC_HMAC
case TPM_CC_HMAC: {
    HMAC_In *in = (HMAC_In *)
            MemoryGetInBuffer(sizeof(HMAC_In));
    HMAC_Out *out = (HMAC_Out *) 
            MemoryGetOutBuffer(sizeof(HMAC_Out));
    in->handle = handles[0];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->buffer, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_HMAC_buffer);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_HMAC_hashAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_HMAC (in, out);
    rSize = sizeof(HMAC_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->outHMAC, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_HMAC
#if CC_MAC
case TPM_CC_MAC: {
    MAC_In *in = (MAC_In *)
            MemoryGetInBuffer(sizeof(MAC_In));
    MAC_Out *out = (MAC_Out *) 
            MemoryGetOutBuffer(sizeof(MAC_Out));
    in->handle = handles[0];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->buffer, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_MAC_buffer);
    result = TPMI_ALG_MAC_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_MAC_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_MAC (in, out);
    rSize = sizeof(MAC_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->outMAC, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_MAC
#if CC_GetRandom
case TPM_CC_GetRandom: {
    GetRandom_In *in = (GetRandom_In *)
            MemoryGetInBuffer(sizeof(GetRandom_In));
    GetRandom_Out *out = (GetRandom_Out *) 
            MemoryGetOutBuffer(sizeof(GetRandom_Out));
    result = UINT16_Unmarshal(&in->bytesRequested, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetRandom_bytesRequested);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetRandom (in, out);
    rSize = sizeof(GetRandom_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->randomBytes, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetRandom
#if CC_StirRandom
case TPM_CC_StirRandom: {
    StirRandom_In *in = (StirRandom_In *)
            MemoryGetInBuffer(sizeof(StirRandom_In));
    result = TPM2B_SENSITIVE_DATA_Unmarshal(&in->inData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_StirRandom_inData);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_StirRandom (in);
break; 
}
#endif     // CC_StirRandom
#if CC_HMAC_Start
case TPM_CC_HMAC_Start: {
    HMAC_Start_In *in = (HMAC_Start_In *)
            MemoryGetInBuffer(sizeof(HMAC_Start_In));
    HMAC_Start_Out *out = (HMAC_Start_Out *) 
            MemoryGetOutBuffer(sizeof(HMAC_Start_Out));
    in->handle = handles[0];
    result = TPM2B_AUTH_Unmarshal(&in->auth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_HMAC_Start_auth);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_HMAC_Start_hashAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_HMAC_Start (in, out);
    rSize = sizeof(HMAC_Start_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->sequenceHandle;
break; 
}
#endif     // CC_HMAC_Start
#if CC_MAC_Start
case TPM_CC_MAC_Start: {
    MAC_Start_In *in = (MAC_Start_In *)
            MemoryGetInBuffer(sizeof(MAC_Start_In));
    MAC_Start_Out *out = (MAC_Start_Out *) 
            MemoryGetOutBuffer(sizeof(MAC_Start_Out));
    in->handle = handles[0];
    result = TPM2B_AUTH_Unmarshal(&in->auth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_MAC_Start_auth);
    result = TPMI_ALG_MAC_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_MAC_Start_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_MAC_Start (in, out);
    rSize = sizeof(MAC_Start_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->sequenceHandle;
break; 
}
#endif     // CC_MAC_Start
#if CC_HashSequenceStart
case TPM_CC_HashSequenceStart: {
    HashSequenceStart_In *in = (HashSequenceStart_In *)
            MemoryGetInBuffer(sizeof(HashSequenceStart_In));
    HashSequenceStart_Out *out = (HashSequenceStart_Out *) 
            MemoryGetOutBuffer(sizeof(HashSequenceStart_Out));
    result = TPM2B_AUTH_Unmarshal(&in->auth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_HashSequenceStart_auth);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_HashSequenceStart_hashAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_HashSequenceStart (in, out);
    rSize = sizeof(HashSequenceStart_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->sequenceHandle;
break; 
}
#endif     // CC_HashSequenceStart
#if CC_SequenceUpdate
case TPM_CC_SequenceUpdate: {
    SequenceUpdate_In *in = (SequenceUpdate_In *)
            MemoryGetInBuffer(sizeof(SequenceUpdate_In));
    in->sequenceHandle = handles[0];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->buffer, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SequenceUpdate_buffer);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SequenceUpdate (in);
break; 
}
#endif     // CC_SequenceUpdate
#if CC_SequenceComplete
case TPM_CC_SequenceComplete: {
    SequenceComplete_In *in = (SequenceComplete_In *)
            MemoryGetInBuffer(sizeof(SequenceComplete_In));
    SequenceComplete_Out *out = (SequenceComplete_Out *) 
            MemoryGetOutBuffer(sizeof(SequenceComplete_Out));
    in->sequenceHandle = handles[0];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->buffer, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SequenceComplete_buffer);
    result = TPMI_RH_HIERARCHY_Unmarshal(&in->hierarchy, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_SequenceComplete_hierarchy);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SequenceComplete (in, out);
    rSize = sizeof(SequenceComplete_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->result, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_HASHCHECK_Marshal(&out->validation, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_SequenceComplete
#if CC_EventSequenceComplete
case TPM_CC_EventSequenceComplete: {
    EventSequenceComplete_In *in = (EventSequenceComplete_In *)
            MemoryGetInBuffer(sizeof(EventSequenceComplete_In));
    EventSequenceComplete_Out *out = (EventSequenceComplete_Out *) 
            MemoryGetOutBuffer(sizeof(EventSequenceComplete_Out));
    in->pcrHandle = handles[0];
    in->sequenceHandle = handles[1];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->buffer, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EventSequenceComplete_buffer);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_EventSequenceComplete (in, out);
    rSize = sizeof(EventSequenceComplete_Out);
    *respParmSize += TPML_DIGEST_VALUES_Marshal(&out->results, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_EventSequenceComplete
#if CC_Certify
case TPM_CC_Certify: {
    Certify_In *in = (Certify_In *)
            MemoryGetInBuffer(sizeof(Certify_In));
    Certify_Out *out = (Certify_Out *) 
            MemoryGetOutBuffer(sizeof(Certify_Out));
    in->objectHandle = handles[0];
    in->signHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Certify_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Certify_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Certify (in, out);
    rSize = sizeof(Certify_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->certifyInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Certify
#if CC_CertifyCreation
case TPM_CC_CertifyCreation: {
    CertifyCreation_In *in = (CertifyCreation_In *)
            MemoryGetInBuffer(sizeof(CertifyCreation_In));
    CertifyCreation_Out *out = (CertifyCreation_Out *) 
            MemoryGetOutBuffer(sizeof(CertifyCreation_Out));
    in->signHandle = handles[0];
    in->objectHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CertifyCreation_qualifyingData);
    result = TPM2B_DIGEST_Unmarshal(&in->creationHash, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CertifyCreation_creationHash);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_CertifyCreation_inScheme);
    result = TPMT_TK_CREATION_Unmarshal(&in->creationTicket, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CertifyCreation_creationTicket);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_CertifyCreation (in, out);
    rSize = sizeof(CertifyCreation_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->certifyInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_CertifyCreation
#if CC_Quote
case TPM_CC_Quote: {
    Quote_In *in = (Quote_In *)
            MemoryGetInBuffer(sizeof(Quote_In));
    Quote_Out *out = (Quote_Out *) 
            MemoryGetOutBuffer(sizeof(Quote_Out));
    in->signHandle = handles[0];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Quote_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Quote_inScheme);
    result = TPML_PCR_SELECTION_Unmarshal(&in->PCRselect, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Quote_PCRselect);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Quote (in, out);
    rSize = sizeof(Quote_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->quoted, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Quote
#if CC_GetSessionAuditDigest
case TPM_CC_GetSessionAuditDigest: {
    GetSessionAuditDigest_In *in = (GetSessionAuditDigest_In *)
            MemoryGetInBuffer(sizeof(GetSessionAuditDigest_In));
    GetSessionAuditDigest_Out *out = (GetSessionAuditDigest_Out *) 
            MemoryGetOutBuffer(sizeof(GetSessionAuditDigest_Out));
    in->privacyAdminHandle = handles[0];
    in->signHandle = handles[1];
    in->sessionHandle = handles[2];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetSessionAuditDigest_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_GetSessionAuditDigest_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetSessionAuditDigest (in, out);
    rSize = sizeof(GetSessionAuditDigest_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->auditInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetSessionAuditDigest
#if CC_GetCommandAuditDigest
case TPM_CC_GetCommandAuditDigest: {
    GetCommandAuditDigest_In *in = (GetCommandAuditDigest_In *)
            MemoryGetInBuffer(sizeof(GetCommandAuditDigest_In));
    GetCommandAuditDigest_Out *out = (GetCommandAuditDigest_Out *) 
            MemoryGetOutBuffer(sizeof(GetCommandAuditDigest_Out));
    in->privacyHandle = handles[0];
    in->signHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetCommandAuditDigest_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_GetCommandAuditDigest_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetCommandAuditDigest (in, out);
    rSize = sizeof(GetCommandAuditDigest_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->auditInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetCommandAuditDigest
#if CC_GetTime
case TPM_CC_GetTime: {
    GetTime_In *in = (GetTime_In *)
            MemoryGetInBuffer(sizeof(GetTime_In));
    GetTime_Out *out = (GetTime_Out *) 
            MemoryGetOutBuffer(sizeof(GetTime_Out));
    in->privacyAdminHandle = handles[0];
    in->signHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetTime_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_GetTime_inScheme);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetTime (in, out);
    rSize = sizeof(GetTime_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->timeInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetTime
#if CC_CertifyX509
case TPM_CC_CertifyX509: {
    CertifyX509_In *in = (CertifyX509_In *)
            MemoryGetInBuffer(sizeof(CertifyX509_In));
    CertifyX509_Out *out = (CertifyX509_Out *) 
            MemoryGetOutBuffer(sizeof(CertifyX509_Out));
    in->objectHandle = handles[0];
    in->signHandle = handles[1];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CertifyX509_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_CertifyX509_inScheme);
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->partialCertificate, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CertifyX509_partialCertificate);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_CertifyX509 (in, out);
    rSize = sizeof(CertifyX509_Out);
    *respParmSize += TPM2B_MAX_BUFFER_Marshal(&out->addedToCertificate, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->tbsDigest, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_CertifyX509
#if CC_Commit
case TPM_CC_Commit: {
    Commit_In *in = (Commit_In *)
            MemoryGetInBuffer(sizeof(Commit_In));
    Commit_Out *out = (Commit_Out *) 
            MemoryGetOutBuffer(sizeof(Commit_Out));
    in->signHandle = handles[0];
    result = TPM2B_ECC_POINT_Unmarshal(&in->P1, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Commit_P1);
    result = TPM2B_SENSITIVE_DATA_Unmarshal(&in->s2, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Commit_s2);
    result = TPM2B_ECC_PARAMETER_Unmarshal(&in->y2, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Commit_y2);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Commit (in, out);
    rSize = sizeof(Commit_Out);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->K, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->L, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->E, 
                                          responseBuffer, &rSize);
    *respParmSize += UINT16_Marshal(&out->counter, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Commit
#if CC_EC_Ephemeral
case TPM_CC_EC_Ephemeral: {
    EC_Ephemeral_In *in = (EC_Ephemeral_In *)
            MemoryGetInBuffer(sizeof(EC_Ephemeral_In));
    EC_Ephemeral_Out *out = (EC_Ephemeral_Out *) 
            MemoryGetOutBuffer(sizeof(EC_Ephemeral_Out));
    result = TPMI_ECC_CURVE_Unmarshal(&in->curveID, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EC_Ephemeral_curveID);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_EC_Ephemeral (in, out);
    rSize = sizeof(EC_Ephemeral_Out);
    *respParmSize += TPM2B_ECC_POINT_Marshal(&out->Q, 
                                          responseBuffer, &rSize);
    *respParmSize += UINT16_Marshal(&out->counter, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_EC_Ephemeral
#if CC_VerifySignature
case TPM_CC_VerifySignature: {
    VerifySignature_In *in = (VerifySignature_In *)
            MemoryGetInBuffer(sizeof(VerifySignature_In));
    VerifySignature_Out *out = (VerifySignature_Out *) 
            MemoryGetOutBuffer(sizeof(VerifySignature_Out));
    in->keyHandle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->digest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_VerifySignature_digest);
    result = TPMT_SIGNATURE_Unmarshal(&in->signature, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_VerifySignature_signature);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_VerifySignature (in, out);
    rSize = sizeof(VerifySignature_Out);
    *respParmSize += TPMT_TK_VERIFIED_Marshal(&out->validation, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_VerifySignature
#if CC_Sign
case TPM_CC_Sign: {
    Sign_In *in = (Sign_In *)
            MemoryGetInBuffer(sizeof(Sign_In));
    Sign_Out *out = (Sign_Out *) 
            MemoryGetOutBuffer(sizeof(Sign_Out));
    in->keyHandle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->digest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Sign_digest);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_Sign_inScheme);
    result = TPMT_TK_HASHCHECK_Unmarshal(&in->validation, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Sign_validation);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Sign (in, out);
    rSize = sizeof(Sign_Out);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Sign
#if CC_SetCommandCodeAuditStatus
case TPM_CC_SetCommandCodeAuditStatus: {
    SetCommandCodeAuditStatus_In *in = (SetCommandCodeAuditStatus_In *)
            MemoryGetInBuffer(sizeof(SetCommandCodeAuditStatus_In));
    in->auth = handles[0];
    result = TPMI_ALG_HASH_Unmarshal(&in->auditAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_SetCommandCodeAuditStatus_auditAlg);
    result = TPML_CC_Unmarshal(&in->setList, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SetCommandCodeAuditStatus_setList);
    result = TPML_CC_Unmarshal(&in->clearList, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SetCommandCodeAuditStatus_clearList);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SetCommandCodeAuditStatus (in);
break; 
}
#endif     // CC_SetCommandCodeAuditStatus
#if CC_PCR_Extend
case TPM_CC_PCR_Extend: {
    PCR_Extend_In *in = (PCR_Extend_In *)
            MemoryGetInBuffer(sizeof(PCR_Extend_In));
    in->pcrHandle = handles[0];
    result = TPML_DIGEST_VALUES_Unmarshal(&in->digests, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_Extend_digests);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_Extend (in);
break; 
}
#endif     // CC_PCR_Extend
#if CC_PCR_Event
case TPM_CC_PCR_Event: {
    PCR_Event_In *in = (PCR_Event_In *)
            MemoryGetInBuffer(sizeof(PCR_Event_In));
    PCR_Event_Out *out = (PCR_Event_Out *) 
            MemoryGetOutBuffer(sizeof(PCR_Event_Out));
    in->pcrHandle = handles[0];
    result = TPM2B_EVENT_Unmarshal(&in->eventData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_Event_eventData);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_Event (in, out);
    rSize = sizeof(PCR_Event_Out);
    *respParmSize += TPML_DIGEST_VALUES_Marshal(&out->digests, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PCR_Event
#if CC_PCR_Read
case TPM_CC_PCR_Read: {
    PCR_Read_In *in = (PCR_Read_In *)
            MemoryGetInBuffer(sizeof(PCR_Read_In));
    PCR_Read_Out *out = (PCR_Read_Out *) 
            MemoryGetOutBuffer(sizeof(PCR_Read_Out));
    result = TPML_PCR_SELECTION_Unmarshal(&in->pcrSelectionIn, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_Read_pcrSelectionIn);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_Read (in, out);
    rSize = sizeof(PCR_Read_Out);
    *respParmSize += UINT32_Marshal(&out->pcrUpdateCounter, 
                                          responseBuffer, &rSize);
    *respParmSize += TPML_PCR_SELECTION_Marshal(&out->pcrSelectionOut, 
                                          responseBuffer, &rSize);
    *respParmSize += TPML_DIGEST_Marshal(&out->pcrValues, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PCR_Read
#if CC_PCR_Allocate
case TPM_CC_PCR_Allocate: {
    PCR_Allocate_In *in = (PCR_Allocate_In *)
            MemoryGetInBuffer(sizeof(PCR_Allocate_In));
    PCR_Allocate_Out *out = (PCR_Allocate_Out *) 
            MemoryGetOutBuffer(sizeof(PCR_Allocate_Out));
    in->authHandle = handles[0];
    result = TPML_PCR_SELECTION_Unmarshal(&in->pcrAllocation, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_Allocate_pcrAllocation);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_Allocate (in, out);
    rSize = sizeof(PCR_Allocate_Out);
    *respParmSize += TPMI_YES_NO_Marshal(&out->allocationSuccess, 
                                          responseBuffer, &rSize);
    *respParmSize += UINT32_Marshal(&out->maxPCR, 
                                          responseBuffer, &rSize);
    *respParmSize += UINT32_Marshal(&out->sizeNeeded, 
                                          responseBuffer, &rSize);
    *respParmSize += UINT32_Marshal(&out->sizeAvailable, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PCR_Allocate
#if CC_PCR_SetAuthPolicy
case TPM_CC_PCR_SetAuthPolicy: {
    PCR_SetAuthPolicy_In *in = (PCR_SetAuthPolicy_In *)
            MemoryGetInBuffer(sizeof(PCR_SetAuthPolicy_In));
    in->authHandle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->authPolicy, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_SetAuthPolicy_authPolicy);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_PCR_SetAuthPolicy_hashAlg);
    result = TPMI_DH_PCR_Unmarshal(&in->pcrNum, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_PCR_SetAuthPolicy_pcrNum);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_SetAuthPolicy (in);
break; 
}
#endif     // CC_PCR_SetAuthPolicy
#if CC_PCR_SetAuthValue
case TPM_CC_PCR_SetAuthValue: {
    PCR_SetAuthValue_In *in = (PCR_SetAuthValue_In *)
            MemoryGetInBuffer(sizeof(PCR_SetAuthValue_In));
    in->pcrHandle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->auth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PCR_SetAuthValue_auth);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_SetAuthValue (in);
break; 
}
#endif     // CC_PCR_SetAuthValue
#if CC_PCR_Reset
case TPM_CC_PCR_Reset: {
    PCR_Reset_In *in = (PCR_Reset_In *)
            MemoryGetInBuffer(sizeof(PCR_Reset_In));
    in->pcrHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PCR_Reset (in);
break; 
}
#endif     // CC_PCR_Reset
#if CC_PolicySigned
case TPM_CC_PolicySigned: {
    PolicySigned_In *in = (PolicySigned_In *)
            MemoryGetInBuffer(sizeof(PolicySigned_In));
    PolicySigned_Out *out = (PolicySigned_Out *) 
            MemoryGetOutBuffer(sizeof(PolicySigned_Out));
    in->authObject = handles[0];
    in->policySession = handles[1];
    result = TPM2B_NONCE_Unmarshal(&in->nonceTPM, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySigned_nonceTPM);
    result = TPM2B_DIGEST_Unmarshal(&in->cpHashA, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySigned_cpHashA);
    result = TPM2B_NONCE_Unmarshal(&in->policyRef, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySigned_policyRef);
    result = INT32_Unmarshal(&in->expiration, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySigned_expiration);
    result = TPMT_SIGNATURE_Unmarshal(&in->auth, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_PolicySigned_auth);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicySigned (in, out);
    rSize = sizeof(PolicySigned_Out);
    *respParmSize += TPM2B_TIMEOUT_Marshal(&out->timeout, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_AUTH_Marshal(&out->policyTicket, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PolicySigned
#if CC_PolicySecret
case TPM_CC_PolicySecret: {
    PolicySecret_In *in = (PolicySecret_In *)
            MemoryGetInBuffer(sizeof(PolicySecret_In));
    PolicySecret_Out *out = (PolicySecret_Out *) 
            MemoryGetOutBuffer(sizeof(PolicySecret_Out));
    in->authHandle = handles[0];
    in->policySession = handles[1];
    result = TPM2B_NONCE_Unmarshal(&in->nonceTPM, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySecret_nonceTPM);
    result = TPM2B_DIGEST_Unmarshal(&in->cpHashA, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySecret_cpHashA);
    result = TPM2B_NONCE_Unmarshal(&in->policyRef, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySecret_policyRef);
    result = INT32_Unmarshal(&in->expiration, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicySecret_expiration);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicySecret (in, out);
    rSize = sizeof(PolicySecret_Out);
    *respParmSize += TPM2B_TIMEOUT_Marshal(&out->timeout, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_AUTH_Marshal(&out->policyTicket, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PolicySecret
#if CC_PolicyTicket
case TPM_CC_PolicyTicket: {
    PolicyTicket_In *in = (PolicyTicket_In *)
            MemoryGetInBuffer(sizeof(PolicyTicket_In));
    in->policySession = handles[0];
    result = TPM2B_TIMEOUT_Unmarshal(&in->timeout, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTicket_timeout);
    result = TPM2B_DIGEST_Unmarshal(&in->cpHashA, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTicket_cpHashA);
    result = TPM2B_NONCE_Unmarshal(&in->policyRef, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTicket_policyRef);
    result = TPM2B_NAME_Unmarshal(&in->authName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTicket_authName);
    result = TPMT_TK_AUTH_Unmarshal(&in->ticket, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTicket_ticket);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyTicket (in);
break; 
}
#endif     // CC_PolicyTicket
#if CC_PolicyOR
case TPM_CC_PolicyOR: {
    PolicyOR_In *in = (PolicyOR_In *)
            MemoryGetInBuffer(sizeof(PolicyOR_In));
    in->policySession = handles[0];
    result = TPML_DIGEST_Unmarshal(&in->pHashList, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyOR_pHashList);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyOR (in);
break; 
}
#endif     // CC_PolicyOR
#if CC_PolicyPCR
case TPM_CC_PolicyPCR: {
    PolicyPCR_In *in = (PolicyPCR_In *)
            MemoryGetInBuffer(sizeof(PolicyPCR_In));
    in->policySession = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->pcrDigest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyPCR_pcrDigest);
    result = TPML_PCR_SELECTION_Unmarshal(&in->pcrs, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyPCR_pcrs);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyPCR (in);
break; 
}
#endif     // CC_PolicyPCR
#if CC_PolicyLocality
case TPM_CC_PolicyLocality: {
    PolicyLocality_In *in = (PolicyLocality_In *)
            MemoryGetInBuffer(sizeof(PolicyLocality_In));
    in->policySession = handles[0];
    result = TPMA_LOCALITY_Unmarshal(&in->locality, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyLocality_locality);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyLocality (in);
break; 
}
#endif     // CC_PolicyLocality
#if CC_PolicyNV
case TPM_CC_PolicyNV: {
    PolicyNV_In *in = (PolicyNV_In *)
            MemoryGetInBuffer(sizeof(PolicyNV_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    in->policySession = handles[2];
    result = TPM2B_OPERAND_Unmarshal(&in->operandB, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyNV_operandB);
    result = UINT16_Unmarshal(&in->offset, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyNV_offset);
    result = TPM_EO_Unmarshal(&in->operation, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyNV_operation);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyNV (in);
break; 
}
#endif     // CC_PolicyNV
#if CC_PolicyCounterTimer
case TPM_CC_PolicyCounterTimer: {
    PolicyCounterTimer_In *in = (PolicyCounterTimer_In *)
            MemoryGetInBuffer(sizeof(PolicyCounterTimer_In));
    in->policySession = handles[0];
    result = TPM2B_OPERAND_Unmarshal(&in->operandB, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyCounterTimer_operandB);
    result = UINT16_Unmarshal(&in->offset, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyCounterTimer_offset);
    result = TPM_EO_Unmarshal(&in->operation, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyCounterTimer_operation);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyCounterTimer (in);
break; 
}
#endif     // CC_PolicyCounterTimer
#if CC_PolicyCommandCode
case TPM_CC_PolicyCommandCode: {
    PolicyCommandCode_In *in = (PolicyCommandCode_In *)
            MemoryGetInBuffer(sizeof(PolicyCommandCode_In));
    in->policySession = handles[0];
    result = TPM_CC_Unmarshal(&in->code, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyCommandCode_code);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyCommandCode (in);
break; 
}
#endif     // CC_PolicyCommandCode
#if CC_PolicyPhysicalPresence
case TPM_CC_PolicyPhysicalPresence: {
    PolicyPhysicalPresence_In *in = (PolicyPhysicalPresence_In *)
            MemoryGetInBuffer(sizeof(PolicyPhysicalPresence_In));
    in->policySession = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyPhysicalPresence (in);
break; 
}
#endif     // CC_PolicyPhysicalPresence
#if CC_PolicyCpHash
case TPM_CC_PolicyCpHash: {
    PolicyCpHash_In *in = (PolicyCpHash_In *)
            MemoryGetInBuffer(sizeof(PolicyCpHash_In));
    in->policySession = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->cpHashA, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyCpHash_cpHashA);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyCpHash (in);
break; 
}
#endif     // CC_PolicyCpHash
#if CC_PolicyNameHash
case TPM_CC_PolicyNameHash: {
    PolicyNameHash_In *in = (PolicyNameHash_In *)
            MemoryGetInBuffer(sizeof(PolicyNameHash_In));
    in->policySession = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->nameHash, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyNameHash_nameHash);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyNameHash (in);
break; 
}
#endif     // CC_PolicyNameHash
#if CC_PolicyDuplicationSelect
case TPM_CC_PolicyDuplicationSelect: {
    PolicyDuplicationSelect_In *in = (PolicyDuplicationSelect_In *)
            MemoryGetInBuffer(sizeof(PolicyDuplicationSelect_In));
    in->policySession = handles[0];
    result = TPM2B_NAME_Unmarshal(&in->objectName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyDuplicationSelect_objectName);
    result = TPM2B_NAME_Unmarshal(&in->newParentName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyDuplicationSelect_newParentName);
    result = TPMI_YES_NO_Unmarshal(&in->includeObject, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyDuplicationSelect_includeObject);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyDuplicationSelect (in);
break; 
}
#endif     // CC_PolicyDuplicationSelect
#if CC_PolicyAuthorize
case TPM_CC_PolicyAuthorize: {
    PolicyAuthorize_In *in = (PolicyAuthorize_In *)
            MemoryGetInBuffer(sizeof(PolicyAuthorize_In));
    in->policySession = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->approvedPolicy, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyAuthorize_approvedPolicy);
    result = TPM2B_NONCE_Unmarshal(&in->policyRef, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyAuthorize_policyRef);
    result = TPM2B_NAME_Unmarshal(&in->keySign, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyAuthorize_keySign);
    result = TPMT_TK_VERIFIED_Unmarshal(&in->checkTicket, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyAuthorize_checkTicket);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyAuthorize (in);
break; 
}
#endif     // CC_PolicyAuthorize
#if CC_PolicyAuthValue
case TPM_CC_PolicyAuthValue: {
    PolicyAuthValue_In *in = (PolicyAuthValue_In *)
            MemoryGetInBuffer(sizeof(PolicyAuthValue_In));
    in->policySession = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyAuthValue (in);
break; 
}
#endif     // CC_PolicyAuthValue
#if CC_PolicyPassword
case TPM_CC_PolicyPassword: {
    PolicyPassword_In *in = (PolicyPassword_In *)
            MemoryGetInBuffer(sizeof(PolicyPassword_In));
    in->policySession = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyPassword (in);
break; 
}
#endif     // CC_PolicyPassword
#if CC_PolicyGetDigest
case TPM_CC_PolicyGetDigest: {
    PolicyGetDigest_In *in = (PolicyGetDigest_In *)
            MemoryGetInBuffer(sizeof(PolicyGetDigest_In));
    PolicyGetDigest_Out *out = (PolicyGetDigest_Out *) 
            MemoryGetOutBuffer(sizeof(PolicyGetDigest_Out));
    in->policySession = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyGetDigest (in, out);
    rSize = sizeof(PolicyGetDigest_Out);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->policyDigest, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_PolicyGetDigest
#if CC_PolicyNvWritten
case TPM_CC_PolicyNvWritten: {
    PolicyNvWritten_In *in = (PolicyNvWritten_In *)
            MemoryGetInBuffer(sizeof(PolicyNvWritten_In));
    in->policySession = handles[0];
    result = TPMI_YES_NO_Unmarshal(&in->writtenSet, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyNvWritten_writtenSet);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyNvWritten (in);
break; 
}
#endif     // CC_PolicyNvWritten
#if CC_PolicyTemplate
case TPM_CC_PolicyTemplate: {
    PolicyTemplate_In *in = (PolicyTemplate_In *)
            MemoryGetInBuffer(sizeof(PolicyTemplate_In));
    in->policySession = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->templateHash, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PolicyTemplate_templateHash);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyTemplate (in);
break; 
}
#endif     // CC_PolicyTemplate
#if CC_PolicyAuthorizeNV
case TPM_CC_PolicyAuthorizeNV: {
    PolicyAuthorizeNV_In *in = (PolicyAuthorizeNV_In *)
            MemoryGetInBuffer(sizeof(PolicyAuthorizeNV_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    in->policySession = handles[2];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PolicyAuthorizeNV (in);
break; 
}
#endif     // CC_PolicyAuthorizeNV
#if CC_CreatePrimary
case TPM_CC_CreatePrimary: {
    CreatePrimary_In *in = (CreatePrimary_In *)
            MemoryGetInBuffer(sizeof(CreatePrimary_In));
    CreatePrimary_Out *out = (CreatePrimary_Out *) 
            MemoryGetOutBuffer(sizeof(CreatePrimary_Out));
    in->primaryHandle = handles[0];
    result = TPM2B_SENSITIVE_CREATE_Unmarshal(&in->inSensitive, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CreatePrimary_inSensitive);
    result = TPM2B_PUBLIC_Unmarshal(&in->inPublic, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_CreatePrimary_inPublic);
    result = TPM2B_DATA_Unmarshal(&in->outsideInfo, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CreatePrimary_outsideInfo);
    result = TPML_PCR_SELECTION_Unmarshal(&in->creationPCR, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_CreatePrimary_creationPCR);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_CreatePrimary (in, out);
    rSize = sizeof(CreatePrimary_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->objectHandle;
    *respParmSize += TPM2B_PUBLIC_Marshal(&out->outPublic, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_CREATION_DATA_Marshal(&out->creationData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_DIGEST_Marshal(&out->creationHash, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_TK_CREATION_Marshal(&out->creationTicket, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_NAME_Marshal(&out->name, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_CreatePrimary
#if CC_HierarchyControl
case TPM_CC_HierarchyControl: {
    HierarchyControl_In *in = (HierarchyControl_In *)
            MemoryGetInBuffer(sizeof(HierarchyControl_In));
    in->authHandle = handles[0];
    result = TPMI_RH_ENABLES_Unmarshal(&in->enable, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_HierarchyControl_enable);
    result = TPMI_YES_NO_Unmarshal(&in->state, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_HierarchyControl_state);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_HierarchyControl (in);
break; 
}
#endif     // CC_HierarchyControl
#if CC_SetPrimaryPolicy
case TPM_CC_SetPrimaryPolicy: {
    SetPrimaryPolicy_In *in = (SetPrimaryPolicy_In *)
            MemoryGetInBuffer(sizeof(SetPrimaryPolicy_In));
    in->authHandle = handles[0];
    result = TPM2B_DIGEST_Unmarshal(&in->authPolicy, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SetPrimaryPolicy_authPolicy);
    result = TPMI_ALG_HASH_Unmarshal(&in->hashAlg, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_SetPrimaryPolicy_hashAlg);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SetPrimaryPolicy (in);
break; 
}
#endif     // CC_SetPrimaryPolicy
#if CC_ChangePPS
case TPM_CC_ChangePPS: {
    ChangePPS_In *in = (ChangePPS_In *)
            MemoryGetInBuffer(sizeof(ChangePPS_In));
    in->authHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ChangePPS (in);
break; 
}
#endif     // CC_ChangePPS
#if CC_ChangeEPS
case TPM_CC_ChangeEPS: {
    ChangeEPS_In *in = (ChangeEPS_In *)
            MemoryGetInBuffer(sizeof(ChangeEPS_In));
    in->authHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ChangeEPS (in);
break; 
}
#endif     // CC_ChangeEPS
#if CC_Clear
case TPM_CC_Clear: {
    Clear_In *in = (Clear_In *)
            MemoryGetInBuffer(sizeof(Clear_In));
    in->authHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Clear (in);
break; 
}
#endif     // CC_Clear
#if CC_ClearControl
case TPM_CC_ClearControl: {
    ClearControl_In *in = (ClearControl_In *)
            MemoryGetInBuffer(sizeof(ClearControl_In));
    in->auth = handles[0];
    result = TPMI_YES_NO_Unmarshal(&in->disable, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ClearControl_disable);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ClearControl (in);
break; 
}
#endif     // CC_ClearControl
#if CC_HierarchyChangeAuth
case TPM_CC_HierarchyChangeAuth: {
    HierarchyChangeAuth_In *in = (HierarchyChangeAuth_In *)
            MemoryGetInBuffer(sizeof(HierarchyChangeAuth_In));
    in->authHandle = handles[0];
    result = TPM2B_AUTH_Unmarshal(&in->newAuth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_HierarchyChangeAuth_newAuth);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_HierarchyChangeAuth (in);
break; 
}
#endif     // CC_HierarchyChangeAuth
#if CC_DictionaryAttackLockReset
case TPM_CC_DictionaryAttackLockReset: {
    DictionaryAttackLockReset_In *in = (DictionaryAttackLockReset_In *)
            MemoryGetInBuffer(sizeof(DictionaryAttackLockReset_In));
    in->lockHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_DictionaryAttackLockReset (in);
break; 
}
#endif     // CC_DictionaryAttackLockReset
#if CC_DictionaryAttackParameters
case TPM_CC_DictionaryAttackParameters: {
    DictionaryAttackParameters_In *in = (DictionaryAttackParameters_In *)
            MemoryGetInBuffer(sizeof(DictionaryAttackParameters_In));
    in->lockHandle = handles[0];
    result = UINT32_Unmarshal(&in->newMaxTries, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_DictionaryAttackParameters_newMaxTries);
    result = UINT32_Unmarshal(&in->newRecoveryTime, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_DictionaryAttackParameters_newRecoveryTime);
    result = UINT32_Unmarshal(&in->lockoutRecovery, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_DictionaryAttackParameters_lockoutRecovery);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_DictionaryAttackParameters (in);
break; 
}
#endif     // CC_DictionaryAttackParameters
#if CC_PP_Commands
case TPM_CC_PP_Commands: {
    PP_Commands_In *in = (PP_Commands_In *)
            MemoryGetInBuffer(sizeof(PP_Commands_In));
    in->auth = handles[0];
    result = TPML_CC_Unmarshal(&in->setList, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PP_Commands_setList);
    result = TPML_CC_Unmarshal(&in->clearList, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_PP_Commands_clearList);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_PP_Commands (in);
break; 
}
#endif     // CC_PP_Commands
#if CC_SetAlgorithmSet
case TPM_CC_SetAlgorithmSet: {
    SetAlgorithmSet_In *in = (SetAlgorithmSet_In *)
            MemoryGetInBuffer(sizeof(SetAlgorithmSet_In));
    in->authHandle = handles[0];
    result = UINT32_Unmarshal(&in->algorithmSet, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_SetAlgorithmSet_algorithmSet);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_SetAlgorithmSet (in);
break; 
}
#endif     // CC_SetAlgorithmSet
#if CC_FieldUpgradeStart
case TPM_CC_FieldUpgradeStart: {
    FieldUpgradeStart_In *in = (FieldUpgradeStart_In *)
            MemoryGetInBuffer(sizeof(FieldUpgradeStart_In));
    in->authorization = handles[0];
    in->keyHandle = handles[1];
    result = TPM2B_DIGEST_Unmarshal(&in->fuDigest, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_FieldUpgradeStart_fuDigest);
    result = TPMT_SIGNATURE_Unmarshal(&in->manifestSignature, paramBuffer, paramBufferSize, FALSE);
        ERROR_IF_EXIT_PLUS(RC_FieldUpgradeStart_manifestSignature);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_FieldUpgradeStart (in);
break; 
}
#endif     // CC_FieldUpgradeStart
#if CC_FieldUpgradeData
case TPM_CC_FieldUpgradeData: {
    FieldUpgradeData_In *in = (FieldUpgradeData_In *)
            MemoryGetInBuffer(sizeof(FieldUpgradeData_In));
    FieldUpgradeData_Out *out = (FieldUpgradeData_Out *) 
            MemoryGetOutBuffer(sizeof(FieldUpgradeData_Out));
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->fuData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_FieldUpgradeData_fuData);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_FieldUpgradeData (in, out);
    rSize = sizeof(FieldUpgradeData_Out);
    *respParmSize += TPMT_HA_Marshal(&out->nextDigest, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_HA_Marshal(&out->firstDigest, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_FieldUpgradeData
#if CC_FirmwareRead
case TPM_CC_FirmwareRead: {
    FirmwareRead_In *in = (FirmwareRead_In *)
            MemoryGetInBuffer(sizeof(FirmwareRead_In));
    FirmwareRead_Out *out = (FirmwareRead_Out *) 
            MemoryGetOutBuffer(sizeof(FirmwareRead_Out));
    result = UINT32_Unmarshal(&in->sequenceNumber, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_FirmwareRead_sequenceNumber);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_FirmwareRead (in, out);
    rSize = sizeof(FirmwareRead_Out);
    *respParmSize += TPM2B_MAX_BUFFER_Marshal(&out->fuData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_FirmwareRead
#if CC_ContextSave
case TPM_CC_ContextSave: {
    ContextSave_In *in = (ContextSave_In *)
            MemoryGetInBuffer(sizeof(ContextSave_In));
    ContextSave_Out *out = (ContextSave_Out *) 
            MemoryGetOutBuffer(sizeof(ContextSave_Out));
    in->saveHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ContextSave (in, out);
    rSize = sizeof(ContextSave_Out);
    *respParmSize += TPMS_CONTEXT_Marshal(&out->context, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ContextSave
#if CC_ContextLoad
case TPM_CC_ContextLoad: {
    ContextLoad_In *in = (ContextLoad_In *)
            MemoryGetInBuffer(sizeof(ContextLoad_In));
    ContextLoad_Out *out = (ContextLoad_Out *) 
            MemoryGetOutBuffer(sizeof(ContextLoad_Out));
    result = TPMS_CONTEXT_Unmarshal(&in->context, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ContextLoad_context);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ContextLoad (in, out);
    rSize = sizeof(ContextLoad_Out);
    if(TPM_RC_SUCCESS != result) goto Exit;
;    command->handles[command->handleNum++] = out->loadedHandle;
break; 
}
#endif     // CC_ContextLoad
#if CC_FlushContext
case TPM_CC_FlushContext: {
    FlushContext_In *in = (FlushContext_In *)
            MemoryGetInBuffer(sizeof(FlushContext_In));
    result = TPMI_DH_CONTEXT_Unmarshal(&in->flushHandle, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_FlushContext_flushHandle);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_FlushContext (in);
break; 
}
#endif     // CC_FlushContext
#if CC_EvictControl
case TPM_CC_EvictControl: {
    EvictControl_In *in = (EvictControl_In *)
            MemoryGetInBuffer(sizeof(EvictControl_In));
    in->auth = handles[0];
    in->objectHandle = handles[1];
    result = TPMI_DH_PERSISTENT_Unmarshal(&in->persistentHandle, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_EvictControl_persistentHandle);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_EvictControl (in);
break; 
}
#endif     // CC_EvictControl
#if CC_ReadClock
case TPM_CC_ReadClock: {
    ReadClock_Out *out = (ReadClock_Out *) 
            MemoryGetOutBuffer(sizeof(ReadClock_Out));
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ReadClock (out);
    rSize = sizeof(ReadClock_Out);
    *respParmSize += TPMS_TIME_INFO_Marshal(&out->currentTime, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_ReadClock
#if CC_ClockSet
case TPM_CC_ClockSet: {
    ClockSet_In *in = (ClockSet_In *)
            MemoryGetInBuffer(sizeof(ClockSet_In));
    in->auth = handles[0];
    result = UINT64_Unmarshal(&in->newTime, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ClockSet_newTime);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ClockSet (in);
break; 
}
#endif     // CC_ClockSet
#if CC_ClockRateAdjust
case TPM_CC_ClockRateAdjust: {
    ClockRateAdjust_In *in = (ClockRateAdjust_In *)
            MemoryGetInBuffer(sizeof(ClockRateAdjust_In));
    in->auth = handles[0];
    result = TPM_CLOCK_ADJUST_Unmarshal(&in->rateAdjust, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_ClockRateAdjust_rateAdjust);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_ClockRateAdjust (in);
break; 
}
#endif     // CC_ClockRateAdjust
#if CC_GetCapability
case TPM_CC_GetCapability: {
    GetCapability_In *in = (GetCapability_In *)
            MemoryGetInBuffer(sizeof(GetCapability_In));
    GetCapability_Out *out = (GetCapability_Out *) 
            MemoryGetOutBuffer(sizeof(GetCapability_Out));
    result = TPM_CAP_Unmarshal(&in->capability, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetCapability_capability);
    result = UINT32_Unmarshal(&in->property, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetCapability_property);
    result = UINT32_Unmarshal(&in->propertyCount, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_GetCapability_propertyCount);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_GetCapability (in, out);
    rSize = sizeof(GetCapability_Out);
    *respParmSize += TPMI_YES_NO_Marshal(&out->moreData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMS_CAPABILITY_DATA_Marshal(&out->capabilityData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_GetCapability
#if CC_TestParms
case TPM_CC_TestParms: {
    TestParms_In *in = (TestParms_In *)
            MemoryGetInBuffer(sizeof(TestParms_In));
    result = TPMT_PUBLIC_PARMS_Unmarshal(&in->parameters, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_TestParms_parameters);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_TestParms (in);
break; 
}
#endif     // CC_TestParms
#if CC_NV_DefineSpace
case TPM_CC_NV_DefineSpace: {
    NV_DefineSpace_In *in = (NV_DefineSpace_In *)
            MemoryGetInBuffer(sizeof(NV_DefineSpace_In));
    in->authHandle = handles[0];
    result = TPM2B_AUTH_Unmarshal(&in->auth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_DefineSpace_auth);
    result = TPM2B_NV_PUBLIC_Unmarshal(&in->publicInfo, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_DefineSpace_publicInfo);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_DefineSpace (in);
break; 
}
#endif     // CC_NV_DefineSpace
#if CC_NV_UndefineSpace
case TPM_CC_NV_UndefineSpace: {
    NV_UndefineSpace_In *in = (NV_UndefineSpace_In *)
            MemoryGetInBuffer(sizeof(NV_UndefineSpace_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_UndefineSpace (in);
break; 
}
#endif     // CC_NV_UndefineSpace
#if CC_NV_UndefineSpaceSpecial
case TPM_CC_NV_UndefineSpaceSpecial: {
    NV_UndefineSpaceSpecial_In *in = (NV_UndefineSpaceSpecial_In *)
            MemoryGetInBuffer(sizeof(NV_UndefineSpaceSpecial_In));
    in->nvIndex = handles[0];
    in->platform = handles[1];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_UndefineSpaceSpecial (in);
break; 
}
#endif     // CC_NV_UndefineSpaceSpecial
#if CC_NV_ReadPublic
case TPM_CC_NV_ReadPublic: {
    NV_ReadPublic_In *in = (NV_ReadPublic_In *)
            MemoryGetInBuffer(sizeof(NV_ReadPublic_In));
    NV_ReadPublic_Out *out = (NV_ReadPublic_Out *) 
            MemoryGetOutBuffer(sizeof(NV_ReadPublic_Out));
    in->nvIndex = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_ReadPublic (in, out);
    rSize = sizeof(NV_ReadPublic_Out);
    *respParmSize += TPM2B_NV_PUBLIC_Marshal(&out->nvPublic, 
                                          responseBuffer, &rSize);
    *respParmSize += TPM2B_NAME_Marshal(&out->nvName, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_NV_ReadPublic
#if CC_NV_Write
case TPM_CC_NV_Write: {
    NV_Write_In *in = (NV_Write_In *)
            MemoryGetInBuffer(sizeof(NV_Write_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    result = TPM2B_MAX_NV_BUFFER_Unmarshal(&in->data, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Write_data);
    result = UINT16_Unmarshal(&in->offset, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Write_offset);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_Write (in);
break; 
}
#endif     // CC_NV_Write
#if CC_NV_Increment
case TPM_CC_NV_Increment: {
    NV_Increment_In *in = (NV_Increment_In *)
            MemoryGetInBuffer(sizeof(NV_Increment_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_Increment (in);
break; 
}
#endif     // CC_NV_Increment
#if CC_NV_Extend
case TPM_CC_NV_Extend: {
    NV_Extend_In *in = (NV_Extend_In *)
            MemoryGetInBuffer(sizeof(NV_Extend_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    result = TPM2B_MAX_NV_BUFFER_Unmarshal(&in->data, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Extend_data);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_Extend (in);
break; 
}
#endif     // CC_NV_Extend
#if CC_NV_SetBits
case TPM_CC_NV_SetBits: {
    NV_SetBits_In *in = (NV_SetBits_In *)
            MemoryGetInBuffer(sizeof(NV_SetBits_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    result = UINT64_Unmarshal(&in->bits, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_SetBits_bits);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_SetBits (in);
break; 
}
#endif     // CC_NV_SetBits
#if CC_NV_WriteLock
case TPM_CC_NV_WriteLock: {
    NV_WriteLock_In *in = (NV_WriteLock_In *)
            MemoryGetInBuffer(sizeof(NV_WriteLock_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_WriteLock (in);
break; 
}
#endif     // CC_NV_WriteLock
#if CC_NV_GlobalWriteLock
case TPM_CC_NV_GlobalWriteLock: {
    NV_GlobalWriteLock_In *in = (NV_GlobalWriteLock_In *)
            MemoryGetInBuffer(sizeof(NV_GlobalWriteLock_In));
    in->authHandle = handles[0];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_GlobalWriteLock (in);
break; 
}
#endif     // CC_NV_GlobalWriteLock
#if CC_NV_Read
case TPM_CC_NV_Read: {
    NV_Read_In *in = (NV_Read_In *)
            MemoryGetInBuffer(sizeof(NV_Read_In));
    NV_Read_Out *out = (NV_Read_Out *) 
            MemoryGetOutBuffer(sizeof(NV_Read_Out));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    result = UINT16_Unmarshal(&in->size, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Read_size);
    result = UINT16_Unmarshal(&in->offset, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Read_offset);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_Read (in, out);
    rSize = sizeof(NV_Read_Out);
    *respParmSize += TPM2B_MAX_NV_BUFFER_Marshal(&out->data, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_NV_Read
#if CC_NV_ReadLock
case TPM_CC_NV_ReadLock: {
    NV_ReadLock_In *in = (NV_ReadLock_In *)
            MemoryGetInBuffer(sizeof(NV_ReadLock_In));
    in->authHandle = handles[0];
    in->nvIndex = handles[1];
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_ReadLock (in);
break; 
}
#endif     // CC_NV_ReadLock
#if CC_NV_ChangeAuth
case TPM_CC_NV_ChangeAuth: {
    NV_ChangeAuth_In *in = (NV_ChangeAuth_In *)
            MemoryGetInBuffer(sizeof(NV_ChangeAuth_In));
    in->nvIndex = handles[0];
    result = TPM2B_AUTH_Unmarshal(&in->newAuth, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_ChangeAuth_newAuth);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_ChangeAuth (in);
break; 
}
#endif     // CC_NV_ChangeAuth
#if CC_NV_Certify
case TPM_CC_NV_Certify: {
    NV_Certify_In *in = (NV_Certify_In *)
            MemoryGetInBuffer(sizeof(NV_Certify_In));
    NV_Certify_Out *out = (NV_Certify_Out *) 
            MemoryGetOutBuffer(sizeof(NV_Certify_Out));
    in->signHandle = handles[0];
    in->authHandle = handles[1];
    in->nvIndex = handles[2];
    result = TPM2B_DATA_Unmarshal(&in->qualifyingData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Certify_qualifyingData);
    result = TPMT_SIG_SCHEME_Unmarshal(&in->inScheme, paramBuffer, paramBufferSize, TRUE);
        ERROR_IF_EXIT_PLUS(RC_NV_Certify_inScheme);
    result = UINT16_Unmarshal(&in->size, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Certify_size);
    result = UINT16_Unmarshal(&in->offset, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_NV_Certify_offset);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_NV_Certify (in, out);
    rSize = sizeof(NV_Certify_Out);
    *respParmSize += TPM2B_ATTEST_Marshal(&out->certifyInfo, 
                                          responseBuffer, &rSize);
    *respParmSize += TPMT_SIGNATURE_Marshal(&out->signature, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_NV_Certify
#if CC_AC_GetCapability
case TPM_CC_AC_GetCapability: {
    AC_GetCapability_In *in = (AC_GetCapability_In *)
            MemoryGetInBuffer(sizeof(AC_GetCapability_In));
    AC_GetCapability_Out *out = (AC_GetCapability_Out *) 
            MemoryGetOutBuffer(sizeof(AC_GetCapability_Out));
    in->ac = handles[0];
    result = TPM_AT_Unmarshal(&in->capability, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_AC_GetCapability_capability);
    result = UINT32_Unmarshal(&in->count, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_AC_GetCapability_count);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_AC_GetCapability (in, out);
    rSize = sizeof(AC_GetCapability_Out);
    *respParmSize += TPMI_YES_NO_Marshal(&out->moreData, 
                                          responseBuffer, &rSize);
    *respParmSize += TPML_AC_CAPABILITIES_Marshal(&out->capabilitiesData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_AC_GetCapability
#if CC_AC_Send
case TPM_CC_AC_Send: {
    AC_Send_In *in = (AC_Send_In *)
            MemoryGetInBuffer(sizeof(AC_Send_In));
    AC_Send_Out *out = (AC_Send_Out *) 
            MemoryGetOutBuffer(sizeof(AC_Send_Out));
    in->sendObject = handles[0];
    in->authHandle = handles[1];
    in->ac = handles[2];
    result = TPM2B_MAX_BUFFER_Unmarshal(&in->acDataIn, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_AC_Send_acDataIn);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_AC_Send (in, out);
    rSize = sizeof(AC_Send_Out);
    *respParmSize += TPMS_AC_OUTPUT_Marshal(&out->acDataOut, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_AC_Send
#if CC_Policy_AC_SendSelect
case TPM_CC_Policy_AC_SendSelect: {
    Policy_AC_SendSelect_In *in = (Policy_AC_SendSelect_In *)
            MemoryGetInBuffer(sizeof(Policy_AC_SendSelect_In));
    in->policySession = handles[0];
    result = TPM2B_NAME_Unmarshal(&in->objectName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Policy_AC_SendSelect_objectName);
    result = TPM2B_NAME_Unmarshal(&in->authHandleName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Policy_AC_SendSelect_authHandleName);
    result = TPM2B_NAME_Unmarshal(&in->acName, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Policy_AC_SendSelect_acName);
    result = TPMI_YES_NO_Unmarshal(&in->includeObject, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Policy_AC_SendSelect_includeObject);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Policy_AC_SendSelect (in);
break; 
}
#endif     // CC_Policy_AC_SendSelect
#if CC_Vendor_TCG_Test
case TPM_CC_Vendor_TCG_Test: {
    Vendor_TCG_Test_In *in = (Vendor_TCG_Test_In *)
            MemoryGetInBuffer(sizeof(Vendor_TCG_Test_In));
    Vendor_TCG_Test_Out *out = (Vendor_TCG_Test_Out *) 
            MemoryGetOutBuffer(sizeof(Vendor_TCG_Test_Out));
    result = TPM2B_DATA_Unmarshal(&in->inputData, paramBuffer, paramBufferSize);
        ERROR_IF_EXIT_PLUS(RC_Vendor_TCG_Test_inputData);
    if(*paramBufferSize != 0) (result = TPM_RC_SIZE; goto Exit; }
result = TPM2_Vendor_TCG_Test (in, out);
    rSize = sizeof(Vendor_TCG_Test_Out);
    *respParmSize += TPM2B_DATA_Marshal(&out->outputData, 
                                          responseBuffer, &rSize);
break; 
}
#endif     // CC_Vendor_TCG_Test
