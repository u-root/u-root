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
 *  Date: Apr  2, 2019  Time: 11:00:48AM
 */

// This file should only be included by CommandCodeAttibutes.c
#ifdef _COMMAND_TABLE_DISPATCH_

   
// Define the stop value
#define END_OF_LIST     0xff
#define ADD_FLAG        0x80

// These macros provide some variability in how the data is encoded. They also make
// the lines a little sorter. ;-)
#   define UNMARSHAL_DISPATCH(name)   (UNMARSHAL_t)name##_Unmarshal
#   define MARSHAL_DISPATCH(name)     (MARSHAL_t)name##_Marshal
#   define _UNMARSHAL_T_    UNMARSHAL_t
#   define _MARSHAL_T_      MARSHAL_t


// The UnmarshalArray contains the dispatch functions for the unmarshaling code.
// The defines in this array are used to make it easier to cross reference the
// unmarshaling values in the types array of each command

const _UNMARSHAL_T_ UnmarshalArray[] = {
#define TPMI_DH_CONTEXT_H_UNMARSHAL         0
            UNMARSHAL_DISPATCH(TPMI_DH_CONTEXT),
#define TPMI_RH_AC_H_UNMARSHAL              (TPMI_DH_CONTEXT_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_AC),
#define TPMI_RH_CLEAR_H_UNMARSHAL           (TPMI_RH_AC_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_CLEAR),
#define TPMI_RH_HIERARCHY_AUTH_H_UNMARSHAL  (TPMI_RH_CLEAR_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_HIERARCHY_AUTH),
#define TPMI_RH_LOCKOUT_H_UNMARSHAL         (TPMI_RH_HIERARCHY_AUTH_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_LOCKOUT),
#define TPMI_RH_NV_AUTH_H_UNMARSHAL         (TPMI_RH_LOCKOUT_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_NV_AUTH),
#define TPMI_RH_NV_INDEX_H_UNMARSHAL        (TPMI_RH_NV_AUTH_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_NV_INDEX),
#define TPMI_RH_PLATFORM_H_UNMARSHAL        (TPMI_RH_NV_INDEX_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_PLATFORM),
#define TPMI_RH_PROVISION_H_UNMARSHAL       (TPMI_RH_PLATFORM_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_PROVISION),
#define TPMI_SH_HMAC_H_UNMARSHAL            (TPMI_RH_PROVISION_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_SH_HMAC),
#define TPMI_SH_POLICY_H_UNMARSHAL          (TPMI_SH_HMAC_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_SH_POLICY),
// HANDLE_FIRST_FLAG_TYPE is the first handle that needs a flag when called.
#define HANDLE_FIRST_FLAG_TYPE              (TPMI_SH_POLICY_H_UNMARSHAL + 1)
#define TPMI_DH_ENTITY_H_UNMARSHAL          (TPMI_SH_POLICY_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_ENTITY),
#define TPMI_DH_OBJECT_H_UNMARSHAL          (TPMI_DH_ENTITY_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_OBJECT),
#define TPMI_DH_PARENT_H_UNMARSHAL          (TPMI_DH_OBJECT_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_PARENT),
#define TPMI_DH_PCR_H_UNMARSHAL             (TPMI_DH_PARENT_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_PCR),
#define TPMI_RH_ENDORSEMENT_H_UNMARSHAL     (TPMI_DH_PCR_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_ENDORSEMENT),
#define TPMI_RH_HIERARCHY_H_UNMARSHAL       (TPMI_RH_ENDORSEMENT_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_HIERARCHY),
// PARAMETER_FIRST_TYPE marks the end of the handle list.
#define PARAMETER_FIRST_TYPE                (TPMI_RH_HIERARCHY_H_UNMARSHAL + 1)
#define TPM2B_DATA_P_UNMARSHAL              (TPMI_RH_HIERARCHY_H_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_DATA),
#define TPM2B_DIGEST_P_UNMARSHAL            (TPM2B_DATA_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_DIGEST),
#define TPM2B_ECC_PARAMETER_P_UNMARSHAL     (TPM2B_DIGEST_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_ECC_PARAMETER),
#define TPM2B_ECC_POINT_P_UNMARSHAL         (TPM2B_ECC_PARAMETER_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_ECC_POINT),
#define TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL  (TPM2B_ECC_POINT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_ENCRYPTED_SECRET),
#define TPM2B_EVENT_P_UNMARSHAL             (TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_EVENT),
#define TPM2B_ID_OBJECT_P_UNMARSHAL         (TPM2B_EVENT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_ID_OBJECT),
#define TPM2B_IV_P_UNMARSHAL                (TPM2B_ID_OBJECT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_IV),
#define TPM2B_MAX_BUFFER_P_UNMARSHAL        (TPM2B_IV_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_MAX_BUFFER),
#define TPM2B_MAX_NV_BUFFER_P_UNMARSHAL     (TPM2B_MAX_BUFFER_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_MAX_NV_BUFFER),
#define TPM2B_NAME_P_UNMARSHAL              (TPM2B_MAX_NV_BUFFER_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_NAME),
#define TPM2B_NV_PUBLIC_P_UNMARSHAL         (TPM2B_NAME_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_NV_PUBLIC),
#define TPM2B_PRIVATE_P_UNMARSHAL           (TPM2B_NV_PUBLIC_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_PRIVATE),
#define TPM2B_PUBLIC_KEY_RSA_P_UNMARSHAL    (TPM2B_PRIVATE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_PUBLIC_KEY_RSA),
#define TPM2B_SENSITIVE_P_UNMARSHAL         (TPM2B_PUBLIC_KEY_RSA_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_SENSITIVE),
#define TPM2B_SENSITIVE_CREATE_P_UNMARSHAL  (TPM2B_SENSITIVE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_SENSITIVE_CREATE),
#define TPM2B_SENSITIVE_DATA_P_UNMARSHAL    (TPM2B_SENSITIVE_CREATE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_SENSITIVE_DATA),
#define TPM2B_TEMPLATE_P_UNMARSHAL          (TPM2B_SENSITIVE_DATA_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_TEMPLATE),
#define TPM2B_TIMEOUT_P_UNMARSHAL           (TPM2B_TEMPLATE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_TIMEOUT),
#define TPMI_DH_CONTEXT_P_UNMARSHAL         (TPM2B_TIMEOUT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_CONTEXT),
#define TPMI_DH_PERSISTENT_P_UNMARSHAL      (TPMI_DH_CONTEXT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_PERSISTENT),
#define TPMI_ECC_CURVE_P_UNMARSHAL          (TPMI_DH_PERSISTENT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_ECC_CURVE),
#define TPMI_YES_NO_P_UNMARSHAL             (TPMI_ECC_CURVE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_YES_NO),
#define TPML_ALG_P_UNMARSHAL                (TPMI_YES_NO_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPML_ALG),
#define TPML_CC_P_UNMARSHAL                 (TPML_ALG_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPML_CC),
#define TPML_DIGEST_P_UNMARSHAL             (TPML_CC_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPML_DIGEST),
#define TPML_DIGEST_VALUES_P_UNMARSHAL      (TPML_DIGEST_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPML_DIGEST_VALUES),
#define TPML_PCR_SELECTION_P_UNMARSHAL      (TPML_DIGEST_VALUES_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPML_PCR_SELECTION),
#define TPMS_CONTEXT_P_UNMARSHAL            (TPML_PCR_SELECTION_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMS_CONTEXT),
#define TPMT_PUBLIC_PARMS_P_UNMARSHAL       (TPMS_CONTEXT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_PUBLIC_PARMS),
#define TPMT_TK_AUTH_P_UNMARSHAL            (TPMT_PUBLIC_PARMS_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_TK_AUTH),
#define TPMT_TK_CREATION_P_UNMARSHAL        (TPMT_TK_AUTH_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_TK_CREATION),
#define TPMT_TK_HASHCHECK_P_UNMARSHAL       (TPMT_TK_CREATION_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_TK_HASHCHECK),
#define TPMT_TK_VERIFIED_P_UNMARSHAL        (TPMT_TK_HASHCHECK_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_TK_VERIFIED),
#define TPM_AT_P_UNMARSHAL                  (TPMT_TK_VERIFIED_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_AT),
#define TPM_CAP_P_UNMARSHAL                 (TPM_AT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_CAP),
#define TPM_CLOCK_ADJUST_P_UNMARSHAL        (TPM_CAP_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_CLOCK_ADJUST),
#define TPM_EO_P_UNMARSHAL                  (TPM_CLOCK_ADJUST_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_EO),
#define TPM_SE_P_UNMARSHAL                  (TPM_EO_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_SE),
#define TPM_SU_P_UNMARSHAL                  (TPM_SE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM_SU),
#define UINT16_P_UNMARSHAL                  (TPM_SU_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(UINT16),
#define UINT32_P_UNMARSHAL                  (UINT16_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(UINT32),
#define UINT64_P_UNMARSHAL                  (UINT32_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(UINT64),
#define UINT8_P_UNMARSHAL                   (UINT64_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(UINT8),
// PARAMETER_FIRST_FLAG_TYPE is the first parameter to need a flag.
#define PARAMETER_FIRST_FLAG_TYPE           (UINT8_P_UNMARSHAL + 1)
#define TPM2B_PUBLIC_P_UNMARSHAL            (UINT8_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPM2B_PUBLIC),
#define TPMI_ALG_CIPHER_MODE_P_UNMARSHAL    (TPM2B_PUBLIC_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_ALG_CIPHER_MODE),
#define TPMI_ALG_HASH_P_UNMARSHAL           (TPMI_ALG_CIPHER_MODE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_ALG_HASH),
#define TPMI_ALG_MAC_SCHEME_P_UNMARSHAL     (TPMI_ALG_HASH_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_ALG_MAC_SCHEME),
#define TPMI_DH_PCR_P_UNMARSHAL             (TPMI_ALG_MAC_SCHEME_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_DH_PCR),
#define TPMI_ECC_KEY_EXCHANGE_P_UNMARSHAL   (TPMI_DH_PCR_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_ECC_KEY_EXCHANGE),
#define TPMI_RH_ENABLES_P_UNMARSHAL         (TPMI_ECC_KEY_EXCHANGE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_ENABLES),
#define TPMI_RH_HIERARCHY_P_UNMARSHAL       (TPMI_RH_ENABLES_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMI_RH_HIERARCHY),
#define TPMT_RSA_DECRYPT_P_UNMARSHAL        (TPMI_RH_HIERARCHY_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_RSA_DECRYPT),
#define TPMT_SIGNATURE_P_UNMARSHAL          (TPMT_RSA_DECRYPT_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_SIGNATURE),
#define TPMT_SIG_SCHEME_P_UNMARSHAL         (TPMT_SIGNATURE_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_SIG_SCHEME),
#define TPMT_SYM_DEF_P_UNMARSHAL            (TPMT_SIG_SCHEME_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_SYM_DEF),
#define TPMT_SYM_DEF_OBJECT_P_UNMARSHAL     (TPMT_SYM_DEF_P_UNMARSHAL + 1)
            UNMARSHAL_DISPATCH(TPMT_SYM_DEF_OBJECT)
// PARAMETER_LAST_TYPE is the end of the command parameter list.
#define PARAMETER_LAST_TYPE                 (TPMT_SYM_DEF_OBJECT_P_UNMARSHAL)
};
   
// The MarshalArray contains the dispatch functions for the marshaling code.
// The defines in this array are used to make it easier to cross reference the
// marshaling values in the types array of each command
const _MARSHAL_T_ MarshalArray[] = {

#define UINT32_H_MARSHAL                    0
            MARSHAL_DISPATCH(UINT32),
// RESPONSE_PARAMETER_FIRST_TYPE marks the end of the response handles.
#define RESPONSE_PARAMETER_FIRST_TYPE       (UINT32_H_MARSHAL + 1)
#define TPM2B_ATTEST_P_MARSHAL              (UINT32_H_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_ATTEST),
#define TPM2B_CREATION_DATA_P_MARSHAL       (TPM2B_ATTEST_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_CREATION_DATA),
#define TPM2B_DATA_P_MARSHAL                (TPM2B_CREATION_DATA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_DATA),
#define TPM2B_DIGEST_P_MARSHAL              (TPM2B_DATA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_DIGEST),
#define TPM2B_ECC_POINT_P_MARSHAL           (TPM2B_DIGEST_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_ECC_POINT),
#define TPM2B_ENCRYPTED_SECRET_P_MARSHAL    (TPM2B_ECC_POINT_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_ENCRYPTED_SECRET),
#define TPM2B_ID_OBJECT_P_MARSHAL           (TPM2B_ENCRYPTED_SECRET_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_ID_OBJECT),
#define TPM2B_IV_P_MARSHAL                  (TPM2B_ID_OBJECT_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_IV),
#define TPM2B_MAX_BUFFER_P_MARSHAL          (TPM2B_IV_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_MAX_BUFFER),
#define TPM2B_MAX_NV_BUFFER_P_MARSHAL       (TPM2B_MAX_BUFFER_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_MAX_NV_BUFFER),
#define TPM2B_NAME_P_MARSHAL                (TPM2B_MAX_NV_BUFFER_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_NAME),
#define TPM2B_NV_PUBLIC_P_MARSHAL           (TPM2B_NAME_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_NV_PUBLIC),
#define TPM2B_PRIVATE_P_MARSHAL             (TPM2B_NV_PUBLIC_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_PRIVATE),
#define TPM2B_PUBLIC_P_MARSHAL              (TPM2B_PRIVATE_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_PUBLIC),
#define TPM2B_PUBLIC_KEY_RSA_P_MARSHAL      (TPM2B_PUBLIC_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_PUBLIC_KEY_RSA),
#define TPM2B_SENSITIVE_DATA_P_MARSHAL      (TPM2B_PUBLIC_KEY_RSA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_SENSITIVE_DATA),
#define TPM2B_TIMEOUT_P_MARSHAL             (TPM2B_SENSITIVE_DATA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPM2B_TIMEOUT),
#define UINT8_P_MARSHAL                     (TPM2B_TIMEOUT_P_MARSHAL + 1)
            MARSHAL_DISPATCH(UINT8),
#define TPML_AC_CAPABILITIES_P_MARSHAL      (UINT8_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPML_AC_CAPABILITIES),
#define TPML_ALG_P_MARSHAL                  (TPML_AC_CAPABILITIES_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPML_ALG),
#define TPML_DIGEST_P_MARSHAL               (TPML_ALG_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPML_DIGEST),
#define TPML_DIGEST_VALUES_P_MARSHAL        (TPML_DIGEST_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPML_DIGEST_VALUES),
#define TPML_PCR_SELECTION_P_MARSHAL        (TPML_DIGEST_VALUES_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPML_PCR_SELECTION),
#define TPMS_AC_OUTPUT_P_MARSHAL            (TPML_PCR_SELECTION_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMS_AC_OUTPUT),
#define TPMS_ALGORITHM_DETAIL_ECC_P_MARSHAL (TPMS_AC_OUTPUT_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMS_ALGORITHM_DETAIL_ECC),
#define TPMS_CAPABILITY_DATA_P_MARSHAL      \
            (TPMS_ALGORITHM_DETAIL_ECC_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMS_CAPABILITY_DATA),
#define TPMS_CONTEXT_P_MARSHAL              (TPMS_CAPABILITY_DATA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMS_CONTEXT),
#define TPMS_TIME_INFO_P_MARSHAL            (TPMS_CONTEXT_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMS_TIME_INFO),
#define TPMT_HA_P_MARSHAL                   (TPMS_TIME_INFO_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_HA),
#define TPMT_SIGNATURE_P_MARSHAL            (TPMT_HA_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_SIGNATURE),
#define TPMT_TK_AUTH_P_MARSHAL              (TPMT_SIGNATURE_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_TK_AUTH),
#define TPMT_TK_CREATION_P_MARSHAL          (TPMT_TK_AUTH_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_TK_CREATION),
#define TPMT_TK_HASHCHECK_P_MARSHAL         (TPMT_TK_CREATION_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_TK_HASHCHECK),
#define TPMT_TK_VERIFIED_P_MARSHAL          (TPMT_TK_HASHCHECK_P_MARSHAL + 1)
            MARSHAL_DISPATCH(TPMT_TK_VERIFIED),
#define UINT32_P_MARSHAL                    (TPMT_TK_VERIFIED_P_MARSHAL + 1)
            MARSHAL_DISPATCH(UINT32),
#define UINT16_P_MARSHAL                    (UINT32_P_MARSHAL + 1)
            MARSHAL_DISPATCH(UINT16)
// RESPONSE_PARAMETER_LAST_TYPE is the end of the response parameter list.
#define RESPONSE_PARAMETER_LAST_TYPE        (UINT16_P_MARSHAL)
};

// This list of aliases allows the types in the _COMMAND_DESCRIPTOR_T to match the
// types in the command/response templates of part 3.
#define INT32_P_UNMARSHAL                   UINT32_P_UNMARSHAL
#define TPM2B_AUTH_P_UNMARSHAL              TPM2B_DIGEST_P_UNMARSHAL
#define TPM2B_NONCE_P_UNMARSHAL             TPM2B_DIGEST_P_UNMARSHAL
#define TPM2B_OPERAND_P_UNMARSHAL           TPM2B_DIGEST_P_UNMARSHAL
#define TPMA_LOCALITY_P_UNMARSHAL           UINT8_P_UNMARSHAL
#define TPM_CC_P_UNMARSHAL                  UINT32_P_UNMARSHAL
#define TPMI_DH_CONTEXT_H_MARSHAL           UINT32_H_MARSHAL
#define TPMI_DH_OBJECT_H_MARSHAL            UINT32_H_MARSHAL
#define TPMI_SH_AUTH_SESSION_H_MARSHAL      UINT32_H_MARSHAL
#define TPM_HANDLE_H_MARSHAL                UINT32_H_MARSHAL
#define TPM2B_NONCE_P_MARSHAL               TPM2B_DIGEST_P_MARSHAL
#define TPMI_YES_NO_P_MARSHAL               UINT8_P_MARSHAL
#define TPM_RC_P_MARSHAL                    UINT32_P_MARSHAL


#if CC_Startup

#include "Startup_fp.h"

typedef TPM_RC  (Startup_Entry)(
    Startup_In                  *in
);

typedef const struct {
    Startup_Entry           *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} Startup_COMMAND_DESCRIPTOR_t;

