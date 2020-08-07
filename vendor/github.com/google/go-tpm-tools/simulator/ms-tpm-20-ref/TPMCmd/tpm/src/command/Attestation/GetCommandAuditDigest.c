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
#include "Tpm.h"
#include "Attest_spt_fp.h"
#include "GetCommandAuditDigest_fp.h"

#if CC_GetCommandAuditDigest  // Conditional expansion of this file

/*(See part 3 specification)
// Get current value of command audit log
*/
//  Return Type: TPM_RC
//      TPM_RC_KEY          key referenced by 'signHandle' is not a signing key
//      TPM_RC_SCHEME       'inScheme' is incompatible with 'signHandle' type; or
//                          both 'scheme' and key's default scheme are empty; or
//                          'scheme' is empty while key's default scheme requires
//                          explicit input scheme (split signing); or
//                          non-empty default key scheme differs from 'scheme'
//      TPM_RC_VALUE        digest generated for the given 'scheme' is greater than
//                          the modulus of 'signHandle' (for an RSA key);
//                          invalid commit status or failed to generate "r" value
//                          (for an ECC key)
TPM_RC
TPM2_GetCommandAuditDigest(
    GetCommandAuditDigest_In    *in,            // IN: input parameter list
    GetCommandAuditDigest_Out   *out            // OUT: output parameter list
    )
{
    TPM_RC                  result;
    TPMS_ATTEST             auditInfo;
    OBJECT                 *signObject = HandleToObject(in->signHandle);
// Input validation
    if(!IsSigningObject(signObject))
        return TPM_RCS_KEY + RC_GetCommandAuditDigest_signHandle;
    if(!CryptSelectSignScheme(signObject, &in->inScheme))
        return TPM_RCS_SCHEME + RC_GetCommandAuditDigest_inScheme;

// Command Output
    // Fill in attest information common fields
    FillInAttestInfo(in->signHandle, &in->inScheme, &in->qualifyingData,
                     &auditInfo);

    // CommandAuditDigest specific fields
    auditInfo.type = TPM_ST_ATTEST_COMMAND_AUDIT;
    auditInfo.attested.commandAudit.digestAlg = gp.auditHashAlg;
    auditInfo.attested.commandAudit.auditCounter = gp.auditCounter;

    // Copy command audit log
    auditInfo.attested.commandAudit.auditDigest = gr.commandAuditDigest;
    CommandAuditGetDigest(&auditInfo.attested.commandAudit.commandDigest);

    // Sign attestation structure.  A NULL signature will be returned if
    // signHandle is TPM_RH_NULL.  A TPM_RC_NV_UNAVAILABLE, TPM_RC_NV_RATE,
    // TPM_RC_VALUE, TPM_RC_SCHEME or TPM_RC_ATTRIBUTES error may be returned at
    // this point
    result = SignAttestInfo(signObject, &in->inScheme, &auditInfo,
                            &in->qualifyingData, &out->auditInfo,
                            &out->signature);
    // Internal Data Update
    if(result == TPM_RC_SUCCESS && in->signHandle != TPM_RH_NULL)
        // Reset log
        gr.commandAuditDigest.t.size = 0;

    return result;
}

#endif // CC_GetCommandAuditDigest