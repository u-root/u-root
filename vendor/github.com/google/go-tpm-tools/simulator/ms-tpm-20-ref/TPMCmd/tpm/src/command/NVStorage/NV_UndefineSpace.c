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
#include "NV_UndefineSpace_fp.h"

#if CC_NV_UndefineSpace  // Conditional expansion of this file

/*(See part 3 specification)
// Delete an NV Index
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES               TPMA_NV_POLICY_DELETE is SET in the Index
//                                      referenced by 'nvIndex' so this command may
//                                      not be used to delete this Index (see
//                                      TPM2_NV_UndefineSpaceSpecial())
//      TPM_RC_NV_AUTHORIZATION         attempt to use ownerAuth to delete an index
//                                      created by the platform
//
TPM_RC
TPM2_NV_UndefineSpace(
    NV_UndefineSpace_In     *in             // IN: input parameter list
    )
{
    NV_REF           locator;
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, &locator);

// Input Validation
    // This command can't be used to delete an index with TPMA_NV_POLICY_DELETE SET
    if(IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, POLICY_DELETE))   
        return TPM_RCS_ATTRIBUTES + RC_NV_UndefineSpace_nvIndex;

    // The owner may only delete an index that was defined with ownerAuth. The
    // platform may delete an index that was created with either authorization.
    if(in->authHandle == TPM_RH_OWNER
       && IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, PLATFORMCREATE))  
        return TPM_RC_NV_AUTHORIZATION;

// Internal Data Update

    // Call implementation dependent internal routine to delete NV index
    return NvDeleteIndex(nvIndex, locator);
}

#endif // CC_NV_UndefineSpace