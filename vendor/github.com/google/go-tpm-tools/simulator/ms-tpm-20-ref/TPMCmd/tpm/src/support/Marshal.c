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
 *  Created by TpmMarshal; Version 4.1 Dec 10, 2018
 *  Date: Apr  2, 2019  Time: 11:00:48AM
 */

#include "Tpm.h"
#include "Marshal_fp.h"

// Table 2:3 - Definition of Base Types
//   UINT8 definition from table 2:3
TPM_RC
UINT8_Unmarshal(UINT8 *target, BYTE **buffer, INT32 *size)
{
    if((*size -= 1) < 0)
        return TPM_RC_INSUFFICIENT;
    *target = BYTE_ARRAY_TO_UINT8(*buffer);
    *buffer += 1;
    return TPM_RC_SUCCESS;
}
UINT16
UINT8_Marshal(UINT8 *source, BYTE **buffer, INT32 *size)
{
    if (buffer != 0)
    {
        if ((size == 0) || ((*size -= 1) >= 0))
        {
            UINT8_TO_BYTE_ARRAY(*source, *buffer);
            *buffer += 1;
        }
        pAssert(size == 0 || (*size >= 0));
    }
    return (1);
}

//   BYTE definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
BYTE_Unmarshal(BYTE *target, BYTE **buffer, INT32 *size)
{
    return UINT8_Unmarshal((UINT8 *)target, buffer, size);
}
UINT16
BYTE_Marshal(BYTE *source, BYTE **buffer, INT32 *size)
{
    return UINT8_Marshal((UINT8 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

//   INT8 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT8_Unmarshal(INT8 *target, BYTE **buffer, INT32 *size)
{
    return UINT8_Unmarshal((UINT8 *)target, buffer, size);
}
UINT16
INT8_Marshal(INT8 *source, BYTE **buffer, INT32 *size)
{
    return UINT8_Marshal((UINT8 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

//   UINT16 definition from table 2:3
TPM_RC
UINT16_Unmarshal(UINT16 *target, BYTE **buffer, INT32 *size)
{
    if((*size -= 2) < 0)
        return TPM_RC_INSUFFICIENT;
    *target = BYTE_ARRAY_TO_UINT16(*buffer);
    *buffer += 2;
    return TPM_RC_SUCCESS;
}
UINT16
UINT16_Marshal(UINT16 *source, BYTE **buffer, INT32 *size)
{
    if (buffer != 0)
    {
        if ((size == 0) || ((*size -= 2) >= 0))
        {
            UINT16_TO_BYTE_ARRAY(*source, *buffer);
            *buffer += 2;
        }
        pAssert(size == 0 || (*size >= 0));
    }
    return (2);
}

//   INT16 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT16_Unmarshal(INT16 *target, BYTE **buffer, INT32 *size)
{
    return UINT16_Unmarshal((UINT16 *)target, buffer, size);
}
UINT16
INT16_Marshal(INT16 *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

//   UINT32 definition from table 2:3
TPM_RC
UINT32_Unmarshal(UINT32 *target, BYTE **buffer, INT32 *size)
{
    if((*size -= 4) < 0)
        return TPM_RC_INSUFFICIENT;
    *target = BYTE_ARRAY_TO_UINT32(*buffer);
    *buffer += 4;
    return TPM_RC_SUCCESS;
}
UINT16
UINT32_Marshal(UINT32 *source, BYTE **buffer, INT32 *size)
{
    if (buffer != 0)
    {
        if ((size == 0) || ((*size -= 4) >= 0))
        {
            UINT32_TO_BYTE_ARRAY(*source, *buffer);
            *buffer += 4;
        }
        pAssert(size == 0 || (*size >= 0));
    }
    return (4);
}

//   INT32 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT32_Unmarshal(INT32 *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
INT32_Marshal(INT32 *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

//   UINT64 definition from table 2:3
TPM_RC
UINT64_Unmarshal(UINT64 *target, BYTE **buffer, INT32 *size)
{
    if((*size -= 8) < 0)
        return TPM_RC_INSUFFICIENT;
    *target = BYTE_ARRAY_TO_UINT64(*buffer);
    *buffer += 8;
    return TPM_RC_SUCCESS;
}
UINT16
UINT64_Marshal(UINT64 *source, BYTE **buffer, INT32 *size)
{
    if (buffer != 0)
    {
        if ((size == 0) || ((*size -= 8) >= 0))
        {
            UINT64_TO_BYTE_ARRAY(*source, *buffer);
            *buffer += 8;
        }
        pAssert(size == 0 || (*size >= 0));
    }
    return (8);
}

//   INT64 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT64_Unmarshal(INT64 *target, BYTE **buffer, INT32 *size)
{
    return UINT64_Unmarshal((UINT64 *)target, buffer, size);
}
UINT16
INT64_Marshal(INT64 *source, BYTE **buffer, INT32 *size)
{
    return UINT64_Marshal((UINT64 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:4 - Defines for Logic Values
// Table 2:5 - Definition of Types for Documentation Clarity
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ALGORITHM_ID_Unmarshal(TPM_ALGORITHM_ID *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_ALGORITHM_ID_Marshal(TPM_ALGORITHM_ID *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
TPM_RC
TPM_MODIFIER_INDICATOR_Unmarshal(TPM_MODIFIER_INDICATOR *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_MODIFIER_INDICATOR_Marshal(TPM_MODIFIER_INDICATOR *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
TPM_RC
TPM_AUTHORIZATION_SIZE_Unmarshal(TPM_AUTHORIZATION_SIZE *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_AUTHORIZATION_SIZE_Marshal(TPM_AUTHORIZATION_SIZE *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
TPM_RC
TPM_PARAMETER_SIZE_Unmarshal(TPM_PARAMETER_SIZE *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_PARAMETER_SIZE_Marshal(TPM_PARAMETER_SIZE *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
TPM_RC
TPM_KEY_SIZE_Unmarshal(TPM_KEY_SIZE *target, BYTE **buffer, INT32 *size)
{
    return UINT16_Unmarshal((UINT16 *)target, buffer, size);
}
UINT16
TPM_KEY_SIZE_Marshal(TPM_KEY_SIZE *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
TPM_RC
TPM_KEY_BITS_Unmarshal(TPM_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    return UINT16_Unmarshal((UINT16 *)target, buffer, size);
}
UINT16
TPM_KEY_BITS_Marshal(TPM_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:6 - Definition of TPM_SPEC Constants
// Table 2:7 - Definition of TPM_GENERATED Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_GENERATED_Marshal(TPM_GENERATED *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:9 - Definition of TPM_ALG_ID Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ALG_ID_Unmarshal(TPM_ALG_ID *target, BYTE **buffer, INT32 *size)
{
    return UINT16_Unmarshal((UINT16 *)target, buffer, size);
}
UINT16
TPM_ALG_ID_Marshal(TPM_ALG_ID *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:10 - Definition of TPM_ECC_CURVE Constants
#if ALG_ECC
TPM_RC
TPM_ECC_CURVE_Unmarshal(TPM_ECC_CURVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_ECC_NIST_P192 :
            case TPM_ECC_NIST_P224 :
            case TPM_ECC_NIST_P256 :
            case TPM_ECC_NIST_P384 :
            case TPM_ECC_NIST_P521 :
            case TPM_ECC_BN_P256 :
            case TPM_ECC_BN_P638 :
            case TPM_ECC_SM2_P256 :
                break;
            default :
                result = TPM_RC_CURVE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPM_ECC_CURVE_Marshal(TPM_ECC_CURVE *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:12 - Definition of TPM_CC Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_CC_Unmarshal(TPM_CC *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_CC_Marshal(TPM_CC *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:16 - Definition of TPM_RC Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_RC_Marshal(TPM_RC *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:17 - Definition of TPM_CLOCK_ADJUST Constants
TPM_RC
TPM_CLOCK_ADJUST_Unmarshal(TPM_CLOCK_ADJUST *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = INT8_Unmarshal((INT8 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_CLOCK_COARSE_SLOWER :
            case TPM_CLOCK_MEDIUM_SLOWER :
            case TPM_CLOCK_FINE_SLOWER :
            case TPM_CLOCK_NO_CHANGE :
            case TPM_CLOCK_FINE_FASTER :
            case TPM_CLOCK_MEDIUM_FASTER :
            case TPM_CLOCK_COARSE_FASTER :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:18 - Definition of TPM_EO Constants
TPM_RC
TPM_EO_Unmarshal(TPM_EO *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_EO_EQ :
            case TPM_EO_NEQ :
            case TPM_EO_SIGNED_GT :
            case TPM_EO_UNSIGNED_GT :
            case TPM_EO_SIGNED_LT :
            case TPM_EO_UNSIGNED_LT :
            case TPM_EO_SIGNED_GE :
            case TPM_EO_UNSIGNED_GE :
            case TPM_EO_SIGNED_LE :
            case TPM_EO_UNSIGNED_LE :
            case TPM_EO_BITSET :
            case TPM_EO_BITCLEAR :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPM_EO_Marshal(TPM_EO *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:19 - Definition of TPM_ST Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ST_Unmarshal(TPM_ST *target, BYTE **buffer, INT32 *size)
{
    return UINT16_Unmarshal((UINT16 *)target, buffer, size);
}
UINT16
TPM_ST_Marshal(TPM_ST *source, BYTE **buffer, INT32 *size)
{
    return UINT16_Marshal((UINT16 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:20 - Definition of TPM_SU Constants
TPM_RC
TPM_SU_Unmarshal(TPM_SU *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_SU_CLEAR :
            case TPM_SU_STATE :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:21 - Definition of TPM_SE Constants
TPM_RC
TPM_SE_Unmarshal(TPM_SE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT8_Unmarshal((UINT8 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_SE_HMAC :
            case TPM_SE_POLICY :
            case TPM_SE_TRIAL :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:22 - Definition of TPM_CAP Constants
TPM_RC
TPM_CAP_Unmarshal(TPM_CAP *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_CAP_ALGS :
            case TPM_CAP_HANDLES :
            case TPM_CAP_COMMANDS :
            case TPM_CAP_PP_COMMANDS :
            case TPM_CAP_AUDIT_COMMANDS :
            case TPM_CAP_PCRS :
            case TPM_CAP_TPM_PROPERTIES :
            case TPM_CAP_PCR_PROPERTIES :
            case TPM_CAP_ECC_CURVES :
            case TPM_CAP_AUTH_POLICIES :
            case TPM_CAP_VENDOR_PROPERTY :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPM_CAP_Marshal(TPM_CAP *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:23 - Definition of TPM_PT Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_PT_Unmarshal(TPM_PT *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_PT_Marshal(TPM_PT *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:24 - Definition of TPM_PT_PCR Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_PT_PCR_Unmarshal(TPM_PT_PCR *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_PT_PCR_Marshal(TPM_PT_PCR *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:25 - Definition of TPM_PS Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_PS_Marshal(TPM_PS *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:26 - Definition of Types for Handles
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HANDLE_Unmarshal(TPM_HANDLE *target, BYTE **buffer, INT32 *size)
{
    return UINT32_Unmarshal((UINT32 *)target, buffer, size);
}
UINT16
TPM_HANDLE_Marshal(TPM_HANDLE *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:27 - Definition of TPM_HT Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HT_Unmarshal(TPM_HT *target, BYTE **buffer, INT32 *size)
{
    return UINT8_Unmarshal((UINT8 *)target, buffer, size);
}
UINT16
TPM_HT_Marshal(TPM_HT *source, BYTE **buffer, INT32 *size)
{
    return UINT8_Marshal((UINT8 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:28 - Definition of TPM_RH Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_RH_Unmarshal(TPM_RH *target, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
}
UINT16
TPM_RH_Marshal(TPM_RH *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:29 - Definition of TPM_HC Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HC_Unmarshal(TPM_HC *target, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
}
UINT16
TPM_HC_Marshal(TPM_HC *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:30 - Definition of TPMA_ALGORITHM Bits
TPM_RC
TPMA_ALGORITHM_Unmarshal(TPMA_ALGORITHM *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*((UINT32 *)target) & (UINT32)0xfffff8f0)
            result = TPM_RC_RESERVED_BITS;
    }
    return result;
}

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_ALGORITHM_Marshal(TPMA_ALGORITHM *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:31 - Definition of TPMA_OBJECT Bits
TPM_RC
TPMA_OBJECT_Unmarshal(TPMA_OBJECT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*((UINT32 *)target) & (UINT32)0xfff0f309)
            result = TPM_RC_RESERVED_BITS;
    }
    return result;
}

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_OBJECT_Marshal(TPMA_OBJECT *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:32 - Definition of TPMA_SESSION Bits
TPM_RC
TPMA_SESSION_Unmarshal(TPMA_SESSION *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT8_Unmarshal((UINT8 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*((UINT8 *)target) & (UINT8)0x18)
            result = TPM_RC_RESERVED_BITS;
    }
    return result;
}

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_SESSION_Marshal(TPMA_SESSION *source, BYTE **buffer, INT32 *size)
{
    return UINT8_Marshal((UINT8 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:33 - Definition of TPMA_LOCALITY Bits
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMA_LOCALITY_Unmarshal(TPMA_LOCALITY *target, BYTE **buffer, INT32 *size)
{
    return UINT8_Unmarshal((UINT8 *)target, buffer, size);
}
UINT16
TPMA_LOCALITY_Marshal(TPMA_LOCALITY *source, BYTE **buffer, INT32 *size)
{
    return UINT8_Marshal((UINT8 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:34 - Definition of TPMA_PERMANENT Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_PERMANENT_Marshal(TPMA_PERMANENT *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:35 - Definition of TPMA_STARTUP_CLEAR Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_STARTUP_CLEAR_Marshal(TPMA_STARTUP_CLEAR *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:36 - Definition of TPMA_MEMORY Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_MEMORY_Marshal(TPMA_MEMORY *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:37 - Definition of TPMA_CC Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_CC_Marshal(TPMA_CC *source, BYTE **buffer, INT32 *size)
{
    return TPM_CC_Marshal((TPM_CC *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:38 - Definition of TPMA_MODES Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_MODES_Marshal(TPMA_MODES *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:39 - Definition of TPMA_X509_KEY_USAGE Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_X509_KEY_USAGE_Marshal(TPMA_X509_KEY_USAGE *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:40 - Definition of TPMI_YES_NO Type
TPM_RC
TPMI_YES_NO_Unmarshal(TPMI_YES_NO *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = BYTE_Unmarshal((BYTE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case NO:
            case YES:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_YES_NO_Marshal(TPMI_YES_NO *source, BYTE **buffer, INT32 *size)
{
    return BYTE_Marshal((BYTE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:41 - Definition of TPMI_DH_OBJECT Type
TPM_RC
TPMI_DH_OBJECT_Unmarshal(TPMI_DH_OBJECT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*target == TPM_RH_NULL)
        {
            if(!flag)
                result = TPM_RC_VALUE;
        }
        else if(  ((*target < TRANSIENT_FIRST) || (*target > TRANSIENT_LAST))
              && ((*target < PERSISTENT_FIRST) || (*target > PERSISTENT_LAST)))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_OBJECT_Marshal(TPMI_DH_OBJECT *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:42 - Definition of TPMI_DH_PARENT Type
TPM_RC
TPMI_DH_PARENT_Unmarshal(TPMI_DH_PARENT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_PLATFORM:
            case TPM_RH_ENDORSEMENT:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                if(  ((*target < TRANSIENT_FIRST) || (*target > TRANSIENT_LAST))
                  && ((*target < PERSISTENT_FIRST) || (*target > PERSISTENT_LAST)))
                    result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_PARENT_Marshal(TPMI_DH_PARENT *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:43 - Definition of TPMI_DH_PERSISTENT Type
TPM_RC
TPMI_DH_PERSISTENT_Unmarshal(TPMI_DH_PERSISTENT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((*target < PERSISTENT_FIRST) || (*target > PERSISTENT_LAST))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_PERSISTENT_Marshal(TPMI_DH_PERSISTENT *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:44 - Definition of TPMI_DH_ENTITY Type
TPM_RC
TPMI_DH_ENTITY_Unmarshal(TPMI_DH_ENTITY *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_ENDORSEMENT:
            case TPM_RH_PLATFORM:
            case TPM_RH_LOCKOUT:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                if(  ((*target < TRANSIENT_FIRST) || (*target > TRANSIENT_LAST))
                  && ((*target < PERSISTENT_FIRST) || (*target > PERSISTENT_LAST))
                  && ((*target < NV_INDEX_FIRST) || (*target > NV_INDEX_LAST))
                  && (*target > PCR_LAST)
                  && ((*target < TPM_RH_AUTH_00) || (*target > TPM_RH_AUTH_FF)))
                    result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:45 - Definition of TPMI_DH_PCR Type
TPM_RC
TPMI_DH_PCR_Unmarshal(TPMI_DH_PCR *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*target == TPM_RH_NULL)
        {
            if(!flag)
                result = TPM_RC_VALUE;
        }
        else if(*target > PCR_LAST)
            result = TPM_RC_VALUE;
    }
    return result;
}

// Table 2:46 - Definition of TPMI_SH_AUTH_SESSION Type
TPM_RC
TPMI_SH_AUTH_SESSION_Unmarshal(TPMI_SH_AUTH_SESSION *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*target == TPM_RS_PW)
        {
            if(!flag)
                result = TPM_RC_VALUE;
        }
        else if(  ((*target < HMAC_SESSION_FIRST) || (*target > HMAC_SESSION_LAST))
              && ((*target < POLICY_SESSION_FIRST) || (*target > POLICY_SESSION_LAST)))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_AUTH_SESSION_Marshal(TPMI_SH_AUTH_SESSION *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:47 - Definition of TPMI_SH_HMAC Type
TPM_RC
TPMI_SH_HMAC_Unmarshal(TPMI_SH_HMAC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((*target < HMAC_SESSION_FIRST) || (*target > HMAC_SESSION_LAST))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_HMAC_Marshal(TPMI_SH_HMAC *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:48 - Definition of TPMI_SH_POLICY Type
TPM_RC
TPMI_SH_POLICY_Unmarshal(TPMI_SH_POLICY *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((*target < POLICY_SESSION_FIRST) || (*target > POLICY_SESSION_LAST))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_POLICY_Marshal(TPMI_SH_POLICY *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:49 - Definition of TPMI_DH_CONTEXT Type
TPM_RC
TPMI_DH_CONTEXT_Unmarshal(TPMI_DH_CONTEXT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(  ((*target < HMAC_SESSION_FIRST) || (*target > HMAC_SESSION_LAST))
          && ((*target < POLICY_SESSION_FIRST) || (*target > POLICY_SESSION_LAST))
          && ((*target < TRANSIENT_FIRST) || (*target > TRANSIENT_LAST)))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_CONTEXT_Marshal(TPMI_DH_CONTEXT *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:50 - Definition of TPMI_DH_SAVED Type
TPM_RC
TPMI_DH_SAVED_Unmarshal(TPMI_DH_SAVED *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case 0x80000000:
            case 0x80000001:
            case 0x80000002:
                break;
            default:
                if(  ((*target < HMAC_SESSION_FIRST) || (*target > HMAC_SESSION_LAST))
                  && ((*target < POLICY_SESSION_FIRST) || (*target > POLICY_SESSION_LAST)))
                    result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_SAVED_Marshal(TPMI_DH_SAVED *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:51 - Definition of TPMI_RH_HIERARCHY Type
TPM_RC
TPMI_RH_HIERARCHY_Unmarshal(TPMI_RH_HIERARCHY *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_PLATFORM:
            case TPM_RH_ENDORSEMENT:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_HIERARCHY_Marshal(TPMI_RH_HIERARCHY *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:52 - Definition of TPMI_RH_ENABLES Type
TPM_RC
TPMI_RH_ENABLES_Unmarshal(TPMI_RH_ENABLES *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_PLATFORM:
            case TPM_RH_ENDORSEMENT:
            case TPM_RH_PLATFORM_NV:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_ENABLES_Marshal(TPMI_RH_ENABLES *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:53 - Definition of TPMI_RH_HIERARCHY_AUTH Type
TPM_RC
TPMI_RH_HIERARCHY_AUTH_Unmarshal(TPMI_RH_HIERARCHY_AUTH *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_PLATFORM:
            case TPM_RH_ENDORSEMENT:
            case TPM_RH_LOCKOUT:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:54 - Definition of TPMI_RH_PLATFORM Type
TPM_RC
TPMI_RH_PLATFORM_Unmarshal(TPMI_RH_PLATFORM *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_PLATFORM:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:55 - Definition of TPMI_RH_OWNER Type
TPM_RC
TPMI_RH_OWNER_Unmarshal(TPMI_RH_OWNER *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:56 - Definition of TPMI_RH_ENDORSEMENT Type
TPM_RC
TPMI_RH_ENDORSEMENT_Unmarshal(TPMI_RH_ENDORSEMENT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_ENDORSEMENT:
                break;
            case TPM_RH_NULL:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:57 - Definition of TPMI_RH_PROVISION Type
TPM_RC
TPMI_RH_PROVISION_Unmarshal(TPMI_RH_PROVISION *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_OWNER:
            case TPM_RH_PLATFORM:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:58 - Definition of TPMI_RH_CLEAR Type
TPM_RC
TPMI_RH_CLEAR_Unmarshal(TPMI_RH_CLEAR *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_LOCKOUT:
            case TPM_RH_PLATFORM:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:59 - Definition of TPMI_RH_NV_AUTH Type
TPM_RC
TPMI_RH_NV_AUTH_Unmarshal(TPMI_RH_NV_AUTH *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_PLATFORM:
            case TPM_RH_OWNER:
                break;
            default:
                if((*target < NV_INDEX_FIRST) || (*target > NV_INDEX_LAST))
                    result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:60 - Definition of TPMI_RH_LOCKOUT Type
TPM_RC
TPMI_RH_LOCKOUT_Unmarshal(TPMI_RH_LOCKOUT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_RH_LOCKOUT:
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}

// Table 2:61 - Definition of TPMI_RH_NV_INDEX Type
TPM_RC
TPMI_RH_NV_INDEX_Unmarshal(TPMI_RH_NV_INDEX *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((*target < NV_INDEX_FIRST) || (*target > NV_INDEX_LAST))
            result = TPM_RC_VALUE;
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_NV_INDEX_Marshal(TPMI_RH_NV_INDEX *source, BYTE **buffer, INT32 *size)
{
    return TPM_HANDLE_Marshal((TPM_HANDLE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:62 - Definition of TPMI_RH_AC Type
TPM_RC
TPMI_RH_AC_Unmarshal(TPMI_RH_AC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_HANDLE_Unmarshal((TPM_HANDLE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((*target < AC_FIRST) || (*target > AC_LAST))
            result = TPM_RC_VALUE;
    }
    return result;
}

// Table 2:63 - Definition of TPMI_ALG_HASH Type
TPM_RC
TPMI_ALG_HASH_Unmarshal(TPMI_ALG_HASH *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_SHA1
            case ALG_SHA1_VALUE:
#endif // ALG_SHA1
#if ALG_SHA256
            case ALG_SHA256_VALUE:
#endif // ALG_SHA256
#if ALG_SHA384
            case ALG_SHA384_VALUE:
#endif // ALG_SHA384
#if ALG_SHA512
            case ALG_SHA512_VALUE:
#endif // ALG_SHA512
#if ALG_SM3_256
            case ALG_SM3_256_VALUE:
#endif // ALG_SM3_256
#if ALG_SHA3_256
            case ALG_SHA3_256_VALUE:
#endif // ALG_SHA3_256
#if ALG_SHA3_384
            case ALG_SHA3_384_VALUE:
#endif // ALG_SHA3_384
#if ALG_SHA3_512
            case ALG_SHA3_512_VALUE:
#endif // ALG_SHA3_512
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_HASH;
                break;
            default:
                result = TPM_RC_HASH;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_HASH_Marshal(TPMI_ALG_HASH *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:64 - Definition of TPMI_ALG_ASYM Type
TPM_RC
TPMI_ALG_ASYM_Unmarshal(TPMI_ALG_ASYM *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_RSA
            case ALG_RSA_VALUE:
#endif // ALG_RSA
#if ALG_ECC
            case ALG_ECC_VALUE:
#endif // ALG_ECC
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_ASYMMETRIC;
                break;
            default:
                result = TPM_RC_ASYMMETRIC;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ASYM_Marshal(TPMI_ALG_ASYM *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:65 - Definition of TPMI_ALG_SYM Type
TPM_RC
TPMI_ALG_SYM_Unmarshal(TPMI_ALG_SYM *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_TDES
            case ALG_TDES_VALUE:
#endif // ALG_TDES
#if ALG_AES
            case ALG_AES_VALUE:
#endif // ALG_AES
#if ALG_SM4
            case ALG_SM4_VALUE:
#endif // ALG_SM4
#if ALG_CAMELLIA
            case ALG_CAMELLIA_VALUE:
#endif // ALG_CAMELLIA
#if ALG_XOR
            case ALG_XOR_VALUE:
#endif // ALG_XOR
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SYMMETRIC;
                break;
            default:
                result = TPM_RC_SYMMETRIC;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_Marshal(TPMI_ALG_SYM *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:66 - Definition of TPMI_ALG_SYM_OBJECT Type
TPM_RC
TPMI_ALG_SYM_OBJECT_Unmarshal(TPMI_ALG_SYM_OBJECT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_TDES
            case ALG_TDES_VALUE:
#endif // ALG_TDES
#if ALG_AES
            case ALG_AES_VALUE:
#endif // ALG_AES
#if ALG_SM4
            case ALG_SM4_VALUE:
#endif // ALG_SM4
#if ALG_CAMELLIA
            case ALG_CAMELLIA_VALUE:
#endif // ALG_CAMELLIA
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SYMMETRIC;
                break;
            default:
                result = TPM_RC_SYMMETRIC;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_OBJECT_Marshal(TPMI_ALG_SYM_OBJECT *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:67 - Definition of TPMI_ALG_SYM_MODE Type
TPM_RC
TPMI_ALG_SYM_MODE_Unmarshal(TPMI_ALG_SYM_MODE *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
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
#if ALG_CMAC
            case ALG_CMAC_VALUE:
#endif // ALG_CMAC
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_MODE;
                break;
            default:
                result = TPM_RC_MODE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_MODE_Marshal(TPMI_ALG_SYM_MODE *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:68 - Definition of TPMI_ALG_KDF Type
TPM_RC
TPMI_ALG_KDF_Unmarshal(TPMI_ALG_KDF *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_MGF1
            case ALG_MGF1_VALUE:
#endif // ALG_MGF1
#if ALG_KDF1_SP800_56A
            case ALG_KDF1_SP800_56A_VALUE:
#endif // ALG_KDF1_SP800_56A
#if ALG_KDF2
            case ALG_KDF2_VALUE:
#endif // ALG_KDF2
#if ALG_KDF1_SP800_108
            case ALG_KDF1_SP800_108_VALUE:
#endif // ALG_KDF1_SP800_108
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_KDF;
                break;
            default:
                result = TPM_RC_KDF;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_KDF_Marshal(TPMI_ALG_KDF *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:69 - Definition of TPMI_ALG_SIG_SCHEME Type
TPM_RC
TPMI_ALG_SIG_SCHEME_Unmarshal(TPMI_ALG_SIG_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_ECDAA
            case ALG_ECDAA_VALUE:
#endif // ALG_ECDAA
#if ALG_RSASSA
            case ALG_RSASSA_VALUE:
#endif // ALG_RSASSA
#if ALG_RSAPSS
            case ALG_RSAPSS_VALUE:
#endif // ALG_RSAPSS
#if ALG_ECDSA
            case ALG_ECDSA_VALUE:
#endif // ALG_ECDSA
#if ALG_SM2
            case ALG_SM2_VALUE:
#endif // ALG_SM2
#if ALG_ECSCHNORR
            case ALG_ECSCHNORR_VALUE:
#endif // ALG_ECSCHNORR
#if ALG_HMAC
            case ALG_HMAC_VALUE:
#endif // ALG_HMAC
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SCHEME;
                break;
            default:
                result = TPM_RC_SCHEME;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SIG_SCHEME_Marshal(TPMI_ALG_SIG_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:70 - Definition of TPMI_ECC_KEY_EXCHANGE Type
#if ALG_ECC
TPM_RC
TPMI_ECC_KEY_EXCHANGE_Unmarshal(TPMI_ECC_KEY_EXCHANGE *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_ECDH
            case ALG_ECDH_VALUE:
#endif // ALG_ECDH
#if ALG_ECMQV
            case ALG_ECMQV_VALUE:
#endif // ALG_ECMQV
#if ALG_SM2
            case ALG_SM2_VALUE:
#endif // ALG_SM2
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SCHEME;
                break;
            default:
                result = TPM_RC_SCHEME;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ECC_KEY_EXCHANGE_Marshal(TPMI_ECC_KEY_EXCHANGE *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:71 - Definition of TPMI_ST_COMMAND_TAG Type
TPM_RC
TPMI_ST_COMMAND_TAG_Unmarshal(TPMI_ST_COMMAND_TAG *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ST_Unmarshal((TPM_ST *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
            case TPM_ST_NO_SESSIONS:
            case TPM_ST_SESSIONS:
                break;
            default:
                result = TPM_RC_BAD_TAG;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ST_COMMAND_TAG_Marshal(TPMI_ST_COMMAND_TAG *source, BYTE **buffer, INT32 *size)
{
    return TPM_ST_Marshal((TPM_ST *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:72 - Definition of TPMI_ALG_MAC_SCHEME Type
TPM_RC
TPMI_ALG_MAC_SCHEME_Unmarshal(TPMI_ALG_MAC_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_CMAC
            case ALG_CMAC_VALUE:
#endif // ALG_CMAC
#if ALG_SHA1
            case ALG_SHA1_VALUE:
#endif // ALG_SHA1
#if ALG_SHA256
            case ALG_SHA256_VALUE:
#endif // ALG_SHA256
#if ALG_SHA384
            case ALG_SHA384_VALUE:
#endif // ALG_SHA384
#if ALG_SHA512
            case ALG_SHA512_VALUE:
#endif // ALG_SHA512
#if ALG_SM3_256
            case ALG_SM3_256_VALUE:
#endif // ALG_SM3_256
#if ALG_SHA3_256
            case ALG_SHA3_256_VALUE:
#endif // ALG_SHA3_256
#if ALG_SHA3_384
            case ALG_SHA3_384_VALUE:
#endif // ALG_SHA3_384
#if ALG_SHA3_512
            case ALG_SHA3_512_VALUE:
#endif // ALG_SHA3_512
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SYMMETRIC;
                break;
            default:
                result = TPM_RC_SYMMETRIC;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_MAC_SCHEME_Marshal(TPMI_ALG_MAC_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:73 - Definition of TPMI_ALG_CIPHER_MODE Type
TPM_RC
TPMI_ALG_CIPHER_MODE_Unmarshal(TPMI_ALG_CIPHER_MODE *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
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
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_MODE;
                break;
            default:
                result = TPM_RC_MODE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_CIPHER_MODE_Marshal(TPMI_ALG_CIPHER_MODE *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:74 - Definition of TPMS_EMPTY Structure
TPM_RC
TPMS_EMPTY_Unmarshal(TPMS_EMPTY *target, BYTE **buffer, INT32 *size)
{
    // to prevent the compiler from complaining
    NOT_REFERENCED(target);
    NOT_REFERENCED(buffer);
    NOT_REFERENCED(size);
    return TPM_RC_SUCCESS;
}
UINT16
TPMS_EMPTY_Marshal(TPMS_EMPTY *source, BYTE **buffer, INT32 *size)
{
    // to prevent the compiler from complaining
    NOT_REFERENCED(source);
    NOT_REFERENCED(buffer);
    NOT_REFERENCED(size);
    return 0;
}

// Table 2:75 - Definition of TPMS_ALGORITHM_DESCRIPTION Structure
UINT16
TPMS_ALGORITHM_DESCRIPTION_Marshal(TPMS_ALGORITHM_DESCRIPTION *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ALG_ID_Marshal((TPM_ALG_ID *)&(source->alg), buffer, size));
    result = (UINT16)(result + TPMA_ALGORITHM_Marshal((TPMA_ALGORITHM *)&(source->attributes), buffer, size));
    return result;
}

// Table 2:76 - Definition of TPMU_HA Union
TPM_RC
TPMU_HA_Unmarshal(TPMU_HA *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha1), buffer, size, (INT32)SHA1_DIGEST_SIZE);
#endif // ALG_SHA1
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha256), buffer, size, (INT32)SHA256_DIGEST_SIZE);
#endif // ALG_SHA256
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha384), buffer, size, (INT32)SHA384_DIGEST_SIZE);
#endif // ALG_SHA384
#if ALG_SHA512
        case ALG_SHA512_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha512), buffer, size, (INT32)SHA512_DIGEST_SIZE);
#endif // ALG_SHA512
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sm3_256), buffer, size, (INT32)SM3_256_DIGEST_SIZE);
#endif // ALG_SM3_256
#if ALG_SHA3_256
        case ALG_SHA3_256_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha3_256), buffer, size, (INT32)SHA3_256_DIGEST_SIZE);
#endif // ALG_SHA3_256
#if ALG_SHA3_384
        case ALG_SHA3_384_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha3_384), buffer, size, (INT32)SHA3_384_DIGEST_SIZE);
#endif // ALG_SHA3_384
#if ALG_SHA3_512
        case ALG_SHA3_512_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->sha3_512), buffer, size, (INT32)SHA3_512_DIGEST_SIZE);
#endif // ALG_SHA3_512
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_HA_Marshal(TPMU_HA *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_SHA1
        case ALG_SHA1_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha1), buffer, size, (INT32)SHA1_DIGEST_SIZE);
#endif // ALG_SHA1
#if ALG_SHA256
        case ALG_SHA256_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha256), buffer, size, (INT32)SHA256_DIGEST_SIZE);
#endif // ALG_SHA256
#if ALG_SHA384
        case ALG_SHA384_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha384), buffer, size, (INT32)SHA384_DIGEST_SIZE);
#endif // ALG_SHA384
#if ALG_SHA512
        case ALG_SHA512_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha512), buffer, size, (INT32)SHA512_DIGEST_SIZE);
#endif // ALG_SHA512
#if ALG_SM3_256
        case ALG_SM3_256_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sm3_256), buffer, size, (INT32)SM3_256_DIGEST_SIZE);
#endif // ALG_SM3_256
#if ALG_SHA3_256
        case ALG_SHA3_256_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha3_256), buffer, size, (INT32)SHA3_256_DIGEST_SIZE);
#endif // ALG_SHA3_256
#if ALG_SHA3_384
        case ALG_SHA3_384_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha3_384), buffer, size, (INT32)SHA3_384_DIGEST_SIZE);
#endif // ALG_SHA3_384
#if ALG_SHA3_512
        case ALG_SHA3_512_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->sha3_512), buffer, size, (INT32)SHA3_512_DIGEST_SIZE);
#endif // ALG_SHA3_512
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:77 - Definition of TPMT_HA Structure
TPM_RC
TPMT_HA_Unmarshal(TPMT_HA *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hashAlg), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_HA_Unmarshal((TPMU_HA *)&(target->digest), buffer, size, (UINT32)target->hashAlg);
    return result;
}
UINT16
TPMT_HA_Marshal(TPMT_HA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hashAlg), buffer, size));
    result = (UINT16)(result + TPMU_HA_Marshal((TPMU_HA *)&(source->digest), buffer, size, (UINT32)source->hashAlg));
    return result;
}

// Table 2:78 - Definition of TPM2B_DIGEST Structure
TPM_RC
TPM2B_DIGEST_Unmarshal(TPM2B_DIGEST *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMU_HA))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_DIGEST_Marshal(TPM2B_DIGEST *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:79 - Definition of TPM2B_DATA Structure
TPM_RC
TPM2B_DATA_Unmarshal(TPM2B_DATA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMT_HA))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_DATA_Marshal(TPM2B_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:80 - Definition of Types for TPM2B_NONCE
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_NONCE_Unmarshal(TPM2B_NONCE *target, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)target, buffer, size);
}
UINT16
TPM2B_NONCE_Marshal(TPM2B_NONCE *source, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:81 - Definition of Types for TPM2B_AUTH
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_AUTH_Unmarshal(TPM2B_AUTH *target, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)target, buffer, size);
}
UINT16
TPM2B_AUTH_Marshal(TPM2B_AUTH *source, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:82 - Definition of Types for TPM2B_OPERAND
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_OPERAND_Unmarshal(TPM2B_OPERAND *target, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)target, buffer, size);
}
UINT16
TPM2B_OPERAND_Marshal(TPM2B_OPERAND *source, BYTE **buffer, INT32 *size)
{
    return TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:83 - Definition of TPM2B_EVENT Structure
TPM_RC
TPM2B_EVENT_Unmarshal(TPM2B_EVENT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > 1024)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_EVENT_Marshal(TPM2B_EVENT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:84 - Definition of TPM2B_MAX_BUFFER Structure
TPM_RC
TPM2B_MAX_BUFFER_Unmarshal(TPM2B_MAX_BUFFER *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_DIGEST_BUFFER)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_MAX_BUFFER_Marshal(TPM2B_MAX_BUFFER *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:85 - Definition of TPM2B_MAX_NV_BUFFER Structure
TPM_RC
TPM2B_MAX_NV_BUFFER_Unmarshal(TPM2B_MAX_NV_BUFFER *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_NV_BUFFER_SIZE)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_MAX_NV_BUFFER_Marshal(TPM2B_MAX_NV_BUFFER *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:86 - Definition of TPM2B_TIMEOUT Structure
TPM_RC
TPM2B_TIMEOUT_Unmarshal(TPM2B_TIMEOUT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(UINT64))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_TIMEOUT_Marshal(TPM2B_TIMEOUT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:87 - Definition of TPM2B_IV Structure
TPM_RC
TPM2B_IV_Unmarshal(TPM2B_IV *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_SYM_BLOCK_SIZE)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_IV_Marshal(TPM2B_IV *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:88 - Definition of TPMU_NAME Union
// Table 2:89 - Definition of TPM2B_NAME Structure
TPM_RC
TPM2B_NAME_Unmarshal(TPM2B_NAME *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMU_NAME))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.name), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_NAME_Marshal(TPM2B_NAME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.name), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:90 - Definition of TPMS_PCR_SELECT Structure
TPM_RC
TPMS_PCR_SELECT_Unmarshal(TPMS_PCR_SELECT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT8_Unmarshal((UINT8 *)&(target->sizeofSelect), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->sizeofSelect < PCR_SELECT_MIN))
        result = TPM_RC_VALUE;
    if(result == TPM_RC_SUCCESS)
    {
        if((target->sizeofSelect) > PCR_SELECT_MAX)
            result = TPM_RC_VALUE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->pcrSelect), buffer, size, (INT32)(target->sizeofSelect));
    }
    return result;
}
UINT16
TPMS_PCR_SELECT_Marshal(TPMS_PCR_SELECT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT8_Marshal((UINT8 *)&(source->sizeofSelect), buffer, size));
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->pcrSelect), buffer, size, (INT32)(source->sizeofSelect)));
    return result;
}

// Table 2:91 - Definition of TPMS_PCR_SELECTION Structure
TPM_RC
TPMS_PCR_SELECTION_Unmarshal(TPMS_PCR_SELECTION *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hash), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = UINT8_Unmarshal((UINT8 *)&(target->sizeofSelect), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->sizeofSelect < PCR_SELECT_MIN))
        result = TPM_RC_VALUE;
    if(result == TPM_RC_SUCCESS)
    {
        if((target->sizeofSelect) > PCR_SELECT_MAX)
            result = TPM_RC_VALUE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->pcrSelect), buffer, size, (INT32)(target->sizeofSelect));
    }
    return result;
}
UINT16
TPMS_PCR_SELECTION_Marshal(TPMS_PCR_SELECTION *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hash), buffer, size));
    result = (UINT16)(result + UINT8_Marshal((UINT8 *)&(source->sizeofSelect), buffer, size));
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->pcrSelect), buffer, size, (INT32)(source->sizeofSelect)));
    return result;
}

// Table 2:94 - Definition of TPMT_TK_CREATION Structure
TPM_RC
TPMT_TK_CREATION_Unmarshal(TPMT_TK_CREATION *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ST_Unmarshal((TPM_ST *)&(target->tag), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->tag != TPM_ST_CREATION))
        result = TPM_RC_TAG;
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RH_HIERARCHY_Unmarshal((TPMI_RH_HIERARCHY *)&(target->hierarchy), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->digest), buffer, size);
    return result;
}
UINT16
TPMT_TK_CREATION_Marshal(TPMT_TK_CREATION *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ST_Marshal((TPM_ST *)&(source->tag), buffer, size));
    result = (UINT16)(result + TPMI_RH_HIERARCHY_Marshal((TPMI_RH_HIERARCHY *)&(source->hierarchy), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->digest), buffer, size));
    return result;
}

// Table 2:95 - Definition of TPMT_TK_VERIFIED Structure
TPM_RC
TPMT_TK_VERIFIED_Unmarshal(TPMT_TK_VERIFIED *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ST_Unmarshal((TPM_ST *)&(target->tag), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->tag != TPM_ST_VERIFIED))
        result = TPM_RC_TAG;
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RH_HIERARCHY_Unmarshal((TPMI_RH_HIERARCHY *)&(target->hierarchy), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->digest), buffer, size);
    return result;
}
UINT16
TPMT_TK_VERIFIED_Marshal(TPMT_TK_VERIFIED *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ST_Marshal((TPM_ST *)&(source->tag), buffer, size));
    result = (UINT16)(result + TPMI_RH_HIERARCHY_Marshal((TPMI_RH_HIERARCHY *)&(source->hierarchy), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->digest), buffer, size));
    return result;
}

// Table 2:96 - Definition of TPMT_TK_AUTH Structure
TPM_RC
TPMT_TK_AUTH_Unmarshal(TPMT_TK_AUTH *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ST_Unmarshal((TPM_ST *)&(target->tag), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->tag != TPM_ST_AUTH_SIGNED)
      && (target->tag != TPM_ST_AUTH_SECRET))
        result = TPM_RC_TAG;
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RH_HIERARCHY_Unmarshal((TPMI_RH_HIERARCHY *)&(target->hierarchy), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->digest), buffer, size);
    return result;
}
UINT16
TPMT_TK_AUTH_Marshal(TPMT_TK_AUTH *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ST_Marshal((TPM_ST *)&(source->tag), buffer, size));
    result = (UINT16)(result + TPMI_RH_HIERARCHY_Marshal((TPMI_RH_HIERARCHY *)&(source->hierarchy), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->digest), buffer, size));
    return result;
}

// Table 2:97 - Definition of TPMT_TK_HASHCHECK Structure
TPM_RC
TPMT_TK_HASHCHECK_Unmarshal(TPMT_TK_HASHCHECK *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ST_Unmarshal((TPM_ST *)&(target->tag), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->tag != TPM_ST_HASHCHECK))
        result = TPM_RC_TAG;
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RH_HIERARCHY_Unmarshal((TPMI_RH_HIERARCHY *)&(target->hierarchy), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->digest), buffer, size);
    return result;
}
UINT16
TPMT_TK_HASHCHECK_Marshal(TPMT_TK_HASHCHECK *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ST_Marshal((TPM_ST *)&(source->tag), buffer, size));
    result = (UINT16)(result + TPMI_RH_HIERARCHY_Marshal((TPMI_RH_HIERARCHY *)&(source->hierarchy), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->digest), buffer, size));
    return result;
}

// Table 2:98 - Definition of TPMS_ALG_PROPERTY Structure
UINT16
TPMS_ALG_PROPERTY_Marshal(TPMS_ALG_PROPERTY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ALG_ID_Marshal((TPM_ALG_ID *)&(source->alg), buffer, size));
    result = (UINT16)(result + TPMA_ALGORITHM_Marshal((TPMA_ALGORITHM *)&(source->algProperties), buffer, size));
    return result;
}

// Table 2:99 - Definition of TPMS_TAGGED_PROPERTY Structure
UINT16
TPMS_TAGGED_PROPERTY_Marshal(TPMS_TAGGED_PROPERTY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_PT_Marshal((TPM_PT *)&(source->property), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->value), buffer, size));
    return result;
}

// Table 2:100 - Definition of TPMS_TAGGED_PCR_SELECT Structure
UINT16
TPMS_TAGGED_PCR_SELECT_Marshal(TPMS_TAGGED_PCR_SELECT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_PT_PCR_Marshal((TPM_PT_PCR *)&(source->tag), buffer, size));
    result = (UINT16)(result + UINT8_Marshal((UINT8 *)&(source->sizeofSelect), buffer, size));
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->pcrSelect), buffer, size, (INT32)(source->sizeofSelect)));
    return result;
}

// Table 2:101 - Definition of TPMS_TAGGED_POLICY Structure
UINT16
TPMS_TAGGED_POLICY_Marshal(TPMS_TAGGED_POLICY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_HANDLE_Marshal((TPM_HANDLE *)&(source->handle), buffer, size));
    result = (UINT16)(result + TPMT_HA_Marshal((TPMT_HA *)&(source->policyHash), buffer, size));
    return result;
}

// Table 2:102 - Definition of TPML_CC Structure
TPM_RC
TPML_CC_Unmarshal(TPML_CC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->count), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->count) > MAX_CAP_CC)
            result = TPM_RC_SIZE;
        else
            result = TPM_CC_Array_Unmarshal((TPM_CC *)(target->commandCodes), buffer, size, (INT32)(target->count));
    }
    return result;
}
UINT16
TPML_CC_Marshal(TPML_CC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPM_CC_Array_Marshal((TPM_CC *)(source->commandCodes), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:103 - Definition of TPML_CCA Structure
UINT16
TPML_CCA_Marshal(TPML_CCA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMA_CC_Array_Marshal((TPMA_CC *)(source->commandAttributes), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:104 - Definition of TPML_ALG Structure
TPM_RC
TPML_ALG_Unmarshal(TPML_ALG *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->count), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->count) > MAX_ALG_LIST_SIZE)
            result = TPM_RC_SIZE;
        else
            result = TPM_ALG_ID_Array_Unmarshal((TPM_ALG_ID *)(target->algorithms), buffer, size, (INT32)(target->count));
    }
    return result;
}
UINT16
TPML_ALG_Marshal(TPML_ALG *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPM_ALG_ID_Array_Marshal((TPM_ALG_ID *)(source->algorithms), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:105 - Definition of TPML_HANDLE Structure
UINT16
TPML_HANDLE_Marshal(TPML_HANDLE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPM_HANDLE_Array_Marshal((TPM_HANDLE *)(source->handle), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:106 - Definition of TPML_DIGEST Structure
TPM_RC
TPML_DIGEST_Unmarshal(TPML_DIGEST *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->count), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->count < 2))
        result = TPM_RC_SIZE;
    if(result == TPM_RC_SUCCESS)
    {
        if((target->count) > 8)
            result = TPM_RC_SIZE;
        else
            result = TPM2B_DIGEST_Array_Unmarshal((TPM2B_DIGEST *)(target->digests), buffer, size, (INT32)(target->count));
    }
    return result;
}
UINT16
TPML_DIGEST_Marshal(TPML_DIGEST *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Array_Marshal((TPM2B_DIGEST *)(source->digests), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:107 - Definition of TPML_DIGEST_VALUES Structure
TPM_RC
TPML_DIGEST_VALUES_Unmarshal(TPML_DIGEST_VALUES *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->count), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->count) > HASH_COUNT)
            result = TPM_RC_SIZE;
        else
            result = TPMT_HA_Array_Unmarshal((TPMT_HA *)(target->digests), buffer, size, 0, (INT32)(target->count));
    }
    return result;
}
UINT16
TPML_DIGEST_VALUES_Marshal(TPML_DIGEST_VALUES *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMT_HA_Array_Marshal((TPMT_HA *)(source->digests), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:108 - Definition of TPML_PCR_SELECTION Structure
TPM_RC
TPML_PCR_SELECTION_Unmarshal(TPML_PCR_SELECTION *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->count), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->count) > HASH_COUNT)
            result = TPM_RC_SIZE;
        else
            result = TPMS_PCR_SELECTION_Array_Unmarshal((TPMS_PCR_SELECTION *)(target->pcrSelections), buffer, size, (INT32)(target->count));
    }
    return result;
}
UINT16
TPML_PCR_SELECTION_Marshal(TPML_PCR_SELECTION *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_PCR_SELECTION_Array_Marshal((TPMS_PCR_SELECTION *)(source->pcrSelections), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:109 - Definition of TPML_ALG_PROPERTY Structure
UINT16
TPML_ALG_PROPERTY_Marshal(TPML_ALG_PROPERTY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_ALG_PROPERTY_Array_Marshal((TPMS_ALG_PROPERTY *)(source->algProperties), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:110 - Definition of TPML_TAGGED_TPM_PROPERTY Structure
UINT16
TPML_TAGGED_TPM_PROPERTY_Marshal(TPML_TAGGED_TPM_PROPERTY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_TAGGED_PROPERTY_Array_Marshal((TPMS_TAGGED_PROPERTY *)(source->tpmProperty), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:111 - Definition of TPML_TAGGED_PCR_PROPERTY Structure
UINT16
TPML_TAGGED_PCR_PROPERTY_Marshal(TPML_TAGGED_PCR_PROPERTY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_TAGGED_PCR_SELECT_Array_Marshal((TPMS_TAGGED_PCR_SELECT *)(source->pcrProperty), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:112 - Definition of TPML_ECC_CURVE Structure
#if ALG_ECC
UINT16
TPML_ECC_CURVE_Marshal(TPML_ECC_CURVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPM_ECC_CURVE_Array_Marshal((TPM_ECC_CURVE *)(source->eccCurves), buffer, size, (INT32)(source->count)));
    return result;
}
#endif // ALG_ECC

// Table 2:113 - Definition of TPML_TAGGED_POLICY Structure
UINT16
TPML_TAGGED_POLICY_Marshal(TPML_TAGGED_POLICY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_TAGGED_POLICY_Array_Marshal((TPMS_TAGGED_POLICY *)(source->policies), buffer, size, (INT32)(source->count)));
    return result;
}

// Table 2:114 - Definition of TPMU_CAPABILITIES Union
UINT16
TPMU_CAPABILITIES_Marshal(TPMU_CAPABILITIES *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
        case TPM_CAP_ALGS:
            return TPML_ALG_PROPERTY_Marshal((TPML_ALG_PROPERTY *)&(source->algorithms), buffer, size);
        case TPM_CAP_HANDLES:
            return TPML_HANDLE_Marshal((TPML_HANDLE *)&(source->handles), buffer, size);
        case TPM_CAP_COMMANDS:
            return TPML_CCA_Marshal((TPML_CCA *)&(source->command), buffer, size);
        case TPM_CAP_PP_COMMANDS:
            return TPML_CC_Marshal((TPML_CC *)&(source->ppCommands), buffer, size);
        case TPM_CAP_AUDIT_COMMANDS:
            return TPML_CC_Marshal((TPML_CC *)&(source->auditCommands), buffer, size);
        case TPM_CAP_PCRS:
            return TPML_PCR_SELECTION_Marshal((TPML_PCR_SELECTION *)&(source->assignedPCR), buffer, size);
        case TPM_CAP_TPM_PROPERTIES:
            return TPML_TAGGED_TPM_PROPERTY_Marshal((TPML_TAGGED_TPM_PROPERTY *)&(source->tpmProperties), buffer, size);
        case TPM_CAP_PCR_PROPERTIES:
            return TPML_TAGGED_PCR_PROPERTY_Marshal((TPML_TAGGED_PCR_PROPERTY *)&(source->pcrProperties), buffer, size);
#if ALG_ECC
        case TPM_CAP_ECC_CURVES:
            return TPML_ECC_CURVE_Marshal((TPML_ECC_CURVE *)&(source->eccCurves), buffer, size);
#endif // ALG_ECC
        case TPM_CAP_AUTH_POLICIES:
            return TPML_TAGGED_POLICY_Marshal((TPML_TAGGED_POLICY *)&(source->authPolicies), buffer, size);
    }
    return 0;
}

// Table 2:115 - Definition of TPMS_CAPABILITY_DATA Structure
UINT16
TPMS_CAPABILITY_DATA_Marshal(TPMS_CAPABILITY_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_CAP_Marshal((TPM_CAP *)&(source->capability), buffer, size));
    result = (UINT16)(result + TPMU_CAPABILITIES_Marshal((TPMU_CAPABILITIES *)&(source->data), buffer, size, (UINT32)source->capability));
    return result;
}

// Table 2:116 - Definition of TPMS_CLOCK_INFO Structure
TPM_RC
TPMS_CLOCK_INFO_Unmarshal(TPMS_CLOCK_INFO *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT64_Unmarshal((UINT64 *)&(target->clock), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = UINT32_Unmarshal((UINT32 *)&(target->resetCount), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = UINT32_Unmarshal((UINT32 *)&(target->restartCount), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_YES_NO_Unmarshal((TPMI_YES_NO *)&(target->safe), buffer, size);
    return result;
}
UINT16
TPMS_CLOCK_INFO_Marshal(TPMS_CLOCK_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->clock), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->resetCount), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->restartCount), buffer, size));
    result = (UINT16)(result + TPMI_YES_NO_Marshal((TPMI_YES_NO *)&(source->safe), buffer, size));
    return result;
}

// Table 2:117 - Definition of TPMS_TIME_INFO Structure
TPM_RC
TPMS_TIME_INFO_Unmarshal(TPMS_TIME_INFO *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT64_Unmarshal((UINT64 *)&(target->time), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMS_CLOCK_INFO_Unmarshal((TPMS_CLOCK_INFO *)&(target->clockInfo), buffer, size);
    return result;
}
UINT16
TPMS_TIME_INFO_Marshal(TPMS_TIME_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->time), buffer, size));
    result = (UINT16)(result + TPMS_CLOCK_INFO_Marshal((TPMS_CLOCK_INFO *)&(source->clockInfo), buffer, size));
    return result;
}

// Table 2:118 - Definition of TPMS_TIME_ATTEST_INFO Structure
UINT16
TPMS_TIME_ATTEST_INFO_Marshal(TPMS_TIME_ATTEST_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMS_TIME_INFO_Marshal((TPMS_TIME_INFO *)&(source->time), buffer, size));
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->firmwareVersion), buffer, size));
    return result;
}

// Table 2:119 - Definition of TPMS_CERTIFY_INFO Structure
UINT16
TPMS_CERTIFY_INFO_Marshal(TPMS_CERTIFY_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->name), buffer, size));
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->qualifiedName), buffer, size));
    return result;
}

// Table 2:120 - Definition of TPMS_QUOTE_INFO Structure
UINT16
TPMS_QUOTE_INFO_Marshal(TPMS_QUOTE_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPML_PCR_SELECTION_Marshal((TPML_PCR_SELECTION *)&(source->pcrSelect), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->pcrDigest), buffer, size));
    return result;
}

// Table 2:121 - Definition of TPMS_COMMAND_AUDIT_INFO Structure
UINT16
TPMS_COMMAND_AUDIT_INFO_Marshal(TPMS_COMMAND_AUDIT_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->auditCounter), buffer, size));
    result = (UINT16)(result + TPM_ALG_ID_Marshal((TPM_ALG_ID *)&(source->digestAlg), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->auditDigest), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->commandDigest), buffer, size));
    return result;
}

// Table 2:122 - Definition of TPMS_SESSION_AUDIT_INFO Structure
UINT16
TPMS_SESSION_AUDIT_INFO_Marshal(TPMS_SESSION_AUDIT_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_YES_NO_Marshal((TPMI_YES_NO *)&(source->exclusiveSession), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->sessionDigest), buffer, size));
    return result;
}

// Table 2:123 - Definition of TPMS_CREATION_INFO Structure
UINT16
TPMS_CREATION_INFO_Marshal(TPMS_CREATION_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->objectName), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->creationHash), buffer, size));
    return result;
}

// Table 2:124 - Definition of TPMS_NV_CERTIFY_INFO Structure
UINT16
TPMS_NV_CERTIFY_INFO_Marshal(TPMS_NV_CERTIFY_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->indexName), buffer, size));
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->offset), buffer, size));
    result = (UINT16)(result + TPM2B_MAX_NV_BUFFER_Marshal((TPM2B_MAX_NV_BUFFER *)&(source->nvContents), buffer, size));
    return result;
}

// Table 2:125 - Definition of TPMS_NV_DIGEST_CERTIFY_INFO Structure
UINT16
TPMS_NV_DIGEST_CERTIFY_INFO_Marshal(TPMS_NV_DIGEST_CERTIFY_INFO *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->indexName), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->nvDigest), buffer, size));
    return result;
}

// Table 2:126 - Definition of TPMI_ST_ATTEST Type
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ST_ATTEST_Marshal(TPMI_ST_ATTEST *source, BYTE **buffer, INT32 *size)
{
    return TPM_ST_Marshal((TPM_ST *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:127 - Definition of TPMU_ATTEST Union
UINT16
TPMU_ATTEST_Marshal(TPMU_ATTEST *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
        case TPM_ST_ATTEST_CERTIFY:
            return TPMS_CERTIFY_INFO_Marshal((TPMS_CERTIFY_INFO *)&(source->certify), buffer, size);
        case TPM_ST_ATTEST_CREATION:
            return TPMS_CREATION_INFO_Marshal((TPMS_CREATION_INFO *)&(source->creation), buffer, size);
        case TPM_ST_ATTEST_QUOTE:
            return TPMS_QUOTE_INFO_Marshal((TPMS_QUOTE_INFO *)&(source->quote), buffer, size);
        case TPM_ST_ATTEST_COMMAND_AUDIT:
            return TPMS_COMMAND_AUDIT_INFO_Marshal((TPMS_COMMAND_AUDIT_INFO *)&(source->commandAudit), buffer, size);
        case TPM_ST_ATTEST_SESSION_AUDIT:
            return TPMS_SESSION_AUDIT_INFO_Marshal((TPMS_SESSION_AUDIT_INFO *)&(source->sessionAudit), buffer, size);
        case TPM_ST_ATTEST_TIME:
            return TPMS_TIME_ATTEST_INFO_Marshal((TPMS_TIME_ATTEST_INFO *)&(source->time), buffer, size);
        case TPM_ST_ATTEST_NV:
            return TPMS_NV_CERTIFY_INFO_Marshal((TPMS_NV_CERTIFY_INFO *)&(source->nv), buffer, size);
        case TPM_ST_ATTEST_NV_DIGEST:
            return TPMS_NV_DIGEST_CERTIFY_INFO_Marshal((TPMS_NV_DIGEST_CERTIFY_INFO *)&(source->nvDigest), buffer, size);
    }
    return 0;
}

// Table 2:128 - Definition of TPMS_ATTEST Structure
UINT16
TPMS_ATTEST_Marshal(TPMS_ATTEST *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_GENERATED_Marshal((TPM_GENERATED *)&(source->magic), buffer, size));
    result = (UINT16)(result + TPMI_ST_ATTEST_Marshal((TPMI_ST_ATTEST *)&(source->type), buffer, size));
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->qualifiedSigner), buffer, size));
    result = (UINT16)(result + TPM2B_DATA_Marshal((TPM2B_DATA *)&(source->extraData), buffer, size));
    result = (UINT16)(result + TPMS_CLOCK_INFO_Marshal((TPMS_CLOCK_INFO *)&(source->clockInfo), buffer, size));
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->firmwareVersion), buffer, size));
    result = (UINT16)(result + TPMU_ATTEST_Marshal((TPMU_ATTEST *)&(source->attested), buffer, size, (UINT32)source->type));
    return result;
}

// Table 2:129 - Definition of TPM2B_ATTEST Structure
UINT16
TPM2B_ATTEST_Marshal(TPM2B_ATTEST *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.attestationData), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:130 - Definition of TPMS_AUTH_COMMAND Structure
TPM_RC
TPMS_AUTH_COMMAND_Unmarshal(TPMS_AUTH_COMMAND *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_SH_AUTH_SESSION_Unmarshal((TPMI_SH_AUTH_SESSION *)&(target->sessionHandle), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_NONCE_Unmarshal((TPM2B_NONCE *)&(target->nonce), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMA_SESSION_Unmarshal((TPMA_SESSION *)&(target->sessionAttributes), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_AUTH_Unmarshal((TPM2B_AUTH *)&(target->hmac), buffer, size);
    return result;
}

// Table 2:131 - Definition of TPMS_AUTH_RESPONSE Structure
UINT16
TPMS_AUTH_RESPONSE_Marshal(TPMS_AUTH_RESPONSE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_NONCE_Marshal((TPM2B_NONCE *)&(source->nonce), buffer, size));
    result = (UINT16)(result + TPMA_SESSION_Marshal((TPMA_SESSION *)&(source->sessionAttributes), buffer, size));
    result = (UINT16)(result + TPM2B_AUTH_Marshal((TPM2B_AUTH *)&(source->hmac), buffer, size));
    return result;
}

// Table 2:132 - Definition of TPMI_TDES_KEY_BITS Type
#if ALG_TDES
TPM_RC
TPMI_TDES_KEY_BITS_Unmarshal(TPMI_TDES_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_KEY_BITS_Unmarshal((TPM_KEY_BITS *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if TDES_128
            case 128:
#endif // TDES_128
#if TDES_192
            case 192:
#endif // TDES_192
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_TDES_KEY_BITS_Marshal(TPMI_TDES_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_TDES

// Table 2:132 - Definition of TPMI_AES_KEY_BITS Type
#if ALG_AES
TPM_RC
TPMI_AES_KEY_BITS_Unmarshal(TPMI_AES_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_KEY_BITS_Unmarshal((TPM_KEY_BITS *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if AES_128
            case 128:
#endif // AES_128
#if AES_192
            case 192:
#endif // AES_192
#if AES_256
            case 256:
#endif // AES_256
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_AES_KEY_BITS_Marshal(TPMI_AES_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_AES

// Table 2:132 - Definition of TPMI_SM4_KEY_BITS Type
#if ALG_SM4
TPM_RC
TPMI_SM4_KEY_BITS_Unmarshal(TPMI_SM4_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_KEY_BITS_Unmarshal((TPM_KEY_BITS *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if SM4_128
            case 128:
#endif // SM4_128
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SM4_KEY_BITS_Marshal(TPMI_SM4_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_SM4

// Table 2:132 - Definition of TPMI_CAMELLIA_KEY_BITS Type
#if ALG_CAMELLIA
TPM_RC
TPMI_CAMELLIA_KEY_BITS_Unmarshal(TPMI_CAMELLIA_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_KEY_BITS_Unmarshal((TPM_KEY_BITS *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if CAMELLIA_128
            case 128:
#endif // CAMELLIA_128
#if CAMELLIA_192
            case 192:
#endif // CAMELLIA_192
#if CAMELLIA_256
            case 256:
#endif // CAMELLIA_256
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_CAMELLIA_KEY_BITS_Marshal(TPMI_CAMELLIA_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_CAMELLIA

// Table 2:133 - Definition of TPMU_SYM_KEY_BITS Union
TPM_RC
TPMU_SYM_KEY_BITS_Unmarshal(TPMU_SYM_KEY_BITS *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_TDES
        case ALG_TDES_VALUE:
            return TPMI_TDES_KEY_BITS_Unmarshal((TPMI_TDES_KEY_BITS *)&(target->tdes), buffer, size);
#endif // ALG_TDES
#if ALG_AES
        case ALG_AES_VALUE:
            return TPMI_AES_KEY_BITS_Unmarshal((TPMI_AES_KEY_BITS *)&(target->aes), buffer, size);
#endif // ALG_AES
#if ALG_SM4
        case ALG_SM4_VALUE:
            return TPMI_SM4_KEY_BITS_Unmarshal((TPMI_SM4_KEY_BITS *)&(target->sm4), buffer, size);
#endif // ALG_SM4
#if ALG_CAMELLIA
        case ALG_CAMELLIA_VALUE:
            return TPMI_CAMELLIA_KEY_BITS_Unmarshal((TPMI_CAMELLIA_KEY_BITS *)&(target->camellia), buffer, size);
#endif // ALG_CAMELLIA
#if ALG_XOR
        case ALG_XOR_VALUE:
            return TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->xor), buffer, size, 0);
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SYM_KEY_BITS_Marshal(TPMU_SYM_KEY_BITS *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_TDES
        case ALG_TDES_VALUE:
            return TPMI_TDES_KEY_BITS_Marshal((TPMI_TDES_KEY_BITS *)&(source->tdes), buffer, size);
#endif // ALG_TDES
#if ALG_AES
        case ALG_AES_VALUE:
            return TPMI_AES_KEY_BITS_Marshal((TPMI_AES_KEY_BITS *)&(source->aes), buffer, size);
#endif // ALG_AES
#if ALG_SM4
        case ALG_SM4_VALUE:
            return TPMI_SM4_KEY_BITS_Marshal((TPMI_SM4_KEY_BITS *)&(source->sm4), buffer, size);
#endif // ALG_SM4
#if ALG_CAMELLIA
        case ALG_CAMELLIA_VALUE:
            return TPMI_CAMELLIA_KEY_BITS_Marshal((TPMI_CAMELLIA_KEY_BITS *)&(source->camellia), buffer, size);
#endif // ALG_CAMELLIA
#if ALG_XOR
        case ALG_XOR_VALUE:
            return TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->xor), buffer, size);
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:134 - Definition of TPMU_SYM_MODE Union
TPM_RC
TPMU_SYM_MODE_Unmarshal(TPMU_SYM_MODE *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_TDES
        case ALG_TDES_VALUE:
            return TPMI_ALG_SYM_MODE_Unmarshal((TPMI_ALG_SYM_MODE *)&(target->tdes), buffer, size, 1);
#endif // ALG_TDES
#if ALG_AES
        case ALG_AES_VALUE:
            return TPMI_ALG_SYM_MODE_Unmarshal((TPMI_ALG_SYM_MODE *)&(target->aes), buffer, size, 1);
#endif // ALG_AES
#if ALG_SM4
        case ALG_SM4_VALUE:
            return TPMI_ALG_SYM_MODE_Unmarshal((TPMI_ALG_SYM_MODE *)&(target->sm4), buffer, size, 1);
#endif // ALG_SM4
#if ALG_CAMELLIA
        case ALG_CAMELLIA_VALUE:
            return TPMI_ALG_SYM_MODE_Unmarshal((TPMI_ALG_SYM_MODE *)&(target->camellia), buffer, size, 1);
#endif // ALG_CAMELLIA
#if ALG_XOR
        case ALG_XOR_VALUE:
            return TPM_RC_SUCCESS;
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SYM_MODE_Marshal(TPMU_SYM_MODE *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_TDES
        case ALG_TDES_VALUE:
            return TPMI_ALG_SYM_MODE_Marshal((TPMI_ALG_SYM_MODE *)&(source->tdes), buffer, size);
#endif // ALG_TDES
#if ALG_AES
        case ALG_AES_VALUE:
            return TPMI_ALG_SYM_MODE_Marshal((TPMI_ALG_SYM_MODE *)&(source->aes), buffer, size);
#endif // ALG_AES
#if ALG_SM4
        case ALG_SM4_VALUE:
            return TPMI_ALG_SYM_MODE_Marshal((TPMI_ALG_SYM_MODE *)&(source->sm4), buffer, size);
#endif // ALG_SM4
#if ALG_CAMELLIA
        case ALG_CAMELLIA_VALUE:
            return TPMI_ALG_SYM_MODE_Marshal((TPMI_ALG_SYM_MODE *)&(source->camellia), buffer, size);
#endif // ALG_CAMELLIA
#if ALG_XOR
        case ALG_XOR_VALUE:
            return 0;
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:136 - Definition of TPMT_SYM_DEF Structure
TPM_RC
TPMT_SYM_DEF_Unmarshal(TPMT_SYM_DEF *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_SYM_Unmarshal((TPMI_ALG_SYM *)&(target->algorithm), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SYM_KEY_BITS_Unmarshal((TPMU_SYM_KEY_BITS *)&(target->keyBits), buffer, size, (UINT32)target->algorithm);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SYM_MODE_Unmarshal((TPMU_SYM_MODE *)&(target->mode), buffer, size, (UINT32)target->algorithm);
    return result;
}
UINT16
TPMT_SYM_DEF_Marshal(TPMT_SYM_DEF *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_SYM_Marshal((TPMI_ALG_SYM *)&(source->algorithm), buffer, size));
    result = (UINT16)(result + TPMU_SYM_KEY_BITS_Marshal((TPMU_SYM_KEY_BITS *)&(source->keyBits), buffer, size, (UINT32)source->algorithm));
    result = (UINT16)(result + TPMU_SYM_MODE_Marshal((TPMU_SYM_MODE *)&(source->mode), buffer, size, (UINT32)source->algorithm));
    return result;
}

// Table 2:137 - Definition of TPMT_SYM_DEF_OBJECT Structure
TPM_RC
TPMT_SYM_DEF_OBJECT_Unmarshal(TPMT_SYM_DEF_OBJECT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_SYM_OBJECT_Unmarshal((TPMI_ALG_SYM_OBJECT *)&(target->algorithm), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SYM_KEY_BITS_Unmarshal((TPMU_SYM_KEY_BITS *)&(target->keyBits), buffer, size, (UINT32)target->algorithm);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SYM_MODE_Unmarshal((TPMU_SYM_MODE *)&(target->mode), buffer, size, (UINT32)target->algorithm);
    return result;
}
UINT16
TPMT_SYM_DEF_OBJECT_Marshal(TPMT_SYM_DEF_OBJECT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_SYM_OBJECT_Marshal((TPMI_ALG_SYM_OBJECT *)&(source->algorithm), buffer, size));
    result = (UINT16)(result + TPMU_SYM_KEY_BITS_Marshal((TPMU_SYM_KEY_BITS *)&(source->keyBits), buffer, size, (UINT32)source->algorithm));
    result = (UINT16)(result + TPMU_SYM_MODE_Marshal((TPMU_SYM_MODE *)&(source->mode), buffer, size, (UINT32)source->algorithm));
    return result;
}

// Table 2:138 - Definition of TPM2B_SYM_KEY Structure
TPM_RC
TPM2B_SYM_KEY_Unmarshal(TPM2B_SYM_KEY *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_SYM_KEY_BYTES)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_SYM_KEY_Marshal(TPM2B_SYM_KEY *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:139 - Definition of TPMS_SYMCIPHER_PARMS Structure
TPM_RC
TPMS_SYMCIPHER_PARMS_Unmarshal(TPMS_SYMCIPHER_PARMS *target, BYTE **buffer, INT32 *size)
{
    return TPMT_SYM_DEF_OBJECT_Unmarshal((TPMT_SYM_DEF_OBJECT *)&(target->sym), buffer, size, 0);
}
UINT16
TPMS_SYMCIPHER_PARMS_Marshal(TPMS_SYMCIPHER_PARMS *source, BYTE **buffer, INT32 *size)
{
    return TPMT_SYM_DEF_OBJECT_Marshal((TPMT_SYM_DEF_OBJECT *)&(source->sym), buffer, size);
}

// Table 2:140 - Definition of TPM2B_LABEL Structure
TPM_RC
TPM2B_LABEL_Unmarshal(TPM2B_LABEL *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > LABEL_MAX_BUFFER)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_LABEL_Marshal(TPM2B_LABEL *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:141 - Definition of TPMS_DERIVE Structure
TPM_RC
TPMS_DERIVE_Unmarshal(TPMS_DERIVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM2B_LABEL_Unmarshal((TPM2B_LABEL *)&(target->label), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_LABEL_Unmarshal((TPM2B_LABEL *)&(target->context), buffer, size);
    return result;
}
UINT16
TPMS_DERIVE_Marshal(TPMS_DERIVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_LABEL_Marshal((TPM2B_LABEL *)&(source->label), buffer, size));
    result = (UINT16)(result + TPM2B_LABEL_Marshal((TPM2B_LABEL *)&(source->context), buffer, size));
    return result;
}

// Table 2:142 - Definition of TPM2B_DERIVE Structure
TPM_RC
TPM2B_DERIVE_Unmarshal(TPM2B_DERIVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMS_DERIVE))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_DERIVE_Marshal(TPM2B_DERIVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:143 - Definition of TPMU_SENSITIVE_CREATE Union
// Table 2:144 - Definition of TPM2B_SENSITIVE_DATA Structure
TPM_RC
TPM2B_SENSITIVE_DATA_Unmarshal(TPM2B_SENSITIVE_DATA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMU_SENSITIVE_CREATE))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_SENSITIVE_DATA_Marshal(TPM2B_SENSITIVE_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:145 - Definition of TPMS_SENSITIVE_CREATE Structure
TPM_RC
TPMS_SENSITIVE_CREATE_Unmarshal(TPMS_SENSITIVE_CREATE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM2B_AUTH_Unmarshal((TPM2B_AUTH *)&(target->userAuth), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_SENSITIVE_DATA_Unmarshal((TPM2B_SENSITIVE_DATA *)&(target->data), buffer, size);
    return result;
}

// Table 2:146 - Definition of TPM2B_SENSITIVE_CREATE Structure
TPM_RC
TPM2B_SENSITIVE_CREATE_Unmarshal(TPM2B_SENSITIVE_CREATE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->size), buffer, size); // =a
    if(result == TPM_RC_SUCCESS)
    {
        // if size is zero, then the required structure is missing
        if(target->size == 0)
            result = TPM_RC_SIZE;
        else
        {
            INT32   startSize = *size;
            result = TPMS_SENSITIVE_CREATE_Unmarshal((TPMS_SENSITIVE_CREATE *)&(target->sensitive), buffer, size); // =b
            if(result == TPM_RC_SUCCESS)
            {
                if(target->size != (startSize - *size))
                    result = TPM_RC_SIZE;
            }
        }
    }
    return result;
}

// Table 2:147 - Definition of TPMS_SCHEME_HASH Structure
TPM_RC
TPMS_SCHEME_HASH_Unmarshal(TPMS_SCHEME_HASH *target, BYTE **buffer, INT32 *size)
{
    return TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hashAlg), buffer, size, 0);
}
UINT16
TPMS_SCHEME_HASH_Marshal(TPMS_SCHEME_HASH *source, BYTE **buffer, INT32 *size)
{
    return TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hashAlg), buffer, size);
}

// Table 2:148 - Definition of TPMS_SCHEME_ECDAA Structure
#if ALG_ECC
TPM_RC
TPMS_SCHEME_ECDAA_Unmarshal(TPMS_SCHEME_ECDAA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hashAlg), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = UINT16_Unmarshal((UINT16 *)&(target->count), buffer, size);
    return result;
}
UINT16
TPMS_SCHEME_ECDAA_Marshal(TPMS_SCHEME_ECDAA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hashAlg), buffer, size));
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->count), buffer, size));
    return result;
}
#endif // ALG_ECC

// Table 2:149 - Definition of TPMI_ALG_KEYEDHASH_SCHEME Type
TPM_RC
TPMI_ALG_KEYEDHASH_SCHEME_Unmarshal(TPMI_ALG_KEYEDHASH_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_HMAC
            case ALG_HMAC_VALUE:
#endif // ALG_HMAC
#if ALG_XOR
            case ALG_XOR_VALUE:
#endif // ALG_XOR
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_KEYEDHASH_SCHEME_Marshal(TPMI_ALG_KEYEDHASH_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:150 - Definition of Types for HMAC_SIG_SCHEME
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_HMAC_Unmarshal(TPMS_SCHEME_HMAC *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SCHEME_HMAC_Marshal(TPMS_SCHEME_HMAC *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:151 - Definition of TPMS_SCHEME_XOR Structure
TPM_RC
TPMS_SCHEME_XOR_Unmarshal(TPMS_SCHEME_XOR *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hashAlg), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_ALG_KDF_Unmarshal((TPMI_ALG_KDF *)&(target->kdf), buffer, size, 1);
    return result;
}
UINT16
TPMS_SCHEME_XOR_Marshal(TPMS_SCHEME_XOR *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hashAlg), buffer, size));
    result = (UINT16)(result + TPMI_ALG_KDF_Marshal((TPMI_ALG_KDF *)&(source->kdf), buffer, size));
    return result;
}

// Table 2:152 - Definition of TPMU_SCHEME_KEYEDHASH Union
TPM_RC
TPMU_SCHEME_KEYEDHASH_Unmarshal(TPMU_SCHEME_KEYEDHASH *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMS_SCHEME_HMAC_Unmarshal((TPMS_SCHEME_HMAC *)&(target->hmac), buffer, size);
#endif // ALG_HMAC
#if ALG_XOR
        case ALG_XOR_VALUE:
            return TPMS_SCHEME_XOR_Unmarshal((TPMS_SCHEME_XOR *)&(target->xor), buffer, size);
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SCHEME_KEYEDHASH_Marshal(TPMU_SCHEME_KEYEDHASH *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMS_SCHEME_HMAC_Marshal((TPMS_SCHEME_HMAC *)&(source->hmac), buffer, size);
#endif // ALG_HMAC
#if ALG_XOR
        case ALG_XOR_VALUE:
            return TPMS_SCHEME_XOR_Marshal((TPMS_SCHEME_XOR *)&(source->xor), buffer, size);
#endif // ALG_XOR
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:153 - Definition of TPMT_KEYEDHASH_SCHEME Structure
TPM_RC
TPMT_KEYEDHASH_SCHEME_Unmarshal(TPMT_KEYEDHASH_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_KEYEDHASH_SCHEME_Unmarshal((TPMI_ALG_KEYEDHASH_SCHEME *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SCHEME_KEYEDHASH_Unmarshal((TPMU_SCHEME_KEYEDHASH *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_KEYEDHASH_SCHEME_Marshal(TPMT_KEYEDHASH_SCHEME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_KEYEDHASH_SCHEME_Marshal((TPMI_ALG_KEYEDHASH_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_SCHEME_KEYEDHASH_Marshal((TPMU_SCHEME_KEYEDHASH *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}

// Table 2:154 - Definition of Types for RSA Signature Schemes
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_RSASSA_Unmarshal(TPMS_SIG_SCHEME_RSASSA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_RSASSA_Marshal(TPMS_SIG_SCHEME_RSASSA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SIG_SCHEME_RSAPSS_Unmarshal(TPMS_SIG_SCHEME_RSAPSS *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_RSAPSS_Marshal(TPMS_SIG_SCHEME_RSAPSS *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:155 - Definition of Types for ECC Signature Schemes
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_ECDSA_Unmarshal(TPMS_SIG_SCHEME_ECDSA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_ECDSA_Marshal(TPMS_SIG_SCHEME_ECDSA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SIG_SCHEME_SM2_Unmarshal(TPMS_SIG_SCHEME_SM2 *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_SM2_Marshal(TPMS_SIG_SCHEME_SM2 *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SIG_SCHEME_ECSCHNORR_Unmarshal(TPMS_SIG_SCHEME_ECSCHNORR *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_ECSCHNORR_Marshal(TPMS_SIG_SCHEME_ECSCHNORR *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SIG_SCHEME_ECDAA_Unmarshal(TPMS_SIG_SCHEME_ECDAA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_ECDAA_Unmarshal((TPMS_SCHEME_ECDAA *)target, buffer, size);
}
UINT16
TPMS_SIG_SCHEME_ECDAA_Marshal(TPMS_SIG_SCHEME_ECDAA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_ECDAA_Marshal((TPMS_SCHEME_ECDAA *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:156 - Definition of TPMU_SIG_SCHEME Union
TPM_RC
TPMU_SIG_SCHEME_Unmarshal(TPMU_SIG_SCHEME *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIG_SCHEME_ECDAA_Unmarshal((TPMS_SIG_SCHEME_ECDAA *)&(target->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIG_SCHEME_RSASSA_Unmarshal((TPMS_SIG_SCHEME_RSASSA *)&(target->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIG_SCHEME_RSAPSS_Unmarshal((TPMS_SIG_SCHEME_RSAPSS *)&(target->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIG_SCHEME_ECDSA_Unmarshal((TPMS_SIG_SCHEME_ECDSA *)&(target->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIG_SCHEME_SM2_Unmarshal((TPMS_SIG_SCHEME_SM2 *)&(target->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIG_SCHEME_ECSCHNORR_Unmarshal((TPMS_SIG_SCHEME_ECSCHNORR *)&(target->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMS_SCHEME_HMAC_Unmarshal((TPMS_SCHEME_HMAC *)&(target->hmac), buffer, size);
#endif // ALG_HMAC
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SIG_SCHEME_Marshal(TPMU_SIG_SCHEME *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIG_SCHEME_ECDAA_Marshal((TPMS_SIG_SCHEME_ECDAA *)&(source->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIG_SCHEME_RSASSA_Marshal((TPMS_SIG_SCHEME_RSASSA *)&(source->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIG_SCHEME_RSAPSS_Marshal((TPMS_SIG_SCHEME_RSAPSS *)&(source->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIG_SCHEME_ECDSA_Marshal((TPMS_SIG_SCHEME_ECDSA *)&(source->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIG_SCHEME_SM2_Marshal((TPMS_SIG_SCHEME_SM2 *)&(source->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIG_SCHEME_ECSCHNORR_Marshal((TPMS_SIG_SCHEME_ECSCHNORR *)&(source->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMS_SCHEME_HMAC_Marshal((TPMS_SCHEME_HMAC *)&(source->hmac), buffer, size);
#endif // ALG_HMAC
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:157 - Definition of TPMT_SIG_SCHEME Structure
TPM_RC
TPMT_SIG_SCHEME_Unmarshal(TPMT_SIG_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_SIG_SCHEME_Unmarshal((TPMI_ALG_SIG_SCHEME *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SIG_SCHEME_Unmarshal((TPMU_SIG_SCHEME *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_SIG_SCHEME_Marshal(TPMT_SIG_SCHEME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_SIG_SCHEME_Marshal((TPMI_ALG_SIG_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_SIG_SCHEME_Marshal((TPMU_SIG_SCHEME *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}

// Table 2:158 - Definition of Types for Encryption Schemes
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_ENC_SCHEME_OAEP_Unmarshal(TPMS_ENC_SCHEME_OAEP *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_ENC_SCHEME_OAEP_Marshal(TPMS_ENC_SCHEME_OAEP *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_ENC_SCHEME_RSAES_Unmarshal(TPMS_ENC_SCHEME_RSAES *target, BYTE **buffer, INT32 *size)
{
    return TPMS_EMPTY_Unmarshal((TPMS_EMPTY *)target, buffer, size);
}
UINT16
TPMS_ENC_SCHEME_RSAES_Marshal(TPMS_ENC_SCHEME_RSAES *source, BYTE **buffer, INT32 *size)
{
    return TPMS_EMPTY_Marshal((TPMS_EMPTY *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:159 - Definition of Types for ECC Key Exchange
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_KEY_SCHEME_ECDH_Unmarshal(TPMS_KEY_SCHEME_ECDH *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_KEY_SCHEME_ECDH_Marshal(TPMS_KEY_SCHEME_ECDH *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_KEY_SCHEME_ECMQV_Unmarshal(TPMS_KEY_SCHEME_ECMQV *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_KEY_SCHEME_ECMQV_Marshal(TPMS_KEY_SCHEME_ECMQV *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:160 - Definition of Types for KDF Schemes
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_MGF1_Unmarshal(TPMS_SCHEME_MGF1 *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SCHEME_MGF1_Marshal(TPMS_SCHEME_MGF1 *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SCHEME_KDF1_SP800_56A_Unmarshal(TPMS_SCHEME_KDF1_SP800_56A *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SCHEME_KDF1_SP800_56A_Marshal(TPMS_SCHEME_KDF1_SP800_56A *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SCHEME_KDF2_Unmarshal(TPMS_SCHEME_KDF2 *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SCHEME_KDF2_Marshal(TPMS_SCHEME_KDF2 *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
TPM_RC
TPMS_SCHEME_KDF1_SP800_108_Unmarshal(TPMS_SCHEME_KDF1_SP800_108 *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)target, buffer, size);
}
UINT16
TPMS_SCHEME_KDF1_SP800_108_Marshal(TPMS_SCHEME_KDF1_SP800_108 *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:161 - Definition of TPMU_KDF_SCHEME Union
TPM_RC
TPMU_KDF_SCHEME_Unmarshal(TPMU_KDF_SCHEME *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_MGF1
        case ALG_MGF1_VALUE:
            return TPMS_SCHEME_MGF1_Unmarshal((TPMS_SCHEME_MGF1 *)&(target->mgf1), buffer, size);
#endif // ALG_MGF1
#if ALG_KDF1_SP800_56A
        case ALG_KDF1_SP800_56A_VALUE:
            return TPMS_SCHEME_KDF1_SP800_56A_Unmarshal((TPMS_SCHEME_KDF1_SP800_56A *)&(target->kdf1_sp800_56a), buffer, size);
#endif // ALG_KDF1_SP800_56A
#if ALG_KDF2
        case ALG_KDF2_VALUE:
            return TPMS_SCHEME_KDF2_Unmarshal((TPMS_SCHEME_KDF2 *)&(target->kdf2), buffer, size);
#endif // ALG_KDF2
#if ALG_KDF1_SP800_108
        case ALG_KDF1_SP800_108_VALUE:
            return TPMS_SCHEME_KDF1_SP800_108_Unmarshal((TPMS_SCHEME_KDF1_SP800_108 *)&(target->kdf1_sp800_108), buffer, size);
#endif // ALG_KDF1_SP800_108
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_KDF_SCHEME_Marshal(TPMU_KDF_SCHEME *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_MGF1
        case ALG_MGF1_VALUE:
            return TPMS_SCHEME_MGF1_Marshal((TPMS_SCHEME_MGF1 *)&(source->mgf1), buffer, size);
#endif // ALG_MGF1
#if ALG_KDF1_SP800_56A
        case ALG_KDF1_SP800_56A_VALUE:
            return TPMS_SCHEME_KDF1_SP800_56A_Marshal((TPMS_SCHEME_KDF1_SP800_56A *)&(source->kdf1_sp800_56a), buffer, size);
#endif // ALG_KDF1_SP800_56A
#if ALG_KDF2
        case ALG_KDF2_VALUE:
            return TPMS_SCHEME_KDF2_Marshal((TPMS_SCHEME_KDF2 *)&(source->kdf2), buffer, size);
#endif // ALG_KDF2
#if ALG_KDF1_SP800_108
        case ALG_KDF1_SP800_108_VALUE:
            return TPMS_SCHEME_KDF1_SP800_108_Marshal((TPMS_SCHEME_KDF1_SP800_108 *)&(source->kdf1_sp800_108), buffer, size);
#endif // ALG_KDF1_SP800_108
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:162 - Definition of TPMT_KDF_SCHEME Structure
TPM_RC
TPMT_KDF_SCHEME_Unmarshal(TPMT_KDF_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_KDF_Unmarshal((TPMI_ALG_KDF *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_KDF_SCHEME_Unmarshal((TPMU_KDF_SCHEME *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_KDF_SCHEME_Marshal(TPMT_KDF_SCHEME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_KDF_Marshal((TPMI_ALG_KDF *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_KDF_SCHEME_Marshal((TPMU_KDF_SCHEME *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}

// Table 2:163 - Definition of TPMI_ALG_ASYM_SCHEME Type
TPM_RC
TPMI_ALG_ASYM_SCHEME_Unmarshal(TPMI_ALG_ASYM_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_ECDH
            case ALG_ECDH_VALUE:
#endif // ALG_ECDH
#if ALG_ECMQV
            case ALG_ECMQV_VALUE:
#endif // ALG_ECMQV
#if ALG_ECDAA
            case ALG_ECDAA_VALUE:
#endif // ALG_ECDAA
#if ALG_RSASSA
            case ALG_RSASSA_VALUE:
#endif // ALG_RSASSA
#if ALG_RSAPSS
            case ALG_RSAPSS_VALUE:
#endif // ALG_RSAPSS
#if ALG_ECDSA
            case ALG_ECDSA_VALUE:
#endif // ALG_ECDSA
#if ALG_SM2
            case ALG_SM2_VALUE:
#endif // ALG_SM2
#if ALG_ECSCHNORR
            case ALG_ECSCHNORR_VALUE:
#endif // ALG_ECSCHNORR
#if ALG_RSAES
            case ALG_RSAES_VALUE:
#endif // ALG_RSAES
#if ALG_OAEP
            case ALG_OAEP_VALUE:
#endif // ALG_OAEP
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ASYM_SCHEME_Marshal(TPMI_ALG_ASYM_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:164 - Definition of TPMU_ASYM_SCHEME Union
TPM_RC
TPMU_ASYM_SCHEME_Unmarshal(TPMU_ASYM_SCHEME *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDH
        case ALG_ECDH_VALUE:
            return TPMS_KEY_SCHEME_ECDH_Unmarshal((TPMS_KEY_SCHEME_ECDH *)&(target->ecdh), buffer, size);
#endif // ALG_ECDH
#if ALG_ECMQV
        case ALG_ECMQV_VALUE:
            return TPMS_KEY_SCHEME_ECMQV_Unmarshal((TPMS_KEY_SCHEME_ECMQV *)&(target->ecmqv), buffer, size);
#endif // ALG_ECMQV
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIG_SCHEME_ECDAA_Unmarshal((TPMS_SIG_SCHEME_ECDAA *)&(target->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIG_SCHEME_RSASSA_Unmarshal((TPMS_SIG_SCHEME_RSASSA *)&(target->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIG_SCHEME_RSAPSS_Unmarshal((TPMS_SIG_SCHEME_RSAPSS *)&(target->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIG_SCHEME_ECDSA_Unmarshal((TPMS_SIG_SCHEME_ECDSA *)&(target->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIG_SCHEME_SM2_Unmarshal((TPMS_SIG_SCHEME_SM2 *)&(target->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIG_SCHEME_ECSCHNORR_Unmarshal((TPMS_SIG_SCHEME_ECSCHNORR *)&(target->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_RSAES
        case ALG_RSAES_VALUE:
            return TPMS_ENC_SCHEME_RSAES_Unmarshal((TPMS_ENC_SCHEME_RSAES *)&(target->rsaes), buffer, size);
#endif // ALG_RSAES
#if ALG_OAEP
        case ALG_OAEP_VALUE:
            return TPMS_ENC_SCHEME_OAEP_Unmarshal((TPMS_ENC_SCHEME_OAEP *)&(target->oaep), buffer, size);
#endif // ALG_OAEP
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_ASYM_SCHEME_Marshal(TPMU_ASYM_SCHEME *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDH
        case ALG_ECDH_VALUE:
            return TPMS_KEY_SCHEME_ECDH_Marshal((TPMS_KEY_SCHEME_ECDH *)&(source->ecdh), buffer, size);
#endif // ALG_ECDH
#if ALG_ECMQV
        case ALG_ECMQV_VALUE:
            return TPMS_KEY_SCHEME_ECMQV_Marshal((TPMS_KEY_SCHEME_ECMQV *)&(source->ecmqv), buffer, size);
#endif // ALG_ECMQV
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIG_SCHEME_ECDAA_Marshal((TPMS_SIG_SCHEME_ECDAA *)&(source->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIG_SCHEME_RSASSA_Marshal((TPMS_SIG_SCHEME_RSASSA *)&(source->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIG_SCHEME_RSAPSS_Marshal((TPMS_SIG_SCHEME_RSAPSS *)&(source->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIG_SCHEME_ECDSA_Marshal((TPMS_SIG_SCHEME_ECDSA *)&(source->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIG_SCHEME_SM2_Marshal((TPMS_SIG_SCHEME_SM2 *)&(source->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIG_SCHEME_ECSCHNORR_Marshal((TPMS_SIG_SCHEME_ECSCHNORR *)&(source->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_RSAES
        case ALG_RSAES_VALUE:
            return TPMS_ENC_SCHEME_RSAES_Marshal((TPMS_ENC_SCHEME_RSAES *)&(source->rsaes), buffer, size);
#endif // ALG_RSAES
#if ALG_OAEP
        case ALG_OAEP_VALUE:
            return TPMS_ENC_SCHEME_OAEP_Marshal((TPMS_ENC_SCHEME_OAEP *)&(source->oaep), buffer, size);
#endif // ALG_OAEP
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:165 - Definition of TPMT_ASYM_SCHEME Structure
// Table 2:166 - Definition of TPMI_ALG_RSA_SCHEME Type
#if ALG_RSA
TPM_RC
TPMI_ALG_RSA_SCHEME_Unmarshal(TPMI_ALG_RSA_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_RSAES
            case ALG_RSAES_VALUE:
#endif // ALG_RSAES
#if ALG_OAEP
            case ALG_OAEP_VALUE:
#endif // ALG_OAEP
#if ALG_RSASSA
            case ALG_RSASSA_VALUE:
#endif // ALG_RSASSA
#if ALG_RSAPSS
            case ALG_RSAPSS_VALUE:
#endif // ALG_RSAPSS
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_RSA_SCHEME_Marshal(TPMI_ALG_RSA_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:167 - Definition of TPMT_RSA_SCHEME Structure
#if ALG_RSA
TPM_RC
TPMT_RSA_SCHEME_Unmarshal(TPMT_RSA_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_RSA_SCHEME_Unmarshal((TPMI_ALG_RSA_SCHEME *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_ASYM_SCHEME_Unmarshal((TPMU_ASYM_SCHEME *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_RSA_SCHEME_Marshal(TPMT_RSA_SCHEME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_RSA_SCHEME_Marshal((TPMI_ALG_RSA_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_ASYM_SCHEME_Marshal((TPMU_ASYM_SCHEME *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}
#endif // ALG_RSA

// Table 2:168 - Definition of TPMI_ALG_RSA_DECRYPT Type
#if ALG_RSA
TPM_RC
TPMI_ALG_RSA_DECRYPT_Unmarshal(TPMI_ALG_RSA_DECRYPT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_RSAES
            case ALG_RSAES_VALUE:
#endif // ALG_RSAES
#if ALG_OAEP
            case ALG_OAEP_VALUE:
#endif // ALG_OAEP
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_VALUE;
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_RSA_DECRYPT_Marshal(TPMI_ALG_RSA_DECRYPT *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:169 - Definition of TPMT_RSA_DECRYPT Structure
#if ALG_RSA
TPM_RC
TPMT_RSA_DECRYPT_Unmarshal(TPMT_RSA_DECRYPT *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_RSA_DECRYPT_Unmarshal((TPMI_ALG_RSA_DECRYPT *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_ASYM_SCHEME_Unmarshal((TPMU_ASYM_SCHEME *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_RSA_DECRYPT_Marshal(TPMT_RSA_DECRYPT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_RSA_DECRYPT_Marshal((TPMI_ALG_RSA_DECRYPT *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_ASYM_SCHEME_Marshal((TPMU_ASYM_SCHEME *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}
#endif // ALG_RSA

// Table 2:170 - Definition of TPM2B_PUBLIC_KEY_RSA Structure
#if ALG_RSA
TPM_RC
TPM2B_PUBLIC_KEY_RSA_Unmarshal(TPM2B_PUBLIC_KEY_RSA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_RSA_KEY_BYTES)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_PUBLIC_KEY_RSA_Marshal(TPM2B_PUBLIC_KEY_RSA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}
#endif // ALG_RSA

// Table 2:171 - Definition of TPMI_RSA_KEY_BITS Type
#if ALG_RSA
TPM_RC
TPMI_RSA_KEY_BITS_Unmarshal(TPMI_RSA_KEY_BITS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_KEY_BITS_Unmarshal((TPM_KEY_BITS *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if RSA_1024
            case 1024:
#endif // RSA_1024
#if RSA_2048
            case 2048:
#endif // RSA_2048
#if RSA_3072
            case 3072:
#endif // RSA_3072
#if RSA_4096
            case 4096:
#endif // RSA_4096
                break;
            default:
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RSA_KEY_BITS_Marshal(TPMI_RSA_KEY_BITS *source, BYTE **buffer, INT32 *size)
{
    return TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:172 - Definition of TPM2B_PRIVATE_KEY_RSA Structure
#if ALG_RSA
TPM_RC
TPM2B_PRIVATE_KEY_RSA_Unmarshal(TPM2B_PRIVATE_KEY_RSA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > RSA_PRIVATE_SIZE)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_PRIVATE_KEY_RSA_Marshal(TPM2B_PRIVATE_KEY_RSA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}
#endif // ALG_RSA

// Table 2:173 - Definition of TPM2B_ECC_PARAMETER Structure
TPM_RC
TPM2B_ECC_PARAMETER_Unmarshal(TPM2B_ECC_PARAMETER *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_ECC_KEY_BYTES)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_ECC_PARAMETER_Marshal(TPM2B_ECC_PARAMETER *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:174 - Definition of TPMS_ECC_POINT Structure
#if ALG_ECC
TPM_RC
TPMS_ECC_POINT_Unmarshal(TPMS_ECC_POINT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM2B_ECC_PARAMETER_Unmarshal((TPM2B_ECC_PARAMETER *)&(target->x), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_ECC_PARAMETER_Unmarshal((TPM2B_ECC_PARAMETER *)&(target->y), buffer, size);
    return result;
}
UINT16
TPMS_ECC_POINT_Marshal(TPMS_ECC_POINT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->x), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->y), buffer, size));
    return result;
}
#endif // ALG_ECC

// Table 2:175 - Definition of TPM2B_ECC_POINT Structure
#if ALG_ECC
TPM_RC
TPM2B_ECC_POINT_Unmarshal(TPM2B_ECC_POINT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->size), buffer, size); // =a
    if(result == TPM_RC_SUCCESS)
    {
        // if size is zero, then the required structure is missing
        if(target->size == 0)
            result = TPM_RC_SIZE;
        else
        {
            INT32   startSize = *size;
            result = TPMS_ECC_POINT_Unmarshal((TPMS_ECC_POINT *)&(target->point), buffer, size); // =b
            if(result == TPM_RC_SUCCESS)
            {
                if(target->size != (startSize - *size))
                    result = TPM_RC_SIZE;
            }
        }
    }
    return result;
}
UINT16
TPM2B_ECC_POINT_Marshal(TPM2B_ECC_POINT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    // Marshal a dummy value of the 2B size. This makes sure that 'buffer'
    // and 'size' are advanced as necessary (i.e., if they are present)
    result = UINT16_Marshal(&result, buffer, size);
    // Marshal the structure
    result = (UINT16)(result + TPMS_ECC_POINT_Marshal((TPMS_ECC_POINT *)&(source->point), buffer, size));
    // if a buffer was provided, go back and fill in the actual size
    if(buffer != NULL)
        UINT16_TO_BYTE_ARRAY((result - 2), (*buffer - result));
    return result;
}
#endif // ALG_ECC

// Table 2:176 - Definition of TPMI_ALG_ECC_SCHEME Type
#if ALG_ECC
TPM_RC
TPMI_ALG_ECC_SCHEME_Unmarshal(TPMI_ALG_ECC_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_ECDAA
            case ALG_ECDAA_VALUE:
#endif // ALG_ECDAA
#if ALG_ECDSA
            case ALG_ECDSA_VALUE:
#endif // ALG_ECDSA
#if ALG_SM2
            case ALG_SM2_VALUE:
#endif // ALG_SM2
#if ALG_ECSCHNORR
            case ALG_ECSCHNORR_VALUE:
#endif // ALG_ECSCHNORR
#if ALG_ECDH
            case ALG_ECDH_VALUE:
#endif // ALG_ECDH
#if ALG_ECMQV
            case ALG_ECMQV_VALUE:
#endif // ALG_ECMQV
                break;
            case ALG_NULL_VALUE:
                if(!flag)
                    result = TPM_RC_SCHEME;
                break;
            default:
                result = TPM_RC_SCHEME;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ECC_SCHEME_Marshal(TPMI_ALG_ECC_SCHEME *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:177 - Definition of TPMI_ECC_CURVE Type
#if ALG_ECC
TPM_RC
TPMI_ECC_CURVE_Unmarshal(TPMI_ECC_CURVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ECC_CURVE_Unmarshal((TPM_ECC_CURVE *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ECC_BN_P256
            case TPM_ECC_BN_P256:
#endif // ECC_BN_P256
#if ECC_BN_P638
            case TPM_ECC_BN_P638:
#endif // ECC_BN_P638
#if ECC_NIST_P192
            case TPM_ECC_NIST_P192:
#endif // ECC_NIST_P192
#if ECC_NIST_P224
            case TPM_ECC_NIST_P224:
#endif // ECC_NIST_P224
#if ECC_NIST_P256
            case TPM_ECC_NIST_P256:
#endif // ECC_NIST_P256
#if ECC_NIST_P384
            case TPM_ECC_NIST_P384:
#endif // ECC_NIST_P384
#if ECC_NIST_P521
            case TPM_ECC_NIST_P521:
#endif // ECC_NIST_P521
#if ECC_SM2_P256
            case TPM_ECC_SM2_P256:
#endif // ECC_SM2_P256
                break;
            default:
                result = TPM_RC_CURVE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ECC_CURVE_Marshal(TPMI_ECC_CURVE *source, BYTE **buffer, INT32 *size)
{
    return TPM_ECC_CURVE_Marshal((TPM_ECC_CURVE *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:178 - Definition of TPMT_ECC_SCHEME Structure
#if ALG_ECC
TPM_RC
TPMT_ECC_SCHEME_Unmarshal(TPMT_ECC_SCHEME *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_ECC_SCHEME_Unmarshal((TPMI_ALG_ECC_SCHEME *)&(target->scheme), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_ASYM_SCHEME_Unmarshal((TPMU_ASYM_SCHEME *)&(target->details), buffer, size, (UINT32)target->scheme);
    return result;
}
UINT16
TPMT_ECC_SCHEME_Marshal(TPMT_ECC_SCHEME *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_ECC_SCHEME_Marshal((TPMI_ALG_ECC_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMU_ASYM_SCHEME_Marshal((TPMU_ASYM_SCHEME *)&(source->details), buffer, size, (UINT32)source->scheme));
    return result;
}
#endif // ALG_ECC

// Table 2:179 - Definition of TPMS_ALGORITHM_DETAIL_ECC Structure
#if ALG_ECC
UINT16
TPMS_ALGORITHM_DETAIL_ECC_Marshal(TPMS_ALGORITHM_DETAIL_ECC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_ECC_CURVE_Marshal((TPM_ECC_CURVE *)&(source->curveID), buffer, size));
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->keySize), buffer, size));
    result = (UINT16)(result + TPMT_KDF_SCHEME_Marshal((TPMT_KDF_SCHEME *)&(source->kdf), buffer, size));
    result = (UINT16)(result + TPMT_ECC_SCHEME_Marshal((TPMT_ECC_SCHEME *)&(source->sign), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->p), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->a), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->b), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->gX), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->gY), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->n), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->h), buffer, size));
    return result;
}
#endif // ALG_ECC

// Table 2:180 - Definition of TPMS_SIGNATURE_RSA Structure
#if ALG_RSA
TPM_RC
TPMS_SIGNATURE_RSA_Unmarshal(TPMS_SIGNATURE_RSA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hash), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_PUBLIC_KEY_RSA_Unmarshal((TPM2B_PUBLIC_KEY_RSA *)&(target->sig), buffer, size);
    return result;
}
UINT16
TPMS_SIGNATURE_RSA_Marshal(TPMS_SIGNATURE_RSA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hash), buffer, size));
    result = (UINT16)(result + TPM2B_PUBLIC_KEY_RSA_Marshal((TPM2B_PUBLIC_KEY_RSA *)&(source->sig), buffer, size));
    return result;
}
#endif // ALG_RSA

// Table 2:181 - Definition of Types for Signature
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_RSASSA_Unmarshal(TPMS_SIGNATURE_RSASSA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_RSA_Unmarshal((TPMS_SIGNATURE_RSA *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_RSASSA_Marshal(TPMS_SIGNATURE_RSASSA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_RSA_Marshal((TPMS_SIGNATURE_RSA *)source, buffer, size);
}
TPM_RC
TPMS_SIGNATURE_RSAPSS_Unmarshal(TPMS_SIGNATURE_RSAPSS *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_RSA_Unmarshal((TPMS_SIGNATURE_RSA *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_RSAPSS_Marshal(TPMS_SIGNATURE_RSAPSS *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_RSA_Marshal((TPMS_SIGNATURE_RSA *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:182 - Definition of TPMS_SIGNATURE_ECC Structure
#if ALG_ECC
TPM_RC
TPMS_SIGNATURE_ECC_Unmarshal(TPMS_SIGNATURE_ECC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->hash), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_ECC_PARAMETER_Unmarshal((TPM2B_ECC_PARAMETER *)&(target->signatureR), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_ECC_PARAMETER_Unmarshal((TPM2B_ECC_PARAMETER *)&(target->signatureS), buffer, size);
    return result;
}
UINT16
TPMS_SIGNATURE_ECC_Marshal(TPMS_SIGNATURE_ECC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->hash), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->signatureR), buffer, size));
    result = (UINT16)(result + TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->signatureS), buffer, size));
    return result;
}
#endif // ALG_ECC

// Table 2:183 - Definition of Types for TPMS_SIGNATURE_ECC
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_ECDAA_Unmarshal(TPMS_SIGNATURE_ECDAA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_ECDAA_Marshal(TPMS_SIGNATURE_ECDAA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)source, buffer, size);
}
TPM_RC
TPMS_SIGNATURE_ECDSA_Unmarshal(TPMS_SIGNATURE_ECDSA *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_ECDSA_Marshal(TPMS_SIGNATURE_ECDSA *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)source, buffer, size);
}
TPM_RC
TPMS_SIGNATURE_SM2_Unmarshal(TPMS_SIGNATURE_SM2 *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_SM2_Marshal(TPMS_SIGNATURE_SM2 *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)source, buffer, size);
}
TPM_RC
TPMS_SIGNATURE_ECSCHNORR_Unmarshal(TPMS_SIGNATURE_ECSCHNORR *target, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)target, buffer, size);
}
UINT16
TPMS_SIGNATURE_ECSCHNORR_Marshal(TPMS_SIGNATURE_ECSCHNORR *source, BYTE **buffer, INT32 *size)
{
    return TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:184 - Definition of TPMU_SIGNATURE Union
TPM_RC
TPMU_SIGNATURE_Unmarshal(TPMU_SIGNATURE *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIGNATURE_ECDAA_Unmarshal((TPMS_SIGNATURE_ECDAA *)&(target->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIGNATURE_RSASSA_Unmarshal((TPMS_SIGNATURE_RSASSA *)&(target->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIGNATURE_RSAPSS_Unmarshal((TPMS_SIGNATURE_RSAPSS *)&(target->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIGNATURE_ECDSA_Unmarshal((TPMS_SIGNATURE_ECDSA *)&(target->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIGNATURE_SM2_Unmarshal((TPMS_SIGNATURE_SM2 *)&(target->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIGNATURE_ECSCHNORR_Unmarshal((TPMS_SIGNATURE_ECSCHNORR *)&(target->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMT_HA_Unmarshal((TPMT_HA *)&(target->hmac), buffer, size, 0);
#endif // ALG_HMAC
        case ALG_NULL_VALUE:
            return TPM_RC_SUCCESS;
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SIGNATURE_Marshal(TPMU_SIGNATURE *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECDAA
        case ALG_ECDAA_VALUE:
            return TPMS_SIGNATURE_ECDAA_Marshal((TPMS_SIGNATURE_ECDAA *)&(source->ecdaa), buffer, size);
#endif // ALG_ECDAA
#if ALG_RSASSA
        case ALG_RSASSA_VALUE:
            return TPMS_SIGNATURE_RSASSA_Marshal((TPMS_SIGNATURE_RSASSA *)&(source->rsassa), buffer, size);
#endif // ALG_RSASSA
#if ALG_RSAPSS
        case ALG_RSAPSS_VALUE:
            return TPMS_SIGNATURE_RSAPSS_Marshal((TPMS_SIGNATURE_RSAPSS *)&(source->rsapss), buffer, size);
#endif // ALG_RSAPSS
#if ALG_ECDSA
        case ALG_ECDSA_VALUE:
            return TPMS_SIGNATURE_ECDSA_Marshal((TPMS_SIGNATURE_ECDSA *)&(source->ecdsa), buffer, size);
#endif // ALG_ECDSA
#if ALG_SM2
        case ALG_SM2_VALUE:
            return TPMS_SIGNATURE_SM2_Marshal((TPMS_SIGNATURE_SM2 *)&(source->sm2), buffer, size);
#endif // ALG_SM2
#if ALG_ECSCHNORR
        case ALG_ECSCHNORR_VALUE:
            return TPMS_SIGNATURE_ECSCHNORR_Marshal((TPMS_SIGNATURE_ECSCHNORR *)&(source->ecschnorr), buffer, size);
#endif // ALG_ECSCHNORR
#if ALG_HMAC
        case ALG_HMAC_VALUE:
            return TPMT_HA_Marshal((TPMT_HA *)&(source->hmac), buffer, size);
#endif // ALG_HMAC
        case ALG_NULL_VALUE:
            return 0;
    }
    return 0;
}

// Table 2:185 - Definition of TPMT_SIGNATURE Structure
TPM_RC
TPMT_SIGNATURE_Unmarshal(TPMT_SIGNATURE *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_SIG_SCHEME_Unmarshal((TPMI_ALG_SIG_SCHEME *)&(target->sigAlg), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SIGNATURE_Unmarshal((TPMU_SIGNATURE *)&(target->signature), buffer, size, (UINT32)target->sigAlg);
    return result;
}
UINT16
TPMT_SIGNATURE_Marshal(TPMT_SIGNATURE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_SIG_SCHEME_Marshal((TPMI_ALG_SIG_SCHEME *)&(source->sigAlg), buffer, size));
    result = (UINT16)(result + TPMU_SIGNATURE_Marshal((TPMU_SIGNATURE *)&(source->signature), buffer, size, (UINT32)source->sigAlg));
    return result;
}

// Table 2:186 - Definition of TPMU_ENCRYPTED_SECRET Union
TPM_RC
TPMU_ENCRYPTED_SECRET_Unmarshal(TPMU_ENCRYPTED_SECRET *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECC
        case ALG_ECC_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->ecc), buffer, size, (INT32)sizeof(TPMS_ECC_POINT));
#endif // ALG_ECC
#if ALG_RSA
        case ALG_RSA_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->rsa), buffer, size, (INT32)MAX_RSA_KEY_BYTES);
#endif // ALG_RSA
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->symmetric), buffer, size, (INT32)sizeof(TPM2B_DIGEST));
#endif // ALG_SYMCIPHER
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return BYTE_Array_Unmarshal((BYTE *)(target->keyedHash), buffer, size, (INT32)sizeof(TPM2B_DIGEST));
#endif // ALG_KEYEDHASH
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_ENCRYPTED_SECRET_Marshal(TPMU_ENCRYPTED_SECRET *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_ECC
        case ALG_ECC_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->ecc), buffer, size, (INT32)sizeof(TPMS_ECC_POINT));
#endif // ALG_ECC
#if ALG_RSA
        case ALG_RSA_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->rsa), buffer, size, (INT32)MAX_RSA_KEY_BYTES);
#endif // ALG_RSA
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->symmetric), buffer, size, (INT32)sizeof(TPM2B_DIGEST));
#endif // ALG_SYMCIPHER
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return BYTE_Array_Marshal((BYTE *)(source->keyedHash), buffer, size, (INT32)sizeof(TPM2B_DIGEST));
#endif // ALG_KEYEDHASH
    }
    return 0;
}

// Table 2:187 - Definition of TPM2B_ENCRYPTED_SECRET Structure
TPM_RC
TPM2B_ENCRYPTED_SECRET_Unmarshal(TPM2B_ENCRYPTED_SECRET *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMU_ENCRYPTED_SECRET))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.secret), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_ENCRYPTED_SECRET_Marshal(TPM2B_ENCRYPTED_SECRET *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.secret), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:188 - Definition of TPMI_ALG_PUBLIC Type
TPM_RC
TPMI_ALG_PUBLIC_Unmarshal(TPMI_ALG_PUBLIC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM_ALG_ID_Unmarshal((TPM_ALG_ID *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch (*target)
        {
#if ALG_RSA
            case ALG_RSA_VALUE:
#endif // ALG_RSA
#if ALG_ECC
            case ALG_ECC_VALUE:
#endif // ALG_ECC
#if ALG_KEYEDHASH
            case ALG_KEYEDHASH_VALUE:
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
            case ALG_SYMCIPHER_VALUE:
#endif // ALG_SYMCIPHER
                break;
            default:
                result = TPM_RC_TYPE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_PUBLIC_Marshal(TPMI_ALG_PUBLIC *source, BYTE **buffer, INT32 *size)
{
    return TPM_ALG_ID_Marshal((TPM_ALG_ID *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:189 - Definition of TPMU_PUBLIC_ID Union
TPM_RC
TPMU_PUBLIC_ID_Unmarshal(TPMU_PUBLIC_ID *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->keyedHash), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->sym), buffer, size);
#endif // ALG_SYMCIPHER
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPM2B_PUBLIC_KEY_RSA_Unmarshal((TPM2B_PUBLIC_KEY_RSA *)&(target->rsa), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPMS_ECC_POINT_Unmarshal((TPMS_ECC_POINT *)&(target->ecc), buffer, size);
#endif // ALG_ECC
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_PUBLIC_ID_Marshal(TPMU_PUBLIC_ID *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->keyedHash), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->sym), buffer, size);
#endif // ALG_SYMCIPHER
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPM2B_PUBLIC_KEY_RSA_Marshal((TPM2B_PUBLIC_KEY_RSA *)&(source->rsa), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPMS_ECC_POINT_Marshal((TPMS_ECC_POINT *)&(source->ecc), buffer, size);
#endif // ALG_ECC
    }
    return 0;
}

// Table 2:190 - Definition of TPMS_KEYEDHASH_PARMS Structure
TPM_RC
TPMS_KEYEDHASH_PARMS_Unmarshal(TPMS_KEYEDHASH_PARMS *target, BYTE **buffer, INT32 *size)
{
    return TPMT_KEYEDHASH_SCHEME_Unmarshal((TPMT_KEYEDHASH_SCHEME *)&(target->scheme), buffer, size, 1);
}
UINT16
TPMS_KEYEDHASH_PARMS_Marshal(TPMS_KEYEDHASH_PARMS *source, BYTE **buffer, INT32 *size)
{
    return TPMT_KEYEDHASH_SCHEME_Marshal((TPMT_KEYEDHASH_SCHEME *)&(source->scheme), buffer, size);
}

// Table 2:191 - Definition of TPMS_ASYM_PARMS Structure
// Table 2:192 - Definition of TPMS_RSA_PARMS Structure
#if ALG_RSA
TPM_RC
TPMS_RSA_PARMS_Unmarshal(TPMS_RSA_PARMS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMT_SYM_DEF_OBJECT_Unmarshal((TPMT_SYM_DEF_OBJECT *)&(target->symmetric), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPMT_RSA_SCHEME_Unmarshal((TPMT_RSA_SCHEME *)&(target->scheme), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RSA_KEY_BITS_Unmarshal((TPMI_RSA_KEY_BITS *)&(target->keyBits), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = UINT32_Unmarshal((UINT32 *)&(target->exponent), buffer, size);
    return result;
}
UINT16
TPMS_RSA_PARMS_Marshal(TPMS_RSA_PARMS *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMT_SYM_DEF_OBJECT_Marshal((TPMT_SYM_DEF_OBJECT *)&(source->symmetric), buffer, size));
    result = (UINT16)(result + TPMT_RSA_SCHEME_Marshal((TPMT_RSA_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMI_RSA_KEY_BITS_Marshal((TPMI_RSA_KEY_BITS *)&(source->keyBits), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->exponent), buffer, size));
    return result;
}
#endif // ALG_RSA

// Table 2:193 - Definition of TPMS_ECC_PARMS Structure
#if ALG_ECC
TPM_RC
TPMS_ECC_PARMS_Unmarshal(TPMS_ECC_PARMS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMT_SYM_DEF_OBJECT_Unmarshal((TPMT_SYM_DEF_OBJECT *)&(target->symmetric), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPMT_ECC_SCHEME_Unmarshal((TPMT_ECC_SCHEME *)&(target->scheme), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_ECC_CURVE_Unmarshal((TPMI_ECC_CURVE *)&(target->curveID), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMT_KDF_SCHEME_Unmarshal((TPMT_KDF_SCHEME *)&(target->kdf), buffer, size, 1);
    return result;
}
UINT16
TPMS_ECC_PARMS_Marshal(TPMS_ECC_PARMS *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMT_SYM_DEF_OBJECT_Marshal((TPMT_SYM_DEF_OBJECT *)&(source->symmetric), buffer, size));
    result = (UINT16)(result + TPMT_ECC_SCHEME_Marshal((TPMT_ECC_SCHEME *)&(source->scheme), buffer, size));
    result = (UINT16)(result + TPMI_ECC_CURVE_Marshal((TPMI_ECC_CURVE *)&(source->curveID), buffer, size));
    result = (UINT16)(result + TPMT_KDF_SCHEME_Marshal((TPMT_KDF_SCHEME *)&(source->kdf), buffer, size));
    return result;
}
#endif // ALG_ECC

// Table 2:194 - Definition of TPMU_PUBLIC_PARMS Union
TPM_RC
TPMU_PUBLIC_PARMS_Unmarshal(TPMU_PUBLIC_PARMS *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPMS_KEYEDHASH_PARMS_Unmarshal((TPMS_KEYEDHASH_PARMS *)&(target->keyedHashDetail), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPMS_SYMCIPHER_PARMS_Unmarshal((TPMS_SYMCIPHER_PARMS *)&(target->symDetail), buffer, size);
#endif // ALG_SYMCIPHER
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPMS_RSA_PARMS_Unmarshal((TPMS_RSA_PARMS *)&(target->rsaDetail), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPMS_ECC_PARMS_Unmarshal((TPMS_ECC_PARMS *)&(target->eccDetail), buffer, size);
#endif // ALG_ECC
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_PUBLIC_PARMS_Marshal(TPMU_PUBLIC_PARMS *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPMS_KEYEDHASH_PARMS_Marshal((TPMS_KEYEDHASH_PARMS *)&(source->keyedHashDetail), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPMS_SYMCIPHER_PARMS_Marshal((TPMS_SYMCIPHER_PARMS *)&(source->symDetail), buffer, size);
#endif // ALG_SYMCIPHER
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPMS_RSA_PARMS_Marshal((TPMS_RSA_PARMS *)&(source->rsaDetail), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPMS_ECC_PARMS_Marshal((TPMS_ECC_PARMS *)&(source->eccDetail), buffer, size);
#endif // ALG_ECC
    }
    return 0;
}

// Table 2:195 - Definition of TPMT_PUBLIC_PARMS Structure
TPM_RC
TPMT_PUBLIC_PARMS_Unmarshal(TPMT_PUBLIC_PARMS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_PUBLIC_Unmarshal((TPMI_ALG_PUBLIC *)&(target->type), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_PUBLIC_PARMS_Unmarshal((TPMU_PUBLIC_PARMS *)&(target->parameters), buffer, size, (UINT32)target->type);
    return result;
}
UINT16
TPMT_PUBLIC_PARMS_Marshal(TPMT_PUBLIC_PARMS *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_PUBLIC_Marshal((TPMI_ALG_PUBLIC *)&(source->type), buffer, size));
    result = (UINT16)(result + TPMU_PUBLIC_PARMS_Marshal((TPMU_PUBLIC_PARMS *)&(source->parameters), buffer, size, (UINT32)source->type));
    return result;
}

// Table 2:196 - Definition of TPMT_PUBLIC Structure
TPM_RC
TPMT_PUBLIC_Unmarshal(TPMT_PUBLIC *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = TPMI_ALG_PUBLIC_Unmarshal((TPMI_ALG_PUBLIC *)&(target->type), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->nameAlg), buffer, size, flag);
    if(result == TPM_RC_SUCCESS)
        result = TPMA_OBJECT_Unmarshal((TPMA_OBJECT *)&(target->objectAttributes), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->authPolicy), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_PUBLIC_PARMS_Unmarshal((TPMU_PUBLIC_PARMS *)&(target->parameters), buffer, size, (UINT32)target->type);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_PUBLIC_ID_Unmarshal((TPMU_PUBLIC_ID *)&(target->unique), buffer, size, (UINT32)target->type);
    return result;
}
UINT16
TPMT_PUBLIC_Marshal(TPMT_PUBLIC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_PUBLIC_Marshal((TPMI_ALG_PUBLIC *)&(source->type), buffer, size));
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->nameAlg), buffer, size));
    result = (UINT16)(result + TPMA_OBJECT_Marshal((TPMA_OBJECT *)&(source->objectAttributes), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->authPolicy), buffer, size));
    result = (UINT16)(result + TPMU_PUBLIC_PARMS_Marshal((TPMU_PUBLIC_PARMS *)&(source->parameters), buffer, size, (UINT32)source->type));
    result = (UINT16)(result + TPMU_PUBLIC_ID_Marshal((TPMU_PUBLIC_ID *)&(source->unique), buffer, size, (UINT32)source->type));
    return result;
}

// Table 2:197 - Definition of TPM2B_PUBLIC Structure
TPM_RC
TPM2B_PUBLIC_Unmarshal(TPM2B_PUBLIC *target, BYTE **buffer, INT32 *size, BOOL flag)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->size), buffer, size); // =a
    if(result == TPM_RC_SUCCESS)
    {
        // if size is zero, then the required structure is missing
        if(target->size == 0)
            result = TPM_RC_SIZE;
        else
        {
            INT32   startSize = *size;
            result = TPMT_PUBLIC_Unmarshal((TPMT_PUBLIC *)&(target->publicArea), buffer, size, flag); // =b
            if(result == TPM_RC_SUCCESS)
            {
                if(target->size != (startSize - *size))
                    result = TPM_RC_SIZE;
            }
        }
    }
    return result;
}
UINT16
TPM2B_PUBLIC_Marshal(TPM2B_PUBLIC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    // Marshal a dummy value of the 2B size. This makes sure that 'buffer'
    // and 'size' are advanced as necessary (i.e., if they are present)
    result = UINT16_Marshal(&result, buffer, size);
    // Marshal the structure
    result = (UINT16)(result + TPMT_PUBLIC_Marshal((TPMT_PUBLIC *)&(source->publicArea), buffer, size));
    // if a buffer was provided, go back and fill in the actual size
    if(buffer != NULL)
        UINT16_TO_BYTE_ARRAY((result - 2), (*buffer - result));
    return result;
}

// Table 2:198 - Definition of TPM2B_TEMPLATE Structure
TPM_RC
TPM2B_TEMPLATE_Unmarshal(TPM2B_TEMPLATE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMT_PUBLIC))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_TEMPLATE_Marshal(TPM2B_TEMPLATE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:199 - Definition of TPM2B_PRIVATE_VENDOR_SPECIFIC Structure
TPM_RC
TPM2B_PRIVATE_VENDOR_SPECIFIC_Unmarshal(TPM2B_PRIVATE_VENDOR_SPECIFIC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > PRIVATE_VENDOR_SPECIFIC_BYTES)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_PRIVATE_VENDOR_SPECIFIC_Marshal(TPM2B_PRIVATE_VENDOR_SPECIFIC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:200 - Definition of TPMU_SENSITIVE_COMPOSITE Union
TPM_RC
TPMU_SENSITIVE_COMPOSITE_Unmarshal(TPMU_SENSITIVE_COMPOSITE *target, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPM2B_PRIVATE_KEY_RSA_Unmarshal((TPM2B_PRIVATE_KEY_RSA *)&(target->rsa), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPM2B_ECC_PARAMETER_Unmarshal((TPM2B_ECC_PARAMETER *)&(target->ecc), buffer, size);
#endif // ALG_ECC
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPM2B_SENSITIVE_DATA_Unmarshal((TPM2B_SENSITIVE_DATA *)&(target->bits), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPM2B_SYM_KEY_Unmarshal((TPM2B_SYM_KEY *)&(target->sym), buffer, size);
#endif // ALG_SYMCIPHER
    }
    return TPM_RC_SELECTOR;
}
UINT16
TPMU_SENSITIVE_COMPOSITE_Marshal(TPMU_SENSITIVE_COMPOSITE *source, BYTE **buffer, INT32 *size, UINT32 selector)
{
    switch(selector) {
#if ALG_RSA
        case ALG_RSA_VALUE:
            return TPM2B_PRIVATE_KEY_RSA_Marshal((TPM2B_PRIVATE_KEY_RSA *)&(source->rsa), buffer, size);
#endif // ALG_RSA
#if ALG_ECC
        case ALG_ECC_VALUE:
            return TPM2B_ECC_PARAMETER_Marshal((TPM2B_ECC_PARAMETER *)&(source->ecc), buffer, size);
#endif // ALG_ECC
#if ALG_KEYEDHASH
        case ALG_KEYEDHASH_VALUE:
            return TPM2B_SENSITIVE_DATA_Marshal((TPM2B_SENSITIVE_DATA *)&(source->bits), buffer, size);
#endif // ALG_KEYEDHASH
#if ALG_SYMCIPHER
        case ALG_SYMCIPHER_VALUE:
            return TPM2B_SYM_KEY_Marshal((TPM2B_SYM_KEY *)&(source->sym), buffer, size);
#endif // ALG_SYMCIPHER
    }
    return 0;
}

// Table 2:201 - Definition of TPMT_SENSITIVE Structure
TPM_RC
TPMT_SENSITIVE_Unmarshal(TPMT_SENSITIVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_ALG_PUBLIC_Unmarshal((TPMI_ALG_PUBLIC *)&(target->sensitiveType), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_AUTH_Unmarshal((TPM2B_AUTH *)&(target->authValue), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->seedValue), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMU_SENSITIVE_COMPOSITE_Unmarshal((TPMU_SENSITIVE_COMPOSITE *)&(target->sensitive), buffer, size, (UINT32)target->sensitiveType);
    return result;
}
UINT16
TPMT_SENSITIVE_Marshal(TPMT_SENSITIVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_ALG_PUBLIC_Marshal((TPMI_ALG_PUBLIC *)&(source->sensitiveType), buffer, size));
    result = (UINT16)(result + TPM2B_AUTH_Marshal((TPM2B_AUTH *)&(source->authValue), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->seedValue), buffer, size));
    result = (UINT16)(result + TPMU_SENSITIVE_COMPOSITE_Marshal((TPMU_SENSITIVE_COMPOSITE *)&(source->sensitive), buffer, size, (UINT32)source->sensitiveType));
    return result;
}

// Table 2:202 - Definition of TPM2B_SENSITIVE Structure
TPM_RC
TPM2B_SENSITIVE_Unmarshal(TPM2B_SENSITIVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->size), buffer, size); // =a
    // if there was an error or if target->size equal to 0,
    // skip unmarshaling of the structure
    if((result == TPM_RC_SUCCESS) && (target->size != 0))
    {
        INT32   startSize = *size;
        result = TPMT_SENSITIVE_Unmarshal((TPMT_SENSITIVE *)&(target->sensitiveArea), buffer, size); // =b
        if(result == TPM_RC_SUCCESS)
        {
            if(target->size != (startSize - *size))
                result = TPM_RC_SIZE;
        }
    }
    return result;
}
UINT16
TPM2B_SENSITIVE_Marshal(TPM2B_SENSITIVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    // Marshal a dummy value of the 2B size. This makes sure that 'buffer'
    // and 'size' are advanced as necessary (i.e., if they are present)
    result = UINT16_Marshal(&result, buffer, size);
    // Marshal the structure
    result = (UINT16)(result + TPMT_SENSITIVE_Marshal((TPMT_SENSITIVE *)&(source->sensitiveArea), buffer, size));
    // if a buffer was provided, go back and fill in the actual size
    if(buffer != NULL)
        UINT16_TO_BYTE_ARRAY((result - 2), (*buffer - result));
    return result;
}

// Table 2:203 - Definition of _PRIVATE Structure
// Table 2:204 - Definition of TPM2B_PRIVATE Structure
TPM_RC
TPM2B_PRIVATE_Unmarshal(TPM2B_PRIVATE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(_PRIVATE))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_PRIVATE_Marshal(TPM2B_PRIVATE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:205 - Definition of TPMS_ID_OBJECT Structure
// Table 2:206 - Definition of TPM2B_ID_OBJECT Structure
TPM_RC
TPM2B_ID_OBJECT_Unmarshal(TPM2B_ID_OBJECT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMS_ID_OBJECT))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.credential), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_ID_OBJECT_Marshal(TPM2B_ID_OBJECT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.credential), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:207 - Definition of TPM_NV_INDEX Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPM_NV_INDEX_Marshal(TPM_NV_INDEX *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:208 - Definition of TPM_NT Constants
// Table 2:209 - Definition of TPMS_NV_PIN_COUNTER_PARAMETERS Structure
TPM_RC
TPMS_NV_PIN_COUNTER_PARAMETERS_Unmarshal(TPMS_NV_PIN_COUNTER_PARAMETERS *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)&(target->pinCount), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = UINT32_Unmarshal((UINT32 *)&(target->pinLimit), buffer, size);
    return result;
}
UINT16
TPMS_NV_PIN_COUNTER_PARAMETERS_Marshal(TPMS_NV_PIN_COUNTER_PARAMETERS *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->pinCount), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->pinLimit), buffer, size));
    return result;
}

// Table 2:210 - Definition of TPMA_NV Bits
TPM_RC
TPMA_NV_Unmarshal(TPMA_NV *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if(*((UINT32 *)target) & (UINT32)0x01f00300)
            result = TPM_RC_RESERVED_BITS;
    }
    return result;
}

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_NV_Marshal(TPMA_NV *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:211 - Definition of TPMS_NV_PUBLIC Structure
TPM_RC
TPMS_NV_PUBLIC_Unmarshal(TPMS_NV_PUBLIC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPMI_RH_NV_INDEX_Unmarshal((TPMI_RH_NV_INDEX *)&(target->nvIndex), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_ALG_HASH_Unmarshal((TPMI_ALG_HASH *)&(target->nameAlg), buffer, size, 0);
    if(result == TPM_RC_SUCCESS)
        result = TPMA_NV_Unmarshal((TPMA_NV *)&(target->attributes), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->authPolicy), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = UINT16_Unmarshal((UINT16 *)&(target->dataSize), buffer, size);
    if(  (result == TPM_RC_SUCCESS)
      && (target->dataSize > MAX_NV_INDEX_SIZE))
        result = TPM_RC_SIZE;
    return result;
}
UINT16
TPMS_NV_PUBLIC_Marshal(TPMS_NV_PUBLIC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPMI_RH_NV_INDEX_Marshal((TPMI_RH_NV_INDEX *)&(source->nvIndex), buffer, size));
    result = (UINT16)(result + TPMI_ALG_HASH_Marshal((TPMI_ALG_HASH *)&(source->nameAlg), buffer, size));
    result = (UINT16)(result + TPMA_NV_Marshal((TPMA_NV *)&(source->attributes), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->authPolicy), buffer, size));
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->dataSize), buffer, size));
    return result;
}

// Table 2:212 - Definition of TPM2B_NV_PUBLIC Structure
TPM_RC
TPM2B_NV_PUBLIC_Unmarshal(TPM2B_NV_PUBLIC *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->size), buffer, size); // =a
    if(result == TPM_RC_SUCCESS)
    {
        // if size is zero, then the required structure is missing
        if(target->size == 0)
            result = TPM_RC_SIZE;
        else
        {
            INT32   startSize = *size;
            result = TPMS_NV_PUBLIC_Unmarshal((TPMS_NV_PUBLIC *)&(target->nvPublic), buffer, size); // =b
            if(result == TPM_RC_SUCCESS)
            {
                if(target->size != (startSize - *size))
                    result = TPM_RC_SIZE;
            }
        }
    }
    return result;
}
UINT16
TPM2B_NV_PUBLIC_Marshal(TPM2B_NV_PUBLIC *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    // Marshal a dummy value of the 2B size. This makes sure that 'buffer'
    // and 'size' are advanced as necessary (i.e., if they are present)
    result = UINT16_Marshal(&result, buffer, size);
    // Marshal the structure
    result = (UINT16)(result + TPMS_NV_PUBLIC_Marshal((TPMS_NV_PUBLIC *)&(source->nvPublic), buffer, size));
    // if a buffer was provided, go back and fill in the actual size
    if(buffer != NULL)
        UINT16_TO_BYTE_ARRAY((result - 2), (*buffer - result));
    return result;
}

// Table 2:213 - Definition of TPM2B_CONTEXT_SENSITIVE Structure
TPM_RC
TPM2B_CONTEXT_SENSITIVE_Unmarshal(TPM2B_CONTEXT_SENSITIVE *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > MAX_CONTEXT_SIZE)
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_CONTEXT_SENSITIVE_Marshal(TPM2B_CONTEXT_SENSITIVE *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:214 - Definition of TPMS_CONTEXT_DATA Structure
TPM_RC
TPMS_CONTEXT_DATA_Unmarshal(TPMS_CONTEXT_DATA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)&(target->integrity), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_CONTEXT_SENSITIVE_Unmarshal((TPM2B_CONTEXT_SENSITIVE *)&(target->encrypted), buffer, size);
    return result;
}
UINT16
TPMS_CONTEXT_DATA_Marshal(TPMS_CONTEXT_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->integrity), buffer, size));
    result = (UINT16)(result + TPM2B_CONTEXT_SENSITIVE_Marshal((TPM2B_CONTEXT_SENSITIVE *)&(source->encrypted), buffer, size));
    return result;
}

// Table 2:215 - Definition of TPM2B_CONTEXT_DATA Structure
TPM_RC
TPM2B_CONTEXT_DATA_Unmarshal(TPM2B_CONTEXT_DATA *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT16_Unmarshal((UINT16 *)&(target->t.size), buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        if((target->t.size) > sizeof(TPMS_CONTEXT_DATA))
            result = TPM_RC_SIZE;
        else
            result = BYTE_Array_Unmarshal((BYTE *)(target->t.buffer), buffer, size, (INT32)(target->t.size));
    }
    return result;
}
UINT16
TPM2B_CONTEXT_DATA_Marshal(TPM2B_CONTEXT_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT16_Marshal((UINT16 *)&(source->t.size), buffer, size));
    // if size equal to 0, the rest of the structure is a zero buffer.  Stop processing
    if(source->t.size == 0)
        return result;
    result = (UINT16)(result + BYTE_Array_Marshal((BYTE *)(source->t.buffer), buffer, size, (INT32)(source->t.size)));
    return result;
}

// Table 2:216 - Definition of TPMS_CONTEXT Structure
TPM_RC
TPMS_CONTEXT_Unmarshal(TPMS_CONTEXT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT64_Unmarshal((UINT64 *)&(target->sequence), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_DH_SAVED_Unmarshal((TPMI_DH_SAVED *)&(target->savedHandle), buffer, size);
    if(result == TPM_RC_SUCCESS)
        result = TPMI_RH_HIERARCHY_Unmarshal((TPMI_RH_HIERARCHY *)&(target->hierarchy), buffer, size, 1);
    if(result == TPM_RC_SUCCESS)
        result = TPM2B_CONTEXT_DATA_Unmarshal((TPM2B_CONTEXT_DATA *)&(target->contextBlob), buffer, size);
    return result;
}
UINT16
TPMS_CONTEXT_Marshal(TPMS_CONTEXT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT64_Marshal((UINT64 *)&(source->sequence), buffer, size));
    result = (UINT16)(result + TPMI_DH_SAVED_Marshal((TPMI_DH_SAVED *)&(source->savedHandle), buffer, size));
    result = (UINT16)(result + TPMI_RH_HIERARCHY_Marshal((TPMI_RH_HIERARCHY *)&(source->hierarchy), buffer, size));
    result = (UINT16)(result + TPM2B_CONTEXT_DATA_Marshal((TPM2B_CONTEXT_DATA *)&(source->contextBlob), buffer, size));
    return result;
}

// Table 2:218 - Definition of TPMS_CREATION_DATA Structure
UINT16
TPMS_CREATION_DATA_Marshal(TPMS_CREATION_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPML_PCR_SELECTION_Marshal((TPML_PCR_SELECTION *)&(source->pcrSelect), buffer, size));
    result = (UINT16)(result + TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)&(source->pcrDigest), buffer, size));
    result = (UINT16)(result + TPMA_LOCALITY_Marshal((TPMA_LOCALITY *)&(source->locality), buffer, size));
    result = (UINT16)(result + TPM_ALG_ID_Marshal((TPM_ALG_ID *)&(source->parentNameAlg), buffer, size));
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->parentName), buffer, size));
    result = (UINT16)(result + TPM2B_NAME_Marshal((TPM2B_NAME *)&(source->parentQualifiedName), buffer, size));
    result = (UINT16)(result + TPM2B_DATA_Marshal((TPM2B_DATA *)&(source->outsideInfo), buffer, size));
    return result;
}

// Table 2:219 - Definition of TPM2B_CREATION_DATA Structure
UINT16
TPM2B_CREATION_DATA_Marshal(TPM2B_CREATION_DATA *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    // Marshal a dummy value of the 2B size. This makes sure that 'buffer'
    // and 'size' are advanced as necessary (i.e., if they are present)
    result = UINT16_Marshal(&result, buffer, size);
    // Marshal the structure
    result = (UINT16)(result + TPMS_CREATION_DATA_Marshal((TPMS_CREATION_DATA *)&(source->creationData), buffer, size));
    // if a buffer was provided, go back and fill in the actual size
    if(buffer != NULL)
        UINT16_TO_BYTE_ARRAY((result - 2), (*buffer - result));
    return result;
}

// Table 2:220 - Definition of TPM_AT Constants
TPM_RC
TPM_AT_Unmarshal(TPM_AT *target, BYTE **buffer, INT32 *size)
{
    TPM_RC    result;
    result = UINT32_Unmarshal((UINT32 *)target, buffer, size);
    if(result == TPM_RC_SUCCESS)
    {
        switch(*target)
        {
            case TPM_AT_ANY :
            case TPM_AT_ERROR :
            case TPM_AT_PV1 :
            case TPM_AT_VEND :
                break;
            default :
                result = TPM_RC_VALUE;
                break;
        }
    }
    return result;
}
#if !USE_MARSHALING_DEFINES
UINT16
TPM_AT_Marshal(TPM_AT *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:221 - Definition of TPM_AE Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_AE_Marshal(TPM_AE *source, BYTE **buffer, INT32 *size)
{
    return UINT32_Marshal((UINT32 *)source, buffer, size);
}
#endif // !USE_MARSHALING_DEFINES

// Table 2:222 - Definition of TPMS_AC_OUTPUT Structure
UINT16
TPMS_AC_OUTPUT_Marshal(TPMS_AC_OUTPUT *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + TPM_AT_Marshal((TPM_AT *)&(source->tag), buffer, size));
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->data), buffer, size));
    return result;
}

// Table 2:223 - Definition of TPML_AC_CAPABILITIES Structure
UINT16
TPML_AC_CAPABILITIES_Marshal(TPML_AC_CAPABILITIES *source, BYTE **buffer, INT32 *size)
{
    UINT16    result = 0;
    result = (UINT16)(result + UINT32_Marshal((UINT32 *)&(source->count), buffer, size));
    result = (UINT16)(result + TPMS_AC_OUTPUT_Array_Marshal((TPMS_AC_OUTPUT *)(source->acCapabilities), buffer, size, (INT32)(source->count)));
    return result;
}

// Array Marshal/Unmarshal for BYTE
TPM_RC
BYTE_Array_Unmarshal(BYTE *target, BYTE **buffer, INT32 *size, INT32 count)
{
    if(*size < count)
        return TPM_RC_INSUFFICIENT;
    memcpy(target, *buffer, count);
    *size -= count;
    *buffer += count;
    return TPM_RC_SUCCESS;
}
UINT16
BYTE_Array_Marshal(BYTE *source, BYTE **buffer, INT32 *size, INT32 count)
{
    if (buffer != 0)
    {
        if ((size == 0) || ((*size -= count) >= 0))
        {
            memcpy(*buffer, source, count);
            *buffer += count;
        }
        pAssert(size == 0 || (*size >= 0));
    }
    pAssert(count < INT16_MAX);
    return ((UINT16)count);
}

// Array Marshal/Unmarshal for TPM2B_DIGEST
TPM_RC
TPM2B_DIGEST_Array_Unmarshal(TPM2B_DIGEST *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPM2B_DIGEST_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPM2B_DIGEST_Array_Marshal(TPM2B_DIGEST *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPM2B_DIGEST_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMA_CC
UINT16
TPMA_CC_Array_Marshal(TPMA_CC *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMA_CC_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMS_AC_OUTPUT
UINT16
TPMS_AC_OUTPUT_Array_Marshal(TPMS_AC_OUTPUT *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_AC_OUTPUT_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMS_ALG_PROPERTY
UINT16
TPMS_ALG_PROPERTY_Array_Marshal(TPMS_ALG_PROPERTY *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_ALG_PROPERTY_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal/Unmarshal for TPMS_PCR_SELECTION
TPM_RC
TPMS_PCR_SELECTION_Array_Unmarshal(TPMS_PCR_SELECTION *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPMS_PCR_SELECTION_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPMS_PCR_SELECTION_Array_Marshal(TPMS_PCR_SELECTION *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_PCR_SELECTION_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMS_TAGGED_PCR_SELECT
UINT16
TPMS_TAGGED_PCR_SELECT_Array_Marshal(TPMS_TAGGED_PCR_SELECT *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_TAGGED_PCR_SELECT_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMS_TAGGED_POLICY
UINT16
TPMS_TAGGED_POLICY_Array_Marshal(TPMS_TAGGED_POLICY *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_TAGGED_POLICY_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal for TPMS_TAGGED_PROPERTY
UINT16
TPMS_TAGGED_PROPERTY_Array_Marshal(TPMS_TAGGED_PROPERTY *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMS_TAGGED_PROPERTY_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal/Unmarshal for TPMT_HA
TPM_RC
TPMT_HA_Array_Unmarshal(TPMT_HA *target, BYTE **buffer, INT32 *size, BOOL flag, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPMT_HA_Unmarshal(&target[i], buffer, size, flag);
    }
    return result;
}
UINT16
TPMT_HA_Array_Marshal(TPMT_HA *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPMT_HA_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal/Unmarshal for TPM_ALG_ID
TPM_RC
TPM_ALG_ID_Array_Unmarshal(TPM_ALG_ID *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPM_ALG_ID_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPM_ALG_ID_Array_Marshal(TPM_ALG_ID *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPM_ALG_ID_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal/Unmarshal for TPM_CC
TPM_RC
TPM_CC_Array_Unmarshal(TPM_CC *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPM_CC_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPM_CC_Array_Marshal(TPM_CC *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPM_CC_Marshal(&source[i], buffer, size));
    }
    return result;
}

// Array Marshal/Unmarshal for TPM_ECC_CURVE
#if ALG_ECC
TPM_RC
TPM_ECC_CURVE_Array_Unmarshal(TPM_ECC_CURVE *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPM_ECC_CURVE_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPM_ECC_CURVE_Array_Marshal(TPM_ECC_CURVE *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPM_ECC_CURVE_Marshal(&source[i], buffer, size));
    }
    return result;
}
#endif // ALG_ECC

// Array Marshal/Unmarshal for TPM_HANDLE
TPM_RC
TPM_HANDLE_Array_Unmarshal(TPM_HANDLE *target, BYTE **buffer, INT32 *size, INT32 count)
{
    TPM_RC    result;
    INT32       i;
    for(result = TPM_RC_SUCCESS, i = 0;
        ((result == TPM_RC_SUCCESS) && (i < count)); i++)
    {
        result = TPM_HANDLE_Unmarshal(&target[i], buffer, size);
    }
    return result;
}
UINT16
TPM_HANDLE_Array_Marshal(TPM_HANDLE *source, BYTE **buffer, INT32 *size, INT32 count)
{
    UINT16    result = 0;
    INT32 i;
    for(i = 0; i < count; i++)
    {
        result = (UINT16)(result + TPM_HANDLE_Marshal(&source[i], buffer, size));
    }
    return result;
}

