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
 *  Date: Mar 28, 2019  Time: 08:25:18PM
 */

#ifndef    _POLICY_SPT_FP_H_
#define    _POLICY_SPT_FP_H_

//** Functions
//*** PolicyParameterChecks()
// This function validates the common parameters of TPM2_PolicySiged()
// and TPM2_PolicySecret(). The common parameters are 'nonceTPM',
// 'expiration', and 'cpHashA'.
TPM_RC
PolicyParameterChecks(
    SESSION         *session,
    UINT64           authTimeout,
    TPM2B_DIGEST    *cpHashA,
    TPM2B_NONCE     *nonce,
    TPM_RC           blameNonce,
    TPM_RC           blameCpHash,
    TPM_RC           blameExpiration
);

//*** PolicyContextUpdate()
// Update policy hash
//      Update the policyDigest in policy session by extending policyRef and
//      objectName to it. This will also update the cpHash if it is present.
//  Return Type: void
void
PolicyContextUpdate(
    TPM_CC           commandCode,   // IN: command code
    TPM2B_NAME      *name,          // IN: name of entity
    TPM2B_NONCE     *ref,           // IN: the reference data
    TPM2B_DIGEST    *cpHash,        // IN: the cpHash (optional)
    UINT64           policyTimeout, // IN: the timeout value for the policy
    SESSION         *session        // IN/OUT: policy session to be updated
);

//*** ComputeAuthTimeout()
// This function is used to determine what the authorization timeout value for
// the session should be.
UINT64
ComputeAuthTimeout(
    SESSION         *session,               // IN: the session containing the time
                                            //     values
    INT32            expiration,            // IN: either the number of seconds from
                                            //     the start of the session or the
                                            //     time in g_timer;
    TPM2B_NONCE     *nonce                  // IN: indicator of the time base
);

//*** PolicyDigestClear()
// Function to reset the policyDigest of a session
void
PolicyDigestClear(
    SESSION         *session
);

BOOL
PolicySptCheckCondition(
    TPM_EO          operation,
    BYTE            *opA,
    BYTE            *opB,
    UINT16           size
);

#endif  // _POLICY_SPT_FP_H_
