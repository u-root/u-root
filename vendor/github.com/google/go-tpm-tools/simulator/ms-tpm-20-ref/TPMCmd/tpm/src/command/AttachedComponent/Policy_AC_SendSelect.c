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
#include "Policy_AC_SendSelect_fp.h"

#if CC_Policy_AC_SendSelect      // Conditional expansion of this file

/*(See part 3 specification)
// allows qualification of attached component and object to be sent.
*/
//  Return Type: TPM_RC
//      TPM_RC_COMMAND_CODE   'commandCode' of 'policySession; is not empty
//      TPM_RC_CPHASH         'cpHash' of 'policySession' is not empty
TPM_RC
TPM2_Policy_AC_SendSelect(
    Policy_AC_SendSelect_In  *in                // IN: input parameter list
    )
{
    SESSION         *session;
    HASH_STATE      hashState;
    TPM_CC          commandCode = TPM_CC_Policy_AC_SendSelect;

// Input Validation

    // Get pointer to the session structure
    session = SessionGet(in->policySession);

    // cpHash in session context must be empty
    if(session->u1.cpHash.t.size != 0)
        return TPM_RC_CPHASH;
    // commandCode in session context must be empty
    if(session->commandCode != 0)
        return TPM_RC_COMMAND_CODE;
// Internal Data Update
    // Update name hash
    session->u1.cpHash.t.size = CryptHashStart(&hashState, session->authHashAlg);

    //  add objectName
    CryptDigestUpdate2B(&hashState, &in->objectName.b);

    // add authHandleName
    CryptDigestUpdate2B(&hashState, &in->authHandleName.b);

    //  add ac name
    CryptDigestUpdate2B(&hashState, &in->acName.b);

    //  complete hash
    CryptHashEnd2B(&hashState, &session->u1.cpHash.b);

    // update policy hash
    // Old policyDigest size should be the same as the new policyDigest size since
    // they are using the same hash algorithm
    session->u2.policyDigest.t.size
        = CryptHashStart(&hashState, session->authHashAlg);
//  add old policy
    CryptDigestUpdate2B(&hashState, &session->u2.policyDigest.b);

    //  add command code
    CryptDigestUpdateInt(&hashState, sizeof(TPM_CC), commandCode);

    //  add objectName
    if(in->includeObject == YES)
        CryptDigestUpdate2B(&hashState, &in->objectName.b);

    //  add authHandleName
    CryptDigestUpdate2B(&hashState, &in->authHandleName.b);

    // add acName
    CryptDigestUpdate2B(&hashState, &in->acName.b);

    //  add includeObject
    CryptDigestUpdateInt(&hashState, sizeof(TPMI_YES_NO), in->includeObject);

    //  complete digest
    CryptHashEnd2B(&hashState, &session->u2.policyDigest.b);

    // set commandCode in session context
    session->commandCode = TPM_CC_AC_Send;

    return TPM_RC_SUCCESS;
}

#endif // CC_PolicyDuplicationSelect