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

#ifndef _MARSHAL_FP_H_
#define _MARSHAL_FP_H_

// Table 2:3 - Definition of Base Types
//   UINT8 definition from table 2:3
TPM_RC
UINT8_Unmarshal(UINT8 *target, BYTE **buffer, INT32 *size);
UINT16
UINT8_Marshal(UINT8 *source, BYTE **buffer, INT32 *size);

//   BYTE definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
BYTE_Unmarshal(BYTE *target, BYTE **buffer, INT32 *size);
#else
#define BYTE_Unmarshal(target, buffer, size)                                       \
            UINT8_Unmarshal((UINT8 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
BYTE_Marshal(BYTE *source, BYTE **buffer, INT32 *size);
#else
#define BYTE_Marshal(source, buffer, size)                                         \
            UINT8_Marshal((UINT8 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

//   INT8 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT8_Unmarshal(INT8 *target, BYTE **buffer, INT32 *size);
#else
#define INT8_Unmarshal(target, buffer, size)                                       \
            UINT8_Unmarshal((UINT8 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
INT8_Marshal(INT8 *source, BYTE **buffer, INT32 *size);
#else
#define INT8_Marshal(source, buffer, size)                                         \
            UINT8_Marshal((UINT8 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

//   UINT16 definition from table 2:3
TPM_RC
UINT16_Unmarshal(UINT16 *target, BYTE **buffer, INT32 *size);
UINT16
UINT16_Marshal(UINT16 *source, BYTE **buffer, INT32 *size);

//   INT16 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT16_Unmarshal(INT16 *target, BYTE **buffer, INT32 *size);
#else
#define INT16_Unmarshal(target, buffer, size)                                      \
            UINT16_Unmarshal((UINT16 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
INT16_Marshal(INT16 *source, BYTE **buffer, INT32 *size);
#else
#define INT16_Marshal(source, buffer, size)                                        \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

//   UINT32 definition from table 2:3
TPM_RC
UINT32_Unmarshal(UINT32 *target, BYTE **buffer, INT32 *size);
UINT16
UINT32_Marshal(UINT32 *source, BYTE **buffer, INT32 *size);

//   INT32 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT32_Unmarshal(INT32 *target, BYTE **buffer, INT32 *size);
#else
#define INT32_Unmarshal(target, buffer, size)                                      \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
INT32_Marshal(INT32 *source, BYTE **buffer, INT32 *size);
#else
#define INT32_Marshal(source, buffer, size)                                        \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

//   UINT64 definition from table 2:3
TPM_RC
UINT64_Unmarshal(UINT64 *target, BYTE **buffer, INT32 *size);
UINT16
UINT64_Marshal(UINT64 *source, BYTE **buffer, INT32 *size);

//   INT64 definition from table 2:3
#if !USE_MARSHALING_DEFINES
TPM_RC
INT64_Unmarshal(INT64 *target, BYTE **buffer, INT32 *size);
#else
#define INT64_Unmarshal(target, buffer, size)                                      \
            UINT64_Unmarshal((UINT64 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
INT64_Marshal(INT64 *source, BYTE **buffer, INT32 *size);
#else
#define INT64_Marshal(source, buffer, size)                                        \
            UINT64_Marshal((UINT64 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:4 - Defines for Logic Values
// Table 2:5 - Definition of Types for Documentation Clarity
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ALGORITHM_ID_Unmarshal(TPM_ALGORITHM_ID *target, BYTE **buffer, INT32 *size);
#else
#define TPM_ALGORITHM_ID_Unmarshal(target, buffer, size)                           \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_ALGORITHM_ID_Marshal(TPM_ALGORITHM_ID *source, BYTE **buffer, INT32 *size);
#else
#define TPM_ALGORITHM_ID_Marshal(source, buffer, size)                             \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_MODIFIER_INDICATOR_Unmarshal(TPM_MODIFIER_INDICATOR *target,
            BYTE **buffer, INT32 *size);
#else
#define TPM_MODIFIER_INDICATOR_Unmarshal(target, buffer, size)                     \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_MODIFIER_INDICATOR_Marshal(TPM_MODIFIER_INDICATOR *source,
            BYTE **buffer, INT32 *size);
#else
#define TPM_MODIFIER_INDICATOR_Marshal(source, buffer, size)                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_AUTHORIZATION_SIZE_Unmarshal(TPM_AUTHORIZATION_SIZE *target,
            BYTE **buffer, INT32 *size);
#else
#define TPM_AUTHORIZATION_SIZE_Unmarshal(target, buffer, size)                     \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_AUTHORIZATION_SIZE_Marshal(TPM_AUTHORIZATION_SIZE *source,
            BYTE **buffer, INT32 *size);
#else
#define TPM_AUTHORIZATION_SIZE_Marshal(source, buffer, size)                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_PARAMETER_SIZE_Unmarshal(TPM_PARAMETER_SIZE *target,
            BYTE **buffer, INT32 *size);
#else
#define TPM_PARAMETER_SIZE_Unmarshal(target, buffer, size)                         \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_PARAMETER_SIZE_Marshal(TPM_PARAMETER_SIZE *source, BYTE **buffer, INT32 *size);
#else
#define TPM_PARAMETER_SIZE_Marshal(source, buffer, size)                           \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_KEY_SIZE_Unmarshal(TPM_KEY_SIZE *target, BYTE **buffer, INT32 *size);
#else
#define TPM_KEY_SIZE_Unmarshal(target, buffer, size)                               \
            UINT16_Unmarshal((UINT16 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_KEY_SIZE_Marshal(TPM_KEY_SIZE *source, BYTE **buffer, INT32 *size);
#else
#define TPM_KEY_SIZE_Marshal(source, buffer, size)                                 \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_KEY_BITS_Unmarshal(TPM_KEY_BITS *target, BYTE **buffer, INT32 *size);
#else
#define TPM_KEY_BITS_Unmarshal(target, buffer, size)                               \
            UINT16_Unmarshal((UINT16 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_KEY_BITS_Marshal(TPM_KEY_BITS *source, BYTE **buffer, INT32 *size);
#else
#define TPM_KEY_BITS_Marshal(source, buffer, size)                                 \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:6 - Definition of TPM_SPEC Constants
// Table 2:7 - Definition of TPM_GENERATED Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_GENERATED_Marshal(TPM_GENERATED *source, BYTE **buffer, INT32 *size);
#else
#define TPM_GENERATED_Marshal(source, buffer, size)                                \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:9 - Definition of TPM_ALG_ID Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ALG_ID_Unmarshal(TPM_ALG_ID *target, BYTE **buffer, INT32 *size);
#else
#define TPM_ALG_ID_Unmarshal(target, buffer, size)                                 \
            UINT16_Unmarshal((UINT16 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_ALG_ID_Marshal(TPM_ALG_ID *source, BYTE **buffer, INT32 *size);
#else
#define TPM_ALG_ID_Marshal(source, buffer, size)                                   \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:10 - Definition of TPM_ECC_CURVE Constants
#if ALG_ECC
TPM_RC
TPM_ECC_CURVE_Unmarshal(TPM_ECC_CURVE *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPM_ECC_CURVE_Marshal(TPM_ECC_CURVE *source, BYTE **buffer, INT32 *size);
#else
#define TPM_ECC_CURVE_Marshal(source, buffer, size)                                \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:12 - Definition of TPM_CC Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_CC_Unmarshal(TPM_CC *target, BYTE **buffer, INT32 *size);
#else
#define TPM_CC_Unmarshal(target, buffer, size)                                     \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_CC_Marshal(TPM_CC *source, BYTE **buffer, INT32 *size);
#else
#define TPM_CC_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:16 - Definition of TPM_RC Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_RC_Marshal(TPM_RC *source, BYTE **buffer, INT32 *size);
#else
#define TPM_RC_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:17 - Definition of TPM_CLOCK_ADJUST Constants
TPM_RC
TPM_CLOCK_ADJUST_Unmarshal(TPM_CLOCK_ADJUST *target, BYTE **buffer, INT32 *size);

// Table 2:18 - Definition of TPM_EO Constants
TPM_RC
TPM_EO_Unmarshal(TPM_EO *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPM_EO_Marshal(TPM_EO *source, BYTE **buffer, INT32 *size);
#else
#define TPM_EO_Marshal(source, buffer, size)                                       \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:19 - Definition of TPM_ST Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_ST_Unmarshal(TPM_ST *target, BYTE **buffer, INT32 *size);
#else
#define TPM_ST_Unmarshal(target, buffer, size)                                     \
            UINT16_Unmarshal((UINT16 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_ST_Marshal(TPM_ST *source, BYTE **buffer, INT32 *size);
#else
#define TPM_ST_Marshal(source, buffer, size)                                       \
            UINT16_Marshal((UINT16 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:20 - Definition of TPM_SU Constants
TPM_RC
TPM_SU_Unmarshal(TPM_SU *target, BYTE **buffer, INT32 *size);

// Table 2:21 - Definition of TPM_SE Constants
TPM_RC
TPM_SE_Unmarshal(TPM_SE *target, BYTE **buffer, INT32 *size);

// Table 2:22 - Definition of TPM_CAP Constants
TPM_RC
TPM_CAP_Unmarshal(TPM_CAP *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPM_CAP_Marshal(TPM_CAP *source, BYTE **buffer, INT32 *size);
#else
#define TPM_CAP_Marshal(source, buffer, size)                                      \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:23 - Definition of TPM_PT Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_PT_Unmarshal(TPM_PT *target, BYTE **buffer, INT32 *size);
#else
#define TPM_PT_Unmarshal(target, buffer, size)                                     \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_PT_Marshal(TPM_PT *source, BYTE **buffer, INT32 *size);
#else
#define TPM_PT_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:24 - Definition of TPM_PT_PCR Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_PT_PCR_Unmarshal(TPM_PT_PCR *target, BYTE **buffer, INT32 *size);
#else
#define TPM_PT_PCR_Unmarshal(target, buffer, size)                                 \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_PT_PCR_Marshal(TPM_PT_PCR *source, BYTE **buffer, INT32 *size);
#else
#define TPM_PT_PCR_Marshal(source, buffer, size)                                   \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:25 - Definition of TPM_PS Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_PS_Marshal(TPM_PS *source, BYTE **buffer, INT32 *size);
#else
#define TPM_PS_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:26 - Definition of Types for Handles
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HANDLE_Unmarshal(TPM_HANDLE *target, BYTE **buffer, INT32 *size);
#else
#define TPM_HANDLE_Unmarshal(target, buffer, size)                                 \
            UINT32_Unmarshal((UINT32 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_HANDLE_Marshal(TPM_HANDLE *source, BYTE **buffer, INT32 *size);
#else
#define TPM_HANDLE_Marshal(source, buffer, size)                                   \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:27 - Definition of TPM_HT Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HT_Unmarshal(TPM_HT *target, BYTE **buffer, INT32 *size);
#else
#define TPM_HT_Unmarshal(target, buffer, size)                                     \
            UINT8_Unmarshal((UINT8 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_HT_Marshal(TPM_HT *source, BYTE **buffer, INT32 *size);
#else
#define TPM_HT_Marshal(source, buffer, size)                                       \
            UINT8_Marshal((UINT8 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:28 - Definition of TPM_RH Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_RH_Unmarshal(TPM_RH *target, BYTE **buffer, INT32 *size);
#else
#define TPM_RH_Unmarshal(target, buffer, size)                                     \
            TPM_HANDLE_Unmarshal((TPM_HANDLE *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_RH_Marshal(TPM_RH *source, BYTE **buffer, INT32 *size);
#else
#define TPM_RH_Marshal(source, buffer, size)                                       \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:29 - Definition of TPM_HC Constants
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM_HC_Unmarshal(TPM_HC *target, BYTE **buffer, INT32 *size);
#else
#define TPM_HC_Unmarshal(target, buffer, size)                                     \
            TPM_HANDLE_Unmarshal((TPM_HANDLE *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM_HC_Marshal(TPM_HC *source, BYTE **buffer, INT32 *size);
#else
#define TPM_HC_Marshal(source, buffer, size)                                       \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:30 - Definition of TPMA_ALGORITHM Bits
TPM_RC
TPMA_ALGORITHM_Unmarshal(TPMA_ALGORITHM *target, BYTE **buffer, INT32 *size);

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_ALGORITHM_Marshal(TPMA_ALGORITHM *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_ALGORITHM_Marshal(source, buffer, size)                               \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:31 - Definition of TPMA_OBJECT Bits
TPM_RC
TPMA_OBJECT_Unmarshal(TPMA_OBJECT *target, BYTE **buffer, INT32 *size);

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_OBJECT_Marshal(TPMA_OBJECT *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_OBJECT_Marshal(source, buffer, size)                                  \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:32 - Definition of TPMA_SESSION Bits
TPM_RC
TPMA_SESSION_Unmarshal(TPMA_SESSION *target, BYTE **buffer, INT32 *size);

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_SESSION_Marshal(TPMA_SESSION *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_SESSION_Marshal(source, buffer, size)                                 \
            UINT8_Marshal((UINT8 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:33 - Definition of TPMA_LOCALITY Bits
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMA_LOCALITY_Unmarshal(TPMA_LOCALITY *target, BYTE **buffer, INT32 *size);
#else
#define TPMA_LOCALITY_Unmarshal(target, buffer, size)                              \
            UINT8_Unmarshal((UINT8 *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_LOCALITY_Marshal(TPMA_LOCALITY *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_LOCALITY_Marshal(source, buffer, size)                                \
            UINT8_Marshal((UINT8 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:34 - Definition of TPMA_PERMANENT Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_PERMANENT_Marshal(TPMA_PERMANENT *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_PERMANENT_Marshal(source, buffer, size)                               \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:35 - Definition of TPMA_STARTUP_CLEAR Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_STARTUP_CLEAR_Marshal(TPMA_STARTUP_CLEAR *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_STARTUP_CLEAR_Marshal(source, buffer, size)                           \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:36 - Definition of TPMA_MEMORY Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_MEMORY_Marshal(TPMA_MEMORY *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_MEMORY_Marshal(source, buffer, size)                                  \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:37 - Definition of TPMA_CC Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_CC_Marshal(TPMA_CC *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_CC_Marshal(source, buffer, size)                                      \
            TPM_CC_Marshal((TPM_CC *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:38 - Definition of TPMA_MODES Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_MODES_Marshal(TPMA_MODES *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_MODES_Marshal(source, buffer, size)                                   \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:39 - Definition of TPMA_X509_KEY_USAGE Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPMA_X509_KEY_USAGE_Marshal(TPMA_X509_KEY_USAGE *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMA_X509_KEY_USAGE_Marshal(source, buffer, size)                          \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:40 - Definition of TPMI_YES_NO Type
TPM_RC
TPMI_YES_NO_Unmarshal(TPMI_YES_NO *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_YES_NO_Marshal(TPMI_YES_NO *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_YES_NO_Marshal(source, buffer, size)                                  \
            BYTE_Marshal((BYTE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:41 - Definition of TPMI_DH_OBJECT Type
TPM_RC
TPMI_DH_OBJECT_Unmarshal(TPMI_DH_OBJECT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_OBJECT_Marshal(TPMI_DH_OBJECT *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_DH_OBJECT_Marshal(source, buffer, size)                               \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:42 - Definition of TPMI_DH_PARENT Type
TPM_RC
TPMI_DH_PARENT_Unmarshal(TPMI_DH_PARENT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_PARENT_Marshal(TPMI_DH_PARENT *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_DH_PARENT_Marshal(source, buffer, size)                               \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:43 - Definition of TPMI_DH_PERSISTENT Type
TPM_RC
TPMI_DH_PERSISTENT_Unmarshal(TPMI_DH_PERSISTENT *target,
            BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_PERSISTENT_Marshal(TPMI_DH_PERSISTENT *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_DH_PERSISTENT_Marshal(source, buffer, size)                           \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:44 - Definition of TPMI_DH_ENTITY Type
TPM_RC
TPMI_DH_ENTITY_Unmarshal(TPMI_DH_ENTITY *target,
            BYTE **buffer, INT32 *size, BOOL flag);

// Table 2:45 - Definition of TPMI_DH_PCR Type
TPM_RC
TPMI_DH_PCR_Unmarshal(TPMI_DH_PCR *target, BYTE **buffer, INT32 *size, BOOL flag);

// Table 2:46 - Definition of TPMI_SH_AUTH_SESSION Type
TPM_RC
TPMI_SH_AUTH_SESSION_Unmarshal(TPMI_SH_AUTH_SESSION *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_AUTH_SESSION_Marshal(TPMI_SH_AUTH_SESSION *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_SH_AUTH_SESSION_Marshal(source, buffer, size)                         \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:47 - Definition of TPMI_SH_HMAC Type
TPM_RC
TPMI_SH_HMAC_Unmarshal(TPMI_SH_HMAC *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_HMAC_Marshal(TPMI_SH_HMAC *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_SH_HMAC_Marshal(source, buffer, size)                                 \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:48 - Definition of TPMI_SH_POLICY Type
TPM_RC
TPMI_SH_POLICY_Unmarshal(TPMI_SH_POLICY *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SH_POLICY_Marshal(TPMI_SH_POLICY *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_SH_POLICY_Marshal(source, buffer, size)                               \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:49 - Definition of TPMI_DH_CONTEXT Type
TPM_RC
TPMI_DH_CONTEXT_Unmarshal(TPMI_DH_CONTEXT *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_CONTEXT_Marshal(TPMI_DH_CONTEXT *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_DH_CONTEXT_Marshal(source, buffer, size)                              \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:50 - Definition of TPMI_DH_SAVED Type
TPM_RC
TPMI_DH_SAVED_Unmarshal(TPMI_DH_SAVED *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_DH_SAVED_Marshal(TPMI_DH_SAVED *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_DH_SAVED_Marshal(source, buffer, size)                                \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:51 - Definition of TPMI_RH_HIERARCHY Type
TPM_RC
TPMI_RH_HIERARCHY_Unmarshal(TPMI_RH_HIERARCHY *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_HIERARCHY_Marshal(TPMI_RH_HIERARCHY *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_RH_HIERARCHY_Marshal(source, buffer, size)                            \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:52 - Definition of TPMI_RH_ENABLES Type
TPM_RC
TPMI_RH_ENABLES_Unmarshal(TPMI_RH_ENABLES *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_ENABLES_Marshal(TPMI_RH_ENABLES *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_RH_ENABLES_Marshal(source, buffer, size)                              \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:53 - Definition of TPMI_RH_HIERARCHY_AUTH Type
TPM_RC
TPMI_RH_HIERARCHY_AUTH_Unmarshal(TPMI_RH_HIERARCHY_AUTH *target,
            BYTE **buffer, INT32 *size);

// Table 2:54 - Definition of TPMI_RH_PLATFORM Type
TPM_RC
TPMI_RH_PLATFORM_Unmarshal(TPMI_RH_PLATFORM *target, BYTE **buffer, INT32 *size);

// Table 2:55 - Definition of TPMI_RH_OWNER Type
TPM_RC
TPMI_RH_OWNER_Unmarshal(TPMI_RH_OWNER *target,
            BYTE **buffer, INT32 *size, BOOL flag);

// Table 2:56 - Definition of TPMI_RH_ENDORSEMENT Type
TPM_RC
TPMI_RH_ENDORSEMENT_Unmarshal(TPMI_RH_ENDORSEMENT *target,
            BYTE **buffer, INT32 *size, BOOL flag);

// Table 2:57 - Definition of TPMI_RH_PROVISION Type
TPM_RC
TPMI_RH_PROVISION_Unmarshal(TPMI_RH_PROVISION *target, BYTE **buffer, INT32 *size);

// Table 2:58 - Definition of TPMI_RH_CLEAR Type
TPM_RC
TPMI_RH_CLEAR_Unmarshal(TPMI_RH_CLEAR *target, BYTE **buffer, INT32 *size);

// Table 2:59 - Definition of TPMI_RH_NV_AUTH Type
TPM_RC
TPMI_RH_NV_AUTH_Unmarshal(TPMI_RH_NV_AUTH *target, BYTE **buffer, INT32 *size);

// Table 2:60 - Definition of TPMI_RH_LOCKOUT Type
TPM_RC
TPMI_RH_LOCKOUT_Unmarshal(TPMI_RH_LOCKOUT *target, BYTE **buffer, INT32 *size);

// Table 2:61 - Definition of TPMI_RH_NV_INDEX Type
TPM_RC
TPMI_RH_NV_INDEX_Unmarshal(TPMI_RH_NV_INDEX *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RH_NV_INDEX_Marshal(TPMI_RH_NV_INDEX *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_RH_NV_INDEX_Marshal(source, buffer, size)                             \
            TPM_HANDLE_Marshal((TPM_HANDLE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:62 - Definition of TPMI_RH_AC Type
TPM_RC
TPMI_RH_AC_Unmarshal(TPMI_RH_AC *target, BYTE **buffer, INT32 *size);

// Table 2:63 - Definition of TPMI_ALG_HASH Type
TPM_RC
TPMI_ALG_HASH_Unmarshal(TPMI_ALG_HASH *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_HASH_Marshal(TPMI_ALG_HASH *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_HASH_Marshal(source, buffer, size)                                \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:64 - Definition of TPMI_ALG_ASYM Type
TPM_RC
TPMI_ALG_ASYM_Unmarshal(TPMI_ALG_ASYM *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ASYM_Marshal(TPMI_ALG_ASYM *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_ASYM_Marshal(source, buffer, size)                                \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:65 - Definition of TPMI_ALG_SYM Type
TPM_RC
TPMI_ALG_SYM_Unmarshal(TPMI_ALG_SYM *target, BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_Marshal(TPMI_ALG_SYM *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_SYM_Marshal(source, buffer, size)                                 \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:66 - Definition of TPMI_ALG_SYM_OBJECT Type
TPM_RC
TPMI_ALG_SYM_OBJECT_Unmarshal(TPMI_ALG_SYM_OBJECT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_OBJECT_Marshal(TPMI_ALG_SYM_OBJECT *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_SYM_OBJECT_Marshal(source, buffer, size)                          \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:67 - Definition of TPMI_ALG_SYM_MODE Type
TPM_RC
TPMI_ALG_SYM_MODE_Unmarshal(TPMI_ALG_SYM_MODE *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SYM_MODE_Marshal(TPMI_ALG_SYM_MODE *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_SYM_MODE_Marshal(source, buffer, size)                            \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:68 - Definition of TPMI_ALG_KDF Type
TPM_RC
TPMI_ALG_KDF_Unmarshal(TPMI_ALG_KDF *target, BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_KDF_Marshal(TPMI_ALG_KDF *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_KDF_Marshal(source, buffer, size)                                 \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:69 - Definition of TPMI_ALG_SIG_SCHEME Type
TPM_RC
TPMI_ALG_SIG_SCHEME_Unmarshal(TPMI_ALG_SIG_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_SIG_SCHEME_Marshal(TPMI_ALG_SIG_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_SIG_SCHEME_Marshal(source, buffer, size)                          \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:70 - Definition of TPMI_ECC_KEY_EXCHANGE Type
#if ALG_ECC
TPM_RC
TPMI_ECC_KEY_EXCHANGE_Unmarshal(TPMI_ECC_KEY_EXCHANGE *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ECC_KEY_EXCHANGE_Marshal(TPMI_ECC_KEY_EXCHANGE *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ECC_KEY_EXCHANGE_Marshal(source, buffer, size)                        \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:71 - Definition of TPMI_ST_COMMAND_TAG Type
TPM_RC
TPMI_ST_COMMAND_TAG_Unmarshal(TPMI_ST_COMMAND_TAG *target,
            BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ST_COMMAND_TAG_Marshal(TPMI_ST_COMMAND_TAG *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ST_COMMAND_TAG_Marshal(source, buffer, size)                          \
            TPM_ST_Marshal((TPM_ST *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:72 - Definition of TPMI_ALG_MAC_SCHEME Type
TPM_RC
TPMI_ALG_MAC_SCHEME_Unmarshal(TPMI_ALG_MAC_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_MAC_SCHEME_Marshal(TPMI_ALG_MAC_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_MAC_SCHEME_Marshal(source, buffer, size)                          \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:73 - Definition of TPMI_ALG_CIPHER_MODE Type
TPM_RC
TPMI_ALG_CIPHER_MODE_Unmarshal(TPMI_ALG_CIPHER_MODE *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_CIPHER_MODE_Marshal(TPMI_ALG_CIPHER_MODE *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_CIPHER_MODE_Marshal(source, buffer, size)                         \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:74 - Definition of TPMS_EMPTY Structure
TPM_RC
TPMS_EMPTY_Unmarshal(TPMS_EMPTY *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_EMPTY_Marshal(TPMS_EMPTY *source, BYTE **buffer, INT32 *size);

// Table 2:75 - Definition of TPMS_ALGORITHM_DESCRIPTION Structure
UINT16
TPMS_ALGORITHM_DESCRIPTION_Marshal(TPMS_ALGORITHM_DESCRIPTION *source,
            BYTE **buffer, INT32 *size);

// Table 2:76 - Definition of TPMU_HA Union
TPM_RC
TPMU_HA_Unmarshal(TPMU_HA *target, BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_HA_Marshal(TPMU_HA *source, BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:77 - Definition of TPMT_HA Structure
TPM_RC
TPMT_HA_Unmarshal(TPMT_HA *target, BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_HA_Marshal(TPMT_HA *source, BYTE **buffer, INT32 *size);

// Table 2:78 - Definition of TPM2B_DIGEST Structure
TPM_RC
TPM2B_DIGEST_Unmarshal(TPM2B_DIGEST *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_DIGEST_Marshal(TPM2B_DIGEST *source, BYTE **buffer, INT32 *size);

// Table 2:79 - Definition of TPM2B_DATA Structure
TPM_RC
TPM2B_DATA_Unmarshal(TPM2B_DATA *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_DATA_Marshal(TPM2B_DATA *source, BYTE **buffer, INT32 *size);

// Table 2:80 - Definition of Types for TPM2B_NONCE
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_NONCE_Unmarshal(TPM2B_NONCE *target, BYTE **buffer, INT32 *size);
#else
#define TPM2B_NONCE_Unmarshal(target, buffer, size)                                \
            TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM2B_NONCE_Marshal(TPM2B_NONCE *source, BYTE **buffer, INT32 *size);
#else
#define TPM2B_NONCE_Marshal(source, buffer, size)                                  \
            TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:81 - Definition of Types for TPM2B_AUTH
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_AUTH_Unmarshal(TPM2B_AUTH *target, BYTE **buffer, INT32 *size);
#else
#define TPM2B_AUTH_Unmarshal(target, buffer, size)                                 \
            TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM2B_AUTH_Marshal(TPM2B_AUTH *source, BYTE **buffer, INT32 *size);
#else
#define TPM2B_AUTH_Marshal(source, buffer, size)                                   \
            TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:82 - Definition of Types for TPM2B_OPERAND
#if !USE_MARSHALING_DEFINES
TPM_RC
TPM2B_OPERAND_Unmarshal(TPM2B_OPERAND *target, BYTE **buffer, INT32 *size);
#else
#define TPM2B_OPERAND_Unmarshal(target, buffer, size)                              \
            TPM2B_DIGEST_Unmarshal((TPM2B_DIGEST *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPM2B_OPERAND_Marshal(TPM2B_OPERAND *source, BYTE **buffer, INT32 *size);
#else
#define TPM2B_OPERAND_Marshal(source, buffer, size)                                \
            TPM2B_DIGEST_Marshal((TPM2B_DIGEST *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:83 - Definition of TPM2B_EVENT Structure
TPM_RC
TPM2B_EVENT_Unmarshal(TPM2B_EVENT *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_EVENT_Marshal(TPM2B_EVENT *source, BYTE **buffer, INT32 *size);

// Table 2:84 - Definition of TPM2B_MAX_BUFFER Structure
TPM_RC
TPM2B_MAX_BUFFER_Unmarshal(TPM2B_MAX_BUFFER *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_MAX_BUFFER_Marshal(TPM2B_MAX_BUFFER *source, BYTE **buffer, INT32 *size);

// Table 2:85 - Definition of TPM2B_MAX_NV_BUFFER Structure
TPM_RC
TPM2B_MAX_NV_BUFFER_Unmarshal(TPM2B_MAX_NV_BUFFER *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_MAX_NV_BUFFER_Marshal(TPM2B_MAX_NV_BUFFER *source,
            BYTE **buffer, INT32 *size);

// Table 2:86 - Definition of TPM2B_TIMEOUT Structure
TPM_RC
TPM2B_TIMEOUT_Unmarshal(TPM2B_TIMEOUT *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_TIMEOUT_Marshal(TPM2B_TIMEOUT *source, BYTE **buffer, INT32 *size);

// Table 2:87 - Definition of TPM2B_IV Structure
TPM_RC
TPM2B_IV_Unmarshal(TPM2B_IV *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_IV_Marshal(TPM2B_IV *source, BYTE **buffer, INT32 *size);

// Table 2:88 - Definition of TPMU_NAME Union
// Table 2:89 - Definition of TPM2B_NAME Structure
TPM_RC
TPM2B_NAME_Unmarshal(TPM2B_NAME *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_NAME_Marshal(TPM2B_NAME *source, BYTE **buffer, INT32 *size);

// Table 2:90 - Definition of TPMS_PCR_SELECT Structure
TPM_RC
TPMS_PCR_SELECT_Unmarshal(TPMS_PCR_SELECT *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_PCR_SELECT_Marshal(TPMS_PCR_SELECT *source, BYTE **buffer, INT32 *size);

// Table 2:91 - Definition of TPMS_PCR_SELECTION Structure
TPM_RC
TPMS_PCR_SELECTION_Unmarshal(TPMS_PCR_SELECTION *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_PCR_SELECTION_Marshal(TPMS_PCR_SELECTION *source, BYTE **buffer, INT32 *size);

// Table 2:94 - Definition of TPMT_TK_CREATION Structure
TPM_RC
TPMT_TK_CREATION_Unmarshal(TPMT_TK_CREATION *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_TK_CREATION_Marshal(TPMT_TK_CREATION *source, BYTE **buffer, INT32 *size);

// Table 2:95 - Definition of TPMT_TK_VERIFIED Structure
TPM_RC
TPMT_TK_VERIFIED_Unmarshal(TPMT_TK_VERIFIED *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_TK_VERIFIED_Marshal(TPMT_TK_VERIFIED *source, BYTE **buffer, INT32 *size);

// Table 2:96 - Definition of TPMT_TK_AUTH Structure
TPM_RC
TPMT_TK_AUTH_Unmarshal(TPMT_TK_AUTH *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_TK_AUTH_Marshal(TPMT_TK_AUTH *source, BYTE **buffer, INT32 *size);

// Table 2:97 - Definition of TPMT_TK_HASHCHECK Structure
TPM_RC
TPMT_TK_HASHCHECK_Unmarshal(TPMT_TK_HASHCHECK *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_TK_HASHCHECK_Marshal(TPMT_TK_HASHCHECK *source, BYTE **buffer, INT32 *size);

// Table 2:98 - Definition of TPMS_ALG_PROPERTY Structure
UINT16
TPMS_ALG_PROPERTY_Marshal(TPMS_ALG_PROPERTY *source, BYTE **buffer, INT32 *size);

// Table 2:99 - Definition of TPMS_TAGGED_PROPERTY Structure
UINT16
TPMS_TAGGED_PROPERTY_Marshal(TPMS_TAGGED_PROPERTY *source,
            BYTE **buffer, INT32 *size);

// Table 2:100 - Definition of TPMS_TAGGED_PCR_SELECT Structure
UINT16
TPMS_TAGGED_PCR_SELECT_Marshal(TPMS_TAGGED_PCR_SELECT *source,
            BYTE **buffer, INT32 *size);

// Table 2:101 - Definition of TPMS_TAGGED_POLICY Structure
UINT16
TPMS_TAGGED_POLICY_Marshal(TPMS_TAGGED_POLICY *source, BYTE **buffer, INT32 *size);

// Table 2:102 - Definition of TPML_CC Structure
TPM_RC
TPML_CC_Unmarshal(TPML_CC *target, BYTE **buffer, INT32 *size);
UINT16
TPML_CC_Marshal(TPML_CC *source, BYTE **buffer, INT32 *size);

// Table 2:103 - Definition of TPML_CCA Structure
UINT16
TPML_CCA_Marshal(TPML_CCA *source, BYTE **buffer, INT32 *size);

// Table 2:104 - Definition of TPML_ALG Structure
TPM_RC
TPML_ALG_Unmarshal(TPML_ALG *target, BYTE **buffer, INT32 *size);
UINT16
TPML_ALG_Marshal(TPML_ALG *source, BYTE **buffer, INT32 *size);

// Table 2:105 - Definition of TPML_HANDLE Structure
UINT16
TPML_HANDLE_Marshal(TPML_HANDLE *source, BYTE **buffer, INT32 *size);

// Table 2:106 - Definition of TPML_DIGEST Structure
TPM_RC
TPML_DIGEST_Unmarshal(TPML_DIGEST *target, BYTE **buffer, INT32 *size);
UINT16
TPML_DIGEST_Marshal(TPML_DIGEST *source, BYTE **buffer, INT32 *size);

// Table 2:107 - Definition of TPML_DIGEST_VALUES Structure
TPM_RC
TPML_DIGEST_VALUES_Unmarshal(TPML_DIGEST_VALUES *target,
            BYTE **buffer, INT32 *size);
UINT16
TPML_DIGEST_VALUES_Marshal(TPML_DIGEST_VALUES *source, BYTE **buffer, INT32 *size);

// Table 2:108 - Definition of TPML_PCR_SELECTION Structure
TPM_RC
TPML_PCR_SELECTION_Unmarshal(TPML_PCR_SELECTION *target,
            BYTE **buffer, INT32 *size);
UINT16
TPML_PCR_SELECTION_Marshal(TPML_PCR_SELECTION *source, BYTE **buffer, INT32 *size);

// Table 2:109 - Definition of TPML_ALG_PROPERTY Structure
UINT16
TPML_ALG_PROPERTY_Marshal(TPML_ALG_PROPERTY *source, BYTE **buffer, INT32 *size);

// Table 2:110 - Definition of TPML_TAGGED_TPM_PROPERTY Structure
UINT16
TPML_TAGGED_TPM_PROPERTY_Marshal(TPML_TAGGED_TPM_PROPERTY *source,
            BYTE **buffer, INT32 *size);

// Table 2:111 - Definition of TPML_TAGGED_PCR_PROPERTY Structure
UINT16
TPML_TAGGED_PCR_PROPERTY_Marshal(TPML_TAGGED_PCR_PROPERTY *source,
            BYTE **buffer, INT32 *size);

// Table 2:112 - Definition of TPML_ECC_CURVE Structure
#if ALG_ECC
UINT16
TPML_ECC_CURVE_Marshal(TPML_ECC_CURVE *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:113 - Definition of TPML_TAGGED_POLICY Structure
UINT16
TPML_TAGGED_POLICY_Marshal(TPML_TAGGED_POLICY *source, BYTE **buffer, INT32 *size);

// Table 2:114 - Definition of TPMU_CAPABILITIES Union
UINT16
TPMU_CAPABILITIES_Marshal(TPMU_CAPABILITIES *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:115 - Definition of TPMS_CAPABILITY_DATA Structure
UINT16
TPMS_CAPABILITY_DATA_Marshal(TPMS_CAPABILITY_DATA *source,
            BYTE **buffer, INT32 *size);

// Table 2:116 - Definition of TPMS_CLOCK_INFO Structure
TPM_RC
TPMS_CLOCK_INFO_Unmarshal(TPMS_CLOCK_INFO *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_CLOCK_INFO_Marshal(TPMS_CLOCK_INFO *source, BYTE **buffer, INT32 *size);

// Table 2:117 - Definition of TPMS_TIME_INFO Structure
TPM_RC
TPMS_TIME_INFO_Unmarshal(TPMS_TIME_INFO *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_TIME_INFO_Marshal(TPMS_TIME_INFO *source, BYTE **buffer, INT32 *size);

// Table 2:118 - Definition of TPMS_TIME_ATTEST_INFO Structure
UINT16
TPMS_TIME_ATTEST_INFO_Marshal(TPMS_TIME_ATTEST_INFO *source,
            BYTE **buffer, INT32 *size);

// Table 2:119 - Definition of TPMS_CERTIFY_INFO Structure
UINT16
TPMS_CERTIFY_INFO_Marshal(TPMS_CERTIFY_INFO *source, BYTE **buffer, INT32 *size);

// Table 2:120 - Definition of TPMS_QUOTE_INFO Structure
UINT16
TPMS_QUOTE_INFO_Marshal(TPMS_QUOTE_INFO *source, BYTE **buffer, INT32 *size);

// Table 2:121 - Definition of TPMS_COMMAND_AUDIT_INFO Structure
UINT16
TPMS_COMMAND_AUDIT_INFO_Marshal(TPMS_COMMAND_AUDIT_INFO *source,
            BYTE **buffer, INT32 *size);

// Table 2:122 - Definition of TPMS_SESSION_AUDIT_INFO Structure
UINT16
TPMS_SESSION_AUDIT_INFO_Marshal(TPMS_SESSION_AUDIT_INFO *source,
            BYTE **buffer, INT32 *size);

// Table 2:123 - Definition of TPMS_CREATION_INFO Structure
UINT16
TPMS_CREATION_INFO_Marshal(TPMS_CREATION_INFO *source, BYTE **buffer, INT32 *size);

// Table 2:124 - Definition of TPMS_NV_CERTIFY_INFO Structure
UINT16
TPMS_NV_CERTIFY_INFO_Marshal(TPMS_NV_CERTIFY_INFO *source,
            BYTE **buffer, INT32 *size);

// Table 2:125 - Definition of TPMS_NV_DIGEST_CERTIFY_INFO Structure
UINT16
TPMS_NV_DIGEST_CERTIFY_INFO_Marshal(TPMS_NV_DIGEST_CERTIFY_INFO *source,
            BYTE **buffer, INT32 *size);

// Table 2:126 - Definition of TPMI_ST_ATTEST Type
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ST_ATTEST_Marshal(TPMI_ST_ATTEST *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ST_ATTEST_Marshal(source, buffer, size)                               \
            TPM_ST_Marshal((TPM_ST *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:127 - Definition of TPMU_ATTEST Union
UINT16
TPMU_ATTEST_Marshal(TPMU_ATTEST *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:128 - Definition of TPMS_ATTEST Structure
UINT16
TPMS_ATTEST_Marshal(TPMS_ATTEST *source, BYTE **buffer, INT32 *size);

// Table 2:129 - Definition of TPM2B_ATTEST Structure
UINT16
TPM2B_ATTEST_Marshal(TPM2B_ATTEST *source, BYTE **buffer, INT32 *size);

// Table 2:130 - Definition of TPMS_AUTH_COMMAND Structure
TPM_RC
TPMS_AUTH_COMMAND_Unmarshal(TPMS_AUTH_COMMAND *target, BYTE **buffer, INT32 *size);

// Table 2:131 - Definition of TPMS_AUTH_RESPONSE Structure
UINT16
TPMS_AUTH_RESPONSE_Marshal(TPMS_AUTH_RESPONSE *source, BYTE **buffer, INT32 *size);

// Table 2:132 - Definition of TPMI_TDES_KEY_BITS Type
#if ALG_TDES
TPM_RC
TPMI_TDES_KEY_BITS_Unmarshal(TPMI_TDES_KEY_BITS *target,
            BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_TDES_KEY_BITS_Marshal(TPMI_TDES_KEY_BITS *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_TDES_KEY_BITS_Marshal(source, buffer, size)                           \
            TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_TDES

// Table 2:132 - Definition of TPMI_AES_KEY_BITS Type
#if ALG_AES
TPM_RC
TPMI_AES_KEY_BITS_Unmarshal(TPMI_AES_KEY_BITS *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_AES_KEY_BITS_Marshal(TPMI_AES_KEY_BITS *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_AES_KEY_BITS_Marshal(source, buffer, size)                            \
            TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_AES

// Table 2:132 - Definition of TPMI_SM4_KEY_BITS Type
#if ALG_SM4
TPM_RC
TPMI_SM4_KEY_BITS_Unmarshal(TPMI_SM4_KEY_BITS *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_SM4_KEY_BITS_Marshal(TPMI_SM4_KEY_BITS *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_SM4_KEY_BITS_Marshal(source, buffer, size)                            \
            TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_SM4

// Table 2:132 - Definition of TPMI_CAMELLIA_KEY_BITS Type
#if ALG_CAMELLIA
TPM_RC
TPMI_CAMELLIA_KEY_BITS_Unmarshal(TPMI_CAMELLIA_KEY_BITS *target,
            BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_CAMELLIA_KEY_BITS_Marshal(TPMI_CAMELLIA_KEY_BITS *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_CAMELLIA_KEY_BITS_Marshal(source, buffer, size)                       \
            TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_CAMELLIA

// Table 2:133 - Definition of TPMU_SYM_KEY_BITS Union
TPM_RC
TPMU_SYM_KEY_BITS_Unmarshal(TPMU_SYM_KEY_BITS *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SYM_KEY_BITS_Marshal(TPMU_SYM_KEY_BITS *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:134 - Definition of TPMU_SYM_MODE Union
TPM_RC
TPMU_SYM_MODE_Unmarshal(TPMU_SYM_MODE *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SYM_MODE_Marshal(TPMU_SYM_MODE *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:136 - Definition of TPMT_SYM_DEF Structure
TPM_RC
TPMT_SYM_DEF_Unmarshal(TPMT_SYM_DEF *target, BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_SYM_DEF_Marshal(TPMT_SYM_DEF *source, BYTE **buffer, INT32 *size);

// Table 2:137 - Definition of TPMT_SYM_DEF_OBJECT Structure
TPM_RC
TPMT_SYM_DEF_OBJECT_Unmarshal(TPMT_SYM_DEF_OBJECT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_SYM_DEF_OBJECT_Marshal(TPMT_SYM_DEF_OBJECT *source,
            BYTE **buffer, INT32 *size);

// Table 2:138 - Definition of TPM2B_SYM_KEY Structure
TPM_RC
TPM2B_SYM_KEY_Unmarshal(TPM2B_SYM_KEY *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_SYM_KEY_Marshal(TPM2B_SYM_KEY *source, BYTE **buffer, INT32 *size);

// Table 2:139 - Definition of TPMS_SYMCIPHER_PARMS Structure
TPM_RC
TPMS_SYMCIPHER_PARMS_Unmarshal(TPMS_SYMCIPHER_PARMS *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_SYMCIPHER_PARMS_Marshal(TPMS_SYMCIPHER_PARMS *source,
            BYTE **buffer, INT32 *size);

// Table 2:140 - Definition of TPM2B_LABEL Structure
TPM_RC
TPM2B_LABEL_Unmarshal(TPM2B_LABEL *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_LABEL_Marshal(TPM2B_LABEL *source, BYTE **buffer, INT32 *size);

// Table 2:141 - Definition of TPMS_DERIVE Structure
TPM_RC
TPMS_DERIVE_Unmarshal(TPMS_DERIVE *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_DERIVE_Marshal(TPMS_DERIVE *source, BYTE **buffer, INT32 *size);

// Table 2:142 - Definition of TPM2B_DERIVE Structure
TPM_RC
TPM2B_DERIVE_Unmarshal(TPM2B_DERIVE *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_DERIVE_Marshal(TPM2B_DERIVE *source, BYTE **buffer, INT32 *size);

// Table 2:143 - Definition of TPMU_SENSITIVE_CREATE Union
// Table 2:144 - Definition of TPM2B_SENSITIVE_DATA Structure
TPM_RC
TPM2B_SENSITIVE_DATA_Unmarshal(TPM2B_SENSITIVE_DATA *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_SENSITIVE_DATA_Marshal(TPM2B_SENSITIVE_DATA *source,
            BYTE **buffer, INT32 *size);

// Table 2:145 - Definition of TPMS_SENSITIVE_CREATE Structure
TPM_RC
TPMS_SENSITIVE_CREATE_Unmarshal(TPMS_SENSITIVE_CREATE *target,
            BYTE **buffer, INT32 *size);

// Table 2:146 - Definition of TPM2B_SENSITIVE_CREATE Structure
TPM_RC
TPM2B_SENSITIVE_CREATE_Unmarshal(TPM2B_SENSITIVE_CREATE *target,
            BYTE **buffer, INT32 *size);

// Table 2:147 - Definition of TPMS_SCHEME_HASH Structure
TPM_RC
TPMS_SCHEME_HASH_Unmarshal(TPMS_SCHEME_HASH *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_SCHEME_HASH_Marshal(TPMS_SCHEME_HASH *source, BYTE **buffer, INT32 *size);

// Table 2:148 - Definition of TPMS_SCHEME_ECDAA Structure
#if ALG_ECC
TPM_RC
TPMS_SCHEME_ECDAA_Unmarshal(TPMS_SCHEME_ECDAA *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_SCHEME_ECDAA_Marshal(TPMS_SCHEME_ECDAA *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:149 - Definition of TPMI_ALG_KEYEDHASH_SCHEME Type
TPM_RC
TPMI_ALG_KEYEDHASH_SCHEME_Unmarshal(TPMI_ALG_KEYEDHASH_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_KEYEDHASH_SCHEME_Marshal(TPMI_ALG_KEYEDHASH_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_KEYEDHASH_SCHEME_Marshal(source, buffer, size)                    \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:150 - Definition of Types for HMAC_SIG_SCHEME
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_HMAC_Unmarshal(TPMS_SCHEME_HMAC *target, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_HMAC_Unmarshal(target, buffer, size)                           \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SCHEME_HMAC_Marshal(TPMS_SCHEME_HMAC *source, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_HMAC_Marshal(source, buffer, size)                             \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:151 - Definition of TPMS_SCHEME_XOR Structure
TPM_RC
TPMS_SCHEME_XOR_Unmarshal(TPMS_SCHEME_XOR *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_SCHEME_XOR_Marshal(TPMS_SCHEME_XOR *source, BYTE **buffer, INT32 *size);

// Table 2:152 - Definition of TPMU_SCHEME_KEYEDHASH Union
TPM_RC
TPMU_SCHEME_KEYEDHASH_Unmarshal(TPMU_SCHEME_KEYEDHASH *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SCHEME_KEYEDHASH_Marshal(TPMU_SCHEME_KEYEDHASH *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:153 - Definition of TPMT_KEYEDHASH_SCHEME Structure
TPM_RC
TPMT_KEYEDHASH_SCHEME_Unmarshal(TPMT_KEYEDHASH_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_KEYEDHASH_SCHEME_Marshal(TPMT_KEYEDHASH_SCHEME *source,
            BYTE **buffer, INT32 *size);

// Table 2:154 - Definition of Types for RSA Signature Schemes
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_RSASSA_Unmarshal(TPMS_SIG_SCHEME_RSASSA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_RSASSA_Unmarshal(target, buffer, size)                     \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_RSASSA_Marshal(TPMS_SIG_SCHEME_RSASSA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_RSASSA_Marshal(source, buffer, size)                       \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_RSAPSS_Unmarshal(TPMS_SIG_SCHEME_RSAPSS *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_RSAPSS_Unmarshal(target, buffer, size)                     \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_RSAPSS_Marshal(TPMS_SIG_SCHEME_RSAPSS *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_RSAPSS_Marshal(source, buffer, size)                       \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:155 - Definition of Types for ECC Signature Schemes
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_ECDSA_Unmarshal(TPMS_SIG_SCHEME_ECDSA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECDSA_Unmarshal(target, buffer, size)                      \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_ECDSA_Marshal(TPMS_SIG_SCHEME_ECDSA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECDSA_Marshal(source, buffer, size)                        \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_SM2_Unmarshal(TPMS_SIG_SCHEME_SM2 *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_SM2_Unmarshal(target, buffer, size)                        \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_SM2_Marshal(TPMS_SIG_SCHEME_SM2 *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_SM2_Marshal(source, buffer, size)                          \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_ECSCHNORR_Unmarshal(TPMS_SIG_SCHEME_ECSCHNORR *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECSCHNORR_Unmarshal(target, buffer, size)                  \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_ECSCHNORR_Marshal(TPMS_SIG_SCHEME_ECSCHNORR *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECSCHNORR_Marshal(source, buffer, size)                    \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIG_SCHEME_ECDAA_Unmarshal(TPMS_SIG_SCHEME_ECDAA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECDAA_Unmarshal(target, buffer, size)                      \
            TPMS_SCHEME_ECDAA_Unmarshal((TPMS_SCHEME_ECDAA *)(target),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIG_SCHEME_ECDAA_Marshal(TPMS_SIG_SCHEME_ECDAA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIG_SCHEME_ECDAA_Marshal(source, buffer, size)                        \
            TPMS_SCHEME_ECDAA_Marshal((TPMS_SCHEME_ECDAA *)(source),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:156 - Definition of TPMU_SIG_SCHEME Union
TPM_RC
TPMU_SIG_SCHEME_Unmarshal(TPMU_SIG_SCHEME *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SIG_SCHEME_Marshal(TPMU_SIG_SCHEME *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:157 - Definition of TPMT_SIG_SCHEME Structure
TPM_RC
TPMT_SIG_SCHEME_Unmarshal(TPMT_SIG_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_SIG_SCHEME_Marshal(TPMT_SIG_SCHEME *source, BYTE **buffer, INT32 *size);

// Table 2:158 - Definition of Types for Encryption Schemes
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_ENC_SCHEME_OAEP_Unmarshal(TPMS_ENC_SCHEME_OAEP *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_ENC_SCHEME_OAEP_Unmarshal(target, buffer, size)                       \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_ENC_SCHEME_OAEP_Marshal(TPMS_ENC_SCHEME_OAEP *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_ENC_SCHEME_OAEP_Marshal(source, buffer, size)                         \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_ENC_SCHEME_RSAES_Unmarshal(TPMS_ENC_SCHEME_RSAES *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_ENC_SCHEME_RSAES_Unmarshal(target, buffer, size)                      \
            TPMS_EMPTY_Unmarshal((TPMS_EMPTY *)(target), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_ENC_SCHEME_RSAES_Marshal(TPMS_ENC_SCHEME_RSAES *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_ENC_SCHEME_RSAES_Marshal(source, buffer, size)                        \
            TPMS_EMPTY_Marshal((TPMS_EMPTY *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:159 - Definition of Types for ECC Key Exchange
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_KEY_SCHEME_ECDH_Unmarshal(TPMS_KEY_SCHEME_ECDH *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_KEY_SCHEME_ECDH_Unmarshal(target, buffer, size)                       \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_KEY_SCHEME_ECDH_Marshal(TPMS_KEY_SCHEME_ECDH *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_KEY_SCHEME_ECDH_Marshal(source, buffer, size)                         \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_KEY_SCHEME_ECMQV_Unmarshal(TPMS_KEY_SCHEME_ECMQV *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_KEY_SCHEME_ECMQV_Unmarshal(target, buffer, size)                      \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_KEY_SCHEME_ECMQV_Marshal(TPMS_KEY_SCHEME_ECMQV *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_KEY_SCHEME_ECMQV_Marshal(source, buffer, size)                        \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:160 - Definition of Types for KDF Schemes
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_MGF1_Unmarshal(TPMS_SCHEME_MGF1 *target, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_MGF1_Unmarshal(target, buffer, size)                           \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SCHEME_MGF1_Marshal(TPMS_SCHEME_MGF1 *source, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_MGF1_Marshal(source, buffer, size)                             \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_KDF1_SP800_56A_Unmarshal(TPMS_SCHEME_KDF1_SP800_56A *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF1_SP800_56A_Unmarshal(target, buffer, size)                 \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SCHEME_KDF1_SP800_56A_Marshal(TPMS_SCHEME_KDF1_SP800_56A *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF1_SP800_56A_Marshal(source, buffer, size)                   \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_KDF2_Unmarshal(TPMS_SCHEME_KDF2 *target, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF2_Unmarshal(target, buffer, size)                           \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SCHEME_KDF2_Marshal(TPMS_SCHEME_KDF2 *source, BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF2_Marshal(source, buffer, size)                             \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SCHEME_KDF1_SP800_108_Unmarshal(TPMS_SCHEME_KDF1_SP800_108 *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF1_SP800_108_Unmarshal(target, buffer, size)                 \
            TPMS_SCHEME_HASH_Unmarshal((TPMS_SCHEME_HASH *)(target),               \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SCHEME_KDF1_SP800_108_Marshal(TPMS_SCHEME_KDF1_SP800_108 *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SCHEME_KDF1_SP800_108_Marshal(source, buffer, size)                   \
            TPMS_SCHEME_HASH_Marshal((TPMS_SCHEME_HASH *)(source),                 \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:161 - Definition of TPMU_KDF_SCHEME Union
TPM_RC
TPMU_KDF_SCHEME_Unmarshal(TPMU_KDF_SCHEME *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_KDF_SCHEME_Marshal(TPMU_KDF_SCHEME *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:162 - Definition of TPMT_KDF_SCHEME Structure
TPM_RC
TPMT_KDF_SCHEME_Unmarshal(TPMT_KDF_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_KDF_SCHEME_Marshal(TPMT_KDF_SCHEME *source, BYTE **buffer, INT32 *size);

// Table 2:163 - Definition of TPMI_ALG_ASYM_SCHEME Type
TPM_RC
TPMI_ALG_ASYM_SCHEME_Unmarshal(TPMI_ALG_ASYM_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ASYM_SCHEME_Marshal(TPMI_ALG_ASYM_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_ASYM_SCHEME_Marshal(source, buffer, size)                         \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:164 - Definition of TPMU_ASYM_SCHEME Union
TPM_RC
TPMU_ASYM_SCHEME_Unmarshal(TPMU_ASYM_SCHEME *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_ASYM_SCHEME_Marshal(TPMU_ASYM_SCHEME *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:165 - Definition of TPMT_ASYM_SCHEME Structure
// Table 2:166 - Definition of TPMI_ALG_RSA_SCHEME Type
#if ALG_RSA
TPM_RC
TPMI_ALG_RSA_SCHEME_Unmarshal(TPMI_ALG_RSA_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_RSA_SCHEME_Marshal(TPMI_ALG_RSA_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_RSA_SCHEME_Marshal(source, buffer, size)                          \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:167 - Definition of TPMT_RSA_SCHEME Structure
#if ALG_RSA
TPM_RC
TPMT_RSA_SCHEME_Unmarshal(TPMT_RSA_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_RSA_SCHEME_Marshal(TPMT_RSA_SCHEME *source, BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:168 - Definition of TPMI_ALG_RSA_DECRYPT Type
#if ALG_RSA
TPM_RC
TPMI_ALG_RSA_DECRYPT_Unmarshal(TPMI_ALG_RSA_DECRYPT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_RSA_DECRYPT_Marshal(TPMI_ALG_RSA_DECRYPT *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_RSA_DECRYPT_Marshal(source, buffer, size)                         \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:169 - Definition of TPMT_RSA_DECRYPT Structure
#if ALG_RSA
TPM_RC
TPMT_RSA_DECRYPT_Unmarshal(TPMT_RSA_DECRYPT *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_RSA_DECRYPT_Marshal(TPMT_RSA_DECRYPT *source, BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:170 - Definition of TPM2B_PUBLIC_KEY_RSA Structure
#if ALG_RSA
TPM_RC
TPM2B_PUBLIC_KEY_RSA_Unmarshal(TPM2B_PUBLIC_KEY_RSA *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_PUBLIC_KEY_RSA_Marshal(TPM2B_PUBLIC_KEY_RSA *source,
            BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:171 - Definition of TPMI_RSA_KEY_BITS Type
#if ALG_RSA
TPM_RC
TPMI_RSA_KEY_BITS_Unmarshal(TPMI_RSA_KEY_BITS *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_RSA_KEY_BITS_Marshal(TPMI_RSA_KEY_BITS *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_RSA_KEY_BITS_Marshal(source, buffer, size)                            \
            TPM_KEY_BITS_Marshal((TPM_KEY_BITS *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:172 - Definition of TPM2B_PRIVATE_KEY_RSA Structure
#if ALG_RSA
TPM_RC
TPM2B_PRIVATE_KEY_RSA_Unmarshal(TPM2B_PRIVATE_KEY_RSA *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_PRIVATE_KEY_RSA_Marshal(TPM2B_PRIVATE_KEY_RSA *source,
            BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:173 - Definition of TPM2B_ECC_PARAMETER Structure
TPM_RC
TPM2B_ECC_PARAMETER_Unmarshal(TPM2B_ECC_PARAMETER *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_ECC_PARAMETER_Marshal(TPM2B_ECC_PARAMETER *source,
            BYTE **buffer, INT32 *size);

// Table 2:174 - Definition of TPMS_ECC_POINT Structure
#if ALG_ECC
TPM_RC
TPMS_ECC_POINT_Unmarshal(TPMS_ECC_POINT *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_ECC_POINT_Marshal(TPMS_ECC_POINT *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:175 - Definition of TPM2B_ECC_POINT Structure
#if ALG_ECC
TPM_RC
TPM2B_ECC_POINT_Unmarshal(TPM2B_ECC_POINT *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_ECC_POINT_Marshal(TPM2B_ECC_POINT *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:176 - Definition of TPMI_ALG_ECC_SCHEME Type
#if ALG_ECC
TPM_RC
TPMI_ALG_ECC_SCHEME_Unmarshal(TPMI_ALG_ECC_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_ECC_SCHEME_Marshal(TPMI_ALG_ECC_SCHEME *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_ECC_SCHEME_Marshal(source, buffer, size)                          \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:177 - Definition of TPMI_ECC_CURVE Type
#if ALG_ECC
TPM_RC
TPMI_ECC_CURVE_Unmarshal(TPMI_ECC_CURVE *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ECC_CURVE_Marshal(TPMI_ECC_CURVE *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ECC_CURVE_Marshal(source, buffer, size)                               \
            TPM_ECC_CURVE_Marshal((TPM_ECC_CURVE *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:178 - Definition of TPMT_ECC_SCHEME Structure
#if ALG_ECC
TPM_RC
TPMT_ECC_SCHEME_Unmarshal(TPMT_ECC_SCHEME *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_ECC_SCHEME_Marshal(TPMT_ECC_SCHEME *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:179 - Definition of TPMS_ALGORITHM_DETAIL_ECC Structure
#if ALG_ECC
UINT16
TPMS_ALGORITHM_DETAIL_ECC_Marshal(TPMS_ALGORITHM_DETAIL_ECC *source,
            BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:180 - Definition of TPMS_SIGNATURE_RSA Structure
#if ALG_RSA
TPM_RC
TPMS_SIGNATURE_RSA_Unmarshal(TPMS_SIGNATURE_RSA *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_SIGNATURE_RSA_Marshal(TPMS_SIGNATURE_RSA *source, BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:181 - Definition of Types for Signature
#if ALG_RSA
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_RSASSA_Unmarshal(TPMS_SIGNATURE_RSASSA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_RSASSA_Unmarshal(target, buffer, size)                      \
            TPMS_SIGNATURE_RSA_Unmarshal((TPMS_SIGNATURE_RSA *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_RSASSA_Marshal(TPMS_SIGNATURE_RSASSA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_RSASSA_Marshal(source, buffer, size)                        \
            TPMS_SIGNATURE_RSA_Marshal((TPMS_SIGNATURE_RSA *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_RSAPSS_Unmarshal(TPMS_SIGNATURE_RSAPSS *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_RSAPSS_Unmarshal(target, buffer, size)                      \
            TPMS_SIGNATURE_RSA_Unmarshal((TPMS_SIGNATURE_RSA *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_RSAPSS_Marshal(TPMS_SIGNATURE_RSAPSS *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_RSAPSS_Marshal(source, buffer, size)                        \
            TPMS_SIGNATURE_RSA_Marshal((TPMS_SIGNATURE_RSA *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_RSA

// Table 2:182 - Definition of TPMS_SIGNATURE_ECC Structure
#if ALG_ECC
TPM_RC
TPMS_SIGNATURE_ECC_Unmarshal(TPMS_SIGNATURE_ECC *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_SIGNATURE_ECC_Marshal(TPMS_SIGNATURE_ECC *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:183 - Definition of Types for TPMS_SIGNATURE_ECC
#if ALG_ECC
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_ECDAA_Unmarshal(TPMS_SIGNATURE_ECDAA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECDAA_Unmarshal(target, buffer, size)                       \
            TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_ECDAA_Marshal(TPMS_SIGNATURE_ECDAA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECDAA_Marshal(source, buffer, size)                         \
            TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_ECDSA_Unmarshal(TPMS_SIGNATURE_ECDSA *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECDSA_Unmarshal(target, buffer, size)                       \
            TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_ECDSA_Marshal(TPMS_SIGNATURE_ECDSA *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECDSA_Marshal(source, buffer, size)                         \
            TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_SM2_Unmarshal(TPMS_SIGNATURE_SM2 *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_SM2_Unmarshal(target, buffer, size)                         \
            TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_SM2_Marshal(TPMS_SIGNATURE_SM2 *source, BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_SM2_Marshal(source, buffer, size)                           \
            TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
TPM_RC
TPMS_SIGNATURE_ECSCHNORR_Unmarshal(TPMS_SIGNATURE_ECSCHNORR *target,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECSCHNORR_Unmarshal(target, buffer, size)                   \
            TPMS_SIGNATURE_ECC_Unmarshal((TPMS_SIGNATURE_ECC *)(target),           \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#if !USE_MARSHALING_DEFINES
UINT16
TPMS_SIGNATURE_ECSCHNORR_Marshal(TPMS_SIGNATURE_ECSCHNORR *source,
            BYTE **buffer, INT32 *size);
#else
#define TPMS_SIGNATURE_ECSCHNORR_Marshal(source, buffer, size)                     \
            TPMS_SIGNATURE_ECC_Marshal((TPMS_SIGNATURE_ECC *)(source),             \
            (buffer),                                                              \
            (size))
#endif // !USE_MARSHALING_DEFINES
#endif // ALG_ECC

// Table 2:184 - Definition of TPMU_SIGNATURE Union
TPM_RC
TPMU_SIGNATURE_Unmarshal(TPMU_SIGNATURE *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SIGNATURE_Marshal(TPMU_SIGNATURE *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:185 - Definition of TPMT_SIGNATURE Structure
TPM_RC
TPMT_SIGNATURE_Unmarshal(TPMT_SIGNATURE *target,
            BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_SIGNATURE_Marshal(TPMT_SIGNATURE *source, BYTE **buffer, INT32 *size);

// Table 2:186 - Definition of TPMU_ENCRYPTED_SECRET Union
TPM_RC
TPMU_ENCRYPTED_SECRET_Unmarshal(TPMU_ENCRYPTED_SECRET *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_ENCRYPTED_SECRET_Marshal(TPMU_ENCRYPTED_SECRET *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:187 - Definition of TPM2B_ENCRYPTED_SECRET Structure
TPM_RC
TPM2B_ENCRYPTED_SECRET_Unmarshal(TPM2B_ENCRYPTED_SECRET *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_ENCRYPTED_SECRET_Marshal(TPM2B_ENCRYPTED_SECRET *source,
            BYTE **buffer, INT32 *size);

// Table 2:188 - Definition of TPMI_ALG_PUBLIC Type
TPM_RC
TPMI_ALG_PUBLIC_Unmarshal(TPMI_ALG_PUBLIC *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPMI_ALG_PUBLIC_Marshal(TPMI_ALG_PUBLIC *source, BYTE **buffer, INT32 *size);
#else
#define TPMI_ALG_PUBLIC_Marshal(source, buffer, size)                              \
            TPM_ALG_ID_Marshal((TPM_ALG_ID *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:189 - Definition of TPMU_PUBLIC_ID Union
TPM_RC
TPMU_PUBLIC_ID_Unmarshal(TPMU_PUBLIC_ID *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_PUBLIC_ID_Marshal(TPMU_PUBLIC_ID *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:190 - Definition of TPMS_KEYEDHASH_PARMS Structure
TPM_RC
TPMS_KEYEDHASH_PARMS_Unmarshal(TPMS_KEYEDHASH_PARMS *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_KEYEDHASH_PARMS_Marshal(TPMS_KEYEDHASH_PARMS *source,
            BYTE **buffer, INT32 *size);

// Table 2:191 - Definition of TPMS_ASYM_PARMS Structure
// Table 2:192 - Definition of TPMS_RSA_PARMS Structure
#if ALG_RSA
TPM_RC
TPMS_RSA_PARMS_Unmarshal(TPMS_RSA_PARMS *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_RSA_PARMS_Marshal(TPMS_RSA_PARMS *source, BYTE **buffer, INT32 *size);
#endif // ALG_RSA

// Table 2:193 - Definition of TPMS_ECC_PARMS Structure
#if ALG_ECC
TPM_RC
TPMS_ECC_PARMS_Unmarshal(TPMS_ECC_PARMS *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_ECC_PARMS_Marshal(TPMS_ECC_PARMS *source, BYTE **buffer, INT32 *size);
#endif // ALG_ECC

// Table 2:194 - Definition of TPMU_PUBLIC_PARMS Union
TPM_RC
TPMU_PUBLIC_PARMS_Unmarshal(TPMU_PUBLIC_PARMS *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_PUBLIC_PARMS_Marshal(TPMU_PUBLIC_PARMS *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:195 - Definition of TPMT_PUBLIC_PARMS Structure
TPM_RC
TPMT_PUBLIC_PARMS_Unmarshal(TPMT_PUBLIC_PARMS *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_PUBLIC_PARMS_Marshal(TPMT_PUBLIC_PARMS *source, BYTE **buffer, INT32 *size);

// Table 2:196 - Definition of TPMT_PUBLIC Structure
TPM_RC
TPMT_PUBLIC_Unmarshal(TPMT_PUBLIC *target, BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPMT_PUBLIC_Marshal(TPMT_PUBLIC *source, BYTE **buffer, INT32 *size);

// Table 2:197 - Definition of TPM2B_PUBLIC Structure
TPM_RC
TPM2B_PUBLIC_Unmarshal(TPM2B_PUBLIC *target, BYTE **buffer, INT32 *size, BOOL flag);
UINT16
TPM2B_PUBLIC_Marshal(TPM2B_PUBLIC *source, BYTE **buffer, INT32 *size);

// Table 2:198 - Definition of TPM2B_TEMPLATE Structure
TPM_RC
TPM2B_TEMPLATE_Unmarshal(TPM2B_TEMPLATE *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_TEMPLATE_Marshal(TPM2B_TEMPLATE *source, BYTE **buffer, INT32 *size);

// Table 2:199 - Definition of TPM2B_PRIVATE_VENDOR_SPECIFIC Structure
TPM_RC
TPM2B_PRIVATE_VENDOR_SPECIFIC_Unmarshal(TPM2B_PRIVATE_VENDOR_SPECIFIC *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_PRIVATE_VENDOR_SPECIFIC_Marshal(TPM2B_PRIVATE_VENDOR_SPECIFIC *source,
            BYTE **buffer, INT32 *size);

// Table 2:200 - Definition of TPMU_SENSITIVE_COMPOSITE Union
TPM_RC
TPMU_SENSITIVE_COMPOSITE_Unmarshal(TPMU_SENSITIVE_COMPOSITE *target,
            BYTE **buffer, INT32 *size, UINT32 selector);
UINT16
TPMU_SENSITIVE_COMPOSITE_Marshal(TPMU_SENSITIVE_COMPOSITE *source,
            BYTE **buffer, INT32 *size, UINT32 selector);

// Table 2:201 - Definition of TPMT_SENSITIVE Structure
TPM_RC
TPMT_SENSITIVE_Unmarshal(TPMT_SENSITIVE *target, BYTE **buffer, INT32 *size);
UINT16
TPMT_SENSITIVE_Marshal(TPMT_SENSITIVE *source, BYTE **buffer, INT32 *size);

// Table 2:202 - Definition of TPM2B_SENSITIVE Structure
TPM_RC
TPM2B_SENSITIVE_Unmarshal(TPM2B_SENSITIVE *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_SENSITIVE_Marshal(TPM2B_SENSITIVE *source, BYTE **buffer, INT32 *size);

// Table 2:203 - Definition of _PRIVATE Structure
// Table 2:204 - Definition of TPM2B_PRIVATE Structure
TPM_RC
TPM2B_PRIVATE_Unmarshal(TPM2B_PRIVATE *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_PRIVATE_Marshal(TPM2B_PRIVATE *source, BYTE **buffer, INT32 *size);

// Table 2:205 - Definition of TPMS_ID_OBJECT Structure
// Table 2:206 - Definition of TPM2B_ID_OBJECT Structure
TPM_RC
TPM2B_ID_OBJECT_Unmarshal(TPM2B_ID_OBJECT *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_ID_OBJECT_Marshal(TPM2B_ID_OBJECT *source, BYTE **buffer, INT32 *size);

// Table 2:207 - Definition of TPM_NV_INDEX Bits
#if !USE_MARSHALING_DEFINES
UINT16
TPM_NV_INDEX_Marshal(TPM_NV_INDEX *source, BYTE **buffer, INT32 *size);
#else
#define TPM_NV_INDEX_Marshal(source, buffer, size)                                 \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:208 - Definition of TPM_NT Constants
// Table 2:209 - Definition of TPMS_NV_PIN_COUNTER_PARAMETERS Structure
TPM_RC
TPMS_NV_PIN_COUNTER_PARAMETERS_Unmarshal(TPMS_NV_PIN_COUNTER_PARAMETERS *target,
            BYTE **buffer, INT32 *size);
UINT16
TPMS_NV_PIN_COUNTER_PARAMETERS_Marshal(TPMS_NV_PIN_COUNTER_PARAMETERS *source,
            BYTE **buffer, INT32 *size);

// Table 2:210 - Definition of TPMA_NV Bits
TPM_RC
TPMA_NV_Unmarshal(TPMA_NV *target, BYTE **buffer, INT32 *size);

#if !USE_MARSHALING_DEFINES
UINT16
TPMA_NV_Marshal(TPMA_NV *source, BYTE **buffer, INT32 *size);
#else
#define TPMA_NV_Marshal(source, buffer, size)                                      \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:211 - Definition of TPMS_NV_PUBLIC Structure
TPM_RC
TPMS_NV_PUBLIC_Unmarshal(TPMS_NV_PUBLIC *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_NV_PUBLIC_Marshal(TPMS_NV_PUBLIC *source, BYTE **buffer, INT32 *size);

// Table 2:212 - Definition of TPM2B_NV_PUBLIC Structure
TPM_RC
TPM2B_NV_PUBLIC_Unmarshal(TPM2B_NV_PUBLIC *target, BYTE **buffer, INT32 *size);
UINT16
TPM2B_NV_PUBLIC_Marshal(TPM2B_NV_PUBLIC *source, BYTE **buffer, INT32 *size);

// Table 2:213 - Definition of TPM2B_CONTEXT_SENSITIVE Structure
TPM_RC
TPM2B_CONTEXT_SENSITIVE_Unmarshal(TPM2B_CONTEXT_SENSITIVE *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_CONTEXT_SENSITIVE_Marshal(TPM2B_CONTEXT_SENSITIVE *source,
            BYTE **buffer, INT32 *size);

// Table 2:214 - Definition of TPMS_CONTEXT_DATA Structure
TPM_RC
TPMS_CONTEXT_DATA_Unmarshal(TPMS_CONTEXT_DATA *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_CONTEXT_DATA_Marshal(TPMS_CONTEXT_DATA *source, BYTE **buffer, INT32 *size);

// Table 2:215 - Definition of TPM2B_CONTEXT_DATA Structure
TPM_RC
TPM2B_CONTEXT_DATA_Unmarshal(TPM2B_CONTEXT_DATA *target,
            BYTE **buffer, INT32 *size);
UINT16
TPM2B_CONTEXT_DATA_Marshal(TPM2B_CONTEXT_DATA *source, BYTE **buffer, INT32 *size);

// Table 2:216 - Definition of TPMS_CONTEXT Structure
TPM_RC
TPMS_CONTEXT_Unmarshal(TPMS_CONTEXT *target, BYTE **buffer, INT32 *size);
UINT16
TPMS_CONTEXT_Marshal(TPMS_CONTEXT *source, BYTE **buffer, INT32 *size);

// Table 2:218 - Definition of TPMS_CREATION_DATA Structure
UINT16
TPMS_CREATION_DATA_Marshal(TPMS_CREATION_DATA *source, BYTE **buffer, INT32 *size);

// Table 2:219 - Definition of TPM2B_CREATION_DATA Structure
UINT16
TPM2B_CREATION_DATA_Marshal(TPM2B_CREATION_DATA *source,
            BYTE **buffer, INT32 *size);

// Table 2:220 - Definition of TPM_AT Constants
TPM_RC
TPM_AT_Unmarshal(TPM_AT *target, BYTE **buffer, INT32 *size);
#if !USE_MARSHALING_DEFINES
UINT16
TPM_AT_Marshal(TPM_AT *source, BYTE **buffer, INT32 *size);
#else
#define TPM_AT_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:221 - Definition of TPM_AE Constants
#if !USE_MARSHALING_DEFINES
UINT16
TPM_AE_Marshal(TPM_AE *source, BYTE **buffer, INT32 *size);
#else
#define TPM_AE_Marshal(source, buffer, size)                                       \
            UINT32_Marshal((UINT32 *)(source), (buffer), (size))
#endif // !USE_MARSHALING_DEFINES

// Table 2:222 - Definition of TPMS_AC_OUTPUT Structure
UINT16
TPMS_AC_OUTPUT_Marshal(TPMS_AC_OUTPUT *source, BYTE **buffer, INT32 *size);

// Table 2:223 - Definition of TPML_AC_CAPABILITIES Structure
UINT16
TPML_AC_CAPABILITIES_Marshal(TPML_AC_CAPABILITIES *source,
            BYTE **buffer, INT32 *size);

// Array Marshal/Unmarshal for BYTE
TPM_RC
BYTE_Array_Unmarshal(BYTE *target, BYTE **buffer, INT32 *size, INT32 count);
UINT16
BYTE_Array_Marshal(BYTE *source, BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPM2B_DIGEST
TPM_RC
TPM2B_DIGEST_Array_Unmarshal(TPM2B_DIGEST *target,
            BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPM2B_DIGEST_Array_Marshal(TPM2B_DIGEST *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMA_CC
UINT16
TPMA_CC_Array_Marshal(TPMA_CC *source, BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMS_AC_OUTPUT
UINT16
TPMS_AC_OUTPUT_Array_Marshal(TPMS_AC_OUTPUT *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMS_ALG_PROPERTY
UINT16
TPMS_ALG_PROPERTY_Array_Marshal(TPMS_ALG_PROPERTY *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPMS_PCR_SELECTION
TPM_RC
TPMS_PCR_SELECTION_Array_Unmarshal(TPMS_PCR_SELECTION *target,
            BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPMS_PCR_SELECTION_Array_Marshal(TPMS_PCR_SELECTION *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMS_TAGGED_PCR_SELECT
UINT16
TPMS_TAGGED_PCR_SELECT_Array_Marshal(TPMS_TAGGED_PCR_SELECT *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMS_TAGGED_POLICY
UINT16
TPMS_TAGGED_POLICY_Array_Marshal(TPMS_TAGGED_POLICY *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal for TPMS_TAGGED_PROPERTY
UINT16
TPMS_TAGGED_PROPERTY_Array_Marshal(TPMS_TAGGED_PROPERTY *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPMT_HA
TPM_RC
TPMT_HA_Array_Unmarshal(TPMT_HA *target,
            BYTE **buffer, INT32 *size, BOOL flag, INT32 count);
UINT16
TPMT_HA_Array_Marshal(TPMT_HA *source, BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPM_ALG_ID
TPM_RC
TPM_ALG_ID_Array_Unmarshal(TPM_ALG_ID *target,
            BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPM_ALG_ID_Array_Marshal(TPM_ALG_ID *source,
            BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPM_CC
TPM_RC
TPM_CC_Array_Unmarshal(TPM_CC *target, BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPM_CC_Array_Marshal(TPM_CC *source, BYTE **buffer, INT32 *size, INT32 count);

// Array Marshal/Unmarshal for TPM_ECC_CURVE
#if ALG_ECC
TPM_RC
TPM_ECC_CURVE_Array_Unmarshal(TPM_ECC_CURVE *target,
            BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPM_ECC_CURVE_Array_Marshal(TPM_ECC_CURVE *source,
            BYTE **buffer, INT32 *size, INT32 count);
#endif // ALG_ECC

// Array Marshal/Unmarshal for TPM_HANDLE
TPM_RC
TPM_HANDLE_Array_Unmarshal(TPM_HANDLE *target,
            BYTE **buffer, INT32 *size, INT32 count);
UINT16
TPM_HANDLE_Array_Marshal(TPM_HANDLE *source,
            BYTE **buffer, INT32 *size, INT32 count);
#endif // _MARSHAL_FP_H_
