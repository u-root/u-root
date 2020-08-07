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
#include "AC_Send_fp.h"
#include "AC_spt_fp.h"


#if CC_AC_Send  // Conditional expansion of this file

/*(See part 3 specification)
// Duplicate a loaded object
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES   key to duplicate has 'fixedParent' SET
//      TPM_RC_HASH         for an RSA key, the nameAlg digest size for the
//                          newParent is not compatible with the key size
//      TPM_RC_HIERARCHY    'encryptedDuplication' is SET and 'newParentHandle'
//                          specifies Null Hierarchy
//      TPM_RC_KEY          'newParentHandle' references invalid ECC key (public
//                          point not on the curve)
//      TPM_RC_SIZE         input encryption key size does not match the
//                          size specified in symmetric algorithm
//      TPM_RC_SYMMETRIC    'encryptedDuplication' is SET but no symmetric
//                          algorithm is provided
//      TPM_RC_TYPE         'newParentHandle' is neither a storage key nor
//                          TPM_RH_NULL; or the object has a NULL nameAlg
//      TPM_RC_VALUE        for an RSA newParent, the sizes of the digest and
//                          the encryption key are too large to be OAEP encoded
TPM_RC
TPM2_AC_Send(
    AC_Send_In    *in,              // IN: input parameter list
    AC_Send_Out   *out              // OUT: output parameter list
)
{
    NV_REF           locator;
    TPM_HANDLE       nvAlias = ((in->ac - AC_FIRST) + NV_AC_FIRST);
    NV_INDEX        *nvIndex = NvGetIndexInfo(nvAlias, &locator);
    OBJECT          *object = HandleToObject(in->sendObject);
    TPM_RC           result;
// Input validation
    // If there is an NV alias, then the index must allow the authorization provided
    if(nvIndex != NULL)
    {
        // Common access checks, NvWriteAccessCheck() may return 
        // TPM_RC_NV_AUTHORIZATION or TPM_RC_NV_LOCKED 
        result = NvWriteAccessChecks(in->authHandle, nvAlias,
                                     nvIndex->publicArea.attributes);
        if(result != TPM_RC_SUCCESS)
            return result;
    }
    // If 'ac' did not have an alias then the authorization had to be with either
    // platform or owner authorization. The type of TPMI_RH_NV_AUTH only allows
    // owner or platform or an NV index. If it was a valid index, it would have had
    // an alias and be processed above, so only success here is if this is a
    // permanent handle.
    else if(HandleGetType(in->authHandle) != TPM_HT_PERMANENT)
        return TPM_RCS_HANDLE + RC_AC_Send_authHandle;
    // Make sure that the object to be duplicated has the right attributes
    if(IS_ATTRIBUTE(object->publicArea.objectAttributes, 
                    TPMA_OBJECT, encryptedDuplication)
       || IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, 
                       fixedParent)
       || IS_ATTRIBUTE(object->publicArea.objectAttributes, TPMA_OBJECT, fixedTPM))
        return TPM_RCS_ATTRIBUTES + RC_AC_Send_sendObject;
// Command output
    // Do the implementation dependent send
    return AcSendObject(in->ac, object, &out->acDataOut);
}

#endif // TPM_CC_AC_Send