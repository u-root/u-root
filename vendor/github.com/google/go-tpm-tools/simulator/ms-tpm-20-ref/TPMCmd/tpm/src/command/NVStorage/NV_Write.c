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
#include "NV_Write_fp.h"

#if CC_NV_Write  // Conditional expansion of this file

/*(See part 3 specification)
// Write to a NV index
*/
//  Return Type: TPM_RC
//      TPM_RC_ATTRIBUTES               Index referenced by 'nvIndex' has either
//                                      TPMA_NV_BITS, TPMA_NV_COUNTER, or
//                                      TPMA_NV_EVENT attribute SET
//      TPM_RC_NV_AUTHORIZATION         the authorization was valid but the
//                                      authorizing entity ('authHandle')
//                                      is not allowed to write to the Index
//                                      referenced by 'nvIndex'
//      TPM_RC_NV_LOCKED                Index referenced by 'nvIndex' is write
//                                      locked
//      TPM_RC_NV_RANGE                 if TPMA_NV_WRITEALL is SET then the write
//                                      is not the size of the Index referenced by
//                                      'nvIndex'; otherwise, the write extends
//                                      beyond the limits of the Index
//
TPM_RC
TPM2_NV_Write(
    NV_Write_In     *in             // IN: input parameter list
    )
{
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, NULL);
    TPMA_NV          attributes = nvIndex->publicArea.attributes;
    TPM_RC           result;

// Input Validation

    // Common access checks, NvWriteAccessCheck() may return TPM_RC_NV_AUTHORIZATION
    // or TPM_RC_NV_LOCKED 
    result = NvWriteAccessChecks(in->authHandle,
                                 in->nvIndex,
                                 attributes);
    if(result != TPM_RC_SUCCESS)
        return result;

    // Bits index, extend index or counter index may not be updated by
    // TPM2_NV_Write
    if(IsNvCounterIndex(attributes)
       || IsNvBitsIndex(attributes)
       || IsNvExtendIndex(attributes))
        return TPM_RC_ATTRIBUTES;

    // Make sure that the offset is not too large
    if(in->offset > nvIndex->publicArea.dataSize)
        return TPM_RCS_VALUE + RC_NV_Write_offset;

    // Make sure that the selection is within the range of the Index
    if(in->data.t.size > (nvIndex->publicArea.dataSize - in->offset))
        return TPM_RC_NV_RANGE;

    // If this index requires a full sized write, make sure that input range is
    // full sized.
    // Note: if the requested size is the same as the Index data size, then offset
    // will have to be zero. Otherwise, the range check above would have failed.
    if(IS_ATTRIBUTE(attributes, TPMA_NV, WRITEALL)
       && in->data.t.size < nvIndex->publicArea.dataSize)   
        return TPM_RC_NV_RANGE;

// Internal Data Update

    // Perform the write.  This called routine will SET the TPMA_NV_WRITTEN
    // attribute if it has not already been SET. If NV isn't available, an error
    // will be returned.
    return NvWriteIndexData(nvIndex, in->offset, in->data.t.size,
                            in->data.t.buffer);
}

#endif // CC_NV_Write