Startup_COMMAND_DESCRIPTOR_t _StartupData = {
    /* entry         */     &TPM2_Startup,
    /* inSize        */     (UINT16)(sizeof(Startup_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(Startup_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPM_SU_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _StartupDataAddress (&_StartupData)
#else
#define _StartupDataAddress 0
#endif // CC_Startup

#if CC_Shutdown

#include "Shutdown_fp.h"

typedef TPM_RC  (Shutdown_Entry)(
    Shutdown_In                 *in
);

typedef const struct {
    Shutdown_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} Shutdown_COMMAND_DESCRIPTOR_t;

Shutdown_COMMAND_DESCRIPTOR_t _ShutdownData = {
    /* entry         */     &TPM2_Shutdown,
    /* inSize        */     (UINT16)(sizeof(Shutdown_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(Shutdown_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPM_SU_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ShutdownDataAddress (&_ShutdownData)
#else
#define _ShutdownDataAddress 0
#endif // CC_Shutdown

#if CC_SelfTest

#include "SelfTest_fp.h"

typedef TPM_RC  (SelfTest_Entry)(
    SelfTest_In                 *in
);

typedef const struct {
    SelfTest_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} SelfTest_COMMAND_DESCRIPTOR_t;

SelfTest_COMMAND_DESCRIPTOR_t _SelfTestData = {
    /* entry         */     &TPM2_SelfTest,
    /* inSize        */     (UINT16)(sizeof(SelfTest_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(SelfTest_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_YES_NO_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _SelfTestDataAddress (&_SelfTestData)
#else
#define _SelfTestDataAddress 0
#endif // CC_SelfTest

#if CC_IncrementalSelfTest

#include "IncrementalSelfTest_fp.h"

typedef TPM_RC  (IncrementalSelfTest_Entry)(
    IncrementalSelfTest_In          *in,
    IncrementalSelfTest_Out         *out
);

typedef const struct {
    IncrementalSelfTest_Entry   *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    BYTE                        types[4];
} IncrementalSelfTest_COMMAND_DESCRIPTOR_t;

IncrementalSelfTest_COMMAND_DESCRIPTOR_t _IncrementalSelfTestData = {
    /* entry         */         &TPM2_IncrementalSelfTest,
    /* inSize        */         (UINT16)(sizeof(IncrementalSelfTest_In)),
    /* outSize       */         (UINT16)(sizeof(IncrementalSelfTest_Out)),
    /* offsetOfTypes */         offsetof(IncrementalSelfTest_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         // No parameter offsets;
    /* types         */         {TPML_ALG_P_UNMARSHAL,
                                 END_OF_LIST,
                                 TPML_ALG_P_MARSHAL,
                                 END_OF_LIST}
};

#define _IncrementalSelfTestDataAddress (&_IncrementalSelfTestData)
#else
#define _IncrementalSelfTestDataAddress 0
#endif // CC_IncrementalSelfTest

#if CC_GetTestResult

#include "GetTestResult_fp.h"

typedef TPM_RC  (GetTestResult_Entry)(
    GetTestResult_Out           *out
);

typedef const struct {
    GetTestResult_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} GetTestResult_COMMAND_DESCRIPTOR_t;

GetTestResult_COMMAND_DESCRIPTOR_t _GetTestResultData = {
    /* entry         */     &TPM2_GetTestResult,
    /* inSize        */     0,
    /* outSize       */     (UINT16)(sizeof(GetTestResult_Out)),
    /* offsetOfTypes */     offsetof(GetTestResult_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(GetTestResult_Out, testResult))},
    /* types         */     {END_OF_LIST,
                             TPM2B_MAX_BUFFER_P_MARSHAL,
                             TPM_RC_P_MARSHAL,
                             END_OF_LIST}
};

#define _GetTestResultDataAddress (&_GetTestResultData)
#else
#define _GetTestResultDataAddress 0
#endif // CC_GetTestResult

#if CC_StartAuthSession

#include "StartAuthSession_fp.h"

typedef TPM_RC  (StartAuthSession_Entry)(
    StartAuthSession_In         *in,
    StartAuthSession_Out        *out
);

typedef const struct {
    StartAuthSession_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[7];
    BYTE                    types[11];
} StartAuthSession_COMMAND_DESCRIPTOR_t;

StartAuthSession_COMMAND_DESCRIPTOR_t _StartAuthSessionData = {
    /* entry         */     &TPM2_StartAuthSession,
    /* inSize        */     (UINT16)(sizeof(StartAuthSession_In)),
    /* outSize       */     (UINT16)(sizeof(StartAuthSession_Out)),
    /* offsetOfTypes */     offsetof(StartAuthSession_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(StartAuthSession_In, bind)),
                             (UINT16)(offsetof(StartAuthSession_In, nonceCaller)),
                             (UINT16)(offsetof(StartAuthSession_In, encryptedSalt)),
                             (UINT16)(offsetof(StartAuthSession_In, sessionType)),
                             (UINT16)(offsetof(StartAuthSession_In, symmetric)),
                             (UINT16)(offsetof(StartAuthSession_In, authHash)),
                             (UINT16)(offsetof(StartAuthSession_Out, nonceTPM))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPMI_DH_ENTITY_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_NONCE_P_UNMARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL,
                             TPM_SE_P_UNMARSHAL,
                             TPMT_SYM_DEF_P_UNMARSHAL + ADD_FLAG,
                             TPMI_ALG_HASH_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMI_SH_AUTH_SESSION_H_MARSHAL,
                             TPM2B_NONCE_P_MARSHAL,
                             END_OF_LIST}
};

#define _StartAuthSessionDataAddress (&_StartAuthSessionData)
#else
#define _StartAuthSessionDataAddress 0
#endif // CC_StartAuthSession

#if CC_PolicyRestart

#include "PolicyRestart_fp.h"

typedef TPM_RC  (PolicyRestart_Entry)(
    PolicyRestart_In            *in
);

typedef const struct {
    PolicyRestart_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} PolicyRestart_COMMAND_DESCRIPTOR_t;

PolicyRestart_COMMAND_DESCRIPTOR_t _PolicyRestartData = {
    /* entry         */     &TPM2_PolicyRestart,
    /* inSize        */     (UINT16)(sizeof(PolicyRestart_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyRestart_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyRestartDataAddress (&_PolicyRestartData)
#else
#define _PolicyRestartDataAddress 0
#endif // CC_PolicyRestart

#if CC_Create

#include "Create_fp.h"

typedef TPM_RC  (Create_Entry)(
    Create_In                   *in,
    Create_Out                  *out
);

typedef const struct {
    Create_Entry            *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[8];
    BYTE                    types[12];
} Create_COMMAND_DESCRIPTOR_t;

Create_COMMAND_DESCRIPTOR_t _CreateData = {
    /* entry         */     &TPM2_Create,
    /* inSize        */     (UINT16)(sizeof(Create_In)),
    /* outSize       */     (UINT16)(sizeof(Create_Out)),
    /* offsetOfTypes */     offsetof(Create_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Create_In, inSensitive)),
                             (UINT16)(offsetof(Create_In, inPublic)),
                             (UINT16)(offsetof(Create_In, outsideInfo)),
                             (UINT16)(offsetof(Create_In, creationPCR)),
                             (UINT16)(offsetof(Create_Out, outPublic)),
                             (UINT16)(offsetof(Create_Out, creationData)),
                             (UINT16)(offsetof(Create_Out, creationHash)),
                             (UINT16)(offsetof(Create_Out, creationTicket))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_SENSITIVE_CREATE_P_UNMARSHAL,
                             TPM2B_PUBLIC_P_UNMARSHAL,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PRIVATE_P_MARSHAL,
                             TPM2B_PUBLIC_P_MARSHAL,
                             TPM2B_CREATION_DATA_P_MARSHAL,
                             TPM2B_DIGEST_P_MARSHAL,
                             TPMT_TK_CREATION_P_MARSHAL,
                             END_OF_LIST}
};

#define _CreateDataAddress (&_CreateData)
#else
#define _CreateDataAddress 0
#endif // CC_Create

#if CC_Load

#include "Load_fp.h"

typedef TPM_RC  (Load_Entry)(
    Load_In                     *in,
    Load_Out                    *out
);

typedef const struct {
    Load_Entry              *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} Load_COMMAND_DESCRIPTOR_t;

Load_COMMAND_DESCRIPTOR_t _LoadData = {
    /* entry         */     &TPM2_Load,
    /* inSize        */     (UINT16)(sizeof(Load_In)),
    /* outSize       */     (UINT16)(sizeof(Load_Out)),
    /* offsetOfTypes */     offsetof(Load_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Load_In, inPrivate)),
                             (UINT16)(offsetof(Load_In, inPublic)),
                             (UINT16)(offsetof(Load_Out, name))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_PRIVATE_P_UNMARSHAL,
                             TPM2B_PUBLIC_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM_HANDLE_H_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _LoadDataAddress (&_LoadData)
#else
#define _LoadDataAddress 0
#endif // CC_Load

#if CC_LoadExternal

#include "LoadExternal_fp.h"

typedef TPM_RC  (LoadExternal_Entry)(
    LoadExternal_In             *in,
    LoadExternal_Out            *out
);

typedef const struct {
    LoadExternal_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} LoadExternal_COMMAND_DESCRIPTOR_t;

LoadExternal_COMMAND_DESCRIPTOR_t _LoadExternalData = {
    /* entry         */     &TPM2_LoadExternal,
    /* inSize        */     (UINT16)(sizeof(LoadExternal_In)),
    /* outSize       */     (UINT16)(sizeof(LoadExternal_Out)),
    /* offsetOfTypes */     offsetof(LoadExternal_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(LoadExternal_In, inPublic)),
                             (UINT16)(offsetof(LoadExternal_In, hierarchy)),
                             (UINT16)(offsetof(LoadExternal_Out, name))},
    /* types         */     {TPM2B_SENSITIVE_P_UNMARSHAL,
                             TPM2B_PUBLIC_P_UNMARSHAL + ADD_FLAG,
                             TPMI_RH_HIERARCHY_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM_HANDLE_H_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _LoadExternalDataAddress (&_LoadExternalData)
#else
#define _LoadExternalDataAddress 0
#endif // CC_LoadExternal

#if CC_ReadPublic

#include "ReadPublic_fp.h"

typedef TPM_RC  (ReadPublic_Entry)(
    ReadPublic_In               *in,
    ReadPublic_Out              *out
);

typedef const struct {
    ReadPublic_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} ReadPublic_COMMAND_DESCRIPTOR_t;

ReadPublic_COMMAND_DESCRIPTOR_t _ReadPublicData = {
    /* entry         */     &TPM2_ReadPublic,
    /* inSize        */     (UINT16)(sizeof(ReadPublic_In)),
    /* outSize       */     (UINT16)(sizeof(ReadPublic_Out)),
    /* offsetOfTypes */     offsetof(ReadPublic_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ReadPublic_Out, name)),
                             (UINT16)(offsetof(ReadPublic_Out, qualifiedName))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PUBLIC_P_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _ReadPublicDataAddress (&_ReadPublicData)
#else
#define _ReadPublicDataAddress 0
#endif // CC_ReadPublic

#if CC_ActivateCredential

#include "ActivateCredential_fp.h"

typedef TPM_RC  (ActivateCredential_Entry)(
    ActivateCredential_In           *in,
    ActivateCredential_Out          *out
);

typedef const struct {
    ActivateCredential_Entry    *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[3];
    BYTE                        types[7];
} ActivateCredential_COMMAND_DESCRIPTOR_t;

ActivateCredential_COMMAND_DESCRIPTOR_t _ActivateCredentialData = {
    /* entry         */         &TPM2_ActivateCredential,
    /* inSize        */         (UINT16)(sizeof(ActivateCredential_In)),
    /* outSize       */         (UINT16)(sizeof(ActivateCredential_Out)),
    /* offsetOfTypes */         offsetof(ActivateCredential_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(ActivateCredential_In, keyHandle)),
                                 (UINT16)(offsetof(ActivateCredential_In, credentialBlob)),
                                 (UINT16)(offsetof(ActivateCredential_In, secret))},
    /* types         */         {TPMI_DH_OBJECT_H_UNMARSHAL,
                                 TPMI_DH_OBJECT_H_UNMARSHAL,
                                 TPM2B_ID_OBJECT_P_UNMARSHAL,
                                 TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL,
                                 END_OF_LIST,
                                 TPM2B_DIGEST_P_MARSHAL,
                                 END_OF_LIST}
};

#define _ActivateCredentialDataAddress (&_ActivateCredentialData)
#else
#define _ActivateCredentialDataAddress 0
#endif // CC_ActivateCredential

#if CC_MakeCredential

#include "MakeCredential_fp.h"

typedef TPM_RC  (MakeCredential_Entry)(
    MakeCredential_In           *in,
    MakeCredential_Out          *out
);

typedef const struct {
    MakeCredential_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} MakeCredential_COMMAND_DESCRIPTOR_t;

MakeCredential_COMMAND_DESCRIPTOR_t _MakeCredentialData = {
    /* entry         */     &TPM2_MakeCredential,
    /* inSize        */     (UINT16)(sizeof(MakeCredential_In)),
    /* outSize       */     (UINT16)(sizeof(MakeCredential_Out)),
    /* offsetOfTypes */     offsetof(MakeCredential_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(MakeCredential_In, credential)),
                             (UINT16)(offsetof(MakeCredential_In, objectName)),
                             (UINT16)(offsetof(MakeCredential_Out, secret))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPM2B_NAME_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ID_OBJECT_P_MARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_MARSHAL,
                             END_OF_LIST}
};

#define _MakeCredentialDataAddress (&_MakeCredentialData)
#else
#define _MakeCredentialDataAddress 0
#endif // CC_MakeCredential

#if CC_Unseal

#include "Unseal_fp.h"

typedef TPM_RC  (Unseal_Entry)(
    Unseal_In                   *in,
    Unseal_Out                  *out
);

typedef const struct {
    Unseal_Entry            *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} Unseal_COMMAND_DESCRIPTOR_t;

Unseal_COMMAND_DESCRIPTOR_t _UnsealData = {
    /* entry         */     &TPM2_Unseal,
    /* inSize        */     (UINT16)(sizeof(Unseal_In)),
    /* outSize       */     (UINT16)(sizeof(Unseal_Out)),
    /* offsetOfTypes */     offsetof(Unseal_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_SENSITIVE_DATA_P_MARSHAL,
                             END_OF_LIST}
};

#define _UnsealDataAddress (&_UnsealData)
#else
#define _UnsealDataAddress 0
#endif // CC_Unseal

#if CC_ObjectChangeAuth

#include "ObjectChangeAuth_fp.h"

typedef TPM_RC  (ObjectChangeAuth_Entry)(
    ObjectChangeAuth_In         *in,
    ObjectChangeAuth_Out        *out
);

typedef const struct {
    ObjectChangeAuth_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} ObjectChangeAuth_COMMAND_DESCRIPTOR_t;

ObjectChangeAuth_COMMAND_DESCRIPTOR_t _ObjectChangeAuthData = {
    /* entry         */     &TPM2_ObjectChangeAuth,
    /* inSize        */     (UINT16)(sizeof(ObjectChangeAuth_In)),
    /* outSize       */     (UINT16)(sizeof(ObjectChangeAuth_Out)),
    /* offsetOfTypes */     offsetof(ObjectChangeAuth_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ObjectChangeAuth_In, parentHandle)),
                             (UINT16)(offsetof(ObjectChangeAuth_In, newAuth))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_AUTH_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PRIVATE_P_MARSHAL,
                             END_OF_LIST}
};

#define _ObjectChangeAuthDataAddress (&_ObjectChangeAuthData)
#else
#define _ObjectChangeAuthDataAddress 0
#endif // CC_ObjectChangeAuth

#if CC_CreateLoaded

#include "CreateLoaded_fp.h"

typedef TPM_RC  (CreateLoaded_Entry)(
    CreateLoaded_In             *in,
    CreateLoaded_Out            *out
);

typedef const struct {
    CreateLoaded_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} CreateLoaded_COMMAND_DESCRIPTOR_t;

CreateLoaded_COMMAND_DESCRIPTOR_t _CreateLoadedData = {
    /* entry         */     &TPM2_CreateLoaded,
    /* inSize        */     (UINT16)(sizeof(CreateLoaded_In)),
    /* outSize       */     (UINT16)(sizeof(CreateLoaded_Out)),
    /* offsetOfTypes */     offsetof(CreateLoaded_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(CreateLoaded_In, inSensitive)),
                             (UINT16)(offsetof(CreateLoaded_In, inPublic)),
                             (UINT16)(offsetof(CreateLoaded_Out, outPrivate)),
                             (UINT16)(offsetof(CreateLoaded_Out, outPublic)),
                             (UINT16)(offsetof(CreateLoaded_Out, name))},
    /* types         */     {TPMI_DH_PARENT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_SENSITIVE_CREATE_P_UNMARSHAL,
                             TPM2B_TEMPLATE_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM_HANDLE_H_MARSHAL,
                             TPM2B_PRIVATE_P_MARSHAL,
                             TPM2B_PUBLIC_P_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _CreateLoadedDataAddress (&_CreateLoadedData)
#else
#define _CreateLoadedDataAddress 0
#endif // CC_CreateLoaded

#if CC_Duplicate

#include "Duplicate_fp.h"

typedef TPM_RC  (Duplicate_Entry)(
    Duplicate_In                *in,
    Duplicate_Out               *out
);

typedef const struct {
    Duplicate_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} Duplicate_COMMAND_DESCRIPTOR_t;

Duplicate_COMMAND_DESCRIPTOR_t _DuplicateData = {
    /* entry         */     &TPM2_Duplicate,
    /* inSize        */     (UINT16)(sizeof(Duplicate_In)),
    /* outSize       */     (UINT16)(sizeof(Duplicate_Out)),
    /* offsetOfTypes */     offsetof(Duplicate_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Duplicate_In, newParentHandle)),
                             (UINT16)(offsetof(Duplicate_In, encryptionKeyIn)),
                             (UINT16)(offsetof(Duplicate_In, symmetricAlg)),
                             (UINT16)(offsetof(Duplicate_Out, duplicate)),
                             (UINT16)(offsetof(Duplicate_Out, outSymSeed))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SYM_DEF_OBJECT_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_DATA_P_MARSHAL,
                             TPM2B_PRIVATE_P_MARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_MARSHAL,
                             END_OF_LIST}
};

#define _DuplicateDataAddress (&_DuplicateData)
#else
#define _DuplicateDataAddress 0
#endif // CC_Duplicate

#if CC_Rewrap

#include "Rewrap_fp.h"

typedef TPM_RC  (Rewrap_Entry)(
    Rewrap_In                   *in,
    Rewrap_Out                  *out
);

typedef const struct {
    Rewrap_Entry            *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} Rewrap_COMMAND_DESCRIPTOR_t;

Rewrap_COMMAND_DESCRIPTOR_t _RewrapData = {
    /* entry         */     &TPM2_Rewrap,
    /* inSize        */     (UINT16)(sizeof(Rewrap_In)),
    /* outSize       */     (UINT16)(sizeof(Rewrap_Out)),
    /* offsetOfTypes */     offsetof(Rewrap_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Rewrap_In, newParent)),
                             (UINT16)(offsetof(Rewrap_In, inDuplicate)),
                             (UINT16)(offsetof(Rewrap_In, name)),
                             (UINT16)(offsetof(Rewrap_In, inSymSeed)),
                             (UINT16)(offsetof(Rewrap_Out, outSymSeed))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_PRIVATE_P_UNMARSHAL,
                             TPM2B_NAME_P_UNMARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PRIVATE_P_MARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_MARSHAL,
                             END_OF_LIST}
};

#define _RewrapDataAddress (&_RewrapData)
#else
#define _RewrapDataAddress 0
#endif // CC_Rewrap

#if CC_Import

#include "Import_fp.h"

typedef TPM_RC  (Import_Entry)(
    Import_In                   *in,
    Import_Out                  *out
);

typedef const struct {
    Import_Entry            *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} Import_COMMAND_DESCRIPTOR_t;

Import_COMMAND_DESCRIPTOR_t _ImportData = {
    /* entry         */     &TPM2_Import,
    /* inSize        */     (UINT16)(sizeof(Import_In)),
    /* outSize       */     (UINT16)(sizeof(Import_Out)),
    /* offsetOfTypes */     offsetof(Import_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Import_In, encryptionKey)),
                             (UINT16)(offsetof(Import_In, objectPublic)),
                             (UINT16)(offsetof(Import_In, duplicate)),
                             (UINT16)(offsetof(Import_In, inSymSeed)),
                             (UINT16)(offsetof(Import_In, symmetricAlg))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPM2B_PUBLIC_P_UNMARSHAL,
                             TPM2B_PRIVATE_P_UNMARSHAL,
                             TPM2B_ENCRYPTED_SECRET_P_UNMARSHAL,
                             TPMT_SYM_DEF_OBJECT_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_PRIVATE_P_MARSHAL,
                             END_OF_LIST}
};

#define _ImportDataAddress (&_ImportData)
#else
#define _ImportDataAddress 0
#endif // CC_Import

#if CC_RSA_Encrypt

#include "RSA_Encrypt_fp.h"

typedef TPM_RC  (RSA_Encrypt_Entry)(
    RSA_Encrypt_In              *in,
    RSA_Encrypt_Out             *out
);

typedef const struct {
    RSA_Encrypt_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} RSA_Encrypt_COMMAND_DESCRIPTOR_t;

RSA_Encrypt_COMMAND_DESCRIPTOR_t _RSA_EncryptData = {
    /* entry         */     &TPM2_RSA_Encrypt,
    /* inSize        */     (UINT16)(sizeof(RSA_Encrypt_In)),
    /* outSize       */     (UINT16)(sizeof(RSA_Encrypt_Out)),
    /* offsetOfTypes */     offsetof(RSA_Encrypt_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(RSA_Encrypt_In, message)),
                             (UINT16)(offsetof(RSA_Encrypt_In, inScheme)),
                             (UINT16)(offsetof(RSA_Encrypt_In, label))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_PUBLIC_KEY_RSA_P_UNMARSHAL,
                             TPMT_RSA_DECRYPT_P_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PUBLIC_KEY_RSA_P_MARSHAL,
                             END_OF_LIST}
};

#define _RSA_EncryptDataAddress (&_RSA_EncryptData)
#else
#define _RSA_EncryptDataAddress 0
#endif // CC_RSA_Encrypt

#if CC_RSA_Decrypt

#include "RSA_Decrypt_fp.h"

typedef TPM_RC  (RSA_Decrypt_Entry)(
    RSA_Decrypt_In              *in,
    RSA_Decrypt_Out             *out
);

typedef const struct {
    RSA_Decrypt_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} RSA_Decrypt_COMMAND_DESCRIPTOR_t;

RSA_Decrypt_COMMAND_DESCRIPTOR_t _RSA_DecryptData = {
    /* entry         */     &TPM2_RSA_Decrypt,
    /* inSize        */     (UINT16)(sizeof(RSA_Decrypt_In)),
    /* outSize       */     (UINT16)(sizeof(RSA_Decrypt_Out)),
    /* offsetOfTypes */     offsetof(RSA_Decrypt_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(RSA_Decrypt_In, cipherText)),
                             (UINT16)(offsetof(RSA_Decrypt_In, inScheme)),
                             (UINT16)(offsetof(RSA_Decrypt_In, label))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_PUBLIC_KEY_RSA_P_UNMARSHAL,
                             TPMT_RSA_DECRYPT_P_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_PUBLIC_KEY_RSA_P_MARSHAL,
                             END_OF_LIST}
};

#define _RSA_DecryptDataAddress (&_RSA_DecryptData)
#else
#define _RSA_DecryptDataAddress 0
#endif // CC_RSA_Decrypt

#if CC_ECDH_KeyGen

#include "ECDH_KeyGen_fp.h"

typedef TPM_RC  (ECDH_KeyGen_Entry)(
    ECDH_KeyGen_In              *in,
    ECDH_KeyGen_Out             *out
);

typedef const struct {
    ECDH_KeyGen_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} ECDH_KeyGen_COMMAND_DESCRIPTOR_t;

ECDH_KeyGen_COMMAND_DESCRIPTOR_t _ECDH_KeyGenData = {
    /* entry         */     &TPM2_ECDH_KeyGen,
    /* inSize        */     (UINT16)(sizeof(ECDH_KeyGen_In)),
    /* outSize       */     (UINT16)(sizeof(ECDH_KeyGen_Out)),
    /* offsetOfTypes */     offsetof(ECDH_KeyGen_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ECDH_KeyGen_Out, pubPoint))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             END_OF_LIST}
};

#define _ECDH_KeyGenDataAddress (&_ECDH_KeyGenData)
#else
#define _ECDH_KeyGenDataAddress 0
#endif // CC_ECDH_KeyGen

#if CC_ECDH_ZGen

