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
#include "NV_ReadLock_fp.h"

#if CC_NV_ReadLock  // Conditional expansion of this file

/*(See part 3 specification)
// Set read lock on a NV index
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES               TPMA_NV_READ_STCLEAR is not SET so
//                                      Index referenced by 'nvIndex' may not be
//                                      write locked
//      TPM_RC_NV_AUTHORIZATION         the authorization was valid but the
//                                      authorizing entity ('authHandle')
//                                      is not allowed to read from the Index
//                                      referenced by 'nvIndex'
TPM_RC
TPM2_NV_ReadLock(
    NV_ReadLock_In  *in             // IN: input parameter list
    )
{
    TPM_RC           result;
    NV_REF           locator;
    // The referenced index has been checked multiple times before this is called
    // so it must be present and will be loaded into cache
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, &locator);
    TPMA_NV          nvAttributes = nvIndex->publicArea.attributes;

// Input Validation
    // Common read access checks. NvReadAccessChecks() may return
    // TPM_RC_NV_AUTHORIZATION, TPM_RC_NV_LOCKED, or TPM_RC_NV_UNINITIALIZED
    result = NvReadAccessChecks(in->authHandle,
                                in->nvIndex,
                                nvAttributes);
    if(result == TPM_RC_NV_AUTHORIZATION)
        return TPM_RC_NV_AUTHORIZATION;
    // Index is already locked for write
    else if(result == TPM_RC_NV_LOCKED)
            return TPM_RC_SUCCESS;

    // If NvReadAccessChecks return TPM_RC_NV_UNINITALIZED, then continue.
    // It is not an error to read lock an uninitialized Index.
    
    // if TPMA_NV_READ_STCLEAR is not set, the index can not be read-locked
    if(!IS_ATTRIBUTE(nvAttributes, TPMA_NV, READ_STCLEAR))   
        return TPM_RCS_ATTRIBUTES + RC_NV_ReadLock_nvIndex;

// Internal Data Update

    // Set the READLOCK attribute
    SET_ATTRIBUTE(nvAttributes, TPMA_NV, READLOCKED);

    // Write NV info back
    return NvWriteIndexAttributes(nvIndex->publicArea.nvIndex,
                                  locator,
                                  nvAttributes);
}

#endif // CC_NV_ReadLock