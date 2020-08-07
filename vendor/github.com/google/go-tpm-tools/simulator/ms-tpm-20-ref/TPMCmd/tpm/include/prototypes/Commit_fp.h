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
 *  Date: Mar 28, 2019  Time: 08:25:17PM
 */

#if CC_Commit // Command must be enabled

#ifndef _Commit_FP_H_
#define _Commit_FP_H_

// Input structure definition
typedef struct {
    TPMI_DH_OBJECT              signHandle;
    TPM2B_ECC_POINT             P1;
    TPM2B_SENSITIVE_DATA        s2;
    TPM2B_ECC_PARAMETER         y2;
} Commit_In;

// Output structure definition
typedef struct {
    TPM2B_ECC_POINT             K;
    TPM2B_ECC_POINT             L;
    TPM2B_ECC_POINT             E;
    UINT16                      counter;
} Commit_Out;

// Response code modifiers
#define RC_Commit_signHandle    (TPM_RC_H + TPM_RC_1)
#define RC_Commit_P1            (TPM_RC_P + TPM_RC_1)
#define RC_Commit_s2            (TPM_RC_P + TPM_RC_2)
#define RC_Commit_y2            (TPM_RC_P + TPM_RC_3)

// Function prototype
TPM_RC
TPM2_Commit(
    Commit_In                   *in,
    Commit_Out                  *out
);

#endif  // _Commit_FP_H_
#endif  // CC_Commit