#include "ECDH_ZGen_fp.h"

typedef TPM_RC  (ECDH_ZGen_Entry)(
    ECDH_ZGen_In                *in,
    ECDH_ZGen_Out               *out
);

typedef const struct {
    ECDH_ZGen_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} ECDH_ZGen_COMMAND_DESCRIPTOR_t;

ECDH_ZGen_COMMAND_DESCRIPTOR_t _ECDH_ZGenData = {
    /* entry         */     &TPM2_ECDH_ZGen,
    /* inSize        */     (UINT16)(sizeof(ECDH_ZGen_In)),
    /* outSize       */     (UINT16)(sizeof(ECDH_ZGen_Out)),
    /* offsetOfTypes */     offsetof(ECDH_ZGen_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ECDH_ZGen_In, inPoint))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_ECC_POINT_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             END_OF_LIST}
};

#define _ECDH_ZGenDataAddress (&_ECDH_ZGenData)
#else
#define _ECDH_ZGenDataAddress 0
#endif // CC_ECDH_ZGen

#if CC_ECC_Parameters

#include "ECC_Parameters_fp.h"

typedef TPM_RC  (ECC_Parameters_Entry)(
    ECC_Parameters_In           *in,
    ECC_Parameters_Out          *out
);

typedef const struct {
    ECC_Parameters_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} ECC_Parameters_COMMAND_DESCRIPTOR_t;

ECC_Parameters_COMMAND_DESCRIPTOR_t _ECC_ParametersData = {
    /* entry         */     &TPM2_ECC_Parameters,
    /* inSize        */     (UINT16)(sizeof(ECC_Parameters_In)),
    /* outSize       */     (UINT16)(sizeof(ECC_Parameters_Out)),
    /* offsetOfTypes */     offsetof(ECC_Parameters_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_ECC_CURVE_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMS_ALGORITHM_DETAIL_ECC_P_MARSHAL,
                             END_OF_LIST}
};

#define _ECC_ParametersDataAddress (&_ECC_ParametersData)
#else
#define _ECC_ParametersDataAddress 0
#endif // CC_ECC_Parameters

#if CC_ZGen_2Phase

#include "ZGen_2Phase_fp.h"

typedef TPM_RC  (ZGen_2Phase_Entry)(
    ZGen_2Phase_In              *in,
    ZGen_2Phase_Out             *out
);

typedef const struct {
    ZGen_2Phase_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} ZGen_2Phase_COMMAND_DESCRIPTOR_t;

ZGen_2Phase_COMMAND_DESCRIPTOR_t _ZGen_2PhaseData = {
    /* entry         */     &TPM2_ZGen_2Phase,
    /* inSize        */     (UINT16)(sizeof(ZGen_2Phase_In)),
    /* outSize       */     (UINT16)(sizeof(ZGen_2Phase_Out)),
    /* offsetOfTypes */     offsetof(ZGen_2Phase_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ZGen_2Phase_In, inQsB)),
                             (UINT16)(offsetof(ZGen_2Phase_In, inQeB)),
                             (UINT16)(offsetof(ZGen_2Phase_In, inScheme)),
                             (UINT16)(offsetof(ZGen_2Phase_In, counter)),
                             (UINT16)(offsetof(ZGen_2Phase_Out, outZ2))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_ECC_POINT_P_UNMARSHAL,
                             TPM2B_ECC_POINT_P_UNMARSHAL,
                             TPMI_ECC_KEY_EXCHANGE_P_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             END_OF_LIST}
};

#define _ZGen_2PhaseDataAddress (&_ZGen_2PhaseData)
#else
#define _ZGen_2PhaseDataAddress 0
#endif // CC_ZGen_2Phase

#if CC_EncryptDecrypt

#include "EncryptDecrypt_fp.h"

typedef TPM_RC  (EncryptDecrypt_Entry)(
    EncryptDecrypt_In           *in,
    EncryptDecrypt_Out          *out
);

typedef const struct {
    EncryptDecrypt_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} EncryptDecrypt_COMMAND_DESCRIPTOR_t;

EncryptDecrypt_COMMAND_DESCRIPTOR_t _EncryptDecryptData = {
    /* entry         */     &TPM2_EncryptDecrypt,
    /* inSize        */     (UINT16)(sizeof(EncryptDecrypt_In)),
    /* outSize       */     (UINT16)(sizeof(EncryptDecrypt_Out)),
    /* offsetOfTypes */     offsetof(EncryptDecrypt_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(EncryptDecrypt_In, decrypt)),
                             (UINT16)(offsetof(EncryptDecrypt_In, mode)),
                             (UINT16)(offsetof(EncryptDecrypt_In, ivIn)),
                             (UINT16)(offsetof(EncryptDecrypt_In, inData)),
                             (UINT16)(offsetof(EncryptDecrypt_Out, ivOut))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_YES_NO_P_UNMARSHAL,
                             TPMI_ALG_CIPHER_MODE_P_UNMARSHAL + ADD_FLAG,
                             TPM2B_IV_P_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_MAX_BUFFER_P_MARSHAL,
                             TPM2B_IV_P_MARSHAL,
                             END_OF_LIST}
};

#define _EncryptDecryptDataAddress (&_EncryptDecryptData)
#else
#define _EncryptDecryptDataAddress 0
#endif // CC_EncryptDecrypt

#if CC_EncryptDecrypt2

#include "EncryptDecrypt2_fp.h"

typedef TPM_RC  (EncryptDecrypt2_Entry)(
    EncryptDecrypt2_In          *in,
    EncryptDecrypt2_Out         *out
);

typedef const struct {
    EncryptDecrypt2_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[9];
} EncryptDecrypt2_COMMAND_DESCRIPTOR_t;

EncryptDecrypt2_COMMAND_DESCRIPTOR_t _EncryptDecrypt2Data = {
    /* entry         */     &TPM2_EncryptDecrypt2,
    /* inSize        */     (UINT16)(sizeof(EncryptDecrypt2_In)),
    /* outSize       */     (UINT16)(sizeof(EncryptDecrypt2_Out)),
    /* offsetOfTypes */     offsetof(EncryptDecrypt2_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(EncryptDecrypt2_In, inData)),
                             (UINT16)(offsetof(EncryptDecrypt2_In, decrypt)),
                             (UINT16)(offsetof(EncryptDecrypt2_In, mode)),
                             (UINT16)(offsetof(EncryptDecrypt2_In, ivIn)),
                             (UINT16)(offsetof(EncryptDecrypt2_Out, ivOut))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             TPMI_YES_NO_P_UNMARSHAL,
                             TPMI_ALG_CIPHER_MODE_P_UNMARSHAL + ADD_FLAG,
                             TPM2B_IV_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_MAX_BUFFER_P_MARSHAL,
                             TPM2B_IV_P_MARSHAL,
                             END_OF_LIST}
};

#define _EncryptDecrypt2DataAddress (&_EncryptDecrypt2Data)
#else
#define _EncryptDecrypt2DataAddress 0
#endif // CC_EncryptDecrypt2

#if CC_Hash

#include "Hash_fp.h"

typedef TPM_RC  (Hash_Entry)(
    Hash_In                     *in,
    Hash_Out                    *out
);

typedef const struct {
    Hash_Entry              *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} Hash_COMMAND_DESCRIPTOR_t;

Hash_COMMAND_DESCRIPTOR_t _HashData = {
    /* entry         */     &TPM2_Hash,
    /* inSize        */     (UINT16)(sizeof(Hash_In)),
    /* outSize       */     (UINT16)(sizeof(Hash_Out)),
    /* offsetOfTypes */     offsetof(Hash_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Hash_In, hashAlg)),
                             (UINT16)(offsetof(Hash_In, hierarchy)),
                             (UINT16)(offsetof(Hash_Out, validation))},
    /* types         */     {TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             TPMI_ALG_HASH_P_UNMARSHAL,
                             TPMI_RH_HIERARCHY_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             TPMT_TK_HASHCHECK_P_MARSHAL,
                             END_OF_LIST}
};

#define _HashDataAddress (&_HashData)
#else
#define _HashDataAddress 0
#endif // CC_Hash

#if CC_HMAC

#include "HMAC_fp.h"

typedef TPM_RC  (HMAC_Entry)(
    HMAC_In                     *in,
    HMAC_Out                    *out
);

typedef const struct {
    HMAC_Entry              *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} HMAC_COMMAND_DESCRIPTOR_t;

HMAC_COMMAND_DESCRIPTOR_t _HMACData = {
    /* entry         */     &TPM2_HMAC,
    /* inSize        */     (UINT16)(sizeof(HMAC_In)),
    /* outSize       */     (UINT16)(sizeof(HMAC_Out)),
    /* offsetOfTypes */     offsetof(HMAC_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(HMAC_In, buffer)),
                             (UINT16)(offsetof(HMAC_In, hashAlg))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             END_OF_LIST}
};

#define _HMACDataAddress (&_HMACData)
#else
#define _HMACDataAddress 0
#endif // CC_HMAC

#if CC_MAC

#include "MAC_fp.h"

typedef TPM_RC  (MAC_Entry)(
    MAC_In                      *in,
    MAC_Out                     *out
);

typedef const struct {
    MAC_Entry               *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} MAC_COMMAND_DESCRIPTOR_t;

MAC_COMMAND_DESCRIPTOR_t _MACData = {
    /* entry         */     &TPM2_MAC,
    /* inSize        */     (UINT16)(sizeof(MAC_In)),
    /* outSize       */     (UINT16)(sizeof(MAC_Out)),
    /* offsetOfTypes */     offsetof(MAC_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(MAC_In, buffer)),
                             (UINT16)(offsetof(MAC_In, inScheme))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             TPMI_ALG_MAC_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             END_OF_LIST}
};

#define _MACDataAddress (&_MACData)
#else
#define _MACDataAddress 0
#endif // CC_MAC

#if CC_GetRandom

#include "GetRandom_fp.h"

typedef TPM_RC  (GetRandom_Entry)(
    GetRandom_In                *in,
    GetRandom_Out               *out
);

typedef const struct {
    GetRandom_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} GetRandom_COMMAND_DESCRIPTOR_t;

