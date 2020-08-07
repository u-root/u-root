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
#include "NV_UndefineSpaceSpecial_fp.h"
#include "SessionProcess_fp.h"

#if CC_NV_UndefineSpaceSpecial  // Conditional expansion of this file

/*(See part 3 specification)
// Delete a NV index that requires policy to delete.
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES               TPMA_NV_POLICY_DELETE is not SET in the
//                                      Index referenced by 'nvIndex'
TPM_RC
TPM2_NV_UndefineSpaceSpecial(
    NV_UndefineSpaceSpecial_In  *in             // IN: input parameter list
    )
{
    TPM_RC           result;
    NV_REF           locator;
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, &locator);
// Input Validation
    // This operation only applies when the TPMA_NV_POLICY_DELETE attribute is SET
    if(!IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, POLICY_DELETE))   
        return TPM_RCS_ATTRIBUTES + RC_NV_UndefineSpaceSpecial_nvIndex;
// Internal Data Update
    // Call implementation dependent internal routine to delete NV index
    result = NvDeleteIndex(nvIndex, locator);

    // If we just removed the index providing the authorization, make sure that the
    // authorization session computation is modified so that it doesn't try to
    // access the authValue of the just deleted index
    if(result == TPM_RC_SUCCESS)
        SessionRemoveAssociationToHandle(in->nvIndex);
    return result;
}

#endif // CC_NV_UndefineSpaceSpecial