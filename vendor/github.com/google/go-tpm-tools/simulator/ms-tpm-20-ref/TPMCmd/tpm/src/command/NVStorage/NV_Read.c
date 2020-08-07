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
#include "NV_Read_fp.h"

#if CC_NV_Read  // Conditional expansion of this file

/*(See part 3 specification)
// Read of an NV index
*/
//  Return Type: TPM_RC
//      TPM_RC_NV_AUTHORIZATION         the authorization was valid but the
//                                      authorizing entity ('authHandle')
//                                      is not allowed to read from the Index
//                                      referenced by 'nvIndex'
//      TPM_RC_NV_LOCKED                the Index referenced by 'nvIndex' is
//                                      read locked
//      TPM_RC_NV_RANGE                 read range defined by 'size' and 'offset'
//                                      is outside the range of the Index referenced
//                                      by 'nvIndex'
//      TPM_RC_NV_UNINITIALIZED         the Index referenced by 'nvIndex' has
//                                      not been initialized (written)
//      TPM_RC_VALUE                    the read size is larger than the
//                                      MAX_NV_BUFFER_SIZE
TPM_RC
TPM2_NV_Read(
    NV_Read_In      *in,            // IN: input parameter list
    NV_Read_Out     *out            // OUT: output parameter list
    )
{
    NV_REF           locator;
    NV_INDEX        *nvIndex = NvGetIndexInfo(in->nvIndex, &locator);
    TPM_RC           result;

// Input Validation
    // Common read access checks. NvReadAccessChecks() may return
    // TPM_RC_NV_AUTHORIZATION, TPM_RC_NV_LOCKED, or TPM_RC_NV_UNINITIALIZED
    result = NvReadAccessChecks(in->authHandle, in->nvIndex,
                                nvIndex->publicArea.attributes);
    if(result != TPM_RC_SUCCESS)
        return result;

    // Make sure the data will fit the return buffer
    if(in->size > MAX_NV_BUFFER_SIZE)
        return TPM_RCS_VALUE + RC_NV_Read_size;

    // Verify that the offset is not too large
    if(in->offset > nvIndex->publicArea.dataSize)
        return TPM_RCS_VALUE + RC_NV_Read_offset;

    // Make sure that the selection is within the range of the Index
    if(in->size > (nvIndex->publicArea.dataSize - in->offset))
        return TPM_RC_NV_RANGE;

// Command Output
    // Set the return size
    out->data.t.size = in->size;

    // Perform the read
    NvGetIndexData(nvIndex, locator, in->offset, in->size, out->data.t.buffer);

    return TPM_RC_SUCCESS;
}

#endif // CC_NV_Read