GetRandom_COMMAND_DESCRIPTOR_t _GetRandomData = {
    /* entry         */     &TPM2_GetRandom,
    /* inSize        */     (UINT16)(sizeof(GetRandom_In)),
    /* outSize       */     (UINT16)(sizeof(GetRandom_Out)),
    /* offsetOfTypes */     offsetof(GetRandom_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {UINT16_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             END_OF_LIST}
};

#define _GetRandomDataAddress (&_GetRandomData)
#else
#define _GetRandomDataAddress 0
#endif // CC_GetRandom

#if CC_StirRandom

#include "StirRandom_fp.h"

typedef TPM_RC  (StirRandom_Entry)(
    StirRandom_In               *in
);

typedef const struct {
    StirRandom_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} StirRandom_COMMAND_DESCRIPTOR_t;

StirRandom_COMMAND_DESCRIPTOR_t _StirRandomData = {
    /* entry         */     &TPM2_StirRandom,
    /* inSize        */     (UINT16)(sizeof(StirRandom_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(StirRandom_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPM2B_SENSITIVE_DATA_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _StirRandomDataAddress (&_StirRandomData)
#else
#define _StirRandomDataAddress 0
#endif // CC_StirRandom

#if CC_HMAC_Start

#include "HMAC_Start_fp.h"

typedef TPM_RC  (HMAC_Start_Entry)(
    HMAC_Start_In               *in,
    HMAC_Start_Out              *out
);

typedef const struct {
    HMAC_Start_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} HMAC_Start_COMMAND_DESCRIPTOR_t;

HMAC_Start_COMMAND_DESCRIPTOR_t _HMAC_StartData = {
    /* entry         */     &TPM2_HMAC_Start,
    /* inSize        */     (UINT16)(sizeof(HMAC_Start_In)),
    /* outSize       */     (UINT16)(sizeof(HMAC_Start_Out)),
    /* offsetOfTypes */     offsetof(HMAC_Start_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(HMAC_Start_In, auth)),
                             (UINT16)(offsetof(HMAC_Start_In, hashAlg))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_AUTH_P_UNMARSHAL,
                             TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPMI_DH_OBJECT_H_MARSHAL,
                             END_OF_LIST}
};

#define _HMAC_StartDataAddress (&_HMAC_StartData)
#else
#define _HMAC_StartDataAddress 0
#endif // CC_HMAC_Start

#if CC_MAC_Start

#include "MAC_Start_fp.h"

typedef TPM_RC  (MAC_Start_Entry)(
    MAC_Start_In                *in,
    MAC_Start_Out               *out
);

typedef const struct {
    MAC_Start_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} MAC_Start_COMMAND_DESCRIPTOR_t;

MAC_Start_COMMAND_DESCRIPTOR_t _MAC_StartData = {
    /* entry         */     &TPM2_MAC_Start,
    /* inSize        */     (UINT16)(sizeof(MAC_Start_In)),
    /* outSize       */     (UINT16)(sizeof(MAC_Start_Out)),
    /* offsetOfTypes */     offsetof(MAC_Start_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(MAC_Start_In, auth)),
                             (UINT16)(offsetof(MAC_Start_In, inScheme))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_AUTH_P_UNMARSHAL,
                             TPMI_ALG_MAC_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPMI_DH_OBJECT_H_MARSHAL,
                             END_OF_LIST}
};

#define _MAC_StartDataAddress (&_MAC_StartData)
#else
#define _MAC_StartDataAddress 0
#endif // CC_MAC_Start

#if CC_HashSequenceStart

#include "HashSequenceStart_fp.h"

typedef TPM_RC  (HashSequenceStart_Entry)(
    HashSequenceStart_In            *in,
    HashSequenceStart_Out           *out
);

typedef const struct {
    HashSequenceStart_Entry     *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[1];
    BYTE                        types[5];
} HashSequenceStart_COMMAND_DESCRIPTOR_t;

HashSequenceStart_COMMAND_DESCRIPTOR_t _HashSequenceStartData = {
    /* entry         */         &TPM2_HashSequenceStart,
    /* inSize        */         (UINT16)(sizeof(HashSequenceStart_In)),
    /* outSize       */         (UINT16)(sizeof(HashSequenceStart_Out)),
    /* offsetOfTypes */         offsetof(HashSequenceStart_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(HashSequenceStart_In, hashAlg))},
    /* types         */         {TPM2B_AUTH_P_UNMARSHAL,
                                 TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                                 END_OF_LIST,
                                 TPMI_DH_OBJECT_H_MARSHAL,
                                 END_OF_LIST}
};

#define _HashSequenceStartDataAddress (&_HashSequenceStartData)
#else
#define _HashSequenceStartDataAddress 0
#endif // CC_HashSequenceStart

#if CC_SequenceUpdate

#include "SequenceUpdate_fp.h"

typedef TPM_RC  (SequenceUpdate_Entry)(
    SequenceUpdate_In           *in
);

typedef const struct {
    SequenceUpdate_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} SequenceUpdate_COMMAND_DESCRIPTOR_t;

SequenceUpdate_COMMAND_DESCRIPTOR_t _SequenceUpdateData = {
    /* entry         */     &TPM2_SequenceUpdate,
    /* inSize        */     (UINT16)(sizeof(SequenceUpdate_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(SequenceUpdate_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(SequenceUpdate_In, buffer))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _SequenceUpdateDataAddress (&_SequenceUpdateData)
#else
#define _SequenceUpdateDataAddress 0
#endif // CC_SequenceUpdate

#if CC_SequenceComplete

#include "SequenceComplete_fp.h"

typedef TPM_RC  (SequenceComplete_Entry)(
    SequenceComplete_In         *in,
    SequenceComplete_Out        *out
);

typedef const struct {
    SequenceComplete_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} SequenceComplete_COMMAND_DESCRIPTOR_t;

SequenceComplete_COMMAND_DESCRIPTOR_t _SequenceCompleteData = {
    /* entry         */     &TPM2_SequenceComplete,
    /* inSize        */     (UINT16)(sizeof(SequenceComplete_In)),
    /* outSize       */     (UINT16)(sizeof(SequenceComplete_Out)),
    /* offsetOfTypes */     offsetof(SequenceComplete_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(SequenceComplete_In, buffer)),
                             (UINT16)(offsetof(SequenceComplete_In, hierarchy)),
                             (UINT16)(offsetof(SequenceComplete_Out, validation))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             TPMI_RH_HIERARCHY_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             TPMT_TK_HASHCHECK_P_MARSHAL,
                             END_OF_LIST}
};

#define _SequenceCompleteDataAddress (&_SequenceCompleteData)
#else
#define _SequenceCompleteDataAddress 0
#endif // CC_SequenceComplete

#if CC_EventSequenceComplete

#include "EventSequenceComplete_fp.h"

typedef TPM_RC  (EventSequenceComplete_Entry)(
    EventSequenceComplete_In            *in,
    EventSequenceComplete_Out           *out
);

typedef const struct {
    EventSequenceComplete_Entry     *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    UINT16                          paramOffsets[2];
    BYTE                            types[6];
} EventSequenceComplete_COMMAND_DESCRIPTOR_t;

EventSequenceComplete_COMMAND_DESCRIPTOR_t _EventSequenceCompleteData = {
    /* entry         */             &TPM2_EventSequenceComplete,
    /* inSize        */             (UINT16)(sizeof(EventSequenceComplete_In)),
    /* outSize       */             (UINT16)(sizeof(EventSequenceComplete_Out)),
    /* offsetOfTypes */             offsetof(EventSequenceComplete_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             {(UINT16)(offsetof(EventSequenceComplete_In, sequenceHandle)),
                                     (UINT16)(offsetof(EventSequenceComplete_In, buffer))},
    /* types         */             {TPMI_DH_PCR_H_UNMARSHAL + ADD_FLAG,
                                     TPMI_DH_OBJECT_H_UNMARSHAL,
                                     TPM2B_MAX_BUFFER_P_UNMARSHAL,
                                     END_OF_LIST,
                                     TPML_DIGEST_VALUES_P_MARSHAL,
                                     END_OF_LIST}
};

#define _EventSequenceCompleteDataAddress (&_EventSequenceCompleteData)
#else
#define _EventSequenceCompleteDataAddress 0
#endif // CC_EventSequenceComplete

#if CC_Certify

#include "Certify_fp.h"

typedef TPM_RC  (Certify_Entry)(
    Certify_In                  *in,
    Certify_Out                 *out
);

typedef const struct {
    Certify_Entry           *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[4];
    BYTE                    types[8];
} Certify_COMMAND_DESCRIPTOR_t;

Certify_COMMAND_DESCRIPTOR_t _CertifyData = {
    /* entry         */     &TPM2_Certify,
    /* inSize        */     (UINT16)(sizeof(Certify_In)),
    /* outSize       */     (UINT16)(sizeof(Certify_Out)),
    /* offsetOfTypes */     offsetof(Certify_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Certify_In, signHandle)),
                             (UINT16)(offsetof(Certify_In, qualifyingData)),
                             (UINT16)(offsetof(Certify_In, inScheme)),
                             (UINT16)(offsetof(Certify_Out, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_ATTEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _CertifyDataAddress (&_CertifyData)
#else
#define _CertifyDataAddress 0
#endif // CC_Certify

#if CC_CertifyCreation

#include "CertifyCreation_fp.h"

typedef TPM_RC  (CertifyCreation_Entry)(
    CertifyCreation_In          *in,
    CertifyCreation_Out         *out
);

typedef const struct {
    CertifyCreation_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[6];
    BYTE                    types[10];
} CertifyCreation_COMMAND_DESCRIPTOR_t;

CertifyCreation_COMMAND_DESCRIPTOR_t _CertifyCreationData = {
    /* entry         */     &TPM2_CertifyCreation,
    /* inSize        */     (UINT16)(sizeof(CertifyCreation_In)),
    /* outSize       */     (UINT16)(sizeof(CertifyCreation_Out)),
    /* offsetOfTypes */     offsetof(CertifyCreation_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(CertifyCreation_In, objectHandle)),
                             (UINT16)(offsetof(CertifyCreation_In, qualifyingData)),
                             (UINT16)(offsetof(CertifyCreation_In, creationHash)),
                             (UINT16)(offsetof(CertifyCreation_In, inScheme)),
                             (UINT16)(offsetof(CertifyCreation_In, creationTicket)),
                             (UINT16)(offsetof(CertifyCreation_Out, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             TPMT_TK_CREATION_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ATTEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _CertifyCreationDataAddress (&_CertifyCreationData)
#else
#define _CertifyCreationDataAddress 0
#endif // CC_CertifyCreation

#if CC_Quote

#include "Quote_fp.h"

typedef TPM_RC  (Quote_Entry)(
    Quote_In                    *in,
    Quote_Out                   *out
);

typedef const struct {
    Quote_Entry             *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[4];
    BYTE                    types[8];
} Quote_COMMAND_DESCRIPTOR_t;

Quote_COMMAND_DESCRIPTOR_t _QuoteData = {
    /* entry         */     &TPM2_Quote,
    /* inSize        */     (UINT16)(sizeof(Quote_In)),
    /* outSize       */     (UINT16)(sizeof(Quote_Out)),
    /* offsetOfTypes */     offsetof(Quote_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Quote_In, qualifyingData)),
                             (UINT16)(offsetof(Quote_In, inScheme)),
                             (UINT16)(offsetof(Quote_In, PCRselect)),
                             (UINT16)(offsetof(Quote_Out, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ATTEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _QuoteDataAddress (&_QuoteData)
#else
#define _QuoteDataAddress 0
#endif // CC_Quote

#if CC_GetSessionAuditDigest

#include "GetSessionAuditDigest_fp.h"

typedef TPM_RC  (GetSessionAuditDigest_Entry)(
    GetSessionAuditDigest_In            *in,
    GetSessionAuditDigest_Out           *out
);

typedef const struct {
    GetSessionAuditDigest_Entry     *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    UINT16                          paramOffsets[5];
    BYTE                            types[9];
} GetSessionAuditDigest_COMMAND_DESCRIPTOR_t;

GetSessionAuditDigest_COMMAND_DESCRIPTOR_t _GetSessionAuditDigestData = {
    /* entry         */             &TPM2_GetSessionAuditDigest,
    /* inSize        */             (UINT16)(sizeof(GetSessionAuditDigest_In)),
    /* outSize       */             (UINT16)(sizeof(GetSessionAuditDigest_Out)),
    /* offsetOfTypes */             offsetof(GetSessionAuditDigest_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             {(UINT16)(offsetof(GetSessionAuditDigest_In, signHandle)),
                                     (UINT16)(offsetof(GetSessionAuditDigest_In, sessionHandle)),
                                     (UINT16)(offsetof(GetSessionAuditDigest_In, qualifyingData)),
                                     (UINT16)(offsetof(GetSessionAuditDigest_In, inScheme)),
                                     (UINT16)(offsetof(GetSessionAuditDigest_Out, signature))},
    /* types         */             {TPMI_RH_ENDORSEMENT_H_UNMARSHAL,
                                     TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                                     TPMI_SH_HMAC_H_UNMARSHAL,
                                     TPM2B_DATA_P_UNMARSHAL,
                                     TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                                     END_OF_LIST,
                                     TPM2B_ATTEST_P_MARSHAL,
                                     TPMT_SIGNATURE_P_MARSHAL,
                                     END_OF_LIST}
};

#define _GetSessionAuditDigestDataAddress (&_GetSessionAuditDigestData)
#else
#define _GetSessionAuditDigestDataAddress 0
#endif // CC_GetSessionAuditDigest

#if CC_GetCommandAuditDigest

#include "GetCommandAuditDigest_fp.h"

typedef TPM_RC  (GetCommandAuditDigest_Entry)(
    GetCommandAuditDigest_In            *in,
    GetCommandAuditDigest_Out           *out
);

typedef const struct {
    GetCommandAuditDigest_Entry     *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    UINT16                          paramOffsets[4];
    BYTE                            types[8];
} GetCommandAuditDigest_COMMAND_DESCRIPTOR_t;

GetCommandAuditDigest_COMMAND_DESCRIPTOR_t _GetCommandAuditDigestData = {
    /* entry         */             &TPM2_GetCommandAuditDigest,
    /* inSize        */             (UINT16)(sizeof(GetCommandAuditDigest_In)),
    /* outSize       */             (UINT16)(sizeof(GetCommandAuditDigest_Out)),
    /* offsetOfTypes */             offsetof(GetCommandAuditDigest_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             {(UINT16)(offsetof(GetCommandAuditDigest_In, signHandle)),
                                     (UINT16)(offsetof(GetCommandAuditDigest_In, qualifyingData)),
                                     (UINT16)(offsetof(GetCommandAuditDigest_In, inScheme)),
                                     (UINT16)(offsetof(GetCommandAuditDigest_Out, signature))},
    /* types         */             {TPMI_RH_ENDORSEMENT_H_UNMARSHAL,
                                     TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                                     TPM2B_DATA_P_UNMARSHAL,
                                     TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                                     END_OF_LIST,
                                     TPM2B_ATTEST_P_MARSHAL,
                                     TPMT_SIGNATURE_P_MARSHAL,
                                     END_OF_LIST}
};

#define _GetCommandAuditDigestDataAddress (&_GetCommandAuditDigestData)
#else
#define _GetCommandAuditDigestDataAddress 0
#endif // CC_GetCommandAuditDigest

#if CC_GetTime

#include "GetTime_fp.h"

typedef TPM_RC  (GetTime_Entry)(
    GetTime_In                  *in,
    GetTime_Out                 *out
);

typedef const struct {
    GetTime_Entry           *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[4];
    BYTE                    types[8];
} GetTime_COMMAND_DESCRIPTOR_t;

GetTime_COMMAND_DESCRIPTOR_t _GetTimeData = {
    /* entry         */     &TPM2_GetTime,
    /* inSize        */     (UINT16)(sizeof(GetTime_In)),
    /* outSize       */     (UINT16)(sizeof(GetTime_Out)),
    /* offsetOfTypes */     offsetof(GetTime_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(GetTime_In, signHandle)),
                             (UINT16)(offsetof(GetTime_In, qualifyingData)),
                             (UINT16)(offsetof(GetTime_In, inScheme)),
                             (UINT16)(offsetof(GetTime_Out, signature))},
    /* types         */     {TPMI_RH_ENDORSEMENT_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             TPM2B_ATTEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _GetTimeDataAddress (&_GetTimeData)
#else
#define _GetTimeDataAddress 0
#endif // CC_GetTime

#if CC_CertifyX509

#include "CertifyX509_fp.h"

typedef TPM_RC  (CertifyX509_Entry)(
    CertifyX509_In              *in,
    CertifyX509_Out             *out
);

typedef const struct {
    CertifyX509_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[6];
    BYTE                    types[10];
} CertifyX509_COMMAND_DESCRIPTOR_t;

CertifyX509_COMMAND_DESCRIPTOR_t _CertifyX509Data = {
    /* entry         */     &TPM2_CertifyX509,
    /* inSize        */     (UINT16)(sizeof(CertifyX509_In)),
    /* outSize       */     (UINT16)(sizeof(CertifyX509_Out)),
    /* offsetOfTypes */     offsetof(CertifyX509_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(CertifyX509_In, signHandle)),
                             (UINT16)(offsetof(CertifyX509_In, qualifyingData)),
                             (UINT16)(offsetof(CertifyX509_In, inScheme)),
                             (UINT16)(offsetof(CertifyX509_In, partialCertificate)),
                             (UINT16)(offsetof(CertifyX509_Out, tbsDigest)),
                             (UINT16)(offsetof(CertifyX509_Out, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_MAX_BUFFER_P_MARSHAL,
                             TPM2B_DIGEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _CertifyX509DataAddress (&_CertifyX509Data)
#else
#define _CertifyX509DataAddress 0
#endif // CC_CertifyX509

#if CC_Commit

#include "Commit_fp.h"

typedef TPM_RC  (Commit_Entry)(
    Commit_In                   *in,
    Commit_Out                  *out
);

typedef const struct {
    Commit_Entry            *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[6];
    BYTE                    types[10];
} Commit_COMMAND_DESCRIPTOR_t;

Commit_COMMAND_DESCRIPTOR_t _CommitData = {
    /* entry         */     &TPM2_Commit,
    /* inSize        */     (UINT16)(sizeof(Commit_In)),
    /* outSize       */     (UINT16)(sizeof(Commit_Out)),
    /* offsetOfTypes */     offsetof(Commit_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Commit_In, P1)),
                             (UINT16)(offsetof(Commit_In, s2)),
                             (UINT16)(offsetof(Commit_In, y2)),
                             (UINT16)(offsetof(Commit_Out, L)),
                             (UINT16)(offsetof(Commit_Out, E)),
                             (UINT16)(offsetof(Commit_Out, counter))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_ECC_POINT_P_UNMARSHAL,
                             TPM2B_SENSITIVE_DATA_P_UNMARSHAL,
                             TPM2B_ECC_PARAMETER_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             UINT16_P_MARSHAL,
                             END_OF_LIST}
};

#define _CommitDataAddress (&_CommitData)
#else
#define _CommitDataAddress 0
#endif // CC_Commit

#if CC_EC_Ephemeral

#include "EC_Ephemeral_fp.h"

typedef TPM_RC  (EC_Ephemeral_Entry)(
    EC_Ephemeral_In             *in,
    EC_Ephemeral_Out            *out
);

typedef const struct {
    EC_Ephemeral_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} EC_Ephemeral_COMMAND_DESCRIPTOR_t;

EC_Ephemeral_COMMAND_DESCRIPTOR_t _EC_EphemeralData = {
    /* entry         */     &TPM2_EC_Ephemeral,
    /* inSize        */     (UINT16)(sizeof(EC_Ephemeral_In)),
    /* outSize       */     (UINT16)(sizeof(EC_Ephemeral_Out)),
    /* offsetOfTypes */     offsetof(EC_Ephemeral_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(EC_Ephemeral_Out, counter))},
    /* types         */     {TPMI_ECC_CURVE_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ECC_POINT_P_MARSHAL,
                             UINT16_P_MARSHAL,
                             END_OF_LIST}
};

#define _EC_EphemeralDataAddress (&_EC_EphemeralData)
#else
#define _EC_EphemeralDataAddress 0
#endif // CC_EC_Ephemeral

#if CC_VerifySignature

#include "VerifySignature_fp.h"

typedef TPM_RC  (VerifySignature_Entry)(
    VerifySignature_In          *in,
    VerifySignature_Out         *out
);

typedef const struct {
    VerifySignature_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} VerifySignature_COMMAND_DESCRIPTOR_t;

VerifySignature_COMMAND_DESCRIPTOR_t _VerifySignatureData = {
    /* entry         */     &TPM2_VerifySignature,
    /* inSize        */     (UINT16)(sizeof(VerifySignature_In)),
    /* outSize       */     (UINT16)(sizeof(VerifySignature_Out)),
    /* offsetOfTypes */     offsetof(VerifySignature_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(VerifySignature_In, digest)),
                             (UINT16)(offsetof(VerifySignature_In, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPMT_SIGNATURE_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMT_TK_VERIFIED_P_MARSHAL,
                             END_OF_LIST}
};

#define _VerifySignatureDataAddress (&_VerifySignatureData)
#else
#define _VerifySignatureDataAddress 0
#endif // CC_VerifySignature

#if CC_Sign

#include "Sign_fp.h"

typedef TPM_RC  (Sign_Entry)(
    Sign_In                     *in,
    Sign_Out                    *out
);

typedef const struct {
    Sign_Entry              *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} Sign_COMMAND_DESCRIPTOR_t;

Sign_COMMAND_DESCRIPTOR_t _SignData = {
    /* entry         */     &TPM2_Sign,
    /* inSize        */     (UINT16)(sizeof(Sign_In)),
    /* outSize       */     (UINT16)(sizeof(Sign_Out)),
    /* offsetOfTypes */     offsetof(Sign_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(Sign_In, digest)),
                             (UINT16)(offsetof(Sign_In, inScheme)),
                             (UINT16)(offsetof(Sign_In, validation))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             TPMT_TK_HASHCHECK_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _SignDataAddress (&_SignData)
#else
#define _SignDataAddress 0
#endif // CC_Sign

#if CC_SetCommandCodeAuditStatus

#include "SetCommandCodeAuditStatus_fp.h"

typedef TPM_RC  (SetCommandCodeAuditStatus_Entry)(
    SetCommandCodeAuditStatus_In            *in
);

typedef const struct {
    SetCommandCodeAuditStatus_Entry     *entry;
    UINT16                              inSize;
    UINT16                              outSize;
    UINT16                              offsetOfTypes;
    UINT16                              paramOffsets[3];
    BYTE                                types[6];
} SetCommandCodeAuditStatus_COMMAND_DESCRIPTOR_t;

SetCommandCodeAuditStatus_COMMAND_DESCRIPTOR_t _SetCommandCodeAuditStatusData = {
    /* entry         */                 &TPM2_SetCommandCodeAuditStatus,
    /* inSize        */                 (UINT16)(sizeof(SetCommandCodeAuditStatus_In)),
    /* outSize       */                 0,
    /* offsetOfTypes */                 offsetof(SetCommandCodeAuditStatus_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */                 {(UINT16)(offsetof(SetCommandCodeAuditStatus_In, auditAlg)),
                                         (UINT16)(offsetof(SetCommandCodeAuditStatus_In, setList)),
                                         (UINT16)(offsetof(SetCommandCodeAuditStatus_In, clearList))},
    /* types         */                 {TPMI_RH_PROVISION_H_UNMARSHAL,
                                         TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                                         TPML_CC_P_UNMARSHAL,
                                         TPML_CC_P_UNMARSHAL,
                                         END_OF_LIST,
                                         END_OF_LIST}
};

#define _SetCommandCodeAuditStatusDataAddress (&_SetCommandCodeAuditStatusData)
#else
#define _SetCommandCodeAuditStatusDataAddress 0
#endif // CC_SetCommandCodeAuditStatus

#if CC_PCR_Extend

#include "PCR_Extend_fp.h"

typedef TPM_RC  (PCR_Extend_Entry)(
    PCR_Extend_In               *in
);

typedef const struct {
    PCR_Extend_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PCR_Extend_COMMAND_DESCRIPTOR_t;

PCR_Extend_COMMAND_DESCRIPTOR_t _PCR_ExtendData = {
    /* entry         */     &TPM2_PCR_Extend,
    /* inSize        */     (UINT16)(sizeof(PCR_Extend_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PCR_Extend_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PCR_Extend_In, digests))},
    /* types         */     {TPMI_DH_PCR_H_UNMARSHAL + ADD_FLAG,
                             TPML_DIGEST_VALUES_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PCR_ExtendDataAddress (&_PCR_ExtendData)
#else
#define _PCR_ExtendDataAddress 0
#endif // CC_PCR_Extend

#if CC_PCR_Event

#include "PCR_Event_fp.h"

typedef TPM_RC  (PCR_Event_Entry)(
    PCR_Event_In                *in,
    PCR_Event_Out               *out
);

typedef const struct {
    PCR_Event_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} PCR_Event_COMMAND_DESCRIPTOR_t;

PCR_Event_COMMAND_DESCRIPTOR_t _PCR_EventData = {
    /* entry         */     &TPM2_PCR_Event,
    /* inSize        */     (UINT16)(sizeof(PCR_Event_In)),
    /* outSize       */     (UINT16)(sizeof(PCR_Event_Out)),
    /* offsetOfTypes */     offsetof(PCR_Event_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PCR_Event_In, eventData))},
    /* types         */     {TPMI_DH_PCR_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_EVENT_P_UNMARSHAL,
                             END_OF_LIST,
                             TPML_DIGEST_VALUES_P_MARSHAL,
                             END_OF_LIST}
};

#define _PCR_EventDataAddress (&_PCR_EventData)
#else
#define _PCR_EventDataAddress 0
#endif // CC_PCR_Event

#if CC_PCR_Read

#include "PCR_Read_fp.h"

typedef TPM_RC  (PCR_Read_Entry)(
    PCR_Read_In                 *in,
    PCR_Read_Out                *out
);

typedef const struct {
    PCR_Read_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[6];
} PCR_Read_COMMAND_DESCRIPTOR_t;

PCR_Read_COMMAND_DESCRIPTOR_t _PCR_ReadData = {
    /* entry         */     &TPM2_PCR_Read,
    /* inSize        */     (UINT16)(sizeof(PCR_Read_In)),
    /* outSize       */     (UINT16)(sizeof(PCR_Read_Out)),
    /* offsetOfTypes */     offsetof(PCR_Read_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PCR_Read_Out, pcrSelectionOut)),
                             (UINT16)(offsetof(PCR_Read_Out, pcrValues))},
    /* types         */     {TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             UINT32_P_MARSHAL,
                             TPML_PCR_SELECTION_P_MARSHAL,
                             TPML_DIGEST_P_MARSHAL,
                             END_OF_LIST}
};

#define _PCR_ReadDataAddress (&_PCR_ReadData)
#else
#define _PCR_ReadDataAddress 0
#endif // CC_PCR_Read

#if CC_PCR_Allocate

#include "PCR_Allocate_fp.h"

typedef TPM_RC  (PCR_Allocate_Entry)(
    PCR_Allocate_In             *in,
    PCR_Allocate_Out            *out
);

typedef const struct {
    PCR_Allocate_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[4];
    BYTE                    types[8];
} PCR_Allocate_COMMAND_DESCRIPTOR_t;

PCR_Allocate_COMMAND_DESCRIPTOR_t _PCR_AllocateData = {
    /* entry         */     &TPM2_PCR_Allocate,
    /* inSize        */     (UINT16)(sizeof(PCR_Allocate_In)),
    /* outSize       */     (UINT16)(sizeof(PCR_Allocate_Out)),
    /* offsetOfTypes */     offsetof(PCR_Allocate_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PCR_Allocate_In, pcrAllocation)),
                             (UINT16)(offsetof(PCR_Allocate_Out, maxPCR)),
                             (UINT16)(offsetof(PCR_Allocate_Out, sizeNeeded)),
                             (UINT16)(offsetof(PCR_Allocate_Out, sizeAvailable))},
    /* types         */     {TPMI_RH_PLATFORM_H_UNMARSHAL,
                             TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMI_YES_NO_P_MARSHAL,
                             UINT32_P_MARSHAL,
                             UINT32_P_MARSHAL,
                             UINT32_P_MARSHAL,
                             END_OF_LIST}
};

#define _PCR_AllocateDataAddress (&_PCR_AllocateData)
#else
#define _PCR_AllocateDataAddress 0
#endif // CC_PCR_Allocate

#if CC_PCR_SetAuthPolicy

#include "PCR_SetAuthPolicy_fp.h"

typedef TPM_RC  (PCR_SetAuthPolicy_Entry)(
    PCR_SetAuthPolicy_In            *in
);

typedef const struct {
    PCR_SetAuthPolicy_Entry     *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[3];
    BYTE                        types[6];
} PCR_SetAuthPolicy_COMMAND_DESCRIPTOR_t;

PCR_SetAuthPolicy_COMMAND_DESCRIPTOR_t _PCR_SetAuthPolicyData = {
    /* entry         */         &TPM2_PCR_SetAuthPolicy,
    /* inSize        */         (UINT16)(sizeof(PCR_SetAuthPolicy_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(PCR_SetAuthPolicy_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(PCR_SetAuthPolicy_In, authPolicy)),
                                 (UINT16)(offsetof(PCR_SetAuthPolicy_In, hashAlg)),
                                 (UINT16)(offsetof(PCR_SetAuthPolicy_In, pcrNum))},
    /* types         */         {TPMI_RH_PLATFORM_H_UNMARSHAL,
                                 TPM2B_DIGEST_P_UNMARSHAL,
                                 TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                                 TPMI_DH_PCR_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _PCR_SetAuthPolicyDataAddress (&_PCR_SetAuthPolicyData)
#else
#define _PCR_SetAuthPolicyDataAddress 0
#endif // CC_PCR_SetAuthPolicy

#if CC_PCR_SetAuthValue

#include "PCR_SetAuthValue_fp.h"

typedef TPM_RC  (PCR_SetAuthValue_Entry)(
    PCR_SetAuthValue_In         *in
);

typedef const struct {
    PCR_SetAuthValue_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PCR_SetAuthValue_COMMAND_DESCRIPTOR_t;

PCR_SetAuthValue_COMMAND_DESCRIPTOR_t _PCR_SetAuthValueData = {
    /* entry         */     &TPM2_PCR_SetAuthValue,
    /* inSize        */     (UINT16)(sizeof(PCR_SetAuthValue_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PCR_SetAuthValue_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PCR_SetAuthValue_In, auth))},
    /* types         */     {TPMI_DH_PCR_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PCR_SetAuthValueDataAddress (&_PCR_SetAuthValueData)
#else
#define _PCR_SetAuthValueDataAddress 0
#endif // CC_PCR_SetAuthValue

#if CC_PCR_Reset

#include "PCR_Reset_fp.h"

typedef TPM_RC  (PCR_Reset_Entry)(
    PCR_Reset_In                *in
);

typedef const struct {
    PCR_Reset_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} PCR_Reset_COMMAND_DESCRIPTOR_t;

PCR_Reset_COMMAND_DESCRIPTOR_t _PCR_ResetData = {
    /* entry         */     &TPM2_PCR_Reset,
    /* inSize        */     (UINT16)(sizeof(PCR_Reset_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PCR_Reset_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_DH_PCR_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PCR_ResetDataAddress (&_PCR_ResetData)
#else
#define _PCR_ResetDataAddress 0
#endif // CC_PCR_Reset

#if CC_PolicySigned

#include "PolicySigned_fp.h"

typedef TPM_RC  (PolicySigned_Entry)(
    PolicySigned_In             *in,
    PolicySigned_Out            *out
);

typedef const struct {
    PolicySigned_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[7];
    BYTE                    types[11];
} PolicySigned_COMMAND_DESCRIPTOR_t;

PolicySigned_COMMAND_DESCRIPTOR_t _PolicySignedData = {
    /* entry         */     &TPM2_PolicySigned,
    /* inSize        */     (UINT16)(sizeof(PolicySigned_In)),
    /* outSize       */     (UINT16)(sizeof(PolicySigned_Out)),
    /* offsetOfTypes */     offsetof(PolicySigned_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicySigned_In, policySession)),
                             (UINT16)(offsetof(PolicySigned_In, nonceTPM)),
                             (UINT16)(offsetof(PolicySigned_In, cpHashA)),
                             (UINT16)(offsetof(PolicySigned_In, policyRef)),
                             (UINT16)(offsetof(PolicySigned_In, expiration)),
                             (UINT16)(offsetof(PolicySigned_In, auth)),
                             (UINT16)(offsetof(PolicySigned_Out, policyTicket))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             INT32_P_UNMARSHAL,
                             TPMT_SIGNATURE_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_TIMEOUT_P_MARSHAL,
                             TPMT_TK_AUTH_P_MARSHAL,
                             END_OF_LIST}
};

#define _PolicySignedDataAddress (&_PolicySignedData)
#else
#define _PolicySignedDataAddress 0
#endif // CC_PolicySigned

#if CC_PolicySecret

#include "PolicySecret_fp.h"

typedef TPM_RC  (PolicySecret_Entry)(
    PolicySecret_In             *in,
    PolicySecret_Out            *out
);

typedef const struct {
    PolicySecret_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[6];
    BYTE                    types[10];
} PolicySecret_COMMAND_DESCRIPTOR_t;

PolicySecret_COMMAND_DESCRIPTOR_t _PolicySecretData = {
    /* entry         */     &TPM2_PolicySecret,
    /* inSize        */     (UINT16)(sizeof(PolicySecret_In)),
    /* outSize       */     (UINT16)(sizeof(PolicySecret_Out)),
    /* offsetOfTypes */     offsetof(PolicySecret_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicySecret_In, policySession)),
                             (UINT16)(offsetof(PolicySecret_In, nonceTPM)),
                             (UINT16)(offsetof(PolicySecret_In, cpHashA)),
                             (UINT16)(offsetof(PolicySecret_In, policyRef)),
                             (UINT16)(offsetof(PolicySecret_In, expiration)),
                             (UINT16)(offsetof(PolicySecret_Out, policyTicket))},
    /* types         */     {TPMI_DH_ENTITY_H_UNMARSHAL,
                             TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             INT32_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_TIMEOUT_P_MARSHAL,
                             TPMT_TK_AUTH_P_MARSHAL,
                             END_OF_LIST}
};

#define _PolicySecretDataAddress (&_PolicySecretData)
#else
#define _PolicySecretDataAddress 0
#endif // CC_PolicySecret

#if CC_PolicyTicket

#include "PolicyTicket_fp.h"

typedef TPM_RC  (PolicyTicket_Entry)(
    PolicyTicket_In             *in
);

typedef const struct {
    PolicyTicket_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[8];
} PolicyTicket_COMMAND_DESCRIPTOR_t;

PolicyTicket_COMMAND_DESCRIPTOR_t _PolicyTicketData = {
    /* entry         */     &TPM2_PolicyTicket,
    /* inSize        */     (UINT16)(sizeof(PolicyTicket_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyTicket_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyTicket_In, timeout)),
                             (UINT16)(offsetof(PolicyTicket_In, cpHashA)),
                             (UINT16)(offsetof(PolicyTicket_In, policyRef)),
                             (UINT16)(offsetof(PolicyTicket_In, authName)),
                             (UINT16)(offsetof(PolicyTicket_In, ticket))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_TIMEOUT_P_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             TPM2B_NAME_P_UNMARSHAL,
                             TPMT_TK_AUTH_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyTicketDataAddress (&_PolicyTicketData)
#else
#define _PolicyTicketDataAddress 0
#endif // CC_PolicyTicket

#if CC_PolicyOR

#include "PolicyOR_fp.h"

typedef TPM_RC  (PolicyOR_Entry)(
    PolicyOR_In                 *in
);

typedef const struct {
    PolicyOR_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyOR_COMMAND_DESCRIPTOR_t;

PolicyOR_COMMAND_DESCRIPTOR_t _PolicyORData = {
    /* entry         */     &TPM2_PolicyOR,
    /* inSize        */     (UINT16)(sizeof(PolicyOR_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyOR_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyOR_In, pHashList))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPML_DIGEST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyORDataAddress (&_PolicyORData)
#else
#define _PolicyORDataAddress 0
#endif // CC_PolicyOR

#if CC_PolicyPCR

#include "PolicyPCR_fp.h"

typedef TPM_RC  (PolicyPCR_Entry)(
    PolicyPCR_In                *in
);

typedef const struct {
    PolicyPCR_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} PolicyPCR_COMMAND_DESCRIPTOR_t;

PolicyPCR_COMMAND_DESCRIPTOR_t _PolicyPCRData = {
    /* entry         */     &TPM2_PolicyPCR,
    /* inSize        */     (UINT16)(sizeof(PolicyPCR_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyPCR_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyPCR_In, pcrDigest)),
                             (UINT16)(offsetof(PolicyPCR_In, pcrs))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyPCRDataAddress (&_PolicyPCRData)
#else
#define _PolicyPCRDataAddress 0
#endif // CC_PolicyPCR

#if CC_PolicyLocality

#include "PolicyLocality_fp.h"

typedef TPM_RC  (PolicyLocality_Entry)(
    PolicyLocality_In           *in
);

typedef const struct {
    PolicyLocality_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyLocality_COMMAND_DESCRIPTOR_t;

PolicyLocality_COMMAND_DESCRIPTOR_t _PolicyLocalityData = {
    /* entry         */     &TPM2_PolicyLocality,
    /* inSize        */     (UINT16)(sizeof(PolicyLocality_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyLocality_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyLocality_In, locality))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPMA_LOCALITY_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyLocalityDataAddress (&_PolicyLocalityData)
#else
#define _PolicyLocalityDataAddress 0
#endif // CC_PolicyLocality

#if CC_PolicyNV

#include "PolicyNV_fp.h"

typedef TPM_RC  (PolicyNV_Entry)(
    PolicyNV_In                 *in
);

typedef const struct {
    PolicyNV_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[5];
    BYTE                    types[8];
} PolicyNV_COMMAND_DESCRIPTOR_t;

PolicyNV_COMMAND_DESCRIPTOR_t _PolicyNVData = {
    /* entry         */     &TPM2_PolicyNV,
    /* inSize        */     (UINT16)(sizeof(PolicyNV_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyNV_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyNV_In, nvIndex)),
                             (UINT16)(offsetof(PolicyNV_In, policySession)),
                             (UINT16)(offsetof(PolicyNV_In, operandB)),
                             (UINT16)(offsetof(PolicyNV_In, offset)),
                             (UINT16)(offsetof(PolicyNV_In, operation))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_OPERAND_P_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             TPM_EO_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyNVDataAddress (&_PolicyNVData)
#else
#define _PolicyNVDataAddress 0
#endif // CC_PolicyNV

#if CC_PolicyCounterTimer

#include "PolicyCounterTimer_fp.h"

typedef TPM_RC  (PolicyCounterTimer_Entry)(
    PolicyCounterTimer_In           *in
);

typedef const struct {
    PolicyCounterTimer_Entry    *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[3];
    BYTE                        types[6];
} PolicyCounterTimer_COMMAND_DESCRIPTOR_t;

PolicyCounterTimer_COMMAND_DESCRIPTOR_t _PolicyCounterTimerData = {
    /* entry         */         &TPM2_PolicyCounterTimer,
    /* inSize        */         (UINT16)(sizeof(PolicyCounterTimer_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(PolicyCounterTimer_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(PolicyCounterTimer_In, operandB)),
                                 (UINT16)(offsetof(PolicyCounterTimer_In, offset)),
                                 (UINT16)(offsetof(PolicyCounterTimer_In, operation))},
    /* types         */         {TPMI_SH_POLICY_H_UNMARSHAL,
                                 TPM2B_OPERAND_P_UNMARSHAL,
                                 UINT16_P_UNMARSHAL,
                                 TPM_EO_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _PolicyCounterTimerDataAddress (&_PolicyCounterTimerData)
#else
#define _PolicyCounterTimerDataAddress 0
#endif // CC_PolicyCounterTimer

#if CC_PolicyCommandCode

#include "PolicyCommandCode_fp.h"

typedef TPM_RC  (PolicyCommandCode_Entry)(
    PolicyCommandCode_In            *in
);

typedef const struct {
    PolicyCommandCode_Entry     *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[1];
    BYTE                        types[4];
} PolicyCommandCode_COMMAND_DESCRIPTOR_t;

PolicyCommandCode_COMMAND_DESCRIPTOR_t _PolicyCommandCodeData = {
    /* entry         */         &TPM2_PolicyCommandCode,
    /* inSize        */         (UINT16)(sizeof(PolicyCommandCode_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(PolicyCommandCode_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(PolicyCommandCode_In, code))},
    /* types         */         {TPMI_SH_POLICY_H_UNMARSHAL,
                                 TPM_CC_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _PolicyCommandCodeDataAddress (&_PolicyCommandCodeData)
#else
#define _PolicyCommandCodeDataAddress 0
#endif // CC_PolicyCommandCode

#if CC_PolicyPhysicalPresence

#include "PolicyPhysicalPresence_fp.h"

typedef TPM_RC  (PolicyPhysicalPresence_Entry)(
    PolicyPhysicalPresence_In           *in
);

typedef const struct {
    PolicyPhysicalPresence_Entry    *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    BYTE                            types[3];
} PolicyPhysicalPresence_COMMAND_DESCRIPTOR_t;

PolicyPhysicalPresence_COMMAND_DESCRIPTOR_t _PolicyPhysicalPresenceData = {
    /* entry         */             &TPM2_PolicyPhysicalPresence,
    /* inSize        */             (UINT16)(sizeof(PolicyPhysicalPresence_In)),
    /* outSize       */             0,
    /* offsetOfTypes */             offsetof(PolicyPhysicalPresence_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             // No parameter offsets;
    /* types         */             {TPMI_SH_POLICY_H_UNMARSHAL,
                                     END_OF_LIST,
                                     END_OF_LIST}
};

#define _PolicyPhysicalPresenceDataAddress (&_PolicyPhysicalPresenceData)
#else
#define _PolicyPhysicalPresenceDataAddress 0
#endif // CC_PolicyPhysicalPresence

#if CC_PolicyCpHash

#include "PolicyCpHash_fp.h"

typedef TPM_RC  (PolicyCpHash_Entry)(
    PolicyCpHash_In             *in
);

typedef const struct {
    PolicyCpHash_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyCpHash_COMMAND_DESCRIPTOR_t;

PolicyCpHash_COMMAND_DESCRIPTOR_t _PolicyCpHashData = {
    /* entry         */     &TPM2_PolicyCpHash,
    /* inSize        */     (UINT16)(sizeof(PolicyCpHash_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyCpHash_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyCpHash_In, cpHashA))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyCpHashDataAddress (&_PolicyCpHashData)
#else
#define _PolicyCpHashDataAddress 0
#endif // CC_PolicyCpHash

#if CC_PolicyNameHash

#include "PolicyNameHash_fp.h"

typedef TPM_RC  (PolicyNameHash_Entry)(
    PolicyNameHash_In           *in
);

typedef const struct {
    PolicyNameHash_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyNameHash_COMMAND_DESCRIPTOR_t;

PolicyNameHash_COMMAND_DESCRIPTOR_t _PolicyNameHashData = {
    /* entry         */     &TPM2_PolicyNameHash,
    /* inSize        */     (UINT16)(sizeof(PolicyNameHash_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyNameHash_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyNameHash_In, nameHash))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyNameHashDataAddress (&_PolicyNameHashData)
#else
#define _PolicyNameHashDataAddress 0
#endif // CC_PolicyNameHash

#if CC_PolicyDuplicationSelect

#include "PolicyDuplicationSelect_fp.h"

typedef TPM_RC  (PolicyDuplicationSelect_Entry)(
    PolicyDuplicationSelect_In          *in
);

typedef const struct {
    PolicyDuplicationSelect_Entry   *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    UINT16                          paramOffsets[3];
    BYTE                            types[6];
} PolicyDuplicationSelect_COMMAND_DESCRIPTOR_t;

PolicyDuplicationSelect_COMMAND_DESCRIPTOR_t _PolicyDuplicationSelectData = {
    /* entry         */             &TPM2_PolicyDuplicationSelect,
    /* inSize        */             (UINT16)(sizeof(PolicyDuplicationSelect_In)),
    /* outSize       */             0,
    /* offsetOfTypes */             offsetof(PolicyDuplicationSelect_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             {(UINT16)(offsetof(PolicyDuplicationSelect_In, objectName)),
                                     (UINT16)(offsetof(PolicyDuplicationSelect_In, newParentName)),
                                     (UINT16)(offsetof(PolicyDuplicationSelect_In, includeObject))},
    /* types         */             {TPMI_SH_POLICY_H_UNMARSHAL,
                                     TPM2B_NAME_P_UNMARSHAL,
                                     TPM2B_NAME_P_UNMARSHAL,
                                     TPMI_YES_NO_P_UNMARSHAL,
                                     END_OF_LIST,
                                     END_OF_LIST}
};

#define _PolicyDuplicationSelectDataAddress (&_PolicyDuplicationSelectData)
#else
#define _PolicyDuplicationSelectDataAddress 0
#endif // CC_PolicyDuplicationSelect

#if CC_PolicyAuthorize

#include "PolicyAuthorize_fp.h"

typedef TPM_RC  (PolicyAuthorize_Entry)(
    PolicyAuthorize_In          *in
);

typedef const struct {
    PolicyAuthorize_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[4];
    BYTE                    types[7];
} PolicyAuthorize_COMMAND_DESCRIPTOR_t;

PolicyAuthorize_COMMAND_DESCRIPTOR_t _PolicyAuthorizeData = {
    /* entry         */     &TPM2_PolicyAuthorize,
    /* inSize        */     (UINT16)(sizeof(PolicyAuthorize_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyAuthorize_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyAuthorize_In, approvedPolicy)),
                             (UINT16)(offsetof(PolicyAuthorize_In, policyRef)),
                             (UINT16)(offsetof(PolicyAuthorize_In, keySign)),
                             (UINT16)(offsetof(PolicyAuthorize_In, checkTicket))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPM2B_NONCE_P_UNMARSHAL,
                             TPM2B_NAME_P_UNMARSHAL,
                             TPMT_TK_VERIFIED_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyAuthorizeDataAddress (&_PolicyAuthorizeData)
#else
#define _PolicyAuthorizeDataAddress 0
#endif // CC_PolicyAuthorize

#if CC_PolicyAuthValue

#include "PolicyAuthValue_fp.h"

typedef TPM_RC  (PolicyAuthValue_Entry)(
    PolicyAuthValue_In          *in
);

typedef const struct {
    PolicyAuthValue_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} PolicyAuthValue_COMMAND_DESCRIPTOR_t;

PolicyAuthValue_COMMAND_DESCRIPTOR_t _PolicyAuthValueData = {
    /* entry         */     &TPM2_PolicyAuthValue,
    /* inSize        */     (UINT16)(sizeof(PolicyAuthValue_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyAuthValue_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyAuthValueDataAddress (&_PolicyAuthValueData)
#else
#define _PolicyAuthValueDataAddress 0
#endif // CC_PolicyAuthValue

#if CC_PolicyPassword

#include "PolicyPassword_fp.h"

typedef TPM_RC  (PolicyPassword_Entry)(
    PolicyPassword_In           *in
);

typedef const struct {
    PolicyPassword_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} PolicyPassword_COMMAND_DESCRIPTOR_t;

PolicyPassword_COMMAND_DESCRIPTOR_t _PolicyPasswordData = {
    /* entry         */     &TPM2_PolicyPassword,
    /* inSize        */     (UINT16)(sizeof(PolicyPassword_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyPassword_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyPasswordDataAddress (&_PolicyPasswordData)
#else
#define _PolicyPasswordDataAddress 0
#endif // CC_PolicyPassword

#if CC_PolicyGetDigest

#include "PolicyGetDigest_fp.h"

typedef TPM_RC  (PolicyGetDigest_Entry)(
    PolicyGetDigest_In          *in,
    PolicyGetDigest_Out         *out
);

typedef const struct {
    PolicyGetDigest_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} PolicyGetDigest_COMMAND_DESCRIPTOR_t;

PolicyGetDigest_COMMAND_DESCRIPTOR_t _PolicyGetDigestData = {
    /* entry         */     &TPM2_PolicyGetDigest,
    /* inSize        */     (UINT16)(sizeof(PolicyGetDigest_In)),
    /* outSize       */     (UINT16)(sizeof(PolicyGetDigest_Out)),
    /* offsetOfTypes */     offsetof(PolicyGetDigest_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_DIGEST_P_MARSHAL,
                             END_OF_LIST}
};

#define _PolicyGetDigestDataAddress (&_PolicyGetDigestData)
#else
#define _PolicyGetDigestDataAddress 0
#endif // CC_PolicyGetDigest

#if CC_PolicyNvWritten

#include "PolicyNvWritten_fp.h"

typedef TPM_RC  (PolicyNvWritten_Entry)(
    PolicyNvWritten_In          *in
);

typedef const struct {
    PolicyNvWritten_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyNvWritten_COMMAND_DESCRIPTOR_t;

PolicyNvWritten_COMMAND_DESCRIPTOR_t _PolicyNvWrittenData = {
    /* entry         */     &TPM2_PolicyNvWritten,
    /* inSize        */     (UINT16)(sizeof(PolicyNvWritten_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyNvWritten_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyNvWritten_In, writtenSet))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPMI_YES_NO_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyNvWrittenDataAddress (&_PolicyNvWrittenData)
#else
#define _PolicyNvWrittenDataAddress 0
#endif // CC_PolicyNvWritten

#if CC_PolicyTemplate

#include "PolicyTemplate_fp.h"

typedef TPM_RC  (PolicyTemplate_Entry)(
    PolicyTemplate_In           *in
);

typedef const struct {
    PolicyTemplate_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} PolicyTemplate_COMMAND_DESCRIPTOR_t;

PolicyTemplate_COMMAND_DESCRIPTOR_t _PolicyTemplateData = {
    /* entry         */     &TPM2_PolicyTemplate,
    /* inSize        */     (UINT16)(sizeof(PolicyTemplate_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PolicyTemplate_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PolicyTemplate_In, templateHash))},
    /* types         */     {TPMI_SH_POLICY_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PolicyTemplateDataAddress (&_PolicyTemplateData)
#else
#define _PolicyTemplateDataAddress 0
#endif // CC_PolicyTemplate

#if CC_PolicyAuthorizeNV

#include "PolicyAuthorizeNV_fp.h"

typedef TPM_RC  (PolicyAuthorizeNV_Entry)(
    PolicyAuthorizeNV_In            *in
);

typedef const struct {
    PolicyAuthorizeNV_Entry     *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[2];
    BYTE                        types[5];
} PolicyAuthorizeNV_COMMAND_DESCRIPTOR_t;

PolicyAuthorizeNV_COMMAND_DESCRIPTOR_t _PolicyAuthorizeNVData = {
    /* entry         */         &TPM2_PolicyAuthorizeNV,
    /* inSize        */         (UINT16)(sizeof(PolicyAuthorizeNV_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(PolicyAuthorizeNV_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(PolicyAuthorizeNV_In, nvIndex)),
                                 (UINT16)(offsetof(PolicyAuthorizeNV_In, policySession))},
    /* types         */         {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                                 TPMI_RH_NV_INDEX_H_UNMARSHAL,
                                 TPMI_SH_POLICY_H_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _PolicyAuthorizeNVDataAddress (&_PolicyAuthorizeNVData)
#else
#define _PolicyAuthorizeNVDataAddress 0
#endif // CC_PolicyAuthorizeNV

#if CC_CreatePrimary

#include "CreatePrimary_fp.h"

typedef TPM_RC  (CreatePrimary_Entry)(
    CreatePrimary_In            *in,
    CreatePrimary_Out           *out
);

typedef const struct {
    CreatePrimary_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[9];
    BYTE                    types[13];
} CreatePrimary_COMMAND_DESCRIPTOR_t;

CreatePrimary_COMMAND_DESCRIPTOR_t _CreatePrimaryData = {
    /* entry         */     &TPM2_CreatePrimary,
    /* inSize        */     (UINT16)(sizeof(CreatePrimary_In)),
    /* outSize       */     (UINT16)(sizeof(CreatePrimary_Out)),
    /* offsetOfTypes */     offsetof(CreatePrimary_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(CreatePrimary_In, inSensitive)),
                             (UINT16)(offsetof(CreatePrimary_In, inPublic)),
                             (UINT16)(offsetof(CreatePrimary_In, outsideInfo)),
                             (UINT16)(offsetof(CreatePrimary_In, creationPCR)),
                             (UINT16)(offsetof(CreatePrimary_Out, outPublic)),
                             (UINT16)(offsetof(CreatePrimary_Out, creationData)),
                             (UINT16)(offsetof(CreatePrimary_Out, creationHash)),
                             (UINT16)(offsetof(CreatePrimary_Out, creationTicket)),
                             (UINT16)(offsetof(CreatePrimary_Out, name))},
    /* types         */     {TPMI_RH_HIERARCHY_H_UNMARSHAL + ADD_FLAG,
                             TPM2B_SENSITIVE_CREATE_P_UNMARSHAL,
                             TPM2B_PUBLIC_P_UNMARSHAL,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPML_PCR_SELECTION_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM_HANDLE_H_MARSHAL,
                             TPM2B_PUBLIC_P_MARSHAL,
                             TPM2B_CREATION_DATA_P_MARSHAL,
                             TPM2B_DIGEST_P_MARSHAL,
                             TPMT_TK_CREATION_P_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _CreatePrimaryDataAddress (&_CreatePrimaryData)
#else
#define _CreatePrimaryDataAddress 0
#endif // CC_CreatePrimary

#if CC_HierarchyControl

#include "HierarchyControl_fp.h"

typedef TPM_RC  (HierarchyControl_Entry)(
    HierarchyControl_In         *in
);

typedef const struct {
    HierarchyControl_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} HierarchyControl_COMMAND_DESCRIPTOR_t;

HierarchyControl_COMMAND_DESCRIPTOR_t _HierarchyControlData = {
    /* entry         */     &TPM2_HierarchyControl,
    /* inSize        */     (UINT16)(sizeof(HierarchyControl_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(HierarchyControl_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(HierarchyControl_In, enable)),
                             (UINT16)(offsetof(HierarchyControl_In, state))},
    /* types         */     {TPMI_RH_HIERARCHY_H_UNMARSHAL,
                             TPMI_RH_ENABLES_P_UNMARSHAL,
                             TPMI_YES_NO_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _HierarchyControlDataAddress (&_HierarchyControlData)
#else
#define _HierarchyControlDataAddress 0
#endif // CC_HierarchyControl

#if CC_SetPrimaryPolicy

#include "SetPrimaryPolicy_fp.h"

typedef TPM_RC  (SetPrimaryPolicy_Entry)(
    SetPrimaryPolicy_In         *in
);

typedef const struct {
    SetPrimaryPolicy_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} SetPrimaryPolicy_COMMAND_DESCRIPTOR_t;

SetPrimaryPolicy_COMMAND_DESCRIPTOR_t _SetPrimaryPolicyData = {
    /* entry         */     &TPM2_SetPrimaryPolicy,
    /* inSize        */     (UINT16)(sizeof(SetPrimaryPolicy_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(SetPrimaryPolicy_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(SetPrimaryPolicy_In, authPolicy)),
                             (UINT16)(offsetof(SetPrimaryPolicy_In, hashAlg))},
    /* types         */     {TPMI_RH_HIERARCHY_AUTH_H_UNMARSHAL,
                             TPM2B_DIGEST_P_UNMARSHAL,
                             TPMI_ALG_HASH_P_UNMARSHAL + ADD_FLAG,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _SetPrimaryPolicyDataAddress (&_SetPrimaryPolicyData)
#else
#define _SetPrimaryPolicyDataAddress 0
#endif // CC_SetPrimaryPolicy

#if CC_ChangePPS

#include "ChangePPS_fp.h"

typedef TPM_RC  (ChangePPS_Entry)(
    ChangePPS_In                *in
);

typedef const struct {
    ChangePPS_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} ChangePPS_COMMAND_DESCRIPTOR_t;

ChangePPS_COMMAND_DESCRIPTOR_t _ChangePPSData = {
    /* entry         */     &TPM2_ChangePPS,
    /* inSize        */     (UINT16)(sizeof(ChangePPS_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(ChangePPS_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_RH_PLATFORM_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ChangePPSDataAddress (&_ChangePPSData)
#else
#define _ChangePPSDataAddress 0
#endif // CC_ChangePPS

#if CC_ChangeEPS

#include "ChangeEPS_fp.h"

typedef TPM_RC  (ChangeEPS_Entry)(
    ChangeEPS_In                *in
);

typedef const struct {
    ChangeEPS_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} ChangeEPS_COMMAND_DESCRIPTOR_t;

ChangeEPS_COMMAND_DESCRIPTOR_t _ChangeEPSData = {
    /* entry         */     &TPM2_ChangeEPS,
    /* inSize        */     (UINT16)(sizeof(ChangeEPS_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(ChangeEPS_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_RH_PLATFORM_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ChangeEPSDataAddress (&_ChangeEPSData)
#else
#define _ChangeEPSDataAddress 0
#endif // CC_ChangeEPS

#if CC_Clear

#include "Clear_fp.h"

typedef TPM_RC  (Clear_Entry)(
    Clear_In                    *in
);

typedef const struct {
    Clear_Entry             *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} Clear_COMMAND_DESCRIPTOR_t;

Clear_COMMAND_DESCRIPTOR_t _ClearData = {
    /* entry         */     &TPM2_Clear,
    /* inSize        */     (UINT16)(sizeof(Clear_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(Clear_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_RH_CLEAR_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ClearDataAddress (&_ClearData)
#else
#define _ClearDataAddress 0
#endif // CC_Clear

#if CC_ClearControl

#include "ClearControl_fp.h"

typedef TPM_RC  (ClearControl_Entry)(
    ClearControl_In             *in
);

typedef const struct {
    ClearControl_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} ClearControl_COMMAND_DESCRIPTOR_t;

ClearControl_COMMAND_DESCRIPTOR_t _ClearControlData = {
    /* entry         */     &TPM2_ClearControl,
    /* inSize        */     (UINT16)(sizeof(ClearControl_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(ClearControl_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ClearControl_In, disable))},
    /* types         */     {TPMI_RH_CLEAR_H_UNMARSHAL,
                             TPMI_YES_NO_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ClearControlDataAddress (&_ClearControlData)
#else
#define _ClearControlDataAddress 0
#endif // CC_ClearControl

#if CC_HierarchyChangeAuth

#include "HierarchyChangeAuth_fp.h"

typedef TPM_RC  (HierarchyChangeAuth_Entry)(
    HierarchyChangeAuth_In          *in
);

typedef const struct {
    HierarchyChangeAuth_Entry   *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[1];
    BYTE                        types[4];
} HierarchyChangeAuth_COMMAND_DESCRIPTOR_t;

HierarchyChangeAuth_COMMAND_DESCRIPTOR_t _HierarchyChangeAuthData = {
    /* entry         */         &TPM2_HierarchyChangeAuth,
    /* inSize        */         (UINT16)(sizeof(HierarchyChangeAuth_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(HierarchyChangeAuth_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(HierarchyChangeAuth_In, newAuth))},
    /* types         */         {TPMI_RH_HIERARCHY_AUTH_H_UNMARSHAL,
                                 TPM2B_AUTH_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _HierarchyChangeAuthDataAddress (&_HierarchyChangeAuthData)
#else
#define _HierarchyChangeAuthDataAddress 0
#endif // CC_HierarchyChangeAuth

#if CC_DictionaryAttackLockReset

#include "DictionaryAttackLockReset_fp.h"

typedef TPM_RC  (DictionaryAttackLockReset_Entry)(
    DictionaryAttackLockReset_In            *in
);

typedef const struct {
    DictionaryAttackLockReset_Entry     *entry;
    UINT16                              inSize;
    UINT16                              outSize;
    UINT16                              offsetOfTypes;
    BYTE                                types[3];
} DictionaryAttackLockReset_COMMAND_DESCRIPTOR_t;

DictionaryAttackLockReset_COMMAND_DESCRIPTOR_t _DictionaryAttackLockResetData = {
    /* entry         */                 &TPM2_DictionaryAttackLockReset,
    /* inSize        */                 (UINT16)(sizeof(DictionaryAttackLockReset_In)),
    /* outSize       */                 0,
    /* offsetOfTypes */                 offsetof(DictionaryAttackLockReset_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */                 // No parameter offsets;
    /* types         */                 {TPMI_RH_LOCKOUT_H_UNMARSHAL,
                                         END_OF_LIST,
                                         END_OF_LIST}
};

#define _DictionaryAttackLockResetDataAddress (&_DictionaryAttackLockResetData)
#else
#define _DictionaryAttackLockResetDataAddress 0
#endif // CC_DictionaryAttackLockReset

#if CC_DictionaryAttackParameters

#include "DictionaryAttackParameters_fp.h"

typedef TPM_RC  (DictionaryAttackParameters_Entry)(
    DictionaryAttackParameters_In           *in
);

typedef const struct {
    DictionaryAttackParameters_Entry    *entry;
    UINT16                              inSize;
    UINT16                              outSize;
    UINT16                              offsetOfTypes;
    UINT16                              paramOffsets[3];
    BYTE                                types[6];
} DictionaryAttackParameters_COMMAND_DESCRIPTOR_t;

DictionaryAttackParameters_COMMAND_DESCRIPTOR_t _DictionaryAttackParametersData = {
    /* entry         */                 &TPM2_DictionaryAttackParameters,
    /* inSize        */                 (UINT16)(sizeof(DictionaryAttackParameters_In)),
    /* outSize       */                 0,
    /* offsetOfTypes */                 offsetof(DictionaryAttackParameters_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */                 {(UINT16)(offsetof(DictionaryAttackParameters_In, newMaxTries)),
                                         (UINT16)(offsetof(DictionaryAttackParameters_In, newRecoveryTime)),
                                         (UINT16)(offsetof(DictionaryAttackParameters_In, lockoutRecovery))},
    /* types         */                 {TPMI_RH_LOCKOUT_H_UNMARSHAL,
                                         UINT32_P_UNMARSHAL,
                                         UINT32_P_UNMARSHAL,
                                         UINT32_P_UNMARSHAL,
                                         END_OF_LIST,
                                         END_OF_LIST}
};

#define _DictionaryAttackParametersDataAddress (&_DictionaryAttackParametersData)
#else
#define _DictionaryAttackParametersDataAddress 0
#endif // CC_DictionaryAttackParameters

#if CC_PP_Commands

#include "PP_Commands_fp.h"

typedef TPM_RC  (PP_Commands_Entry)(
    PP_Commands_In              *in
);

typedef const struct {
    PP_Commands_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} PP_Commands_COMMAND_DESCRIPTOR_t;

PP_Commands_COMMAND_DESCRIPTOR_t _PP_CommandsData = {
    /* entry         */     &TPM2_PP_Commands,
    /* inSize        */     (UINT16)(sizeof(PP_Commands_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(PP_Commands_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(PP_Commands_In, setList)),
                             (UINT16)(offsetof(PP_Commands_In, clearList))},
    /* types         */     {TPMI_RH_PLATFORM_H_UNMARSHAL,
                             TPML_CC_P_UNMARSHAL,
                             TPML_CC_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _PP_CommandsDataAddress (&_PP_CommandsData)
#else
#define _PP_CommandsDataAddress 0
#endif // CC_PP_Commands

#if CC_SetAlgorithmSet

#include "SetAlgorithmSet_fp.h"

typedef TPM_RC  (SetAlgorithmSet_Entry)(
    SetAlgorithmSet_In          *in
);

typedef const struct {
    SetAlgorithmSet_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} SetAlgorithmSet_COMMAND_DESCRIPTOR_t;

SetAlgorithmSet_COMMAND_DESCRIPTOR_t _SetAlgorithmSetData = {
    /* entry         */     &TPM2_SetAlgorithmSet,
    /* inSize        */     (UINT16)(sizeof(SetAlgorithmSet_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(SetAlgorithmSet_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(SetAlgorithmSet_In, algorithmSet))},
    /* types         */     {TPMI_RH_PLATFORM_H_UNMARSHAL,
                             UINT32_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _SetAlgorithmSetDataAddress (&_SetAlgorithmSetData)
#else
#define _SetAlgorithmSetDataAddress 0
#endif // CC_SetAlgorithmSet

#if CC_FieldUpgradeStart

#include "FieldUpgradeStart_fp.h"

typedef TPM_RC  (FieldUpgradeStart_Entry)(
    FieldUpgradeStart_In            *in
);

typedef const struct {
    FieldUpgradeStart_Entry     *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[3];
    BYTE                        types[6];
} FieldUpgradeStart_COMMAND_DESCRIPTOR_t;

FieldUpgradeStart_COMMAND_DESCRIPTOR_t _FieldUpgradeStartData = {
    /* entry         */         &TPM2_FieldUpgradeStart,
    /* inSize        */         (UINT16)(sizeof(FieldUpgradeStart_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(FieldUpgradeStart_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(FieldUpgradeStart_In, keyHandle)),
                                 (UINT16)(offsetof(FieldUpgradeStart_In, fuDigest)),
                                 (UINT16)(offsetof(FieldUpgradeStart_In, manifestSignature))},
    /* types         */         {TPMI_RH_PLATFORM_H_UNMARSHAL,
                                 TPMI_DH_OBJECT_H_UNMARSHAL,
                                 TPM2B_DIGEST_P_UNMARSHAL,
                                 TPMT_SIGNATURE_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _FieldUpgradeStartDataAddress (&_FieldUpgradeStartData)
#else
#define _FieldUpgradeStartDataAddress 0
#endif // CC_FieldUpgradeStart

#if CC_FieldUpgradeData

#include "FieldUpgradeData_fp.h"

typedef TPM_RC  (FieldUpgradeData_Entry)(
    FieldUpgradeData_In         *in,
    FieldUpgradeData_Out        *out
);

typedef const struct {
    FieldUpgradeData_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} FieldUpgradeData_COMMAND_DESCRIPTOR_t;

FieldUpgradeData_COMMAND_DESCRIPTOR_t _FieldUpgradeDataData = {
    /* entry         */     &TPM2_FieldUpgradeData,
    /* inSize        */     (UINT16)(sizeof(FieldUpgradeData_In)),
    /* outSize       */     (UINT16)(sizeof(FieldUpgradeData_Out)),
    /* offsetOfTypes */     offsetof(FieldUpgradeData_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(FieldUpgradeData_Out, firstDigest))},
    /* types         */     {TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMT_HA_P_MARSHAL,
                             TPMT_HA_P_MARSHAL,
                             END_OF_LIST}
};

#define _FieldUpgradeDataDataAddress (&_FieldUpgradeDataData)
#else
#define _FieldUpgradeDataDataAddress 0
#endif // CC_FieldUpgradeData

#if CC_FirmwareRead

#include "FirmwareRead_fp.h"

typedef TPM_RC  (FirmwareRead_Entry)(
    FirmwareRead_In             *in,
    FirmwareRead_Out            *out
);

typedef const struct {
    FirmwareRead_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} FirmwareRead_COMMAND_DESCRIPTOR_t;

FirmwareRead_COMMAND_DESCRIPTOR_t _FirmwareReadData = {
    /* entry         */     &TPM2_FirmwareRead,
    /* inSize        */     (UINT16)(sizeof(FirmwareRead_In)),
    /* outSize       */     (UINT16)(sizeof(FirmwareRead_Out)),
    /* offsetOfTypes */     offsetof(FirmwareRead_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {UINT32_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_MAX_BUFFER_P_MARSHAL,
                             END_OF_LIST}
};

#define _FirmwareReadDataAddress (&_FirmwareReadData)
#else
#define _FirmwareReadDataAddress 0
#endif // CC_FirmwareRead

#if CC_ContextSave

#include "ContextSave_fp.h"

typedef TPM_RC  (ContextSave_Entry)(
    ContextSave_In              *in,
    ContextSave_Out             *out
);

typedef const struct {
    ContextSave_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} ContextSave_COMMAND_DESCRIPTOR_t;

ContextSave_COMMAND_DESCRIPTOR_t _ContextSaveData = {
    /* entry         */     &TPM2_ContextSave,
    /* inSize        */     (UINT16)(sizeof(ContextSave_In)),
    /* outSize       */     (UINT16)(sizeof(ContextSave_Out)),
    /* offsetOfTypes */     offsetof(ContextSave_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_DH_CONTEXT_H_UNMARSHAL,
                             END_OF_LIST,
                             TPMS_CONTEXT_P_MARSHAL,
                             END_OF_LIST}
};

#define _ContextSaveDataAddress (&_ContextSaveData)
#else
#define _ContextSaveDataAddress 0
#endif // CC_ContextSave

#if CC_ContextLoad

#include "ContextLoad_fp.h"

typedef TPM_RC  (ContextLoad_Entry)(
    ContextLoad_In              *in,
    ContextLoad_Out             *out
);

typedef const struct {
    ContextLoad_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} ContextLoad_COMMAND_DESCRIPTOR_t;

ContextLoad_COMMAND_DESCRIPTOR_t _ContextLoadData = {
    /* entry         */     &TPM2_ContextLoad,
    /* inSize        */     (UINT16)(sizeof(ContextLoad_In)),
    /* outSize       */     (UINT16)(sizeof(ContextLoad_Out)),
    /* offsetOfTypes */     offsetof(ContextLoad_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMS_CONTEXT_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMI_DH_CONTEXT_H_MARSHAL,
                             END_OF_LIST}
};

#define _ContextLoadDataAddress (&_ContextLoadData)
#else
#define _ContextLoadDataAddress 0
#endif // CC_ContextLoad

#if CC_FlushContext

#include "FlushContext_fp.h"

typedef TPM_RC  (FlushContext_Entry)(
    FlushContext_In             *in
);

typedef const struct {
    FlushContext_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} FlushContext_COMMAND_DESCRIPTOR_t;

FlushContext_COMMAND_DESCRIPTOR_t _FlushContextData = {
    /* entry         */     &TPM2_FlushContext,
    /* inSize        */     (UINT16)(sizeof(FlushContext_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(FlushContext_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMI_DH_CONTEXT_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _FlushContextDataAddress (&_FlushContextData)
#else
#define _FlushContextDataAddress 0
#endif // CC_FlushContext

#if CC_EvictControl

#include "EvictControl_fp.h"

typedef TPM_RC  (EvictControl_Entry)(
    EvictControl_In             *in
);

typedef const struct {
    EvictControl_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} EvictControl_COMMAND_DESCRIPTOR_t;

EvictControl_COMMAND_DESCRIPTOR_t _EvictControlData = {
    /* entry         */     &TPM2_EvictControl,
    /* inSize        */     (UINT16)(sizeof(EvictControl_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(EvictControl_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(EvictControl_In, objectHandle)),
                             (UINT16)(offsetof(EvictControl_In, persistentHandle))},
    /* types         */     {TPMI_RH_PROVISION_H_UNMARSHAL,
                             TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_DH_PERSISTENT_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _EvictControlDataAddress (&_EvictControlData)
#else
#define _EvictControlDataAddress 0
#endif // CC_EvictControl

#if CC_ReadClock

#include "ReadClock_fp.h"

typedef TPM_RC  (ReadClock_Entry)(
    ReadClock_Out               *out
);

typedef const struct {
    ReadClock_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} ReadClock_COMMAND_DESCRIPTOR_t;

ReadClock_COMMAND_DESCRIPTOR_t _ReadClockData = {
    /* entry         */     &TPM2_ReadClock,
    /* inSize        */     0,
    /* outSize       */     (UINT16)(sizeof(ReadClock_Out)),
    /* offsetOfTypes */     offsetof(ReadClock_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {END_OF_LIST,
                             TPMS_TIME_INFO_P_MARSHAL,
                             END_OF_LIST}
};

#define _ReadClockDataAddress (&_ReadClockData)
#else
#define _ReadClockDataAddress 0
#endif // CC_ReadClock

#if CC_ClockSet

#include "ClockSet_fp.h"

typedef TPM_RC  (ClockSet_Entry)(
    ClockSet_In                 *in
);

typedef const struct {
    ClockSet_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} ClockSet_COMMAND_DESCRIPTOR_t;

ClockSet_COMMAND_DESCRIPTOR_t _ClockSetData = {
    /* entry         */     &TPM2_ClockSet,
    /* inSize        */     (UINT16)(sizeof(ClockSet_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(ClockSet_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ClockSet_In, newTime))},
    /* types         */     {TPMI_RH_PROVISION_H_UNMARSHAL,
                             UINT64_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ClockSetDataAddress (&_ClockSetData)
#else
#define _ClockSetDataAddress 0
#endif // CC_ClockSet

#if CC_ClockRateAdjust

#include "ClockRateAdjust_fp.h"

typedef TPM_RC  (ClockRateAdjust_Entry)(
    ClockRateAdjust_In          *in
);

typedef const struct {
    ClockRateAdjust_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} ClockRateAdjust_COMMAND_DESCRIPTOR_t;

ClockRateAdjust_COMMAND_DESCRIPTOR_t _ClockRateAdjustData = {
    /* entry         */     &TPM2_ClockRateAdjust,
    /* inSize        */     (UINT16)(sizeof(ClockRateAdjust_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(ClockRateAdjust_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(ClockRateAdjust_In, rateAdjust))},
    /* types         */     {TPMI_RH_PROVISION_H_UNMARSHAL,
                             TPM_CLOCK_ADJUST_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _ClockRateAdjustDataAddress (&_ClockRateAdjustData)
#else
#define _ClockRateAdjustDataAddress 0
#endif // CC_ClockRateAdjust

#if CC_GetCapability

#include "GetCapability_fp.h"

typedef TPM_RC  (GetCapability_Entry)(
    GetCapability_In            *in,
    GetCapability_Out           *out
);

typedef const struct {
    GetCapability_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} GetCapability_COMMAND_DESCRIPTOR_t;

GetCapability_COMMAND_DESCRIPTOR_t _GetCapabilityData = {
    /* entry         */     &TPM2_GetCapability,
    /* inSize        */     (UINT16)(sizeof(GetCapability_In)),
    /* outSize       */     (UINT16)(sizeof(GetCapability_Out)),
    /* offsetOfTypes */     offsetof(GetCapability_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(GetCapability_In, property)),
                             (UINT16)(offsetof(GetCapability_In, propertyCount)),
                             (UINT16)(offsetof(GetCapability_Out, capabilityData))},
    /* types         */     {TPM_CAP_P_UNMARSHAL,
                             UINT32_P_UNMARSHAL,
                             UINT32_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMI_YES_NO_P_MARSHAL,
                             TPMS_CAPABILITY_DATA_P_MARSHAL,
                             END_OF_LIST}
};

#define _GetCapabilityDataAddress (&_GetCapabilityData)
#else
#define _GetCapabilityDataAddress 0
#endif // CC_GetCapability

#if CC_TestParms

#include "TestParms_fp.h"

typedef TPM_RC  (TestParms_Entry)(
    TestParms_In                *in
);

typedef const struct {
    TestParms_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[3];
} TestParms_COMMAND_DESCRIPTOR_t;

TestParms_COMMAND_DESCRIPTOR_t _TestParmsData = {
    /* entry         */     &TPM2_TestParms,
    /* inSize        */     (UINT16)(sizeof(TestParms_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(TestParms_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPMT_PUBLIC_PARMS_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _TestParmsDataAddress (&_TestParmsData)
#else
#define _TestParmsDataAddress 0
#endif // CC_TestParms

#if CC_NV_DefineSpace

#include "NV_DefineSpace_fp.h"

typedef TPM_RC  (NV_DefineSpace_Entry)(
    NV_DefineSpace_In           *in
);

typedef const struct {
    NV_DefineSpace_Entry    *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} NV_DefineSpace_COMMAND_DESCRIPTOR_t;

NV_DefineSpace_COMMAND_DESCRIPTOR_t _NV_DefineSpaceData = {
    /* entry         */     &TPM2_NV_DefineSpace,
    /* inSize        */     (UINT16)(sizeof(NV_DefineSpace_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_DefineSpace_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_DefineSpace_In, auth)),
                             (UINT16)(offsetof(NV_DefineSpace_In, publicInfo))},
    /* types         */     {TPMI_RH_PROVISION_H_UNMARSHAL,
                             TPM2B_AUTH_P_UNMARSHAL,
                             TPM2B_NV_PUBLIC_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_DefineSpaceDataAddress (&_NV_DefineSpaceData)
#else
#define _NV_DefineSpaceDataAddress 0
#endif // CC_NV_DefineSpace

#if CC_NV_UndefineSpace

#include "NV_UndefineSpace_fp.h"

typedef TPM_RC  (NV_UndefineSpace_Entry)(
    NV_UndefineSpace_In         *in
);

typedef const struct {
    NV_UndefineSpace_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} NV_UndefineSpace_COMMAND_DESCRIPTOR_t;

NV_UndefineSpace_COMMAND_DESCRIPTOR_t _NV_UndefineSpaceData = {
    /* entry         */     &TPM2_NV_UndefineSpace,
    /* inSize        */     (UINT16)(sizeof(NV_UndefineSpace_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_UndefineSpace_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_UndefineSpace_In, nvIndex))},
    /* types         */     {TPMI_RH_PROVISION_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_UndefineSpaceDataAddress (&_NV_UndefineSpaceData)
#else
#define _NV_UndefineSpaceDataAddress 0
#endif // CC_NV_UndefineSpace

#if CC_NV_UndefineSpaceSpecial

#include "NV_UndefineSpaceSpecial_fp.h"

typedef TPM_RC  (NV_UndefineSpaceSpecial_Entry)(
    NV_UndefineSpaceSpecial_In          *in
);

typedef const struct {
    NV_UndefineSpaceSpecial_Entry   *entry;
    UINT16                          inSize;
    UINT16                          outSize;
    UINT16                          offsetOfTypes;
    UINT16                          paramOffsets[1];
    BYTE                            types[4];
} NV_UndefineSpaceSpecial_COMMAND_DESCRIPTOR_t;

NV_UndefineSpaceSpecial_COMMAND_DESCRIPTOR_t _NV_UndefineSpaceSpecialData = {
    /* entry         */             &TPM2_NV_UndefineSpaceSpecial,
    /* inSize        */             (UINT16)(sizeof(NV_UndefineSpaceSpecial_In)),
    /* outSize       */             0,
    /* offsetOfTypes */             offsetof(NV_UndefineSpaceSpecial_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */             {(UINT16)(offsetof(NV_UndefineSpaceSpecial_In, platform))},
    /* types         */             {TPMI_RH_NV_INDEX_H_UNMARSHAL,
                                     TPMI_RH_PLATFORM_H_UNMARSHAL,
                                     END_OF_LIST,
                                     END_OF_LIST}
};

#define _NV_UndefineSpaceSpecialDataAddress (&_NV_UndefineSpaceSpecialData)
#else
#define _NV_UndefineSpaceSpecialDataAddress 0
#endif // CC_NV_UndefineSpaceSpecial

#if CC_NV_ReadPublic

#include "NV_ReadPublic_fp.h"

typedef TPM_RC  (NV_ReadPublic_Entry)(
    NV_ReadPublic_In            *in,
    NV_ReadPublic_Out           *out
);

typedef const struct {
    NV_ReadPublic_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[5];
} NV_ReadPublic_COMMAND_DESCRIPTOR_t;

NV_ReadPublic_COMMAND_DESCRIPTOR_t _NV_ReadPublicData = {
    /* entry         */     &TPM2_NV_ReadPublic,
    /* inSize        */     (UINT16)(sizeof(NV_ReadPublic_In)),
    /* outSize       */     (UINT16)(sizeof(NV_ReadPublic_Out)),
    /* offsetOfTypes */     offsetof(NV_ReadPublic_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_ReadPublic_Out, nvName))},
    /* types         */     {TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_NV_PUBLIC_P_MARSHAL,
                             TPM2B_NAME_P_MARSHAL,
                             END_OF_LIST}
};

#define _NV_ReadPublicDataAddress (&_NV_ReadPublicData)
#else
#define _NV_ReadPublicDataAddress 0
#endif // CC_NV_ReadPublic

#if CC_NV_Write

#include "NV_Write_fp.h"

typedef TPM_RC  (NV_Write_Entry)(
    NV_Write_In                 *in
);

typedef const struct {
    NV_Write_Entry          *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[6];
} NV_Write_COMMAND_DESCRIPTOR_t;

NV_Write_COMMAND_DESCRIPTOR_t _NV_WriteData = {
    /* entry         */     &TPM2_NV_Write,
    /* inSize        */     (UINT16)(sizeof(NV_Write_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_Write_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_Write_In, nvIndex)),
                             (UINT16)(offsetof(NV_Write_In, data)),
                             (UINT16)(offsetof(NV_Write_In, offset))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             TPM2B_MAX_NV_BUFFER_P_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_WriteDataAddress (&_NV_WriteData)
#else
#define _NV_WriteDataAddress 0
#endif // CC_NV_Write

#if CC_NV_Increment

#include "NV_Increment_fp.h"

typedef TPM_RC  (NV_Increment_Entry)(
    NV_Increment_In             *in
);

typedef const struct {
    NV_Increment_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} NV_Increment_COMMAND_DESCRIPTOR_t;

NV_Increment_COMMAND_DESCRIPTOR_t _NV_IncrementData = {
    /* entry         */     &TPM2_NV_Increment,
    /* inSize        */     (UINT16)(sizeof(NV_Increment_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_Increment_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_Increment_In, nvIndex))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_IncrementDataAddress (&_NV_IncrementData)
#else
#define _NV_IncrementDataAddress 0
#endif // CC_NV_Increment

#if CC_NV_Extend

#include "NV_Extend_fp.h"

typedef TPM_RC  (NV_Extend_Entry)(
    NV_Extend_In                *in
);

typedef const struct {
    NV_Extend_Entry         *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} NV_Extend_COMMAND_DESCRIPTOR_t;

NV_Extend_COMMAND_DESCRIPTOR_t _NV_ExtendData = {
    /* entry         */     &TPM2_NV_Extend,
    /* inSize        */     (UINT16)(sizeof(NV_Extend_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_Extend_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_Extend_In, nvIndex)),
                             (UINT16)(offsetof(NV_Extend_In, data))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             TPM2B_MAX_NV_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_ExtendDataAddress (&_NV_ExtendData)
#else
#define _NV_ExtendDataAddress 0
#endif // CC_NV_Extend

#if CC_NV_SetBits

#include "NV_SetBits_fp.h"

typedef TPM_RC  (NV_SetBits_Entry)(
    NV_SetBits_In               *in
);

typedef const struct {
    NV_SetBits_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[2];
    BYTE                    types[5];
} NV_SetBits_COMMAND_DESCRIPTOR_t;

NV_SetBits_COMMAND_DESCRIPTOR_t _NV_SetBitsData = {
    /* entry         */     &TPM2_NV_SetBits,
    /* inSize        */     (UINT16)(sizeof(NV_SetBits_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_SetBits_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_SetBits_In, nvIndex)),
                             (UINT16)(offsetof(NV_SetBits_In, bits))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             UINT64_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_SetBitsDataAddress (&_NV_SetBitsData)
#else
#define _NV_SetBitsDataAddress 0
#endif // CC_NV_SetBits

#if CC_NV_WriteLock

#include "NV_WriteLock_fp.h"

typedef TPM_RC  (NV_WriteLock_Entry)(
    NV_WriteLock_In             *in
);

typedef const struct {
    NV_WriteLock_Entry      *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} NV_WriteLock_COMMAND_DESCRIPTOR_t;

NV_WriteLock_COMMAND_DESCRIPTOR_t _NV_WriteLockData = {
    /* entry         */     &TPM2_NV_WriteLock,
    /* inSize        */     (UINT16)(sizeof(NV_WriteLock_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_WriteLock_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_WriteLock_In, nvIndex))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_WriteLockDataAddress (&_NV_WriteLockData)
#else
#define _NV_WriteLockDataAddress 0
#endif // CC_NV_WriteLock

#if CC_NV_GlobalWriteLock

#include "NV_GlobalWriteLock_fp.h"

typedef TPM_RC  (NV_GlobalWriteLock_Entry)(
    NV_GlobalWriteLock_In           *in
);

typedef const struct {
    NV_GlobalWriteLock_Entry    *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    BYTE                        types[3];
} NV_GlobalWriteLock_COMMAND_DESCRIPTOR_t;

NV_GlobalWriteLock_COMMAND_DESCRIPTOR_t _NV_GlobalWriteLockData = {
    /* entry         */         &TPM2_NV_GlobalWriteLock,
    /* inSize        */         (UINT16)(sizeof(NV_GlobalWriteLock_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(NV_GlobalWriteLock_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         // No parameter offsets;
    /* types         */         {TPMI_RH_PROVISION_H_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _NV_GlobalWriteLockDataAddress (&_NV_GlobalWriteLockData)
#else
#define _NV_GlobalWriteLockDataAddress 0
#endif // CC_NV_GlobalWriteLock

#if CC_NV_Read

#include "NV_Read_fp.h"

typedef TPM_RC  (NV_Read_Entry)(
    NV_Read_In                  *in,
    NV_Read_Out                 *out
);

typedef const struct {
    NV_Read_Entry           *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} NV_Read_COMMAND_DESCRIPTOR_t;

NV_Read_COMMAND_DESCRIPTOR_t _NV_ReadData = {
    /* entry         */     &TPM2_NV_Read,
    /* inSize        */     (UINT16)(sizeof(NV_Read_In)),
    /* outSize       */     (UINT16)(sizeof(NV_Read_Out)),
    /* offsetOfTypes */     offsetof(NV_Read_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_Read_In, nvIndex)),
                             (UINT16)(offsetof(NV_Read_In, size)),
                             (UINT16)(offsetof(NV_Read_In, offset))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_MAX_NV_BUFFER_P_MARSHAL,
                             END_OF_LIST}
};

#define _NV_ReadDataAddress (&_NV_ReadData)
#else
#define _NV_ReadDataAddress 0
#endif // CC_NV_Read

#if CC_NV_ReadLock

#include "NV_ReadLock_fp.h"

typedef TPM_RC  (NV_ReadLock_Entry)(
    NV_ReadLock_In              *in
);

typedef const struct {
    NV_ReadLock_Entry       *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} NV_ReadLock_COMMAND_DESCRIPTOR_t;

NV_ReadLock_COMMAND_DESCRIPTOR_t _NV_ReadLockData = {
    /* entry         */     &TPM2_NV_ReadLock,
    /* inSize        */     (UINT16)(sizeof(NV_ReadLock_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_ReadLock_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_ReadLock_In, nvIndex))},
    /* types         */     {TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_ReadLockDataAddress (&_NV_ReadLockData)
#else
#define _NV_ReadLockDataAddress 0
#endif // CC_NV_ReadLock

#if CC_NV_ChangeAuth

#include "NV_ChangeAuth_fp.h"

typedef TPM_RC  (NV_ChangeAuth_Entry)(
    NV_ChangeAuth_In            *in
);

typedef const struct {
    NV_ChangeAuth_Entry     *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[1];
    BYTE                    types[4];
} NV_ChangeAuth_COMMAND_DESCRIPTOR_t;

NV_ChangeAuth_COMMAND_DESCRIPTOR_t _NV_ChangeAuthData = {
    /* entry         */     &TPM2_NV_ChangeAuth,
    /* inSize        */     (UINT16)(sizeof(NV_ChangeAuth_In)),
    /* outSize       */     0,
    /* offsetOfTypes */     offsetof(NV_ChangeAuth_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_ChangeAuth_In, newAuth))},
    /* types         */     {TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             TPM2B_AUTH_P_UNMARSHAL,
                             END_OF_LIST,
                             END_OF_LIST}
};

#define _NV_ChangeAuthDataAddress (&_NV_ChangeAuthData)
#else
#define _NV_ChangeAuthDataAddress 0
#endif // CC_NV_ChangeAuth

#if CC_NV_Certify

#include "NV_Certify_fp.h"

typedef TPM_RC  (NV_Certify_Entry)(
    NV_Certify_In               *in,
    NV_Certify_Out              *out
);

typedef const struct {
    NV_Certify_Entry        *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[7];
    BYTE                    types[11];
} NV_Certify_COMMAND_DESCRIPTOR_t;

NV_Certify_COMMAND_DESCRIPTOR_t _NV_CertifyData = {
    /* entry         */     &TPM2_NV_Certify,
    /* inSize        */     (UINT16)(sizeof(NV_Certify_In)),
    /* outSize       */     (UINT16)(sizeof(NV_Certify_Out)),
    /* offsetOfTypes */     offsetof(NV_Certify_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(NV_Certify_In, authHandle)),
                             (UINT16)(offsetof(NV_Certify_In, nvIndex)),
                             (UINT16)(offsetof(NV_Certify_In, qualifyingData)),
                             (UINT16)(offsetof(NV_Certify_In, inScheme)),
                             (UINT16)(offsetof(NV_Certify_In, size)),
                             (UINT16)(offsetof(NV_Certify_In, offset)),
                             (UINT16)(offsetof(NV_Certify_Out, signature))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL + ADD_FLAG,
                             TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_NV_INDEX_H_UNMARSHAL,
                             TPM2B_DATA_P_UNMARSHAL,
                             TPMT_SIG_SCHEME_P_UNMARSHAL + ADD_FLAG,
                             UINT16_P_UNMARSHAL,
                             UINT16_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_ATTEST_P_MARSHAL,
                             TPMT_SIGNATURE_P_MARSHAL,
                             END_OF_LIST}
};

#define _NV_CertifyDataAddress (&_NV_CertifyData)
#else
#define _NV_CertifyDataAddress 0
#endif // CC_NV_Certify

#if CC_AC_GetCapability

#include "AC_GetCapability_fp.h"

typedef TPM_RC  (AC_GetCapability_Entry)(
    AC_GetCapability_In         *in,
    AC_GetCapability_Out        *out
);

typedef const struct {
    AC_GetCapability_Entry  *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} AC_GetCapability_COMMAND_DESCRIPTOR_t;

AC_GetCapability_COMMAND_DESCRIPTOR_t _AC_GetCapabilityData = {
    /* entry         */     &TPM2_AC_GetCapability,
    /* inSize        */     (UINT16)(sizeof(AC_GetCapability_In)),
    /* outSize       */     (UINT16)(sizeof(AC_GetCapability_Out)),
    /* offsetOfTypes */     offsetof(AC_GetCapability_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(AC_GetCapability_In, capability)),
                             (UINT16)(offsetof(AC_GetCapability_In, count)),
                             (UINT16)(offsetof(AC_GetCapability_Out, capabilitiesData))},
    /* types         */     {TPMI_RH_AC_H_UNMARSHAL,
                             TPM_AT_P_UNMARSHAL,
                             UINT32_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMI_YES_NO_P_MARSHAL,
                             TPML_AC_CAPABILITIES_P_MARSHAL,
                             END_OF_LIST}
};

#define _AC_GetCapabilityDataAddress (&_AC_GetCapabilityData)
#else
#define _AC_GetCapabilityDataAddress 0
#endif // CC_AC_GetCapability

#if CC_AC_Send

#include "AC_Send_fp.h"

typedef TPM_RC  (AC_Send_Entry)(
    AC_Send_In                  *in,
    AC_Send_Out                 *out
);

typedef const struct {
    AC_Send_Entry           *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    UINT16                  paramOffsets[3];
    BYTE                    types[7];
} AC_Send_COMMAND_DESCRIPTOR_t;

AC_Send_COMMAND_DESCRIPTOR_t _AC_SendData = {
    /* entry         */     &TPM2_AC_Send,
    /* inSize        */     (UINT16)(sizeof(AC_Send_In)),
    /* outSize       */     (UINT16)(sizeof(AC_Send_Out)),
    /* offsetOfTypes */     offsetof(AC_Send_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     {(UINT16)(offsetof(AC_Send_In, authHandle)),
                             (UINT16)(offsetof(AC_Send_In, ac)),
                             (UINT16)(offsetof(AC_Send_In, acDataIn))},
    /* types         */     {TPMI_DH_OBJECT_H_UNMARSHAL,
                             TPMI_RH_NV_AUTH_H_UNMARSHAL,
                             TPMI_RH_AC_H_UNMARSHAL,
                             TPM2B_MAX_BUFFER_P_UNMARSHAL,
                             END_OF_LIST,
                             TPMS_AC_OUTPUT_P_MARSHAL,
                             END_OF_LIST}
};

#define _AC_SendDataAddress (&_AC_SendData)
#else
#define _AC_SendDataAddress 0
#endif // CC_AC_Send

#if CC_Policy_AC_SendSelect

#include "Policy_AC_SendSelect_fp.h"

typedef TPM_RC  (Policy_AC_SendSelect_Entry)(
    Policy_AC_SendSelect_In         *in
);

typedef const struct {
    Policy_AC_SendSelect_Entry  *entry;
    UINT16                      inSize;
    UINT16                      outSize;
    UINT16                      offsetOfTypes;
    UINT16                      paramOffsets[4];
    BYTE                        types[7];
} Policy_AC_SendSelect_COMMAND_DESCRIPTOR_t;

Policy_AC_SendSelect_COMMAND_DESCRIPTOR_t _Policy_AC_SendSelectData = {
    /* entry         */         &TPM2_Policy_AC_SendSelect,
    /* inSize        */         (UINT16)(sizeof(Policy_AC_SendSelect_In)),
    /* outSize       */         0,
    /* offsetOfTypes */         offsetof(Policy_AC_SendSelect_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */         {(UINT16)(offsetof(Policy_AC_SendSelect_In, objectName)),
                                 (UINT16)(offsetof(Policy_AC_SendSelect_In, authHandleName)),
                                 (UINT16)(offsetof(Policy_AC_SendSelect_In, acName)),
                                 (UINT16)(offsetof(Policy_AC_SendSelect_In, includeObject))},
    /* types         */         {TPMI_SH_POLICY_H_UNMARSHAL,
                                 TPM2B_NAME_P_UNMARSHAL,
                                 TPM2B_NAME_P_UNMARSHAL,
                                 TPM2B_NAME_P_UNMARSHAL,
                                 TPMI_YES_NO_P_UNMARSHAL,
                                 END_OF_LIST,
                                 END_OF_LIST}
};

#define _Policy_AC_SendSelectDataAddress (&_Policy_AC_SendSelectData)
#else
#define _Policy_AC_SendSelectDataAddress 0
#endif // CC_Policy_AC_SendSelect

#if CC_Vendor_TCG_Test

#include "Vendor_TCG_Test_fp.h"

typedef TPM_RC  (Vendor_TCG_Test_Entry)(
    Vendor_TCG_Test_In          *in,
    Vendor_TCG_Test_Out         *out
);

typedef const struct {
    Vendor_TCG_Test_Entry   *entry;
    UINT16                  inSize;
    UINT16                  outSize;
    UINT16                  offsetOfTypes;
    BYTE                    types[4];
} Vendor_TCG_Test_COMMAND_DESCRIPTOR_t;

Vendor_TCG_Test_COMMAND_DESCRIPTOR_t _Vendor_TCG_TestData = {
    /* entry         */     &TPM2_Vendor_TCG_Test,
    /* inSize        */     (UINT16)(sizeof(Vendor_TCG_Test_In)),
    /* outSize       */     (UINT16)(sizeof(Vendor_TCG_Test_Out)),
    /* offsetOfTypes */     offsetof(Vendor_TCG_Test_COMMAND_DESCRIPTOR_t, types),
    /* offsets       */     // No parameter offsets;
    /* types         */     {TPM2B_DATA_P_UNMARSHAL,
                             END_OF_LIST,
                             TPM2B_DATA_P_MARSHAL,
                             END_OF_LIST}
};

#define _Vendor_TCG_TestDataAddress (&_Vendor_TCG_TestData)
#else
#define _Vendor_TCG_TestDataAddress 0
#endif // CC_Vendor_TCG_Test

COMMAND_DESCRIPTOR_t *s_CommandDataArray[] = {
#if (PAD_LIST || CC_NV_UndefineSpaceSpecial)
        (COMMAND_DESCRIPTOR_t *)_NV_UndefineSpaceSpecialDataAddress,
#endif // CC_NV_UndefineSpaceSpecial
#if (PAD_LIST || CC_EvictControl)
        (COMMAND_DESCRIPTOR_t *)_EvictControlDataAddress,
#endif // CC_EvictControl
#if (PAD_LIST || CC_HierarchyControl)
        (COMMAND_DESCRIPTOR_t *)_HierarchyControlDataAddress,
#endif // CC_HierarchyControl
#if (PAD_LIST || CC_NV_UndefineSpace)
        (COMMAND_DESCRIPTOR_t *)_NV_UndefineSpaceDataAddress,
#endif // CC_NV_UndefineSpace
#if (PAD_LIST)
        (COMMAND_DESCRIPTOR_t *)0,
#endif //
#if (PAD_LIST || CC_ChangeEPS)
        (COMMAND_DESCRIPTOR_t *)_ChangeEPSDataAddress,
#endif // CC_ChangeEPS
#if (PAD_LIST || CC_ChangePPS)
        (COMMAND_DESCRIPTOR_t *)_ChangePPSDataAddress,
#endif // CC_ChangePPS
#if (PAD_LIST || CC_Clear)
        (COMMAND_DESCRIPTOR_t *)_ClearDataAddress,
#endif // CC_Clear
#if (PAD_LIST || CC_ClearControl)
        (COMMAND_DESCRIPTOR_t *)_ClearControlDataAddress,
#endif // CC_ClearControl
#if (PAD_LIST || CC_ClockSet)
        (COMMAND_DESCRIPTOR_t *)_ClockSetDataAddress,
#endif // CC_ClockSet
#if (PAD_LIST || CC_HierarchyChangeAuth)
        (COMMAND_DESCRIPTOR_t *)_HierarchyChangeAuthDataAddress,
#endif // CC_HierarchyChangeAuth
#if (PAD_LIST || CC_NV_DefineSpace)
        (COMMAND_DESCRIPTOR_t *)_NV_DefineSpaceDataAddress,
#endif // CC_NV_DefineSpace
#if (PAD_LIST || CC_PCR_Allocate)
        (COMMAND_DESCRIPTOR_t *)_PCR_AllocateDataAddress,
#endif // CC_PCR_Allocate
#if (PAD_LIST || CC_PCR_SetAuthPolicy)
        (COMMAND_DESCRIPTOR_t *)_PCR_SetAuthPolicyDataAddress,
#endif // CC_PCR_SetAuthPolicy
#if (PAD_LIST || CC_PP_Commands)
        (COMMAND_DESCRIPTOR_t *)_PP_CommandsDataAddress,
#endif // CC_PP_Commands
#if (PAD_LIST || CC_SetPrimaryPolicy)
        (COMMAND_DESCRIPTOR_t *)_SetPrimaryPolicyDataAddress,
#endif // CC_SetPrimaryPolicy
#if (PAD_LIST || CC_FieldUpgradeStart)
        (COMMAND_DESCRIPTOR_t *)_FieldUpgradeStartDataAddress,
#endif // CC_FieldUpgradeStart
#if (PAD_LIST || CC_ClockRateAdjust)
        (COMMAND_DESCRIPTOR_t *)_ClockRateAdjustDataAddress,
#endif // CC_ClockRateAdjust
#if (PAD_LIST || CC_CreatePrimary)
        (COMMAND_DESCRIPTOR_t *)_CreatePrimaryDataAddress,
#endif // CC_CreatePrimary
#if (PAD_LIST || CC_NV_GlobalWriteLock)
        (COMMAND_DESCRIPTOR_t *)_NV_GlobalWriteLockDataAddress,
#endif // CC_NV_GlobalWriteLock
#if (PAD_LIST || CC_GetCommandAuditDigest)
        (COMMAND_DESCRIPTOR_t *)_GetCommandAuditDigestDataAddress,
#endif // CC_GetCommandAuditDigest
#if (PAD_LIST || CC_NV_Increment)
        (COMMAND_DESCRIPTOR_t *)_NV_IncrementDataAddress,
#endif // CC_NV_Increment
#if (PAD_LIST || CC_NV_SetBits)
        (COMMAND_DESCRIPTOR_t *)_NV_SetBitsDataAddress,
#endif // CC_NV_SetBits
#if (PAD_LIST || CC_NV_Extend)
        (COMMAND_DESCRIPTOR_t *)_NV_ExtendDataAddress,
#endif // CC_NV_Extend
#if (PAD_LIST || CC_NV_Write)
        (COMMAND_DESCRIPTOR_t *)_NV_WriteDataAddress,
#endif // CC_NV_Write
#if (PAD_LIST || CC_NV_WriteLock)
        (COMMAND_DESCRIPTOR_t *)_NV_WriteLockDataAddress,
#endif // CC_NV_WriteLock
#if (PAD_LIST || CC_DictionaryAttackLockReset)
        (COMMAND_DESCRIPTOR_t *)_DictionaryAttackLockResetDataAddress,
#endif // CC_DictionaryAttackLockReset
#if (PAD_LIST || CC_DictionaryAttackParameters)
        (COMMAND_DESCRIPTOR_t *)_DictionaryAttackParametersDataAddress,
#endif // CC_DictionaryAttackParameters
#if (PAD_LIST || CC_NV_ChangeAuth)
        (COMMAND_DESCRIPTOR_t *)_NV_ChangeAuthDataAddress,
#endif // CC_NV_ChangeAuth
#if (PAD_LIST || CC_PCR_Event)
        (COMMAND_DESCRIPTOR_t *)_PCR_EventDataAddress,
#endif // CC_PCR_Event
#if (PAD_LIST || CC_PCR_Reset)
        (COMMAND_DESCRIPTOR_t *)_PCR_ResetDataAddress,
#endif // CC_PCR_Reset
#if (PAD_LIST || CC_SequenceComplete)
        (COMMAND_DESCRIPTOR_t *)_SequenceCompleteDataAddress,
#endif // CC_SequenceComplete
#if (PAD_LIST || CC_SetAlgorithmSet)
        (COMMAND_DESCRIPTOR_t *)_SetAlgorithmSetDataAddress,
#endif // CC_SetAlgorithmSet
#if (PAD_LIST || CC_SetCommandCodeAuditStatus)
        (COMMAND_DESCRIPTOR_t *)_SetCommandCodeAuditStatusDataAddress,
#endif // CC_SetCommandCodeAuditStatus
#if (PAD_LIST || CC_FieldUpgradeData)
        (COMMAND_DESCRIPTOR_t *)_FieldUpgradeDataDataAddress,
#endif // CC_FieldUpgradeData
#if (PAD_LIST || CC_IncrementalSelfTest)
        (COMMAND_DESCRIPTOR_t *)_IncrementalSelfTestDataAddress,
#endif // CC_IncrementalSelfTest
#if (PAD_LIST || CC_SelfTest)
        (COMMAND_DESCRIPTOR_t *)_SelfTestDataAddress,
#endif // CC_SelfTest
#if (PAD_LIST || CC_Startup)
        (COMMAND_DESCRIPTOR_t *)_StartupDataAddress,
#endif // CC_Startup
#if (PAD_LIST || CC_Shutdown)
        (COMMAND_DESCRIPTOR_t *)_ShutdownDataAddress,
#endif // CC_Shutdown
#if (PAD_LIST || CC_StirRandom)
        (COMMAND_DESCRIPTOR_t *)_StirRandomDataAddress,
#endif // CC_StirRandom
#if (PAD_LIST || CC_ActivateCredential)
        (COMMAND_DESCRIPTOR_t *)_ActivateCredentialDataAddress,
#endif // CC_ActivateCredential
#if (PAD_LIST || CC_Certify)
        (COMMAND_DESCRIPTOR_t *)_CertifyDataAddress,
#endif // CC_Certify
#if (PAD_LIST || CC_PolicyNV)
        (COMMAND_DESCRIPTOR_t *)_PolicyNVDataAddress,
#endif // CC_PolicyNV
#if (PAD_LIST || CC_CertifyCreation)
        (COMMAND_DESCRIPTOR_t *)_CertifyCreationDataAddress,
#endif // CC_CertifyCreation
#if (PAD_LIST || CC_Duplicate)
        (COMMAND_DESCRIPTOR_t *)_DuplicateDataAddress,
#endif // CC_Duplicate
#if (PAD_LIST || CC_GetTime)
        (COMMAND_DESCRIPTOR_t *)_GetTimeDataAddress,
#endif // CC_GetTime
#if (PAD_LIST || CC_GetSessionAuditDigest)
        (COMMAND_DESCRIPTOR_t *)_GetSessionAuditDigestDataAddress,
#endif // CC_GetSessionAuditDigest
#if (PAD_LIST || CC_NV_Read)
        (COMMAND_DESCRIPTOR_t *)_NV_ReadDataAddress,
#endif // CC_NV_Read
#if (PAD_LIST || CC_NV_ReadLock)
        (COMMAND_DESCRIPTOR_t *)_NV_ReadLockDataAddress,
#endif // CC_NV_ReadLock
#if (PAD_LIST || CC_ObjectChangeAuth)
        (COMMAND_DESCRIPTOR_t *)_ObjectChangeAuthDataAddress,
#endif // CC_ObjectChangeAuth
#if (PAD_LIST || CC_PolicySecret)
        (COMMAND_DESCRIPTOR_t *)_PolicySecretDataAddress,
#endif // CC_PolicySecret
#if (PAD_LIST || CC_Rewrap)
        (COMMAND_DESCRIPTOR_t *)_RewrapDataAddress,
#endif // CC_Rewrap
#if (PAD_LIST || CC_Create)
        (COMMAND_DESCRIPTOR_t *)_CreateDataAddress,
#endif // CC_Create
#if (PAD_LIST || CC_ECDH_ZGen)
        (COMMAND_DESCRIPTOR_t *)_ECDH_ZGenDataAddress,
#endif // CC_ECDH_ZGen
#if (PAD_LIST || (CC_HMAC || CC_MAC))
#    if CC_HMAC
        (COMMAND_DESCRIPTOR_t *)_HMACDataAddress,
#    endif
#    if CC_MAC
        (COMMAND_DESCRIPTOR_t *)_MACDataAddress,
#    endif
#    if (CC_HMAC || CC_MAC) > 1
#        error "More than one aliased command defined"
#    endif
#endif // CC_HMAC CC_MAC
#if (PAD_LIST || CC_Import)
        (COMMAND_DESCRIPTOR_t *)_ImportDataAddress,
#endif // CC_Import
#if (PAD_LIST || CC_Load)
        (COMMAND_DESCRIPTOR_t *)_LoadDataAddress,
#endif // CC_Load
#if (PAD_LIST || CC_Quote)
        (COMMAND_DESCRIPTOR_t *)_QuoteDataAddress,
#endif // CC_Quote
#if (PAD_LIST || CC_RSA_Decrypt)
        (COMMAND_DESCRIPTOR_t *)_RSA_DecryptDataAddress,
#endif // CC_RSA_Decrypt
#if (PAD_LIST)
        (COMMAND_DESCRIPTOR_t *)0,
#endif //
#if (PAD_LIST || (CC_HMAC_Start || CC_MAC_Start))
#    if CC_HMAC_Start
        (COMMAND_DESCRIPTOR_t *)_HMAC_StartDataAddress,
#    endif
#    if CC_MAC_Start
        (COMMAND_DESCRIPTOR_t *)_MAC_StartDataAddress,
#    endif
#    if (CC_HMAC_Start || CC_MAC_Start) > 1
#        error "More than one aliased command defined"
#    endif
#endif // CC_HMAC_Start CC_MAC_Start
#if (PAD_LIST || CC_SequenceUpdate)
        (COMMAND_DESCRIPTOR_t *)_SequenceUpdateDataAddress,
#endif // CC_SequenceUpdate
#if (PAD_LIST || CC_Sign)
        (COMMAND_DESCRIPTOR_t *)_SignDataAddress,
#endif // CC_Sign
#if (PAD_LIST || CC_Unseal)
        (COMMAND_DESCRIPTOR_t *)_UnsealDataAddress,
#endif // CC_Unseal
#if (PAD_LIST)
        (COMMAND_DESCRIPTOR_t *)0,
#endif //
#if (PAD_LIST || CC_PolicySigned)
        (COMMAND_DESCRIPTOR_t *)_PolicySignedDataAddress,
#endif // CC_PolicySigned
#if (PAD_LIST || CC_ContextLoad)
        (COMMAND_DESCRIPTOR_t *)_ContextLoadDataAddress,
#endif // CC_ContextLoad
#if (PAD_LIST || CC_ContextSave)
        (COMMAND_DESCRIPTOR_t *)_ContextSaveDataAddress,
#endif // CC_ContextSave
#if (PAD_LIST || CC_ECDH_KeyGen)
        (COMMAND_DESCRIPTOR_t *)_ECDH_KeyGenDataAddress,
#endif // CC_ECDH_KeyGen
#if (PAD_LIST || CC_EncryptDecrypt)
        (COMMAND_DESCRIPTOR_t *)_EncryptDecryptDataAddress,
#endif // CC_EncryptDecrypt
#if (PAD_LIST || CC_FlushContext)
        (COMMAND_DESCRIPTOR_t *)_FlushContextDataAddress,
#endif // CC_FlushContext
#if (PAD_LIST)
        (COMMAND_DESCRIPTOR_t *)0,
#endif //
#if (PAD_LIST || CC_LoadExternal)
        (COMMAND_DESCRIPTOR_t *)_LoadExternalDataAddress,
#endif // CC_LoadExternal
#if (PAD_LIST || CC_MakeCredential)
        (COMMAND_DESCRIPTOR_t *)_MakeCredentialDataAddress,
#endif // CC_MakeCredential
#if (PAD_LIST || CC_NV_ReadPublic)
        (COMMAND_DESCRIPTOR_t *)_NV_ReadPublicDataAddress,
#endif // CC_NV_ReadPublic
#if (PAD_LIST || CC_PolicyAuthorize)
        (COMMAND_DESCRIPTOR_t *)_PolicyAuthorizeDataAddress,
#endif // CC_PolicyAuthorize
#if (PAD_LIST || CC_PolicyAuthValue)
        (COMMAND_DESCRIPTOR_t *)_PolicyAuthValueDataAddress,
#endif // CC_PolicyAuthValue
#if (PAD_LIST || CC_PolicyCommandCode)
        (COMMAND_DESCRIPTOR_t *)_PolicyCommandCodeDataAddress,
#endif // CC_PolicyCommandCode
#if (PAD_LIST || CC_PolicyCounterTimer)
        (COMMAND_DESCRIPTOR_t *)_PolicyCounterTimerDataAddress,
#endif // CC_PolicyCounterTimer
#if (PAD_LIST || CC_PolicyCpHash)
        (COMMAND_DESCRIPTOR_t *)_PolicyCpHashDataAddress,
#endif // CC_PolicyCpHash
#if (PAD_LIST || CC_PolicyLocality)
        (COMMAND_DESCRIPTOR_t *)_PolicyLocalityDataAddress,
#endif // CC_PolicyLocality
#if (PAD_LIST || CC_PolicyNameHash)
        (COMMAND_DESCRIPTOR_t *)_PolicyNameHashDataAddress,
#endif // CC_PolicyNameHash
#if (PAD_LIST || CC_PolicyOR)
        (COMMAND_DESCRIPTOR_t *)_PolicyORDataAddress,
#endif // CC_PolicyOR
#if (PAD_LIST || CC_PolicyTicket)
        (COMMAND_DESCRIPTOR_t *)_PolicyTicketDataAddress,
#endif // CC_PolicyTicket
#if (PAD_LIST || CC_ReadPublic)
        (COMMAND_DESCRIPTOR_t *)_ReadPublicDataAddress,
#endif // CC_ReadPublic
#if (PAD_LIST || CC_RSA_Encrypt)
        (COMMAND_DESCRIPTOR_t *)_RSA_EncryptDataAddress,
#endif // CC_RSA_Encrypt
#if (PAD_LIST)
        (COMMAND_DESCRIPTOR_t *)0,
#endif //
#if (PAD_LIST || CC_StartAuthSession)
        (COMMAND_DESCRIPTOR_t *)_StartAuthSessionDataAddress,
#endif // CC_StartAuthSession
#if (PAD_LIST || CC_VerifySignature)
        (COMMAND_DESCRIPTOR_t *)_VerifySignatureDataAddress,
#endif // CC_VerifySignature
#if (PAD_LIST || CC_ECC_Parameters)
        (COMMAND_DESCRIPTOR_t *)_ECC_ParametersDataAddress,
#endif // CC_ECC_Parameters
#if (PAD_LIST || CC_FirmwareRead)
        (COMMAND_DESCRIPTOR_t *)_FirmwareReadDataAddress,
#endif // CC_FirmwareRead
#if (PAD_LIST || CC_GetCapability)
        (COMMAND_DESCRIPTOR_t *)_GetCapabilityDataAddress,
#endif // CC_GetCapability
#if (PAD_LIST || CC_GetRandom)
        (COMMAND_DESCRIPTOR_t *)_GetRandomDataAddress,
#endif // CC_GetRandom
#if (PAD_LIST || CC_GetTestResult)
        (COMMAND_DESCRIPTOR_t *)_GetTestResultDataAddress,
#endif // CC_GetTestResult
#if (PAD_LIST || CC_Hash)
        (COMMAND_DESCRIPTOR_t *)_HashDataAddress,
#endif // CC_Hash
#if (PAD_LIST || CC_PCR_Read)
        (COMMAND_DESCRIPTOR_t *)_PCR_ReadDataAddress,
#endif // CC_PCR_Read
#if (PAD_LIST || CC_PolicyPCR)
        (COMMAND_DESCRIPTOR_t *)_PolicyPCRDataAddress,
#endif // CC_PolicyPCR
#if (PAD_LIST || CC_PolicyRestart)
        (COMMAND_DESCRIPTOR_t *)_PolicyRestartDataAddress,
#endif // CC_PolicyRestart
#if (PAD_LIST || CC_ReadClock)
        (COMMAND_DESCRIPTOR_t *)_ReadClockDataAddress,
#endif // CC_ReadClock
#if (PAD_LIST || CC_PCR_Extend)
        (COMMAND_DESCRIPTOR_t *)_PCR_ExtendDataAddress,
#endif // CC_PCR_Extend
#if (PAD_LIST || CC_PCR_SetAuthValue)
        (COMMAND_DESCRIPTOR_t *)_PCR_SetAuthValueDataAddress,
#endif // CC_PCR_SetAuthValue
#if (PAD_LIST || CC_NV_Certify)
        (COMMAND_DESCRIPTOR_t *)_NV_CertifyDataAddress,
#endif // CC_NV_Certify
#if (PAD_LIST || CC_EventSequenceComplete)
        (COMMAND_DESCRIPTOR_t *)_EventSequenceCompleteDataAddress,
#endif // CC_EventSequenceComplete
#if (PAD_LIST || CC_HashSequenceStart)
        (COMMAND_DESCRIPTOR_t *)_HashSequenceStartDataAddress,
#endif // CC_HashSequenceStart
#if (PAD_LIST || CC_PolicyPhysicalPresence)
        (COMMAND_DESCRIPTOR_t *)_PolicyPhysicalPresenceDataAddress,
#endif // CC_PolicyPhysicalPresence
#if (PAD_LIST || CC_PolicyDuplicationSelect)
        (COMMAND_DESCRIPTOR_t *)_PolicyDuplicationSelectDataAddress,
#endif // CC_PolicyDuplicationSelect
#if (PAD_LIST || CC_PolicyGetDigest)
        (COMMAND_DESCRIPTOR_t *)_PolicyGetDigestDataAddress,
#endif // CC_PolicyGetDigest
#if (PAD_LIST || CC_TestParms)
        (COMMAND_DESCRIPTOR_t *)_TestParmsDataAddress,
#endif // CC_TestParms
#if (PAD_LIST || CC_Commit)
        (COMMAND_DESCRIPTOR_t *)_CommitDataAddress,
#endif // CC_Commit
#if (PAD_LIST || CC_PolicyPassword)
        (COMMAND_DESCRIPTOR_t *)_PolicyPasswordDataAddress,
#endif // CC_PolicyPassword
#if (PAD_LIST || CC_ZGen_2Phase)
        (COMMAND_DESCRIPTOR_t *)_ZGen_2PhaseDataAddress,
#endif // CC_ZGen_2Phase
#if (PAD_LIST || CC_EC_Ephemeral)
        (COMMAND_DESCRIPTOR_t *)_EC_EphemeralDataAddress,
#endif // CC_EC_Ephemeral
#if (PAD_LIST || CC_PolicyNvWritten)
        (COMMAND_DESCRIPTOR_t *)_PolicyNvWrittenDataAddress,
#endif // CC_PolicyNvWritten
#if (PAD_LIST || CC_PolicyTemplate)
        (COMMAND_DESCRIPTOR_t *)_PolicyTemplateDataAddress,
#endif // CC_PolicyTemplate
#if (PAD_LIST || CC_CreateLoaded)
        (COMMAND_DESCRIPTOR_t *)_CreateLoadedDataAddress,
#endif // CC_CreateLoaded
#if (PAD_LIST || CC_PolicyAuthorizeNV)
        (COMMAND_DESCRIPTOR_t *)_PolicyAuthorizeNVDataAddress,
#endif // CC_PolicyAuthorizeNV
#if (PAD_LIST || CC_EncryptDecrypt2)
        (COMMAND_DESCRIPTOR_t *)_EncryptDecrypt2DataAddress,
#endif // CC_EncryptDecrypt2
#if (PAD_LIST || CC_AC_GetCapability)
        (COMMAND_DESCRIPTOR_t *)_AC_GetCapabilityDataAddress,
#endif // CC_AC_GetCapability
#if (PAD_LIST || CC_AC_Send)
        (COMMAND_DESCRIPTOR_t *)_AC_SendDataAddress,
#endif // CC_AC_Send
#if (PAD_LIST || CC_Policy_AC_SendSelect)
        (COMMAND_DESCRIPTOR_t *)_Policy_AC_SendSelectDataAddress,
#endif // CC_Policy_AC_SendSelect
#if (PAD_LIST || CC_CertifyX509)
        (COMMAND_DESCRIPTOR_t *)_CertifyX509DataAddress,
#endif // CC_CertifyX509
#if (PAD_LIST || CC_Vendor_TCG_Test)
        (COMMAND_DESCRIPTOR_t *)_Vendor_TCG_TestDataAddress,
#endif // CC_Vendor_TCG_Test
        0
};


#endif  // _COMMAND_TABLE_DISPATCH_
