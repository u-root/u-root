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
#include "NV_SetBits_fp.h"

#if CC_NV_SetBits  // Conditional expansion of this file

/*(See part 3 specification)
// Set bits in a NV index
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES               the TPMA_NV_BITS attribute is not SET in the
//                                      Index referenced by 'nvIndex'
//      TPM_RC_NV_AUTHORIZATION         the authorization was valid but the
//                                      authorizing entity ('authHandle')
//                                      is not allowed to write to the Index
//                                      referenced by 'nvIndex'
//      TPM_RC_NV_LOCKED                the Index referenced by 'nvIndex' is locked
//                                      for writing
TPM_RC
TPM2_NV_SetBits(
    NV_SetBits_In   *in             // IN: input parameter list
    )
{
    TPM_RC           result;
    NV_REF           locator;
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, &locator);
    UINT64           oldValue;
    UINT64           newValue;

// Input Validation

    // Common access checks, NvWriteAccessCheck() may return TPM_RC_NV_AUTHORIZATION
    // or TPM_RC_NV_LOCKED 
    result = NvWriteAccessChecks(in->authHandle,
                                 in->nvIndex,
                                 nvIndex->publicArea.attributes);
    if(result != TPM_RC_SUCCESS)
        return result;

    // Make sure that this is a bit field
    if(!IsNvBitsIndex(nvIndex->publicArea.attributes))
        return TPM_RCS_ATTRIBUTES + RC_NV_SetBits_nvIndex;

    // If index is not been written, initialize it
    if(!IS_ATTRIBUTE(nvIndex->publicArea.attributes, TPMA_NV, WRITTEN))   
        oldValue = 0;
    else
        // Read index data
        oldValue = NvGetUINT64Data(nvIndex, locator);

    // Figure out what the new value is going to be
    newValue = oldValue | in->bits;

// Internal Data Update
    return  NvWriteUINT64Data(nvIndex, newValue);
}

#endif // CC_NV_SetBits