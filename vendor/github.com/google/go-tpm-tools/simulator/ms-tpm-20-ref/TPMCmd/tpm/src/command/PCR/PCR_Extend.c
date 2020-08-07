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
#include "PCR_Extend_fp.h"

#if CC_PCR_Extend  // Conditional expansion of this file

/*(See part 3 specification)
// Update PCR
*/
//  Return Type: TPM_RC
//      TPM_RC_LOCALITY             current command locality is not allowed to
//                                  extend the PCR referenced by 'pcrHandle'
TPM_RC
TPM2_PCR_Extend(
    PCR_Extend_In   *in             // IN: input parameter list
    )
{
    UINT32              i;

// Input Validation

    // NOTE: This function assumes that the unmarshaling function for 'digests' will
    // have validated that all of the indicated hash algorithms are valid. If the
    // hash algorithms are correct, the unmarshaling code will unmarshal a digest
    // of the size indicated by the hash algorithm. If the overall size is not
    // consistent, the unmarshaling code will run out of input data or have input
    // data left over. In either case, it will cause an unmarshaling error and this
    // function will not be called.

    // For NULL handle, do nothing and return success
    if(in->pcrHandle == TPM_RH_NULL)
        return TPM_RC_SUCCESS;

    // Check if the extend operation is allowed by the current command locality
    if(!PCRIsExtendAllowed(in->pcrHandle))
        return TPM_RC_LOCALITY;

    // If PCR is state saved and we need to update orderlyState, check NV
    // availability
    if(PCRIsStateSaved(in->pcrHandle))
        RETURN_IF_ORDERLY;

// Internal Data Update

    // Iterate input digest list to extend
    for(i = 0; i < in->digests.count; i++)
    {
        PCRExtend(in->pcrHandle, in->digests.digests[i].hashAlg,
                  CryptHashGetDigestSize(in->digests.digests[i].hashAlg),
                  (BYTE *)&in->digests.digests[i].digest);
    }

    return TPM_RC_SUCCESS;
}

#endif // CC_PCR_Extend